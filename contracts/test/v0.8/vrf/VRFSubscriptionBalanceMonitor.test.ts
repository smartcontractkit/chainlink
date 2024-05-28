import { ethers } from 'hardhat'
import { assert, expect } from 'chai'
import type { SignerWithAddress } from '@nomiclabs/hardhat-ethers/signers'
import {
  LinkToken,
  VRFSubscriptionBalanceMonitorExposed,
} from '../../../typechain'
import * as h from '../../test-helpers/helpers'
import { BigNumber, Contract } from 'ethers'

const OWNABLE_ERR = 'Only callable by owner'
const INVALID_WATCHLIST_ERR = `InvalidWatchList`
const PAUSED_ERR = 'Pausable: paused'
const ONLY_KEEPER_ERR = `OnlyKeeperRegistry`

const zeroLINK = ethers.utils.parseEther('0')
const oneLINK = ethers.utils.parseEther('1')
const twoLINK = ethers.utils.parseEther('2')
const threeLINK = ethers.utils.parseEther('3')
const fiveLINK = ethers.utils.parseEther('5')
const sixLINK = ethers.utils.parseEther('6')
const tenLINK = ethers.utils.parseEther('10')
const oneHundredLINK = ethers.utils.parseEther('100')

let lt: LinkToken
let coordinator: Contract
let bm: VRFSubscriptionBalanceMonitorExposed
let owner: SignerWithAddress
let stranger: SignerWithAddress
let keeperRegistry: SignerWithAddress

const sub1 = BigNumber.from(1)
const sub2 = BigNumber.from(2)
const sub3 = BigNumber.from(3)
const sub4 = BigNumber.from(4)
const sub5 = BigNumber.from(5)
const sub6 = BigNumber.from(6)

const toNums = (bigNums: BigNumber[]) => bigNums.map((n) => n.toNumber())

async function assertWatchlistBalances(
  balance1: BigNumber,
  balance2: BigNumber,
  balance3: BigNumber,
  balance4: BigNumber,
  balance5: BigNumber,
  balance6: BigNumber,
) {
  await h.assertSubscriptionBalance(coordinator, sub1, balance1, 'sub 1')
  await h.assertSubscriptionBalance(coordinator, sub2, balance2, 'sub 2')
  await h.assertSubscriptionBalance(coordinator, sub3, balance3, 'sub 3')
  await h.assertSubscriptionBalance(coordinator, sub4, balance4, 'sub 4')
  await h.assertSubscriptionBalance(coordinator, sub5, balance5, 'sub 5')
  await h.assertSubscriptionBalance(coordinator, sub6, balance6, 'sub 6')
}

