import { ethers } from 'hardhat'
import { assert } from 'chai'
import { KeeperRegistry2_1__factory as KeeperRegistryFactory } from '../../../typechain/factories/KeeperRegistry2_1__factory'
import { KeeperRegistryLogicA2_1__factory as KeeperRegistryLogicAFactory } from '../../../typechain/factories/KeeperRegistryLogicA2_1__factory'
import { KeeperRegistryLogicB2_1__factory as KeeperRegistryLogicBFactory } from '../../../typechain/factories/KeeperRegistryLogicB2_1__factory'
import { AutomationForwarderLogic__factory as AutomationForwarderLogicFactory } from '../../../typechain/factories/AutomationForwarderLogic__factory'

// const describeMaybe = process.env.SKIP_SLOW ? describe.skip : describe
// const itMaybe = process.env.SKIP_SLOW ? it.skip : it

//////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////

/*********************************** REGISTRY v2.1 IS FROZEN ************************************/

// We are leaving the original tests enabled, however as 2.1 is still actively being deployed

describe('KeeperRegistry2_1 - Frozen [ @skip-coverage ]', () => {
  it('has not changed', () => {
    assert.equal(
      ethers.utils.id(KeeperRegistryFactory.bytecode),
      '0x05aaa1024d7400e9c4824dde093b96edf5888fa6e6be2c2fc4dca7ae47cc9de9',
      'KeeperRegistry bytecode has changed',
    )
    assert.equal(
      ethers.utils.id(KeeperRegistryLogicAFactory.bytecode),
      '0xdcc8805e88c550b2a25b972bee9f4e4c3649f01e26f8dda6b25d7a9c5da8ab2f',
      'KeeperRegistryLogicA bytecode has changed',
    )
    assert.equal(
      ethers.utils.id(KeeperRegistryLogicBFactory.bytecode),
      '0x891c26ba35b9b13afc9400fac5471d15842828ab717cbdc70ee263210c542563',
      'KeeperRegistryLogicB bytecode has changed',
    )
    assert.equal(
      ethers.utils.id(AutomationForwarderLogicFactory.bytecode),
      '0x6b89065111e9236407329fae3d68b33c311b7d3b6c2ae3dd15c1691a28b1aca7',
      'AutomationForwarderLogic bytecode has changed',
    )
  })
})

//////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////

