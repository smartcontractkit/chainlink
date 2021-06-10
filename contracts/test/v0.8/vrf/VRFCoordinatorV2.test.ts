import { ethers } from "hardhat";
import { Signer, Contract } from "ethers";

describe("VRFCoordinatorV2", () => {
    let vrfCoordinatorV2: Contract;
    let owner: Signer;
    let subOwner: Signer;
    let consumer: Signer;

    beforeEach(async () => {
        let accounts = await ethers.getSigners();
        owner = accounts[0]
        subOwner = accounts[1]
        consumer = accounts[2]
        let ltFactory = await ethers.getContractFactory("LinkToken", accounts[0]);
        let lt = await ltFactory.deploy();
        lt = await lt.deployed();
        let bhFactory = await ethers.getContractFactory("BlockhashStore", accounts[0]);
        let bh = await bhFactory.deploy();
        bh = await bh.deployed();
        let vrfCoordinatorV2Factory = await ethers.getContractFactory("VRFCoordinatorV2", accounts[0]);
        vrfCoordinatorV2 = await vrfCoordinatorV2Factory.deploy(lt.address, bh.address);
        vrfCoordinatorV2 = await vrfCoordinatorV2.deployed();
    });

    it("#subscription", async () => {
        // Lets create a subscription
        console.log("address", vrfCoordinatorV2.address)
        const response = await vrfCoordinatorV2.owner()
        console.log(response, await owner.getAddress())
        let consumers: string[] = [await consumer.getAddress()]
        await vrfCoordinatorV2.createSubscription(consumers)
        console.log(await vrfCoordinatorV2.subscriptions(1))
    });
})