describe('VRFSubscriptionBalanceMonitor', () => {
  beforeEach(async () => {
    const accounts = await ethers.getSigners()
    owner = accounts[0]
    stranger = accounts[1]
    keeperRegistry = accounts[2]

    const bmFactory = await ethers.getContractFactory(
      'VRFSubscriptionBalanceMonitorExposed',
      owner,
    )
    const ltFactory = await ethers.getContractFactory(
      'src/v0.8/shared/test/helpers/LinkTokenTestHelper.sol:LinkTokenTestHelper',
      owner,
    )

    const coordinatorFactory = await ethers.getContractFactory(
      'src/v0.8/vrf/VRFCoordinatorV2.sol:VRFCoordinatorV2',
      owner,
    )

    lt = await ltFactory.deploy()
    coordinator = await coordinatorFactory.deploy(
      lt.address,
      lt.address,
      lt.address,
    ) // we don't use BHS or LinkEthFeed
    bm = await bmFactory.deploy(
      lt.address,
      coordinator.address,
      keeperRegistry.address,
      0,
    )

    for (let i = 0; i <= 5; i++) {
      await coordinator.connect(owner).createSubscription()
    }

    // Transfer LINK to stranger.
    await lt.transfer(stranger.address, oneHundredLINK)

    // Fund sub 5.
    await lt
      .connect(owner)
      .transferAndCall(
        coordinator.address,
        oneHundredLINK,
        ethers.utils.defaultAbiCoder.encode(['uint256'], ['5']),
      )

    // Fun sub 6.
    await lt
      .connect(owner)
      .transferAndCall(
        coordinator.address,
        oneHundredLINK,
        ethers.utils.defaultAbiCoder.encode(['uint256'], ['6']),
      )

    await Promise.all([bm.deployed(), coordinator.deployed(), lt.deployed()])
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
      assert.isFalse((await bm.getSubscriptionInfo(sub1)).isActive)
      // add first watchlist
      let setTx = await bm
        .connect(owner)
        .setWatchList([sub1], [oneLINK], [twoLINK])
      await setTx.wait()
      let watchList = await bm.getWatchList()
      assert.deepEqual(toNums(watchList), toNums([sub1]))
      const subInfo = await bm.getSubscriptionInfo(1)
      assert.isTrue(subInfo.isActive)
      expect(subInfo.minBalanceJuels).to.equal(oneLINK)
      expect(subInfo.topUpAmountJuels).to.equal(twoLINK)
      // add more to watchlist
      setTx = await bm
        .connect(owner)
        .setWatchList(
          [1, 2, 3],
          [oneLINK, twoLINK, threeLINK],
          [twoLINK, threeLINK, fiveLINK],
        )
      await setTx.wait()
      watchList = await bm.getWatchList()
      assert.deepEqual(toNums(watchList), toNums([sub1, sub2, sub3]))
      let subInfo1 = await bm.getSubscriptionInfo(sub1)
      let subInfo2 = await bm.getSubscriptionInfo(sub2)
      let subInfo3 = await bm.getSubscriptionInfo(sub3)
      expect(subInfo1.isActive).to.be.true
      expect(subInfo1.minBalanceJuels).to.equal(oneLINK)
      expect(subInfo1.topUpAmountJuels).to.equal(twoLINK)
      expect(subInfo2.isActive).to.be.true
      expect(subInfo2.minBalanceJuels).to.equal(twoLINK)
      expect(subInfo2.topUpAmountJuels).to.equal(threeLINK)
      expect(subInfo3.isActive).to.be.true
      expect(subInfo3.minBalanceJuels).to.equal(threeLINK)
      expect(subInfo3.topUpAmountJuels).to.equal(fiveLINK)
      // remove some from watchlist
      setTx = await bm
        .connect(owner)
        .setWatchList([sub3, sub1], [threeLINK, oneLINK], [fiveLINK, twoLINK])
      await setTx.wait()
      watchList = await bm.getWatchList()
      assert.deepEqual(toNums(watchList), toNums([sub3, sub1]))
      subInfo1 = await bm.getSubscriptionInfo(sub1)
      subInfo2 = await bm.getSubscriptionInfo(sub2)
      subInfo3 = await bm.getSubscriptionInfo(sub3)
      expect(subInfo1.isActive).to.be.true
      expect(subInfo2.isActive).to.be.false
      expect(subInfo3.isActive).to.be.true
    })

    it('Should not allow duplicates in the watchlist', async () => {
      const errMsg = `DuplicateSubcriptionId`
      const setTx = bm
        .connect(owner)
        .setWatchList(
          [sub1, sub2, sub1],
          [oneLINK, twoLINK, threeLINK],
          [twoLINK, threeLINK, fiveLINK],
        )
      await expect(setTx)
        .to.be.revertedWithCustomError(bm, errMsg)
        .withArgs(sub1)
    })

    it('Should not allow a topUpAmountJuels les than or equal to minBalance in the watchlist', async () => {
      const setTx = bm
        .connect(owner)
        .setWatchList(
          [sub1, sub2, sub1],
          [oneLINK, twoLINK, threeLINK],
          [zeroLINK, twoLINK, threeLINK],
        )
      await expect(setTx).to.be.revertedWithCustomError(
        bm,
        INVALID_WATCHLIST_ERR,
      )
    })

    it('Should not allow strangers to set the watchlist', async () => {
      const setTxStranger = bm
        .connect(stranger)
        .setWatchList([sub1], [oneLINK], [twoLINK])
      await expect(setTxStranger).to.be.revertedWith(OWNABLE_ERR)
    })

    it('Should revert if the list lengths differ', async () => {
      let tx = bm.connect(owner).setWatchList([sub1], [], [twoLINK])
      await expect(tx).to.be.revertedWithCustomError(bm, INVALID_WATCHLIST_ERR)
      tx = bm.connect(owner).setWatchList([sub1], [oneLINK], [])
      await expect(tx).to.be.revertedWithCustomError(bm, INVALID_WATCHLIST_ERR)
      tx = bm.connect(owner).setWatchList([], [oneLINK], [twoLINK])
      await expect(tx).to.be.revertedWithCustomError(bm, INVALID_WATCHLIST_ERR)
    })

    it('Should revert if any of the subIDs are zero', async () => {
      let tx = bm
        .connect(owner)
        .setWatchList([sub1, 0], [oneLINK, oneLINK], [twoLINK, twoLINK])
      await expect(tx).to.be.revertedWithCustomError(bm, INVALID_WATCHLIST_ERR)
    })

    it('Should revert if any of the top up amounts are 0', async () => {
      const tx = bm
        .connect(owner)
        .setWatchList([sub1, sub2], [oneLINK, oneLINK], [twoLINK, zeroLINK])
      await expect(tx).to.be.revertedWithCustomError(bm, INVALID_WATCHLIST_ERR)
    })
  })

  describe('getKeeperRegistryAddress() / setKeeperRegistryAddress()', () => {
    const newAddress = ethers.Wallet.createRandom().address

    it('Should initialize with the registry address provided to the constructor', async () => {
      const address = await bm.s_keeperRegistryAddress()
      assert.equal(address, keeperRegistry.address)
    })

    it('Should allow the owner to set the registry address', async () => {
      const setTx = await bm.connect(owner).setKeeperRegistryAddress(newAddress)
      await setTx.wait()
      const address = await bm.s_keeperRegistryAddress()
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
      const minWaitPeriod = await bm.s_minWaitPeriodSeconds()
      expect(minWaitPeriod).to.equal(0)
    })

    it('Should allow owner to set the wait period', async () => {
      const setTx = await bm
        .connect(owner)
        .setMinWaitPeriodSeconds(newWaitPeriod)
      await setTx.wait()
      const minWaitPeriod = await bm.s_minWaitPeriodSeconds()
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

  describe('checkUpkeep() / getUnderfundedSubscriptions()', () => {
    beforeEach(async () => {
      const setTx = await bm.connect(owner).setWatchList(
        [
          sub1, // needs funds
          sub5, // funded
          sub2, // needs funds
          sub6, // funded
          sub3, // needs funds
        ],
        new Array(5).fill(oneLINK),
        new Array(5).fill(twoLINK),
      )
      await setTx.wait()
    })

    it('Should return list of subscriptions that are underfunded', async () => {
      const fundTx = await lt.connect(owner).transfer(
        bm.address,
        sixLINK, // needs 6 total
      )
      await fundTx.wait()
      const [should, payload] = await bm.checkUpkeep('0x')
      assert.isTrue(should)
      let [subs] = ethers.utils.defaultAbiCoder.decode(['uint64[]'], payload)
      assert.deepEqual(toNums(subs), toNums([sub1, sub2, sub3]))
      // checkUpkeep payload should match getUnderfundedSubscriptions()
      subs = await bm.getUnderfundedSubscriptions()
      assert.deepEqual(toNums(subs), toNums([sub1, sub2, sub3]))
    })

    it('Should return some results even if contract cannot fund all eligible targets', async () => {
      const fundTx = await lt.connect(owner).transfer(
        bm.address,
        fiveLINK, // needs 6 total
      )
      await fundTx.wait()
      const [should, payload] = await bm.checkUpkeep('0x')
      assert.isTrue(should)
      const [subs] = ethers.utils.defaultAbiCoder.decode(['uint64[]'], payload)
      assert.deepEqual(toNums(subs), toNums([sub1, sub2]))
    })

    it('Should omit subscriptions that have been funded recently', async () => {
      const setWaitPdTx = await bm.setMinWaitPeriodSeconds(3600) // 1 hour
      const fundTx = await lt.connect(owner).transfer(bm.address, sixLINK)
      await Promise.all([setWaitPdTx.wait(), fundTx.wait()])
      const block = await ethers.provider.getBlock('latest')
      const setTopUpTx = await bm.setLastTopUpXXXTestOnly(
        sub2,
        block.timestamp - 100,
      )
      await setTopUpTx.wait()
      const [should, payload] = await bm.checkUpkeep('0x')
      assert.isTrue(should)
      const [subs] = ethers.utils.defaultAbiCoder.decode(['uint64[]'], payload)
      assert.deepEqual(toNums(subs), toNums([sub1, sub3]))
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
        ['uint64[]'],
        [[sub1, sub2, sub3]],
      )
      invalidPayload = ethers.utils.defaultAbiCoder.encode(
        ['uint64[]'],
        [[sub1, sub2, sub4, sub5]],
      )
      const setTx = await bm.connect(owner).setWatchList(
        [
          sub1, // needs funds
          sub5, // funded
          sub2, // needs funds
          sub6, // funded
          sub3, // needs funds
          // sub4 - omitted
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
      it('Should fund as many subscriptions as possible', async () => {
        const fundTx = await lt.connect(owner).transfer(
          bm.address,
          fiveLINK, // only enough LINK to fund 2 subscriptions
        )
        await fundTx.wait()
        console.log((await lt.balanceOf(bm.address)).toString())
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

        await assertWatchlistBalances(
          twoLINK,
          twoLINK,
          zeroLINK,
          zeroLINK,
          oneHundredLINK,
          oneHundredLINK,
        )
        await expect(performTx).to.emit(bm, 'TopUpSucceeded').withArgs(sub1)
        await expect(performTx).to.emit(bm, 'TopUpSucceeded').withArgs(sub1)
      })
    })

    context('when fully funded', () => {
      beforeEach(async () => {
        const fundTx = await lt.connect(owner).transfer(bm.address, tenLINK)
        await fundTx.wait()
      })

      it('Should fund the appropriate subscriptions', async () => {
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

      it('Should only fund active, underfunded subscriptions', async () => {
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

      it('Should not fund subscriptions that have been funded recently', async () => {
        const setWaitPdTx = await bm.setMinWaitPeriodSeconds(3600) // 1 hour
        await setWaitPdTx.wait()
        const block = await ethers.provider.getBlock('latest')
        const setTopUpTx = await bm.setLastTopUpXXXTestOnly(
          sub2,
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
        const balance1 = (await coordinator.getSubscription(sub1)).balance
        const balance2 = (await coordinator.getSubscription(sub2)).balance
        const balance3 = (await coordinator.getSubscription(sub3)).balance
        const balances = [balance1, balance2, balance3].map((n) => n.toString())
        expect(balances)
          .to.include(twoLINK.toString()) // expect at least 1 transfer
          .to.include(zeroLINK.toString()) // expect at least 1 out of funds
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
