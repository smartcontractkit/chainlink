import { ethers } from 'hardhat'
import { assert, expect } from 'chai'
import { BigNumber, Signer, Wallet } from 'ethers'
import { evmRevert } from '../../test-helpers/matchers'
import { getUsers, Personas } from '../../test-helpers/setup'
import { toWei } from '../../test-helpers/helpers'
import { LinkToken__factory as LinkTokenFactory } from '../../../typechain/factories/LinkToken__factory'
import { MockV3Aggregator__factory as MockV3AggregatorFactory } from '../../../typechain/factories/MockV3Aggregator__factory'
import { UpkeepMock__factory as UpkeepMockFactory } from '../../../typechain/factories/UpkeepMock__factory'
import { UpkeepAutoFunder__factory as UpkeepAutoFunderFactory } from '../../../typechain/factories/UpkeepAutoFunder__factory'
import { UpkeepTranscoder__factory as UpkeepTranscoderFactory } from '../../../typechain/factories/UpkeepTranscoder__factory'
import { KeeperRegistry20__factory as KeeperRegistryFactory } from '../../../typechain/factories/KeeperRegistry20__factory'
import { MockArbGasInfo__factory as MockArbGasInfoFactory } from '../../../typechain/factories/MockArbGasInfo__factory'
import { MockOVMGasPriceOracle__factory as MockOVMGasPriceOracleFactory } from '../../../typechain/factories/MockOVMGasPriceOracle__factory'
import { KeeperRegistryLogic20__factory as KeeperRegistryLogicFactory } from '../../../typechain/factories/KeeperRegistryLogic20__factory'
import { MockArbSys__factory as MockArbSysFactory } from '../../../typechain/factories/MockArbSys__factory'
import { KeeperRegistry20 as KeeperRegistry } from '../../../typechain/KeeperRegistry20'
import { KeeperRegistryLogic20 as KeeperRegistryLogic } from '../../../typechain/KeeperRegistryLogic20'
import { MockV3Aggregator } from '../../../typechain/MockV3Aggregator'
import { LinkToken } from '../../../typechain/LinkToken'
import { UpkeepMock } from '../../../typechain/UpkeepMock'
import { MockArbGasInfo } from '../../../typechain/MockArbGasInfo'
import { MockOVMGasPriceOracle } from '../../../typechain/MockOVMGasPriceOracle'
import { UpkeepTranscoder } from '../../../typechain/UpkeepTranscoder'

// copied from AutomationRegistryInterface2_0.sol
enum UpkeepFailureReason {
  NONE,
  UPKEEP_CANCELLED,
  UPKEEP_PAUSED,
  TARGET_CHECK_REVERTED,
  UPKEEP_NOT_NEEDED,
  PERFORM_DATA_EXCEEDS_LIMIT,
  INSUFFICIENT_BALANCE,
}

// copied from AutomationRegistryInterface2_0.sol
enum Mode {
  DEFAULT,
  ARBITRUM,
  OPTIMISM,
}

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
const transmitGasOverhead = BigNumber.from(800000)
const checkGasOverhead = BigNumber.from(400000)

// These values should match the constants declared in registry
const registryGasOverhead = BigNumber.from(70_000)
const registryPerSignerGasOverhead = BigNumber.from(7500)
const registryPerPerformByteGasOverhead = BigNumber.from(20)
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

const linkEth = BigNumber.from(5000000000000000) // 1 Link = 0.005 Eth
const gasWei = BigNumber.from(1000000000) // 1 gwei
const encodeReport = (
  upkeeps: any,
  gasWeiReport = gasWei,
  linkEthReport = linkEth,
) => {
  const upkeepIds = upkeeps.map((u: any) => u.Id)
  const performDataTuples = upkeeps.map((u: any) => [
    u.checkBlockNum,
    u.checkBlockHash,
    u.performData,
  ])
  return ethers.utils.defaultAbiCoder.encode(
    ['uint256', 'uint256', 'uint256[]', 'tuple(uint32,bytes32,bytes)[]'],
    [gasWeiReport, linkEthReport, upkeepIds, performDataTuples],
  )
}

