import { ethers } from "hardhat";
import { assert, expect } from "chai";
import { Signer, Contract } from "ethers";

describe("AggregatorV2V3Interface Solidity version tests", () => {
  let accounts: Signer[];
  let owner: Signer;

  let mockAggregator: Contract;
  let consumerV6: Contract;
  let consumerV7: Contract;
  let consumerV8: Contract;

  beforeEach(async () => {
    accounts = await ethers.getSigners();
    owner = accounts[0];

    const mockAggregatorFactory = await ethers.getContractFactory("MockAggregator", owner);
    const v6Factory = await ethers.getContractFactory("AggregatorInterfaceConsumerTest6");
    const v7Factory = await ethers.getContractFactory("AggregatorInterfaceConsumerTest7");
    const v8Factory = await ethers.getContractFactory("AggregatorInterfaceConsumerTest8");

    mockAggregator = await mockAggregatorFactory.deploy();
    consumerV6 = await v6Factory.deploy(mockAggregator.address);
    consumerV7 = await v7Factory.deploy(mockAggregator.address);
    consumerV8 = await v8Factory.deploy(mockAggregator.address);
  })

  it('should work using a Solidity 0.6.x consumer', async () => {
    const answer = 123;
    mockAggregator.connect(owner).setLatestAnswer(answer);
    const response = await consumerV6.getLatestPrice();
    assert.equal(response.toNumber(), answer);
  })

  it('should work using a Solidity 0.7.x consumer', async () => {
    const answer = 234;
    mockAggregator.connect(owner).setLatestAnswer(answer);
    const response = await consumerV7.getLatestPrice();
    assert.equal(response.toNumber(), answer);
  })

  it('should work using a Solidity 0.8.x consumer', async () => {
    const answer = 345;
    mockAggregator.connect(owner).setLatestAnswer(answer);
    const response = await consumerV8.getLatestPrice();
    assert.equal(response.toNumber(), answer);
  })
})