import { ethers } from 'hardhat'
import { loadFixture } from '@nomicfoundation/hardhat-network-helpers'
import { assert, expect } from 'chai'
import {
  BigNumber,
  BigNumberish,
  BytesLike,
  Contract,
  ContractFactory,
  ContractReceipt,
  ContractTransaction,
  Signer,
  Wallet,
} from 'ethers'
import { evmRevert, evmRevertCustomError } from '../../test-helpers/matchers'
import { getUsers, Personas } from '../../test-helpers/setup'
import { randomAddress, toWei } from '../../test-helpers/helpers'
import { StreamsLookupUpkeep__factory as StreamsLookupUpkeepFactory } from '../../../typechain/factories/StreamsLookupUpkeep__factory'
import { MockV3Aggregator__factory as MockV3AggregatorFactory } from '../../../typechain/factories/MockV3Aggregator__factory'
import { UpkeepMock__factory as UpkeepMockFactory } from '../../../typechain/factories/UpkeepMock__factory'
import { UpkeepAutoFunder__factory as UpkeepAutoFunderFactory } from '../../../typechain/factories/UpkeepAutoFunder__factory'
import { MockZKSyncSystemContext__factory as MockZKSyncSystemContextFactory } from '../../../typechain/factories/MockZKSyncSystemContext__factory'
import { ChainModuleBase__factory as ChainModuleBaseFactory } from '../../../typechain/factories/ChainModuleBase__factory'
import { MockGasBoundCaller__factory as MockGasBoundCallerFactory } from '../../../typechain/factories/MockGasBoundCaller__factory'
import { ILogAutomation__factory as ILogAutomationactory } from '../../../typechain/factories/ILogAutomation__factory'
import { AutomationCompatibleUtils } from '../../../typechain/AutomationCompatibleUtils'
import { StreamsLookupUpkeep } from '../../../typechain/StreamsLookupUpkeep'
import { MockV3Aggregator } from '../../../typechain/MockV3Aggregator'
import { MockGasBoundCaller } from '../../../typechain/MockGasBoundCaller'
import { UpkeepMock } from '../../../typechain/UpkeepMock'
import { ChainModuleBase } from '../../../typechain/ChainModuleBase'
import { UpkeepTranscoder } from '../../../typechain/UpkeepTranscoder'
import { MockZKSyncSystemContext } from '../../../typechain/MockZKSyncSystemContext'
import { IChainModule, UpkeepAutoFunder } from '../../../typechain'
import {
  CancelledUpkeepReportEvent,
  IAutomationRegistryMaster2_3 as IAutomationRegistry,
  ReorgedUpkeepReportEvent,
  StaleUpkeepReportEvent,
  UpkeepPerformedEvent,
} from '../../../typechain/IAutomationRegistryMaster2_3'
import {
  deployMockContract,
  MockContract,
} from '@ethereum-waffle/mock-contract'
import { deployZKSyncRegistry23 } from './helpers'
import { AutomationUtils2_3 } from '../../../typechain/AutomationUtils2_3'

const describeMaybe = process.env.SKIP_SLOW ? describe.skip : describe
const itMaybe = process.env.SKIP_SLOW ? it.skip : it

// copied from AutomationRegistryInterface2_3.sol
enum UpkeepFailureReason {
  NONE,
  UPKEEP_CANCELLED,
  UPKEEP_PAUSED,
  TARGET_CHECK_REVERTED,
  UPKEEP_NOT_NEEDED,
  PERFORM_DATA_EXCEEDS_LIMIT,
  INSUFFICIENT_BALANCE,
  CHECK_CALLBACK_REVERTED,
  REVERT_DATA_EXCEEDS_LIMIT,
  REGISTRY_PAUSED,
}

// copied from AutomationRegistryBase2_3.sol
enum Trigger {
  CONDITION,
  LOG,
}

// un-exported types that must be extracted from the utils contract
type Report = Parameters<AutomationUtils2_3['_report']>[0]
type LogTrigger = Parameters<AutomationCompatibleUtils['_logTrigger']>[0]
type ConditionalTrigger = Parameters<
  AutomationCompatibleUtils['_conditionalTrigger']
>[0]
type Log = Parameters<AutomationCompatibleUtils['_log']>[0]
type OnChainConfig = Parameters<IAutomationRegistry['setConfigTypeSafe']>[3]

// -----------------------------------------------------------------------------------------------

// These values should match the constants declared in registry
let registryConditionalOverhead: BigNumber
let registryLogOverhead: BigNumber
let registryPerSignerGasOverhead: BigNumber
// let registryPerPerformByteGasOverhead: BigNumber
// let registryTransmitCalldataFixedBytesOverhead: BigNumber
// let registryTransmitCalldataPerSignerBytesOverhead: BigNumber
let cancellationDelay: number

// This is the margin for gas that we test for. Gas charged should always be greater
// than total gas used in tx but should not increase beyond this margin
// const gasCalculationMargin = BigNumber.from(50_000)
// This is the margin for gas overhead estimation in checkUpkeep. The estimated gas
// overhead should be larger than actual gas overhead but should not increase beyond this margin
// const gasEstimationMargin = BigNumber.from(50_000)

// 1 Link = 0.005 Eth
const linkUSD = BigNumber.from('2000000000') // 1 LINK = $20
const nativeUSD = BigNumber.from('400000000000') // 1 ETH = $4000
const gasWei = BigNumber.from(1000000000) // 1 gwei
// -----------------------------------------------------------------------------------------------
// test-wide configs for upkeeps
const performGas = BigNumber.from('1000000')
const paymentPremiumBase = BigNumber.from('1000000000')
const paymentPremiumPPB = BigNumber.from('250000000')
const flatFeeMilliCents = BigNumber.from(0)

const randomBytes = '0x1234abcd'
const emptyBytes = '0x'
const emptyBytes32 =
  '0x0000000000000000000000000000000000000000000000000000000000000000'

const pubdataGas = BigNumber.from(500000)
const transmitGasOverhead = 1_040_000
const checkGasOverhead = 600_000

const stalenessSeconds = BigNumber.from(43820)
const gasCeilingMultiplier = BigNumber.from(2)
const checkGasLimit = BigNumber.from(10000000)
const fallbackGasPrice = gasWei.mul(BigNumber.from('2'))
const fallbackLinkPrice = linkUSD.div(BigNumber.from('2'))
const fallbackNativePrice = nativeUSD.div(BigNumber.from('2'))
const maxCheckDataSize = BigNumber.from(1000)
const maxPerformDataSize = BigNumber.from(1000)
const maxRevertDataSize = BigNumber.from(1000)
const maxPerformGas = BigNumber.from(5000000)
const minUpkeepSpend = BigNumber.from(0)
const f = 1
const offchainVersion = 1
const offchainBytes = '0x'
const zeroAddress = ethers.constants.AddressZero
const wrappedNativeTokenAddress = '0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2'
const epochAndRound5_1 =
  '0x0000000000000000000000000000000000000000000000000000000000000501'

let logTriggerConfig: string

// -----------------------------------------------------------------------------------------------

// Smart contract factories
let linkTokenFactory: ContractFactory
let mockV3AggregatorFactory: MockV3AggregatorFactory
let mockGasBoundCallerFactory: MockGasBoundCallerFactory
let upkeepMockFactory: UpkeepMockFactory
let upkeepAutoFunderFactory: UpkeepAutoFunderFactory
let moduleBaseFactory: ChainModuleBaseFactory
let mockZKSyncSystemContextFactory: MockZKSyncSystemContextFactory
let streamsLookupUpkeepFactory: StreamsLookupUpkeepFactory
let personas: Personas

// contracts
let linkToken: Contract
let linkUSDFeed: MockV3Aggregator
let nativeUSDFeed: MockV3Aggregator
let gasPriceFeed: MockV3Aggregator
let registry: IAutomationRegistry // default registry, used for most tests
let mgRegistry: IAutomationRegistry // "migrate registry" used in migration tests
let mock: UpkeepMock
let autoFunderUpkeep: UpkeepAutoFunder
let ltUpkeep: MockContract
let transcoder: UpkeepTranscoder
let moduleBase: ChainModuleBase
let mockGasBoundCaller: MockGasBoundCaller
let mockZKSyncSystemContext: MockZKSyncSystemContext
let streamsLookupUpkeep: StreamsLookupUpkeep
let automationUtils: AutomationCompatibleUtils
let automationUtils2_3: AutomationUtils2_3

function now() {
  return Math.floor(Date.now() / 1000)
}

async function getUpkeepID(tx: ContractTransaction): Promise<BigNumber> {
  const receipt = await tx.wait()
  for (const event of receipt.events || []) {
    if (
      event.args &&
      event.eventSignature == 'UpkeepRegistered(uint256,uint32,address)'
    ) {
      return event.args[0]
    }
  }
  throw new Error('could not find upkeep ID in tx event logs')
}

const getTriggerType = (upkeepId: BigNumber): Trigger => {
  const hexBytes = ethers.utils.defaultAbiCoder.encode(['uint256'], [upkeepId])
  const bytes = ethers.utils.arrayify(hexBytes)
  for (let idx = 4; idx < 15; idx++) {
    if (bytes[idx] != 0) {
      return Trigger.CONDITION
    }
  }
  return bytes[15] as Trigger
}

const encodeBlockTrigger = (conditionalTrigger: ConditionalTrigger) => {
  return (
    '0x' +
    automationUtils.interface
      .encodeFunctionData('_conditionalTrigger', [conditionalTrigger])
      .slice(10)
  )
}

const encodeLogTrigger = (logTrigger: LogTrigger) => {
  return (
    '0x' +
    automationUtils.interface
      .encodeFunctionData('_logTrigger', [logTrigger])
      .slice(10)
  )
}

const encodeLog = (log: Log) => {
  return (
    '0x' + automationUtils.interface.encodeFunctionData('_log', [log]).slice(10)
  )
}

const encodeReport = (report: Report) => {
  return (
    '0x' +
    automationUtils2_3.interface
      .encodeFunctionData('_report', [report])
      .slice(10)
  )
}

type UpkeepData = {
  Id: BigNumberish
  performGas: BigNumberish
  performData: BytesLike
  trigger: BytesLike
}

const makeReport = (upkeeps: UpkeepData[]) => {
  const upkeepIds = upkeeps.map((u) => u.Id)
  const performGases = upkeeps.map((u) => u.performGas)
  const triggers = upkeeps.map((u) => u.trigger)
  const performDatas = upkeeps.map((u) => u.performData)
  return encodeReport({
    fastGasWei: gasWei,
    linkUSD,
    upkeepIds,
    gasLimits: performGases,
    triggers,
    performDatas,
  })
}

