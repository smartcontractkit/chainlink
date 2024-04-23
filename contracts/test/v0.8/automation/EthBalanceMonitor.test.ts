import { ethers } from 'hardhat'
import { assert, expect } from 'chai'
import type { SignerWithAddress } from '@nomiclabs/hardhat-ethers/signers'
import { EthBalanceMonitorExposed } from '../../../typechain/EthBalanceMonitorExposed'
import { ReceiveReverter } from '../../../typechain/ReceiveReverter'
import { ReceiveEmitter } from '../../../typechain/ReceiveEmitter'
import { ReceiveFallbackEmitter } from '../../../typechain/ReceiveFallbackEmitter'
import { BigNumber } from 'ethers'
import * as h from '../../test-helpers/helpers'

const OWNABLE_ERR = 'Only callable by owner'
const INVALID_WATCHLIST_ERR = `InvalidWatchList`
const PAUSED_ERR = 'Pausable: paused'
const ONLY_KEEPER_ERR = `OnlyKeeperRegistry`

const zeroEth = ethers.utils.parseEther('0')
const oneEth = ethers.utils.parseEther('1')
const twoEth = ethers.utils.parseEther('2')
const threeEth = ethers.utils.parseEther('3')
const fiveEth = ethers.utils.parseEther('5')
const sixEth = ethers.utils.parseEther('6')
const tenEth = ethers.utils.parseEther('10')

const watchAddress1 = ethers.Wallet.createRandom().address
const watchAddress2 = ethers.Wallet.createRandom().address
const watchAddress3 = ethers.Wallet.createRandom().address
const watchAddress4 = ethers.Wallet.createRandom().address
let watchAddress5: string
let watchAddress6: string

async function assertWatchlistBalances(
  balance1: number,
  balance2: number,
  balance3: number,
  balance4: number,
  balance5: number,
  balance6: number,
) {
  const toEth = (n: number) => ethers.utils.parseUnits(n.toString(), 'ether')
  await h.assertBalance(watchAddress1, toEth(balance1), 'address 1')
  await h.assertBalance(watchAddress2, toEth(balance2), 'address 2')
  await h.assertBalance(watchAddress3, toEth(balance3), 'address 3')
  await h.assertBalance(watchAddress4, toEth(balance4), 'address 4')
  await h.assertBalance(watchAddress5, toEth(balance5), 'address 5')
  await h.assertBalance(watchAddress6, toEth(balance6), 'address 6')
}

let bm: EthBalanceMonitorExposed
let receiveReverter: ReceiveReverter
let receiveEmitter: ReceiveEmitter
let receiveFallbackEmitter: ReceiveFallbackEmitter
let owner: SignerWithAddress
let stranger: SignerWithAddress
let keeperRegistry: SignerWithAddress

