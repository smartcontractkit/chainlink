import { ethers } from 'hardhat'
import { assert, expect } from 'chai'
import type { SignerWithAddress } from '@nomiclabs/hardhat-ethers/signers'
import { ReceiveEmitter } from '../../../typechain/ReceiveEmitter'
import { ReceiveFallbackEmitter } from '../../../typechain/ReceiveFallbackEmitter'
import * as h from '../../test-helpers/helpers'
import { ERC20BalanceMonitorExposed, LinkToken } from '../../../typechain'
import { BigNumber } from 'ethers'

const OWNABLE_ERR = 'Only callable by owner'
const INVALID_WATCHLIST_ERR = `InvalidWatchList()`
const PAUSED_ERR = 'Pausable: paused'
const ONLY_KEEPER_ERR = `OnlyKeeperRegistry()`

const zeroLINK = ethers.utils.parseEther('0')
const oneLINK = ethers.utils.parseEther('1')
const twoLINK = ethers.utils.parseEther('2')
const threeLINK = ethers.utils.parseEther('3')
const fiveLINK = ethers.utils.parseEther('5')
const sixLINK = ethers.utils.parseEther('6')
const tenLINK = ethers.utils.parseEther('10')

const oneHundredLINK = ethers.utils.parseEther('100')

const watchAddress1 = ethers.Wallet.createRandom().address
const watchAddress2 = ethers.Wallet.createRandom().address
const watchAddress3 = ethers.Wallet.createRandom().address
const watchAddress4 = ethers.Wallet.createRandom().address
let watchAddress5: string
let watchAddress6: string

let bm: ERC20BalanceMonitorExposed
let lt: LinkToken
let receiveEmitter: ReceiveEmitter
let receiveFallbackEmitter: ReceiveFallbackEmitter
let owner: SignerWithAddress
let stranger: SignerWithAddress
let keeperRegistry: SignerWithAddress

async function assertWatchlistBalances(
  balance1: BigNumber,
  balance2: BigNumber,
  balance3: BigNumber,
  balance4: BigNumber,
  balance5: BigNumber,
  balance6: BigNumber,
) {
  await h.assertLinkTokenBalance(lt, watchAddress1, balance1, 'address 1')
  await h.assertLinkTokenBalance(lt, watchAddress2, balance2, 'address 2')
  await h.assertLinkTokenBalance(lt, watchAddress3, balance3, 'address 3')
  await h.assertLinkTokenBalance(lt, watchAddress4, balance4, 'address 4')
  await h.assertLinkTokenBalance(lt, watchAddress5, balance5, 'address 5')
  await h.assertLinkTokenBalance(lt, watchAddress6, balance6, 'address 6')
}

