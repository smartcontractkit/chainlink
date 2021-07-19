import { ethers } from "hardhat";
import { Signer, Contract, BigNumber } from "ethers";
import { assert, expect } from "chai";
import {defaultAbiCoder} from "ethers/utils";

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
    const s = defaultAbiCoder.encode(["uint64"], [subId])
    await expect(linkToken.connect(random).transferAndCall(vrfCoordinatorV2.address, BigNumber.from("1000000000000000000"),
        s)).to.be.revertedWith(`MustBeSubOwner("${subOwnerAddress}")`);

    // Fund the subscription
    await expect(linkToken.connect(subOwner).transferAndCall(vrfCoordinatorV2.address, BigNumber.from("1000000000000000000"),
        defaultAbiCoder.encode(["uint64"], [subId])))
          .to.emit(vrfCoordinatorV2, "SubscriptionFundsAdded")
          .withArgs(subId, BigNumber.from(0), BigNumber.from("1000000000000000000"));

    // Non-owners cannot withdraw
    await expect(
      vrfCoordinatorV2.connect(random).defundSubscription(subId, randomAddress, BigNumber.from("1000000000000000000")),
    ).to.be.revertedWith(`MustBeSubOwner("${subOwnerAddress}")`);

    // Withdraw from the subscription
    await expect(vrfCoordinatorV2.connect(subOwner).defundSubscription(subId, randomAddress, BigNumber.from("100")))
      .to.emit(vrfCoordinatorV2, "SubscriptionFundsWithdrawn")
      .withArgs(subId, BigNumber.from("1000000000000000000"), BigNumber.from("999999999999999900"));
    const randomBalance = await linkToken.balanceOf(randomAddress);
    assert.equal(randomBalance.toString(), "1000000000000000100");

    // Non-owners cannot change the consumers
    await expect(vrfCoordinatorV2.connect(random).updateSubscription(subId, consumers)).to.be.revertedWith(
      `MustBeSubOwner("${subOwnerAddress}")`,
    );
    // Owners can update
    await vrfCoordinatorV2.connect(subOwner).updateSubscription(subId, consumers);

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
    await linkToken.connect(subOwner).transferAndCall(vrfCoordinatorV2.address, BigNumber.from("1000000000000000000"),
        defaultAbiCoder.encode(["uint64"], [subId]));

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
              0,
            ),
          ).to.be.revertedWith(`InvalidConsumer("${v.toString()}")`);
        },
    );

    // Should respect the minconfs
    await expect(
      vrfCoordinatorV2.connect(consumer).requestRandomWords(
        kh, // keyhash
        subId, // subId
        0, // minReqConf
        1000, // callbackGasLimit
        1, // numWords
        0,
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
        0,
      ),
    ).to.be.revertedWith("InvalidSubscription()");

    const reqTx = await vrfCoordinatorV2.connect(consumer).requestRandomWords(
      kh, // keyhash
      subId, // subId
      1, // minReqConf
      1000, // callbackGasLimit
      1, // numWords
      0,
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
