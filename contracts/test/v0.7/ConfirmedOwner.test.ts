import { ethers } from "hardhat";
import { publicAbi } from "../helpers";
import { assert, expect } from "chai";
import { Contract, ContractFactory, Signer } from "ethers";
import { Personas, getUsers } from "../setup";
import { evmRevert } from "../matchers";

let personas: Personas;
let owner: Signer;
let nonOwner: Signer;
let newOwner: Signer;

let confirmedOwnerTestHelperFactory: ContractFactory;

before(async () => {
  const users = await getUsers();
  personas = users.personas;
  owner = personas.Carol;
  nonOwner = personas.Neil;
  newOwner = personas.Ned;
  confirmedOwnerTestHelperFactory = await ethers.getContractFactory("ConfirmedOwnerTestHelper", owner);
});

describe("ConfirmedOwner", () => {
  let confirmedOwner: Contract;

  beforeEach(async () => {
    confirmedOwner = await confirmedOwnerTestHelperFactory.connect(owner).deploy();
  });

  it("has a limited public interface", () => {
    publicAbi(confirmedOwner, [
      "acceptOwnership",
      "owner",
      "transferOwnership",
      // test helper public methods
      "modifierOnlyOwner",
    ]);
  });

  describe("#constructor", () => {
    it("assigns ownership to the deployer", async () => {
      const [actual, expected] = await Promise.all([owner.getAddress(), confirmedOwner.owner()]);

      assert.equal(actual, expected);
    });
  });

  describe("#transferOwnership", () => {
    describe("when called by an owner", () => {
      it("emits a log", async () => {
        const tx = await confirmedOwner.connect(owner).transferOwnership(await newOwner.getAddress());

        await expect(tx)
          .to.emit(confirmedOwner, "OwnershipTransferRequested")
          .withArgs(await owner.getAddress(), await newOwner.getAddress());
      });

      it("does not allow ownership transfer to self", async () => {
        await evmRevert(
          confirmedOwner.connect(owner).transferOwnership(await owner.getAddress()),
          "Cannot transfer to self",
        );
      });
    });
  });

  describe("when called by anyone but the owner", () => {
    it("successfully calls the method", async () =>
      evmRevert(confirmedOwner.connect(nonOwner).transferOwnership(await newOwner.getAddress())));
  });

  describe("#acceptOwnership", () => {
    describe("after #transferOwnership has been called", () => {
      beforeEach(async () => {
        await confirmedOwner.connect(owner).transferOwnership(await newOwner.getAddress());
      });

      it("allows the recipient to call it", async () => {
        const tx = await confirmedOwner.connect(newOwner).acceptOwnership();
        await expect(tx)
          .to.emit(confirmedOwner, "OwnershipTransferred")
          .withArgs(await owner.getAddress(), await newOwner.getAddress());
      });

      it("does not allow a non-recipient to call it", () =>
        evmRevert(confirmedOwner.connect(nonOwner).acceptOwnership()));
    });
  });
});
