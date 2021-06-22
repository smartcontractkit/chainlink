import { ethers } from "hardhat";
import { publicAbi } from "../test-helpers/helpers";
import { assert, expect } from "chai";
import { Contract, ContractFactory, ContractReceipt } from "ethers";
import { getUsers, Roles } from "../test-helpers/setup";
import { evmRevert } from "../test-helpers/matchers";

let getterSetterFactory: ContractFactory;
let forwarderFactory: ContractFactory;
let linkTokenFactory: ContractFactory;

let roles: Roles;
const zeroAddress = ethers.constants.AddressZero;

before(async () => {
  const users = await getUsers();

  roles = users.roles;
  getterSetterFactory = await ethers.getContractFactory(
    "src/v0.4/tests/GetterSetter.sol:GetterSetter",
    roles.defaultAccount,
  );
  forwarderFactory = await ethers.getContractFactory("AuthorizedForwarder", roles.defaultAccount);
  linkTokenFactory = await ethers.getContractFactory("LinkToken", roles.defaultAccount);
});

describe("AuthorizedForwarder", () => {
  let link: Contract;
  let forwarder: Contract;

  beforeEach(async () => {
    link = await linkTokenFactory.connect(roles.defaultAccount).deploy();
    forwarder = await forwarderFactory
      .connect(roles.defaultAccount)
      .deploy(link.address, await roles.defaultAccount.getAddress(), zeroAddress, "0x");
  });

  it("has a limited public interface", () => {
    publicAbi(forwarder, [
      "forward",
      "getAuthorizedSenders",
      "getChainlinkToken",
      "isAuthorizedSender",
      "ownerForward",
      "setAuthorizedSenders",
      "transferOwnershipWithMessage",
      // ConfirmedOwner
      "transferOwnership",
      "acceptOwnership",
      "owner",
    ]);
  });

  describe("deployment", () => {
    it("sets the correct link token", async () => {
      assert.equal(await forwarder.getChainlinkToken(), link.address);
    });

    it("sets no authorized senders", async () => {
      const senders = await forwarder.getAuthorizedSenders();
      assert.equal(senders.length, 0);
    });
  });

  describe("#setAuthorizedSenders", () => {
    let newSenders: string[];
    let receipt: ContractReceipt;
    describe("when called by the owner", () => {
      describe("setting 3 authorized senders", () => {
        beforeEach(async () => {
          newSenders = [
            await roles.oracleNode1.getAddress(),
            await roles.oracleNode2.getAddress(),
            await roles.oracleNode3.getAddress(),
          ];
          const tx = await forwarder.connect(roles.defaultAccount).setAuthorizedSenders(newSenders);
          receipt = await tx.wait();
        });

        it("adds the authorized nodes", async () => {
          const authorizedSenders = await forwarder.getAuthorizedSenders();
          assert.equal(newSenders.length, authorizedSenders.length);
          for (let i = 0; i < authorizedSenders.length; i++) {
            assert.equal(authorizedSenders[i], newSenders[i]);
          }
        });

        it("emits an event", async () => {
          assert.equal(receipt.events?.length, 1);
          const responseEvent = receipt.events?.[0];
          assert.equal(responseEvent?.event, "AuthorizedSendersChanged");
          const encodedSenders = ethers.utils.defaultAbiCoder.encode(
            ["address[]", "address"],
            [newSenders, await roles.defaultAccount.getAddress()],
          );
          assert.equal(responseEvent?.data, encodedSenders);
        });

        it("replaces the authorized nodes", async () => {
          const newSenders = await forwarder.connect(roles.defaultAccount).getAuthorizedSenders();
          assert.notIncludeOrderedMembers(newSenders, [await roles.oracleNode.getAddress()]);
        });

        after(async () => {
          await forwarder.connect(roles.defaultAccount).setAuthorizedSenders([await roles.oracleNode.getAddress()]);
        });
      });

      describe("setting 0 authorized senders", () => {
        beforeEach(async () => {
          newSenders = [];
        });

        it("reverts with a minimum senders message", async () => {
          await evmRevert(
            forwarder.connect(roles.defaultAccount).setAuthorizedSenders(newSenders),
            "Must have at least 1 authorized sender",
          );
        });
      });
    });

    describe("when called by a non-owner", () => {
      it("cannot add an authorized node", async () => {
        await evmRevert(
          forwarder.connect(roles.stranger).setAuthorizedSenders([await roles.stranger.getAddress()]),
          "Cannot set authorized senders",
        );
      });
    });
  });

  describe("#forward", () => {
    let bytes: string;
    let payload: string;
    let mock: Contract;

    beforeEach(async () => {
      mock = await getterSetterFactory.connect(roles.defaultAccount).deploy();
      bytes = ethers.utils.hexlify(ethers.utils.randomBytes(100));
      payload = getterSetterFactory.interface.encodeFunctionData(
        getterSetterFactory.interface.getFunction("setBytes"),
        [bytes],
      );
    });

    describe("when called by an unauthorized node", () => {
      it("reverts", async () => {
        await evmRevert(forwarder.connect(roles.stranger).forward(mock.address, payload));
      });
    });

    describe("when called by an authorized node", () => {
      beforeEach(async () => {
        await forwarder.connect(roles.defaultAccount).setAuthorizedSenders([await roles.defaultAccount.getAddress()]);
      });

      describe("when sending to a non-contract address", () => {
        it("reverts", async () => {
          await evmRevert(
            forwarder.connect(roles.defaultAccount).forward(zeroAddress, payload),
            "Must forward to a contract",
          );
        });
      });

      describe("when attempting to forward to the link token", () => {
        it("reverts", async () => {
          const sighash = linkTokenFactory.interface.getSighash("name"); // any Link Token function
          await evmRevert(forwarder.connect(roles.defaultAccount).forward(link.address, sighash));
        });
      });

      describe("when forwarding to any other address", () => {
        it("forwards the data", async () => {
          const tx = await forwarder.connect(roles.defaultAccount).forward(mock.address, payload);
          await tx.wait();
          assert.equal(await mock.getBytes(), bytes);
        });

        it("perceives the message is sent by the AuthorizedForwarder", async () => {
          const tx = await forwarder.connect(roles.defaultAccount).forward(mock.address, payload);
          await expect(tx).to.emit(mock, "SetBytes").withArgs(forwarder.address, bytes);
        });
      });
    });
  });

  describe("#transferOwnershipWithMessage", () => {
    const message = "0x42";

    describe("when called by a non-owner", () => {
      it("reverts", async () => {
        await evmRevert(
          forwarder.connect(roles.stranger).transferOwnershipWithMessage(await roles.stranger.getAddress(), message),
          "Only callable by owner",
        );
      });
    });

    describe("when called by the owner", () => {
      it("calls the normal ownership transfer proposal", async () => {
        const tx = await forwarder
          .connect(roles.defaultAccount)
          .transferOwnershipWithMessage(await roles.stranger.getAddress(), message);
        const receipt = await tx.wait();

        assert.equal(receipt?.events?.[0]?.event, "OwnershipTransferRequested");
        assert.equal(receipt?.events?.[0]?.address, forwarder.address);
        assert.equal(receipt?.events?.[0]?.args?.[0], await roles.defaultAccount.getAddress());
        assert.equal(receipt?.events?.[0]?.args?.[1], await roles.stranger.getAddress());
      });

      it("calls the normal ownership transfer proposal", async () => {
        const tx = await forwarder
          .connect(roles.defaultAccount)
          .transferOwnershipWithMessage(await roles.stranger.getAddress(), message);
        const receipt = await tx.wait();

        assert.equal(receipt?.events?.[1]?.event, "OwnershipTransferRequestedWithMessage");
        assert.equal(receipt?.events?.[1]?.address, forwarder.address);
        assert.equal(receipt?.events?.[1]?.args?.[0], await roles.defaultAccount.getAddress());
        assert.equal(receipt?.events?.[1]?.args?.[1], await roles.stranger.getAddress());
        assert.equal(receipt?.events?.[1]?.args?.[2], message);
      });
    });
  });

  describe("#ownerForward", () => {
    let bytes: string;
    let payload: string;
    let mock: Contract;

    beforeEach(async () => {
      mock = await getterSetterFactory.connect(roles.defaultAccount).deploy();
      bytes = ethers.utils.hexlify(ethers.utils.randomBytes(100));
      payload = getterSetterFactory.interface.encodeFunctionData(
        getterSetterFactory.interface.getFunction("setBytes"),
        [bytes],
      );
    });

    describe("when called by a non-owner", () => {
      it("reverts", async () => {
        await evmRevert(forwarder.connect(roles.stranger).ownerForward(mock.address, payload));
      });
    });

    describe("when called by owner", () => {
      describe("when attempting to forward to the link token", () => {
        it("does not revert", async () => {
          const sighash = linkTokenFactory.interface.getSighash("name"); // any Link Token function

          await forwarder.connect(roles.defaultAccount).ownerForward(link.address, sighash);
        });
      });

      describe("when forwarding to any other address", () => {
        it("forwards the data", async () => {
          const tx = await forwarder.connect(roles.defaultAccount).ownerForward(mock.address, payload);
          await tx.wait();
          assert.equal(await mock.getBytes(), bytes);
        });

        it("reverts when sending to a non-contract address", async () => {
          await evmRevert(
            forwarder.connect(roles.defaultAccount).ownerForward(zeroAddress, payload),
            "Must forward to a contract",
          );
        });

        it("perceives the message is sent by the Operator", async () => {
          const tx = await forwarder.connect(roles.defaultAccount).ownerForward(mock.address, payload);
          await expect(tx).to.emit(mock, "SetBytes").withArgs(forwarder.address, bytes);
        });
      });
    });
  });
});
