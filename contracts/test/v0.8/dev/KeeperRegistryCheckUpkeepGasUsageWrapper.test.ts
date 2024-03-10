import { ethers } from 'hardhat'
import { BigNumber, Signer } from 'ethers'
import { assert } from 'chai'
import { KeeperRegistryCheckUpkeepGasUsageWrapper12 as GasWrapper } from '../../../typechain/KeeperRegistryCheckUpkeepGasUsageWrapper12'
import { KeeperRegistryCheckUpkeepGasUsageWrapper1_2__factory as GasWrapperFactory } from '../../../typechain/factories/KeeperRegistryCheckUpkeepGasUsageWrapper1_2__factory'
import { getUsers, Personas } from '../../test-helpers/setup'
import {
  deployMockContract,
  MockContract,
} from '@ethereum-waffle/mock-contract'
import { KeeperRegistry1_2__factory as KeeperRegistryFactory } from '../../../typechain/factories/KeeperRegistry1_2__factory'

let personas: Personas
let owner: Signer
let caller: Signer
let nelly: Signer
let registryMockContract: MockContract
let gasWrapper: GasWrapper
let gasWrapperFactory: GasWrapperFactory

const upkeepId = 123

describe('KeeperRegistryCheckUpkeepGasUsageWrapper1_2', () => {
  before(async () => {
    personas = (await getUsers()).personas
    owner = personas.Default
    caller = personas.Carol
    nelly = personas.Nelly

    registryMockContract = await deployMockContract(
      owner as any,
      KeeperRegistryFactory.abi,
    )
    // @ts-ignore bug in autogen file
    gasWrapperFactory = await ethers.getContractFactory(
      'KeeperRegistryCheckUpkeepGasUsageWrapper1_2',
    )
    gasWrapper = await gasWrapperFactory
      .connect(owner)
      .deploy(registryMockContract.address)
    await gasWrapper.deployed()
  })

  describe('measureCheckGas()', () => {
    it("returns gas used when registry's checkUpkeep executes successfully", async () => {
      await registryMockContract.mock.checkUpkeep
        .withArgs(upkeepId, await nelly.getAddress())
        .returns(
          '0xabcd' /* performData */,
          BigNumber.from(1000) /* maxLinkPayment */,
          BigNumber.from(2000) /* gasLimit */,
          BigNumber.from(3000) /* adjustedGasWei */,
          BigNumber.from(4000) /* linkEth */,
        )

      const response = await gasWrapper
        .connect(caller)
        .callStatic.measureCheckGas(
          BigNumber.from(upkeepId),
          await nelly.getAddress(),
        )

      assert.isTrue(response[0], 'The checkUpkeepSuccess should be true')
      assert.equal(
        response[1],
        '0xabcd',
        'The performData should be forwarded correctly',
      )
      assert.isTrue(
        response[2] > BigNumber.from(0),
        'The gasUsed value must be larger than 0',
      )
    })

    it("returns gas used when registry's checkUpkeep reverts", async () => {
      await registryMockContract.mock.checkUpkeep
        .withArgs(upkeepId, await nelly.getAddress())
        .revertsWithReason('Error')

      const response = await gasWrapper
        .connect(caller)
        .callStatic.measureCheckGas(
          BigNumber.from(upkeepId),
          await nelly.getAddress(),
        )

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
      const registry = await gasWrapper.connect(caller).getKeeperRegistry()
      assert.equal(
        registry,
        registryMockContract.address,
        'The underlying keeper registry is incorrect',
      )
    })
  })
})
