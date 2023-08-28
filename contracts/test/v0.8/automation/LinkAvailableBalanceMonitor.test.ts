import { ethers } from 'hardhat'
import chai, { assert, expect } from 'chai'
import type { SignerWithAddress } from '@nomiclabs/hardhat-ethers/signers'
import { loadFixture } from '@nomicfoundation/hardhat-network-helpers'
import * as h from '../../test-helpers/helpers'
import { mineBlock } from '../../test-helpers/helpers'
import { IAggregatorProxy__factory as IAggregatorProxyFactory } from '../../../typechain/factories/IAggregatorProxy__factory'
import { ILinkAvailable__factory as ILinkAvailableFactory } from '../../../typechain/factories/ILinkAvailable__factory'
import { LinkAvailableBalanceMonitor, LinkToken } from '../../../typechain'
import { BigNumber } from 'ethers'
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

let labm: LinkAvailableBalanceMonitor
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

let directTarget1: MockContract // Contracts which are direct target of balance monitoring without proxy
let directTarget2: MockContract

let watchListAddresses: string[]
let watchListMinBalances: BigNumber[]

async function assertContractLinkBalances(
  balance1: BigNumber,
  balance2: BigNumber,
  balance3: BigNumber,
  balance4: BigNumber,
  balance5: BigNumber,
) {
  await h.assertLinkTokenBalance(lt, aggregator1.address, balance1, 'address 1')
  await h.assertLinkTokenBalance(lt, aggregator2.address, balance2, 'address 2')
  await h.assertLinkTokenBalance(lt, aggregator3.address, balance3, 'address 3')
  await h.assertLinkTokenBalance(
    lt,
    directTarget1.address,
    balance4,
    'address 4',
  )
  await h.assertLinkTokenBalance(
    lt,
    directTarget2.address,
    balance5,
    'address 5',
  )
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
  directTarget1 = await deployMockContract(owner, ILinkAvailableFactory.abi)
  directTarget2 = await deployMockContract(owner, ILinkAvailableFactory.abi)

  await proxy1.deployed()
  await proxy2.deployed()
  await proxy3.deployed()
  await proxy4.deployed()
  await aggregator1.deployed()
  await aggregator2.deployed()
  await aggregator3.deployed()
  await aggregator4.deployed()
  await directTarget1.deployed()
  await directTarget2.deployed()

  watchListAddresses = [
    proxy1.address,
    proxy2.address,
    proxy3.address,
    directTarget1.address,
    directTarget2.address,
  ]
  watchListMinBalances = [oneLINK, oneLINK, oneLINK, twoLINK, twoLINK]

  await proxy1.mock.aggregator.returns(aggregator1.address)
  await proxy2.mock.aggregator.returns(aggregator2.address)
  await proxy3.mock.aggregator.returns(aggregator3.address)

  await aggregator1.mock.linkAvailableForPayment.returns(0)
  await aggregator2.mock.linkAvailableForPayment.returns(0)
  await aggregator3.mock.linkAvailableForPayment.returns(0)

  await directTarget1.mock.linkAvailableForPayment.returns(0)
  await directTarget2.mock.linkAvailableForPayment.returns(0)

  const labmFactory = await ethers.getContractFactory(
    'LinkAvailableBalanceMonitor',
    owner,
  )
  const ltFactory = await ethers.getContractFactory(
    'src/v0.4/LinkToken.sol:LinkToken',
    owner,
  )

  lt = await ltFactory.deploy()
  labm = await labmFactory.deploy(lt.address, twoLINK)
  await labm.deployed()

  for (let i = 1; i <= 4; i++) {
    const recipient = await accounts[i].getAddress()
    await lt.connect(owner).transfer(recipient, oneHundredLINK)
  }

  const setTx = await labm
    .connect(owner)
    .setWatchList(watchListAddresses, watchListMinBalances)
  await setTx.wait()
}

