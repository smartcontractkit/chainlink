import { ethers } from "hardhat";
import { publicAbi } from "../test-helpers/helpers";
import { assert, expect } from "chai";
import { Contract, ContractFactory } from "ethers";
import { Personas, getUsers } from "../test-helpers/setup";

let personas: Personas;

let controllerFactory: ContractFactory;
let l2ExplicitFlagFactory: ContractFactory;
let consumerFactory: ContractFactory;

let controller: Contract;
let flags: Contract;
let consumer: Contract;

before(async () => {
  personas = (await getUsers()).personas;
  controllerFactory = await ethers.getContractFactory("SimpleWriteAccessController", personas.Nelly);
  consumerFactory = await ethers.getContractFactory("L2ExplicitFlagTestHelper", personas.Nelly);
  l2ExplicitFlagFactory = await ethers.getContractFactory("L2ExplicitFlag", personas.Nelly);
});

describe("L2 Explicit Flag", () => {
  beforeEach(async () => {
    controller = await controllerFactory.deploy();
    flags = await l2ExplicitFlagFactory.deploy(controller.address);
    await flags.disableAccessCheck();
    consumer = await consumerFactory.deploy(flags.address);
  });

  it("has a limited public interface", async () => {
    publicAbi(flags, [
      "isRaised",
      "lowerFlag",
      "raiseFlag",
      "raisingAccessController",
      "setRaisingAccessController",
      // Ownable methods:
      "acceptOwnership",
      "owner",
      "transferOwnership",
      // AccessControl methods:
      "addAccess",
      "disableAccessCheck",
      "enableAccessCheck",
      "removeAccess",
      "checkEnabled",
      "hasAccess",
    ]);
  });

  describe("#raiseFlag", () => {
    describe("when called by the owner", () => {
      it("updates the warning flag", async () => {
        assert.equal(false, await flags.isRaised());

        await flags.connect(personas.Nelly).raiseFlag();

        assert.equal(true, await flags.isRaised());
      });

      it("emits an event log", async () => {
        await expect(flags.connect(personas.Nelly).raiseFlag()).to.emit(flags, "FlagRaised").withArgs(true);
      });

      describe("if the flag has already been raised", () => {
        beforeEach(async () => {
          await flags.connect(personas.Nelly).raiseFlag();
        });

        it("emits an event log", async () => {
          const tx = await flags.connect(personas.Nelly).raiseFlag();
          const receipt = await tx.wait();
          assert.equal(0, receipt.events?.length);
        });
      });
    });

    describe("when called by an enabled setter", () => {
      beforeEach(async () => {
        await controller.connect(personas.Nelly).addAccess(await personas.Neil.getAddress());
      });

      it("raise the flag", async () => {
        await flags.connect(personas.Neil).raiseFlag(), assert.equal(true, await flags.isRaised());
      });
    });

    describe("when called by a non-enabled setter", () => {
      it("reverts", async () => {
        await expect(flags.connect(personas.Neil).raiseFlag()).to.be.revertedWith("Not allowed to raise the flag");
      });
    });

    describe("when called when there is no raisingAccessController", () => {
      beforeEach(async () => {
        await expect(
          flags.connect(personas.Nelly).setRaisingAccessController("0x0000000000000000000000000000000000000000"),
        ).to.emit(flags, "RaisingAccessControllerUpdated");
        assert.equal("0x0000000000000000000000000000000000000000", await flags.raisingAccessController());
      });

      it("succeeds for the owner", async () => {
        await flags.connect(personas.Nelly).raiseFlag();
        assert.equal(true, await flags.isRaised());
      });

      it("reverts for non-owner", async () => {
        await expect(flags.connect(personas.Neil).raiseFlag()).to.be.reverted;
      });
    });
  });

  describe("#lowerFlag", () => {
    beforeEach(async () => {
      await flags.connect(personas.Nelly).raiseFlag();
    });

    describe("when called by the owner", () => {
      it("updates the warning flag", async () => {
        assert.equal(true, await flags.isRaised());

        await flags.connect(personas.Nelly).lowerFlag();

        assert.equal(false, await flags.isRaised());
      });

      it("emits an event log", async () => {
        await expect(flags.connect(personas.Nelly).lowerFlag()).to.emit(flags, "FlagLowered").withArgs(false);
      });

      describe("if the flag has already been lowered", () => {
        beforeEach(async () => {
          await flags.connect(personas.Nelly).lowerFlag();
        });

        it("emits an event log", async () => {
          const tx = await flags.connect(personas.Nelly).lowerFlag();
          const receipt = await tx.wait();
          assert.equal(0, receipt.events?.length);
        });
      });
    });

    describe("when called by an enabled setter", () => {
      beforeEach(async () => {
        await controller.connect(personas.Nelly).addAccess(await personas.Neil.getAddress());
      });

      it("lower the flag", async () => {
        await flags.connect(personas.Neil).lowerFlag(), assert.equal(false, await flags.isRaised());
      });
    });

    describe("when called by a non-enabled setter", () => {
      it("reverts", async () => {
        await expect(flags.connect(personas.Neil).lowerFlag()).to.be.revertedWith("Not allowed to lower the flag");
      });
    });

    describe("when called when there is no raisingAccessController", () => {
      beforeEach(async () => {
        await expect(
          flags.connect(personas.Nelly).setRaisingAccessController("0x0000000000000000000000000000000000000000"),
        ).to.emit(flags, "RaisingAccessControllerUpdated");
        assert.equal("0x0000000000000000000000000000000000000000", await flags.raisingAccessController());
      });

      it("succeeds for the owner", async () => {
        await flags.connect(personas.Nelly).lowerFlag();
        assert.equal(false, await flags.isRaised());
      });

      it("reverts for non-owner", async () => {
        await expect(flags.connect(personas.Neil).lowerFlag()).to.be.reverted;
      });
    });
  });

  describe("#isRaised", () => {
    describe("when called by a consumer", () => {
      it("gets the correct status", async () => {
        assert.equal(false, await consumer.isRaised());
        await flags.connect(personas.Nelly).raiseFlag();
        assert.equal(true, await consumer.isRaised());
      });
    });
  });

  describe("#setRaisingAccessController", () => {
    let controller2: Contract;

    beforeEach(async () => {
      controller2 = await controllerFactory.connect(personas.Nelly).deploy();
      await controller2.connect(personas.Nelly).enableAccessCheck();
    });

    it("updates access control rules", async () => {
      const neilAddress = await personas.Neil.getAddress();
      await controller.connect(personas.Nelly).addAccess(neilAddress);
      await flags.connect(personas.Neil).raiseFlag(); // doesn't raise

      await flags.connect(personas.Nelly).setRaisingAccessController(controller2.address);

      await expect(flags.connect(personas.Neil).raiseFlag()).to.be.revertedWith("Not allowed to raise the flag");
    });

    it("emits a log announcing the change", async () => {
      await expect(flags.connect(personas.Nelly).setRaisingAccessController(controller2.address))
        .to.emit(flags, "RaisingAccessControllerUpdated")
        .withArgs(controller.address, controller2.address);
    });

    it("does not emit a log when there is no change", async () => {
      await flags.connect(personas.Nelly).setRaisingAccessController(controller2.address);

      await expect(flags.connect(personas.Nelly).setRaisingAccessController(controller2.address)).to.not.emit(
        flags,
        "RaisingAccessControllerUpdated",
      );
    });

    describe("when called by a non-owner", () => {
      it("reverts", async () => {
        await expect(flags.connect(personas.Neil).setRaisingAccessController(controller2.address)).to.be.revertedWith(
          "Only callable by owner",
        );
      });
    });
  });
});
