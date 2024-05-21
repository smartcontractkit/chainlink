import { ethers } from 'hardhat'
import { assert, expect } from 'chai'
import { UpkeepTranscoder4_0 as UpkeepTranscoder } from '../../../typechain/UpkeepTranscoder4_0'
import { KeeperRegistry2_0__factory as KeeperRegistry2_0Factory } from '../../../typechain/factories/KeeperRegistry2_0__factory'
import { LinkToken__factory as LinkTokenFactory } from '../../../typechain/factories/LinkToken__factory'
import { MockV3Aggregator__factory as MockV3AggregatorFactory } from '../../../typechain/factories/MockV3Aggregator__factory'
import { evmRevert } from '../../test-helpers/matchers'
import { BigNumber, Signer } from 'ethers'
import { getUsers, Personas } from '../../test-helpers/setup'
import { KeeperRegistryLogic2_0__factory as KeeperRegistryLogic20Factory } from '../../../typechain/factories/KeeperRegistryLogic2_0__factory'
import { KeeperRegistry1_3__factory as KeeperRegistry1_3Factory } from '../../../typechain/factories/KeeperRegistry1_3__factory'
import { KeeperRegistryLogic1_3__factory as KeeperRegistryLogicFactory } from '../../../typechain/factories/KeeperRegistryLogic1_3__factory'
import { UpkeepTranscoder4_0__factory as UpkeepTranscoderFactory } from '../../../typechain/factories/UpkeepTranscoder4_0__factory'
import { toWei } from '../../test-helpers/helpers'
import { loadFixture } from '@nomicfoundation/hardhat-network-helpers'
import {
  IKeeperRegistryMaster,
  KeeperRegistry1_2,
  KeeperRegistry1_3,
  KeeperRegistry2_0,
  LinkToken,
  MockV3Aggregator,
  UpkeepMock,
} from '../../../typechain'
import { deployRegistry21 } from './helpers'

//////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////

/*********************************** TRANSCODER v4.0 IS FROZEN ************************************/

// We are leaving the original tests enabled, however as automation v2.1 is still actively being deployed

describe('UpkeepTranscoder v4.0 - Frozen [ @skip-coverage ]', () => {
  it('has not changed', () => {
    assert.equal(
      ethers.utils.id(UpkeepTranscoderFactory.bytecode),
      '0xf22c4701b0088e6e69c389a34a22041a69f00890a89246e3c2a6d38172222dae',
      'UpkeepTranscoder bytecode has changed',
    )
  })
})

//////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////

let transcoder: UpkeepTranscoder
let linkTokenFactory: LinkTokenFactory
let keeperRegistryFactory20: KeeperRegistry2_0Factory
let keeperRegistryFactory13: KeeperRegistry1_3Factory
let keeperRegistryLogicFactory20: KeeperRegistryLogic20Factory
let keeperRegistryLogicFactory13: KeeperRegistryLogicFactory
let linkToken: LinkToken
let registry12: KeeperRegistry1_2
let registry13: KeeperRegistry1_3
let registry20: KeeperRegistry2_0
let registry21: IKeeperRegistryMaster
let gasPriceFeed: MockV3Aggregator
let linkEthFeed: MockV3Aggregator
let mock: UpkeepMock
let personas: Personas
let owner: Signer
let upkeepsV12: any[]
let upkeepsV13: any[]
let upkeepsV21: any[]
let admins: string[]
let admin0: Signer
let admin1: Signer
let id12: BigNumber
let id13: BigNumber
let id20: BigNumber
const executeGas = BigNumber.from('100000')
const paymentPremiumPPB = BigNumber.from('250000000')
const flatFeeMicroLink = BigNumber.from(0)
const blockCountPerTurn = BigNumber.from(3)
const randomBytes = '0x1234abcd'
const stalenessSeconds = BigNumber.from(43820)
const gasCeilingMultiplier = BigNumber.from(1)
const checkGasLimit = BigNumber.from(20000000)
const fallbackGasPrice = BigNumber.from(200)
const fallbackLinkPrice = BigNumber.from(200000000)
const maxPerformGas = BigNumber.from(5000000)
const minUpkeepSpend = BigNumber.from(0)
const maxCheckDataSize = BigNumber.from(1000)
const maxPerformDataSize = BigNumber.from(1000)
const mode = BigNumber.from(0)
const linkEth = BigNumber.from(300000000)
const gasWei = BigNumber.from(100)
const registryGasOverhead = BigNumber.from('80000')
const balance = 50000000000000
const amountSpent = 200000000000000
const { AddressZero } = ethers.constants
const target0 = '0xffffffffffffffffffffffffffffffffffffffff'
const target1 = '0xfffffffffffffffffffffffffffffffffffffffe'
const lastKeeper0 = '0x233a95ccebf3c9f934482c637c08b4015cdd6ddd'
const lastKeeper1 = '0x233a95ccebf3c9f934482c637c08b4015cdd6ddc'

