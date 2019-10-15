import { debug, helpers } from 'chainlink'
import {
  getArgs,
  registerPromiseHandler,
  createProvider,
  DEVNET_ADDRESS,
} from './common'
import { CoordinatorFactory } from './generated/CoordinatorFactory'
import agreementJson from './fixtures/agreement.json'
import { ethers } from 'ethers'
import { Coordinator } from './generated/Coordinator'

const _d = debug.makeDebug('initiateServiceAgreement')
async function main() {
  const d = _d.extend('main')
  registerPromiseHandler()
  const args = getArgs([
    'COORDINATOR_ADDRESS',
    'MEAN_AGGREGATOR_ADDRESS',
    'ORACLE_SIGNATURE',
    'NORMALIZED_REQUEST',
  ])

  d(args)

  await initiateServiceAgreement({
    coordinatorAddress: args.COORDINATOR_ADDRESS,
    meanAggregatorAddress: args.MEAN_AGGREGATOR_ADDRESS,
    normalizedRequest: args.NORMALIZED_REQUEST,
    oracleSignature: args.ORACLE_SIGNATURE,
  })
}
main()

interface Args {
  coordinatorAddress: string
  meanAggregatorAddress: string
  oracleSignature: string
  normalizedRequest: string
}

async function initiateServiceAgreement({
  coordinatorAddress,
  meanAggregatorAddress,
  normalizedRequest,
  oracleSignature,
}: Args) {
  const provider = createProvider()
  const signer = provider.getSigner(DEVNET_ADDRESS)
  const coordinatorFactory = new CoordinatorFactory(signer)
  const coordinator = coordinatorFactory.attach(coordinatorAddress)

  type CoordinatorParams = Parameters<Coordinator['initiateServiceAgreement']>
  type ServiceAgreement = CoordinatorParams[0]
  type OracleSignatures = CoordinatorParams[1]

  const fieldTypes = helpers.serviceAgreementFieldTypes().map(f => f.type).join(',')
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

  console.log('agreement', agreement)

  const sig = ethers.utils.splitSignature(oracleSignature)
  if (!sig.v) {
    throw Error(`Could not extract v from signature`)
  }
  const oracleSignatures: OracleSignatures = {
    rs: [sig.r],
    ss: [sig.s],
    vs: [sig.v],
  }

  console.log('Attempting to initiate service agreement')
  console.log('oracle signatures', oracleSignatures)

  const said = helpers.calculateSAID2(agreement)
  console.log('our said', said, 'solidity\'s said', await coordinator.getId(agreement))

  console.log('apparent address', ethers.utils.recoverAddress(said, oracleSignature))
  console.log('actual address', agreement.oracles)

  throw Error('foo')

  try {
    provider.on(
      { topics: [ ethers.utils.id('SignatureCheck(address,address)')] },
      r => console.log('SignatureCheck event', r),
    )
  } catch(e) {
    console.log('provider.on error', e)
    throw e
  }
  console.log('provider.on worked')
  const tx = await coordinator.initiateServiceAgreement(
    agreement,
    oracleSignatures,
  )
  const iSAreceipt = await tx.wait()
  console.log('initiateServiceAgreement receipt', iSAreceipt)

  console.log('geetting to here, said is', said)

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
    receipt 
  )
}
