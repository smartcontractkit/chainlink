import { ethers } from "hardhat";
import { assert, expect } from "chai";
import { BigNumber, constants, Contract, ContractFactory, ContractTransaction } from "ethers";
import { getUsers, Personas, Roles } from "../test-helpers/setup";
import {
  evmWordToAddress,
  getLog,
  publicAbi,
  toBytes32String,
  toWei,
  numToBytes32,
  getLogs,
} from "../test-helpers/helpers";

let roles: Roles;
let personas: Personas;
let linkTokenFactory: ContractFactory;
let vrfCoordinatorMockFactory: ContractFactory;
let vrfD20Factory: ContractFactory;

before(async () => {
  const users = await getUsers();

  roles = users.roles;
  personas = users.personas;
  linkTokenFactory = await ethers.getContractFactory("LinkToken", roles.defaultAccount);
  vrfCoordinatorMockFactory = await ethers.getContractFactory("VRFCoordinatorMock", roles.defaultAccount);
  vrfD20Factory = await ethers.getContractFactory("VRFD20", roles.defaultAccount);
});

describe("VRFD20", () => {
  const deposit = toWei("1");
  const fee = toWei("0.1");
  const keyHash = toBytes32String("keyHash");

  let link: Contract;
  let vrfCoordinator: Contract;
  let vrfD20: Contract;

  beforeEach(async () => {
    link = await linkTokenFactory.connect(roles.defaultAccount).deploy();
    vrfCoordinator = await vrfCoordinatorMockFactory.connect(roles.defaultAccount).deploy(link.address);
    vrfD20 = await vrfD20Factory
      .connect(roles.defaultAccount)
      .deploy(vrfCoordinator.address, link.address, keyHash, fee);
    await link.transfer(vrfD20.address, deposit);
  });

  it("has a limited public interface", () => {
    publicAbi(vrfD20, [
      // Owned
      "acceptOwnership",
      "owner",
      "transferOwnership",
      //VRFConsumerBase
      "rawFulfillRandomness",
      // VRFD20
      "rollDice",
      "house",
      "withdrawLINK",
      "keyHash",
      "fee",
      "setKeyHash",
      "setFee",
    ]);
  });

  describe("#withdrawLINK", () => {
    describe("failure", () => {
      it("reverts when called by a non-owner", async () => {
        await expect(
          vrfD20.connect(roles.stranger).withdrawLINK(await roles.stranger.getAddress(), deposit),
        ).to.be.revertedWith("Only callable by owner");
      });

      it("reverts when not enough LINK in the contract", async () => {
        const withdrawAmount = deposit.mul(2);
        await expect(
          vrfD20.connect(roles.defaultAccount).withdrawLINK(await roles.defaultAccount.getAddress(), withdrawAmount),
        ).to.be.reverted;
      });
    });

    describe("success", () => {
      it("withdraws LINK", async () => {
        const startingAmount = await link.balanceOf(await roles.defaultAccount.getAddress());
        const expectedAmount = BigNumber.from(startingAmount).add(deposit);
        await vrfD20.connect(roles.defaultAccount).withdrawLINK(await roles.defaultAccount.getAddress(), deposit);
        const actualAmount = await link.balanceOf(await roles.defaultAccount.getAddress());
        assert.equal(actualAmount.toString(), expectedAmount.toString());
      });
    });
  });

  describe("#setKeyHash", () => {
    const newHash = toBytes32String("newhash");

    describe("failure", () => {
      it("reverts when called by a non-owner", async () => {
        await expect(vrfD20.connect(roles.stranger).setKeyHash(newHash)).to.be.revertedWith("Only callable by owner");
      });
    });

    describe("success", () => {
      it("sets the key hash", async () => {
        await vrfD20.setKeyHash(newHash);
        const actualHash = await vrfD20.keyHash();
        assert.equal(actualHash, newHash);
      });
    });
  });

  describe("#setFee", () => {
    const newFee = 1234;

    describe("failure", () => {
      it("reverts when called by a non-owner", async () => {
        await expect(vrfD20.connect(roles.stranger).setFee(newFee)).to.be.revertedWith("Only callable by owner");
      });
    });

    describe("success", () => {
      it("sets the fee", async () => {
        await vrfD20.setFee(newFee);
        const actualFee = await vrfD20.fee();
        assert.equal(actualFee.toString(), newFee.toString());
      });
    });
  });

  describe("#house", () => {
    describe("failure", () => {
      it("reverts when dice not rolled", async () => {
        await expect(vrfD20.house(await personas.Nancy.getAddress())).to.be.revertedWith("Dice not rolled");
      });

      it("reverts when dice roll is in progress", async () => {
        await vrfD20.rollDice(await personas.Nancy.getAddress());
        await expect(vrfD20.house(await personas.Nancy.getAddress())).to.be.revertedWith("Roll in progress");
      });
    });

    describe("success", () => {
      it("returns the correct house", async () => {
        const randomness = 98765;
        const expectedHouse = "Martell";
        const tx = await vrfD20.rollDice(await personas.Nancy.getAddress());
        const log = await getLog(tx, 3);
        const eventRequestId = log?.topics?.[1];
        await vrfCoordinator.callBackWithRandomness(eventRequestId, randomness, vrfD20.address);
        const response = await vrfD20.house(await personas.Nancy.getAddress());
        assert.equal(response.toString(), expectedHouse);
      });
    });
  });

  describe("#rollDice", () => {
    describe("success", () => {
      let tx: ContractTransaction;
      beforeEach(async () => {
        tx = await vrfD20.rollDice(await personas.Nancy.getAddress());
      });

      it("emits a RandomnessRequest event from the VRFCoordinator", async () => {
        const log = await getLog(tx, 2);
        const topics = log?.topics;
        assert.equal(evmWordToAddress(topics?.[1]), vrfD20.address);
        assert.equal(topics?.[2], keyHash);
        assert.equal(topics?.[3], constants.HashZero);
      });
    });

    describe("failure", () => {
      it("reverts when LINK balance is zero", async () => {
        const vrfD202 = await vrfD20Factory
          .connect(roles.defaultAccount)
          .deploy(vrfCoordinator.address, link.address, keyHash, fee);
        await expect(vrfD202.rollDice(await personas.Nancy.getAddress())).to.be.revertedWith(
          "Not enough LINK to pay fee",
        );
      });

      it("reverts when called by a non-owner", async () => {
        await expect(vrfD20.connect(roles.stranger).rollDice(await personas.Nancy.getAddress())).to.be.revertedWith(
          "Only callable by owner",
        );
      });

      it("reverts when the roller rolls more than once", async () => {
        await vrfD20.rollDice(await personas.Nancy.getAddress());
        await expect(vrfD20.rollDice(await personas.Nancy.getAddress())).to.be.revertedWith("Already rolled");
      });
    });
  });

  describe("#fulfillRandomness", () => {
    const randomness = 98765;
    const expectedModResult = (randomness % 20) + 1;
    const expectedHouse = "Martell";
    let eventRequestId: string;
    beforeEach(async () => {
      const tx = await vrfD20.rollDice(await personas.Nancy.getAddress());
      const log = await getLog(tx, 3);
      eventRequestId = log?.topics?.[1];
    });

    describe("success", () => {
      let tx: ContractTransaction;
      beforeEach(async () => {
        tx = await vrfCoordinator.callBackWithRandomness(eventRequestId, randomness, vrfD20.address);
      });

      it("emits a DiceLanded event", async () => {
        const log = await getLog(tx, 0);
        assert.equal(log?.topics[1], eventRequestId);
        assert.equal(log?.topics[2], numToBytes32(expectedModResult));
      });

      it("sets the correct dice roll result", async () => {
        const response = await vrfD20.house(await personas.Nancy.getAddress());
        assert.equal(response.toString(), expectedHouse);
      });

      it("allows someone else to roll", async () => {
        const secondRandomness = 55555;
        tx = await vrfD20.rollDice(await personas.Ned.getAddress());
        const log = await getLog(tx, 3);
        eventRequestId = log?.topics?.[1];
        tx = await vrfCoordinator.callBackWithRandomness(eventRequestId, secondRandomness, vrfD20.address);
      });
    });

    describe("failure", () => {
      it("does not fulfill when fulfilled by the wrong VRFcoordinator", async () => {
        const vrfCoordinator2 = await vrfCoordinatorMockFactory.connect(roles.defaultAccount).deploy(link.address);

        const tx = await vrfCoordinator2.callBackWithRandomness(eventRequestId, randomness, vrfD20.address);
        const logs = await getLogs(tx);
        assert.equal(logs.length, 0);
      });
    });
  });
});
