import { helpers } from 'chainlink'
import {
  getArgs,
  registerPromiseHandler,
  
} from './common'
import { CoordinatorFactory } from './generated/CoordinatorFactory'
import agreementJson from './fixtures/agreement.json'
import { ethers } from 'ethers'
import { Coordinator } from './generated/Coordinator'
import { MeanAggregatorFactory } from './generated/MeanAggregatorFactory'
import { createTraceProvider } from './common'
import { deployContracts } from './deployV0.5Contracts'

async function main() {
  registerPromiseHandler()
  const { defaultFromAddress, provider } = await createTraceProvider()
  const { coordinator, meanAggregator } = await deployContracts(
    provider,
    defaultFromAddress,
  )

  process.env.COORDINATOR_ADDRESS = coordinator.address
  process.env.MEAN_AGGREGATOR_ADDRESS = meanAggregator.address
  process.env.ORACLE_SIGNATURE =
    '0xc846280320ffef933ce090706c61945865e3407cbf35b6a3edd63cf11e2190206f531499c7d3b748a3538ed41bf0df76ad421704d7ab89131ae3b11654ce62b701'
  process.env.NORMALIZED_REQUEST =
    '{"aggFulfillSelector":"0xbadc0de5","aggInitiateJobSelector":"0xd0771e55","aggregator":"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF","endAt":"2019-10-19T22:17:19Z","expiration":3.000000e+02,"initiators":[{"type":"execagreement"}],"oracles":["0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f"],"payment":"1000000000000000000","tasks":[{"params":{"get":"https://bitstamp.net/api/ticker/"},"type":"HttpGet"},{"params":{"path":["last"]},"type":"JsonParse"},{"type":"EthBytes32"},{"params":{"address":"0x356a04bce728ba4c62a30294a55e6a8600a320b3","functionSelector":"0x609ff1bd"},"type":"EthTx"}]}'

  const args = getArgs([
    'COORDINATOR_ADDRESS',
    'MEAN_AGGREGATOR_ADDRESS',
    'ORACLE_SIGNATURE',
    'NORMALIZED_REQUEST',
  ])

  await initiateServiceAgreement({
    coordinatorAddress: args.COORDINATOR_ADDRESS,
    meanAggregatorAddress: args.MEAN_AGGREGATOR_ADDRESS,
    normalizedRequest: args.NORMALIZED_REQUEST,
    oracleSignature: args.ORACLE_SIGNATURE,
    provider,
    DEVNET_ADDRESS: defaultFromAddress
  })
}
main()

interface Args {
  coordinatorAddress: string
  meanAggregatorAddress: string
  oracleSignature: string
  normalizedRequest: string
  provider: ethers.providers.JsonRpcProvider
  DEVNET_ADDRESS: string
}

async function initiateServiceAgreement({
  coordinatorAddress,
  meanAggregatorAddress,
  normalizedRequest,
  oracleSignature,
  provider,
  DEVNET_ADDRESS
}: Args) {
  const signer = provider.getSigner(DEVNET_ADDRESS)
  const coordinatorFactory = new CoordinatorFactory(signer)
  const coordinator = coordinatorFactory.attach(coordinatorAddress)

  // Monkey-patches coordinator-contract interface to claim that all calls are
  // static / constant, so that all its methods can be called rather than
  // transacted upon, and their return values can be inspected.
  const coordinatorStaticfactory = new ethers.ContractFactory(
    coordinatorFactory.interface.abi.map(a => ({ ...a, constant: true })),
    coordinatorFactory.bytecode,
    signer,
  )
  const coordinatorStatic = coordinatorStaticfactory.attach(coordinatorAddress)

  type CoordinatorParams = Parameters<Coordinator['initiateServiceAgreement']>
  type ServiceAgreement = CoordinatorParams[0]
  type OracleSignatures = CoordinatorParams[1]

  const fieldTypes = helpers
    .serviceAgreementFieldTypes()
    .map(f => f.type)
    .join(',')
  const aggInitiateJobSelector = helpers.functionSelector(
    `initiateJob(bytes32,tuple(${fieldTypes}))`,
  )
  if (agreementJson.aggInitiateJobSelector !== aggInitiateJobSelector) {
    throw Error('Unexpected aggInitiateJobSelector')
  }
  // Must be equal because creation of the job on the CL node is done elsewhere
  const aggFulfillSelector = helpers.functionSelector(
    'fulfull(bytes32,bytes32,bytes32,bytes32)',
  )
  if (agreementJson.aggFulfillSelector !== aggFulfillSelector) {
    throw Error('Unexpected aggFulfillSelector')
  }

  const agreement: ServiceAgreement = {
    aggFulfillSelector: aggFulfillSelector,
    aggInitiateJobSelector: aggInitiateJobSelector,
    aggregator: meanAggregatorAddress,
    payment: agreementJson.payment,
    expiration: agreementJson.expiration,
    endAt: Math.round(new Date(agreementJson.endAt).getTime() / 1000), // end date in seconds
    oracles: agreementJson.oracles,
    requestDigest: ethers.utils.keccak256(
      ethers.utils.toUtf8Bytes(normalizedRequest),
    ),
  }

  const sig = ethers.utils.splitSignature(oracleSignature)
  if (!sig.v) {
    throw Error(`Could not extract v from signature`)
  }
  const oracleSignatures: OracleSignatures = {
    rs: [sig.r],
    ss: [sig.s],
    vs: [sig.v],
  }

  const said = helpers.calculateSAID2(agreement)
  const ssaid = await coordinator.getId(agreement)
  if (said != ssaid) {
    throw Error(`sAId mismatch. javascript: ${said} solidity: ${ssaid}`)
  }

  const meanAggregator = new MeanAggregatorFactory(signer).attach(
    meanAggregatorAddress,
  )

  meanAggregator.initiateJob(said, agreement)

  console.log('meanAggregator call worked...')

  // call initiateServiceAgreement instead of sending it a transaction, so we
  // can inspect the return value.
  const callVal = await coordinatorStatic.initiateServiceAgreement(
    agreement,
    oracleSignatures,
  )
  // if callVal !=  // XXX: Make this an assertion
  console.log(
    'call value of coordinatorStatic.initiateServiceAgreement',
    callVal,
  )

  const tx = await coordinator.initiateServiceAgreement(
    agreement,
    oracleSignatures,
  )
  const iSAreceipt = await tx.wait()
  console.log('initiateServiceAgreement receipt', iSAreceipt)

  const reqId = await coordinator.oracleRequest(
    '0x0101010101010101010101010101010101010101',
    10000000000000,
    said as any, // XXX:
    '0x0101010101010101010101010101010101010101', // Receiving contract address
    '0x12345678', // receiving method selector
    1, // nonce
    1, // data version
    '0x0', // data for initialization of request
  )
  const receipt = await reqId.wait()
  console.log(
    '************************************************************************ oracleRequest',
    receipt,
  )
}
