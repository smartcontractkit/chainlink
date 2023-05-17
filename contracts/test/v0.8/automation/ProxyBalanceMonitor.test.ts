import { ethers } from 'hardhat'
import chai, { assert, expect } from 'chai'
import type { SignerWithAddress } from '@nomiclabs/hardhat-ethers/signers'
import { loadFixture } from '@nomicfoundation/hardhat-network-helpers'
import * as h from '../../test-helpers/helpers'
import { IAggregatorProxy__factory as IAggregatorProxyFactory } from '../../../typechain/factories/IAggregatorProxy__factory'
import { ILinkAvailable__factory as ILinkAvailableFactory } from '../../../typechain/factories/ILinkAvailable__factory'
import { ProxyBalanceMonitor, LinkToken } from '../../../typechain'
import { BigNumber } from 'ethers'
import { mineBlock } from '../../test-helpers/helpers'
import deepEqualInAnyOrder from 'deep-equal-in-any-order'
import {
  deployMockContract,
  MockContract,
} from '@ethereum-waffle/mock-contract'
chai.use(deepEqualInAnyOrder)

//////////////////////////////// GAS USAGE LIMITS - CHANGE WITH CAUTION //////////////////////////
//                                                                                              //
// we try to keep gas usage under this amount (max is 5M)                                       //
const TARGET_PERFORM_GAS_LIMIT = 2_000_000
// we try to keep gas usage under this amount (max is 5M) the test is not a perfectly accurate  //
// measurement of gas usage because it relies on mocks which may do fewer storage reads         //
// therefore, we keep a healthy margin to avoid running over the limit!                         //
const TARGET_CHECK_GAS_LIMIT = 3_500_000
//                                                                                              //
//////////////////////////////////////////////////////////////////////////////////////////////////

const OWNABLE_ERR = 'Only callable by owner'
const INVALID_WATCHLIST_ERR = `InvalidWatchList()`
const PAUSED_ERR = 'Pausable: paused'

const zeroLINK = ethers.utils.parseEther('0')
const oneLINK = ethers.utils.parseEther('1')
const twoLINK = ethers.utils.parseEther('2')
const fiveLINK = ethers.utils.parseEther('5')
const sixLINK = ethers.utils.parseEther('6')
const tenLINK = ethers.utils.parseEther('10')
const oneHundredLINK = ethers.utils.parseEther('100')

const randAddr = () => ethers.Wallet.createRandom().address

let pm: ProxyBalanceMonitor
let lt: LinkToken
let owner: SignerWithAddress
let stranger: SignerWithAddress
let keeperRegistry: SignerWithAddress
let proxy1: MockContract
let proxy2: MockContract
let proxy3: MockContract
let proxy4: MockContract // leave this proxy / aggregator unconfigured for topUp() testing
let aggregator1: MockContract
let aggregator2: MockContract
let aggregator3: MockContract
let aggregator4: MockContract // leave this proxy / aggregator unconfigured for topUp() testing

let proxies: string[]

async function assertAggregatorBalances(
  balance1: BigNumber,
  balance2: BigNumber,
  balance3: BigNumber,
) {
  await h.assertLinkTokenBalance(lt, aggregator1.address, balance1, 'address 1')
  await h.assertLinkTokenBalance(lt, aggregator2.address, balance2, 'address 2')
  await h.assertLinkTokenBalance(lt, aggregator3.address, balance3, 'address 3')
}

