import { ethers } from 'hardhat'
import { assert, expect } from 'chai'
import type { SignerWithAddress } from '@nomiclabs/hardhat-ethers/signers'
import { EthBalanceMonitorExposed } from '../../typechain/EthBalanceMonitorExposed'
import { ReceiveReverter } from '../../typechain/ReceiveReverter'
import { ReceiveEmitter } from '../../typechain/ReceiveEmitter'
import { ReceiveFallbackEmitter } from '../../typechain/ReceiveFallbackEmitter'
import * as h from '../test-helpers/helpers'

const OWNABLE_ERR = 'Only callable by owner'
const INVALID_WATCHLIST_ERR = `InvalidWatchList()`

const zeroEth = ethers.utils.parseEther('0')
const oneEth = ethers.utils.parseEther('1')
const twoEth = ethers.utils.parseEther('2')
const threeEth = ethers.utils.parseEther('3')
const watchAddress1 = ethers.Wallet.createRandom().address
const watchAddress2 = ethers.Wallet.createRandom().address
const watchAddress3 = ethers.Wallet.createRandom().address

let bm: EthBalanceMonitorExposed
let receiveReverter: ReceiveReverter
let receiveEmitter: ReceiveEmitter
let receiveFallbackEmitter: ReceiveFallbackEmitter
let owner: SignerWithAddress
let stranger: SignerWithAddress
let keeperRegistry: SignerWithAddress

describe('EthBalanceMonitor 1/2', () => {
  beforeEach(async () => {
    const accounts = await ethers.getSigners()
    owner = accounts[0]
    stranger = accounts[1]
    keeperRegistry = accounts[2]

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
      const errMsg = `DuplicateAddress("${watchAddress1}")`
      const setTx = bm
        .connect(owner)
        .setWatchList(
          [watchAddress1, watchAddress2, watchAddress1],
          [oneEth, twoEth, threeEth],
          [oneEth, twoEth, threeEth],
        )
      await expect(setTx).to.be.revertedWith(errMsg)
    })

    it('Should not allow strangers to set the watchlist', async () => {
      const setTxStranger = bm
        .connect(stranger)
        .setWatchList([watchAddress1], [oneEth], [twoEth])
      await expect(setTxStranger).to.be.revertedWith(OWNABLE_ERR)
    })

    it('Should revert if the list lengths differ', async () => {
      let tx = bm.connect(owner).setWatchList([watchAddress1], [], [twoEth])
      await expect(tx).to.be.revertedWith(INVALID_WATCHLIST_ERR)
      tx = bm.connect(owner).setWatchList([watchAddress1], [oneEth], [])
      await expect(tx).to.be.revertedWith(INVALID_WATCHLIST_ERR)
      tx = bm.connect(owner).setWatchList([], [oneEth], [twoEth])
      await expect(tx).to.be.revertedWith(INVALID_WATCHLIST_ERR)
    })

    it('Should revert if any of the addresses are empty', async () => {
      let tx = bm
        .connect(owner)
        .setWatchList(
          [watchAddress1, ethers.constants.AddressZero],
          [oneEth, oneEth],
          [twoEth, twoEth],
        )
      await expect(tx).to.be.revertedWith(INVALID_WATCHLIST_ERR)
    })

    it('Should revert if any of the top up amounts are 0', async () => {
      const tx = bm
        .connect(owner)
        .setWatchList(
          [watchAddress1, watchAddress2],
          [oneEth, oneEth],
          [twoEth, zeroEth],
        )
      await expect(tx).to.be.revertedWith(INVALID_WATCHLIST_ERR)
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
})
