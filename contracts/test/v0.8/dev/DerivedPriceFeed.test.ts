import { ethers } from 'hardhat'
import { BigNumber, ContractFactory } from 'ethers'
import { expect } from 'chai'
import { describe } from 'mocha'

describe('DerivedPriceFeed', () => {
  let mockAggFactory: ContractFactory
  let derivedFeedFactory: ContractFactory
  before(async () => {
    const accounts = await ethers.getSigners()
    mockAggFactory = await ethers.getContractFactory(
      'src/v0.7/tests/MockV3Aggregator.sol:MockV3Aggregator',
      accounts[0],
    )
    derivedFeedFactory = await ethers.getContractFactory(
      'src/v0.8/dev/DerivedPriceFeed.sol:DerivedPriceFeed',
      accounts[0],
    )
  })

  it('reverts on getRoundData', async () => {
    let base = await mockAggFactory.deploy(8, 10e8) // Price = 10
    let quote = await mockAggFactory.deploy(8, 5e8) // Price = 5

    let derived = await derivedFeedFactory.deploy(
      base.address,
      quote.address,
      8,
    )

    await expect(derived.getRoundData(1)).to.be.reverted
  })

  it('returns decimals', async () => {
    let base = await mockAggFactory.deploy(8, 10e8) // Price = 10
    let quote = await mockAggFactory.deploy(8, 5e8) // Price = 5

    let derived = await derivedFeedFactory.deploy(
      base.address,
      quote.address,
      9,
    )

    await expect(await derived.decimals()).to.equal(9)
  })

  describe('calculates price', async () => {
    it('when all decimals are the same', async () => {
      let base = await mockAggFactory.deploy(8, 10e8) // 10
      let quote = await mockAggFactory.deploy(8, 5e8) // 5

      let derived = await derivedFeedFactory.deploy(
        base.address,
        quote.address,
        8,
      )

      await expect((await derived.latestRoundData()).answer).to.equal(
        2e8 /* 2 */,
      )
    })

    it('when all decimals are the same 2', async () => {
      let base = await mockAggFactory.deploy(8, 3e8) // 3
      let quote = await mockAggFactory.deploy(8, 15e8) // 15

      let derived = await derivedFeedFactory.deploy(
        base.address,
        quote.address,
        8,
      )

      await expect((await derived.latestRoundData()).answer).to.equal(
        0.2e8 /* 0.2 */,
      )
    })

    it('when result decimals are higher', async () => {
      let base = await mockAggFactory.deploy(8, 10e8) // Price = 10
      let quote = await mockAggFactory.deploy(8, 5e8) // Price = 5

      let derived = await derivedFeedFactory.deploy(
        base.address,
        quote.address,
        12,
      )

      await expect((await derived.latestRoundData()).answer).to.equal(
        2e12 /* 2 */,
      )
    })

    it('when result decimals are lower', async () => {
      let base = await mockAggFactory.deploy(8, 10e8) // Price = 10
      let quote = await mockAggFactory.deploy(8, 5e8) // Price = 5

      let derived = await derivedFeedFactory.deploy(
        base.address,
        quote.address,
        6,
      )

      await expect((await derived.latestRoundData()).answer).to.equal(
        2e6 /* 2 */,
      )
    })

    it('base decimals are higher', async () => {
      let base = await mockAggFactory.deploy(
        16,
        BigNumber.from('100000000000000000'),
      ) // Price = 10
      let quote = await mockAggFactory.deploy(8, 5e8) // Price = 5

      let derived = await derivedFeedFactory.deploy(
        base.address,
        quote.address,
        10,
      )

      await expect((await derived.latestRoundData()).answer).to.equal(
        2e10 /* 2 */,
      )
    })

    it('base decimals are lower', async () => {
      let base = await mockAggFactory.deploy(6, 10e6) // Price = 10
      let quote = await mockAggFactory.deploy(8, 5e8) // Price = 5

      let derived = await derivedFeedFactory.deploy(
        base.address,
        quote.address,
        10,
      )

      await expect((await derived.latestRoundData()).answer).to.equal(
        2e10 /* 2 */,
      )
    })
  })
})
