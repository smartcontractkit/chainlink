import { ethers } from 'hardhat'
import { assert, expect } from 'chai'
import { evmRevert } from '../test-helpers/matchers'
import { getUsers, Personas } from '../test-helpers/setup'
import { BigNumber, Signer, BigNumberish } from 'ethers'
import { LinkToken__factory as LinkTokenFactory } from '../../typechain/factories/LinkToken__factory'
import { KeeperRegistry11__factory as KeeperRegistryFactory } from '../../typechain/factories/KeeperRegistry11__factory'
import { MockV3Aggregator__factory as MockV3AggregatorFactory } from '../../typechain/factories/MockV3Aggregator__factory'
import { UpkeepMock__factory as UpkeepMockFactory } from '../../typechain/factories/UpkeepMock__factory'
import { UpkeepReverter__factory as UpkeepReverterFactory } from '../../typechain/factories/UpkeepReverter__factory'
import { KeeperRegistry11 as KeeperRegistry } from '../../typechain/KeeperRegistry11'
import { MockV3Aggregator } from '../../typechain/MockV3Aggregator'
import { LinkToken } from '../../typechain/LinkToken'
import { UpkeepMock } from '../../typechain/UpkeepMock'
import { toWei } from '../test-helpers/helpers'

async function getUpkeepID(tx: any) {
  const receipt = await tx.wait()
  return receipt.events[0].args.id
}

// -----------------------------------------------------------------------------------------------
// DEV: these *should* match the perform/check gas overhead values in the contract and on the node
const PERFORM_GAS_OVERHEAD = BigNumber.from(90000)
const CHECK_GAS_OVERHEAD = BigNumber.from(170000)
// -----------------------------------------------------------------------------------------------

// Smart contract factories
let linkTokenFactory: LinkTokenFactory
let mockV3AggregatorFactory: MockV3AggregatorFactory
let keeperRegistryFactory: KeeperRegistryFactory
let upkeepMockFactory: UpkeepMockFactory
let upkeepReverterFactory: UpkeepReverterFactory

let personas: Personas

before(async () => {
  personas = (await getUsers()).personas

  linkTokenFactory = await ethers.getContractFactory('LinkToken')
  // need full path because there are two contracts with name MockV3Aggregator
  mockV3AggregatorFactory = (await ethers.getContractFactory(
    'src/v0.7/tests/MockV3Aggregator.sol:MockV3Aggregator',
  )) as unknown as MockV3AggregatorFactory
  // @ts-ignore bug in autogen file
  keeperRegistryFactory = await ethers.getContractFactory('KeeperRegistry1_1')
  upkeepMockFactory = await ethers.getContractFactory('UpkeepMock')
  upkeepReverterFactory = await ethers.getContractFactory('UpkeepReverter')
})

