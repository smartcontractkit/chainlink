import { ethers } from 'hardhat'
import { assert, expect } from 'chai'
import { BigNumber, Signer, Wallet } from 'ethers'
import { evmRevert } from '../../../test-helpers/matchers'
import { getUsers, Personas } from '../../../test-helpers/setup'
import { toWei } from '../../../test-helpers/helpers'
import { LinkToken__factory as LinkTokenFactory } from '../../../../typechain/factories/LinkToken__factory'
import { MockV3Aggregator__factory as MockV3AggregatorFactory } from '../../../../typechain/factories/MockV3Aggregator__factory'
import { UpkeepMock__factory as UpkeepMockFactory } from '../../../../typechain/factories/UpkeepMock__factory'
import { UpkeepAutoFunder__factory as UpkeepAutoFunderFactory } from '../../../../typechain/factories/UpkeepAutoFunder__factory'
import { UpkeepTranscoder__factory as UpkeepTranscoderFactory } from '../../../../typechain/factories/UpkeepTranscoder__factory'
import { KeeperRegistry20__factory as KeeperRegistryFactory } from '../../../../typechain/factories/KeeperRegistry20__factory'
import { MockArbGasInfo__factory as MockArbGasInfoFactory } from '../../../../typechain/factories/MockArbGasInfo__factory'
import { MockOVMGasPriceOracle__factory as MockOVMGasPriceOracleFactory } from '../../../../typechain/factories/MockOVMGasPriceOracle__factory'
import { KeeperRegistryLogic20__factory as KeeperRegistryLogicFactory } from '../../../../typechain/factories/KeeperRegistryLogic20__factory'
import { KeeperRegistry20 as KeeperRegistry } from '../../../../typechain/KeeperRegistry20'
import { KeeperRegistryLogic20 as KeeperRegistryLogic } from '../../../../typechain/KeeperRegistryLogic20'
import { MockV3Aggregator } from '../../../../typechain/MockV3Aggregator'
import { LinkToken } from '../../../../typechain/LinkToken'
import { UpkeepMock } from '../../../../typechain/UpkeepMock'
import { MockArbGasInfo } from '../../../../typechain/MockArbGasInfo'
import { MockOVMGasPriceOracle } from '../../../../typechain/MockOVMGasPriceOracle'
import { UpkeepTranscoder } from '../../../../typechain/UpkeepTranscoder'

async function getUpkeepID(tx: any) {
  const receipt = await tx.wait()
  return receipt.events[0].args.id
}

function randomAddress() {
  return ethers.Wallet.createRandom().address
}

// -----------------------------------------------------------------------------------------------
// These are the gas overheads that off chain systems should provide to check upkeep / transmit
// These overheads are not actually charged for
const transmitGasOverhead = BigNumber.from(400000)
const checkGasOverhead = BigNumber.from(400000)

// These values should match the constants declared in registry
const registryGasOverhead = BigNumber.from(100000)
const verifySigOverhead = BigNumber.from(7500)
const cancellationDelay = 50

// This is the margin for gas that we test for. Gas charged should always be greater
// than total gas used in tx but should not increase beyond this margin
const gasCalculationMargin = BigNumber.from(4000)
// -----------------------------------------------------------------------------------------------

// Smart contract factories
let linkTokenFactory: LinkTokenFactory
let mockV3AggregatorFactory: MockV3AggregatorFactory
let keeperRegistryFactory: KeeperRegistryFactory
let keeperRegistryLogicFactory: KeeperRegistryLogicFactory
let upkeepMockFactory: UpkeepMockFactory
let upkeepAutoFunderFactory: UpkeepAutoFunderFactory
let upkeepTranscoderFactory: UpkeepTranscoderFactory
let mockArbGasInfoFactory: MockArbGasInfoFactory
let mockOVMGasPriceOracleFactory: MockOVMGasPriceOracleFactory
let personas: Personas

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

const encodeReport = (upkeeps: any) => {
  let upkeepIds = []
  let wrappedPerformDatas = []
  for (let i = 0; i < upkeeps.length; i++) {
    upkeepIds.push(upkeeps[i].Id)
    wrappedPerformDatas.push(
      ethers.utils.defaultAbiCoder.encode(
        ['tuple(uint32,bytes32,bytes)'],
        [
          [
            upkeeps[i].checkBlockNum,
            upkeeps[i].checkBlockHash,
            upkeeps[i].performData,
          ],
        ],
      ),
    )
  }
  return ethers.utils.defaultAbiCoder.encode(
    ['uint256[]', 'bytes[]'],
    [upkeepIds, wrappedPerformDatas],
  )
}

const encodeLatestBlockReport = async (upkeeps: any) => {
  let latestBlock = await ethers.provider.getBlock('latest')
  for (let i = 0; i < upkeeps.length; i++) {
    upkeeps[i].checkBlockNum = latestBlock.number + 1
    upkeeps[i].checkBlockHash = latestBlock.hash
    upkeeps[i].performData = '0x'
  }
  return encodeReport(upkeeps)
}

const signReport = (
  reportContext: string[],
  report: any,
  signers: Wallet[],
) => {
  const reportDigest = ethers.utils.keccak256(report)
  const packedArgs = ethers.utils.solidityPack(
    ['bytes32', 'bytes32[3]'],
    [reportDigest, reportContext],
  )
  const packedDigest = ethers.utils.keccak256(packedArgs)

  let signatures = []
  for (let signer of signers) {
    signatures.push(signer._signingKey().signDigest(packedDigest))
  }
  const vs = signatures.map((i) => '0' + (i.v - 27).toString(16)).join('')
  return {
    vs: '0x' + vs.padEnd(64, '0'),
    rs: signatures.map((i) => i.r),
    ss: signatures.map((i) => i.s),
  }
}

const parseUpkeepPerformedLogs = (receipt: any) => {
  let upkeepPerformedABI = [
    'event UpkeepPerformed(uint256 indexed id,bool indexed success, \
  uint32 checkBlockNumber,uint256 gasUsed,uint256 gasOverhead,uint96 totalPayment)',
  ]
  let iface = new ethers.utils.Interface(upkeepPerformedABI)

  let parsedLogs = []
  for (let i = 0; i < receipt.logs.length; i++) {
    const log = receipt.logs[i]
    try {
      parsedLogs.push(iface.parseLog(log))
    } catch (e) {
      // ignore log
    }
  }
  return parsedLogs
}

before(async () => {
  personas = (await getUsers()).personas

  linkTokenFactory = await ethers.getContractFactory('LinkToken')
  // need full path because there are two contracts with name MockV3Aggregator
  mockV3AggregatorFactory = (await ethers.getContractFactory(
    'src/v0.8/tests/MockV3Aggregator.sol:MockV3Aggregator',
  )) as unknown as MockV3AggregatorFactory
  // @ts-ignore bug in autogen file
  keeperRegistryFactory = await ethers.getContractFactory('KeeperRegistry2_0')
  // @ts-ignore bug in autogen file
  keeperRegistryLogicFactory = await ethers.getContractFactory(
    'KeeperRegistryLogic2_0',
  )
  upkeepMockFactory = await ethers.getContractFactory('UpkeepMock')
  upkeepAutoFunderFactory = await ethers.getContractFactory('UpkeepAutoFunder')
  upkeepTranscoderFactory = await ethers.getContractFactory('UpkeepTranscoder')
  mockArbGasInfoFactory = await ethers.getContractFactory('MockArbGasInfo')
  mockOVMGasPriceOracleFactory = await ethers.getContractFactory(
    'MockOVMGasPriceOracle',
  )
})

