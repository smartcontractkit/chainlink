import { ethers } from 'hardhat'
import { BigNumber, Signer } from 'ethers'
import moment from 'moment'
import { assert } from 'chai'
import { CanaryUpkeep12 as CanaryUpkeep } from '../../../typechain/CanaryUpkeep12'
import { CanaryUpkeep12__factory as CanaryUpkeepFactory } from '../../../typechain/factories/CanaryUpkeep12__factory'
import { KeeperRegistry12 as KeeperRegistry } from '../../../typechain/KeeperRegistry12'
import { KeeperRegistry12__factory as KeeperRegistryFactory } from '../../../typechain/factories/KeeperRegistry12__factory'
import { fastForward, reset } from '../../test-helpers/helpers'
import { getUsers, Personas } from '../../test-helpers/setup'
import { evmRevert } from '../../test-helpers/matchers'

let personas: Personas
let canaryUpkeep: CanaryUpkeep
let canaryUpkeepFactory: CanaryUpkeepFactory
let owner: Signer
let nelly: Signer
let nancy: Signer
let ned: Signer
let keeperAddresses: string[]
let keeperRegistry: KeeperRegistry
let keeperRegistryFactory: KeeperRegistryFactory

const defaultInterval = 300
const paymentPremiumPPB = BigNumber.from(250000000)
const flatFeeMicroLink = BigNumber.from(0)
const blockCountPerTurn = BigNumber.from(3)
const stalenessSeconds = BigNumber.from(43820)
const gasCeilingMultiplier = BigNumber.from(1)
const checkGasLimit = BigNumber.from(20000000)
const fallbackGasPrice = BigNumber.from(200)
const fallbackLinkPrice = BigNumber.from(200000000)
const maxPerformGas = BigNumber.from(5000000)
const minUpkeepSpend = BigNumber.from('1000000000000000000')
const transcoder = ethers.constants.AddressZero
const registrar = ethers.constants.AddressZero
const config = {
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
  transcoder,
  registrar,
}