// copied from AutomationRegistryInterface2_1.sol
// enum UpkeepFailureReason {
//   NONE,
//   UPKEEP_CANCELLED,
//   UPKEEP_PAUSED,
//   TARGET_CHECK_REVERTED,
//   UPKEEP_NOT_NEEDED,
//   PERFORM_DATA_EXCEEDS_LIMIT,
//   INSUFFICIENT_BALANCE,
//   CHECK_CALLBACK_REVERTED,
//   REVERT_DATA_EXCEEDS_LIMIT,
//   REGISTRY_PAUSED,
// }
//
// // copied from AutomationRegistryInterface2_1.sol
// enum Mode {
//   DEFAULT,
//   ARBITRUM,
//   OPTIMISM,
// }
//
// // copied from KeeperRegistryBase2_1.sol
// enum Trigger {
//   CONDITION,
//   LOG,
// }
//
// // un-exported types that must be extracted from the utils contract
// type Report = Parameters<AutomationUtils['_report']>[0]
// type OnChainConfig = Parameters<AutomationUtils['_onChainConfig']>[0]
// type LogTrigger = Parameters<AutomationUtils['_logTrigger']>[0]
// type ConditionalTrigger = Parameters<AutomationUtils['_conditionalTrigger']>[0]
// type Log = Parameters<AutomationUtils['_log']>[0]
//
// // -----------------------------------------------------------------------------------------------
//
// // These values should match the constants declared in registry
// let registryConditionalOverhead: BigNumber
// let registryLogOverhead: BigNumber
// let registryPerSignerGasOverhead: BigNumber
// let registryPerPerformByteGasOverhead: BigNumber
// let cancellationDelay: number
//
// // This is the margin for gas that we test for. Gas charged should always be greater
// // than total gas used in tx but should not increase beyond this margin
// const gasCalculationMargin = BigNumber.from(8000)
//
// const linkEth = BigNumber.from(5000000000000000) // 1 Link = 0.005 Eth
// const gasWei = BigNumber.from(1000000000) // 1 gwei
// // -----------------------------------------------------------------------------------------------
// // test-wide configs for upkeeps
// const linkDivisibility = BigNumber.from('1000000000000000000')
// const performGas = BigNumber.from('1000000')
// const paymentPremiumBase = BigNumber.from('1000000000')
// const paymentPremiumPPB = BigNumber.from('250000000')
// const flatFeeMicroLink = BigNumber.from(0)
//
// const randomBytes = '0x1234abcd'
// const emptyBytes = '0x'
// const emptyBytes32 =
//   '0x0000000000000000000000000000000000000000000000000000000000000000'
//
// const transmitGasOverhead = 1_000_000
// const checkGasOverhead = 400_000
//
// const stalenessSeconds = BigNumber.from(43820)
// const gasCeilingMultiplier = BigNumber.from(2)
// const checkGasLimit = BigNumber.from(10000000)
// const fallbackGasPrice = gasWei.mul(BigNumber.from('2'))
// const fallbackLinkPrice = linkEth.div(BigNumber.from('2'))
// const maxCheckDataSize = BigNumber.from(1000)
// const maxPerformDataSize = BigNumber.from(1000)
// const maxRevertDataSize = BigNumber.from(1000)
// const maxPerformGas = BigNumber.from(5000000)
// const minUpkeepSpend = BigNumber.from(0)
// const f = 1
// const offchainVersion = 1
// const offchainBytes = '0x'
// const zeroAddress = ethers.constants.AddressZero
// const epochAndRound5_1 =
//   '0x0000000000000000000000000000000000000000000000000000000000000501'
//
// let logTriggerConfig: string
//
// // -----------------------------------------------------------------------------------------------
//
// // Smart contract factories
// let linkTokenFactory: LinkTokenFactory
// let mockV3AggregatorFactory: MockV3AggregatorFactory
// let upkeepMockFactory: UpkeepMockFactory
// let upkeepAutoFunderFactory: UpkeepAutoFunderFactory
// let mockArbGasInfoFactory: MockArbGasInfoFactory
// let mockOVMGasPriceOracleFactory: MockOVMGasPriceOracleFactory
// let streamsLookupUpkeepFactory: StreamsLookupUpkeepFactory
// let personas: Personas
//
// // contracts
// let linkToken: LinkToken
// let linkEthFeed: MockV3Aggregator
// let gasPriceFeed: MockV3Aggregator
// let registry: IKeeperRegistry // default registry, used for most tests
// let arbRegistry: IKeeperRegistry // arbitrum registry
// let opRegistry: IKeeperRegistry // optimism registry
// let mgRegistry: IKeeperRegistry // "migrate registry" used in migration tests
// let blankRegistry: IKeeperRegistry // used to test initial configurations
// let mock: UpkeepMock
// let autoFunderUpkeep: UpkeepAutoFunder
// let ltUpkeep: MockContract
// let transcoder: UpkeepTranscoder
// let mockArbGasInfo: MockArbGasInfo
// let mockOVMGasPriceOracle: MockOVMGasPriceOracle
// let streamsLookupUpkeep: StreamsLookupUpkeep
// let automationUtils: AutomationUtils
//
// function now() {
//   return Math.floor(Date.now() / 1000)
// }
//
// async function getUpkeepID(tx: ContractTransaction): Promise<BigNumber> {
//   const receipt = await tx.wait()
//   for (const event of receipt.events || []) {
//     if (
//       event.args &&
//       event.eventSignature == 'UpkeepRegistered(uint256,uint32,address)'
//     ) {
//       return event.args[0]
//     }
//   }
//   throw new Error('could not find upkeep ID in tx event logs')
// }
//
// const getTriggerType = (upkeepId: BigNumber): Trigger => {
//   const hexBytes = ethers.utils.defaultAbiCoder.encode(['uint256'], [upkeepId])
//   const bytes = ethers.utils.arrayify(hexBytes)
//   for (let idx = 4; idx < 15; idx++) {
//     if (bytes[idx] != 0) {
//       return Trigger.CONDITION
//     }
//   }
//   return bytes[15] as Trigger
// }
//
// const encodeConfig = (onchainConfig: OnChainConfig) => {
//   return (
//     '0x' +
//     automationUtils.interface
//       .encodeFunctionData('_onChainConfig', [onchainConfig])
//       .slice(10)
//   )
// }
//
// const encodeBlockTrigger = (conditionalTrigger: ConditionalTrigger) => {
//   return (
//     '0x' +
//     automationUtils.interface
//       .encodeFunctionData('_conditionalTrigger', [conditionalTrigger])
//       .slice(10)
//   )
// }
//
// const encodeLogTrigger = (logTrigger: LogTrigger) => {
//   return (
//     '0x' +
//     automationUtils.interface
//       .encodeFunctionData('_logTrigger', [logTrigger])
//       .slice(10)
//   )
// }
//
// const encodeLog = (log: Log) => {
//   return (
//     '0x' + automationUtils.interface.encodeFunctionData('_log', [log]).slice(10)
//   )
// }
//
// const encodeReport = (report: Report) => {
//   return (
//     '0x' +
//     automationUtils.interface.encodeFunctionData('_report', [report]).slice(10)
//   )
// }
//
// type UpkeepData = {
//   Id: BigNumberish
//   performGas: BigNumberish
//   performData: BytesLike
//   trigger: BytesLike
// }
//
// const makeReport = (upkeeps: UpkeepData[]) => {
//   const upkeepIds = upkeeps.map((u) => u.Id)
//   const performGases = upkeeps.map((u) => u.performGas)
//   const triggers = upkeeps.map((u) => u.trigger)
//   const performDatas = upkeeps.map((u) => u.performData)
//   return encodeReport({
//     fastGasWei: gasWei,
//     linkNative: linkEth,
//     upkeepIds,
//     gasLimits: performGases,
//     triggers,
//     performDatas,
//   })
// }
//
// const makeLatestBlockReport = async (upkeepsIDs: BigNumberish[]) => {
//   const latestBlock = await ethers.provider.getBlock('latest')
//   const upkeeps: UpkeepData[] = []
//   for (let i = 0; i < upkeepsIDs.length; i++) {
//     upkeeps.push({
//       Id: upkeepsIDs[i],
//       performGas,
//       trigger: encodeBlockTrigger({
//         blockNum: latestBlock.number,
//         blockHash: latestBlock.hash,
//       }),
//       performData: '0x',
//     })
//   }
//   return makeReport(upkeeps)
// }
//
// const signReport = (
//   reportContext: string[],
//   report: any,
//   signers: Wallet[],
// ) => {
//   const reportDigest = ethers.utils.keccak256(report)
//   const packedArgs = ethers.utils.solidityPack(
//     ['bytes32', 'bytes32[3]'],
//     [reportDigest, reportContext],
//   )
//   const packedDigest = ethers.utils.keccak256(packedArgs)
//
//   const signatures = []
//   for (const signer of signers) {
//     signatures.push(signer._signingKey().signDigest(packedDigest))
//   }
//   const vs = signatures.map((i) => '0' + (i.v - 27).toString(16)).join('')
//   return {
//     vs: '0x' + vs.padEnd(64, '0'),
//     rs: signatures.map((i) => i.r),
//     ss: signatures.map((i) => i.s),
//   }
// }
//
// const parseUpkeepPerformedLogs = (receipt: ContractReceipt) => {
//   const parsedLogs = []
//   for (const rawLog of receipt.logs) {
//     try {
//       const log = registry.interface.parseLog(rawLog)
//       if (
//         log.name ==
//         registry.interface.events[
//           'UpkeepPerformed(uint256,bool,uint96,uint256,uint256,bytes)'
//         ].name
//       ) {
//         parsedLogs.push(log as unknown as UpkeepPerformedEvent)
//       }
//     } catch {
//       continue
//     }
//   }
//   return parsedLogs
// }
//
// const parseReorgedUpkeepReportLogs = (receipt: ContractReceipt) => {
//   const parsedLogs = []
//   for (const rawLog of receipt.logs) {
//     try {
//       const log = registry.interface.parseLog(rawLog)
//       if (
//         log.name ==
//         registry.interface.events['ReorgedUpkeepReport(uint256,bytes)'].name
//       ) {
//         parsedLogs.push(log as unknown as ReorgedUpkeepReportEvent)
//       }
//     } catch {
//       continue
//     }
//   }
//   return parsedLogs
// }
//
// const parseStaleUpkeepReportLogs = (receipt: ContractReceipt) => {
//   const parsedLogs = []
//   for (const rawLog of receipt.logs) {
//     try {
//       const log = registry.interface.parseLog(rawLog)
//       if (
//         log.name ==
//         registry.interface.events['StaleUpkeepReport(uint256,bytes)'].name
//       ) {
//         parsedLogs.push(log as unknown as StaleUpkeepReportEvent)
//       }
//     } catch {
//       continue
//     }
//   }
//   return parsedLogs
// }
//
// const parseInsufficientFundsUpkeepReportLogs = (receipt: ContractReceipt) => {
//   const parsedLogs = []
//   for (const rawLog of receipt.logs) {
//     try {
//       const log = registry.interface.parseLog(rawLog)
//       if (
//         log.name ==
//         registry.interface.events[
//           'InsufficientFundsUpkeepReport(uint256,bytes)'
//         ].name
//       ) {
//         parsedLogs.push(log as unknown as InsufficientFundsUpkeepReportEvent)
//       }
//     } catch {
//       continue
//     }
//   }
//   return parsedLogs
// }
//
// const parseCancelledUpkeepReportLogs = (receipt: ContractReceipt) => {
//   const parsedLogs = []
//   for (const rawLog of receipt.logs) {
//     try {
//       const log = registry.interface.parseLog(rawLog)
//       if (
//         log.name ==
//         registry.interface.events['CancelledUpkeepReport(uint256,bytes)'].name
//       ) {
//         parsedLogs.push(log as unknown as CancelledUpkeepReportEvent)
//       }
//     } catch {
//       continue
//     }
//   }
//   return parsedLogs
// }
//
// describe('KeeperRegistry2_1', () => {
//   let owner: Signer
//   let keeper1: Signer
//   let keeper2: Signer
//   let keeper3: Signer
//   let keeper4: Signer
//   let keeper5: Signer
//   let nonkeeper: Signer
//   let signer1: Wallet
//   let signer2: Wallet
//   let signer3: Wallet
//   let signer4: Wallet
//   let signer5: Wallet
//   let admin: Signer
//   let payee1: Signer
//   let payee2: Signer
//   let payee3: Signer
//   let payee4: Signer
//   let payee5: Signer
//
//   let upkeepId: BigNumber // conditional upkeep
//   let afUpkeepId: BigNumber // auto funding upkeep
//   let logUpkeepId: BigNumber // log trigger upkeepID
//   let streamsLookupUpkeepId: BigNumber // streams lookup upkeep
//   const numUpkeeps = 4 // see above
//   let keeperAddresses: string[]
//   let payees: string[]
//   let signers: Wallet[]
//   let signerAddresses: string[]
//   let config: any
//   let baseConfig: Parameters<IKeeperRegistry['setConfig']>
//   let upkeepManager: string
//
//   before(async () => {
//     personas = (await getUsers()).personas
//
//     const utilsFactory = await ethers.getContractFactory('AutomationUtils2_1')
//     automationUtils = await utilsFactory.deploy()
//
//     linkTokenFactory = await ethers.getContractFactory(
//       'src/v0.8/shared/test/helpers/LinkTokenTestHelper.sol:LinkTokenTestHelper',
//     )
//     // need full path because there are two contracts with name MockV3Aggregator
//     mockV3AggregatorFactory = (await ethers.getContractFactory(
//       'src/v0.8/tests/MockV3Aggregator.sol:MockV3Aggregator',
//     )) as unknown as MockV3AggregatorFactory
//     upkeepMockFactory = await ethers.getContractFactory('UpkeepMock')
//     upkeepAutoFunderFactory =
//       await ethers.getContractFactory('UpkeepAutoFunder')
//     mockArbGasInfoFactory = await ethers.getContractFactory('MockArbGasInfo')
//     mockOVMGasPriceOracleFactory = await ethers.getContractFactory(
//       'MockOVMGasPriceOracle',
//     )
//     streamsLookupUpkeepFactory = await ethers.getContractFactory(
//       'StreamsLookupUpkeep',
//     )
//
//     owner = personas.Default
//     keeper1 = personas.Carol
//     keeper2 = personas.Eddy
//     keeper3 = personas.Nancy
//     keeper4 = personas.Norbert
//     keeper5 = personas.Nick
//     nonkeeper = personas.Ned
//     admin = personas.Neil
//     payee1 = personas.Nelly
//     payee2 = personas.Norbert
//     payee3 = personas.Nick
//     payee4 = personas.Eddy
//     payee5 = personas.Carol
//     upkeepManager = await personas.Norbert.getAddress()
//     // signers
//     signer1 = new ethers.Wallet(
//       '0x7777777000000000000000000000000000000000000000000000000000000001',
//     )
//     signer2 = new ethers.Wallet(
//       '0x7777777000000000000000000000000000000000000000000000000000000002',
//     )
//     signer3 = new ethers.Wallet(
//       '0x7777777000000000000000000000000000000000000000000000000000000003',
//     )
//     signer4 = new ethers.Wallet(
//       '0x7777777000000000000000000000000000000000000000000000000000000004',
//     )
//     signer5 = new ethers.Wallet(
//       '0x7777777000000000000000000000000000000000000000000000000000000005',
//     )
//
//     keeperAddresses = [
//       await keeper1.getAddress(),
//       await keeper2.getAddress(),
//       await keeper3.getAddress(),
//       await keeper4.getAddress(),
//       await keeper5.getAddress(),
//     ]
//     payees = [
//       await payee1.getAddress(),
//       await payee2.getAddress(),
//       await payee3.getAddress(),
//       await payee4.getAddress(),
//       await payee5.getAddress(),
//     ]
//     signers = [signer1, signer2, signer3, signer4, signer5]
//
//     // We append 26 random addresses to keepers, payees and signers to get a system of 31 oracles
//     // This allows f value of 1 - 10
//     for (let i = 0; i < 26; i++) {
//       keeperAddresses.push(randomAddress())
//       payees.push(randomAddress())
//       signers.push(ethers.Wallet.createRandom())
//     }
//     signerAddresses = []
//     for (const signer of signers) {
//       signerAddresses.push(await signer.getAddress())
//     }
//
//     logTriggerConfig =
//       '0x' +
//       automationUtils.interface
//         .encodeFunctionData('_logTriggerConfig', [
//           {
//             contractAddress: randomAddress(),
//             filterSelector: 0,
//             topic0: ethers.utils.randomBytes(32),
//             topic1: ethers.utils.randomBytes(32),
//             topic2: ethers.utils.randomBytes(32),
//             topic3: ethers.utils.randomBytes(32),
//           },
//         ])
//         .slice(10)
//   })
//
//   const linkForGas = (
//     upkeepGasSpent: BigNumber,
//     gasOverhead: BigNumber,
//     gasMultiplier: BigNumber,
//     premiumPPB: BigNumber,
//     flatFee: BigNumber,
//     l1CostWei?: BigNumber,
//     numUpkeepsBatch?: BigNumber,
//   ) => {
//     l1CostWei = l1CostWei === undefined ? BigNumber.from(0) : l1CostWei
//     numUpkeepsBatch =
//       numUpkeepsBatch === undefined ? BigNumber.from(1) : numUpkeepsBatch
//
//     const gasSpent = gasOverhead.add(BigNumber.from(upkeepGasSpent))
//     const base = gasWei
//       .mul(gasMultiplier)
//       .mul(gasSpent)
//       .mul(linkDivisibility)
//       .div(linkEth)
//     const l1Fee = l1CostWei
//       .mul(gasMultiplier)
//       .div(numUpkeepsBatch)
//       .mul(linkDivisibility)
//       .div(linkEth)
//     const gasPayment = base.add(l1Fee)
//
//     const premium = gasWei
//       .mul(gasMultiplier)
//       .mul(upkeepGasSpent)
//       .add(l1CostWei.mul(gasMultiplier).div(numUpkeepsBatch))
//       .mul(linkDivisibility)
//       .div(linkEth)
//       .mul(premiumPPB)
//       .div(paymentPremiumBase)
//       .add(BigNumber.from(flatFee).mul('1000000000000'))
//
//     return {
//       total: gasPayment.add(premium),
//       gasPaymemnt: gasPayment,
//       premium,
//     }
//   }
//
//   const verifyMaxPayment = async (
//     registry: IKeeperRegistry,
//     l1CostWei?: BigNumber,
//   ) => {
//     type TestCase = {
//       name: string
//       multiplier: number
//       gas: number
//       premium: number
//       flatFee: number
//     }
//
//     const tests: TestCase[] = [
//       {
//         name: 'no fees',
//         multiplier: 1,
//         gas: 100000,
//         premium: 0,
//         flatFee: 0,
//       },
//       {
//         name: 'basic fees',
//         multiplier: 1,
//         gas: 100000,
//         premium: 250000000,
//         flatFee: 1000000,
//       },
//       {
//         name: 'max fees',
//         multiplier: 3,
//         gas: 10000000,
//         premium: 250000000,
//         flatFee: 1000000,
//       },
//     ]
//
//     const fPlusOne = BigNumber.from(f + 1)
//     const totalConditionalOverhead = registryConditionalOverhead
//       .add(registryPerSignerGasOverhead.mul(fPlusOne))
//       .add(registryPerPerformByteGasOverhead.mul(maxPerformDataSize))
//     const totalLogOverhead = registryLogOverhead
//       .add(registryPerSignerGasOverhead.mul(fPlusOne))
//       .add(registryPerPerformByteGasOverhead.mul(maxPerformDataSize))
//
//     for (const test of tests) {
//       await registry.connect(owner).setConfig(
//         signerAddresses,
//         keeperAddresses,
//         f,
//         encodeConfig({
//           paymentPremiumPPB: test.premium,
//           flatFeeMicroLink: test.flatFee,
//           checkGasLimit,
//           stalenessSeconds,
//           gasCeilingMultiplier: test.multiplier,
//           minUpkeepSpend,
//           maxCheckDataSize,
//           maxPerformDataSize,
//           maxRevertDataSize,
//           maxPerformGas,
//           fallbackGasPrice,
//           fallbackLinkPrice,
//           transcoder: transcoder.address,
//           registrars: [],
//           upkeepPrivilegeManager: upkeepManager,
//         }),
//         offchainVersion,
//         offchainBytes,
//       )
//
//       const conditionalPrice = await registry.getMaxPaymentForGas(
//         Trigger.CONDITION,
//         test.gas,
//       )
//       expect(conditionalPrice).to.equal(
//         linkForGas(
//           BigNumber.from(test.gas),
//           totalConditionalOverhead,
//           BigNumber.from(test.multiplier),
//           BigNumber.from(test.premium),
//           BigNumber.from(test.flatFee),
//           l1CostWei,
//         ).total,
//       )
//
//       const logPrice = await registry.getMaxPaymentForGas(Trigger.LOG, test.gas)
//       expect(logPrice).to.equal(
//         linkForGas(
//           BigNumber.from(test.gas),
//           totalLogOverhead,
//           BigNumber.from(test.multiplier),
//           BigNumber.from(test.premium),
//           BigNumber.from(test.flatFee),
//           l1CostWei,
//         ).total,
//       )
//     }
//   }
//
//   const verifyConsistentAccounting = async (
//     maxAllowedSpareChange: BigNumber,
//   ) => {
//     const expectedLinkBalance = (await registry.getState()).state
//       .expectedLinkBalance
//     const linkTokenBalance = await linkToken.balanceOf(registry.address)
//     const upkeepIdBalance = (await registry.getUpkeep(upkeepId)).balance
//     let totalKeeperBalance = BigNumber.from(0)
//     for (let i = 0; i < keeperAddresses.length; i++) {
//       totalKeeperBalance = totalKeeperBalance.add(
//         (await registry.getTransmitterInfo(keeperAddresses[i])).balance,
//       )
//     }
//     const ownerBalance = (await registry.getState()).state.ownerLinkBalance
//     assert.isTrue(expectedLinkBalance.eq(linkTokenBalance))
//     assert.isTrue(
//       upkeepIdBalance
//         .add(totalKeeperBalance)
//         .add(ownerBalance)
//         .lte(expectedLinkBalance),
//     )
//     assert.isTrue(
//       expectedLinkBalance
//         .sub(upkeepIdBalance)
//         .sub(totalKeeperBalance)
//         .sub(ownerBalance)
//         .lte(maxAllowedSpareChange),
//     )
//   }
//
//   interface GetTransmitTXOptions {
//     numSigners?: number
//     startingSignerIndex?: number
//     gasLimit?: BigNumberish
//     gasPrice?: BigNumberish
//     performGas?: BigNumberish
//     performData?: string
//     checkBlockNum?: number
//     checkBlockHash?: string
//     logBlockHash?: BytesLike
//     txHash?: BytesLike
//     logIndex?: number
//     timestamp?: number
//   }
//
//   const getTransmitTx = async (
//     registry: IKeeperRegistry,
//     transmitter: Signer,
//     upkeepIds: BigNumber[],
//     overrides: GetTransmitTXOptions = {},
//   ) => {
//     const latestBlock = await ethers.provider.getBlock('latest')
//     const configDigest = (await registry.getState()).state.latestConfigDigest
//     const config = {
//       numSigners: f + 1,
//       startingSignerIndex: 0,
//       performData: '0x',
//       performGas,
//       checkBlockNum: latestBlock.number,
//       checkBlockHash: latestBlock.hash,
//       logIndex: 0,
//       txHash: undefined, // assigned uniquely below
//       logBlockHash: undefined, // assigned uniquely below
//       timestamp: now(),
//       gasLimit: undefined,
//       gasPrice: undefined,
//     }
//     Object.assign(config, overrides)
//     const upkeeps: UpkeepData[] = []
//     for (let i = 0; i < upkeepIds.length; i++) {
//       let trigger: string
//       switch (getTriggerType(upkeepIds[i])) {
//         case Trigger.CONDITION:
//           trigger = encodeBlockTrigger({
//             blockNum: config.checkBlockNum,
//             blockHash: config.checkBlockHash,
//           })
//           break
//         case Trigger.LOG:
//           trigger = encodeLogTrigger({
//             logBlockHash: config.logBlockHash || ethers.utils.randomBytes(32),
//             txHash: config.txHash || ethers.utils.randomBytes(32),
//             logIndex: config.logIndex,
//             blockNum: config.checkBlockNum,
//             blockHash: config.checkBlockHash,
//           })
//           break
//       }
//       upkeeps.push({
//         Id: upkeepIds[i],
//         performGas: config.performGas,
//         trigger,
//         performData: config.performData,
//       })
//     }
//
//     const report = makeReport(upkeeps)
//     const reportContext = [configDigest, epochAndRound5_1, emptyBytes32]
//     const sigs = signReport(
//       reportContext,
//       report,
//       signers.slice(
//         config.startingSignerIndex,
//         config.startingSignerIndex + config.numSigners,
//       ),
//     )
//
//     type txOverride = {
//       gasLimit?: BigNumberish | Promise<BigNumberish>
//       gasPrice?: BigNumberish | Promise<BigNumberish>
//     }
//     const txOverrides: txOverride = {}
//     if (config.gasLimit) {
//       txOverrides.gasLimit = config.gasLimit
//     }
//     if (config.gasPrice) {
//       txOverrides.gasPrice = config.gasPrice
//     }
//
//     return registry
//       .connect(transmitter)
//       .transmit(
//         [configDigest, epochAndRound5_1, emptyBytes32],
//         report,
//         sigs.rs,
//         sigs.ss,
//         sigs.vs,
//         txOverrides,
//       )
//   }
//
//   const getTransmitTxWithReport = async (
//     registry: IKeeperRegistry,
//     transmitter: Signer,
//     report: BytesLike,
//   ) => {
//     const configDigest = (await registry.getState()).state.latestConfigDigest
//     const reportContext = [configDigest, epochAndRound5_1, emptyBytes32]
//     const sigs = signReport(reportContext, report, signers.slice(0, f + 1))
//
//     return registry
//       .connect(transmitter)
//       .transmit(
//         [configDigest, epochAndRound5_1, emptyBytes32],
//         report,
//         sigs.rs,
//         sigs.ss,
//         sigs.vs,
//       )
//   }
//
//   const setup = async () => {
//     linkToken = await linkTokenFactory.connect(owner).deploy()
//     gasPriceFeed = await mockV3AggregatorFactory
//       .connect(owner)
//       .deploy(0, gasWei)
//     linkEthFeed = await mockV3AggregatorFactory
//       .connect(owner)
//       .deploy(9, linkEth)
//     const upkeepTranscoderFactory = await ethers.getContractFactory(
//       'UpkeepTranscoder4_0',
//     )
//     transcoder = await upkeepTranscoderFactory.connect(owner).deploy()
//     mockArbGasInfo = await mockArbGasInfoFactory.connect(owner).deploy()
//     mockOVMGasPriceOracle = await mockOVMGasPriceOracleFactory
//       .connect(owner)
//       .deploy()
//     streamsLookupUpkeep = await streamsLookupUpkeepFactory
//       .connect(owner)
//       .deploy(
//         BigNumber.from('10000'),
//         BigNumber.from('100'),
//         false /* useArbBlock */,
//         true /* staging */,
//         false /* verify mercury response */,
//       )
//
//     const arbOracleCode = await ethers.provider.send('eth_getCode', [
//       mockArbGasInfo.address,
//     ])
//     await ethers.provider.send('hardhat_setCode', [
//       '0x000000000000000000000000000000000000006C',
//       arbOracleCode,
//     ])
//
//     const optOracleCode = await ethers.provider.send('eth_getCode', [
//       mockOVMGasPriceOracle.address,
//     ])
//     await ethers.provider.send('hardhat_setCode', [
//       '0x420000000000000000000000000000000000000F',
//       optOracleCode,
//     ])
//
//     const mockArbSys = await new MockArbSysFactory(owner).deploy()
//     const arbSysCode = await ethers.provider.send('eth_getCode', [
//       mockArbSys.address,
//     ])
//     await ethers.provider.send('hardhat_setCode', [
//       '0x0000000000000000000000000000000000000064',
//       arbSysCode,
//     ])
//
//     config = {
//       paymentPremiumPPB,
//       flatFeeMicroLink,
//       checkGasLimit,
//       stalenessSeconds,
//       gasCeilingMultiplier,
//       minUpkeepSpend,
//       maxCheckDataSize,
//       maxPerformDataSize,
//       maxRevertDataSize,
//       maxPerformGas,
//       fallbackGasPrice,
//       fallbackLinkPrice,
//       transcoder: transcoder.address,
//       registrars: [],
//       upkeepPrivilegeManager: upkeepManager,
//     }
//
//     baseConfig = [
//       signerAddresses,
//       keeperAddresses,
//       f,
//       encodeConfig(config),
//       offchainVersion,
//       offchainBytes,
//     ]
//
//     registry = await deployRegistry21(
//       owner,
//       Mode.DEFAULT,
//       linkToken.address,
//       linkEthFeed.address,
//       gasPriceFeed.address,
//     )
//
//     arbRegistry = await deployRegistry21(
//       owner,
//       Mode.ARBITRUM,
//       linkToken.address,
//       linkEthFeed.address,
//       gasPriceFeed.address,
//     )
//
//     opRegistry = await deployRegistry21(
//       owner,
//       Mode.OPTIMISM,
//       linkToken.address,
//       linkEthFeed.address,
//       gasPriceFeed.address,
//     )
//
//     mgRegistry = await deployRegistry21(
//       owner,
//       Mode.DEFAULT,
//       linkToken.address,
//       linkEthFeed.address,
//       gasPriceFeed.address,
//     )
//
//     blankRegistry = await deployRegistry21(
//       owner,
//       Mode.DEFAULT,
//       linkToken.address,
//       linkEthFeed.address,
//       gasPriceFeed.address,
//     )
//
//     registryConditionalOverhead = await registry.getConditionalGasOverhead()
//     registryLogOverhead = await registry.getLogGasOverhead()
//     registryPerSignerGasOverhead = await registry.getPerSignerGasOverhead()
//     registryPerPerformByteGasOverhead =
//       await registry.getPerPerformByteGasOverhead()
//     cancellationDelay = (await registry.getCancellationDelay()).toNumber()
//
//     for (const reg of [registry, arbRegistry, opRegistry, mgRegistry]) {
//       await reg.connect(owner).setConfig(...baseConfig)
//       await reg.connect(owner).setPayees(payees)
//       await linkToken.connect(admin).approve(reg.address, toWei('1000'))
//       await linkToken.connect(owner).approve(reg.address, toWei('1000'))
//     }
//
//     mock = await upkeepMockFactory.deploy()
//     await linkToken
//       .connect(owner)
//       .transfer(await admin.getAddress(), toWei('1000'))
//     let tx = await registry
//       .connect(owner)
//       [
//         'registerUpkeep(address,uint32,address,bytes,bytes)'
//       ](mock.address, performGas, await admin.getAddress(), randomBytes, '0x')
//     upkeepId = await getUpkeepID(tx)
//
//     autoFunderUpkeep = await upkeepAutoFunderFactory
//       .connect(owner)
//       .deploy(linkToken.address, registry.address)
//     tx = await registry
//       .connect(owner)
//       [
//         'registerUpkeep(address,uint32,address,bytes,bytes)'
//       ](autoFunderUpkeep.address, performGas, autoFunderUpkeep.address, randomBytes, '0x')
//     afUpkeepId = await getUpkeepID(tx)
//
//     ltUpkeep = await deployMockContract(owner, ILogAutomationactory.abi)
//     tx = await registry
//       .connect(owner)
//       [
//         'registerUpkeep(address,uint32,address,uint8,bytes,bytes,bytes)'
//       ](ltUpkeep.address, performGas, await admin.getAddress(), Trigger.LOG, '0x', logTriggerConfig, emptyBytes)
//     logUpkeepId = await getUpkeepID(tx)
//
//     await autoFunderUpkeep.setUpkeepId(afUpkeepId)
//     // Give enough funds for upkeep as well as to the upkeep contract
//     await linkToken
//       .connect(owner)
//       .transfer(autoFunderUpkeep.address, toWei('1000'))
//
//     tx = await registry
//       .connect(owner)
//       [
//         'registerUpkeep(address,uint32,address,bytes,bytes)'
//       ](streamsLookupUpkeep.address, performGas, await admin.getAddress(), randomBytes, '0x')
//     streamsLookupUpkeepId = await getUpkeepID(tx)
//   }
//
//   const getMultipleUpkeepsDeployedAndFunded = async (
//     numPassingConditionalUpkeeps: number,
//     numPassingLogUpkeeps: number,
//     numFailingUpkeeps: number,
//   ) => {
//     const passingConditionalUpkeepIds = []
//     const passingLogUpkeepIds = []
//     const failingUpkeepIds = []
//     for (let i = 0; i < numPassingConditionalUpkeeps; i++) {
//       const mock = await upkeepMockFactory.deploy()
//       await mock.setCanPerform(true)
//       await mock.setPerformGasToBurn(BigNumber.from('0'))
//       const tx = await registry
//         .connect(owner)
//         [
//           'registerUpkeep(address,uint32,address,bytes,bytes)'
//         ](mock.address, performGas, await admin.getAddress(), randomBytes, '0x')
//       const condUpkeepId = await getUpkeepID(tx)
//       passingConditionalUpkeepIds.push(condUpkeepId)
//
//       // Add funds to passing upkeeps
//       await registry.connect(admin).addFunds(condUpkeepId, toWei('100'))
//     }
//     for (let i = 0; i < numPassingLogUpkeeps; i++) {
//       const mock = await upkeepMockFactory.deploy()
//       await mock.setCanPerform(true)
//       await mock.setPerformGasToBurn(BigNumber.from('0'))
//       const tx = await registry
//         .connect(owner)
//         [
//           'registerUpkeep(address,uint32,address,uint8,bytes,bytes,bytes)'
//         ](mock.address, performGas, await admin.getAddress(), Trigger.LOG, '0x', logTriggerConfig, emptyBytes)
//       const logUpkeepId = await getUpkeepID(tx)
//       passingLogUpkeepIds.push(logUpkeepId)
//
//       // Add funds to passing upkeeps
//       await registry.connect(admin).addFunds(logUpkeepId, toWei('100'))
//     }
//     for (let i = 0; i < numFailingUpkeeps; i++) {
//       const mock = await upkeepMockFactory.deploy()
//       await mock.setCanPerform(true)
//       await mock.setPerformGasToBurn(BigNumber.from('0'))
//       const tx = await registry
//         .connect(owner)
//         [
//           'registerUpkeep(address,uint32,address,bytes,bytes)'
//         ](mock.address, performGas, await admin.getAddress(), randomBytes, '0x')
//       const failingUpkeepId = await getUpkeepID(tx)
//       failingUpkeepIds.push(failingUpkeepId)
//     }
//     return {
//       passingConditionalUpkeepIds,
//       passingLogUpkeepIds,
//       failingUpkeepIds,
//     }
//   }
//
//   beforeEach(async () => {
//     await loadFixture(setup)
//   })
//
//   describe('#transmit', () => {
//     const fArray = [1, 5, 10]
//
//     it('reverts when registry is paused', async () => {
//       await registry.connect(owner).pause()
//       await evmRevert(
//         getTransmitTx(registry, keeper1, [upkeepId]),
//         'RegistryPaused()',
//       )
//     })
//
//     it('reverts when called by non active transmitter', async () => {
//       await evmRevert(
//         getTransmitTx(registry, payee1, [upkeepId]),
//         'OnlyActiveTransmitters()',
//       )
//     })
//
//     it('reverts when report data lengths mismatches', async () => {
//       const upkeepIds = []
//       const gasLimits: BigNumber[] = []
//       const triggers: string[] = []
//       const performDatas = []
//
//       upkeepIds.push(upkeepId)
//       gasLimits.push(performGas)
//       triggers.push('0x')
//       performDatas.push('0x')
//       // Push an extra perform data
//       performDatas.push('0x')
//
//       const report = encodeReport({
//         fastGasWei: 0,
//         linkNative: 0,
//         upkeepIds,
//         gasLimits,
//         triggers,
//         performDatas,
//       })
//
//       await evmRevert(
//         getTransmitTxWithReport(registry, keeper1, report),
//         'InvalidReport()',
//       )
//     })
//
//     it('returns early when invalid upkeepIds are included in report', async () => {
//       const tx = await getTransmitTx(registry, keeper1, [
//         upkeepId.add(BigNumber.from('1')),
//       ])
//
//       const receipt = await tx.wait()
//       const cancelledUpkeepReportLogs = parseCancelledUpkeepReportLogs(receipt)
//       // exactly 1 CancelledUpkeepReport log should be emitted
//       assert.equal(cancelledUpkeepReportLogs.length, 1)
//     })
//
//     it('returns early when upkeep has insufficient funds', async () => {
//       const tx = await getTransmitTx(registry, keeper1, [upkeepId])
//       const receipt = await tx.wait()
//       const insufficientFundsUpkeepReportLogs =
//         parseInsufficientFundsUpkeepReportLogs(receipt)
//       // exactly 1 InsufficientFundsUpkeepReportLogs log should be emitted
//       assert.equal(insufficientFundsUpkeepReportLogs.length, 1)
//     })
//
//     it('permits retrying log triggers after funds are added', async () => {
//       const txHash = ethers.utils.randomBytes(32)
//       let tx = await getTransmitTx(registry, keeper1, [logUpkeepId], {
//         txHash,
//         logIndex: 0,
//       })
//       let receipt = await tx.wait()
//       const insufficientFundsLogs =
//         parseInsufficientFundsUpkeepReportLogs(receipt)
//       assert.equal(insufficientFundsLogs.length, 1)
//       registry.connect(admin).addFunds(logUpkeepId, toWei('100'))
//       tx = await getTransmitTx(registry, keeper1, [logUpkeepId], {
//         txHash,
//         logIndex: 0,
//       })
//       receipt = await tx.wait()
//       const performedLogs = parseUpkeepPerformedLogs(receipt)
//       assert.equal(performedLogs.length, 1)
//     })
//
//     context('When the upkeep is funded', async () => {
//       beforeEach(async () => {
//         // Fund the upkeep
//         await Promise.all([
//           registry.connect(admin).addFunds(upkeepId, toWei('100')),
//           registry.connect(admin).addFunds(logUpkeepId, toWei('100')),
//         ])
//       })
//
//       it('handles duplicate upkeepIDs', async () => {
//         const tests: [string, BigNumber, number, number][] = [
//           // [name, upkeep, num stale, num performed]
//           ['conditional', upkeepId, 1, 1], // checkBlocks must be sequential
//           ['log-trigger', logUpkeepId, 0, 2], // logs are deduped based on the "trigger ID"
//         ]
//         for (const [type, id, nStale, nPerformed] of tests) {
//           const tx = await getTransmitTx(registry, keeper1, [id, id])
//           const receipt = await tx.wait()
//           const staleUpkeepReport = parseStaleUpkeepReportLogs(receipt)
//           const upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
//           assert.equal(
//             staleUpkeepReport.length,
//             nStale,
//             `wrong log count for ${type} upkeep`,
//           )
//           assert.equal(
//             upkeepPerformedLogs.length,
//             nPerformed,
//             `wrong log count for ${type} upkeep`,
//           )
//         }
//       })
//
//       it('handles duplicate log triggers', async () => {
//         const logBlockHash = ethers.utils.randomBytes(32)
//         const txHash = ethers.utils.randomBytes(32)
//         const logIndex = 0
//         const expectedDedupKey = ethers.utils.solidityKeccak256(
//           ['uint256', 'bytes32', 'bytes32', 'uint32'],
//           [logUpkeepId, logBlockHash, txHash, logIndex],
//         )
//         assert.isFalse(await registry.hasDedupKey(expectedDedupKey))
//         const tx = await getTransmitTx(
//           registry,
//           keeper1,
//           [logUpkeepId, logUpkeepId],
//           { logBlockHash, txHash, logIndex }, // will result in the same dedup key
//         )
//         const receipt = await tx.wait()
//         const staleUpkeepReport = parseStaleUpkeepReportLogs(receipt)
//         const upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
//         assert.equal(staleUpkeepReport.length, 1)
//         assert.equal(upkeepPerformedLogs.length, 1)
//         assert.isTrue(await registry.hasDedupKey(expectedDedupKey))
//         await expect(tx)
//           .to.emit(registry, 'DedupKeyAdded')
//           .withArgs(expectedDedupKey)
//       })
//
//       it('returns early when check block number is less than last perform (block)', async () => {
//         // First perform an upkeep to put last perform block number on upkeep state
//         const tx = await getTransmitTx(registry, keeper1, [upkeepId])
//         await tx.wait()
//         const lastPerformed = (await registry.getUpkeep(upkeepId))
//           .lastPerformedBlockNumber
//         const lastPerformBlock = await ethers.provider.getBlock(lastPerformed)
//         assert.equal(lastPerformed.toString(), tx.blockNumber?.toString())
//         // Try to transmit a report which has checkBlockNumber = lastPerformed-1, should result in stale report
//         const transmitTx = await getTransmitTx(registry, keeper1, [upkeepId], {
//           checkBlockNum: lastPerformBlock.number - 1,
//           checkBlockHash: lastPerformBlock.parentHash,
//         })
//         const receipt = await transmitTx.wait()
//         const staleUpkeepReportLogs = parseStaleUpkeepReportLogs(receipt)
//         // exactly 1 StaleUpkeepReportLogs log should be emitted
//         assert.equal(staleUpkeepReportLogs.length, 1)
//       })
//
//       it('handles case when check block hash does not match', async () => {
//         const tests: [string, BigNumber][] = [
//           ['conditional', upkeepId],
//           ['log-trigger', logUpkeepId],
//         ]
//         for (const [type, id] of tests) {
//           const latestBlock = await ethers.provider.getBlock('latest')
//           // Try to transmit a report which has incorrect checkBlockHash
//           const tx = await getTransmitTx(registry, keeper1, [id], {
//             checkBlockNum: latestBlock.number - 1,
//             checkBlockHash: latestBlock.hash, // should be latestBlock.parentHash
//           })
//
//           const receipt = await tx.wait()
//           const reorgedUpkeepReportLogs = parseReorgedUpkeepReportLogs(receipt)
//           // exactly 1 ReorgedUpkeepReportLogs log should be emitted
//           assert.equal(
//             reorgedUpkeepReportLogs.length,
//             1,
//             `wrong log count for ${type} upkeep`,
//           )
//         }
//       })
//
//       it('handles case when check block number is older than 256 blocks', async () => {
//         for (let i = 0; i < 256; i++) {
//           await ethers.provider.send('evm_mine', [])
//         }
//         const tests: [string, BigNumber][] = [
//           ['conditional', upkeepId],
//           ['log-trigger', logUpkeepId],
//         ]
//         for (const [type, id] of tests) {
//           const latestBlock = await ethers.provider.getBlock('latest')
//           const old = await ethers.provider.getBlock(latestBlock.number - 256)
//           // Try to transmit a report which has incorrect checkBlockHash
//           const tx = await getTransmitTx(registry, keeper1, [id], {
//             checkBlockNum: old.number,
//             checkBlockHash: old.hash,
//           })
//
//           const receipt = await tx.wait()
//           const reorgedUpkeepReportLogs = parseReorgedUpkeepReportLogs(receipt)
//           // exactly 1 ReorgedUpkeepReportLogs log should be emitted
//           assert.equal(
//             reorgedUpkeepReportLogs.length,
//             1,
//             `wrong log count for ${type} upkeep`,
//           )
//         }
//       })
//
//       it('allows bypassing reorg protection with empty blockhash', async () => {
//         const tests: [string, BigNumber][] = [
//           ['conditional', upkeepId],
//           ['log-trigger', logUpkeepId],
//         ]
//         for (const [type, id] of tests) {
//           const latestBlock = await ethers.provider.getBlock('latest')
//           const tx = await getTransmitTx(registry, keeper1, [id], {
//             checkBlockNum: latestBlock.number,
//             checkBlockHash: emptyBytes32,
//           })
//           const receipt = await tx.wait()
//           const upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
//           assert.equal(
//             upkeepPerformedLogs.length,
//             1,
//             `wrong log count for ${type} upkeep`,
//           )
//         }
//       })
//
//       it('allows very old trigger block numbers when bypassing reorg protection with empty blockhash', async () => {
//         // mine enough blocks so that blockhash(1) is unavailable
//         for (let i = 0; i <= 256; i++) {
//           await ethers.provider.send('evm_mine', [])
//         }
//         const tests: [string, BigNumber][] = [
//           ['conditional', upkeepId],
//           ['log-trigger', logUpkeepId],
//         ]
//         for (const [type, id] of tests) {
//           const tx = await getTransmitTx(registry, keeper1, [id], {
//             checkBlockNum: 1,
//             checkBlockHash: emptyBytes32,
//           })
//           const receipt = await tx.wait()
//           const upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
//           assert.equal(
//             upkeepPerformedLogs.length,
//             1,
//             `wrong log count for ${type} upkeep`,
//           )
//         }
//       })
//
//       it('returns early when future block number is provided as trigger, irrespective of blockhash being present', async () => {
//         const tests: [string, BigNumber][] = [
//           ['conditional', upkeepId],
//           ['log-trigger', logUpkeepId],
//         ]
//         for (const [type, id] of tests) {
//           const latestBlock = await ethers.provider.getBlock('latest')
//
//           // Should fail when blockhash is empty
//           let tx = await getTransmitTx(registry, keeper1, [id], {
//             checkBlockNum: latestBlock.number + 100,
//             checkBlockHash: emptyBytes32,
//           })
//           let receipt = await tx.wait()
//           let reorgedUpkeepReportLogs = parseReorgedUpkeepReportLogs(receipt)
//           // exactly 1 ReorgedUpkeepReportLogs log should be emitted
//           assert.equal(
//             reorgedUpkeepReportLogs.length,
//             1,
//             `wrong log count for ${type} upkeep`,
//           )
//
//           // Should also fail when blockhash is not empty
//           tx = await getTransmitTx(registry, keeper1, [id], {
//             checkBlockNum: latestBlock.number + 100,
//             checkBlockHash: latestBlock.hash,
//           })
//           receipt = await tx.wait()
//           reorgedUpkeepReportLogs = parseReorgedUpkeepReportLogs(receipt)
//           // exactly 1 ReorgedUpkeepReportLogs log should be emitted
//           assert.equal(
//             reorgedUpkeepReportLogs.length,
//             1,
//             `wrong log count for ${type} upkeep`,
//           )
//         }
//       })
//
//       it('returns early when upkeep is cancelled and cancellation delay has gone', async () => {
//         const latestBlockReport = await makeLatestBlockReport([upkeepId])
//         await registry.connect(admin).cancelUpkeep(upkeepId)
//
//         for (let i = 0; i < cancellationDelay; i++) {
//           await ethers.provider.send('evm_mine', [])
//         }
//
//         const tx = await getTransmitTxWithReport(
//           registry,
//           keeper1,
//           latestBlockReport,
//         )
//
//         const receipt = await tx.wait()
//         const cancelledUpkeepReportLogs =
//           parseCancelledUpkeepReportLogs(receipt)
//         // exactly 1 CancelledUpkeepReport log should be emitted
//         assert.equal(cancelledUpkeepReportLogs.length, 1)
//       })
//
//       it('does not revert if the target cannot execute', async () => {
//         await mock.setCanPerform(false)
//         const tx = await getTransmitTx(registry, keeper1, [upkeepId])
//
//         const receipt = await tx.wait()
//         const upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
//         // exactly 1 Upkeep Performed should be emitted
//         assert.equal(upkeepPerformedLogs.length, 1)
//         const upkeepPerformedLog = upkeepPerformedLogs[0]
//
//         const success = upkeepPerformedLog.args.success
//         assert.equal(success, false)
//       })
//
//       it('does not revert if the target runs out of gas', async () => {
//         await mock.setCanPerform(false)
//
//         const tx = await getTransmitTx(registry, keeper1, [upkeepId], {
//           performGas: 10, // too little gas
//         })
//
//         const receipt = await tx.wait()
//         const upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
//         // exactly 1 Upkeep Performed should be emitted
//         assert.equal(upkeepPerformedLogs.length, 1)
//         const upkeepPerformedLog = upkeepPerformedLogs[0]
//
//         const success = upkeepPerformedLog.args.success
//         assert.equal(success, false)
//       })
//
//       it('reverts if not enough gas supplied', async () => {
//         await evmRevert(
//           getTransmitTx(registry, keeper1, [upkeepId], {
//             gasLimit: performGas,
//           }),
//         )
//       })
//
//       it('executes the data passed to the registry', async () => {
//         await mock.setCanPerform(true)
//
//         const tx = await getTransmitTx(registry, keeper1, [upkeepId], {
//           performData: randomBytes,
//         })
//         const receipt = await tx.wait()
//
//         const upkeepPerformedWithABI = [
//           'event UpkeepPerformedWith(bytes upkeepData)',
//         ]
//         const iface = new ethers.utils.Interface(upkeepPerformedWithABI)
//         const parsedLogs = []
//         for (let i = 0; i < receipt.logs.length; i++) {
//           const log = receipt.logs[i]
//           try {
//             parsedLogs.push(iface.parseLog(log))
//           } catch (e) {
//             // ignore log
//           }
//         }
//         assert.equal(parsedLogs.length, 1)
//         assert.equal(parsedLogs[0].args.upkeepData, randomBytes)
//       })
//
//       it('uses actual execution price for payment and premium calculation', async () => {
//         // Actual multiplier is 2, but we set gasPrice to be 1x gasWei
//         const gasPrice = gasWei.mul(BigNumber.from('1'))
//         await mock.setCanPerform(true)
//         const registryPremiumBefore = (await registry.getState()).state
//           .totalPremium
//         const tx = await getTransmitTx(registry, keeper1, [upkeepId], {
//           gasPrice,
//         })
//         const receipt = await tx.wait()
//         const registryPremiumAfter = (await registry.getState()).state
//           .totalPremium
//         const premium = registryPremiumAfter.sub(registryPremiumBefore)
//
//         const upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
//         // exactly 1 Upkeep Performed should be emitted
//         assert.equal(upkeepPerformedLogs.length, 1)
//         const upkeepPerformedLog = upkeepPerformedLogs[0]
//
//         const gasUsed = upkeepPerformedLog.args.gasUsed
//         const gasOverhead = upkeepPerformedLog.args.gasOverhead
//         const totalPayment = upkeepPerformedLog.args.totalPayment
//
//         assert.equal(
//           linkForGas(
//             gasUsed,
//             gasOverhead,
//             BigNumber.from('1'), // Not the config multiplier, but the actual gas used
//             paymentPremiumPPB,
//             flatFeeMicroLink,
//           ).total.toString(),
//           totalPayment.toString(),
//         )
//
//         assert.equal(
//           linkForGas(
//             gasUsed,
//             gasOverhead,
//             BigNumber.from('1'), // Not the config multiplier, but the actual gas used
//             paymentPremiumPPB,
//             flatFeeMicroLink,
//           ).premium.toString(),
//           premium.toString(),
//         )
//       })
//
//       it('only pays at a rate up to the gas ceiling [ @skip-coverage ]', async () => {
//         // Actual multiplier is 2, but we set gasPrice to be 10x
//         const gasPrice = gasWei.mul(BigNumber.from('10'))
//         await mock.setCanPerform(true)
//
//         const tx = await getTransmitTx(registry, keeper1, [upkeepId], {
//           gasPrice,
//         })
//         const receipt = await tx.wait()
//         const upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
//         // exactly 1 Upkeep Performed should be emitted
//         assert.equal(upkeepPerformedLogs.length, 1)
//         const upkeepPerformedLog = upkeepPerformedLogs[0]
//
//         const gasUsed = upkeepPerformedLog.args.gasUsed
//         const gasOverhead = upkeepPerformedLog.args.gasOverhead
//         const totalPayment = upkeepPerformedLog.args.totalPayment
//
//         assert.equal(
//           linkForGas(
//             gasUsed,
//             gasOverhead,
//             gasCeilingMultiplier, // Should be same with exisitng multiplier
//             paymentPremiumPPB,
//             flatFeeMicroLink,
//           ).total.toString(),
//           totalPayment.toString(),
//         )
//       })
//
//       it('correctly accounts for l payment', async () => {
//         await mock.setCanPerform(true)
//         // Same as MockArbGasInfo.sol
//         const l1CostWeiArb = BigNumber.from(1000000)
//
//         let tx = await arbRegistry
//           .connect(owner)
//           [
//             'registerUpkeep(address,uint32,address,bytes,bytes)'
//           ](mock.address, performGas, await admin.getAddress(), randomBytes, '0x')
//         const testUpkeepId = await getUpkeepID(tx)
//         await arbRegistry.connect(owner).addFunds(testUpkeepId, toWei('100'))
//
//         // Do the thing
//         tx = await getTransmitTx(
//           arbRegistry,
//           keeper1,
//           [testUpkeepId],
//
//           { gasPrice: gasWei.mul('5') }, // High gas price so that it gets capped
//         )
//         const receipt = await tx.wait()
//         const upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
//         // exactly 1 Upkeep Performed should be emitted
//         assert.equal(upkeepPerformedLogs.length, 1)
//         const upkeepPerformedLog = upkeepPerformedLogs[0]
//
//         const gasUsed = upkeepPerformedLog.args.gasUsed
//         const gasOverhead = upkeepPerformedLog.args.gasOverhead
//         const totalPayment = upkeepPerformedLog.args.totalPayment
//
//         assert.equal(
//           linkForGas(
//             gasUsed,
//             gasOverhead,
//             gasCeilingMultiplier,
//             paymentPremiumPPB,
//             flatFeeMicroLink,
//             l1CostWeiArb.div(gasCeilingMultiplier), // Dividing by gasCeilingMultiplier as it gets multiplied later
//           ).total.toString(),
//           totalPayment.toString(),
//         )
//       })
//
//       itMaybe('can self fund', async () => {
//         const maxPayment = await registry.getMaxPaymentForGas(
//           Trigger.CONDITION,
//           performGas,
//         )
//
//         // First set auto funding amount to 0 and verify that balance is deducted upon performUpkeep
//         let initialBalance = toWei('100')
//         await registry.connect(owner).addFunds(afUpkeepId, initialBalance)
//         await autoFunderUpkeep.setAutoFundLink(0)
//         await autoFunderUpkeep.setIsEligible(true)
//         await getTransmitTx(registry, keeper1, [afUpkeepId])
//
//         let postUpkeepBalance = (await registry.getUpkeep(afUpkeepId)).balance
//         assert.isTrue(postUpkeepBalance.lt(initialBalance)) // Balance should be deducted
//         assert.isTrue(postUpkeepBalance.gte(initialBalance.sub(maxPayment))) // Balance should not be deducted more than maxPayment
//
//         // Now set auto funding amount to 100 wei and verify that the balance increases
//         initialBalance = postUpkeepBalance
//         const autoTopupAmount = toWei('100')
//         await autoFunderUpkeep.setAutoFundLink(autoTopupAmount)
//         await autoFunderUpkeep.setIsEligible(true)
//         await getTransmitTx(registry, keeper1, [afUpkeepId])
//
//         postUpkeepBalance = (await registry.getUpkeep(afUpkeepId)).balance
//         // Balance should increase by autoTopupAmount and decrease by max maxPayment
//         assert.isTrue(
//           postUpkeepBalance.gte(
//             initialBalance.add(autoTopupAmount).sub(maxPayment),
//           ),
//         )
//       })
//
//       it('can self cancel', async () => {
//         await registry.connect(owner).addFunds(afUpkeepId, toWei('100'))
//
//         await autoFunderUpkeep.setIsEligible(true)
//         await autoFunderUpkeep.setShouldCancel(true)
//
//         let registration = await registry.getUpkeep(afUpkeepId)
//         const oldExpiration = registration.maxValidBlocknumber
//
//         // Do the thing
//         await getTransmitTx(registry, keeper1, [afUpkeepId])
//
//         // Verify upkeep gets cancelled
//         registration = await registry.getUpkeep(afUpkeepId)
//         const newExpiration = registration.maxValidBlocknumber
//         assert.isTrue(newExpiration.lt(oldExpiration))
//       })
//
//       it('reverts when configDigest mismatches', async () => {
//         const report = await makeLatestBlockReport([upkeepId])
//         const reportContext = [emptyBytes32, epochAndRound5_1, emptyBytes32] // wrong config digest
//         const sigs = signReport(reportContext, report, signers.slice(0, f + 1))
//         await evmRevert(
//           registry
//             .connect(keeper1)
//             .transmit(
//               [reportContext[0], reportContext[1], reportContext[2]],
//               report,
//               sigs.rs,
//               sigs.ss,
//               sigs.vs,
//             ),
//           'ConfigDigestMismatch()',
//         )
//       })
//
//       it('reverts with incorrect number of signatures', async () => {
//         const configDigest = (await registry.getState()).state
//           .latestConfigDigest
//         const report = await makeLatestBlockReport([upkeepId])
//         const reportContext = [configDigest, epochAndRound5_1, emptyBytes32] // wrong config digest
//         const sigs = signReport(reportContext, report, signers.slice(0, f + 2))
//         await evmRevert(
//           registry
//             .connect(keeper1)
//             .transmit(
//               [reportContext[0], reportContext[1], reportContext[2]],
//               report,
//               sigs.rs,
//               sigs.ss,
//               sigs.vs,
//             ),
//           'IncorrectNumberOfSignatures()',
//         )
//       })
//
//       it('reverts with invalid signature for inactive signers', async () => {
//         const configDigest = (await registry.getState()).state
//           .latestConfigDigest
//         const report = await makeLatestBlockReport([upkeepId])
//         const reportContext = [configDigest, epochAndRound5_1, emptyBytes32] // wrong config digest
//         const sigs = signReport(reportContext, report, [
//           new ethers.Wallet(ethers.Wallet.createRandom()),
//           new ethers.Wallet(ethers.Wallet.createRandom()),
//         ])
//         await evmRevert(
//           registry
//             .connect(keeper1)
//             .transmit(
//               [reportContext[0], reportContext[1], reportContext[2]],
//               report,
//               sigs.rs,
//               sigs.ss,
//               sigs.vs,
//             ),
//           'OnlyActiveSigners()',
//         )
//       })
//
//       it('reverts with invalid signature for duplicated signers', async () => {
//         const configDigest = (await registry.getState()).state
//           .latestConfigDigest
//         const report = await makeLatestBlockReport([upkeepId])
//         const reportContext = [configDigest, epochAndRound5_1, emptyBytes32] // wrong config digest
//         const sigs = signReport(reportContext, report, [signer1, signer1])
//         await evmRevert(
//           registry
//             .connect(keeper1)
//             .transmit(
//               [reportContext[0], reportContext[1], reportContext[2]],
//               report,
//               sigs.rs,
//               sigs.ss,
//               sigs.vs,
//             ),
//           'DuplicateSigners()',
//         )
//       })
//
//       itMaybe(
//         'has a large enough gas overhead to cover upkeep that use all its gas [ @skip-coverage ]',
//         async () => {
//           await registry.connect(owner).setConfigTypeSafe(
//             signerAddresses,
//             keeperAddresses,
//             10, // maximise f to maximise overhead
//             config,
//             offchainVersion,
//             offchainBytes,
//           )
//           const tx = await registry
//             .connect(owner)
//             ['registerUpkeep(address,uint32,address,bytes,bytes)'](
//               mock.address,
//               maxPerformGas, // max allowed gas
//               await admin.getAddress(),
//               randomBytes,
//               '0x',
//             )
//           const testUpkeepId = await getUpkeepID(tx)
//           await registry.connect(admin).addFunds(testUpkeepId, toWei('100'))
//
//           let performData = '0x'
//           for (let i = 0; i < maxPerformDataSize.toNumber(); i++) {
//             performData += '11'
//           } // max allowed performData
//
//           await mock.setCanPerform(true)
//           await mock.setPerformGasToBurn(maxPerformGas)
//
//           await getTransmitTx(registry, keeper1, [testUpkeepId], {
//             gasLimit: maxPerformGas.add(transmitGasOverhead),
//             numSigners: 11,
//             performData,
//           }) // Should not revert
//         },
//       )
//
//       itMaybe(
//         'performs upkeep, deducts payment, updates lastPerformed and emits events',
//         async () => {
//           await mock.setCanPerform(true)
//
//           for (const i in fArray) {
//             const newF = fArray[i]
//             await registry
//               .connect(owner)
//               .setConfigTypeSafe(
//                 signerAddresses,
//                 keeperAddresses,
//                 newF,
//                 config,
//                 offchainVersion,
//                 offchainBytes,
//               )
//             const checkBlock = await ethers.provider.getBlock('latest')
//
//             const keeperBefore = await registry.getTransmitterInfo(
//               await keeper1.getAddress(),
//             )
//             const registrationBefore = await registry.getUpkeep(upkeepId)
//             const registryPremiumBefore = (await registry.getState()).state
//               .totalPremium
//             const keeperLinkBefore = await linkToken.balanceOf(
//               await keeper1.getAddress(),
//             )
//             const registryLinkBefore = await linkToken.balanceOf(
//               registry.address,
//             )
//
//             // Do the thing
//             const tx = await getTransmitTx(registry, keeper1, [upkeepId], {
//               checkBlockNum: checkBlock.number,
//               checkBlockHash: checkBlock.hash,
//               numSigners: newF + 1,
//             })
//
//             const receipt = await tx.wait()
//
//             const upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
//             // exactly 1 Upkeep Performed should be emitted
//             assert.equal(upkeepPerformedLogs.length, 1)
//             const upkeepPerformedLog = upkeepPerformedLogs[0]
//
//             const id = upkeepPerformedLog.args.id
//             const success = upkeepPerformedLog.args.success
//             const trigger = upkeepPerformedLog.args.trigger
//             const gasUsed = upkeepPerformedLog.args.gasUsed
//             const gasOverhead = upkeepPerformedLog.args.gasOverhead
//             const totalPayment = upkeepPerformedLog.args.totalPayment
//             assert.equal(id.toString(), upkeepId.toString())
//             assert.equal(success, true)
//             assert.equal(
//               trigger,
//               encodeBlockTrigger({
//                 blockNum: checkBlock.number,
//                 blockHash: checkBlock.hash,
//               }),
//             )
//             assert.isTrue(gasUsed.gt(BigNumber.from('0')))
//             assert.isTrue(gasOverhead.gt(BigNumber.from('0')))
//             assert.isTrue(totalPayment.gt(BigNumber.from('0')))
//
//             const keeperAfter = await registry.getTransmitterInfo(
//               await keeper1.getAddress(),
//             )
//             const registrationAfter = await registry.getUpkeep(upkeepId)
//             const keeperLinkAfter = await linkToken.balanceOf(
//               await keeper1.getAddress(),
//             )
//             const registryLinkAfter = await linkToken.balanceOf(
//               registry.address,
//             )
//             const registryPremiumAfter = (await registry.getState()).state
//               .totalPremium
//             const premium = registryPremiumAfter.sub(registryPremiumBefore)
//             // Keeper payment is gasPayment + premium / num keepers
//             const keeperPayment = totalPayment
//               .sub(premium)
//               .add(premium.div(BigNumber.from(keeperAddresses.length)))
//
//             assert.equal(
//               keeperAfter.balance.sub(keeperPayment).toString(),
//               keeperBefore.balance.toString(),
//             )
//             assert.equal(
//               registrationBefore.balance.sub(totalPayment).toString(),
//               registrationAfter.balance.toString(),
//             )
//             assert.isTrue(keeperLinkAfter.eq(keeperLinkBefore))
//             assert.isTrue(registryLinkBefore.eq(registryLinkAfter))
//
//             // Amount spent should be updated correctly
//             assert.equal(
//               registrationAfter.amountSpent.sub(totalPayment).toString(),
//               registrationBefore.amountSpent.toString(),
//             )
//             assert.isTrue(
//               registrationAfter.amountSpent
//                 .sub(registrationBefore.amountSpent)
//                 .eq(registrationBefore.balance.sub(registrationAfter.balance)),
//             )
//             // Last perform block number should be updated
//             assert.equal(
//               registrationAfter.lastPerformedBlockNumber.toString(),
//               tx.blockNumber?.toString(),
//             )
//
//             // Latest epoch should be 5
//             assert.equal((await registry.getState()).state.latestEpoch, 5)
//           }
//         },
//       )
//
//       describeMaybe(
//         'Gas benchmarking conditional upkeeps [ @skip-coverage ]',
//         function () {
//           const fs = [1, 10]
//           fs.forEach(function (newF) {
//             it(
//               'When f=' +
//                 newF +
//                 ' calculates gas overhead appropriately within a margin for different scenarios',
//               async () => {
//                 // Perform the upkeep once to remove non-zero storage slots and have predictable gas measurement
//                 let tx = await getTransmitTx(registry, keeper1, [upkeepId])
//                 await tx.wait()
//
//                 // Different test scenarios
//                 let longBytes = '0x'
//                 for (let i = 0; i < maxPerformDataSize.toNumber(); i++) {
//                   longBytes += '11'
//                 }
//                 const upkeepSuccessArray = [true, false]
//                 const performGasArray = [5000, performGas]
//                 const performDataArray = ['0x', longBytes]
//
//                 for (const i in upkeepSuccessArray) {
//                   for (const j in performGasArray) {
//                     for (const k in performDataArray) {
//                       const upkeepSuccess = upkeepSuccessArray[i]
//                       const performGas = performGasArray[j]
//                       const performData = performDataArray[k]
//
//                       await mock.setCanPerform(upkeepSuccess)
//                       await mock.setPerformGasToBurn(performGas)
//                       await registry
//                         .connect(owner)
//                         .setConfigTypeSafe(
//                           signerAddresses,
//                           keeperAddresses,
//                           newF,
//                           config,
//                           offchainVersion,
//                           offchainBytes,
//                         )
//                       tx = await getTransmitTx(registry, keeper1, [upkeepId], {
//                         numSigners: newF + 1,
//                         performData,
//                       })
//                       const receipt = await tx.wait()
//                       const upkeepPerformedLogs =
//                         parseUpkeepPerformedLogs(receipt)
//                       // exactly 1 Upkeep Performed should be emitted
//                       assert.equal(upkeepPerformedLogs.length, 1)
//                       const upkeepPerformedLog = upkeepPerformedLogs[0]
//
//                       const upkeepGasUsed = upkeepPerformedLog.args.gasUsed
//                       const chargedGasOverhead =
//                         upkeepPerformedLog.args.gasOverhead
//                       const actualGasOverhead =
//                         receipt.gasUsed.sub(upkeepGasUsed)
//
//                       assert.isTrue(upkeepGasUsed.gt(BigNumber.from('0')))
//                       assert.isTrue(chargedGasOverhead.gt(BigNumber.from('0')))
//
//                       console.log(
//                         'Gas Benchmarking conditional upkeeps:',
//                         'upkeepSuccess=',
//                         upkeepSuccess,
//                         'performGas=',
//                         performGas.toString(),
//                         'performData length=',
//                         performData.length / 2 - 1,
//                         'sig verification ( f =',
//                         newF,
//                         '): calculated overhead: ',
//                         chargedGasOverhead.toString(),
//                         ' actual overhead: ',
//                         actualGasOverhead.toString(),
//                         ' margin over gasUsed: ',
//                         chargedGasOverhead.sub(actualGasOverhead).toString(),
//                       )
//
//                       // Overhead should not get capped
//                       const gasOverheadCap = registryConditionalOverhead
//                         .add(
//                           registryPerSignerGasOverhead.mul(
//                             BigNumber.from(newF + 1),
//                           ),
//                         )
//                         .add(
//                           BigNumber.from(
//                             registryPerPerformByteGasOverhead.toNumber() *
//                               performData.length,
//                           ),
//                         )
//                       const gasCapMinusOverhead =
//                         gasOverheadCap.sub(chargedGasOverhead)
//                       assert.isTrue(
//                         gasCapMinusOverhead.gt(BigNumber.from(0)),
//                         'Gas overhead got capped. Verify gas overhead variables in test match those in the registry. To not have the overheads capped increase REGISTRY_GAS_OVERHEAD by atleast ' +
//                           gasCapMinusOverhead.toString(),
//                       )
//                       // total gas charged should be greater than tx gas but within gasCalculationMargin
//                       assert.isTrue(
//                         chargedGasOverhead.gt(actualGasOverhead),
//                         'Gas overhead calculated is too low, increase account gas variables (ACCOUNTING_FIXED_GAS_OVERHEAD/ACCOUNTING_PER_SIGNER_GAS_OVERHEAD) by atleast ' +
//                           actualGasOverhead.sub(chargedGasOverhead).toString(),
//                       )
//
//                       assert.isTrue(
//                         chargedGasOverhead
//                           .sub(actualGasOverhead)
//                           .lt(gasCalculationMargin),
//                       ),
//                         'Gas overhead calculated is too high, decrease account gas variables (ACCOUNTING_FIXED_GAS_OVERHEAD/ACCOUNTING_PER_SIGNER_GAS_OVERHEAD)  by atleast ' +
//                           chargedGasOverhead
//                             .sub(chargedGasOverhead)
//                             .sub(gasCalculationMargin)
//                             .toString()
//                     }
//                   }
//                 }
//               },
//             )
//           })
//         },
//       )
//
//       describeMaybe(
//         'Gas benchmarking log upkeeps [ @skip-coverage ]',
//         function () {
//           const fs = [1, 10]
//           fs.forEach(function (newF) {
//             it(
//               'When f=' +
//                 newF +
//                 ' calculates gas overhead appropriately within a margin',
//               async () => {
//                 // Perform the upkeep once to remove non-zero storage slots and have predictable gas measurement
//                 let tx = await getTransmitTx(registry, keeper1, [logUpkeepId])
//                 await tx.wait()
//                 const performData = '0x'
//                 await mock.setCanPerform(true)
//                 await mock.setPerformGasToBurn(performGas)
//                 await registry.setConfigTypeSafe(
//                   signerAddresses,
//                   keeperAddresses,
//                   newF,
//                   config,
//                   offchainVersion,
//                   offchainBytes,
//                 )
//                 tx = await getTransmitTx(registry, keeper1, [logUpkeepId], {
//                   numSigners: newF + 1,
//                   performData,
//                 })
//                 const receipt = await tx.wait()
//                 const upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
//                 // exactly 1 Upkeep Performed should be emitted
//                 assert.equal(upkeepPerformedLogs.length, 1)
//                 const upkeepPerformedLog = upkeepPerformedLogs[0]
//
//                 const upkeepGasUsed = upkeepPerformedLog.args.gasUsed
//                 const chargedGasOverhead = upkeepPerformedLog.args.gasOverhead
//                 const actualGasOverhead = receipt.gasUsed.sub(upkeepGasUsed)
//
//                 assert.isTrue(upkeepGasUsed.gt(BigNumber.from('0')))
//                 assert.isTrue(chargedGasOverhead.gt(BigNumber.from('0')))
//
//                 console.log(
//                   'Gas Benchmarking log upkeeps:',
//                   'upkeepSuccess=',
//                   true,
//                   'performGas=',
//                   performGas.toString(),
//                   'performData length=',
//                   performData.length / 2 - 1,
//                   'sig verification ( f =',
//                   newF,
//                   '): calculated overhead: ',
//                   chargedGasOverhead.toString(),
//                   ' actual overhead: ',
//                   actualGasOverhead.toString(),
//                   ' margin over gasUsed: ',
//                   chargedGasOverhead.sub(actualGasOverhead).toString(),
//                 )
//
//                 // Overhead should not get capped
//                 const gasOverheadCap = registryLogOverhead
//                   .add(
//                     registryPerSignerGasOverhead.mul(BigNumber.from(newF + 1)),
//                   )
//                   .add(
//                     BigNumber.from(
//                       registryPerPerformByteGasOverhead.toNumber() *
//                         performData.length,
//                     ),
//                   )
//                 const gasCapMinusOverhead =
//                   gasOverheadCap.sub(chargedGasOverhead)
//                 assert.isTrue(
//                   gasCapMinusOverhead.gt(BigNumber.from(0)),
//                   'Gas overhead got capped. Verify gas overhead variables in test match those in the registry. To not have the overheads capped increase REGISTRY_GAS_OVERHEAD by atleast ' +
//                     gasCapMinusOverhead.toString(),
//                 )
//                 // total gas charged should be greater than tx gas but within gasCalculationMargin
//                 assert.isTrue(
//                   chargedGasOverhead.gt(actualGasOverhead),
//                   'Gas overhead calculated is too low, increase account gas variables (ACCOUNTING_FIXED_GAS_OVERHEAD/ACCOUNTING_PER_SIGNER_GAS_OVERHEAD) by atleast ' +
//                     actualGasOverhead.sub(chargedGasOverhead).toString(),
//                 )
//
//                 assert.isTrue(
//                   chargedGasOverhead
//                     .sub(actualGasOverhead)
//                     .lt(gasCalculationMargin),
//                 ),
//                   'Gas overhead calculated is too high, decrease account gas variables (ACCOUNTING_FIXED_GAS_OVERHEAD/ACCOUNTING_PER_SIGNER_GAS_OVERHEAD)  by atleast ' +
//                     chargedGasOverhead
//                       .sub(chargedGasOverhead)
//                       .sub(gasCalculationMargin)
//                       .toString()
//               },
//             )
//           })
//         },
//       )
//     })
//   })
//
//   describeMaybe(
//     '#transmit with upkeep batches [ @skip-coverage ]',
//     function () {
//       const numPassingConditionalUpkeepsArray = [0, 1, 5]
//       const numPassingLogUpkeepsArray = [0, 1, 5]
//       const numFailingUpkeepsArray = [0, 3]
//
//       for (let idx = 0; idx < numPassingConditionalUpkeepsArray.length; idx++) {
//         for (let jdx = 0; jdx < numPassingLogUpkeepsArray.length; jdx++) {
//           for (let kdx = 0; kdx < numFailingUpkeepsArray.length; kdx++) {
//             const numPassingConditionalUpkeeps =
//               numPassingConditionalUpkeepsArray[idx]
//             const numPassingLogUpkeeps = numPassingLogUpkeepsArray[jdx]
//             const numFailingUpkeeps = numFailingUpkeepsArray[kdx]
//             if (
//               numPassingConditionalUpkeeps == 0 &&
//               numPassingLogUpkeeps == 0
//             ) {
//               continue
//             }
//             it(
//               '[Conditional:' +
//                 numPassingConditionalUpkeeps +
//                 ',Log:' +
//                 numPassingLogUpkeeps +
//                 ',Failures:' +
//                 numFailingUpkeeps +
//                 '] performs successful upkeeps and does not charge failing upkeeps',
//               async () => {
//                 const allUpkeeps = await getMultipleUpkeepsDeployedAndFunded(
//                   numPassingConditionalUpkeeps,
//                   numPassingLogUpkeeps,
//                   numFailingUpkeeps,
//                 )
//                 const passingConditionalUpkeepIds =
//                   allUpkeeps.passingConditionalUpkeepIds
//                 const passingLogUpkeepIds = allUpkeeps.passingLogUpkeepIds
//                 const failingUpkeepIds = allUpkeeps.failingUpkeepIds
//
//                 const keeperBefore = await registry.getTransmitterInfo(
//                   await keeper1.getAddress(),
//                 )
//                 const keeperLinkBefore = await linkToken.balanceOf(
//                   await keeper1.getAddress(),
//                 )
//                 const registryLinkBefore = await linkToken.balanceOf(
//                   registry.address,
//                 )
//                 const registryPremiumBefore = (await registry.getState()).state
//                   .totalPremium
//                 const registrationConditionalPassingBefore = await Promise.all(
//                   passingConditionalUpkeepIds.map(async (id) => {
//                     const reg = await registry.getUpkeep(BigNumber.from(id))
//                     assert.equal(reg.lastPerformedBlockNumber.toString(), '0')
//                     return reg
//                   }),
//                 )
//                 const registrationLogPassingBefore = await Promise.all(
//                   passingLogUpkeepIds.map(async (id) => {
//                     const reg = await registry.getUpkeep(BigNumber.from(id))
//                     assert.equal(reg.lastPerformedBlockNumber.toString(), '0')
//                     return reg
//                   }),
//                 )
//                 const registrationFailingBefore = await Promise.all(
//                   failingUpkeepIds.map(async (id) => {
//                     const reg = await registry.getUpkeep(BigNumber.from(id))
//                     assert.equal(reg.lastPerformedBlockNumber.toString(), '0')
//                     return reg
//                   }),
//                 )
//
//                 const tx = await getTransmitTx(
//                   registry,
//                   keeper1,
//                   passingConditionalUpkeepIds.concat(
//                     passingLogUpkeepIds.concat(failingUpkeepIds),
//                   ),
//                 )
//
//                 const receipt = await tx.wait()
//                 const upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
//                 // exactly numPassingUpkeeps Upkeep Performed should be emitted
//                 assert.equal(
//                   upkeepPerformedLogs.length,
//                   numPassingConditionalUpkeeps + numPassingLogUpkeeps,
//                 )
//                 const insufficientFundsLogs =
//                   parseInsufficientFundsUpkeepReportLogs(receipt)
//                 // exactly numFailingUpkeeps Upkeep Performed should be emitted
//                 assert.equal(insufficientFundsLogs.length, numFailingUpkeeps)
//
//                 const keeperAfter = await registry.getTransmitterInfo(
//                   await keeper1.getAddress(),
//                 )
//                 const keeperLinkAfter = await linkToken.balanceOf(
//                   await keeper1.getAddress(),
//                 )
//                 const registryLinkAfter = await linkToken.balanceOf(
//                   registry.address,
//                 )
//                 const registrationConditionalPassingAfter = await Promise.all(
//                   passingConditionalUpkeepIds.map(async (id) => {
//                     return await registry.getUpkeep(BigNumber.from(id))
//                   }),
//                 )
//                 const registrationLogPassingAfter = await Promise.all(
//                   passingLogUpkeepIds.map(async (id) => {
//                     return await registry.getUpkeep(BigNumber.from(id))
//                   }),
//                 )
//                 const registrationFailingAfter = await Promise.all(
//                   failingUpkeepIds.map(async (id) => {
//                     return await registry.getUpkeep(BigNumber.from(id))
//                   }),
//                 )
//                 const registryPremiumAfter = (await registry.getState()).state
//                   .totalPremium
//                 const premium = registryPremiumAfter.sub(registryPremiumBefore)
//
//                 let netPayment = BigNumber.from('0')
//                 for (let i = 0; i < numPassingConditionalUpkeeps; i++) {
//                   const id = upkeepPerformedLogs[i].args.id
//                   const gasUsed = upkeepPerformedLogs[i].args.gasUsed
//                   const gasOverhead = upkeepPerformedLogs[i].args.gasOverhead
//                   const totalPayment = upkeepPerformedLogs[i].args.totalPayment
//
//                   expect(id).to.equal(passingConditionalUpkeepIds[i])
//                   assert.isTrue(gasUsed.gt(BigNumber.from('0')))
//                   assert.isTrue(gasOverhead.gt(BigNumber.from('0')))
//                   assert.isTrue(totalPayment.gt(BigNumber.from('0')))
//
//                   // Balance should be deducted
//                   assert.equal(
//                     registrationConditionalPassingBefore[i].balance
//                       .sub(totalPayment)
//                       .toString(),
//                     registrationConditionalPassingAfter[i].balance.toString(),
//                   )
//
//                   // Amount spent should be updated correctly
//                   assert.equal(
//                     registrationConditionalPassingAfter[i].amountSpent
//                       .sub(totalPayment)
//                       .toString(),
//                     registrationConditionalPassingBefore[
//                       i
//                     ].amountSpent.toString(),
//                   )
//
//                   // Last perform block number should be updated
//                   assert.equal(
//                     registrationConditionalPassingAfter[
//                       i
//                     ].lastPerformedBlockNumber.toString(),
//                     tx.blockNumber?.toString(),
//                   )
//
//                   netPayment = netPayment.add(totalPayment)
//                 }
//
//                 for (let i = 0; i < numPassingLogUpkeeps; i++) {
//                   const id =
//                     upkeepPerformedLogs[numPassingConditionalUpkeeps + i].args
//                       .id
//                   const gasUsed =
//                     upkeepPerformedLogs[numPassingConditionalUpkeeps + i].args
//                       .gasUsed
//                   const gasOverhead =
//                     upkeepPerformedLogs[numPassingConditionalUpkeeps + i].args
//                       .gasOverhead
//                   const totalPayment =
//                     upkeepPerformedLogs[numPassingConditionalUpkeeps + i].args
//                       .totalPayment
//
//                   expect(id).to.equal(passingLogUpkeepIds[i])
//                   assert.isTrue(gasUsed.gt(BigNumber.from('0')))
//                   assert.isTrue(gasOverhead.gt(BigNumber.from('0')))
//                   assert.isTrue(totalPayment.gt(BigNumber.from('0')))
//
//                   // Balance should be deducted
//                   assert.equal(
//                     registrationLogPassingBefore[i].balance
//                       .sub(totalPayment)
//                       .toString(),
//                     registrationLogPassingAfter[i].balance.toString(),
//                   )
//
//                   // Amount spent should be updated correctly
//                   assert.equal(
//                     registrationLogPassingAfter[i].amountSpent
//                       .sub(totalPayment)
//                       .toString(),
//                     registrationLogPassingBefore[i].amountSpent.toString(),
//                   )
//
//                   // Last perform block number should not be updated for log triggers
//                   assert.equal(
//                     registrationLogPassingAfter[
//                       i
//                     ].lastPerformedBlockNumber.toString(),
//                     '0',
//                   )
//
//                   netPayment = netPayment.add(totalPayment)
//                 }
//
//                 for (let i = 0; i < numFailingUpkeeps; i++) {
//                   // InsufficientFunds log should be emitted
//                   const id = insufficientFundsLogs[i].args.id
//                   expect(id).to.equal(failingUpkeepIds[i])
//
//                   // Balance and amount spent should be same
//                   assert.equal(
//                     registrationFailingBefore[i].balance.toString(),
//                     registrationFailingAfter[i].balance.toString(),
//                   )
//                   assert.equal(
//                     registrationFailingBefore[i].amountSpent.toString(),
//                     registrationFailingAfter[i].amountSpent.toString(),
//                   )
//
//                   // Last perform block number should not be updated
//                   assert.equal(
//                     registrationFailingAfter[
//                       i
//                     ].lastPerformedBlockNumber.toString(),
//                     '0',
//                   )
//                 }
//
//                 // Keeper payment is gasPayment + premium / num keepers
//                 const keeperPayment = netPayment
//                   .sub(premium)
//                   .add(premium.div(BigNumber.from(keeperAddresses.length)))
//
//                 // Keeper should be paid net payment for all passed upkeeps
//                 assert.equal(
//                   keeperAfter.balance.sub(keeperPayment).toString(),
//                   keeperBefore.balance.toString(),
//                 )
//
//                 assert.isTrue(keeperLinkAfter.eq(keeperLinkBefore))
//                 assert.isTrue(registryLinkBefore.eq(registryLinkAfter))
//               },
//             )
//
//             it(
//               '[Conditional:' +
//                 numPassingConditionalUpkeeps +
//                 ',Log' +
//                 numPassingLogUpkeeps +
//                 ',Failures:' +
//                 numFailingUpkeeps +
//                 '] splits gas overhead appropriately among performed upkeeps [ @skip-coverage ]',
//               async () => {
//                 const allUpkeeps = await getMultipleUpkeepsDeployedAndFunded(
//                   numPassingConditionalUpkeeps,
//                   numPassingLogUpkeeps,
//                   numFailingUpkeeps,
//                 )
//                 const passingConditionalUpkeepIds =
//                   allUpkeeps.passingConditionalUpkeepIds
//                 const passingLogUpkeepIds = allUpkeeps.passingLogUpkeepIds
//                 const failingUpkeepIds = allUpkeeps.failingUpkeepIds
//
//                 // Perform the upkeeps once to remove non-zero storage slots and have predictable gas measurement
//                 let tx = await getTransmitTx(
//                   registry,
//                   keeper1,
//                   passingConditionalUpkeepIds.concat(
//                     passingLogUpkeepIds.concat(failingUpkeepIds),
//                   ),
//                 )
//
//                 await tx.wait()
//
//                 // Do the actual thing
//
//                 tx = await getTransmitTx(
//                   registry,
//                   keeper1,
//                   passingConditionalUpkeepIds.concat(
//                     passingLogUpkeepIds.concat(failingUpkeepIds),
//                   ),
//                 )
//
//                 const receipt = await tx.wait()
//                 const upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
//                 // exactly numPassingUpkeeps Upkeep Performed should be emitted
//                 assert.equal(
//                   upkeepPerformedLogs.length,
//                   numPassingConditionalUpkeeps + numPassingLogUpkeeps,
//                 )
//
//                 const gasConditionalOverheadCap =
//                   registryConditionalOverhead.add(
//                     registryPerSignerGasOverhead.mul(BigNumber.from(f + 1)),
//                   )
//                 const gasLogOverheadCap = registryLogOverhead.add(
//                   registryPerSignerGasOverhead.mul(BigNumber.from(f + 1)),
//                 )
//
//                 const overheadCanGetCapped =
//                   numFailingUpkeeps > 0 &&
//                   numPassingConditionalUpkeeps <= 1 &&
//                   numPassingLogUpkeeps <= 1
//                 // Can happen if there are failing upkeeps and only 1 successful upkeep of each type
//                 let netGasUsedPlusOverhead = BigNumber.from('0')
//
//                 for (let i = 0; i < numPassingConditionalUpkeeps; i++) {
//                   const gasUsed = upkeepPerformedLogs[i].args.gasUsed
//                   const gasOverhead = upkeepPerformedLogs[i].args.gasOverhead
//
//                   assert.isTrue(gasUsed.gt(BigNumber.from('0')))
//                   assert.isTrue(gasOverhead.gt(BigNumber.from('0')))
//
//                   // Overhead should not exceed capped
//                   assert.isTrue(gasOverhead.lte(gasConditionalOverheadCap))
//
//                   // Overhead should be same for every upkeep since they have equal performData, hence same caps
//                   assert.isTrue(
//                     gasOverhead.eq(upkeepPerformedLogs[0].args.gasOverhead),
//                   )
//
//                   netGasUsedPlusOverhead = netGasUsedPlusOverhead
//                     .add(gasUsed)
//                     .add(gasOverhead)
//                 }
//                 for (let i = 0; i < numPassingLogUpkeeps; i++) {
//                   const gasUsed =
//                     upkeepPerformedLogs[numPassingConditionalUpkeeps + i].args
//                       .gasUsed
//                   const gasOverhead =
//                     upkeepPerformedLogs[numPassingConditionalUpkeeps + i].args
//                       .gasOverhead
//
//                   assert.isTrue(gasUsed.gt(BigNumber.from('0')))
//                   assert.isTrue(gasOverhead.gt(BigNumber.from('0')))
//
//                   // Overhead should not exceed capped
//                   assert.isTrue(gasOverhead.lte(gasLogOverheadCap))
//
//                   // Overhead should be same for every upkeep since they have equal performData, hence same caps
//                   assert.isTrue(
//                     gasOverhead.eq(
//                       upkeepPerformedLogs[numPassingConditionalUpkeeps].args
//                         .gasOverhead,
//                     ),
//                   )
//
//                   netGasUsedPlusOverhead = netGasUsedPlusOverhead
//                     .add(gasUsed)
//                     .add(gasOverhead)
//                 }
//
//                 const overheadsGotCapped =
//                   (numPassingConditionalUpkeeps > 0 &&
//                     upkeepPerformedLogs[0].args.gasOverhead.eq(
//                       gasConditionalOverheadCap,
//                     )) ||
//                   (numPassingLogUpkeeps > 0 &&
//                     upkeepPerformedLogs[
//                       numPassingConditionalUpkeeps
//                     ].args.gasOverhead.eq(gasLogOverheadCap))
//                 // Should only get capped in certain scenarios
//                 if (overheadsGotCapped) {
//                   assert.isTrue(
//                     overheadCanGetCapped,
//                     'Gas overhead got capped. Verify gas overhead variables in test match those in the registry. To not have the overheads capped increase REGISTRY_GAS_OVERHEAD',
//                   )
//                 }
//
//                 console.log(
//                   'Gas Benchmarking - batching (passedConditionalUpkeeps: ',
//                   numPassingConditionalUpkeeps,
//                   'passedLogUpkeeps:',
//                   numPassingLogUpkeeps,
//                   'failedUpkeeps:',
//                   numFailingUpkeeps,
//                   '): ',
//                   'overheadsGotCapped',
//                   overheadsGotCapped,
//                   numPassingConditionalUpkeeps > 0
//                     ? 'calculated conditional overhead'
//                     : '',
//                   numPassingConditionalUpkeeps > 0
//                     ? upkeepPerformedLogs[0].args.gasOverhead.toString()
//                     : '',
//                   numPassingLogUpkeeps > 0 ? 'calculated log overhead' : '',
//                   numPassingLogUpkeeps > 0
//                     ? upkeepPerformedLogs[
//                         numPassingConditionalUpkeeps
//                       ].args.gasOverhead.toString()
//                     : '',
//                   ' margin over gasUsed',
//                   netGasUsedPlusOverhead.sub(receipt.gasUsed).toString(),
//                 )
//
//                 // If overheads dont get capped then total gas charged should be greater than tx gas
//                 // We don't check whether the net is within gasMargin as the margin changes with numFailedUpkeeps
//                 // Which is ok, as long as individual gas overhead is capped
//                 if (!overheadsGotCapped) {
//                   assert.isTrue(
//                     netGasUsedPlusOverhead.gt(receipt.gasUsed),
//                     'Gas overhead is too low, increase ACCOUNTING_PER_UPKEEP_GAS_OVERHEAD',
//                   )
//                 }
//               },
//             )
//           }
//         }
//       }
//
//       it('has enough perform gas overhead for large batches [ @skip-coverage ]', async () => {
//         const numUpkeeps = 20
//         const upkeepIds: BigNumber[] = []
//         let totalPerformGas = BigNumber.from('0')
//         for (let i = 0; i < numUpkeeps; i++) {
//           const mock = await upkeepMockFactory.deploy()
//           const tx = await registry
//             .connect(owner)
//             [
//               'registerUpkeep(address,uint32,address,bytes,bytes)'
//             ](mock.address, performGas, await admin.getAddress(), randomBytes, '0x')
//           const testUpkeepId = await getUpkeepID(tx)
//           upkeepIds.push(testUpkeepId)
//
//           // Add funds to passing upkeeps
//           await registry.connect(owner).addFunds(testUpkeepId, toWei('10'))
//
//           await mock.setCanPerform(true)
//           await mock.setPerformGasToBurn(performGas)
//
//           totalPerformGas = totalPerformGas.add(performGas)
//         }
//
//         // Should revert with no overhead added
//         await evmRevert(
//           getTransmitTx(registry, keeper1, upkeepIds, {
//             gasLimit: totalPerformGas,
//           }),
//         )
//         // Should not revert with overhead added
//         await getTransmitTx(registry, keeper1, upkeepIds, {
//           gasLimit: totalPerformGas.add(transmitGasOverhead),
//         })
//       })
//
//       it('splits l2 payment among performed upkeeps', async () => {
//         const numUpkeeps = 7
//         const upkeepIds: BigNumber[] = []
//         // Same as MockArbGasInfo.sol
//         const l1CostWeiArb = BigNumber.from(1000000)
//
//         for (let i = 0; i < numUpkeeps; i++) {
//           const mock = await upkeepMockFactory.deploy()
//           const tx = await arbRegistry
//             .connect(owner)
//             [
//               'registerUpkeep(address,uint32,address,bytes,bytes)'
//             ](mock.address, performGas, await admin.getAddress(), randomBytes, '0x')
//           const testUpkeepId = await getUpkeepID(tx)
//           upkeepIds.push(testUpkeepId)
//
//           // Add funds to passing upkeeps
//           await arbRegistry.connect(owner).addFunds(testUpkeepId, toWei('100'))
//         }
//
//         // Do the thing
//         const tx = await getTransmitTx(
//           arbRegistry,
//           keeper1,
//           upkeepIds,
//
//           { gasPrice: gasWei.mul('5') }, // High gas price so that it gets capped
//         )
//
//         const receipt = await tx.wait()
//         const upkeepPerformedLogs = parseUpkeepPerformedLogs(receipt)
//         // exactly numPassingUpkeeps Upkeep Performed should be emitted
//         assert.equal(upkeepPerformedLogs.length, numUpkeeps)
//
//         // Verify the payment calculation in upkeepPerformed[0]
//         const upkeepPerformedLog = upkeepPerformedLogs[0]
//
//         const gasUsed = upkeepPerformedLog.args.gasUsed
//         const gasOverhead = upkeepPerformedLog.args.gasOverhead
//         const totalPayment = upkeepPerformedLog.args.totalPayment
//
//         assert.equal(
//           linkForGas(
//             gasUsed,
//             gasOverhead,
//             gasCeilingMultiplier,
//             paymentPremiumPPB,
//             flatFeeMicroLink,
//             l1CostWeiArb.div(gasCeilingMultiplier), // Dividing by gasCeilingMultiplier as it gets multiplied later
//             BigNumber.from(numUpkeeps),
//           ).total.toString(),
//           totalPayment.toString(),
//         )
//       })
//     },
//   )
//
//   describe('#recoverFunds', () => {
//     const sent = toWei('7')
//
//     beforeEach(async () => {
//       await linkToken.connect(admin).approve(registry.address, toWei('100'))
//       await linkToken
//         .connect(owner)
//         .transfer(await keeper1.getAddress(), toWei('1000'))
//
//       // add funds to upkeep 1 and perform and withdraw some payment
//       const tx = await registry
//         .connect(owner)
//         [
//           'registerUpkeep(address,uint32,address,bytes,bytes)'
//         ](mock.address, performGas, await admin.getAddress(), emptyBytes, emptyBytes)
//
//       const id1 = await getUpkeepID(tx)
//       await registry.connect(admin).addFunds(id1, toWei('5'))
//
//       await getTransmitTx(registry, keeper1, [id1])
//       await getTransmitTx(registry, keeper2, [id1])
//       await getTransmitTx(registry, keeper3, [id1])
//
//       await registry
//         .connect(payee1)
//         .withdrawPayment(
//           await keeper1.getAddress(),
//           await nonkeeper.getAddress(),
//         )
//
//       // transfer funds directly to the registry
//       await linkToken.connect(keeper1).transfer(registry.address, sent)
//
//       // add funds to upkeep 2 and perform and withdraw some payment
//       const tx2 = await registry
//         .connect(owner)
//         [
//           'registerUpkeep(address,uint32,address,bytes,bytes)'
//         ](mock.address, performGas, await admin.getAddress(), emptyBytes, emptyBytes)
//       const id2 = await getUpkeepID(tx2)
//       await registry.connect(admin).addFunds(id2, toWei('5'))
//
//       await getTransmitTx(registry, keeper1, [id2])
//       await getTransmitTx(registry, keeper2, [id2])
//       await getTransmitTx(registry, keeper3, [id2])
//
//       await registry
//         .connect(payee2)
//         .withdrawPayment(
//           await keeper2.getAddress(),
//           await nonkeeper.getAddress(),
//         )
//
//       // transfer funds using onTokenTransfer
//       const data = ethers.utils.defaultAbiCoder.encode(['uint256'], [id2])
//       await linkToken
//         .connect(owner)
//         .transferAndCall(registry.address, toWei('1'), data)
//
//       // withdraw some funds
//       await registry.connect(owner).cancelUpkeep(id1)
//       await registry
//         .connect(admin)
//         .withdrawFunds(id1, await nonkeeper.getAddress())
//     })
//
//     it('reverts if not called by owner', async () => {
//       await evmRevert(
//         registry.connect(keeper1).recoverFunds(),
//         'Only callable by owner',
//       )
//     })
//
//     it('allows any funds that have been accidentally transfered to be moved', async () => {
//       const balanceBefore = await linkToken.balanceOf(registry.address)
//       const ownerBefore = await linkToken.balanceOf(await owner.getAddress())
//
//       await registry.connect(owner).recoverFunds()
//
//       const balanceAfter = await linkToken.balanceOf(registry.address)
//       const ownerAfter = await linkToken.balanceOf(await owner.getAddress())
//
//       assert.isTrue(balanceBefore.eq(balanceAfter.add(sent)))
//       assert.isTrue(ownerAfter.eq(ownerBefore.add(sent)))
//     })
//   })
//
//   describe('#getMinBalanceForUpkeep / #checkUpkeep / #transmit', () => {
//     it('calculates the minimum balance appropriately', async () => {
//       await mock.setCanCheck(true)
//
//       const oneWei = BigNumber.from(1)
//       const minBalance = await registry.getMinBalanceForUpkeep(upkeepId)
//       const tooLow = minBalance.sub(oneWei)
//
//       await registry.connect(admin).addFunds(upkeepId, tooLow)
//       let checkUpkeepResult = await registry
//         .connect(zeroAddress)
//         .callStatic['checkUpkeep(uint256)'](upkeepId)
//
//       assert.equal(checkUpkeepResult.upkeepNeeded, false)
//       assert.equal(
//         checkUpkeepResult.upkeepFailureReason,
//         UpkeepFailureReason.INSUFFICIENT_BALANCE,
//       )
//
//       await registry.connect(admin).addFunds(upkeepId, oneWei)
//       checkUpkeepResult = await registry
//         .connect(zeroAddress)
//         .callStatic['checkUpkeep(uint256)'](upkeepId)
//       assert.equal(checkUpkeepResult.upkeepNeeded, true)
//     })
//
//     it('uses maxPerformData size in checkUpkeep but actual performDataSize in transmit', async () => {
//       const tx1 = await registry
//         .connect(owner)
//         [
//           'registerUpkeep(address,uint32,address,bytes,bytes)'
//         ](mock.address, performGas, await admin.getAddress(), randomBytes, '0x')
//       const upkeepID1 = await getUpkeepID(tx1)
//       const tx2 = await registry
//         .connect(owner)
//         [
//           'registerUpkeep(address,uint32,address,bytes,bytes)'
//         ](mock.address, performGas, await admin.getAddress(), randomBytes, '0x')
//       const upkeepID2 = await getUpkeepID(tx2)
//       await mock.setCanCheck(true)
//       await mock.setCanPerform(true)
//
//       // upkeep 1 is underfunded, 2 is fully funded
//       const minBalance1 = (
//         await registry.getMinBalanceForUpkeep(upkeepID1)
//       ).sub(1)
//       const minBalance2 = await registry.getMinBalanceForUpkeep(upkeepID2)
//       await registry.connect(owner).addFunds(upkeepID1, minBalance1)
//       await registry.connect(owner).addFunds(upkeepID2, minBalance2)
//
//       // upkeep 1 check should return false, 2 should return true
//       let checkUpkeepResult = await registry
//         .connect(zeroAddress)
//         .callStatic['checkUpkeep(uint256)'](upkeepID1)
//       assert.equal(checkUpkeepResult.upkeepNeeded, false)
//       assert.equal(
//         checkUpkeepResult.upkeepFailureReason,
//         UpkeepFailureReason.INSUFFICIENT_BALANCE,
//       )
//
//       checkUpkeepResult = await registry
//         .connect(zeroAddress)
//         .callStatic['checkUpkeep(uint256)'](upkeepID2)
//       assert.equal(checkUpkeepResult.upkeepNeeded, true)
//
//       // upkeep 1 perform should return with insufficient balance using max performData size
//       let maxPerformData = '0x'
//       for (let i = 0; i < maxPerformDataSize.toNumber(); i++) {
//         maxPerformData += '11'
//       }
//
//       const tx = await getTransmitTx(registry, keeper1, [upkeepID1], {
//         gasPrice: gasWei.mul(gasCeilingMultiplier),
//         performData: maxPerformData,
//       })
//
//       const receipt = await tx.wait()
//       const insufficientFundsUpkeepReportLogs =
//         parseInsufficientFundsUpkeepReportLogs(receipt)
//       // exactly 1 InsufficientFundsUpkeepReportLogs log should be emitted
//       assert.equal(insufficientFundsUpkeepReportLogs.length, 1)
//
//       // upkeep 1 perform should succeed with empty performData
//       await getTransmitTx(registry, keeper1, [upkeepID1], {
//         gasPrice: gasWei.mul(gasCeilingMultiplier),
//       }),
//         // upkeep 2 perform should succeed with max performData size
//         await getTransmitTx(registry, keeper1, [upkeepID2], {
//           gasPrice: gasWei.mul(gasCeilingMultiplier),
//           performData: maxPerformData,
//         })
//     })
//   })
//
//   describe('#withdrawFunds', () => {
//     let upkeepId2: BigNumber
//
//     beforeEach(async () => {
//       const tx = await registry
//         .connect(owner)
//         [
//           'registerUpkeep(address,uint32,address,bytes,bytes)'
//         ](mock.address, performGas, await admin.getAddress(), randomBytes, '0x')
//       upkeepId2 = await getUpkeepID(tx)
//
//       await registry.connect(admin).addFunds(upkeepId, toWei('100'))
//       await registry.connect(admin).addFunds(upkeepId2, toWei('100'))
//
//       // Do a perform so that upkeep is charged some amount
//       await getTransmitTx(registry, keeper1, [upkeepId])
//       await getTransmitTx(registry, keeper1, [upkeepId2])
//     })
//
//     it('reverts if called on a non existing ID', async () => {
//       await evmRevert(
//         registry
//           .connect(admin)
//           .withdrawFunds(upkeepId.add(1), await payee1.getAddress()),
//         'OnlyCallableByAdmin()',
//       )
//     })
//
//     it('reverts if called by anyone but the admin', async () => {
//       await evmRevert(
//         registry
//           .connect(owner)
//           .withdrawFunds(upkeepId, await payee1.getAddress()),
//         'OnlyCallableByAdmin()',
//       )
//     })
//
//     it('reverts if called on an uncanceled upkeep', async () => {
//       await evmRevert(
//         registry
//           .connect(admin)
//           .withdrawFunds(upkeepId, await payee1.getAddress()),
//         'UpkeepNotCanceled()',
//       )
//     })
//
//     it('reverts if called with the 0 address', async () => {
//       await evmRevert(
//         registry.connect(admin).withdrawFunds(upkeepId, zeroAddress),
//         'InvalidRecipient()',
//       )
//     })
//
//     describe('after the registration is paused, then cancelled', () => {
//       it('allows the admin to withdraw', async () => {
//         const balance = await registry.getBalance(upkeepId)
//         const payee = await payee1.getAddress()
//         await registry.connect(admin).pauseUpkeep(upkeepId)
//         await registry.connect(owner).cancelUpkeep(upkeepId)
//         await expect(() =>
//           registry.connect(admin).withdrawFunds(upkeepId, payee),
//         ).to.changeTokenBalance(linkToken, payee1, balance)
//       })
//     })
//
//     describe('after the registration is cancelled', () => {
//       beforeEach(async () => {
//         await registry.connect(owner).cancelUpkeep(upkeepId)
//         await registry.connect(owner).cancelUpkeep(upkeepId2)
//       })
//
//       it('can be called successively on two upkeeps', async () => {
//         await registry
//           .connect(admin)
//           .withdrawFunds(upkeepId, await payee1.getAddress())
//         await registry
//           .connect(admin)
//           .withdrawFunds(upkeepId2, await payee1.getAddress())
//       })
//
//       it('moves the funds out and updates the balance and emits an event', async () => {
//         const payee1Before = await linkToken.balanceOf(
//           await payee1.getAddress(),
//         )
//         const registryBefore = await linkToken.balanceOf(registry.address)
//
//         let registration = await registry.getUpkeep(upkeepId)
//         const previousBalance = registration.balance
//
//         const tx = await registry
//           .connect(admin)
//           .withdrawFunds(upkeepId, await payee1.getAddress())
//         await expect(tx)
//           .to.emit(registry, 'FundsWithdrawn')
//           .withArgs(upkeepId, previousBalance, await payee1.getAddress())
//
//         const payee1After = await linkToken.balanceOf(await payee1.getAddress())
//         const registryAfter = await linkToken.balanceOf(registry.address)
//
//         assert.isTrue(payee1Before.add(previousBalance).eq(payee1After))
//         assert.isTrue(registryBefore.sub(previousBalance).eq(registryAfter))
//
//         registration = await registry.getUpkeep(upkeepId)
//         assert.equal(0, registration.balance.toNumber())
//       })
//     })
//   })
//
//   describe('#simulatePerformUpkeep', () => {
//     it('reverts if called by non zero address', async () => {
//       await evmRevert(
//         registry
//           .connect(await owner.getAddress())
//           .callStatic.simulatePerformUpkeep(upkeepId, '0x'),
//         'OnlySimulatedBackend()',
//       )
//     })
//
//     it('reverts when registry is paused', async () => {
//       await registry.connect(owner).pause()
//       await evmRevert(
//         registry
//           .connect(zeroAddress)
//           .callStatic.simulatePerformUpkeep(upkeepId, '0x'),
//         'RegistryPaused()',
//       )
//     })
//
//     it('returns false and gasUsed when perform fails', async () => {
//       await mock.setCanPerform(false)
//
//       const simulatePerformResult = await registry
//         .connect(zeroAddress)
//         .callStatic.simulatePerformUpkeep(upkeepId, '0x')
//
//       assert.equal(simulatePerformResult.success, false)
//       assert.isTrue(simulatePerformResult.gasUsed.gt(BigNumber.from('0'))) // Some gas should be used
//     })
//
//     it('returns true, gasUsed, and performGas when perform succeeds', async () => {
//       await mock.setCanPerform(true)
//
//       const simulatePerformResult = await registry
//         .connect(zeroAddress)
//         .callStatic.simulatePerformUpkeep(upkeepId, '0x')
//
//       assert.equal(simulatePerformResult.success, true)
//       assert.isTrue(simulatePerformResult.gasUsed.gt(BigNumber.from('0'))) // Some gas should be used
//     })
//
//     it('returns correct amount of gasUsed when perform succeeds', async () => {
//       await mock.setCanPerform(true)
//       await mock.setPerformGasToBurn(performGas)
//
//       const simulatePerformResult = await registry
//         .connect(zeroAddress)
//         .callStatic.simulatePerformUpkeep(upkeepId, '0x')
//
//       assert.equal(simulatePerformResult.success, true)
//       // Full execute gas should be used, with some performGasBuffer(1000)
//       assert.isTrue(
//         simulatePerformResult.gasUsed.gt(
//           performGas.sub(BigNumber.from('1000')),
//         ),
//       )
//     })
//   })
//
//   describe('#checkUpkeep', () => {
//     it('reverts if called by non zero address', async () => {
//       await evmRevert(
//         registry
//           .connect(await owner.getAddress())
//           .callStatic['checkUpkeep(uint256)'](upkeepId),
//         'OnlySimulatedBackend()',
//       )
//     })
//
//     it('returns false and error code if the upkeep is cancelled by admin', async () => {
//       await registry.connect(admin).cancelUpkeep(upkeepId)
//
//       const checkUpkeepResult = await registry
//         .connect(zeroAddress)
//         .callStatic['checkUpkeep(uint256)'](upkeepId)
//
//       assert.equal(checkUpkeepResult.upkeepNeeded, false)
//       assert.equal(checkUpkeepResult.performData, '0x')
//       assert.equal(
//         checkUpkeepResult.upkeepFailureReason,
//         UpkeepFailureReason.UPKEEP_CANCELLED,
//       )
//       expect(checkUpkeepResult.gasUsed).to.equal(0)
//       expect(checkUpkeepResult.gasLimit).to.equal(performGas)
//     })
//
//     it('returns false and error code if the upkeep is cancelled by owner', async () => {
//       await registry.connect(owner).cancelUpkeep(upkeepId)
//
//       const checkUpkeepResult = await registry
//         .connect(zeroAddress)
//         .callStatic['checkUpkeep(uint256)'](upkeepId)
//
//       assert.equal(checkUpkeepResult.upkeepNeeded, false)
//       assert.equal(checkUpkeepResult.performData, '0x')
//       assert.equal(
//         checkUpkeepResult.upkeepFailureReason,
//         UpkeepFailureReason.UPKEEP_CANCELLED,
//       )
//       expect(checkUpkeepResult.gasUsed).to.equal(0)
//       expect(checkUpkeepResult.gasLimit).to.equal(performGas)
//     })
//
//     it('returns false and error code if the registry is paused', async () => {
//       await registry.connect(owner).pause()
//
//       const checkUpkeepResult = await registry
//         .connect(zeroAddress)
//         .callStatic['checkUpkeep(uint256)'](upkeepId)
//
//       assert.equal(checkUpkeepResult.upkeepNeeded, false)
//       assert.equal(checkUpkeepResult.performData, '0x')
//       assert.equal(
//         checkUpkeepResult.upkeepFailureReason,
//         UpkeepFailureReason.REGISTRY_PAUSED,
//       )
//       expect(checkUpkeepResult.gasUsed).to.equal(0)
//       expect(checkUpkeepResult.gasLimit).to.equal(performGas)
//     })
//
//     it('returns false and error code if the upkeep is paused', async () => {
//       await registry.connect(admin).pauseUpkeep(upkeepId)
//
//       const checkUpkeepResult = await registry
//         .connect(zeroAddress)
//         .callStatic['checkUpkeep(uint256)'](upkeepId)
//
//       assert.equal(checkUpkeepResult.upkeepNeeded, false)
//       assert.equal(checkUpkeepResult.performData, '0x')
//       assert.equal(
//         checkUpkeepResult.upkeepFailureReason,
//         UpkeepFailureReason.UPKEEP_PAUSED,
//       )
//       expect(checkUpkeepResult.gasUsed).to.equal(0)
//       expect(checkUpkeepResult.gasLimit).to.equal(performGas)
//     })
//
//     it('returns false and error code if user is out of funds', async () => {
//       const checkUpkeepResult = await registry
//         .connect(zeroAddress)
//         .callStatic['checkUpkeep(uint256)'](upkeepId)
//
//       assert.equal(checkUpkeepResult.upkeepNeeded, false)
//       assert.equal(checkUpkeepResult.performData, '0x')
//       assert.equal(
//         checkUpkeepResult.upkeepFailureReason,
//         UpkeepFailureReason.INSUFFICIENT_BALANCE,
//       )
//       expect(checkUpkeepResult.gasUsed).to.equal(0)
//       expect(checkUpkeepResult.gasLimit).to.equal(performGas)
//     })
//
//     context('when the registration is funded', () => {
//       beforeEach(async () => {
//         await linkToken.connect(admin).approve(registry.address, toWei('200'))
//         await registry.connect(admin).addFunds(upkeepId, toWei('100'))
//         await registry.connect(admin).addFunds(logUpkeepId, toWei('100'))
//       })
//
//       it('returns false, error code, and revert data if the target check reverts', async () => {
//         await mock.setShouldRevertCheck(true)
//         await mock.setCheckRevertReason(
//           'custom revert error, clever way to insert offchain data',
//         )
//         const checkUpkeepResult = await registry
//           .connect(zeroAddress)
//           .callStatic['checkUpkeep(uint256)'](upkeepId)
//         assert.equal(checkUpkeepResult.upkeepNeeded, false)
//
//         const revertReasonBytes = `0x${checkUpkeepResult.performData.slice(10)}` // remove sighash
//         assert.equal(
//           ethers.utils.defaultAbiCoder.decode(['string'], revertReasonBytes)[0],
//           'custom revert error, clever way to insert offchain data',
//         )
//         assert.equal(
//           checkUpkeepResult.upkeepFailureReason,
//           UpkeepFailureReason.TARGET_CHECK_REVERTED,
//         )
//         assert.isTrue(checkUpkeepResult.gasUsed.gt(BigNumber.from('0'))) // Some gas should be used
//         expect(checkUpkeepResult.gasLimit).to.equal(performGas)
//         // Feed data should be returned here
//         assert.isTrue(checkUpkeepResult.fastGasWei.gt(BigNumber.from('0')))
//         assert.isTrue(checkUpkeepResult.linkNative.gt(BigNumber.from('0')))
//       })
//
//       it('returns false, error code, and no revert data if the target check revert data exceeds maxRevertDataSize', async () => {
//         await mock.setShouldRevertCheck(true)
//         let longRevertReason = ''
//         for (let i = 0; i <= maxRevertDataSize.toNumber(); i++) {
//           longRevertReason += 'x'
//         }
//         await mock.setCheckRevertReason(longRevertReason)
//         const checkUpkeepResult = await registry
//           .connect(zeroAddress)
//           .callStatic['checkUpkeep(uint256)'](upkeepId)
//         assert.equal(checkUpkeepResult.upkeepNeeded, false)
//
//         assert.equal(checkUpkeepResult.performData, '0x')
//         assert.equal(
//           checkUpkeepResult.upkeepFailureReason,
//           UpkeepFailureReason.REVERT_DATA_EXCEEDS_LIMIT,
//         )
//         assert.isTrue(checkUpkeepResult.gasUsed.gt(BigNumber.from('0'))) // Some gas should be used
//         expect(checkUpkeepResult.gasLimit).to.equal(performGas)
//       })
//
//       it('returns false and error code if the upkeep is not needed', async () => {
//         await mock.setCanCheck(false)
//         const checkUpkeepResult = await registry
//           .connect(zeroAddress)
//           .callStatic['checkUpkeep(uint256)'](upkeepId)
//
//         assert.equal(checkUpkeepResult.upkeepNeeded, false)
//         assert.equal(checkUpkeepResult.performData, '0x')
//         assert.equal(
//           checkUpkeepResult.upkeepFailureReason,
//           UpkeepFailureReason.UPKEEP_NOT_NEEDED,
//         )
//         assert.isTrue(checkUpkeepResult.gasUsed.gt(BigNumber.from('0'))) // Some gas should be used
//         expect(checkUpkeepResult.gasLimit).to.equal(performGas)
//       })
//
//       it('returns false and error code if the performData exceeds limit', async () => {
//         let longBytes = '0x'
//         for (let i = 0; i < 5000; i++) {
//           longBytes += '1'
//         }
//         await mock.setCanCheck(true)
//         await mock.setPerformData(longBytes)
//
//         const checkUpkeepResult = await registry
//           .connect(zeroAddress)
//           .callStatic['checkUpkeep(uint256)'](upkeepId)
//
//         assert.equal(checkUpkeepResult.upkeepNeeded, false)
//         assert.equal(checkUpkeepResult.performData, '0x')
//         assert.equal(
//           checkUpkeepResult.upkeepFailureReason,
//           UpkeepFailureReason.PERFORM_DATA_EXCEEDS_LIMIT,
//         )
//         assert.isTrue(checkUpkeepResult.gasUsed.gt(BigNumber.from('0'))) // Some gas should be used
//         expect(checkUpkeepResult.gasLimit).to.equal(performGas)
//       })
//
//       it('returns true with gas used if the target can execute', async () => {
//         await mock.setCanCheck(true)
//         await mock.setPerformData(randomBytes)
//
//         const latestBlock = await ethers.provider.getBlock('latest')
//
//         const checkUpkeepResult = await registry
//           .connect(zeroAddress)
//           .callStatic['checkUpkeep(uint256)'](upkeepId, {
//             blockTag: latestBlock.number,
//           })
//
//         assert.equal(checkUpkeepResult.upkeepNeeded, true)
//         assert.equal(checkUpkeepResult.performData, randomBytes)
//         assert.equal(
//           checkUpkeepResult.upkeepFailureReason,
//           UpkeepFailureReason.NONE,
//         )
//         assert.isTrue(checkUpkeepResult.gasUsed.gt(BigNumber.from('0'))) // Some gas should be used
//         expect(checkUpkeepResult.gasLimit).to.equal(performGas)
//         assert.isTrue(checkUpkeepResult.fastGasWei.eq(gasWei))
//         assert.isTrue(checkUpkeepResult.linkNative.eq(linkEth))
//       })
//
//       it('calls checkLog for log-trigger upkeeps', async () => {
//         const log: Log = {
//           index: 0,
//           timestamp: 0,
//           txHash: ethers.utils.randomBytes(32),
//           blockNumber: 100,
//           blockHash: ethers.utils.randomBytes(32),
//           source: randomAddress(),
//           topics: [ethers.utils.randomBytes(32), ethers.utils.randomBytes(32)],
//           data: ethers.utils.randomBytes(1000),
//         }
//
//         await ltUpkeep.mock.checkLog.withArgs(log, '0x').returns(true, '0x1234')
//
//         const checkData = encodeLog(log)
//
//         const checkUpkeepResult = await registry
//           .connect(zeroAddress)
//           .callStatic['checkUpkeep(uint256,bytes)'](logUpkeepId, checkData)
//
//         expect(checkUpkeepResult.upkeepNeeded).to.be.true
//         expect(checkUpkeepResult.performData).to.equal('0x1234')
//       })
//
//       itMaybe(
//         'has a large enough gas overhead to cover upkeeps that use all their gas [ @skip-coverage ]',
//         async () => {
//           await mock.setCanCheck(true)
//           await mock.setCheckGasToBurn(checkGasLimit)
//           const gas = checkGasLimit.add(checkGasOverhead)
//           const checkUpkeepResult = await registry
//             .connect(zeroAddress)
//             .callStatic['checkUpkeep(uint256)'](upkeepId, {
//               gasLimit: gas,
//             })
//
//           assert.equal(checkUpkeepResult.upkeepNeeded, true)
//         },
//       )
//     })
//   })
//
//   describe('#addFunds', () => {
//     const amount = toWei('1')
//
//     it('reverts if the registration does not exist', async () => {
//       await evmRevert(
//         registry.connect(keeper1).addFunds(upkeepId.add(1), amount),
//         'UpkeepCancelled()',
//       )
//     })
//
//     it('adds to the balance of the registration', async () => {
//       await registry.connect(admin).addFunds(upkeepId, amount)
//       const registration = await registry.getUpkeep(upkeepId)
//       assert.isTrue(amount.eq(registration.balance))
//     })
//
//     it('lets anyone add funds to an upkeep not just admin', async () => {
//       await linkToken.connect(owner).transfer(await payee1.getAddress(), amount)
//       await linkToken.connect(payee1).approve(registry.address, amount)
//
//       await registry.connect(payee1).addFunds(upkeepId, amount)
//       const registration = await registry.getUpkeep(upkeepId)
//       assert.isTrue(amount.eq(registration.balance))
//     })
//
//     it('emits a log', async () => {
//       const tx = await registry.connect(admin).addFunds(upkeepId, amount)
//       await expect(tx)
//         .to.emit(registry, 'FundsAdded')
//         .withArgs(upkeepId, await admin.getAddress(), amount)
//     })
//
//     it('reverts if the upkeep is canceled', async () => {
//       await registry.connect(admin).cancelUpkeep(upkeepId)
//       await evmRevert(
//         registry.connect(keeper1).addFunds(upkeepId, amount),
//         'UpkeepCancelled()',
//       )
//     })
//   })
//
//   describe('#getActiveUpkeepIDs', () => {
//     it('reverts if startIndex is out of bounds ', async () => {
//       await evmRevert(
//         registry.getActiveUpkeepIDs(numUpkeeps, 0),
//         'IndexOutOfRange()',
//       )
//       await evmRevert(
//         registry.getActiveUpkeepIDs(numUpkeeps + 1, 0),
//         'IndexOutOfRange()',
//       )
//     })
//
//     it('returns upkeep IDs bounded by maxCount', async () => {
//       let upkeepIds = await registry.getActiveUpkeepIDs(0, 1)
//       assert(upkeepIds.length == 1)
//       assert(upkeepIds[0].eq(upkeepId))
//       upkeepIds = await registry.getActiveUpkeepIDs(1, 3)
//       assert(upkeepIds.length == 3)
//       expect(upkeepIds).to.deep.equal([
//         afUpkeepId,
//         logUpkeepId,
//         streamsLookupUpkeepId,
//       ])
//     })
//
//     it('returns as many ids as possible if maxCount > num available', async () => {
//       const upkeepIds = await registry.getActiveUpkeepIDs(1, numUpkeeps + 100)
//       assert(upkeepIds.length == numUpkeeps - 1)
//     })
//
//     it('returns all upkeep IDs if maxCount is 0', async () => {
//       let upkeepIds = await registry.getActiveUpkeepIDs(0, 0)
//       assert(upkeepIds.length == numUpkeeps)
//       upkeepIds = await registry.getActiveUpkeepIDs(2, 0)
//       assert(upkeepIds.length == numUpkeeps - 2)
//     })
//   })
//
//   describe('#getMaxPaymentForGas', () => {
//     const arbL1PriceinWei = BigNumber.from(1000) // Same as MockArbGasInfo.sol
//     const l1CostWeiArb = arbL1PriceinWei.mul(16).mul(maxPerformDataSize)
//     const l1CostWeiOpt = BigNumber.from(2000000) // Same as MockOVMGasPriceOracle.sol
//     itMaybe('calculates the max fee appropriately', async () => {
//       await verifyMaxPayment(registry)
//     })
//
//     itMaybe('calculates the max fee appropriately for Arbitrum', async () => {
//       await verifyMaxPayment(arbRegistry, l1CostWeiArb)
//     })
//
//     itMaybe('calculates the max fee appropriately for Optimism', async () => {
//       await verifyMaxPayment(opRegistry, l1CostWeiOpt)
//     })
//
//     it('uses the fallback gas price if the feed has issues', async () => {
//       const expectedFallbackMaxPayment = linkForGas(
//         performGas,
//         registryConditionalOverhead
//           .add(registryPerSignerGasOverhead.mul(f + 1))
//           .add(maxPerformDataSize.mul(registryPerPerformByteGasOverhead)),
//         gasCeilingMultiplier.mul('2'), // fallbackGasPrice is 2x gas price
//         paymentPremiumPPB,
//         flatFeeMicroLink,
//       ).total
//
//       // Stale feed
//       let roundId = 99
//       const answer = 100
//       let updatedAt = 946684800 // New Years 2000 🥳
//       let startedAt = 946684799
//       await gasPriceFeed
//         .connect(owner)
//         .updateRoundData(roundId, answer, updatedAt, startedAt)
//
//       assert.equal(
//         expectedFallbackMaxPayment.toString(),
//         (
//           await registry.getMaxPaymentForGas(Trigger.CONDITION, performGas)
//         ).toString(),
//       )
//
//       // Negative feed price
//       roundId = 100
//       updatedAt = now()
//       startedAt = 946684799
//       await gasPriceFeed
//         .connect(owner)
//         .updateRoundData(roundId, -100, updatedAt, startedAt)
//
//       assert.equal(
//         expectedFallbackMaxPayment.toString(),
//         (
//           await registry.getMaxPaymentForGas(Trigger.CONDITION, performGas)
//         ).toString(),
//       )
//
//       // Zero feed price
//       roundId = 101
//       updatedAt = now()
//       startedAt = 946684799
//       await gasPriceFeed
//         .connect(owner)
//         .updateRoundData(roundId, 0, updatedAt, startedAt)
//
//       assert.equal(
//         expectedFallbackMaxPayment.toString(),
//         (
//           await registry.getMaxPaymentForGas(Trigger.CONDITION, performGas)
//         ).toString(),
//       )
//     })
//
//     it('uses the fallback link price if the feed has issues', async () => {
//       const expectedFallbackMaxPayment = linkForGas(
//         performGas,
//         registryConditionalOverhead
//           .add(registryPerSignerGasOverhead.mul(f + 1))
//           .add(maxPerformDataSize.mul(registryPerPerformByteGasOverhead)),
//         gasCeilingMultiplier.mul('2'), // fallbackLinkPrice is 1/2 link price, so multiply by 2
//         paymentPremiumPPB,
//         flatFeeMicroLink,
//       ).total
//
//       // Stale feed
//       let roundId = 99
//       const answer = 100
//       let updatedAt = 946684800 // New Years 2000 🥳
//       let startedAt = 946684799
//       await linkEthFeed
//         .connect(owner)
//         .updateRoundData(roundId, answer, updatedAt, startedAt)
//
//       assert.equal(
//         expectedFallbackMaxPayment.toString(),
//         (
//           await registry.getMaxPaymentForGas(Trigger.CONDITION, performGas)
//         ).toString(),
//       )
//
//       // Negative feed price
//       roundId = 100
//       updatedAt = now()
//       startedAt = 946684799
//       await linkEthFeed
//         .connect(owner)
//         .updateRoundData(roundId, -100, updatedAt, startedAt)
//
//       assert.equal(
//         expectedFallbackMaxPayment.toString(),
//         (
//           await registry.getMaxPaymentForGas(Trigger.CONDITION, performGas)
//         ).toString(),
//       )
//
//       // Zero feed price
//       roundId = 101
//       updatedAt = now()
//       startedAt = 946684799
//       await linkEthFeed
//         .connect(owner)
//         .updateRoundData(roundId, 0, updatedAt, startedAt)
//
//       assert.equal(
//         expectedFallbackMaxPayment.toString(),
//         (
//           await registry.getMaxPaymentForGas(Trigger.CONDITION, performGas)
//         ).toString(),
//       )
//     })
//   })
//
//   describe('#typeAndVersion', () => {
//     it('uses the correct type and version', async () => {
//       const typeAndVersion = await registry.typeAndVersion()
//       assert.equal(typeAndVersion, 'KeeperRegistry 2.1.0')
//     })
//   })
//
//   describe('#onTokenTransfer', () => {
//     const amount = toWei('1')
//
//     it('reverts if not called by the LINK token', async () => {
//       const data = ethers.utils.defaultAbiCoder.encode(['uint256'], [upkeepId])
//
//       await evmRevert(
//         registry
//           .connect(keeper1)
//           .onTokenTransfer(await keeper1.getAddress(), amount, data),
//         'OnlyCallableByLINKToken()',
//       )
//     })
//
//     it('reverts if not called with more or less than 32 bytes', async () => {
//       const longData = ethers.utils.defaultAbiCoder.encode(
//         ['uint256', 'uint256'],
//         ['33', '34'],
//       )
//       const shortData = '0x12345678'
//
//       await evmRevert(
//         linkToken
//           .connect(owner)
//           .transferAndCall(registry.address, amount, longData),
//       )
//       await evmRevert(
//         linkToken
//           .connect(owner)
//           .transferAndCall(registry.address, amount, shortData),
//       )
//     })
//
//     it('reverts if the upkeep is canceled', async () => {
//       await registry.connect(admin).cancelUpkeep(upkeepId)
//       await evmRevert(
//         registry.connect(keeper1).addFunds(upkeepId, amount),
//         'UpkeepCancelled()',
//       )
//     })
//
//     it('updates the funds of the job id passed', async () => {
//       const data = ethers.utils.defaultAbiCoder.encode(['uint256'], [upkeepId])
//
//       const before = (await registry.getUpkeep(upkeepId)).balance
//       await linkToken
//         .connect(owner)
//         .transferAndCall(registry.address, amount, data)
//       const after = (await registry.getUpkeep(upkeepId)).balance
//
//       assert.isTrue(before.add(amount).eq(after))
//     })
//   })
//
//   describeMaybe('#setConfig - onchain', () => {
//     const payment = BigNumber.from(1)
//     const flatFee = BigNumber.from(2)
//     const maxGas = BigNumber.from(6)
//     const staleness = BigNumber.from(4)
//     const ceiling = BigNumber.from(5)
//     const newMinUpkeepSpend = BigNumber.from(9)
//     const newMaxCheckDataSize = BigNumber.from(10000)
//     const newMaxPerformDataSize = BigNumber.from(10000)
//     const newMaxRevertDataSize = BigNumber.from(10000)
//     const newMaxPerformGas = BigNumber.from(10000000)
//     const fbGasEth = BigNumber.from(7)
//     const fbLinkEth = BigNumber.from(8)
//     const newTranscoder = randomAddress()
//     const newRegistrars = [randomAddress(), randomAddress()]
//     const upkeepManager = randomAddress()
//
//     const newConfig: OnChainConfig = {
//       paymentPremiumPPB: payment,
//       flatFeeMicroLink: flatFee,
//       checkGasLimit: maxGas,
//       stalenessSeconds: staleness,
//       gasCeilingMultiplier: ceiling,
//       minUpkeepSpend: newMinUpkeepSpend,
//       maxCheckDataSize: newMaxCheckDataSize,
//       maxPerformDataSize: newMaxPerformDataSize,
//       maxRevertDataSize: newMaxRevertDataSize,
//       maxPerformGas: newMaxPerformGas,
//       fallbackGasPrice: fbGasEth,
//       fallbackLinkPrice: fbLinkEth,
//       transcoder: newTranscoder,
//       registrars: newRegistrars,
//       upkeepPrivilegeManager: upkeepManager,
//     }
//
//     it('reverts when called by anyone but the proposed owner', async () => {
//       await evmRevert(
//         registry
//           .connect(payee1)
//           .setConfigTypeSafe(
//             signerAddresses,
//             keeperAddresses,
//             f,
//             newConfig,
//             offchainVersion,
//             offchainBytes,
//           ),
//         'Only callable by owner',
//       )
//     })
//
//     it('reverts if signers or transmitters are the zero address', async () => {
//       await evmRevert(
//         registry
//           .connect(owner)
//           .setConfigTypeSafe(
//             [randomAddress(), randomAddress(), randomAddress(), zeroAddress],
//             [
//               randomAddress(),
//               randomAddress(),
//               randomAddress(),
//               randomAddress(),
//             ],
//             f,
//             newConfig,
//             offchainVersion,
//             offchainBytes,
//           ),
//         'InvalidSigner()',
//       )
//
//       await evmRevert(
//         registry
//           .connect(owner)
//           .setConfigTypeSafe(
//             [
//               randomAddress(),
//               randomAddress(),
//               randomAddress(),
//               randomAddress(),
//             ],
//             [randomAddress(), randomAddress(), randomAddress(), zeroAddress],
//             f,
//             newConfig,
//             offchainVersion,
//             offchainBytes,
//           ),
//         'InvalidTransmitter()',
//       )
//     })
//
//     it('updates the onchainConfig and configDigest', async () => {
//       const old = await registry.getState()
//       const oldConfig = old.config
//       const oldState = old.state
//       assert.isTrue(paymentPremiumPPB.eq(oldConfig.paymentPremiumPPB))
//       assert.isTrue(flatFeeMicroLink.eq(oldConfig.flatFeeMicroLink))
//       assert.isTrue(stalenessSeconds.eq(oldConfig.stalenessSeconds))
//       assert.isTrue(gasCeilingMultiplier.eq(oldConfig.gasCeilingMultiplier))
//
//       await registry
//         .connect(owner)
//         .setConfigTypeSafe(
//           signerAddresses,
//           keeperAddresses,
//           f,
//           newConfig,
//           offchainVersion,
//           offchainBytes,
//         )
//
//       const updated = await registry.getState()
//       const updatedConfig = updated.config
//       const updatedState = updated.state
//       assert.equal(updatedConfig.paymentPremiumPPB, payment.toNumber())
//       assert.equal(updatedConfig.flatFeeMicroLink, flatFee.toNumber())
//       assert.equal(updatedConfig.stalenessSeconds, staleness.toNumber())
//       assert.equal(updatedConfig.gasCeilingMultiplier, ceiling.toNumber())
//       assert.equal(
//         updatedConfig.minUpkeepSpend.toString(),
//         newMinUpkeepSpend.toString(),
//       )
//       assert.equal(
//         updatedConfig.maxCheckDataSize,
//         newMaxCheckDataSize.toNumber(),
//       )
//       assert.equal(
//         updatedConfig.maxPerformDataSize,
//         newMaxPerformDataSize.toNumber(),
//       )
//       assert.equal(
//         updatedConfig.maxRevertDataSize,
//         newMaxRevertDataSize.toNumber(),
//       )
//       assert.equal(updatedConfig.maxPerformGas, newMaxPerformGas.toNumber())
//       assert.equal(updatedConfig.checkGasLimit, maxGas.toNumber())
//       assert.equal(
//         updatedConfig.fallbackGasPrice.toNumber(),
//         fbGasEth.toNumber(),
//       )
//       assert.equal(
//         updatedConfig.fallbackLinkPrice.toNumber(),
//         fbLinkEth.toNumber(),
//       )
//       assert.equal(updatedState.latestEpoch, 0)
//
//       assert(oldState.configCount + 1 == updatedState.configCount)
//       assert(
//         oldState.latestConfigBlockNumber !=
//           updatedState.latestConfigBlockNumber,
//       )
//       assert(oldState.latestConfigDigest != updatedState.latestConfigDigest)
//
//       assert.equal(updatedConfig.transcoder, newTranscoder)
//       assert.deepEqual(updatedConfig.registrars, newRegistrars)
//       assert.equal(updatedConfig.upkeepPrivilegeManager, upkeepManager)
//     })
//
//     it('maintains paused state when config is changed', async () => {
//       await registry.pause()
//       const old = await registry.getState()
//       assert.isTrue(old.state.paused)
//
//       await registry
//         .connect(owner)
//         .setConfigTypeSafe(
//           signerAddresses,
//           keeperAddresses,
//           f,
//           newConfig,
//           offchainVersion,
//           offchainBytes,
//         )
//
//       const updated = await registry.getState()
//       assert.isTrue(updated.state.paused)
//     })
//
//     it('emits an event', async () => {
//       const tx = await registry
//         .connect(owner)
//         .setConfigTypeSafe(
//           signerAddresses,
//           keeperAddresses,
//           f,
//           newConfig,
//           offchainVersion,
//           offchainBytes,
//         )
//       await expect(tx).to.emit(registry, 'ConfigSet')
//     })
//   })
//
//   describe('#setConfig - offchain', () => {
//     let newKeepers: string[]
//
//     beforeEach(async () => {
//       newKeepers = [
//         await personas.Eddy.getAddress(),
//         await personas.Nick.getAddress(),
//         await personas.Neil.getAddress(),
//         await personas.Carol.getAddress(),
//       ]
//     })
//
//     it('reverts when called by anyone but the owner', async () => {
//       await evmRevert(
//         registry
//           .connect(payee1)
//           .setConfigTypeSafe(
//             newKeepers,
//             newKeepers,
//             f,
//             config,
//             offchainVersion,
//             offchainBytes,
//           ),
//         'Only callable by owner',
//       )
//     })
//
//     it('reverts if too many keeperAddresses set', async () => {
//       for (let i = 0; i < 40; i++) {
//         newKeepers.push(randomAddress())
//       }
//       await evmRevert(
//         registry
//           .connect(owner)
//           .setConfigTypeSafe(
//             newKeepers,
//             newKeepers,
//             f,
//             config,
//             offchainVersion,
//             offchainBytes,
//           ),
//         'TooManyOracles()',
//       )
//     })
//
//     it('reverts if f=0', async () => {
//       await evmRevert(
//         registry
//           .connect(owner)
//           .setConfigTypeSafe(
//             newKeepers,
//             newKeepers,
//             0,
//             config,
//             offchainVersion,
//             offchainBytes,
//           ),
//         'IncorrectNumberOfFaultyOracles()',
//       )
//     })
//
//     it('reverts if signers != transmitters length', async () => {
//       const signers = [randomAddress()]
//       await evmRevert(
//         registry
//           .connect(owner)
//           .setConfigTypeSafe(
//             signers,
//             newKeepers,
//             f,
//             config,
//             offchainVersion,
//             offchainBytes,
//           ),
//         'IncorrectNumberOfSigners()',
//       )
//     })
//
//     it('reverts if signers <= 3f', async () => {
//       newKeepers.pop()
//       await evmRevert(
//         registry
//           .connect(owner)
//           .setConfigTypeSafe(
//             newKeepers,
//             newKeepers,
//             f,
//             config,
//             offchainVersion,
//             offchainBytes,
//           ),
//         'IncorrectNumberOfSigners()',
//       )
//     })
//
//     it('reverts on repeated signers', async () => {
//       const newSigners = [
//         await personas.Eddy.getAddress(),
//         await personas.Eddy.getAddress(),
//         await personas.Eddy.getAddress(),
//         await personas.Eddy.getAddress(),
//       ]
//       await evmRevert(
//         registry
//           .connect(owner)
//           .setConfigTypeSafe(
//             newSigners,
//             newKeepers,
//             f,
//             config,
//             offchainVersion,
//             offchainBytes,
//           ),
//         'RepeatedSigner()',
//       )
//     })
//
//     it('reverts on repeated transmitters', async () => {
//       const newTransmitters = [
//         await personas.Eddy.getAddress(),
//         await personas.Eddy.getAddress(),
//         await personas.Eddy.getAddress(),
//         await personas.Eddy.getAddress(),
//       ]
//       await evmRevert(
//         registry
//           .connect(owner)
//           .setConfigTypeSafe(
//             newKeepers,
//             newTransmitters,
//             f,
//             config,
//             offchainVersion,
//             offchainBytes,
//           ),
//         'RepeatedTransmitter()',
//       )
//     })
//
//     itMaybe('stores new config and emits event', async () => {
//       // Perform an upkeep so that totalPremium is updated
//       await registry.connect(admin).addFunds(upkeepId, toWei('100'))
//       let tx = await getTransmitTx(registry, keeper1, [upkeepId])
//       await tx.wait()
//
//       const newOffChainVersion = BigNumber.from('2')
//       const newOffChainConfig = '0x1122'
//
//       const old = await registry.getState()
//       const oldState = old.state
//       assert(oldState.totalPremium.gt(BigNumber.from('0')))
//
//       const newSigners = newKeepers
//       tx = await registry
//         .connect(owner)
//         .setConfigTypeSafe(
//           newSigners,
//           newKeepers,
//           f,
//           config,
//           newOffChainVersion,
//           newOffChainConfig,
//         )
//
//       const updated = await registry.getState()
//       const updatedState = updated.state
//       assert(oldState.totalPremium.eq(updatedState.totalPremium))
//
//       // Old signer addresses which are not in new signers should be non active
//       for (let i = 0; i < signerAddresses.length; i++) {
//         const signer = signerAddresses[i]
//         if (!newSigners.includes(signer)) {
//           assert((await registry.getSignerInfo(signer)).active == false)
//           assert((await registry.getSignerInfo(signer)).index == 0)
//         }
//       }
//       // New signer addresses should be active
//       for (let i = 0; i < newSigners.length; i++) {
//         const signer = newSigners[i]
//         assert((await registry.getSignerInfo(signer)).active == true)
//         assert((await registry.getSignerInfo(signer)).index == i)
//       }
//       // Old transmitter addresses which are not in new transmitter should be non active, update lastCollected but retain other info
//       for (let i = 0; i < keeperAddresses.length; i++) {
//         const transmitter = keeperAddresses[i]
//         if (!newKeepers.includes(transmitter)) {
//           assert(
//             (await registry.getTransmitterInfo(transmitter)).active == false,
//           )
//           assert((await registry.getTransmitterInfo(transmitter)).index == i)
//           assert(
//             (await registry.getTransmitterInfo(transmitter)).lastCollected.eq(
//               oldState.totalPremium.sub(
//                 oldState.totalPremium.mod(keeperAddresses.length),
//               ),
//             ),
//           )
//         }
//       }
//       // New transmitter addresses should be active
//       for (let i = 0; i < newKeepers.length; i++) {
//         const transmitter = newKeepers[i]
//         assert((await registry.getTransmitterInfo(transmitter)).active == true)
//         assert((await registry.getTransmitterInfo(transmitter)).index == i)
//         assert(
//           (await registry.getTransmitterInfo(transmitter)).lastCollected.eq(
//             oldState.totalPremium,
//           ),
//         )
//       }
//
//       // config digest should be updated
//       assert(oldState.configCount + 1 == updatedState.configCount)
//       assert(
//         oldState.latestConfigBlockNumber !=
//           updatedState.latestConfigBlockNumber,
//       )
//       assert(oldState.latestConfigDigest != updatedState.latestConfigDigest)
//
//       //New config should be updated
//       assert.deepEqual(updated.signers, newKeepers)
//       assert.deepEqual(updated.transmitters, newKeepers)
//
//       // Event should have been emitted
//       await expect(tx).to.emit(registry, 'ConfigSet')
//     })
//   })
//
//   describe('#setPeerRegistryMigrationPermission() / #getPeerRegistryMigrationPermission()', () => {
//     const peer = randomAddress()
//     it('allows the owner to set the peer registries', async () => {
//       let permission = await registry.getPeerRegistryMigrationPermission(peer)
//       expect(permission).to.equal(0)
//       await registry.setPeerRegistryMigrationPermission(peer, 1)
//       permission = await registry.getPeerRegistryMigrationPermission(peer)
//       expect(permission).to.equal(1)
//       await registry.setPeerRegistryMigrationPermission(peer, 2)
//       permission = await registry.getPeerRegistryMigrationPermission(peer)
//       expect(permission).to.equal(2)
//       await registry.setPeerRegistryMigrationPermission(peer, 0)
//       permission = await registry.getPeerRegistryMigrationPermission(peer)
//       expect(permission).to.equal(0)
//     })
//     it('reverts if passed an unsupported permission', async () => {
//       await expect(
//         registry.connect(admin).setPeerRegistryMigrationPermission(peer, 10),
//       ).to.be.reverted
//     })
//     it('reverts if not called by the owner', async () => {
//       await expect(
//         registry.connect(admin).setPeerRegistryMigrationPermission(peer, 1),
//       ).to.be.revertedWith('Only callable by owner')
//     })
//   })
//
//   describe('#registerUpkeep', () => {
//     it('reverts when registry is paused', async () => {
//       await registry.connect(owner).pause()
//       await evmRevert(
//         registry
//           .connect(owner)
//           [
//             'registerUpkeep(address,uint32,address,bytes,bytes)'
//           ](mock.address, performGas, await admin.getAddress(), emptyBytes, '0x'),
//         'RegistryPaused()',
//       )
//     })
//
//     it('reverts if the target is not a contract', async () => {
//       await evmRevert(
//         registry
//           .connect(owner)
//           [
//             'registerUpkeep(address,uint32,address,bytes,bytes)'
//           ](zeroAddress, performGas, await admin.getAddress(), emptyBytes, '0x'),
//         'NotAContract()',
//       )
//     })
//
//     it('reverts if called by a non-owner', async () => {
//       await evmRevert(
//         registry
//           .connect(keeper1)
//           [
//             'registerUpkeep(address,uint32,address,bytes,bytes)'
//           ](mock.address, performGas, await admin.getAddress(), emptyBytes, '0x'),
//         'OnlyCallableByOwnerOrRegistrar()',
//       )
//     })
//
//     it('reverts if execute gas is too low', async () => {
//       await evmRevert(
//         registry
//           .connect(owner)
//           [
//             'registerUpkeep(address,uint32,address,bytes,bytes)'
//           ](mock.address, 2299, await admin.getAddress(), emptyBytes, '0x'),
//         'GasLimitOutsideRange()',
//       )
//     })
//
//     it('reverts if execute gas is too high', async () => {
//       await evmRevert(
//         registry
//           .connect(owner)
//           [
//             'registerUpkeep(address,uint32,address,bytes,bytes)'
//           ](mock.address, 5000001, await admin.getAddress(), emptyBytes, '0x'),
//         'GasLimitOutsideRange()',
//       )
//     })
//
//     it('reverts if checkData is too long', async () => {
//       let longBytes = '0x'
//       for (let i = 0; i < 10000; i++) {
//         longBytes += '1'
//       }
//       await evmRevert(
//         registry
//           .connect(owner)
//           [
//             'registerUpkeep(address,uint32,address,bytes,bytes)'
//           ](mock.address, performGas, await admin.getAddress(), longBytes, '0x'),
//         'CheckDataExceedsLimit()',
//       )
//     })
//
//     it('creates a record of the registration', async () => {
//       const performGases = [100000, 500000]
//       const checkDatas = [emptyBytes, '0x12']
//
//       for (let jdx = 0; jdx < performGases.length; jdx++) {
//         const performGas = performGases[jdx]
//         for (let kdx = 0; kdx < checkDatas.length; kdx++) {
//           const checkData = checkDatas[kdx]
//           const tx = await registry
//             .connect(owner)
//             [
//               'registerUpkeep(address,uint32,address,bytes,bytes)'
//             ](mock.address, performGas, await admin.getAddress(), checkData, '0x')
//
//           //confirm the upkeep details and verify emitted events
//           const testUpkeepId = await getUpkeepID(tx)
//           await expect(tx)
//             .to.emit(registry, 'UpkeepRegistered')
//             .withArgs(testUpkeepId, performGas, await admin.getAddress())
//
//           await expect(tx)
//             .to.emit(registry, 'UpkeepCheckDataSet')
//             .withArgs(testUpkeepId, checkData)
//           await expect(tx)
//             .to.emit(registry, 'UpkeepTriggerConfigSet')
//             .withArgs(testUpkeepId, '0x')
//
//           const registration = await registry.getUpkeep(testUpkeepId)
//
//           assert.equal(mock.address, registration.target)
//           assert.notEqual(
//             ethers.constants.AddressZero,
//             await registry.getForwarder(testUpkeepId),
//           )
//           assert.equal(
//             performGas.toString(),
//             registration.performGas.toString(),
//           )
//           assert.equal(await admin.getAddress(), registration.admin)
//           assert.equal(0, registration.balance.toNumber())
//           assert.equal(0, registration.amountSpent.toNumber())
//           assert.equal(0, registration.lastPerformedBlockNumber)
//           assert.equal(checkData, registration.checkData)
//           assert.equal(registration.paused, false)
//           assert.equal(registration.offchainConfig, '0x')
//           assert(registration.maxValidBlocknumber.eq('0xffffffff'))
//         }
//       }
//     })
//   })
//
//   describe('#pauseUpkeep', () => {
//     it('reverts if the registration does not exist', async () => {
//       await evmRevert(
//         registry.connect(keeper1).pauseUpkeep(upkeepId.add(1)),
//         'OnlyCallableByAdmin()',
//       )
//     })
//
//     it('reverts if the upkeep is already canceled', async () => {
//       await registry.connect(admin).cancelUpkeep(upkeepId)
//
//       await evmRevert(
//         registry.connect(admin).pauseUpkeep(upkeepId),
//         'UpkeepCancelled()',
//       )
//     })
//
//     it('reverts if the upkeep is already paused', async () => {
//       await registry.connect(admin).pauseUpkeep(upkeepId)
//
//       await evmRevert(
//         registry.connect(admin).pauseUpkeep(upkeepId),
//         'OnlyUnpausedUpkeep()',
//       )
//     })
//
//     it('reverts if the caller is not the upkeep admin', async () => {
//       await evmRevert(
//         registry.connect(keeper1).pauseUpkeep(upkeepId),
//         'OnlyCallableByAdmin()',
//       )
//     })
//
//     it('pauses the upkeep and emits an event', async () => {
//       const tx = await registry.connect(admin).pauseUpkeep(upkeepId)
//       await expect(tx).to.emit(registry, 'UpkeepPaused').withArgs(upkeepId)
//
//       const registration = await registry.getUpkeep(upkeepId)
//       assert.equal(registration.paused, true)
//     })
//   })
//
//   describe('#unpauseUpkeep', () => {
//     it('reverts if the registration does not exist', async () => {
//       await evmRevert(
//         registry.connect(keeper1).unpauseUpkeep(upkeepId.add(1)),
//         'OnlyCallableByAdmin()',
//       )
//     })
//
//     it('reverts if the upkeep is already canceled', async () => {
//       await registry.connect(owner).cancelUpkeep(upkeepId)
//
//       await evmRevert(
//         registry.connect(admin).unpauseUpkeep(upkeepId),
//         'UpkeepCancelled()',
//       )
//     })
//
//     it('marks the contract as paused', async () => {
//       assert.isFalse((await registry.getState()).state.paused)
//
//       await registry.connect(owner).pause()
//
//       assert.isTrue((await registry.getState()).state.paused)
//     })
//
//     it('reverts if the upkeep is not paused', async () => {
//       await evmRevert(
//         registry.connect(admin).unpauseUpkeep(upkeepId),
//         'OnlyPausedUpkeep()',
//       )
//     })
//
//     it('reverts if the caller is not the upkeep admin', async () => {
//       await registry.connect(admin).pauseUpkeep(upkeepId)
//
//       const registration = await registry.getUpkeep(upkeepId)
//
//       assert.equal(registration.paused, true)
//
//       await evmRevert(
//         registry.connect(keeper1).unpauseUpkeep(upkeepId),
//         'OnlyCallableByAdmin()',
//       )
//     })
//
//     it('unpauses the upkeep and emits an event', async () => {
//       const originalCount = (await registry.getActiveUpkeepIDs(0, 0)).length
//
//       await registry.connect(admin).pauseUpkeep(upkeepId)
//
//       const tx = await registry.connect(admin).unpauseUpkeep(upkeepId)
//
//       await expect(tx).to.emit(registry, 'UpkeepUnpaused').withArgs(upkeepId)
//
//       const registration = await registry.getUpkeep(upkeepId)
//       assert.equal(registration.paused, false)
//
//       const upkeepIds = await registry.getActiveUpkeepIDs(0, 0)
//       assert.equal(upkeepIds.length, originalCount)
//     })
//   })
//
//   describe('#setUpkeepCheckData', () => {
//     it('reverts if the registration does not exist', async () => {
//       await evmRevert(
//         registry
//           .connect(keeper1)
//           .setUpkeepCheckData(upkeepId.add(1), randomBytes),
//         'OnlyCallableByAdmin()',
//       )
//     })
//
//     it('reverts if the caller is not upkeep admin', async () => {
//       await evmRevert(
//         registry.connect(keeper1).setUpkeepCheckData(upkeepId, randomBytes),
//         'OnlyCallableByAdmin()',
//       )
//     })
//
//     it('reverts if the upkeep is cancelled', async () => {
//       await registry.connect(admin).cancelUpkeep(upkeepId)
//
//       await evmRevert(
//         registry.connect(admin).setUpkeepCheckData(upkeepId, randomBytes),
//         'UpkeepCancelled()',
//       )
//     })
//
//     it('is allowed to update on paused upkeep', async () => {
//       await registry.connect(admin).pauseUpkeep(upkeepId)
//       await registry.connect(admin).setUpkeepCheckData(upkeepId, randomBytes)
//
//       const registration = await registry.getUpkeep(upkeepId)
//       assert.equal(randomBytes, registration.checkData)
//     })
//
//     it('reverts if new data exceeds limit', async () => {
//       let longBytes = '0x'
//       for (let i = 0; i < 10000; i++) {
//         longBytes += '1'
//       }
//
//       await evmRevert(
//         registry.connect(admin).setUpkeepCheckData(upkeepId, longBytes),
//         'CheckDataExceedsLimit()',
//       )
//     })
//
//     it('updates the upkeep check data and emits an event', async () => {
//       const tx = await registry
//         .connect(admin)
//         .setUpkeepCheckData(upkeepId, randomBytes)
//       await expect(tx)
//         .to.emit(registry, 'UpkeepCheckDataSet')
//         .withArgs(upkeepId, randomBytes)
//
//       const registration = await registry.getUpkeep(upkeepId)
//       assert.equal(randomBytes, registration.checkData)
//     })
//   })
//
//   describe('#setUpkeepGasLimit', () => {
//     const newGasLimit = BigNumber.from('300000')
//
//     it('reverts if the registration does not exist', async () => {
//       await evmRevert(
//         registry.connect(admin).setUpkeepGasLimit(upkeepId.add(1), newGasLimit),
//         'OnlyCallableByAdmin()',
//       )
//     })
//
//     it('reverts if the upkeep is canceled', async () => {
//       await registry.connect(admin).cancelUpkeep(upkeepId)
//       await evmRevert(
//         registry.connect(admin).setUpkeepGasLimit(upkeepId, newGasLimit),
//         'UpkeepCancelled()',
//       )
//     })
//
//     it('reverts if called by anyone but the admin', async () => {
//       await evmRevert(
//         registry.connect(owner).setUpkeepGasLimit(upkeepId, newGasLimit),
//         'OnlyCallableByAdmin()',
//       )
//     })
//
//     it('reverts if new gas limit is out of bounds', async () => {
//       await evmRevert(
//         registry
//           .connect(admin)
//           .setUpkeepGasLimit(upkeepId, BigNumber.from('100')),
//         'GasLimitOutsideRange()',
//       )
//       await evmRevert(
//         registry
//           .connect(admin)
//           .setUpkeepGasLimit(upkeepId, BigNumber.from('6000000')),
//         'GasLimitOutsideRange()',
//       )
//     })
//
//     it('updates the gas limit successfully', async () => {
//       const initialGasLimit = (await registry.getUpkeep(upkeepId)).performGas
//       assert.equal(initialGasLimit, performGas.toNumber())
//       await registry.connect(admin).setUpkeepGasLimit(upkeepId, newGasLimit)
//       const updatedGasLimit = (await registry.getUpkeep(upkeepId)).performGas
//       assert.equal(updatedGasLimit, newGasLimit.toNumber())
//     })
//
//     it('emits a log', async () => {
//       const tx = await registry
//         .connect(admin)
//         .setUpkeepGasLimit(upkeepId, newGasLimit)
//       await expect(tx)
//         .to.emit(registry, 'UpkeepGasLimitSet')
//         .withArgs(upkeepId, newGasLimit)
//     })
//   })
//
//   describe('#setUpkeepOffchainConfig', () => {
//     const newConfig = '0xc0ffeec0ffee'
//
//     it('reverts if the registration does not exist', async () => {
//       await evmRevert(
//         registry
//           .connect(admin)
//           .setUpkeepOffchainConfig(upkeepId.add(1), newConfig),
//         'OnlyCallableByAdmin()',
//       )
//     })
//
//     it('reverts if the upkeep is canceled', async () => {
//       await registry.connect(admin).cancelUpkeep(upkeepId)
//       await evmRevert(
//         registry.connect(admin).setUpkeepOffchainConfig(upkeepId, newConfig),
//         'UpkeepCancelled()',
//       )
//     })
//
//     it('reverts if called by anyone but the admin', async () => {
//       await evmRevert(
//         registry.connect(owner).setUpkeepOffchainConfig(upkeepId, newConfig),
//         'OnlyCallableByAdmin()',
//       )
//     })
//
//     it('updates the config successfully', async () => {
//       const initialConfig = (await registry.getUpkeep(upkeepId)).offchainConfig
//       assert.equal(initialConfig, '0x')
//       await registry.connect(admin).setUpkeepOffchainConfig(upkeepId, newConfig)
//       const updatedConfig = (await registry.getUpkeep(upkeepId)).offchainConfig
//       assert.equal(newConfig, updatedConfig)
//     })
//
//     it('emits a log', async () => {
//       const tx = await registry
//         .connect(admin)
//         .setUpkeepOffchainConfig(upkeepId, newConfig)
//       await expect(tx)
//         .to.emit(registry, 'UpkeepOffchainConfigSet')
//         .withArgs(upkeepId, newConfig)
//     })
//   })
//
//   describe('#setUpkeepTriggerConfig', () => {
//     const newConfig = '0xdeadbeef'
//
//     it('reverts if the registration does not exist', async () => {
//       await evmRevert(
//         registry
//           .connect(admin)
//           .setUpkeepTriggerConfig(upkeepId.add(1), newConfig),
//         'OnlyCallableByAdmin()',
//       )
//     })
//
//     it('reverts if the upkeep is canceled', async () => {
//       await registry.connect(admin).cancelUpkeep(upkeepId)
//       await evmRevert(
//         registry.connect(admin).setUpkeepTriggerConfig(upkeepId, newConfig),
//         'UpkeepCancelled()',
//       )
//     })
//
//     it('reverts if called by anyone but the admin', async () => {
//       await evmRevert(
//         registry.connect(owner).setUpkeepTriggerConfig(upkeepId, newConfig),
//         'OnlyCallableByAdmin()',
//       )
//     })
//
//     it('emits a log', async () => {
//       const tx = await registry
//         .connect(admin)
//         .setUpkeepTriggerConfig(upkeepId, newConfig)
//       await expect(tx)
//         .to.emit(registry, 'UpkeepTriggerConfigSet')
//         .withArgs(upkeepId, newConfig)
//     })
//   })
//
//   describe('#transferUpkeepAdmin', () => {
//     it('reverts when called by anyone but the current upkeep admin', async () => {
//       await evmRevert(
//         registry
//           .connect(payee1)
//           .transferUpkeepAdmin(upkeepId, await payee2.getAddress()),
//         'OnlyCallableByAdmin()',
//       )
//     })
//
//     it('reverts when transferring to self', async () => {
//       await evmRevert(
//         registry
//           .connect(admin)
//           .transferUpkeepAdmin(upkeepId, await admin.getAddress()),
//         'ValueNotChanged()',
//       )
//     })
//
//     it('reverts when the upkeep is cancelled', async () => {
//       await registry.connect(admin).cancelUpkeep(upkeepId)
//
//       await evmRevert(
//         registry
//           .connect(admin)
//           .transferUpkeepAdmin(upkeepId, await keeper1.getAddress()),
//         'UpkeepCancelled()',
//       )
//     })
//
//     it('allows cancelling transfer by reverting to zero address', async () => {
//       await registry
//         .connect(admin)
//         .transferUpkeepAdmin(upkeepId, await payee1.getAddress())
//       const tx = await registry
//         .connect(admin)
//         .transferUpkeepAdmin(upkeepId, ethers.constants.AddressZero)
//
//       await expect(tx)
//         .to.emit(registry, 'UpkeepAdminTransferRequested')
//         .withArgs(
//           upkeepId,
//           await admin.getAddress(),
//           ethers.constants.AddressZero,
//         )
//     })
//
//     it('does not change the upkeep admin', async () => {
//       await registry
//         .connect(admin)
//         .transferUpkeepAdmin(upkeepId, await payee1.getAddress())
//
//       const upkeep = await registry.getUpkeep(upkeepId)
//       assert.equal(await admin.getAddress(), upkeep.admin)
//     })
//
//     it('emits an event announcing the new upkeep admin', async () => {
//       const tx = await registry
//         .connect(admin)
//         .transferUpkeepAdmin(upkeepId, await payee1.getAddress())
//
//       await expect(tx)
//         .to.emit(registry, 'UpkeepAdminTransferRequested')
//         .withArgs(upkeepId, await admin.getAddress(), await payee1.getAddress())
//     })
//
//     it('does not emit an event when called with the same proposed upkeep admin', async () => {
//       await registry
//         .connect(admin)
//         .transferUpkeepAdmin(upkeepId, await payee1.getAddress())
//
//       const tx = await registry
//         .connect(admin)
//         .transferUpkeepAdmin(upkeepId, await payee1.getAddress())
//       const receipt = await tx.wait()
//       assert.equal(0, receipt.logs.length)
//     })
//   })
//
//   describe('#acceptUpkeepAdmin', () => {
//     beforeEach(async () => {
//       // Start admin transfer to payee1
//       await registry
//         .connect(admin)
//         .transferUpkeepAdmin(upkeepId, await payee1.getAddress())
//     })
//
//     it('reverts when not called by the proposed upkeep admin', async () => {
//       await evmRevert(
//         registry.connect(payee2).acceptUpkeepAdmin(upkeepId),
//         'OnlyCallableByProposedAdmin()',
//       )
//     })
//
//     it('reverts when the upkeep is cancelled', async () => {
//       await registry.connect(admin).cancelUpkeep(upkeepId)
//
//       await evmRevert(
//         registry.connect(payee1).acceptUpkeepAdmin(upkeepId),
//         'UpkeepCancelled()',
//       )
//     })
//
//     it('does change the admin', async () => {
//       await registry.connect(payee1).acceptUpkeepAdmin(upkeepId)
//
//       const upkeep = await registry.getUpkeep(upkeepId)
//       assert.equal(await payee1.getAddress(), upkeep.admin)
//     })
//
//     it('emits an event announcing the new upkeep admin', async () => {
//       const tx = await registry.connect(payee1).acceptUpkeepAdmin(upkeepId)
//       await expect(tx)
//         .to.emit(registry, 'UpkeepAdminTransferred')
//         .withArgs(upkeepId, await admin.getAddress(), await payee1.getAddress())
//     })
//   })
//
//   describe('#withdrawOwnerFunds', () => {
//     it('can only be called by owner', async () => {
//       await evmRevert(
//         registry.connect(keeper1).withdrawOwnerFunds(),
//         'Only callable by owner',
//       )
//     })
//
//     itMaybe('withdraws the collected fees to owner', async () => {
//       await registry.connect(admin).addFunds(upkeepId, toWei('100'))
//       // Very high min spend, whole balance as cancellation fees
//       const minUpkeepSpend = toWei('1000')
//       await registry.connect(owner).setConfigTypeSafe(
//         signerAddresses,
//         keeperAddresses,
//         f,
//         {
//           paymentPremiumPPB,
//           flatFeeMicroLink,
//           checkGasLimit,
//           stalenessSeconds,
//           gasCeilingMultiplier,
//           minUpkeepSpend,
//           maxCheckDataSize,
//           maxPerformDataSize,
//           maxRevertDataSize,
//           maxPerformGas,
//           fallbackGasPrice,
//           fallbackLinkPrice,
//           transcoder: transcoder.address,
//           registrars: [],
//           upkeepPrivilegeManager: upkeepManager,
//         },
//         offchainVersion,
//         offchainBytes,
//       )
//       const upkeepBalance = (await registry.getUpkeep(upkeepId)).balance
//       const ownerBefore = await linkToken.balanceOf(await owner.getAddress())
//
//       await registry.connect(owner).cancelUpkeep(upkeepId)
//
//       // Transfered to owner balance on registry
//       let ownerRegistryBalance = (await registry.getState()).state
//         .ownerLinkBalance
//       assert.isTrue(ownerRegistryBalance.eq(upkeepBalance))
//
//       // Now withdraw
//       await registry.connect(owner).withdrawOwnerFunds()
//
//       ownerRegistryBalance = (await registry.getState()).state.ownerLinkBalance
//       const ownerAfter = await linkToken.balanceOf(await owner.getAddress())
//
//       // Owner registry balance should be changed to 0
//       assert.isTrue(ownerRegistryBalance.eq(BigNumber.from('0')))
//
//       // Owner should be credited with the balance
//       assert.isTrue(ownerBefore.add(upkeepBalance).eq(ownerAfter))
//     })
//   })
//
//   describe('#transferPayeeship', () => {
//     it('reverts when called by anyone but the current payee', async () => {
//       await evmRevert(
//         registry
//           .connect(payee2)
//           .transferPayeeship(
//             await keeper1.getAddress(),
//             await payee2.getAddress(),
//           ),
//         'OnlyCallableByPayee()',
//       )
//     })
//
//     it('reverts when transferring to self', async () => {
//       await evmRevert(
//         registry
//           .connect(payee1)
//           .transferPayeeship(
//             await keeper1.getAddress(),
//             await payee1.getAddress(),
//           ),
//         'ValueNotChanged()',
//       )
//     })
//
//     it('does not change the payee', async () => {
//       await registry
//         .connect(payee1)
//         .transferPayeeship(
//           await keeper1.getAddress(),
//           await payee2.getAddress(),
//         )
//
//       const info = await registry.getTransmitterInfo(await keeper1.getAddress())
//       assert.equal(await payee1.getAddress(), info.payee)
//     })
//
//     it('emits an event announcing the new payee', async () => {
//       const tx = await registry
//         .connect(payee1)
//         .transferPayeeship(
//           await keeper1.getAddress(),
//           await payee2.getAddress(),
//         )
//       await expect(tx)
//         .to.emit(registry, 'PayeeshipTransferRequested')
//         .withArgs(
//           await keeper1.getAddress(),
//           await payee1.getAddress(),
//           await payee2.getAddress(),
//         )
//     })
//
//     it('does not emit an event when called with the same proposal', async () => {
//       await registry
//         .connect(payee1)
//         .transferPayeeship(
//           await keeper1.getAddress(),
//           await payee2.getAddress(),
//         )
//
//       const tx = await registry
//         .connect(payee1)
//         .transferPayeeship(
//           await keeper1.getAddress(),
//           await payee2.getAddress(),
//         )
//       const receipt = await tx.wait()
//       assert.equal(0, receipt.logs.length)
//     })
//   })
//
//   describe('#acceptPayeeship', () => {
//     beforeEach(async () => {
//       await registry
//         .connect(payee1)
//         .transferPayeeship(
//           await keeper1.getAddress(),
//           await payee2.getAddress(),
//         )
//     })
//
//     it('reverts when called by anyone but the proposed payee', async () => {
//       await evmRevert(
//         registry.connect(payee1).acceptPayeeship(await keeper1.getAddress()),
//         'OnlyCallableByProposedPayee()',
//       )
//     })
//
//     it('emits an event announcing the new payee', async () => {
//       const tx = await registry
//         .connect(payee2)
//         .acceptPayeeship(await keeper1.getAddress())
//       await expect(tx)
//         .to.emit(registry, 'PayeeshipTransferred')
//         .withArgs(
//           await keeper1.getAddress(),
//           await payee1.getAddress(),
//           await payee2.getAddress(),
//         )
//     })
//
//     it('does change the payee', async () => {
//       await registry.connect(payee2).acceptPayeeship(await keeper1.getAddress())
//
//       const info = await registry.getTransmitterInfo(await keeper1.getAddress())
//       assert.equal(await payee2.getAddress(), info.payee)
//     })
//   })
//
//   describe('#pause', () => {
//     it('reverts if called by a non-owner', async () => {
//       await evmRevert(
//         registry.connect(keeper1).pause(),
//         'Only callable by owner',
//       )
//     })
//
//     it('marks the contract as paused', async () => {
//       assert.isFalse((await registry.getState()).state.paused)
//
//       await registry.connect(owner).pause()
//
//       assert.isTrue((await registry.getState()).state.paused)
//     })
//
//     it('Does not allow transmits when paused', async () => {
//       await registry.connect(owner).pause()
//
//       await evmRevert(
//         getTransmitTx(registry, keeper1, [upkeepId]),
//         'RegistryPaused()',
//       )
//     })
//
//     it('Does not allow creation of new upkeeps when paused', async () => {
//       await registry.connect(owner).pause()
//
//       await evmRevert(
//         registry
//           .connect(owner)
//           [
//             'registerUpkeep(address,uint32,address,bytes,bytes)'
//           ](mock.address, performGas, await admin.getAddress(), emptyBytes, '0x'),
//         'RegistryPaused()',
//       )
//     })
//   })
//
//   describe('#unpause', () => {
//     beforeEach(async () => {
//       await registry.connect(owner).pause()
//     })
//
//     it('reverts if called by a non-owner', async () => {
//       await evmRevert(
//         registry.connect(keeper1).unpause(),
//         'Only callable by owner',
//       )
//     })
//
//     it('marks the contract as not paused', async () => {
//       assert.isTrue((await registry.getState()).state.paused)
//
//       await registry.connect(owner).unpause()
//
//       assert.isFalse((await registry.getState()).state.paused)
//     })
//   })
//
//   describe('#migrateUpkeeps() / #receiveUpkeeps()', async () => {
//     context('when permissions are set', () => {
//       beforeEach(async () => {
//         await linkToken.connect(owner).approve(registry.address, toWei('100'))
//         await registry.connect(owner).addFunds(upkeepId, toWei('100'))
//         await registry.setPeerRegistryMigrationPermission(mgRegistry.address, 1)
//         await mgRegistry.setPeerRegistryMigrationPermission(registry.address, 2)
//       })
//
//       it('migrates an upkeep', async () => {
//         const offchainBytes = '0x987654abcd'
//         await registry
//           .connect(admin)
//           .setUpkeepOffchainConfig(upkeepId, offchainBytes)
//         const reg1Upkeep = await registry.getUpkeep(upkeepId)
//         const forwarderAddress = await registry.getForwarder(upkeepId)
//         expect(reg1Upkeep.balance).to.equal(toWei('100'))
//         expect(reg1Upkeep.checkData).to.equal(randomBytes)
//         expect(forwarderAddress).to.not.equal(ethers.constants.AddressZero)
//         expect(reg1Upkeep.offchainConfig).to.equal(offchainBytes)
//         expect((await registry.getState()).state.numUpkeeps).to.equal(
//           numUpkeeps,
//         )
//         const forwarder = await IAutomationForwarderFactory.connect(
//           forwarderAddress,
//           owner,
//         )
//         expect(await forwarder.getRegistry()).to.equal(registry.address)
//         // Set an upkeep admin transfer in progress too
//         await registry
//           .connect(admin)
//           .transferUpkeepAdmin(upkeepId, await payee1.getAddress())
//
//         // migrate
//         await registry
//           .connect(admin)
//           .migrateUpkeeps([upkeepId], mgRegistry.address)
//         expect((await registry.getState()).state.numUpkeeps).to.equal(
//           numUpkeeps - 1,
//         )
//         expect((await mgRegistry.getState()).state.numUpkeeps).to.equal(1)
//         expect((await registry.getUpkeep(upkeepId)).balance).to.equal(0)
//         expect((await registry.getUpkeep(upkeepId)).checkData).to.equal('0x')
//         expect((await mgRegistry.getUpkeep(upkeepId)).balance).to.equal(
//           toWei('100'),
//         )
//         expect(
//           (await mgRegistry.getState()).state.expectedLinkBalance,
//         ).to.equal(toWei('100'))
//         expect((await mgRegistry.getUpkeep(upkeepId)).checkData).to.equal(
//           randomBytes,
//         )
//         expect((await mgRegistry.getUpkeep(upkeepId)).offchainConfig).to.equal(
//           offchainBytes,
//         )
//         expect(await mgRegistry.getForwarder(upkeepId)).to.equal(
//           forwarderAddress,
//         )
//         // test that registry is updated on forwarder
//         expect(await forwarder.getRegistry()).to.equal(mgRegistry.address)
//         // migration will delete the upkeep and nullify admin transfer
//         await expect(
//           registry.connect(payee1).acceptUpkeepAdmin(upkeepId),
//         ).to.be.revertedWith('UpkeepCancelled()')
//         await expect(
//           mgRegistry.connect(payee1).acceptUpkeepAdmin(upkeepId),
//         ).to.be.revertedWith('OnlyCallableByProposedAdmin()')
//       })
//
//       it('migrates a paused upkeep', async () => {
//         expect((await registry.getUpkeep(upkeepId)).balance).to.equal(
//           toWei('100'),
//         )
//         expect((await registry.getUpkeep(upkeepId)).checkData).to.equal(
//           randomBytes,
//         )
//         expect((await registry.getState()).state.numUpkeeps).to.equal(
//           numUpkeeps,
//         )
//         await registry.connect(admin).pauseUpkeep(upkeepId)
//         // verify the upkeep is paused
//         expect((await registry.getUpkeep(upkeepId)).paused).to.equal(true)
//         // migrate
//         await registry
//           .connect(admin)
//           .migrateUpkeeps([upkeepId], mgRegistry.address)
//         expect((await registry.getState()).state.numUpkeeps).to.equal(
//           numUpkeeps - 1,
//         )
//         expect((await mgRegistry.getState()).state.numUpkeeps).to.equal(1)
//         expect((await registry.getUpkeep(upkeepId)).balance).to.equal(0)
//         expect((await mgRegistry.getUpkeep(upkeepId)).balance).to.equal(
//           toWei('100'),
//         )
//         expect((await registry.getUpkeep(upkeepId)).checkData).to.equal('0x')
//         expect((await mgRegistry.getUpkeep(upkeepId)).checkData).to.equal(
//           randomBytes,
//         )
//         expect(
//           (await mgRegistry.getState()).state.expectedLinkBalance,
//         ).to.equal(toWei('100'))
//         // verify the upkeep is still paused after migration
//         expect((await mgRegistry.getUpkeep(upkeepId)).paused).to.equal(true)
//       })
//
//       it('emits an event on both contracts', async () => {
//         expect((await registry.getUpkeep(upkeepId)).balance).to.equal(
//           toWei('100'),
//         )
//         expect((await registry.getUpkeep(upkeepId)).checkData).to.equal(
//           randomBytes,
//         )
//         expect((await registry.getState()).state.numUpkeeps).to.equal(
//           numUpkeeps,
//         )
//         const tx = registry
//           .connect(admin)
//           .migrateUpkeeps([upkeepId], mgRegistry.address)
//         await expect(tx)
//           .to.emit(registry, 'UpkeepMigrated')
//           .withArgs(upkeepId, toWei('100'), mgRegistry.address)
//         await expect(tx)
//           .to.emit(mgRegistry, 'UpkeepReceived')
//           .withArgs(upkeepId, toWei('100'), registry.address)
//       })
//
//       it('is only migratable by the admin', async () => {
//         await expect(
//           registry
//             .connect(owner)
//             .migrateUpkeeps([upkeepId], mgRegistry.address),
//         ).to.be.revertedWith('OnlyCallableByAdmin()')
//         await registry
//           .connect(admin)
//           .migrateUpkeeps([upkeepId], mgRegistry.address)
//       })
//     })
//
//     context('when permissions are not set', () => {
//       it('reverts', async () => {
//         // no permissions
//         await registry.setPeerRegistryMigrationPermission(mgRegistry.address, 0)
//         await mgRegistry.setPeerRegistryMigrationPermission(registry.address, 0)
//         await expect(registry.migrateUpkeeps([upkeepId], mgRegistry.address)).to
//           .be.reverted
//         // only outgoing permissions
//         await registry.setPeerRegistryMigrationPermission(mgRegistry.address, 1)
//         await mgRegistry.setPeerRegistryMigrationPermission(registry.address, 0)
//         await expect(registry.migrateUpkeeps([upkeepId], mgRegistry.address)).to
//           .be.reverted
//         // only incoming permissions
//         await registry.setPeerRegistryMigrationPermission(mgRegistry.address, 0)
//         await mgRegistry.setPeerRegistryMigrationPermission(registry.address, 2)
//         await expect(registry.migrateUpkeeps([upkeepId], mgRegistry.address)).to
//           .be.reverted
//         // permissions opposite direction
//         await registry.setPeerRegistryMigrationPermission(mgRegistry.address, 2)
//         await mgRegistry.setPeerRegistryMigrationPermission(registry.address, 1)
//         await expect(registry.migrateUpkeeps([upkeepId], mgRegistry.address)).to
//           .be.reverted
//       })
//     })
//   })
//
//   describe('#setPayees', () => {
//     const IGNORE_ADDRESS = '0xFFfFfFffFFfffFFfFFfFFFFFffFFFffffFfFFFfF'
//
//     it('reverts when not called by the owner', async () => {
//       await evmRevert(
//         registry.connect(keeper1).setPayees(payees),
//         'Only callable by owner',
//       )
//     })
//
//     it('reverts with different numbers of payees than transmitters', async () => {
//       await evmRevert(
//         registry.connect(owner).setPayees([...payees, randomAddress()]),
//         'ParameterLengthError()',
//       )
//     })
//
//     it('reverts if the payee is the zero address', async () => {
//       await blankRegistry.connect(owner).setConfig(...baseConfig) // used to test initial config
//
//       await evmRevert(
//         blankRegistry // used to test initial config
//           .connect(owner)
//           .setPayees([ethers.constants.AddressZero, ...payees.slice(1)]),
//         'InvalidPayee()',
//       )
//     })
//
//     itMaybe(
//       'sets the payees when exisitng payees are zero address',
//       async () => {
//         //Initial payees should be zero address
//         await blankRegistry.connect(owner).setConfig(...baseConfig) // used to test initial config
//
//         for (let i = 0; i < keeperAddresses.length; i++) {
//           const payee = (
//             await blankRegistry.getTransmitterInfo(keeperAddresses[i])
//           ).payee // used to test initial config
//           assert.equal(payee, zeroAddress)
//         }
//
//         await blankRegistry.connect(owner).setPayees(payees) // used to test initial config
//
//         for (let i = 0; i < keeperAddresses.length; i++) {
//           const payee = (
//             await blankRegistry.getTransmitterInfo(keeperAddresses[i])
//           ).payee
//           assert.equal(payee, payees[i])
//         }
//       },
//     )
//
//     it('does not change the payee if IGNORE_ADDRESS is used as payee', async () => {
//       const signers = Array.from({ length: 5 }, randomAddress)
//       const keepers = Array.from({ length: 5 }, randomAddress)
//       const payees = Array.from({ length: 5 }, randomAddress)
//       const newTransmitter = randomAddress()
//       const newPayee = randomAddress()
//       const ignoreAddresses = new Array(payees.length).fill(IGNORE_ADDRESS)
//       const newPayees = [...ignoreAddresses, newPayee]
//       // arbitrum registry
//       // configure registry with 5 keepers // optimism registry
//       await blankRegistry // used to test initial configurations
//         .connect(owner)
//         .setConfigTypeSafe(
//           signers,
//           keepers,
//           f,
//           config,
//           offchainVersion,
//           offchainBytes,
//         )
//       // arbitrum registry
//       // set initial payees // optimism registry
//       await blankRegistry.connect(owner).setPayees(payees) // used to test initial configurations
//       // arbitrum registry
//       // add another keeper // optimism registry
//       await blankRegistry // used to test initial configurations
//         .connect(owner)
//         .setConfigTypeSafe(
//           [...signers, randomAddress()],
//           [...keepers, newTransmitter],
//           f,
//           config,
//           offchainVersion,
//           offchainBytes,
//         )
//       // arbitrum registry
//       // update payee list // optimism registry // arbitrum registry
//       await blankRegistry.connect(owner).setPayees(newPayees) // used to test initial configurations // optimism registry
//       const ignored = await blankRegistry.getTransmitterInfo(newTransmitter) // used to test initial configurations
//       assert.equal(newPayee, ignored.payee)
//       assert.equal(true, ignored.active)
//     })
//
//     it('reverts if payee is non zero and owner tries to change payee', async () => {
//       const newPayees = [randomAddress(), ...payees.slice(1)]
//
//       await evmRevert(
//         registry.connect(owner).setPayees(newPayees),
//         'InvalidPayee()',
//       )
//     })
//
//     it('emits events for every payee added and removed', async () => {
//       const tx = await registry.connect(owner).setPayees(payees)
//       await expect(tx)
//         .to.emit(registry, 'PayeesUpdated')
//         .withArgs(keeperAddresses, payees)
//     })
//   })
//
//   describe('#cancelUpkeep', () => {
//     it('reverts if the ID is not valid', async () => {
//       await evmRevert(
//         registry.connect(owner).cancelUpkeep(upkeepId.add(1)),
//         'CannotCancel()',
//       )
//     })
//
//     it('reverts if called by a non-owner/non-admin', async () => {
//       await evmRevert(
//         registry.connect(keeper1).cancelUpkeep(upkeepId),
//         'OnlyCallableByOwnerOrAdmin()',
//       )
//     })
//
//     describe('when called by the owner', async () => {
//       it('sets the registration to invalid immediately', async () => {
//         const tx = await registry.connect(owner).cancelUpkeep(upkeepId)
//         const receipt = await tx.wait()
//         const registration = await registry.getUpkeep(upkeepId)
//         assert.equal(
//           registration.maxValidBlocknumber.toNumber(),
//           receipt.blockNumber,
//         )
//       })
//
//       it('emits an event', async () => {
//         const tx = await registry.connect(owner).cancelUpkeep(upkeepId)
//         const receipt = await tx.wait()
//         await expect(tx)
//           .to.emit(registry, 'UpkeepCanceled')
//           .withArgs(upkeepId, BigNumber.from(receipt.blockNumber))
//       })
//
//       it('immediately prevents upkeep', async () => {
//         await registry.connect(owner).cancelUpkeep(upkeepId)
//
//         const tx = await getTransmitTx(registry, keeper1, [upkeepId])
//         const receipt = await tx.wait()
//         const cancelledUpkeepReportLogs =
//           parseCancelledUpkeepReportLogs(receipt)
//         // exactly 1 CancelledUpkeepReport log should be emitted
//         assert.equal(cancelledUpkeepReportLogs.length, 1)
//       })
//
//       it('does not revert if reverts if called multiple times', async () => {
//         await registry.connect(owner).cancelUpkeep(upkeepId)
//         await evmRevert(
//           registry.connect(owner).cancelUpkeep(upkeepId),
//           'CannotCancel()',
//         )
//       })
//
//       describe('when called by the owner when the admin has just canceled', () => {
//         let oldExpiration: BigNumber
//
//         beforeEach(async () => {
//           await registry.connect(admin).cancelUpkeep(upkeepId)
//           const registration = await registry.getUpkeep(upkeepId)
//           oldExpiration = registration.maxValidBlocknumber
//         })
//
//         it('allows the owner to cancel it more quickly', async () => {
//           await registry.connect(owner).cancelUpkeep(upkeepId)
//
//           const registration = await registry.getUpkeep(upkeepId)
//           const newExpiration = registration.maxValidBlocknumber
//           assert.isTrue(newExpiration.lt(oldExpiration))
//         })
//       })
//     })
//
//     describe('when called by the admin', async () => {
//       it('reverts if called again by the admin', async () => {
//         await registry.connect(admin).cancelUpkeep(upkeepId)
//
//         await evmRevert(
//           registry.connect(admin).cancelUpkeep(upkeepId),
//           'CannotCancel()',
//         )
//       })
//
//       it('reverts if called by the owner after the timeout', async () => {
//         await registry.connect(admin).cancelUpkeep(upkeepId)
//
//         for (let i = 0; i < cancellationDelay; i++) {
//           await ethers.provider.send('evm_mine', [])
//         }
//
//         await evmRevert(
//           registry.connect(owner).cancelUpkeep(upkeepId),
//           'CannotCancel()',
//         )
//       })
//
//       it('sets the registration to invalid in 50 blocks', async () => {
//         const tx = await registry.connect(admin).cancelUpkeep(upkeepId)
//         const receipt = await tx.wait()
//         const registration = await registry.getUpkeep(upkeepId)
//         assert.equal(
//           registration.maxValidBlocknumber.toNumber(),
//           receipt.blockNumber + 50,
//         )
//       })
//
//       it('emits an event', async () => {
//         const tx = await registry.connect(admin).cancelUpkeep(upkeepId)
//         const receipt = await tx.wait()
//         await expect(tx)
//           .to.emit(registry, 'UpkeepCanceled')
//           .withArgs(
//             upkeepId,
//             BigNumber.from(receipt.blockNumber + cancellationDelay),
//           )
//       })
//
//       it('immediately prevents upkeep', async () => {
//         await linkToken.connect(owner).approve(registry.address, toWei('100'))
//         await registry.connect(owner).addFunds(upkeepId, toWei('100'))
//         await registry.connect(admin).cancelUpkeep(upkeepId)
//
//         await getTransmitTx(registry, keeper1, [upkeepId])
//
//         for (let i = 0; i < cancellationDelay; i++) {
//           await ethers.provider.send('evm_mine', [])
//         }
//
//         const tx = await getTransmitTx(registry, keeper1, [upkeepId])
//
//         const receipt = await tx.wait()
//         const cancelledUpkeepReportLogs =
//           parseCancelledUpkeepReportLogs(receipt)
//         // exactly 1 CancelledUpkeepReport log should be emitted
//         assert.equal(cancelledUpkeepReportLogs.length, 1)
//       })
//
//       describeMaybe('when an upkeep has been performed', async () => {
//         beforeEach(async () => {
//           await linkToken.connect(owner).approve(registry.address, toWei('100'))
//           await registry.connect(owner).addFunds(upkeepId, toWei('100'))
//           await getTransmitTx(registry, keeper1, [upkeepId])
//         })
//
//         it('deducts a cancellation fee from the upkeep and gives to owner', async () => {
//           const minUpkeepSpend = toWei('10')
//
//           await registry.connect(owner).setConfigTypeSafe(
//             signerAddresses,
//             keeperAddresses,
//             f,
//             {
//               paymentPremiumPPB,
//               flatFeeMicroLink,
//               checkGasLimit,
//               stalenessSeconds,
//               gasCeilingMultiplier,
//               minUpkeepSpend,
//               maxCheckDataSize,
//               maxPerformDataSize,
//               maxRevertDataSize,
//               maxPerformGas,
//               fallbackGasPrice,
//               fallbackLinkPrice,
//               transcoder: transcoder.address,
//               registrars: [],
//               upkeepPrivilegeManager: upkeepManager,
//             },
//             offchainVersion,
//             offchainBytes,
//           )
//
//           const payee1Before = await linkToken.balanceOf(
//             await payee1.getAddress(),
//           )
//           const upkeepBefore = (await registry.getUpkeep(upkeepId)).balance
//           const ownerBefore = (await registry.getState()).state.ownerLinkBalance
//
//           const amountSpent = toWei('100').sub(upkeepBefore)
//           const cancellationFee = minUpkeepSpend.sub(amountSpent)
//
//           await registry.connect(admin).cancelUpkeep(upkeepId)
//
//           const payee1After = await linkToken.balanceOf(
//             await payee1.getAddress(),
//           )
//           const upkeepAfter = (await registry.getUpkeep(upkeepId)).balance
//           const ownerAfter = (await registry.getState()).state.ownerLinkBalance
//
//           // post upkeep balance should be previous balance minus cancellation fee
//           assert.isTrue(upkeepBefore.sub(cancellationFee).eq(upkeepAfter))
//           // payee balance should not change
//           assert.isTrue(payee1Before.eq(payee1After))
//           // owner should receive the cancellation fee
//           assert.isTrue(ownerAfter.sub(ownerBefore).eq(cancellationFee))
//         })
//
//         it('deducts up to balance as cancellation fee', async () => {
//           // Very high min spend, should deduct whole balance as cancellation fees
//           const minUpkeepSpend = toWei('1000')
//           await registry.connect(owner).setConfigTypeSafe(
//             signerAddresses,
//             keeperAddresses,
//             f,
//             {
//               paymentPremiumPPB,
//               flatFeeMicroLink,
//               checkGasLimit,
//               stalenessSeconds,
//               gasCeilingMultiplier,
//               minUpkeepSpend,
//               maxCheckDataSize,
//               maxPerformDataSize,
//               maxRevertDataSize,
//               maxPerformGas,
//               fallbackGasPrice,
//               fallbackLinkPrice,
//               transcoder: transcoder.address,
//               registrars: [],
//               upkeepPrivilegeManager: upkeepManager,
//             },
//             offchainVersion,
//             offchainBytes,
//           )
//           const payee1Before = await linkToken.balanceOf(
//             await payee1.getAddress(),
//           )
//           const upkeepBefore = (await registry.getUpkeep(upkeepId)).balance
//           const ownerBefore = (await registry.getState()).state.ownerLinkBalance
//
//           await registry.connect(admin).cancelUpkeep(upkeepId)
//           const payee1After = await linkToken.balanceOf(
//             await payee1.getAddress(),
//           )
//           const ownerAfter = (await registry.getState()).state.ownerLinkBalance
//           const upkeepAfter = (await registry.getUpkeep(upkeepId)).balance
//
//           // all upkeep balance is deducted for cancellation fee
//           assert.equal(0, upkeepAfter.toNumber())
//           // payee balance should not change
//           assert.isTrue(payee1After.eq(payee1Before))
//           // all upkeep balance is transferred to the owner
//           assert.isTrue(ownerAfter.sub(ownerBefore).eq(upkeepBefore))
//         })
//
//         it('does not deduct cancellation fee if more than minUpkeepSpend is spent', async () => {
//           // Very low min spend, already spent in one perform upkeep
//           const minUpkeepSpend = BigNumber.from(420)
//           await registry.connect(owner).setConfigTypeSafe(
//             signerAddresses,
//             keeperAddresses,
//             f,
//             {
//               paymentPremiumPPB,
//               flatFeeMicroLink,
//               checkGasLimit,
//               stalenessSeconds,
//               gasCeilingMultiplier,
//               minUpkeepSpend,
//               maxCheckDataSize,
//               maxPerformDataSize,
//               maxRevertDataSize,
//               maxPerformGas,
//               fallbackGasPrice,
//               fallbackLinkPrice,
//               transcoder: transcoder.address,
//               registrars: [],
//               upkeepPrivilegeManager: upkeepManager,
//             },
//             offchainVersion,
//             offchainBytes,
//           )
//           const payee1Before = await linkToken.balanceOf(
//             await payee1.getAddress(),
//           )
//           const upkeepBefore = (await registry.getUpkeep(upkeepId)).balance
//           const ownerBefore = (await registry.getState()).state.ownerLinkBalance
//
//           await registry.connect(admin).cancelUpkeep(upkeepId)
//           const payee1After = await linkToken.balanceOf(
//             await payee1.getAddress(),
//           )
//           const ownerAfter = (await registry.getState()).state.ownerLinkBalance
//           const upkeepAfter = (await registry.getUpkeep(upkeepId)).balance
//
//           // upkeep does not pay cancellation fee after cancellation because minimum upkeep spent is met
//           assert.isTrue(upkeepBefore.eq(upkeepAfter))
//           // owner balance does not change
//           assert.isTrue(ownerAfter.eq(ownerBefore))
//           // payee balance does not change
//           assert.isTrue(payee1Before.eq(payee1After))
//         })
//       })
//     })
//   })
//
//   describe('#withdrawPayment', () => {
//     beforeEach(async () => {
//       await linkToken.connect(owner).approve(registry.address, toWei('100'))
//       await registry.connect(owner).addFunds(upkeepId, toWei('100'))
//       await getTransmitTx(registry, keeper1, [upkeepId])
//     })
//
//     it('reverts if called by anyone but the payee', async () => {
//       await evmRevert(
//         registry
//           .connect(payee2)
//           .withdrawPayment(
//             await keeper1.getAddress(),
//             await nonkeeper.getAddress(),
//           ),
//         'OnlyCallableByPayee()',
//       )
//     })
//
//     it('reverts if called with the 0 address', async () => {
//       await evmRevert(
//         registry
//           .connect(payee2)
//           .withdrawPayment(await keeper1.getAddress(), zeroAddress),
//         'InvalidRecipient()',
//       )
//     })
//
//     it('updates the balances', async () => {
//       const to = await nonkeeper.getAddress()
//       const keeperBefore = await registry.getTransmitterInfo(
//         await keeper1.getAddress(),
//       )
//       const registrationBefore = (await registry.getUpkeep(upkeepId)).balance
//       const toLinkBefore = await linkToken.balanceOf(to)
//       const registryLinkBefore = await linkToken.balanceOf(registry.address)
//       const registryPremiumBefore = (await registry.getState()).state
//         .totalPremium
//       const ownerBefore = (await registry.getState()).state.ownerLinkBalance
//
//       // Withdrawing for first time, last collected = 0
//       assert.equal(keeperBefore.lastCollected.toString(), '0')
//
//       //// Do the thing
//       await registry
//         .connect(payee1)
//         .withdrawPayment(await keeper1.getAddress(), to)
//
//       const keeperAfter = await registry.getTransmitterInfo(
//         await keeper1.getAddress(),
//       )
//       const registrationAfter = (await registry.getUpkeep(upkeepId)).balance
//       const toLinkAfter = await linkToken.balanceOf(to)
//       const registryLinkAfter = await linkToken.balanceOf(registry.address)
//       const registryPremiumAfter = (await registry.getState()).state
//         .totalPremium
//       const ownerAfter = (await registry.getState()).state.ownerLinkBalance
//
//       // registry total premium should not change
//       assert.isTrue(registryPremiumBefore.eq(registryPremiumAfter))
//
//       // Last collected should be updated to premium-change
//       assert.isTrue(
//         keeperAfter.lastCollected.eq(
//           registryPremiumBefore.sub(
//             registryPremiumBefore.mod(keeperAddresses.length),
//           ),
//         ),
//       )
//
//       // owner balance should remain unchanged
//       assert.isTrue(ownerAfter.eq(ownerBefore))
//
//       assert.isTrue(keeperAfter.balance.eq(BigNumber.from(0)))
//       assert.isTrue(registrationBefore.eq(registrationAfter))
//       assert.isTrue(toLinkBefore.add(keeperBefore.balance).eq(toLinkAfter))
//       assert.isTrue(
//         registryLinkBefore.sub(keeperBefore.balance).eq(registryLinkAfter),
//       )
//     })
//
//     it('emits a log announcing the withdrawal', async () => {
//       const balance = (
//         await registry.getTransmitterInfo(await keeper1.getAddress())
//       ).balance
//       const tx = await registry
//         .connect(payee1)
//         .withdrawPayment(
//           await keeper1.getAddress(),
//           await nonkeeper.getAddress(),
//         )
//       await expect(tx)
//         .to.emit(registry, 'PaymentWithdrawn')
//         .withArgs(
//           await keeper1.getAddress(),
//           balance,
//           await nonkeeper.getAddress(),
//           await payee1.getAddress(),
//         )
//     })
//   })
//
//   describe('#checkCallback', () => {
//     it('returns false with appropriate failure reason when target callback reverts', async () => {
//       await streamsLookupUpkeep.setShouldRevertCallback(true)
//
//       const values: any[] = ['0x1234', '0xabcd']
//       const res = await registry
//         .connect(zeroAddress)
//         .callStatic.checkCallback(streamsLookupUpkeepId, values, '0x')
//
//       assert.isFalse(res.upkeepNeeded)
//       assert.equal(res.performData, '0x')
//       assert.equal(
//         res.upkeepFailureReason,
//         UpkeepFailureReason.CHECK_CALLBACK_REVERTED,
//       )
//       assert.isTrue(res.gasUsed.gt(BigNumber.from('0'))) // Some gas should be used
//     })
//
//     it('returns false with appropriate failure reason when target callback returns big performData', async () => {
//       let longBytes = '0x'
//       for (let i = 0; i <= maxPerformDataSize.toNumber(); i++) {
//         longBytes += '11'
//       }
//       const values: any[] = [longBytes, longBytes]
//       const res = await registry
//         .connect(zeroAddress)
//         .callStatic.checkCallback(streamsLookupUpkeepId, values, '0x')
//
//       assert.isFalse(res.upkeepNeeded)
//       assert.equal(res.performData, '0x')
//       assert.equal(
//         res.upkeepFailureReason,
//         UpkeepFailureReason.PERFORM_DATA_EXCEEDS_LIMIT,
//       )
//       assert.isTrue(res.gasUsed.gt(BigNumber.from('0'))) // Some gas should be used
//     })
//
//     it('returns false with appropriate failure reason when target callback returns false', async () => {
//       await streamsLookupUpkeep.setCallbackReturnBool(false)
//       const values: any[] = ['0x1234', '0xabcd']
//       const res = await registry
//         .connect(zeroAddress)
//         .callStatic.checkCallback(streamsLookupUpkeepId, values, '0x')
//
//       assert.isFalse(res.upkeepNeeded)
//       assert.equal(res.performData, '0x')
//       assert.equal(
//         res.upkeepFailureReason,
//         UpkeepFailureReason.UPKEEP_NOT_NEEDED,
//       )
//       assert.isTrue(res.gasUsed.gt(BigNumber.from('0'))) // Some gas should be used
//     })
//
//     it('succeeds with upkeep needed', async () => {
//       const values: any[] = ['0x1234', '0xabcd']
//
//       const res = await registry
//         .connect(zeroAddress)
//         .callStatic.checkCallback(streamsLookupUpkeepId, values, '0x')
//       const expectedPerformData = ethers.utils.defaultAbiCoder.encode(
//         ['bytes[]', 'bytes'],
//         [values, '0x'],
//       )
//
//       assert.isTrue(res.upkeepNeeded)
//       assert.equal(res.performData, expectedPerformData)
//       assert.equal(res.upkeepFailureReason, UpkeepFailureReason.NONE)
//       assert.isTrue(res.gasUsed.gt(BigNumber.from('0'))) // Some gas should be used
//     })
//   })
//
//   describe('#setUpkeepPrivilegeConfig() / #getUpkeepPrivilegeConfig()', () => {
//     it('reverts when non manager tries to set privilege config', async () => {
//       await evmRevert(
//         registry.connect(payee3).setUpkeepPrivilegeConfig(upkeepId, '0x1234'),
//         'OnlyCallableByUpkeepPrivilegeManager()',
//       )
//     })
//
//     it('returns empty bytes for upkeep privilege config before setting', async () => {
//       const cfg = await registry.getUpkeepPrivilegeConfig(upkeepId)
//       assert.equal(cfg, '0x')
//     })
//
//     it('allows upkeep manager to set privilege config', async () => {
//       const tx = await registry
//         .connect(personas.Norbert)
//         .setUpkeepPrivilegeConfig(upkeepId, '0x1234')
//       await expect(tx)
//         .to.emit(registry, 'UpkeepPrivilegeConfigSet')
//         .withArgs(upkeepId, '0x1234')
//
//       const cfg = await registry.getUpkeepPrivilegeConfig(upkeepId)
//       assert.equal(cfg, '0x1234')
//     })
//   })
//
//   describe('#setAdminPrivilegeConfig() / #getAdminPrivilegeConfig()', () => {
//     const admin = randomAddress()
//
//     it('reverts when non manager tries to set privilege config', async () => {
//       await evmRevert(
//         registry.connect(payee3).setAdminPrivilegeConfig(admin, '0x1234'),
//         'OnlyCallableByUpkeepPrivilegeManager()',
//       )
//     })
//
//     it('returns empty bytes for upkeep privilege config before setting', async () => {
//       const cfg = await registry.getAdminPrivilegeConfig(admin)
//       assert.equal(cfg, '0x')
//     })
//
//     it('allows upkeep manager to set privilege config', async () => {
//       const tx = await registry
//         .connect(personas.Norbert)
//         .setAdminPrivilegeConfig(admin, '0x1234')
//       await expect(tx)
//         .to.emit(registry, 'AdminPrivilegeConfigSet')
//         .withArgs(admin, '0x1234')
//
//       const cfg = await registry.getAdminPrivilegeConfig(admin)
//       assert.equal(cfg, '0x1234')
//     })
//   })
//
//   describe('transmitterPremiumSplit [ @skip-coverage ]', () => {
//     beforeEach(async () => {
//       await linkToken.connect(owner).approve(registry.address, toWei('100'))
//       await registry.connect(owner).addFunds(upkeepId, toWei('100'))
//     })
//
//     it('splits premium evenly across transmitters', async () => {
//       // Do a transmit from keeper1
//       await getTransmitTx(registry, keeper1, [upkeepId])
//
//       const registryPremium = (await registry.getState()).state.totalPremium
//       assert.isTrue(registryPremium.gt(BigNumber.from(0)))
//
//       const premiumPerTransmitter = registryPremium.div(
//         BigNumber.from(keeperAddresses.length),
//       )
//       const k1Balance = (
//         await registry.getTransmitterInfo(await keeper1.getAddress())
//       ).balance
//       // transmitter should be reimbursed for gas and get the premium
//       assert.isTrue(k1Balance.gt(premiumPerTransmitter))
//       const k1GasReimbursement = k1Balance.sub(premiumPerTransmitter)
//
//       const k2Balance = (
//         await registry.getTransmitterInfo(await keeper2.getAddress())
//       ).balance
//       // non transmitter should get its share of premium
//       assert.isTrue(k2Balance.eq(premiumPerTransmitter))
//
//       // Now do a transmit from keeper 2
//       await getTransmitTx(registry, keeper2, [upkeepId])
//       const registryPremiumNew = (await registry.getState()).state.totalPremium
//       assert.isTrue(registryPremiumNew.gt(registryPremium))
//       const premiumPerTransmitterNew = registryPremiumNew.div(
//         BigNumber.from(keeperAddresses.length),
//       )
//       const additionalPremium = premiumPerTransmitterNew.sub(
//         premiumPerTransmitter,
//       )
//
//       const k1BalanceNew = (
//         await registry.getTransmitterInfo(await keeper1.getAddress())
//       ).balance
//       // k1 should get the new premium
//       assert.isTrue(
//         k1BalanceNew.eq(k1GasReimbursement.add(premiumPerTransmitterNew)),
//       )
//
//       const k2BalanceNew = (
//         await registry.getTransmitterInfo(await keeper2.getAddress())
//       ).balance
//       // k2 should get gas reimbursement in addition to new premium
//       assert.isTrue(k2BalanceNew.gt(k2Balance.add(additionalPremium)))
//     })
//
//     it('updates last collected upon payment withdrawn', async () => {
//       // Do a transmit from keeper1
//       await getTransmitTx(registry, keeper1, [upkeepId])
//
//       const registryPremium = (await registry.getState()).state.totalPremium
//       const k1 = await registry.getTransmitterInfo(await keeper1.getAddress())
//       const k2 = await registry.getTransmitterInfo(await keeper2.getAddress())
//
//       // Withdrawing for first time, last collected = 0
//       assert.isTrue(k1.lastCollected.eq(BigNumber.from(0)))
//       assert.isTrue(k2.lastCollected.eq(BigNumber.from(0)))
//
//       //// Do the thing
//       await registry
//         .connect(payee1)
//         .withdrawPayment(
//           await keeper1.getAddress(),
//           await nonkeeper.getAddress(),
//         )
//
//       const k1New = await registry.getTransmitterInfo(
//         await keeper1.getAddress(),
//       )
//       const k2New = await registry.getTransmitterInfo(
//         await keeper2.getAddress(),
//       )
//
//       // transmitter info lastCollected should be updated for k1, not for k2
//       assert.isTrue(
//         k1New.lastCollected.eq(
//           registryPremium.sub(registryPremium.mod(keeperAddresses.length)),
//         ),
//       )
//       assert.isTrue(k2New.lastCollected.eq(BigNumber.from(0)))
//     })
//
//     itMaybe(
//       'maintains consistent balance information across all parties',
//       async () => {
//         // throughout transmits, withdrawals, setConfigs total claim on balances should remain less than expected balance
//         // some spare change can get lost but it should be less than maxAllowedSpareChange
//
//         let maxAllowedSpareChange = BigNumber.from('0')
//         await verifyConsistentAccounting(maxAllowedSpareChange)
//
//         await getTransmitTx(registry, keeper1, [upkeepId])
//         maxAllowedSpareChange = maxAllowedSpareChange.add(BigNumber.from('31'))
//         await verifyConsistentAccounting(maxAllowedSpareChange)
//
//         await registry
//           .connect(payee1)
//           .withdrawPayment(
//             await keeper1.getAddress(),
//             await nonkeeper.getAddress(),
//           )
//         await verifyConsistentAccounting(maxAllowedSpareChange)
//
//         await registry
//           .connect(payee2)
//           .withdrawPayment(
//             await keeper2.getAddress(),
//             await nonkeeper.getAddress(),
//           )
//         await verifyConsistentAccounting(maxAllowedSpareChange)
//
//         await getTransmitTx(registry, keeper1, [upkeepId])
//         maxAllowedSpareChange = maxAllowedSpareChange.add(BigNumber.from('31'))
//         await verifyConsistentAccounting(maxAllowedSpareChange)
//
//         await registry.connect(owner).setConfigTypeSafe(
//           signerAddresses.slice(2, 15), // only use 2-14th index keepers
//           keeperAddresses.slice(2, 15),
//           f,
//           config,
//           offchainVersion,
//           offchainBytes,
//         )
//         await verifyConsistentAccounting(maxAllowedSpareChange)
//
//         await getTransmitTx(registry, keeper3, [upkeepId], {
//           startingSignerIndex: 2,
//         })
//         maxAllowedSpareChange = maxAllowedSpareChange.add(BigNumber.from('13'))
//         await verifyConsistentAccounting(maxAllowedSpareChange)
//
//         await registry
//           .connect(payee1)
//           .withdrawPayment(
//             await keeper1.getAddress(),
//             await nonkeeper.getAddress(),
//           )
//         await verifyConsistentAccounting(maxAllowedSpareChange)
//
//         await registry
//           .connect(payee3)
//           .withdrawPayment(
//             await keeper3.getAddress(),
//             await nonkeeper.getAddress(),
//           )
//         await verifyConsistentAccounting(maxAllowedSpareChange)
//
//         await registry.connect(owner).setConfigTypeSafe(
//           signerAddresses.slice(0, 4), // only use 0-3rd index keepers
//           keeperAddresses.slice(0, 4),
//           f,
//           config,
//           offchainVersion,
//           offchainBytes,
//         )
//         await verifyConsistentAccounting(maxAllowedSpareChange)
//         await getTransmitTx(registry, keeper1, [upkeepId])
//         maxAllowedSpareChange = maxAllowedSpareChange.add(BigNumber.from('4'))
//         await getTransmitTx(registry, keeper3, [upkeepId])
//         maxAllowedSpareChange = maxAllowedSpareChange.add(BigNumber.from('4'))
//
//         await verifyConsistentAccounting(maxAllowedSpareChange)
//         await registry
//           .connect(payee5)
//           .withdrawPayment(
//             await keeper5.getAddress(),
//             await nonkeeper.getAddress(),
//           )
//         await verifyConsistentAccounting(maxAllowedSpareChange)
//
//         await registry
//           .connect(payee1)
//           .withdrawPayment(
//             await keeper1.getAddress(),
//             await nonkeeper.getAddress(),
//           )
//         await verifyConsistentAccounting(maxAllowedSpareChange)
//       },
//     )
//   })
// })