describe('KeeperRegistry2_0', () => {
  const linkEth = BigNumber.from(5000000000000000) // 1 Link = 0.005 Eth
  const gasWei = BigNumber.from(1000000000) // 1 gwei
  const linkDivisibility = BigNumber.from('1000000000000000000')
  const executeGas = BigNumber.from('1000000')
  const paymentPremiumBase = BigNumber.from('1000000000')
  const paymentPremiumPPB = BigNumber.from('250000000')
  const flatFeeMicroLink = BigNumber.from(0)

  const randomBytes = '0x1234abcd'
  const emptyBytes = '0x'
  const emptyBytes32 =
    '0x0000000000000000000000000000000000000000000000000000000000000000'

  const stalenessSeconds = BigNumber.from(43820)
  const gasCeilingMultiplier = BigNumber.from(2)
  const checkGasLimit = BigNumber.from(10000000)
  const fallbackGasPrice = gasWei.mul(BigNumber.from('2'))
  const fallbackLinkPrice = linkEth.div(BigNumber.from('2'))
  const maxCheckDataSize = BigNumber.from(1000)
  const maxPerformDataSize = BigNumber.from(1000)
  const maxPerformGas = BigNumber.from(5000000)
  const minUpkeepSpend = BigNumber.from(0)
  const f = 1
  const offchainVersion = 1
  const offchainBytes = '0x'
  const zeroAddress = ethers.constants.AddressZero
  const epochAndRound5_1 =
    '0x0000000000000000000000000000000000000000000000000000000000000501'

  let owner: Signer
  let keeper1: Signer
  let keeper2: Signer
  let keeper3: Signer
  let keeper4: Signer
  let keeper5: Signer
  let nonkeeper: Signer
  let signer1: Wallet
  let signer2: Wallet
  let signer3: Wallet
  let signer4: Wallet
  let signer5: Wallet
  let admin: Signer
  let payee1: Signer
  let payee2: Signer
  let payee3: Signer
  let payee4: Signer
  let payee5: Signer

  let linkToken: LinkToken
  let linkEthFeed: MockV3Aggregator
  let gasPriceFeed: MockV3Aggregator
  let registry: KeeperRegistry
  let registryLogic: KeeperRegistryLogic
  let mock: UpkeepMock
  let transcoder: UpkeepTranscoder
  let mockArbGasInfo: MockArbGasInfo
  let mockOVMGasPriceOracle: MockOVMGasPriceOracle

  let upkeepId: BigNumber
  let keeperAddresses: string[]
  let payees: string[]
  let signers: Wallet[]
  let signerAddresses: string[]
  let config: any

  const linkForGas = (
    upkeepGasSpent: BigNumber,
    gasOverhead: BigNumber,
    gasMultiplier: BigNumber,
    premiumPPB: BigNumber,
    flatFee: BigNumber,
    l1CostWei?: BigNumber,
    numUpkeepsBatch?: BigNumber,
  ) => {
    l1CostWei = l1CostWei === undefined ? BigNumber.from(0) : l1CostWei
    numUpkeepsBatch =
      numUpkeepsBatch === undefined ? BigNumber.from(1) : numUpkeepsBatch

    const gasSpent = gasOverhead.add(BigNumber.from(upkeepGasSpent))
    const base = gasWei
      .mul(gasMultiplier)
      .mul(gasSpent)
      .mul(linkDivisibility)
      .div(linkEth)
    const l1Fee = l1CostWei
      .mul(gasMultiplier)
      .div(numUpkeepsBatch)
      .mul(linkDivisibility)
      .div(linkEth)
    const premium = base.add(l1Fee).mul(premiumPPB).div(paymentPremiumBase)
    const flatFeeJules = BigNumber.from(flatFee).mul('1000000000000')
    return base.add(premium).add(flatFeeJules).add(l1Fee)
  }

  const verifyMaxPayment = async (
    paymentModel: number,
    multipliers: BigNumber[],
    gasAmounts: number[],
    premiums: number[],
    flatFees: number[],
    l1CostWei?: BigNumber,
  ) => {
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

    // Deploy a new registry since we change payment model
    let registryLogic = await keeperRegistryLogicFactory
      .connect(owner)
      .deploy(
        paymentModel,
        linkToken.address,
        linkEthFeed.address,
        gasPriceFeed.address,
      )
    // Deploy a new registry since we change payment model
    let registry = await keeperRegistryFactory
      .connect(owner)
      .deploy(registryLogic.address)
    await registry
      .connect(owner)
      .setConfig(
        signerAddresses,
        keeperAddresses,
        f,
        encodeConfig(config),
        offchainVersion,
        offchainBytes,
      )

    let fPlusOne = BigNumber.from(f + 1)
    let totalGasOverhead = registryGasOverhead
      .add(verifySigOverhead.mul(fPlusOne))
      .add(BigNumber.from('16').mul(maxPerformDataSize))

    for (let idx = 0; idx < gasAmounts.length; idx++) {
      const gas = gasAmounts[idx]
      for (let jdx = 0; jdx < premiums.length; jdx++) {
        const premium = premiums[jdx]
        for (let kdx = 0; kdx < flatFees.length; kdx++) {
          const flatFee = flatFees[kdx]
          for (let ldx = 0; ldx < multipliers.length; ldx++) {
            const multiplier = multipliers[ldx]

            await registry.connect(owner).setConfig(
              signerAddresses,
              keeperAddresses,
              f,
              encodeConfig({
                paymentPremiumPPB: premium,
                flatFeeMicroLink: flatFee,
                checkGasLimit,
                stalenessSeconds,
                gasCeilingMultiplier: multiplier,
                minUpkeepSpend,
                maxCheckDataSize,
                maxPerformDataSize,
                maxPerformGas,
                fallbackGasPrice,
                fallbackLinkPrice,
                transcoder: transcoder.address,
                registrar: ethers.constants.AddressZero,
              }),
              offchainVersion,
              offchainBytes,
            )

            const price = await registry.getMaxPaymentForGas(gas)
            expect(price).to.equal(
              linkForGas(
                BigNumber.from(gas),
                totalGasOverhead,
                multiplier,
                BigNumber.from(premium),
                BigNumber.from(flatFee),
                l1CostWei,
              ),
            )
          }
        }
      }
    }
  }

  beforeEach(async () => {
    // Deploys a registry, setups of initial configuration
    // Registers an upkeep which is unfunded to start with
    owner = personas.Default
    keeper1 = personas.Carol
    keeper2 = personas.Eddy
    keeper3 = personas.Nancy
    keeper4 = personas.Norbert
    keeper5 = personas.Nick
    nonkeeper = personas.Ned
    admin = personas.Neil
    payee1 = personas.Nelly
    payee2 = personas.Norbert
    payee3 = personas.Nick
    payee4 = personas.Eddy
    payee5 = personas.Carol
    // signers
    signer1 = new ethers.Wallet(
      '0x7777777000000000000000000000000000000000000000000000000000000001',
    )
    signer2 = new ethers.Wallet(
      '0x7777777000000000000000000000000000000000000000000000000000000002',
    )
    signer3 = new ethers.Wallet(
      '0x7777777000000000000000000000000000000000000000000000000000000003',
    )
    signer4 = new ethers.Wallet(
      '0x7777777000000000000000000000000000000000000000000000000000000004',
    )
    signer5 = new ethers.Wallet(
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
    signers = [signer1, signer2, signer3, signer4, signer5]

    // We append 26 random addresses to keepers, payees and signers to get a system of 31 oracles
    // This allows f value of 1 - 10
    for (let i = 0; i < 26; i++) {
      keeperAddresses.push(randomAddress())
      payees.push(randomAddress())
      signers.push(new ethers.Wallet(randomAddress()))
    }
    signerAddresses = []
    for (let signer of signers) {
      signerAddresses.push(await signer.getAddress())
    }

    linkToken = await linkTokenFactory.connect(owner).deploy()
    gasPriceFeed = await mockV3AggregatorFactory
      .connect(owner)
      .deploy(0, gasWei)
    linkEthFeed = await mockV3AggregatorFactory
      .connect(owner)
      .deploy(9, linkEth)
    transcoder = await upkeepTranscoderFactory.connect(owner).deploy()
    mockArbGasInfo = await mockArbGasInfoFactory.connect(owner).deploy()
    mockOVMGasPriceOracle = await mockOVMGasPriceOracleFactory
      .connect(owner)
      .deploy()

    const arbOracleCode = await ethers.provider.send('eth_getCode', [
      mockArbGasInfo.address,
    ])
    await ethers.provider.send('hardhat_setCode', [
      '0x000000000000000000000000000000000000006C',
      arbOracleCode,
    ])

    const optOracleCode = await ethers.provider.send('eth_getCode', [
      mockOVMGasPriceOracle.address,
    ])
    await ethers.provider.send('hardhat_setCode', [
      '0x420000000000000000000000000000000000000F',
      optOracleCode,
    ])

    registryLogic = await keeperRegistryLogicFactory
      .connect(owner)
      .deploy(0, linkToken.address, linkEthFeed.address, gasPriceFeed.address)

    config = {
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
    registry = await keeperRegistryFactory
      .connect(owner)
      .deploy(registryLogic.address)

    await registry
      .connect(owner)
      .setConfig(
        signerAddresses,
        keeperAddresses,
        f,
        encodeConfig(config),
        offchainVersion,
        offchainBytes,
      )
    await registry.connect(owner).setPayees(payees)

    mock = await upkeepMockFactory.deploy()
    await linkToken
      .connect(owner)
      .transfer(await admin.getAddress(), toWei('1000'))
    await linkToken.connect(admin).approve(registry.address, toWei('1000'))
    await linkToken.connect(owner).approve(registry.address, toWei('1000'))

    const tx = await registry
      .connect(owner)
      .registerUpkeep(
        mock.address,
        executeGas,
        await admin.getAddress(),
        true,
        randomBytes,
      )
    upkeepId = await getUpkeepID(tx)
  })

  describe('#transmit', () => {
    let sigVerificationUpkeepId: BigNumber
    const fArray = [1, 5, 10]
    beforeEach(async () => {
      const tx = await registry
        .connect(owner)
        .registerUpkeep(
          mock.address,
          executeGas,
          await admin.getAddress(),
          false,
          randomBytes,
        )
      sigVerificationUpkeepId = await getUpkeepID(tx)
    })

    it('reverts when registry is paused', async () => {
      await registry.connect(owner).pause()
      await evmRevert(
        registry
          .connect(keeper1)
          .transmit(
            [emptyBytes32, emptyBytes32, emptyBytes32],
            await encodeLatestBlockReport([{ Id: upkeepId.toString() }]),
            [],
            [],
            emptyBytes32,
          ),
        'RegistryPaused()',
      )
    })

    it('reverts when called by non active transmitter', async () => {
      await evmRevert(
        registry
          .connect(payee1)
          .transmit(
            [emptyBytes32, emptyBytes32, emptyBytes32],
            await encodeLatestBlockReport([{ Id: upkeepId.toString() }]),
            [],
            [],
            emptyBytes32,
          ),
        'OnlyActiveTransmitters()',
      )
    })

    it('reverts when upkeeps and performData length mismatches', async () => {
      let upkeepIds = []
      let wrappedPerformDatas = []
      let latestBlock = await ethers.provider.getBlock('latest')

      upkeepIds.push(upkeepId)
      wrappedPerformDatas.push(
        ethers.utils.defaultAbiCoder.encode(
          ['tuple(uint32,bytes32,bytes)'],
          [[latestBlock.number + 1, latestBlock.hash, '0x']],
        ),
      )
      // Push an extra perform data
      wrappedPerformDatas.push(
        ethers.utils.defaultAbiCoder.encode(
          ['tuple(uint32,bytes32,bytes)'],
          [[latestBlock.number + 1, latestBlock.hash, '0x']],
        ),
      )

      let report = ethers.utils.defaultAbiCoder.encode(
        ['uint256[]', 'bytes[]'],
        [upkeepIds, wrappedPerformDatas],
      )

      await evmRevert(
        registry
          .connect(keeper1)
          .transmit(
            [emptyBytes32, emptyBytes32, emptyBytes32],
            report,
            [],
            [],
            emptyBytes32,
          ),
        'InvalidReport()',
      )
    })

    it('reverts when wrappedPerformData is incorrectly encoded', async () => {
      let upkeepIds = []
      let wrappedPerformDatas = []
      let latestBlock = await ethers.provider.getBlock('latest')

      upkeepIds.push(upkeepId)
      wrappedPerformDatas.push(
        ethers.utils.defaultAbiCoder.encode(
          ['tuple(uint32,bytes32)'], // missing performData
          [[latestBlock.number + 1, latestBlock.hash]],
        ),
      )

      let report = ethers.utils.defaultAbiCoder.encode(
        ['uint256[]', 'bytes[]'],
        [upkeepIds, wrappedPerformDatas],
      )

      await evmRevert(
        registry
          .connect(keeper1)
          .transmit(
            [emptyBytes32, emptyBytes32, emptyBytes32],
            report,
            [],
            [],
            emptyBytes32,
          ),
      )
    })

    it('reverts when non sig verification and sig verification upkeeps are included in batch', async () => {
      await evmRevert(
        registry
          .connect(keeper1)
          .transmit(
            [emptyBytes32, emptyBytes32, emptyBytes32],
            await encodeLatestBlockReport([
              { Id: upkeepId.toString() },
              { Id: sigVerificationUpkeepId.toString() },
            ]),
            [],
            [],
            emptyBytes32,
          ),
        'InvalidReport()',
      )
    })

    it('reverts when no upkeeps are included in report', async () => {
      let upkeepIds: string[]
      let wrappedPerformDatas: string[]
      upkeepIds = []
      wrappedPerformDatas = []
      let report = ethers.utils.defaultAbiCoder.encode(
        ['uint256[]', 'bytes[]'],
        [upkeepIds, wrappedPerformDatas],
      )

      await evmRevert(
        registry
          .connect(keeper1)
          .transmit(
            [emptyBytes32, emptyBytes32, emptyBytes32],
            report,
            [],
            [],
            emptyBytes32,
          ),
        'StaleReport()',
      )
    })

    it('reverts when invalid upkeepIds are included in report', async () => {
      await evmRevert(
        registry
          .connect(keeper1)
          .transmit(
            [emptyBytes32, emptyBytes32, emptyBytes32],
            await encodeLatestBlockReport([
              { Id: upkeepId.add(BigNumber.from('1')).toString() },
            ]),
            [],
            [],
            emptyBytes32,
          ),
        'StaleReport()',
      )
    })

    it('reverts when duplicated upkeepIds are included in report', async () => {
      // Fund the upkeep so that pre-checks pass
      await registry.connect(admin).addFunds(upkeepId, toWei('100'))
      await evmRevert(
        registry
          .connect(keeper1)
          .transmit(
            [emptyBytes32, emptyBytes32, emptyBytes32],
            await encodeLatestBlockReport([
              { Id: upkeepId.toString() },
              { Id: upkeepId.toString() },
            ]),
            [],
            [],
            emptyBytes32,
          ),
        'InvalidReport()',
      )
    })

    it('reverts when upkeep has insufficient funds', async () => {
      await evmRevert(
        registry
          .connect(keeper1)
          .transmit(
            [emptyBytes32, emptyBytes32, emptyBytes32],
            await encodeLatestBlockReport([{ Id: upkeepId.toString() }]),
            [],
            [],
            emptyBytes32,
          ),
        'StaleReport()',
      )
    })

    context('When the upkeep is funded', async () => {
      beforeEach(async () => {
        // Fund the upkeep
        await registry.connect(admin).addFunds(upkeepId, toWei('100'))

        await registry
          .connect(admin)
          .addFunds(sigVerificationUpkeepId, toWei('100'))
      })

      it('reverts when check block number is less than last perform', async () => {
        // First perform an upkeep to put last perform block number on upkeep state

        const tx = await registry
          .connect(keeper1)
          .transmit(
            [emptyBytes32, emptyBytes32, emptyBytes32],
            await encodeLatestBlockReport([{ Id: upkeepId.toString() }]),
            [],
            [],
            emptyBytes32,
          )
        await tx.wait()

        let lastPerformBlockNumber = (await registry.getUpkeep(upkeepId))
          .lastPerformBlockNumber
        let lastPerformBlock = await ethers.provider.getBlock(
          lastPerformBlockNumber,
        )
        assert.equal(
          lastPerformBlockNumber.toString(),
          tx.blockNumber?.toString(),
        )

        // Try to transmit a report which has checkBlockNumber = lastPerformBlockNumber, should result in staleReport
        await evmRevert(
          registry.connect(keeper1).transmit(
            [emptyBytes32, emptyBytes32, emptyBytes32],
            await encodeReport([
              {
                Id: upkeepId.toString(),
                checkBlockNum: lastPerformBlockNumber,
                checkBlockHash: lastPerformBlock.parentHash,
                performData: '0x',
              },
            ]),
            [],
            [],
            emptyBytes32,
          ),
          'StaleReport()',
        )
      })

      it('reverts when check block hash does not match', async () => {
        await registry.connect(admin).addFunds(upkeepId, toWei('100'))
        let latestBlock = await ethers.provider.getBlock('latest')
        // Try to transmit a report which has incorrect checkBlockHash
        await evmRevert(
          registry.connect(keeper1).transmit(
            [emptyBytes32, emptyBytes32, emptyBytes32],
            await encodeReport([
              {
                Id: upkeepId.toString(),
                checkBlockNum: latestBlock.number,
                checkBlockHash: latestBlock.hash,
                performData: '0x',
              }, // should be latestBlock.parentHash
            ]),
            [],
            [],
            emptyBytes32,
          ),
          'StaleReport()',
        )
      })

      it('reverts when check block number is older than 256 blocks', async () => {
        let latestBlockReport = await encodeLatestBlockReport([
          { Id: upkeepId.toString() },
        ])

        await registry.connect(admin).cancelUpkeep(upkeepId)

        for (let i = 0; i < cancellationDelay; i++) {
          await ethers.provider.send('evm_mine', [])
        }

        await evmRevert(
          registry
            .connect(keeper1)
            .transmit(
              [emptyBytes32, emptyBytes32, emptyBytes32],
              latestBlockReport,
              [],
              [],
              emptyBytes32,
            ),
          'StaleReport()',
        )
      })

      it('reverts when upkeep is cancelled and cancellation delay has gone', async () => {
        let latestBlockReport = await encodeLatestBlockReport([
          { Id: upkeepId.toString() },
        ])
        await registry.connect(admin).cancelUpkeep(upkeepId)

        for (let i = 0; i < cancellationDelay; i++) {
          await ethers.provider.send('evm_mine', [])
        }

        // Try to transmit a report which is older than 256 blocks so block hash cannot be matched
        await evmRevert(
          registry
            .connect(keeper1)
            .transmit(
              [emptyBytes32, emptyBytes32, emptyBytes32],
              latestBlockReport,
              [],
              [],
              emptyBytes32,
            ),
          'StaleReport()',
        )
      })

      describe('When signatures are NOT validated', () => {
        it('does not revert if the target cannot execute', async () => {
          mock.setCanPerform(false)
          const tx = await registry.connect(keeper1).transmit(
            [emptyBytes32, emptyBytes32, emptyBytes32],
            await encodeLatestBlockReport([
              {
                Id: upkeepId.toString(),
              },
            ]),
            [],
            [],
            emptyBytes32,
          )
          const receipt = await tx.wait()
          let upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
          // exactly 1 Upkeep Performed should be emitted
          assert.equal(upkeepPerformedLogs.length, 1)
          let upkeepPerformedLog = upkeepPerformedLogs[0]

          let success = upkeepPerformedLog.args.success
          assert.equal(success, false)
        })

        it('reverts if not enough gas supplied', async () => {
          mock.setPerformGasToBurn(executeGas)
          await evmRevert(
            registry.connect(keeper1).transmit(
              [emptyBytes32, emptyBytes32, emptyBytes32],
              await encodeLatestBlockReport([
                {
                  Id: upkeepId.toString(),
                },
              ]),
              [],
              [],
              emptyBytes32,
              { gasLimit: executeGas },
            ),
          )
        })

        it('executes the data passed to the registry', async () => {
          mock.setCanPerform(true)
          let checkBlock = await ethers.provider.getBlock('latest')

          const tx = await registry.connect(keeper1).transmit(
            [emptyBytes32, emptyBytes32, emptyBytes32],
            encodeReport([
              {
                Id: upkeepId.toString(),
                checkBlockNum: checkBlock.number,
                checkBlockHash: checkBlock.parentHash,
                performData: randomBytes,
              },
            ]),
            [],
            [],
            emptyBytes32,
          )
          const receipt = await tx.wait()

          let upkeepPerformedWithABI = [
            'event UpkeepPerformedWith(bytes upkeepData)',
          ]
          let iface = new ethers.utils.Interface(upkeepPerformedWithABI)
          let parsedLogs = []
          for (let i = 0; i < receipt.logs.length; i++) {
            const log = receipt.logs[i]
            try {
              parsedLogs.push(iface.parseLog(log))
            } catch (e) {
              // ignore log
            }
          }
          assert.equal(parsedLogs.length, 1)
          assert.equal(parsedLogs[0].args.upkeepData, randomBytes)
        })

        it('uses actual execution price for payment', async () => {
          // Actual multiplier is 2, but we set gasPrice to be 1x gasWei
          const gasPrice = gasWei.mul(BigNumber.from('1'))
          mock.setCanPerform(true)

          const tx = await registry.connect(keeper1).transmit(
            [emptyBytes32, emptyBytes32, emptyBytes32],
            await encodeLatestBlockReport([
              {
                Id: upkeepId.toString(),
              },
            ]),
            [],
            [],
            emptyBytes32,
            { gasPrice },
          )
          const receipt = await tx.wait()
          let upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
          // exactly 1 Upkeep Performed should be emitted
          assert.equal(upkeepPerformedLogs.length, 1)
          let upkeepPerformedLog = upkeepPerformedLogs[0]

          let gasUsed = upkeepPerformedLog.args.gasUsed
          let gasOverhead = upkeepPerformedLog.args.gasOverhead
          let totalPayment = upkeepPerformedLog.args.totalPayment

          assert.equal(
            linkForGas(
              gasUsed,
              gasOverhead,
              BigNumber.from('1'), // Not the config multiplier, but the actual gas used
              paymentPremiumPPB,
              flatFeeMicroLink,
            ).toString(),
            totalPayment.toString(),
          )
        })

        it('only pays at a rate up to the gas ceiling', async () => {
          // Actual multiplier is 2, but we set gasPrice to be 100x
          const gasPrice = gasWei.mul(BigNumber.from('100'))
          mock.setCanPerform(true)

          const tx = await registry.connect(keeper1).transmit(
            [emptyBytes32, emptyBytes32, emptyBytes32],
            await encodeLatestBlockReport([
              {
                Id: upkeepId.toString(),
              },
            ]),
            [],
            [],
            emptyBytes32,
            { gasPrice },
          )
          const receipt = await tx.wait()
          let upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
          // exactly 1 Upkeep Performed should be emitted
          assert.equal(upkeepPerformedLogs.length, 1)
          let upkeepPerformedLog = upkeepPerformedLogs[0]

          let gasUsed = upkeepPerformedLog.args.gasUsed
          let gasOverhead = upkeepPerformedLog.args.gasOverhead
          let totalPayment = upkeepPerformedLog.args.totalPayment

          assert.equal(
            linkForGas(
              gasUsed,
              gasOverhead,
              gasCeilingMultiplier, // Should be same with exisitng multiplier
              paymentPremiumPPB,
              flatFeeMicroLink,
            ).toString(),
            totalPayment.toString(),
          )
        })

        it('uses the fallback gas price if the feed price is stale', async () => {
          const roundId = 99
          const answer = 100
          const updatedAt = 946684800 // New Years 2000 🥳
          const startedAt = 946684799
          await gasPriceFeed
            .connect(owner)
            .updateRoundData(roundId, answer, updatedAt, startedAt)

          const tx = await registry.connect(keeper1).transmit(
            [emptyBytes32, emptyBytes32, emptyBytes32],
            await encodeLatestBlockReport([
              {
                Id: upkeepId.toString(),
              },
            ]),
            [],
            [],
            emptyBytes32,
            { gasPrice: gasWei.mul('5') }, // High gas price so that it gets capped
          )
          const receipt = await tx.wait()
          let upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
          // exactly 1 Upkeep Performed should be emitted
          assert.equal(upkeepPerformedLogs.length, 1)
          let upkeepPerformedLog = upkeepPerformedLogs[0]

          let gasUsed = upkeepPerformedLog.args.gasUsed
          let gasOverhead = upkeepPerformedLog.args.gasOverhead
          let totalPayment = upkeepPerformedLog.args.totalPayment

          assert.equal(
            linkForGas(
              gasUsed,
              gasOverhead,
              gasCeilingMultiplier.mul('2'), // fallbackGasPrice is 2x gas price
              paymentPremiumPPB,
              flatFeeMicroLink,
            ).toString(),
            totalPayment.toString(),
          )
        })

        it('uses the fallback gas price if the feed price is non-sensical (negative)', async () => {
          const roundId = 99
          const updatedAt = Math.floor(Date.now() / 1000)
          const startedAt = 946684799
          await gasPriceFeed
            .connect(owner)
            .updateRoundData(roundId, -100, updatedAt, startedAt)

          // Negative feed value
          let tx = await registry.connect(keeper1).transmit(
            [emptyBytes32, emptyBytes32, emptyBytes32],
            await encodeLatestBlockReport([
              {
                Id: upkeepId.toString(),
              },
            ]),
            [],
            [],
            emptyBytes32,
            { gasPrice: gasWei.mul('5') }, // High gas price so that it gets capped
          )
          let receipt = await tx.wait()
          let upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
          // exactly 1 Upkeep Performed should be emitted
          assert.equal(upkeepPerformedLogs.length, 1)
          let upkeepPerformedLog = upkeepPerformedLogs[0]

          let gasUsed = upkeepPerformedLog.args.gasUsed
          let gasOverhead = upkeepPerformedLog.args.gasOverhead
          let totalPayment = upkeepPerformedLog.args.totalPayment

          assert.equal(
            linkForGas(
              gasUsed,
              gasOverhead,
              gasCeilingMultiplier.mul('2'), // fallbackGasPrice is 2x gas price
              paymentPremiumPPB,
              flatFeeMicroLink,
            ).toString(),
            totalPayment.toString(),
          )
        })

        it('uses the fallback gas price if the feed price is non-sensical (zero)', async () => {
          const roundId = 99
          const updatedAt = Math.floor(Date.now() / 1000)
          const startedAt = 946684799
          await gasPriceFeed
            .connect(owner)
            .updateRoundData(roundId, 0, updatedAt, startedAt)

          // Negative feed value
          let tx = await registry.connect(keeper1).transmit(
            [emptyBytes32, emptyBytes32, emptyBytes32],
            await encodeLatestBlockReport([
              {
                Id: upkeepId.toString(),
              },
            ]),
            [],
            [],
            emptyBytes32,
            { gasPrice: gasWei.mul('5') }, // High gas price so that it gets capped
          )
          let receipt = await tx.wait()
          let upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
          // exactly 1 Upkeep Performed should be emitted
          assert.equal(upkeepPerformedLogs.length, 1)
          let upkeepPerformedLog = upkeepPerformedLogs[0]

          let gasUsed = upkeepPerformedLog.args.gasUsed
          let gasOverhead = upkeepPerformedLog.args.gasOverhead
          let totalPayment = upkeepPerformedLog.args.totalPayment

          assert.equal(
            linkForGas(
              gasUsed,
              gasOverhead,
              gasCeilingMultiplier.mul('2'), // fallbackGasPrice is 2x gas price
              paymentPremiumPPB,
              flatFeeMicroLink,
            ).toString(),
            totalPayment.toString(),
          )
        })

        it('uses the fallback link price if the link price is stale', async () => {
          const roundId = 99
          const answer = 100
          const updatedAt = 946684800 // New Years 2000 🥳
          const startedAt = 946684799
          await linkEthFeed
            .connect(owner)
            .updateRoundData(roundId, answer, updatedAt, startedAt)

          const tx = await registry.connect(keeper1).transmit(
            [emptyBytes32, emptyBytes32, emptyBytes32],
            await encodeLatestBlockReport([
              {
                Id: upkeepId.toString(),
              },
            ]),
            [],
            [],
            emptyBytes32,
            { gasPrice: gasWei.mul('5') }, // High gas price so that it gets capped
          )
          const receipt = await tx.wait()
          let upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
          // exactly 1 Upkeep Performed should be emitted
          assert.equal(upkeepPerformedLogs.length, 1)
          let upkeepPerformedLog = upkeepPerformedLogs[0]

          let gasUsed = upkeepPerformedLog.args.gasUsed
          let gasOverhead = upkeepPerformedLog.args.gasOverhead
          let totalPayment = upkeepPerformedLog.args.totalPayment

          assert.equal(
            linkForGas(
              gasUsed,
              gasOverhead,
              gasCeilingMultiplier.mul('2'), // fallbackLinkPrice is 1/2 link price, so multiply by 2
              paymentPremiumPPB,
              flatFeeMicroLink,
            ).toString(),
            totalPayment.toString(),
          )
        })

        it('uses the fallback link price if the link price is non-sensical (negative)', async () => {
          const roundId = 99
          const updatedAt = Math.floor(Date.now() / 1000)
          const startedAt = 946684799
          await linkEthFeed
            .connect(owner)
            .updateRoundData(roundId, -100, updatedAt, startedAt)

          const tx = await registry.connect(keeper1).transmit(
            [emptyBytes32, emptyBytes32, emptyBytes32],
            await encodeLatestBlockReport([
              {
                Id: upkeepId.toString(),
              },
            ]),
            [],
            [],
            emptyBytes32,
            { gasPrice: gasWei.mul('5') }, // High gas price so that it gets capped
          )
          const receipt = await tx.wait()
          let upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
          // exactly 1 Upkeep Performed should be emitted
          assert.equal(upkeepPerformedLogs.length, 1)
          let upkeepPerformedLog = upkeepPerformedLogs[0]

          let gasUsed = upkeepPerformedLog.args.gasUsed
          let gasOverhead = upkeepPerformedLog.args.gasOverhead
          let totalPayment = upkeepPerformedLog.args.totalPayment

          assert.equal(
            linkForGas(
              gasUsed,
              gasOverhead,
              gasCeilingMultiplier.mul('2'), // fallbackLinkPrice is 1/2 link price, so multiply by 2
              paymentPremiumPPB,
              flatFeeMicroLink,
            ).toString(),
            totalPayment.toString(),
          )
        })

        it('uses the fallback link price if the link price is non-sensical (zero)', async () => {
          const roundId = 99
          const updatedAt = Math.floor(Date.now() / 1000)
          const startedAt = 946684799
          await linkEthFeed
            .connect(owner)
            .updateRoundData(roundId, 0, updatedAt, startedAt)

          const tx = await registry.connect(keeper1).transmit(
            [emptyBytes32, emptyBytes32, emptyBytes32],
            await encodeLatestBlockReport([
              {
                Id: upkeepId.toString(),
              },
            ]),
            [],
            [],
            emptyBytes32,
            { gasPrice: gasWei.mul('5') }, // High gas price so that it gets capped
          )
          const receipt = await tx.wait()
          let upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
          // exactly 1 Upkeep Performed should be emitted
          assert.equal(upkeepPerformedLogs.length, 1)
          let upkeepPerformedLog = upkeepPerformedLogs[0]

          let gasUsed = upkeepPerformedLog.args.gasUsed
          let gasOverhead = upkeepPerformedLog.args.gasOverhead
          let totalPayment = upkeepPerformedLog.args.totalPayment

          assert.equal(
            linkForGas(
              gasUsed,
              gasOverhead,
              gasCeilingMultiplier.mul('2'), // fallbackLinkPrice is 1/2 link price, so multiply by 2
              paymentPremiumPPB,
              flatFeeMicroLink,
            ).toString(),
            totalPayment.toString(),
          )
        })

        it('correctly accounts for l1 payment')

        it('can self fund', async () => {
          const autoFunderUpkeep = await upkeepAutoFunderFactory
            .connect(owner)
            .deploy(linkToken.address, registry.address)
          let tx = await registry
            .connect(owner)
            .registerUpkeep(
              autoFunderUpkeep.address,
              executeGas,
              autoFunderUpkeep.address,
              true,
              randomBytes,
            )
          upkeepId = await getUpkeepID(tx)

          await autoFunderUpkeep.setUpkeepId(upkeepId)
          // Give enough funds for upkeep as well as to the upkeep contract
          await linkToken
            .connect(owner)
            .transfer(autoFunderUpkeep.address, toWei('1000'))
          let maxPayment = await registry.getMaxPaymentForGas(executeGas)

          // First set auto funding amount to 0 and verify that balance is deducted upon performUpkeep
          let initialBalance = toWei('100')
          await registry.connect(owner).addFunds(upkeepId, initialBalance)
          await autoFunderUpkeep.setAutoFundLink(0)
          await autoFunderUpkeep.setIsEligible(true)
          await registry.connect(keeper1).transmit(
            [emptyBytes32, emptyBytes32, emptyBytes32],
            await encodeLatestBlockReport([
              {
                Id: upkeepId.toString(),
              },
            ]),
            [],
            [],
            emptyBytes32,
          )

          let postUpkeepBalance = (await registry.getUpkeep(upkeepId)).balance
          assert.isTrue(postUpkeepBalance.lt(initialBalance)) // Balance should be deducted
          assert.isTrue(postUpkeepBalance.gte(initialBalance.sub(maxPayment))) // Balance should not be deducted more than maxPayment

          // Now set auto funding amount to 100 wei and verify that the balance increases
          initialBalance = postUpkeepBalance
          let autoTopupAmount = toWei('100')
          await autoFunderUpkeep.setAutoFundLink(autoTopupAmount)
          await autoFunderUpkeep.setIsEligible(true)
          await registry.connect(keeper1).transmit(
            [emptyBytes32, emptyBytes32, emptyBytes32],
            await encodeLatestBlockReport([
              {
                Id: upkeepId.toString(),
              },
            ]),
            [],
            [],
            emptyBytes32,
          )

          postUpkeepBalance = (await registry.getUpkeep(upkeepId)).balance
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
          let tx = await registry
            .connect(owner)
            .registerUpkeep(
              autoFunderUpkeep.address,
              executeGas,
              autoFunderUpkeep.address,
              true,
              randomBytes,
            )
          upkeepId = await getUpkeepID(tx)

          await autoFunderUpkeep.setUpkeepId(upkeepId)
          await registry.connect(owner).addFunds(upkeepId, toWei('100'))

          await autoFunderUpkeep.setIsEligible(true)
          await autoFunderUpkeep.setShouldCancel(true)

          let registration = await registry.getUpkeep(upkeepId)
          const oldExpiration = registration.maxValidBlocknumber

          // Do the thing
          await registry.connect(keeper1).transmit(
            [emptyBytes32, emptyBytes32, emptyBytes32],
            await encodeLatestBlockReport([
              {
                Id: upkeepId.toString(),
              },
            ]),
            [],
            [],
            emptyBytes32,
          )

          // Verify upkeep gets cancelled
          registration = await registry.getUpkeep(upkeepId)
          const newExpiration = registration.maxValidBlocknumber
          assert.isTrue(newExpiration.lt(oldExpiration))
        })

        it('performs upkeep, deducts payment, updates lastPerformBlockNumber and emits event', async () => {
          mock.setCanPerform(true)
          let checkBlock = await ethers.provider.getBlock('latest')

          const keeperBefore = await registry.getTransmitterInfo(
            await keeper1.getAddress(),
          )
          const registrationBefore = await registry.getUpkeep(upkeepId)
          const keeperLinkBefore = await linkToken.balanceOf(
            await keeper1.getAddress(),
          )
          const registryLinkBefore = await linkToken.balanceOf(registry.address)
          assert.equal(
            registrationBefore.lastPerformBlockNumber.toString(),
            '0',
          )

          const tx = await registry.connect(keeper1).transmit(
            [emptyBytes32, emptyBytes32, emptyBytes32],
            encodeReport([
              {
                Id: upkeepId.toString(),
                checkBlockNum: checkBlock.number,
                checkBlockHash: checkBlock.parentHash,
                performData: '0x',
              },
            ]),
            [],
            [],
            emptyBytes32,
          )
          const receipt = await tx.wait()
          let upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
          // exactly 1 Upkeep Performed should be emitted
          assert.equal(upkeepPerformedLogs.length, 1)
          let upkeepPerformedLog = upkeepPerformedLogs[0]

          let id = upkeepPerformedLog.args.id
          let success = upkeepPerformedLog.args.success
          let checkBlockNumber = upkeepPerformedLog.args.checkBlockNumber
          let gasUsed = upkeepPerformedLog.args.gasUsed
          let gasOverhead = upkeepPerformedLog.args.gasOverhead
          let totalPayment = upkeepPerformedLog.args.totalPayment

          assert.equal(id.toString(), upkeepId.toString())
          assert.equal(success, true)
          assert.equal(
            checkBlockNumber.toString(),
            checkBlock.number.toString(),
          )
          assert.isTrue(gasUsed.gt(BigNumber.from('0')))
          assert.isTrue(gasOverhead.gt(BigNumber.from('0')))
          assert.isTrue(totalPayment.gt(BigNumber.from('0')))

          const keeperAfter = await registry.getTransmitterInfo(
            await keeper1.getAddress(),
          )
          const registrationAfter = await registry.getUpkeep(upkeepId)
          const keeperLinkAfter = await linkToken.balanceOf(
            await keeper1.getAddress(),
          )
          const registryLinkAfter = await linkToken.balanceOf(registry.address)

          assert.equal(
            keeperAfter.balance.sub(totalPayment).toString(),
            keeperBefore.balance.toString(),
          )
          assert.equal(
            registrationBefore.balance.sub(totalPayment).toString(),
            registrationAfter.balance.toString(),
          )
          assert.isTrue(keeperLinkAfter.eq(keeperLinkBefore))
          assert.isTrue(registryLinkBefore.eq(registryLinkAfter))

          // Amount spent should be updated correctly
          assert.equal(
            registrationAfter.amountSpent.sub(totalPayment).toString(),
            registrationBefore.amountSpent.toString(),
          )
          assert.isTrue(
            registrationAfter.amountSpent
              .sub(registrationBefore.amountSpent)
              .eq(registrationBefore.balance.sub(registrationAfter.balance)),
          )
          // Last perform block number should be updated
          assert.equal(
            registrationAfter.lastPerformBlockNumber.toString(),
            tx.blockNumber?.toString(),
          )

          // Transmitted should not be emitted for skip sig validation
          expect(tx).to.not.emit(registry, 'Transmitted')
        })

        it('calculates gas overhead appropriately within a margin for different scenarios', async () => {
          // Perform the upkeep once to remove non-zero storage slots and have predictable gas measurement
          let tx = await registry.connect(keeper1).transmit(
            [emptyBytes32, emptyBytes32, emptyBytes32],
            await encodeLatestBlockReport([
              {
                Id: upkeepId.toString(),
              },
            ]),
            [],
            [],
            emptyBytes32,
          )
          await tx.wait()

          // Different test scenarios
          let longBytes = '0x'
          for (let i = 0; i < maxPerformDataSize.toNumber(); i++) {
            longBytes += '11'
          }
          let upkeepSuccessArray = [true, false]
          let performGasArray = [5000, 100000, executeGas]
          let performDataArray = ['0x', randomBytes, longBytes]

          for (let i in upkeepSuccessArray) {
            for (let j in performGasArray) {
              for (let k in performDataArray) {
                for (let l in fArray) {
                  const upkeepSuccess = upkeepSuccessArray[i]
                  const performGas = performGasArray[j]
                  const performData = performDataArray[k]
                  let newF = fArray[l]

                  mock.setCanPerform(upkeepSuccess)
                  mock.setPerformGasToBurn(performGas)
                  await registry
                    .connect(owner)
                    .setConfig(
                      signerAddresses,
                      keeperAddresses,
                      newF,
                      encodeConfig(config),
                      offchainVersion,
                      offchainBytes,
                    )
                  let latestBlock = await ethers.provider.getBlock('latest')

                  tx = await registry.connect(keeper1).transmit(
                    [emptyBytes32, emptyBytes32, emptyBytes32],
                    encodeReport([
                      {
                        Id: upkeepId.toString(),
                        checkBlockNum: latestBlock.number + 1,
                        checkBlockHash: latestBlock.hash,
                        performData: performData,
                      },
                    ]),
                    [],
                    [],
                    emptyBytes32,
                  )
                  const receipt = await tx.wait()
                  let upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
                  // exactly 1 Upkeep Performed should be emitted
                  assert.equal(upkeepPerformedLogs.length, 1)
                  let upkeepPerformedLog = upkeepPerformedLogs[0]

                  let gasUsed = upkeepPerformedLog.args.gasUsed
                  let gasOverhead = upkeepPerformedLog.args.gasOverhead

                  assert.isTrue(gasUsed.gt(BigNumber.from('0')))
                  assert.isTrue(gasOverhead.gt(BigNumber.from('0')))

                  if (i == '0' && j == '0' && k == '0' && l == '0') {
                    console.log(
                      'Gas Benchmarking - non sig verification: calculated overhead: ',
                      gasOverhead.toString(),
                      ' margin over gasUsed: ',
                      gasUsed.add(gasOverhead).sub(receipt.gasUsed).toString(),
                    )
                  }

                  // Overhead should not get capped
                  assert.isTrue(
                    gasOverhead.lt(
                      registryGasOverhead.add(
                        BigNumber.from(16 * performData.length),
                      ),
                    ),
                  )
                  // total gas charged should be greater than tx gas but within gasCalculationMargin
                  assert.isTrue(gasUsed.add(gasOverhead).gt(receipt.gasUsed))
                  assert.isTrue(
                    gasUsed
                      .add(gasOverhead)
                      .lt(
                        receipt.gasUsed.add(
                          BigNumber.from(gasCalculationMargin),
                        ),
                      ),
                  )
                }
              }
            }
          }
        })
      })

      describe('When signatures are validated', () => {
        it('reverts when configDigest mismatches', async () => {
          const report = await encodeLatestBlockReport([
            {
              Id: sigVerificationUpkeepId.toString(),
            },
          ])
          const reportContext = [emptyBytes32, epochAndRound5_1, emptyBytes32] // wrong config digest
          const sigs = signReport(
            reportContext,
            report,
            signers.slice(0, f + 1),
          )
          await evmRevert(
            registry
              .connect(keeper1)
              .transmit(
                [reportContext[0], reportContext[1], reportContext[2]],
                report,
                sigs.rs,
                sigs.ss,
                sigs.vs,
              ),
            'ConfigDigestMismatch()',
          )
        })

        it('reverts with incorrect number of signatures', async () => {
          const configDigest = (await registry.getState()).state
            .latestConfigDigest
          const report = await encodeLatestBlockReport([
            {
              Id: sigVerificationUpkeepId.toString(),
            },
          ])
          const reportContext = [configDigest, epochAndRound5_1, emptyBytes32] // wrong config digest
          const sigs = signReport(
            reportContext,
            report,
            signers.slice(0, f + 2),
          )
          await evmRevert(
            registry
              .connect(keeper1)
              .transmit(
                [reportContext[0], reportContext[1], reportContext[2]],
                report,
                sigs.rs,
                sigs.ss,
                sigs.vs,
              ),
            'IncorrectNumberOfSignatures()',
          )
        })

        it('reverts with invalid signature for inactive signers', async () => {
          const configDigest = (await registry.getState()).state
            .latestConfigDigest
          const report = await encodeLatestBlockReport([
            {
              Id: sigVerificationUpkeepId.toString(),
            },
          ])
          const reportContext = [configDigest, epochAndRound5_1, emptyBytes32] // wrong config digest
          const sigs = signReport(reportContext, report, [
            new ethers.Wallet(randomAddress()),
            new ethers.Wallet(randomAddress()),
          ])
          await evmRevert(
            registry
              .connect(keeper1)
              .transmit(
                [reportContext[0], reportContext[1], reportContext[2]],
                report,
                sigs.rs,
                sigs.ss,
                sigs.vs,
              ),
            'OnlyActiveSigners()',
          )
        })

        it('reverts with invalid signature for duplicated signers', async () => {
          const configDigest = (await registry.getState()).state
            .latestConfigDigest
          const report = await encodeLatestBlockReport([
            {
              Id: sigVerificationUpkeepId.toString(),
            },
          ])
          const reportContext = [configDigest, epochAndRound5_1, emptyBytes32] // wrong config digest
          const sigs = signReport(reportContext, report, [signer1, signer1])
          await evmRevert(
            registry
              .connect(keeper1)
              .transmit(
                [reportContext[0], reportContext[1], reportContext[2]],
                report,
                sigs.rs,
                sigs.ss,
                sigs.vs,
              ),
            'DuplicateSigners()',
          )
        })

        it('has a large enough gas overhead to cover upkeep that use all its gas', async () => {
          await registry.connect(owner).setConfig(
            signerAddresses,
            keeperAddresses,
            10, // maximise f to maximise overhead
            encodeConfig(config),
            offchainVersion,
            offchainBytes,
          )
          const tx = await registry.connect(owner).registerUpkeep(
            mock.address,
            maxPerformGas, // max allowed gas
            await admin.getAddress(),
            true,
            randomBytes,
          )
          sigVerificationUpkeepId = await getUpkeepID(tx)
          await registry
            .connect(admin)
            .addFunds(sigVerificationUpkeepId, toWei('100'))

          let performData = '0x'
          for (let i = 0; i < maxPerformDataSize.toNumber(); i++) {
            performData += '11'
          } // max allowed performData

          mock.setCanPerform(true)
          mock.setPerformGasToBurn(maxPerformGas)

          let latestBlock = await ethers.provider.getBlock('latest')
          let configDigest = (await registry.getState()).state
            .latestConfigDigest
          let report = encodeReport([
            {
              Id: sigVerificationUpkeepId.toString(),
              checkBlockNum: latestBlock.number + 1,
              checkBlockHash: latestBlock.hash,
              performData: performData,
            },
          ])
          let reportContext = [configDigest, epochAndRound5_1, emptyBytes32]
          let sigs = signReport(reportContext, report, signers.slice(0, 11))

          await registry
            .connect(keeper1)
            .transmit(
              [reportContext[0], reportContext[1], reportContext[2]],
              report,
              sigs.rs,
              sigs.ss,
              sigs.vs,
              { gasLimit: maxPerformGas.add(transmitGasOverhead) },
            ) // Should not revert
        })

        it('performs upkeep, deducts payment, updates lastPerformBlockNumber and emits events', async () => {
          for (let i in fArray) {
            let newF = fArray[i]
            await registry
              .connect(owner)
              .setConfig(
                signerAddresses,
                keeperAddresses,
                newF,
                encodeConfig(config),
                offchainVersion,
                offchainBytes,
              )
            mock.setCanPerform(true)
            let checkBlock = await ethers.provider.getBlock('latest')

            const keeperBefore = await registry.getTransmitterInfo(
              await keeper1.getAddress(),
            )
            const registrationBefore = await registry.getUpkeep(
              sigVerificationUpkeepId,
            )
            const keeperLinkBefore = await linkToken.balanceOf(
              await keeper1.getAddress(),
            )
            const registryLinkBefore = await linkToken.balanceOf(
              registry.address,
            )

            // Do the thing
            const configDigest = (await registry.getState()).state
              .latestConfigDigest
            const report = encodeReport([
              {
                Id: sigVerificationUpkeepId.toString(),
                checkBlockNum: checkBlock.number,
                checkBlockHash: checkBlock.parentHash,
                performData: '0x',
              },
            ])
            const reportContext = [configDigest, epochAndRound5_1, emptyBytes32]
            const sigs = signReport(
              reportContext,
              report,
              signers.slice(0, newF + 1),
            )
            const tx = await registry
              .connect(keeper1)
              .transmit(
                [reportContext[0], reportContext[1], reportContext[2]],
                report,
                sigs.rs,
                sigs.ss,
                sigs.vs,
              )
            const receipt = await tx.wait()

            let upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
            // exactly 1 Upkeep Performed should be emitted
            assert.equal(upkeepPerformedLogs.length, 1)
            let upkeepPerformedLog = upkeepPerformedLogs[0]

            let id = upkeepPerformedLog.args.id
            let success = upkeepPerformedLog.args.success
            let checkBlockNumber = upkeepPerformedLog.args.checkBlockNumber
            let gasUsed = upkeepPerformedLog.args.gasUsed
            let gasOverhead = upkeepPerformedLog.args.gasOverhead
            let totalPayment = upkeepPerformedLog.args.totalPayment

            assert.equal(id.toString(), sigVerificationUpkeepId.toString())
            assert.equal(success, true)
            assert.equal(
              checkBlockNumber.toString(),
              checkBlock.number.toString(),
            )
            assert.isTrue(gasUsed.gt(BigNumber.from('0')))
            assert.isTrue(gasOverhead.gt(BigNumber.from('0')))
            assert.isTrue(totalPayment.gt(BigNumber.from('0')))

            const keeperAfter = await registry.getTransmitterInfo(
              await keeper1.getAddress(),
            )
            const registrationAfter = await registry.getUpkeep(
              sigVerificationUpkeepId,
            )
            const keeperLinkAfter = await linkToken.balanceOf(
              await keeper1.getAddress(),
            )
            const registryLinkAfter = await linkToken.balanceOf(
              registry.address,
            )

            assert.equal(
              keeperAfter.balance.sub(totalPayment).toString(),
              keeperBefore.balance.toString(),
            )
            assert.equal(
              registrationBefore.balance.sub(totalPayment).toString(),
              registrationAfter.balance.toString(),
            )
            assert.isTrue(keeperLinkAfter.eq(keeperLinkBefore))
            assert.isTrue(registryLinkBefore.eq(registryLinkAfter))

            // Amount spent should be updated correctly
            assert.equal(
              registrationAfter.amountSpent.sub(totalPayment).toString(),
              registrationBefore.amountSpent.toString(),
            )
            assert.isTrue(
              registrationAfter.amountSpent
                .sub(registrationBefore.amountSpent)
                .eq(registrationBefore.balance.sub(registrationAfter.balance)),
            )
            // Last perform block number should be updated
            assert.equal(
              registrationAfter.lastPerformBlockNumber.toString(),
              tx.blockNumber?.toString(),
            )

            // Epoch should be 5 as we used epochAndRound5_1
            await expect(tx)
              .to.emit(registry, 'Transmitted')
              .withArgs(configDigest, 5)
          }
        })

        it('calculates gas overhead appropriately within a margin for different scenarios', async () => {
          // Perform the upkeep once to remove non-zero storage slots and have predictable gas measurement
          let configDigest = (await registry.getState()).state
            .latestConfigDigest
          let report = await encodeLatestBlockReport([
            {
              Id: sigVerificationUpkeepId.toString(),
            },
          ])
          let reportContext = [configDigest, epochAndRound5_1, emptyBytes32]
          let sigs = signReport(reportContext, report, signers.slice(0, f + 1))
          let tx = await registry
            .connect(keeper1)
            .transmit(
              [reportContext[0], reportContext[1], reportContext[2]],
              report,
              sigs.rs,
              sigs.ss,
              sigs.vs,
            )
          await tx.wait()

          // Different test scenarios
          let longBytes = '0x'
          for (let i = 0; i < maxPerformDataSize.toNumber(); i++) {
            longBytes += '11'
          }
          let upkeepSuccessArray = [true, false]
          let performGasArray = [5000, 100000, executeGas]
          let performDataArray = ['0x', randomBytes, longBytes]

          for (let i in upkeepSuccessArray) {
            for (let j in performGasArray) {
              for (let k in performDataArray) {
                for (let l in fArray) {
                  const upkeepSuccess = upkeepSuccessArray[i]
                  const performGas = performGasArray[j]
                  const performData = performDataArray[k]
                  let newF = fArray[l]

                  mock.setCanPerform(upkeepSuccess)
                  mock.setPerformGasToBurn(performGas)
                  await registry
                    .connect(owner)
                    .setConfig(
                      signerAddresses,
                      keeperAddresses,
                      newF,
                      encodeConfig(config),
                      offchainVersion,
                      offchainBytes,
                    )

                  let latestBlock = await ethers.provider.getBlock('latest')
                  configDigest = (await registry.getState()).state
                    .latestConfigDigest
                  report = encodeReport([
                    {
                      Id: sigVerificationUpkeepId.toString(),
                      checkBlockNum: latestBlock.number + 1,
                      checkBlockHash: latestBlock.hash,
                      performData: performData,
                    },
                  ])
                  reportContext = [configDigest, epochAndRound5_1, emptyBytes32]
                  sigs = signReport(
                    reportContext,
                    report,
                    signers.slice(0, newF + 1),
                  )
                  tx = await registry
                    .connect(keeper1)
                    .transmit(
                      [reportContext[0], reportContext[1], reportContext[2]],
                      report,
                      sigs.rs,
                      sigs.ss,
                      sigs.vs,
                    )
                  const receipt = await tx.wait()
                  let upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
                  // exactly 1 Upkeep Performed should be emitted
                  assert.equal(upkeepPerformedLogs.length, 1)
                  let upkeepPerformedLog = upkeepPerformedLogs[0]

                  let gasUsed = upkeepPerformedLog.args.gasUsed
                  let gasOverhead = upkeepPerformedLog.args.gasOverhead

                  assert.isTrue(gasUsed.gt(BigNumber.from('0')))
                  assert.isTrue(gasOverhead.gt(BigNumber.from('0')))

                  if (i == '0' && j == '0' && k == '0') {
                    console.log(
                      'Gas Benchmarking - sig verification ( f =',
                      newF,
                      '): calculated overhead: ',
                      gasOverhead.toString(),
                      ' margin over gasUsed: ',
                      gasUsed.add(gasOverhead).sub(receipt.gasUsed).toString(),
                    )
                  }

                  // Overhead should not get capped
                  assert.isTrue(
                    gasOverhead.lt(
                      registryGasOverhead
                        .add(verifySigOverhead.mul(BigNumber.from(newF + 1)))
                        .add(BigNumber.from(16 * performData.length)),
                    ),
                  )
                  // total gas charged should be greater than tx gas but within gasCalculationMargin
                  assert.isTrue(gasUsed.add(gasOverhead).gt(receipt.gasUsed))
                  assert.isTrue(
                    gasUsed
                      .add(gasOverhead)
                      .lt(
                        receipt.gasUsed.add(
                          BigNumber.from(gasCalculationMargin),
                        ),
                      ),
                  )
                }
              }
            }
          }
        })
      })
    })

    describe('When upkeeps are batched (non sig verification)', () => {
      const numPassingUpkeepsArray = [1, 3, 5]
      const numFailingUpkeepsArray = [0, 1, 5]

      numPassingUpkeepsArray.forEach(function (numPassingUpkeeps) {
        numFailingUpkeepsArray.forEach(function (numFailingUpkeeps) {
          describe(
            numPassingUpkeeps.toString() +
              ' passing upkeeps, ' +
              numFailingUpkeeps.toString() +
              ' failing upkeeps',
            () => {
              let passingUpkeepIds: string[]
              let failingUpkeepIds: string[]

              beforeEach(async () => {
                passingUpkeepIds = []
                failingUpkeepIds = []
                for (let i = 0; i < numPassingUpkeeps; i++) {
                  let tx = await registry
                    .connect(owner)
                    .registerUpkeep(
                      mock.address,
                      executeGas,
                      await admin.getAddress(),
                      true,
                      randomBytes,
                    )
                  upkeepId = await getUpkeepID(tx)
                  passingUpkeepIds.push(upkeepId.toString())

                  // Add funds to passing upkeeps
                  await registry.connect(admin).addFunds(upkeepId, toWei('100'))
                }
                for (let i = 0; i < numFailingUpkeeps; i++) {
                  let tx = await registry
                    .connect(owner)
                    .registerUpkeep(
                      mock.address,
                      executeGas,
                      await admin.getAddress(),
                      true,
                      randomBytes,
                    )
                  upkeepId = await getUpkeepID(tx)
                  failingUpkeepIds.push(upkeepId.toString())
                }
              })

              it('performs successful upkeeps and does not change failing upkeeps', async () => {
                const keeperBefore = await registry.getTransmitterInfo(
                  await keeper1.getAddress(),
                )
                const keeperLinkBefore = await linkToken.balanceOf(
                  await keeper1.getAddress(),
                )
                const registryLinkBefore = await linkToken.balanceOf(
                  registry.address,
                )
                const registrationPassingBefore = await Promise.all(
                  passingUpkeepIds.map(async (id) => {
                    let reg = await registry.getUpkeep(BigNumber.from(id))
                    assert.equal(reg.lastPerformBlockNumber.toString(), '0')
                    return reg
                  }),
                )
                const registrationFailingBefore = await await Promise.all(
                  failingUpkeepIds.map(async (id) => {
                    let reg = await registry.getUpkeep(BigNumber.from(id))
                    assert.equal(reg.lastPerformBlockNumber.toString(), '0')
                    return reg
                  }),
                )

                const tx = await registry.connect(keeper1).transmit(
                  [emptyBytes32, emptyBytes32, emptyBytes32],
                  await encodeLatestBlockReport(
                    passingUpkeepIds.concat(failingUpkeepIds).map((id) => {
                      return {
                        Id: id,
                      }
                    }),
                  ),
                  [],
                  [],
                  emptyBytes32,
                )
                const receipt = await tx.wait()
                let upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
                // exactly numPassingUpkeeps Upkeep Performed should be emitted
                assert.equal(upkeepPerformedLogs.length, numPassingUpkeeps)

                const keeperAfter = await registry.getTransmitterInfo(
                  await keeper1.getAddress(),
                )
                const keeperLinkAfter = await linkToken.balanceOf(
                  await keeper1.getAddress(),
                )
                const registryLinkAfter = await linkToken.balanceOf(
                  registry.address,
                )
                const registrationPassingAfter = await Promise.all(
                  passingUpkeepIds.map(async (id) => {
                    return await registry.getUpkeep(BigNumber.from(id))
                  }),
                )
                const registrationFailingAfter = await await Promise.all(
                  failingUpkeepIds.map(async (id) => {
                    return await registry.getUpkeep(BigNumber.from(id))
                  }),
                )

                let netPayment = BigNumber.from('0')
                for (let i = 0; i < numPassingUpkeeps; i++) {
                  let id = upkeepPerformedLogs[i].args.id
                  let gasUsed = upkeepPerformedLogs[i].args.gasUsed
                  let gasOverhead = upkeepPerformedLogs[i].args.gasOverhead
                  let totalPayment = upkeepPerformedLogs[i].args.totalPayment

                  assert.equal(id.toString(), passingUpkeepIds[i])
                  assert.isTrue(gasUsed.gt(BigNumber.from('0')))
                  assert.isTrue(gasOverhead.gt(BigNumber.from('0')))
                  assert.isTrue(totalPayment.gt(BigNumber.from('0')))

                  // Balance should be deducted
                  assert.equal(
                    registrationPassingBefore[i].balance
                      .sub(totalPayment)
                      .toString(),
                    registrationPassingAfter[i].balance.toString(),
                  )

                  // Amount spent should be updated correctly
                  assert.equal(
                    registrationPassingAfter[i].amountSpent
                      .sub(totalPayment)
                      .toString(),
                    registrationPassingBefore[i].amountSpent.toString(),
                  )

                  // Last perform block number should be updated
                  assert.equal(
                    registrationPassingAfter[
                      i
                    ].lastPerformBlockNumber.toString(),
                    tx.blockNumber?.toString(),
                  )

                  netPayment = netPayment.add(totalPayment)
                }

                for (let i = 0; i < numFailingUpkeeps; i++) {
                  // Balance and amount spent should be same
                  assert.equal(
                    registrationFailingBefore[i].balance.toString(),
                    registrationFailingAfter[i].balance.toString(),
                  )
                  assert.equal(
                    registrationFailingBefore[i].amountSpent.toString(),
                    registrationFailingAfter[i].amountSpent.toString(),
                  )

                  // Last perform block number should not be updated
                  assert.equal(
                    registrationFailingAfter[
                      i
                    ].lastPerformBlockNumber.toString(),
                    '0',
                  )
                }

                // Keeper should be paid net payment for all passed upkeeps
                assert.equal(
                  keeperAfter.balance.sub(netPayment).toString(),
                  keeperBefore.balance.toString(),
                )

                assert.isTrue(keeperLinkAfter.eq(keeperLinkBefore))
                assert.isTrue(registryLinkBefore.eq(registryLinkAfter))
              })

              it('splits gas overhead appropriately among performed upkeeps')

              it('benchmarks gas usage and verifies that its within margin')

              it('caps gasOverhead when other upkeeps are failing')
            },
          )
        })
      })

      it('splits l1 payment among performed upkeeps')
    })

    describe('When upkeeps are batched (sig verification)', () => {
      it('performs multiple upkeeps in batch')
      it('benchmarks gas usage and verifies that its within margin')
    })
  })

  describe('#recoverFunds', () => {
    const sent = toWei('7')

    beforeEach(async () => {
      await linkToken.connect(admin).approve(registry.address, toWei('100'))
      await linkToken
        .connect(owner)
        .transfer(await keeper1.getAddress(), toWei('1000'))

      // add funds to upkeep 1 and perform and withdraw some payment
      const tx = await registry
        .connect(owner)
        .registerUpkeep(
          mock.address,
          executeGas,
          await admin.getAddress(),
          true,
          emptyBytes,
        )

      const id1 = await getUpkeepID(tx)
      await registry.connect(admin).addFunds(id1, toWei('5'))

      await registry
        .connect(keeper1)
        .transmit(
          [emptyBytes32, emptyBytes32, emptyBytes32],
          await encodeLatestBlockReport([{ Id: id1.toString() }]),
          [],
          [],
          emptyBytes32,
        )
      await registry
        .connect(keeper2)
        .transmit(
          [emptyBytes32, emptyBytes32, emptyBytes32],
          await encodeLatestBlockReport([{ Id: id1.toString() }]),
          [],
          [],
          emptyBytes32,
        )
      await registry
        .connect(keeper3)
        .transmit(
          [emptyBytes32, emptyBytes32, emptyBytes32],
          await encodeLatestBlockReport([{ Id: id1.toString() }]),
          [],
          [],
          emptyBytes32,
        )
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
          true,
          emptyBytes,
        )
      const id2 = await getUpkeepID(tx2)
      await registry.connect(admin).addFunds(id2, toWei('5'))

      await registry
        .connect(keeper1)
        .transmit(
          [emptyBytes32, emptyBytes32, emptyBytes32],
          await encodeLatestBlockReport([{ Id: id2.toString() }]),
          [],
          [],
          emptyBytes32,
        )
      await registry
        .connect(keeper2)
        .transmit(
          [emptyBytes32, emptyBytes32, emptyBytes32],
          await encodeLatestBlockReport([{ Id: id2.toString() }]),
          [],
          [],
          emptyBytes32,
        )
      await registry
        .connect(keeper3)
        .transmit(
          [emptyBytes32, emptyBytes32, emptyBytes32],
          await encodeLatestBlockReport([{ Id: id2.toString() }]),
          [],
          [],
          emptyBytes32,
        )

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

      // withdraw some funds
      await registry.connect(owner).cancelUpkeep(id1)
      await registry
        .connect(admin)
        .withdrawFunds(id1, await nonkeeper.getAddress())
    })

    it('reverts if not called by owner', async () => {
      await evmRevert(
        registry.connect(keeper1).recoverFunds(),
        'Only callable by owner',
      )
    })

    it('allows any funds that have been accidentally transfered to be moved', async () => {
      const balanceBefore = await linkToken.balanceOf(registry.address)
      const ownerBefore = await linkToken.balanceOf(await owner.getAddress())

      await registry.connect(owner).recoverFunds()

      const balanceAfter = await linkToken.balanceOf(registry.address)
      const ownerAfter = await linkToken.balanceOf(await owner.getAddress())

      assert.isTrue(balanceBefore.eq(balanceAfter.add(sent)))
      assert.isTrue(ownerAfter.eq(ownerBefore.add(sent)))
    })
  })

  describe('#getMinBalanceForUpkeep / #checkUpkeep / #transmit', () => {
    it('calculates the minimum balance appropriately for skip verification upkeep', async () => {
      await mock.setCanCheck(true)
      const oneWei = BigNumber.from(1)
      const minBalance = await registry.getMinBalanceForUpkeep(upkeepId)

      await registry.connect(admin).addFunds(upkeepId, oneWei)
      let checkUpkeepResult = await registry
        .connect(zeroAddress)
        .callStatic.checkUpkeep(upkeepId)

      assert.equal(checkUpkeepResult.upkeepNeeded, false)
      assert.equal(checkUpkeepResult.upkeepFailureReason, 6)

      await registry.connect(admin).addFunds(upkeepId, minBalance.sub(oneWei))
      checkUpkeepResult = await registry
        .connect(zeroAddress)
        .callStatic.checkUpkeep(upkeepId)
      assert.equal(checkUpkeepResult.upkeepNeeded, true)
    })

    it('calculates the minimum balance appropriately for non skip verification upkeep', async () => {
      await mock.setCanCheck(true)
      const tx = await registry
        .connect(owner)
        .registerUpkeep(
          mock.address,
          executeGas,
          await admin.getAddress(),
          false,
          randomBytes,
        )
      upkeepId = await getUpkeepID(tx)

      const oneWei = BigNumber.from(1)
      const minBalance = await registry.getMinBalanceForUpkeep(upkeepId)
      const tooLow = minBalance.sub(oneWei)

      await registry.connect(admin).addFunds(upkeepId, tooLow)
      let checkUpkeepResult = await registry
        .connect(zeroAddress)
        .callStatic.checkUpkeep(upkeepId)

      assert.equal(checkUpkeepResult.upkeepNeeded, false)
      assert.equal(checkUpkeepResult.upkeepFailureReason, 6)

      await registry.connect(admin).addFunds(upkeepId, oneWei)
      checkUpkeepResult = await registry
        .connect(zeroAddress)
        .callStatic.checkUpkeep(upkeepId)
      assert.equal(checkUpkeepResult.upkeepNeeded, true)
    })

    it('uses maxPerformData size in checkUpkeep but actual performDataSize in transmit', async () => {
      const tx1 = await registry
        .connect(owner)
        .registerUpkeep(
          mock.address,
          executeGas,
          await admin.getAddress(),
          false,
          randomBytes,
        )
      const upkeepID1 = await getUpkeepID(tx1)
      const tx2 = await registry
        .connect(owner)
        .registerUpkeep(
          mock.address,
          executeGas,
          await admin.getAddress(),
          false,
          randomBytes,
        )
      const upkeepID2 = await getUpkeepID(tx2)
      await mock.setCanCheck(true)
      await mock.setCanPerform(true)

      // upkeep 1 is underfunded, 2 is fully funded
      const minBalance1 = (
        await registry.getMinBalanceForUpkeep(upkeepID1)
      ).sub(1)
      const minBalance2 = await registry.getMinBalanceForUpkeep(upkeepID2)
      await registry.connect(owner).addFunds(upkeepID1, minBalance1)
      await registry.connect(owner).addFunds(upkeepID2, minBalance2)

      // upkeep 1 check should return false, 2 should return true
      let checkUpkeepResult = await registry
        .connect(zeroAddress)
        .callStatic.checkUpkeep(upkeepID1)
      assert.equal(checkUpkeepResult.upkeepNeeded, false)
      assert.equal(checkUpkeepResult.upkeepFailureReason, 6)

      checkUpkeepResult = await registry
        .connect(zeroAddress)
        .callStatic.checkUpkeep(upkeepID2)
      assert.equal(checkUpkeepResult.upkeepNeeded, true)

      // upkeep 1 perform should revert with max performData size'
      let maxPerformData = '0x'
      for (let i = 0; i < maxPerformDataSize.toNumber(); i++) {
        maxPerformData += '11'
      }
      let configDigest = (await registry.getState()).state.latestConfigDigest
      let checkBlock = await ethers.provider.getBlock('latest')
      let report = encodeReport([
        {
          Id: upkeepID1.toString(),
          checkBlockNum: checkBlock.number,
          checkBlockHash: checkBlock.parentHash,
          performData: maxPerformData,
        },
      ])
      let reportContext = [configDigest, epochAndRound5_1, emptyBytes32]
      let sigs = signReport(reportContext, report, signers.slice(0, f + 1))
      await evmRevert(
        registry
          .connect(keeper1)
          .transmit(
            [reportContext[0], reportContext[1], reportContext[2]],
            report,
            sigs.rs,
            sigs.ss,
            sigs.vs,
            { gasPrice: gasWei.mul(gasCeilingMultiplier) },
          ),
        'StaleReport()',
      )

      // upkeep 1 perform should succeed with empty performData
      checkBlock = await ethers.provider.getBlock('latest')
      report = encodeReport([
        {
          Id: upkeepID1.toString(),
          checkBlockNum: checkBlock.number,
          checkBlockHash: checkBlock.parentHash,
          performData: '0x',
        },
      ])
      reportContext = [configDigest, epochAndRound5_1, emptyBytes32]
      sigs = signReport(reportContext, report, signers.slice(0, f + 1))
      await registry
        .connect(keeper1)
        .transmit(
          [reportContext[0], reportContext[1], reportContext[2]],
          report,
          sigs.rs,
          sigs.ss,
          sigs.vs,
        )

      // upkeep 2 perform should succeed with max performData size
      checkBlock = await ethers.provider.getBlock('latest')
      report = encodeReport([
        {
          Id: upkeepID2.toString(),
          checkBlockNum: checkBlock.number,
          checkBlockHash: checkBlock.parentHash,
          performData: maxPerformData,
        },
      ])
      reportContext = [configDigest, epochAndRound5_1, emptyBytes32]
      sigs = signReport(reportContext, report, signers.slice(0, f + 1))
      await registry
        .connect(keeper1)
        .transmit(
          [reportContext[0], reportContext[1], reportContext[2]],
          report,
          sigs.rs,
          sigs.ss,
          sigs.vs,
        )
    })
  })

  describe('#withdrawFunds', () => {
    beforeEach(async () => {
      await registry.connect(admin).addFunds(upkeepId, toWei('100'))

      // Do a perform so that upkeep is charged some amount
      await registry
        .connect(keeper1)
        .transmit(
          [emptyBytes32, emptyBytes32, emptyBytes32],
          await encodeLatestBlockReport([{ Id: upkeepId.toString() }]),
          [],
          [],
          emptyBytes32,
        )
    })

    it('reverts if called on a non existing ID', async () => {
      await evmRevert(
        registry
          .connect(admin)
          .withdrawFunds(upkeepId.add(1), await payee1.getAddress()),
        'OnlyCallableByAdmin()',
      )
    })

    it('reverts if called by anyone but the admin', async () => {
      await evmRevert(
        registry
          .connect(owner)
          .withdrawFunds(upkeepId, await payee1.getAddress()),
        'OnlyCallableByAdmin()',
      )
    })

    it('reverts if called on an uncanceled upkeep', async () => {
      await evmRevert(
        registry
          .connect(admin)
          .withdrawFunds(upkeepId, await payee1.getAddress()),
        'UpkeepNotCanceled()',
      )
    })

    it('reverts if called with the 0 address', async () => {
      await evmRevert(
        registry.connect(admin).withdrawFunds(upkeepId, zeroAddress),
        'InvalidRecipient()',
      )
    })

    describe('after the registration is cancelled', () => {
      beforeEach(async () => {
        await registry.connect(owner).cancelUpkeep(upkeepId)
      })

      it('moves the funds out and updates the balance and emits an event', async () => {
        const payee1Before = await linkToken.balanceOf(
          await payee1.getAddress(),
        )
        const registryBefore = await linkToken.balanceOf(registry.address)

        let registration = await registry.getUpkeep(upkeepId)
        let previousBalance = registration.balance

        const tx = await registry
          .connect(admin)
          .withdrawFunds(upkeepId, await payee1.getAddress())
        await expect(tx)
          .to.emit(registry, 'FundsWithdrawn')
          .withArgs(upkeepId, previousBalance, await payee1.getAddress())

        const payee1After = await linkToken.balanceOf(await payee1.getAddress())
        const registryAfter = await linkToken.balanceOf(registry.address)

        assert.isTrue(payee1Before.add(previousBalance).eq(payee1After))
        assert.isTrue(registryBefore.sub(previousBalance).eq(registryAfter))

        registration = await registry.getUpkeep(upkeepId)
        assert.equal(0, registration.balance.toNumber())
      })
    })
  })

  describe('#simulatePerformUpkeep', () => {
    it('reverts if called by non zero address', async () => {
      await evmRevert(
        registry
          .connect(await owner.getAddress())
          .callStatic.simulatePerformUpkeep(upkeepId, '0x'),
        'OnlySimulatedBackend()',
      )
    })

    it('reverts when registry is paused', async () => {
      await registry.connect(owner).pause()
      await evmRevert(
        registry
          .connect(zeroAddress)
          .callStatic.simulatePerformUpkeep(upkeepId, '0x'),
        'RegistryPaused()',
      )
    })

    it('returns false and gasUsed when perform fails', async () => {
      await mock.setCanPerform(false)

      let simulatePerformResult = await registry
        .connect(zeroAddress)
        .callStatic.simulatePerformUpkeep(upkeepId, '0x')

      assert.equal(simulatePerformResult.success, false)
      assert.isTrue(simulatePerformResult.gasUsed.gt(BigNumber.from('0'))) // Some gas should be used
    })

    it('returns true and gasUsed when perform succeeds', async () => {
      await mock.setCanPerform(true)

      let simulatePerformResult = await registry
        .connect(zeroAddress)
        .callStatic.simulatePerformUpkeep(upkeepId, '0x')

      assert.equal(simulatePerformResult.success, true)
      assert.isTrue(simulatePerformResult.gasUsed.gt(BigNumber.from('0'))) // Some gas should be used
    })

    it('returns correct amount of gasUsed when perform succeeds', async () => {
      await mock.setCanPerform(true)
      await mock.setPerformGasToBurn(executeGas)

      let simulatePerformResult = await registry
        .connect(zeroAddress)
        .callStatic.simulatePerformUpkeep(upkeepId, '0x')

      assert.equal(simulatePerformResult.success, true)
      // Full execute gas should be used, with some performGasBuffer(1000)
      assert.isTrue(
        simulatePerformResult.gasUsed.gt(
          executeGas.sub(BigNumber.from('1000')),
        ),
      )
    })
  })

  describe('#checkUpkeep', () => {
    it('reverts if called by non zero address', async () => {
      await evmRevert(
        registry
          .connect(await owner.getAddress())
          .callStatic.checkUpkeep(upkeepId),
        'OnlySimulatedBackend()',
      )
    })

    it('returns false and error code if the upkeep is cancelled by admin', async () => {
      await registry.connect(admin).cancelUpkeep(upkeepId)

      let checkUpkeepResult = await registry
        .connect(zeroAddress)
        .callStatic.checkUpkeep(upkeepId)

      assert.equal(checkUpkeepResult.upkeepNeeded, false)
      assert.equal(checkUpkeepResult.performData, '0x')
      assert.equal(checkUpkeepResult.upkeepFailureReason, 1)
      assert.equal(checkUpkeepResult.gasUsed.toString(), '0')
    })

    it('returns false and error code if the upkeep is cancelled by owner', async () => {
      await registry.connect(owner).cancelUpkeep(upkeepId)

      let checkUpkeepResult = await registry
        .connect(zeroAddress)
        .callStatic.checkUpkeep(upkeepId)

      assert.equal(checkUpkeepResult.upkeepNeeded, false)
      assert.equal(checkUpkeepResult.performData, '0x')
      assert.equal(checkUpkeepResult.upkeepFailureReason, 1)
      assert.equal(checkUpkeepResult.gasUsed.toString(), '0')
    })

    it('returns false and error code if the upkeep is paused', async () => {
      await registry.connect(admin).pauseUpkeep(upkeepId)

      let checkUpkeepResult = await registry
        .connect(zeroAddress)
        .callStatic.checkUpkeep(upkeepId)

      assert.equal(checkUpkeepResult.upkeepNeeded, false)
      assert.equal(checkUpkeepResult.performData, '0x')
      assert.equal(checkUpkeepResult.upkeepFailureReason, 2)
      assert.equal(checkUpkeepResult.gasUsed.toString(), '0')
    })

    it('returns false and error code if user is out of funds', async () => {
      let checkUpkeepResult = await registry
        .connect(zeroAddress)
        .callStatic.checkUpkeep(upkeepId)

      assert.equal(checkUpkeepResult.upkeepNeeded, false)
      assert.equal(checkUpkeepResult.performData, '0x')
      assert.equal(checkUpkeepResult.upkeepFailureReason, 6)
      assert.equal(checkUpkeepResult.gasUsed.toString(), '0')
    })

    context('when the registration is funded', () => {
      beforeEach(async () => {
        await linkToken.connect(admin).approve(registry.address, toWei('100'))
        await registry.connect(admin).addFunds(upkeepId, toWei('100'))
      })

      it('returns false and error code if the target check reverts', async () => {
        await mock.setShouldRevertCheck(true)
        let checkUpkeepResult = await registry
          .connect(zeroAddress)
          .callStatic.checkUpkeep(upkeepId)

        assert.equal(checkUpkeepResult.upkeepNeeded, false)
        assert.equal(checkUpkeepResult.performData, '0x')
        assert.equal(checkUpkeepResult.upkeepFailureReason, 3)
        assert.isTrue(checkUpkeepResult.gasUsed.gt(BigNumber.from('0'))) // Some gas should be used
      })

      it('returns false and error code if the upkeep is not needed', async () => {
        await mock.setCanCheck(false)
        let checkUpkeepResult = await registry
          .connect(zeroAddress)
          .callStatic.checkUpkeep(upkeepId)

        assert.equal(checkUpkeepResult.upkeepNeeded, false)
        assert.equal(checkUpkeepResult.performData, '0x')
        assert.equal(checkUpkeepResult.upkeepFailureReason, 4)
        assert.isTrue(checkUpkeepResult.gasUsed.gt(BigNumber.from('0'))) // Some gas should be used
      })

      it('returns false and error code if the performData exceeds limit', async () => {
        let longBytes = '0x'
        for (let i = 0; i < 5000; i++) {
          longBytes += '1'
        }
        await mock.setCanCheck(true)
        await mock.setPerformData(longBytes)

        let checkUpkeepResult = await registry
          .connect(zeroAddress)
          .callStatic.checkUpkeep(upkeepId)

        assert.equal(checkUpkeepResult.upkeepNeeded, false)
        assert.equal(checkUpkeepResult.performData, '0x')
        assert.equal(checkUpkeepResult.upkeepFailureReason, 5)
        assert.isTrue(checkUpkeepResult.gasUsed.gt(BigNumber.from('0'))) // Some gas should be used
      })

      it('returns true with wrapped perform data and gas used if the target can execute', async () => {
        await mock.setCanCheck(true)
        await mock.setPerformData(randomBytes)

        const latestBlock = await ethers.provider.getBlock('latest')

        let checkUpkeepResult = await registry
          .connect(zeroAddress)
          .callStatic.checkUpkeep(upkeepId, {
            blockTag: latestBlock.number,
          })

        let wrappedPerfromData = ethers.utils.defaultAbiCoder.decode(
          [
            'tuple(uint32 checkBlockNum, bytes32 checkBlockHash, bytes performData)',
          ],
          checkUpkeepResult.performData,
        )

        assert.equal(checkUpkeepResult.upkeepNeeded, true)
        assert.equal(wrappedPerfromData[0].checkBlockNum, latestBlock.number)
        assert.equal(wrappedPerfromData[0].performData, randomBytes)
        assert.equal(checkUpkeepResult.upkeepFailureReason, 0)
        assert.isTrue(checkUpkeepResult.gasUsed.gt(BigNumber.from('0'))) // Some gas should be used
      })

      it('has a large enough gas overhead to cover upkeeps that use all their gas', async () => {
        await mock.setCanCheck(true)
        await mock.setCheckGasToBurn(checkGasLimit)
        const gas = checkGasLimit.add(checkGasOverhead)
        let checkUpkeepResult = await registry
          .connect(zeroAddress)
          .callStatic.checkUpkeep(upkeepId, {
            gasLimit: gas,
          })

        assert.equal(checkUpkeepResult.upkeepNeeded, true)
      })
    })
  })

  describe('#addFunds', () => {
    const amount = toWei('1')

    it('reverts if the registration does not exist', async () => {
      await evmRevert(
        registry.connect(keeper1).addFunds(upkeepId.add(1), amount),
        'UpkeepCancelled()',
      )
    })

    it('adds to the balance of the registration', async () => {
      await registry.connect(admin).addFunds(upkeepId, amount)
      const registration = await registry.getUpkeep(upkeepId)
      assert.isTrue(amount.eq(registration.balance))
    })

    it('lets anyone add funds to an upkeep not just admin', async () => {
      await linkToken.connect(owner).transfer(await payee1.getAddress(), amount)
      await linkToken.connect(payee1).approve(registry.address, amount)

      await registry.connect(payee1).addFunds(upkeepId, amount)
      const registration = await registry.getUpkeep(upkeepId)
      assert.isTrue(amount.eq(registration.balance))
    })

    it('emits a log', async () => {
      const tx = await registry.connect(admin).addFunds(upkeepId, amount)
      await expect(tx)
        .to.emit(registry, 'FundsAdded')
        .withArgs(upkeepId, await admin.getAddress(), amount)
    })

    it('reverts if the upkeep is canceled', async () => {
      await registry.connect(admin).cancelUpkeep(upkeepId)
      await evmRevert(
        registry.connect(keeper1).addFunds(upkeepId, amount),
        'UpkeepCancelled()',
      )
    })
  })

  describe('#getActiveUpkeepIDs', () => {
    let upkeepId2: BigNumber

    beforeEach(async () => {
      // Register another upkeep so that we have 2
      const tx = await registry
        .connect(owner)
        .registerUpkeep(
          mock.address,
          executeGas,
          await admin.getAddress(),
          false,
          randomBytes,
        )
      upkeepId2 = await getUpkeepID(tx)
    })

    it('reverts if startIndex is out of bounds ', async () => {
      await evmRevert(registry.getActiveUpkeepIDs(4, 0), 'IndexOutOfRange()')
    })

    it('reverts if startIndex + maxCount is out of bounds', async () => {
      await evmRevert(registry.getActiveUpkeepIDs(0, 4))
    })

    it('returns upkeep IDs bounded by maxCount', async () => {
      let upkeepIds = await registry.getActiveUpkeepIDs(0, 1)
      assert(
        upkeepIds.length == 1,
        'Only maxCount number of upkeeps should be returned',
      )
      assert(
        upkeepIds[0].toString() == upkeepId.toString(),
        'Correct upkeep ID should be returned',
      )

      upkeepIds = await registry.getActiveUpkeepIDs(1, 1)
      assert(
        upkeepIds.length == 1,
        'Only maxCount number of upkeeps should be returned',
      )
      assert(
        upkeepIds[0].toString() == upkeepId2.toString(),
        'Correct upkeep ID should be returned',
      )
    })

    it('returns all upkeep IDs if maxCount is 0', async () => {
      let upkeepIds = await registry.getActiveUpkeepIDs(0, 0)
      assert(upkeepIds.length == 2, 'All upkeeps should be returned')
      assert(
        upkeepIds[0].toString() == upkeepId.toString(),
        'Correct upkeep ID should be returned',
      )
      assert(
        upkeepIds[1].toString() == upkeepId2.toString(),
        'Correct upkeep ID should be returned',
      )
    })
  })

  describe('#getMaxPaymentForGas', () => {
    const multipliers = [BigNumber.from(1), BigNumber.from(3)]
    const gasAmounts = [100000, 10000000]
    const premiums = [0, 250000000]
    const flatFees = [0, 1000000]
    // Same as MockArbGasInfo.sol
    const l1CostWeiArb = BigNumber.from(1000000)
    // Same as MockOVMGasPriceOracle.sol
    const l1CostWeiOpt = BigNumber.from(2000000)

    it('calculates the max fee appropriately', async () => {
      await verifyMaxPayment(0, multipliers, gasAmounts, premiums, flatFees)
    })

    it('calculates the max fee appropriately for Arbitrum', async () => {
      await verifyMaxPayment(
        1,
        multipliers,
        gasAmounts,
        premiums,
        flatFees,
        l1CostWeiArb,
      )
    })

    it('calculates the max fee appropriately for Optimism', async () => {
      await verifyMaxPayment(
        2,
        multipliers,
        gasAmounts,
        premiums,
        flatFees,
        l1CostWeiOpt,
      )
    })
  })

  describe('#typeAndVersion', () => {
    it('uses the correct type and version', async () => {
      const typeAndVersion = await registry.typeAndVersion()
      assert.equal(typeAndVersion, 'KeeperRegistry 2.0.0')
    })
  })

  describe('#onTokenTransfer', () => {
    const amount = toWei('1')

    it('reverts if not called by the LINK token', async () => {
      const data = ethers.utils.defaultAbiCoder.encode(['uint256'], [upkeepId])

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
      await registry.connect(admin).cancelUpkeep(upkeepId)
      await evmRevert(
        registry.connect(keeper1).addFunds(upkeepId, amount),
        'UpkeepCancelled()',
      )
    })

    it('updates the funds of the job id passed', async () => {
      const data = ethers.utils.defaultAbiCoder.encode(['uint256'], [upkeepId])

      const before = (await registry.getUpkeep(upkeepId)).balance
      await linkToken
        .connect(owner)
        .transferAndCall(registry.address, amount, data)
      const after = (await registry.getUpkeep(upkeepId)).balance

      assert.isTrue(before.add(amount).eq(after))
    })
  })

  describe('#setConfig - onchain', () => {
    const payment = BigNumber.from(1)
    const flatFee = BigNumber.from(2)
    const staleness = BigNumber.from(4)
    const ceiling = BigNumber.from(5)
    const maxGas = BigNumber.from(6)
    const fbGasEth = BigNumber.from(7)
    const fbLinkEth = BigNumber.from(8)
    const newMinUpkeepSpend = BigNumber.from(9)
    const newMaxCheckDataSize = BigNumber.from(10000)
    const newMaxPerformDataSize = BigNumber.from(10000)
    const newMaxPerformGas = BigNumber.from(10000000)

    it('reverts when called by anyone but the proposed owner', async () => {
      await evmRevert(
        registry.connect(payee1).setConfig(
          signerAddresses,
          keeperAddresses,
          f,
          encodeConfig({
            paymentPremiumPPB: payment,
            flatFeeMicroLink: flatFee,
            checkGasLimit: maxGas,
            stalenessSeconds: staleness,
            gasCeilingMultiplier: ceiling,
            minUpkeepSpend: newMinUpkeepSpend,
            maxCheckDataSize: newMaxCheckDataSize,
            maxPerformDataSize: newMaxPerformDataSize,
            maxPerformGas: newMaxPerformGas,
            fallbackGasPrice: fbGasEth,
            fallbackLinkPrice: fbLinkEth,
            transcoder: transcoder.address,
            registrar: ethers.constants.AddressZero,
          }),
          offchainVersion,
          offchainBytes,
        ),
        'Only callable by owner',
      )
    })

    it('updates the onchainConfig and configDigest', async () => {
      const old = await registry.getState()
      const oldConfig = old.config
      const oldState = old.state
      assert.isTrue(paymentPremiumPPB.eq(oldConfig.paymentPremiumPPB))
      assert.isTrue(flatFeeMicroLink.eq(oldConfig.flatFeeMicroLink))
      assert.isTrue(stalenessSeconds.eq(oldConfig.stalenessSeconds))
      assert.isTrue(gasCeilingMultiplier.eq(oldConfig.gasCeilingMultiplier))

      await registry.connect(owner).setConfig(
        signerAddresses,
        keeperAddresses,
        f,
        encodeConfig({
          paymentPremiumPPB: payment,
          flatFeeMicroLink: flatFee,
          checkGasLimit: maxGas,
          stalenessSeconds: staleness,
          gasCeilingMultiplier: ceiling,
          minUpkeepSpend: newMinUpkeepSpend,
          maxCheckDataSize: newMaxCheckDataSize,
          maxPerformDataSize: newMaxPerformDataSize,
          maxPerformGas: newMaxPerformGas,
          fallbackGasPrice: fbGasEth,
          fallbackLinkPrice: fbLinkEth,
          transcoder: transcoder.address,
          registrar: ethers.constants.AddressZero,
        }),
        offchainVersion,
        offchainBytes,
      )

      const updated = await registry.getState()
      const updatedConfig = updated.config
      const updatedState = updated.state
      assert.equal(updatedConfig.paymentPremiumPPB, payment.toNumber())
      assert.equal(updatedConfig.flatFeeMicroLink, flatFee.toNumber())
      assert.equal(updatedConfig.stalenessSeconds, staleness.toNumber())
      assert.equal(updatedConfig.gasCeilingMultiplier, ceiling.toNumber())
      assert.equal(
        updatedConfig.minUpkeepSpend.toString(),
        newMinUpkeepSpend.toString(),
      )
      assert.equal(
        updatedConfig.maxCheckDataSize,
        newMaxCheckDataSize.toNumber(),
      )
      assert.equal(
        updatedConfig.maxPerformDataSize,
        newMaxPerformDataSize.toNumber(),
      )
      assert.equal(updatedConfig.maxPerformGas, newMaxPerformGas.toNumber())
      assert.equal(updatedConfig.checkGasLimit, maxGas.toNumber())
      assert.equal(
        updatedConfig.fallbackGasPrice.toNumber(),
        fbGasEth.toNumber(),
      )
      assert.equal(
        updatedConfig.fallbackLinkPrice.toNumber(),
        fbLinkEth.toNumber(),
      )

      assert(oldState.configCount + 1 == updatedState.configCount)
      assert(
        oldState.latestConfigBlockNumber !=
          updatedState.latestConfigBlockNumber,
      )
      assert(oldState.latestConfigDigest != updatedState.latestConfigDigest)
    })

    it('emits an event', async () => {
      const tx = await registry.connect(owner).setConfig(
        signerAddresses,
        keeperAddresses,
        f,
        encodeConfig({
          paymentPremiumPPB: payment,
          flatFeeMicroLink: flatFee,
          checkGasLimit: maxGas,
          stalenessSeconds: staleness,
          gasCeilingMultiplier: ceiling,
          minUpkeepSpend: newMinUpkeepSpend,
          maxCheckDataSize: newMaxCheckDataSize,
          maxPerformDataSize: newMaxPerformDataSize,
          maxPerformGas: newMaxPerformGas,
          fallbackGasPrice: fbGasEth,
          fallbackLinkPrice: fbLinkEth,
          transcoder: transcoder.address,
          registrar: ethers.constants.AddressZero,
        }),
        offchainVersion,
        offchainBytes,
      )
      await expect(tx).to.emit(registry, 'ConfigSet')
    })

    it('reverts upon decreasing max limits', async () => {
      await evmRevert(
        registry.connect(owner).setConfig(
          signerAddresses,
          keeperAddresses,
          f,
          encodeConfig({
            paymentPremiumPPB: payment,
            flatFeeMicroLink: flatFee,
            checkGasLimit: maxGas,
            stalenessSeconds: staleness,
            gasCeilingMultiplier: ceiling,
            minUpkeepSpend: newMinUpkeepSpend,
            maxCheckDataSize: BigNumber.from(1),
            maxPerformDataSize: newMaxPerformDataSize,
            maxPerformGas: newMaxPerformGas,
            fallbackGasPrice: fbGasEth,
            fallbackLinkPrice: fbLinkEth,
            transcoder: transcoder.address,
            registrar: ethers.constants.AddressZero,
          }),
          offchainVersion,
          offchainBytes,
        ),
        'MaxCheckDataSizeCanOnlyIncrease()',
      )
      await evmRevert(
        registry.connect(owner).setConfig(
          signerAddresses,
          keeperAddresses,
          f,
          encodeConfig({
            paymentPremiumPPB: payment,
            flatFeeMicroLink: flatFee,
            checkGasLimit: maxGas,
            stalenessSeconds: staleness,
            gasCeilingMultiplier: ceiling,
            minUpkeepSpend: newMinUpkeepSpend,
            maxCheckDataSize: newMaxCheckDataSize,
            maxPerformDataSize: BigNumber.from(1),
            maxPerformGas: newMaxPerformGas,
            fallbackGasPrice: fbGasEth,
            fallbackLinkPrice: fbLinkEth,
            transcoder: transcoder.address,
            registrar: ethers.constants.AddressZero,
          }),
          offchainVersion,
          offchainBytes,
        ),
        'MaxPerformDataSizeCanOnlyIncrease()',
      )
      await evmRevert(
        registry.connect(owner).setConfig(
          signerAddresses,
          keeperAddresses,
          f,
          encodeConfig({
            paymentPremiumPPB: payment,
            flatFeeMicroLink: flatFee,
            checkGasLimit: maxGas,
            stalenessSeconds: staleness,
            gasCeilingMultiplier: ceiling,
            minUpkeepSpend: newMinUpkeepSpend,
            maxCheckDataSize: newMaxCheckDataSize,
            maxPerformDataSize: newMaxPerformDataSize,
            maxPerformGas: BigNumber.from(1),
            fallbackGasPrice: fbGasEth,
            fallbackLinkPrice: fbLinkEth,
            transcoder: transcoder.address,
            registrar: ethers.constants.AddressZero,
          }),
          offchainVersion,
          offchainBytes,
        ),
        'GasLimitCanOnlyIncrease()',
      )
    })
  })

  describe('#setConfig - offchain', () => {
    let newKeepers: string[]

    beforeEach(async () => {
      newKeepers = [
        await personas.Eddy.getAddress(),
        await personas.Nick.getAddress(),
        await personas.Neil.getAddress(),
        await personas.Carol.getAddress(),
      ]
    })

    it('reverts when called by anyone but the owner', async () => {
      await evmRevert(
        registry
          .connect(payee1)
          .setConfig(
            newKeepers,
            newKeepers,
            f,
            encodeConfig(config),
            offchainVersion,
            offchainBytes,
          ),
        'Only callable by owner',
      )
    })

    it('reverts if too many keeperAddresses set', async () => {
      for (let i = 0; i < 40; i++) {
        newKeepers.push(randomAddress())
      }
      await evmRevert(
        registry
          .connect(owner)
          .setConfig(
            newKeepers,
            newKeepers,
            f,
            encodeConfig(config),
            offchainVersion,
            offchainBytes,
          ),
        'TooManyOracles()',
      )
    })

    it('reverts if f=0', async () => {
      await evmRevert(
        registry
          .connect(owner)
          .setConfig(
            newKeepers,
            newKeepers,
            0,
            encodeConfig(config),
            offchainVersion,
            offchainBytes,
          ),
        'IncorrectNumberOfFaultyOracles()',
      )
    })

    it('reverts if signers != transmitters length', async () => {
      let signers = [randomAddress()]
      await evmRevert(
        registry
          .connect(owner)
          .setConfig(
            signers,
            newKeepers,
            f,
            encodeConfig(config),
            offchainVersion,
            offchainBytes,
          ),
        'IncorrectNumberOfSigners()',
      )
    })

    it('reverts if signers <= 3f', async () => {
      newKeepers.pop()
      await evmRevert(
        registry
          .connect(owner)
          .setConfig(
            newKeepers,
            newKeepers,
            f,
            encodeConfig(config),
            offchainVersion,
            offchainBytes,
          ),
        'IncorrectNumberOfSigners()',
      )
    })

    it('reverts on repeated signers', async () => {
      let newSigners = [
        await personas.Eddy.getAddress(),
        await personas.Eddy.getAddress(),
        await personas.Eddy.getAddress(),
        await personas.Eddy.getAddress(),
      ]
      await evmRevert(
        registry
          .connect(owner)
          .setConfig(
            newSigners,
            newKeepers,
            f,
            encodeConfig(config),
            offchainVersion,
            offchainBytes,
          ),
        'RepeatedSigner()',
      )
    })

    it('reverts on repeated transmitters', async () => {
      let newTransmitters = [
        await personas.Eddy.getAddress(),
        await personas.Eddy.getAddress(),
        await personas.Eddy.getAddress(),
        await personas.Eddy.getAddress(),
      ]
      await evmRevert(
        registry
          .connect(owner)
          .setConfig(
            newKeepers,
            newTransmitters,
            f,
            encodeConfig(config),
            offchainVersion,
            offchainBytes,
          ),
        'RepeatedTransmitter()',
      )
    })

    it('stores new config and emits event', async () => {
      let newOffChainVersion = BigNumber.from('2')
      let newOffChainConfig = '0x1122'

      const old = await registry.getState()
      const oldState = old.state

      const tx = await registry
        .connect(owner)
        .setConfig(
          newKeepers,
          newKeepers,
          f,
          encodeConfig(config),
          newOffChainVersion,
          newOffChainConfig,
        )

      const updated = await registry.getState()

      const updatedState = updated.state

      // Old signer addresses which are not in new signers should be non active
      for (var i = 0; i < keeperAddresses.length; i++) {
        let signer = keeperAddresses[i]
        if (!newKeepers.includes(signer)) {
          assert((await registry.getSignerInfo(signer)).active == false)
          assert((await registry.getSignerInfo(signer)).index == 0)
        }
      }
      // New signer addresses should be active
      for (var i = 0; i < newKeepers.length; i++) {
        let signer = newKeepers[i]
        assert((await registry.getSignerInfo(signer)).active == true)
        assert((await registry.getSignerInfo(signer)).index == i)
      }
      // Old transmitter addresses which are not in new transmitter should be non active, but retain other info
      for (var i = 0; i < keeperAddresses.length; i++) {
        let transmitter = keeperAddresses[i]
        if (!newKeepers.includes(transmitter)) {
          assert(
            (await registry.getTransmitterInfo(transmitter)).active == false,
          )
          assert((await registry.getTransmitterInfo(transmitter)).index == i)
        }
      }
      // New transmitter addresses should be active
      for (var i = 0; i < newKeepers.length; i++) {
        let transmitter = newKeepers[i]
        assert((await registry.getTransmitterInfo(transmitter)).active == true)
        assert((await registry.getTransmitterInfo(transmitter)).index == i)
      }

      // config digest should be updated
      assert(oldState.configCount + 1 == updatedState.configCount)
      assert(
        oldState.latestConfigBlockNumber !=
          updatedState.latestConfigBlockNumber,
      )
      assert(oldState.latestConfigDigest != updatedState.latestConfigDigest)

      //New config should be updated
      assert.deepEqual(updated.signers, newKeepers)
      assert.deepEqual(updated.transmitters, newKeepers)

      // Event should have been emitted
      await expect(tx).to.emit(registry, 'ConfigSet')
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

  describe('#registerUpkeep', () => {
    it('reverts when registry is paused', async () => {
      await registry.connect(owner).pause()
      await evmRevert(
        registry
          .connect(owner)
          .registerUpkeep(
            mock.address,
            executeGas,
            await admin.getAddress(),
            false,
            emptyBytes,
          ),
        'RegistryPaused()',
      )
    })

    it('reverts if the target is not a contract', async () => {
      await evmRevert(
        registry
          .connect(owner)
          .registerUpkeep(
            zeroAddress,
            executeGas,
            await admin.getAddress(),
            false,
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
            false,
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
            false,
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
            false,
            emptyBytes,
          ),
        'GasLimitOutsideRange()',
      )
    })

    it('reverts if checkData is too long', async () => {
      let longBytes = '0x'
      for (let i = 0; i < 10000; i++) {
        longBytes += '1'
      }
      await evmRevert(
        registry
          .connect(owner)
          .registerUpkeep(
            mock.address,
            executeGas,
            await admin.getAddress(),
            false,
            longBytes,
          ),
        'CheckDataExceedsLimit()',
      )
    })

    it('creates a record of the registration', async () => {
      let skipSigVerifications = [true, false]
      let executeGases = [100000, 500000]
      let checkDatas = [emptyBytes, '0x12']

      for (let idx = 0; idx < skipSigVerifications.length; idx++) {
        const skipSigVerification = skipSigVerifications[idx]
        for (let jdx = 0; jdx < executeGases.length; jdx++) {
          const executeGas = executeGases[jdx]
          for (let kdx = 0; kdx < checkDatas.length; kdx++) {
            const checkData = checkDatas[kdx]
            const tx = await registry
              .connect(owner)
              .registerUpkeep(
                mock.address,
                executeGas,
                await admin.getAddress(),
                skipSigVerification,
                checkData,
              )

            //confirm the upkeep details
            upkeepId = await getUpkeepID(tx)
            await expect(tx)
              .to.emit(registry, 'UpkeepRegistered')
              .withArgs(
                upkeepId,
                executeGas,
                await admin.getAddress(),
                skipSigVerification,
              )
            const registration = await registry.getUpkeep(upkeepId)

            assert.equal(mock.address, registration.target)
            assert.equal(
              executeGas.toString(),
              registration.executeGas.toString(),
            )
            assert.equal(await admin.getAddress(), registration.admin)
            assert.equal(0, registration.balance.toNumber())
            assert.equal(0, registration.amountSpent.toNumber())
            assert.equal(0, registration.lastPerformBlockNumber)
            assert.equal(checkData, registration.checkData)
            assert.equal(registration.paused, false)
            assert.equal(registration.skipSigVerification, skipSigVerification)
            assert(registration.maxValidBlocknumber.eq('0xffffffff'))
          }
        }
      }
    })
  })

  describe('#pauseUpkeep', () => {
    it('reverts if the registration does not exist', async () => {
      await evmRevert(
        registry.connect(keeper1).pauseUpkeep(upkeepId.add(1)),
        'OnlyCallableByAdmin()',
      )
    })

    it('reverts if the upkeep is already canceled', async () => {
      await registry.connect(admin).cancelUpkeep(upkeepId)

      await evmRevert(
        registry.connect(admin).pauseUpkeep(upkeepId),
        'UpkeepCancelled()',
      )
    })

    it('reverts if the upkeep is already paused', async () => {
      await registry.connect(admin).pauseUpkeep(upkeepId)

      await evmRevert(
        registry.connect(admin).pauseUpkeep(upkeepId),
        'OnlyUnpausedUpkeep()',
      )
    })

    it('reverts if the caller is not the upkeep admin', async () => {
      await evmRevert(
        registry.connect(keeper1).pauseUpkeep(upkeepId),
        'OnlyCallableByAdmin()',
      )
    })

    it('pauses the upkeep and emits an event', async () => {
      const tx = await registry.connect(admin).pauseUpkeep(upkeepId)
      await expect(tx).to.emit(registry, 'UpkeepPaused').withArgs(upkeepId)

      const registration = await registry.getUpkeep(upkeepId)
      assert.equal(registration.paused, true)
    })
  })

  describe('#unpauseUpkeep', () => {
    it('reverts if the registration does not exist', async () => {
      await evmRevert(
        registry.connect(keeper1).unpauseUpkeep(upkeepId.add(1)),
        'OnlyCallableByAdmin()',
      )
    })

    it('reverts if the upkeep is already canceled', async () => {
      await registry.connect(owner).cancelUpkeep(upkeepId)

      await evmRevert(
        registry.connect(admin).unpauseUpkeep(upkeepId),
        'UpkeepCancelled()',
      )
    })

    it('marks the contract as paused', async () => {
      assert.isFalse((await registry.getState()).state.paused)

      await registry.connect(owner).pause()

      assert.isTrue((await registry.getState()).state.paused)
    })

    it('reverts if the upkeep is not paused', async () => {
      await evmRevert(
        registry.connect(admin).unpauseUpkeep(upkeepId),
        'OnlyPausedUpkeep()',
      )
    })

    it('reverts if the caller is not the upkeep admin', async () => {
      await registry.connect(admin).pauseUpkeep(upkeepId)

      const registration = await registry.getUpkeep(upkeepId)

      assert.equal(registration.paused, true)

      await evmRevert(
        registry.connect(keeper1).unpauseUpkeep(upkeepId),
        'OnlyCallableByAdmin()',
      )
    })

    it('unpauses the upkeep and emits an event', async () => {
      await registry.connect(admin).pauseUpkeep(upkeepId)

      const tx = await registry.connect(admin).unpauseUpkeep(upkeepId)

      await expect(tx).to.emit(registry, 'UpkeepUnpaused').withArgs(upkeepId)

      const registration = await registry.getUpkeep(upkeepId)
      assert.equal(registration.paused, false)

      const upkeepIds = await registry.getActiveUpkeepIDs(0, 0)
      assert.equal(upkeepIds.length, 1)
    })
  })

  describe('#updateCheckData', () => {
    it('reverts if the registration does not exist', async () => {
      await evmRevert(
        registry.connect(keeper1).updateCheckData(upkeepId.add(1), randomBytes),
        'OnlyCallableByAdmin()',
      )
    })

    it('reverts if the caller is not upkeep admin', async () => {
      await evmRevert(
        registry.connect(keeper1).updateCheckData(upkeepId, randomBytes),
        'OnlyCallableByAdmin()',
      )
    })

    it('reverts if the upkeep is cancelled', async () => {
      await registry.connect(admin).cancelUpkeep(upkeepId)

      await evmRevert(
        registry.connect(admin).updateCheckData(upkeepId, randomBytes),
        'UpkeepCancelled()',
      )
    })

    it('is allowed to update on paused upkeep', async () => {
      await registry.connect(admin).pauseUpkeep(upkeepId)
      await registry.connect(admin).updateCheckData(upkeepId, randomBytes)

      const registration = await registry.getUpkeep(upkeepId)
      assert.equal(randomBytes, registration.checkData)
    })

    it('reverts if newCheckData exceeds limit', async () => {
      let longBytes = '0x'
      for (let i = 0; i < 10000; i++) {
        longBytes += '1'
      }

      await evmRevert(
        registry.connect(admin).updateCheckData(upkeepId, longBytes),
        'CheckDataExceedsLimit()',
      )
    })

    it('updates the upkeep check data and emits an event', async () => {
      const tx = await registry
        .connect(admin)
        .updateCheckData(upkeepId, randomBytes)
      await expect(tx)
        .to.emit(registry, 'UpkeepCheckDataUpdated')
        .withArgs(upkeepId, randomBytes)

      const registration = await registry.getUpkeep(upkeepId)
      assert.equal(randomBytes, registration.checkData)
    })
  })

  describe('#setUpkeepGasLimit', () => {
    const newGasLimit = BigNumber.from('300000')

    it('reverts if the registration does not exist', async () => {
      await evmRevert(
        registry.connect(admin).setUpkeepGasLimit(upkeepId.add(1), newGasLimit),
        'OnlyCallableByAdmin()',
      )
    })

    it('reverts if the upkeep is canceled', async () => {
      await registry.connect(admin).cancelUpkeep(upkeepId)
      await evmRevert(
        registry.connect(admin).setUpkeepGasLimit(upkeepId, newGasLimit),
        'UpkeepCancelled()',
      )
    })

    it('reverts if called by anyone but the admin', async () => {
      await evmRevert(
        registry.connect(owner).setUpkeepGasLimit(upkeepId, newGasLimit),
        'OnlyCallableByAdmin()',
      )
    })

    it('reverts if new gas limit is out of bounds', async () => {
      await evmRevert(
        registry
          .connect(admin)
          .setUpkeepGasLimit(upkeepId, BigNumber.from('100')),
        'GasLimitOutsideRange()',
      )
      await evmRevert(
        registry
          .connect(admin)
          .setUpkeepGasLimit(upkeepId, BigNumber.from('6000000')),
        'GasLimitOutsideRange()',
      )
    })

    it('updates the gas limit successfully', async () => {
      const initialGasLimit = (await registry.getUpkeep(upkeepId)).executeGas
      assert.equal(initialGasLimit, executeGas.toNumber())
      await registry.connect(admin).setUpkeepGasLimit(upkeepId, newGasLimit)
      const updatedGasLimit = (await registry.getUpkeep(upkeepId)).executeGas
      assert.equal(updatedGasLimit, newGasLimit.toNumber())
    })

    it('emits a log', async () => {
      const tx = await registry
        .connect(admin)
        .setUpkeepGasLimit(upkeepId, newGasLimit)
      await expect(tx)
        .to.emit(registry, 'UpkeepGasLimitSet')
        .withArgs(upkeepId, newGasLimit)
    })
  })

  describe('#transferUpkeepAdmin', () => {
    it('reverts when called by anyone but the current upkeep admin', async () => {
      await evmRevert(
        registry
          .connect(payee1)
          .transferUpkeepAdmin(upkeepId, await payee2.getAddress()),
        'OnlyCallableByAdmin()',
      )
    })

    it('reverts when transferring to self', async () => {
      await evmRevert(
        registry
          .connect(admin)
          .transferUpkeepAdmin(upkeepId, await admin.getAddress()),
        'ValueNotChanged()',
      )
    })

    it('reverts when the upkeep is cancelled', async () => {
      await registry.connect(admin).cancelUpkeep(upkeepId)

      await evmRevert(
        registry
          .connect(admin)
          .transferUpkeepAdmin(upkeepId, await keeper1.getAddress()),
        'UpkeepCancelled()',
      )
    })

    it('reverts when transferring to zero address', async () => {
      await evmRevert(
        registry
          .connect(admin)
          .transferUpkeepAdmin(upkeepId, ethers.constants.AddressZero),
        'InvalidRecipient()',
      )
    })

    it('does not change the upkeep admin', async () => {
      await registry
        .connect(admin)
        .transferUpkeepAdmin(upkeepId, await payee1.getAddress())

      const upkeep = await registry.getUpkeep(upkeepId)
      assert.equal(await admin.getAddress(), upkeep.admin)
    })

    it('emits an event announcing the new upkeep admin', async () => {
      const tx = await registry
        .connect(admin)
        .transferUpkeepAdmin(upkeepId, await payee1.getAddress())

      await expect(tx)
        .to.emit(registry, 'UpkeepAdminTransferRequested')
        .withArgs(upkeepId, await admin.getAddress(), await payee1.getAddress())
    })

    it('does not emit an event when called with the same proposed upkeep admin', async () => {
      await registry
        .connect(admin)
        .transferUpkeepAdmin(upkeepId, await payee1.getAddress())

      const tx = await registry
        .connect(admin)
        .transferUpkeepAdmin(upkeepId, await payee1.getAddress())
      const receipt = await tx.wait()
      assert.equal(0, receipt.logs.length)
    })
  })

  describe('#acceptUpkeepAdmin', () => {
    beforeEach(async () => {
      // Start admin transfer to payee1
      await registry
        .connect(admin)
        .transferUpkeepAdmin(upkeepId, await payee1.getAddress())
    })

    it('reverts when not called by the proposed upkeep admin', async () => {
      await evmRevert(
        registry.connect(payee2).acceptUpkeepAdmin(upkeepId),
        'OnlyCallableByProposedAdmin()',
      )
    })

    it('reverts when the upkeep is cancelled', async () => {
      await registry.connect(admin).cancelUpkeep(upkeepId)

      await evmRevert(
        registry.connect(payee1).acceptUpkeepAdmin(upkeepId),
        'UpkeepCancelled()',
      )
    })

    it('does change the admin', async () => {
      await registry.connect(payee1).acceptUpkeepAdmin(upkeepId)

      const upkeep = await registry.getUpkeep(upkeepId)
      assert.equal(await payee1.getAddress(), upkeep.admin)
    })

    it('emits an event announcing the new upkeep admin', async () => {
      const tx = await registry.connect(payee1).acceptUpkeepAdmin(upkeepId)
      await expect(tx)
        .to.emit(registry, 'UpkeepAdminTransferred')
        .withArgs(upkeepId, await admin.getAddress(), await payee1.getAddress())
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
      await registry.connect(admin).addFunds(upkeepId, toWei('100'))
      // Very high min spend, whole balance as cancellation fees
      let minUpkeepSpend = toWei('1000')
      await registry.connect(owner).setConfig(
        signerAddresses,
        keeperAddresses,
        f,
        encodeConfig({
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
        }),
        offchainVersion,
        offchainBytes,
      )
      let upkeepBalance = (await registry.getUpkeep(upkeepId)).balance
      const ownerBefore = await linkToken.balanceOf(await owner.getAddress())

      await registry.connect(owner).cancelUpkeep(upkeepId)

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

      const info = await registry.getTransmitterInfo(await keeper1.getAddress())
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

      const info = await registry.getTransmitterInfo(await keeper1.getAddress())
      assert.equal(await payee2.getAddress(), info.payee)
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
      assert.isFalse((await registry.getState()).state.paused)

      await registry.connect(owner).pause()

      assert.isTrue((await registry.getState()).state.paused)
    })

    it('Does not allow transmits when paused', async () => {
      await registry.connect(owner).pause()

      await evmRevert(
        registry.transmit(
          [emptyBytes32, emptyBytes32, emptyBytes32],
          emptyBytes,
          [],
          [],
          emptyBytes32,
        ),
        'RegistryPaused()',
      )
    })

    it('Does not allow creation of new upkeeps when paused', async () => {
      await registry.connect(owner).pause()

      await evmRevert(
        registry
          .connect(owner)
          .registerUpkeep(
            mock.address,
            executeGas,
            await admin.getAddress(),
            false,
            emptyBytes,
          ),
        'RegistryPaused()',
      )
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
      assert.isTrue((await registry.getState()).state.paused)

      await registry.connect(owner).unpause()

      assert.isFalse((await registry.getState()).state.paused)
    })
  })

  describe('migrateUpkeeps() / #receiveUpkeeps()', async () => {
    let registry2: KeeperRegistry
    let registryLogic2: KeeperRegistryLogic

    beforeEach(async () => {
      registryLogic2 = await keeperRegistryLogicFactory
        .connect(owner)
        .deploy(0, linkToken.address, linkEthFeed.address, gasPriceFeed.address)

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
      registry2 = await keeperRegistryFactory
        .connect(owner)
        .deploy(registryLogic2.address)
      await registry2
        .connect(owner)
        .setConfig(
          signerAddresses,
          keeperAddresses,
          f,
          encodeConfig(config),
          1,
          '0x',
        )
    })

    context('when permissions are set', () => {
      beforeEach(async () => {
        await linkToken.connect(owner).approve(registry.address, toWei('100'))
        await registry.connect(owner).addFunds(upkeepId, toWei('100'))
        await registry.setPeerRegistryMigrationPermission(registry2.address, 1)
        await registry2.setPeerRegistryMigrationPermission(registry.address, 2)
      })

      it('migrates an upkeep', async () => {
        expect((await registry.getUpkeep(upkeepId)).balance).to.equal(
          toWei('100'),
        )
        expect((await registry.getUpkeep(upkeepId)).checkData).to.equal(
          randomBytes,
        )
        expect((await registry.getState()).state.numUpkeeps).to.equal(1)
        // Set an upkeep admin transfer in progress too
        await registry
          .connect(admin)
          .transferUpkeepAdmin(upkeepId, await payee1.getAddress())

        // migrate
        await registry
          .connect(admin)
          .migrateUpkeeps([upkeepId], registry2.address)
        expect((await registry.getState()).state.numUpkeeps).to.equal(0)
        expect((await registry2.getState()).state.numUpkeeps).to.equal(1)
        expect((await registry.getUpkeep(upkeepId)).balance).to.equal(0)
        expect((await registry.getUpkeep(upkeepId)).checkData).to.equal('0x')
        expect((await registry2.getUpkeep(upkeepId)).balance).to.equal(
          toWei('100'),
        )
        expect((await registry2.getState()).state.expectedLinkBalance).to.equal(
          toWei('100'),
        )
        expect((await registry2.getUpkeep(upkeepId)).checkData).to.equal(
          randomBytes,
        )
        // migration will delete the upkeep and nullify admin transfer
        await expect(
          registry.connect(payee1).acceptUpkeepAdmin(upkeepId),
        ).to.be.revertedWith('UpkeepCancelled()')
        await expect(
          registry2.connect(payee1).acceptUpkeepAdmin(upkeepId),
        ).to.be.revertedWith('OnlyCallableByProposedAdmin()')
      })

      it('migrates a paused upkeep', async () => {
        expect((await registry.getUpkeep(upkeepId)).balance).to.equal(
          toWei('100'),
        )
        expect((await registry.getUpkeep(upkeepId)).checkData).to.equal(
          randomBytes,
        )
        expect((await registry.getState()).state.numUpkeeps).to.equal(1)
        await registry.connect(admin).pauseUpkeep(upkeepId)
        // verify the upkeep is paused
        expect((await registry.getUpkeep(upkeepId)).paused).to.equal(true)
        // migrate
        await registry
          .connect(admin)
          .migrateUpkeeps([upkeepId], registry2.address)
        expect((await registry.getState()).state.numUpkeeps).to.equal(0)
        expect((await registry2.getState()).state.numUpkeeps).to.equal(1)
        expect((await registry.getUpkeep(upkeepId)).balance).to.equal(0)
        expect((await registry2.getUpkeep(upkeepId)).balance).to.equal(
          toWei('100'),
        )
        expect((await registry.getUpkeep(upkeepId)).checkData).to.equal('0x')
        expect((await registry2.getUpkeep(upkeepId)).checkData).to.equal(
          randomBytes,
        )
        expect((await registry2.getState()).state.expectedLinkBalance).to.equal(
          toWei('100'),
        )
        // verify the upkeep is still paused after migration
        expect((await registry2.getUpkeep(upkeepId)).paused).to.equal(true)
      })

      it('emits an event on both contracts', async () => {
        expect((await registry.getUpkeep(upkeepId)).balance).to.equal(
          toWei('100'),
        )
        expect((await registry.getUpkeep(upkeepId)).checkData).to.equal(
          randomBytes,
        )
        expect((await registry.getState()).state.numUpkeeps).to.equal(1)
        const tx = registry
          .connect(admin)
          .migrateUpkeeps([upkeepId], registry2.address)
        await expect(tx)
          .to.emit(registry, 'UpkeepMigrated')
          .withArgs(upkeepId, toWei('100'), registry2.address)
        await expect(tx)
          .to.emit(registry2, 'UpkeepReceived')
          .withArgs(upkeepId, toWei('100'), registry.address)
      })

      it('is only migratable by the admin', async () => {
        await expect(
          registry.connect(owner).migrateUpkeeps([upkeepId], registry2.address),
        ).to.be.revertedWith('OnlyCallableByAdmin()')
        await registry
          .connect(admin)
          .migrateUpkeeps([upkeepId], registry2.address)
      })
    })

    context('when permissions are not set', () => {
      it('reverts', async () => {
        // no permissions
        await registry.setPeerRegistryMigrationPermission(registry2.address, 0)
        await registry2.setPeerRegistryMigrationPermission(registry.address, 0)
        await expect(registry.migrateUpkeeps([upkeepId], registry2.address)).to
          .be.reverted
        // only outgoing permissions
        await registry.setPeerRegistryMigrationPermission(registry2.address, 1)
        await registry2.setPeerRegistryMigrationPermission(registry.address, 0)
        await expect(registry.migrateUpkeeps([upkeepId], registry2.address)).to
          .be.reverted
        // only incoming permissions
        await registry.setPeerRegistryMigrationPermission(registry2.address, 0)
        await registry2.setPeerRegistryMigrationPermission(registry.address, 2)
        await expect(registry.migrateUpkeeps([upkeepId], registry2.address)).to
          .be.reverted
        // permissions opposite direction
        await registry.setPeerRegistryMigrationPermission(registry2.address, 2)
        await registry2.setPeerRegistryMigrationPermission(registry.address, 1)
        await expect(registry.migrateUpkeeps([upkeepId], registry2.address)).to
          .be.reverted
      })
    })
  })

  describe('#setPayees', () => {
    const IGNORE_ADDRESS = '0xFFfFfFffFFfffFFfFFfFFFFFffFFFffffFfFFFfF'

    beforeEach(async () => {
      keeperAddresses = keeperAddresses.slice(0, 4)
      signerAddresses = signerAddresses.slice(0, 4)
      payees = payees.slice(0, 4)

      // Redeploy registry with zero address payees (non set)
      registry = await keeperRegistryFactory
        .connect(owner)
        .deploy(registryLogic.address)

      await registry
        .connect(owner)
        .setConfig(
          signerAddresses,
          keeperAddresses,
          f,
          encodeConfig(config),
          offchainVersion,
          offchainBytes,
        )
    })

    it('reverts when not called by the owner', async () => {
      await evmRevert(
        registry.connect(keeper1).setPayees([]),
        'Only callable by owner',
      )
    })

    it('reverts with different numbers of payees than transmitters', async () => {
      // 4 transmitters are set, so exactly 4 payess should be added
      await evmRevert(
        registry.connect(owner).setPayees([await payee1.getAddress()]),
        'ParameterLengthError()',
      )
      await evmRevert(
        registry
          .connect(owner)
          .setPayees([
            await payee1.getAddress(),
            await payee1.getAddress(),
            await payee1.getAddress(),
            await payee1.getAddress(),
            await payee1.getAddress(),
          ]),
        'ParameterLengthError()',
      )
    })

    it('reverts if the payee is the zero address', async () => {
      await evmRevert(
        registry
          .connect(owner)
          .setPayees([
            await payee1.getAddress(),
            '0x0000000000000000000000000000000000000000',
            await payee3.getAddress(),
            await payee4.getAddress(),
          ]),
        'InvalidPayee()',
      )
    })

    it('sets the payees when exisitng payees are zero address', async () => {
      //Initial payees should be zero address
      for (let i = 0; i < keeperAddresses.length; i++) {
        let payee = (await registry.getTransmitterInfo(keeperAddresses[i]))
          .payee
        assert.equal(payee, zeroAddress)
      }

      await registry.connect(owner).setPayees(payees)

      for (let i = 0; i < keeperAddresses.length; i++) {
        let payee = (await registry.getTransmitterInfo(keeperAddresses[i]))
          .payee
        assert.equal(payee, payees[i])
      }
    })

    it('does not change the payee if IGNORE_ADDRESS is used as payee', async () => {
      // Set initial payees
      await registry.connect(owner).setPayees(payees)

      const newPayees = [
        await payee1.getAddress(),
        IGNORE_ADDRESS,
        await payee3.getAddress(),
        await payee4.getAddress(),
      ]
      await registry.connect(owner).setPayees(newPayees)

      const ignored = await registry.getTransmitterInfo(
        await keeper2.getAddress(),
      )
      assert.equal(await payee2.getAddress(), ignored.payee)
      assert.equal(true, ignored.active)
    })

    it('reverts if payee is non zero and owner tries to change payee', async () => {
      // Set initial payees
      await registry.connect(owner).setPayees(payees)

      const newPayees = [
        await payee1.getAddress(),
        await owner.getAddress(),
        await payee3.getAddress(),
        await payee4.getAddress(),
      ]
      await evmRevert(
        registry.connect(owner).setPayees(newPayees),
        'InvalidPayee()',
      )
    })

    it('emits events for every payee added and removed', async () => {
      const tx = await registry.connect(owner).setPayees(payees)
      await expect(tx)
        .to.emit(registry, 'PayeesUpdated')
        .withArgs(keeperAddresses, payees)
    })
  })

  describe('#cancelUpkeep', () => {
    it('reverts if the ID is not valid', async () => {
      await evmRevert(
        registry.connect(owner).cancelUpkeep(upkeepId.add(1)),
        'CannotCancel()',
      )
    })

    it('reverts if called by a non-owner/non-admin', async () => {
      await evmRevert(
        registry.connect(keeper1).cancelUpkeep(upkeepId),
        'OnlyCallableByOwnerOrAdmin()',
      )
    })

    describe('when called by the owner', async () => {
      it('sets the registration to invalid immediately', async () => {
        const tx = await registry.connect(owner).cancelUpkeep(upkeepId)
        const receipt = await tx.wait()
        const registration = await registry.getUpkeep(upkeepId)
        assert.equal(
          registration.maxValidBlocknumber.toNumber(),
          receipt.blockNumber,
        )
      })

      it('emits an event', async () => {
        const tx = await registry.connect(owner).cancelUpkeep(upkeepId)
        const receipt = await tx.wait()
        await expect(tx)
          .to.emit(registry, 'UpkeepCanceled')
          .withArgs(upkeepId, BigNumber.from(receipt.blockNumber))
      })

      it('immediately prevents upkeep', async () => {
        await registry.connect(owner).cancelUpkeep(upkeepId)

        await evmRevert(
          registry
            .connect(keeper1)
            .transmit(
              [emptyBytes32, emptyBytes32, emptyBytes32],
              await encodeLatestBlockReport([{ Id: upkeepId.toString() }]),
              [],
              [],
              emptyBytes32,
            ),
          'StaleReport()',
        )
      })

      it('does not revert if reverts if called multiple times', async () => {
        await registry.connect(owner).cancelUpkeep(upkeepId)
        await evmRevert(
          registry.connect(owner).cancelUpkeep(upkeepId),
          'CannotCancel()',
        )
      })

      describe('when called by the owner when the admin has just canceled', () => {
        let oldExpiration: BigNumber

        beforeEach(async () => {
          await registry.connect(admin).cancelUpkeep(upkeepId)
          const registration = await registry.getUpkeep(upkeepId)
          oldExpiration = registration.maxValidBlocknumber
        })

        it('allows the owner to cancel it more quickly', async () => {
          await registry.connect(owner).cancelUpkeep(upkeepId)

          const registration = await registry.getUpkeep(upkeepId)
          const newExpiration = registration.maxValidBlocknumber
          assert.isTrue(newExpiration.lt(oldExpiration))
        })
      })
    })

    describe('when called by the admin', async () => {
      it('reverts if called again by the admin', async () => {
        await registry.connect(admin).cancelUpkeep(upkeepId)

        await evmRevert(
          registry.connect(admin).cancelUpkeep(upkeepId),
          'CannotCancel()',
        )
      })

      it('reverts if called by the owner after the timeout', async () => {
        await registry.connect(admin).cancelUpkeep(upkeepId)

        for (let i = 0; i < cancellationDelay; i++) {
          await ethers.provider.send('evm_mine', [])
        }

        await evmRevert(
          registry.connect(owner).cancelUpkeep(upkeepId),
          'CannotCancel()',
        )
      })

      it('sets the registration to invalid in 50 blocks', async () => {
        const tx = await registry.connect(admin).cancelUpkeep(upkeepId)
        const receipt = await tx.wait()
        const registration = await registry.getUpkeep(upkeepId)
        assert.equal(
          registration.maxValidBlocknumber.toNumber(),
          receipt.blockNumber + 50,
        )
      })

      it('emits an event', async () => {
        const tx = await registry.connect(admin).cancelUpkeep(upkeepId)
        const receipt = await tx.wait()
        await expect(tx)
          .to.emit(registry, 'UpkeepCanceled')
          .withArgs(
            upkeepId,
            BigNumber.from(receipt.blockNumber + cancellationDelay),
          )
      })

      it('immediately prevents upkeep', async () => {
        await linkToken.connect(owner).approve(registry.address, toWei('100'))
        await registry.connect(owner).addFunds(upkeepId, toWei('100'))
        await registry.connect(admin).cancelUpkeep(upkeepId)
        registry
          .connect(keeper1)
          .transmit(
            [emptyBytes32, emptyBytes32, emptyBytes32],
            await encodeLatestBlockReport([{ Id: upkeepId.toString() }]),
            [],
            [],
            emptyBytes32,
          ) // still works

        for (let i = 0; i < cancellationDelay; i++) {
          await ethers.provider.send('evm_mine', [])
        }

        await evmRevert(
          registry
            .connect(keeper1)
            .transmit(
              [emptyBytes32, emptyBytes32, emptyBytes32],
              await encodeLatestBlockReport([{ Id: upkeepId.toString() }]),
              [],
              [],
              emptyBytes32,
            ),
          'StaleReport()',
        )
      })

      describe('when an upkeep has been performed', async () => {
        beforeEach(async () => {
          await linkToken.connect(owner).approve(registry.address, toWei('100'))
          await registry.connect(owner).addFunds(upkeepId, toWei('100'))
          registry
            .connect(keeper1)
            .transmit(
              [emptyBytes32, emptyBytes32, emptyBytes32],
              await encodeLatestBlockReport([{ Id: upkeepId.toString() }]),
              [],
              [],
              emptyBytes32,
            )
        })

        it('deducts a cancellation fee from the upkeep and gives to owner', async () => {
          let minUpkeepSpend = toWei('10')

          await registry.connect(owner).setConfig(
            signerAddresses,
            keeperAddresses,
            f,
            encodeConfig({
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
            }),
            offchainVersion,
            offchainBytes,
          )

          const payee1Before = await linkToken.balanceOf(
            await payee1.getAddress(),
          )
          let upkeepBefore = (await registry.getUpkeep(upkeepId)).balance
          let ownerBefore = (await registry.getState()).state.ownerLinkBalance
          assert.equal(0, ownerBefore.toNumber())

          let amountSpent = toWei('100').sub(upkeepBefore)
          let cancellationFee = minUpkeepSpend.sub(amountSpent)

          await registry.connect(admin).cancelUpkeep(upkeepId)

          const payee1After = await linkToken.balanceOf(
            await payee1.getAddress(),
          )
          let upkeepAfter = (await registry.getUpkeep(upkeepId)).balance
          let ownerAfter = (await registry.getState()).state.ownerLinkBalance

          // post upkeep balance should be previous balance minus cancellation fee
          assert.isTrue(upkeepBefore.sub(cancellationFee).eq(upkeepAfter))
          // payee balance should not change
          assert.isTrue(payee1Before.eq(payee1After))
          // owner should receive the cancellation fee
          assert.isTrue(ownerAfter.eq(cancellationFee))
        })

        it('deducts up to balance as cancellation fee', async () => {
          // Very high min spend, should deduct whole balance as cancellation fees
          let minUpkeepSpend = toWei('1000')
          await registry.connect(owner).setConfig(
            signerAddresses,
            keeperAddresses,
            f,
            encodeConfig({
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
            }),
            offchainVersion,
            offchainBytes,
          )
          const payee1Before = await linkToken.balanceOf(
            await payee1.getAddress(),
          )
          let upkeepBefore = (await registry.getUpkeep(upkeepId)).balance
          let ownerBefore = (await registry.getState()).state.ownerLinkBalance
          assert.equal(0, ownerBefore.toNumber())

          await registry.connect(admin).cancelUpkeep(upkeepId)
          const payee1After = await linkToken.balanceOf(
            await payee1.getAddress(),
          )
          let ownerAfter = (await registry.getState()).state.ownerLinkBalance
          let upkeepAfter = (await registry.getUpkeep(upkeepId)).balance

          // all upkeep balance is deducted for cancellation fee
          assert.equal(0, upkeepAfter.toNumber())
          // payee balance should not change
          assert.isTrue(payee1After.eq(payee1Before))
          // all upkeep balance is transferred to the owner
          assert.isTrue(ownerAfter.eq(upkeepBefore))
        })

        it('does not deduct cancellation fee if more than minUpkeepSpend is spent', async () => {
          // Very low min spend, already spent in one perform upkeep
          let minUpkeepSpend = BigNumber.from(420)
          await registry.connect(owner).setConfig(
            signerAddresses,
            keeperAddresses,
            f,
            encodeConfig({
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
            }),
            offchainVersion,
            offchainBytes,
          )
          const payee1Before = await linkToken.balanceOf(
            await payee1.getAddress(),
          )
          let upkeepBefore = (await registry.getUpkeep(upkeepId)).balance
          let ownerBefore = (await registry.getState()).state.ownerLinkBalance
          assert.equal(0, ownerBefore.toNumber())

          await registry.connect(admin).cancelUpkeep(upkeepId)
          const payee1After = await linkToken.balanceOf(
            await payee1.getAddress(),
          )
          let ownerAfter = (await registry.getState()).state.ownerLinkBalance
          let upkeepAfter = (await registry.getUpkeep(upkeepId)).balance

          // upkeep does not pay cancellation fee after cancellation because minimum upkeep spent is met
          assert.isTrue(upkeepBefore.eq(upkeepAfter))
          // owner balance does not change
          assert.equal(0, ownerAfter.toNumber())
          // payee balance does not change
          assert.isTrue(payee1Before.eq(payee1After))
        })
      })
    })
  })

  describe('#withdrawPayment', () => {
    beforeEach(async () => {
      await linkToken.connect(owner).approve(registry.address, toWei('100'))
      await registry.connect(owner).addFunds(upkeepId, toWei('100'))
      registry
        .connect(keeper1)
        .transmit(
          [emptyBytes32, emptyBytes32, emptyBytes32],
          await encodeLatestBlockReport([{ Id: upkeepId.toString() }]),
          [],
          [],
          emptyBytes32,
        )
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
        await registry.getTransmitterInfo(await keeper1.getAddress())
      ).balance
      const registrationBefore = (await registry.getUpkeep(upkeepId)).balance
      const toLinkBefore = await linkToken.balanceOf(to)
      const registryLinkBefore = await linkToken.balanceOf(registry.address)

      //// Do the thing
      await registry
        .connect(payee1)
        .withdrawPayment(await keeper1.getAddress(), to)

      const keeperAfter = (
        await registry.getTransmitterInfo(await keeper1.getAddress())
      ).balance
      const registrationAfter = (await registry.getUpkeep(upkeepId)).balance
      const toLinkAfter = await linkToken.balanceOf(to)
      const registryLinkAfter = await linkToken.balanceOf(registry.address)

      assert.isTrue(keeperAfter.eq(BigNumber.from(0)))
      assert.isTrue(registrationBefore.eq(registrationAfter))
      assert.isTrue(toLinkBefore.add(keeperBefore).eq(toLinkAfter))
      assert.isTrue(registryLinkBefore.sub(keeperBefore).eq(registryLinkAfter))
    })

    it('emits a log announcing the withdrawal', async () => {
      const balance = (
        await registry.getTransmitterInfo(await keeper1.getAddress())
      ).balance
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
})
