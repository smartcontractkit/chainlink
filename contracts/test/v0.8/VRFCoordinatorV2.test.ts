import { ethers } from "hardhat";
import { Signer, Contract, BigNumber } from "ethers";
import { assert, expect } from "chai";

describe("VRFCoordinatorV2", () => {
  let vrfCoordinatorV2: Contract;
  let linkToken: Contract;
  let blockHashStore: Contract;
  let mockLinkEth: Contract;
  let owner: Signer;
  let subOwner: Signer;
  let consumer: Signer;
  let random: Signer;
  let oracle: Signer;
  const linkEth = BigNumber.from(300000000);
  type config = {
    minimumRequestBlockConfirmations: number;
    maxConsumersPerSubscription: number;
    stalenessSeconds: number;
    gasAfterPaymentCalculation: number;
    fallbackLinkPrice: BigNumber;
    minimumSubscriptionBalance: BigNumber;
  };
  let c: config;

  beforeEach(async () => {
    let accounts = await ethers.getSigners();
    owner = accounts[0];
    subOwner = accounts[1];
    consumer = accounts[2];
    random = accounts[3];
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
    await linkToken.transfer(await subOwner.getAddress(), BigNumber.from("1000000000000000000")); // 1 link
    c = {
      minimumRequestBlockConfirmations: 1,
      maxConsumersPerSubscription: 10,
      stalenessSeconds: 86400,
      gasAfterPaymentCalculation: 21000 + 5000 + 2100 + 20000 + 2 * 2100 - 15000 + 7315,
      fallbackLinkPrice: BigNumber.from("10000000000000000"),
      minimumSubscriptionBalance: BigNumber.from("100000000000000000"), // 0.1 link
    };
    await vrfCoordinatorV2
      .connect(owner)
      .setConfig(
        c.minimumRequestBlockConfirmations,
        c.maxConsumersPerSubscription,
        c.stalenessSeconds,
        c.gasAfterPaymentCalculation,
        c.fallbackLinkPrice,
        c.minimumSubscriptionBalance,
      );
  });

  it("configuration management", async function () {
    await expect(
      vrfCoordinatorV2
        .connect(subOwner)
        .setConfig(
          c.minimumRequestBlockConfirmations,
          c.maxConsumersPerSubscription,
          c.stalenessSeconds,
          c.gasAfterPaymentCalculation,
          c.fallbackLinkPrice,
          c.minimumSubscriptionBalance,
        ),
    ).to.be.revertedWith("Only callable by owner");
    // Anyone can read the config.
    const resp = await vrfCoordinatorV2.connect(random).getConfig();
    assert(resp[0] == c.minimumRequestBlockConfirmations);
    assert(resp[1] == c.maxConsumersPerSubscription);
    assert(resp[2] == c.stalenessSeconds);
    assert(resp[3].toString() == c.gasAfterPaymentCalculation.toString());
    assert(resp[4].toString() == c.fallbackLinkPrice.toString());
  });

  it("subscription lifecycle", async function () {
    // Create subscription with more than max consumers should revert.
    let tooManyConsumers: string[] = new Array(c.maxConsumersPerSubscription + 1).fill(await random.getAddress());
    await expect(vrfCoordinatorV2.connect(subOwner).createSubscription(tooManyConsumers)).to.be.revertedWith(
      "InvalidNumberOfConsumers(11, 10)",
    );

    // Create subscription.
    let consumers: string[] = [await consumer.getAddress()];
    const tx = await vrfCoordinatorV2.connect(subOwner).createSubscription(consumers);
    const receipt = await tx.wait();
    assert(receipt.events[0].event == "SubscriptionCreated");
    assert(receipt.events[0].args["owner"] == (await subOwner.getAddress()), "sub owner");
    assert(receipt.events[0].args["consumers"][0] == consumers[0], "wrong consumers");
    const subId = receipt.events[0].args["subId"];

    // Subscription owner cannot fund
    await expect(
      vrfCoordinatorV2.connect(random).fundSubscription(subId, BigNumber.from("1000000000000000000")),
    ).to.be.revertedWith("MustBeSubOwner()");

    // Fund the subscription
    await linkToken.connect(subOwner).approve(vrfCoordinatorV2.address, BigNumber.from("1000000000000000000"));
    await linkToken.allowance(await subOwner.getAddress(), vrfCoordinatorV2.address);
    await expect(vrfCoordinatorV2.connect(subOwner).fundSubscription(subId, BigNumber.from("1000000000000000000")))
      .to.emit(vrfCoordinatorV2, "SubscriptionFundsAdded")
      .withArgs(subId, BigNumber.from(0), BigNumber.from("1000000000000000000"));

    // Non-owners cannot withdraw
    await expect(
      vrfCoordinatorV2
        .connect(random)
        .withdrawFromSubscription(subId, await random.getAddress(), BigNumber.from("1000000000000000000")),
    ).to.be.revertedWith("MustBeSubOwner()");

    // Withdraw from the subscription
    await expect(
      vrfCoordinatorV2
        .connect(subOwner)
        .withdrawFromSubscription(subId, await random.getAddress(), BigNumber.from("100")),
    )
      .to.emit(vrfCoordinatorV2, "SubscriptionFundsWithdrawn")
      .withArgs(subId, BigNumber.from("1000000000000000000"), BigNumber.from("999999999999999900"));
    const randomBalance = await linkToken.balanceOf(await random.getAddress());
    assert.equal(randomBalance.toString(), "100");

    // Non-owners cannot change the consumers
    await expect(vrfCoordinatorV2.connect(random).updateSubscription(subId, consumers)).to.be.revertedWith(
      "MustBeSubOwner()",
    );
    // Owners can
    await vrfCoordinatorV2.connect(subOwner).updateSubscription(subId, consumers);

    // Non-owners cannot cancel
    await expect(
      vrfCoordinatorV2.connect(random).cancelSubscription(subId, await random.getAddress()),
    ).to.be.revertedWith("MustBeSubOwner()");

    const randomAddress = await random.getAddress();
    await expect(vrfCoordinatorV2.connect(subOwner).cancelSubscription(subId, randomAddress))
      .to.emit(vrfCoordinatorV2, "SubscriptionCanceled")
      .withArgs(subId, randomAddress, BigNumber.from("999999999999999900"));
    const random2Balance = await linkToken.balanceOf(randomAddress);
    assert.equal(random2Balance.toString(), "1000000000000000000");
  });

  it("request random words", async () => {
    // Create and fund subscription.
    let consumers: string[] = [await consumer.getAddress()];
    const tx = await vrfCoordinatorV2.connect(subOwner).createSubscription(consumers);
    const receipt = await tx.wait();
    const subId = receipt.events[0].args["subId"];
    await linkToken.connect(subOwner).approve(vrfCoordinatorV2.address, BigNumber.from("1000000000000000000"));
    await linkToken.allowance(await subOwner.getAddress(), vrfCoordinatorV2.address);
    await vrfCoordinatorV2.connect(subOwner).fundSubscription(subId, BigNumber.from("1000000000000000000"));

    // Should fail without a key registered
    const testKey = [BigNumber.from("1"), BigNumber.from("2")];
    let kh = await vrfCoordinatorV2.hashOfKey(testKey);
    await expect(vrfCoordinatorV2.connect(consumer).requestRandomWords(kh, 1, 1000, subId, 1, 0)).to.be.revertedWith(
      `UnregisteredKeyHash("${kh.toString()}")`,
    );
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
    ).to.be.revertedWith("RequestBlockConfsTooLow(0, 1)");

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