const f = 1
const offchainVersion = 1
const offchainBytes = '0x'
let keeperAddresses: string[]
let signerAddresses: string[]
let payees: string[]

enum UpkeepFormat {
  V12,
  V13,
  V20,
  V21,
  V30, // Does not exist
}
const idx = [123, 124]

async function getUpkeepID(tx: any): Promise<BigNumber> {
  const receipt = await tx.wait()
  return receipt.events[0].args.id
}

const encodeConfig20 = (config: any) => {
  return ethers.utils.defaultAbiCoder.encode(
    [
      'tuple(uint32 paymentPremiumPPB,uint32 flatFeeMicroLink,uint32 checkGasLimit,uint24 stalenessSeconds\
        ,uint16 gasCeilingMultiplier,uint96 minUpkeepSpend,uint32 maxPerformGas,uint32 maxCheckDataSize,\
        uint32 maxPerformDataSize,uint256 fallbackGasPrice,uint256 fallbackLinkPrice,address transcoder,\
        address registrar)',
    ],
    [config],
  )
}

const encodeUpkeepV12 = (ids: number[], upkeeps: any[], checkDatas: any[]) => {
  return ethers.utils.defaultAbiCoder.encode(
    [
      'uint256[]',
      'tuple(uint96,address,uint32,uint64,address,uint96,address)[]',
      'bytes[]',
    ],
    [ids, upkeeps, checkDatas],
  )
}

async function deployRegistry1_2(): Promise<[BigNumber, KeeperRegistry1_2]> {
  const keeperRegistryFactory =
    await ethers.getContractFactory('KeeperRegistry1_2')
  const registry12 = await keeperRegistryFactory
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
  const tx = await registry12
    .connect(owner)
    .registerUpkeep(
      mock.address,
      executeGas,
      await admin0.getAddress(),
      randomBytes,
    )
  const id = await getUpkeepID(tx)
  return [id, registry12]
}

async function deployRegistry1_3(): Promise<[BigNumber, KeeperRegistry1_3]> {
  keeperRegistryFactory13 = await ethers.getContractFactory('KeeperRegistry1_3')
  keeperRegistryLogicFactory13 = await ethers.getContractFactory(
    'KeeperRegistryLogic1_3',
  )

  const registryLogic13 = await keeperRegistryLogicFactory13
    .connect(owner)
    .deploy(
      0,
      registryGasOverhead,
      linkToken.address,
      linkEthFeed.address,
      gasPriceFeed.address,
    )

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
    transcoder: transcoder.address,
    registrar: ethers.constants.AddressZero,
  }
  const registry13 = await keeperRegistryFactory13
    .connect(owner)
    .deploy(registryLogic13.address, config)

  const tx = await registry13
    .connect(owner)
    .registerUpkeep(
      mock.address,
      executeGas,
      await admin0.getAddress(),
      randomBytes,
    )
  const id = await getUpkeepID(tx)

  return [id, registry13]
}

