import { ethers } from 'hardhat'
import { assert, expect } from 'chai'
import { UpkeepTranscoder30__factory as UpkeepTranscoderFactory } from '../../../typechain/factories/UpkeepTranscoder30__factory'
import { UpkeepTranscoder30 as UpkeepTranscoder } from '../../../typechain/UpkeepTranscoder30'
import { KeeperRegistry20__factory as KeeperRegistry20Factory } from '../../../typechain/factories/KeeperRegistry20__factory'
import { LinkToken__factory as LinkTokenFactory } from '../../../typechain/factories/LinkToken__factory'
import { MockV3Aggregator__factory as MockV3AggregatorFactory } from '../../../typechain/factories/MockV3Aggregator__factory'
import { UpkeepMock__factory as UpkeepMockFactory } from '../../../typechain/factories/UpkeepMock__factory'
import { evmRevert } from '../../test-helpers/matchers'
import { BigNumber, Signer } from 'ethers'
import { getUsers, Personas } from '../../test-helpers/setup'
import { KeeperRegistryLogic20__factory as KeeperRegistryLogic20Factory } from '../../../typechain/factories/KeeperRegistryLogic20__factory'
import { KeeperRegistry13__factory as KeeperRegistry13Factory } from '../../../typechain/factories/KeeperRegistry13__factory'
import { KeeperRegistryLogic13__factory as KeeperRegistryLogicFactory } from '../../../typechain/factories/KeeperRegistryLogic13__factory'
import { toWei } from '../../test-helpers/helpers'
import { LinkToken } from '../../../typechain'

let upkeepMockFactory: UpkeepMockFactory
let upkeepTranscoderFactory: UpkeepTranscoderFactory
let transcoder: UpkeepTranscoder
let linkTokenFactory: LinkTokenFactory
let mockV3AggregatorFactory: MockV3AggregatorFactory
let keeperRegistryFactory20: KeeperRegistry20Factory
let keeperRegistryFactory13: KeeperRegistry13Factory
let keeperRegistryLogicFactory20: KeeperRegistryLogic20Factory
let keeperRegistryLogicFactory13: KeeperRegistryLogicFactory
let personas: Personas
let owner: Signer
let upkeepsV1: any[]
let upkeepsV2: any[]
let upkeepsV3: any[]
let admins: string[]
let admin0: Signer
let admin1: Signer
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
const paymentModel = BigNumber.from(0)
const linkEth = BigNumber.from(300000000)
const gasWei = BigNumber.from(100)
const registryGasOverhead = BigNumber.from('80000')
const balance = 50000000000000
const amountSpent = 200000000000000
const target0 = '0xffffffffffffffffffffffffffffffffffffffff'
const target1 = '0xfffffffffffffffffffffffffffffffffffffffe'
const lastKeeper0 = '0x233a95ccebf3c9f934482c637c08b4015cdd6ddd'
const lastKeeper1 = '0x233a95ccebf3c9f934482c637c08b4015cdd6ddc'
enum UpkeepFormat {
  V1,
  V2,
  V3,
  V4,
}
const idx = [123, 124]

async function getUpkeepID(tx: any) {
  const receipt = await tx.wait()
  return receipt.events[0].args.id
}