describe('LinkAvailableBalanceMonitor', () => {
  beforeEach(async () => {
    await loadFixture(setup)
  })

  describe('add funds', () => {
    it('Should allow anyone to add funds', async () => {
      await lt.transfer(labm.address, oneLINK)
      await lt.connect(stranger).transfer(labm.address, oneLINK)
    })
  })

  describe('setTopUpAmount()', () => {
    it('configures the top-up amount', async () => {
      await labm.connect(owner).setTopUpAmount(100)
      assert.equal((await labm.getTopUpAmount()).toNumber(), 100)
    })

    it('configuresis only callable by the owner', async () => {
      await expect(labm.connect(stranger).setTopUpAmount(100)).to.be.reverted
    })
  })

  describe('setMinBalance()', () => {
    it('configures the min balance', async () => {
      await labm.connect(owner).setMinBalance(proxy1.address, 100)
      assert.equal((await labm.getMinBalance(proxy1.address)).toNumber(), 100)
    })

    it('reverts if address is not in the watchlist', async () => {
      await expect(labm.connect(owner).setMinBalance(proxy4.address, 100)).to.be
        .reverted
    })

    it('is only callable by the owner', async () => {
      await expect(labm.connect(stranger).setMinBalance(proxy1.address, 100)).to
        .be.reverted
    })
  })

  describe('withdraw()', () => {
    beforeEach(async () => {
      const tx = await lt.connect(owner).transfer(labm.address, oneLINK)
      await tx.wait()
    })

    it('Should allow the owner to withdraw', async () => {
      const beforeBalance = await lt.balanceOf(owner.address)
      const tx = await labm.connect(owner).withdraw(oneLINK, owner.address)
      await tx.wait()
      const afterBalance = await lt.balanceOf(owner.address)
      assert.isTrue(
        afterBalance.gt(beforeBalance),
        'balance did not increase after withdraw',
      )
    })

    it('Should emit an event', async () => {
      const tx = await labm.connect(owner).withdraw(oneLINK, owner.address)
      await expect(tx)
        .to.emit(labm, 'FundsWithdrawn')
        .withArgs(oneLINK, owner.address)
    })

    it('Should allow the owner to withdraw to anyone', async () => {
      const beforeBalance = await lt.balanceOf(stranger.address)
      const tx = await labm.connect(owner).withdraw(oneLINK, stranger.address)
      await tx.wait()
      const afterBalance = await lt.balanceOf(stranger.address)
      assert.isTrue(
        beforeBalance.add(oneLINK).eq(afterBalance),
        'balance did not increase after withdraw',
      )
    })

    it('Should not allow strangers to withdraw', async () => {
      const tx = labm.connect(stranger).withdraw(oneLINK, owner.address)
      await expect(tx).to.be.revertedWith(OWNABLE_ERR)
    })
  })

  describe('pause() / unpause()', () => {
    it('Should allow owner to pause / unpause', async () => {
      const pauseTx = await labm.connect(owner).pause()
      await pauseTx.wait()
      const unpauseTx = await labm.connect(owner).unpause()
      await unpauseTx.wait()
    })

    it('Should not allow strangers to pause / unpause', async () => {
      const pauseTxStranger = labm.connect(stranger).pause()
      await expect(pauseTxStranger).to.be.revertedWith(OWNABLE_ERR)
      const pauseTxOwner = await labm.connect(owner).pause()
      await pauseTxOwner.wait()
      const unpauseTxStranger = labm.connect(stranger).unpause()
      await expect(unpauseTxStranger).to.be.revertedWith(OWNABLE_ERR)
    })
  })

  describe('setWatchList() / addToWatchList() / removeFromWatchlist() / getWatchList()', () => {
    const watchAddress1 = randAddr()
    const watchAddress2 = randAddr()
    const watchAddress3 = randAddr()

    beforeEach(async () => {
      // reset watchlist to empty before running these tests
      await labm.connect(owner).setWatchList([], [])
      let watchList = await labm.getWatchList()
      assert.deepEqual(watchList, [[], []])
    })

    it('Should allow owner to adjust the watchlist', async () => {
      // add first watchlist
      let tx = await labm
        .connect(owner)
        .setWatchList([watchAddress1], [oneLINK])
      let watchList = await labm.getWatchList()
      assert.deepEqual(watchList[0], [watchAddress1])
      assert.deepEqual(
        watchList[1].map((x) => x.toString()),
        [oneLINK].map((x) => x.toString()),
      )
      // add more to watchlist
      tx = await labm
        .connect(owner)
        .setWatchList(
          [watchAddress1, watchAddress2, watchAddress3],
          [oneLINK, oneLINK, oneLINK],
        )
      await tx.wait()
      watchList = await labm.getWatchList()
      assert.deepEqual(watchList[0], [
        watchAddress1,
        watchAddress2,
        watchAddress3,
      ])
      assert.deepEqual(
        watchList[1].map((x) => x.toString()),
        [oneLINK, oneLINK, oneLINK].map((x) => x.toString()),
      )
      // remove some from watchlist
      tx = await labm
        .connect(owner)
        .removeFromWatchlist([watchAddress3, watchAddress1])
      await tx.wait()
      watchList = await labm.getWatchList()
      assert.deepEqual(watchList[0], [watchAddress2])
      assert.deepEqual(
        watchList[1].map((x) => x.toString()),
        [oneLINK].map((x) => x.toString()),
      )
      // add some to watchlist
      tx = await labm
        .connect(owner)
        .addToWatchList([watchAddress1, watchAddress3], [twoLINK, twoLINK])
      await tx.wait()
      watchList = await labm.getWatchList()
      assert.deepEqual(watchList[0], [
        watchAddress2,
        watchAddress1,
        watchAddress3,
      ])
      assert.deepEqual(
        watchList[1].map((x) => x.toString()),
        [oneLINK, twoLINK, twoLINK].map((x) => x.toString()),
      )
    })

    it('Should not allow different length arrays in the watchlist', async () => {
      const errMsg = `InvalidWatchList()`
      let tx = labm
        .connect(owner)
        .setWatchList(
          [watchAddress1, watchAddress2, watchAddress1],
          [oneLINK, oneLINK],
        )
      await expect(tx).to.be.revertedWith(errMsg)
    })

    it('Should not allow duplicates in the watchlist', async () => {
      const errMsg = `DuplicateAddress("${watchAddress1}")`
      let tx = labm
        .connect(owner)
        .setWatchList(
          [watchAddress1, watchAddress2, watchAddress1],
          [oneLINK, oneLINK, oneLINK],
        )
      await expect(tx).to.be.revertedWith(errMsg)
      tx = labm
        .connect(owner)
        .addToWatchList(
          [watchAddress1, watchAddress2, watchAddress1],
          [oneLINK, oneLINK, oneLINK],
        )
      await expect(tx).to.be.revertedWith(errMsg)
    })

    it('Should not allow strangers to set the watchlist', async () => {
      const setTxStranger = labm
        .connect(stranger)
        .setWatchList([watchAddress1], [oneLINK])
      await expect(setTxStranger).to.be.revertedWith(OWNABLE_ERR)
      const addTxStranger = labm
        .connect(stranger)
        .addToWatchList([watchAddress1], [oneLINK])
      await expect(addTxStranger).to.be.revertedWith(OWNABLE_ERR)
      const removeTxStranger = labm.connect(stranger).removeFromWatchlist([])
      await expect(removeTxStranger).to.be.revertedWith(OWNABLE_ERR)
    })

    it('Should revert if any of the addresses are empty', async () => {
      let tx = labm
        .connect(owner)
        .setWatchList(
          [watchAddress1, ethers.constants.AddressZero],
          [oneLINK, oneLINK],
        )
      await expect(tx).to.be.revertedWith(INVALID_WATCHLIST_ERR)
      tx = labm
        .connect(owner)
        .addToWatchList(
          [watchAddress1, ethers.constants.AddressZero],
          [oneLINK, oneLINK],
        )
      await expect(tx).to.be.revertedWith(INVALID_WATCHLIST_ERR)
    })
  })

  describe('checkUpkeep() / sampleUnderfundedAddresses() [ @skip-coverage ]', () => {
    it('Should return list of address that are underfunded', async () => {
      const fundTx = await lt.connect(owner).transfer(
        labm.address,
        tenLINK, // needs 10 total
      )
      await fundTx.wait()
      const [should, payload] = await labm.checkUpkeep('0x')
      assert.isTrue(should)
      let [addresses] = ethers.utils.defaultAbiCoder.decode(
        ['address[]'],
        payload,
      )
      expect(addresses).to.deep.equalInAnyOrder(watchListAddresses)
      // checkUpkeep payload should match sampleUnderfundedAddresses()
      addresses = await labm.sampleUnderfundedAddresses()
      expect(addresses).to.deep.equalInAnyOrder(watchListAddresses)
    })

    it('Should return some results even if contract cannot fund all eligible targets', async () => {
      const fundTx = await lt.connect(owner).transfer(
        labm.address,
        fiveLINK, // needs 2Link per contract, so can fund 2 max
      )
      await fundTx.wait()
      const [should, payload] = await labm.checkUpkeep('0x')
      assert.isTrue(should)
      let [addresses] = ethers.utils.defaultAbiCoder.decode(
        ['address[]'],
        payload,
      )
      assert.equal(addresses.length, 2)
      assert.notEqual(addresses[0], addresses[1])
      assert(watchListAddresses.includes(addresses[0]))
      assert(watchListAddresses.includes(addresses[1]))
      // underfunded sample should still match list
      addresses = await labm.sampleUnderfundedAddresses()
      expect(addresses).to.deep.equalInAnyOrder(watchListAddresses)
    })

    it('Should omit aggregators that have sufficient funding', async () => {
      let addresses = await labm.sampleUnderfundedAddresses()
      expect(addresses).to.deep.equalInAnyOrder(watchListAddresses)
      await aggregator2.mock.linkAvailableForPayment.returns(oneLINK) // aggregator2 is enough funded
      await directTarget1.mock.linkAvailableForPayment.returns(oneLINK) // directTarget1 is NOT enough funded
      await directTarget2.mock.linkAvailableForPayment.returns(twoLINK) // directTarget2 is enough funded
      addresses = await labm.sampleUnderfundedAddresses()
      expect(addresses).to.deep.equalInAnyOrder([
        proxy1.address,
        proxy3.address,
        directTarget1.address,
      ])

      await aggregator1.mock.linkAvailableForPayment.returns(tenLINK)
      addresses = await labm.sampleUnderfundedAddresses()
      expect(addresses).to.deep.equalInAnyOrder([
        proxy3.address,
        directTarget1.address,
      ])

      await aggregator3.mock.linkAvailableForPayment.returns(tenLINK)
      addresses = await labm.sampleUnderfundedAddresses()
      expect(addresses).to.deep.equalInAnyOrder([directTarget1.address])

      await directTarget1.mock.linkAvailableForPayment.returns(tenLINK)
      addresses = await labm.sampleUnderfundedAddresses()
      expect(addresses).to.deep.equalInAnyOrder([])
    })

    it('Should revert when paused', async () => {
      const tx = await labm.connect(owner).pause()
      await tx.wait()
      const ethCall = labm.checkUpkeep('0x')
      await expect(ethCall).to.be.revertedWith(PAUSED_ERR)
    })

    context('with a large set of proxies', async () => {
      // in this test, we cheat a little bit and point each proxy to the same aggregator,
      // which helps cut down on test time
      let MAX_PERFORM: number
      let MAX_CHECK: number
      let proxyAddresses: string[]
      let minBalances: BigNumber[]
      let aggregators: MockContract[]

      beforeEach(async () => {
        MAX_PERFORM = (await labm.MAX_PERFORM()).toNumber()
        MAX_CHECK = (await labm.MAX_CHECK()).toNumber()
        proxyAddresses = []
        minBalances = []
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
          proxyAddresses.push(proxy.address)
          minBalances.push(oneLINK)
          aggregators.push(aggregator)
        }
        await labm.setWatchList(proxyAddresses, minBalances)
        expect(await labm.getWatchList()).to.deep.equalInAnyOrder([
          proxyAddresses,
          minBalances,
        ])
      })

      it('Should not include more than MAX_PERFORM addresses', async () => {
        const addresses = await labm.sampleUnderfundedAddresses()
        assert.equal(addresses.length, MAX_PERFORM)
      })

      it('Should sample from the list of addresses pseudorandomly', async () => {
        const firstAddress: string[] = []
        for (let idx = 0; idx < 10; idx++) {
          const addresses = await labm.sampleUnderfundedAddresses()
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
        await labm.checkUpkeep('0x', { gasLimit: TARGET_CHECK_GAS_LIMIT })
      })
    })
  })

  describe('performUpkeep()', () => {
    let validPayload: string

    beforeEach(async () => {
      validPayload = ethers.utils.defaultAbiCoder.encode(
        ['address[]'],
        [watchListAddresses],
      )
      await labm
        .connect(owner)
        .setWatchList(watchListAddresses, watchListMinBalances)
    })

    it('Should revert when paused', async () => {
      await labm.connect(owner).pause()
      const performTx = labm.connect(keeperRegistry).performUpkeep(validPayload)
      await expect(performTx).to.be.revertedWith(PAUSED_ERR)
    })

    it('Should fund the appropriate addresses', async () => {
      await lt.connect(owner).transfer(labm.address, tenLINK)
      await assertContractLinkBalances(
        zeroLINK,
        zeroLINK,
        zeroLINK,
        zeroLINK,
        zeroLINK,
      )
      const performTx = await labm
        .connect(keeperRegistry)
        .performUpkeep(validPayload, { gasLimit: 2_500_000 })
      await performTx.wait()
      await assertContractLinkBalances(
        twoLINK,
        twoLINK,
        twoLINK,
        twoLINK,
        twoLINK,
      )
    })

    it('Can handle MAX_PERFORM proxies within gas limit', async () => {
      // add MAX_PERFORM number of proxies
      const MAX_PERFORM = (await labm.MAX_PERFORM()).toNumber()
      const proxyAddresses = []
      const minBalances = []
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
        proxyAddresses.push(proxy.address)
        minBalances.push(oneLINK)
      }
      await labm.setWatchList(proxyAddresses, minBalances)
      expect(await labm.getWatchList()).to.deep.equalInAnyOrder([
        proxyAddresses,
        minBalances,
      ])
      // add funds
      const fundsNeeded = (await labm.getTopUpAmount()).mul(MAX_PERFORM)
      await lt.connect(owner).transfer(labm.address, fundsNeeded)
      // encode payload
      const payload = ethers.utils.defaultAbiCoder.encode(
        ['address[]'],
        [proxyAddresses],
      )
      // do the thing
      await labm
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
          await labm.connect(user).topUp([])
        }
      })
    })

    context('when paused', () => {
      it('Should be callable by no one', async () => {
        await labm.connect(owner).pause()
        const users = [owner, keeperRegistry, stranger]
        for (let idx = 0; idx < users.length; idx++) {
          const user = users[idx]
          const tx = labm.connect(user).topUp([])
          await expect(tx).to.be.revertedWith(PAUSED_ERR)
        }
      })
    })

    context('when fully funded', () => {
      beforeEach(async () => {
        await lt.connect(owner).transfer(labm.address, tenLINK)
        await assertContractLinkBalances(
          zeroLINK,
          zeroLINK,
          zeroLINK,
          zeroLINK,
          zeroLINK,
        )
      })

      it('Should fund the appropriate addresses', async () => {
        const tx = await labm.connect(keeperRegistry).topUp(watchListAddresses)
        await assertContractLinkBalances(
          twoLINK,
          twoLINK,
          twoLINK,
          twoLINK,
          twoLINK,
        )
        await expect(tx)
          .to.emit(labm, 'TopUpSucceeded')
          .withArgs(proxy1.address)
        await expect(tx)
          .to.emit(labm, 'TopUpSucceeded')
          .withArgs(proxy2.address)
        await expect(tx)
          .to.emit(labm, 'TopUpSucceeded')
          .withArgs(proxy3.address)
        await expect(tx)
          .to.emit(labm, 'TopUpSucceeded')
          .withArgs(directTarget1.address)
        await expect(tx)
          .to.emit(labm, 'TopUpSucceeded')
          .withArgs(directTarget2.address)
      })

      it('Should only fund the addresses provided', async () => {
        await labm
          .connect(keeperRegistry)
          .topUp([proxy1.address, directTarget1.address])
        await assertContractLinkBalances(
          twoLINK,
          zeroLINK,
          zeroLINK,
          twoLINK,
          zeroLINK,
        )
      })

      it('Should skip un-approved addresses', async () => {
        await labm
          .connect(owner)
          .setWatchList(
            [proxy1.address, directTarget1.address],
            [oneLINK, oneLINK],
          )
        const tx = await labm
          .connect(keeperRegistry)
          .topUp([
            proxy1.address,
            proxy2.address,
            proxy3.address,
            directTarget1.address,
            directTarget2.address,
          ])
        await assertContractLinkBalances(
          twoLINK,
          zeroLINK,
          zeroLINK,
          twoLINK,
          zeroLINK,
        )
        await expect(tx)
          .to.emit(labm, 'TopUpSucceeded')
          .withArgs(proxy1.address)
        await expect(tx)
          .to.emit(labm, 'TopUpSucceeded')
          .withArgs(directTarget1.address)
        await expect(tx).to.emit(labm, 'TopUpBlocked').withArgs(proxy2.address)
        await expect(tx).to.emit(labm, 'TopUpBlocked').withArgs(proxy3.address)
        await expect(tx)
          .to.emit(labm, 'TopUpBlocked')
          .withArgs(directTarget2.address)
      })

      it('Should skip an address if the proxy is invalid and it is not a direct target', async () => {
        await labm
          .connect(owner)
          .setWatchList([proxy1.address, proxy4.address], [oneLINK, oneLINK])
        const tx = await labm
          .connect(keeperRegistry)
          .topUp([proxy1.address, proxy4.address])
        await expect(tx)
          .to.emit(labm, 'TopUpSucceeded')
          .withArgs(proxy1.address)
        await expect(tx).to.emit(labm, 'TopUpBlocked').withArgs(proxy4.address)
      })

      it('Should skip an address if the aggregator is invalid', async () => {
        await proxy4.mock.aggregator.returns(aggregator4.address)
        await labm
          .connect(owner)
          .setWatchList([proxy1.address, proxy4.address], [oneLINK, oneLINK])
        const tx = await labm
          .connect(keeperRegistry)
          .topUp([proxy1.address, proxy4.address])
        await expect(tx)
          .to.emit(labm, 'TopUpSucceeded')
          .withArgs(proxy1.address)
        await expect(tx).to.emit(labm, 'TopUpBlocked').withArgs(proxy4.address)
      })

      it('Should skip an address if the aggregator has sufficient funding', async () => {
        await proxy4.mock.aggregator.returns(aggregator4.address)
        await aggregator4.mock.linkAvailableForPayment.returns(tenLINK)
        await labm
          .connect(owner)
          .setWatchList([proxy1.address, proxy4.address], [oneLINK, oneLINK])
        const tx = await labm
          .connect(keeperRegistry)
          .topUp([proxy1.address, proxy4.address])
        await expect(tx)
          .to.emit(labm, 'TopUpSucceeded')
          .withArgs(proxy1.address)
        await expect(tx).to.emit(labm, 'TopUpBlocked').withArgs(proxy4.address)
      })

      it('Should skip an address if the direct target has sufficient funding', async () => {
        await directTarget1.mock.linkAvailableForPayment.returns(tenLINK)
        await labm
          .connect(owner)
          .setWatchList(
            [proxy1.address, directTarget1.address],
            [oneLINK, oneLINK],
          )
        const tx = await labm
          .connect(keeperRegistry)
          .topUp([proxy1.address, directTarget1.address])
        await expect(tx)
          .to.emit(labm, 'TopUpSucceeded')
          .withArgs(proxy1.address)
        await expect(tx)
          .to.emit(labm, 'TopUpBlocked')
          .withArgs(directTarget1.address)
      })
    })

    context('when partially funded', () => {
      it('Should fund as many addresses as possible', async () => {
        await lt.connect(owner).transfer(
          labm.address,
          fiveLINK, // only enough LINK to fund 2 addresses
        )
        const tx = await labm.connect(keeperRegistry).topUp(watchListAddresses)
        await assertContractLinkBalances(
          twoLINK,
          twoLINK,
          zeroLINK,
          zeroLINK,
          zeroLINK,
        )
        await expect(tx)
          .to.emit(labm, 'TopUpSucceeded')
          .withArgs(proxy1.address)
        await expect(tx)
          .to.emit(labm, 'TopUpSucceeded')
          .withArgs(proxy2.address)
      })
    })
  })
})
