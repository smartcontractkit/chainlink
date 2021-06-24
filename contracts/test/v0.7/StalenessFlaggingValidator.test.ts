import { ethers } from "hardhat";
import { evmWordToAddress, getLog, getLogs, numToBytes32, publicAbi } from "../test-helpers/helpers";
import { assert, expect } from "chai";
import { BigNumber, Contract, ContractFactory } from "ethers";
import { Personas, getUsers } from "../test-helpers/setup";
import { evmRevert } from "../test-helpers/matchers";

let personas: Personas;
let validatorFactory: ContractFactory;
let flagsFactory: ContractFactory;
let acFactory: ContractFactory;
let aggregatorFactory: ContractFactory;

before(async () => {
  personas = (await getUsers()).personas;

  validatorFactory = await ethers.getContractFactory("StalenessFlaggingValidator", personas.Carol);
  flagsFactory = await ethers.getContractFactory("Flags", personas.Carol);
  acFactory = await ethers.getContractFactory("SimpleWriteAccessController", personas.Carol);
  aggregatorFactory = await ethers.getContractFactory("MockV3Aggregator", personas.Carol);
});

describe("StalenessFlaggingValidator", () => {
  let validator: Contract;
  let flags: Contract;
  let ac: Contract;

  const flaggingThreshold1 = 10000;
  const flaggingThreshold2 = 20000;

  beforeEach(async () => {
    ac = await acFactory.connect(personas.Carol).deploy();
    flags = await flagsFactory.connect(personas.Carol).deploy(ac.address);
    validator = await validatorFactory.connect(personas.Carol).deploy(flags.address);

    await ac.connect(personas.Carol).addAccess(validator.address);
  });

  it("has a limited public interface", () => {
    publicAbi(validator, [
      "update",
      "check",
      "setThresholds",
      "setFlagsAddress",
      "threshold",
      "flags",
      // Upkeep methods:
      "checkUpkeep",
      "performUpkeep",
      // Owned methods:
      "acceptOwnership",
      "owner",
      "transferOwnership",
    ]);
  });

  describe("#constructor", () => {
    it("sets the arguments passed in", async () => {
      assert.equal(await validator.flags(), flags.address);
    });

    it("sets the owner", async () => {
      assert.equal(await validator.owner(), await personas.Carol.getAddress());
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
      const tx = await validator.connect(personas.Carol).setFlagsAddress(newFlagsAddress);
      await expect(tx).to.emit(validator, "FlagsAddressUpdated").withArgs(flags.address, newFlagsAddress);

      const sameChangeTx = await validator.connect(personas.Carol).setFlagsAddress(newFlagsAddress);

      await expect(sameChangeTx).to.not.emit(validator, "FlagsAddressUpdated");
    });

    describe("when called by a non-owner", () => {
      it("reverts", async () => {
        await evmRevert(validator.connect(personas.Neil).setFlagsAddress(newFlagsAddress), "Only callable by owner");
      });
    });
  });

  describe("#setThresholds", () => {
    let agg1: Contract;
    let agg2: Contract;
    let aggregators: Array<string>;
    let thresholds: Array<number>;

    beforeEach(async () => {
      const decimals = 8;
      const initialAnswer = 10000000000;
      agg1 = await aggregatorFactory.connect(personas.Carol).deploy(decimals, initialAnswer);
      agg2 = await aggregatorFactory.connect(personas.Carol).deploy(decimals, initialAnswer);
    });

    describe("failure", () => {
      beforeEach(() => {
        aggregators = [agg1.address, agg2.address];
        thresholds = [flaggingThreshold1];
      });

      it("reverts when called by a non-owner", async () => {
        await evmRevert(
          validator.connect(personas.Neil).setThresholds(aggregators, thresholds),
          "Only callable by owner",
        );
      });

      it("reverts when passed uneven arrays", async () => {
        await evmRevert(
          validator.connect(personas.Carol).setThresholds(aggregators, thresholds),
          "Different sized arrays",
        );
      });
    });

    describe("success", () => {
      let tx: any;

      beforeEach(() => {
        aggregators = [agg1.address, agg2.address];
        thresholds = [flaggingThreshold1, flaggingThreshold2];
      });

      describe("when called with 2 new thresholds", () => {
        beforeEach(async () => {
          tx = await validator.connect(personas.Carol).setThresholds(aggregators, thresholds);
        });

        it("sets the thresholds", async () => {
          const first = await validator.threshold(agg1.address);
          const second = await validator.threshold(agg2.address);
          assert.equal(first.toString(), flaggingThreshold1.toString());
          assert.equal(second.toString(), flaggingThreshold2.toString());
        });

        it("emits events", async () => {
          const firstEvent = await getLog(tx, 0);
          assert.equal(evmWordToAddress(firstEvent.topics[1]), agg1.address);
          assert.equal(firstEvent.topics[3], numToBytes32(flaggingThreshold1));
          const secondEvent = await getLog(tx, 1);
          assert.equal(evmWordToAddress(secondEvent.topics[1]), agg2.address);
          assert.equal(secondEvent.topics[3], numToBytes32(flaggingThreshold2));
        });
      });

      describe("when called with 2, but 1 has not changed", () => {
        it("emits only 1 event", async () => {
          tx = await validator.connect(personas.Carol).setThresholds(aggregators, thresholds);

          const newThreshold = flaggingThreshold2 + 1;
          tx = await validator.connect(personas.Carol).setThresholds(aggregators, [flaggingThreshold1, newThreshold]);
          const logs = await getLogs(tx);
          assert.equal(logs.length, 1);
          const log = logs[0];
          assert.equal(evmWordToAddress(log.topics[1]), agg2.address);
          assert.equal(log.topics[2], numToBytes32(flaggingThreshold2));
          assert.equal(log.topics[3], numToBytes32(newThreshold));
        });
      });
    });
  });

  describe("#check", () => {
    let agg1: Contract;
    let agg2: Contract;
    let aggregators: Array<string>;
    let thresholds: Array<number>;
    const decimals = 8;
    const initialAnswer = 10000000000;
    beforeEach(async () => {
      agg1 = await aggregatorFactory.connect(personas.Carol).deploy(decimals, initialAnswer);
      agg2 = await aggregatorFactory.connect(personas.Carol).deploy(decimals, initialAnswer);
      aggregators = [agg1.address, agg2.address];
      thresholds = [flaggingThreshold1, flaggingThreshold2];
      await validator.setThresholds(aggregators, thresholds);
    });

    describe("when neither are stale", () => {
      it("returns an empty array", async () => {
        const response = await validator.check(aggregators);
        assert.equal(response.length, 0);
      });
    });

    describe("when threshold is not set in the validator", () => {
      it("returns an empty array", async () => {
        const agg3 = await aggregatorFactory.connect(personas.Carol).deploy(decimals, initialAnswer);
        const response = await validator.check([agg3.address]);
        assert.equal(response.length, 0);
      });
    });

    describe("when one of the aggregators is stale", () => {
      it("returns an array with one stale aggregator", async () => {
        const currentTimestamp = await agg1.latestTimestamp();
        const staleTimestamp = currentTimestamp.sub(BigNumber.from(flaggingThreshold1 + 1));
        await agg1.updateRoundData(99, initialAnswer, staleTimestamp, staleTimestamp);
        const response = await validator.check(aggregators);

        assert.equal(response.length, 1);
        assert.equal(response[0], agg1.address);
      });
    });

    describe("When both aggregators are stale", () => {
      it("returns an array with both aggregators", async () => {
        let currentTimestamp = await agg1.latestTimestamp();
        let staleTimestamp = currentTimestamp.sub(BigNumber.from(flaggingThreshold1 + 1));
        await agg1.updateRoundData(99, initialAnswer, staleTimestamp, staleTimestamp);

        currentTimestamp = await agg2.latestTimestamp();
        staleTimestamp = currentTimestamp.sub(BigNumber.from(flaggingThreshold2 + 1));
        await agg2.updateRoundData(99, initialAnswer, staleTimestamp, staleTimestamp);

        const response = await validator.check(aggregators);

        assert.equal(response.length, 2);
        assert.equal(response[0], agg1.address);
        assert.equal(response[1], agg2.address);
      });
    });
  });

  describe("#update", () => {
    let agg1: Contract;
    let agg2: Contract;
    let aggregators: Array<string>;
    let thresholds: Array<number>;
    const decimals = 8;
    const initialAnswer = 10000000000;
    beforeEach(async () => {
      agg1 = await aggregatorFactory.connect(personas.Carol).deploy(decimals, initialAnswer);
      agg2 = await aggregatorFactory.connect(personas.Carol).deploy(decimals, initialAnswer);
      aggregators = [agg1.address, agg2.address];
      thresholds = [flaggingThreshold1, flaggingThreshold2];
      await validator.setThresholds(aggregators, thresholds);
    });

    describe("when neither are stale", () => {
      it("does not raise a flag", async () => {
        const tx = await validator.update(aggregators);
        const logs = await getLogs(tx);
        assert.equal(logs.length, 0);
      });
    });

    describe("when threshold is not set in the validator", () => {
      it("does not raise a flag", async () => {
        const agg3 = await aggregatorFactory.connect(personas.Carol).deploy(decimals, initialAnswer);
        const tx = await validator.update([agg3.address]);
        const logs = await getLogs(tx);
        assert.equal(logs.length, 0);
      });
    });

    describe("when one is stale", () => {
      it("raises a flag for that aggregator", async () => {
        const currentTimestamp = await agg1.latestTimestamp();
        const staleTimestamp = currentTimestamp.sub(BigNumber.from(flaggingThreshold1 + 1));
        await agg1.updateRoundData(99, initialAnswer, staleTimestamp, staleTimestamp);

        const tx = await validator.update(aggregators);
        const logs = await getLogs(tx);
        assert.equal(logs.length, 1);
        assert.equal(evmWordToAddress(logs[0].topics[1]), agg1.address);
      });
    });

    describe("when both are stale", () => {
      it("raises 2 flags, one for each aggregator", async () => {
        let currentTimestamp = await agg1.latestTimestamp();
        let staleTimestamp = currentTimestamp.sub(BigNumber.from(flaggingThreshold1 + 1));
        await agg1.updateRoundData(99, initialAnswer, staleTimestamp, staleTimestamp);

        currentTimestamp = await agg2.latestTimestamp();
        staleTimestamp = currentTimestamp.sub(BigNumber.from(flaggingThreshold2 + 1));
        await agg2.updateRoundData(99, initialAnswer, staleTimestamp, staleTimestamp);

        const tx = await validator.update(aggregators);
        const logs = await getLogs(tx);
        assert.equal(logs.length, 2);
        assert.equal(evmWordToAddress(logs[0].topics[1]), agg1.address);
        assert.equal(evmWordToAddress(logs[1].topics[1]), agg2.address);
      });
    });
  });

  describe("#checkUpkeep", () => {
    let agg1: Contract;
    let agg2: Contract;
    let aggregators: Array<string>;
    let thresholds: Array<number>;
    const decimals = 8;
    const initialAnswer = 10000000000;
    beforeEach(async () => {
      agg1 = await aggregatorFactory.connect(personas.Carol).deploy(decimals, initialAnswer);
      agg2 = await aggregatorFactory.connect(personas.Carol).deploy(decimals, initialAnswer);
      aggregators = [agg1.address, agg2.address];
      thresholds = [flaggingThreshold1, flaggingThreshold2];
      await validator.setThresholds(aggregators, thresholds);
    });

    describe("when neither are stale", () => {
      it("returns false and an empty array", async () => {
        const bytesData = ethers.utils.defaultAbiCoder.encode(["address[]"], [aggregators]);
        const response = await validator.checkUpkeep(bytesData);

        assert.equal(response[0], false);
        const decodedResponse = ethers.utils.defaultAbiCoder.decode(["address[]"], response?.[1]);
        assert.equal(decodedResponse[0].length, 0);
      });
    });

    describe("when threshold is not set in the validator", () => {
      it("returns flase and an empty array", async () => {
        const agg3 = await aggregatorFactory.connect(personas.Carol).deploy(decimals, initialAnswer);
        const bytesData = ethers.utils.defaultAbiCoder.encode(["address[]"], [[agg3.address]]);
        const response = await validator.checkUpkeep(bytesData);

        assert.equal(response[0], false);
        const decodedResponse = ethers.utils.defaultAbiCoder.decode(["address[]"], response?.[1]);
        assert.equal(decodedResponse[0].length, 0);
      });
    });

    describe("when one of the aggregators is stale", () => {
      it("returns true with an array with one stale aggregator", async () => {
        const currentTimestamp = await agg1.latestTimestamp();
        const staleTimestamp = currentTimestamp.sub(BigNumber.from(flaggingThreshold1 + 1));
        await agg1.updateRoundData(99, initialAnswer, staleTimestamp, staleTimestamp);

        const bytesData = ethers.utils.defaultAbiCoder.encode(["address[]"], [aggregators]);
        const response = await validator.checkUpkeep(bytesData);

        assert.equal(response[0], true);
        const decodedResponse = ethers.utils.defaultAbiCoder.decode(["address[]"], response?.[1]);
        const decodedArray = decodedResponse[0];
        assert.equal(decodedArray.length, 1);
        assert.equal(decodedArray[0], agg1.address);
      });
    });

    describe("When both aggregators are stale", () => {
      it("returns true with an array with both aggregators", async () => {
        let currentTimestamp = await agg1.latestTimestamp();
        let staleTimestamp = currentTimestamp.sub(BigNumber.from(flaggingThreshold1 + 1));
        await agg1.updateRoundData(99, initialAnswer, staleTimestamp, staleTimestamp);

        currentTimestamp = await agg2.latestTimestamp();
        staleTimestamp = currentTimestamp.sub(BigNumber.from(flaggingThreshold2 + 1));
        await agg2.updateRoundData(99, initialAnswer, staleTimestamp, staleTimestamp);

        const bytesData = ethers.utils.defaultAbiCoder.encode(["address[]"], [aggregators]);
        const response = await validator.checkUpkeep(bytesData);

        assert.equal(response[0], true);
        const decodedResponse = ethers.utils.defaultAbiCoder.decode(["address[]"], response?.[1]);
        const decodedArray = decodedResponse[0];
        assert.equal(decodedArray.length, 2);
        assert.equal(decodedArray[0], agg1.address);
        assert.equal(decodedArray[1], agg2.address);
      });
    });
  });

  describe("#performUpkeep", () => {
    let agg1: Contract;
    let agg2: Contract;
    let aggregators: Array<string>;
    let thresholds: Array<number>;
    const decimals = 8;
    const initialAnswer = 10000000000;
    beforeEach(async () => {
      agg1 = await aggregatorFactory.connect(personas.Carol).deploy(decimals, initialAnswer);
      agg2 = await aggregatorFactory.connect(personas.Carol).deploy(decimals, initialAnswer);
      aggregators = [agg1.address, agg2.address];
      thresholds = [flaggingThreshold1, flaggingThreshold2];
      await validator.setThresholds(aggregators, thresholds);
    });

    describe("when neither are stale", () => {
      it("does not raise a flag", async () => {
        const bytesData = ethers.utils.defaultAbiCoder.encode(["address[]"], [aggregators]);
        const tx = await validator.performUpkeep(bytesData);
        const logs = await getLogs(tx);
        assert.equal(logs.length, 0);
      });
    });

    describe("when threshold is not set in the validator", () => {
      it("does not raise a flag", async () => {
        const agg3 = await aggregatorFactory.connect(personas.Carol).deploy(decimals, initialAnswer);
        const bytesData = ethers.utils.defaultAbiCoder.encode(["address[]"], [[agg3.address]]);
        const tx = await validator.performUpkeep(bytesData);
        const logs = await getLogs(tx);
        assert.equal(logs.length, 0);
      });
    });

    describe("when one is stale", () => {
      it("raises a flag for that aggregator", async () => {
        const currentTimestamp = await agg1.latestTimestamp();
        const staleTimestamp = currentTimestamp.sub(BigNumber.from(flaggingThreshold1 + 1));
        await agg1.updateRoundData(99, initialAnswer, staleTimestamp, staleTimestamp);

        const bytesData = ethers.utils.defaultAbiCoder.encode(["address[]"], [aggregators]);
        const tx = await validator.performUpkeep(bytesData);
        const logs = await getLogs(tx);
        assert.equal(logs.length, 1);
        assert.equal(evmWordToAddress(logs[0].topics[1]), agg1.address);
      });
    });

    describe("when both are stale", () => {
      it("raises 2 flags, one for each aggregator", async () => {
        let currentTimestamp = await agg1.latestTimestamp();
        let staleTimestamp = currentTimestamp.sub(BigNumber.from(flaggingThreshold1 + 1));
        await agg1.updateRoundData(99, initialAnswer, staleTimestamp, staleTimestamp);

        currentTimestamp = await agg2.latestTimestamp();
        staleTimestamp = currentTimestamp.sub(BigNumber.from(flaggingThreshold2 + 1));
        await agg2.updateRoundData(99, initialAnswer, staleTimestamp, staleTimestamp);

        const bytesData = ethers.utils.defaultAbiCoder.encode(["address[]"], [aggregators]);
        const tx = await validator.performUpkeep(bytesData);
        const logs = await getLogs(tx);
        assert.equal(logs.length, 2);
        assert.equal(evmWordToAddress(logs[0].topics[1]), agg1.address);
        assert.equal(evmWordToAddress(logs[1].topics[1]), agg2.address);
      });
    });
  });
});