async function deployRegistry2_0(): Promise<[BigNumber, KeeperRegistry2_0]> {
  keeperRegistryFactory20 = await ethers.getContractFactory('KeeperRegistry2_0')
  keeperRegistryLogicFactory20 = await ethers.getContractFactory(
    'KeeperRegistryLogic2_0',
  )

  const config = {
    paymentPremiumPPB,
    flatFeeMicroLink,
    checkGasLimit,
    stalenessSeconds,
    gasCeilingMultiplier,
    minUpkeepSpend,
    maxCheckDataSize,
    maxPerformDataSize,
    maxPerformGas,
    fallbackGasPrice,
    fallbackLinkPrice,
    transcoder: transcoder.address,
    registrar: ethers.constants.AddressZero,
  }

  const registryLogic = await keeperRegistryLogicFactory20
    .connect(owner)
    .deploy(mode, linkToken.address, linkEthFeed.address, gasPriceFeed.address)

  const registry20 = await keeperRegistryFactory20
    .connect(owner)
    .deploy(registryLogic.address)

  await registry20
    .connect(owner)
    .setConfig(
      signerAddresses,
      keeperAddresses,
      f,
      encodeConfig20(config),
      offchainVersion,
      offchainBytes,
    )
  await registry20.connect(owner).setPayees(payees)

  const tx = await registry20
    .connect(owner)
    .registerUpkeep(
      mock.address,
      executeGas,
      await admin0.getAddress(),
      randomBytes,
      randomBytes,
    )
  const id = await getUpkeepID(tx)

  return [id, registry20]
}

async function deployRegistry2_1() {
  const registry = await deployRegistry21(
    owner,
    mode,
    linkToken.address,
    linkEthFeed.address,
    gasPriceFeed.address,
  )

  const onchainConfig = {
    paymentPremiumPPB,
    flatFeeMicroLink,
    checkGasLimit,
    stalenessSeconds,
    gasCeilingMultiplier,
    minUpkeepSpend,
    maxCheckDataSize,
    maxPerformDataSize,
    maxRevertDataSize: 1000,
    maxPerformGas,
    fallbackGasPrice,
    fallbackLinkPrice,
    transcoder: ethers.constants.AddressZero,
    registrars: [],
    upkeepPrivilegeManager: await owner.getAddress(),
  }

  await registry
    .connect(owner)
    .setConfigTypeSafe(
      signerAddresses,
      keeperAddresses,
      f,
      onchainConfig,
      offchainVersion,
      offchainBytes,
    )

  return registry
}

const setup = async () => {
  personas = (await getUsers()).personas
  owner = personas.Norbert
  admin0 = personas.Neil
  admin1 = personas.Nick
  admins = [
    (await admin0.getAddress()).toLowerCase(),
    (await admin1.getAddress()).toLowerCase(),
  ]

  const upkeepTranscoderFactory = await ethers.getContractFactory(
    'UpkeepTranscoder4_0',
  )
  transcoder = await upkeepTranscoderFactory.connect(owner).deploy()

  linkTokenFactory = await ethers.getContractFactory(
    'src/v0.8/shared/test/helpers/LinkTokenTestHelper.sol:LinkTokenTestHelper',
  )
  linkToken = await linkTokenFactory.connect(owner).deploy()
  // need full path because there are two contracts with name MockV3Aggregator
  const mockV3AggregatorFactory = (await ethers.getContractFactory(
    'src/v0.8/tests/MockV3Aggregator.sol:MockV3Aggregator',
  )) as unknown as MockV3AggregatorFactory

  gasPriceFeed = await mockV3AggregatorFactory.connect(owner).deploy(0, gasWei)
  linkEthFeed = await mockV3AggregatorFactory.connect(owner).deploy(9, linkEth)

  const upkeepMockFactory = await ethers.getContractFactory('UpkeepMock')
  mock = await upkeepMockFactory.deploy()

  const keeper1 = personas.Carol
  const keeper2 = personas.Eddy
  const keeper3 = personas.Nancy
  const keeper4 = personas.Norbert
  const keeper5 = personas.Nick
  const payee1 = personas.Nelly
  const payee2 = personas.Norbert
  const payee3 = personas.Nick
  const payee4 = personas.Eddy
  const payee5 = personas.Carol
  // signers
  const signer1 = new ethers.Wallet(
    '0x7777777000000000000000000000000000000000000000000000000000000001',
  )
  const signer2 = new ethers.Wallet(
    '0x7777777000000000000000000000000000000000000000000000000000000002',
  )
  const signer3 = new ethers.Wallet(
    '0x7777777000000000000000000000000000000000000000000000000000000003',
  )
  const signer4 = new ethers.Wallet(
    '0x7777777000000000000000000000000000000000000000000000000000000004',
  )
  const signer5 = new ethers.Wallet(
    '0x7777777000000000000000000000000000000000000000000000000000000005',
  )

  keeperAddresses = [
    await keeper1.getAddress(),
    await keeper2.getAddress(),
    await keeper3.getAddress(),
    await keeper4.getAddress(),
    await keeper5.getAddress(),
  ]

  payees = [
    await payee1.getAddress(),
    await payee2.getAddress(),
    await payee3.getAddress(),
    await payee4.getAddress(),
    await payee5.getAddress(),
  ]
  const signers = [signer1, signer2, signer3, signer4, signer5]

  signerAddresses = signers.map((signer) => signer.address)
  ;[id12, registry12] = await deployRegistry1_2()
  ;[id13, registry13] = await deployRegistry1_3()
  ;[id20, registry20] = await deployRegistry2_0()
  registry21 = await deployRegistry2_1()

  upkeepsV12 = [
    [
      balance,
      lastKeeper0,
      executeGas,
      2 ** 32,
      target0,
      amountSpent,
      await admin0.getAddress(),
    ],
    [
      balance,
      lastKeeper1,
      executeGas,
      2 ** 32,
      target1,
      amountSpent,
      await admin1.getAddress(),
    ],
  ]

  upkeepsV13 = [
    [
      balance,
      lastKeeper0,
      amountSpent,
      await admin0.getAddress(),
      executeGas,
      2 ** 32 - 1,
      target0,
      false,
    ],
    [
      balance,
      lastKeeper1,
      amountSpent,
      await admin1.getAddress(),
      executeGas,
      2 ** 32 - 1,
      target1,
      false,
    ],
  ]

  upkeepsV21 = [
    [
      false,
      executeGas,
      2 ** 32 - 1,
      AddressZero, // forwarder will always be zero
      amountSpent,
      balance,
      0,
      target0,
    ],
    [
      false,
      executeGas,
      2 ** 32 - 1,
      AddressZero, // forwarder will always be zero
      amountSpent,
      balance,
      0,
      target1,
    ],
  ]
}

