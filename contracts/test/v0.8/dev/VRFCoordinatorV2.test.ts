import { ethers } from "hardhat";
import { Signer, Contract, BigNumber } from "ethers";
import { assert, expect } from "chai";
import { defaultAbiCoder } from "ethers/utils";
import { publicAbi } from "../../test-helpers/helpers";
import { randomAddressString } from "hardhat/internal/hardhat-network/provider/fork/random";

describe("VRFCoordinatorV2", () => {
  let vrfCoordinatorV2: Contract;
  let linkToken: Contract;
  let blockHashStore: Contract;
  let mockLinkEth: Contract;
  let owner: Signer;
  let subOwner: Signer;
  let subOwnerAddress: string;
  let consumer: Signer;
  let random: Signer;
  let randomAddress: string;
  let oracle: Signer;
  const linkEth = BigNumber.from(300000000);
  type config = {
    minimumRequestBlockConfirmations: number;
    fulfillmentFlatFeePPM: number;
    maxGasLimit: number;
    stalenessSeconds: number;
    gasAfterPaymentCalculation: number;
    weiPerUnitLink: BigNumber;
    minimumSubscriptionBalance: BigNumber;
  };
  let c: config;

  beforeEach(async () => {
    let accounts = await ethers.getSigners();
    owner = accounts[0];
    subOwner = accounts[1];
    subOwnerAddress = await subOwner.getAddress();
    consumer = accounts[2];
    random = accounts[3];
    randomAddress = await random.getAddress();
    oracle = accounts[4];
    let ltFactory = await ethers.getContractFactory("LinkToken", accounts[0]);
    linkToken = await ltFactory.deploy();
    let bhFactory = await ethers.getContractFactory("BlockhashStore", accounts[0]);
    blockHashStore = await bhFactory.deploy();
    let mockAggregatorV3Factory = await ethers.getContractFactory("MockV3Aggregator", accounts[0]);
    mockLinkEth = await mockAggregatorV3Factory.deploy(0, linkEth);
    let vrfCoordinatorV2Factory = await ethers.getContractFactory("VRFCoordinatorV2", accounts[0]);
    vrfCoordinatorV2 = await vrfCoordinatorV2Factory.deploy(
      linkToken.address,
      blockHashStore.address,
      mockLinkEth.address,
    );
    await linkToken.transfer(subOwnerAddress, BigNumber.from("1000000000000000000")); // 1 link
    await linkToken.transfer(randomAddress, BigNumber.from("1000000000000000000")); // 1 link
    c = {
      minimumRequestBlockConfirmations: 1,
      fulfillmentFlatFeePPM: 0,
      maxGasLimit: 1000000,
      stalenessSeconds: 86400,
      gasAfterPaymentCalculation: 21000 + 5000 + 2100 + 20000 + 2 * 2100 - 15000 + 7315,
      weiPerUnitLink: BigNumber.from("10000000000000000"),
      minimumSubscriptionBalance: BigNumber.from("100000000000000000"), // 0.1 link
    };
    await vrfCoordinatorV2
      .connect(owner)
      .setConfig(
        c.minimumRequestBlockConfirmations,
        c.fulfillmentFlatFeePPM,
        c.maxGasLimit,
        c.stalenessSeconds,
        c.gasAfterPaymentCalculation,
        c.weiPerUnitLink,
        c.minimumSubscriptionBalance,
      );
  });

  it("has a limited public interface", async () => {
    publicAbi(vrfCoordinatorV2, [
      // Owner
      "acceptOwnership",
      "transferOwnership",
      "owner",
      "getConfig",
      "setConfig",
      "recoverFunds",
      "s_totalBalance",
      // Oracle
      "requestRandomWords",
      "getCommitment", // Note we use this to check if a request is already fulfilled.
      "hashOfKey",
      "fulfillRandomWords",
      "registerProvingKey",
      "oracleWithdraw",
      // Subscription management
      "createSubscription",
      "addConsumer",
      "removeConsumer",
      "getSubscription",
      "onTokenTransfer", // Effectively the fundSubscription.
      "defundSubscription",
      "cancelSubscription",
      "requestSubscriptionOwnerTransfer",
      "acceptSubscriptionOwnerTransfer",
      // Misc
      "typeAndVersion",
      "BLOCKHASH_STORE",
      "LINK",
      "LINK_ETH_FEED",
      "PROOF_LENGTH", // Inherited from VRF.sol as public.
    ]);
  });

  it("configuration management", async function () {
    await expect(
      vrfCoordinatorV2
        .connect(subOwner)
        .setConfig(
          c.minimumRequestBlockConfirmations,
          c.fulfillmentFlatFeePPM,
          c.maxGasLimit,
          c.stalenessSeconds,
          c.gasAfterPaymentCalculation,
          c.weiPerUnitLink,
          c.minimumSubscriptionBalance,
        ),
    ).to.be.revertedWith("Only callable by owner");
    // Anyone can read the config.
    const resp = await vrfCoordinatorV2.connect(random).getConfig();
    assert(resp[0] == c.minimumRequestBlockConfirmations);
    assert(resp[1] == c.fulfillmentFlatFeePPM);
    assert(resp[2] == c.maxGasLimit);
    assert(resp[3] == c.stalenessSeconds);
    assert(resp[4].toString() == c.gasAfterPaymentCalculation.toString());
    assert(resp[5].toString() == c.weiPerUnitLink.toString());
  });

  async function createSubscription(): Promise<number> {
    let consumers: string[] = [await consumer.getAddress()];
    const tx = await vrfCoordinatorV2.connect(subOwner).createSubscription(consumers);
    const receipt = await tx.wait();
    return receipt.events[0].args["subId"];
  }

  describe("#createSubscription", async function () {
    it("can create a subscription", async function () {
      let consumers: string[] = [await consumer.getAddress()];
      await expect(vrfCoordinatorV2.connect(subOwner).createSubscription(consumers))
        .to.emit(vrfCoordinatorV2, "SubscriptionCreated")
        .withArgs(1, subOwnerAddress, consumers);
      const s = await vrfCoordinatorV2.getSubscription(1);
      assert(s.balance.toString() == "0", "invalid balance");
      assert(s.owner == subOwnerAddress, "invalid address");
    });
    it("subscription id increments", async function () {
      let consumers: string[] = [await consumer.getAddress()];
      await expect(vrfCoordinatorV2.connect(subOwner).createSubscription(consumers))
        .to.emit(vrfCoordinatorV2, "SubscriptionCreated")
        .withArgs(1, subOwnerAddress, consumers);
      await expect(vrfCoordinatorV2.connect(subOwner).createSubscription(consumers))
        .to.emit(vrfCoordinatorV2, "SubscriptionCreated")
        .withArgs(2, subOwnerAddress, consumers);
    });
  });

  describe("#requestSubscriptionOwnerTransfer", async function () {
    let subId: number;
    beforeEach(async () => {
      subId = await createSubscription();
    });
    it("rejects non-owner", async function () {
      await expect(
        vrfCoordinatorV2.connect(random).requestSubscriptionOwnerTransfer(subId, randomAddress),
      ).to.be.revertedWith(`MustBeSubOwner("${subOwnerAddress}")`);
    });
    it("owner can request transfer", async function () {
      await expect(vrfCoordinatorV2.connect(subOwner).requestSubscriptionOwnerTransfer(subId, randomAddress))
        .to.emit(vrfCoordinatorV2, "SubscriptionOwnerTransferRequested")
        .withArgs(subId, subOwnerAddress, randomAddress);
      // Same request is a noop
      await expect(
        vrfCoordinatorV2.connect(subOwner).requestSubscriptionOwnerTransfer(subId, randomAddress),
      ).to.not.emit(vrfCoordinatorV2, "SubscriptionOwnerTransferRequested");
    });
  });

  describe("#acceptSubscriptionOwnerTransfer", async function () {
    let subId: number;
    beforeEach(async () => {
      subId = await createSubscription();
    });
    it("subscription must exist", async function () {
      await expect(vrfCoordinatorV2.connect(subOwner).acceptSubscriptionOwnerTransfer(1203123123)).to.be.revertedWith(
        `InvalidSubscription`,
      );
    });
    it("must be requested owner to accept", async function () {
      await expect(vrfCoordinatorV2.connect(subOwner).requestSubscriptionOwnerTransfer(subId, randomAddress));
      await expect(vrfCoordinatorV2.connect(subOwner).acceptSubscriptionOwnerTransfer(subId)).to.be.revertedWith(
        `MustBeRequestedOwner("${randomAddress}")`,
      );
    });
    it("requested owner can accept", async function () {
      await expect(vrfCoordinatorV2.connect(subOwner).requestSubscriptionOwnerTransfer(subId, randomAddress))
        .to.emit(vrfCoordinatorV2, "SubscriptionOwnerTransferRequested")
        .withArgs(subId, subOwnerAddress, randomAddress);
      await expect(vrfCoordinatorV2.connect(random).acceptSubscriptionOwnerTransfer(subId))
        .to.emit(vrfCoordinatorV2, "SubscriptionOwnerTransferred")
        .withArgs(subId, subOwnerAddress, randomAddress);
    });
  });

  describe("#addConsumer", async function () {
    let subId: number;
    beforeEach(async () => {
      subId = await createSubscription();
    });
    it("subscription must exist", async function () {
      await expect(vrfCoordinatorV2.connect(subOwner).addConsumer(1203123123, randomAddress)).to.be.revertedWith(
        `InvalidSubscription`,
      );
    });
    it("must be owner", async function () {
      await expect(vrfCoordinatorV2.connect(random).addConsumer(subId, randomAddress)).to.be.revertedWith(
        `MustBeSubOwner("${subOwnerAddress}")`,
      );
    });
    it("cannot overwrite", async function () {
      await vrfCoordinatorV2.connect(subOwner).addConsumer(subId, randomAddress);
      await expect(vrfCoordinatorV2.connect(subOwner).addConsumer(subId, randomAddress)).to.be.revertedWith(
        `AlreadySubscribed(${subId}, "${randomAddress}")`,
      );
    });
    it("cannot add more than maximum", async function () {
      // There is one consumer, add another 99 to hit the max
      for (let i = 0; i < 99; i++) {
        await vrfCoordinatorV2.connect(subOwner).addConsumer(subId, randomAddressString());
      }
      // Adding one more should fail
      // await vrfCoordinatorV2.connect(subOwner).addConsumer(subId, randomAddress);
      await expect(vrfCoordinatorV2.connect(subOwner).addConsumer(subId, randomAddress)).to.be.revertedWith(
        `TooManyConsumers()`,
      );
    });
    it("owner can update", async function () {
      await expect(vrfCoordinatorV2.connect(subOwner).addConsumer(subId, randomAddress))
        .to.emit(vrfCoordinatorV2, "SubscriptionConsumerAdded")
        .withArgs(subId, randomAddress);
    });
  });

  describe("#removeConsumer", async function () {
    let subId: number;
    beforeEach(async () => {
      subId = await createSubscription();
    });
    it("subscription must exist", async function () {
      await expect(vrfCoordinatorV2.connect(subOwner).removeConsumer(1203123123, randomAddress)).to.be.revertedWith(
        `InvalidSubscription`,
      );
    });
    it("must be owner", async function () {
      await expect(vrfCoordinatorV2.connect(random).removeConsumer(subId, randomAddress)).to.be.revertedWith(
        `MustBeSubOwner("${subOwnerAddress}")`,
      );
    });
    it("owner can update", async function () {
      const subBefore = await vrfCoordinatorV2.getSubscription(subId);
      await vrfCoordinatorV2.connect(subOwner).addConsumer(subId, randomAddress);
      await expect(vrfCoordinatorV2.connect(subOwner).removeConsumer(subId, randomAddress))
        .to.emit(vrfCoordinatorV2, "SubscriptionConsumerRemoved")
        .withArgs(subId, randomAddress);
      const subAfter = await vrfCoordinatorV2.getSubscription(subId);
      // Subscription should NOT contain the removed consumer
      assert.deepEqual(subBefore.consumers, subAfter.consumers);
    });
  });

  describe("#defundSubscription", async function () {
    let subId: number;
    beforeEach(async () => {
      subId = await createSubscription();
    });
    it("subscription must exist", async function () {
      await expect(
        vrfCoordinatorV2.connect(subOwner).defundSubscription(1203123123, subOwnerAddress, BigNumber.from("1000")),
      ).to.be.revertedWith(`InvalidSubscription`);
    });
    it("must be owner", async function () {
      await expect(
        vrfCoordinatorV2.connect(random).defundSubscription(subId, subOwnerAddress, BigNumber.from("1000")),
      ).to.be.revertedWith(`MustBeSubOwner("${subOwnerAddress}")`);
    });
    it("insufficient balance", async function () {
      await linkToken
        .connect(subOwner)
        .transferAndCall(vrfCoordinatorV2.address, BigNumber.from("1000"), defaultAbiCoder.encode(["uint64"], [subId]));
      await expect(
        vrfCoordinatorV2.connect(subOwner).defundSubscription(subId, subOwnerAddress, BigNumber.from("1001")),
      ).to.be.revertedWith(`InsufficientBalance()`);
    });
    it("can defund", async function () {
      await linkToken
        .connect(subOwner)
        .transferAndCall(vrfCoordinatorV2.address, BigNumber.from("1000"), defaultAbiCoder.encode(["uint64"], [subId]));
      await expect(vrfCoordinatorV2.connect(subOwner).defundSubscription(subId, randomAddress, BigNumber.from("999")))
        .to.emit(vrfCoordinatorV2, "SubscriptionDefunded")
        .withArgs(subId, BigNumber.from("1000"), BigNumber.from("1"));
      const randomBalance = await linkToken.balanceOf(randomAddress);
      assert.equal(randomBalance.toString(), "1000000000000000999");
    });
  });

  describe("#cancelSubscription", async function () {
    let subId: number;
    beforeEach(async () => {
      subId = await createSubscription();
    });
    it("subscription must exist", async function () {
      await expect(
        vrfCoordinatorV2.connect(subOwner).cancelSubscription(1203123123, subOwnerAddress),
      ).to.be.revertedWith(`InvalidSubscription`);
    });
    it("must be owner", async function () {
      await expect(vrfCoordinatorV2.connect(random).cancelSubscription(subId, subOwnerAddress)).to.be.revertedWith(
        `MustBeSubOwner("${subOwnerAddress}")`,
      );
    });
    it("can cancel", async function () {
      await linkToken
        .connect(subOwner)
        .transferAndCall(vrfCoordinatorV2.address, BigNumber.from("1000"), defaultAbiCoder.encode(["uint64"], [subId]));
      await expect(vrfCoordinatorV2.connect(subOwner).cancelSubscription(subId, randomAddress))
        .to.emit(vrfCoordinatorV2, "SubscriptionCanceled")
        .withArgs(subId, randomAddress, BigNumber.from("1000"));
      const randomBalance = await linkToken.balanceOf(randomAddress);
      assert.equal(randomBalance.toString(), "1000000000000001000");
      await expect(vrfCoordinatorV2.connect(subOwner).getSubscription(subId)).to.be.revertedWith("InvalidSubscription");
    });
  });

  describe("#recoverFunds", async function () {
    let subId: number;
    beforeEach(async () => {
      subId = await createSubscription();
    });

    // Note we can't test the oracleWithdraw without fulfilling a request, so leave
    // that coverage to the go tests.
    it("function that should change internal balance do", async function () {
      type bf = [() => Promise<any>, BigNumber];
      const balanceChangingFns: Array<bf> = [
        [
          async function () {
            const s = defaultAbiCoder.encode(["uint64"], [subId]);
            await linkToken.connect(subOwner).transferAndCall(vrfCoordinatorV2.address, BigNumber.from("1000"), s);
          },
          BigNumber.from("1000"),
        ],
        [
          async function () {
            await vrfCoordinatorV2.connect(subOwner).defundSubscription(subId, randomAddress, BigNumber.from("100"));
          },
          BigNumber.from("-100"),
        ],
        [
          async function () {
            await vrfCoordinatorV2.connect(subOwner).cancelSubscription(subId, randomAddress);
          },
          BigNumber.from("-900"),
        ],
      ];
      for (let [fn, expectedBalanceChange] of balanceChangingFns) {
        const startingBalance = await vrfCoordinatorV2.s_totalBalance();
        await fn();
        const endingBalance = await vrfCoordinatorV2.s_totalBalance();
        assert(endingBalance.sub(startingBalance).toString() == expectedBalanceChange.toString());
      }
    });
    it("only owner can recover", async function () {
      await expect(vrfCoordinatorV2.connect(subOwner).recoverFunds(randomAddress)).to.be.revertedWith(
        `Only callable by owner`,
      );
    });

    it("owner can recover link transferred", async function () {
      // Set the internal balance
      assert(BigNumber.from("0"), linkToken.balanceOf(randomAddress));
      const s = defaultAbiCoder.encode(["uint64"], [subId]);
      await linkToken.connect(subOwner).transferAndCall(vrfCoordinatorV2.address, BigNumber.from("1000"), s);
      // Circumvent internal balance
      await linkToken.connect(subOwner).transfer(vrfCoordinatorV2.address, BigNumber.from("1000"));
      // Should recover this 1000
      await expect(vrfCoordinatorV2.connect(owner).recoverFunds(randomAddress))
        .to.emit(vrfCoordinatorV2, "FundsRecovered")
        .withArgs(randomAddress, BigNumber.from("1000"));
      assert(BigNumber.from("1000"), linkToken.balanceOf(randomAddress));
    });
  });

  it("subscription lifecycle", async function () {
    // Create subscription.
    let consumers: string[] = [await consumer.getAddress()];
    const tx = await vrfCoordinatorV2.connect(subOwner).createSubscription(consumers);
    const receipt = await tx.wait();
    assert(receipt.events[0].event == "SubscriptionCreated");
    assert(receipt.events[0].args["owner"] == subOwnerAddress, "sub owner");
    assert(receipt.events[0].args["consumers"][0] == consumers[0], "wrong consumers");
    const subId = receipt.events[0].args["subId"];

    // Subscription owner cannot fund
    const s = defaultAbiCoder.encode(["uint64"], [subId]);
    await expect(
      linkToken.connect(random).transferAndCall(vrfCoordinatorV2.address, BigNumber.from("1000000000000000000"), s),
    ).to.be.revertedWith(`MustBeSubOwner("${subOwnerAddress}")`);

    // Fund the subscription
    await expect(
      linkToken
        .connect(subOwner)
        .transferAndCall(
          vrfCoordinatorV2.address,
          BigNumber.from("1000000000000000000"),
          defaultAbiCoder.encode(["uint64"], [subId]),
        ),
    )
      .to.emit(vrfCoordinatorV2, "SubscriptionFunded")
      .withArgs(subId, BigNumber.from(0), BigNumber.from("1000000000000000000"));

    // Non-owners cannot withdraw
    await expect(
      vrfCoordinatorV2.connect(random).defundSubscription(subId, randomAddress, BigNumber.from("1000000000000000000")),
    ).to.be.revertedWith(`MustBeSubOwner("${subOwnerAddress}")`);

    // Withdraw from the subscription
    await expect(vrfCoordinatorV2.connect(subOwner).defundSubscription(subId, randomAddress, BigNumber.from("100")))
      .to.emit(vrfCoordinatorV2, "SubscriptionDefunded")
      .withArgs(subId, BigNumber.from("1000000000000000000"), BigNumber.from("999999999999999900"));
    const randomBalance = await linkToken.balanceOf(randomAddress);
    assert.equal(randomBalance.toString(), "1000000000000000100");

    // Non-owners cannot change the consumers
    await expect(vrfCoordinatorV2.connect(random).addConsumer(subId, randomAddress)).to.be.revertedWith(
      `MustBeSubOwner("${subOwnerAddress}")`,
    );
    await expect(vrfCoordinatorV2.connect(random).removeConsumer(subId, randomAddress)).to.be.revertedWith(
      `MustBeSubOwner("${subOwnerAddress}")`,
    );

    // Non-owners cannot ask to transfer ownership
    await expect(
      vrfCoordinatorV2.connect(random).requestSubscriptionOwnerTransfer(subId, randomAddress),
    ).to.be.revertedWith(`MustBeSubOwner("${subOwnerAddress}")`);

    // Owners can request ownership transfership
    await expect(vrfCoordinatorV2.connect(subOwner).requestSubscriptionOwnerTransfer(subId, randomAddress))
      .to.emit(vrfCoordinatorV2, "SubscriptionOwnerTransferRequested")
      .withArgs(subId, subOwnerAddress, randomAddress);

    // Non-requested owners cannot accept
    await expect(vrfCoordinatorV2.connect(subOwner).acceptSubscriptionOwnerTransfer(subId)).to.be.revertedWith(
      `MustBeRequestedOwner("${randomAddress}")`,
    );

    // Requested owners can accept
    await expect(vrfCoordinatorV2.connect(random).acceptSubscriptionOwnerTransfer(subId))
      .to.emit(vrfCoordinatorV2, "SubscriptionOwnerTransferred")
      .withArgs(subId, subOwnerAddress, randomAddress);

    // Transfer it back to subOwner
    vrfCoordinatorV2.connect(random).requestSubscriptionOwnerTransfer(subId, subOwnerAddress);
    vrfCoordinatorV2.connect(subOwner).acceptSubscriptionOwnerTransfer(subId);

    // Non-owners cannot cancel
    await expect(vrfCoordinatorV2.connect(random).cancelSubscription(subId, randomAddress)).to.be.revertedWith(
      `MustBeSubOwner("${subOwnerAddress}")`,
    );

    await expect(vrfCoordinatorV2.connect(subOwner).cancelSubscription(subId, randomAddress))
      .to.emit(vrfCoordinatorV2, "SubscriptionCanceled")
      .withArgs(subId, randomAddress, BigNumber.from("999999999999999900"));
    const random2Balance = await linkToken.balanceOf(randomAddress);
    assert.equal(random2Balance.toString(), "2000000000000000000");
  });

  it("request random words", async () => {
    // Create and fund subscription.
    let consumers: string[] = [await consumer.getAddress()];
    const tx = await vrfCoordinatorV2.connect(subOwner).createSubscription(consumers);
    const receipt = await tx.wait();
    const subId = receipt.events[0].args["subId"];
    await linkToken
      .connect(subOwner)
      .transferAndCall(
        vrfCoordinatorV2.address,
        BigNumber.from("1000000000000000000"),
        defaultAbiCoder.encode(["uint64"], [subId]),
      );

    // Should fail without a key registered
    const testKey = [BigNumber.from("1"), BigNumber.from("2")];
    let kh = await vrfCoordinatorV2.hashOfKey(testKey);
    // Non-owner cannot register a proving key
    await expect(
      vrfCoordinatorV2.connect(random).registerProvingKey(await oracle.getAddress(), [1, 2]),
    ).to.be.revertedWith("Only callable by owner");

    // Register a proving key
    await vrfCoordinatorV2.connect(owner).registerProvingKey(await oracle.getAddress(), [1, 2]);
    const realkh = await vrfCoordinatorV2.hashOfKey(testKey);
    // Cannot register the same key twice
    await expect(
      vrfCoordinatorV2.connect(owner).registerProvingKey(await oracle.getAddress(), [1, 2]),
    ).to.be.revertedWith(`KeyHashAlreadyRegistered("${realkh.toString()}")`);

    // IMPORTANT: Only registered consumers can use the subscription
    // Should fail for contract owner, sub owner, random address
    const invalidConsumers = [owner, subOwner, random];
    invalidConsumers.forEach(
      v =>
        async function () {
          await expect(
            vrfCoordinatorV2.connect(v).requestRandomWords(
              kh, // keyhash
              subId, // subId
              1, // minReqConf
              1000, // callbackGasLimit
              1, // numWords
            ),
          ).to.be.revertedWith(`InvalidConsumer(${subId}, "${v.toString()}")`);
        },
    );

    // Adding and removing a consumer should NOT allow that consumer to request
    // Non-owners cannot change the consumers
    await vrfCoordinatorV2.connect(subOwner).addConsumer(subId, randomAddress);
    await vrfCoordinatorV2.connect(subOwner).removeConsumer(subId, randomAddress);
    await expect(
      vrfCoordinatorV2.connect(random).requestRandomWords(
        kh, // keyhash
        subId, // subId
        1, // minReqConf
        1000, // callbackGasLimit
        1, // numWords
      ),
    ).to.be.revertedWith(`InvalidConsumer(${subId}, "${randomAddress.toString()}")`);

    // Should respect the minconfs
    await expect(
      vrfCoordinatorV2.connect(consumer).requestRandomWords(
        kh, // keyhash
        subId, // subId
        0, // minReqConf
        1000, // callbackGasLimit
        1, // numWords
      ),
    ).to.be.revertedWith("InvalidRequestBlockConfs(0, 1, 200)");

    // SubId must be valid
    await expect(
      vrfCoordinatorV2.connect(consumer).requestRandomWords(
        kh, // keyhash
        12398, // subId
        0, // minReqConf
        1000, // callbackGasLimit
        1, // numWords
      ),
    ).to.be.revertedWith("InvalidSubscription()");

    const reqTx = await vrfCoordinatorV2.connect(consumer).requestRandomWords(
      kh, // keyhash
      subId, // subId
      1, // minReqConf
      1000, // callbackGasLimit
      1, // numWords
    );
    const reqReceipt = await reqTx.wait();
    assert(reqReceipt.events.length == 1);
    const reqEvent = reqReceipt.events[0];
    assert(reqEvent.args["keyHash"] == kh, "wrong key hash");
    assert(reqEvent.args["subId"].toString() == subId.toString(), "wrong subId");
    assert(
      reqEvent.args["minimumRequestConfirmations"].toString() == BigNumber.from(1).toString(),
      "wrong minRequestConf",
    );
    assert(reqEvent.args["callbackGasLimit"] == 1000, "wrong callbackGasLimit");
    assert(reqEvent.args["numWords"] == 1, "wrong numWords");
    assert(reqEvent.args["sender"] == (await consumer.getAddress()), "wrong sender address");
  });

  /*
    Note that all the fulfillment testing is done in Go, to make use of the existing go code to produce
    proofs offchain.
   */
});