const makeLatestBlockReport = async (upkeepsIDs: BigNumberish[]) => {
  const latestBlock = await ethers.provider.getBlock('latest')
  const upkeeps: UpkeepData[] = []
  for (let i = 0; i < upkeepsIDs.length; i++) {
    upkeeps.push({
      Id: upkeepsIDs[i],
      performGas,
      trigger: encodeBlockTrigger({
        blockNum: latestBlock.number,
        blockHash: latestBlock.hash,
      }),
      performData: '0x',
    })
  }
  return makeReport(upkeeps)
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

const parseUpkeepPerformedLogs = (receipt: ContractReceipt) => {
  const parsedLogs = []
  for (const rawLog of receipt.logs) {
    try {
      const log = registry.interface.parseLog(rawLog)
      if (
        log.name ==
        registry.interface.events[
          'UpkeepPerformed(uint256,bool,uint96,uint256,uint256,bytes)'
        ].name
      ) {
        parsedLogs.push(log as unknown as UpkeepPerformedEvent)
      }
    } catch {
      continue
    }
  }
  return parsedLogs
}

const parseReorgedUpkeepReportLogs = (receipt: ContractReceipt) => {
  const parsedLogs = []
  for (const rawLog of receipt.logs) {
    try {
      const log = registry.interface.parseLog(rawLog)
      if (
        log.name ==
        registry.interface.events['ReorgedUpkeepReport(uint256,bytes)'].name
      ) {
        parsedLogs.push(log as unknown as ReorgedUpkeepReportEvent)
      }
    } catch {
      continue
    }
  }
  return parsedLogs
}

const parseStaleUpkeepReportLogs = (receipt: ContractReceipt) => {
  const parsedLogs = []
  for (const rawLog of receipt.logs) {
    try {
      const log = registry.interface.parseLog(rawLog)
      if (
        log.name ==
        registry.interface.events['StaleUpkeepReport(uint256,bytes)'].name
      ) {
        parsedLogs.push(log as unknown as StaleUpkeepReportEvent)
      }
    } catch {
      continue
    }
  }
  return parsedLogs
}

const parseCancelledUpkeepReportLogs = (receipt: ContractReceipt) => {
  const parsedLogs = []
  for (const rawLog of receipt.logs) {
    try {
      const log = registry.interface.parseLog(rawLog)
      if (
        log.name ==
        registry.interface.events['CancelledUpkeepReport(uint256,bytes)'].name
      ) {
        parsedLogs.push(log as unknown as CancelledUpkeepReportEvent)
      }
    } catch {
      continue
    }
  }
  return parsedLogs
}

describe('ZKSyncAutomationRegistry2_3', () => {
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
  let financeAdmin: Signer

  let upkeepId: BigNumber // conditional upkeep
  let afUpkeepId: BigNumber // auto funding upkeep
  let logUpkeepId: BigNumber // log trigger upkeepID
  let streamsLookupUpkeepId: BigNumber // streams lookup upkeep
  // const numUpkeeps = 4 // see above
  let keeperAddresses: string[]
  let payees: string[]
  let signers: Wallet[]
  let signerAddresses: string[]
  let config: OnChainConfig
  let baseConfig: Parameters<IAutomationRegistry['setConfigTypeSafe']>
  let upkeepManager: string

  before(async () => {
    personas = (await getUsers()).personas

    const compatibleUtilsFactory = await ethers.getContractFactory(
      'AutomationCompatibleUtils',
    )
    automationUtils = await compatibleUtilsFactory.deploy()

    const utilsFactory = await ethers.getContractFactory('AutomationUtils2_3')
    automationUtils2_3 = await utilsFactory.deploy()

    linkTokenFactory = await ethers.getContractFactory(
      'src/v0.8/shared/test/helpers/LinkTokenTestHelper.sol:LinkTokenTestHelper',
    )
    // need full path because there are two contracts with name MockV3Aggregator
    mockV3AggregatorFactory = (await ethers.getContractFactory(
      'src/v0.8/tests/MockV3Aggregator.sol:MockV3Aggregator',
    )) as unknown as MockV3AggregatorFactory
    mockZKSyncSystemContextFactory = await ethers.getContractFactory(
      'MockZKSyncSystemContext',
    )
    mockGasBoundCallerFactory =
      await ethers.getContractFactory('MockGasBoundCaller')
    upkeepMockFactory = await ethers.getContractFactory('UpkeepMock')
    upkeepAutoFunderFactory =
      await ethers.getContractFactory('UpkeepAutoFunder')
    moduleBaseFactory = await ethers.getContractFactory('ChainModuleBase')
    streamsLookupUpkeepFactory = await ethers.getContractFactory(
      'StreamsLookupUpkeep',
    )

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
    upkeepManager = await personas.Norbert.getAddress()
    financeAdmin = personas.Nick
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

    logTriggerConfig =
      '0x' +
      automationUtils.interface
        .encodeFunctionData('_logTriggerConfig', [
          {
            contractAddress: randomAddress(),
            filterSelector: 0,
            topic0: ethers.utils.randomBytes(32),
            topic1: ethers.utils.randomBytes(32),
            topic2: ethers.utils.randomBytes(32),
            topic3: ethers.utils.randomBytes(32),
          },
        ])
        .slice(10)
  })

  // This function is similar to registry's _calculatePaymentAmount
  // It uses global fastGasWei, linkEth, and assumes isExecution = false (gasFee = fastGasWei*multiplier)
  // rest of the parameters are the same
  const linkForGas = (
    upkeepGasSpent: BigNumber,
    gasOverhead: BigNumber,
    gasMultiplier: BigNumber,
    premiumPPB: BigNumber,
    flatFee: BigNumber, // in millicents
  ) => {
    const gasSpent = gasOverhead.add(BigNumber.from(upkeepGasSpent))
    const gasPayment = gasWei
      .mul(gasMultiplier)
      .mul(gasSpent)
      .mul(nativeUSD)
      .div(linkUSD)

    const premium = gasWei
      .mul(gasMultiplier)
      .mul(upkeepGasSpent)
      .mul(premiumPPB)
      .mul(nativeUSD)
      .div(paymentPremiumBase)
      .add(flatFee.mul(BigNumber.from(10).pow(21)))
      .div(linkUSD)

    return {
      total: gasPayment.add(premium),
      gasPayment,
      premium,
    }
  }

  const verifyMaxPayment = async (
    registry: IAutomationRegistry,
    chainModule: IChainModule,
  ) => {
    type TestCase = {
      name: string
      multiplier: number
      gas: number
      premium: number
      flatFee: number
    }

    const tests: TestCase[] = [
      {
        name: 'no fees',
        multiplier: 1,
        gas: 100000,
        premium: 0,
        flatFee: 0,
      },
      {
        name: 'basic fees',
        multiplier: 1,
        gas: 100000,
        premium: 250000000,
        flatFee: 1000000,
      },
      {
        name: 'max fees',
        multiplier: 3,
        gas: 10000000,
        premium: 250000000,
        flatFee: 1000000,
      },
    ]

    const fPlusOne = BigNumber.from(f + 1)
    const chainModuleOverheads = await chainModule.getGasOverhead()
    const totalConditionalOverhead = registryConditionalOverhead
      .add(registryPerSignerGasOverhead.mul(fPlusOne))
      .add(chainModuleOverheads.chainModuleFixedOverhead)

    const totalLogOverhead = registryLogOverhead
      .add(registryPerSignerGasOverhead.mul(fPlusOne))
      .add(chainModuleOverheads.chainModuleFixedOverhead)

    const financeAdminAddress = await financeAdmin.getAddress()

    for (const test of tests) {
      await registry.connect(owner).setConfigTypeSafe(
        signerAddresses,
        keeperAddresses,
        f,
        {
          checkGasLimit,
          stalenessSeconds,
          gasCeilingMultiplier: test.multiplier,
          maxCheckDataSize,
          maxPerformDataSize,
          maxRevertDataSize,
          maxPerformGas,
          fallbackGasPrice,
          fallbackLinkPrice,
          fallbackNativePrice,
          transcoder: transcoder.address,
          registrars: [],
          upkeepPrivilegeManager: upkeepManager,
          chainModule: chainModule.address,
          reorgProtectionEnabled: true,
          financeAdmin: financeAdminAddress,
        },
        offchainVersion,
        offchainBytes,
        [linkToken.address],
        [
          {
            gasFeePPB: test.premium,
            flatFeeMilliCents: test.flatFee,
            priceFeed: linkUSDFeed.address,
            fallbackPrice: fallbackLinkPrice,
            minSpend: minUpkeepSpend,
            decimals: 18,
          },
        ],
      )

      const conditionalPrice = await registry.getMaxPaymentForGas(
        upkeepId,
        Trigger.CONDITION,
        test.gas,
        linkToken.address,
      )
      expect(conditionalPrice).to.equal(
        linkForGas(
          BigNumber.from(test.gas),
          totalConditionalOverhead,
          BigNumber.from(test.multiplier),
          BigNumber.from(test.premium),
          BigNumber.from(test.flatFee),
        ).total,
      )

      const logPrice = await registry.getMaxPaymentForGas(
        upkeepId,
        Trigger.LOG,
        test.gas,
        linkToken.address,
      )
      expect(logPrice).to.equal(
        linkForGas(
          BigNumber.from(test.gas),
          totalLogOverhead,
          BigNumber.from(test.multiplier),
          BigNumber.from(test.premium),
          BigNumber.from(test.flatFee),
        ).total,
      )
    }
  }

  const verifyConsistentAccounting = async (
    maxAllowedSpareChange: BigNumber,
  ) => {
    const expectedLinkBalance = await registry.getReserveAmount(
      linkToken.address,
    )
    const linkTokenBalance = await linkToken.balanceOf(registry.address)
    const upkeepIdBalance = (await registry.getUpkeep(upkeepId)).balance
    let totalKeeperBalance = BigNumber.from(0)
    for (let i = 0; i < keeperAddresses.length; i++) {
      totalKeeperBalance = totalKeeperBalance.add(
        (await registry.getTransmitterInfo(keeperAddresses[i])).balance,
      )
    }

    const linkAvailableForPayment = await registry.linkAvailableForPayment()
    assert.isTrue(expectedLinkBalance.eq(linkTokenBalance))
    assert.isTrue(
      upkeepIdBalance
        .add(totalKeeperBalance)
        .add(linkAvailableForPayment)
        .lte(expectedLinkBalance),
    )
    assert.isTrue(
      expectedLinkBalance
        .sub(upkeepIdBalance)
        .sub(totalKeeperBalance)
        .sub(linkAvailableForPayment)
        .lte(maxAllowedSpareChange),
    )
  }

  interface GetTransmitTXOptions {
    numSigners?: number
    startingSignerIndex?: number
    gasLimit?: BigNumberish
    gasPrice?: BigNumberish
    performGas?: BigNumberish
    performDatas?: string[]
    checkBlockNum?: number
    checkBlockHash?: string
    logBlockHash?: BytesLike
    txHash?: BytesLike
    logIndex?: number
    timestamp?: number
  }

  const getTransmitTx = async (
    registry: IAutomationRegistry,
    transmitter: Signer,
    upkeepIds: BigNumber[],
    overrides: GetTransmitTXOptions = {},
  ) => {
    const latestBlock = await ethers.provider.getBlock('latest')
    const configDigest = (await registry.getState()).state.latestConfigDigest
    const config = {
      numSigners: f + 1,
      startingSignerIndex: 0,
      performDatas: undefined,
      performGas,
      checkBlockNum: latestBlock.number,
      checkBlockHash: latestBlock.hash,
      logIndex: 0,
      txHash: undefined, // assigned uniquely below
      logBlockHash: undefined, // assigned uniquely below
      timestamp: now(),
      gasLimit: undefined,
      gasPrice: undefined,
    }
    Object.assign(config, overrides)
    const upkeeps: UpkeepData[] = []
    for (let i = 0; i < upkeepIds.length; i++) {
      let trigger: string
      switch (getTriggerType(upkeepIds[i])) {
        case Trigger.CONDITION:
          trigger = encodeBlockTrigger({
            blockNum: config.checkBlockNum,
            blockHash: config.checkBlockHash,
          })
          break
        case Trigger.LOG:
          trigger = encodeLogTrigger({
            logBlockHash: config.logBlockHash || ethers.utils.randomBytes(32),
            txHash: config.txHash || ethers.utils.randomBytes(32),
            logIndex: config.logIndex,
            blockNum: config.checkBlockNum,
            blockHash: config.checkBlockHash,
          })
          break
      }
      upkeeps.push({
        Id: upkeepIds[i],
        performGas: config.performGas,
        trigger,
        performData: config.performDatas ? config.performDatas[i] : '0x',
      })
    }

    const report = makeReport(upkeeps)
    const reportContext = [configDigest, epochAndRound5_1, emptyBytes32]
    const sigs = signReport(
      reportContext,
      report,
      signers.slice(
        config.startingSignerIndex,
        config.startingSignerIndex + config.numSigners,
      ),
    )

    type txOverride = {
      gasLimit?: BigNumberish | Promise<BigNumberish>
      gasPrice?: BigNumberish | Promise<BigNumberish>
    }
    const txOverrides: txOverride = {}
    if (config.gasLimit) {
      txOverrides.gasLimit = config.gasLimit
    }
    if (config.gasPrice) {
      txOverrides.gasPrice = config.gasPrice
    }

    return registry
      .connect(transmitter)
      .transmit(
        [configDigest, epochAndRound5_1, emptyBytes32],
        report,
        sigs.rs,
        sigs.ss,
        sigs.vs,
        txOverrides,
      )
  }

  const getTransmitTxWithReport = async (
    registry: IAutomationRegistry,
    transmitter: Signer,
    report: BytesLike,
  ) => {
    const configDigest = (await registry.getState()).state.latestConfigDigest
    const reportContext = [configDigest, epochAndRound5_1, emptyBytes32]
    const sigs = signReport(reportContext, report, signers.slice(0, f + 1))

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

  const setup = async () => {
    linkToken = await linkTokenFactory.connect(owner).deploy()
    gasPriceFeed = await mockV3AggregatorFactory
      .connect(owner)
      .deploy(0, gasWei)
    linkUSDFeed = await mockV3AggregatorFactory
      .connect(owner)
      .deploy(8, linkUSD)
    nativeUSDFeed = await mockV3AggregatorFactory
      .connect(owner)
      .deploy(8, nativeUSD)
    const upkeepTranscoderFactory = await ethers.getContractFactory(
      'UpkeepTranscoder5_0',
    )
    transcoder = await upkeepTranscoderFactory.connect(owner).deploy()
    mockZKSyncSystemContext = await mockZKSyncSystemContextFactory
      .connect(owner)
      .deploy()
    mockGasBoundCaller = await mockGasBoundCallerFactory.connect(owner).deploy()
    moduleBase = await moduleBaseFactory.connect(owner).deploy()
    streamsLookupUpkeep = await streamsLookupUpkeepFactory
      .connect(owner)
      .deploy(
        BigNumber.from('10000'),
        BigNumber.from('100'),
        false /* useArbBlock */,
        true /* staging */,
        false /* verify mercury response */,
      )

    const zksyncSystemContextCode = await ethers.provider.send('eth_getCode', [
      mockZKSyncSystemContext.address,
    ])
    await ethers.provider.send('hardhat_setCode', [
      '0x000000000000000000000000000000000000800B',
      zksyncSystemContextCode,
    ])

    const gasBoundCallerCode = await ethers.provider.send('eth_getCode', [
      mockGasBoundCaller.address,
    ])
    await ethers.provider.send('hardhat_setCode', [
      '0xc706EC7dfA5D4Dc87f29f859094165E8290530f5',
      gasBoundCallerCode,
    ])

    const financeAdminAddress = await financeAdmin.getAddress()

    config = {
      checkGasLimit,
      stalenessSeconds,
      gasCeilingMultiplier,
      maxCheckDataSize,
      maxPerformDataSize,
      maxRevertDataSize,
      maxPerformGas,
      fallbackGasPrice,
      fallbackLinkPrice,
      fallbackNativePrice,
      transcoder: transcoder.address,
      registrars: [],
      upkeepPrivilegeManager: upkeepManager,
      chainModule: moduleBase.address,
      reorgProtectionEnabled: true,
      financeAdmin: financeAdminAddress,
    }

    baseConfig = [
      signerAddresses,
      keeperAddresses,
      f,
      config,
      offchainVersion,
      offchainBytes,
      [linkToken.address],
      [
        {
          gasFeePPB: paymentPremiumPPB,
          flatFeeMilliCents,
          priceFeed: linkUSDFeed.address,
          fallbackPrice: fallbackLinkPrice,
          minSpend: minUpkeepSpend,
          decimals: 18,
        },
      ],
    ]

    const registryParams: Parameters<typeof deployZKSyncRegistry23> = [
      owner,
      linkToken.address,
      linkUSDFeed.address,
      nativeUSDFeed.address,
      gasPriceFeed.address,
      zeroAddress,
      0, // onchain payout mode
      wrappedNativeTokenAddress,
    ]

    registry = await deployZKSyncRegistry23(...registryParams)
    mgRegistry = await deployZKSyncRegistry23(...registryParams)

    registryConditionalOverhead = await registry.getConditionalGasOverhead()
    registryLogOverhead = await registry.getLogGasOverhead()
    registryPerSignerGasOverhead = await registry.getPerSignerGasOverhead()
    // registryPerPerformByteGasOverhead =
    //   await registry.getPerPerformByteGasOverhead()
    // registryTransmitCalldataFixedBytesOverhead =
    //   await registry.getTransmitCalldataFixedBytesOverhead()
    // registryTransmitCalldataPerSignerBytesOverhead =
    //   await registry.getTransmitCalldataPerSignerBytesOverhead()
    cancellationDelay = (await registry.getCancellationDelay()).toNumber()

    await registry.connect(owner).setConfigTypeSafe(...baseConfig)
    await mgRegistry.connect(owner).setConfigTypeSafe(...baseConfig)
    for (const reg of [registry, mgRegistry]) {
      await reg.connect(owner).setPayees(payees)
      await linkToken.connect(admin).approve(reg.address, toWei('1000'))
      await linkToken.connect(owner).approve(reg.address, toWei('1000'))
    }

    mock = await upkeepMockFactory.deploy()
    await linkToken
      .connect(owner)
      .transfer(await admin.getAddress(), toWei('1000'))
    let tx = await registry
      .connect(owner)
      .registerUpkeep(
        mock.address,
        performGas,
        await admin.getAddress(),
        Trigger.CONDITION,
        linkToken.address,
        randomBytes,
        '0x',
        '0x',
      )
    upkeepId = await getUpkeepID(tx)

    autoFunderUpkeep = await upkeepAutoFunderFactory
      .connect(owner)
      .deploy(linkToken.address, registry.address)
    tx = await registry
      .connect(owner)
      .registerUpkeep(
        autoFunderUpkeep.address,
        performGas,
        autoFunderUpkeep.address,
        Trigger.CONDITION,
        linkToken.address,
        '0x',
        '0x',
        '0x',
      )
    afUpkeepId = await getUpkeepID(tx)

    ltUpkeep = await deployMockContract(owner, ILogAutomationactory.abi)
    tx = await registry
      .connect(owner)
      .registerUpkeep(
        ltUpkeep.address,
        performGas,
        await admin.getAddress(),
        Trigger.LOG,
        linkToken.address,
        '0x',
        logTriggerConfig,
        emptyBytes,
      )
    logUpkeepId = await getUpkeepID(tx)

    await autoFunderUpkeep.setUpkeepId(afUpkeepId)
    // Give enough funds for upkeep as well as to the upkeep contract
    await linkToken
      .connect(owner)
      .transfer(autoFunderUpkeep.address, toWei('1000'))

    tx = await registry
      .connect(owner)
      .registerUpkeep(
        streamsLookupUpkeep.address,
        performGas,
        await admin.getAddress(),
        Trigger.CONDITION,
        linkToken.address,
        '0x',
        '0x',
        '0x',
      )
    streamsLookupUpkeepId = await getUpkeepID(tx)
  }

  const getMultipleUpkeepsDeployedAndFunded = async (
    numPassingConditionalUpkeeps: number,
    numPassingLogUpkeeps: number,
    numFailingUpkeeps: number,
  ) => {
    const passingConditionalUpkeepIds = []
    const passingLogUpkeepIds = []
    const failingUpkeepIds = []
    for (let i = 0; i < numPassingConditionalUpkeeps; i++) {
      const mock = await upkeepMockFactory.deploy()
      await mock.setCanPerform(true)
      await mock.setPerformGasToBurn(BigNumber.from('0'))
      const tx = await registry
        .connect(owner)
        .registerUpkeep(
          mock.address,
          performGas,
          await admin.getAddress(),
          Trigger.CONDITION,
          linkToken.address,
          '0x',
          '0x',
          '0x',
        )
      const condUpkeepId = await getUpkeepID(tx)
      passingConditionalUpkeepIds.push(condUpkeepId)

      // Add funds to passing upkeeps
      await registry.connect(admin).addFunds(condUpkeepId, toWei('100'))
    }
    for (let i = 0; i < numPassingLogUpkeeps; i++) {
      const mock = await upkeepMockFactory.deploy()
      await mock.setCanPerform(true)
      await mock.setPerformGasToBurn(BigNumber.from('0'))
      const tx = await registry
        .connect(owner)
        .registerUpkeep(
          mock.address,
          performGas,
          await admin.getAddress(),
          Trigger.LOG,
          linkToken.address,
          '0x',
          logTriggerConfig,
          emptyBytes,
        )
      const logUpkeepId = await getUpkeepID(tx)
      passingLogUpkeepIds.push(logUpkeepId)

      // Add funds to passing upkeeps
      await registry.connect(admin).addFunds(logUpkeepId, toWei('100'))
    }
    for (let i = 0; i < numFailingUpkeeps; i++) {
      const mock = await upkeepMockFactory.deploy()
      await mock.setCanPerform(true)
      await mock.setPerformGasToBurn(BigNumber.from('0'))
      const tx = await registry
        .connect(owner)
        .registerUpkeep(
          mock.address,
          performGas,
          await admin.getAddress(),
          Trigger.CONDITION,
          linkToken.address,
          '0x',
          '0x',
          '0x',
        )
      const failingUpkeepId = await getUpkeepID(tx)
      failingUpkeepIds.push(failingUpkeepId)
    }
    return {
      passingConditionalUpkeepIds,
      passingLogUpkeepIds,
      failingUpkeepIds,
    }
  }

  beforeEach(async () => {
    await loadFixture(setup)
  })

  describe('#transmit', () => {
    const fArray = [1, 5, 10]

    it('reverts when registry is paused', async () => {
      await registry.connect(owner).pause()
      await evmRevertCustomError(
        getTransmitTx(registry, keeper1, [upkeepId]),
        registry,
        'RegistryPaused',
      )
    })

    it('reverts when called by non active transmitter', async () => {
      await evmRevertCustomError(
        getTransmitTx(registry, payee1, [upkeepId]),
        registry,
        'OnlyActiveTransmitters',
      )
    })

    it('reverts when report data lengths mismatches', async () => {
      const upkeepIds = []
      const gasLimits: BigNumber[] = []
      const triggers: string[] = []
      const performDatas = []

      upkeepIds.push(upkeepId)
      gasLimits.push(performGas)
      triggers.push('0x')
      performDatas.push('0x')
      // Push an extra perform data
      performDatas.push('0x')

      const report = encodeReport({
        fastGasWei: 0,
        linkUSD: 0,
        upkeepIds,
        gasLimits,
        triggers,
        performDatas,
      })

      await evmRevertCustomError(
        getTransmitTxWithReport(registry, keeper1, report),
        registry,
        'InvalidReport',
      )
    })

    it('returns early when invalid upkeepIds are included in report', async () => {
      const tx = await getTransmitTx(registry, keeper1, [
        upkeepId.add(BigNumber.from('1')),
      ])

      const receipt = await tx.wait()
      const cancelledUpkeepReportLogs = parseCancelledUpkeepReportLogs(receipt)
      // exactly 1 CancelledUpkeepReport log should be emitted
      assert.equal(cancelledUpkeepReportLogs.length, 1)
    })

    it('performs even when the upkeep has insufficient funds and the upkeep pays out all the remaining balance', async () => {
      // add very little fund to this upkeep
      await registry.connect(admin).addFunds(upkeepId, BigNumber.from(10))
      const tx = await getTransmitTx(registry, keeper1, [upkeepId])
      const receipt = await tx.wait()
      // the upkeep is underfunded in transmit but still performed
      const upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
      assert.equal(upkeepPerformedLogs.length, 1)
      const balance = (await registry.getUpkeep(upkeepId)).balance
      assert.equal(balance.toNumber(), 0)
    })

    context('When the upkeep is funded', async () => {
      beforeEach(async () => {
        // Fund the upkeep
        await Promise.all([
          registry.connect(admin).addFunds(upkeepId, toWei('100')),
          registry.connect(admin).addFunds(logUpkeepId, toWei('100')),
        ])
      })

      it('handles duplicate upkeepIDs', async () => {
        const tests: [string, BigNumber, number, number][] = [
          // [name, upkeep, num stale, num performed]
          ['conditional', upkeepId, 1, 1], // checkBlocks must be sequential
          ['log-trigger', logUpkeepId, 0, 2], // logs are deduped based on the "trigger ID"
        ]
        for (const [type, id, nStale, nPerformed] of tests) {
          const tx = await getTransmitTx(registry, keeper1, [id, id])
          const receipt = await tx.wait()
          const staleUpkeepReport = parseStaleUpkeepReportLogs(receipt)
          const upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
          assert.equal(
            staleUpkeepReport.length,
            nStale,
            `wrong log count for ${type} upkeep`,
          )
          assert.equal(
            upkeepPerformedLogs.length,
            nPerformed,
            `wrong log count for ${type} upkeep`,
          )
        }
      })

      it('handles duplicate log triggers', async () => {
        const logBlockHash = ethers.utils.randomBytes(32)
        const txHash = ethers.utils.randomBytes(32)
        const logIndex = 0
        const expectedDedupKey = ethers.utils.solidityKeccak256(
          ['uint256', 'bytes32', 'bytes32', 'uint32'],
          [logUpkeepId, logBlockHash, txHash, logIndex],
        )
        assert.isFalse(await registry.hasDedupKey(expectedDedupKey))
        const tx = await getTransmitTx(
          registry,
          keeper1,
          [logUpkeepId, logUpkeepId],
          { logBlockHash, txHash, logIndex }, // will result in the same dedup key
        )
        const receipt = await tx.wait()
        const staleUpkeepReport = parseStaleUpkeepReportLogs(receipt)
        const upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
        assert.equal(staleUpkeepReport.length, 1)
        assert.equal(upkeepPerformedLogs.length, 1)
        assert.isTrue(await registry.hasDedupKey(expectedDedupKey))
        await expect(tx)
          .to.emit(registry, 'DedupKeyAdded')
          .withArgs(expectedDedupKey)
      })

      it('returns early when check block number is less than last perform (block)', async () => {
        // First perform an upkeep to put last perform block number on upkeep state
        const tx = await getTransmitTx(registry, keeper1, [upkeepId])
        await tx.wait()
        const lastPerformed = (await registry.getUpkeep(upkeepId))
          .lastPerformedBlockNumber
        const lastPerformBlock = await ethers.provider.getBlock(lastPerformed)
        assert.equal(lastPerformed.toString(), tx.blockNumber?.toString())
        // Try to transmit a report which has checkBlockNumber = lastPerformed-1, should result in stale report
        const transmitTx = await getTransmitTx(registry, keeper1, [upkeepId], {
          checkBlockNum: lastPerformBlock.number - 1,
          checkBlockHash: lastPerformBlock.parentHash,
        })
        const receipt = await transmitTx.wait()
        const staleUpkeepReportLogs = parseStaleUpkeepReportLogs(receipt)
        // exactly 1 StaleUpkeepReportLogs log should be emitted
        assert.equal(staleUpkeepReportLogs.length, 1)
      })

      it('handles case when check block hash does not match', async () => {
        const tests: [string, BigNumber][] = [
          ['conditional', upkeepId],
          ['log-trigger', logUpkeepId],
        ]
        for (const [type, id] of tests) {
          const latestBlock = await ethers.provider.getBlock('latest')
          // Try to transmit a report which has incorrect checkBlockHash
          const tx = await getTransmitTx(registry, keeper1, [id], {
            checkBlockNum: latestBlock.number - 1,
            checkBlockHash: latestBlock.hash, // should be latestBlock.parentHash
          })

          const receipt = await tx.wait()
          const reorgedUpkeepReportLogs = parseReorgedUpkeepReportLogs(receipt)
          // exactly 1 ReorgedUpkeepReportLogs log should be emitted
          assert.equal(
            reorgedUpkeepReportLogs.length,
            1,
            `wrong log count for ${type} upkeep`,
          )
        }
      })

      it('handles case when check block number is older than 256 blocks', async () => {
        for (let i = 0; i < 256; i++) {
          await ethers.provider.send('evm_mine', [])
        }
        const tests: [string, BigNumber][] = [
          ['conditional', upkeepId],
          ['log-trigger', logUpkeepId],
        ]
        for (const [type, id] of tests) {
          const latestBlock = await ethers.provider.getBlock('latest')
          const old = await ethers.provider.getBlock(latestBlock.number - 256)
          // Try to transmit a report which has incorrect checkBlockHash
          const tx = await getTransmitTx(registry, keeper1, [id], {
            checkBlockNum: old.number,
            checkBlockHash: old.hash,
          })

          const receipt = await tx.wait()
          const reorgedUpkeepReportLogs = parseReorgedUpkeepReportLogs(receipt)
          // exactly 1 ReorgedUpkeepReportLogs log should be emitted
          assert.equal(
            reorgedUpkeepReportLogs.length,
            1,
            `wrong log count for ${type} upkeep`,
          )
        }
      })

      it('allows bypassing reorg protection with empty blockhash', async () => {
        const tests: [string, BigNumber][] = [
          ['conditional', upkeepId],
          ['log-trigger', logUpkeepId],
        ]
        for (const [type, id] of tests) {
          const latestBlock = await ethers.provider.getBlock('latest')
          const tx = await getTransmitTx(registry, keeper1, [id], {
            checkBlockNum: latestBlock.number,
            checkBlockHash: emptyBytes32,
          })
          const receipt = await tx.wait()
          const upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
          assert.equal(
            upkeepPerformedLogs.length,
            1,
            `wrong log count for ${type} upkeep`,
          )
        }
      })

      it('allows bypassing reorg protection with reorgProtectionEnabled false config', async () => {
        const tests: [string, BigNumber][] = [
          ['conditional', upkeepId],
          ['log-trigger', logUpkeepId],
        ]
        const newConfig = config
        newConfig.reorgProtectionEnabled = false
        await registry // used to test initial configurations
          .connect(owner)
          .setConfigTypeSafe(
            signerAddresses,
            keeperAddresses,
            f,
            newConfig,
            offchainVersion,
            offchainBytes,
            baseConfig[6],
            baseConfig[7],
          )

        for (const [type, id] of tests) {
          const latestBlock = await ethers.provider.getBlock('latest')
          // Try to transmit a report which has incorrect checkBlockHash
          const tx = await getTransmitTx(registry, keeper1, [id], {
            checkBlockNum: latestBlock.number - 1,
            checkBlockHash: latestBlock.hash, // should be latestBlock.parentHash
          })

          const receipt = await tx.wait()
          const upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
          assert.equal(
            upkeepPerformedLogs.length,
            1,
            `wrong log count for ${type} upkeep`,
          )
        }
      })

      it('allows very old trigger block numbers when bypassing reorg protection with reorgProtectionEnabled config', async () => {
        const newConfig = config
        newConfig.reorgProtectionEnabled = false
        await registry // used to test initial configurations
          .connect(owner)
          .setConfigTypeSafe(
            signerAddresses,
            keeperAddresses,
            f,
            newConfig,
            offchainVersion,
            offchainBytes,
            baseConfig[6],
            baseConfig[7],
          )
        for (let i = 0; i < 256; i++) {
          await ethers.provider.send('evm_mine', [])
        }
        const tests: [string, BigNumber][] = [
          ['conditional', upkeepId],
          ['log-trigger', logUpkeepId],
        ]
        for (const [type, id] of tests) {
          const latestBlock = await ethers.provider.getBlock('latest')
          const old = await ethers.provider.getBlock(latestBlock.number - 256)
          // Try to transmit a report which has incorrect checkBlockHash
          const tx = await getTransmitTx(registry, keeper1, [id], {
            checkBlockNum: old.number,
            checkBlockHash: old.hash,
          })

          const receipt = await tx.wait()
          const upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
          assert.equal(
            upkeepPerformedLogs.length,
            1,
            `wrong log count for ${type} upkeep`,
          )
        }
      })

      it('allows very old trigger block numbers when bypassing reorg protection with empty blockhash', async () => {
        // mine enough blocks so that blockhash(1) is unavailable
        for (let i = 0; i <= 256; i++) {
          await ethers.provider.send('evm_mine', [])
        }
        const tests: [string, BigNumber][] = [
          ['conditional', upkeepId],
          ['log-trigger', logUpkeepId],
        ]
        for (const [type, id] of tests) {
          const tx = await getTransmitTx(registry, keeper1, [id], {
            checkBlockNum: 1,
            checkBlockHash: emptyBytes32,
          })
          const receipt = await tx.wait()
          const upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
          assert.equal(
            upkeepPerformedLogs.length,
            1,
            `wrong log count for ${type} upkeep`,
          )
        }
      })

      it('returns early when future block number is provided as trigger, irrespective of blockhash being present', async () => {
        const tests: [string, BigNumber][] = [
          ['conditional', upkeepId],
          ['log-trigger', logUpkeepId],
        ]
        for (const [type, id] of tests) {
          const latestBlock = await ethers.provider.getBlock('latest')

          // Should fail when blockhash is empty
          let tx = await getTransmitTx(registry, keeper1, [id], {
            checkBlockNum: latestBlock.number + 100,
            checkBlockHash: emptyBytes32,
          })
          let receipt = await tx.wait()
          let reorgedUpkeepReportLogs = parseReorgedUpkeepReportLogs(receipt)
          // exactly 1 ReorgedUpkeepReportLogs log should be emitted
          assert.equal(
            reorgedUpkeepReportLogs.length,
            1,
            `wrong log count for ${type} upkeep`,
          )

          // Should also fail when blockhash is not empty
          tx = await getTransmitTx(registry, keeper1, [id], {
            checkBlockNum: latestBlock.number + 100,
            checkBlockHash: latestBlock.hash,
          })
          receipt = await tx.wait()
          reorgedUpkeepReportLogs = parseReorgedUpkeepReportLogs(receipt)
          // exactly 1 ReorgedUpkeepReportLogs log should be emitted
          assert.equal(
            reorgedUpkeepReportLogs.length,
            1,
            `wrong log count for ${type} upkeep`,
          )
        }
      })

      it('returns early when future block number is provided as trigger, irrespective of reorgProtectionEnabled config', async () => {
        const newConfig = config
        newConfig.reorgProtectionEnabled = false
        await registry // used to test initial configurations
          .connect(owner)
          .setConfigTypeSafe(
            signerAddresses,
            keeperAddresses,
            f,
            newConfig,
            offchainVersion,
            offchainBytes,
            baseConfig[6],
            baseConfig[7],
          )
        const tests: [string, BigNumber][] = [
          ['conditional', upkeepId],
          ['log-trigger', logUpkeepId],
        ]
        for (const [type, id] of tests) {
          const latestBlock = await ethers.provider.getBlock('latest')

          // Should fail when blockhash is empty
          let tx = await getTransmitTx(registry, keeper1, [id], {
            checkBlockNum: latestBlock.number + 100,
            checkBlockHash: emptyBytes32,
          })
          let receipt = await tx.wait()
          let reorgedUpkeepReportLogs = parseReorgedUpkeepReportLogs(receipt)
          // exactly 1 ReorgedUpkeepReportLogs log should be emitted
          assert.equal(
            reorgedUpkeepReportLogs.length,
            1,
            `wrong log count for ${type} upkeep`,
          )

          // Should also fail when blockhash is not empty
          tx = await getTransmitTx(registry, keeper1, [id], {
            checkBlockNum: latestBlock.number + 100,
            checkBlockHash: latestBlock.hash,
          })
          receipt = await tx.wait()
          reorgedUpkeepReportLogs = parseReorgedUpkeepReportLogs(receipt)
          // exactly 1 ReorgedUpkeepReportLogs log should be emitted
          assert.equal(
            reorgedUpkeepReportLogs.length,
            1,
            `wrong log count for ${type} upkeep`,
          )
        }
      })

      it('returns early when upkeep is cancelled and cancellation delay has gone', async () => {
        const latestBlockReport = await makeLatestBlockReport([upkeepId])
        await registry.connect(admin).cancelUpkeep(upkeepId)

        for (let i = 0; i < cancellationDelay; i++) {
          await ethers.provider.send('evm_mine', [])
        }

        const tx = await getTransmitTxWithReport(
          registry,
          keeper1,
          latestBlockReport,
        )

        const receipt = await tx.wait()
        const cancelledUpkeepReportLogs =
          parseCancelledUpkeepReportLogs(receipt)
        // exactly 1 CancelledUpkeepReport log should be emitted
        assert.equal(cancelledUpkeepReportLogs.length, 1)
      })

      it('does not revert if the target cannot execute', async () => {
        await mock.setCanPerform(false)
        const tx = await getTransmitTx(registry, keeper1, [upkeepId])

        const receipt = await tx.wait()
        const upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
        // exactly 1 Upkeep Performed should be emitted
        assert.equal(upkeepPerformedLogs.length, 1)
        const upkeepPerformedLog = upkeepPerformedLogs[0]

        const success = upkeepPerformedLog.args.success
        assert.equal(success, false)
      })

      it('does not revert if the target runs out of gas', async () => {
        await mock.setCanPerform(false)

        const tx = await getTransmitTx(registry, keeper1, [upkeepId], {
          performGas: 10, // too little gas
        })

        const receipt = await tx.wait()
        const upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
        // exactly 1 Upkeep Performed should be emitted
        assert.equal(upkeepPerformedLogs.length, 1)
        const upkeepPerformedLog = upkeepPerformedLogs[0]

        const success = upkeepPerformedLog.args.success
        assert.equal(success, false)
      })

      it('reverts if not enough gas supplied', async () => {
        await mock.setCanPerform(true)
        await evmRevert(
          getTransmitTx(registry, keeper1, [upkeepId], {
            gasLimit: BigNumber.from(150000),
          }),
        )
      })

      it('executes the data passed to the registry', async () => {
        await mock.setCanPerform(true)

        const tx = await getTransmitTx(registry, keeper1, [upkeepId], {
          performDatas: [randomBytes],
        })
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
        // Actual multiplier is 2, but we set gasPrice to be == gasWei
        const gasPrice = gasWei
        await mock.setCanPerform(true)
        const registryPremiumBefore = (await registry.getState()).state
          .totalPremium
        const tx = await getTransmitTx(registry, keeper1, [upkeepId], {
          gasPrice,
        })
        const receipt = await tx.wait()
        const registryPremiumAfter = (await registry.getState()).state
          .totalPremium
        const premium = registryPremiumAfter.sub(registryPremiumBefore)

        const upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
        // exactly 1 Upkeep Performed should be emitted
        assert.equal(upkeepPerformedLogs.length, 1)
        const upkeepPerformedLog = upkeepPerformedLogs[0]

        const gasUsed = upkeepPerformedLog.args.gasUsed // 14657 gasUsed
        const gasOverhead = upkeepPerformedLog.args.gasOverhead // 137230 gasOverhead
        const totalPayment = upkeepPerformedLog.args.totalPayment

        assert.equal(
          linkForGas(
            gasUsed,
            gasOverhead,
            BigNumber.from('1'), // Not the config multiplier, but the actual gas used
            paymentPremiumPPB,
            flatFeeMilliCents,
            // pubdataGas.mul(gasPrice),
          ).total.toString(),
          totalPayment.toString(),
        )

        assert.equal(
          linkForGas(
            gasUsed,
            gasOverhead,
            BigNumber.from('1'), // Not the config multiplier, but the actual gas used
            paymentPremiumPPB,
            flatFeeMilliCents,
            // pubdataGas.mul(gasPrice),
          ).premium.toString(),
          premium.toString(),
        )
      })

      it('only pays at a rate up to the gas ceiling [ @skip-coverage ]', async () => {
        // Actual multiplier is 2, but we set gasPrice to be 10x
        const gasPrice = gasWei.mul(BigNumber.from('10'))
        await mock.setCanPerform(true)

        const tx = await getTransmitTx(registry, keeper1, [upkeepId], {
          gasPrice,
        })
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
            flatFeeMilliCents,
            // pubdataGas.mul(gasPrice),
          ).total.toString(),
          totalPayment.toString(),
        )
      })

      itMaybe('can self fund', async () => {
        const maxPayment = await registry.getMaxPaymentForGas(
          upkeepId,
          Trigger.CONDITION,
          performGas,
          linkToken.address,
        )

        // First set auto funding amount to 0 and verify that balance is deducted upon performUpkeep
        let initialBalance = toWei('100')
        await registry.connect(owner).addFunds(afUpkeepId, initialBalance)
        await autoFunderUpkeep.setAutoFundLink(0)
        await autoFunderUpkeep.setIsEligible(true)
        await getTransmitTx(registry, keeper1, [afUpkeepId])

        let postUpkeepBalance = (await registry.getUpkeep(afUpkeepId)).balance
        assert.isTrue(postUpkeepBalance.lt(initialBalance)) // Balance should be deducted
        assert.isTrue(postUpkeepBalance.gte(initialBalance.sub(maxPayment))) // Balance should not be deducted more than maxPayment

        // Now set auto funding amount to 100 wei and verify that the balance increases
        initialBalance = postUpkeepBalance
        const autoTopupAmount = toWei('100')
        await autoFunderUpkeep.setAutoFundLink(autoTopupAmount)
        await autoFunderUpkeep.setIsEligible(true)
        await getTransmitTx(registry, keeper1, [afUpkeepId])

        postUpkeepBalance = (await registry.getUpkeep(afUpkeepId)).balance
        // Balance should increase by autoTopupAmount and decrease by max maxPayment
        assert.isTrue(
          postUpkeepBalance.gte(
            initialBalance.add(autoTopupAmount).sub(maxPayment),
          ),
        )
      })

      it('can self cancel', async () => {
        await registry.connect(owner).addFunds(afUpkeepId, toWei('100'))

        await autoFunderUpkeep.setIsEligible(true)
        await autoFunderUpkeep.setShouldCancel(true)

        let registration = await registry.getUpkeep(afUpkeepId)
        const oldExpiration = registration.maxValidBlocknumber

        // Do the thing
        await getTransmitTx(registry, keeper1, [afUpkeepId])

        // Verify upkeep gets cancelled
        registration = await registry.getUpkeep(afUpkeepId)
        const newExpiration = registration.maxValidBlocknumber
        assert.isTrue(newExpiration.lt(oldExpiration))
      })

      it('reverts when configDigest mismatches', async () => {
        const report = await makeLatestBlockReport([upkeepId])
        const reportContext = [emptyBytes32, epochAndRound5_1, emptyBytes32] // wrong config digest
        const sigs = signReport(reportContext, report, signers.slice(0, f + 1))
        await evmRevertCustomError(
          registry
            .connect(keeper1)
            .transmit(
              [reportContext[0], reportContext[1], reportContext[2]],
              report,
              sigs.rs,
              sigs.ss,
              sigs.vs,
            ),
          registry,
          'ConfigDigestMismatch',
        )
      })

      it('reverts with incorrect number of signatures', async () => {
        const configDigest = (await registry.getState()).state
          .latestConfigDigest
        const report = await makeLatestBlockReport([upkeepId])
        const reportContext = [configDigest, epochAndRound5_1, emptyBytes32] // wrong config digest
        const sigs = signReport(reportContext, report, signers.slice(0, f + 2))
        await evmRevertCustomError(
          registry
            .connect(keeper1)
            .transmit(
              [reportContext[0], reportContext[1], reportContext[2]],
              report,
              sigs.rs,
              sigs.ss,
              sigs.vs,
            ),
          registry,
          'IncorrectNumberOfSignatures',
        )
      })

      it('reverts with invalid signature for inactive signers', async () => {
        const configDigest = (await registry.getState()).state
          .latestConfigDigest
        const report = await makeLatestBlockReport([upkeepId])
        const reportContext = [configDigest, epochAndRound5_1, emptyBytes32] // wrong config digest
        const sigs = signReport(reportContext, report, [
          new ethers.Wallet(ethers.Wallet.createRandom()),
          new ethers.Wallet(ethers.Wallet.createRandom()),
        ])
        await evmRevertCustomError(
          registry
            .connect(keeper1)
            .transmit(
              [reportContext[0], reportContext[1], reportContext[2]],
              report,
              sigs.rs,
              sigs.ss,
              sigs.vs,
            ),
          registry,
          'OnlyActiveSigners',
        )
      })

      it('reverts with invalid signature for duplicated signers', async () => {
        const configDigest = (await registry.getState()).state
          .latestConfigDigest
        const report = await makeLatestBlockReport([upkeepId])
        const reportContext = [configDigest, epochAndRound5_1, emptyBytes32] // wrong config digest
        const sigs = signReport(reportContext, report, [signer1, signer1])
        await evmRevertCustomError(
          registry
            .connect(keeper1)
            .transmit(
              [reportContext[0], reportContext[1], reportContext[2]],
              report,
              sigs.rs,
              sigs.ss,
              sigs.vs,
            ),
          registry,
          'DuplicateSigners',
        )
      })

      itMaybe(
        'has a large enough gas overhead to cover upkeep that use all its gas [ @skip-coverage ]',
        async () => {
          await registry.connect(owner).setConfigTypeSafe(
            signerAddresses,
            keeperAddresses,
            10, // maximise f to maximise overhead
            config,
            offchainVersion,
            offchainBytes,
            baseConfig[6],
            baseConfig[7],
          )
          const tx = await registry.connect(owner).registerUpkeep(
            mock.address,
            maxPerformGas, // max allowed gas
            await admin.getAddress(),
            Trigger.CONDITION,
            linkToken.address,
            '0x',
            '0x',
            '0x',
          )
          const testUpkeepId = await getUpkeepID(tx)
          await registry.connect(admin).addFunds(testUpkeepId, toWei('100'))

          let performData = '0x'
          for (let i = 0; i < maxPerformDataSize.toNumber(); i++) {
            performData += '11'
          } // max allowed performData

          await mock.setCanPerform(true)
          await mock.setPerformGasToBurn(maxPerformGas)

          await getTransmitTx(registry, keeper1, [testUpkeepId], {
            gasLimit: maxPerformGas.add(transmitGasOverhead),
            numSigners: 11,
            performDatas: [performData],
          }) // Should not revert
        },
      )

      itMaybe(
        'performs upkeep, deducts payment, updates lastPerformed and emits events',
        async () => {
          await mock.setCanPerform(true)

          for (const i in fArray) {
            const newF = fArray[i]
            await registry
              .connect(owner)
              .setConfigTypeSafe(
                signerAddresses,
                keeperAddresses,
                newF,
                config,
                offchainVersion,
                offchainBytes,
                baseConfig[6],
                baseConfig[7],
              )
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
            const registryLinkBefore = await linkToken.balanceOf(
              registry.address,
            )

            // Do the thing
            const tx = await getTransmitTx(registry, keeper1, [upkeepId], {
              checkBlockNum: checkBlock.number,
              checkBlockHash: checkBlock.hash,
              numSigners: newF + 1,
            })

            const receipt = await tx.wait()

            const upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
            // exactly 1 Upkeep Performed should be emitted
            assert.equal(upkeepPerformedLogs.length, 1)
            const upkeepPerformedLog = upkeepPerformedLogs[0]

            const id = upkeepPerformedLog.args.id
            const success = upkeepPerformedLog.args.success
            const trigger = upkeepPerformedLog.args.trigger
            const gasUsed = upkeepPerformedLog.args.gasUsed
            const gasOverhead = upkeepPerformedLog.args.gasOverhead
            const totalPayment = upkeepPerformedLog.args.totalPayment
            assert.equal(id.toString(), upkeepId.toString())
            assert.equal(success, true)
            assert.equal(
              trigger,
              encodeBlockTrigger({
                blockNum: checkBlock.number,
                blockHash: checkBlock.hash,
              }),
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
            const registryLinkAfter = await linkToken.balanceOf(
              registry.address,
            )
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
              registrationAfter.lastPerformedBlockNumber.toString(),
              tx.blockNumber?.toString(),
            )

            // Latest epoch should be 5
            assert.equal((await registry.getState()).state.latestEpoch, 5)
          }
        },
      )

      // describe.only('Gas benchmarking conditional upkeeps [ @skip-coverage ]', function () {
      //   const fs = [1]
      //   fs.forEach(function (newF) {
      //     it(
      //       'When f=' +
      //         newF +
      //         ' calculates gas overhead appropriately within a margin for different scenarios',
      //       async () => {
      //         // Perform the upkeep once to remove non-zero storage slots and have predictable gas measurement
      //         let tx = await getTransmitTx(registry, keeper1, [upkeepId])
      //         await tx.wait()
      //
      //         await registry
      //           .connect(admin)
      //           .setUpkeepGasLimit(upkeepId, performGas.mul(3))
      //
      //         // Different test scenarios
      //         let longBytes = '0x'
      //         for (let i = 0; i < maxPerformDataSize.toNumber(); i++) {
      //           longBytes += '11'
      //         }
      //         const upkeepSuccessArray = [true, false]
      //         const performGasArray = [5000, performGas]
      //         const performDataArray = ['0x', longBytes]
      //         const chainModuleOverheads = await moduleBase.getGasOverhead()
      //
      //         for (const i in upkeepSuccessArray) {
      //           for (const j in performGasArray) {
      //             for (const k in performDataArray) {
      //               const upkeepSuccess = upkeepSuccessArray[i]
      //               const performGas = performGasArray[j]
      //               const performData = performDataArray[k]
      //
      //               await mock.setCanPerform(upkeepSuccess)
      //               await mock.setPerformGasToBurn(performGas)
      //               await registry
      //                 .connect(owner)
      //                 .setConfigTypeSafe(
      //                   signerAddresses,
      //                   keeperAddresses,
      //                   newF,
      //                   config,
      //                   offchainVersion,
      //                   offchainBytes,
      //                   baseConfig[6],
      //                   baseConfig[7],
      //                 )
      //               tx = await getTransmitTx(registry, keeper1, [upkeepId], {
      //                 numSigners: newF + 1,
      //                 performDatas: [performData],
      //               })
      //               const receipt = await tx.wait()
      //               const upkeepPerformedLogs =
      //                 parseUpkeepPerformedLogs(receipt)
      //               // exactly 1 Upkeep Performed should be emitted
      //               assert.equal(upkeepPerformedLogs.length, 1)
      //               const upkeepPerformedLog = upkeepPerformedLogs[0]
      //
      //               const upkeepGasUsed = upkeepPerformedLog.args.gasUsed
      //               const chargedGasOverhead =
      //                 upkeepPerformedLog.args.gasOverhead
      //               const actualGasOverhead = receipt.gasUsed
      //                 .sub(upkeepGasUsed)
      //                 .add(500000) // the amount of pubdataGas used returned by mock gas bound caller
      //               const estimatedGasOverhead = registryConditionalOverhead
      //                 .add(
      //                   registryPerSignerGasOverhead.mul(
      //                     BigNumber.from(newF + 1),
      //                   ),
      //                 )
      //                 .add(chainModuleOverheads.chainModuleFixedOverhead)
      //                 .add(65_400)
      //
      //               assert.isTrue(upkeepGasUsed.gt(BigNumber.from('0')))
      //               assert.isTrue(chargedGasOverhead.gt(BigNumber.from('0')))
      //               assert.isTrue(actualGasOverhead.gt(BigNumber.from('0')))
      //
      //               console.log(
      //                 'Gas Benchmarking conditional upkeeps:',
      //                 'upkeepSuccess=',
      //                 upkeepSuccess,
      //                 'performGas=',
      //                 performGas.toString(),
      //                 'performData length=',
      //                 performData.length / 2 - 1,
      //                 'sig verification ( f =',
      //                 newF,
      //                 '): estimated overhead: ',
      //                 estimatedGasOverhead.toString(), // 179800
      //                 ' charged overhead: ',
      //                 chargedGasOverhead.toString(), // 180560
      //                 ' actual overhead: ',
      //                 actualGasOverhead.toString(), // 632949
      //                 ' calculation margin over gasUsed: ',
      //                 chargedGasOverhead.sub(actualGasOverhead).toString(), // 18456
      //                 ' estimation margin over gasUsed: ',
      //                 estimatedGasOverhead.sub(actualGasOverhead).toString(), // -27744
      //                 ' upkeepGasUsed: ',
      //                 upkeepGasUsed, // 988620
      //                 ' receipt.gasUsed: ',
      //                 receipt.gasUsed, // 1121569
      //               )
      //
      //               // The actual gas overhead should be less than charged gas overhead, but not by a lot
      //               // The charged gas overhead is controlled by ACCOUNTING_FIXED_GAS_OVERHEAD and
      //               // ACCOUNTING_PER_UPKEEP_GAS_OVERHEAD, and their correct values should be set to
      //               // satisfy constraints in multiple places
      //               assert.isTrue(
      //                 chargedGasOverhead.gt(actualGasOverhead),
      //                 'Gas overhead calculated is too low, increase account gas variables (ACCOUNTING_FIXED_GAS_OVERHEAD/ACCOUNTING_PER_UPKEEP_GAS_OVERHEAD) by at least ' +
      //                   actualGasOverhead.sub(chargedGasOverhead).toString(),
      //               )
      //               assert.isTrue(
      //                 chargedGasOverhead // 180560
      //                   .sub(actualGasOverhead) // 132940
      //                   .lt(gasCalculationMargin),
      //                 'Gas overhead calculated is too high, decrease account gas variables (ACCOUNTING_FIXED_GAS_OVERHEAD/ACCOUNTING_PER_SIGNER_GAS_OVERHEAD)  by at least ' +
      //                   chargedGasOverhead
      //                     .sub(actualGasOverhead)
      //                     .sub(gasCalculationMargin)
      //                     .toString(),
      //               )
      //
      //               // The estimated overhead during checkUpkeep should be close to the actual overhead in transaction
      //               // It should be greater than the actual overhead but not by a lot
      //               // The estimated overhead is controlled by variables
      //               // REGISTRY_CONDITIONAL_OVERHEAD, REGISTRY_LOG_OVERHEAD, REGISTRY_PER_SIGNER_GAS_OVERHEAD
      //               // REGISTRY_PER_PERFORM_BYTE_GAS_OVERHEAD
      //               assert.isTrue(
      //                 estimatedGasOverhead.gt(actualGasOverhead),
      //                 'Gas overhead estimated in check upkeep is too low, increase estimation gas variables (REGISTRY_CONDITIONAL_OVERHEAD/REGISTRY_LOG_OVERHEAD/REGISTRY_PER_SIGNER_GAS_OVERHEAD/REGISTRY_PER_PERFORM_BYTE_GAS_OVERHEAD) by at least ' +
      //                   estimatedGasOverhead.sub(chargedGasOverhead).toString(),
      //               )
      //               assert.isTrue(
      //                 estimatedGasOverhead
      //                   .sub(actualGasOverhead)
      //                   .lt(gasEstimationMargin),
      //                 'Gas overhead estimated is too high, decrease estimation gas variables (REGISTRY_CONDITIONAL_OVERHEAD/REGISTRY_LOG_OVERHEAD/REGISTRY_PER_SIGNER_GAS_OVERHEAD/REGISTRY_PER_PERFORM_BYTE_GAS_OVERHEAD)  by at least ' +
      //                   estimatedGasOverhead
      //                     .sub(actualGasOverhead)
      //                     .sub(gasEstimationMargin)
      //                     .toString(),
      //               )
      //             }
      //           }
      //         }
      //       },
      //     )
      //   })
      // })

      // describe.only('Gas benchmarking log upkeeps [ @skip-coverage ]', function () {
      //   const fs = [1]
      //   fs.forEach(function (newF) {
      //     it(
      //       'When f=' +
      //         newF +
      //         ' calculates gas overhead appropriately within a margin',
      //       async () => {
      //         // Perform the upkeep once to remove non-zero storage slots and have predictable gas measurement
      //         let tx = await getTransmitTx(registry, keeper1, [logUpkeepId])
      //         await tx.wait()
      //         const performData = '0x'
      //         await mock.setCanPerform(true)
      //         await mock.setPerformGasToBurn(performGas)
      //         await registry.setConfigTypeSafe(
      //           signerAddresses,
      //           keeperAddresses,
      //           newF,
      //           config,
      //           offchainVersion,
      //           offchainBytes,
      //           baseConfig[6],
      //           baseConfig[7],
      //         )
      //         tx = await getTransmitTx(registry, keeper1, [logUpkeepId], {
      //           numSigners: newF + 1,
      //           performDatas: [performData],
      //         })
      //         const receipt = await tx.wait()
      //         const upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
      //         // exactly 1 Upkeep Performed should be emitted
      //         assert.equal(upkeepPerformedLogs.length, 1)
      //         const upkeepPerformedLog = upkeepPerformedLogs[0]
      //         const chainModuleOverheads = await moduleBase.getGasOverhead()
      //
      //         const upkeepGasUsed = upkeepPerformedLog.args.gasUsed
      //         const chargedGasOverhead = upkeepPerformedLog.args.gasOverhead
      //         const actualGasOverhead = receipt.gasUsed
      //           .sub(upkeepGasUsed)
      //           .add(500000) // the amount of pubdataGas used returned by mock gas bound caller
      //         const estimatedGasOverhead = registryLogOverhead
      //           .add(registryPerSignerGasOverhead.mul(BigNumber.from(newF + 1)))
      //           .add(chainModuleOverheads.chainModuleFixedOverhead)
      //           .add(65_400)
      //
      //         assert.isTrue(upkeepGasUsed.gt(BigNumber.from('0')))
      //         assert.isTrue(chargedGasOverhead.gt(BigNumber.from('0')))
      //         assert.isTrue(actualGasOverhead.gt(BigNumber.from('0')))
      //
      //         console.log(
      //           'Gas Benchmarking log upkeeps:',
      //           'upkeepSuccess=',
      //           true,
      //           'performGas=',
      //           performGas.toString(),
      //           'performData length=',
      //           performData.length / 2 - 1,
      //           'sig verification ( f =',
      //           newF,
      //           '): estimated overhead: ',
      //           estimatedGasOverhead.toString(),
      //           ' charged overhead: ',
      //           chargedGasOverhead.toString(),
      //           ' actual overhead: ',
      //           actualGasOverhead.toString(),
      //           ' calculation margin over gasUsed: ',
      //           chargedGasOverhead.sub(actualGasOverhead).toString(),
      //           ' estimation margin over gasUsed: ',
      //           estimatedGasOverhead.sub(actualGasOverhead).toString(),
      //           ' upkeepGasUsed: ',
      //           upkeepGasUsed,
      //           ' receipt.gasUsed: ',
      //           receipt.gasUsed,
      //         )
      //
      //         assert.isTrue(
      //           chargedGasOverhead.gt(actualGasOverhead),
      //           'Gas overhead calculated is too low, increase account gas variables (ACCOUNTING_FIXED_GAS_OVERHEAD/ACCOUNTING_PER_UPKEEP_GAS_OVERHEAD) by at least ' +
      //             actualGasOverhead.sub(chargedGasOverhead).toString(),
      //         )
      //         assert.isTrue(
      //           chargedGasOverhead
      //             .sub(actualGasOverhead)
      //             .lt(gasCalculationMargin),
      //           'Gas overhead calculated is too high, decrease account gas variables (ACCOUNTING_FIXED_GAS_OVERHEAD/ACCOUNTING_PER_SIGNER_GAS_OVERHEAD)  by at least ' +
      //             chargedGasOverhead
      //               .sub(actualGasOverhead)
      //               .sub(gasCalculationMargin)
      //               .toString(),
      //         )
      //
      //         assert.isTrue(
      //           estimatedGasOverhead.gt(actualGasOverhead),
      //           'Gas overhead estimated in check upkeep is too low, increase estimation gas variables (REGISTRY_CONDITIONAL_OVERHEAD/REGISTRY_LOG_OVERHEAD/REGISTRY_PER_SIGNER_GAS_OVERHEAD/REGISTRY_PER_PERFORM_BYTE_GAS_OVERHEAD) by at least ' +
      //             estimatedGasOverhead.sub(chargedGasOverhead).toString(),
      //         )
      //         assert.isTrue(
      //           estimatedGasOverhead
      //             .sub(actualGasOverhead)
      //             .lt(gasEstimationMargin),
      //           'Gas overhead estimated is too high, decrease estimation gas variables (REGISTRY_CONDITIONAL_OVERHEAD/REGISTRY_LOG_OVERHEAD/REGISTRY_PER_SIGNER_GAS_OVERHEAD/REGISTRY_PER_PERFORM_BYTE_GAS_OVERHEAD)  by at least ' +
      //             estimatedGasOverhead
      //               .sub(actualGasOverhead)
      //               .sub(gasEstimationMargin)
      //               .toString(),
      //         )
      //       },
      //     )
      //   })
      // })
    })
  })

  describeMaybe(
    '#transmit with upkeep batches [ @skip-coverage ]',
    function () {
      const numPassingConditionalUpkeepsArray = [0, 1, 5]
      const numPassingLogUpkeepsArray = [0, 1, 5]
      const numFailingUpkeepsArray = [0, 3]

      for (let idx = 0; idx < numPassingConditionalUpkeepsArray.length; idx++) {
        for (let jdx = 0; jdx < numPassingLogUpkeepsArray.length; jdx++) {
          for (let kdx = 0; kdx < numFailingUpkeepsArray.length; kdx++) {
            const numPassingConditionalUpkeeps =
              numPassingConditionalUpkeepsArray[idx]
            const numPassingLogUpkeeps = numPassingLogUpkeepsArray[jdx]
            const numFailingUpkeeps = numFailingUpkeepsArray[kdx]
            if (
              numPassingConditionalUpkeeps == 0 &&
              numPassingLogUpkeeps == 0
            ) {
              continue
            }
            it(
              '[Conditional:' +
                numPassingConditionalUpkeeps +
                ',Log:' +
                numPassingLogUpkeeps +
                ',Failures:' +
                numFailingUpkeeps +
                '] performs successful upkeeps and does not charge failing upkeeps',
              async () => {
                const allUpkeeps = await getMultipleUpkeepsDeployedAndFunded(
                  numPassingConditionalUpkeeps,
                  numPassingLogUpkeeps,
                  numFailingUpkeeps,
                )
                const passingConditionalUpkeepIds =
                  allUpkeeps.passingConditionalUpkeepIds
                const passingLogUpkeepIds = allUpkeeps.passingLogUpkeepIds
                const failingUpkeepIds = allUpkeeps.failingUpkeepIds

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
                const registrationConditionalPassingBefore = await Promise.all(
                  passingConditionalUpkeepIds.map(async (id) => {
                    const reg = await registry.getUpkeep(BigNumber.from(id))
                    assert.equal(reg.lastPerformedBlockNumber.toString(), '0')
                    return reg
                  }),
                )
                const registrationLogPassingBefore = await Promise.all(
                  passingLogUpkeepIds.map(async (id) => {
                    const reg = await registry.getUpkeep(BigNumber.from(id))
                    assert.equal(reg.lastPerformedBlockNumber.toString(), '0')
                    return reg
                  }),
                )
                const registrationFailingBefore = await Promise.all(
                  failingUpkeepIds.map(async (id) => {
                    const reg = await registry.getUpkeep(BigNumber.from(id))
                    assert.equal(reg.lastPerformedBlockNumber.toString(), '0')
                    return reg
                  }),
                )

                // cancel upkeeps so they will fail in the transmit process
                // must call the cancel upkeep as the owner to avoid the CANCELLATION_DELAY
                for (let ldx = 0; ldx < failingUpkeepIds.length; ldx++) {
                  await registry
                    .connect(owner)
                    .cancelUpkeep(failingUpkeepIds[ldx])
                }

                const tx = await getTransmitTx(
                  registry,
                  keeper1,
                  passingConditionalUpkeepIds.concat(
                    passingLogUpkeepIds.concat(failingUpkeepIds),
                  ),
                )

                const receipt = await tx.wait()
                const upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
                // exactly numPassingUpkeeps Upkeep Performed should be emitted
                assert.equal(
                  upkeepPerformedLogs.length,
                  numPassingConditionalUpkeeps + numPassingLogUpkeeps,
                )
                const cancelledUpkeepReportLogs =
                  parseCancelledUpkeepReportLogs(receipt)
                // exactly numFailingUpkeeps Upkeep Performed should be emitted
                assert.equal(
                  cancelledUpkeepReportLogs.length,
                  numFailingUpkeeps,
                )

                const keeperAfter = await registry.getTransmitterInfo(
                  await keeper1.getAddress(),
                )
                const keeperLinkAfter = await linkToken.balanceOf(
                  await keeper1.getAddress(),
                )
                const registryLinkAfter = await linkToken.balanceOf(
                  registry.address,
                )
                const registrationConditionalPassingAfter = await Promise.all(
                  passingConditionalUpkeepIds.map(async (id) => {
                    return await registry.getUpkeep(BigNumber.from(id))
                  }),
                )
                const registrationLogPassingAfter = await Promise.all(
                  passingLogUpkeepIds.map(async (id) => {
                    return await registry.getUpkeep(BigNumber.from(id))
                  }),
                )
                const registrationFailingAfter = await Promise.all(
                  failingUpkeepIds.map(async (id) => {
                    return await registry.getUpkeep(BigNumber.from(id))
                  }),
                )
                const registryPremiumAfter = (await registry.getState()).state
                  .totalPremium
                const premium = registryPremiumAfter.sub(registryPremiumBefore)

                let netPayment = BigNumber.from('0')
                for (let i = 0; i < numPassingConditionalUpkeeps; i++) {
                  const id = upkeepPerformedLogs[i].args.id
                  const gasUsed = upkeepPerformedLogs[i].args.gasUsed
                  const gasOverhead = upkeepPerformedLogs[i].args.gasOverhead
                  const totalPayment = upkeepPerformedLogs[i].args.totalPayment

                  expect(id).to.equal(passingConditionalUpkeepIds[i])
                  assert.isTrue(gasUsed.gt(BigNumber.from('0')))
                  assert.isTrue(gasOverhead.gt(BigNumber.from('0')))
                  assert.isTrue(totalPayment.gt(BigNumber.from('0')))

                  // Balance should be deducted
                  assert.equal(
                    registrationConditionalPassingBefore[i].balance
                      .sub(totalPayment)
                      .toString(),
                    registrationConditionalPassingAfter[i].balance.toString(),
                  )

                  // Amount spent should be updated correctly
                  assert.equal(
                    registrationConditionalPassingAfter[i].amountSpent
                      .sub(totalPayment)
                      .toString(),
                    registrationConditionalPassingBefore[
                      i
                    ].amountSpent.toString(),
                  )

                  // Last perform block number should be updated
                  assert.equal(
                    registrationConditionalPassingAfter[
                      i
                    ].lastPerformedBlockNumber.toString(),
                    tx.blockNumber?.toString(),
                  )

                  netPayment = netPayment.add(totalPayment)
                }

                for (let i = 0; i < numPassingLogUpkeeps; i++) {
                  const id =
                    upkeepPerformedLogs[numPassingConditionalUpkeeps + i].args
                      .id
                  const gasUsed =
                    upkeepPerformedLogs[numPassingConditionalUpkeeps + i].args
                      .gasUsed
                  const gasOverhead =
                    upkeepPerformedLogs[numPassingConditionalUpkeeps + i].args
                      .gasOverhead
                  const totalPayment =
                    upkeepPerformedLogs[numPassingConditionalUpkeeps + i].args
                      .totalPayment

                  expect(id).to.equal(passingLogUpkeepIds[i])
                  assert.isTrue(gasUsed.gt(BigNumber.from('0')))
                  assert.isTrue(gasOverhead.gt(BigNumber.from('0')))
                  assert.isTrue(totalPayment.gt(BigNumber.from('0')))

                  // Balance should be deducted
                  assert.equal(
                    registrationLogPassingBefore[i].balance
                      .sub(totalPayment)
                      .toString(),
                    registrationLogPassingAfter[i].balance.toString(),
                  )

                  // Amount spent should be updated correctly
                  assert.equal(
                    registrationLogPassingAfter[i].amountSpent
                      .sub(totalPayment)
                      .toString(),
                    registrationLogPassingBefore[i].amountSpent.toString(),
                  )

                  // Last perform block number should not be updated for log triggers
                  assert.equal(
                    registrationLogPassingAfter[
                      i
                    ].lastPerformedBlockNumber.toString(),
                    '0',
                  )

                  netPayment = netPayment.add(totalPayment)
                }

                for (let i = 0; i < numFailingUpkeeps; i++) {
                  // CancelledUpkeep log should be emitted
                  const id = cancelledUpkeepReportLogs[i].args.id
                  expect(id).to.equal(failingUpkeepIds[i])

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
                    ].lastPerformedBlockNumber.toString(),
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
              },
            )

            it(
              '[Conditional:' +
                numPassingConditionalUpkeeps +
                ',Log' +
                numPassingLogUpkeeps +
                ',Failures:' +
                numFailingUpkeeps +
                '] splits gas overhead appropriately among performed upkeeps [ @skip-coverage ]',
              async () => {
                const allUpkeeps = await getMultipleUpkeepsDeployedAndFunded(
                  numPassingConditionalUpkeeps,
                  numPassingLogUpkeeps,
                  numFailingUpkeeps,
                )
                const passingConditionalUpkeepIds =
                  allUpkeeps.passingConditionalUpkeepIds
                const passingLogUpkeepIds = allUpkeeps.passingLogUpkeepIds
                const failingUpkeepIds = allUpkeeps.failingUpkeepIds

                // Perform the upkeeps once to remove non-zero storage slots and have predictable gas measurement
                let tx = await getTransmitTx(
                  registry,
                  keeper1,
                  passingConditionalUpkeepIds.concat(
                    passingLogUpkeepIds.concat(failingUpkeepIds),
                  ),
                )

                await tx.wait()

                // cancel upkeeps so they will fail in the transmit process
                // must call the cancel upkeep as the owner to avoid the CANCELLATION_DELAY
                for (let ldx = 0; ldx < failingUpkeepIds.length; ldx++) {
                  await registry
                    .connect(owner)
                    .cancelUpkeep(failingUpkeepIds[ldx])
                }

                // Do the actual thing

                tx = await getTransmitTx(
                  registry,
                  keeper1,
                  passingConditionalUpkeepIds.concat(
                    passingLogUpkeepIds.concat(failingUpkeepIds),
                  ),
                )

                const receipt = await tx.wait()
                const upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
                // exactly numPassingUpkeeps Upkeep Performed should be emitted
                assert.equal(
                  upkeepPerformedLogs.length,
                  numPassingConditionalUpkeeps + numPassingLogUpkeeps,
                )

                let netGasUsedPlusChargedOverhead = BigNumber.from('0')
                for (let i = 0; i < numPassingConditionalUpkeeps; i++) {
                  const gasUsed = upkeepPerformedLogs[i].args.gasUsed
                  const chargedGasOverhead =
                    upkeepPerformedLogs[i].args.gasOverhead

                  assert.isTrue(gasUsed.gt(BigNumber.from('0')))
                  assert.isTrue(chargedGasOverhead.gt(BigNumber.from('0')))

                  // Overhead should be same for every upkeep
                  assert.isTrue(
                    chargedGasOverhead.eq(
                      upkeepPerformedLogs[0].args.gasOverhead,
                    ),
                  )
                  netGasUsedPlusChargedOverhead = netGasUsedPlusChargedOverhead
                    .add(gasUsed)
                    .add(chargedGasOverhead)
                }

                for (let i = 0; i < numPassingLogUpkeeps; i++) {
                  const gasUsed =
                    upkeepPerformedLogs[numPassingConditionalUpkeeps + i].args
                      .gasUsed
                  const chargedGasOverhead =
                    upkeepPerformedLogs[numPassingConditionalUpkeeps + i].args
                      .gasOverhead

                  assert.isTrue(gasUsed.gt(BigNumber.from('0')))
                  assert.isTrue(chargedGasOverhead.gt(BigNumber.from('0')))

                  // Overhead should be same for every upkeep
                  assert.isTrue(
                    chargedGasOverhead.eq(
                      upkeepPerformedLogs[numPassingConditionalUpkeeps].args
                        .gasOverhead,
                    ),
                  )
                  netGasUsedPlusChargedOverhead = netGasUsedPlusChargedOverhead
                    .add(gasUsed)
                    .add(chargedGasOverhead)
                }

                console.log(
                  'Gas Benchmarking - batching (passedConditionalUpkeeps: ',
                  numPassingConditionalUpkeeps,
                  'passedLogUpkeeps:',
                  numPassingLogUpkeeps,
                  'failedUpkeeps:',
                  numFailingUpkeeps,
                  '): ',
                  numPassingConditionalUpkeeps > 0
                    ? 'charged conditional overhead'
                    : '',
                  numPassingConditionalUpkeeps > 0
                    ? upkeepPerformedLogs[0].args.gasOverhead.toString()
                    : '',
                  numPassingLogUpkeeps > 0 ? 'charged log overhead' : '',
                  numPassingLogUpkeeps > 0
                    ? upkeepPerformedLogs[
                        numPassingConditionalUpkeeps
                      ].args.gasOverhead.toString()
                    : '',
                  ' margin over gasUsed',
                  netGasUsedPlusChargedOverhead.sub(receipt.gasUsed).toString(),
                )

                // The total gas charged should be greater than tx gas
                assert.isTrue(
                  netGasUsedPlusChargedOverhead.gt(receipt.gasUsed),
                  'Charged gas overhead is too low for batch upkeeps, increase ACCOUNTING_PER_UPKEEP_GAS_OVERHEAD',
                )
              },
            )
          }
        }
      }

      it('has enough perform gas overhead for large batches [ @skip-coverage ]', async () => {
        const numUpkeeps = 20
        const upkeepIds: BigNumber[] = []
        let totalPerformGas = BigNumber.from('0')
        for (let i = 0; i < numUpkeeps; i++) {
          const mock = await upkeepMockFactory.deploy()
          const tx = await registry
            .connect(owner)
            .registerUpkeep(
              mock.address,
              performGas,
              await admin.getAddress(),
              Trigger.CONDITION,
              linkToken.address,
              '0x',
              '0x',
              '0x',
            )
          const testUpkeepId = await getUpkeepID(tx)
          upkeepIds.push(testUpkeepId)

          // Add funds to passing upkeeps
          await registry.connect(owner).addFunds(testUpkeepId, toWei('10'))

          await mock.setCanPerform(true)
          await mock.setPerformGasToBurn(performGas)

          totalPerformGas = totalPerformGas.add(performGas)
        }

        // Should revert with no overhead added
        await evmRevert(
          getTransmitTx(registry, keeper1, upkeepIds, {
            gasLimit: totalPerformGas,
          }),
        )
        // Should not revert with overhead added
        await getTransmitTx(registry, keeper1, upkeepIds, {
          gasLimit: totalPerformGas.add(transmitGasOverhead),
        })
      })
    },
  )

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
          performGas,
          await admin.getAddress(),
          Trigger.CONDITION,
          linkToken.address,
          '0x',
          '0x',
          '0x',
        )

      const id1 = await getUpkeepID(tx)
      await registry.connect(admin).addFunds(id1, toWei('5'))

      await getTransmitTx(registry, keeper1, [id1])
      await getTransmitTx(registry, keeper2, [id1])
      await getTransmitTx(registry, keeper3, [id1])

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
          performGas,
          await admin.getAddress(),
          Trigger.CONDITION,
          linkToken.address,
          '0x',
          '0x',
          '0x',
        )
      const id2 = await getUpkeepID(tx2)
      await registry.connect(admin).addFunds(id2, toWei('5'))

      await getTransmitTx(registry, keeper1, [id2])
      await getTransmitTx(registry, keeper2, [id2])
      await getTransmitTx(registry, keeper3, [id2])

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
        .callStatic['checkUpkeep(uint256)'](upkeepId)

      assert.equal(checkUpkeepResult.upkeepNeeded, false)
      assert.equal(
        checkUpkeepResult.upkeepFailureReason,
        UpkeepFailureReason.INSUFFICIENT_BALANCE,
      )

      await registry.connect(admin).addFunds(upkeepId, oneWei)
      checkUpkeepResult = await registry
        .connect(zeroAddress)
        .callStatic['checkUpkeep(uint256)'](upkeepId)
      assert.equal(checkUpkeepResult.upkeepNeeded, true)
    })

    it('uses maxPerformData size in checkUpkeep but actual performDataSize in transmit', async () => {
      const tx = await registry
        .connect(owner)
        .registerUpkeep(
          mock.address,
          performGas,
          await admin.getAddress(),
          Trigger.CONDITION,
          linkToken.address,
          '0x',
          '0x',
          '0x',
        )
      const upkeepID = await getUpkeepID(tx)
      await mock.setCanCheck(true)
      await mock.setCanPerform(true)

      // upkeep is underfunded by 1 wei
      const minBalance1 = (await registry.getMinBalanceForUpkeep(upkeepID)).sub(
        1,
      )
      await registry.connect(owner).addFunds(upkeepID, minBalance1)

      // upkeep check should return false, 2 should return true
      const checkUpkeepResult = await registry
        .connect(zeroAddress)
        .callStatic['checkUpkeep(uint256)'](upkeepID)
      assert.equal(checkUpkeepResult.upkeepNeeded, false)
      assert.equal(
        checkUpkeepResult.upkeepFailureReason,
        UpkeepFailureReason.INSUFFICIENT_BALANCE,
      )

      // however upkeep should perform and pay all the remaining balance
      let maxPerformData = '0x'
      for (let i = 0; i < maxPerformDataSize.toNumber(); i++) {
        maxPerformData += '11'
      }

      const tx2 = await getTransmitTx(registry, keeper1, [upkeepID], {
        gasPrice: gasWei.mul(gasCeilingMultiplier),
        performDatas: [maxPerformData],
      })

      const receipt = await tx2.wait()
      const upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
      assert.equal(upkeepPerformedLogs.length, 1)
    })
  })

  describe('#withdrawFunds', () => {
    let upkeepId2: BigNumber

    beforeEach(async () => {
      const tx = await registry
        .connect(owner)
        .registerUpkeep(
          mock.address,
          performGas,
          await admin.getAddress(),
          Trigger.CONDITION,
          linkToken.address,
          '0x',
          '0x',
          '0x',
        )
      upkeepId2 = await getUpkeepID(tx)

      await registry.connect(admin).addFunds(upkeepId, toWei('100'))
      await registry.connect(admin).addFunds(upkeepId2, toWei('100'))

      // Do a perform so that upkeep is charged some amount
      await getTransmitTx(registry, keeper1, [upkeepId])
      await getTransmitTx(registry, keeper1, [upkeepId2])
    })

    describe('after the registration is paused, then cancelled', () => {
      it('allows the admin to withdraw', async () => {
        const balance = await registry.getBalance(upkeepId)
        const payee = await payee1.getAddress()
        await registry.connect(admin).pauseUpkeep(upkeepId)
        await registry.connect(owner).cancelUpkeep(upkeepId)
        await expect(() =>
          registry.connect(admin).withdrawFunds(upkeepId, payee),
        ).to.changeTokenBalance(linkToken, payee1, balance)
      })
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
        assert.equal(registration.balance.toNumber(), 0)
      })
    })
  })

  describe('#simulatePerformUpkeep', () => {
    it('reverts if called by non zero address', async () => {
      await evmRevertCustomError(
        registry
          .connect(await owner.getAddress())
          .callStatic.simulatePerformUpkeep(upkeepId, '0x'),
        registry,
        'OnlySimulatedBackend',
      )
    })

    it('reverts when registry is paused', async () => {
      await registry.connect(owner).pause()
      await evmRevertCustomError(
        registry
          .connect(zeroAddress)
          .callStatic.simulatePerformUpkeep(upkeepId, '0x'),
        registry,
        'RegistryPaused',
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

    it('returns true, gasUsed, and performGas when perform succeeds', async () => {
      await mock.setCanPerform(true)

      const simulatePerformResult = await registry
        .connect(zeroAddress)
        .callStatic.simulatePerformUpkeep(upkeepId, '0x')

      assert.equal(simulatePerformResult.success, true)
      assert.isTrue(simulatePerformResult.gasUsed.gt(BigNumber.from('0'))) // Some gas should be used
    })

    it('returns correct amount of gasUsed when perform succeeds', async () => {
      await mock.setCanPerform(true)
      await mock.setPerformGasToBurn(performGas) // 1,000,000

      // increase upkeep gas limit because the mock gas bound caller will always return 500,000 as the L1 gas used
      // that brings the total gas used to about 1M + 0.5M = 1.5M
      await registry
        .connect(admin)
        .setUpkeepGasLimit(upkeepId, BigNumber.from(2000000))

      const simulatePerformResult = await registry
        .connect(zeroAddress)
        .callStatic.simulatePerformUpkeep(upkeepId, '0x')

      // Full execute gas should be used, with some performGasBuffer(1000)
      assert.isTrue(
        simulatePerformResult.gasUsed.gt(
          performGas.add(pubdataGas).sub(BigNumber.from('1000')),
        ),
      )
    })
  })

  describe('#checkUpkeep', () => {
    it('reverts if called by non zero address', async () => {
      await evmRevertCustomError(
        registry
          .connect(await owner.getAddress())
          .callStatic['checkUpkeep(uint256)'](upkeepId),
        registry,
        'OnlySimulatedBackend',
      )
    })

    it('returns false and error code if the upkeep is cancelled by admin', async () => {
      await registry.connect(admin).cancelUpkeep(upkeepId)

      const checkUpkeepResult = await registry
        .connect(zeroAddress)
        .callStatic['checkUpkeep(uint256)'](upkeepId)

      assert.equal(checkUpkeepResult.upkeepNeeded, false)
      assert.equal(checkUpkeepResult.performData, '0x')
      assert.equal(
        checkUpkeepResult.upkeepFailureReason,
        UpkeepFailureReason.UPKEEP_CANCELLED,
      )
      expect(checkUpkeepResult.gasUsed).to.equal(0)
      expect(checkUpkeepResult.gasLimit).to.equal(performGas)
    })

    it('returns false and error code if the upkeep is cancelled by owner', async () => {
      await registry.connect(owner).cancelUpkeep(upkeepId)

      const checkUpkeepResult = await registry
        .connect(zeroAddress)
        .callStatic['checkUpkeep(uint256)'](upkeepId)

      assert.equal(checkUpkeepResult.upkeepNeeded, false)
      assert.equal(checkUpkeepResult.performData, '0x')
      assert.equal(
        checkUpkeepResult.upkeepFailureReason,
        UpkeepFailureReason.UPKEEP_CANCELLED,
      )
      expect(checkUpkeepResult.gasUsed).to.equal(0)
      expect(checkUpkeepResult.gasLimit).to.equal(performGas)
    })

    it('returns false and error code if the registry is paused', async () => {
      await registry.connect(owner).pause()

      const checkUpkeepResult = await registry
        .connect(zeroAddress)
        .callStatic['checkUpkeep(uint256)'](upkeepId)

      assert.equal(checkUpkeepResult.upkeepNeeded, false)
      assert.equal(checkUpkeepResult.performData, '0x')
      assert.equal(
        checkUpkeepResult.upkeepFailureReason,
        UpkeepFailureReason.REGISTRY_PAUSED,
      )
      expect(checkUpkeepResult.gasUsed).to.equal(0)
      expect(checkUpkeepResult.gasLimit).to.equal(performGas)
    })

    it('returns false and error code if the upkeep is paused', async () => {
      await registry.connect(admin).pauseUpkeep(upkeepId)

      const checkUpkeepResult = await registry
        .connect(zeroAddress)
        .callStatic['checkUpkeep(uint256)'](upkeepId)

      assert.equal(checkUpkeepResult.upkeepNeeded, false)
      assert.equal(checkUpkeepResult.performData, '0x')
      assert.equal(
        checkUpkeepResult.upkeepFailureReason,
        UpkeepFailureReason.UPKEEP_PAUSED,
      )
      expect(checkUpkeepResult.gasUsed).to.equal(0)
      expect(checkUpkeepResult.gasLimit).to.equal(performGas)
    })

    it('returns false and error code if user is out of funds', async () => {
      const checkUpkeepResult = await registry
        .connect(zeroAddress)
        .callStatic['checkUpkeep(uint256)'](upkeepId)

      assert.equal(checkUpkeepResult.upkeepNeeded, false)
      assert.equal(checkUpkeepResult.performData, '0x')
      assert.equal(
        checkUpkeepResult.upkeepFailureReason,
        UpkeepFailureReason.INSUFFICIENT_BALANCE,
      )
      expect(checkUpkeepResult.gasUsed).to.equal(0)
      expect(checkUpkeepResult.gasLimit).to.equal(performGas)
    })

    context('when the registration is funded', () => {
      beforeEach(async () => {
        await linkToken.connect(admin).approve(registry.address, toWei('200'))
        await registry.connect(admin).addFunds(upkeepId, toWei('100'))
        await registry.connect(admin).addFunds(logUpkeepId, toWei('100'))
      })

      it('returns false, error code, and revert data if the target check reverts', async () => {
        await mock.setShouldRevertCheck(true)
        await mock.setCheckRevertReason(
          'custom revert error, clever way to insert offchain data',
        )
        const checkUpkeepResult = await registry
          .connect(zeroAddress)
          .callStatic['checkUpkeep(uint256)'](upkeepId)
        assert.equal(checkUpkeepResult.upkeepNeeded, false)

        const revertReasonBytes = `0x${checkUpkeepResult.performData.slice(10)}` // remove sighash
        assert.equal(
          ethers.utils.defaultAbiCoder.decode(['string'], revertReasonBytes)[0],
          'custom revert error, clever way to insert offchain data',
        )
        assert.equal(
          checkUpkeepResult.upkeepFailureReason,
          UpkeepFailureReason.TARGET_CHECK_REVERTED,
        )
        assert.isTrue(checkUpkeepResult.gasUsed.gt(BigNumber.from('0'))) // Some gas should be used
        expect(checkUpkeepResult.gasLimit).to.equal(performGas)
        // Feed data should be returned here
        assert.isTrue(checkUpkeepResult.fastGasWei.gt(BigNumber.from('0')))
        assert.isTrue(checkUpkeepResult.linkUSD.gt(BigNumber.from('0')))
      })

      it('returns false, error code, and no revert data if the target check revert data exceeds maxRevertDataSize', async () => {
        await mock.setShouldRevertCheck(true)
        let longRevertReason = ''
        for (let i = 0; i <= maxRevertDataSize.toNumber(); i++) {
          longRevertReason += 'x'
        }
        await mock.setCheckRevertReason(longRevertReason)
        const checkUpkeepResult = await registry
          .connect(zeroAddress)
          .callStatic['checkUpkeep(uint256)'](upkeepId)
        assert.equal(checkUpkeepResult.upkeepNeeded, false)

        assert.equal(checkUpkeepResult.performData, '0x')
        assert.equal(
          checkUpkeepResult.upkeepFailureReason,
          UpkeepFailureReason.REVERT_DATA_EXCEEDS_LIMIT,
        )
        assert.isTrue(checkUpkeepResult.gasUsed.gt(BigNumber.from('0'))) // Some gas should be used
        expect(checkUpkeepResult.gasLimit).to.equal(performGas)
      })

      it('returns false and error code if the upkeep is not needed', async () => {
        await mock.setCanCheck(false)
        const checkUpkeepResult = await registry
          .connect(zeroAddress)
          .callStatic['checkUpkeep(uint256)'](upkeepId)

        assert.equal(checkUpkeepResult.upkeepNeeded, false)
        assert.equal(checkUpkeepResult.performData, '0x')
        assert.equal(
          checkUpkeepResult.upkeepFailureReason,
          UpkeepFailureReason.UPKEEP_NOT_NEEDED,
        )
        assert.isTrue(checkUpkeepResult.gasUsed.gt(BigNumber.from('0'))) // Some gas should be used
        expect(checkUpkeepResult.gasLimit).to.equal(performGas)
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
          .callStatic['checkUpkeep(uint256)'](upkeepId)

        assert.equal(checkUpkeepResult.upkeepNeeded, false)
        assert.equal(checkUpkeepResult.performData, '0x')
        assert.equal(
          checkUpkeepResult.upkeepFailureReason,
          UpkeepFailureReason.PERFORM_DATA_EXCEEDS_LIMIT,
        )
        assert.isTrue(checkUpkeepResult.gasUsed.gt(BigNumber.from('0'))) // Some gas should be used
        expect(checkUpkeepResult.gasLimit).to.equal(performGas)
      })

      it('returns true with gas used if the target can execute', async () => {
        await mock.setCanCheck(true)
        await mock.setPerformData(randomBytes)

        const latestBlock = await ethers.provider.getBlock('latest')

        const checkUpkeepResult = await registry
          .connect(zeroAddress)
          .callStatic['checkUpkeep(uint256)'](upkeepId, {
            blockTag: latestBlock.number,
          })

        assert.equal(checkUpkeepResult.upkeepNeeded, true)
        assert.equal(checkUpkeepResult.performData, randomBytes)
        assert.equal(
          checkUpkeepResult.upkeepFailureReason,
          UpkeepFailureReason.NONE,
        )
        assert.isTrue(checkUpkeepResult.gasUsed.gt(BigNumber.from('0'))) // Some gas should be used
        expect(checkUpkeepResult.gasLimit).to.equal(performGas)
        assert.isTrue(checkUpkeepResult.fastGasWei.eq(gasWei))
        assert.isTrue(checkUpkeepResult.linkUSD.eq(linkUSD))
      })

      it('calls checkLog for log-trigger upkeeps', async () => {
        const log: Log = {
          index: 0,
          timestamp: 0,
          txHash: ethers.utils.randomBytes(32),
          blockNumber: 100,
          blockHash: ethers.utils.randomBytes(32),
          source: randomAddress(),
          topics: [ethers.utils.randomBytes(32), ethers.utils.randomBytes(32)],
          data: ethers.utils.randomBytes(1000),
        }

        await ltUpkeep.mock.checkLog.withArgs(log, '0x').returns(true, '0x1234')

        const checkData = encodeLog(log)

        const checkUpkeepResult = await registry
          .connect(zeroAddress)
          .callStatic['checkUpkeep(uint256,bytes)'](logUpkeepId, checkData)

        expect(checkUpkeepResult.upkeepNeeded).to.be.true
        expect(checkUpkeepResult.performData).to.equal('0x1234')
      })

      itMaybe(
        'has a large enough gas overhead to cover upkeeps that use all their gas [ @skip-coverage ]',
        async () => {
          await mock.setCanCheck(true)
          await mock.setCheckGasToBurn(checkGasLimit)
          const gas = checkGasLimit.add(checkGasOverhead)
          const checkUpkeepResult = await registry
            .connect(zeroAddress)
            .callStatic['checkUpkeep(uint256)'](upkeepId, {
              gasLimit: gas,
            })

          assert.equal(checkUpkeepResult.upkeepNeeded, true)
        },
      )
    })
  })

  describe('#getMaxPaymentForGas', () => {
    itMaybe('calculates the max fee appropriately in ZKSync', async () => {
      await verifyMaxPayment(registry, moduleBase)
    })

    it('uses the fallback gas price if the feed has issues in ZKSync', async () => {
      const chainModuleOverheads = await moduleBase.getGasOverhead()
      const expectedFallbackMaxPayment = linkForGas(
        performGas,
        registryConditionalOverhead
          .add(registryPerSignerGasOverhead.mul(f + 1))
          .add(chainModuleOverheads.chainModuleFixedOverhead),
        gasCeilingMultiplier.mul('2'), // fallbackGasPrice is 2x gas price
        paymentPremiumPPB,
        flatFeeMilliCents,
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
        (
          await registry.getMaxPaymentForGas(
            upkeepId,
            Trigger.CONDITION,
            performGas,
            linkToken.address,
          )
        ).toString(),
      )

      // Negative feed price
      roundId = 100
      updatedAt = now()
      startedAt = 946684799
      await gasPriceFeed
        .connect(owner)
        .updateRoundData(roundId, -100, updatedAt, startedAt)

      assert.equal(
        expectedFallbackMaxPayment.toString(),
        (
          await registry.getMaxPaymentForGas(
            upkeepId,
            Trigger.CONDITION,
            performGas,
            linkToken.address,
          )
        ).toString(),
      )

      // Zero feed price
      roundId = 101
      updatedAt = now()
      startedAt = 946684799
      await gasPriceFeed
        .connect(owner)
        .updateRoundData(roundId, 0, updatedAt, startedAt)

      assert.equal(
        expectedFallbackMaxPayment.toString(),
        (
          await registry.getMaxPaymentForGas(
            upkeepId,
            Trigger.CONDITION,
            performGas,
            linkToken.address,
          )
        ).toString(),
      )
    })

    it('uses the fallback link price if the feed has issues in ZKSync', async () => {
      const chainModuleOverheads = await moduleBase.getGasOverhead()
      const expectedFallbackMaxPayment = linkForGas(
        performGas,
        registryConditionalOverhead
          .add(registryPerSignerGasOverhead.mul(f + 1))
          .add(chainModuleOverheads.chainModuleFixedOverhead),
        gasCeilingMultiplier.mul('2'), // fallbackLinkPrice is 1/2 link price, so multiply by 2
        paymentPremiumPPB,
        flatFeeMilliCents,
      ).total

      // Stale feed
      let roundId = 99
      const answer = 100
      let updatedAt = 946684800 // New Years 2000 🥳
      let startedAt = 946684799
      await linkUSDFeed
        .connect(owner)
        .updateRoundData(roundId, answer, updatedAt, startedAt)

      assert.equal(
        expectedFallbackMaxPayment.toString(),
        (
          await registry.getMaxPaymentForGas(
            upkeepId,
            Trigger.CONDITION,
            performGas,
            linkToken.address,
          )
        ).toString(),
      )

      // Negative feed price
      roundId = 100
      updatedAt = now()
      startedAt = 946684799
      await linkUSDFeed
        .connect(owner)
        .updateRoundData(roundId, -100, updatedAt, startedAt)

      assert.equal(
        expectedFallbackMaxPayment.toString(),
        (
          await registry.getMaxPaymentForGas(
            upkeepId,
            Trigger.CONDITION,
            performGas,
            linkToken.address,
          )
        ).toString(),
      )

      // Zero feed price
      roundId = 101
      updatedAt = now()
      startedAt = 946684799
      await linkUSDFeed
        .connect(owner)
        .updateRoundData(roundId, 0, updatedAt, startedAt)

      assert.equal(
        expectedFallbackMaxPayment.toString(),
        (
          await registry.getMaxPaymentForGas(
            upkeepId,
            Trigger.CONDITION,
            performGas,
            linkToken.address,
          )
        ).toString(),
      )
    })
  })

  describe('#typeAndVersion', () => {
    it('uses the correct type and version', async () => {
      const typeAndVersion = await registry.typeAndVersion()
      assert.equal(typeAndVersion, 'AutomationRegistry 2.3.0')
    })
  })

  describeMaybe('#setConfig - onchain', async () => {
    const maxGas = BigNumber.from(6)
    const staleness = BigNumber.from(4)
    const ceiling = BigNumber.from(5)
    const newMaxCheckDataSize = BigNumber.from(10000)
    const newMaxPerformDataSize = BigNumber.from(10000)
    const newMaxRevertDataSize = BigNumber.from(10000)
    const newMaxPerformGas = BigNumber.from(10000000)
    const fbGasEth = BigNumber.from(7)
    const fbLinkEth = BigNumber.from(8)
    const fbNativeEth = BigNumber.from(100)
    const newTranscoder = randomAddress()
    const newRegistrars = [randomAddress(), randomAddress()]
    const upkeepManager = randomAddress()
    const financeAdminAddress = randomAddress()

    const newConfig: OnChainConfig = {
      checkGasLimit: maxGas,
      stalenessSeconds: staleness,
      gasCeilingMultiplier: ceiling,
      maxCheckDataSize: newMaxCheckDataSize,
      maxPerformDataSize: newMaxPerformDataSize,
      maxRevertDataSize: newMaxRevertDataSize,
      maxPerformGas: newMaxPerformGas,
      fallbackGasPrice: fbGasEth,
      fallbackLinkPrice: fbLinkEth,
      fallbackNativePrice: fbNativeEth,
      transcoder: newTranscoder,
      registrars: newRegistrars,
      upkeepPrivilegeManager: upkeepManager,
      chainModule: moduleBase.address,
      reorgProtectionEnabled: true,
      financeAdmin: financeAdminAddress,
    }

    it('reverts when called by anyone but the proposed owner', async () => {
      await evmRevert(
        registry
          .connect(payee1)
          .setConfigTypeSafe(
            signerAddresses,
            keeperAddresses,
            f,
            newConfig,
            offchainVersion,
            offchainBytes,
            baseConfig[6],
            baseConfig[7],
          ),
        'Only callable by owner',
      )
    })

    it('reverts if signers or transmitters are the zero address', async () => {
      await evmRevertCustomError(
        registry
          .connect(owner)
          .setConfigTypeSafe(
            [randomAddress(), randomAddress(), randomAddress(), zeroAddress],
            [
              randomAddress(),
              randomAddress(),
              randomAddress(),
              randomAddress(),
            ],
            f,
            newConfig,
            offchainVersion,
            offchainBytes,
            baseConfig[6],
            baseConfig[7],
          ),
        registry,
        'InvalidSigner',
      )

      await evmRevertCustomError(
        registry
          .connect(owner)
          .setConfigTypeSafe(
            [
              randomAddress(),
              randomAddress(),
              randomAddress(),
              randomAddress(),
            ],
            [randomAddress(), randomAddress(), randomAddress(), zeroAddress],
            f,
            newConfig,
            offchainVersion,
            offchainBytes,
            baseConfig[6],
            baseConfig[7],
          ),
        registry,
        'InvalidTransmitter',
      )
    })

    it('updates the onchainConfig and configDigest', async () => {
      const old = await registry.getState()
      const oldConfig = await registry.getConfig()
      const oldState = old.state
      assert.isTrue(stalenessSeconds.eq(oldConfig.stalenessSeconds))
      assert.isTrue(gasCeilingMultiplier.eq(oldConfig.gasCeilingMultiplier))

      await registry
        .connect(owner)
        .setConfigTypeSafe(
          signerAddresses,
          keeperAddresses,
          f,
          newConfig,
          offchainVersion,
          offchainBytes,
          [],
          [],
        )

      const updated = await registry.getState()
      const updatedConfig = updated.config
      const updatedState = updated.state
      assert.equal(updatedConfig.stalenessSeconds, staleness.toNumber())
      assert.equal(updatedConfig.gasCeilingMultiplier, ceiling.toNumber())
      assert.equal(
        updatedConfig.maxCheckDataSize,
        newMaxCheckDataSize.toNumber(),
      )
      assert.equal(
        updatedConfig.maxPerformDataSize,
        newMaxPerformDataSize.toNumber(),
      )
      assert.equal(
        updatedConfig.maxRevertDataSize,
        newMaxRevertDataSize.toNumber(),
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

      assert.equal(updatedConfig.transcoder, newTranscoder)
      assert.deepEqual(updatedConfig.registrars, newRegistrars)
      assert.equal(updatedConfig.upkeepPrivilegeManager, upkeepManager)
    })

    it('maintains paused state when config is changed', async () => {
      await registry.pause()
      const old = await registry.getState()
      assert.isTrue(old.state.paused)

      await registry
        .connect(owner)
        .setConfigTypeSafe(
          signerAddresses,
          keeperAddresses,
          f,
          newConfig,
          offchainVersion,
          offchainBytes,
          [],
          [],
        )

      const updated = await registry.getState()
      assert.isTrue(updated.state.paused)
    })

    it('emits an event', async () => {
      const tx = await registry
        .connect(owner)
        .setConfigTypeSafe(
          signerAddresses,
          keeperAddresses,
          f,
          newConfig,
          offchainVersion,
          offchainBytes,
          [],
          [],
        )
      await expect(tx).to.emit(registry, 'ConfigSet')
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
          .setConfigTypeSafe(
            newKeepers,
            newKeepers,
            f,
            config,
            offchainVersion,
            offchainBytes,
            baseConfig[6],
            baseConfig[7],
          ),
        'Only callable by owner',
      )
    })

    it('reverts if too many keeperAddresses set', async () => {
      for (let i = 0; i < 40; i++) {
        newKeepers.push(randomAddress())
      }
      await evmRevertCustomError(
        registry
          .connect(owner)
          .setConfigTypeSafe(
            newKeepers,
            newKeepers,
            f,
            config,
            offchainVersion,
            offchainBytes,
            baseConfig[6],
            baseConfig[7],
          ),
        registry,
        'TooManyOracles',
      )
    })

    it('reverts if f=0', async () => {
      await evmRevertCustomError(
        registry
          .connect(owner)
          .setConfigTypeSafe(
            newKeepers,
            newKeepers,
            0,
            config,
            offchainVersion,
            offchainBytes,
            baseConfig[6],
            baseConfig[7],
          ),
        registry,
        'IncorrectNumberOfFaultyOracles',
      )
    })

    it('reverts if signers != transmitters length', async () => {
      const signers = [randomAddress()]
      await evmRevertCustomError(
        registry
          .connect(owner)
          .setConfigTypeSafe(
            signers,
            newKeepers,
            f,
            config,
            offchainVersion,
            offchainBytes,
            baseConfig[6],
            baseConfig[7],
          ),
        registry,
        'IncorrectNumberOfSigners',
      )
    })

    it('reverts if signers <= 3f', async () => {
      newKeepers.pop()
      await evmRevertCustomError(
        registry
          .connect(owner)
          .setConfigTypeSafe(
            newKeepers,
            newKeepers,
            f,
            config,
            offchainVersion,
            offchainBytes,
            baseConfig[6],
            baseConfig[7],
          ),
        registry,
        'IncorrectNumberOfSigners',
      )
    })

    it('reverts on repeated signers', async () => {
      const newSigners = [
        await personas.Eddy.getAddress(),
        await personas.Eddy.getAddress(),
        await personas.Eddy.getAddress(),
        await personas.Eddy.getAddress(),
      ]
      await evmRevertCustomError(
        registry
          .connect(owner)
          .setConfigTypeSafe(
            newSigners,
            newKeepers,
            f,
            config,
            offchainVersion,
            offchainBytes,
            baseConfig[6],
            baseConfig[7],
          ),
        registry,
        'RepeatedSigner',
      )
    })

    it('reverts on repeated transmitters', async () => {
      const newTransmitters = [
        await personas.Eddy.getAddress(),
        await personas.Eddy.getAddress(),
        await personas.Eddy.getAddress(),
        await personas.Eddy.getAddress(),
      ]
      await evmRevertCustomError(
        registry
          .connect(owner)
          .setConfigTypeSafe(
            newKeepers,
            newTransmitters,
            f,
            config,
            offchainVersion,
            offchainBytes,
            baseConfig[6],
            baseConfig[7],
          ),
        registry,
        'RepeatedTransmitter',
      )
    })

    itMaybe('stores new config and emits event', async () => {
      // Perform an upkeep so that totalPremium is updated
      await registry.connect(admin).addFunds(upkeepId, toWei('100'))
      let tx = await getTransmitTx(registry, keeper1, [upkeepId])
      await tx.wait()

      const newOffChainVersion = BigNumber.from('2')
      const newOffChainConfig = '0x1122'

      const old = await registry.getState()
      const oldState = old.state
      assert(oldState.totalPremium.gt(BigNumber.from('0')))

      const newSigners = newKeepers
      tx = await registry
        .connect(owner)
        .setConfigTypeSafe(
          newSigners,
          newKeepers,
          f,
          config,
          newOffChainVersion,
          newOffChainConfig,
          [],
          [],
        )

      const updated = await registry.getState()
      const updatedState = updated.state
      assert(oldState.totalPremium.eq(updatedState.totalPremium))

      // Old signer addresses which are not in new signers should be non active
      for (let i = 0; i < signerAddresses.length; i++) {
        const signer = signerAddresses[i]
        if (!newSigners.includes(signer)) {
          assert(!(await registry.getSignerInfo(signer)).active)
          assert((await registry.getSignerInfo(signer)).index == 0)
        }
      }
      // New signer addresses should be active
      for (let i = 0; i < newSigners.length; i++) {
        const signer = newSigners[i]
        assert((await registry.getSignerInfo(signer)).active)
        assert((await registry.getSignerInfo(signer)).index == i)
      }
      // Old transmitter addresses which are not in new transmitter should be non active, update lastCollected but retain other info
      for (let i = 0; i < keeperAddresses.length; i++) {
        const transmitter = keeperAddresses[i]
        if (!newKeepers.includes(transmitter)) {
          assert(!(await registry.getTransmitterInfo(transmitter)).active)
          assert((await registry.getTransmitterInfo(transmitter)).index == i)
          assert(
            (await registry.getTransmitterInfo(transmitter)).lastCollected.eq(
              oldState.totalPremium.sub(
                oldState.totalPremium.mod(keeperAddresses.length),
              ),
            ),
          )
        }
      }
      // New transmitter addresses should be active
      for (let i = 0; i < newKeepers.length; i++) {
        const transmitter = newKeepers[i]
        assert((await registry.getTransmitterInfo(transmitter)).active)
        assert((await registry.getTransmitterInfo(transmitter)).index == i)
        assert(
          (await registry.getTransmitterInfo(transmitter)).lastCollected.eq(
            oldState.totalPremium,
          ),
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

  describe('#cancelUpkeep', () => {
    describe('when called by the admin', async () => {
      describeMaybe('when an upkeep has been performed', async () => {
        beforeEach(async () => {
          await linkToken.connect(owner).approve(registry.address, toWei('100'))
          await registry.connect(owner).addFunds(upkeepId, toWei('100'))
          await getTransmitTx(registry, keeper1, [upkeepId])
        })

        it('deducts a cancellation fee from the upkeep and adds to reserve', async () => {
          const newMinUpkeepSpend = toWei('10')
          const financeAdminAddress = await financeAdmin.getAddress()

          await registry.connect(owner).setConfigTypeSafe(
            signerAddresses,
            keeperAddresses,
            f,
            {
              checkGasLimit,
              stalenessSeconds,
              gasCeilingMultiplier,
              maxCheckDataSize,
              maxPerformDataSize,
              maxRevertDataSize,
              maxPerformGas,
              fallbackGasPrice,
              fallbackLinkPrice,
              fallbackNativePrice,
              transcoder: transcoder.address,
              registrars: [],
              upkeepPrivilegeManager: upkeepManager,
              chainModule: moduleBase.address,
              reorgProtectionEnabled: true,
              financeAdmin: financeAdminAddress,
            },
            offchainVersion,
            offchainBytes,
            [linkToken.address],
            [
              {
                gasFeePPB: paymentPremiumPPB,
                flatFeeMilliCents,
                priceFeed: linkUSDFeed.address,
                fallbackPrice: fallbackLinkPrice,
                minSpend: newMinUpkeepSpend,
                decimals: 18,
              },
            ],
          )

          const payee1Before = await linkToken.balanceOf(
            await payee1.getAddress(),
          )
          const upkeepBefore = (await registry.getUpkeep(upkeepId)).balance
          const ownerBefore = await registry.linkAvailableForPayment()

          const amountSpent = toWei('100').sub(upkeepBefore)
          const cancellationFee = newMinUpkeepSpend.sub(amountSpent)

          await registry.connect(admin).cancelUpkeep(upkeepId)

          const payee1After = await linkToken.balanceOf(
            await payee1.getAddress(),
          )
          const upkeepAfter = (await registry.getUpkeep(upkeepId)).balance
          const ownerAfter = await registry.linkAvailableForPayment()

          // post upkeep balance should be previous balance minus cancellation fee
          assert.isTrue(upkeepBefore.sub(cancellationFee).eq(upkeepAfter))
          // payee balance should not change
          assert.isTrue(payee1Before.eq(payee1After))
          // owner should receive the cancellation fee
          assert.isTrue(ownerAfter.sub(ownerBefore).eq(cancellationFee))
        })

        it('deducts up to balance as cancellation fee', async () => {
          // Very high min spend, should deduct whole balance as cancellation fees
          const newMinUpkeepSpend = toWei('1000')
          const financeAdminAddress = await financeAdmin.getAddress()

          await registry.connect(owner).setConfigTypeSafe(
            signerAddresses,
            keeperAddresses,
            f,
            {
              checkGasLimit,
              stalenessSeconds,
              gasCeilingMultiplier,
              maxCheckDataSize,
              maxPerformDataSize,
              maxRevertDataSize,
              maxPerformGas,
              fallbackGasPrice,
              fallbackLinkPrice,
              fallbackNativePrice,
              transcoder: transcoder.address,
              registrars: [],
              upkeepPrivilegeManager: upkeepManager,
              chainModule: moduleBase.address,
              reorgProtectionEnabled: true,
              financeAdmin: financeAdminAddress,
            },
            offchainVersion,
            offchainBytes,
            [linkToken.address],
            [
              {
                gasFeePPB: paymentPremiumPPB,
                flatFeeMilliCents,
                priceFeed: linkUSDFeed.address,
                fallbackPrice: fallbackLinkPrice,
                minSpend: newMinUpkeepSpend,
                decimals: 18,
              },
            ],
          )
          const payee1Before = await linkToken.balanceOf(
            await payee1.getAddress(),
          )
          const upkeepBefore = (await registry.getUpkeep(upkeepId)).balance
          const ownerBefore = await registry.linkAvailableForPayment()

          await registry.connect(admin).cancelUpkeep(upkeepId)
          const payee1After = await linkToken.balanceOf(
            await payee1.getAddress(),
          )
          const ownerAfter = await registry.linkAvailableForPayment()
          const upkeepAfter = (await registry.getUpkeep(upkeepId)).balance

          // all upkeep balance is deducted for cancellation fee
          assert.equal(upkeepAfter.toNumber(), 0)
          // payee balance should not change
          assert.isTrue(payee1After.eq(payee1Before))
          // all upkeep balance is transferred to the owner
          assert.isTrue(ownerAfter.sub(ownerBefore).eq(upkeepBefore))
        })

        it('does not deduct cancellation fee if more than minUpkeepSpendDollars is spent', async () => {
          // Very low min spend, already spent in one perform upkeep
          const newMinUpkeepSpend = BigNumber.from(420)
          const financeAdminAddress = await financeAdmin.getAddress()

          await registry.connect(owner).setConfigTypeSafe(
            signerAddresses,
            keeperAddresses,
            f,
            {
              checkGasLimit,
              stalenessSeconds,
              gasCeilingMultiplier,
              maxCheckDataSize,
              maxPerformDataSize,
              maxRevertDataSize,
              maxPerformGas,
              fallbackGasPrice,
              fallbackLinkPrice,
              fallbackNativePrice,
              transcoder: transcoder.address,
              registrars: [],
              upkeepPrivilegeManager: upkeepManager,
              chainModule: moduleBase.address,
              reorgProtectionEnabled: true,
              financeAdmin: financeAdminAddress,
            },
            offchainVersion,
            offchainBytes,
            [linkToken.address],
            [
              {
                gasFeePPB: paymentPremiumPPB,
                flatFeeMilliCents,
                priceFeed: linkUSDFeed.address,
                fallbackPrice: fallbackLinkPrice,
                minSpend: newMinUpkeepSpend,
                decimals: 18,
              },
            ],
          )
          const payee1Before = await linkToken.balanceOf(
            await payee1.getAddress(),
          )
          const upkeepBefore = (await registry.getUpkeep(upkeepId)).balance
          const ownerBefore = await registry.linkAvailableForPayment()

          await registry.connect(admin).cancelUpkeep(upkeepId)
          const payee1After = await linkToken.balanceOf(
            await payee1.getAddress(),
          )
          const ownerAfter = await registry.linkAvailableForPayment()
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
      await getTransmitTx(registry, keeper1, [upkeepId])
    })

    it('reverts if called by anyone but the payee', async () => {
      await evmRevertCustomError(
        registry
          .connect(payee2)
          .withdrawPayment(
            await keeper1.getAddress(),
            await nonkeeper.getAddress(),
          ),
        registry,
        'OnlyCallableByPayee',
      )
    })

    it('reverts if called with the 0 address', async () => {
      await evmRevertCustomError(
        registry
          .connect(payee2)
          .withdrawPayment(await keeper1.getAddress(), zeroAddress),
        registry,
        'InvalidRecipient',
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
      const ownerBefore = await registry.linkAvailableForPayment()

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
      const ownerAfter = await registry.linkAvailableForPayment()

      // registry total premium should not change
      assert.isTrue(registryPremiumBefore.eq(registryPremiumAfter))

      // Last collected should be updated to premium-change
      assert.isTrue(
        keeperAfter.lastCollected.eq(
          registryPremiumBefore.sub(
            registryPremiumBefore.mod(keeperAddresses.length),
          ),
        ),
      )

      // owner balance should remain unchanged
      assert.isTrue(ownerAfter.eq(ownerBefore))

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

  describe('#checkCallback', () => {
    it('returns false with appropriate failure reason when target callback reverts', async () => {
      await streamsLookupUpkeep.setShouldRevertCallback(true)

      const values: any[] = ['0x1234', '0xabcd']
      const res = await registry
        .connect(zeroAddress)
        .callStatic.checkCallback(streamsLookupUpkeepId, values, '0x')

      assert.isFalse(res.upkeepNeeded)
      assert.equal(res.performData, '0x')
      assert.equal(
        res.upkeepFailureReason,
        UpkeepFailureReason.CHECK_CALLBACK_REVERTED,
      )
      assert.isTrue(res.gasUsed.gt(BigNumber.from('0'))) // Some gas should be used
    })

    it('returns false with appropriate failure reason when target callback returns big performData', async () => {
      let longBytes = '0x'
      for (let i = 0; i <= maxPerformDataSize.toNumber(); i++) {
        longBytes += '11'
      }
      const values: any[] = [longBytes, longBytes]
      const res = await registry
        .connect(zeroAddress)
        .callStatic.checkCallback(streamsLookupUpkeepId, values, '0x')

      assert.isFalse(res.upkeepNeeded)
      assert.equal(res.performData, '0x')
      assert.equal(
        res.upkeepFailureReason,
        UpkeepFailureReason.PERFORM_DATA_EXCEEDS_LIMIT,
      )
      assert.isTrue(res.gasUsed.gt(BigNumber.from('0'))) // Some gas should be used
    })

    it('returns false with appropriate failure reason when target callback returns false', async () => {
      await streamsLookupUpkeep.setCallbackReturnBool(false)
      const values: any[] = ['0x1234', '0xabcd']
      const res = await registry
        .connect(zeroAddress)
        .callStatic.checkCallback(streamsLookupUpkeepId, values, '0x')

      assert.isFalse(res.upkeepNeeded)
      assert.equal(res.performData, '0x')
      assert.equal(
        res.upkeepFailureReason,
        UpkeepFailureReason.UPKEEP_NOT_NEEDED,
      )
      assert.isTrue(res.gasUsed.gt(BigNumber.from('0'))) // Some gas should be used
    })

    it('succeeds with upkeep needed', async () => {
      const values: any[] = ['0x1234', '0xabcd']

      const res = await registry
        .connect(zeroAddress)
        .callStatic.checkCallback(streamsLookupUpkeepId, values, '0x')
      const expectedPerformData = ethers.utils.defaultAbiCoder.encode(
        ['bytes[]', 'bytes'],
        [values, '0x'],
      )

      assert.isTrue(res.upkeepNeeded)
      assert.equal(res.performData, expectedPerformData)
      assert.equal(res.upkeepFailureReason, UpkeepFailureReason.NONE)
      assert.isTrue(res.gasUsed.gt(BigNumber.from('0'))) // Some gas should be used
    })
  })

  describe('transmitterPremiumSplit [ @skip-coverage ]', () => {
    beforeEach(async () => {
      await linkToken.connect(owner).approve(registry.address, toWei('100'))
      await registry.connect(owner).addFunds(upkeepId, toWei('100'))
    })

    it('splits premium evenly across transmitters', async () => {
      // Do a transmit from keeper1
      await getTransmitTx(registry, keeper1, [upkeepId])

      const registryPremium = (await registry.getState()).state.totalPremium
      assert.isTrue(registryPremium.gt(BigNumber.from(0)))

      const premiumPerTransmitter = registryPremium.div(
        BigNumber.from(keeperAddresses.length),
      )
      const k1Balance = (
        await registry.getTransmitterInfo(await keeper1.getAddress())
      ).balance
      // transmitter should be reimbursed for gas and get the premium
      assert.isTrue(k1Balance.gt(premiumPerTransmitter))
      const k1GasReimbursement = k1Balance.sub(premiumPerTransmitter)

      const k2Balance = (
        await registry.getTransmitterInfo(await keeper2.getAddress())
      ).balance
      // non transmitter should get its share of premium
      assert.isTrue(k2Balance.eq(premiumPerTransmitter))

      // Now do a transmit from keeper 2
      await getTransmitTx(registry, keeper2, [upkeepId])
      const registryPremiumNew = (await registry.getState()).state.totalPremium
      assert.isTrue(registryPremiumNew.gt(registryPremium))
      const premiumPerTransmitterNew = registryPremiumNew.div(
        BigNumber.from(keeperAddresses.length),
      )
      const additionalPremium = premiumPerTransmitterNew.sub(
        premiumPerTransmitter,
      )

      const k1BalanceNew = (
        await registry.getTransmitterInfo(await keeper1.getAddress())
      ).balance
      // k1 should get the new premium
      assert.isTrue(
        k1BalanceNew.eq(k1GasReimbursement.add(premiumPerTransmitterNew)),
      )

      const k2BalanceNew = (
        await registry.getTransmitterInfo(await keeper2.getAddress())
      ).balance
      // k2 should get gas reimbursement in addition to new premium
      assert.isTrue(k2BalanceNew.gt(k2Balance.add(additionalPremium)))
    })

    it('updates last collected upon payment withdrawn', async () => {
      // Do a transmit from keeper1
      await getTransmitTx(registry, keeper1, [upkeepId])

      const registryPremium = (await registry.getState()).state.totalPremium
      const k1 = await registry.getTransmitterInfo(await keeper1.getAddress())
      const k2 = await registry.getTransmitterInfo(await keeper2.getAddress())

      // Withdrawing for first time, last collected = 0
      assert.isTrue(k1.lastCollected.eq(BigNumber.from(0)))
      assert.isTrue(k2.lastCollected.eq(BigNumber.from(0)))

      //// Do the thing
      await registry
        .connect(payee1)
        .withdrawPayment(
          await keeper1.getAddress(),
          await nonkeeper.getAddress(),
        )

      const k1New = await registry.getTransmitterInfo(
        await keeper1.getAddress(),
      )
      const k2New = await registry.getTransmitterInfo(
        await keeper2.getAddress(),
      )

      // transmitter info lastCollected should be updated for k1, not for k2
      assert.isTrue(
        k1New.lastCollected.eq(
          registryPremium.sub(registryPremium.mod(keeperAddresses.length)),
        ),
      )
      assert.isTrue(k2New.lastCollected.eq(BigNumber.from(0)))
    })

    // itMaybe(
    it('maintains consistent balance information across all parties', async () => {
      // throughout transmits, withdrawals, setConfigs total claim on balances should remain less than expected balance
      // some spare change can get lost but it should be less than maxAllowedSpareChange

      let maxAllowedSpareChange = BigNumber.from('0')
      await verifyConsistentAccounting(maxAllowedSpareChange)

      await getTransmitTx(registry, keeper1, [upkeepId])
      maxAllowedSpareChange = maxAllowedSpareChange.add(BigNumber.from('31'))
      await verifyConsistentAccounting(maxAllowedSpareChange)

      await registry
        .connect(payee1)
        .withdrawPayment(
          await keeper1.getAddress(),
          await nonkeeper.getAddress(),
        )
      await verifyConsistentAccounting(maxAllowedSpareChange)

      await registry
        .connect(payee2)
        .withdrawPayment(
          await keeper2.getAddress(),
          await nonkeeper.getAddress(),
        )
      await verifyConsistentAccounting(maxAllowedSpareChange)

      await getTransmitTx(registry, keeper1, [upkeepId])
      maxAllowedSpareChange = maxAllowedSpareChange.add(BigNumber.from('31'))
      await verifyConsistentAccounting(maxAllowedSpareChange)

      await registry.connect(owner).setConfigTypeSafe(
        signerAddresses.slice(2, 15), // only use 2-14th index keepers
        keeperAddresses.slice(2, 15),
        f,
        config,
        offchainVersion,
        offchainBytes,
        baseConfig[6],
        baseConfig[7],
      )
      await verifyConsistentAccounting(maxAllowedSpareChange)

      await getTransmitTx(registry, keeper3, [upkeepId], {
        startingSignerIndex: 2,
      })
      maxAllowedSpareChange = maxAllowedSpareChange.add(BigNumber.from('13'))
      await verifyConsistentAccounting(maxAllowedSpareChange)

      await registry
        .connect(payee1)
        .withdrawPayment(
          await keeper1.getAddress(),
          await nonkeeper.getAddress(),
        )
      await verifyConsistentAccounting(maxAllowedSpareChange)

      await registry
        .connect(payee3)
        .withdrawPayment(
          await keeper3.getAddress(),
          await nonkeeper.getAddress(),
        )
      await verifyConsistentAccounting(maxAllowedSpareChange)

      await registry.connect(owner).setConfigTypeSafe(
        signerAddresses.slice(0, 4), // only use 0-3rd index keepers
        keeperAddresses.slice(0, 4),
        f,
        config,
        offchainVersion,
        offchainBytes,
        baseConfig[6],
        baseConfig[7],
      )
      await verifyConsistentAccounting(maxAllowedSpareChange)
      await getTransmitTx(registry, keeper1, [upkeepId])
      maxAllowedSpareChange = maxAllowedSpareChange.add(BigNumber.from('4'))
      await getTransmitTx(registry, keeper3, [upkeepId])
      maxAllowedSpareChange = maxAllowedSpareChange.add(BigNumber.from('4'))

      await verifyConsistentAccounting(maxAllowedSpareChange)
      await registry
        .connect(payee5)
        .withdrawPayment(
          await keeper5.getAddress(),
          await nonkeeper.getAddress(),
        )
      await verifyConsistentAccounting(maxAllowedSpareChange)

      await registry
        .connect(payee1)
        .withdrawPayment(
          await keeper1.getAddress(),
          await nonkeeper.getAddress(),
        )
      await verifyConsistentAccounting(maxAllowedSpareChange)
    })
  })
})
