import { assert } from 'chai'
import { ContractTransaction, ethers, utils } from 'ethers'
import { BigNumberish } from 'ethers/utils'
import { bigNum, sixMonthsFromNow, stripHexPrefix, toHex } from '../helpers'
import * as matchers from '../matchers'

export interface ServiceAgreement {
  payment: ethers.utils.BigNumberish // uint256
  expiration: ethers.utils.BigNumberish // uint256
  endAt: ethers.utils.BigNumberish // uint256
  oracles: string[] | ethers.Wallet[] // 0x hex representation of oracle addresses (uint160's), or wallet instances to map to addresses
  requestDigest: string // 0x hex representation of bytes32
  aggregator: string // 0x hex representation of aggregator address
  aggInitiateJobSelector: string // 0x hex representation of aggregator.initiateAggregatorForJob function selector (uint32)
  aggFulfillSelector: string // function selector for aggregator.fulfill
}

export interface Signature extends utils.Signature {
  v: number
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

/**
 * Abi encode the oracleRequest() function
 *
 * @param sAID The service agreement ID
 * @param callbackAddr The callback contract address for the response
 * @param callbackFunctionId The callback function id for the response
 * @param nonce The nonce sent by the requester
 * @param data The CBOR payload of the request
 */
export function encodeOracleRequest(
  specId: string,
  to: string,
  fHash: string,
  nonce: BigNumberish,
  dataBytes: string,
): string {
  // 'oracleRequest(address,uint256,bytes32,address,bytes4,uint256,uint256,bytes)'
  const oracleRequestSighash = '0x40429946'
  const oracleRequestInputs = [
    { name: '_sender', type: 'address' },
    { name: '_amount', type: 'uint256' },
    { name: '_sAId', type: 'bytes32' },
    { name: '_callbackAddress', type: 'address' },
    { name: '_callbackFunctionId', type: 'bytes4' },
    { name: '_nonce', type: 'uint256' },
    { name: '_dataVersion', type: 'uint256' },
    { name: '_data', type: 'bytes' },
  ]

  const encodedParams = ethers.utils.defaultAbiCoder.encode(
    oracleRequestInputs.map(i => i.type),
    [ethers.constants.AddressZero, 0, specId, to, fHash, nonce, 1, dataBytes],
  )

  return `${oracleRequestSighash}${stripHexPrefix(encodedParams)}`
}

export async function newServiceAgreement(
  params: Partial<ServiceAgreement>,
): Promise<ServiceAgreement> {
  const agreement: ServiceAgreement = {
    payment: bigNum('1000000000000000000'),
    expiration: bigNum(300),
    endAt: sixMonthsFromNow(),
    oracles: [],
    requestDigest:
      '0xbadc0de5badc0de5badc0de5badc0de5badc0de5badc0de5badc0de5badc0de5',
    aggregator: '0x3141592653589793238462643383279502884197',
    aggInitiateJobSelector: '0xd43a12f6', // initiateJob()
    aggFulfillSelector: '0x9760168f', // fulfill()
    ...params,
  }
  return agreement
}

interface InitiateServiceAgreementer {
  initiateServiceAgreement(
    _serviceAgreementData: utils.Arrayish,
    _oracleSignaturesData: utils.Arrayish,
  ): Promise<ContractTransaction>
}
export async function initiateServiceAgreement(
  coordinator: InitiateServiceAgreementer,
  serviceAgreementParams: Partial<ServiceAgreement>,
) {
  const serviceAgreement = await newServiceAgreement(serviceAgreementParams)
  const signatures = await generateOracleSignatures(serviceAgreement)
  return coordinator.initiateServiceAgreement(
    encodeServiceAgreement(serviceAgreement),
    encodeOracleSignatures(signatures),
  )
}

const serviceAgreementValues = (sa: ServiceAgreement) => {
  return [
    sa.payment,
    sa.expiration,
    sa.endAt,
    oracleAddresses(sa.oracles),
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

/**
 * Generates the oracle signatures on a ServiceAgreement
 * @param serviceAgreement The service agreement to sign
 * @param signers The list oracles that will sign the service agreement
 */
export async function generateOracleSignatures(
  serviceAgreement: ServiceAgreement,
): Promise<OracleSignatures> {
  const sAID = generateSAID(serviceAgreement)
  const signatures = []

  for (let i = 0; i < serviceAgreement.oracles.length; i++) {
    const oracle = serviceAgreement.oracles[i]
    if (!(oracle instanceof ethers.Wallet)) {
      throw Error('cannot generate signatures without oracle wallets')
    }
    const oracleSignature = await personalSign(sAID, oracle)
    const requestDigestAddr = recoverAddressFromSignature(sAID, oracleSignature)
    assert.equal(oracle.address, requestDigestAddr)
    signatures.push(oracleSignature)
  }

  return combineOracleSignatures(signatures)
}

function oracleAddresses(oracles: string[] | ethers.Wallet[]): string[] {
  const oracleAddresses: string[] = []
  oracles.forEach((oracle: string | ethers.Wallet) => {
    if (oracle instanceof ethers.Wallet) {
      oracleAddresses.push(oracle.address)
    } else {
      oracleAddresses.push(oracle)
    }
  })
  return oracleAddresses
}

interface ServiceAgreementer {
  serviceAgreements(
    id: string,
  ): Promise<{
    payment: utils.BigNumber
    expiration: utils.BigNumber
    endAt: utils.BigNumber
    requestDigest: string
    aggregator: string
    aggInitiateJobSelector: string
    aggFulfillSelector: string
  }>
}
/**
 * Check that the given service agreement was stored at the correct location
 * @param coordinator The coordinator contract
 * @param serviceAgreement The service agreement
 */
export async function assertServiceAgreementPresent(
  coordinator: ServiceAgreementer,
  serviceAgreement: ServiceAgreement,
) {
  const sAID = generateSAID(serviceAgreement)
  const sa = await coordinator.serviceAgreements(sAID)

  matchers.bigNum(
    sa.payment,
    bigNum(serviceAgreement.payment),
    'expected payment',
  )
  matchers.bigNum(
    sa.expiration,
    bigNum(serviceAgreement.expiration),
    'expected expiration',
  )
  matchers.bigNum(
    sa.endAt,
    bigNum(serviceAgreement.endAt),
    'expected endAt date',
  )
  assert.equal(
    sa.requestDigest,
    serviceAgreement.requestDigest,
    'expected requestDigest',
  )
}

/**
 * Check that all values for the struct at this SAID have default values. I.e.
 * nothing was changed due to invalid request
 * @param coordinator The coordinator contract
 * @param serviceAgreementID The service agreement ID
 */
export async function assertServiceAgreementEmpty(
  coordinator: ServiceAgreementer,
  serviceAgreementID: string,
) {
  const sa = await coordinator.serviceAgreements(
    toHex(serviceAgreementID).slice(0, 66), // serviceAggrementId is contained within the highest 32 bytes
  )
  matchers.bigNum(sa.payment, bigNum(0), 'service agreement is not absent')
  matchers.bigNum(sa.expiration, bigNum(0), 'service agreement is not absent')
  matchers.bigNum(sa.endAt, bigNum(0), 'service agreement is not absent')
  assert.equal(
    sa.requestDigest,
    '0x0000000000000000000000000000000000000000000000000000000000000000',
  )
}

/**
 * Combine v, r, and s params of multiple signatures into format expected by contracts
 * @param signatures The list of signatures to combine
 */
export function combineOracleSignatures(
  signatures: Signature[],
): OracleSignatures {
  return signatures.reduce<OracleSignatures>(
    (prev, { v, r, s }) => {
      prev.vs.push(v)
      prev.rs.push(r)
      prev.ss.push(s)

      return prev
    },
    { vs: [], rs: [], ss: [] },
  )
}

/**
 * Signs a message according to ethereum specs by first appending
 * "\x19Ethereum Signed Message:\n' + <message.length>" to the message
 * @param message The message to sign - either a Buffer or a hex string
 * @param wallet The wallet of the signer
 */
export async function personalSign(
  message: Buffer | string,
  wallet: ethers.Wallet,
): Promise<Signature> {
  function assertIsSignature(sig: utils.Signature): asserts sig is Signature {
    if (!sig.v) throw Error(`Could not extract v from signature`)
  }
  if (message instanceof String && !utils.isHexString(message)) {
    throw Error(`The message ${message} is not a valid hex string`)
  }
  const flatSig = await wallet.signMessage(utils.arrayify(message))
  const splitSignature = utils.splitSignature(flatSig)
  assertIsSignature(splitSignature)
  return splitSignature
}

/**
 * Recovers the address of the signer of a message
 * @param message The message that was signed
 * @param signature The signature on the message
 */
export function recoverAddressFromSignature(
  message: string | Buffer,
  signature: Signature,
): string {
  const messageBuff = utils.arrayify(message)
  return utils.verifyMessage(messageBuff, signature)
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
