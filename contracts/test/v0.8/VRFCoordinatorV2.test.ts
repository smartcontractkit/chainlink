import { ethers } from "hardhat";
import { Signer, Contract, BigNumber } from "ethers";
import { assert, expect } from "chai";

describe("VRFCoordinatorV2", () => {
  let vrfCoordinatorV2: Contract;
  let linkToken: Contract;
  let blockHashStore: Contract;
  let mockGasPrice: Contract;
  let mockLinkEth: Contract;
  let owner: Signer;
  let subOwner: Signer;
  let consumer: Signer;
  let random: Signer;
  let oracle: Signer;
  const linkEth = BigNumber.from(300000000);
  const gasWei = BigNumber.from(1e9);
  type config = {
    minimumRequestBlockConfirmations: number;
    maxConsumersPerSubscription: number;
    stalenessSeconds: number;
    gasAfterPaymentCalculation: number;
    fallbackGasPrice: BigNumber;
    fallbackLinkPrice: BigNumber;
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
    mockGasPrice = await mockAggregatorV3Factory.deploy(0, gasWei);
    mockLinkEth = await mockAggregatorV3Factory.deploy(0, linkEth);
    let vrfCoordinatorV2Factory = await ethers.getContractFactory("VRFCoordinatorV2", accounts[0]);
    vrfCoordinatorV2 = await vrfCoordinatorV2Factory.deploy(
      linkToken.address,
      blockHashStore.address,
      mockGasPrice.address,
      mockLinkEth.address,
    );
    await linkToken.transfer(await subOwner.getAddress(), BigNumber.from("1000000000000000000")); // 1 link
    c = {
      minimumRequestBlockConfirmations: 1,
      maxConsumersPerSubscription: 10,
      stalenessSeconds: 86400,
      gasAfterPaymentCalculation: 21000 + 5000 + 2100 + 20000 + 2 * 2100 - 15000 + 7315,
      fallbackGasPrice: BigNumber.from(1e9),
      fallbackLinkPrice: BigNumber.from(1e9).mul(BigNumber.from(1e7)),
    };
    // TODO: confirm only owner
    await vrfCoordinatorV2
      .connect(owner)
      .setConfig(
        c.minimumRequestBlockConfirmations,
        c.maxConsumersPerSubscription,
        c.stalenessSeconds,
        c.gasAfterPaymentCalculation,
        c.fallbackGasPrice,
        c.fallbackLinkPrice,
      );
  });

  it("subscription lifecycle", async function () {
    // Create subscription with more than max consumers should revert.
    let tooManyConsumers: string[] = new Array(c.maxConsumersPerSubscription + 1).fill(await random.getAddress());
    await expect(vrfCoordinatorV2.connect(subOwner).createSubscription(tooManyConsumers)).to.be.revertedWith(
      "" + ">max consumers per sub",
    );

    // Create subscription.
    let consumers: string[] = [await consumer.getAddress()];
    const tx = await vrfCoordinatorV2.connect(subOwner).createSubscription(consumers);
    const receipt = await tx.wait();
    const subId = receipt.events[0].args["subId"];

    // Subscription owner cannot fund
    await expect(
      vrfCoordinatorV2.connect(random).fundSubscription(subId, BigNumber.from("1000000000000000000")),
    ).to.be.revertedWith("sub owner must fund");

    // Fund the subscription
    await linkToken.connect(subOwner).approve(vrfCoordinatorV2.address, BigNumber.from("1000000000000000000"));
    await linkToken.allowance(await subOwner.getAddress(), vrfCoordinatorV2.address);
    await vrfCoordinatorV2.connect(subOwner).fundSubscription(subId, BigNumber.from("1000000000000000000"));

    // Non-owners cannot withdraw
    await expect(
      vrfCoordinatorV2
        .connect(random)
        .withdrawFromSubscription(subId, await random.getAddress(), BigNumber.from("1000000000000000000")),
    ).to.be.revertedWith("sub owner must withdraw");

    // Withdraw from the subscription
    const withdrawTx = await vrfCoordinatorV2
      .connect(subOwner)
      .withdrawFromSubscription(subId, await random.getAddress(), BigNumber.from("100"));
    await withdrawTx.wait();
    const randomBalance = await linkToken.balanceOf(await random.getAddress());
    assert.equal(randomBalance.toString(), "100");

    // Non-owners cannot cancel
    await expect(vrfCoordinatorV2.connect(random).cancelSubscription(subId)).to.be.revertedWith(
      "sub owner must cancel",
    );
    // Cannot cancel sub with funds
    await expect(vrfCoordinatorV2.connect(subOwner).cancelSubscription(subId)).to.be.revertedWith("balance != 0");

    // Withdraw remaining balance then cancel
    let sub = await vrfCoordinatorV2.connect(subOwner).getSubscription(subId);
    const withdraw2Tx = await vrfCoordinatorV2
      .connect(subOwner)
      .withdrawFromSubscription(subId, await random.getAddress(), sub.balance);
    await withdraw2Tx.wait();
    const random2Balance = await linkToken.balanceOf(await random.getAddress());
    assert.equal(random2Balance.toString(), "1000000000000000000");
    await vrfCoordinatorV2.connect(subOwner).cancelSubscription(subId);
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
    await expect(vrfCoordinatorV2.connect(consumer).requestRandomWords(kh, 1, 1000, subId, 1)).to.be.revertedWith(
      "must be a registered key",
    );

    // Non-owner cannot register a proving key
    await expect(
      vrfCoordinatorV2.connect(random).registerProvingKey(await oracle.getAddress(), [1, 2]),
    ).to.be.revertedWith("caller is not the owner");

    // Register a proving key
    await vrfCoordinatorV2.connect(owner).registerProvingKey(await oracle.getAddress(), [1, 2]);
    // Cannot register the same key twice
    await expect(
      vrfCoordinatorV2.connect(owner).registerProvingKey(await oracle.getAddress(), [1, 2]),
    ).to.be.revertedWith("key already registered");

    const reqTx = await vrfCoordinatorV2.connect(consumer).requestRandomWords(
      kh, // keyhash
      1, // minReqConf
      1000, // callbackGasLimit
      subId, // subId
      1, // numWords
    );
    const reqReceipt = await reqTx.wait();
    const reqId = reqReceipt.events[0].args["preSeed"];
    console.log(reqId);

    // TODO: Should see the request event
    // TODO: Should respect minReqConfs
  });

  /*
    Note that all the fulfillment testing is done in Go, to make use of the existing go code to produce
    proofs offchain.
   */
});
