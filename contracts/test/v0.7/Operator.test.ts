import { ethers } from "hardhat";
import { publicAbi, toBytes32String, toWei, stringToBytes, increaseTime5Minutes } from "../test-helpers/helpers";
import { assert, expect } from "chai";
import { BigNumber, constants, Contract, ContractFactory, ContractReceipt, ContractTransaction, Signer } from "ethers";
import { getUsers, Roles } from "../test-helpers/setup";
import { bigNumEquals, evmRevert } from "../test-helpers/matchers";
import type { providers } from "ethers";
import {
  convertCancelParams,
  convertFufillParams,
  convertFulfill2Params,
  decodeRunRequest,
  encodeOracleRequest,
  encodeRequestOracleData,
  RunRequest,
} from "../test-helpers/oracle";

let v7ConsumerFactory: ContractFactory;
let basicConsumerFactory: ContractFactory;
let multiWordConsumerFactory: ContractFactory;
let gasGuzzlingConsumerFactory: ContractFactory;
let getterSetterFactory: ContractFactory;
let maliciousRequesterFactory: ContractFactory;
let maliciousConsumerFactory: ContractFactory;
let maliciousMultiWordConsumerFactory: ContractFactory;
let operatorFactory: ContractFactory;
let forwarderFactory: ContractFactory;
let linkTokenFactory: ContractFactory;
const zeroAddress = ethers.constants.AddressZero;

let roles: Roles;

before(async () => {
  const users = await getUsers();

  roles = users.roles;
  v7ConsumerFactory = await ethers.getContractFactory("src/v0.7/tests/Consumer.sol:Consumer");
  basicConsumerFactory = await ethers.getContractFactory("src/v0.6/tests/BasicConsumer.sol:BasicConsumer");
  multiWordConsumerFactory = await ethers.getContractFactory("src/v0.6/tests/MultiWordConsumer.sol:MultiWordConsumer");
  gasGuzzlingConsumerFactory = await ethers.getContractFactory(
    "src/v0.6/tests/GasGuzzlingConsumer.sol:GasGuzzlingConsumer",
  );
  getterSetterFactory = await ethers.getContractFactory("src/v0.4/tests/GetterSetter.sol:GetterSetter");
  maliciousRequesterFactory = await ethers.getContractFactory(
    "src/v0.4/tests/MaliciousRequester.sol:MaliciousRequester",
  );
  maliciousConsumerFactory = await ethers.getContractFactory("src/v0.4/tests/MaliciousConsumer.sol:MaliciousConsumer");
  maliciousMultiWordConsumerFactory = await ethers.getContractFactory("MaliciousMultiWordConsumer");
  operatorFactory = await ethers.getContractFactory("Operator");
  forwarderFactory = await ethers.getContractFactory("AuthorizedForwarder");
  linkTokenFactory = await ethers.getContractFactory("LinkToken");
});

