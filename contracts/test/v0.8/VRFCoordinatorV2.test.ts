import { ethers } from "hardhat";
import { Signer, Contract, BigNumber } from "ethers";

describe("VRFCoordinatorV2", () => {
  let vrfCoordinatorV2: Contract;
  let linkToken: Contract;
  let blockHashStore: Contract;
  let mockGasPrice: Contract;
  let mockLinkEth: Contract;
  let owner: Signer;
  let subOwner: Signer;
  let consumer: Signer;
  // let random: Signer;
  let oracle: Signer;
  const linkEth = BigNumber.from(300000000);
  const gasWei = BigNumber.from(1e9);

  beforeEach(async () => {
    let accounts = await ethers.getSigners();
    owner = accounts[0];
    subOwner = accounts[1];
    consumer = accounts[2];
    // random = accounts[3]
    oracle = accounts[4];
    let ltFactory = await ethers.getContractFactory("LinkToken", accounts[0]);
    linkToken = await ltFactory.deploy();
    let bhFactory = await ethers.getContractFactory("BlockhashStore", accounts[0]);
    blockHashStore = await bhFactory.deploy();
    let mockAggregatorV3Factory = await ethers.getContractFactory("MockV3Aggregator", accounts[0]);
    mockGasPrice = await mockAggregatorV3Factory.deploy(0, gasWei);
    mockLinkEth = await mockAggregatorV3Factory.deploy(0, linkEth);
    let vrfCoordinatorV2Factory = await ethers.getContractFactory("VRFCoordinatorV2", accounts[0]);
    console.log("link address", linkToken.address);
    vrfCoordinatorV2 = await vrfCoordinatorV2Factory.deploy(
      linkToken.address,
      blockHashStore.address,
      mockGasPrice.address,
      mockLinkEth.address,
    );
    console.log("link balance", await linkToken.balanceOf(await owner.getAddress()));
    await linkToken.transfer(await subOwner.getAddress(), BigNumber.from("1000000000000000000")); // 1 link
    console.log("link balance", await linkToken.balanceOf(await subOwner.getAddress()));
    //        uint16 minimumRequestBlockConfirmations,
    //         uint16 maxConsumersPerSubscription,
    //         uint32 stalenessSeconds,
    //         uint32 gasAfterPaymentCalculation,
    //         int256 fallbackGasPrice,
    //         int256 fallbackLinkPrice
    await vrfCoordinatorV2.setConfig(
      1,
      10,
      86400,
      21000 + 5000 + 2100 + 20000 + 2 * 2100 - 15000 + 7315,
      BigNumber.from(1e9),
      BigNumber.from(1e9).mul(BigNumber.from(1e7)),
    );
  });

  it("subscription lifecycle", async () => {
    // Create subscription.
    // const response = await vrfCoordinatorV2.owner()
    let consumers: string[] = [await consumer.getAddress()];
    const tx = await vrfCoordinatorV2.connect(subOwner).createSubscription(consumers);
    const receipt = await tx.wait();
    const subId = receipt.events[0].args["subId"];

    // Fund the subscription
    await linkToken.connect(subOwner).approve(vrfCoordinatorV2.address, BigNumber.from("1000000000000000000"));
    const resp = await linkToken.allowance(await subOwner.getAddress(), vrfCoordinatorV2.address);
    console.log(resp);
    await vrfCoordinatorV2.connect(subOwner).fundSubscription(subId, BigNumber.from("1000000000000000000"));
    // TODO: non-owners cannot fund
    // TODO: withdraw funds, non-owners cannot
    // TODO: cancel sub, non-owners cannot
  });

  // TODO: request words, check logs, non consumer cannot request
  it("request random words", async () => {
    // Create subscription.
    // const response = await vrfCoordinatorV2.owner()
    let consumers: string[] = [await consumer.getAddress()];
    const tx = await vrfCoordinatorV2.connect(subOwner).createSubscription(consumers);
    const receipt = await tx.wait();
    const subId = receipt.events[0].args["subId"];

    // Fund the subscription
    await linkToken.connect(subOwner).approve(vrfCoordinatorV2.address, BigNumber.from("1000000000000000000"));
    const resp = await linkToken.allowance(await subOwner.getAddress(), vrfCoordinatorV2.address);
    console.log(resp);
    await vrfCoordinatorV2.connect(subOwner).fundSubscription(subId, BigNumber.from("1000000000000000000"));

    // Request random words
    //        bytes32 keyHash,  // Corresponds to a particular offchain job which uses that key for the proofs
    //         uint16  minimumRequestConfirmations,
    //         uint16  callbackGasLimit,
    //         uint256 subId,   // A data structure for billing
    //         uint256 numWords  // Desired number of random words
    const testKey = [BigNumber.from("1"), BigNumber.from("2")];
    let kh = await vrfCoordinatorV2.hashOfKey(testKey);
    console.log("key hash", kh);
    await vrfCoordinatorV2.registerProvingKey(await oracle.getAddress(), [1, 2]);
    const reqTx = await vrfCoordinatorV2.connect(consumer).requestRandomWords(kh, 1, 1000, subId, 1);
    const reqReceipt = await reqTx.wait();
    const reqId = reqReceipt.events[0].args["preSeed"];
    console.log(reqId);
    //Should see the callback
    // console.log(await vrfCoordinatorV2.s_callbacks(reqId))
    // // 265747905000000
    // const r = await vrfCoordinatorV2.calculatePaymentAmount(100000000);
    // const rr = await r.wait()
    // const s = rr .events[0].args['seed']
    // console.log("payment", s.integerValue())
    // console.log("payment", r);
  });
});
