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
import { KeeperRegistryLogic20__factory as KeeperRegistryLogicFactory } from '../../../typechain/factories/KeeperRegistryLogic20__factory'
import { toWei } from '../../test-helpers/helpers'

let upkeepMockFactory: UpkeepMockFactory
let upkeepTranscoderFactory: UpkeepTranscoderFactory
let transcoder: UpkeepTranscoder
let linkTokenFactory: LinkTokenFactory
let mockV3AggregatorFactory: MockV3AggregatorFactory
let keeperRegistryFactory20: KeeperRegistry20Factory
let keeperRegistryLogicFactory: KeeperRegistryLogicFactory
let personas: Personas
let owner: Signer
let upkeepsV1: any[]
let upkeepsV2: any[]
let upkeepsV3: any[]
let admins: string[]
let admin0: Signer
let admin1: Signer
const balance = 50000000000000
const executeGas = 200000
const amountSpent = 200000000000000
const target0 = '0xffffffffffffffffffffffffffffffffffffffff'
const target1 = '0xfffffffffffffffffffffffffffffffffffffffe'
const lastKeeper0 = '0x233a95ccebf3c9f934482c637c08b4015cdd6ddd'
const lastKeeper1 = '0x233a95ccebf3c9f934482c637c08b4015cdd6ddc'
const UpkeepFormatV1 = 0
const UpkeepFormatV2 = 1
const UpkeepFormatV3 = 2
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