const setup = async () => {
  const accounts = await ethers.getSigners()
  owner = accounts[0]
  stranger = accounts[1]
  keeperRegistry = accounts[2]

  proxy1 = await deployMockContract(owner, IAggregatorProxyFactory.abi)
  proxy2 = await deployMockContract(owner, IAggregatorProxyFactory.abi)
  proxy3 = await deployMockContract(owner, IAggregatorProxyFactory.abi)
  proxy4 = await deployMockContract(owner, IAggregatorProxyFactory.abi)
  aggregator1 = await deployMockContract(owner, ILinkAvailableFactory.abi)
  aggregator2 = await deployMockContract(owner, ILinkAvailableFactory.abi)
  aggregator3 = await deployMockContract(owner, ILinkAvailableFactory.abi)
  aggregator4 = await deployMockContract(owner, ILinkAvailableFactory.abi)

  await proxy1.deployed()
  await proxy2.deployed()
  await proxy3.deployed()
  await proxy4.deployed()
  await aggregator1.deployed()
  await aggregator2.deployed()
  await aggregator3.deployed()
  await aggregator4.deployed()

  proxies = [proxy1.address, proxy2.address, proxy3.address]

  await proxy1.mock.aggregator.returns(aggregator1.address)
  await proxy2.mock.aggregator.returns(aggregator2.address)
  await proxy3.mock.aggregator.returns(aggregator3.address)

  await aggregator1.mock.linkAvailableForPayment.returns(0)
  await aggregator2.mock.linkAvailableForPayment.returns(0)
  await aggregator3.mock.linkAvailableForPayment.returns(0)

  await aggregator1.mock.transmitters.returns([randAddr()])
  await aggregator2.mock.transmitters.returns([randAddr()])
  await aggregator3.mock.transmitters.returns([randAddr()])

  const pmFactory = await ethers.getContractFactory(
    'ProxyBalanceMonitor',
    owner,
  )
  const ltFactory = await ethers.getContractFactory('LinkToken', owner)

  lt = await ltFactory.deploy()
  pm = await pmFactory.deploy(lt.address, oneLINK, twoLINK)
  await pm.deployed()

  for (let i = 1; i <= 4; i++) {
    const recipient = await accounts[i].getAddress()
    await lt.connect(owner).transfer(recipient, oneHundredLINK)
  }

  const setTx = await pm.connect(owner).setWatchList(proxies)
  await setTx.wait()
}

