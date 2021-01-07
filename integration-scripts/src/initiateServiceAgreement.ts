import { Coordinator__factory } from '@chainlink/contracts/ethers/v0.5/factories/Coordinator__factory'
import * as t from '@chainlink/test-helpers'
import { ethers } from 'ethers'
import {
  createProvider,
  DEVNET_ADDRESS,
  getArgs,
  registerPromiseHandler,
} from './common'

type OracleSignatures = t.coordinator.OracleSignatures
type ServiceAgreement = t.coordinator.ServiceAgreement

/**
 * This json definition may be missing types, it was generated from a fixture.
 */
interface ServiceAgreementJson {
  initiators: Initiator[]
  tasks: Task[]
  payment: string
  expiration: number
  oracles: string[]
  endAt: string
  aggregator: string
  aggInitiateJobSelector: string
  aggFulfillSelector: string
}

interface Task {
  type: string
  params?: TaskParams
}

interface TaskParams {
  get?: string
  path?: string[]
  address?: string
  functionSelector?: string
}

interface Initiator {
  type: string
}

async function main() {
  registerPromiseHandler()

  const args = getArgs([
    'COORDINATOR_ADDRESS',
    'MEAN_AGGREGATOR_ADDRESS',
    'ORACLE_SIGNATURE',
    'NORMALIZED_REQUEST',
    'SERVICE_AGREEMENT',
    'CHAINLINK_NODE_ADDRESS',
  ])

  await initiateServiceAgreement({
    coordinatorAddress: args.COORDINATOR_ADDRESS,
    meanAggregatorAddress: args.MEAN_AGGREGATOR_ADDRESS,
    normalizedRequest: args.NORMALIZED_REQUEST,
    oracleSignature: args.ORACLE_SIGNATURE,
    saJson: JSON.parse(args.SERVICE_AGREEMENT),
    expectedAddress: args.CHAINLINK_NODE_ADDRESS,
  })
}
main()

interface Args {
  coordinatorAddress: string
  meanAggregatorAddress: string
  oracleSignature: string
  normalizedRequest: string
  saJson: ServiceAgreementJson
  expectedAddress: string
}

async function initiateServiceAgreement({
  coordinatorAddress,
  normalizedRequest,
  oracleSignature,
  saJson,
  expectedAddress,
}: Args) {
  const provider = createProvider()
  const signer = provider.getSigner(DEVNET_ADDRESS)
  const coordinatorFactory = new Coordinator__factory(signer)
  const coordinator = coordinatorFactory.attach(coordinatorAddress)

  console.log('Creating service agreement to initiate with...')
  const sa: ServiceAgreement = {
    aggFulfillSelector: saJson.aggFulfillSelector,
    aggInitiateJobSelector: saJson.aggInitiateJobSelector,
    aggregator: saJson.aggregator,
    expiration: saJson.expiration,
    oracles: saJson.oracles,
    payment: saJson.payment,
    endAt: Math.round(new Date(saJson.endAt).getTime() / 1000), // end date in seconds
    requestDigest: ethers.utils.keccak256(
      ethers.utils.toUtf8Bytes(normalizedRequest),
    ),
  }

  const sig = ethers.utils.splitSignature(oracleSignature)
  if (!sig.v) {
    throw Error(`Could not extract v from signature`)
  }

  const oracleSignatures: OracleSignatures = {
    vs: [sig.v],
    rs: [sig.r],
    ss: [sig.s],
  }
  const encodedSignatures = t.coordinator.encodeOracleSignatures(
    oracleSignatures,
  )

  const said = t.coordinator.generateSAID(sa)
  const encodedSA = t.coordinator.encodeServiceAgreement(sa)
  const ssaid = await coordinator.getId(encodedSA)
  if (said != ssaid) {
    throw Error(`sAId mismatch. javascript: ${said} solidity: ${ssaid}`)
  }

  console.log('Initiating service agreement...')
  const recoveredAddresss = ethers.utils.recoverAddress(said, sig)

  console.log({
    recoveredAddresss,
    expectedAddress,
    oracleSignature,
    ...sig,
  })

  const tx = await coordinator.initiateServiceAgreement(
    encodedSA,
    encodedSignatures,
  )
  console.log(tx)
  const iSAreceipt = await tx.wait()
  console.log('initiateServiceAgreement receipt', iSAreceipt)
}