describe("Operator", () => {
  let fHash: string;
  let specId: string;
  let to: string;
  let link: Contract;
  let operator: Contract;
  let forwarder1: Contract;
  let forwarder2: Contract;
  let owner: Signer;

  beforeEach(async () => {
    fHash = getterSetterFactory.interface.getSighash("requestedBytes32");
    specId = "0x4c7b7ffb66b344fbaa64995af81e355a00000000000000000000000000000000";
    to = "0x80e29acb842498fe6591f020bd82766dce619d43";
    link = await linkTokenFactory.connect(roles.defaultAccount).deploy();
    owner = roles.defaultAccount;
    operator = await operatorFactory.connect(owner).deploy(link.address, await owner.getAddress());
    await operator.connect(roles.defaultAccount).setAuthorizedSenders([await roles.oracleNode.getAddress()]);
  });

  it("has a limited public interface", () => {
    publicAbi(operator, [
      "acceptAuthorizedReceivers",
      "acceptOwnableContracts",
      "cancelOracleRequest",
      "distributeFunds",
      "fulfillOracleRequest",
      "fulfillOracleRequest2",
      "getAuthorizedSenders",
      "getChainlinkToken",
      "getExpiryTime",
      "isAuthorizedSender",
      "onTokenTransfer",
      "oracleRequest",
      "ownerForward",
      "ownerTransferAndCall",
      "requestOracleData",
      "setAuthorizedSenders",
      "setAuthorizedSendersOn",
      "transferOwnableContracts",
      "withdraw",
      "withdrawable",
      // Ownable methods:
      "acceptOwnership",
      "owner",
      "transferOwnership",
    ]);
  });

  describe("#transferOwnableContracts", () => {
    beforeEach(async () => {
      forwarder1 = await forwarderFactory.connect(owner).deploy(link.address, operator.address, zeroAddress, "0x");
      forwarder2 = await forwarderFactory.connect(owner).deploy(link.address, operator.address, zeroAddress, "0x");
    });

    describe("being called by the owner", () => {
      it("cannot transfer to self", async () => {
        await evmRevert(
          operator.connect(owner).transferOwnableContracts([forwarder1.address], operator.address),
          "Cannot transfer to self",
        );
      });

      it("emits an ownership transfer request event", async () => {
        const tx = await operator
          .connect(owner)
          .transferOwnableContracts([forwarder1.address, forwarder2.address], await roles.oracleNode1.getAddress());
        const receipt = await tx.wait();
        assert.equal(receipt?.events?.length, 2);
        const log1 = receipt?.events?.[0];
        assert.equal(log1?.event, "OwnershipTransferRequested");
        assert.equal(log1?.address, forwarder1.address);
        assert.equal(log1?.args?.[0], operator.address);
        assert.equal(log1?.args?.[1], await roles.oracleNode1.getAddress());
        const log2 = receipt?.events?.[1];
        assert.equal(log2?.event, "OwnershipTransferRequested");
        assert.equal(log2?.address, forwarder2.address);
        assert.equal(log2?.args?.[0], operator.address);
        assert.equal(log2?.args?.[1], await roles.oracleNode1.getAddress());
      });
    });

    describe("being called by a non-owner", () => {
      it("reverts with message", async () => {
        await evmRevert(
          operator
            .connect(roles.stranger)
            .transferOwnableContracts([forwarder1.address], await roles.oracleNode2.getAddress()),
          "Only callable by owner",
        );
      });
    });
  });

  describe("#acceptOwnableContracts", () => {
    describe("being called by the owner", () => {
      let operator2: Contract;
      let receipt: ContractReceipt;

      beforeEach(async () => {
        operator2 = await operatorFactory
          .connect(roles.defaultAccount)
          .deploy(link.address, await roles.defaultAccount.getAddress());
        forwarder1 = await forwarderFactory
          .connect(roles.defaultAccount)
          .deploy(link.address, operator.address, zeroAddress, "0x");
        forwarder2 = await forwarderFactory
          .connect(roles.defaultAccount)
          .deploy(link.address, operator.address, zeroAddress, "0x");
        await operator
          .connect(roles.defaultAccount)
          .transferOwnableContracts([forwarder1.address, forwarder2.address], operator2.address);
        const tx = await operator2
          .connect(roles.defaultAccount)
          .acceptOwnableContracts([forwarder1.address, forwarder2.address]);
        receipt = await tx.wait();
      });

      it("sets the new owner on the forwarder", async () => {
        assert.equal(await forwarder1.owner(), operator2.address);
      });

      it("emits ownership transferred events", async () => {
        assert.equal(receipt?.events?.[0]?.event, "OwnershipTransferred");
        assert.equal(receipt?.events?.[0]?.address, forwarder1.address);
        assert.equal(receipt?.events?.[0]?.args?.[0], operator.address);
        assert.equal(receipt?.events?.[0]?.args?.[1], operator2.address);

        assert.equal(receipt?.events?.[1]?.event, "OwnableContractAccepted");
        assert.equal(receipt?.events?.[1]?.args?.[0], forwarder1.address);

        assert.equal(receipt?.events?.[2]?.event, "OwnershipTransferred");
        assert.equal(receipt?.events?.[2]?.address, forwarder2.address);
        assert.equal(receipt?.events?.[2]?.args?.[0], operator.address);
        assert.equal(receipt?.events?.[2]?.args?.[1], operator2.address);

        assert.equal(receipt?.events?.[3]?.event, "OwnableContractAccepted");
        assert.equal(receipt?.events?.[3]?.args?.[0], forwarder2.address);
      });
    });

    describe("being called by a non-owner authorized sender", () => {
      it("does not revert", async () => {
        await operator.connect(roles.defaultAccount).setAuthorizedSenders([await roles.oracleNode1.getAddress()]);

        await operator.connect(roles.oracleNode1).acceptOwnableContracts([]);
      });
    });

    describe("being called by a non owner", () => {
      it("reverts with message", async () => {
        await evmRevert(
          operator.connect(roles.stranger).acceptOwnableContracts([await roles.oracleNode2.getAddress()]),
          "Cannot set authorized senders",
        );
      });
    });
  });

  describe("#distributeFunds", () => {
    describe("when called with empty arrays", () => {
      it("reverts with invalid array message", async () => {
        await evmRevert(operator.connect(roles.defaultAccount).distributeFunds([], []), "Invalid array length(s)");
      });
    });

    describe("when called with unequal array lengths", () => {
      it("reverts with invalid array message", async () => {
        const receivers = [await roles.oracleNode2.getAddress(), await roles.oracleNode3.getAddress()];
        const amounts = [1, 2, 3];
        await evmRevert(
          operator.connect(roles.defaultAccount).distributeFunds(receivers, amounts),
          "Invalid array length(s)",
        );
      });
    });

    describe("when called with not enough ETH", () => {
      it("reverts with subtraction overflow message", async () => {
        const amountToSend = toWei("2");
        const ethSent = toWei("1");
        await evmRevert(
          operator
            .connect(roles.defaultAccount)
            .distributeFunds([await roles.oracleNode2.getAddress()], [amountToSend], {
              value: ethSent,
            }),
          "SafeMath: subtraction overflow",
        );
      });
    });

    describe("when called with too much ETH", () => {
      it("reverts with too much ETH message", async () => {
        const amountToSend = toWei("2");
        const ethSent = toWei("3");
        await evmRevert(
          operator
            .connect(roles.defaultAccount)
            .distributeFunds([await roles.oracleNode2.getAddress()], [amountToSend], {
              value: ethSent,
            }),
          "Too much ETH sent",
        );
      });
    });

    describe("when called with correct values", () => {
      it("updates the balances", async () => {
        const node2BalanceBefore = await roles.oracleNode2.getBalance();
        const node3BalanceBefore = await roles.oracleNode3.getBalance();
        const receivers = [await roles.oracleNode2.getAddress(), await roles.oracleNode3.getAddress()];
        const sendNode2 = toWei("2");
        const sendNode3 = toWei("3");
        const totalAmount = toWei("5");
        const amounts = [sendNode2, sendNode3];

        await operator.connect(roles.defaultAccount).distributeFunds(receivers, amounts, { value: totalAmount });

        const node2BalanceAfter = await roles.oracleNode2.getBalance();
        const node3BalanceAfter = await roles.oracleNode3.getBalance();

        assert.equal(node2BalanceAfter.sub(node2BalanceBefore).toString(), sendNode2.toString());

        assert.equal(node3BalanceAfter.sub(node3BalanceBefore).toString(), sendNode3.toString());
      });
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
          const tx = await operator.connect(roles.defaultAccount).setAuthorizedSenders(newSenders);
          receipt = await tx.wait();
        });

        it("adds the authorized nodes", async () => {
          const authorizedSenders = await operator.getAuthorizedSenders();
          assert.equal(newSenders.length, authorizedSenders.length);
          for (let i = 0; i < authorizedSenders.length; i++) {
            assert.equal(authorizedSenders[i], newSenders[i]);
          }
        });

        it("emits an event on the Operator", async () => {
          assert.equal(receipt.events?.length, 1);

          const encodedSenders1 = ethers.utils.defaultAbiCoder.encode(
            ["address[]", "address"],
            [newSenders, await roles.defaultAccount.getAddress()],
          );

          const responseEvent1 = receipt.events?.[0];
          assert.equal(responseEvent1?.event, "AuthorizedSendersChanged");
          assert.equal(responseEvent1?.data, encodedSenders1);
        });

        it("replaces the authorized nodes", async () => {
          const originalAuthorization = await operator
            .connect(roles.defaultAccount)
            .isAuthorizedSender(await roles.oracleNode.getAddress());
          assert.isFalse(originalAuthorization);
        });

        after(async () => {
          await operator.connect(roles.defaultAccount).setAuthorizedSenders([await roles.oracleNode.getAddress()]);
        });
      });

      describe("setting 0 authorized senders", () => {
        beforeEach(async () => {
          newSenders = [];
        });

        it("reverts with a minimum senders message", async () => {
          await evmRevert(
            operator.connect(roles.defaultAccount).setAuthorizedSenders(newSenders),
            "Must have at least 1 authorized sender",
          );
        });
      });
    });

    describe("when called by an authorized sender", () => {
      beforeEach(async () => {
        newSenders = [await roles.oracleNode1.getAddress()];
        await operator.connect(roles.defaultAccount).setAuthorizedSenders(newSenders);
      });

      it("succeeds", async () => {
        await operator.connect(roles.defaultAccount).setAuthorizedSenders([await roles.stranger.getAddress()]);
      });
    });

    describe("when called by a non-owner", () => {
      it("cannot add an authorized node", async () => {
        await evmRevert(
          operator.connect(roles.stranger).setAuthorizedSenders([await roles.stranger.getAddress()]),
          "Cannot set authorized senders",
        );
      });
    });
  });

  describe("#setAuthorizedSendersOn", () => {
    let newSenders: string[];

    beforeEach(async () => {
      await operator.connect(roles.defaultAccount).setAuthorizedSenders([await roles.oracleNode1.getAddress()]);
      newSenders = [await roles.oracleNode2.getAddress(), await roles.oracleNode3.getAddress()];

      forwarder1 = await forwarderFactory.connect(owner).deploy(link.address, operator.address, zeroAddress, "0x");
      forwarder2 = await forwarderFactory.connect(owner).deploy(link.address, operator.address, zeroAddress, "0x");
    });

    describe("when called by a non-authorized sender", () => {
      it("reverts", async () => {
        await evmRevert(
          operator.connect(roles.stranger).setAuthorizedSendersOn(newSenders, [forwarder1.address]),
          "Cannot set authorized senders",
        );
      });
    });

    describe("when called by an owner", () => {
      it("does not revert", async () => {
        await operator
          .connect(roles.defaultAccount)
          .setAuthorizedSendersOn([forwarder1.address, forwarder2.address], newSenders);
      });
    });

    describe("when called by an authorized sender", () => {
      it("does not revert", async () => {
        await operator
          .connect(roles.oracleNode1)
          .setAuthorizedSendersOn([forwarder1.address, forwarder2.address], newSenders);
      });

      it("does revert with 0 senders", async () => {
        await operator
          .connect(roles.oracleNode1)
          .setAuthorizedSendersOn([forwarder1.address, forwarder2.address], newSenders);
      });

      it("emits a log announcing the change and who made it", async () => {
        const targets = [forwarder1.address, forwarder2.address];
        const tx = await operator.connect(roles.oracleNode1).setAuthorizedSendersOn(targets, newSenders);

        const receipt = await tx.wait();
        const encodedArgs = ethers.utils.defaultAbiCoder.encode(
          ["address[]", "address[]", "address"],
          [targets, newSenders, await roles.oracleNode1.getAddress()],
        );

        const event1 = receipt.events?.[0];
        assert.equal(event1?.event, "TargetsUpdatedAuthorizedSenders");
        assert.equal(event1?.address, operator.address);
        assert.equal(event1?.data, encodedArgs);
      });

      it("updates the sender list on each of the targets", async () => {
        const tx = await operator
          .connect(roles.oracleNode1)
          .setAuthorizedSendersOn([forwarder1.address, forwarder2.address], newSenders);

        const receipt = await tx.wait();
        assert.equal(receipt.events?.length, 3, receipt.toString());
        const encodedSenders = ethers.utils.defaultAbiCoder.encode(
          ["address[]", "address"],
          [newSenders, operator.address],
        );

        const event1 = receipt.events?.[1];
        assert.equal(event1?.event, "AuthorizedSendersChanged");
        assert.equal(event1?.address, forwarder1.address);
        assert.equal(event1?.data, encodedSenders);

        const event2 = receipt.events?.[2];
        assert.equal(event2?.event, "AuthorizedSendersChanged");
        assert.equal(event2?.address, forwarder2.address);
        assert.equal(event2?.data, encodedSenders);
      });
    });
  });

  describe("#acceptAuthorizedReceivers", () => {
    let newSenders: string[];

    describe("being called by the owner", () => {
      let operator2: Contract;
      let receipt: ContractReceipt;

      beforeEach(async () => {
        operator2 = await operatorFactory
          .connect(roles.defaultAccount)
          .deploy(link.address, await roles.defaultAccount.getAddress());
        forwarder1 = await forwarderFactory
          .connect(roles.defaultAccount)
          .deploy(link.address, operator.address, zeroAddress, "0x");
        forwarder2 = await forwarderFactory
          .connect(roles.defaultAccount)
          .deploy(link.address, operator.address, zeroAddress, "0x");
        await operator
          .connect(roles.defaultAccount)
          .transferOwnableContracts([forwarder1.address, forwarder2.address], operator2.address);
        newSenders = [await roles.oracleNode2.getAddress(), await roles.oracleNode3.getAddress()];

        const tx = await operator2
          .connect(roles.defaultAccount)
          .acceptAuthorizedReceivers([forwarder1.address, forwarder2.address], newSenders);
        receipt = await tx.wait();
      });

      it("sets the new owner on the forwarder", async () => {
        assert.equal(await forwarder1.owner(), operator2.address);
      });

      it("emits ownership transferred events", async () => {
        assert.equal(receipt?.events?.[0]?.event, "OwnershipTransferred");
        assert.equal(receipt?.events?.[0]?.address, forwarder1.address);
        assert.equal(receipt?.events?.[0]?.args?.[0], operator.address);
        assert.equal(receipt?.events?.[0]?.args?.[1], operator2.address);

        assert.equal(receipt?.events?.[1]?.event, "OwnableContractAccepted");
        assert.equal(receipt?.events?.[1]?.args?.[0], forwarder1.address);

        assert.equal(receipt?.events?.[2]?.event, "OwnershipTransferred");
        assert.equal(receipt?.events?.[2]?.address, forwarder2.address);
        assert.equal(receipt?.events?.[2]?.args?.[0], operator.address);
        assert.equal(receipt?.events?.[2]?.args?.[1], operator2.address);

        assert.equal(receipt?.events?.[3]?.event, "OwnableContractAccepted");
        assert.equal(receipt?.events?.[3]?.args?.[0], forwarder2.address);

        assert.equal(receipt?.events?.[4]?.event, "TargetsUpdatedAuthorizedSenders");

        const encodedSenders = ethers.utils.defaultAbiCoder.encode(
          ["address[]", "address"],
          [newSenders, operator2.address],
        );
        assert.equal(receipt?.events?.[5]?.event, "AuthorizedSendersChanged");
        assert.equal(receipt?.events?.[5]?.address, forwarder1.address);
        assert.equal(receipt?.events?.[5]?.data, encodedSenders);

        assert.equal(receipt?.events?.[6]?.event, "AuthorizedSendersChanged");
        assert.equal(receipt?.events?.[6]?.address, forwarder2.address);
        assert.equal(receipt?.events?.[6]?.data, encodedSenders);
      });
    });

    describe("being called by a non owner", () => {
      it("reverts with message", async () => {
        await evmRevert(
          operator
            .connect(roles.stranger)
            .acceptAuthorizedReceivers([forwarder1.address, forwarder2.address], newSenders),
          "Cannot set authorized senders",
        );
      });
    });
  });

  describe("#onTokenTransfer", () => {
    describe("when called from any address but the LINK token", () => {
      it("triggers the intended method", async () => {
        const callData = encodeOracleRequest(specId, to, fHash, 0, constants.HashZero);

        await evmRevert(operator.onTokenTransfer(await roles.defaultAccount.getAddress(), 0, callData));
      });
    });

    describe("when called from the LINK token", () => {
      it("triggers the intended method", async () => {
        const callData = encodeOracleRequest(specId, to, fHash, 0, constants.HashZero);

        const tx = await link.transferAndCall(operator.address, 0, callData, {
          value: 0,
        });
        const receipt = await tx.wait();

        assert.equal(3, receipt.logs?.length);
      });

      describe("with no data", () => {
        it("reverts", async () => {
          await evmRevert(
            link.transferAndCall(operator.address, 0, "0x", {
              value: 0,
            }),
          );
        });
      });
    });

    describe("malicious requester", () => {
      let mock: Contract;
      let requester: Contract;
      const paymentAmount = toWei("1");

      beforeEach(async () => {
        mock = await maliciousRequesterFactory.connect(roles.defaultAccount).deploy(link.address, operator.address);
        await link.transfer(mock.address, paymentAmount);
      });

      it("cannot withdraw from oracle", async () => {
        const operatorOriginalBalance = await link.balanceOf(operator.address);
        const mockOriginalBalance = await link.balanceOf(mock.address);

        await evmRevert(mock.maliciousWithdraw());

        const operatorNewBalance = await link.balanceOf(operator.address);
        const mockNewBalance = await link.balanceOf(mock.address);

        bigNumEquals(operatorOriginalBalance, operatorNewBalance);
        bigNumEquals(mockNewBalance, mockOriginalBalance);
      });

      describe("if the requester tries to create a requestId for another contract", () => {
        it("the requesters ID will not match with the oracle contract", async () => {
          const tx = await mock.maliciousTargetConsumer(to);
          const receipt = await tx.wait();

          const mockRequestId = receipt.logs?.[0].data;
          const requestId = (receipt.events?.[0].args as any).requestId;
          assert.notEqual(mockRequestId, requestId);
        });

        it("the target requester can still create valid requests", async () => {
          requester = await basicConsumerFactory
            .connect(roles.defaultAccount)
            .deploy(link.address, operator.address, specId);
          await link.transfer(requester.address, paymentAmount);
          await mock.maliciousTargetConsumer(requester.address);
          await requester.requestEthereumPrice("USD", paymentAmount);
        });
      });
    });

    it("does not allow recursive calls of onTokenTransfer", async () => {
      const requestPayload = encodeOracleRequest(specId, to, fHash, 0, constants.HashZero);

      const ottSelector = operatorFactory.interface.getSighash("onTokenTransfer");
      const header =
        "000000000000000000000000c5fdf4076b8f3a5357c5e395ab970b5b54098fef" + // to
        "0000000000000000000000000000000000000000000000000000000000000539" + // amount
        "0000000000000000000000000000000000000000000000000000000000000060" + // offset
        "0000000000000000000000000000000000000000000000000000000000000136"; //   length

      const maliciousPayload = ottSelector + header + requestPayload.slice(2);

      await evmRevert(
        link.transferAndCall(operator.address, 0, maliciousPayload, {
          value: 0,
        }),
      );
    });
  });

  describe("#oracleRequest", () => {
    describe("when called through the LINK token", () => {
      const paid = 100;
      let log: providers.Log | undefined;
      let receipt: providers.TransactionReceipt;

      beforeEach(async () => {
        const args = encodeOracleRequest(specId, to, fHash, 1, constants.HashZero);
        const tx = await link.transferAndCall(operator.address, paid, args);
        receipt = await tx.wait();
        assert.equal(3, receipt?.logs?.length);

        log = receipt.logs && receipt.logs[2];
      });

      it("logs an event", async () => {
        assert.equal(operator.address, log?.address);

        assert.equal(log?.topics?.[1], specId);

        const req = decodeRunRequest(receipt?.logs?.[2]);
        assert.equal(await roles.defaultAccount.getAddress(), req.requester);
        bigNumEquals(paid, req.payment);
      });

      it("uses the expected event signature", async () => {
        // If updating this test, be sure to update models.RunLogTopic.
        const eventSignature = "0xd8d7ecc4800d25fa53ce0372f13a416d98907a7ef3d8d3bdd79cf4fe75529c65";
        assert.equal(eventSignature, log?.topics?.[0]);
      });

      it("does not allow the same requestId to be used twice", async () => {
        const args2 = encodeOracleRequest(specId, to, fHash, 1, constants.HashZero);
        await evmRevert(link.transferAndCall(operator.address, paid, args2));
      });

      describe("when called with a payload less than 2 EVM words + function selector", () => {
        it("throws an error", async () => {
          const funcSelector = operatorFactory.interface.getSighash("oracleRequest");
          const maliciousData = funcSelector + "0000000000000000000000000000000000000000000000000000000000000000000";
          await evmRevert(link.transferAndCall(operator.address, paid, maliciousData));
        });
      });

      describe("when called with a payload between 3 and 9 EVM words", () => {
        it("throws an error", async () => {
          const funcSelector = operatorFactory.interface.getSighash("oracleRequest");
          const maliciousData =
            funcSelector +
            "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001";
          await evmRevert(link.transferAndCall(operator.address, paid, maliciousData));
        });
      });
    });

    describe("when dataVersion is higher than 255", () => {
      it("throws an error", async () => {
        const paid = 100;
        const args = encodeOracleRequest(specId, to, fHash, 1, constants.HashZero, 256);
        await evmRevert(link.transferAndCall(operator.address, paid, args));
      });
    });

    describe("when not called through the LINK token", () => {
      it("reverts", async () => {
        await evmRevert(
          operator
            .connect(roles.oracleNode)
            .oracleRequest("0x0000000000000000000000000000000000000000", 0, specId, to, fHash, 1, 1, "0x"),
        );
      });
    });
  });

  describe("#requestOracleData", () => {
    describe("when called through the LINK token", () => {
      const paid = 100;
      let log: providers.Log | undefined;
      let receipt: providers.TransactionReceipt;

      beforeEach(async () => {
        const args = encodeRequestOracleData(specId, to, fHash, 1, constants.HashZero);
        const tx = await link.transferAndCall(operator.address, paid, args);
        receipt = await tx.wait();
        assert.equal(3, receipt?.logs?.length);

        log = receipt.logs && receipt.logs[2];
      });

      it("logs an event", async () => {
        assert.equal(operator.address, log?.address);

        assert.equal(log?.topics?.[1], specId);

        const req = decodeRunRequest(receipt?.logs?.[2]);
        assert.equal(await roles.defaultAccount.getAddress(), req.requester);
        bigNumEquals(paid, req.payment);
      });

      it("uses the expected event signature", async () => {
        // If updating this test, be sure to update models.RunLogTopic.
        const eventSignature = "0xd8d7ecc4800d25fa53ce0372f13a416d98907a7ef3d8d3bdd79cf4fe75529c65";
        assert.equal(eventSignature, log?.topics?.[0]);
      });

      it("does not allow the same requestId to be used twice", async () => {
        const args2 = encodeRequestOracleData(specId, to, fHash, 1, constants.HashZero);
        await evmRevert(link.transferAndCall(operator.address, paid, args2));
      });

      describe("when called with a payload less than 2 EVM words + function selector", () => {
        it("throws an error", async () => {
          const funcSelector = operatorFactory.interface.getSighash("oracleRequest");
          const maliciousData = funcSelector + "0000000000000000000000000000000000000000000000000000000000000000000";
          await evmRevert(link.transferAndCall(operator.address, paid, maliciousData));
        });
      });

      describe("when called with a payload between 3 and 9 EVM words", () => {
        it("throws an error", async () => {
          const funcSelector = operatorFactory.interface.getSighash("oracleRequest");
          const maliciousData =
            funcSelector +
            "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001";
          await evmRevert(link.transferAndCall(operator.address, paid, maliciousData));
        });
      });
    });

    describe("when dataVersion is higher than 255", () => {
      it("throws an error", async () => {
        const paid = 100;
        const args = encodeRequestOracleData(specId, to, fHash, 1, constants.HashZero, 256);
        await evmRevert(link.transferAndCall(operator.address, paid, args));
      });
    });

    describe("when not called through the LINK token", () => {
      it("reverts", async () => {
        await evmRevert(
          operator
            .connect(roles.oracleNode)
            .oracleRequest("0x0000000000000000000000000000000000000000", 0, specId, to, fHash, 1, 1, "0x"),
        );
      });
    });
  });

  describe("#fulfillOracleRequest", () => {
    const response = "Hi Mom!";
    let maliciousRequester: Contract;
    let basicConsumer: Contract;
    let maliciousConsumer: Contract;
    let gasGuzzlingConsumer: Contract;
    let request: ReturnType<typeof decodeRunRequest>;

    describe("gas guzzling consumer", () => {
      beforeEach(async () => {
        gasGuzzlingConsumer = await gasGuzzlingConsumerFactory
          .connect(roles.consumer)
          .deploy(link.address, operator.address, specId);
        const paymentAmount = toWei("1");
        await link.transfer(gasGuzzlingConsumer.address, paymentAmount);
        const tx = await gasGuzzlingConsumer.gassyRequestEthereumPrice(paymentAmount);
        const receipt = await tx.wait();
        request = decodeRunRequest(receipt.logs?.[3]);
      });

      it("emits an OracleResponse event", async () => {
        const fulfillParams = convertFufillParams(request, response);
        const tx = await operator.connect(roles.oracleNode).fulfillOracleRequest(...fulfillParams);
        const receipt = await tx.wait();
        assert.equal(receipt.events?.length, 1);
        const responseEvent = receipt.events?.[0];
        assert.equal(responseEvent?.event, "OracleResponse");
        assert.equal(responseEvent?.args?.[0], request.requestId);
      });
    });

    describe("cooperative consumer", () => {
      beforeEach(async () => {
        basicConsumer = await basicConsumerFactory
          .connect(roles.defaultAccount)
          .deploy(link.address, operator.address, specId);
        const paymentAmount = toWei("1");
        await link.transfer(basicConsumer.address, paymentAmount);
        const currency = "USD";
        const tx = await basicConsumer.requestEthereumPrice(currency, paymentAmount);
        const receipt = await tx.wait();
        request = decodeRunRequest(receipt.logs?.[3]);
      });

      describe("when called by an unauthorized node", () => {
        beforeEach(async () => {
          assert.equal(false, await operator.isAuthorizedSender(await roles.stranger.getAddress()));
        });

        it("raises an error", async () => {
          await evmRevert(
            operator.connect(roles.stranger).fulfillOracleRequest(...convertFufillParams(request, response)),
          );
        });
      });

      describe("when fulfilled with the wrong function", () => {
        let v7Consumer;
        beforeEach(async () => {
          v7Consumer = await v7ConsumerFactory
            .connect(roles.defaultAccount)
            .deploy(link.address, operator.address, specId);
          const paymentAmount = toWei("1");
          await link.transfer(v7Consumer.address, paymentAmount);
          const currency = "USD";
          const tx = await v7Consumer.requestEthereumPrice(currency, paymentAmount);
          const receipt = await tx.wait();
          request = decodeRunRequest(receipt.logs?.[3]);
        });

        it("raises an error", async () => {
          await evmRevert(
            operator.connect(roles.stranger).fulfillOracleRequest(...convertFufillParams(request, response)),
          );
        });
      });

      describe("when called by an authorized node", () => {
        it("raises an error if the request ID does not exist", async () => {
          request.requestId = ethers.utils.formatBytes32String("DOESNOTEXIST");
          await evmRevert(
            operator.connect(roles.oracleNode).fulfillOracleRequest(...convertFufillParams(request, response)),
          );
        });

        it("sets the value on the requested contract", async () => {
          await operator.connect(roles.oracleNode).fulfillOracleRequest(...convertFufillParams(request, response));

          const currentValue = await basicConsumer.currentPrice();
          assert.equal(response, ethers.utils.parseBytes32String(currentValue));
        });

        it("emits an OracleResponse event", async () => {
          const fulfillParams = convertFufillParams(request, response);
          const tx = await operator.connect(roles.oracleNode).fulfillOracleRequest(...fulfillParams);
          const receipt = await tx.wait();
          assert.equal(receipt.events?.length, 3);
          const responseEvent = receipt.events?.[0];
          assert.equal(responseEvent?.event, "OracleResponse");
          assert.equal(responseEvent?.args?.[0], request.requestId);
        });

        it("does not allow a request to be fulfilled twice", async () => {
          const response2 = response + " && Hello World!!";

          await operator.connect(roles.oracleNode).fulfillOracleRequest(...convertFufillParams(request, response));

          await evmRevert(
            operator.connect(roles.oracleNode).fulfillOracleRequest(...convertFufillParams(request, response2)),
          );

          const currentValue = await basicConsumer.currentPrice();
          assert.equal(response, ethers.utils.parseBytes32String(currentValue));
        });
      });

      describe("when the oracle does not provide enough gas", () => {
        // if updating this defaultGasLimit, be sure it matches with the
        // defaultGasLimit specified in store/tx_manager.go
        const defaultGasLimit = 500000;

        beforeEach(async () => {
          bigNumEquals(0, await operator.withdrawable());
        });

        it("does not allow the oracle to withdraw the payment", async () => {
          await evmRevert(
            operator.connect(roles.oracleNode).fulfillOracleRequest(
              ...convertFufillParams(request, response, {
                gasLimit: 70000,
              }),
            ),
          );

          bigNumEquals(0, await operator.withdrawable());
        });

        it(`${defaultGasLimit} is enough to pass the gas requirement`, async () => {
          await operator.connect(roles.oracleNode).fulfillOracleRequest(
            ...convertFufillParams(request, response, {
              gasLimit: defaultGasLimit,
            }),
          );

          bigNumEquals(request.payment, await operator.withdrawable());
        });
      });
    });

    describe("with a malicious requester", () => {
      beforeEach(async () => {
        const paymentAmount = toWei("1");
        maliciousRequester = await maliciousRequesterFactory
          .connect(roles.defaultAccount)
          .deploy(link.address, operator.address);
        await link.transfer(maliciousRequester.address, paymentAmount);
      });

      it("cannot cancel before the expiration", async () => {
        await evmRevert(
          maliciousRequester.maliciousRequestCancel(specId, ethers.utils.toUtf8Bytes("doesNothing(bytes32,bytes32)")),
        );
      });

      it("cannot call functions on the LINK token through callbacks", async () => {
        await evmRevert(
          maliciousRequester.request(specId, link.address, ethers.utils.toUtf8Bytes("transfer(address,uint256)")),
        );
      });

      describe("requester lies about amount of LINK sent", () => {
        it("the oracle uses the amount of LINK actually paid", async () => {
          const tx = await maliciousRequester.maliciousPrice(specId);
          const receipt = await tx.wait();
          const req = decodeRunRequest(receipt.logs?.[3]);

          assert(toWei("1").eq(req.payment));
        });
      });
    });

    describe("with a malicious consumer", () => {
      const paymentAmount = toWei("1");

      beforeEach(async () => {
        maliciousConsumer = await maliciousConsumerFactory
          .connect(roles.defaultAccount)
          .deploy(link.address, operator.address);
        await link.transfer(maliciousConsumer.address, paymentAmount);
      });

      describe("fails during fulfillment", () => {
        beforeEach(async () => {
          const tx = await maliciousConsumer.requestData(
            specId,
            ethers.utils.toUtf8Bytes("assertFail(bytes32,bytes32)"),
          );
          const receipt = await tx.wait();
          request = decodeRunRequest(receipt.logs?.[3]);
        });

        it("allows the oracle node to receive their payment", async () => {
          await operator.connect(roles.oracleNode).fulfillOracleRequest(...convertFufillParams(request, response));

          const balance = await link.balanceOf(await roles.oracleNode.getAddress());
          bigNumEquals(balance, 0);

          await operator.connect(roles.defaultAccount).withdraw(await roles.oracleNode.getAddress(), paymentAmount);

          const newBalance = await link.balanceOf(await roles.oracleNode.getAddress());
          bigNumEquals(paymentAmount, newBalance);
        });

        it("can't fulfill the data again", async () => {
          const response2 = "hack the planet 102";

          await operator.connect(roles.oracleNode).fulfillOracleRequest(...convertFufillParams(request, response));

          await evmRevert(
            operator.connect(roles.oracleNode).fulfillOracleRequest(...convertFufillParams(request, response2)),
          );
        });
      });

      describe("calls selfdestruct", () => {
        beforeEach(async () => {
          const tx = await maliciousConsumer.requestData(
            specId,
            ethers.utils.toUtf8Bytes("doesNothing(bytes32,bytes32)"),
          );
          const receipt = await tx.wait();
          request = decodeRunRequest(receipt.logs?.[3]);
          await maliciousConsumer.remove();
        });

        it("allows the oracle node to receive their payment", async () => {
          await operator.connect(roles.oracleNode).fulfillOracleRequest(...convertFufillParams(request, response));

          const balance = await link.balanceOf(await roles.oracleNode.getAddress());
          bigNumEquals(balance, 0);

          await operator.connect(roles.defaultAccount).withdraw(await roles.oracleNode.getAddress(), paymentAmount);
          const newBalance = await link.balanceOf(await roles.oracleNode.getAddress());
          bigNumEquals(paymentAmount, newBalance);
        });
      });

      describe("request is canceled during fulfillment", () => {
        beforeEach(async () => {
          const tx = await maliciousConsumer.requestData(
            specId,
            ethers.utils.toUtf8Bytes("cancelRequestOnFulfill(bytes32,bytes32)"),
          );
          const receipt = await tx.wait();
          request = decodeRunRequest(receipt.logs?.[3]);

          bigNumEquals(0, await link.balanceOf(maliciousConsumer.address));
        });

        it("allows the oracle node to receive their payment", async () => {
          await operator.connect(roles.oracleNode).fulfillOracleRequest(...convertFufillParams(request, response));

          const mockBalance = await link.balanceOf(maliciousConsumer.address);
          bigNumEquals(mockBalance, 0);

          const balance = await link.balanceOf(await roles.oracleNode.getAddress());
          bigNumEquals(balance, 0);

          await operator.connect(roles.defaultAccount).withdraw(await roles.oracleNode.getAddress(), paymentAmount);
          const newBalance = await link.balanceOf(await roles.oracleNode.getAddress());
          bigNumEquals(paymentAmount, newBalance);
        });

        it("can't fulfill the data again", async () => {
          const response2 = "hack the planet 102";

          await operator.connect(roles.oracleNode).fulfillOracleRequest(...convertFufillParams(request, response));

          await evmRevert(
            operator.connect(roles.oracleNode).fulfillOracleRequest(...convertFufillParams(request, response2)),
          );
        });
      });

      describe("tries to steal funds from node", () => {
        it("is not successful with call", async () => {
          const tx = await maliciousConsumer.requestData(
            specId,
            ethers.utils.toUtf8Bytes("stealEthCall(bytes32,bytes32)"),
          );
          const receipt = await tx.wait();
          request = decodeRunRequest(receipt.logs?.[3]);

          await operator.connect(roles.oracleNode).fulfillOracleRequest(...convertFufillParams(request, response));

          bigNumEquals(0, await ethers.provider.getBalance(maliciousConsumer.address));
        });

        it("is not successful with send", async () => {
          const tx = await maliciousConsumer.requestData(
            specId,
            ethers.utils.toUtf8Bytes("stealEthSend(bytes32,bytes32)"),
          );
          const receipt = await tx.wait();
          request = decodeRunRequest(receipt.logs?.[3]);

          await operator.connect(roles.oracleNode).fulfillOracleRequest(...convertFufillParams(request, response));
          bigNumEquals(0, await ethers.provider.getBalance(maliciousConsumer.address));
        });

        it("is not successful with transfer", async () => {
          const tx = await maliciousConsumer.requestData(
            specId,
            ethers.utils.toUtf8Bytes("stealEthTransfer(bytes32,bytes32)"),
          );
          const receipt = await tx.wait();
          request = decodeRunRequest(receipt.logs?.[3]);

          await operator.connect(roles.oracleNode).fulfillOracleRequest(...convertFufillParams(request, response));
          bigNumEquals(0, await ethers.provider.getBalance(maliciousConsumer.address));
        });
      });
    });
  });

  describe("#fulfillOracleRequest2", () => {
    describe("single word fulfils", () => {
      const response = "Hi mom!";
      const responseTypes = ["bytes32"];
      const responseValues = [toBytes32String(response)];
      let maliciousRequester: Contract;
      let basicConsumer: Contract;
      let maliciousConsumer: Contract;
      let gasGuzzlingConsumer: Contract;
      let request: ReturnType<typeof decodeRunRequest>;

      describe("gas guzzling consumer", () => {
        beforeEach(async () => {
          gasGuzzlingConsumer = await gasGuzzlingConsumerFactory
            .connect(roles.consumer)
            .deploy(link.address, operator.address, specId);
          const paymentAmount = toWei("1");
          await link.transfer(gasGuzzlingConsumer.address, paymentAmount);
          const tx = await gasGuzzlingConsumer.gassyRequestEthereumPrice(paymentAmount);
          const receipt = await tx.wait();
          request = decodeRunRequest(receipt.logs?.[3]);
        });

        it("emits an OracleResponse2 event", async () => {
          const fulfillParams = convertFulfill2Params(request, responseTypes, responseValues);
          const tx = await operator.connect(roles.oracleNode).fulfillOracleRequest2(...fulfillParams);
          const receipt = await tx.wait();
          assert.equal(receipt.events?.length, 1);
          const responseEvent = receipt.events?.[0];
          assert.equal(responseEvent?.event, "OracleResponse");
          assert.equal(responseEvent?.args?.[0], request.requestId);
        });
      });

      describe("cooperative consumer", () => {
        beforeEach(async () => {
          basicConsumer = await basicConsumerFactory
            .connect(roles.defaultAccount)
            .deploy(link.address, operator.address, specId);
          const paymentAmount = toWei("1");
          await link.transfer(basicConsumer.address, paymentAmount);
          const currency = "USD";
          const tx = await basicConsumer.requestEthereumPrice(currency, paymentAmount);
          const receipt = await tx.wait();
          request = decodeRunRequest(receipt.logs?.[3]);
        });

        describe("when called by an unauthorized node", () => {
          beforeEach(async () => {
            assert.equal(false, await operator.isAuthorizedSender(await roles.stranger.getAddress()));
          });

          it("raises an error", async () => {
            await evmRevert(
              operator
                .connect(roles.stranger)
                .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, responseValues)),
            );
          });
        });

        describe("when called by an authorized node", () => {
          it("raises an error if the request ID does not exist", async () => {
            request.requestId = ethers.utils.formatBytes32String("DOESNOTEXIST");
            await evmRevert(
              operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, responseValues)),
            );
          });

          it("sets the value on the requested contract", async () => {
            await operator
              .connect(roles.oracleNode)
              .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, responseValues));

            const currentValue = await basicConsumer.currentPrice();
            assert.equal(response, ethers.utils.parseBytes32String(currentValue));
          });

          it("emits an OracleResponse2 event", async () => {
            const fulfillParams = convertFulfill2Params(request, responseTypes, responseValues);
            const tx = await operator.connect(roles.oracleNode).fulfillOracleRequest2(...fulfillParams);
            const receipt = await tx.wait();
            assert.equal(receipt.events?.length, 3);
            const responseEvent = receipt.events?.[0];
            assert.equal(responseEvent?.event, "OracleResponse");
            assert.equal(responseEvent?.args?.[0], request.requestId);
          });

          it("does not allow a request to be fulfilled twice", async () => {
            const response2 = response + " && Hello World!!";
            const response2Values = [toBytes32String(response2)];
            await operator
              .connect(roles.oracleNode)
              .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, responseValues));

            await evmRevert(
              operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, response2Values)),
            );

            const currentValue = await basicConsumer.currentPrice();
            assert.equal(response, ethers.utils.parseBytes32String(currentValue));
          });
        });

        describe("when the oracle does not provide enough gas", () => {
          // if updating this defaultGasLimit, be sure it matches with the
          // defaultGasLimit specified in store/tx_manager.go
          const defaultGasLimit = 500000;

          beforeEach(async () => {
            bigNumEquals(0, await operator.withdrawable());
          });

          it("does not allow the oracle to withdraw the payment", async () => {
            await evmRevert(
              operator.connect(roles.oracleNode).fulfillOracleRequest2(
                ...convertFulfill2Params(request, responseTypes, responseValues, {
                  gasLimit: 70000,
                }),
              ),
            );

            bigNumEquals(0, await operator.withdrawable());
          });

          it(`${defaultGasLimit} is enough to pass the gas requirement`, async () => {
            await operator.connect(roles.oracleNode).fulfillOracleRequest2(
              ...convertFulfill2Params(request, responseTypes, responseValues, {
                gasLimit: defaultGasLimit,
              }),
            );

            bigNumEquals(request.payment, await operator.withdrawable());
          });
        });
      });

      describe("with a malicious requester", () => {
        beforeEach(async () => {
          const paymentAmount = toWei("1");
          maliciousRequester = await maliciousRequesterFactory
            .connect(roles.defaultAccount)
            .deploy(link.address, operator.address);
          await link.transfer(maliciousRequester.address, paymentAmount);
        });

        it("cannot cancel before the expiration", async () => {
          await evmRevert(
            maliciousRequester.maliciousRequestCancel(specId, ethers.utils.toUtf8Bytes("doesNothing(bytes32,bytes32)")),
          );
        });

        it("cannot call functions on the LINK token through callbacks", async () => {
          await evmRevert(
            maliciousRequester.request(specId, link.address, ethers.utils.toUtf8Bytes("transfer(address,uint256)")),
          );
        });

        describe("requester lies about amount of LINK sent", () => {
          it("the oracle uses the amount of LINK actually paid", async () => {
            const tx = await maliciousRequester.maliciousPrice(specId);
            const receipt = await tx.wait();
            const req = decodeRunRequest(receipt.logs?.[3]);

            assert(toWei("1").eq(req.payment));
          });
        });
      });

      describe("with a malicious consumer", () => {
        const paymentAmount = toWei("1");

        beforeEach(async () => {
          maliciousConsumer = await maliciousConsumerFactory
            .connect(roles.defaultAccount)
            .deploy(link.address, operator.address);
          await link.transfer(maliciousConsumer.address, paymentAmount);
        });

        describe("fails during fulfillment", () => {
          beforeEach(async () => {
            const tx = await maliciousConsumer.requestData(
              specId,
              ethers.utils.toUtf8Bytes("assertFail(bytes32,bytes32)"),
            );
            const receipt = await tx.wait();
            request = decodeRunRequest(receipt.logs?.[3]);
          });

          it("allows the oracle node to receive their payment", async () => {
            await operator
              .connect(roles.oracleNode)
              .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, responseValues));

            const balance = await link.balanceOf(await roles.oracleNode.getAddress());
            bigNumEquals(balance, 0);

            await operator.connect(roles.defaultAccount).withdraw(await roles.oracleNode.getAddress(), paymentAmount);

            const newBalance = await link.balanceOf(await roles.oracleNode.getAddress());
            bigNumEquals(paymentAmount, newBalance);
          });

          it("can't fulfill the data again", async () => {
            const response2 = "hack the planet 102";
            const response2Values = [toBytes32String(response2)];
            await operator
              .connect(roles.oracleNode)
              .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, responseValues));

            await evmRevert(
              operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, response2Values)),
            );
          });
        });

        describe("calls selfdestruct", () => {
          beforeEach(async () => {
            const tx = await maliciousConsumer.requestData(
              specId,
              ethers.utils.toUtf8Bytes("doesNothing(bytes32,bytes32)"),
            );
            const receipt = await tx.wait();
            request = decodeRunRequest(receipt.logs?.[3]);
            await maliciousConsumer.remove();
          });

          it("allows the oracle node to receive their payment", async () => {
            await operator
              .connect(roles.oracleNode)
              .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, responseValues));

            const balance = await link.balanceOf(await roles.oracleNode.getAddress());
            bigNumEquals(balance, 0);

            await operator.connect(roles.defaultAccount).withdraw(await roles.oracleNode.getAddress(), paymentAmount);
            const newBalance = await link.balanceOf(await roles.oracleNode.getAddress());
            bigNumEquals(paymentAmount, newBalance);
          });
        });

        describe("request is canceled during fulfillment", () => {
          beforeEach(async () => {
            const tx = await maliciousConsumer.requestData(
              specId,
              ethers.utils.toUtf8Bytes("cancelRequestOnFulfill(bytes32,bytes32)"),
            );
            const receipt = await tx.wait();
            request = decodeRunRequest(receipt.logs?.[3]);

            bigNumEquals(0, await link.balanceOf(maliciousConsumer.address));
          });

          it("allows the oracle node to receive their payment", async () => {
            await operator
              .connect(roles.oracleNode)
              .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, responseValues));

            const mockBalance = await link.balanceOf(maliciousConsumer.address);
            bigNumEquals(mockBalance, 0);

            const balance = await link.balanceOf(await roles.oracleNode.getAddress());
            bigNumEquals(balance, 0);

            await operator.connect(roles.defaultAccount).withdraw(await roles.oracleNode.getAddress(), paymentAmount);
            const newBalance = await link.balanceOf(await roles.oracleNode.getAddress());
            bigNumEquals(paymentAmount, newBalance);
          });

          it("can't fulfill the data again", async () => {
            const response2 = "hack the planet 102";
            const response2Values = [toBytes32String(response2)];

            await operator
              .connect(roles.oracleNode)
              .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, responseValues));

            await evmRevert(
              operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, response2Values)),
            );
          });
        });

        describe("tries to steal funds from node", () => {
          it("is not successful with call", async () => {
            const tx = await maliciousConsumer.requestData(
              specId,
              ethers.utils.toUtf8Bytes("stealEthCall(bytes32,bytes32)"),
            );
            const receipt = await tx.wait();
            request = decodeRunRequest(receipt.logs?.[3]);

            await operator
              .connect(roles.oracleNode)
              .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, responseValues));

            bigNumEquals(0, await ethers.provider.getBalance(maliciousConsumer.address));
          });

          it("is not successful with send", async () => {
            const tx = await maliciousConsumer.requestData(
              specId,
              ethers.utils.toUtf8Bytes("stealEthSend(bytes32,bytes32)"),
            );
            const receipt = await tx.wait();
            request = decodeRunRequest(receipt.logs?.[3]);

            await operator
              .connect(roles.oracleNode)
              .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, responseValues));
            bigNumEquals(0, await ethers.provider.getBalance(maliciousConsumer.address));
          });

          it("is not successful with transfer", async () => {
            const tx = await maliciousConsumer.requestData(
              specId,
              ethers.utils.toUtf8Bytes("stealEthTransfer(bytes32,bytes32)"),
            );
            const receipt = await tx.wait();
            request = decodeRunRequest(receipt.logs?.[3]);

            await operator
              .connect(roles.oracleNode)
              .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, responseValues));
            bigNumEquals(0, await ethers.provider.getBalance(maliciousConsumer.address));
          });
        });
      });
    });

    describe("multi word fulfils", () => {
      describe("one bytes parameter", () => {
        const response =
          "Lorem ipsum dolor sit amet, consectetur adipiscing elit.\
          Fusce euismod malesuada ligula, eget semper metus ultrices sit amet.";
        const responseTypes = ["bytes"];
        const responseValues = [stringToBytes(response)];
        let maliciousRequester: Contract;
        let multiConsumer: Contract;
        let maliciousConsumer: Contract;
        let gasGuzzlingConsumer: Contract;
        let request: ReturnType<typeof decodeRunRequest>;

        describe("gas guzzling consumer", () => {
          beforeEach(async () => {
            gasGuzzlingConsumer = await gasGuzzlingConsumerFactory
              .connect(roles.consumer)
              .deploy(link.address, operator.address, specId);
            const paymentAmount = toWei("1");
            await link.transfer(gasGuzzlingConsumer.address, paymentAmount);
            const tx = await gasGuzzlingConsumer.gassyMultiWordRequest(paymentAmount);
            const receipt = await tx.wait();
            request = decodeRunRequest(receipt.logs?.[3]);
          });

          it("emits an OracleResponse2 event", async () => {
            const fulfillParams = convertFulfill2Params(request, responseTypes, responseValues);
            const tx = await operator.connect(roles.oracleNode).fulfillOracleRequest2(...fulfillParams);
            const receipt = await tx.wait();
            assert.equal(receipt.events?.length, 1);
            const responseEvent = receipt.events?.[0];
            assert.equal(responseEvent?.event, "OracleResponse");
            assert.equal(responseEvent?.args?.[0], request.requestId);
          });
        });

        describe("cooperative consumer", () => {
          beforeEach(async () => {
            multiConsumer = await multiWordConsumerFactory
              .connect(roles.defaultAccount)
              .deploy(link.address, operator.address, specId);
            const paymentAmount = toWei("1");
            await link.transfer(multiConsumer.address, paymentAmount);
            const currency = "USD";
            const tx = await multiConsumer.requestEthereumPrice(currency, paymentAmount);
            const receipt = await tx.wait();
            request = decodeRunRequest(receipt.logs?.[3]);
          });

          describe("when called by an unauthorized node", () => {
            beforeEach(async () => {
              assert.equal(false, await operator.isAuthorizedSender(await roles.stranger.getAddress()));
            });

            it("raises an error", async () => {
              await evmRevert(
                operator
                  .connect(roles.stranger)
                  .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, responseValues)),
              );
            });
          });

          describe("when called by an authorized node", () => {
            it("raises an error if the request ID does not exist", async () => {
              request.requestId = ethers.utils.formatBytes32String("DOESNOTEXIST");
              await evmRevert(
                operator
                  .connect(roles.oracleNode)
                  .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, responseValues)),
              );
            });

            it("sets the value on the requested contract", async () => {
              await operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, responseValues));

              const currentValue = await multiConsumer.currentPrice();
              assert.equal(response, ethers.utils.toUtf8String(currentValue));
            });

            it("emits an OracleResponse2 event", async () => {
              const fulfillParams = convertFulfill2Params(request, responseTypes, responseValues);
              const tx = await operator.connect(roles.oracleNode).fulfillOracleRequest2(...fulfillParams);
              const receipt = await tx.wait();
              assert.equal(receipt.events?.length, 3);
              const responseEvent = receipt.events?.[0];
              assert.equal(responseEvent?.event, "OracleResponse");
              assert.equal(responseEvent?.args?.[0], request.requestId);
            });

            it("does not allow a request to be fulfilled twice", async () => {
              const response2 = response + " && Hello World!!";
              const response2Values = [stringToBytes(response2)];

              await operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, responseValues));

              await evmRevert(
                operator
                  .connect(roles.oracleNode)
                  .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, response2Values)),
              );

              const currentValue = await multiConsumer.currentPrice();
              assert.equal(response, ethers.utils.toUtf8String(currentValue));
            });
          });

          describe("when the oracle does not provide enough gas", () => {
            // if updating this defaultGasLimit, be sure it matches with the
            // defaultGasLimit specified in store/tx_manager.go
            const defaultGasLimit = 500000;

            beforeEach(async () => {
              bigNumEquals(0, await operator.withdrawable());
            });

            it("does not allow the oracle to withdraw the payment", async () => {
              await evmRevert(
                operator.connect(roles.oracleNode).fulfillOracleRequest2(
                  ...convertFulfill2Params(request, responseTypes, responseValues, {
                    gasLimit: 70000,
                  }),
                ),
              );

              bigNumEquals(0, await operator.withdrawable());
            });

            it(`${defaultGasLimit} is enough to pass the gas requirement`, async () => {
              await operator.connect(roles.oracleNode).fulfillOracleRequest2(
                ...convertFulfill2Params(request, responseTypes, responseValues, {
                  gasLimit: defaultGasLimit,
                }),
              );

              bigNumEquals(request.payment, await operator.withdrawable());
            });
          });
        });

        describe("with a malicious requester", () => {
          beforeEach(async () => {
            const paymentAmount = toWei("1");
            maliciousRequester = await maliciousRequesterFactory
              .connect(roles.defaultAccount)
              .deploy(link.address, operator.address);
            await link.transfer(maliciousRequester.address, paymentAmount);
          });

          it("cannot cancel before the expiration", async () => {
            await evmRevert(
              maliciousRequester.maliciousRequestCancel(
                specId,
                ethers.utils.toUtf8Bytes("doesNothing(bytes32,bytes32)"),
              ),
            );
          });

          it("cannot call functions on the LINK token through callbacks", async () => {
            await evmRevert(
              maliciousRequester.request(specId, link.address, ethers.utils.toUtf8Bytes("transfer(address,uint256)")),
            );
          });

          describe("requester lies about amount of LINK sent", () => {
            it("the oracle uses the amount of LINK actually paid", async () => {
              const tx = await maliciousRequester.maliciousPrice(specId);
              const receipt = await tx.wait();
              const req = decodeRunRequest(receipt.logs?.[3]);

              assert(toWei("1").eq(req.payment));
            });
          });
        });

        describe("with a malicious consumer", () => {
          const paymentAmount = toWei("1");

          beforeEach(async () => {
            maliciousConsumer = await maliciousMultiWordConsumerFactory
              .connect(roles.defaultAccount)
              .deploy(link.address, operator.address);
            await link.transfer(maliciousConsumer.address, paymentAmount);
          });

          describe("fails during fulfillment", () => {
            beforeEach(async () => {
              const tx = await maliciousConsumer.requestData(
                specId,
                ethers.utils.toUtf8Bytes("assertFail(bytes32,bytes32)"),
              );
              const receipt = await tx.wait();
              request = decodeRunRequest(receipt.logs?.[3]);
            });

            it("allows the oracle node to receive their payment", async () => {
              await operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, responseValues));

              const balance = await link.balanceOf(await roles.oracleNode.getAddress());
              bigNumEquals(balance, 0);

              await operator.connect(roles.defaultAccount).withdraw(await roles.oracleNode.getAddress(), paymentAmount);

              const newBalance = await link.balanceOf(await roles.oracleNode.getAddress());
              bigNumEquals(paymentAmount, newBalance);
            });

            it("can't fulfill the data again", async () => {
              const response2 = "hack the planet 102";
              const response2Values = [stringToBytes(response2)];
              await operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, responseValues));

              await evmRevert(
                operator
                  .connect(roles.oracleNode)
                  .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, response2Values)),
              );
            });
          });

          describe("calls selfdestruct", () => {
            beforeEach(async () => {
              const tx = await maliciousConsumer.requestData(
                specId,
                ethers.utils.toUtf8Bytes("doesNothing(bytes32,bytes32)"),
              );
              const receipt = await tx.wait();
              request = decodeRunRequest(receipt.logs?.[3]);
              await maliciousConsumer.remove();
            });

            it("allows the oracle node to receive their payment", async () => {
              await operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, responseValues));

              const balance = await link.balanceOf(await roles.oracleNode.getAddress());
              bigNumEquals(balance, 0);

              await operator.connect(roles.defaultAccount).withdraw(await roles.oracleNode.getAddress(), paymentAmount);
              const newBalance = await link.balanceOf(await roles.oracleNode.getAddress());
              bigNumEquals(paymentAmount, newBalance);
            });
          });

          describe("request is canceled during fulfillment", () => {
            beforeEach(async () => {
              const tx = await maliciousConsumer.requestData(
                specId,
                ethers.utils.toUtf8Bytes("cancelRequestOnFulfill(bytes32,bytes32)"),
              );
              const receipt = await tx.wait();
              request = decodeRunRequest(receipt.logs?.[3]);

              bigNumEquals(0, await link.balanceOf(maliciousConsumer.address));
            });

            it("allows the oracle node to receive their payment", async () => {
              await operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, responseValues));

              const mockBalance = await link.balanceOf(maliciousConsumer.address);
              bigNumEquals(mockBalance, 0);

              const balance = await link.balanceOf(await roles.oracleNode.getAddress());
              bigNumEquals(balance, 0);

              await operator.connect(roles.defaultAccount).withdraw(await roles.oracleNode.getAddress(), paymentAmount);
              const newBalance = await link.balanceOf(await roles.oracleNode.getAddress());
              bigNumEquals(paymentAmount, newBalance);
            });

            it("can't fulfill the data again", async () => {
              const response2 = "hack the planet 102";
              const response2Values = [stringToBytes(response2)];
              await operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, responseValues));

              await evmRevert(
                operator
                  .connect(roles.oracleNode)
                  .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, response2Values)),
              );
            });
          });

          describe("tries to steal funds from node", () => {
            it("is not successful with call", async () => {
              const tx = await maliciousConsumer.requestData(
                specId,
                ethers.utils.toUtf8Bytes("stealEthCall(bytes32,bytes32)"),
              );
              const receipt = await tx.wait();
              request = decodeRunRequest(receipt.logs?.[3]);

              await operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, responseValues));

              bigNumEquals(0, await ethers.provider.getBalance(maliciousConsumer.address));
            });

            it("is not successful with send", async () => {
              const tx = await maliciousConsumer.requestData(
                specId,
                ethers.utils.toUtf8Bytes("stealEthSend(bytes32,bytes32)"),
              );
              const receipt = await tx.wait();
              request = decodeRunRequest(receipt.logs?.[3]);

              await operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, responseValues));
              bigNumEquals(0, await ethers.provider.getBalance(maliciousConsumer.address));
            });

            it("is not successful with transfer", async () => {
              const tx = await maliciousConsumer.requestData(
                specId,
                ethers.utils.toUtf8Bytes("stealEthTransfer(bytes32,bytes32)"),
              );
              const receipt = await tx.wait();
              request = decodeRunRequest(receipt.logs?.[3]);

              await operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, responseValues));
              bigNumEquals(0, await ethers.provider.getBalance(maliciousConsumer.address));
            });
          });
        });
      });

      describe("multiple bytes32 parameters", () => {
        const response1 = "Hi mom!";
        const response2 = "Its me!";
        const responseTypes = ["bytes32", "bytes32"];
        const responseValues = [toBytes32String(response1), toBytes32String(response2)];
        let maliciousRequester: Contract;
        let multiConsumer: Contract;
        let maliciousConsumer: Contract;
        let gasGuzzlingConsumer: Contract;
        let request: ReturnType<typeof decodeRunRequest>;

        describe("gas guzzling consumer", () => {
          beforeEach(async () => {
            gasGuzzlingConsumer = await gasGuzzlingConsumerFactory
              .connect(roles.consumer)
              .deploy(link.address, operator.address, specId);
            const paymentAmount = toWei("1");
            await link.transfer(gasGuzzlingConsumer.address, paymentAmount);
            const tx = await gasGuzzlingConsumer.gassyMultiWordRequest(paymentAmount);
            const receipt = await tx.wait();
            request = decodeRunRequest(receipt.logs?.[3]);
          });

          it("emits an OracleResponse2 event", async () => {
            const fulfillParams = convertFulfill2Params(request, responseTypes, responseValues);
            const tx = await operator.connect(roles.oracleNode).fulfillOracleRequest2(...fulfillParams);
            const receipt = await tx.wait();
            assert.equal(receipt.events?.length, 1);
            const responseEvent = receipt.events?.[0];
            assert.equal(responseEvent?.event, "OracleResponse");
            assert.equal(responseEvent?.args?.[0], request.requestId);
          });
        });

        describe("cooperative consumer", () => {
          beforeEach(async () => {
            multiConsumer = await multiWordConsumerFactory
              .connect(roles.defaultAccount)
              .deploy(link.address, operator.address, specId);
            const paymentAmount = toWei("1");
            await link.transfer(multiConsumer.address, paymentAmount);
            const currency = "USD";
            const tx = await multiConsumer.requestMultipleParameters(currency, paymentAmount);
            const receipt = await tx.wait();
            request = decodeRunRequest(receipt.logs?.[3]);
          });

          describe("when called by an unauthorized node", () => {
            beforeEach(async () => {
              assert.equal(false, await operator.isAuthorizedSender(await roles.stranger.getAddress()));
            });

            it("raises an error", async () => {
              await evmRevert(
                operator
                  .connect(roles.stranger)
                  .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, responseValues)),
              );
            });
          });

          describe("when called by an authorized node", () => {
            it("raises an error if the request ID does not exist", async () => {
              request.requestId = ethers.utils.formatBytes32String("DOESNOTEXIST");
              await evmRevert(
                operator
                  .connect(roles.oracleNode)
                  .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, responseValues)),
              );
            });

            it("sets the value on the requested contract", async () => {
              await operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, responseValues));

              const firstValue = await multiConsumer.first();
              const secondValue = await multiConsumer.second();
              assert.equal(response1, ethers.utils.parseBytes32String(firstValue));
              assert.equal(response2, ethers.utils.parseBytes32String(secondValue));
            });

            it("emits an OracleResponse2 event", async () => {
              const fulfillParams = convertFulfill2Params(request, responseTypes, responseValues);
              const tx = await operator.connect(roles.oracleNode).fulfillOracleRequest2(...fulfillParams);
              const receipt = await tx.wait();
              assert.equal(receipt.events?.length, 3);
              const responseEvent = receipt.events?.[0];
              assert.equal(responseEvent?.event, "OracleResponse");
              assert.equal(responseEvent?.args?.[0], request.requestId);
            });

            it("does not allow a request to be fulfilled twice", async () => {
              const response3 = response2 + " && Hello World!!";
              const repeatedResponseValues = [toBytes32String(response2), toBytes32String(response3)];

              await operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, responseValues));

              await evmRevert(
                operator
                  .connect(roles.oracleNode)
                  .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, repeatedResponseValues)),
              );

              const firstValue = await multiConsumer.first();
              const secondValue = await multiConsumer.second();
              assert.equal(response1, ethers.utils.parseBytes32String(firstValue));
              assert.equal(response2, ethers.utils.parseBytes32String(secondValue));
            });
          });

          describe("when the oracle does not provide enough gas", () => {
            // if updating this defaultGasLimit, be sure it matches with the
            // defaultGasLimit specified in store/tx_manager.go
            const defaultGasLimit = 500000;

            beforeEach(async () => {
              bigNumEquals(0, await operator.withdrawable());
            });

            it("does not allow the oracle to withdraw the payment", async () => {
              await evmRevert(
                operator.connect(roles.oracleNode).fulfillOracleRequest2(
                  ...convertFulfill2Params(request, responseTypes, responseValues, {
                    gasLimit: 70000,
                  }),
                ),
              );

              bigNumEquals(0, await operator.withdrawable());
            });

            it(`${defaultGasLimit} is enough to pass the gas requirement`, async () => {
              await operator.connect(roles.oracleNode).fulfillOracleRequest2(
                ...convertFulfill2Params(request, responseTypes, responseValues, {
                  gasLimit: defaultGasLimit,
                }),
              );

              bigNumEquals(request.payment, await operator.withdrawable());
            });
          });
        });

        describe("with a malicious requester", () => {
          beforeEach(async () => {
            const paymentAmount = toWei("1");
            maliciousRequester = await maliciousRequesterFactory
              .connect(roles.defaultAccount)
              .deploy(link.address, operator.address);
            await link.transfer(maliciousRequester.address, paymentAmount);
          });

          it("cannot cancel before the expiration", async () => {
            await evmRevert(
              maliciousRequester.maliciousRequestCancel(
                specId,
                ethers.utils.toUtf8Bytes("doesNothing(bytes32,bytes32)"),
              ),
            );
          });

          it("cannot call functions on the LINK token through callbacks", async () => {
            await evmRevert(
              maliciousRequester.request(specId, link.address, ethers.utils.toUtf8Bytes("transfer(address,uint256)")),
            );
          });

          describe("requester lies about amount of LINK sent", () => {
            it("the oracle uses the amount of LINK actually paid", async () => {
              const tx = await maliciousRequester.maliciousPrice(specId);
              const receipt = await tx.wait();
              const req = decodeRunRequest(receipt.logs?.[3]);

              assert(toWei("1").eq(req.payment));
            });
          });
        });

        describe("with a malicious consumer", () => {
          const paymentAmount = toWei("1");

          beforeEach(async () => {
            maliciousConsumer = await maliciousMultiWordConsumerFactory
              .connect(roles.defaultAccount)
              .deploy(link.address, operator.address);
            await link.transfer(maliciousConsumer.address, paymentAmount);
          });

          describe("fails during fulfillment", () => {
            beforeEach(async () => {
              const tx = await maliciousConsumer.requestData(
                specId,
                ethers.utils.toUtf8Bytes("assertFail(bytes32,bytes32)"),
              );
              const receipt = await tx.wait();
              request = decodeRunRequest(receipt.logs?.[3]);
            });

            it("allows the oracle node to receive their payment", async () => {
              await operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, responseValues));

              const balance = await link.balanceOf(await roles.oracleNode.getAddress());
              bigNumEquals(balance, 0);

              await operator.connect(roles.defaultAccount).withdraw(await roles.oracleNode.getAddress(), paymentAmount);

              const newBalance = await link.balanceOf(await roles.oracleNode.getAddress());
              bigNumEquals(paymentAmount, newBalance);
            });

            it("can't fulfill the data again", async () => {
              const response3 = "hack the planet 102";
              const repeatedResponseValues = [toBytes32String(response2), toBytes32String(response3)];
              await operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, responseValues));

              await evmRevert(
                operator
                  .connect(roles.oracleNode)
                  .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, repeatedResponseValues)),
              );
            });
          });

          describe("calls selfdestruct", () => {
            beforeEach(async () => {
              const tx = await maliciousConsumer.requestData(
                specId,
                ethers.utils.toUtf8Bytes("doesNothing(bytes32,bytes32)"),
              );
              const receipt = await tx.wait();
              request = decodeRunRequest(receipt.logs?.[3]);
              await maliciousConsumer.remove();
            });

            it("allows the oracle node to receive their payment", async () => {
              await operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, responseValues));

              const balance = await link.balanceOf(await roles.oracleNode.getAddress());
              bigNumEquals(balance, 0);

              await operator.connect(roles.defaultAccount).withdraw(await roles.oracleNode.getAddress(), paymentAmount);
              const newBalance = await link.balanceOf(await roles.oracleNode.getAddress());
              bigNumEquals(paymentAmount, newBalance);
            });
          });

          describe("request is canceled during fulfillment", () => {
            beforeEach(async () => {
              const tx = await maliciousConsumer.requestData(
                specId,
                ethers.utils.toUtf8Bytes("cancelRequestOnFulfill(bytes32,bytes32)"),
              );
              const receipt = await tx.wait();
              request = decodeRunRequest(receipt.logs?.[3]);

              bigNumEquals(0, await link.balanceOf(maliciousConsumer.address));
            });

            it("allows the oracle node to receive their payment", async () => {
              await operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, responseValues));

              const mockBalance = await link.balanceOf(maliciousConsumer.address);
              bigNumEquals(mockBalance, 0);

              const balance = await link.balanceOf(await roles.oracleNode.getAddress());
              bigNumEquals(balance, 0);

              await operator.connect(roles.defaultAccount).withdraw(await roles.oracleNode.getAddress(), paymentAmount);
              const newBalance = await link.balanceOf(await roles.oracleNode.getAddress());
              bigNumEquals(paymentAmount, newBalance);
            });

            it("can't fulfill the data again", async () => {
              const response3 = "hack the planet 102";
              const repeatedResponseValues = [toBytes32String(response2), toBytes32String(response3)];
              await operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, responseValues));

              await evmRevert(
                operator
                  .connect(roles.oracleNode)
                  .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, repeatedResponseValues)),
              );
            });
          });

          describe("tries to steal funds from node", () => {
            it("is not successful with call", async () => {
              const tx = await maliciousConsumer.requestData(
                specId,
                ethers.utils.toUtf8Bytes("stealEthCall(bytes32,bytes32)"),
              );
              const receipt = await tx.wait();
              request = decodeRunRequest(receipt.logs?.[3]);

              await operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, responseValues));

              bigNumEquals(0, await ethers.provider.getBalance(maliciousConsumer.address));
            });

            it("is not successful with send", async () => {
              const tx = await maliciousConsumer.requestData(
                specId,
                ethers.utils.toUtf8Bytes("stealEthSend(bytes32,bytes32)"),
              );
              const receipt = await tx.wait();
              request = decodeRunRequest(receipt.logs?.[3]);

              await operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, responseValues));
              bigNumEquals(0, await ethers.provider.getBalance(maliciousConsumer.address));
            });

            it("is not successful with transfer", async () => {
              const tx = await maliciousConsumer.requestData(
                specId,
                ethers.utils.toUtf8Bytes("stealEthTransfer(bytes32,bytes32)"),
              );
              const receipt = await tx.wait();
              request = decodeRunRequest(receipt.logs?.[3]);

              await operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(...convertFulfill2Params(request, responseTypes, responseValues));
              bigNumEquals(0, await ethers.provider.getBalance(maliciousConsumer.address));
            });
          });
        });
      });
    });
  });

  describe("#withdraw", () => {
    describe("without reserving funds via oracleRequest", () => {
      it("does nothing", async () => {
        let balance = await link.balanceOf(await roles.oracleNode.getAddress());
        assert.equal(0, balance.toNumber());
        await evmRevert(
          operator.connect(roles.defaultAccount).withdraw(await roles.oracleNode.getAddress(), toWei("1")),
        );
        balance = await link.balanceOf(await roles.oracleNode.getAddress());
        assert.equal(0, balance.toNumber());
      });

      describe("recovering funds that were mistakenly sent", () => {
        const paid = 1;
        beforeEach(async () => {
          await link.transfer(operator.address, paid);
        });

        it("withdraws funds", async () => {
          const operatorBalanceBefore = await link.balanceOf(operator.address);
          const accountBalanceBefore = await link.balanceOf(await roles.defaultAccount.getAddress());

          await operator.connect(roles.defaultAccount).withdraw(await roles.defaultAccount.getAddress(), paid);

          const operatorBalanceAfter = await link.balanceOf(operator.address);
          const accountBalanceAfter = await link.balanceOf(await roles.defaultAccount.getAddress());

          const accountDifference = accountBalanceAfter.sub(accountBalanceBefore);
          const operatorDifference = operatorBalanceBefore.sub(operatorBalanceAfter);

          bigNumEquals(operatorDifference, paid);
          bigNumEquals(accountDifference, paid);
        });
      });
    });

    describe("reserving funds via oracleRequest", () => {
      const payment = 15;
      let request: ReturnType<typeof decodeRunRequest>;

      beforeEach(async () => {
        const mock = await getterSetterFactory.connect(roles.defaultAccount).deploy();
        const args = encodeOracleRequest(specId, mock.address, fHash, 0, constants.HashZero);
        const tx = await link.transferAndCall(operator.address, payment, args);
        const receipt = await tx.wait();
        assert.equal(3, receipt.logs?.length);
        request = decodeRunRequest(receipt.logs?.[2]);
      });

      describe("but not freeing funds w fulfillOracleRequest", () => {
        it("does not transfer funds", async () => {
          await evmRevert(
            operator.connect(roles.defaultAccount).withdraw(await roles.oracleNode.getAddress(), payment),
          );
          const balance = await link.balanceOf(await roles.oracleNode.getAddress());
          assert.equal(0, balance.toNumber());
        });

        describe("recovering funds that were mistakenly sent", () => {
          const paid = 1;
          beforeEach(async () => {
            await link.transfer(operator.address, paid);
          });

          it("withdraws funds", async () => {
            const operatorBalanceBefore = await link.balanceOf(operator.address);
            const accountBalanceBefore = await link.balanceOf(await roles.defaultAccount.getAddress());

            await operator.connect(roles.defaultAccount).withdraw(await roles.defaultAccount.getAddress(), paid);

            const operatorBalanceAfter = await link.balanceOf(operator.address);
            const accountBalanceAfter = await link.balanceOf(await roles.defaultAccount.getAddress());

            const accountDifference = accountBalanceAfter.sub(accountBalanceBefore);
            const operatorDifference = operatorBalanceBefore.sub(operatorBalanceAfter);

            bigNumEquals(operatorDifference, paid);
            bigNumEquals(accountDifference, paid);
          });
        });
      });

      describe("and freeing funds", () => {
        beforeEach(async () => {
          await operator
            .connect(roles.oracleNode)
            .fulfillOracleRequest(...convertFufillParams(request, "Hello World!"));
        });

        it("does not allow input greater than the balance", async () => {
          const originalOracleBalance = await link.balanceOf(operator.address);
          const originalStrangerBalance = await link.balanceOf(await roles.stranger.getAddress());
          const withdrawalAmount = payment + 1;

          assert.isAbove(withdrawalAmount, originalOracleBalance.toNumber());
          await evmRevert(
            operator.connect(roles.defaultAccount).withdraw(await roles.stranger.getAddress(), withdrawalAmount),
          );

          const newOracleBalance = await link.balanceOf(operator.address);
          const newStrangerBalance = await link.balanceOf(await roles.stranger.getAddress());

          assert.equal(originalOracleBalance.toNumber(), newOracleBalance.toNumber());
          assert.equal(originalStrangerBalance.toNumber(), newStrangerBalance.toNumber());
        });

        it("allows transfer of partial balance by owner to specified address", async () => {
          const partialAmount = 6;
          const difference = payment - partialAmount;
          await operator.connect(roles.defaultAccount).withdraw(await roles.stranger.getAddress(), partialAmount);
          const strangerBalance = await link.balanceOf(await roles.stranger.getAddress());
          const oracleBalance = await link.balanceOf(operator.address);
          assert.equal(partialAmount, strangerBalance.toNumber());
          assert.equal(difference, oracleBalance.toNumber());
        });

        it("allows transfer of entire balance by owner to specified address", async () => {
          await operator.connect(roles.defaultAccount).withdraw(await roles.stranger.getAddress(), payment);
          const balance = await link.balanceOf(await roles.stranger.getAddress());
          assert.equal(payment, balance.toNumber());
        });

        it("does not allow a transfer of funds by non-owner", async () => {
          await evmRevert(operator.connect(roles.stranger).withdraw(await roles.stranger.getAddress(), payment));
          const balance = await link.balanceOf(await roles.stranger.getAddress());
          assert.isTrue(ethers.constants.Zero.eq(balance));
        });

        describe("recovering funds that were mistakenly sent", () => {
          const paid = 1;
          beforeEach(async () => {
            await link.transfer(operator.address, paid);
          });

          it("withdraws funds", async () => {
            const operatorBalanceBefore = await link.balanceOf(operator.address);
            const accountBalanceBefore = await link.balanceOf(await roles.defaultAccount.getAddress());

            await operator.connect(roles.defaultAccount).withdraw(await roles.defaultAccount.getAddress(), paid);

            const operatorBalanceAfter = await link.balanceOf(operator.address);
            const accountBalanceAfter = await link.balanceOf(await roles.defaultAccount.getAddress());

            const accountDifference = accountBalanceAfter.sub(accountBalanceBefore);
            const operatorDifference = operatorBalanceBefore.sub(operatorBalanceAfter);

            bigNumEquals(operatorDifference, paid);
            bigNumEquals(accountDifference, paid);
          });
        });
      });
    });
  });

  describe("#withdrawable", () => {
    let request: ReturnType<typeof decodeRunRequest>;
    const amount = toWei("1");

    beforeEach(async () => {
      const mock = await getterSetterFactory.connect(roles.defaultAccount).deploy();
      const args = encodeOracleRequest(specId, mock.address, fHash, 0, constants.HashZero);
      const tx = await link.transferAndCall(operator.address, amount, args);
      const receipt = await tx.wait();
      assert.equal(3, receipt.logs?.length);
      request = decodeRunRequest(receipt.logs?.[2]);
      await operator.connect(roles.oracleNode).fulfillOracleRequest(...convertFufillParams(request, "Hello World!"));
    });

    it("returns the correct value", async () => {
      const withdrawAmount = await operator.withdrawable();
      bigNumEquals(withdrawAmount, request.payment);
    });

    describe("funds that were mistakenly sent", () => {
      const paid = 1;
      beforeEach(async () => {
        await link.transfer(operator.address, paid);
      });

      it("returns the correct value", async () => {
        const withdrawAmount = await operator.withdrawable();

        const expectedAmount = amount.add(paid);
        bigNumEquals(withdrawAmount, expectedAmount);
      });
    });
  });

  describe("#ownerTransferAndCall", () => {
    let operator2: Contract;
    let args: string;
    let to: string;
    const startingBalance = 1000;
    const payment = 20;

    beforeEach(async () => {
      operator2 = await operatorFactory
        .connect(roles.oracleNode2)
        .deploy(link.address, await roles.oracleNode2.getAddress());
      to = operator2.address;
      args = encodeOracleRequest(
        specId,
        operator.address,
        operatorFactory.interface.getSighash("fulfillOracleRequest"),
        1,
        constants.HashZero,
      );
    });

    describe("when called by a non-owner", () => {
      it("reverts with owner error message", async () => {
        await link.transfer(operator.address, startingBalance);
        await evmRevert(
          operator.connect(roles.stranger).ownerTransferAndCall(to, payment, args),
          "Only callable by owner",
        );
      });
    });

    describe("when called by the owner", () => {
      beforeEach(async () => {
        await link.transfer(operator.address, startingBalance);
      });

      describe("without sufficient funds in contract", () => {
        it("reverts with funds message", async () => {
          const tooMuch = startingBalance * 2;
          await evmRevert(
            operator.connect(roles.defaultAccount).ownerTransferAndCall(to, tooMuch, args),
            "Amount requested is greater than withdrawable balance",
          );
        });
      });

      describe("with sufficient funds", () => {
        let tx: ContractTransaction;
        let receipt: ContractReceipt;
        let requesterBalanceBefore: BigNumber;
        let requesterBalanceAfter: BigNumber;
        let receiverBalanceBefore: BigNumber;
        let receiverBalanceAfter: BigNumber;

        before(async () => {
          requesterBalanceBefore = await link.balanceOf(operator.address);
          receiverBalanceBefore = await link.balanceOf(operator2.address);
          tx = await operator.connect(roles.defaultAccount).ownerTransferAndCall(to, payment, args);
          receipt = await tx.wait();
          requesterBalanceAfter = await link.balanceOf(operator.address);
          receiverBalanceAfter = await link.balanceOf(operator2.address);
        });

        it("emits an event", async () => {
          assert.equal(3, receipt.logs?.length);
          expect(tx).to.emit(link, "Transfer").withArgs(operator.address, to, payment, args);
        });

        it("transfers the tokens", async () => {
          bigNumEquals(requesterBalanceBefore.sub(requesterBalanceAfter), payment);
          bigNumEquals(receiverBalanceAfter.sub(receiverBalanceBefore), payment);
        });
      });
    });
  });

  describe("#cancelOracleRequest", () => {
    describe("with no pending requests", () => {
      it("fails", async () => {
        const fakeRequest: RunRequest = {
          requestId: ethers.utils.formatBytes32String("1337"),
          payment: "0",
          callbackFunc: getterSetterFactory.interface.getSighash("requestedBytes32"),
          expiration: "999999999999",

          callbackAddr: "",
          data: Buffer.from(""),
          dataVersion: 0,
          specId: "",
          requester: "",
          topic: "",
        };
        await increaseTime5Minutes(ethers.provider);

        await evmRevert(operator.connect(roles.stranger).cancelOracleRequest(...convertCancelParams(fakeRequest)));
      });
    });

    describe("with a pending request", () => {
      const startingBalance = 100;
      let request: ReturnType<typeof decodeRunRequest>;
      let receipt: providers.TransactionReceipt;

      beforeEach(async () => {
        const requestAmount = 20;

        await link.transfer(await roles.consumer.getAddress(), startingBalance);

        const args = encodeOracleRequest(specId, await roles.consumer.getAddress(), fHash, 1, constants.HashZero);
        const tx = await link.connect(roles.consumer).transferAndCall(operator.address, requestAmount, args);
        receipt = await tx.wait();

        assert.equal(3, receipt.logs?.length);
        request = decodeRunRequest(receipt.logs?.[2]);
      });

      it("has correct initial balances", async () => {
        const oracleBalance = await link.balanceOf(operator.address);
        bigNumEquals(request.payment, oracleBalance);

        const consumerAmount = await link.balanceOf(await roles.consumer.getAddress());
        assert.equal(startingBalance - Number(request.payment), consumerAmount.toNumber());
      });

      describe("from a stranger", () => {
        it("fails", async () => {
          await evmRevert(operator.connect(roles.consumer).cancelOracleRequest(...convertCancelParams(request)));
        });
      });

      describe("from the requester", () => {
        it("refunds the correct amount", async () => {
          await increaseTime5Minutes(ethers.provider);
          await operator.connect(roles.consumer).cancelOracleRequest(...convertCancelParams(request));
          const balance = await link.balanceOf(await roles.consumer.getAddress());

          assert.equal(startingBalance, balance.toNumber()); // 100
        });

        it("triggers a cancellation event", async () => {
          await increaseTime5Minutes(ethers.provider);
          const tx = await operator.connect(roles.consumer).cancelOracleRequest(...convertCancelParams(request));
          const receipt = await tx.wait();

          assert.equal(receipt.logs?.length, 2);
          assert.equal(request.requestId, receipt.logs?.[0].topics[1]);
        });

        it("fails when called twice", async () => {
          await increaseTime5Minutes(ethers.provider);
          await operator.connect(roles.consumer).cancelOracleRequest(...convertCancelParams(request));

          await evmRevert(operator.connect(roles.consumer).cancelOracleRequest(...convertCancelParams(request)));
        });
      });
    });
  });

  describe("#ownerForward", () => {
    let bytes: string;
    let payload: string;
    let mock: Contract;

    beforeEach(async () => {
      bytes = ethers.utils.hexlify(ethers.utils.randomBytes(100));
      payload = getterSetterFactory.interface.encodeFunctionData(
        getterSetterFactory.interface.getFunction("setBytes"),
        [bytes],
      );
      mock = await getterSetterFactory.connect(roles.defaultAccount).deploy();
    });

    describe("when called by a non-owner", () => {
      it("reverts", async () => {
        await evmRevert(operator.connect(roles.stranger).ownerForward(mock.address, payload));
      });
    });

    describe("when called by owner", () => {
      describe("when attempting to forward to the link token", () => {
        it("reverts", async () => {
          const sighash = linkTokenFactory.interface.getSighash("name");
          await evmRevert(
            operator.connect(roles.defaultAccount).ownerForward(link.address, sighash),
            "Cannot call to LINK",
          );
        });
      });

      describe("when forwarding to any other address", () => {
        it("forwards the data", async () => {
          const tx = await operator.connect(roles.defaultAccount).ownerForward(mock.address, payload);
          await tx.wait();
          assert.equal(await mock.getBytes(), bytes);
        });

        it("reverts when sending to a non-contract address", async () => {
          await evmRevert(
            operator.connect(roles.defaultAccount).ownerForward(zeroAddress, payload),
            "Must forward to a contract",
          );
        });

        it("perceives the message is sent by the Operator", async () => {
          const tx = await operator.connect(roles.defaultAccount).ownerForward(mock.address, payload);
          const receipt = await tx.wait();
          const log: any = receipt.logs?.[0];
          const logData = mock.interface.decodeEventLog(mock.interface.getEvent("SetBytes"), log.data, log.topics);
          assert.equal(ethers.utils.getAddress(logData.from), operator.address);
        });
      });
    });
  });
});
