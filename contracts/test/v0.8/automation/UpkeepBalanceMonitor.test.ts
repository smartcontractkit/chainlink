import { ethers } from 'hardhat'
import { expect } from 'chai'
import type { SignerWithAddress } from '@nomiclabs/hardhat-ethers/signers'
import { randomAddress } from '../../test-helpers/helpers'
import { loadFixture } from '@nomicfoundation/hardhat-network-helpers'
import { IKeeperRegistryMaster__factory as RegistryFactory } from '../../../typechain/factories/IKeeperRegistryMaster__factory'
import { IAutomationForwarder__factory as ForwarderFactory } from '../../../typechain/factories/IAutomationForwarder__factory'
import { UpkeepBalanceMonitor } from '../../../typechain/UpkeepBalanceMonitor'
import { LinkToken } from '../../../typechain/LinkToken'
import { BigNumber } from 'ethers'
import {
  deployMockContract,
  MockContract,
} from '@ethereum-waffle/mock-contract'

let owner: SignerWithAddress
let stranger: SignerWithAddress
let registry: MockContract
let registry2: MockContract
let forwarder: MockContract
let linkToken: LinkToken
let upkeepBalanceMonitor: UpkeepBalanceMonitor

const setup = async () => {
  const accounts = await ethers.getSigners()
  owner = accounts[0]
  stranger = accounts[1]

  const ltFactory = await ethers.getContractFactory(
    'src/v0.8/shared/test/helpers/LinkTokenTestHelper.sol:LinkTokenTestHelper',
    owner,
  )
  linkToken = (await ltFactory.deploy()) as LinkToken
  const bmFactory = await ethers.getContractFactory(
    'UpkeepBalanceMonitor',
    owner,
  )
  upkeepBalanceMonitor = await bmFactory.deploy(linkToken.address, {
    maxBatchSize: 10,
    minPercentage: 120,
    targetPercentage: 300,
    maxTopUpAmount: ethers.utils.parseEther('100'),
  })
  registry = await deployMockContract(owner, RegistryFactory.abi)
  registry2 = await deployMockContract(owner, RegistryFactory.abi)
  forwarder = await deployMockContract(owner, ForwarderFactory.abi)
  await forwarder.mock.getRegistry.returns(registry.address)
  await upkeepBalanceMonitor.setForwarder(forwarder.address)
  await linkToken
    .connect(owner)
    .transfer(upkeepBalanceMonitor.address, ethers.utils.parseEther('10000'))
  await upkeepBalanceMonitor
    .connect(owner)
    .setWatchList(registry.address, [0, 1, 2, 3, 4, 5, 6, 7, 8])
  await upkeepBalanceMonitor
    .connect(owner)
    .setWatchList(registry2.address, [9, 10, 11])
  for (let i = 0; i < 9; i++) {
    await registry.mock.getMinBalance.withArgs(i).returns(100)
    await registry.mock.getBalance.withArgs(i).returns(121) // all upkeeps are sufficiently funded
  }
  for (let i = 9; i < 12; i++) {
    await registry2.mock.getMinBalance.withArgs(i).returns(100)
    await registry2.mock.getBalance.withArgs(i).returns(121) // all upkeeps are sufficiently funded
  }
}