const encodeLatestBlockReport = async (upkeeps: any) => {
  const latestBlock = await ethers.provider.getBlock('latest')
  for (let i = 0; i < upkeeps.length; i++) {
    upkeeps[i].checkBlockNum = latestBlock.number
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

  const signatures = []
  for (const signer of signers) {
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
  const upkeepPerformedABI = [
    'event UpkeepPerformed(uint256 indexed id,bool indexed success, \
  uint32 checkBlockNumber,uint256 gasUsed,uint256 gasOverhead,uint96 totalPayment)',
  ]
  const iface = new ethers.utils.Interface(upkeepPerformedABI)

  const parsedLogs = []
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

const parseReorgedUpkeepReportLogs = (receipt: any) => {
  const logABI = ['  event ReorgedUpkeepReport(uint256 indexed id)']
  const iface = new ethers.utils.Interface(logABI)

  const parsedLogs = []
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

const parseStaleUpkeepReportLogs = (receipt: any) => {
  const logABI = ['  event StaleUpkeepReport(uint256 indexed id)']
  const iface = new ethers.utils.Interface(logABI)

  const parsedLogs = []
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

const parseInsufficientFundsUpkeepReportLogs = (receipt: any) => {
  const logABI = ['  event InsufficientFundsUpkeepReport(uint256 indexed id)']
  const iface = new ethers.utils.Interface(logABI)

  const parsedLogs = []
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

const parseCancelledUpkeepReportLogs = (receipt: any) => {
  const logABI = ['  event CancelledUpkeepReport(uint256 indexed id)']
  const iface = new ethers.utils.Interface(logABI)

  const parsedLogs = []
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
  keeperRegistryFactory = (await ethers.getContractFactory(
    'KeeperRegistry2_0',
  )) as unknown as KeeperRegistryFactory // bug in typechain requires force casting
  keeperRegistryLogicFactory = (await ethers.getContractFactory(
    'KeeperRegistryLogic2_0',
  )) as unknown as KeeperRegistryLogicFactory // bug in typechain requires force casting
  upkeepMockFactory = await ethers.getContractFactory('UpkeepMock')
  upkeepAutoFunderFactory = await ethers.getContractFactory('UpkeepAutoFunder')
  upkeepTranscoderFactory = await ethers.getContractFactory('UpkeepTranscoder')
  mockArbGasInfoFactory = await ethers.getContractFactory('MockArbGasInfo')
  mockOVMGasPriceOracleFactory = await ethers.getContractFactory(
    'MockOVMGasPriceOracle',
  )
})

describe('KeeperRegistry2_0', () => {
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
    const gasPayment = base.add(l1Fee)

    const premium = gasWei
      .mul(gasMultiplier)
      .mul(upkeepGasSpent)
      .add(l1CostWei.mul(gasMultiplier).div(numUpkeepsBatch))
      .mul(linkDivisibility)
      .div(linkEth)
      .mul(premiumPPB)
      .div(paymentPremiumBase)
      .add(BigNumber.from(flatFee).mul('1000000000000'))

    return {
      total: gasPayment.add(premium),
      gasPaymemnt: gasPayment,
      premium,
    }
  }

  const verifyMaxPayment = async (
    mode: number,
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
    const registryLogic = await keeperRegistryLogicFactory
      .connect(owner)
      .deploy(
        mode,
        linkToken.address,
        linkEthFeed.address,
        gasPriceFeed.address,
      )
    // Deploy a new registry since we change payment model
    const registry = await keeperRegistryFactory
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

    const fPlusOne = BigNumber.from(f + 1)
    const totalGasOverhead = registryGasOverhead
      .add(registryPerSignerGasOverhead.mul(fPlusOne))
      .add(registryPerPerformByteGasOverhead.mul(maxPerformDataSize))

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
              ).total,
            )
          }
        }
      }
    }
  }

  const getTransmitTx = async (
    registry: KeeperRegistry,
    transmitter: any,
    upkeepIds: any,
    numSigners: any,
    extraParams?: any,
    performData?: any,
    checkBlockNum?: any,
    checkBlockHash?: any,
  ) => {
    const latestBlock = await ethers.provider.getBlock('latest')
    const configDigest = (await registry.getState()).state.latestConfigDigest

    const upkeeps = []
    for (let i = 0; i < upkeepIds.length; i++) {
      upkeeps.push({
        Id: upkeepIds[i],
        checkBlockNum: checkBlockNum ? checkBlockNum : latestBlock.number,
        checkBlockHash: checkBlockHash ? checkBlockHash : latestBlock.hash,
        performData: performData ? performData : '0x',
      })
    }

    const report = encodeReport(upkeeps)
    const reportContext = [configDigest, epochAndRound5_1, emptyBytes32]
    const sigs = signReport(reportContext, report, signers.slice(0, numSigners))

    return registry
      .connect(transmitter)
      .transmit(
        [configDigest, epochAndRound5_1, emptyBytes32],
        report,
        sigs.rs,
        sigs.ss,
        sigs.vs,
        { gasLimit: extraParams?.gasLimit, gasPrice: extraParams?.gasPrice },
      )
  }

  const getTransmitTxWithReport = async (
    registry: KeeperRegistry,
    transmitter: any,
    report: any,
    numSigners: any,
  ) => {
    const configDigest = (await registry.getState()).state.latestConfigDigest
    const reportContext = [configDigest, epochAndRound5_1, emptyBytes32]
    const sigs = signReport(reportContext, report, signers.slice(0, numSigners))

    return registry
      .connect(transmitter)
      .transmit(
        [configDigest, epochAndRound5_1, emptyBytes32],
        report,
        sigs.rs,
        sigs.ss,
        sigs.vs,
      )
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
      signers.push(ethers.Wallet.createRandom())
    }
    signerAddresses = []
    for (const signer of signers) {
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

    const mockArbSys = await new MockArbSysFactory(owner).deploy()
    const arbSysCode = await ethers.provider.send('eth_getCode', [
      mockArbSys.address,
    ])
    await ethers.provider.send('hardhat_setCode', [
      '0x0000000000000000000000000000000000000064',
      arbSysCode,
    ])

    registryLogic = await keeperRegistryLogicFactory
      .connect(owner)
      .deploy(
        Mode.DEFAULT,
        linkToken.address,
        linkEthFeed.address,
        gasPriceFeed.address,
      )

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
        randomBytes,
        emptyBytes,
      )
    upkeepId = await getUpkeepID(tx)
  })

  describe('#transmit', () => {
    const fArray = [1, 5, 10]

    it('reverts when registry is paused', async () => {
      await registry.connect(owner).pause()
      await evmRevert(
        getTransmitTx(registry, keeper1, [upkeepId.toString()], f + 1),
        'RegistryPaused()',
      )
    })

    it('reverts when called by non active transmitter', async () => {
      await evmRevert(
        getTransmitTx(registry, payee1, [upkeepId.toString()], f + 1),
        'OnlyActiveTransmitters()',
      )
    })

    it('reverts when upkeeps and performData length mismatches', async () => {
      const upkeepIds = []
      const performDataTuples = []
      const latestBlock = await ethers.provider.getBlock('latest')

      upkeepIds.push(upkeepId)
      performDataTuples.push([latestBlock.number + 1, latestBlock.hash, '0x'])
      // Push an extra perform data
      performDataTuples.push([latestBlock.number + 1, latestBlock.hash, '0x'])

      const report = ethers.utils.defaultAbiCoder.encode(
        ['uint256', 'uint256', 'uint256[]', 'tuple(uint32,bytes32,bytes)[]'],
        [0, 0, upkeepIds, performDataTuples],
      )

      await evmRevert(
        getTransmitTxWithReport(registry, keeper1, report, f + 1),
        'InvalidReport()',
      )
    })

    it('reverts when wrappedPerformData is incorrectly encoded', async () => {
      const upkeepIds = []
      const wrappedPerformDatas = []
      const latestBlock = await ethers.provider.getBlock('latest')

      upkeepIds.push(upkeepId)
      wrappedPerformDatas.push(
        ethers.utils.defaultAbiCoder.encode(
          ['tuple(uint32,bytes32)'], // missing performData
          [[latestBlock.number + 1, latestBlock.hash]],
        ),
      )

      const report = ethers.utils.defaultAbiCoder.encode(
        ['uint256[]', 'bytes[]'],
        [upkeepIds, wrappedPerformDatas],
      )

      await evmRevert(getTransmitTxWithReport(registry, keeper1, report, f + 1))
    })

    it('returns early when no upkeeps are included in report', async () => {
      const upkeepIds: string[] = []
      const wrappedPerformDatas: string[] = []
      const report = ethers.utils.defaultAbiCoder.encode(
        ['uint256', 'uint256', 'uint256[]', 'bytes[]'],
        [0, 0, upkeepIds, wrappedPerformDatas],
      )

      await getTransmitTxWithReport(registry, keeper1, report, f + 1)
    })

    it('returns early when invalid upkeepIds are included in report', async () => {
      const tx = await getTransmitTx(
        registry,
        keeper1,
        [upkeepId.add(BigNumber.from('1')).toString()],
        f + 1,
      )

      const receipt = await tx.wait()
      const cancelledUpkeepReportLogs = parseCancelledUpkeepReportLogs(receipt)
      // exactly 1 CancelledUpkeepReport log should be emitted
      assert.equal(cancelledUpkeepReportLogs.length, 1)
    })

    it('reverts when duplicated upkeepIds are included in report', async () => {
      // Fund the upkeep so that pre-checks pass
      await registry.connect(admin).addFunds(upkeepId, toWei('100'))
      await evmRevert(
        getTransmitTx(
          registry,
          keeper1,
          [upkeepId.toString(), upkeepId.toString()],
          f + 1,
        ),
        'InvalidReport()',
      )
    })

    it('returns early when upkeep has insufficient funds', async () => {
      const tx = await getTransmitTx(
        registry,
        keeper1,
        [upkeepId.toString()],
        f + 1,
      )

      const receipt = await tx.wait()
      const insufficientFundsUpkeepReportLogs =
        parseInsufficientFundsUpkeepReportLogs(receipt)
      // exactly 1 InsufficientFundsUpkeepReportLogs log should be emitted
      assert.equal(insufficientFundsUpkeepReportLogs.length, 1)
    })

    context('When the upkeep is funded', async () => {
      beforeEach(async () => {
        // Fund the upkeep
        await registry.connect(admin).addFunds(upkeepId, toWei('100'))
      })

      it('returns early when check block number is less than last perform', async () => {
        // First perform an upkeep to put last perform block number on upkeep state

        const tx = await getTransmitTx(
          registry,
          keeper1,
          [upkeepId.toString()],
          f + 1,
        )
        await tx.wait()

        const lastPerformBlockNumber = (await registry.getUpkeep(upkeepId))
          .lastPerformBlockNumber
        const lastPerformBlock = await ethers.provider.getBlock(
          lastPerformBlockNumber,
        )
        assert.equal(
          lastPerformBlockNumber.toString(),
          tx.blockNumber?.toString(),
        )

        // Try to transmit a report which has checkBlockNumber = lastPerformBlockNumber-1, should result in stale report
        const transmitTx = await getTransmitTx(
          registry,
          keeper1,
          [upkeepId.toString()],
          f + 1,
          {},
          '0x',
          lastPerformBlock.number - 1,
          lastPerformBlock.parentHash,
        )

        const receipt = await transmitTx.wait()
        const staleUpkeepReportLogs = parseStaleUpkeepReportLogs(receipt)
        // exactly 1 StaleUpkeepReportLogs log should be emitted
        assert.equal(staleUpkeepReportLogs.length, 1)
      })

      it('returns early when check block hash does not match', async () => {
        await registry.connect(admin).addFunds(upkeepId, toWei('100'))
        const latestBlock = await ethers.provider.getBlock('latest')
        // Try to transmit a report which has incorrect checkBlockHash
        const tx = await getTransmitTx(
          registry,
          keeper1,
          [upkeepId.toString()],
          f + 1,
          {},
          '0x',
          latestBlock.number - 1,
          latestBlock.hash,
        ) // should be latestBlock.parentHash

        const receipt = await tx.wait()
        const reorgedUpkeepReportLogs = parseReorgedUpkeepReportLogs(receipt)
        // exactly 1 ReorgedUpkeepReportLogs log should be emitted
        assert.equal(reorgedUpkeepReportLogs.length, 1)
      })

      it('returns early when check block number is older than 256 blocks', async () => {
        const latestBlockReport = await encodeLatestBlockReport([
          { Id: upkeepId.toString() },
        ])

        for (let i = 0; i < 256; i++) {
          await ethers.provider.send('evm_mine', [])
        }

        // Try to transmit a report which is older than 256 blocks so block hash cannot be matched
        const tx = await registry
          .connect(keeper1)
          .transmit(
            [emptyBytes32, emptyBytes32, emptyBytes32],
            latestBlockReport,
            [],
            [],
            emptyBytes32,
          )

        const receipt = await tx.wait()
        const reorgedUpkeepReportLogs = parseReorgedUpkeepReportLogs(receipt)
        // exactly 1 ReorgedUpkeepReportLogs log should be emitted
        assert.equal(reorgedUpkeepReportLogs.length, 1)
      })

      it('returns early when upkeep is cancelled and cancellation delay has gone', async () => {
        const latestBlockReport = await encodeLatestBlockReport([
          { Id: upkeepId.toString() },
        ])
        await registry.connect(admin).cancelUpkeep(upkeepId)

        for (let i = 0; i < cancellationDelay; i++) {
          await ethers.provider.send('evm_mine', [])
        }

        const tx = await getTransmitTxWithReport(
          registry,
          keeper1,
          latestBlockReport,
          f + 1,
        )

        const receipt = await tx.wait()
        const cancelledUpkeepReportLogs =
          parseCancelledUpkeepReportLogs(receipt)
        // exactly 1 CancelledUpkeepReport log should be emitted
        assert.equal(cancelledUpkeepReportLogs.length, 1)
      })

      it('does not revert if the target cannot execute', async () => {
        mock.setCanPerform(false)
        const tx = await getTransmitTx(
          registry,
          keeper1,
          [upkeepId.toString()],
          f + 1,
        )

        const receipt = await tx.wait()
        const upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
        // exactly 1 Upkeep Performed should be emitted
        assert.equal(upkeepPerformedLogs.length, 1)
        const upkeepPerformedLog = upkeepPerformedLogs[0]

        const success = upkeepPerformedLog.args.success
        assert.equal(success, false)
      })

      it('reverts if not enough gas supplied', async () => {
        mock.setPerformGasToBurn(executeGas)
        await evmRevert(
          getTransmitTx(registry, keeper1, [upkeepId.toString()], f + 1, {
            gasLimit: executeGas,
          }),
        )
      })

      it('executes the data passed to the registry', async () => {
        mock.setCanPerform(true)

        const tx = await getTransmitTx(
          registry,
          keeper1,
          [upkeepId.toString()],
          f + 1,
          {},
          randomBytes,
        )
        const receipt = await tx.wait()

        const upkeepPerformedWithABI = [
          'event UpkeepPerformedWith(bytes upkeepData)',
        ]
        const iface = new ethers.utils.Interface(upkeepPerformedWithABI)
        const parsedLogs = []
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

      it('uses actual execution price for payment and premium calculation', async () => {
        // Actual multiplier is 2, but we set gasPrice to be 1x gasWei
        const gasPrice = gasWei.mul(BigNumber.from('1'))
        mock.setCanPerform(true)
        const registryPremiumBefore = (await registry.getState()).state
          .totalPremium
        const tx = await getTransmitTx(
          registry,
          keeper1,
          [upkeepId.toString()],
          f + 1,
          { gasPrice },
        )
        const receipt = await tx.wait()
        const registryPremiumAfter = (await registry.getState()).state
          .totalPremium
        const premium = registryPremiumAfter.sub(registryPremiumBefore)

        const upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
        // exactly 1 Upkeep Performed should be emitted
        assert.equal(upkeepPerformedLogs.length, 1)
        const upkeepPerformedLog = upkeepPerformedLogs[0]

        const gasUsed = upkeepPerformedLog.args.gasUsed
        const gasOverhead = upkeepPerformedLog.args.gasOverhead
        const totalPayment = upkeepPerformedLog.args.totalPayment

        assert.equal(
          linkForGas(
            gasUsed,
            gasOverhead,
            BigNumber.from('1'), // Not the config multiplier, but the actual gas used
            paymentPremiumPPB,
            flatFeeMicroLink,
          ).total.toString(),
          totalPayment.toString(),
        )

        assert.equal(
          linkForGas(
            gasUsed,
            gasOverhead,
            BigNumber.from('1'), // Not the config multiplier, but the actual gas used
            paymentPremiumPPB,
            flatFeeMicroLink,
          ).premium.toString(),
          premium.toString(),
        )
      })

      it('only pays at a rate up to the gas ceiling [ @skip-coverage ]', async () => {
        // Actual multiplier is 2, but we set gasPrice to be 10x
        const gasPrice = gasWei.mul(BigNumber.from('10'))
        mock.setCanPerform(true)

        const tx = await getTransmitTx(
          registry,
          keeper1,
          [upkeepId.toString()],
          f + 1,
          { gasPrice },
        )
        const receipt = await tx.wait()
        const upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
        // exactly 1 Upkeep Performed should be emitted
        assert.equal(upkeepPerformedLogs.length, 1)
        const upkeepPerformedLog = upkeepPerformedLogs[0]

        const gasUsed = upkeepPerformedLog.args.gasUsed
        const gasOverhead = upkeepPerformedLog.args.gasOverhead
        const totalPayment = upkeepPerformedLog.args.totalPayment

        assert.equal(
          linkForGas(
            gasUsed,
            gasOverhead,
            gasCeilingMultiplier, // Should be same with exisitng multiplier
            paymentPremiumPPB,
            flatFeeMicroLink,
          ).total.toString(),
          totalPayment.toString(),
        )
      })

      it('correctly accounts for l1 payment', async () => {
        mock.setCanPerform(true)
        // Same as MockArbGasInfo.sol
        const l1CostWeiArb = BigNumber.from(1000000)

        // Deploy a new registry since we change payment model
        const registryLogic = await keeperRegistryLogicFactory
          .connect(owner)
          .deploy(
            Mode.ARBITRUM,
            linkToken.address,
            linkEthFeed.address,
            gasPriceFeed.address,
          )
        // Deploy a new registry since we change payment model
        const registry = await keeperRegistryFactory
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
        let tx = await registry
          .connect(owner)
          .registerUpkeep(
            mock.address,
            executeGas,
            await admin.getAddress(),
            randomBytes,
            emptyBytes,
          )
        upkeepId = await getUpkeepID(tx)
        await linkToken.connect(owner).approve(registry.address, toWei('1000'))
        await registry.connect(owner).addFunds(upkeepId, toWei('100'))

        // Do the thing
        tx = await getTransmitTx(
          registry,
          keeper1,
          [upkeepId.toString()],
          f + 1,
          { gasPrice: gasWei.mul('5') }, // High gas price so that it gets capped
        )
        const receipt = await tx.wait()
        const upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
        // exactly 1 Upkeep Performed should be emitted
        assert.equal(upkeepPerformedLogs.length, 1)
        const upkeepPerformedLog = upkeepPerformedLogs[0]

        const gasUsed = upkeepPerformedLog.args.gasUsed
        const gasOverhead = upkeepPerformedLog.args.gasOverhead
        const totalPayment = upkeepPerformedLog.args.totalPayment

        assert.equal(
          linkForGas(
            gasUsed,
            gasOverhead,
            gasCeilingMultiplier,
            paymentPremiumPPB,
            flatFeeMicroLink,
            l1CostWeiArb.div(gasCeilingMultiplier), // Dividing by gasCeilingMultiplier as it gets multiplied later
          ).total.toString(),
          totalPayment.toString(),
        )
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
            randomBytes,
            emptyBytes,
          )
        upkeepId = await getUpkeepID(tx)

        await autoFunderUpkeep.setUpkeepId(upkeepId)
        // Give enough funds for upkeep as well as to the upkeep contract
        await linkToken
          .connect(owner)
          .transfer(autoFunderUpkeep.address, toWei('1000'))
        const maxPayment = await registry.getMaxPaymentForGas(executeGas)

        // First set auto funding amount to 0 and verify that balance is deducted upon performUpkeep
        let initialBalance = toWei('100')
        await registry.connect(owner).addFunds(upkeepId, initialBalance)
        await autoFunderUpkeep.setAutoFundLink(0)
        await autoFunderUpkeep.setIsEligible(true)
        await getTransmitTx(registry, keeper1, [upkeepId.toString()], f + 1)

        let postUpkeepBalance = (await registry.getUpkeep(upkeepId)).balance
        assert.isTrue(postUpkeepBalance.lt(initialBalance)) // Balance should be deducted
        assert.isTrue(postUpkeepBalance.gte(initialBalance.sub(maxPayment))) // Balance should not be deducted more than maxPayment

        // Now set auto funding amount to 100 wei and verify that the balance increases
        initialBalance = postUpkeepBalance
        const autoTopupAmount = toWei('100')
        await autoFunderUpkeep.setAutoFundLink(autoTopupAmount)
        await autoFunderUpkeep.setIsEligible(true)
        await getTransmitTx(registry, keeper1, [upkeepId.toString()], f + 1)

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
        const tx = await registry
          .connect(owner)
          .registerUpkeep(
            autoFunderUpkeep.address,
            executeGas,
            autoFunderUpkeep.address,
            randomBytes,
            emptyBytes,
          )
        upkeepId = await getUpkeepID(tx)

        await autoFunderUpkeep.setUpkeepId(upkeepId)
        await registry.connect(owner).addFunds(upkeepId, toWei('100'))

        await autoFunderUpkeep.setIsEligible(true)
        await autoFunderUpkeep.setShouldCancel(true)

        let registration = await registry.getUpkeep(upkeepId)
        const oldExpiration = registration.maxValidBlocknumber

        // Do the thing
        await getTransmitTx(registry, keeper1, [upkeepId.toString()], f + 1)

        // Verify upkeep gets cancelled
        registration = await registry.getUpkeep(upkeepId)
        const newExpiration = registration.maxValidBlocknumber
        assert.isTrue(newExpiration.lt(oldExpiration))
      })

      it('reverts when configDigest mismatches', async () => {
        const report = await encodeLatestBlockReport([
          {
            Id: upkeepId.toString(),
          },
        ])
        const reportContext = [emptyBytes32, epochAndRound5_1, emptyBytes32] // wrong config digest
        const sigs = signReport(reportContext, report, signers.slice(0, f + 1))
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
            Id: upkeepId.toString(),
          },
        ])
        const reportContext = [configDigest, epochAndRound5_1, emptyBytes32] // wrong config digest
        const sigs = signReport(reportContext, report, signers.slice(0, f + 2))
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
            Id: upkeepId.toString(),
          },
        ])
        const reportContext = [configDigest, epochAndRound5_1, emptyBytes32] // wrong config digest
        const sigs = signReport(reportContext, report, [
          new ethers.Wallet(ethers.Wallet.createRandom()),
          new ethers.Wallet(ethers.Wallet.createRandom()),
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
            Id: upkeepId.toString(),
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

      it('has a large enough gas overhead to cover upkeep that use all its gas [ @skip-coverage ]', async () => {
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
          randomBytes,
          emptyBytes,
        )
        upkeepId = await getUpkeepID(tx)
        await registry.connect(admin).addFunds(upkeepId, toWei('100'))

        let performData = '0x'
        for (let i = 0; i < maxPerformDataSize.toNumber(); i++) {
          performData += '11'
        } // max allowed performData

        mock.setCanPerform(true)
        mock.setPerformGasToBurn(maxPerformGas)

        await getTransmitTx(
          registry,
          keeper1,
          [upkeepId.toString()],
          11,
          { gasLimit: maxPerformGas.add(transmitGasOverhead) },
          performData,
        ) // Should not revert
      })

      it('performs upkeep, deducts payment, updates lastPerformBlockNumber and emits events', async () => {
        for (const i in fArray) {
          const newF = fArray[i]
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
          const checkBlock = await ethers.provider.getBlock('latest')

          const keeperBefore = await registry.getTransmitterInfo(
            await keeper1.getAddress(),
          )
          const registrationBefore = await registry.getUpkeep(upkeepId)
          const registryPremiumBefore = (await registry.getState()).state
            .totalPremium
          const keeperLinkBefore = await linkToken.balanceOf(
            await keeper1.getAddress(),
          )
          const registryLinkBefore = await linkToken.balanceOf(registry.address)

          // Do the thing
          const tx = await getTransmitTx(
            registry,
            keeper1,
            [upkeepId.toString()],
            newF + 1,
            {},
            '0x',
            checkBlock.number - 1,
            checkBlock.parentHash,
          )

          const receipt = await tx.wait()

          const upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
          // exactly 1 Upkeep Performed should be emitted
          assert.equal(upkeepPerformedLogs.length, 1)
          const upkeepPerformedLog = upkeepPerformedLogs[0]

          const id = upkeepPerformedLog.args.id
          const success = upkeepPerformedLog.args.success
          const checkBlockNumber = upkeepPerformedLog.args.checkBlockNumber
          const gasUsed = upkeepPerformedLog.args.gasUsed
          const gasOverhead = upkeepPerformedLog.args.gasOverhead
          const totalPayment = upkeepPerformedLog.args.totalPayment

          assert.equal(id.toString(), upkeepId.toString())
          assert.equal(success, true)
          assert.equal(
            checkBlockNumber.toString(),
            (checkBlock.number - 1).toString(),
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
          const registryPremiumAfter = (await registry.getState()).state
            .totalPremium
          const premium = registryPremiumAfter.sub(registryPremiumBefore)
          // Keeper payment is gasPayment + premium / num keepers
          const keeperPayment = totalPayment
            .sub(premium)
            .add(premium.div(BigNumber.from(keeperAddresses.length)))

          assert.equal(
            keeperAfter.balance.sub(keeperPayment).toString(),
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

          // Latest epoch should be 5
          assert.equal((await registry.getState()).state.latestEpoch, 5)
        }
      })

      it('calculates gas overhead appropriately within a margin for different scenarios [ @skip-coverage ]', async () => {
        // Perform the upkeep once to remove non-zero storage slots and have predictable gas measurement

        let tx = await getTransmitTx(
          registry,
          keeper1,
          [upkeepId.toString()],
          f + 1,
        )

        await tx.wait()

        // Different test scenarios
        let longBytes = '0x'
        for (let i = 0; i < maxPerformDataSize.toNumber(); i++) {
          longBytes += '11'
        }
        const upkeepSuccessArray = [true, false]
        const performGasArray = [5000, 100000, executeGas]
        const performDataArray = ['0x', randomBytes, longBytes]

        for (const i in upkeepSuccessArray) {
          for (const j in performGasArray) {
            for (const k in performDataArray) {
              for (const l in fArray) {
                const upkeepSuccess = upkeepSuccessArray[i]
                const performGas = performGasArray[j]
                const performData = performDataArray[k]
                const newF = fArray[l]

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
                tx = await getTransmitTx(
                  registry,
                  keeper1,
                  [upkeepId.toString()],
                  newF + 1,
                  {},
                  performData,
                )
                const receipt = await tx.wait()
                const upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
                // exactly 1 Upkeep Performed should be emitted
                assert.equal(upkeepPerformedLogs.length, 1)
                const upkeepPerformedLog = upkeepPerformedLogs[0]

                const upkeepGasUsed = upkeepPerformedLog.args.gasUsed
                const chargedGasOverhead = upkeepPerformedLog.args.gasOverhead
                const actualGasOverhead = receipt.gasUsed.sub(upkeepGasUsed)

                assert.isTrue(upkeepGasUsed.gt(BigNumber.from('0')))
                assert.isTrue(chargedGasOverhead.gt(BigNumber.from('0')))

                if (i == '0' && j == '0' && k == '0') {
                  console.log(
                    'Gas Benchmarking - sig verification ( f =',
                    newF,
                    '): calculated overhead: ',
                    chargedGasOverhead.toString(),
                    ' actual overhead: ',
                    actualGasOverhead.toString(),
                    ' margin over gasUsed: ',
                    chargedGasOverhead.sub(actualGasOverhead).toString(),
                  )
                }

                // Overhead should not get capped
                const gasOverheadCap = registryGasOverhead
                  .add(
                    registryPerSignerGasOverhead.mul(BigNumber.from(newF + 1)),
                  )
                  .add(
                    BigNumber.from(
                      registryPerPerformByteGasOverhead.toNumber() *
                        performData.length,
                    ),
                  )
                const gasCapMinusOverhead =
                  gasOverheadCap.sub(chargedGasOverhead)
                assert.isTrue(
                  gasCapMinusOverhead.gt(BigNumber.from(0)),
                  'Gas overhead got capped. Verify gas overhead variables in test match those in the registry. To not have the overheads capped increase REGISTRY_GAS_OVERHEAD by atleast ' +
                    gasCapMinusOverhead.toString(),
                )
                // total gas charged should be greater than tx gas but within gasCalculationMargin
                assert.isTrue(
                  chargedGasOverhead.gt(actualGasOverhead),
                  'Gas overhead calculated is too low, increase account gas variables (ACCOUNTING_FIXED_GAS_OVERHEAD/ACCOUNTING_PER_SIGNER_GAS_OVERHEAD) by atleast ' +
                    actualGasOverhead.sub(chargedGasOverhead).toString(),
                )

                assert.isTrue(
                  chargedGasOverhead
                    .sub(actualGasOverhead)
                    .lt(BigNumber.from(gasCalculationMargin)),
                ),
                  'Gas overhead calculated is too high, decrease account gas variables (ACCOUNTING_FIXED_GAS_OVERHEAD/ACCOUNTING_PER_SIGNER_GAS_OVERHEAD)  by atleast ' +
                    chargedGasOverhead
                      .sub(chargedGasOverhead)
                      .sub(BigNumber.from(gasCalculationMargin))
                      .toString()
              }
            }
          }
        }
      })
    })

    describe('When upkeeps are batched', () => {
      const numPassingUpkeepsArray = [1, 2, 10]
      const numFailingUpkeepsArray = [0, 1, 3]

      numPassingUpkeepsArray.forEach(function (numPassingUpkeeps) {
        numFailingUpkeepsArray.forEach(function (numFailingUpkeeps) {
          describe(
            'passing upkeeps ' +
              numPassingUpkeeps.toString() +
              ', failing upkeeps ' +
              numFailingUpkeeps.toString(),
            () => {
              let passingUpkeepIds: string[]
              let failingUpkeepIds: string[]

              beforeEach(async () => {
                passingUpkeepIds = []
                failingUpkeepIds = []
                for (let i = 0; i < numPassingUpkeeps; i++) {
                  mock = await upkeepMockFactory.deploy()
                  const tx = await registry
                    .connect(owner)
                    .registerUpkeep(
                      mock.address,
                      executeGas,
                      await admin.getAddress(),
                      randomBytes,
                      emptyBytes,
                    )
                  upkeepId = await getUpkeepID(tx)
                  passingUpkeepIds.push(upkeepId.toString())

                  // Add funds to passing upkeeps
                  await registry.connect(admin).addFunds(upkeepId, toWei('100'))
                }
                for (let i = 0; i < numFailingUpkeeps; i++) {
                  mock = await upkeepMockFactory.deploy()
                  const tx = await registry
                    .connect(owner)
                    .registerUpkeep(
                      mock.address,
                      executeGas,
                      await admin.getAddress(),
                      randomBytes,
                      emptyBytes,
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
                const registryPremiumBefore = (await registry.getState()).state
                  .totalPremium
                const registrationPassingBefore = await Promise.all(
                  passingUpkeepIds.map(async (id) => {
                    const reg = await registry.getUpkeep(BigNumber.from(id))
                    assert.equal(reg.lastPerformBlockNumber.toString(), '0')
                    return reg
                  }),
                )
                const registrationFailingBefore = await await Promise.all(
                  failingUpkeepIds.map(async (id) => {
                    const reg = await registry.getUpkeep(BigNumber.from(id))
                    assert.equal(reg.lastPerformBlockNumber.toString(), '0')
                    return reg
                  }),
                )

                const tx = await getTransmitTx(
                  registry,
                  keeper1,
                  passingUpkeepIds.concat(failingUpkeepIds),
                  f + 1,
                )

                const receipt = await tx.wait()
                const upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
                // exactly numPassingUpkeeps Upkeep Performed should be emitted
                assert.equal(upkeepPerformedLogs.length, numPassingUpkeeps)
                const insufficientFundsLogs =
                  parseInsufficientFundsUpkeepReportLogs(receipt)
                // exactly numFailingUpkeeps Upkeep Performed should be emitted
                assert.equal(insufficientFundsLogs.length, numFailingUpkeeps)

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
                const registryPremiumAfter = (await registry.getState()).state
                  .totalPremium
                const premium = registryPremiumAfter.sub(registryPremiumBefore)

                let netPayment = BigNumber.from('0')
                for (let i = 0; i < numPassingUpkeeps; i++) {
                  const id = upkeepPerformedLogs[i].args.id
                  const gasUsed = upkeepPerformedLogs[i].args.gasUsed
                  const gasOverhead = upkeepPerformedLogs[i].args.gasOverhead
                  const totalPayment = upkeepPerformedLogs[i].args.totalPayment

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
                  // InsufficientFunds log should be emitted
                  const id = insufficientFundsLogs[i].args.id
                  assert.equal(id.toString(), failingUpkeepIds[i])

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

                // Keeper payment is gasPayment + premium / num keepers
                const keeperPayment = netPayment
                  .sub(premium)
                  .add(premium.div(BigNumber.from(keeperAddresses.length)))

                // Keeper should be paid net payment for all passed upkeeps
                assert.equal(
                  keeperAfter.balance.sub(keeperPayment).toString(),
                  keeperBefore.balance.toString(),
                )

                assert.isTrue(keeperLinkAfter.eq(keeperLinkBefore))
                assert.isTrue(registryLinkBefore.eq(registryLinkAfter))
              })

              it('splits gas overhead appropriately among performed upkeeps [ @skip-coverage ]', async () => {
                // Perform the upkeeps once to remove non-zero storage slots and have predictable gas measurement
                let tx = await getTransmitTx(
                  registry,
                  keeper1,
                  passingUpkeepIds.concat(failingUpkeepIds),
                  f + 1,
                )

                await tx.wait()

                // Do the actual thing

                tx = await getTransmitTx(
                  registry,
                  keeper1,
                  passingUpkeepIds.concat(failingUpkeepIds),
                  f + 1,
                )

                const receipt = await tx.wait()
                const upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
                // exactly numPassingUpkeeps Upkeep Performed should be emitted
                assert.equal(upkeepPerformedLogs.length, numPassingUpkeeps)

                const gasOverheadCap = registryGasOverhead.add(
                  registryPerSignerGasOverhead.mul(BigNumber.from(f + 1)),
                )

                const overheadCanGetCapped =
                  numPassingUpkeeps == 1 && numFailingUpkeeps > 0
                // Should only happen with 1 successful upkeep and some failing upkeeps.
                // With 2 successful upkeeps and upto 3 failing upkeeps, overhead should be small enough to not get capped
                let netGasUsedPlusOverhead = BigNumber.from('0')

                for (let i = 0; i < numPassingUpkeeps; i++) {
                  const gasUsed = upkeepPerformedLogs[i].args.gasUsed
                  const gasOverhead = upkeepPerformedLogs[i].args.gasOverhead

                  assert.isTrue(gasUsed.gt(BigNumber.from('0')))
                  assert.isTrue(gasOverhead.gt(BigNumber.from('0')))

                  // Overhead should not exceed capped
                  assert.isTrue(gasOverhead.lte(gasOverheadCap))

                  // Overhead should be same for every upkeep since they have equal performData, hence same caps
                  assert.isTrue(
                    gasOverhead.eq(upkeepPerformedLogs[0].args.gasOverhead),
                  )

                  netGasUsedPlusOverhead = netGasUsedPlusOverhead
                    .add(gasUsed)
                    .add(gasOverhead)
                }

                const overheadsGotCapped =
                  upkeepPerformedLogs[0].args.gasOverhead.eq(gasOverheadCap)
                // Should only get capped in certain scenarios
                if (overheadsGotCapped) {
                  assert.isTrue(
                    overheadCanGetCapped,
                    'Gas overhead got capped. Verify gas overhead variables in test match those in the registry. To not have the overheads capped increase REGISTRY_GAS_OVERHEAD',
                  )
                }

                console.log(
                  'Gas Benchmarking - batching (passedUpkeeps: ',
                  numPassingUpkeeps,
                  'failedUpkeeps:',
                  numFailingUpkeeps,
                  '): ',
                  'overheadsGotCapped',
                  overheadsGotCapped,
                  'calculated overhead',
                  upkeepPerformedLogs[0].args.gasOverhead.toString(),
                  ' margin over gasUsed',
                  netGasUsedPlusOverhead.sub(receipt.gasUsed).toString(),
                )

                // If overheads dont get capped then total gas charged should be greater than tx gas
                // We don't check whether the net is within gasMargin as the margin changes with numFailedUpkeeps
                // Which is ok, as long as individual gas overhead is capped
                if (!overheadsGotCapped) {
                  assert.isTrue(
                    netGasUsedPlusOverhead.gt(receipt.gasUsed),
                    'Gas overhead is too low, increase ACCOUNTING_PER_UPKEEP_GAS_OVERHEAD',
                  )
                }
              })
            },
          )
        })
      })

      it('has enough perform gas overhead for large batches [ @skip-coverage ]', async () => {
        const numUpkeeps = 20
        const upkeepIds: string[] = []
        let totalExecuteGas = BigNumber.from('0')
        for (let i = 0; i < numUpkeeps; i++) {
          mock = await upkeepMockFactory.deploy()
          const tx = await registry
            .connect(owner)
            .registerUpkeep(
              mock.address,
              executeGas,
              await admin.getAddress(),
              randomBytes,
              emptyBytes,
            )
          upkeepId = await getUpkeepID(tx)
          upkeepIds.push(upkeepId.toString())

          // Add funds to passing upkeeps
          await registry.connect(owner).addFunds(upkeepId, toWei('10'))

          mock.setCanPerform(true)
          mock.setPerformGasToBurn(executeGas)

          totalExecuteGas = totalExecuteGas.add(executeGas)
        }

        // Should revert with no overhead added
        await evmRevert(
          getTransmitTx(registry, keeper1, upkeepIds, f + 1, {
            gasLimit: totalExecuteGas,
          }),
        )
        // Should not revert with overhead added
        await getTransmitTx(registry, keeper1, upkeepIds, f + 1, {
          gasLimit: totalExecuteGas.add(transmitGasOverhead),
        })
      })

      it('splits l2 payment among performed upkeeps', async () => {
        const numUpkeeps = 7
        const upkeepIds: string[] = []
        // Same as MockArbGasInfo.sol
        const l1CostWeiArb = BigNumber.from(1000000)

        // Deploy a new registry since we change payment model
        const registryLogic = await keeperRegistryLogicFactory
          .connect(owner)
          .deploy(
            Mode.ARBITRUM,
            linkToken.address,
            linkEthFeed.address,
            gasPriceFeed.address,
          )
        // Deploy a new registry since we change payment model
        const registry = await keeperRegistryFactory
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
        await linkToken.connect(owner).approve(registry.address, toWei('10000'))
        for (let i = 0; i < numUpkeeps; i++) {
          mock = await upkeepMockFactory.deploy()
          const tx = await registry
            .connect(owner)
            .registerUpkeep(
              mock.address,
              executeGas,
              await admin.getAddress(),
              randomBytes,
              emptyBytes,
            )
          upkeepId = await getUpkeepID(tx)
          upkeepIds.push(upkeepId.toString())

          // Add funds to passing upkeeps
          await registry.connect(owner).addFunds(upkeepId, toWei('100'))
        }

        // Do the thing
        const tx = await getTransmitTx(
          registry,
          keeper1,
          upkeepIds,
          f + 1,
          { gasPrice: gasWei.mul('5') }, // High gas price so that it gets capped
        )

        const receipt = await tx.wait()
        const upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
        // exactly numPassingUpkeeps Upkeep Performed should be emitted
        assert.equal(upkeepPerformedLogs.length, numUpkeeps)

        // Verify the payment calculation in upkeepPerformed[0]
        const upkeepPerformedLog = upkeepPerformedLogs[0]

        const gasUsed = upkeepPerformedLog.args.gasUsed
        const gasOverhead = upkeepPerformedLog.args.gasOverhead
        const totalPayment = upkeepPerformedLog.args.totalPayment

        assert.equal(
          linkForGas(
            gasUsed,
            gasOverhead,
            gasCeilingMultiplier,
            paymentPremiumPPB,
            flatFeeMicroLink,
            l1CostWeiArb.div(gasCeilingMultiplier), // Dividing by gasCeilingMultiplier as it gets multiplied later
            BigNumber.from(numUpkeeps),
          ).total.toString(),
          totalPayment.toString(),
        )
      })
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
          emptyBytes,
          emptyBytes,
        )

      const id1 = await getUpkeepID(tx)
      await registry.connect(admin).addFunds(id1, toWei('5'))

      await getTransmitTx(registry, keeper1, [id1.toString()], f + 1)
      await getTransmitTx(registry, keeper2, [id1.toString()], f + 1)
      await getTransmitTx(registry, keeper3, [id1.toString()], f + 1)

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
          emptyBytes,
        )
      const id2 = await getUpkeepID(tx2)
      await registry.connect(admin).addFunds(id2, toWei('5'))

      await getTransmitTx(registry, keeper1, [id2.toString()], f + 1)
      await getTransmitTx(registry, keeper2, [id2.toString()], f + 1)
      await getTransmitTx(registry, keeper3, [id2.toString()], f + 1)

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
    it('calculates the minimum balance appropriately', async () => {
      await mock.setCanCheck(true)

      const oneWei = BigNumber.from(1)
      const minBalance = await registry.getMinBalanceForUpkeep(upkeepId)
      const tooLow = minBalance.sub(oneWei)

      await registry.connect(admin).addFunds(upkeepId, tooLow)
      let checkUpkeepResult = await registry
        .connect(zeroAddress)
        .callStatic.checkUpkeep(upkeepId)

      assert.equal(checkUpkeepResult.upkeepNeeded, false)
      assert.equal(
        checkUpkeepResult.upkeepFailureReason,
        UpkeepFailureReason.INSUFFICIENT_BALANCE,
      )

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
          randomBytes,
          emptyBytes,
        )
      const upkeepID1 = await getUpkeepID(tx1)
      const tx2 = await registry
        .connect(owner)
        .registerUpkeep(
          mock.address,
          executeGas,
          await admin.getAddress(),
          randomBytes,
          emptyBytes,
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
      assert.equal(
        checkUpkeepResult.upkeepFailureReason,
        UpkeepFailureReason.INSUFFICIENT_BALANCE,
      )

      checkUpkeepResult = await registry
        .connect(zeroAddress)
        .callStatic.checkUpkeep(upkeepID2)
      assert.equal(checkUpkeepResult.upkeepNeeded, true)

      // upkeep 1 perform should return with insufficient balance using max performData size
      let maxPerformData = '0x'
      for (let i = 0; i < maxPerformDataSize.toNumber(); i++) {
        maxPerformData += '11'
      }

      const tx = await getTransmitTx(
        registry,
        keeper1,
        [upkeepID1.toString()],
        f + 1,
        { gasPrice: gasWei.mul(gasCeilingMultiplier) },
        maxPerformData,
      )

      const receipt = await tx.wait()
      const insufficientFundsUpkeepReportLogs =
        parseInsufficientFundsUpkeepReportLogs(receipt)
      // exactly 1 InsufficientFundsUpkeepReportLogs log should be emitted
      assert.equal(insufficientFundsUpkeepReportLogs.length, 1)

      // upkeep 1 perform should succeed with empty performData
      await getTransmitTx(
        registry,
        keeper1,
        [upkeepID1.toString()],
        f + 1,
        { gasPrice: gasWei.mul(gasCeilingMultiplier) },
        '0x',
      ),
        // upkeep 2 perform should succeed with max performData size
        await getTransmitTx(
          registry,
          keeper1,
          [upkeepID2.toString()],
          f + 1,
          { gasPrice: gasWei.mul(gasCeilingMultiplier) },
          maxPerformData,
        )
    })
  })

  describe('#withdrawFunds', () => {
    let upkeepId2: BigNumber

    beforeEach(async () => {
      const tx = await registry
        .connect(owner)
        .registerUpkeep(
          mock.address,
          executeGas,
          await admin.getAddress(),
          randomBytes,
          emptyBytes,
        )
      upkeepId2 = await getUpkeepID(tx)

      await registry.connect(admin).addFunds(upkeepId, toWei('100'))
      await registry.connect(admin).addFunds(upkeepId2, toWei('100'))

      // Do a perform so that upkeep is charged some amount
      await getTransmitTx(registry, keeper1, [upkeepId.toString()], f + 1)
      await getTransmitTx(registry, keeper1, [upkeepId2.toString()], f + 1)
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
        await registry.connect(owner).cancelUpkeep(upkeepId2)
      })

      it('can be called successively on two upkeeps', async () => {
        await registry
          .connect(admin)
          .withdrawFunds(upkeepId, await payee1.getAddress())
        await registry
          .connect(admin)
          .withdrawFunds(upkeepId2, await payee1.getAddress())
      })

      it('moves the funds out and updates the balance and emits an event', async () => {
        const payee1Before = await linkToken.balanceOf(
          await payee1.getAddress(),
        )
        const registryBefore = await linkToken.balanceOf(registry.address)

        let registration = await registry.getUpkeep(upkeepId)
        const previousBalance = registration.balance

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

      const simulatePerformResult = await registry
        .connect(zeroAddress)
        .callStatic.simulatePerformUpkeep(upkeepId, '0x')

      assert.equal(simulatePerformResult.success, false)
      assert.isTrue(simulatePerformResult.gasUsed.gt(BigNumber.from('0'))) // Some gas should be used
    })

    it('returns true and gasUsed when perform succeeds', async () => {
      await mock.setCanPerform(true)

      const simulatePerformResult = await registry
        .connect(zeroAddress)
        .callStatic.simulatePerformUpkeep(upkeepId, '0x')

      assert.equal(simulatePerformResult.success, true)
      assert.isTrue(simulatePerformResult.gasUsed.gt(BigNumber.from('0'))) // Some gas should be used
    })

    it('returns correct amount of gasUsed when perform succeeds', async () => {
      await mock.setCanPerform(true)
      await mock.setPerformGasToBurn(executeGas)

      const simulatePerformResult = await registry
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

      const checkUpkeepResult = await registry
        .connect(zeroAddress)
        .callStatic.checkUpkeep(upkeepId)

      assert.equal(checkUpkeepResult.upkeepNeeded, false)
      assert.equal(checkUpkeepResult.performData, '0x')
      assert.equal(
        checkUpkeepResult.upkeepFailureReason,
        UpkeepFailureReason.UPKEEP_CANCELLED,
      )
      assert.equal(checkUpkeepResult.gasUsed.toString(), '0')
    })

    it('returns false and error code if the upkeep is cancelled by owner', async () => {
      await registry.connect(owner).cancelUpkeep(upkeepId)

      const checkUpkeepResult = await registry
        .connect(zeroAddress)
        .callStatic.checkUpkeep(upkeepId)

      assert.equal(checkUpkeepResult.upkeepNeeded, false)
      assert.equal(checkUpkeepResult.performData, '0x')
      assert.equal(
        checkUpkeepResult.upkeepFailureReason,
        UpkeepFailureReason.UPKEEP_CANCELLED,
      )
      assert.equal(checkUpkeepResult.gasUsed.toString(), '0')
    })

    it('returns false and error code if the upkeep is paused', async () => {
      await registry.connect(admin).pauseUpkeep(upkeepId)

      const checkUpkeepResult = await registry
        .connect(zeroAddress)
        .callStatic.checkUpkeep(upkeepId)

      assert.equal(checkUpkeepResult.upkeepNeeded, false)
      assert.equal(checkUpkeepResult.performData, '0x')
      assert.equal(
        checkUpkeepResult.upkeepFailureReason,
        UpkeepFailureReason.UPKEEP_PAUSED,
      )
      assert.equal(checkUpkeepResult.gasUsed.toString(), '0')
    })

    it('returns false and error code if user is out of funds', async () => {
      const checkUpkeepResult = await registry
        .connect(zeroAddress)
        .callStatic.checkUpkeep(upkeepId)

      assert.equal(checkUpkeepResult.upkeepNeeded, false)
      assert.equal(checkUpkeepResult.performData, '0x')
      assert.equal(
        checkUpkeepResult.upkeepFailureReason,
        UpkeepFailureReason.INSUFFICIENT_BALANCE,
      )
      assert.equal(checkUpkeepResult.gasUsed.toString(), '0')
    })

    context('when the registration is funded', () => {
      beforeEach(async () => {
        await linkToken.connect(admin).approve(registry.address, toWei('100'))
        await registry.connect(admin).addFunds(upkeepId, toWei('100'))
      })

      it('returns false, error code, and revert data if the target check reverts', async () => {
        await mock.setShouldRevertCheck(true)
        const checkUpkeepResult = await registry
          .connect(zeroAddress)
          .callStatic.checkUpkeep(upkeepId)
        assert.equal(checkUpkeepResult.upkeepNeeded, false)

        const wrappedPerfromData = ethers.utils.defaultAbiCoder.decode(
          [
            'tuple(uint32 checkBlockNum, bytes32 checkBlockHash, bytes performData)',
          ],
          checkUpkeepResult.performData,
        )
        const revertReasonBytes = `0x${wrappedPerfromData[0][2].slice(10)}` // remove sighash
        assert.equal(
          ethers.utils.defaultAbiCoder.decode(['string'], revertReasonBytes)[0],
          'shouldRevertCheck should be false',
        )
        assert.equal(
          checkUpkeepResult.upkeepFailureReason,
          UpkeepFailureReason.TARGET_CHECK_REVERTED,
        )
        assert.isTrue(checkUpkeepResult.gasUsed.gt(BigNumber.from('0'))) // Some gas should be used
      })

      it('returns false and error code if the upkeep is not needed', async () => {
        await mock.setCanCheck(false)
        const checkUpkeepResult = await registry
          .connect(zeroAddress)
          .callStatic.checkUpkeep(upkeepId)

        assert.equal(checkUpkeepResult.upkeepNeeded, false)
        assert.equal(checkUpkeepResult.performData, '0x')
        assert.equal(
          checkUpkeepResult.upkeepFailureReason,
          UpkeepFailureReason.UPKEEP_NOT_NEEDED,
        )
        assert.isTrue(checkUpkeepResult.gasUsed.gt(BigNumber.from('0'))) // Some gas should be used
      })

      it('returns false and error code if the performData exceeds limit', async () => {
        let longBytes = '0x'
        for (let i = 0; i < 5000; i++) {
          longBytes += '1'
        }
        await mock.setCanCheck(true)
        await mock.setPerformData(longBytes)

        const checkUpkeepResult = await registry
          .connect(zeroAddress)
          .callStatic.checkUpkeep(upkeepId)

        assert.equal(checkUpkeepResult.upkeepNeeded, false)
        assert.equal(checkUpkeepResult.performData, '0x')
        assert.equal(
          checkUpkeepResult.upkeepFailureReason,
          UpkeepFailureReason.PERFORM_DATA_EXCEEDS_LIMIT,
        )
        assert.isTrue(checkUpkeepResult.gasUsed.gt(BigNumber.from('0'))) // Some gas should be used
      })

      it('returns true with wrapped perform data and gas used if the target can execute', async () => {
        await mock.setCanCheck(true)
        await mock.setPerformData(randomBytes)

        const latestBlock = await ethers.provider.getBlock('latest')

        const checkUpkeepResult = await registry
          .connect(zeroAddress)
          .callStatic.checkUpkeep(upkeepId, {
            blockTag: latestBlock.number,
          })

        const wrappedPerfromData = ethers.utils.defaultAbiCoder.decode(
          [
            'tuple(uint32 checkBlockNum, bytes32 checkBlockHash, bytes performData)',
          ],
          checkUpkeepResult.performData,
        )

        assert.equal(checkUpkeepResult.upkeepNeeded, true)
        assert.equal(
          wrappedPerfromData[0].checkBlockNum,
          latestBlock.number - 1,
        )
        assert.equal(
          wrappedPerfromData[0].checkBlockHash,
          latestBlock.parentHash,
        )
        assert.equal(wrappedPerfromData[0].performData, randomBytes)
        assert.equal(
          checkUpkeepResult.upkeepFailureReason,
          UpkeepFailureReason.NONE,
        )
        assert.isTrue(checkUpkeepResult.gasUsed.gt(BigNumber.from('0'))) // Some gas should be used
        assert.isTrue(checkUpkeepResult.fastGasWei.eq(gasWei))
        assert.isTrue(checkUpkeepResult.linkNative.eq(linkEth))
      })

      it('has a large enough gas overhead to cover upkeeps that use all their gas [ @skip-coverage ]', async () => {
        await mock.setCanCheck(true)
        await mock.setCheckGasToBurn(checkGasLimit)
        const gas = checkGasLimit.add(checkGasOverhead)
        const checkUpkeepResult = await registry
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
          randomBytes,
          emptyBytes,
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
      const upkeepIds = await registry.getActiveUpkeepIDs(0, 0)
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
      await verifyMaxPayment(
        Mode.DEFAULT,
        multipliers,
        gasAmounts,
        premiums,
        flatFees,
      )
    })

    it('calculates the max fee appropriately for Arbitrum', async () => {
      await verifyMaxPayment(
        Mode.ARBITRUM,
        multipliers,
        gasAmounts,
        premiums,
        flatFees,
        l1CostWeiArb,
      )
    })

    it('calculates the max fee appropriately for Optimism', async () => {
      await verifyMaxPayment(
        Mode.OPTIMISM,
        multipliers,
        gasAmounts,
        premiums,
        flatFees,
        l1CostWeiOpt,
      )
    })

    it('uses the fallback gas price if the feed has issues', async () => {
      const expectedFallbackMaxPayment = linkForGas(
        executeGas,
        registryGasOverhead
          .add(registryPerSignerGasOverhead.mul(f + 1))
          .add(maxPerformDataSize.mul(registryPerPerformByteGasOverhead)),
        gasCeilingMultiplier.mul('2'), // fallbackGasPrice is 2x gas price
        paymentPremiumPPB,
        flatFeeMicroLink,
      ).total

      // Stale feed
      let roundId = 99
      const answer = 100
      let updatedAt = 946684800 // New Years 2000 🥳
      let startedAt = 946684799
      await gasPriceFeed
        .connect(owner)
        .updateRoundData(roundId, answer, updatedAt, startedAt)

      assert.equal(
        expectedFallbackMaxPayment.toString(),
        (await registry.getMaxPaymentForGas(executeGas)).toString(),
      )

      // Negative feed price
      roundId = 100
      updatedAt = Math.floor(Date.now() / 1000)
      startedAt = 946684799
      await gasPriceFeed
        .connect(owner)
        .updateRoundData(roundId, -100, updatedAt, startedAt)

      assert.equal(
        expectedFallbackMaxPayment.toString(),
        (await registry.getMaxPaymentForGas(executeGas)).toString(),
      )

      // Zero feed price
      roundId = 101
      updatedAt = Math.floor(Date.now() / 1000)
      startedAt = 946684799
      await gasPriceFeed
        .connect(owner)
        .updateRoundData(roundId, 0, updatedAt, startedAt)

      assert.equal(
        expectedFallbackMaxPayment.toString(),
        (await registry.getMaxPaymentForGas(executeGas)).toString(),
      )
    })

    it('uses the fallback link price if the feed has issues', async () => {
      const expectedFallbackMaxPayment = linkForGas(
        executeGas,
        registryGasOverhead
          .add(registryPerSignerGasOverhead.mul(f + 1))
          .add(maxPerformDataSize.mul(registryPerPerformByteGasOverhead)),
        gasCeilingMultiplier.mul('2'), // fallbackLinkPrice is 1/2 link price, so multiply by 2
        paymentPremiumPPB,
        flatFeeMicroLink,
      ).total

      // Stale feed
      let roundId = 99
      const answer = 100
      let updatedAt = 946684800 // New Years 2000 🥳
      let startedAt = 946684799
      await linkEthFeed
        .connect(owner)
        .updateRoundData(roundId, answer, updatedAt, startedAt)

      assert.equal(
        expectedFallbackMaxPayment.toString(),
        (await registry.getMaxPaymentForGas(executeGas)).toString(),
      )

      // Negative feed price
      roundId = 100
      updatedAt = Math.floor(Date.now() / 1000)
      startedAt = 946684799
      await linkEthFeed
        .connect(owner)
        .updateRoundData(roundId, -100, updatedAt, startedAt)

      assert.equal(
        expectedFallbackMaxPayment.toString(),
        (await registry.getMaxPaymentForGas(executeGas)).toString(),
      )

      // Zero feed price
      roundId = 101
      updatedAt = Math.floor(Date.now() / 1000)
      startedAt = 946684799
      await linkEthFeed
        .connect(owner)
        .updateRoundData(roundId, 0, updatedAt, startedAt)

      assert.equal(
        expectedFallbackMaxPayment.toString(),
        (await registry.getMaxPaymentForGas(executeGas)).toString(),
      )
    })
  })

  describe('#typeAndVersion', () => {
    it('uses the correct type and version', async () => {
      const typeAndVersion = await registry.typeAndVersion()
      assert.equal(typeAndVersion, 'KeeperRegistry 2.0.2')
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
      assert.equal(updatedState.latestEpoch, 0)

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
      const signers = [randomAddress()]
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
      const newSigners = [
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
      const newTransmitters = [
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
      // Perform an upkeep so that totalPremium is updated
      await registry.connect(admin).addFunds(upkeepId, toWei('100'))
      let tx = await getTransmitTx(
        registry,
        keeper1,
        [upkeepId.toString()],
        f + 1,
      )
      await tx.wait()

      const newOffChainVersion = BigNumber.from('2')
      const newOffChainConfig = '0x1122'

      const old = await registry.getState()
      const oldState = old.state
      assert(oldState.totalPremium.gt(BigNumber.from('0')))

      const newSigners = newKeepers
      tx = await registry
        .connect(owner)
        .setConfig(
          newSigners,
          newKeepers,
          f,
          encodeConfig(config),
          newOffChainVersion,
          newOffChainConfig,
        )

      const updated = await registry.getState()
      const updatedState = updated.state
      assert(oldState.totalPremium.eq(updatedState.totalPremium))

      // Old signer addresses which are not in new signers should be non active
      for (let i = 0; i < signerAddresses.length; i++) {
        const signer = signerAddresses[i]
        if (!newSigners.includes(signer)) {
          assert((await registry.getSignerInfo(signer)).active == false)
          assert((await registry.getSignerInfo(signer)).index == 0)
        }
      }
      // New signer addresses should be active
      for (let i = 0; i < newSigners.length; i++) {
        const signer = newSigners[i]
        assert((await registry.getSignerInfo(signer)).active == true)
        assert((await registry.getSignerInfo(signer)).index == i)
      }
      // Old transmitter addresses which are not in new transmitter should be non active, update lastCollected but retain other info
      for (let i = 0; i < keeperAddresses.length; i++) {
        const transmitter = keeperAddresses[i]
        if (!newKeepers.includes(transmitter)) {
          assert(
            (await registry.getTransmitterInfo(transmitter)).active == false,
          )
          assert((await registry.getTransmitterInfo(transmitter)).index == i)
          assert(
            (
              await registry.getTransmitterInfo(transmitter)
            ).lastCollected.toString() == oldState.totalPremium.toString(),
          )
        }
      }
      // New transmitter addresses should be active
      for (let i = 0; i < newKeepers.length; i++) {
        const transmitter = newKeepers[i]
        assert((await registry.getTransmitterInfo(transmitter)).active == true)
        assert((await registry.getTransmitterInfo(transmitter)).index == i)
        assert(
          (
            await registry.getTransmitterInfo(transmitter)
          ).lastCollected.toString() == oldState.totalPremium.toString(),
        )
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
            emptyBytes,
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
            emptyBytes,
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
            longBytes,
            emptyBytes,
          ),
        'CheckDataExceedsLimit()',
      )
    })

    it('creates a record of the registration', async () => {
      const executeGases = [100000, 500000]
      const checkDatas = [emptyBytes, '0x12']
      const offchainConfig = '0x1234567890'

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
              checkData,
              offchainConfig,
            )

          //confirm the upkeep details
          upkeepId = await getUpkeepID(tx)
          await expect(tx)
            .to.emit(registry, 'UpkeepRegistered')
            .withArgs(upkeepId, executeGas, await admin.getAddress())
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
          assert.equal(registration.offchainConfig, offchainConfig)
          assert(registration.maxValidBlocknumber.eq('0xffffffff'))
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

  describe('#setUpkeepOffchainConfig', () => {
    const newConfig = '0xc0ffeec0ffee'

    it('reverts if the registration does not exist', async () => {
      await evmRevert(
        registry
          .connect(admin)
          .setUpkeepOffchainConfig(upkeepId.add(1), newConfig),
        'OnlyCallableByAdmin()',
      )
    })

    it('reverts if the upkeep is canceled', async () => {
      await registry.connect(admin).cancelUpkeep(upkeepId)
      await evmRevert(
        registry.connect(admin).setUpkeepOffchainConfig(upkeepId, newConfig),
        'UpkeepCancelled()',
      )
    })

    it('reverts if called by anyone but the admin', async () => {
      await evmRevert(
        registry.connect(owner).setUpkeepOffchainConfig(upkeepId, newConfig),
        'OnlyCallableByAdmin()',
      )
    })

    it('updates the config successfully', async () => {
      const initialConfig = (await registry.getUpkeep(upkeepId)).offchainConfig
      assert.equal(initialConfig, '0x')
      await registry.connect(admin).setUpkeepOffchainConfig(upkeepId, newConfig)
      const updatedConfig = (await registry.getUpkeep(upkeepId)).offchainConfig
      assert.equal(newConfig, updatedConfig)
    })

    it('emits a log', async () => {
      const tx = await registry
        .connect(admin)
        .setUpkeepOffchainConfig(upkeepId, newConfig)
      await expect(tx)
        .to.emit(registry, 'UpkeepOffchainConfigSet')
        .withArgs(upkeepId, newConfig)
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
      const minUpkeepSpend = toWei('1000')
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
      const upkeepBalance = (await registry.getUpkeep(upkeepId)).balance
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
        getTransmitTx(registry, keeper1, [upkeepId.toString()], f + 1),
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
            emptyBytes,
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
        .deploy(
          Mode.DEFAULT,
          linkToken.address,
          linkEthFeed.address,
          gasPriceFeed.address,
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
        const payee = (await registry.getTransmitterInfo(keeperAddresses[i]))
          .payee
        assert.equal(payee, zeroAddress)
      }

      await registry.connect(owner).setPayees(payees)

      for (let i = 0; i < keeperAddresses.length; i++) {
        const payee = (await registry.getTransmitterInfo(keeperAddresses[i]))
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

        const tx = await getTransmitTx(
          registry,
          keeper1,
          [upkeepId.toString()],
          f + 1,
        )
        const receipt = await tx.wait()
        const cancelledUpkeepReportLogs =
          parseCancelledUpkeepReportLogs(receipt)
        // exactly 1 CancelledUpkeepReport log should be emitted
        assert.equal(cancelledUpkeepReportLogs.length, 1)
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

        await getTransmitTx(registry, keeper1, [upkeepId.toString()], f + 1)

        for (let i = 0; i < cancellationDelay; i++) {
          await ethers.provider.send('evm_mine', [])
        }

        const tx = await getTransmitTx(
          registry,
          keeper1,
          [upkeepId.toString()],
          f + 1,
        )

        const receipt = await tx.wait()
        const cancelledUpkeepReportLogs =
          parseCancelledUpkeepReportLogs(receipt)
        // exactly 1 CancelledUpkeepReport log should be emitted
        assert.equal(cancelledUpkeepReportLogs.length, 1)
      })

      describe('when an upkeep has been performed', async () => {
        beforeEach(async () => {
          await linkToken.connect(owner).approve(registry.address, toWei('100'))
          await registry.connect(owner).addFunds(upkeepId, toWei('100'))
          await getTransmitTx(registry, keeper1, [upkeepId.toString()], f + 1)
        })

        it('deducts a cancellation fee from the upkeep and gives to owner', async () => {
          const minUpkeepSpend = toWei('10')

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
          const upkeepBefore = (await registry.getUpkeep(upkeepId)).balance
          const ownerBefore = (await registry.getState()).state.ownerLinkBalance

          const amountSpent = toWei('100').sub(upkeepBefore)
          const cancellationFee = minUpkeepSpend.sub(amountSpent)

          await registry.connect(admin).cancelUpkeep(upkeepId)

          const payee1After = await linkToken.balanceOf(
            await payee1.getAddress(),
          )
          const upkeepAfter = (await registry.getUpkeep(upkeepId)).balance
          const ownerAfter = (await registry.getState()).state.ownerLinkBalance

          // post upkeep balance should be previous balance minus cancellation fee
          assert.isTrue(upkeepBefore.sub(cancellationFee).eq(upkeepAfter))
          // payee balance should not change
          assert.isTrue(payee1Before.eq(payee1After))
          // owner should receive the cancellation fee
          assert.isTrue(ownerAfter.sub(ownerBefore).eq(cancellationFee))
        })

        it('deducts up to balance as cancellation fee', async () => {
          // Very high min spend, should deduct whole balance as cancellation fees
          const minUpkeepSpend = toWei('1000')
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
          const upkeepBefore = (await registry.getUpkeep(upkeepId)).balance
          const ownerBefore = (await registry.getState()).state.ownerLinkBalance

          await registry.connect(admin).cancelUpkeep(upkeepId)
          const payee1After = await linkToken.balanceOf(
            await payee1.getAddress(),
          )
          const ownerAfter = (await registry.getState()).state.ownerLinkBalance
          const upkeepAfter = (await registry.getUpkeep(upkeepId)).balance

          // all upkeep balance is deducted for cancellation fee
          assert.equal(0, upkeepAfter.toNumber())
          // payee balance should not change
          assert.isTrue(payee1After.eq(payee1Before))
          // all upkeep balance is transferred to the owner
          assert.isTrue(ownerAfter.sub(ownerBefore).eq(upkeepBefore))
        })

        it('does not deduct cancellation fee if more than minUpkeepSpend is spent', async () => {
          // Very low min spend, already spent in one perform upkeep
          const minUpkeepSpend = BigNumber.from(420)
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
          const upkeepBefore = (await registry.getUpkeep(upkeepId)).balance
          const ownerBefore = (await registry.getState()).state.ownerLinkBalance

          await registry.connect(admin).cancelUpkeep(upkeepId)
          const payee1After = await linkToken.balanceOf(
            await payee1.getAddress(),
          )
          const ownerAfter = (await registry.getState()).state.ownerLinkBalance
          const upkeepAfter = (await registry.getUpkeep(upkeepId)).balance

          // upkeep does not pay cancellation fee after cancellation because minimum upkeep spent is met
          assert.isTrue(upkeepBefore.eq(upkeepAfter))
          // owner balance does not change
          assert.isTrue(ownerAfter.eq(ownerBefore))
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
      await getTransmitTx(registry, keeper1, [upkeepId.toString()], f + 1)
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
      const keeperBefore = await registry.getTransmitterInfo(
        await keeper1.getAddress(),
      )
      const registrationBefore = (await registry.getUpkeep(upkeepId)).balance
      const toLinkBefore = await linkToken.balanceOf(to)
      const registryLinkBefore = await linkToken.balanceOf(registry.address)
      const registryPremiumBefore = (await registry.getState()).state
        .totalPremium
      const ownerBefore = (await registry.getState()).state.ownerLinkBalance

      // Withdrawing for first time, last collected = 0
      assert.equal(keeperBefore.lastCollected.toString(), '0')

      //// Do the thing
      await registry
        .connect(payee1)
        .withdrawPayment(await keeper1.getAddress(), to)

      const keeperAfter = await registry.getTransmitterInfo(
        await keeper1.getAddress(),
      )
      const registrationAfter = (await registry.getUpkeep(upkeepId)).balance
      const toLinkAfter = await linkToken.balanceOf(to)
      const registryLinkAfter = await linkToken.balanceOf(registry.address)
      const registryPremiumAfter = (await registry.getState()).state
        .totalPremium
      const ownerAfter = (await registry.getState()).state.ownerLinkBalance

      // registry total premium should not change
      assert.isTrue(registryPremiumBefore.eq(registryPremiumAfter))
      // Last collected should be updated
      assert.equal(
        keeperAfter.lastCollected.toString(),
        registryPremiumBefore.toString(),
      )

      const spareChange = registryPremiumBefore.mod(
        BigNumber.from(keeperAddresses.length),
      )
      // spare change should go to owner
      assert.isTrue(ownerAfter.sub(spareChange).eq(ownerBefore))

      assert.isTrue(keeperAfter.balance.eq(BigNumber.from(0)))
      assert.isTrue(registrationBefore.eq(registrationAfter))
      assert.isTrue(toLinkBefore.add(keeperBefore.balance).eq(toLinkAfter))
      assert.isTrue(
        registryLinkBefore.sub(keeperBefore.balance).eq(registryLinkAfter),
      )
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