describe('UpkeepTranscoder4_0', () => {
  beforeEach(async () => {
    await loadFixture(setup)
  })

  describe('#typeAndVersion', () => {
    it('uses the correct type and version', async () => {
      const typeAndVersion = await transcoder.typeAndVersion()
      assert.equal(typeAndVersion, 'UpkeepTranscoder 4.0.0')
    })
  })

  describe('#transcodeUpkeeps', () => {
    const encodedData = '0xabcd'

    it('reverts if the from type is not v1.2, v1.3, v2.0, or v2.1', async () => {
      await evmRevert(
        transcoder.transcodeUpkeeps(
          UpkeepFormat.V30,
          UpkeepFormat.V12,
          encodedData,
        ),
      )
    })

    context('when from version is correct', () => {
      // note this is a bugfix - the "to" version should be accounted for in
      // future versions of the transcoder
      it('transcodes to v2.1, regardless of toVersion value', async () => {
        const data1 = await transcoder.transcodeUpkeeps(
          UpkeepFormat.V12,
          UpkeepFormat.V12,
          encodeUpkeepV12(idx, upkeepsV12, ['0xabcd', '0xffff']),
        )
        const data2 = await transcoder.transcodeUpkeeps(
          UpkeepFormat.V12,
          UpkeepFormat.V13,
          encodeUpkeepV12(idx, upkeepsV12, ['0xabcd', '0xffff']),
        )
        const data3 = await transcoder.transcodeUpkeeps(
          UpkeepFormat.V12,
          100,
          encodeUpkeepV12(idx, upkeepsV12, ['0xabcd', '0xffff']),
        )
        assert.equal(data1, data2)
        assert.equal(data1, data3)
      })

      it('migrates upkeeps from 1.2 registry to 2.1', async () => {
        await linkToken
          .connect(owner)
          .approve(registry12.address, toWei('1000'))
        await registry12.connect(owner).addFunds(id12, toWei('1000'))

        await registry12.setPeerRegistryMigrationPermission(
          registry21.address,
          1,
        )
        await registry21.setPeerRegistryMigrationPermission(
          registry12.address,
          2,
        )

        expect((await registry12.getUpkeep(id12)).balance).to.equal(
          toWei('1000'),
        )
        expect((await registry12.getUpkeep(id12)).checkData).to.equal(
          randomBytes,
        )
        expect((await registry12.getState()).state.numUpkeeps).to.equal(1)

        await registry12
          .connect(admin0)
          .migrateUpkeeps([id12], registry21.address)

        expect((await registry12.getState()).state.numUpkeeps).to.equal(0)
        expect((await registry21.getState()).state.numUpkeeps).to.equal(1)
        expect((await registry12.getUpkeep(id12)).balance).to.equal(0)
        expect((await registry12.getUpkeep(id12)).checkData).to.equal('0x')
        expect((await registry21.getUpkeep(id12)).balance).to.equal(
          toWei('1000'),
        )
        expect(
          (await registry21.getState()).state.expectedLinkBalance,
        ).to.equal(toWei('1000'))
        expect(await linkToken.balanceOf(registry21.address)).to.equal(
          toWei('1000'),
        )
        expect((await registry21.getUpkeep(id12)).checkData).to.equal(
          randomBytes,
        )
        expect((await registry21.getUpkeep(id12)).offchainConfig).to.equal('0x')
        expect(await registry21.getUpkeepTriggerConfig(id12)).to.equal('0x')
      })

      it('migrates upkeeps from 1.3 registry to 2.1', async () => {
        await linkToken
          .connect(owner)
          .approve(registry13.address, toWei('1000'))
        await registry13.connect(owner).addFunds(id13, toWei('1000'))

        await registry13.setPeerRegistryMigrationPermission(
          registry21.address,
          1,
        )
        await registry21.setPeerRegistryMigrationPermission(
          registry13.address,
          2,
        )

        expect((await registry13.getUpkeep(id13)).balance).to.equal(
          toWei('1000'),
        )
        expect((await registry13.getUpkeep(id13)).checkData).to.equal(
          randomBytes,
        )
        expect((await registry13.getState()).state.numUpkeeps).to.equal(1)

        await registry13
          .connect(admin0)
          .migrateUpkeeps([id13], registry21.address)

        expect((await registry13.getState()).state.numUpkeeps).to.equal(0)
        expect((await registry21.getState()).state.numUpkeeps).to.equal(1)
        expect((await registry13.getUpkeep(id13)).balance).to.equal(0)
        expect((await registry13.getUpkeep(id13)).checkData).to.equal('0x')
        expect((await registry21.getUpkeep(id13)).balance).to.equal(
          toWei('1000'),
        )
        expect(
          (await registry21.getState()).state.expectedLinkBalance,
        ).to.equal(toWei('1000'))
        expect(await linkToken.balanceOf(registry21.address)).to.equal(
          toWei('1000'),
        )
        expect((await registry21.getUpkeep(id13)).checkData).to.equal(
          randomBytes,
        )
        expect((await registry21.getUpkeep(id13)).offchainConfig).to.equal('0x')
        expect(await registry21.getUpkeepTriggerConfig(id13)).to.equal('0x')
      })

      it('migrates upkeeps from 2.0 registry to 2.1', async () => {
        await linkToken
          .connect(owner)
          .approve(registry20.address, toWei('1000'))
        await registry20.connect(owner).addFunds(id20, toWei('1000'))

        await registry20.setPeerRegistryMigrationPermission(
          registry21.address,
          1,
        )
        await registry21.setPeerRegistryMigrationPermission(
          registry20.address,
          2,
        )

        expect((await registry20.getUpkeep(id20)).balance).to.equal(
          toWei('1000'),
        )
        expect((await registry20.getUpkeep(id20)).checkData).to.equal(
          randomBytes,
        )
        expect((await registry20.getState()).state.numUpkeeps).to.equal(1)

        await registry20
          .connect(admin0)
          .migrateUpkeeps([id20], registry21.address)

        expect((await registry20.getState()).state.numUpkeeps).to.equal(0)
        expect((await registry21.getState()).state.numUpkeeps).to.equal(1)
        expect((await registry20.getUpkeep(id20)).balance).to.equal(0)
        expect((await registry20.getUpkeep(id20)).checkData).to.equal('0x')
        expect((await registry21.getUpkeep(id20)).balance).to.equal(
          toWei('1000'),
        )
        expect(
          (await registry21.getState()).state.expectedLinkBalance,
        ).to.equal(toWei('1000'))
        expect(await linkToken.balanceOf(registry21.address)).to.equal(
          toWei('1000'),
        )
        expect((await registry21.getUpkeep(id20)).checkData).to.equal(
          randomBytes,
        )
        expect(await registry21.getUpkeepTriggerConfig(id20)).to.equal('0x')
      })
    })
  })
})