describe('KeeperRegistry1_1', () => {
  const linkEth = BigNumber.from(300000000)
  const gasWei = BigNumber.from(100)
  const linkDivisibility = BigNumber.from('1000000000000000000')
  const executeGas = BigNumber.from('100000')
  const paymentPremiumBase = BigNumber.from('1000000000')
  const paymentPremiumPPB = BigNumber.from('250000000')
  const flatFeeMicroLink = BigNumber.from(0)
  const blockCountPerTurn = BigNumber.from(3)
  const emptyBytes = '0x00'
  const zeroAddress = ethers.constants.AddressZero
  const extraGas = BigNumber.from('250000')
  const registryGasOverhead = BigNumber.from('80000')
  const stalenessSeconds = BigNumber.from(43820)
  const gasCeilingMultiplier = BigNumber.from(1)
  const maxCheckGas = BigNumber.from(20000000)
  const fallbackGasPrice = BigNumber.from(200)
  const fallbackLinkPrice = BigNumber.from(200000000)

  let owner: Signer
  let keeper1: Signer
  let keeper2: Signer
  let keeper3: Signer
  let nonkeeper: Signer
  let admin: Signer
  let payee1: Signer
  let payee2: Signer
  let payee3: Signer

  let linkToken: LinkToken
  let linkEthFeed: MockV3Aggregator
  let gasPriceFeed: MockV3Aggregator
  let registry: KeeperRegistry
  let mock: UpkeepMock

  let id: BigNumber
  let keepers: string[]
  let payees: string[]

  beforeEach(async () => {
    owner = personas.Default
    keeper1 = personas.Carol
    keeper2 = personas.Eddy
    keeper3 = personas.Nancy
    nonkeeper = personas.Ned
    admin = personas.Neil
    payee1 = personas.Nelly
    payee2 = personas.Norbert
    payee3 = personas.Nick

    keepers = [
      await keeper1.getAddress(),
      await keeper2.getAddress(),
      await keeper3.getAddress(),
    ]
    payees = [
      await payee1.getAddress(),
      await payee2.getAddress(),
      await payee3.getAddress(),
    ]

    linkToken = await linkTokenFactory.connect(owner).deploy()
    gasPriceFeed = await mockV3AggregatorFactory
      .connect(owner)
      .deploy(0, gasWei)
    linkEthFeed = await mockV3AggregatorFactory
      .connect(owner)
      .deploy(9, linkEth)
    registry = await keeperRegistryFactory
      .connect(owner)
      .deploy(
        linkToken.address,
        linkEthFeed.address,
        gasPriceFeed.address,
        paymentPremiumPPB,
        flatFeeMicroLink,
        blockCountPerTurn,
        maxCheckGas,
        stalenessSeconds,
        gasCeilingMultiplier,
        fallbackGasPrice,
        fallbackLinkPrice,
      )

    mock = await upkeepMockFactory.deploy()
    await linkToken
      .connect(owner)
      .transfer(await keeper1.getAddress(), toWei('1000'))
    await linkToken
      .connect(owner)
      .transfer(await keeper2.getAddress(), toWei('1000'))
    await linkToken
      .connect(owner)
      .transfer(await keeper3.getAddress(), toWei('1000'))

    await registry.connect(owner).setKeepers(keepers, payees)
    const tx = await registry
      .connect(owner)
      .registerUpkeep(
        mock.address,
        executeGas,
        await admin.getAddress(),
        emptyBytes,
      )
    id = await getUpkeepID(tx)
  })

  const linkForGas = (
    upkeepGasSpent: BigNumberish,
    premiumPPB?: BigNumberish,
    flatFee?: BigNumberish,
  ) => {
    premiumPPB = premiumPPB === undefined ? paymentPremiumPPB : premiumPPB
    flatFee = flatFee === undefined ? flatFeeMicroLink : flatFee
    const gasSpent = registryGasOverhead.add(BigNumber.from(upkeepGasSpent))
    const base = gasWei.mul(gasSpent).mul(linkDivisibility).div(linkEth)
    const premium = base.mul(premiumPPB).div(paymentPremiumBase)
    const flatFeeJules = BigNumber.from(flatFee).mul('1000000000000')
    return base.add(premium).add(flatFeeJules)
  }

  describe('#setKeepers', () => {
    const IGNORE_ADDRESS = '0xFFfFfFffFFfffFFfFFfFFFFFffFFFffffFfFFFfF'
    it('reverts when not called by the owner', async () => {
      await evmRevert(
        registry.connect(keeper1).setKeepers([], []),
        'Only callable by owner',
      )
    })

    it('reverts when adding the same keeper twice', async () => {
      await evmRevert(
        registry
          .connect(owner)
          .setKeepers(
            [await keeper1.getAddress(), await keeper1.getAddress()],
            [await payee1.getAddress(), await payee1.getAddress()],
          ),
        'cannot add keeper twice',
      )
    })

    it('reverts with different numbers of keepers/payees', async () => {
      await evmRevert(
        registry
          .connect(owner)
          .setKeepers(
            [await keeper1.getAddress(), await keeper2.getAddress()],
            [await payee1.getAddress()],
          ),
        'address lists not the same length',
      )
      await evmRevert(
        registry
          .connect(owner)
          .setKeepers(
            [await keeper1.getAddress()],
            [await payee1.getAddress(), await payee2.getAddress()],
          ),
        'address lists not the same length',
      )
    })

    it('reverts if the payee is the zero address', async () => {
      await evmRevert(
        registry
          .connect(owner)
          .setKeepers(
            [await keeper1.getAddress(), await keeper2.getAddress()],
            [
              await payee1.getAddress(),
              '0x0000000000000000000000000000000000000000',
            ],
          ),
        'cannot set payee to the zero address',
      )
    })

    it('emits events for every keeper added and removed', async () => {
      const oldKeepers = [
        await keeper1.getAddress(),
        await keeper2.getAddress(),
      ]
      const oldPayees = [await payee1.getAddress(), await payee2.getAddress()]
      await registry.connect(owner).setKeepers(oldKeepers, oldPayees)
      assert.deepEqual(oldKeepers, await registry.getKeeperList())

      // remove keepers
      const newKeepers = [
        await keeper2.getAddress(),
        await keeper3.getAddress(),
      ]
      const newPayees = [await payee2.getAddress(), await payee3.getAddress()]
      const tx = await registry.connect(owner).setKeepers(newKeepers, newPayees)
      assert.deepEqual(newKeepers, await registry.getKeeperList())

      await expect(tx)
        .to.emit(registry, 'KeepersUpdated')
        .withArgs(newKeepers, newPayees)
    })

    it('updates the keeper to inactive when removed', async () => {
      await registry.connect(owner).setKeepers(keepers, payees)
      await registry
        .connect(owner)
        .setKeepers(
          [await keeper1.getAddress(), await keeper3.getAddress()],
          [await payee1.getAddress(), await payee3.getAddress()],
        )
      const added = await registry.getKeeperInfo(await keeper1.getAddress())
      assert.isTrue(added.active)
      const removed = await registry.getKeeperInfo(await keeper2.getAddress())
      assert.isFalse(removed.active)
    })

    it('does not change the payee if IGNORE_ADDRESS is used as payee', async () => {
      const oldKeepers = [
        await keeper1.getAddress(),
        await keeper2.getAddress(),
      ]
      const oldPayees = [await payee1.getAddress(), await payee2.getAddress()]
      await registry.connect(owner).setKeepers(oldKeepers, oldPayees)
      assert.deepEqual(oldKeepers, await registry.getKeeperList())

      const newKeepers = [
        await keeper2.getAddress(),
        await keeper3.getAddress(),
      ]
      const newPayees = [IGNORE_ADDRESS, await payee3.getAddress()]
      const tx = await registry.connect(owner).setKeepers(newKeepers, newPayees)
      assert.deepEqual(newKeepers, await registry.getKeeperList())

      const ignored = await registry.getKeeperInfo(await keeper2.getAddress())
      assert.equal(await payee2.getAddress(), ignored.payee)
      assert.equal(true, ignored.active)

      await expect(tx)
        .to.emit(registry, 'KeepersUpdated')
        .withArgs(newKeepers, newPayees)
    })

    it('reverts if the owner changes the payee', async () => {
      await registry.connect(owner).setKeepers(keepers, payees)
      await evmRevert(
        registry
          .connect(owner)
          .setKeepers(keepers, [
            await payee1.getAddress(),
            await payee2.getAddress(),
            await owner.getAddress(),
          ]),
        'cannot change payee',
      )
    })
  })

  describe('#registerUpkeep', () => {
    it('reverts if the target is not a contract', async () => {
      await evmRevert(
        registry
          .connect(owner)
          .registerUpkeep(
            zeroAddress,
            executeGas,
            await admin.getAddress(),
            emptyBytes,
          ),
        'target is not a contract',
      )
    })

    it('reverts if called by a non-owner', async () => {
      await evmRevert(
        registry
          .connect(keeper1)
          .registerUpkeep(
            mock.address,
            executeGas,
            await admin.getAddress(),
            emptyBytes,
          ),
        'Only callable by owner or registrar',
      )
    })

    it('reverts if execute gas is too low', async () => {
      await evmRevert(
        registry
          .connect(owner)
          .registerUpkeep(
            mock.address,
            2299,
            await admin.getAddress(),
            emptyBytes,
          ),
        'min gas is 2300',
      )
    })

    it('reverts if execute gas is too high', async () => {
      await evmRevert(
        registry
          .connect(owner)
          .registerUpkeep(
            mock.address,
            5000001,
            await admin.getAddress(),
            emptyBytes,
          ),
        'max gas is 5000000',
      )
    })

    it('creates a record of the registration', async () => {
      const tx = await registry
        .connect(owner)
        .registerUpkeep(
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
        )
      id = await getUpkeepID(tx)
      await expect(tx)
        .to.emit(registry, 'UpkeepRegistered')
        .withArgs(id, executeGas, await admin.getAddress())
      const registration = await registry.getUpkeep(id)
      assert.equal(mock.address, registration.target)
      assert.equal(0, registration.balance.toNumber())
      assert.equal(emptyBytes, registration.checkData)
      assert(registration.maxValidBlocknumber.eq('0xffffffffffffffff'))
    })
  })

  describe('#addFunds', () => {
    const amount = toWei('1')

    beforeEach(async () => {
      await linkToken.connect(keeper1).approve(registry.address, toWei('100'))
    })

    it('reverts if the registration does not exist', async () => {
      await evmRevert(
        registry.connect(keeper1).addFunds(id.add(1), amount),
        'upkeep must be active',
      )
    })

    it('adds to the balance of the registration', async () => {
      await registry.connect(keeper1).addFunds(id, amount)
      const registration = await registry.getUpkeep(id)
      assert.isTrue(amount.eq(registration.balance))
    })

    it('emits a log', async () => {
      const tx = await registry.connect(keeper1).addFunds(id, amount)
      await expect(tx)
        .to.emit(registry, 'FundsAdded')
        .withArgs(id, await keeper1.getAddress(), amount)
    })

    it('reverts if the upkeep is canceled', async () => {
      await registry.connect(admin).cancelUpkeep(id)
      await evmRevert(
        registry.connect(keeper1).addFunds(id, amount),
        'upkeep must be active',
      )
    })
  })

  describe('#checkUpkeep', () => {
    it('reverts if the upkeep is not funded', async () => {
      await mock.setCanPerform(true)
      await mock.setCanCheck(true)
      await evmRevert(
        registry
          .connect(zeroAddress)
          .callStatic.checkUpkeep(id, await keeper1.getAddress()),
        'insufficient funds',
      )
    })

    context('when the registration is funded', () => {
      beforeEach(async () => {
        await linkToken.connect(keeper1).approve(registry.address, toWei('100'))
        await registry.connect(keeper1).addFunds(id, toWei('100'))
      })

      it('reverts if executed', async () => {
        await mock.setCanPerform(true)
        await mock.setCanCheck(true)
        await evmRevert(
          registry.checkUpkeep(id, await keeper1.getAddress()),
          'only for simulated backend',
        )
      })

      it('reverts if the specified keeper is not valid', async () => {
        await mock.setCanPerform(true)
        await mock.setCanCheck(true)
        await evmRevert(
          registry.checkUpkeep(id, await owner.getAddress()),
          'only for simulated backend',
        )
      })

      context('and upkeep is not needed', () => {
        beforeEach(async () => {
          await mock.setCanCheck(false)
        })

        it('reverts', async () => {
          await evmRevert(
            registry
              .connect(zeroAddress)
              .callStatic.checkUpkeep(id, await keeper1.getAddress()),
            'upkeep not needed',
          )
        })
      })

      context('and the upkeep check fails', () => {
        beforeEach(async () => {
          const reverter = await upkeepReverterFactory.deploy()
          const tx = await registry
            .connect(owner)
            .registerUpkeep(
              reverter.address,
              2500000,
              await admin.getAddress(),
              emptyBytes,
            )
          id = await getUpkeepID(tx)
          await linkToken
            .connect(keeper1)
            .approve(registry.address, toWei('100'))
          await registry.connect(keeper1).addFunds(id, toWei('100'))
        })

        it('reverts', async () => {
          await evmRevert(
            registry
              .connect(zeroAddress)
              .callStatic.checkUpkeep(id, await keeper1.getAddress()),
            'call to check target failed',
          )
        })
      })

      context('and upkeep check simulations succeeds', () => {
        beforeEach(async () => {
          await mock.setCanCheck(true)
          await mock.setCanPerform(true)
        })

        context('and the registry is paused', () => {
          beforeEach(async () => {
            await registry.connect(owner).pause()
          })

          it('reverts', async () => {
            await evmRevert(
              registry
                .connect(zeroAddress)
                .callStatic.checkUpkeep(id, await keeper1.getAddress()),
              'Pausable: paused',
            )

            await registry.connect(owner).unpause()

            await registry
              .connect(zeroAddress)
              .callStatic.checkUpkeep(id, await keeper1.getAddress())
          })
        })

        it('returns true with pricing info if the target can execute', async () => {
          const newGasMultiplier = BigNumber.from(10)
          await registry
            .connect(owner)
            .setConfig(
              paymentPremiumPPB,
              flatFeeMicroLink,
              blockCountPerTurn,
              maxCheckGas,
              stalenessSeconds,
              newGasMultiplier,
              fallbackGasPrice,
              fallbackLinkPrice,
            )
          const response = await registry
            .connect(zeroAddress)
            .callStatic.checkUpkeep(id, await keeper1.getAddress())
          assert.isTrue(response.gasLimit.eq(executeGas))
          assert.isTrue(response.linkEth.eq(linkEth))
          assert.isTrue(
            response.adjustedGasWei.eq(gasWei.mul(newGasMultiplier)),
          )
          assert.isTrue(
            response.maxLinkPayment.eq(
              linkForGas(executeGas.toNumber()).mul(newGasMultiplier),
            ),
          )
        })

        it('has a large enough gas overhead to cover upkeeps that use all their gas [ @skip-coverage ]', async () => {
          await mock.setCheckGasToBurn(maxCheckGas)
          await mock.setPerformGasToBurn(executeGas)
          const gas = maxCheckGas
            .add(executeGas)
            .add(PERFORM_GAS_OVERHEAD)
            .add(CHECK_GAS_OVERHEAD)
          await registry
            .connect(zeroAddress)
            .callStatic.checkUpkeep(id, await keeper1.getAddress(), {
              gasLimit: gas,
            })
        })
      })
    })
  })

  describe('#performUpkeep', () => {
    let _lastKeeper = keeper1
    async function getPerformPaymentAmount() {
      _lastKeeper = _lastKeeper === keeper1 ? keeper2 : keeper1
      const before = (
        await registry.getKeeperInfo(await _lastKeeper.getAddress())
      ).balance
      await registry.connect(_lastKeeper).performUpkeep(id, '0x')
      const after = (
        await registry.getKeeperInfo(await _lastKeeper.getAddress())
      ).balance
      const difference = after.sub(before)
      return difference
    }

    it('reverts if the registration is not funded', async () => {
      await evmRevert(
        registry.connect(keeper2).performUpkeep(id, '0x'),
        'insufficient funds',
      )
    })

    context('when the registration is funded', () => {
      beforeEach(async () => {
        await linkToken.connect(owner).approve(registry.address, toWei('100'))
        await registry.connect(owner).addFunds(id, toWei('100'))
      })

      it('does not revert if the target cannot execute', async () => {
        const mockResponse = await mock
          .connect(zeroAddress)
          .callStatic.checkUpkeep('0x')
        assert.isFalse(mockResponse.callable)

        await registry.connect(keeper3).performUpkeep(id, '0x')
      })

      it('returns false if the target cannot execute', async () => {
        const mockResponse = await mock
          .connect(zeroAddress)
          .callStatic.checkUpkeep('0x')
        assert.isFalse(mockResponse.callable)

        assert.isFalse(
          await registry.connect(keeper1).callStatic.performUpkeep(id, '0x'),
        )
      })

      it('returns true if called', async () => {
        await mock.setCanPerform(true)

        const response = await registry
          .connect(keeper1)
          .callStatic.performUpkeep(id, '0x')
        assert.isTrue(response)
      })

      it('reverts if not enough gas supplied', async () => {
        await mock.setCanPerform(true)

        await evmRevert(
          registry
            .connect(keeper1)
            .performUpkeep(id, '0x', { gasLimit: BigNumber.from('120000') }),
        )
      })

      it('executes the data passed to the registry', async () => {
        await mock.setCanPerform(true)

        const performData = '0xc0ffeec0ffee'
        const tx = await registry
          .connect(keeper1)
          .performUpkeep(id, performData, { gasLimit: extraGas })
        const receipt = await tx.wait()
        const eventLog = receipt?.events

        assert.equal(eventLog?.length, 2)
        assert.equal(eventLog?.[1].event, 'UpkeepPerformed')
        assert.equal(eventLog?.[1].args?.[0].toNumber(), id.toNumber())
        assert.equal(eventLog?.[1].args?.[1], true)
        assert.equal(eventLog?.[1].args?.[2], await keeper1.getAddress())
        assert.isNotEmpty(eventLog?.[1].args?.[3])
        assert.equal(eventLog?.[1].args?.[4], performData)
      })

      it('updates payment balances', async () => {
        const keeperBefore = await registry.getKeeperInfo(
          await keeper1.getAddress(),
        )
        const registrationBefore = await registry.getUpkeep(id)
        const keeperLinkBefore = await linkToken.balanceOf(
          await keeper1.getAddress(),
        )
        const registryLinkBefore = await linkToken.balanceOf(registry.address)

        // Do the thing
        await registry.connect(keeper1).performUpkeep(id, '0x')

        const keeperAfter = await registry.getKeeperInfo(
          await keeper1.getAddress(),
        )
        const registrationAfter = await registry.getUpkeep(id)
        const keeperLinkAfter = await linkToken.balanceOf(
          await keeper1.getAddress(),
        )
        const registryLinkAfter = await linkToken.balanceOf(registry.address)

        assert.isTrue(keeperAfter.balance.gt(keeperBefore.balance))
        assert.isTrue(registrationBefore.balance.gt(registrationAfter.balance))
        assert.isTrue(keeperLinkAfter.eq(keeperLinkBefore))
        assert.isTrue(registryLinkBefore.eq(registryLinkAfter))
      })

      it('only pays for gas used [ @skip-coverage ]', async () => {
        const before = (
          await registry.getKeeperInfo(await keeper1.getAddress())
        ).balance
        const tx = await registry.connect(keeper1).performUpkeep(id, '0x')
        const receipt = await tx.wait()
        const after = (await registry.getKeeperInfo(await keeper1.getAddress()))
          .balance

        const max = linkForGas(executeGas.toNumber())
        const totalTx = linkForGas(receipt.gasUsed.toNumber())
        const difference = after.sub(before)
        assert.isTrue(max.gt(totalTx))
        assert.isTrue(totalTx.gt(difference))
        assert.isTrue(linkForGas(5700).lt(difference)) // exact number is flaky
        assert.isTrue(linkForGas(6000).gt(difference)) // instead test a range
      })

      it('only pays at a rate up to the gas ceiling [ @skip-coverage ]', async () => {
        const multiplier = BigNumber.from(10)
        const gasPrice = BigNumber.from('1000000000') // 10M x the gas feed's rate
        await registry
          .connect(owner)
          .setConfig(
            paymentPremiumPPB,
            flatFeeMicroLink,
            blockCountPerTurn,
            maxCheckGas,
            stalenessSeconds,
            multiplier,
            fallbackGasPrice,
            fallbackLinkPrice,
          )

        const before = (
          await registry.getKeeperInfo(await keeper1.getAddress())
        ).balance
        const tx = await registry
          .connect(keeper1)
          .performUpkeep(id, '0x', { gasPrice })
        const receipt = await tx.wait()
        const after = (await registry.getKeeperInfo(await keeper1.getAddress()))
          .balance

        const max = linkForGas(executeGas).mul(multiplier)
        const totalTx = linkForGas(receipt.gasUsed).mul(multiplier)
        const difference = after.sub(before)
        assert.isTrue(max.gt(totalTx))
        assert.isTrue(totalTx.gt(difference))
        assert.isTrue(linkForGas(5700).mul(multiplier).lt(difference))
        assert.isTrue(linkForGas(6000).mul(multiplier).gt(difference))
      })

      it('only pays as much as the node spent [ @skip-coverage ]', async () => {
        const multiplier = BigNumber.from(10)
        const gasPrice = BigNumber.from(200) // 2X the gas feed's rate
        const effectiveMultiplier = BigNumber.from(2)
        await registry
          .connect(owner)
          .setConfig(
            paymentPremiumPPB,
            flatFeeMicroLink,
            blockCountPerTurn,
            maxCheckGas,
            stalenessSeconds,
            multiplier,
            fallbackGasPrice,
            fallbackLinkPrice,
          )

        const before = (
          await registry.getKeeperInfo(await keeper1.getAddress())
        ).balance
        const tx = await registry
          .connect(keeper1)
          .performUpkeep(id, '0x', { gasPrice })
        const receipt = await tx.wait()
        const after = (await registry.getKeeperInfo(await keeper1.getAddress()))
          .balance

        const max = linkForGas(executeGas.toNumber()).mul(effectiveMultiplier)
        const totalTx = linkForGas(receipt.gasUsed).mul(effectiveMultiplier)
        const difference = after.sub(before)
        assert.isTrue(max.gt(totalTx))
        assert.isTrue(totalTx.gt(difference))
        assert.isTrue(linkForGas(5700).mul(effectiveMultiplier).lt(difference))
        assert.isTrue(linkForGas(6000).mul(effectiveMultiplier).gt(difference))
      })

      it('pays the caller even if the target function fails', async () => {
        const tx = await registry
          .connect(owner)
          .registerUpkeep(
            mock.address,
            executeGas,
            await admin.getAddress(),
            emptyBytes,
          )
        const id = await getUpkeepID(tx)
        await linkToken.connect(owner).approve(registry.address, toWei('100'))
        await registry.connect(owner).addFunds(id, toWei('100'))
        const keeperBalanceBefore = (
          await registry.getKeeperInfo(await keeper1.getAddress())
        ).balance

        // Do the thing
        await registry.connect(keeper1).performUpkeep(id, '0x')

        const keeperBalanceAfter = (
          await registry.getKeeperInfo(await keeper1.getAddress())
        ).balance
        assert.isTrue(keeperBalanceAfter.gt(keeperBalanceBefore))
      })

      it('reverts if called by a non-keeper', async () => {
        await evmRevert(
          registry.connect(nonkeeper).performUpkeep(id, '0x'),
          'only active keepers',
        )
      })

      it('reverts if the upkeep has been canceled', async () => {
        await mock.setCanPerform(true)

        await registry.connect(owner).cancelUpkeep(id)

        await evmRevert(
          registry.connect(keeper1).performUpkeep(id, '0x'),
          'invalid upkeep id',
        )
      })

      it('uses the fallback gas price if the feed price is stale [ @skip-coverage ]', async () => {
        const normalAmount = await getPerformPaymentAmount()
        const roundId = 99
        const answer = 100
        const updatedAt = 946684800 // New Years 2000 🥳
        const startedAt = 946684799
        await gasPriceFeed
          .connect(owner)
          .updateRoundData(roundId, answer, updatedAt, startedAt)
        const amountWithStaleFeed = await getPerformPaymentAmount()
        assert.isTrue(normalAmount.lt(amountWithStaleFeed))
      })

      it('uses the fallback gas price if the feed price is non-sensical [ @skip-coverage ]', async () => {
        const normalAmount = await getPerformPaymentAmount()
        const roundId = 99
        const updatedAt = Math.floor(Date.now() / 1000)
        const startedAt = 946684799
        await gasPriceFeed
          .connect(owner)
          .updateRoundData(roundId, -100, updatedAt, startedAt)
        const amountWithNegativeFeed = await getPerformPaymentAmount()
        await gasPriceFeed
          .connect(owner)
          .updateRoundData(roundId, 0, updatedAt, startedAt)
        const amountWithZeroFeed = await getPerformPaymentAmount()
        assert.isTrue(normalAmount.lt(amountWithNegativeFeed))
        assert.isTrue(normalAmount.lt(amountWithZeroFeed))
      })

      it('uses the fallback if the link price feed is stale', async () => {
        const normalAmount = await getPerformPaymentAmount()
        const roundId = 99
        const answer = 100
        const updatedAt = 946684800 // New Years 2000 🥳
        const startedAt = 946684799
        await linkEthFeed
          .connect(owner)
          .updateRoundData(roundId, answer, updatedAt, startedAt)
        const amountWithStaleFeed = await getPerformPaymentAmount()
        assert.isTrue(normalAmount.lt(amountWithStaleFeed))
      })

      it('uses the fallback link price if the feed price is non-sensical', async () => {
        const normalAmount = await getPerformPaymentAmount()
        const roundId = 99
        const updatedAt = Math.floor(Date.now() / 1000)
        const startedAt = 946684799
        await linkEthFeed
          .connect(owner)
          .updateRoundData(roundId, -100, updatedAt, startedAt)
        const amountWithNegativeFeed = await getPerformPaymentAmount()
        await linkEthFeed
          .connect(owner)
          .updateRoundData(roundId, 0, updatedAt, startedAt)
        const amountWithZeroFeed = await getPerformPaymentAmount()
        assert.isTrue(normalAmount.lt(amountWithNegativeFeed))
        assert.isTrue(normalAmount.lt(amountWithZeroFeed))
      })

      it('reverts if the same caller calls twice in a row', async () => {
        await registry.connect(keeper1).performUpkeep(id, '0x')
        await evmRevert(
          registry.connect(keeper1).performUpkeep(id, '0x'),
          'keepers must take turns',
        )
        await registry.connect(keeper2).performUpkeep(id, '0x')
        await evmRevert(
          registry.connect(keeper2).performUpkeep(id, '0x'),
          'keepers must take turns',
        )
        await registry.connect(keeper1).performUpkeep(id, '0x')
      })

      it('has a large enough gas overhead to cover upkeeps that use all their gas [ @skip-coverage ]', async () => {
        await mock.setPerformGasToBurn(executeGas)
        await mock.setCanPerform(true)
        const gas = executeGas.add(PERFORM_GAS_OVERHEAD)
        const performData = '0xc0ffeec0ffee'
        const tx = await registry
          .connect(keeper1)
          .performUpkeep(id, performData, { gasLimit: gas })
        const receipt = await tx.wait()
        const eventLog = receipt?.events

        assert.equal(eventLog?.length, 2)
        assert.equal(eventLog?.[1].event, 'UpkeepPerformed')
        assert.equal(eventLog?.[1].args?.[0].toNumber(), id.toNumber())
        assert.equal(eventLog?.[1].args?.[1], true)
        assert.equal(eventLog?.[1].args?.[2], await keeper1.getAddress())
        assert.isNotEmpty(eventLog?.[1].args?.[3])
        assert.equal(eventLog?.[1].args?.[4], performData)
      })
    })
  })

  describe('#withdrawFunds', () => {
    beforeEach(async () => {
      await linkToken.connect(keeper1).approve(registry.address, toWei('100'))
      await registry.connect(keeper1).addFunds(id, toWei('1'))
    })

    it('reverts if called by anyone but the admin', async () => {
      await evmRevert(
        registry
          .connect(owner)
          .withdrawFunds(id.add(1).toNumber(), await payee1.getAddress()),
        'only callable by admin',
      )
    })

    it('reverts if called on an uncanceled upkeep', async () => {
      await evmRevert(
        registry.connect(admin).withdrawFunds(id, await payee1.getAddress()),
        'upkeep must be canceled',
      )
    })

    it('reverts if called with the 0 address', async () => {
      await evmRevert(
        registry.connect(admin).withdrawFunds(id, zeroAddress),
        'cannot send to zero address',
      )
    })

    describe('after the registration is cancelled', () => {
      beforeEach(async () => {
        await registry.connect(owner).cancelUpkeep(id)
      })

      it('moves the funds out and updates the balance', async () => {
        const payee1Before = await linkToken.balanceOf(
          await payee1.getAddress(),
        )
        const registryBefore = await linkToken.balanceOf(registry.address)

        let registration = await registry.getUpkeep(id)
        assert.isTrue(toWei('1').eq(registration.balance))

        await registry
          .connect(admin)
          .withdrawFunds(id, await payee1.getAddress())

        const payee1After = await linkToken.balanceOf(await payee1.getAddress())
        const registryAfter = await linkToken.balanceOf(registry.address)

        assert.isTrue(payee1Before.add(toWei('1')).eq(payee1After))
        assert.isTrue(registryBefore.sub(toWei('1')).eq(registryAfter))

        registration = await registry.getUpkeep(id)
        assert.equal(0, registration.balance.toNumber())
      })
    })
  })

  describe('#cancelUpkeep', () => {
    it('reverts if the ID is not valid', async () => {
      await evmRevert(
        registry.connect(owner).cancelUpkeep(id.add(1).toNumber()),
        'too late to cancel upkeep',
      )
    })

    it('reverts if called by a non-owner/non-admin', async () => {
      await evmRevert(
        registry.connect(keeper1).cancelUpkeep(id),
        'only owner or admin',
      )
    })

    describe('when called by the owner', async () => {
      it('sets the registration to invalid immediately', async () => {
        const tx = await registry.connect(owner).cancelUpkeep(id)
        const receipt = await tx.wait()
        const registration = await registry.getUpkeep(id)
        assert.equal(
          registration.maxValidBlocknumber.toNumber(),
          receipt.blockNumber,
        )
      })

      it('emits an event', async () => {
        const tx = await registry.connect(owner).cancelUpkeep(id)
        const receipt = await tx.wait()
        await expect(tx)
          .to.emit(registry, 'UpkeepCanceled')
          .withArgs(id, BigNumber.from(receipt.blockNumber))
      })

      it('updates the canceled registrations list', async () => {
        let canceled = await registry.callStatic.getCanceledUpkeepList()
        assert.deepEqual([], canceled)

        await registry.connect(owner).cancelUpkeep(id)

        canceled = await registry.callStatic.getCanceledUpkeepList()
        assert.deepEqual([id], canceled)
      })

      it('immediately prevents upkeep', async () => {
        await registry.connect(owner).cancelUpkeep(id)

        await evmRevert(
          registry.connect(keeper2).performUpkeep(id, '0x'),
          'invalid upkeep id',
        )
      })

      it('does not revert if reverts if called multiple times', async () => {
        await registry.connect(owner).cancelUpkeep(id)
        await evmRevert(
          registry.connect(owner).cancelUpkeep(id),
          'too late to cancel upkeep',
        )
      })

      describe('when called by the owner when the admin has just canceled', () => {
        let oldExpiration: BigNumber

        beforeEach(async () => {
          await registry.connect(admin).cancelUpkeep(id)
          const registration = await registry.getUpkeep(id)
          oldExpiration = registration.maxValidBlocknumber
        })

        it('allows the owner to cancel it more quickly', async () => {
          await registry.connect(owner).cancelUpkeep(id)

          const registration = await registry.getUpkeep(id)
          const newExpiration = registration.maxValidBlocknumber
          assert.isTrue(newExpiration.lt(oldExpiration))
        })
      })
    })

    describe('when called by the admin', async () => {
      const delay = 50

      it('sets the registration to invalid in 50 blocks', async () => {
        const tx = await registry.connect(admin).cancelUpkeep(id)
        const receipt = await tx.wait()
        const registration = await registry.getUpkeep(id)
        assert.equal(
          registration.maxValidBlocknumber.toNumber(),
          receipt.blockNumber + 50,
        )
      })

      it('emits an event', async () => {
        const tx = await registry.connect(admin).cancelUpkeep(id)
        const receipt = await tx.wait()
        await expect(tx)
          .to.emit(registry, 'UpkeepCanceled')
          .withArgs(id, BigNumber.from(receipt.blockNumber + delay))
      })

      it('updates the canceled registrations list', async () => {
        let canceled = await registry.callStatic.getCanceledUpkeepList()
        assert.deepEqual([], canceled)

        await registry.connect(admin).cancelUpkeep(id)

        canceled = await registry.callStatic.getCanceledUpkeepList()
        assert.deepEqual([id], canceled)
      })

      it('immediately prevents upkeep', async () => {
        await linkToken.connect(owner).approve(registry.address, toWei('100'))
        await registry.connect(owner).addFunds(id, toWei('100'))
        await registry.connect(admin).cancelUpkeep(id)
        await registry.connect(keeper2).performUpkeep(id, '0x') // still works

        for (let i = 0; i < delay; i++) {
          await ethers.provider.send('evm_mine', [])
        }

        await evmRevert(
          registry.connect(keeper2).performUpkeep(id, '0x'),
          'invalid upkeep id',
        )
      })

      it('reverts if called again by the admin', async () => {
        await registry.connect(admin).cancelUpkeep(id)

        await evmRevert(
          registry.connect(admin).cancelUpkeep(id),
          'too late to cancel upkeep',
        )
      })

      it('does not revert or double add the cancellation record if called by the owner immediately after', async () => {
        await registry.connect(admin).cancelUpkeep(id)

        await registry.connect(owner).cancelUpkeep(id)

        const canceled = await registry.callStatic.getCanceledUpkeepList()
        assert.deepEqual([id], canceled)
      })

      it('reverts if called by the owner after the timeout', async () => {
        await registry.connect(admin).cancelUpkeep(id)

        for (let i = 0; i < delay; i++) {
          await ethers.provider.send('evm_mine', [])
        }

        await evmRevert(
          registry.connect(owner).cancelUpkeep(id),
          'too late to cancel upkeep',
        )
      })
    })
  })

  describe('#withdrawPayment', () => {
    beforeEach(async () => {
      await linkToken.connect(owner).approve(registry.address, toWei('100'))
      await registry.connect(owner).addFunds(id, toWei('100'))
      await registry.connect(keeper1).performUpkeep(id, '0x')
    })

    it('reverts if called by anyone but the payee', async () => {
      await evmRevert(
        registry
          .connect(payee2)
          .withdrawPayment(
            await keeper1.getAddress(),
            await nonkeeper.getAddress(),
          ),
        'only callable by payee',
      )
    })

    it('reverts if called with the 0 address', async () => {
      await evmRevert(
        registry
          .connect(payee2)
          .withdrawPayment(await keeper1.getAddress(), zeroAddress),
        'cannot send to zero address',
      )
    })

    it('updates the balances', async () => {
      const to = await nonkeeper.getAddress()
      const keeperBefore = (
        await registry.getKeeperInfo(await keeper1.getAddress())
      ).balance
      const registrationBefore = (await registry.getUpkeep(id)).balance
      const toLinkBefore = await linkToken.balanceOf(to)
      const registryLinkBefore = await linkToken.balanceOf(registry.address)

      //// Do the thing
      await registry
        .connect(payee1)
        .withdrawPayment(await keeper1.getAddress(), to)

      const keeperAfter = (
        await registry.getKeeperInfo(await keeper1.getAddress())
      ).balance
      const registrationAfter = (await registry.getUpkeep(id)).balance
      const toLinkAfter = await linkToken.balanceOf(to)
      const registryLinkAfter = await linkToken.balanceOf(registry.address)

      assert.isTrue(keeperAfter.eq(BigNumber.from(0)))
      assert.isTrue(registrationBefore.eq(registrationAfter))
      assert.isTrue(toLinkBefore.add(keeperBefore).eq(toLinkAfter))
      assert.isTrue(registryLinkBefore.sub(keeperBefore).eq(registryLinkAfter))
    })

    it('emits a log announcing the withdrawal', async () => {
      const balance = (await registry.getKeeperInfo(await keeper1.getAddress()))
        .balance
      const tx = await registry
        .connect(payee1)
        .withdrawPayment(
          await keeper1.getAddress(),
          await nonkeeper.getAddress(),
        )
      await expect(tx)
        .to.emit(registry, 'PaymentWithdrawn')
        .withArgs(
          await keeper1.getAddress(),
          balance,
          await nonkeeper.getAddress(),
          await payee1.getAddress(),
        )
    })
  })

  describe('#transferPayeeship', () => {
    it('reverts when called by anyone but the current payee', async () => {
      await evmRevert(
        registry
          .connect(payee2)
          .transferPayeeship(
            await keeper1.getAddress(),
            await payee2.getAddress(),
          ),
        'only callable by payee',
      )
    })

    it('reverts when transferring to self', async () => {
      await evmRevert(
        registry
          .connect(payee1)
          .transferPayeeship(
            await keeper1.getAddress(),
            await payee1.getAddress(),
          ),
        'cannot transfer to self',
      )
    })

    it('does not change the payee', async () => {
      await registry
        .connect(payee1)
        .transferPayeeship(
          await keeper1.getAddress(),
          await payee2.getAddress(),
        )

      const info = await registry.getKeeperInfo(await keeper1.getAddress())
      assert.equal(await payee1.getAddress(), info.payee)
    })

    it('emits an event announcing the new payee', async () => {
      const tx = await registry
        .connect(payee1)
        .transferPayeeship(
          await keeper1.getAddress(),
          await payee2.getAddress(),
        )
      await expect(tx)
        .to.emit(registry, 'PayeeshipTransferRequested')
        .withArgs(
          await keeper1.getAddress(),
          await payee1.getAddress(),
          await payee2.getAddress(),
        )
    })

    it('does not emit an event when called with the same proposal', async () => {
      await registry
        .connect(payee1)
        .transferPayeeship(
          await keeper1.getAddress(),
          await payee2.getAddress(),
        )

      const tx = await registry
        .connect(payee1)
        .transferPayeeship(
          await keeper1.getAddress(),
          await payee2.getAddress(),
        )
      const receipt = await tx.wait()
      assert.equal(0, receipt.logs.length)
    })
  })

  describe('#acceptPayeeship', () => {
    beforeEach(async () => {
      await registry
        .connect(payee1)
        .transferPayeeship(
          await keeper1.getAddress(),
          await payee2.getAddress(),
        )
    })

    it('reverts when called by anyone but the proposed payee', async () => {
      await evmRevert(
        registry.connect(payee1).acceptPayeeship(await keeper1.getAddress()),
        'only callable by proposed payee',
      )
    })

    it('emits an event announcing the new payee', async () => {
      const tx = await registry
        .connect(payee2)
        .acceptPayeeship(await keeper1.getAddress())
      await expect(tx)
        .to.emit(registry, 'PayeeshipTransferred')
        .withArgs(
          await keeper1.getAddress(),
          await payee1.getAddress(),
          await payee2.getAddress(),
        )
    })

    it('does change the payee', async () => {
      await registry.connect(payee2).acceptPayeeship(await keeper1.getAddress())

      const info = await registry.getKeeperInfo(await keeper1.getAddress())
      assert.equal(await payee2.getAddress(), info.payee)
    })
  })

  describe('#setConfig', () => {
    const payment = BigNumber.from(1)
    const flatFee = BigNumber.from(2)
    const checks = BigNumber.from(3)
    const staleness = BigNumber.from(4)
    const ceiling = BigNumber.from(5)
    const maxGas = BigNumber.from(6)
    const fbGasEth = BigNumber.from(7)
    const fbLinkEth = BigNumber.from(8)

    it('reverts when called by anyone but the proposed owner', async () => {
      await evmRevert(
        registry
          .connect(payee1)
          .setConfig(
            payment,
            flatFee,
            checks,
            maxGas,
            staleness,
            gasCeilingMultiplier,
            fbGasEth,
            fbLinkEth,
          ),
        'Only callable by owner',
      )
    })

    it('updates the config', async () => {
      const old = await registry.getConfig()
      const oldFlatFee = await registry.getFlatFee()
      assert.isTrue(paymentPremiumPPB.eq(old.paymentPremiumPPB))
      assert.isTrue(flatFeeMicroLink.eq(oldFlatFee))
      assert.isTrue(blockCountPerTurn.eq(old.blockCountPerTurn))
      assert.isTrue(stalenessSeconds.eq(old.stalenessSeconds))
      assert.isTrue(gasCeilingMultiplier.eq(old.gasCeilingMultiplier))

      await registry
        .connect(owner)
        .setConfig(
          payment,
          flatFee,
          checks,
          maxGas,
          staleness,
          ceiling,
          fbGasEth,
          fbLinkEth,
        )

      const updated = await registry.getConfig()
      const newFlatFee = await registry.getFlatFee()
      assert.equal(updated.paymentPremiumPPB, payment.toNumber())
      assert.equal(newFlatFee, flatFee.toNumber())
      assert.equal(updated.blockCountPerTurn, checks.toNumber())
      assert.equal(updated.stalenessSeconds, staleness.toNumber())
      assert.equal(updated.gasCeilingMultiplier, ceiling.toNumber())
      assert.equal(updated.checkGasLimit, maxGas.toNumber())
      assert.equal(updated.fallbackGasPrice.toNumber(), fbGasEth.toNumber())
      assert.equal(updated.fallbackLinkPrice.toNumber(), fbLinkEth.toNumber())
    })

    it('emits an event', async () => {
      const tx = await registry
        .connect(owner)
        .setConfig(
          payment,
          flatFee,
          checks,
          maxGas,
          staleness,
          ceiling,
          fbGasEth,
          fbLinkEth,
        )
      await expect(tx)
        .to.emit(registry, 'ConfigSet')
        .withArgs(
          payment,
          checks,
          maxGas,
          staleness,
          ceiling,
          fbGasEth,
          fbLinkEth,
        )
    })
  })

  describe('#onTokenTransfer', () => {
    const amount = toWei('1')

    it('reverts if not called by the LINK token', async () => {
      const data = ethers.utils.defaultAbiCoder.encode(
        ['uint256'],
        [id.toNumber().toString()],
      )

      await evmRevert(
        registry
          .connect(keeper1)
          .onTokenTransfer(await keeper1.getAddress(), amount, data),
        'only callable through LINK',
      )
    })

    it('reverts if not called with more or less than 32 bytes', async () => {
      const longData = ethers.utils.defaultAbiCoder.encode(
        ['uint256', 'uint256'],
        ['33', '34'],
      )
      const shortData = '0x12345678'

      await evmRevert(
        linkToken
          .connect(owner)
          .transferAndCall(registry.address, amount, longData),
      )
      await evmRevert(
        linkToken
          .connect(owner)
          .transferAndCall(registry.address, amount, shortData),
      )
    })

    it('reverts if the upkeep is canceled', async () => {
      await registry.connect(admin).cancelUpkeep(id)
      await evmRevert(
        registry.connect(keeper1).addFunds(id, amount),
        'upkeep must be active',
      )
    })

    it('updates the funds of the job id passed', async () => {
      const data = ethers.utils.defaultAbiCoder.encode(
        ['uint256'],
        [id.toNumber().toString()],
      )

      const before = (await registry.getUpkeep(id)).balance
      await linkToken
        .connect(owner)
        .transferAndCall(registry.address, amount, data)
      const after = (await registry.getUpkeep(id)).balance

      assert.isTrue(before.add(amount).eq(after))
    })
  })

  describe('#recoverFunds', () => {
    const sent = toWei('7')

    beforeEach(async () => {
      await linkToken.connect(keeper1).approve(registry.address, toWei('100'))

      // add funds to upkeep 1 and perform and withdraw some payment
      const tx = await registry
        .connect(owner)
        .registerUpkeep(
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
        )
      const id1 = await getUpkeepID(tx)
      await registry.connect(keeper1).addFunds(id1, toWei('5'))
      await registry.connect(keeper1).performUpkeep(id1, '0x')
      await registry.connect(keeper2).performUpkeep(id1, '0x')
      await registry.connect(keeper3).performUpkeep(id1, '0x')
      await registry
        .connect(payee1)
        .withdrawPayment(
          await keeper1.getAddress(),
          await nonkeeper.getAddress(),
        )

      // transfer funds directly to the registry
      await linkToken.connect(keeper1).transfer(registry.address, sent)

      // add funds to upkeep 2 and perform and withdraw some payment
      const tx2 = await registry
        .connect(owner)
        .registerUpkeep(
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
        )
      const id2 = await getUpkeepID(tx2)
      await registry.connect(keeper1).addFunds(id2, toWei('5'))
      await registry.connect(keeper1).performUpkeep(id2, '0x')
      await registry.connect(keeper2).performUpkeep(id2, '0x')
      await registry.connect(keeper3).performUpkeep(id2, '0x')
      await registry
        .connect(payee2)
        .withdrawPayment(
          await keeper2.getAddress(),
          await nonkeeper.getAddress(),
        )

      // transfer funds using onTokenTransfer
      const data = ethers.utils.defaultAbiCoder.encode(
        ['uint256'],
        [id2.toNumber().toString()],
      )
      await linkToken
        .connect(owner)
        .transferAndCall(registry.address, toWei('1'), data)

      // remove a keeper
      await registry
        .connect(owner)
        .setKeepers(
          [await keeper1.getAddress(), await keeper2.getAddress()],
          [await payee1.getAddress(), await payee2.getAddress()],
        )

      // withdraw some funds
      await registry.connect(owner).cancelUpkeep(id1)
      await registry.connect(admin).withdrawFunds(id1, await admin.getAddress())
    })

    it('reverts if not called by owner', async () => {
      await evmRevert(
        registry.connect(keeper1).recoverFunds(),
        'Only callable by owner',
      )
    })

    it('allows any funds that have been accidentally transfered to be moved', async () => {
      const balanceBefore = await linkToken.balanceOf(registry.address)

      await linkToken.balanceOf(registry.address)

      await registry.connect(owner).recoverFunds()
      const balanceAfter = await linkToken.balanceOf(registry.address)
      assert.isTrue(balanceBefore.eq(balanceAfter.add(sent)))
    })
  })

  describe('#pause', () => {
    it('reverts if called by a non-owner', async () => {
      await evmRevert(
        registry.connect(keeper1).pause(),
        'Only callable by owner',
      )
    })

    it('marks the contract as paused', async () => {
      assert.isFalse(await registry.paused())

      await registry.connect(owner).pause()

      assert.isTrue(await registry.paused())
    })
  })

  describe('#unpause', () => {
    beforeEach(async () => {
      await registry.connect(owner).pause()
    })

    it('reverts if called by a non-owner', async () => {
      await evmRevert(
        registry.connect(keeper1).unpause(),
        'Only callable by owner',
      )
    })

    it('marks the contract as not paused', async () => {
      assert.isTrue(await registry.paused())

      await registry.connect(owner).unpause()

      assert.isFalse(await registry.paused())
    })
  })

  describe('#getMaxPaymentForGas', () => {
    const gasAmounts = [100000, 10000000]
    const premiums = [0, 250000000]
    const flatFees = [0, 1000000]
    it('calculates the max fee approptiately', async () => {
      for (let idx = 0; idx < gasAmounts.length; idx++) {
        const gas = gasAmounts[idx]
        for (let jdx = 0; jdx < premiums.length; jdx++) {
          const premium = premiums[jdx]
          for (let kdx = 0; kdx < flatFees.length; kdx++) {
            const flatFee = flatFees[kdx]
            await registry
              .connect(owner)
              .setConfig(
                premium,
                flatFee,
                blockCountPerTurn,
                maxCheckGas,
                stalenessSeconds,
                gasCeilingMultiplier,
                fallbackGasPrice,
                fallbackLinkPrice,
              )
            const price = await registry.getMaxPaymentForGas(gas)
            expect(price).to.equal(linkForGas(gas, premium, flatFee))
          }
        }
      }
    })
  })

  describe('#checkUpkeep / #performUpkeep', () => {
    const performData = '0xc0ffeec0ffee'
    const multiplier = BigNumber.from(10)
    const flatFee = BigNumber.from('100000') //0.1 LINK
    const callGasPrice = 1

    it('uses the same minimum balance calculation [ @skip-coverage ]', async () => {
      await registry
        .connect(owner)
        .setConfig(
          paymentPremiumPPB,
          flatFee,
          blockCountPerTurn,
          maxCheckGas,
          stalenessSeconds,
          multiplier,
          fallbackGasPrice,
          fallbackLinkPrice,
        )
      await linkToken.connect(owner).approve(registry.address, toWei('100'))

      const tx1 = await registry
        .connect(owner)
        .registerUpkeep(
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
        )
      const upkeepID1 = await getUpkeepID(tx1)
      const tx2 = await registry
        .connect(owner)
        .registerUpkeep(
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
        )
      const upkeepID2 = await getUpkeepID(tx2)
      await mock.setCanCheck(true)
      await mock.setCanPerform(true)
      // upkeep 1 is underfunded, 2 is funded
      const minBalance1 = (await registry.getMaxPaymentForGas(executeGas)).sub(
        1,
      )
      const minBalance2 = await registry.getMaxPaymentForGas(executeGas)
      await registry.connect(owner).addFunds(upkeepID1, minBalance1)
      await registry.connect(owner).addFunds(upkeepID2, minBalance2)
      // upkeep 1 check should revert, 2 should succeed
      await evmRevert(
        registry
          .connect(zeroAddress)
          .callStatic.checkUpkeep(upkeepID1, await keeper1.getAddress(), {
            gasPrice: callGasPrice,
          }),
      )
      await registry
        .connect(zeroAddress)
        .callStatic.checkUpkeep(upkeepID2, await keeper1.getAddress(), {
          gasPrice: callGasPrice,
        })
      // upkeep 1 perform should revert, 2 should succeed
      await evmRevert(
        registry
          .connect(keeper1)
          .performUpkeep(upkeepID1, performData, { gasLimit: extraGas }),
        'insufficient funds',
      )
      await registry
        .connect(keeper1)
        .performUpkeep(upkeepID2, performData, { gasLimit: extraGas })
    })
  })

  describe('#getMinBalanceForUpkeep / #checkUpkeep', () => {
    it('calculates the minimum balance appropriately', async () => {
      const oneWei = BigNumber.from('1')
      await linkToken.connect(keeper1).approve(registry.address, toWei('100'))
      await mock.setCanCheck(true)
      await mock.setCanPerform(true)
      const minBalance = await registry.getMinBalanceForUpkeep(id)
      const tooLow = minBalance.sub(oneWei)
      await registry.connect(keeper1).addFunds(id, tooLow)
      await evmRevert(
        registry
          .connect(zeroAddress)
          .callStatic.checkUpkeep(id, await keeper1.getAddress()),
        'insufficient funds',
      )
      await registry.connect(keeper1).addFunds(id, oneWei)
      await registry
        .connect(zeroAddress)
        .callStatic.checkUpkeep(id, await keeper1.getAddress())
    })
  })
})