describe('UpkeepBalanceMonitor', () => {
  beforeEach(async () => {
    await loadFixture(setup)
  })

  describe('constructor()', () => {
    it('should set the initial values correctly', async () => {
      const config = await upkeepBalanceMonitor.getConfig()
      expect(config.maxBatchSize).to.equal(10)
      expect(config.minPercentage).to.equal(120)
      expect(config.targetPercentage).to.equal(300)
      expect(config.maxTopUpAmount).to.equal(ethers.utils.parseEther('100'))
    })
  })

  describe('setConfig()', () => {
    const newConfig = {
      maxBatchSize: 100,
      minPercentage: 150,
      targetPercentage: 500,
      maxTopUpAmount: 1,
    }

    it('should set config correctly', async () => {
      await upkeepBalanceMonitor.connect(owner).setConfig(newConfig)
      const config = await upkeepBalanceMonitor.getConfig()
      expect(config.maxBatchSize).to.equal(newConfig.maxBatchSize)
      expect(config.minPercentage).to.equal(newConfig.minPercentage)
      expect(config.targetPercentage).to.equal(newConfig.targetPercentage)
      expect(config.maxTopUpAmount).to.equal(newConfig.maxTopUpAmount)
    })

    it('cannot be called by a non-owner', async () => {
      await expect(
        upkeepBalanceMonitor.connect(stranger).setConfig(newConfig),
      ).to.be.revertedWith('Only callable by owner')
    })

    it('should emit an event', async () => {
      await expect(
        upkeepBalanceMonitor.connect(owner).setConfig(newConfig),
      ).to.emit(upkeepBalanceMonitor, 'ConfigSet')
    })
  })

  describe('setForwarder()', () => {
    const newForwarder = randomAddress()

    it('should set the forwarder correctly', async () => {
      await upkeepBalanceMonitor.connect(owner).setForwarder(newForwarder)
      const forwarderAddress = await upkeepBalanceMonitor.getForwarder()
      expect(forwarderAddress).to.equal(newForwarder)
    })

    it('cannot be called by a non-owner', async () => {
      await expect(
        upkeepBalanceMonitor.connect(stranger).setForwarder(randomAddress()),
      ).to.be.revertedWith('Only callable by owner')
    })

    it('should emit an event', async () => {
      await expect(
        upkeepBalanceMonitor.connect(owner).setForwarder(newForwarder),
      )
        .to.emit(upkeepBalanceMonitor, 'ForwarderSet')
        .withArgs(newForwarder)
    })
  })

  describe('setWatchList()', () => {
    const newWatchList = [
      BigNumber.from(1),
      BigNumber.from(2),
      BigNumber.from(10),
    ]

    it('should add addresses to the watchlist', async () => {
      await upkeepBalanceMonitor
        .connect(owner)
        .setWatchList(registry.address, newWatchList)
      const [_, upkeepIDs] = await upkeepBalanceMonitor.getWatchList()
      expect(upkeepIDs[0]).to.deep.equal(newWatchList)
    })

    it('cannot be called by a non-owner', async () => {
      await expect(
        upkeepBalanceMonitor
          .connect(stranger)
          .setWatchList(registry.address, [1, 2, 3]),
      ).to.be.revertedWith('Only callable by owner')
    })

    it('should emit an event', async () => {
      await expect(
        upkeepBalanceMonitor
          .connect(owner)
          .setWatchList(registry.address, newWatchList),
      )
        .to.emit(upkeepBalanceMonitor, 'WatchListSet')
        .withArgs(registry.address)
    })
  })

  describe('withdraw()', () => {
    const payee = randomAddress()
    const withdrawAmount = 100

    it('should withdraw funds to a payee', async () => {
      const initialBalance = await linkToken.balanceOf(
        upkeepBalanceMonitor.address,
      )
      await upkeepBalanceMonitor.connect(owner).withdraw(withdrawAmount, payee)
      const finalBalance = await linkToken.balanceOf(
        upkeepBalanceMonitor.address,
      )
      const payeeBalance = await linkToken.balanceOf(payee)
      expect(finalBalance).to.equal(initialBalance.sub(withdrawAmount))
      expect(payeeBalance).to.equal(withdrawAmount)
    })

    it('cannot be called by a non-owner', async () => {
      await expect(
        upkeepBalanceMonitor.connect(stranger).withdraw(withdrawAmount, payee),
      ).to.be.revertedWith('Only callable by owner')
    })

    it('should emit an event', async () => {
      await expect(
        upkeepBalanceMonitor.connect(owner).withdraw(withdrawAmount, payee),
      )
        .to.emit(upkeepBalanceMonitor, 'FundsWithdrawn')
        .withArgs(100, payee)
    })
  })

  describe('pause() and unpause()', () => {
    it('should pause and unpause the contract', async () => {
      await upkeepBalanceMonitor.connect(owner).pause()
      expect(await upkeepBalanceMonitor.paused()).to.be.true
      await upkeepBalanceMonitor.connect(owner).unpause()
      expect(await upkeepBalanceMonitor.paused()).to.be.false
    })

    it('cannot be called by a non-owner', async () => {
      await expect(
        upkeepBalanceMonitor.connect(stranger).pause(),
      ).to.be.revertedWith('Only callable by owner')
      await upkeepBalanceMonitor.connect(owner).pause()
      await expect(
        upkeepBalanceMonitor.connect(stranger).unpause(),
      ).to.be.revertedWith('Only callable by owner')
    })
  })

  describe('checkUpkeep() / getUnderfundedUpkeeps()', () => {
    it('should find the underfunded upkeeps', async () => {
      let [upkeepIDs, registries, topUpAmounts] =
        await upkeepBalanceMonitor.getUnderfundedUpkeeps()
      expect(upkeepIDs.length).to.equal(0)
      expect(registries.length).to.equal(0)
      expect(topUpAmounts.length).to.equal(0)
      let [upkeepNeeded, performData] =
        await upkeepBalanceMonitor.checkUpkeep('0x')
      expect(upkeepNeeded).to.be.false
      expect(performData).to.equal('0x')
      // update the balance for some upkeeps
      await registry.mock.getBalance.withArgs(2).returns(120)
      await registry.mock.getBalance.withArgs(4).returns(15)
      await registry.mock.getBalance.withArgs(5).returns(0)
      ;[upkeepIDs, registries, topUpAmounts] =
        await upkeepBalanceMonitor.getUnderfundedUpkeeps()
      expect(upkeepIDs.map((v) => v.toNumber())).to.deep.equal([2, 4, 5])
      expect(registries).to.deep.equal([
        registry.address,
        registry.address,
        registry.address,
      ])
      expect(topUpAmounts.map((v) => v.toNumber())).to.deep.equal([
        180, 285, 300,
      ])
      ;[upkeepNeeded, performData] =
        await upkeepBalanceMonitor.checkUpkeep('0x')
      expect(upkeepNeeded).to.be.true
      expect(performData).to.equal(
        ethers.utils.defaultAbiCoder.encode(
          ['uint256[]', 'address[]', 'uint256[]'],
          [
            [2, 4, 5],
            [registry.address, registry.address, registry.address],
            [180, 285, 300],
          ],
        ),
      )
      // update all to need funding
      for (let i = 0; i < 9; i++) {
        await registry.mock.getBalance.withArgs(i).returns(0)
      }
      for (let i = 9; i < 12; i++) {
        await registry2.mock.getBalance.withArgs(i).returns(0)
      }
      // only the max batch size are included in the list
      ;[upkeepIDs, registries, topUpAmounts] =
        await upkeepBalanceMonitor.getUnderfundedUpkeeps()
      expect(upkeepIDs.length).to.equal(10)
      expect(topUpAmounts.length).to.equal(10)
      expect(upkeepIDs.map((v) => v.toNumber())).to.deep.equal([
        0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
      ])
      expect(registries).to.deep.equal([
        ...Array(9).fill(registry.address),
        registry2.address,
      ])
      expect(topUpAmounts.map((v) => v.toNumber())).to.deep.equal([
        ...Array(10).fill(300),
      ])
      // update the balance for some upkeeps
      await registry.mock.getBalance.withArgs(0).returns(300)
      await registry.mock.getBalance.withArgs(5).returns(300)
      ;[upkeepIDs, registries, topUpAmounts] =
        await upkeepBalanceMonitor.getUnderfundedUpkeeps()
      expect(upkeepIDs.length).to.equal(10)
      expect(topUpAmounts.length).to.equal(10)
      expect(upkeepIDs.map((v) => v.toNumber())).to.deep.equal([
        1, 2, 3, 4, 6, 7, 8, 9, 10, 11,
      ])
      expect(registries).to.deep.equal([
        ...Array(7).fill(registry.address),
        ...Array(3).fill(registry2.address),
      ])
      expect(topUpAmounts.map((v) => v.toNumber())).to.deep.equal([
        ...Array(10).fill(300),
      ])
    })
  })

  describe('topUp()', () => {
    beforeEach(async () => {
      await registry.mock.onTokenTransfer
        .withArgs(
          upkeepBalanceMonitor.address,
          100,
          ethers.utils.defaultAbiCoder.encode(['uint256'], [1]),
        )
        .returns()
      await registry.mock.onTokenTransfer
        .withArgs(
          upkeepBalanceMonitor.address,
          50,
          ethers.utils.defaultAbiCoder.encode(['uint256'], [7]),
        )
        .returns()
    })

    it('cannot be called by a non-owner', async () => {
      await expect(
        upkeepBalanceMonitor.connect(stranger).topUp([], [], []),
      ).to.be.revertedWithCustomError(
        upkeepBalanceMonitor,
        'OnlyForwarderOrOwner',
      )
    })

    it('should revert if the contract is paused', async () => {
      await upkeepBalanceMonitor.connect(owner).pause()
      await expect(
        upkeepBalanceMonitor.connect(owner).topUp([], [], []),
      ).to.be.revertedWith('Pausable: paused')
    })

    it('tops up the upkeeps by the amounts provided', async () => {
      const initialBalance = await linkToken.balanceOf(registry.address)
      const tx = await upkeepBalanceMonitor
        .connect(owner)
        .topUp([1, 7], [registry.address, registry.address], [100, 50])
      const finalBalance = await linkToken.balanceOf(registry.address)
      expect(finalBalance).to.equal(initialBalance.add(150))
      await expect(tx)
        .to.emit(upkeepBalanceMonitor, 'TopUpSucceeded')
        .withArgs(1, 100)
      await expect(tx)
        .to.emit(upkeepBalanceMonitor, 'TopUpSucceeded')
        .withArgs(7, 50)
    })

    it('does not abort if one top-up fails', async () => {
      const initialBalance = await linkToken.balanceOf(registry.address)
      const tx = await upkeepBalanceMonitor
        .connect(owner)
        .topUp(
          [1, 7, 100],
          [registry.address, registry.address, registry.address],
          [100, 50, 100],
        )
      const finalBalance = await linkToken.balanceOf(registry.address)
      expect(finalBalance).to.equal(initialBalance.add(150))
      await expect(tx)
        .to.emit(upkeepBalanceMonitor, 'TopUpSucceeded')
        .withArgs(1, 100)
      await expect(tx)
        .to.emit(upkeepBalanceMonitor, 'TopUpSucceeded')
        .withArgs(7, 50)
      await expect(tx)
        .to.emit(upkeepBalanceMonitor, 'TopUpFailed')
        .withArgs(100)
    })
  })

  describe('checkUpkeep() / performUpkeep()', () => {
    it('works round-trip', async () => {
      await registry.mock.getBalance.withArgs(1).returns(100) // needs 200
      await registry.mock.getBalance.withArgs(7).returns(0) // needs 300
      await registry.mock.onTokenTransfer
        .withArgs(
          upkeepBalanceMonitor.address,
          200,
          ethers.utils.defaultAbiCoder.encode(['uint256'], [1]),
        )
        .returns()
      await registry.mock.onTokenTransfer
        .withArgs(
          upkeepBalanceMonitor.address,
          300,
          ethers.utils.defaultAbiCoder.encode(['uint256'], [7]),
        )
        .returns()
      const [upkeepNeeded, performData] =
        await upkeepBalanceMonitor.checkUpkeep('0x')
      expect(upkeepNeeded).to.be.true
      const initialBalance = await linkToken.balanceOf(registry.address)
      await upkeepBalanceMonitor.connect(owner).performUpkeep(performData)
      const finalBalance = await linkToken.balanceOf(registry.address)
      expect(finalBalance).to.equal(initialBalance.add(500))
    })
  })
})
