const { ethers } = require("hardhat");
const { publicAbi } = require("./helpers");
const { assert, expect } = require("chai");
const { constants } = require("@openzeppelin/test-helpers");

describe("ValidatorProxy", () => {
  // let accounts
  let owner, aggregator, validator;
  let validatorProxy;
  let accounts;

  beforeEach(async () => {
    accounts = await ethers.getSigners();
    owner = accounts[0];
    aggregator = accounts[1];
    validator = accounts[2];
    const vpf = await ethers.getContractFactory("ValidatorProxy", owner.address);
    validatorProxy = await vpf.deploy(aggregator.address, validator.address);
    validatorProxy = await validatorProxy.deployed();
  });

  it("has a limited public interface", async () => {
    publicAbi(validatorProxy, [
      // ConfirmedOwner functions
      "acceptOwnership",
      "owner",
      "transferOwnership",
      // ValidatorProxy functions
      "validate",
      "proposeNewAggregator",
      "upgradeAggregator",
      "getAggregators",
      "proposeNewValidator",
      "upgradeValidator",
      "getValidators",
      "typeAndVersion",
    ]);
  });

  describe("#constructor", () => {
    it("should set the aggregator addresses correctly", async () => {
      const response = await validatorProxy.getAggregators();
      assert.equal(response.current, aggregator.address);
      assert.equal(response.hasProposal, false);
      assert.equal(response.proposed, constants.ZERO_ADDRESS);
    });

    it("should set the validator addresses conrrectly", async () => {
      const response = await validatorProxy.getValidators();
      assert.equal(response.current, validator.address);
      assert.equal(response.hasProposal, false);
      assert.equal(response.proposed, constants.ZERO_ADDRESS);
    });

    it("should set the owner correctly", async () => {
      const response = await validatorProxy.owner();
      assert.equal(response, owner.address);
    });
  });

  describe("#proposeNewAggregator", () => {
    let newAggregator;
    beforeEach(async () => {
      newAggregator = accounts[3].address;
    });

    describe("failure", () => {
      it("should only be called by the owner", async () => {
        const stranger = accounts[4];
        await expect(validatorProxy.connect(stranger).proposeNewAggregator(newAggregator)).to.be.revertedWith(
          "Only callable by owner",
        );
      });

      it("should revert if no change in proposal", async () => {
        await validatorProxy.proposeNewAggregator(newAggregator);
        await expect(validatorProxy.proposeNewAggregator(newAggregator)).to.be.revertedWith("Invalid proposal");
      });

      it("should revert if the proposal is the same as the current", async () => {
        await expect(validatorProxy.proposeNewAggregator(aggregator.address)).to.be.revertedWith("Invalid proposal");
      });
    });

    describe("success", () => {
      it("should emit an event", async () => {
        await expect(validatorProxy.proposeNewAggregator(newAggregator))
          .to.emit(validatorProxy, "AggregatorProposed")
          .withArgs(newAggregator);
      });

      it("should set the correct address and hasProposal is true", async () => {
        await validatorProxy.proposeNewAggregator(newAggregator);
        const response = await validatorProxy.getAggregators();
        assert.equal(response.current, aggregator.address);
        assert.equal(response.hasProposal, true);
        assert.equal(response.proposed, newAggregator);
      });

      it("should set a zero address and hasProposal is false", async () => {
        await validatorProxy.proposeNewAggregator(newAggregator);
        await validatorProxy.proposeNewAggregator(constants.ZERO_ADDRESS);
        const response = await validatorProxy.getAggregators();
        assert.equal(response.current, aggregator.address);
        assert.equal(response.hasProposal, false);
        assert.equal(response.proposed, constants.ZERO_ADDRESS);
      });
    });
  });

  describe("#upgradeAggregator", () => {
    describe("failure", () => {
      it("should only be called by the owner", async () => {
        const stranger = accounts[4];
        await expect(validatorProxy.connect(stranger).upgradeAggregator()).to.be.revertedWith("Only callable by owner");
      });

      it("should revert if there is no proposal", async () => {
        await expect(validatorProxy.upgradeAggregator()).to.be.revertedWith("No proposal");
      });
    });

    describe("success", () => {
      let newAggregator;
      beforeEach(async () => {
        newAggregator = accounts[3].address;
        await validatorProxy.proposeNewAggregator(newAggregator);
      });

      it("should emit an event", async () => {
        await expect(validatorProxy.upgradeAggregator())
          .to.emit(validatorProxy, "AggregatorUpgraded")
          .withArgs(aggregator.address, newAggregator);
      });

      it("should upgrade the addresses", async () => {
        await validatorProxy.upgradeAggregator();
        const response = await validatorProxy.getAggregators();
        assert.equal(response.current, newAggregator);
        assert.equal(response.hasProposal, false);
        assert.equal(response.proposed, constants.ZERO_ADDRESS);
      });
    });
  });

  describe("#proposeNewValidator", () => {
    let newValidator;

    beforeEach(() => {
      newValidator = accounts[3].address;
    });

    describe("failure", () => {
      it("should only be called by the owner", async () => {
        const stranger = accounts[4];
        await expect(validatorProxy.connect(stranger).proposeNewAggregator(newValidator)).to.be.revertedWith(
          "Only callable by owner",
        );
      });

      it("should revert if no change in proposal", async () => {
        await validatorProxy.proposeNewValidator(newValidator);
        await expect(validatorProxy.proposeNewValidator(newValidator)).to.be.revertedWith("Invalid proposal");
      });

      it("should revert if the proposal is the same as the current", async () => {
        await expect(validatorProxy.proposeNewValidator(validator.address)).to.be.revertedWith("Invalid proposal");
      });
    });

    describe("success", () => {
      it("should emit an event", async () => {
        await expect(validatorProxy.proposeNewValidator(newValidator))
          .to.emit(validatorProxy, "ValidatorProposed")
          .withArgs(newValidator);
      });

      it("should set the correct address and hasProposal is true", async () => {
        await validatorProxy.proposeNewValidator(newValidator);
        const response = await validatorProxy.getValidators();
        assert.equal(response.current, validator.address);
        assert.equal(response.hasProposal, true);
        assert.equal(response.proposed, newValidator);
      });

      it("should set a zero address and hasProposal is false", async () => {
        await validatorProxy.proposeNewValidator(newValidator);
        await validatorProxy.proposeNewValidator(constants.ZERO_ADDRESS);
        const response = await validatorProxy.getValidators();
        assert.equal(response.current, validator.address);
        assert.equal(response.hasProposal, false);
        assert.equal(response.proposed, constants.ZERO_ADDRESS);
      });
    });
  });

  describe("#upgradeValidator", () => {
    describe("failure", () => {
      it("should only be called by the owner", async () => {
        const stranger = accounts[4];
        await expect(validatorProxy.connect(stranger).upgradeValidator()).to.be.revertedWith("Only callable by owner");
      });

      it("should revert if there is no proposal", async () => {
        await expect(validatorProxy.upgradeValidator()).to.be.revertedWith("No proposal");
      });
    });

    describe("success", () => {
      let newValidator;
      beforeEach(async () => {
        newValidator = accounts[3].address;
        await validatorProxy.proposeNewValidator(newValidator);
      });

      it("should emit an event", async () => {
        await expect(validatorProxy.upgradeValidator())
          .to.emit(validatorProxy, "ValidatorUpgraded")
          .withArgs(validator.address, newValidator);
      });

      it("should upgrade the addresses", async () => {
        await validatorProxy.upgradeValidator();
        const response = await validatorProxy.getValidators();
        assert.equal(response.current, newValidator);
        assert.equal(response.hasProposal, false);
        assert.equal(response.proposed, constants.ZERO_ADDRESS);
      });
    });
  });

  describe("#validate", () => {
    describe("failure", () => {
      it("reverts when not called by aggregator or proposed aggregator", async () => {
        const stranger = accounts[9];
        await expect(validatorProxy.connect(stranger).validate(99, 88, 77, 66)).to.be.revertedWith(
          "Not a configured aggregator",
        );
      });

      it("reverts when there is no validator set", async () => {
        const vpf = await ethers.getContractFactory("ValidatorProxy", owner.address);
        validatorProxy = await vpf.deploy(aggregator.address, constants.ZERO_ADDRESS);
        await validatorProxy.deployed();
        await expect(validatorProxy.connect(aggregator).validate(99, 88, 77, 66)).to.be.revertedWith(
          "No validator set",
        );
      });
    });

    describe("success", () => {
      describe("from the aggregator", () => {
        let mockValidator1;
        beforeEach(async () => {
          const mvf = await ethers.getContractFactory("MockAggregatorValidator", owner.address);
          mockValidator1 = await mvf.deploy(1);
          mockValidator1 = await mockValidator1.deployed();
          const vpf = await ethers.getContractFactory("ValidatorProxy", owner.address);
          validatorProxy = await vpf.deploy(aggregator.address, mockValidator1.address);
          validatorProxy = await validatorProxy.deployed();
        });

        describe("for a single validator", () => {
          it("calls validate on the validator", async () => {
            await expect(validatorProxy.connect(aggregator).validate(200, 300, 400, 500))
              .to.emit(mockValidator1, "ValidateCalled")
              .withArgs(1, 200, 300, 400, 500);
          });

          it("uses a specific amount of gas", async () => {
            const resp = await validatorProxy.connect(aggregator).validate(200, 300, 400, 500);
            const receipt = await resp.wait();
            assert.equal(receipt.gasUsed.toString(), "35371");
          });
        });

        describe("for a validator and a proposed validator", () => {
          let mockValidator2;

          beforeEach(async () => {
            const mvf = await ethers.getContractFactory("MockAggregatorValidator", owner.address);
            mockValidator2 = await mvf.deploy(2);
            mockValidator2 = await mockValidator2.deployed();
            await validatorProxy.proposeNewValidator(mockValidator2.address);
          });

          it("calls validate on the validator", async () => {
            await expect(validatorProxy.connect(aggregator).validate(2000, 3000, 4000, 5000))
              .to.emit(mockValidator1, "ValidateCalled")
              .withArgs(1, 2000, 3000, 4000, 5000);
          });

          it("also calls validate on the proposed validator", async () => {
            await expect(validatorProxy.connect(aggregator).validate(2000, 3000, 4000, 5000))
              .to.emit(mockValidator2, "ValidateCalled")
              .withArgs(2, 2000, 3000, 4000, 5000);
          });

          it("uses a specific amount of gas", async () => {
            const resp = await validatorProxy.connect(aggregator).validate(2000, 3000, 4000, 5000);
            const receipt = await resp.wait();
            assert.equal(receipt.gasUsed.toString(), "45318");
          });
        });
      });

      describe("from the proposed aggregator", () => {
        let newAggregator;
        beforeEach(async () => {
          newAggregator = accounts[3];
          await validatorProxy.connect(owner).proposeNewAggregator(newAggregator.address);
        });

        it("emits an event", async () => {
          await expect(validatorProxy.connect(newAggregator).validate(555, 666, 777, 888))
            .to.emit(validatorProxy, "ProposedAggregatorValidateCall")
            .withArgs(newAggregator.address, 555, 666, 777, 888);
        });
      });
    });
  });
});
