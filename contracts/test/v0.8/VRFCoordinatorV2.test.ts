import { ethers } from "hardhat";
import { Signer, Contract } from "ethers";

describe("VRFCoordinatorV2", () => {
    beforeEach(async () => {
        let accounts = await ethers.getSigners();
        console.log(accounts[0])
        const vpf = await ethers.getContractFactory("VRFCoordinatorV2", accounts[0]);
        let vrfCoordinatorV2 = await vpf.deploy();
        vrfCoordinatorV2 = await vrfCoordinatorV2.deployed();
        console.log(vrfCoordinatorV2.address)
    });
})
