import { ethers } from 'hardhat'
import { assert, expect } from 'chai'
import { evmRevert } from '../test-helpers/matchers'
import { getUsers, Personas } from '../test-helpers/setup'
import { BigNumber, BigNumberish, Signer } from 'ethers'
import { LinkToken__factory as LinkTokenFactory } from '../../typechain/factories/LinkToken__factory'
import { MockV3Aggregator__factory as MockV3AggregatorFactory } from '../../typechain/factories/MockV3Aggregator__factory'
import { UpkeepMock__factory as UpkeepMockFactory } from '../../typechain/factories/UpkeepMock__factory'
import { UpkeepReverter__factory as UpkeepReverterFactory } from '../../typechain/factories/UpkeepReverter__factory'
import { UpkeepAutoFunder__factory as UpkeepAutoFunderFactory } from '../../typechain/factories/UpkeepAutoFunder__factory'
import { UpkeepTranscoder__factory as UpkeepTranscoderFactory } from '../../typechain/factories/UpkeepTranscoder__factory'
import { KeeperRegistry__factory as KeeperRegistryFactory } from '../../typechain/factories/KeeperRegistry__factory'
import { KeeperRegistry } from '../../typechain/KeeperRegistry'

import { MockV3Aggregator } from '../../typechain/MockV3Aggregator'
import { LinkToken } from '../../typechain/LinkToken'
import { UpkeepMock } from '../../typechain/UpkeepMock'
import { UpkeepTranscoder } from '../../typechain'
import { toWei } from '../test-helpers/helpers'

async function getUpkeepID(tx: any) {
  const receipt = await tx.wait()
  return receipt.events[0].args.id
}

function randomAddress() {
  return ethers.Wallet.createRandom().address
}

// -----------------------------------------------------------------------------------------------
// DEV: these *should* match the perform/check gas overhead values in the contract and on the node
const PERFORM_GAS_OVERHEAD = BigNumber.from(160000)
const CHECK_GAS_OVERHEAD = BigNumber.from(170000)
// -----------------------------------------------------------------------------------------------

// Smart contract factories
let linkTokenFactory: LinkTokenFactory
let mockV3AggregatorFactory: MockV3AggregatorFactory
let keeperRegistryFactory: KeeperRegistryFactory
let upkeepMockFactory: UpkeepMockFactory
let upkeepReverterFactory: UpkeepReverterFactory
let upkeepAutoFunderFactory: UpkeepAutoFunderFactory
let upkeepTranscoderFactory: UpkeepTranscoderFactory
let personas: Personas

before(async () => {
  personas = (await getUsers()).personas

  linkTokenFactory = await ethers.getContractFactory('LinkToken')
  // need full path because there are two contracts with name MockV3Aggregator
  mockV3AggregatorFactory = (await ethers.getContractFactory(
    'src/v0.8/tests/MockV3Aggregator.sol:MockV3Aggregator',
  )) as unknown as MockV3AggregatorFactory
  keeperRegistryFactory = await ethers.getContractFactory('KeeperRegistry')
  upkeepMockFactory = await ethers.getContractFactory('UpkeepMock')
  upkeepReverterFactory = await ethers.getContractFactory('UpkeepReverter')
  upkeepAutoFunderFactory = await ethers.getContractFactory('UpkeepAutoFunder')
  upkeepTranscoderFactory = await ethers.getContractFactory('UpkeepTranscoder')
})