const encodeConfig = (config: any) => {
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

const encodeUpkeepV1 = (ids: number[], upkeeps: any[], checkDatas: any[]) => {
  return ethers.utils.defaultAbiCoder.encode(
    [
      'uint256[]',
      'tuple(uint96,address,uint32,uint64,address,uint96,address)[]',
      'bytes[]',
    ],
    [ids, upkeeps, checkDatas],
  )
}

const encodeUpkeepV2 = (ids: number[], upkeeps: any[], checkDatas: any[]) => {
  return ethers.utils.defaultAbiCoder.encode(
    [
      'uint256[]',
      'tuple(uint96,address,uint96,address,uint32,uint32,address,bool)[]',
      'bytes[]',
    ],
    [ids, upkeeps, checkDatas],
  )
}

const encodeUpkeepV3 = (
  ids: number[],
  upkeeps: any[],
  checkDatas: any[],
  admins: string[],
) => {
  return ethers.utils.defaultAbiCoder.encode(
    [
      'uint256[]',
      'tuple(uint32,uint32,bool,address,uint96,uint96,uint32)[]',
      'bytes[]',
      'address[]',
    ],
    [ids, upkeeps, checkDatas, admins],
  )
}

before(async () => {
  // @ts-ignore bug in autogen file
  upkeepTranscoderFactory = await ethers.getContractFactory(
    'UpkeepTranscoder3_0',
  )
  personas = (await getUsers()).personas

  linkTokenFactory = await ethers.getContractFactory('LinkToken')
  // need full path because there are two contracts with name MockV3Aggregator
  mockV3AggregatorFactory = (await ethers.getContractFactory(
    'src/v0.8/tests/MockV3Aggregator.sol:MockV3Aggregator',
  )) as unknown as MockV3AggregatorFactory

  upkeepMockFactory = await ethers.getContractFactory('UpkeepMock')

  owner = personas.Norbert
  admin0 = personas.Neil
  admin1 = personas.Nick
  admins = [
    (await admin0.getAddress()).toLowerCase(),
    (await admin1.getAddress()).toLowerCase(),
  ]
})

async function deployLinkToken() {
  return await linkTokenFactory.connect(owner).deploy()
}

async function deployFeeds() {
  return [
    await mockV3AggregatorFactory.connect(owner).deploy(0, gasWei),
    await mockV3AggregatorFactory.connect(owner).deploy(9, linkEth),
  ]
}

async function deployLegacyRegistry12(
  linkToken: LinkToken,
  gasPriceFeed: any,
  linkEthFeed: any,
) {
  const mock = await upkeepMockFactory.deploy()
  // @ts-ignore bug in autogen file
  const keeperRegistryFactory = await ethers.getContractFactory(
    'KeeperRegistry1_2',
  )
  transcoder = await upkeepTranscoderFactory.connect(owner).deploy()
  const legacyRegistry = await keeperRegistryFactory
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
  const tx = await legacyRegistry
    .connect(owner)
    .registerUpkeep(
      mock.address,
      executeGas,
      await admin0.getAddress(),
      randomBytes,
    )
  const id = await getUpkeepID(tx)
  return [id, legacyRegistry]
}

async function deployLegacyRegistry13(
  linkToken: LinkToken,
  gasPriceFeed: any,
  linkEthFeed: any,
) {
  const mock = await upkeepMockFactory.deploy()
  // @ts-ignore bug in autogen file
  keeperRegistryFactory13 = await ethers.getContractFactory('KeeperRegistry1_3')
  // @ts-ignore bug in autogen file
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

async function deployRegistry20(
  linkToken: LinkToken,
  gasPriceFeed: any,
  linkEthFeed: any,
) {
  // @ts-ignore bug in autogen file
  keeperRegistryFactory20 = await ethers.getContractFactory('KeeperRegistry2_0')
  // @ts-ignore bug in autogen file
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
    .deploy(
      paymentModel,
      linkToken.address,
      linkEthFeed.address,
      gasPriceFeed.address,
    )

  const registry20 = await keeperRegistryFactory20
    .connect(owner)
    .deploy(registryLogic.address)

  // deploys a registry, setups of initial configuration, registers an upkeep
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

  const keeperAddresses = [
    await keeper1.getAddress(),
    await keeper2.getAddress(),
    await keeper3.getAddress(),
    await keeper4.getAddress(),
    await keeper5.getAddress(),
  ]
  const payees = [
    await payee1.getAddress(),
    await payee2.getAddress(),
    await payee3.getAddress(),
    await payee4.getAddress(),
    await payee5.getAddress(),
  ]
  const signers = [signer1, signer2, signer3, signer4, signer5]

  const signerAddresses = []
  for (const signer of signers) {
    signerAddresses.push(await signer.getAddress())
  }

  const f = 1
  const offchainVersion = 1
  const offchainBytes = '0x'

  await registry20
    .connect(owner)
    .setConfig(
      signerAddresses,
      keeperAddresses,
      f,
      encodeConfig(config),
      offchainVersion,
      offchainBytes,
    )
  await registry20.connect(owner).setPayees(payees)
  return registry20
}

describe('UpkeepTranscoder3_0', () => {
  beforeEach(async () => {
    transcoder = await upkeepTranscoderFactory.connect(owner).deploy()
  })

  describe('#typeAndVersion', () => {
    it('uses the correct type and version', async () => {
      const typeAndVersion = await transcoder.typeAndVersion()
      assert.equal(typeAndVersion, 'UpkeepTranscoder 3.0.0')
    })
  })

  describe('#transcodeUpkeeps', () => {
    const encodedData = '0xabcd'

    it('reverts if the from type is not V1 or V2', async () => {
      await evmRevert(
        transcoder.transcodeUpkeeps(
          UpkeepFormat.V3,
          UpkeepFormat.V1,
          encodedData,
        ),
      )
      await evmRevert(
        transcoder.transcodeUpkeeps(
          UpkeepFormat.V4,
          UpkeepFormat.V1,
          encodedData,
        ),
      )
    })

    context('when from and to versions are correct', () => {
      upkeepsV3 = [
        [executeGas, 2 ** 32 - 1, false, target0, amountSpent, balance, 0],
        [executeGas, 2 ** 32 - 1, false, target1, amountSpent, balance, 0],
      ]

      it('transcodes V1 upkeeps to V3 properly, regardless of toVersion value', async () => {
        upkeepsV1 = [
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

        const data = await transcoder.transcodeUpkeeps(
          UpkeepFormat.V1,
          UpkeepFormat.V1,
          encodeUpkeepV1(idx, upkeepsV1, ['0xabcd', '0xffff']),
        )
        assert.equal(
          encodeUpkeepV3(idx, upkeepsV3, ['0xabcd', '0xffff'], admins),
          data,
        )
      })

      it('transcodes V2 upkeeps to V3 properly, regardless of toVersion value', async () => {
        upkeepsV2 = [
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

        const data = await transcoder.transcodeUpkeeps(
          UpkeepFormat.V2,
          UpkeepFormat.V2,
          encodeUpkeepV2(idx, upkeepsV2, ['0xabcd', '0xffff']),
        )
        assert.equal(
          encodeUpkeepV3(idx, upkeepsV3, ['0xabcd', '0xffff'], admins),
          data,
        )
      })

      it('migrates upkeeps from 1.2 registry to 2.0', async () => {
        const linkToken = await deployLinkToken()
        const [gasPriceFeed, linkEthFeed] = await deployFeeds()
        const [id, legacyRegistry] = await deployLegacyRegistry12(
          linkToken,
          gasPriceFeed,
          linkEthFeed,
        )
        const registry20 = await deployRegistry20(
          linkToken,
          gasPriceFeed,
          linkEthFeed,
        )

        await linkToken
          .connect(owner)
          .approve(legacyRegistry.address, toWei('1000'))
        await legacyRegistry.connect(owner).addFunds(id, toWei('1000'))

        // set outgoing permission to registry 2_0 and incoming permission for registry 1_2
        await legacyRegistry.setPeerRegistryMigrationPermission(
          registry20.address,
          1,
        )
        await registry20.setPeerRegistryMigrationPermission(
          legacyRegistry.address,
          2,
        )

        expect((await legacyRegistry.getUpkeep(id)).balance).to.equal(
          toWei('1000'),
        )
        expect((await legacyRegistry.getUpkeep(id)).checkData).to.equal(
          randomBytes,
        )
        expect((await legacyRegistry.getState()).state.numUpkeeps).to.equal(1)

        await legacyRegistry
          .connect(admin0)
          .migrateUpkeeps([id], registry20.address)

        expect((await legacyRegistry.getState()).state.numUpkeeps).to.equal(0)
        expect((await registry20.getState()).state.numUpkeeps).to.equal(1)
        expect((await legacyRegistry.getUpkeep(id)).balance).to.equal(0)
        expect((await legacyRegistry.getUpkeep(id)).checkData).to.equal('0x')
        expect((await registry20.getUpkeep(id)).balance).to.equal(toWei('1000'))
        expect(
          (await registry20.getState()).state.expectedLinkBalance,
        ).to.equal(toWei('1000'))
        expect(await linkToken.balanceOf(registry20.address)).to.equal(
          toWei('1000'),
        )
        expect((await registry20.getUpkeep(id)).checkData).to.equal(randomBytes)
      })

      it('migrates upkeeps from 1.3 registry to 2.0', async () => {
        const linkToken = await deployLinkToken()
        const [gasPriceFeed, linkEthFeed] = await deployFeeds()
        const [id, legacyRegistry] = await deployLegacyRegistry13(
          linkToken,
          gasPriceFeed,
          linkEthFeed,
        )
        const registry20 = await deployRegistry20(
          linkToken,
          gasPriceFeed,
          linkEthFeed,
        )

        await linkToken
          .connect(owner)
          .approve(legacyRegistry.address, toWei('1000'))
        await legacyRegistry.connect(owner).addFunds(id, toWei('1000'))

        // set outgoing permission to registry 2_0 and incoming permission for registry 1_3
        await legacyRegistry.setPeerRegistryMigrationPermission(
          registry20.address,
          1,
        )
        await registry20.setPeerRegistryMigrationPermission(
          legacyRegistry.address,
          2,
        )

        expect((await legacyRegistry.getUpkeep(id)).balance).to.equal(
          toWei('1000'),
        )
        expect((await legacyRegistry.getUpkeep(id)).checkData).to.equal(
          randomBytes,
        )
        expect((await legacyRegistry.getState()).state.numUpkeeps).to.equal(1)

        await legacyRegistry
          .connect(admin0)
          .migrateUpkeeps([id], registry20.address)

        expect((await legacyRegistry.getState()).state.numUpkeeps).to.equal(0)
        expect((await registry20.getState()).state.numUpkeeps).to.equal(1)
        expect((await legacyRegistry.getUpkeep(id)).balance).to.equal(0)
        expect((await legacyRegistry.getUpkeep(id)).checkData).to.equal('0x')
        expect((await registry20.getUpkeep(id)).balance).to.equal(toWei('1000'))
        expect(
          (await registry20.getState()).state.expectedLinkBalance,
        ).to.equal(toWei('1000'))
        expect(await linkToken.balanceOf(registry20.address)).to.equal(
          toWei('1000'),
        )
        expect((await registry20.getUpkeep(id)).checkData).to.equal(randomBytes)
      })
    })
  })
})