describe.only('UpkeepTranscoder3_0', () => {
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

    it('reverts if the from type is not an enum value', async () => {
      await evmRevert(
        transcoder.transcodeUpkeeps(3, 2, encodedData),
        'function was called with incorrect parameters',
      )
    })

    it('reverts if the to version is not 2', async () => {
      await evmRevert(
        transcoder.transcodeUpkeeps(1, 3, encodedData),
        'function was called with incorrect parameters',
      )
    })

    context('when from and to versions equal', () => {
      it('returns the data that was passed in', async () => {
        let response = await transcoder.transcodeUpkeeps(
          UpkeepFormatV1,
          UpkeepFormatV1,
          encodedData,
        )
        assert.equal(encodedData, response)

        response = await transcoder.transcodeUpkeeps(
          UpkeepFormatV2,
          UpkeepFormatV2,
          encodedData,
        )
        assert.equal(encodedData, response)

        response = await transcoder.transcodeUpkeeps(
          UpkeepFormatV3,
          UpkeepFormatV3,
          encodedData,
        )
        assert.equal(encodedData, response)
      })
    })

    context('when from and to versions are correct', () => {
      upkeepsV3 = [
        [executeGas, 2 ** 32 - 1, false, target0, amountSpent, balance, 0],
        [executeGas, 2 ** 32 - 1, false, target1, amountSpent, balance, 0],
      ]

      it('transcodes V1 upkeeps to V3 properly', async () => {
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

        let data = await transcoder.transcodeUpkeeps(
          UpkeepFormatV1,
          UpkeepFormatV3,
          encodeUpkeepV1(idx, upkeepsV1, ['0xabcd', '0xffff']),
        )
        assert.equal(
          encodeUpkeepV3(idx, upkeepsV3, ['0xabcd', '0xffff'], admins),
          data,
        )
      })

      it('transcodes V2 upkeeps to V3 properly', async () => {
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

        let data = await transcoder.transcodeUpkeeps(
          UpkeepFormatV2,
          UpkeepFormatV3,
          encodeUpkeepV2(idx, upkeepsV2, ['0xabcd', '0xffff']),
        )
        assert.equal(
          encodeUpkeepV3(idx, upkeepsV3, ['0xabcd', '0xffff'], admins),
          data,
        )
      })

      it('migrates upkeeps from one registry to another', async () => {
        let mock = await upkeepMockFactory.deploy()
        let executeGas = BigNumber.from('100000')
        let paymentPremiumPPB = BigNumber.from('250000000')
        let flatFeeMicroLink = BigNumber.from(0)
        let blockCountPerTurn = BigNumber.from(3)
        let randomBytes = '0x1234abcd'
        let stalenessSeconds = BigNumber.from(43820)
        let gasCeilingMultiplier = BigNumber.from(1)
        let checkGasLimit = BigNumber.from(20000000)
        let fallbackGasPrice = BigNumber.from(200)
        let fallbackLinkPrice = BigNumber.from(200000000)
        let maxPerformGas = BigNumber.from(5000000)
        let minUpkeepSpend = BigNumber.from(0)
        let maxCheckDataSize = BigNumber.from(1000)
        let maxPerformDataSize = BigNumber.from(1000)
        let paymentModel = BigNumber.from(0)
        let linkEth = BigNumber.from(300000000)
        let gasWei = BigNumber.from(100)
        // @ts-ignore bug in autogen file
        let keeperRegistryFactory = await ethers.getContractFactory(
          'KeeperRegistry1_2',
        )
        let linkToken = await linkTokenFactory.connect(owner).deploy()
        let gasPriceFeed = await mockV3AggregatorFactory
          .connect(owner)
          .deploy(0, gasWei)
        let linkEthFeed = await mockV3AggregatorFactory
          .connect(owner)
          .deploy(9, linkEth)
        transcoder = await upkeepTranscoderFactory.connect(owner).deploy()
        let registry12 = await keeperRegistryFactory
          .connect(owner)
          .deploy(
            linkToken.address,
            linkEthFeed.address,
            gasPriceFeed.address,
            {
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
            },
          )
        const tx = await registry12
          .connect(owner)
          .registerUpkeep(
            mock.address,
            executeGas,
            await admin0.getAddress(),
            randomBytes,
          )
        const id = await getUpkeepID(tx)

        // @ts-ignore bug in autogen file
        keeperRegistryFactory20 = await ethers.getContractFactory(
          'KeeperRegistry2_0',
        )
        // @ts-ignore bug in autogen file
        keeperRegistryLogicFactory = await ethers.getContractFactory(
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

        let registryLogic = await keeperRegistryLogicFactory
          .connect(owner)
          .deploy(
            paymentModel,
            linkToken.address,
            linkEthFeed.address,
            gasPriceFeed.address,
          )

        let registry20 = await keeperRegistryFactory20
          .connect(owner)
          .deploy(registryLogic.address)

        // deploys a registry, setups of initial configuration, registers an upkeep
        let keeper1 = personas.Carol
        let keeper2 = personas.Eddy
        let keeper3 = personas.Nancy
        let keeper4 = personas.Norbert
        let keeper5 = personas.Nick
        let payee1 = personas.Nelly
        let payee2 = personas.Norbert
        let payee3 = personas.Nick
        let payee4 = personas.Eddy
        let payee5 = personas.Carol
        // signers
        let signer1 = new ethers.Wallet(
          '0x7777777000000000000000000000000000000000000000000000000000000001',
        )
        let signer2 = new ethers.Wallet(
          '0x7777777000000000000000000000000000000000000000000000000000000002',
        )
        let signer3 = new ethers.Wallet(
          '0x7777777000000000000000000000000000000000000000000000000000000003',
        )
        let signer4 = new ethers.Wallet(
          '0x7777777000000000000000000000000000000000000000000000000000000004',
        )
        let signer5 = new ethers.Wallet(
          '0x7777777000000000000000000000000000000000000000000000000000000005',
        )

        let keeperAddresses = [
          await keeper1.getAddress(),
          await keeper2.getAddress(),
          await keeper3.getAddress(),
          await keeper4.getAddress(),
          await keeper5.getAddress(),
        ]
        let payees = [
          await payee1.getAddress(),
          await payee2.getAddress(),
          await payee3.getAddress(),
          await payee4.getAddress(),
          await payee5.getAddress(),
        ]
        let signers = [signer1, signer2, signer3, signer4, signer5]

        let signerAddresses = []
        for (let signer of signers) {
          signerAddresses.push(await signer.getAddress())
        }

        let f = 1
        let offchainVersion = 1
        let offchainBytes = '0x'

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
        await linkToken
          .connect(owner)
          .approve(registry12.address, toWei('1000'))
        await registry12.connect(owner).addFunds(id, toWei('1000'))

        // set outgoing permission to registry 2_0 and incoming permission for registry 1_2
        await registry12.setPeerRegistryMigrationPermission(
          registry20.address,
          1,
        )
        await registry20.setPeerRegistryMigrationPermission(
          registry12.address,
          2,
        )

        expect((await registry12.getUpkeep(id)).balance).to.equal(toWei('1000'))
        expect((await registry12.getUpkeep(id)).checkData).to.equal(randomBytes)
        expect((await registry12.getState()).state.numUpkeeps).to.equal(1)

        await registry12
          .connect(admin0)
          .migrateUpkeeps([id], registry20.address)

        expect((await registry12.getState()).state.numUpkeeps).to.equal(0)
        expect((await registry20.getState()).state.numUpkeeps).to.equal(1)
        expect((await registry12.getUpkeep(id)).balance).to.equal(0)
        expect((await registry12.getUpkeep(id)).checkData).to.equal('0x')
        expect((await registry20.getUpkeep(id)).balance).to.equal(toWei('1000'))
        expect(
          (await registry20.getState()).state.expectedLinkBalance,
        ).to.equal(toWei('1000'))
        expect((await registry20.getUpkeep(id)).checkData).to.equal(randomBytes)
      })
    })
  })
})