describe('ProxyBalanceMonitor', () => {
  beforeEach(async () => {
    await loadFixture(setup)
  })

  describe('add funds', () => {
    it('Should allow anyone to add funds', async () => {
      await lt.transfer(pm.address, oneLINK)
      await lt.connect(stranger).transfer(pm.address, oneLINK)
    })
  })

  describe('setTopUpAmount()', () => {
    it('configures the top-up amount', async () => {
      await pm.connect(owner).setTopUpAmount(100)
      assert.equal((await pm.getTopUpAmount()).toNumber(), 100)
    })

    it('configuresis only callable by the owner', async () => {
      await expect(pm.connect(stranger).setTopUpAmount(100)).to.be.reverted
    })
  })

  describe('setMinBalance()', () => {
    it('configures the min balance', async () => {
      await pm.connect(owner).setMinBalance(100)
      assert.equal((await pm.getMinBalance()).toNumber(), 100)
    })

    it('is only callable by the owner', async () => {
      await expect(pm.connect(stranger).setMinBalance(100)).to.be.reverted
    })
  })

  describe('withdraw()', () => {
    beforeEach(async () => {
      const tx = await lt.connect(owner).transfer(pm.address, oneLINK)
      await tx.wait()
    })

    it('Should allow the owner to withdraw', async () => {
      const beforeBalance = await lt.balanceOf(owner.address)
      const tx = await pm.connect(owner).withdraw(oneLINK, owner.address)
      await tx.wait()
      const afterBalance = await lt.balanceOf(owner.address)
      assert.isTrue(
        afterBalance.gt(beforeBalance),
        'balance did not increase after withdraw',
      )
    })

    it('Should emit an event', async () => {
      const tx = await pm.connect(owner).withdraw(oneLINK, owner.address)
      await expect(tx)
        .to.emit(pm, 'FundsWithdrawn')
        .withArgs(oneLINK, owner.address)
    })

    it('Should allow the owner to withdraw to anyone', async () => {
      const beforeBalance = await lt.balanceOf(stranger.address)
      const tx = await pm.connect(owner).withdraw(oneLINK, stranger.address)
      await tx.wait()
      const afterBalance = await lt.balanceOf(stranger.address)
      assert.isTrue(
        beforeBalance.add(oneLINK).eq(afterBalance),
        'balance did not increase after withdraw',
      )
    })

    it('Should not allow strangers to withdraw', async () => {
      const tx = pm.connect(stranger).withdraw(oneLINK, owner.address)
      await expect(tx).to.be.revertedWith(OWNABLE_ERR)
    })
  })

  describe('pause() / unpause()', () => {
    it('Should allow owner to pause / unpause', async () => {
      const pauseTx = await pm.connect(owner).pause()
      await pauseTx.wait()
      const unpauseTx = await pm.connect(owner).unpause()
      await unpauseTx.wait()
    })

    it('Should not allow strangers to pause / unpause', async () => {
      const pauseTxStranger = pm.connect(stranger).pause()
      await expect(pauseTxStranger).to.be.revertedWith(OWNABLE_ERR)
      const pauseTxOwner = await pm.connect(owner).pause()
      await pauseTxOwner.wait()
      const unpauseTxStranger = pm.connect(stranger).unpause()
      await expect(unpauseTxStranger).to.be.revertedWith(OWNABLE_ERR)
    })
  })

  describe('setWatchList() / addToWatchList() / getWatchList()', () => {
    const watchAddress1 = randAddr()
    const watchAddress2 = randAddr()
    const watchAddress3 = randAddr()

    beforeEach(async () => {
      // reset watchlist to empty before running these tests
      await pm.connect(owner).setWatchList([])
      let watchList = await pm.getWatchList()
      assert.deepEqual(watchList, [])
    })

    it('Should allow owner to set the watchlist', async () => {
      // add first watchlist
      let tx = await pm.connect(owner).setWatchList([watchAddress1])
      let watchList = await pm.getWatchList()
      assert.deepEqual(watchList, [watchAddress1])
      // add more to watchlist
      tx = await pm
        .connect(owner)
        .setWatchList([watchAddress1, watchAddress2, watchAddress3])
      await tx.wait()
      watchList = await pm.getWatchList()
      assert.deepEqual(watchList, [watchAddress1, watchAddress2, watchAddress3])
      // remove some from watchlist
      tx = await pm.connect(owner).setWatchList([watchAddress3, watchAddress1])
      await tx.wait()
      watchList = await pm.getWatchList()
      assert.deepEqual(watchList, [watchAddress3, watchAddress1])
      // add some to watchlist
      tx = await pm.connect(owner).addToWatchList([watchAddress2])
      await tx.wait()
      watchList = await pm.getWatchList()
      assert.deepEqual(watchList, [watchAddress3, watchAddress1, watchAddress2])
    })

    it('Should not allow duplicates in the watchlist', async () => {
      const errMsg = `DuplicateAddress("${watchAddress1}")`
      let tx = pm
        .connect(owner)
        .setWatchList([watchAddress1, watchAddress2, watchAddress1])
      await expect(tx).to.be.revertedWith(errMsg)
      tx = pm
        .connect(owner)
        .addToWatchList([watchAddress1, watchAddress2, watchAddress1])
      await expect(tx).to.be.revertedWith(errMsg)
    })

    it('Should not allow strangers to set the watchlist', async () => {
      const setTxStranger = pm.connect(stranger).setWatchList([watchAddress1])
      await expect(setTxStranger).to.be.revertedWith(OWNABLE_ERR)
      const addTxStranger = pm.connect(stranger).addToWatchList([watchAddress1])
      await expect(addTxStranger).to.be.revertedWith(OWNABLE_ERR)
    })

    it('Should revert if any of the addresses are empty', async () => {
      let tx = pm
        .connect(owner)
        .setWatchList([watchAddress1, ethers.constants.AddressZero])
      await expect(tx).to.be.revertedWith(INVALID_WATCHLIST_ERR)
      tx = pm
        .connect(owner)
        .addToWatchList([watchAddress1, ethers.constants.AddressZero])
      await expect(tx).to.be.revertedWith(INVALID_WATCHLIST_ERR)
    })
  })

  describe('checkUpkeep() / sampleUnderfundedAddresses()', () => {
    it('Should return list of address that are underfunded', async () => {
      const fundTx = await lt.connect(owner).transfer(
        pm.address,
        sixLINK, // needs 6 total
      )
      await fundTx.wait()
      const [should, payload] = await pm.checkUpkeep('0x')
      assert.isTrue(should)
      let [addresses] = ethers.utils.defaultAbiCoder.decode(
        ['address[]'],
        payload,
      )
      expect(addresses).to.deep.equalInAnyOrder(proxies)
      // checkUpkeep payload should match sampleUnderfundedAddresses()
      addresses = await pm.sampleUnderfundedAddresses()
      expect(addresses).to.deep.equalInAnyOrder(proxies)
    })

    it('Should return some results even if contract cannot fund all eligible targets', async () => {
      const fundTx = await lt.connect(owner).transfer(
        pm.address,
        fiveLINK, // needs 6 total
      )
      await fundTx.wait()
      const [should, payload] = await pm.checkUpkeep('0x')
      assert.isTrue(should)
      let [addresses] = ethers.utils.defaultAbiCoder.decode(
        ['address[]'],
        payload,
      )
      assert.equal(addresses.length, 2)
      assert.notEqual(addresses[0], addresses[1])
      assert(proxies.includes(addresses[0]))
      assert(proxies.includes(addresses[1]))
      // underfunded sample should still match list
      addresses = await pm.sampleUnderfundedAddresses()
      expect(addresses).to.deep.equalInAnyOrder(proxies)
    })

    it('Should omit aggregators that have sufficient funding', async () => {
      let addresses = await pm.sampleUnderfundedAddresses()
      expect(addresses).to.deep.equalInAnyOrder(proxies)
      await aggregator2.mock.linkAvailableForPayment.returns(tenLINK)
      addresses = await pm.sampleUnderfundedAddresses()
      expect(addresses).to.deep.equalInAnyOrder([
        proxy1.address,
        proxy3.address,
      ])
      await aggregator1.mock.linkAvailableForPayment.returns(tenLINK)
      addresses = await pm.sampleUnderfundedAddresses()
      expect(addresses).to.deep.equalInAnyOrder([proxy3.address])
      await aggregator3.mock.linkAvailableForPayment.returns(tenLINK)
      addresses = await pm.sampleUnderfundedAddresses()
      expect(addresses).to.deep.equalInAnyOrder([])
    })

    it('Should revert when paused', async () => {
      const tx = await pm.connect(owner).pause()
      await tx.wait()
      const ethCall = pm.checkUpkeep('0x')
      await expect(ethCall).to.be.revertedWith(PAUSED_ERR)
    })

    context('with a large set of proxies', async () => {
      // in this test, we cheat a little bit and point each proxy to the same aggregator,
      // which helps cut down on test time
      let MAX_PERFORM: number
      let MAX_CHECK: number
      let proxyAddresses: string[]
      let aggregators: MockContract[]

      beforeEach(async () => {
        MAX_PERFORM = (await pm.MAX_PERFORM()).toNumber()
        MAX_CHECK = (await pm.MAX_CHECK()).toNumber()
        proxyAddresses = []
        aggregators = []
        const numAggregators = MAX_CHECK + 50
        for (let idx = 0; idx < numAggregators; idx++) {
          const proxy = await deployMockContract(
            owner,
            IAggregatorProxyFactory.abi,
          )
          const aggregator = await deployMockContract(
            owner,
            ILinkAvailableFactory.abi,
          )
          await proxy.mock.aggregator.returns(aggregator.address)
          await aggregator.mock.linkAvailableForPayment.returns(0)
          await aggregator.mock.transmitters.returns([randAddr()])
          proxyAddresses.push(proxy.address)
          aggregators.push(aggregator)
        }
        await pm.setWatchList(proxyAddresses)
        expect(await pm.getWatchList()).to.deep.equalInAnyOrder(proxyAddresses)
      })

      it('Should not include more than MAX_PERFORM addresses', async () => {
        const addresses = await pm.sampleUnderfundedAddresses()
        assert.equal(addresses.length, MAX_PERFORM)
      })

      it('Should sample from the list of addresses pseudorandomly', async () => {
        const firstAddress: string[] = []
        for (let idx = 0; idx < 10; idx++) {
          const addresses = await pm.sampleUnderfundedAddresses()
          assert.equal(addresses.length, MAX_PERFORM)
          assert.equal(
            new Set(addresses).size,
            MAX_PERFORM,
            'duplicate address found',
          )
          firstAddress.push(addresses[0])
          await mineBlock(ethers.provider)
        }
        assert(
          new Set(firstAddress).size > 1,
          'sample did not shuffle starting index',
        )
      })

      it('Can check MAX_CHECK upkeeps within the allotted gas limit', async () => {
        for (const aggregator of aggregators) {
          // here we make no aggregators eligible for funding, requiring the function to
          // traverse the whole list
          await aggregator.mock.linkAvailableForPayment.returns(tenLINK)
        }
        await pm.checkUpkeep('0x', { gasLimit: TARGET_CHECK_GAS_LIMIT })
      })
    })
  })

  describe('performUpkeep()', () => {
    let validPayload: string

    beforeEach(async () => {
      validPayload = ethers.utils.defaultAbiCoder.encode(
        ['address[]'],
        [proxies],
      )
      await pm.connect(owner).setWatchList(proxies)
    })

    it('Should revert when paused', async () => {
      await pm.connect(owner).pause()
      const performTx = pm.connect(keeperRegistry).performUpkeep(validPayload)
      await expect(performTx).to.be.revertedWith(PAUSED_ERR)
    })

    it('Should fund the appropriate addresses', async () => {
      await lt.connect(owner).transfer(pm.address, tenLINK)
      await assertAggregatorBalances(zeroLINK, zeroLINK, zeroLINK)
      const performTx = await pm
        .connect(keeperRegistry)
        .performUpkeep(validPayload, { gasLimit: 2_500_000 })
      await performTx.wait()
      await assertAggregatorBalances(twoLINK, twoLINK, twoLINK)
    })

    it('Can handle MAX_PERFORM proxies within gas limit', async () => {
      // add MAX_PERFORM number of proxies
      const MAX_PERFORM = (await pm.MAX_PERFORM()).toNumber()
      const proxyAddresses = []
      for (let idx = 0; idx < MAX_PERFORM; idx++) {
        const proxy = await deployMockContract(
          owner,
          IAggregatorProxyFactory.abi,
        )
        const aggregator = await deployMockContract(
          owner,
          ILinkAvailableFactory.abi,
        )
        await proxy.mock.aggregator.returns(aggregator.address)
        await aggregator.mock.linkAvailableForPayment.returns(0)
        await aggregator.mock.transmitters.returns([randAddr()])
        proxyAddresses.push(proxy.address)
      }
      await pm.setWatchList(proxyAddresses)
      expect(await pm.getWatchList()).to.deep.equalInAnyOrder(proxyAddresses)
      // add funds
      const fundsNeeded = (await pm.getTopUpAmount()).mul(MAX_PERFORM)
      await lt.connect(owner).transfer(pm.address, fundsNeeded)
      // encode payload
      const payload = ethers.utils.defaultAbiCoder.encode(
        ['address[]'],
        [proxyAddresses],
      )
      // do the thing
      await pm
        .connect(keeperRegistry)
        .performUpkeep(payload, { gasLimit: TARGET_PERFORM_GAS_LIMIT })
    })
  })

  describe('topUp()', () => {
    context('when not paused', () => {
      it('Should be callable by anyone', async () => {
        const users = [owner, keeperRegistry, stranger]
        for (let idx = 0; idx < users.length; idx++) {
          const user = users[idx]
          await pm.connect(user).topUp([])
        }
      })
    })

    context('when paused', () => {
      it('Should be callable by no one', async () => {
        await pm.connect(owner).pause()
        const users = [owner, keeperRegistry, stranger]
        for (let idx = 0; idx < users.length; idx++) {
          const user = users[idx]
          const tx = pm.connect(user).topUp([])
          await expect(tx).to.be.revertedWith(PAUSED_ERR)
        }
      })
    })

    context('when fully funded', () => {
      beforeEach(async () => {
        await lt.connect(owner).transfer(pm.address, tenLINK)
        await assertAggregatorBalances(zeroLINK, zeroLINK, zeroLINK)
      })

      it('Should fund the appropriate addresses', async () => {
        const tx = await pm.connect(keeperRegistry).topUp(proxies)
        await assertAggregatorBalances(twoLINK, twoLINK, twoLINK)
        await expect(tx).to.emit(pm, 'TopUpSucceeded').withArgs(proxy1.address)
        await expect(tx).to.emit(pm, 'TopUpSucceeded').withArgs(proxy2.address)
        await expect(tx).to.emit(pm, 'TopUpSucceeded').withArgs(proxy3.address)
      })

      it('Should only fund the addresses provided', async () => {
        await pm.connect(keeperRegistry).topUp([proxy1.address, proxy3.address])
        await assertAggregatorBalances(twoLINK, zeroLINK, twoLINK)
      })

      it('Should skip un-approved addresses', async () => {
        await pm.connect(owner).setWatchList([proxy1.address, proxy2.address])
        const tx = await pm
          .connect(keeperRegistry)
          .topUp([proxy1.address, proxy2.address, proxy3.address])
        await assertAggregatorBalances(twoLINK, twoLINK, zeroLINK)
        await expect(tx).to.emit(pm, 'TopUpSucceeded').withArgs(proxy1.address)
        await expect(tx).to.emit(pm, 'TopUpSucceeded').withArgs(proxy2.address)
        await expect(tx).to.emit(pm, 'TopUpBlocked').withArgs(proxy3.address)
      })

      it('Should skip an address if the proxy is invalid', async () => {
        await pm.connect(owner).setWatchList([proxy1.address, proxy4.address])
        const tx = await pm
          .connect(keeperRegistry)
          .topUp([proxy1.address, proxy4.address])
        await expect(tx).to.emit(pm, 'TopUpSucceeded').withArgs(proxy1.address)
        await expect(tx).to.emit(pm, 'TopUpBlocked').withArgs(proxy4.address)
      })

      it('Should skip an address if the aggregator is invalid', async () => {
        await proxy4.mock.aggregator.returns(aggregator4.address)
        await pm.connect(owner).setWatchList([proxy1.address, proxy4.address])
        const tx = await pm
          .connect(keeperRegistry)
          .topUp([proxy1.address, proxy4.address])
        await expect(tx).to.emit(pm, 'TopUpSucceeded').withArgs(proxy1.address)
        await expect(tx).to.emit(pm, 'TopUpBlocked').withArgs(proxy4.address)
      })

      it('Should skip an address if the aggregator has no transmitters', async () => {
        await proxy4.mock.aggregator.returns(aggregator4.address)
        await aggregator4.mock.linkAvailableForPayment.returns(0)
        await aggregator4.mock.transmitters.returns([])
        await pm.connect(owner).setWatchList([proxy1.address, proxy4.address])
        const tx = await pm
          .connect(keeperRegistry)
          .topUp([proxy1.address, proxy4.address])
        await expect(tx).to.emit(pm, 'TopUpSucceeded').withArgs(proxy1.address)
        await expect(tx).to.emit(pm, 'TopUpBlocked').withArgs(proxy4.address)
      })

      it('Should skip an address if the aggregator has sufficient funding', async () => {
        await proxy4.mock.aggregator.returns(aggregator4.address)
        await aggregator4.mock.linkAvailableForPayment.returns(tenLINK)
        await aggregator4.mock.transmitters.returns([randAddr()])
        await pm.connect(owner).setWatchList([proxy1.address, proxy4.address])
        const tx = await pm
          .connect(keeperRegistry)
          .topUp([proxy1.address, proxy4.address])
        await expect(tx).to.emit(pm, 'TopUpSucceeded').withArgs(proxy1.address)
        await expect(tx).to.emit(pm, 'TopUpBlocked').withArgs(proxy4.address)
      })
    })

    context('when partially funded', () => {
      it('Should fund as many addresses as possible', async () => {
        await lt.connect(owner).transfer(
          pm.address,
          fiveLINK, // only enough LINK to fund 2 addresses
        )
        const tx = await pm.connect(keeperRegistry).topUp(proxies)
        await assertAggregatorBalances(twoLINK, twoLINK, zeroLINK)
        await expect(tx).to.emit(pm, 'TopUpSucceeded').withArgs(proxy1.address)
        await expect(tx).to.emit(pm, 'TopUpSucceeded').withArgs(proxy2.address)
      })
    })
  })
})