describe('EthBalanceMonitor', () => {
  beforeEach(async () => {
    const accounts = await ethers.getSigners()
    owner = accounts[0]
    stranger = accounts[1]
    keeperRegistry = accounts[2]
    watchAddress5 = accounts[3].address
    watchAddress6 = accounts[4].address

    const bmFactory = await ethers.getContractFactory(
      'EthBalanceMonitorExposed',
      owner,
    )
    const rrFactory = await ethers.getContractFactory('ReceiveReverter', owner)
    const reFactory = await ethers.getContractFactory('ReceiveEmitter', owner)
    const rfeFactory = await ethers.getContractFactory(
      'ReceiveFallbackEmitter',
      owner,
    )

    bm = await bmFactory.deploy(keeperRegistry.address, 0)
    receiveReverter = await rrFactory.deploy()
    receiveEmitter = await reFactory.deploy()
    receiveFallbackEmitter = await rfeFactory.deploy()
    await Promise.all([
      bm.deployed(),
      receiveReverter.deployed(),
      receiveEmitter.deployed(),
      receiveFallbackEmitter.deployed(),
    ])
  })

  afterEach(async () => {
    await h.reset()
  })

  describe('receive()', () => {
    it('Should allow anyone to add funds', async () => {
      await owner.sendTransaction({
        to: bm.address,
        value: oneEth,
      })
      await stranger.sendTransaction({
        to: bm.address,
        value: oneEth,
      })
    })

    it('Should emit an event', async () => {
      await owner.sendTransaction({
        to: bm.address,
        value: oneEth,
      })
      const tx = stranger.sendTransaction({
        to: bm.address,
        value: oneEth,
      })
      await expect(tx)
        .to.emit(bm, 'FundsAdded')
        .withArgs(oneEth, twoEth, stranger.address)
    })
  })

  describe('withdraw()', () => {
    beforeEach(async () => {
      const tx = await owner.sendTransaction({
        to: bm.address,
        value: oneEth,
      })
      await tx.wait()
    })

    it('Should allow the owner to withdraw', async () => {
      const beforeBalance = await owner.getBalance()
      const tx = await bm.connect(owner).withdraw(oneEth, owner.address)
      await tx.wait()
      const afterBalance = await owner.getBalance()
      assert.isTrue(
        afterBalance.gt(beforeBalance),
        'balance did not increase after withdraw',
      )
    })

    it('Should emit an event', async () => {
      const tx = await bm.connect(owner).withdraw(oneEth, owner.address)
      await expect(tx)
        .to.emit(bm, 'FundsWithdrawn')
        .withArgs(oneEth, owner.address)
    })

    it('Should allow the owner to withdraw to anyone', async () => {
      const beforeBalance = await stranger.getBalance()
      const tx = await bm.connect(owner).withdraw(oneEth, stranger.address)
      await tx.wait()
      const afterBalance = await stranger.getBalance()
      assert.isTrue(
        beforeBalance.add(oneEth).eq(afterBalance),
        'balance did not increase after withdraw',
      )
    })

    it('Should not allow strangers to withdraw', async () => {
      const tx = bm.connect(stranger).withdraw(oneEth, owner.address)
      await expect(tx).to.be.revertedWith(OWNABLE_ERR)
    })
  })

  describe('pause() / unpause()', () => {
    it('Should allow owner to pause / unpause', async () => {
      const pauseTx = await bm.connect(owner).pause()
      await pauseTx.wait()
      const unpauseTx = await bm.connect(owner).unpause()
      await unpauseTx.wait()
    })

    it('Should not allow strangers to pause / unpause', async () => {
      const pauseTxStranger = bm.connect(stranger).pause()
      await expect(pauseTxStranger).to.be.revertedWith(OWNABLE_ERR)
      const pauseTxOwner = await bm.connect(owner).pause()
      await pauseTxOwner.wait()
      const unpauseTxStranger = bm.connect(stranger).unpause()
      await expect(unpauseTxStranger).to.be.revertedWith(OWNABLE_ERR)
    })
  })

  describe('setWatchList() / getWatchList() / getAccountInfo()', () => {
    it('Should allow owner to set the watchlist', async () => {
      // should start unactive
      assert.isFalse((await bm.getAccountInfo(watchAddress1)).isActive)
      // add first watchlist
      let setTx = await bm
        .connect(owner)
        .setWatchList([watchAddress1], [oneEth], [twoEth])
      await setTx.wait()
      let watchList = await bm.getWatchList()
      assert.deepEqual(watchList, [watchAddress1])
      const accountInfo = await bm.getAccountInfo(watchAddress1)
      assert.isTrue(accountInfo.isActive)
      expect(accountInfo.minBalanceWei).to.equal(oneEth)
      expect(accountInfo.topUpAmountWei).to.equal(twoEth)
      // add more to watchlist
      setTx = await bm
        .connect(owner)
        .setWatchList(
          [watchAddress1, watchAddress2, watchAddress3],
          [oneEth, twoEth, threeEth],
          [oneEth, twoEth, threeEth],
        )
      await setTx.wait()
      watchList = await bm.getWatchList()
      assert.deepEqual(watchList, [watchAddress1, watchAddress2, watchAddress3])
      let accountInfo1 = await bm.getAccountInfo(watchAddress1)
      let accountInfo2 = await bm.getAccountInfo(watchAddress2)
      let accountInfo3 = await bm.getAccountInfo(watchAddress3)
      expect(accountInfo1.isActive).to.be.true
      expect(accountInfo1.minBalanceWei).to.equal(oneEth)
      expect(accountInfo1.topUpAmountWei).to.equal(oneEth)
      expect(accountInfo2.isActive).to.be.true
      expect(accountInfo2.minBalanceWei).to.equal(twoEth)
      expect(accountInfo2.topUpAmountWei).to.equal(twoEth)
      expect(accountInfo3.isActive).to.be.true
      expect(accountInfo3.minBalanceWei).to.equal(threeEth)
      expect(accountInfo3.topUpAmountWei).to.equal(threeEth)
      // remove some from watchlist
      setTx = await bm
        .connect(owner)
        .setWatchList(
          [watchAddress3, watchAddress1],
          [threeEth, oneEth],
          [threeEth, oneEth],
        )
      await setTx.wait()
      watchList = await bm.getWatchList()
      assert.deepEqual(watchList, [watchAddress3, watchAddress1])
      accountInfo1 = await bm.getAccountInfo(watchAddress1)
      accountInfo2 = await bm.getAccountInfo(watchAddress2)
      accountInfo3 = await bm.getAccountInfo(watchAddress3)
      expect(accountInfo1.isActive).to.be.true
      expect(accountInfo2.isActive).to.be.false
      expect(accountInfo3.isActive).to.be.true
    })

    it('Should not allow duplicates in the watchlist', async () => {
      const errMsg = `DuplicateAddress`
      const setTx = bm
        .connect(owner)
        .setWatchList(
          [watchAddress1, watchAddress2, watchAddress1],
          [oneEth, twoEth, threeEth],
          [oneEth, twoEth, threeEth],
        )
      await expect(setTx)
        .to.be.revertedWithCustomError(bm, errMsg)
        .withArgs(watchAddress1)
    })

    it('Should not allow strangers to set the watchlist', async () => {
      const setTxStranger = bm
        .connect(stranger)
        .setWatchList([watchAddress1], [oneEth], [twoEth])
      await expect(setTxStranger).to.be.revertedWith(OWNABLE_ERR)
    })

    it('Should revert if the list lengths differ', async () => {
      let tx = bm.connect(owner).setWatchList([watchAddress1], [], [twoEth])
      await expect(tx).to.be.revertedWithCustomError(bm, INVALID_WATCHLIST_ERR)
      tx = bm.connect(owner).setWatchList([watchAddress1], [oneEth], [])
      await expect(tx).to.be.revertedWithCustomError(bm, INVALID_WATCHLIST_ERR)
      tx = bm.connect(owner).setWatchList([], [oneEth], [twoEth])
      await expect(tx).to.be.revertedWithCustomError(bm, INVALID_WATCHLIST_ERR)
    })

    it('Should revert if any of the addresses are empty', async () => {
      let tx = bm
        .connect(owner)
        .setWatchList(
          [watchAddress1, ethers.constants.AddressZero],
          [oneEth, oneEth],
          [twoEth, twoEth],
        )
      await expect(tx).to.be.revertedWithCustomError(bm, INVALID_WATCHLIST_ERR)
    })

    it('Should revert if any of the top up amounts are 0', async () => {
      const tx = bm
        .connect(owner)
        .setWatchList(
          [watchAddress1, watchAddress2],
          [oneEth, oneEth],
          [twoEth, zeroEth],
        )
      await expect(tx).to.be.revertedWithCustomError(bm, INVALID_WATCHLIST_ERR)
    })
  })

  describe('getKeeperRegistryAddress() / setKeeperRegistryAddress()', () => {
    const newAddress = ethers.Wallet.createRandom().address

    it('Should initialize with the registry address provided to the constructor', async () => {
      const address = await bm.getKeeperRegistryAddress()
      assert.equal(address, keeperRegistry.address)
    })

    it('Should allow the owner to set the registry address', async () => {
      const setTx = await bm.connect(owner).setKeeperRegistryAddress(newAddress)
      await setTx.wait()
      const address = await bm.getKeeperRegistryAddress()
      assert.equal(address, newAddress)
    })

    it('Should not allow strangers to set the registry address', async () => {
      const setTx = bm.connect(stranger).setKeeperRegistryAddress(newAddress)
      await expect(setTx).to.be.revertedWith(OWNABLE_ERR)
    })

    it('Should emit an event', async () => {
      const setTx = await bm.connect(owner).setKeeperRegistryAddress(newAddress)
      await expect(setTx)
        .to.emit(bm, 'KeeperRegistryAddressUpdated')
        .withArgs(keeperRegistry.address, newAddress)
    })
  })

  describe('getMinWaitPeriodSeconds / setMinWaitPeriodSeconds()', () => {
    const newWaitPeriod = BigNumber.from(1)

    it('Should initialize with the wait period provided to the constructor', async () => {
      const minWaitPeriod = await bm.getMinWaitPeriodSeconds()
      expect(minWaitPeriod).to.equal(0)
    })

    it('Should allow owner to set the wait period', async () => {
      const setTx = await bm
        .connect(owner)
        .setMinWaitPeriodSeconds(newWaitPeriod)
      await setTx.wait()
      const minWaitPeriod = await bm.getMinWaitPeriodSeconds()
      expect(minWaitPeriod).to.equal(newWaitPeriod)
    })

    it('Should not allow strangers to set the wait period', async () => {
      const setTx = bm.connect(stranger).setMinWaitPeriodSeconds(newWaitPeriod)
      await expect(setTx).to.be.revertedWith(OWNABLE_ERR)
    })

    it('Should emit an event', async () => {
      const setTx = await bm
        .connect(owner)
        .setMinWaitPeriodSeconds(newWaitPeriod)
      await expect(setTx)
        .to.emit(bm, 'MinWaitPeriodUpdated')
        .withArgs(0, newWaitPeriod)
    })
  })

  describe('checkUpkeep() / getUnderfundedAddresses()', () => {
    beforeEach(async () => {
      const setTx = await bm.connect(owner).setWatchList(
        [
          watchAddress1, // needs funds
          watchAddress5, // funded
          watchAddress2, // needs funds
          watchAddress6, // funded
          watchAddress3, // needs funds
        ],
        new Array(5).fill(oneEth),
        new Array(5).fill(twoEth),
      )
      await setTx.wait()
    })

    it('Should return list of address that are underfunded', async () => {
      const fundTx = await owner.sendTransaction({
        to: bm.address,
        value: sixEth, // needs 6 total
      })
      await fundTx.wait()
      const [should, payload] = await bm.checkUpkeep('0x')
      assert.isTrue(should)
      let [addresses] = ethers.utils.defaultAbiCoder.decode(
        ['address[]'],
        payload,
      )
      assert.deepEqual(addresses, [watchAddress1, watchAddress2, watchAddress3])
      // checkUpkeep payload should match getUnderfundedAddresses()
      addresses = await bm.getUnderfundedAddresses()
      assert.deepEqual(addresses, [watchAddress1, watchAddress2, watchAddress3])
    })

    it('Should return some results even if contract cannot fund all eligible targets', async () => {
      const fundTx = await owner.sendTransaction({
        to: bm.address,
        value: fiveEth, // needs 6 total
      })
      await fundTx.wait()
      const [should, payload] = await bm.checkUpkeep('0x')
      assert.isTrue(should)
      const [addresses] = ethers.utils.defaultAbiCoder.decode(
        ['address[]'],
        payload,
      )
      assert.deepEqual(addresses, [watchAddress1, watchAddress2])
    })

    it('Should omit addresses that have been funded recently', async () => {
      const setWaitPdTx = await bm.setMinWaitPeriodSeconds(3600) // 1 hour
      const fundTx = await owner.sendTransaction({
        to: bm.address,
        value: sixEth,
      })
      await Promise.all([setWaitPdTx.wait(), fundTx.wait()])
      const block = await ethers.provider.getBlock('latest')
      const setTopUpTx = await bm.setLastTopUpXXXTestOnly(
        watchAddress2,
        block.timestamp - 100,
      )
      await setTopUpTx.wait()
      const [should, payload] = await bm.checkUpkeep('0x')
      assert.isTrue(should)
      const [addresses] = ethers.utils.defaultAbiCoder.decode(
        ['address[]'],
        payload,
      )
      assert.deepEqual(addresses, [watchAddress1, watchAddress3])
    })

    it('Should revert when paused', async () => {
      const tx = await bm.connect(owner).pause()
      await tx.wait()
      const ethCall = bm.checkUpkeep('0x')
      await expect(ethCall).to.be.revertedWith(PAUSED_ERR)
    })
  })

  describe('performUpkeep()', () => {
    let validPayload: string
    let invalidPayload: string

    beforeEach(async () => {
      validPayload = ethers.utils.defaultAbiCoder.encode(
        ['address[]'],
        [[watchAddress1, watchAddress2, watchAddress3]],
      )
      invalidPayload = ethers.utils.defaultAbiCoder.encode(
        ['address[]'],
        [[watchAddress1, watchAddress2, watchAddress4, watchAddress5]],
      )
      const setTx = await bm.connect(owner).setWatchList(
        [
          watchAddress1, // needs funds
          watchAddress5, // funded
          watchAddress2, // needs funds
          watchAddress6, // funded
          watchAddress3, // needs funds
          // watchAddress4 - omitted
        ],
        new Array(5).fill(oneEth),
        new Array(5).fill(twoEth),
      )
      await setTx.wait()
    })

    it('Should revert when paused', async () => {
      const pauseTx = await bm.connect(owner).pause()
      await pauseTx.wait()
      const performTx = bm.connect(keeperRegistry).performUpkeep(validPayload)
      await expect(performTx).to.be.revertedWith(PAUSED_ERR)
    })

    context('when partially funded', () => {
      it('Should fund as many addresses as possible', async () => {
        const fundTx = await owner.sendTransaction({
          to: bm.address,
          value: fiveEth, // only enough eth to fund 2 addresses
        })
        await fundTx.wait()
        await assertWatchlistBalances(0, 0, 0, 0, 10_000, 10_000)
        const performTx = await bm
          .connect(keeperRegistry)
          .performUpkeep(validPayload)
        await assertWatchlistBalances(2, 2, 0, 0, 10_000, 10_000)
        await expect(performTx)
          .to.emit(bm, 'TopUpSucceeded')
          .withArgs(watchAddress1)
        await expect(performTx)
          .to.emit(bm, 'TopUpSucceeded')
          .withArgs(watchAddress2)
      })
    })

    context('when fully funded', () => {
      beforeEach(async () => {
        const fundTx = await owner.sendTransaction({
          to: bm.address,
          value: tenEth,
        })
        await fundTx.wait()
      })

      it('Should fund the appropriate addresses', async () => {
        await assertWatchlistBalances(0, 0, 0, 0, 10_000, 10_000)
        const performTx = await bm
          .connect(keeperRegistry)
          .performUpkeep(validPayload, { gasLimit: 2_500_000 })
        await performTx.wait()
        await assertWatchlistBalances(2, 2, 2, 0, 10_000, 10_000)
      })

      it('Should only fund active, underfunded addresses', async () => {
        await assertWatchlistBalances(0, 0, 0, 0, 10_000, 10_000)
        const performTx = await bm
          .connect(keeperRegistry)
          .performUpkeep(invalidPayload, { gasLimit: 2_500_000 })
        await performTx.wait()
        await assertWatchlistBalances(2, 2, 0, 0, 10_000, 10_000)
      })

      it('Should continue funding addresses even if one reverts', async () => {
        await assertWatchlistBalances(0, 0, 0, 0, 10_000, 10_000)
        const addresses = [
          watchAddress1,
          receiveReverter.address,
          watchAddress2,
        ]
        const setTx = await bm
          .connect(owner)
          .setWatchList(
            addresses,
            new Array(3).fill(oneEth),
            new Array(3).fill(twoEth),
          )
        await setTx.wait()
        const payload = ethers.utils.defaultAbiCoder.encode(
          ['address[]'],
          [addresses],
        )
        const performTx = await bm
          .connect(keeperRegistry)
          .performUpkeep(payload, { gasLimit: 2_500_000 })
        await performTx.wait()
        await assertWatchlistBalances(2, 2, 0, 0, 10_000, 10_000)
        await h.assertBalance(receiveReverter.address, 0)
        await expect(performTx)
          .to.emit(bm, 'TopUpSucceeded')
          .withArgs(watchAddress1)
        await expect(performTx)
          .to.emit(bm, 'TopUpSucceeded')
          .withArgs(watchAddress2)
        await expect(performTx)
          .to.emit(bm, 'TopUpFailed')
          .withArgs(receiveReverter.address)
      })

      it('Should not fund addresses that have been funded recently', async () => {
        const setWaitPdTx = await bm.setMinWaitPeriodSeconds(3600) // 1 hour
        await setWaitPdTx.wait()
        const block = await ethers.provider.getBlock('latest')
        const setTopUpTx = await bm.setLastTopUpXXXTestOnly(
          watchAddress2,
          block.timestamp - 100,
        )
        await setTopUpTx.wait()
        await assertWatchlistBalances(0, 0, 0, 0, 10_000, 10_000)
        const performTx = await bm
          .connect(keeperRegistry)
          .performUpkeep(validPayload, { gasLimit: 2_500_000 })
        await performTx.wait()
        await assertWatchlistBalances(2, 0, 2, 0, 10_000, 10_000)
      })

      it('Should only be callable by the keeper registry contract', async () => {
        let performTx = bm.connect(owner).performUpkeep(validPayload)
        await expect(performTx).to.be.revertedWithCustomError(
          bm,
          ONLY_KEEPER_ERR,
        )
        performTx = bm.connect(stranger).performUpkeep(validPayload)
        await expect(performTx).to.be.revertedWithCustomError(
          bm,
          ONLY_KEEPER_ERR,
        )
      })

      it('Should protect against running out of gas', async () => {
        await assertWatchlistBalances(0, 0, 0, 0, 10_000, 10_000)
        const performTx = await bm
          .connect(keeperRegistry)
          .performUpkeep(validPayload, { gasLimit: 130_000 }) // too little for all 3 transfers
        await performTx.wait()
        const balance1 = await ethers.provider.getBalance(watchAddress1)
        const balance2 = await ethers.provider.getBalance(watchAddress2)
        const balance3 = await ethers.provider.getBalance(watchAddress3)
        const balances = [balance1, balance2, balance3].map((n) => n.toString())
        expect(balances)
          .to.include(twoEth.toString()) // expect at least 1 transfer
          .to.include(zeroEth.toString()) // expect at least 1 out of funds
      })

      it('Should provide enough gas to support receive and fallback functions', async () => {
        const addresses = [
          receiveEmitter.address,
          receiveFallbackEmitter.address,
        ]
        const payload = ethers.utils.defaultAbiCoder.encode(
          ['address[]'],
          [addresses],
        )
        const setTx = await bm
          .connect(owner)
          .setWatchList(
            addresses,
            new Array(2).fill(oneEth),
            new Array(2).fill(twoEth),
          )
        await setTx.wait()

        const reBalanceBefore = await ethers.provider.getBalance(
          receiveEmitter.address,
        )
        const rfeBalanceBefore = await ethers.provider.getBalance(
          receiveFallbackEmitter.address,
        )

        const performTx = await bm
          .connect(keeperRegistry)
          .performUpkeep(payload, { gasLimit: 2_500_000 })
        await h.assertBalance(
          receiveEmitter.address,
          reBalanceBefore.add(twoEth),
        )
        await h.assertBalance(
          receiveFallbackEmitter.address,
          rfeBalanceBefore.add(twoEth),
        )

        await expect(performTx)
          .to.emit(bm, 'TopUpSucceeded')
          .withArgs(receiveEmitter.address)
        await expect(performTx)
          .to.emit(bm, 'TopUpSucceeded')
          .withArgs(receiveFallbackEmitter.address)
      })
    })
  })

  describe('topUp()', () => {
    context('when not paused', () => {
      it('Should be callable by anyone', async () => {
        const users = [owner, keeperRegistry, stranger]
        for (let idx = 0; idx < users.length; idx++) {
          const user = users[idx]
          await bm.connect(user).topUp([])
        }
      })
    })
    context('when paused', () => {
      it('Should be callable by no one', async () => {
        await bm.connect(owner).pause()
        const users = [owner, keeperRegistry, stranger]
        for (let idx = 0; idx < users.length; idx++) {
          const user = users[idx]
          const tx = bm.connect(user).topUp([])
          await expect(tx).to.be.revertedWith(PAUSED_ERR)
        }
      })
    })
  })
})