describe('KeeperRegistry', () => {
  const linkEth = BigNumber.from(300000000)
  const gasWei = BigNumber.from(100)
  const linkDivisibility = BigNumber.from('1000000000000000000')
  const executeGas = BigNumber.from('100000')
  const paymentPremiumBase = BigNumber.from('1000000000')
  const paymentPremiumPPB = BigNumber.from('250000000')
  const flatFeeMicroLink = BigNumber.from(0)
  const blockCountPerTurn = BigNumber.from(3)
  const emptyBytes = '0x00'
  const randomBytes = '0x1234abcd'
  const zeroAddress = ethers.constants.AddressZero
  const extraGas = BigNumber.from('250000')
  const registryGasOverhead = BigNumber.from('80000')
  const stalenessSeconds = BigNumber.from(43820)
  const gasCeilingMultiplier = BigNumber.from(1)
  const checkGasLimit = BigNumber.from(20000000)
  const fallbackGasPrice = BigNumber.from(200)
  const fallbackLinkPrice = BigNumber.from(200000000)
  const maxPerformGas = BigNumber.from(5000000)
  const minUpkeepSpend = BigNumber.from(0)

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
  let registry2: KeeperRegistry
  let mock: UpkeepMock
  let transcoder: UpkeepTranscoder

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
    transcoder = await upkeepTranscoderFactory.connect(owner).deploy()
    registry = await keeperRegistryFactory
      .connect(owner)
      .deploy(linkToken.address, linkEthFeed.address, gasPriceFeed.address, {
        paymentPremiumPPB,
        flatFeeMicroLink,
        blockCountPerTurn,
        checkGasLimit,
        stalenessSeconds,
        gasCeilingMultiplier,
        minUpkeepSpend,
        maxPerformGas,
        fallbackGasPrice,
        fallbackLinkPrice,
        transcoder: transcoder.address,
        registrar: ethers.constants.AddressZero,
      })
    registry2 = await keeperRegistryFactory
      .connect(owner)
      .deploy(linkToken.address, linkEthFeed.address, gasPriceFeed.address, {
        paymentPremiumPPB,
        flatFeeMicroLink,
        blockCountPerTurn,
        checkGasLimit,
        stalenessSeconds,
        gasCeilingMultiplier,
        minUpkeepSpend,
        maxPerformGas,
        fallbackGasPrice,
        fallbackLinkPrice,
        transcoder: transcoder.address,
        registrar: ethers.constants.AddressZero,
      })
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
        randomBytes,
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

  describe('#typeAndVersion', () => {
    it('uses the correct type and version', async () => {
      const typeAndVersion = await registry.typeAndVersion()
      assert.equal(typeAndVersion, 'KeeperRegistry 1.2.0')
    })
  })

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
        'DuplicateEntry()',
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
        'ParameterLengthError()',
      )
      await evmRevert(
        registry
          .connect(owner)
          .setKeepers(
            [await keeper1.getAddress()],
            [await payee1.getAddress(), await payee2.getAddress()],
          ),
        'ParameterLengthError()',
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
        'InvalidPayee()',
      )
    })

    it('emits events for every keeper added and removed', async () => {
      const oldKeepers = [
        await keeper1.getAddress(),
        await keeper2.getAddress(),
      ]
      const oldPayees = [await payee1.getAddress(), await payee2.getAddress()]
      await registry.connect(owner).setKeepers(oldKeepers, oldPayees)
      assert.deepEqual(oldKeepers, (await registry.getState()).keepers)

      // remove keepers
      const newKeepers = [
        await keeper2.getAddress(),
        await keeper3.getAddress(),
      ]
      const newPayees = [await payee2.getAddress(), await payee3.getAddress()]
      const tx = await registry.connect(owner).setKeepers(newKeepers, newPayees)
      assert.deepEqual(newKeepers, (await registry.getState()).keepers)

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
      assert.deepEqual(oldKeepers, (await registry.getState()).keepers)

      const newKeepers = [
        await keeper2.getAddress(),
        await keeper3.getAddress(),
      ]
      const newPayees = [IGNORE_ADDRESS, await payee3.getAddress()]
      const tx = await registry.connect(owner).setKeepers(newKeepers, newPayees)
      assert.deepEqual(newKeepers, (await registry.getState()).keepers)

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
        'InvalidPayee()',
      )
    })
  })

  describe('#registerUpkeep', () => {
    context('and the registry is paused', () => {
      beforeEach(async () => {
        await registry.connect(owner).pause()
      })
      it('reverts', async () => {
        await evmRevert(
          registry
            .connect(owner)
            .registerUpkeep(
              zeroAddress,
              executeGas,
              await admin.getAddress(),
              emptyBytes,
            ),
          'Pausable: paused',
        )
      })
    })

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
        'NotAContract()',
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
        'OnlyCallableByOwnerOrRegistrar()',
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
        'GasLimitOutsideRange()',
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
        'GasLimitOutsideRange()',
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
        'UpkeepNotActive()',
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
        'UpkeepNotActive()',
      )
    })
  })

  describe('#setUpkeepGasLimit', () => {
    const newGasLimit = BigNumber.from('500000')

    it('reverts if the registration does not exist', async () => {
      await evmRevert(
        registry.connect(keeper1).setUpkeepGasLimit(id.add(1), newGasLimit),
        'UpkeepNotActive()',
      )
    })

    it('reverts if the upkeep is canceled', async () => {
      await registry.connect(admin).cancelUpkeep(id)
      await evmRevert(
        registry.connect(keeper1).setUpkeepGasLimit(id, newGasLimit),
        'UpkeepNotActive()',
      )
    })

    it('reverts if called by anyone but the admin', async () => {
      await evmRevert(
        registry.connect(owner).setUpkeepGasLimit(id, newGasLimit),
        'OnlyCallableByAdmin()',
      )
    })

    it('reverts if new gas limit is out of bounds', async () => {
      await evmRevert(
        registry.connect(admin).setUpkeepGasLimit(id, BigNumber.from('100')),
        'GasLimitOutsideRange()',
      )
      await evmRevert(
        registry
          .connect(admin)
          .setUpkeepGasLimit(id, BigNumber.from('6000000')),
        'GasLimitOutsideRange()',
      )
    })

    it('updates the gas limit successfully', async () => {
      const initialGasLimit = (await registry.getUpkeep(id)).executeGas
      assert.equal(initialGasLimit, executeGas.toNumber())
      await registry.connect(admin).setUpkeepGasLimit(id, newGasLimit)
      const updatedGasLimit = (await registry.getUpkeep(id)).executeGas
      assert.equal(updatedGasLimit, newGasLimit.toNumber())
    })

    it('emits a log', async () => {
      const tx = await registry
        .connect(admin)
        .setUpkeepGasLimit(id, newGasLimit)
      await expect(tx)
        .to.emit(registry, 'UpkeepGasLimitSet')
        .withArgs(id, newGasLimit)
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
        'InsufficientFunds()',
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
          'OnlySimulatedBackend()',
        )
      })

      it('reverts if the specified keeper is not valid', async () => {
        await mock.setCanPerform(true)
        await mock.setCanCheck(true)
        await evmRevert(
          registry.checkUpkeep(id, await owner.getAddress()),
          'OnlySimulatedBackend()',
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
            'UpkeepNotNeeded()',
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
            'TargetCheckReverted',
          )
        })
      })

      context('and upkeep check simulations succeeds', () => {
        beforeEach(async () => {
          await mock.setCanCheck(true)
          await mock.setCanPerform(true)
        })

        it('returns true with pricing info if the target can execute', async () => {
          const newGasMultiplier = BigNumber.from(10)
          await registry.connect(owner).setConfig({
            paymentPremiumPPB,
            flatFeeMicroLink,
            blockCountPerTurn,
            checkGasLimit,
            stalenessSeconds,
            gasCeilingMultiplier: newGasMultiplier,
            minUpkeepSpend,
            maxPerformGas,
            fallbackGasPrice,
            fallbackLinkPrice,
            transcoder: transcoder.address,
            registrar: ethers.constants.AddressZero,
          })
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
          await mock.setCheckGasToBurn(checkGasLimit)
          await mock.setPerformGasToBurn(executeGas)
          const gas = checkGasLimit
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
        'InsufficientFunds()',
      )
    })

    context('and the registry is paused', () => {
      beforeEach(async () => {
        await registry.connect(owner).pause()
      })

      it('reverts', async () => {
        await evmRevert(
          registry.connect(keeper2).performUpkeep(id, '0x'),
          'Pausable: paused',
        )
      })
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
        expect(eventLog?.[1].args?.[0]).to.equal(id)
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

      it('updates amount spent correctly', async () => {
        const registrationBefore = await registry.getUpkeep(id)
        const balanceBefore = registrationBefore.balance
        const amountSpentBefore = registrationBefore.amountSpent

        // Do the thing
        await registry.connect(keeper1).performUpkeep(id, '0x')

        const registrationAfter = await registry.getUpkeep(id)
        const balanceAfter = registrationAfter.balance
        const amountSpentAfter = registrationAfter.amountSpent

        assert.isTrue(balanceAfter.lt(balanceBefore))
        assert.isTrue(amountSpentAfter.gt(amountSpentBefore))
        assert.isTrue(
          amountSpentAfter
            .sub(amountSpentBefore)
            .eq(balanceBefore.sub(balanceAfter)),
        )
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
        await registry.connect(owner).setConfig({
          paymentPremiumPPB,
          flatFeeMicroLink,
          blockCountPerTurn,
          checkGasLimit,
          stalenessSeconds,
          gasCeilingMultiplier: multiplier,
          minUpkeepSpend,
          maxPerformGas,
          fallbackGasPrice,
          fallbackLinkPrice,
          transcoder: transcoder.address,
          registrar: ethers.constants.AddressZero,
        })

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
        await registry.connect(owner).setConfig({
          paymentPremiumPPB,
          flatFeeMicroLink,
          blockCountPerTurn,
          checkGasLimit,
          stalenessSeconds,
          gasCeilingMultiplier: multiplier,
          minUpkeepSpend,
          maxPerformGas,
          fallbackGasPrice,
          fallbackLinkPrice,
          transcoder: transcoder.address,
          registrar: ethers.constants.AddressZero,
        })

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
          'OnlyActiveKeepers()',
        )
      })

      it('reverts if the upkeep has been canceled', async () => {
        await mock.setCanPerform(true)

        await registry.connect(owner).cancelUpkeep(id)

        await evmRevert(
          registry.connect(keeper1).performUpkeep(id, '0x'),
          'UpkeepNotActive()',
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

      it('uses the fallback link price if the feed price is non-sensical [ @skip-coverage ]', async () => {
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
          'KeepersMustTakeTurns()',
        )
        await registry.connect(keeper2).performUpkeep(id, '0x')
        await evmRevert(
          registry.connect(keeper2).performUpkeep(id, '0x'),
          'KeepersMustTakeTurns()',
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
        expect(eventLog?.[1].args?.[0]).to.equal(id)
        assert.equal(eventLog?.[1].args?.[1], true)
        assert.equal(eventLog?.[1].args?.[2], await keeper1.getAddress())
        assert.isNotEmpty(eventLog?.[1].args?.[3])
        assert.equal(eventLog?.[1].args?.[4], performData)
      })

      it('can self fund', async () => {
        const autoFunderUpkeep = await upkeepAutoFunderFactory
          .connect(owner)
          .deploy(linkToken.address, registry.address)
        const tx = await registry
          .connect(owner)
          .registerUpkeep(
            autoFunderUpkeep.address,
            executeGas,
            autoFunderUpkeep.address,
            emptyBytes,
          )
        const upkeepID = await getUpkeepID(tx)
        await autoFunderUpkeep.setUpkeepId(upkeepID)
        // Give enough funds for upkeep as well as to the upkeep contract
        await linkToken.connect(owner).approve(registry.address, toWei('1000'))
        await linkToken
          .connect(owner)
          .transfer(autoFunderUpkeep.address, toWei('1000'))
        let maxPayment = await registry.getMaxPaymentForGas(executeGas)

        // First set auto funding amount to 0 and verify that balance is deducted upon performUpkeep
        let initialBalance = toWei('100')
        await registry.connect(owner).addFunds(upkeepID, initialBalance)
        await autoFunderUpkeep.setAutoFundLink(0)
        await autoFunderUpkeep.setIsEligible(true)
        await registry.connect(keeper1).performUpkeep(upkeepID, '0x')

        let postUpkeepBalance = (await registry.getUpkeep(upkeepID)).balance
        assert.isTrue(postUpkeepBalance.lt(initialBalance)) // Balance should be deducted
        assert.isTrue(postUpkeepBalance.gte(initialBalance.sub(maxPayment))) // Balance should not be deducted more than maxPayment

        // Now set auto funding amount to 100 wei and verify that the balance increases
        initialBalance = postUpkeepBalance
        let autoTopupAmount = toWei('100')
        await autoFunderUpkeep.setAutoFundLink(autoTopupAmount)
        await autoFunderUpkeep.setIsEligible(true)
        await registry.connect(keeper2).performUpkeep(upkeepID, '0x')

        postUpkeepBalance = (await registry.getUpkeep(upkeepID)).balance
        // Balance should increase by autoTopupAmount and decrease by max maxPayment
        assert.isTrue(
          postUpkeepBalance.gte(
            initialBalance.add(autoTopupAmount).sub(maxPayment),
          ),
        )
      })

      it('can self cancel', async () => {
        const autoFunderUpkeep = await upkeepAutoFunderFactory
          .connect(owner)
          .deploy(linkToken.address, registry.address)
        const tx = await registry
          .connect(owner)
          .registerUpkeep(
            autoFunderUpkeep.address,
            executeGas,
            autoFunderUpkeep.address,
            emptyBytes,
          )
        const upkeepID = await getUpkeepID(tx)
        await autoFunderUpkeep.setUpkeepId(upkeepID)

        await linkToken.connect(owner).approve(registry.address, toWei('1000'))
        await registry.connect(owner).addFunds(upkeepID, toWei('100'))
        await autoFunderUpkeep.setIsEligible(true)
        await autoFunderUpkeep.setShouldCancel(true)

        let registration = await registry.getUpkeep(upkeepID)
        const oldExpiration = registration.maxValidBlocknumber

        // Do the thing
        await registry.connect(keeper1).performUpkeep(upkeepID, '0x')

        // Verify upkeep gets cancelled
        registration = await registry.getUpkeep(upkeepID)
        const newExpiration = registration.maxValidBlocknumber
        assert.isTrue(newExpiration.lt(oldExpiration))
      })
    })
  })

  describe('#withdrawFunds', () => {
    beforeEach(async () => {
      await linkToken.connect(keeper1).approve(registry.address, toWei('100'))
      await registry.connect(keeper1).addFunds(id, toWei('100'))
      await registry.connect(keeper1).performUpkeep(id, '0x')
    })

    it('reverts if called by anyone but the admin', async () => {
      await evmRevert(
        registry
          .connect(owner)
          .withdrawFunds(id.add(1), await payee1.getAddress()),
        'OnlyCallableByAdmin()',
      )
    })

    it('reverts if called on an uncanceled upkeep', async () => {
      await evmRevert(
        registry.connect(admin).withdrawFunds(id, await payee1.getAddress()),
        'UpkeepNotCanceled()',
      )
    })

    it('reverts if called with the 0 address', async () => {
      await evmRevert(
        registry.connect(admin).withdrawFunds(id, zeroAddress),
        'InvalidRecipient()',
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
        let previousBalance = registration.balance

        await registry
          .connect(admin)
          .withdrawFunds(id, await payee1.getAddress())

        const payee1After = await linkToken.balanceOf(await payee1.getAddress())
        const registryAfter = await linkToken.balanceOf(registry.address)

        assert.isTrue(payee1Before.add(previousBalance).eq(payee1After))
        assert.isTrue(registryBefore.sub(previousBalance).eq(registryAfter))

        registration = await registry.getUpkeep(id)
        assert.equal(0, registration.balance.toNumber())
      })

      it('deducts cancellation fees from upkeep and gives to owner', async () => {
        let minUpkeepSpend = toWei('10')
        await registry.connect(owner).setConfig({
          paymentPremiumPPB,
          flatFeeMicroLink,
          blockCountPerTurn,
          checkGasLimit,
          stalenessSeconds,
          gasCeilingMultiplier,
          minUpkeepSpend,
          maxPerformGas,
          fallbackGasPrice,
          fallbackLinkPrice,
          transcoder: transcoder.address,
          registrar: ethers.constants.AddressZero,
        })

        const payee1Before = await linkToken.balanceOf(
          await payee1.getAddress(),
        )
        let upkeepBefore = (await registry.getUpkeep(id)).balance
        let ownerBefore = (await registry.getState()).state.ownerLinkBalance
        assert.equal(0, ownerBefore.toNumber())

        let amountSpent = toWei('100').sub(upkeepBefore)
        let cancellationFee = minUpkeepSpend.sub(amountSpent)

        await registry
          .connect(admin)
          .withdrawFunds(id, await payee1.getAddress())

        const payee1After = await linkToken.balanceOf(await payee1.getAddress())
        let upkeepAfter = (await registry.getUpkeep(id)).balance
        let ownerAfter = (await registry.getState()).state.ownerLinkBalance

        // Post upkeep balance should be 0
        assert.equal(0, upkeepAfter.toNumber())
        // balance - cancellationFee should be transferred to payee
        assert.isTrue(
          payee1Before.add(upkeepBefore.sub(cancellationFee)).eq(payee1After),
        )
        assert.isTrue(ownerAfter.eq(cancellationFee))
      })

      it('deducts max upto balance as cancellation fees', async () => {
        // Very high min spend, should deduct whole balance as cancellation fees
        let minUpkeepSpend = toWei('1000')
        await registry.connect(owner).setConfig({
          paymentPremiumPPB,
          flatFeeMicroLink,
          blockCountPerTurn,
          checkGasLimit,
          stalenessSeconds,
          gasCeilingMultiplier,
          minUpkeepSpend,
          maxPerformGas,
          fallbackGasPrice,
          fallbackLinkPrice,
          transcoder: transcoder.address,
          registrar: ethers.constants.AddressZero,
        })
        const payee1Before = await linkToken.balanceOf(
          await payee1.getAddress(),
        )
        let upkeepBefore = (await registry.getUpkeep(id)).balance
        let ownerBefore = (await registry.getState()).state.ownerLinkBalance
        assert.equal(0, ownerBefore.toNumber())

        await registry
          .connect(admin)
          .withdrawFunds(id, await payee1.getAddress())
        const payee1After = await linkToken.balanceOf(await payee1.getAddress())
        let ownerAfter = (await registry.getState()).state.ownerLinkBalance
        let upkeepAfter = (await registry.getUpkeep(id)).balance

        assert.equal(0, upkeepAfter.toNumber())
        // No funds should be transferred, all of upkeepBefore should be given to owner
        assert.isTrue(payee1After.eq(payee1Before))
        assert.isTrue(ownerAfter.eq(upkeepBefore))
      })

      it('does not deduct cancellation fees if enough is spent', async () => {
        // Very low min spend, already spent in one upkeep
        let minUpkeepSpend = BigNumber.from(420)
        await registry.connect(owner).setConfig({
          paymentPremiumPPB,
          flatFeeMicroLink,
          blockCountPerTurn,
          checkGasLimit,
          stalenessSeconds,
          gasCeilingMultiplier,
          minUpkeepSpend,
          maxPerformGas,
          fallbackGasPrice,
          fallbackLinkPrice,
          transcoder: transcoder.address,
          registrar: ethers.constants.AddressZero,
        })
        const payee1Before = await linkToken.balanceOf(
          await payee1.getAddress(),
        )
        let upkeepBefore = (await registry.getUpkeep(id)).balance
        let ownerBefore = (await registry.getState()).state.ownerLinkBalance
        assert.equal(0, ownerBefore.toNumber())

        await registry
          .connect(admin)
          .withdrawFunds(id, await payee1.getAddress())
        const payee1After = await linkToken.balanceOf(await payee1.getAddress())
        let ownerAfter = (await registry.getState()).state.ownerLinkBalance
        let upkeepAfter = (await registry.getUpkeep(id)).balance

        assert.equal(0, upkeepAfter.toNumber())
        // No cancellation fees for owner
        assert.equal(0, ownerAfter.toNumber())
        // Whole balance transferred to payee
        assert.isTrue(payee1Before.add(upkeepBefore).eq(payee1After))
      })
    })
  })

  describe('#withdrawOwnerFunds', () => {
    it('can only be called by owner', async () => {
      await evmRevert(
        registry.connect(keeper1).withdrawOwnerFunds(),
        'Only callable by owner',
      )
    })

    it('withdraws the collected fees to owner', async () => {
      await linkToken.connect(keeper1).approve(registry.address, toWei('100'))
      await registry.connect(keeper1).addFunds(id, toWei('100'))
      // Very high min spend, whole balance as cancellation fees
      let minUpkeepSpend = toWei('1000')
      await registry.connect(owner).setConfig({
        paymentPremiumPPB,
        flatFeeMicroLink,
        blockCountPerTurn,
        checkGasLimit,
        stalenessSeconds,
        gasCeilingMultiplier,
        minUpkeepSpend,
        maxPerformGas,
        fallbackGasPrice,
        fallbackLinkPrice,
        transcoder: transcoder.address,
        registrar: ethers.constants.AddressZero,
      })
      let upkeepBalance = (await registry.getUpkeep(id)).balance
      const ownerBefore = await linkToken.balanceOf(await owner.getAddress())

      await registry.connect(owner).cancelUpkeep(id)
      await registry.connect(admin).withdrawFunds(id, await payee1.getAddress())
      // Transfered to owner balance on registry
      let ownerRegistryBalance = (await registry.getState()).state
        .ownerLinkBalance
      assert.isTrue(ownerRegistryBalance.eq(upkeepBalance))

      // Now withdraw
      await registry.connect(owner).withdrawOwnerFunds()

      ownerRegistryBalance = (await registry.getState()).state.ownerLinkBalance
      const ownerAfter = await linkToken.balanceOf(await owner.getAddress())

      // Owner registry balance should be changed to 0
      assert.isTrue(ownerRegistryBalance.eq(BigNumber.from('0')))

      // Owner should be credited with the balance
      assert.isTrue(ownerBefore.add(upkeepBalance).eq(ownerAfter))
    })
  })

  describe('#cancelUpkeep', () => {
    it('reverts if the ID is not valid', async () => {
      await evmRevert(
        registry.connect(owner).cancelUpkeep(id.add(1)),
        'CannotCancel()',
      )
    })

    it('reverts if called by a non-owner/non-admin', async () => {
      await evmRevert(
        registry.connect(keeper1).cancelUpkeep(id),
        'OnlyCallableByOwnerOrAdmin()',
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

      it('immediately prevents upkeep', async () => {
        await registry.connect(owner).cancelUpkeep(id)

        await evmRevert(
          registry.connect(keeper2).performUpkeep(id, '0x'),
          'UpkeepNotActive()',
        )
      })

      it('does not revert if reverts if called multiple times', async () => {
        await registry.connect(owner).cancelUpkeep(id)
        await evmRevert(
          registry.connect(owner).cancelUpkeep(id),
          'CannotCancel()',
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

      // it('updates the canceled registrations list', async () => {
      //   let canceled = await registry.callStatic.getCanceledUpkeepList()
      //   assert.deepEqual([], canceled)

      //   await registry.connect(admin).cancelUpkeep(id)

      //   canceled = await registry.callStatic.getCanceledUpkeepList()
      //   assert.deepEqual([id], canceled)
      // })

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
          'UpkeepNotActive()',
        )
      })

      it('reverts if called again by the admin', async () => {
        await registry.connect(admin).cancelUpkeep(id)

        await evmRevert(
          registry.connect(admin).cancelUpkeep(id),
          'CannotCancel()',
        )
      })

      // it('does not revert or double add the cancellation record if called by the owner immediately after', async () => {
      //   await registry.connect(admin).cancelUpkeep(id)

      //   await registry.connect(owner).cancelUpkeep(id)

      //   const canceled = await registry.callStatic.getCanceledUpkeepList()
      //   assert.deepEqual([id], canceled)
      // })

      it('reverts if called by the owner after the timeout', async () => {
        await registry.connect(admin).cancelUpkeep(id)

        for (let i = 0; i < delay; i++) {
          await ethers.provider.send('evm_mine', [])
        }

        await evmRevert(
          registry.connect(owner).cancelUpkeep(id),
          'CannotCancel()',
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
        'OnlyCallableByPayee()',
      )
    })

    it('reverts if called with the 0 address', async () => {
      await evmRevert(
        registry
          .connect(payee2)
          .withdrawPayment(await keeper1.getAddress(), zeroAddress),
        'InvalidRecipient()',
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
        'OnlyCallableByPayee()',
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
        'ValueNotChanged()',
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
        'OnlyCallableByProposedPayee()',
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
        registry.connect(payee1).setConfig({
          paymentPremiumPPB: payment,
          flatFeeMicroLink: flatFee,
          blockCountPerTurn: checks,
          checkGasLimit: maxGas,
          stalenessSeconds: staleness,
          gasCeilingMultiplier,
          minUpkeepSpend,
          maxPerformGas,
          fallbackGasPrice: fbGasEth,
          fallbackLinkPrice: fbLinkEth,
          transcoder: transcoder.address,
          registrar: ethers.constants.AddressZero,
        }),
        'Only callable by owner',
      )
    })

    it('updates the config', async () => {
      const old = (await registry.getState()).config
      assert.isTrue(paymentPremiumPPB.eq(old.paymentPremiumPPB))
      assert.isTrue(flatFeeMicroLink.eq(old.flatFeeMicroLink))
      assert.isTrue(blockCountPerTurn.eq(old.blockCountPerTurn))
      assert.isTrue(stalenessSeconds.eq(old.stalenessSeconds))
      assert.isTrue(gasCeilingMultiplier.eq(old.gasCeilingMultiplier))

      await registry.connect(owner).setConfig({
        paymentPremiumPPB: payment,
        flatFeeMicroLink: flatFee,
        blockCountPerTurn: checks,
        checkGasLimit: maxGas,
        stalenessSeconds: staleness,
        gasCeilingMultiplier: ceiling,
        minUpkeepSpend,
        maxPerformGas,
        fallbackGasPrice: fbGasEth,
        fallbackLinkPrice: fbLinkEth,
        transcoder: transcoder.address,
        registrar: ethers.constants.AddressZero,
      })

      const updated = (await registry.getState()).config
      assert.equal(updated.paymentPremiumPPB, payment.toNumber())
      assert.equal(updated.flatFeeMicroLink, flatFee.toNumber())
      assert.equal(updated.blockCountPerTurn, checks.toNumber())
      assert.equal(updated.stalenessSeconds, staleness.toNumber())
      assert.equal(updated.gasCeilingMultiplier, ceiling.toNumber())
      assert.equal(updated.checkGasLimit, maxGas.toNumber())
      assert.equal(updated.fallbackGasPrice.toNumber(), fbGasEth.toNumber())
      assert.equal(updated.fallbackLinkPrice.toNumber(), fbLinkEth.toNumber())
    })

    it('emits an event', async () => {
      const tx = await registry.connect(owner).setConfig({
        paymentPremiumPPB: payment,
        flatFeeMicroLink: flatFee,
        blockCountPerTurn: checks,
        checkGasLimit: maxGas,
        stalenessSeconds: staleness,
        gasCeilingMultiplier: ceiling,
        minUpkeepSpend,
        maxPerformGas,
        fallbackGasPrice: fbGasEth,
        fallbackLinkPrice: fbLinkEth,
        transcoder: transcoder.address,
        registrar: ethers.constants.AddressZero,
      })
      await expect(tx)
        .to.emit(registry, 'ConfigSet')
        .withArgs([
          payment,
          flatFee,
          checks,
          maxGas,
          staleness,
          ceiling,
          minUpkeepSpend,
          maxPerformGas,
          fbGasEth,
          fbLinkEth,
        ])
    })
  })

  describe('#onTokenTransfer', () => {
    const amount = toWei('1')

    it('reverts if not called by the LINK token', async () => {
      const data = ethers.utils.defaultAbiCoder.encode(['uint256'], [id])

      await evmRevert(
        registry
          .connect(keeper1)
          .onTokenTransfer(await keeper1.getAddress(), amount, data),
        'OnlyCallableByLINKToken()',
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
        'UpkeepNotActive()',
      )
    })

    it('updates the funds of the job id passed', async () => {
      const data = ethers.utils.defaultAbiCoder.encode(['uint256'], [id])

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
      const data = ethers.utils.defaultAbiCoder.encode(['uint256'], [id2])
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
            await registry.connect(owner).setConfig({
              paymentPremiumPPB: premium,
              flatFeeMicroLink: flatFee,
              blockCountPerTurn,
              checkGasLimit,
              stalenessSeconds,
              gasCeilingMultiplier,
              minUpkeepSpend,
              maxPerformGas,
              fallbackGasPrice,
              fallbackLinkPrice,
              transcoder: transcoder.address,
              registrar: ethers.constants.AddressZero,
            })
            const price = await registry.getMaxPaymentForGas(gas)
            expect(price).to.equal(linkForGas(gas, premium, flatFee))
          }
        }
      }
    })
  })

  describe('#setPeerRegistryMigrationPermission() / #getPeerRegistryMigrationPermission()', () => {
    const peer = randomAddress()
    it('allows the owner to set the peer registries', async () => {
      let permission = await registry.getPeerRegistryMigrationPermission(peer)
      expect(permission).to.equal(0)
      await registry.setPeerRegistryMigrationPermission(peer, 1)
      permission = await registry.getPeerRegistryMigrationPermission(peer)
      expect(permission).to.equal(1)
      await registry.setPeerRegistryMigrationPermission(peer, 2)
      permission = await registry.getPeerRegistryMigrationPermission(peer)
      expect(permission).to.equal(2)
      await registry.setPeerRegistryMigrationPermission(peer, 0)
      permission = await registry.getPeerRegistryMigrationPermission(peer)
      expect(permission).to.equal(0)
    })
    it('reverts if passed an unsupported permission', async () => {
      await expect(
        registry.connect(admin).setPeerRegistryMigrationPermission(peer, 10),
      ).to.be.reverted
    })
    it('reverts if not called by the owner', async () => {
      await expect(
        registry.connect(admin).setPeerRegistryMigrationPermission(peer, 1),
      ).to.be.revertedWith('Only callable by owner')
    })
  })

  describe('migrateUpkeeps() / #receiveUpkeeps()', async () => {
    context('when permissions are set', () => {
      beforeEach(async () => {
        await linkToken.connect(owner).approve(registry.address, toWei('100'))
        await registry.connect(owner).addFunds(id, toWei('100'))
        await registry.setPeerRegistryMigrationPermission(registry2.address, 1)
        await registry2.setPeerRegistryMigrationPermission(registry.address, 2)
      })

      it('migrates an upkeep', async () => {
        expect((await registry.getUpkeep(id)).balance).to.equal(toWei('100'))
        expect((await registry.getUpkeep(id)).checkData).to.equal(randomBytes)
        expect((await registry.getState()).state.numUpkeeps).to.equal(1)
        // migrate
        await registry.connect(admin).migrateUpkeeps([id], registry2.address)
        expect((await registry.getState()).state.numUpkeeps).to.equal(0)
        expect((await registry2.getState()).state.numUpkeeps).to.equal(1)
        expect((await registry.getUpkeep(id)).balance).to.equal(0)
        expect((await registry.getUpkeep(id)).checkData).to.equal('0x')
        expect((await registry2.getUpkeep(id)).balance).to.equal(toWei('100'))
        expect((await registry2.getState()).state.expectedLinkBalance).to.equal(
          toWei('100'),
        )
        expect((await registry2.getUpkeep(id)).checkData).to.equal(randomBytes)
      })
      it('emits an event on both contracts', async () => {
        expect((await registry.getUpkeep(id)).balance).to.equal(toWei('100'))
        expect((await registry.getUpkeep(id)).checkData).to.equal(randomBytes)
        expect((await registry.getState()).state.numUpkeeps).to.equal(1)
        const tx = registry
          .connect(admin)
          .migrateUpkeeps([id], registry2.address)
        await expect(tx)
          .to.emit(registry, 'UpkeepMigrated')
          .withArgs(id, toWei('100'), registry2.address)
        await expect(tx)
          .to.emit(registry2, 'UpkeepReceived')
          .withArgs(id, toWei('100'), registry.address)
      })
    })

    context('when permissions are not set', () => {
      it('reverts', async () => {
        // no permissions
        await registry.setPeerRegistryMigrationPermission(registry2.address, 0)
        await registry2.setPeerRegistryMigrationPermission(registry.address, 0)
        await expect(registry.migrateUpkeeps([id], registry2.address)).to.be
          .reverted
        // only outgoing permissions
        await registry.setPeerRegistryMigrationPermission(registry2.address, 1)
        await registry2.setPeerRegistryMigrationPermission(registry.address, 0)
        await expect(registry.migrateUpkeeps([id], registry2.address)).to.be
          .reverted
        // only incoming permissions
        await registry.setPeerRegistryMigrationPermission(registry2.address, 0)
        await registry2.setPeerRegistryMigrationPermission(registry.address, 2)
        await expect(registry.migrateUpkeeps([id], registry2.address)).to.be
          .reverted
        // permissions opposite direction
        await registry.setPeerRegistryMigrationPermission(registry2.address, 2)
        await registry2.setPeerRegistryMigrationPermission(registry.address, 1)
        await expect(registry.migrateUpkeeps([id], registry2.address)).to.be
          .reverted
      })
    })
  })

  describe('#checkUpkeep / #performUpkeep', () => {
    const performData = '0xc0ffeec0ffee'
    const multiplier = BigNumber.from(10)
    const flatFee = BigNumber.from('100000') //0.1 LINK
    const callGasPrice = 1

    it('uses the same minimum balance calculation [ @skip-coverage ]', async () => {
      await registry.connect(owner).setConfig({
        paymentPremiumPPB,
        flatFeeMicroLink: flatFee,
        blockCountPerTurn,
        checkGasLimit,
        stalenessSeconds,
        gasCeilingMultiplier: multiplier,
        minUpkeepSpend,
        maxPerformGas,
        fallbackGasPrice,
        fallbackLinkPrice,
        transcoder: transcoder.address,
        registrar: ethers.constants.AddressZero,
      })
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
        'InsufficientFunds()',
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
        'InsufficientFunds()',
      )
      await registry.connect(keeper1).addFunds(id, oneWei)
      await registry
        .connect(zeroAddress)
        .callStatic.checkUpkeep(id, await keeper1.getAddress())
    })
  })
})
