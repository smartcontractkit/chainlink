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
    'ORACLE_SIGNATURE',
    'NORMALIZED_REQUEST',
  ])

  d(args)

  await initiateServiceAgreement({
    coordinatorAddress: args.COORDINATOR_ADDRESS,
    normalizedRequest: args.NORMALIZED_REQUEST,
    oracleSignature: args.ORACLE_SIGNATURE,
  })
}
main()

interface Args {
  coordinatorAddress: string
  oracleSignature: string
  normalizedRequest: string
}

async function initiateServiceAgreement({
  coordinatorAddress,
  normalizedRequest,
  oracleSignature,
}: Args) {
  const d = _d.extend('initiateServiceAgreement')
  const provider = createProvider()
  const signer = provider.getSigner(DEVNET_ADDRESS)
  const coordinatorFactory = new CoordinatorFactory(signer)
  const coordinator = coordinatorFactory.attach(coordinatorAddress)

  type CoordinatorParams = Parameters<Coordinator['initiateServiceAgreement']>
  type ServiceAgreement = CoordinatorParams[0]
  type OracleSignatures = CoordinatorParams[1]

  const agreement: ServiceAgreement = {
    aggFulfillSelector: agreementJson.aggFulfillSelector,
    aggInitiateJobSelector: agreementJson.aggInitiateJobSelector,
    aggregator: agreementJson.aggregator,
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

  console.log('Attempting to initiate service agreement')
  const tx = await coordinator.initiateServiceAgreement(
    agreement,
    oracleSignatures,
  )
  const iSAreceipt = await tx.wait()
  console.log('initiateServiceAgreement receipt', iSAreceipt)
  
  const said = helpers.calculateSAID2(agreement)

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