describe('CanaryUpkeep1_2', () => {
  before(async () => {
    personas = (await getUsers()).personas
    owner = personas.Default
    nelly = personas.Nelly
    nancy = personas.Nancy
    ned = personas.Ned
    keeperAddresses = [
      await nelly.getAddress(),
      await nancy.getAddress(),
      await ned.getAddress(),
    ]
  })
  beforeEach(async () => {
    // @ts-ignore bug in autogen file
    keeperRegistryFactory = await ethers.getContractFactory('KeeperRegistry1_2')
    keeperRegistry = await keeperRegistryFactory
      .connect(owner)
      .deploy(
        ethers.constants.AddressZero,
        ethers.constants.AddressZero,
        ethers.constants.AddressZero,
        config,
      )
    await keeperRegistry.deployed()

    // @ts-ignore bug in autogen file
    canaryUpkeepFactory = await ethers.getContractFactory('CanaryUpkeep1_2')
    canaryUpkeep = await canaryUpkeepFactory
      .connect(owner)
      .deploy(keeperRegistry.address, defaultInterval)
    await canaryUpkeep.deployed()
  })

  afterEach(async () => {
    await reset()
  })

  describe('setInterval()', () => {
    it('allows the owner setting interval', async () => {
      await canaryUpkeep.connect(owner).setInterval(400)
      const newInterval = await canaryUpkeep.getInterval()
      assert.equal(
        newInterval.toNumber(),
        400,
        'The interval is not updated correctly',
      )
    })

    it('does not allow someone who is not an owner setting interval', async () => {
      await evmRevert(
        canaryUpkeep.connect(ned).setInterval(400),
        'Only callable by owner',
      )
    })
  })

  describe('checkUpkeep()', () => {
    it('returns true when sufficient time passes', async () => {
      await fastForward(moment.duration(6, 'minutes').asSeconds())
      await keeperRegistry.setKeepers(keeperAddresses, keeperAddresses)
      const [needsUpkeep] = await canaryUpkeep.checkUpkeep('0x')
      assert.isTrue(needsUpkeep)
    })

    it('returns false when insufficient time passes', async () => {
      await fastForward(moment.duration(2, 'minutes').asSeconds())
      await keeperRegistry.setKeepers(keeperAddresses, keeperAddresses)
      const [needsUpkeep] = await canaryUpkeep.checkUpkeep('0x')
      assert.isFalse(needsUpkeep)
    })

    it('returns false when keeper array is empty', async () => {
      await fastForward(moment.duration(6, 'minutes').asSeconds())
      const [needsUpkeep] = await canaryUpkeep.checkUpkeep('0x')
      assert.isTrue(needsUpkeep)
    })
  })

  describe('performUpkeep()', () => {
    it('enforces that transaction origin is the anticipated keeper', async () => {
      await keeperRegistry.setKeepers(keeperAddresses, keeperAddresses)

      const oldTimestamp = await canaryUpkeep.connect(nelly).getTimestamp()
      const oldKeeperIndex = await canaryUpkeep.connect(nelly).getKeeperIndex()
      await fastForward(moment.duration(6, 'minutes').asSeconds())
      await canaryUpkeep.connect(nelly).performUpkeep('0x')
      const newKeeperIndex = await canaryUpkeep.connect(nelly).getKeeperIndex()
      assert.equal(
        newKeeperIndex.toNumber() - oldKeeperIndex.toNumber(),
        1,
        'keeper index needs to increment by 1 after performUpkeep',
      )

      const newTimestamp = await canaryUpkeep.connect(nelly).getTimestamp()
      const interval = await canaryUpkeep.connect(nelly).getInterval()
      assert.isAtLeast(
        newTimestamp.toNumber() - oldTimestamp.toNumber(),
        interval.toNumber(),
        'timestamp needs to be updated after performUpkeep',
      )
    })

    it('enforces that keeper index will reset to zero after visiting the last keeper', async () => {
      await keeperRegistry.setKeepers(keeperAddresses, keeperAddresses)

      await fastForward(moment.duration(6, 'minutes').asSeconds())
      await canaryUpkeep.connect(nelly).performUpkeep('0x')

      await fastForward(moment.duration(6, 'minutes').asSeconds())
      await canaryUpkeep.connect(nancy).performUpkeep('0x')

      await fastForward(moment.duration(6, 'minutes').asSeconds())
      await canaryUpkeep.connect(ned).performUpkeep('0x')

      const keeperIndex = await canaryUpkeep.connect(ned).getKeeperIndex()
      assert.equal(
        keeperIndex.toNumber(),
        0,
        'Keeper index is not updated properly',
      )
    })

    it('updates the keeper index after the keepers array is shortened', async () => {
      await keeperRegistry.setKeepers(keeperAddresses, keeperAddresses)

      await fastForward(moment.duration(6, 'minutes').asSeconds())
      await canaryUpkeep.connect(nelly).performUpkeep('0x')

      await fastForward(moment.duration(6, 'minutes').asSeconds())
      await canaryUpkeep.connect(nancy).performUpkeep('0x')

      let shortAddresses: string[] = [
        await nelly.getAddress(),
        await nancy.getAddress(),
      ]
      await keeperRegistry.setKeepers(shortAddresses, shortAddresses)

      await fastForward(moment.duration(6, 'minutes').asSeconds())
      await canaryUpkeep.connect(nelly).performUpkeep('0x')
      const keeperIndex = await canaryUpkeep.getKeeperIndex()
      assert.equal(
        keeperIndex.toNumber(),
        1,
        'Keeper index is not updated properly',
      )
    })

    it('reverts if the keeper array is empty', async () => {
      await evmRevert(
        canaryUpkeep.connect(nelly).performUpkeep('0x'),
        `NoKeeperNodes`,
      )
    })

    it('reverts if not enough time has passed', async () => {
      await keeperRegistry.setKeepers(keeperAddresses, keeperAddresses)
      await fastForward(moment.duration(3, 'minutes').asSeconds())
      await evmRevert(
        canaryUpkeep.connect(nelly).performUpkeep('0x'),
        `InsufficientInterval`,
      )
    })

    it('reverts if an incorrect keeper tries to perform upkeep', async () => {
      await keeperRegistry.setKeepers(keeperAddresses, keeperAddresses)
      await fastForward(moment.duration(6, 'minutes').asSeconds())
      await evmRevert(
        canaryUpkeep.connect(nancy).performUpkeep('0x'),
        'transaction origin is not the anticipated keeper.',
      )
    })
  })
})
