import { ethers } from 'ethers'

export interface ServiceAgreement {
  payment: ethers.utils.BigNumberish // uint256
  expiration: ethers.utils.BigNumberish // uint256
  endAt: ethers.utils.BigNumberish // uint256
  oracles: string[] // 0x hex representation of oracle addresses (uint160's)
  requestDigest: string // 0x hex representation of bytes32
  aggregator: string // 0x hex representation of aggregator address
  aggInitiateJobSelector: string // 0x hex representation of aggregator.initiateAggregatorForJob function selector (uint32)
  aggFulfillSelector: string // function selector for aggregator.fulfill
}

export interface OracleSignatures {
  vs: ethers.utils.BigNumberish[] // uint8[]
  rs: string[] // bytes32[]
  ss: string[] // bytes32[]
}

const SERVICE_AGREEMENT_TYPES = [
  'uint256',
  'uint256',
  'uint256',
  'address[]',
  'bytes32',
  'address',
  'bytes4',
  'bytes4',
]

const ORACLE_SIGNATURES_TYPES = ['uint8[]', 'bytes32[]', 'bytes32[]']

const serviceAgreementValues = (sa: ServiceAgreement) => {
  return [
    sa.payment,
    sa.expiration,
    sa.endAt,
    sa.oracles,
    sa.requestDigest,
    sa.aggregator,
    sa.aggInitiateJobSelector,
    sa.aggFulfillSelector,
  ]
}

export function encodeServiceAgreement(sa: ServiceAgreement) {
  return ethers.utils.defaultAbiCoder.encode(
    SERVICE_AGREEMENT_TYPES,
    serviceAgreementValues(sa),
  )
}

export function encodeOracleSignatures(os: OracleSignatures) {
  const osValues = [os.vs, os.rs, os.ss]
  return ethers.utils.defaultAbiCoder.encode(ORACLE_SIGNATURES_TYPES, osValues)
}

export async function computeOracleSignature(
  agreement: ServiceAgreement,
  oracle: ethers.Wallet,
): Promise<OracleSignatures> {
  const said = generateSAID(agreement)
  const oracleSignatures: OracleSignatures[] = []

  for (let i = 0; i < agreement.oracles.length; i++) {
    const oracleSignature = await oracle.signMessage(
      ethers.utils.arrayify(said),
    )

    const sig = ethers.utils.splitSignature(oracleSignature)
    if (!sig.v) {
      throw Error(`Could not extract v from signature`)
    }
    const convertedOracleSignature: OracleSignatures = {
      vs: [sig.v],
      rs: [sig.r],
      ss: [sig.s],
    }
    oracleSignatures.push(convertedOracleSignature)
  }

  // TODO: this should be an array!
  return oracleSignatures[0]
}

type Hash = ReturnType<typeof ethers.utils.keccak256>

/**
 * Digest of the ServiceAgreement.
 */
export function generateSAID(sa: ServiceAgreement): Hash {
  return ethers.utils.solidityKeccak256(
    SERVICE_AGREEMENT_TYPES,
    serviceAgreementValues(sa),
  )
}
