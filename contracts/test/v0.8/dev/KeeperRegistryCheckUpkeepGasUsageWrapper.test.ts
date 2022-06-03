import { ethers } from 'hardhat'
import { BigNumber, Signer } from 'ethers'
import { assert } from 'chai'
import { KeeperRegistryCheckUpkeepGasUsageWrapper } from '../../../typechain/KeeperRegistryCheckUpkeepGasUsageWrapper'
import { getUsers, Personas } from '../../test-helpers/setup'
import {
  deployMockContract,
  MockContract,
} from '@ethereum-waffle/mock-contract'
import { abi as registryAbi } from '../../../artifacts/src/v0.8/KeeperRegistry.sol/KeeperRegistry.json'

let personas: Personas
let owner: Signer
let caller: Signer
let nelly: Signer
let ned: Signer
let nancy: Signer
let registryMockContract: MockContract
let gasUsageWrapper: KeeperRegistryCheckUpkeepGasUsageWrapper
let keeperAddresses: string[]

const nonce = BigNumber.from(200)
const ownerLinkBalance = BigNumber.from(200)
const expectedLinkBalance = BigNumber.from(200)
const numUpkeeps = BigNumber.from(200)
const state = {
  nonce,
  ownerLinkBalance,
  expectedLinkBalance,
  numUpkeeps,
}

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

const upkeepId = 123
const lastKeeperIndex = 0

describe('KeeperRegistryCheckUpkeepGasUsageWrapper', () => {
  before(async () => {
    personas = (await getUsers()).personas
    owner = personas.Default
    caller = personas.Carol
    nelly = personas.Nelly
    ned = personas.Ned
    nancy = personas.Nancy

    keeperAddresses = [
      await nelly.getAddress(),
      await ned.getAddress(),
      await nancy.getAddress(),
    ]

    registryMockContract = await deployMockContract(owner as any, registryAbi)
    const gasUsageWrapperFactory = await ethers.getContractFactory(
      'KeeperRegistryCheckUpkeepGasUsageWrapper',
    )
    gasUsageWrapper = await gasUsageWrapperFactory
      .connect(owner)
      .deploy(registryMockContract.address)
    await gasUsageWrapper.deployed()

    await registryMockContract.mock.getUpkeep
      .withArgs(upkeepId)
      .returns(
        ethers.constants.AddressZero,
        0,
        '0x',
        0,
        keeperAddresses[lastKeeperIndex],
        ethers.constants.AddressZero,
        0,
        0,
      )

    await registryMockContract.mock.getState.returns(
      state,
      config,
      keeperAddresses,
    )
  })

  describe('measureCheckGas()', () => {
    it("returns gas used when registry's checkUpkeep executes successfully", async () => {
      const n = 15
      await mineNBlocks(n)

      let nextKeeperIndex = (n + 5) % keeperAddresses.length
      if (nextKeeperIndex == lastKeeperIndex) {
        nextKeeperIndex = (lastKeeperIndex + 1) % keeperAddresses.length
      }

      await registryMockContract.mock.checkUpkeep
        .withArgs(upkeepId, keeperAddresses[nextKeeperIndex])
        .returns(
          '0x' /* performData */,
          BigNumber.from(1000) /* maxLinkPayment */,
          BigNumber.from(2000) /* gasLimit */,
          BigNumber.from(3000) /* adjustedGasWei */,
          BigNumber.from(4000) /* linkEth */,
        )

      const response = await gasUsageWrapper
        .connect(caller)
        .callStatic.measureCheckGas(BigNumber.from(upkeepId))

      assert.isTrue(response[0], 'The checkUpkeepSuccess should be true')
      assert.equal(
        response[1],
        '0x',
        'The performData should be forwarded correctly',
      )
      assert.isTrue(
        response[2] > BigNumber.from(0),
        'The gasUsed value must be larger than 0',
      )
    })

    it("returns gas used when registry's checkUpkeep reverts", async () => {
      const n = 15
      await mineNBlocks(n)

      let nextKeeperIndex = (n + 5) % keeperAddresses.length
      if (nextKeeperIndex == lastKeeperIndex) {
        nextKeeperIndex = (lastKeeperIndex + 1) % keeperAddresses.length
      }

      await registryMockContract.mock.checkUpkeep
        .withArgs(upkeepId, keeperAddresses[nextKeeperIndex])
        .revertsWithReason('Error')

      const response = await gasUsageWrapper
        .connect(caller)
        .callStatic.measureCheckGas(BigNumber.from(upkeepId))

      assert.isFalse(response[0], 'The checkUpkeepSuccess should be false')
      assert.equal(
        response[1],
        '0x',
        'The performData should be forwarded correctly',
      )
      assert.isTrue(
        response[2] > BigNumber.from(0),
        'The gasUsed value must be larger than 0',
      )
    })
  })

  describe('getKeeperRegistry()', () => {
    it('returns the underlying keeper registry', async () => {
      const registry = await gasUsageWrapper.connect(caller).getKeeperRegistry()
      assert.equal(
        registry,
        registryMockContract.address,
        'The underlying keeper registry is incorrect',
      )
    })
  })
})

async function mineNBlocks(n: number) {
  for (let index = 0; index < n; index++) {
    await ethers.provider.send('evm_mine', [])
  }
}
