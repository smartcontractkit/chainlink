import { debug } from 'chainlink'
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

  console.log(
    '########################################################################',
    coordinatorFactory.interface.abi,
  )

  console.log('coordinator_address', coordinatorAddress)
  const amount = ethers.utils.parseEther('1000')

  type CoordinatorParams = Parameters<Coordinator['initiateServiceAgreement']>

  const agreement: CoordinatorParams[0] = {
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
    throw Error(
      `THIS SHOULD NOT HAPPEN: sig.v is undefined when it should always exist`,
    )
  }
  const oracleSignatures: CoordinatorParams[1] = {
    rs: [sig.r],
    ss: [sig.s],
    vs: [sig.v],
  }

  console.log('coordinator_address', coordinatorAddress)

  console.log(
    '!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!',
    await coordinator.dummyMethodXXX(),
  )

  // const tx = await coordinator.initiateServiceAgreement(
  //   agreement,
  //   oracleSignatures,
  // )
  // await tx.wait()
  // console.log('initiateServiceAgreement', tx)

  // console.log(
  //   'oracleRequest',
  //   await Coordinator.methods.oracleRequest(
  //     agreement.sAID,
  //     '0x0101010101010101010101010101010101010101', // Receiving contract address
  //     '0x12345678', // receiving method selector
  //     1, // nonce
  //     '', // data for initialization of request
  // ),
  // )
}
