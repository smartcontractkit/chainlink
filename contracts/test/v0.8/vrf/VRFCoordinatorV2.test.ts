import { ethers } from "hardhat";
import {Signer, Contract, BigNumber} from "ethers";

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
    const linkEth = BigNumber.from(300000000);
    const gasWei = BigNumber.from(100);

    beforeEach(async () => {
        let accounts = await ethers.getSigners();
        owner = accounts[0]
        subOwner = accounts[1]
        consumer = accounts[2]
        random = accounts[3]
        let ltFactory = await ethers.getContractFactory("LinkToken", accounts[0]);
        linkToken = await ltFactory.deploy()
        let bhFactory = await ethers.getContractFactory("BlockhashStore", accounts[0]);
        blockHashStore = await bhFactory.deploy();
        let mockAggregatorV3Factory = await ethers.getContractFactory("MockV3Aggregator", accounts[0])
        mockGasPrice = await mockAggregatorV3Factory.deploy(0, gasWei)
        mockLinkEth = await mockAggregatorV3Factory.deploy(0, linkEth)
        let vrfCoordinatorV2Factory = await ethers.getContractFactory("VRFCoordinatorV2", accounts[0]);
        console.log("link address", linkToken.address)
        vrfCoordinatorV2 = await vrfCoordinatorV2Factory.deploy(linkToken.address, blockHashStore.address, mockGasPrice.address, mockLinkEth.address);
        console.log("link balance", await linkToken.balanceOf(await owner.getAddress()))
        await linkToken.transfer(await subOwner.getAddress(), BigNumber.from("1000000000000000000")) // 1 link
        console.log("link balance", await linkToken.balanceOf(await subOwner.getAddress()))
    });

    it("subscription lifecycle", async () => {
        // Create subscription.
        const response = await vrfCoordinatorV2.owner()
        let consumers: string[] = [await consumer.getAddress()]
        const tx = await vrfCoordinatorV2.connect(subOwner).createSubscription(consumers)
        const receipt = await tx.wait()
        const subId = receipt.events[0].args['subId']

        // Fund the subscription
        await linkToken.connect(subOwner).approve(vrfCoordinatorV2.address, BigNumber.from("1000000000000000000"))
        const resp = await linkToken.allowance(await subOwner.getAddress(), vrfCoordinatorV2.address)
        console.log(resp)
        await vrfCoordinatorV2.connect(subOwner).fundSubscription(subId, BigNumber.from("1000000000000000000"))
        // TODO: non-owners cannot fund
        // TODO: withdraw funds, non-owners cannot
        // TODO: cancel sub, non-owners cannot
    });

    // TODO: request words, check logs, non consumer cannot request
})