describe('ERC20BalanceMonitor', () => {
  beforeEach(async () => {
    const accounts = await ethers.getSigners()
    owner = accounts[0]
    stranger = accounts[1]
    keeperRegistry = accounts[2]
    watchAddress5 = accounts[3].address
    watchAddress6 = accounts[4].address

    const bmFactory = await ethers.getContractFactory(
      'ERC20BalanceMonitorExposed',
      owner,
    )
    const ltFactory = await ethers.getContractFactory('LinkToken', owner)
    const reFactory = await ethers.getContractFactory('ReceiveEmitter', owner)
    const rfeFactory = await ethers.getContractFactory(
      'ReceiveFallbackEmitter',
      owner,
    )

    lt = await ltFactory.deploy()
    bm = await bmFactory.deploy(lt.address, keeperRegistry.address, 0)

    for (let i = 1; i <= 4; i++) {
      const recipient = await accounts[i].getAddress()
      await lt.connect(owner).transfer(recipient, oneHundredLINK)
    }

    receiveEmitter = await reFactory.deploy()
    receiveFallbackEmitter = await rfeFactory.deploy()
    await Promise.all([
      bm.deployed(),
      receiveEmitter.deployed(),
      receiveFallbackEmitter.deployed(),
    ])
  })

  afterEach(async () => {
    await h.reset()
  })

  describe('add funds', () => {
    it('Should allow anyone to add funds', async () => {
      await lt.transfer(bm.address, oneLINK)
      await lt.connect(stranger).transfer(bm.address, oneLINK)
    })
  })

  describe('withdraw()', () => {
    beforeEach(async () => {
      const tx = await lt.connect(owner).transfer(bm.address, oneLINK)
      await tx.wait()
    })

    it('Should allow the owner to withdraw', async () => {
      const beforeBalance = await lt.balanceOf(owner.address)
      const tx = await bm.connect(owner).withdraw(oneLINK, owner.address)
      await tx.wait()
      const afterBalance = await lt.balanceOf(owner.address)
      assert.isTrue(
        afterBalance.gt(beforeBalance),
        'balance did not increase after withdraw',
      )
    })

    it('Should emit an event', async () => {
      const tx = await bm.connect(owner).withdraw(oneLINK, owner.address)
      await expect(tx)
        .to.emit(bm, 'FundsWithdrawn')
        .withArgs(oneLINK, owner.address)
    })

    it('Should allow the owner to withdraw to anyone', async () => {
      const beforeBalance = await lt.balanceOf(stranger.address)
      const tx = await bm.connect(owner).withdraw(oneLINK, stranger.address)
      await tx.wait()
      const afterBalance = await lt.balanceOf(stranger.address)
      assert.isTrue(
        beforeBalance.add(oneLINK).eq(afterBalance),
        'balance did not increase after withdraw',
      )
    })

    it('Should not allow strangers to withdraw', async () => {
      const tx = bm.connect(stranger).withdraw(oneLINK, owner.address)
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
        .setWatchList([watchAddress1], [oneLINK], [twoLINK])
      await setTx.wait()
      let watchList = await bm.getWatchList()
      assert.deepEqual(watchList, [watchAddress1])
      const accountInfo = await bm.getAccountInfo(watchAddress1)
      assert.isTrue(accountInfo.isActive)
      expect(accountInfo.minBalance).to.equal(oneLINK)
      expect(accountInfo.topUpLevel).to.equal(twoLINK)
      // add more to watchlist
      setTx = await bm
        .connect(owner)
        .setWatchList(
          [watchAddress1, watchAddress2, watchAddress3],
          [oneLINK, twoLINK, threeLINK],
          [twoLINK, threeLINK, fiveLINK],
        )
      await setTx.wait()
      watchList = await bm.getWatchList()
      assert.deepEqual(watchList, [watchAddress1, watchAddress2, watchAddress3])
      let accountInfo1 = await bm.getAccountInfo(watchAddress1)
      let accountInfo2 = await bm.getAccountInfo(watchAddress2)
      let accountInfo3 = await bm.getAccountInfo(watchAddress3)
      expect(accountInfo1.isActive).to.be.true
      expect(accountInfo1.minBalance).to.equal(oneLINK)
      expect(accountInfo1.topUpLevel).to.equal(twoLINK)
      expect(accountInfo2.isActive).to.be.true
      expect(accountInfo2.minBalance).to.equal(twoLINK)
      expect(accountInfo2.topUpLevel).to.equal(threeLINK)
      expect(accountInfo3.isActive).to.be.true
      expect(accountInfo3.minBalance).to.equal(threeLINK)
      expect(accountInfo3.topUpLevel).to.equal(fiveLINK)
      // remove some from watchlist
      setTx = await bm
        .connect(owner)
        .setWatchList(
          [watchAddress3, watchAddress1],
          [threeLINK, oneLINK],
          [fiveLINK, twoLINK],
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
          [oneLINK, twoLINK, threeLINK],
          [twoLINK, threeLINK, fiveLINK],
        )
      await expect(setTx).to.be.revertedWith(errMsg)
    })

    it('Should not allow a topUpLevel les than or equal to minBalance in the watchlist', async () => {
      const setTx = bm
        .connect(owner)
        .setWatchList(
          [watchAddress1, watchAddress2, watchAddress1],
          [oneLINK, twoLINK, threeLINK],
          [zeroLINK, twoLINK, threeLINK],
        )
      await expect(setTx).to.be.revertedWith(INVALID_WATCHLIST_ERR)
    })

    it('Should not allow larger than maximum watchlist size', async () => {
      const watchlist: any[][] = [[], [], []]
      Array.from(Array(301).keys()).forEach(() => {
        watchlist[0].push(owner.address)
        watchlist[1].push(oneLINK)
        watchlist[2].push(twoLINK)
      })
      const tx = bm
        .connect(owner)
        .setWatchList(watchlist[0], watchlist[1], watchlist[2])
      await expect(tx).to.be.revertedWith(INVALID_WATCHLIST_ERR)
    })

    it('Should not allow strangers to set the watchlist', async () => {
      const setTxStranger = bm
        .connect(stranger)
        .setWatchList([watchAddress1], [oneLINK], [twoLINK])
      await expect(setTxStranger).to.be.revertedWith(OWNABLE_ERR)
    })

    it('Should revert if the list lengths differ', async () => {
      let tx = bm.connect(owner).setWatchList([watchAddress1], [], [twoLINK])
      await expect(tx).to.be.revertedWith(INVALID_WATCHLIST_ERR)
      tx = bm.connect(owner).setWatchList([watchAddress1], [oneLINK], [])
      await expect(tx).to.be.revertedWith(INVALID_WATCHLIST_ERR)
      tx = bm.connect(owner).setWatchList([], [oneLINK], [twoLINK])
      await expect(tx).to.be.revertedWith(INVALID_WATCHLIST_ERR)
    })

    it('Should revert if any of the addresses are empty', async () => {
      let tx = bm
        .connect(owner)
        .setWatchList(
          [watchAddress1, ethers.constants.AddressZero],
          [oneLINK, oneLINK],
          [twoLINK, twoLINK],
        )
      await expect(tx).to.be.revertedWith(INVALID_WATCHLIST_ERR)
    })

    it('Should revert if any of the top up amounts are 0', async () => {
      const tx = bm
        .connect(owner)
        .setWatchList(
          [watchAddress1, watchAddress2],
          [oneLINK, oneLINK],
          [twoLINK, zeroLINK],
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
        new Array(5).fill(oneLINK),
        new Array(5).fill(twoLINK),
      )
      await setTx.wait()
    })

    it('Should return list of address that are underfunded', async () => {
      const fundTx = await lt.connect(owner).transfer(
        bm.address,
        sixLINK, // needs 6 total
      )
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
      const fundTx = await lt.connect(owner).transfer(
        bm.address,
        fiveLINK, // needs 6 total
      )
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
      const fundTx = await lt.connect(owner).transfer(bm.address, sixLINK)
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
        new Array(5).fill(oneLINK),
        new Array(5).fill(twoLINK),
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
        const fundTx = await lt.connect(owner).transfer(
          bm.address,
          fiveLINK, // only enough LINK to fund 2 addresses
        )
        await fundTx.wait()
        await assertWatchlistBalances(
          zeroLINK,
          zeroLINK,
          zeroLINK,
          zeroLINK,
          oneHundredLINK,
          oneHundredLINK,
        )
        const performTx = await bm
          .connect(keeperRegistry)
          .performUpkeep(validPayload)
        await assertWatchlistBalances(
          twoLINK,
          twoLINK,
          zeroLINK,
          zeroLINK,
          oneHundredLINK,
          oneHundredLINK,
        )
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
        const fundTx = await lt.connect(owner).transfer(bm.address, tenLINK)
        await fundTx.wait()
      })

      it('Should fund the appropriate addresses', async () => {
        await assertWatchlistBalances(
          zeroLINK,
          zeroLINK,
          zeroLINK,
          zeroLINK,
          oneHundredLINK,
          oneHundredLINK,
        )
        const performTx = await bm
          .connect(keeperRegistry)
          .performUpkeep(validPayload, { gasLimit: 2_500_000 })
        await performTx.wait()
        await assertWatchlistBalances(
          twoLINK,
          twoLINK,
          twoLINK,
          zeroLINK,
          oneHundredLINK,
          oneHundredLINK,
        )
      })

      it('Should only fund active, underfunded addresses', async () => {
        await assertWatchlistBalances(
          zeroLINK,
          zeroLINK,
          zeroLINK,
          zeroLINK,
          oneHundredLINK,
          oneHundredLINK,
        )
        const performTx = await bm
          .connect(keeperRegistry)
          .performUpkeep(invalidPayload, { gasLimit: 2_500_000 })
        await performTx.wait()
        await assertWatchlistBalances(
          twoLINK,
          twoLINK,
          zeroLINK,
          zeroLINK,
          oneHundredLINK,
          oneHundredLINK,
        )
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
        await assertWatchlistBalances(
          zeroLINK,
          zeroLINK,
          zeroLINK,
          zeroLINK,
          oneHundredLINK,
          oneHundredLINK,
        )
        const performTx = await bm
          .connect(keeperRegistry)
          .performUpkeep(validPayload, { gasLimit: 2_500_000 })
        await performTx.wait()
        await assertWatchlistBalances(
          twoLINK,
          zeroLINK,
          twoLINK,
          zeroLINK,
          oneHundredLINK,
          oneHundredLINK,
        )
      })

      it('Should only be callable by the keeper registry contract', async () => {
        let performTx = bm.connect(owner).performUpkeep(validPayload)
        await expect(performTx).to.be.revertedWith(ONLY_KEEPER_ERR)
        performTx = bm.connect(stranger).performUpkeep(validPayload)
        await expect(performTx).to.be.revertedWith(ONLY_KEEPER_ERR)
      })

      it('Should protect against running out of gas', async () => {
        await assertWatchlistBalances(
          zeroLINK,
          zeroLINK,
          zeroLINK,
          zeroLINK,
          oneHundredLINK,
          oneHundredLINK,
        )
        const performTx = await bm
          .connect(keeperRegistry)
          .performUpkeep(validPayload, { gasLimit: 130_000 }) // too little for all 3 transfers
        await performTx.wait()
        const balance1 = await lt.balanceOf(watchAddress1)
        const balance2 = await lt.balanceOf(watchAddress2)
        const balance3 = await lt.balanceOf(watchAddress3)
        const balances = [balance1, balance2, balance3].map((n) => n.toString())
        expect(balances)
          .to.include(twoLINK.toString()) // expect at least 1 transfer
          .to.include(zeroLINK.toString()) // expect at least 1 out of funds
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
            new Array(2).fill(oneLINK),
            new Array(2).fill(twoLINK),
          )
        await setTx.wait()

        const reBalanceBefore = await lt.balanceOf(receiveEmitter.address)
        const rfeBalanceBefore = await lt.balanceOf(
          receiveFallbackEmitter.address,
        )
        const performTx = await bm
          .connect(keeperRegistry)
          .performUpkeep(payload, { gasLimit: 2_500_000 })
        await h.assertLinkTokenBalance(
          lt,
          receiveEmitter.address,
          reBalanceBefore.add(twoLINK),
        )
        await h.assertLinkTokenBalance(
          lt,
          receiveFallbackEmitter.address,
          rfeBalanceBefore.add(twoLINK),
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
