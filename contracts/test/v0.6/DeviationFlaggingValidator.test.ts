import { ethers } from "hardhat";
import { publicAbi } from "../test-helpers/helpers";
import { assert, expect } from "chai";
import { BigNumber, Contract, ContractFactory } from "ethers";
import { Personas, getUsers } from "../test-helpers/setup";
import { bigNumEquals } from "../test-helpers/matchers";

let personas: Personas;
let validatorFactory: ContractFactory;
let flagsFactory: ContractFactory;
let acFactory: ContractFactory;

before(async () => {
  personas = (await getUsers()).personas;
  validatorFactory = await ethers.getContractFactory("DeviationFlaggingValidator", personas.Carol);
  flagsFactory = await ethers.getContractFactory("Flags", personas.Carol);
  acFactory = await ethers.getContractFactory("SimpleWriteAccessController", personas.Carol);
});

describe("DeviationFlaggingValidator", () => {
  let validator: Contract;
  let flags: Contract;
  let ac: Contract;
  const flaggingThreshold = 10000; // 10%
  const previousRoundId = 2;
  const previousValue = 1000000;
  const currentRoundId = 3;
  const currentValue = 1000000;

  beforeEach(async () => {
    ac = await acFactory.connect(personas.Carol).deploy();
    flags = await flagsFactory.connect(personas.Carol).deploy(ac.address);
    validator = await validatorFactory.connect(personas.Carol).deploy(flags.address, flaggingThreshold);
    await ac.connect(personas.Carol).addAccess(validator.address);
  });

  it("has a limited public interface", () => {
    publicAbi(validator, [
      "THRESHOLD_MULTIPLIER",
      "flaggingThreshold",
      "flags",
      "isValid",
      "setFlagsAddress",
      "setFlaggingThreshold",
      "validate",
      // Owned methods:
      "acceptOwnership",
      "owner",
      "transferOwnership",
    ]);
  });

  describe("#constructor", () => {
    it("sets the arguments passed in", async () => {
      assert.equal(flags.address, await validator.flags());
      bigNumEquals(flaggingThreshold, await validator.flaggingThreshold());
    });
  });

  describe("#validate", () => {
    describe("when the deviation is greater than the threshold", () => {
      const currentValue = 1100010;

      it("does raises a flag for the calling address", async () => {
        await expect(
          validator.connect(personas.Nelly).validate(previousRoundId, previousValue, currentRoundId, currentValue),
        )
          .to.emit(flags, "FlagRaised")
          .withArgs(await personas.Nelly.getAddress());
      });

      it("uses less than the gas allotted by the aggregator", async () => {
        const tx = await validator
          .connect(personas.Nelly)
          .validate(previousRoundId, previousValue, currentRoundId, currentValue);
        const receipt = await tx.wait();
        assert(receipt);
        if (receipt && receipt.gasUsed) {
          assert.isAbove(receipt.gasUsed.toNumber(), 60000);
        }
      });
    });

    describe("when the deviation is less than or equal to the threshold", () => {
      const currentValue = 1100009;

      it("does raises a flag for the calling address", async () => {
        await expect(
          validator.connect(personas.Nelly).validate(previousRoundId, previousValue, currentRoundId, currentValue),
        ).to.not.emit(flags, "FlagRaised");
      });

      it("uses less than the gas allotted by the aggregator", async () => {
        const tx = await validator
          .connect(personas.Nelly)
          .validate(previousRoundId, previousValue, currentRoundId, currentValue);
        const receipt = await tx.wait();
        assert(receipt);
        if (receipt && receipt.gasUsed) {
          assert.isAbove(receipt.gasUsed.toNumber(), 24000);
        }
      });
    });

    describe("when called with a previous value of zero", () => {
      const previousValue = 0;

      it("does not raise any flags", async () => {
        const tx = await validator
          .connect(personas.Nelly)
          .validate(previousRoundId, previousValue, currentRoundId, currentValue);
        const receipt = await tx.wait();
        assert.equal(0, receipt.events?.length);
      });
    });
  });

  describe("#isValid", () => {
    const previousValue = 1000000;

    describe("with a validation larger than the deviation", () => {
      const currentValue = 1100010;
      it("is not valid", async () => {
        assert.isFalse(await validator.isValid(0, previousValue, 1, currentValue));
      });
    });

    describe("with a validation smaller than the deviation", () => {
      const currentValue = 1100009;
      it("is valid", async () => {
        assert.isTrue(await validator.isValid(0, previousValue, 1, currentValue));
      });
    });

    describe("with positive previous and negative current", () => {
      const previousValue = 1000000;
      const currentValue = -900000;
      it("correctly detects the difference", async () => {
        assert.isFalse(await validator.isValid(0, previousValue, 1, currentValue));
      });
    });

    describe("with negative previous and positive current", () => {
      const previousValue = -900000;
      const currentValue = 1000000;
      it("correctly detects the difference", async () => {
        assert.isFalse(await validator.isValid(0, previousValue, 1, currentValue));
      });
    });

    describe("when the difference overflows", () => {
      const previousValue = BigNumber.from(2).pow(255).sub(1);
      const currentValue = BigNumber.from(-1);

      it("does not revert and returns false", async () => {
        assert.isFalse(await validator.isValid(0, previousValue, 1, currentValue));
      });
    });

    describe("when the rounding overflows", () => {
      const previousValue = BigNumber.from(2).pow(255).div(10000);
      const currentValue = BigNumber.from(1);

      it("does not revert and returns false", async () => {
        assert.isFalse(await validator.isValid(0, previousValue, 1, currentValue));
      });
    });

    describe("when the division overflows", () => {
      const previousValue = BigNumber.from(2).pow(255).sub(1);
      const currentValue = BigNumber.from(-1);

      it("does not revert and returns false", async () => {
        assert.isFalse(await validator.isValid(0, previousValue, 1, currentValue));
      });
    });
  });

  describe("#setFlaggingThreshold", () => {
    const newThreshold = 777;

    it("changes the flagging thresold", async () => {
      assert.equal(flaggingThreshold, await validator.flaggingThreshold());

      await validator.connect(personas.Carol).setFlaggingThreshold(newThreshold);

      assert.equal(newThreshold, await validator.flaggingThreshold());
    });

    it("emits a log event only when actually changed", async () => {
      await expect(validator.connect(personas.Carol).setFlaggingThreshold(newThreshold))
        .to.emit(validator, "FlaggingThresholdUpdated")
        .withArgs(flaggingThreshold, newThreshold);

      await expect(validator.connect(personas.Carol).setFlaggingThreshold(newThreshold)).to.not.emit(
        validator,
        "FlaggingThresholdUpdated",
      );
    });

    describe("when called by a non-owner", () => {
      it("reverts", async () => {
        await expect(validator.connect(personas.Neil).setFlaggingThreshold(newThreshold)).to.be.revertedWith(
          "Only callable by owner",
        );
      });
    });
  });

  describe("#setFlagsAddress", () => {
    const newFlagsAddress = "0x0123456789012345678901234567890123456789";

    it("changes the flags address", async () => {
      assert.equal(flags.address, await validator.flags());

      await validator.connect(personas.Carol).setFlagsAddress(newFlagsAddress);

      assert.equal(newFlagsAddress, await validator.flags());
    });

    it("emits a log event only when actually changed", async () => {
      await expect(validator.connect(personas.Carol).setFlagsAddress(newFlagsAddress))
        .to.emit(validator, "FlagsAddressUpdated")
        .withArgs(flags.address, newFlagsAddress);

      await expect(validator.connect(personas.Carol).setFlagsAddress(newFlagsAddress)).to.not.emit(
        validator,
        "FlagsAddressUpdated",
      );
    });

    describe("when called by a non-owner", () => {
      it("reverts", async () => {
        await expect(validator.connect(personas.Neil).setFlagsAddress(newFlagsAddress)).to.be.revertedWith(
          "Only callable by owner",
        );
      });
    });
  });
});
