/**
 * @packageDocumentation
 *
 * This file provides convenience functions to interact with existing solidity contract abstraction libraries, such as
 * @truffle/contract and ethers.js specifically for our `Coordinator.sol` solidity smart contract.
 */
import { assert } from 'chai'
import { ethers, utils } from 'ethers'
import { BigNumberish } from 'ethers/utils'
import { bigNum, sixMonthsFromNow, stripHexPrefix } from '../helpers'
import * as matchers from '../matchers'

export interface ServiceAgreement {
  /**
   * Price in LINK to request a report based on this agreement
   *
   * @solformat uint256
   */
  payment: ethers.utils.BigNumberish
  /**
   * Expiration is the amount of time an oracle has to answer a request
   *
   * @solformat uint256
   */
  expiration: ethers.utils.BigNumberish
  /**
   * The service agreement is valid until this time
   *
   * @solformat uint256
   */
  endAt: ethers.utils.BigNumberish
  /**
   * An array of oracle addresses to use within the process of aggregation
   *
   * @solformat address[]
   */
  oracles: (string | ethers.Wallet)[]
  /**
   * This effectively functions as an ID tag for the off-chain job of the
   * service agreement. It is calculated as the keccak256 hash of the
   * normalized JSON request to create the ServiceAgreement, but that identity
   * is unused, and its value is essentially arbitrary.
   *
   * @solformat bytes32
   */
  requestDigest: string

  /**
   *  Specification of aggregator interface. See ../../../evm/contracts/tests/MeanAggregator.sol
   *  for example.
   */

  /**
   * Address of where the aggregator instance is held
   *
   * @solformat address
   */
  aggregator: string
  /**
   * Selectors for the interface methods must be specified, because their
   * arguments can vary from aggregator to aggregator.
   *
   * Function selector for aggregator initiateJob method
   *
   * @solformat bytes4
   */
  aggInitiateJobSelector: string
  /**
   * Function selector for aggregator fulfill method
   *
   * @solformat bytes4
   */
  aggFulfillSelector: string
}

/**
 * A collection of multiple oracle signatures stored via parallel arrays
 */
export interface OracleSignatures {
  /**
   * The recovery parameters normalized for Solidity, either 27 or 28
   *
   * @solformat uint8[]
   */
  vs: ethers.utils.BigNumberish[]
  /**
   * the r coordinate within (r, s) public point of a signature
   *
   * @solformat bytes32[]
   */
  rs: string[]
  /**
   * the s coordinate within (r, s) public point of a signature
   *
   * @solformat  bytes32[]
   */
  ss: string[]
}

/**
 * Create a service agreement with sane testing defaults
 *
 * @param overrides Values to override service agreement defaults
 */
export function serviceAgreement(
  overrides: Partial<ServiceAgreement>,
): ServiceAgreement {
  const agreement: ServiceAgreement = {
    payment: bigNum('1000000000000000000'),
    expiration: bigNum(300),
    endAt: sixMonthsFromNow(),
    oracles: [],
    requestDigest:
      '0xbadc0de5badc0de5badc0de5badc0de5badc0de5badc0de5badc0de5badc0de5',
    aggregator: '0x3141592653589793238462643383279502884197',
    aggInitiateJobSelector: '0xd43a12f6',
    aggFulfillSelector: '0x9760168f',
    ...overrides,
  }

  return agreement
}

/**
 * Check that all values for the struct at this SAID have default values.
 *
 * For example, when an invalid service agreement initialization request is made to a `Coordinator`, we want to make sure that
 * it did not initialize its service agreement struct to any value, hence checking for it being empty.
 *
 * @param coordinator The coordinator contract
 * @param serviceAgreementID The service agreement ID
 *
 * @throws when any of payment, expiration, endAt, requestDigest are non-empty
 */
export function assertServiceAgreementEmpty(
  sa: Omit<ServiceAgreement, 'oracles'>,
) {
  matchers.bigNum(sa.payment, bigNum(0), 'service agreement is not absent')
  matchers.bigNum(sa.expiration, bigNum(0), 'service agreement is not absent')
  matchers.bigNum(sa.endAt, bigNum(0), 'service agreement is not absent')
  assert.equal(
    sa.requestDigest,
    '0x0000000000000000000000000000000000000000000000000000000000000000',
  )
}

/**
 * Create parameters needed for the
 * ```solidity
 *   function initiateServiceAgreement(
 *    bytes memory _serviceAgreementData,
 *    bytes memory _oracleSignaturesData
 *  )
 * ```
 * method of the `Coordinator.sol` contract
 *
 * @param overrides Values to override the defaults for creating a service agreement
 */
export async function initiateSAParams(
  overrides: Partial<ServiceAgreement>,
): Promise<[string, string]> {
  const sa = serviceAgreement(overrides)
  const signatures = await generateOracleSignatures(sa)

  return [encodeServiceAgreement(sa), encodeOracleSignatures(signatures)]
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

/**
 * ABI encode a service agreement object
 *
 * @param sa The service agreement to encode
 */
export function encodeServiceAgreement(sa: ServiceAgreement) {
  return ethers.utils.defaultAbiCoder.encode(
    SERVICE_AGREEMENT_TYPES,
    serviceAgreementValues(sa),
  )
}

/**
 * Generate the unique identifier of a service agreement by computing its
 * digest.
 *
 * @param sa The service agreement to compute the digest of
 */
export function generateSAID(
  sa: ServiceAgreement,
): ReturnType<typeof ethers.utils.keccak256> {
  return ethers.utils.solidityKeccak256(
    SERVICE_AGREEMENT_TYPES,
    serviceAgreementValues(sa),
  )
}

/**
 * ABI encode the javascript representation of OracleSignatures
 *```solidity
 *  struct OracleSignatures {
 *    uint8[] vs;
 *    bytes32[] rs;
 *    bytes32[] ss;
 *  }
 * ```
 *
 * @param os The oracle signatures to ABI encode
 */
export function encodeOracleSignatures(os: OracleSignatures) {
  const ORACLE_SIGNATURES_TYPES = ['uint8[]', 'bytes32[]', 'bytes32[]']
  const osValues = [os.vs, os.rs, os.ss]

  return ethers.utils.defaultAbiCoder.encode(ORACLE_SIGNATURES_TYPES, osValues)
}

/**
 * Abi encode the oracleRequest() method for `Coordinator.sol`
 * ```solidity
 *  function oracleRequest(
 *    address _sender,
 *    uint256 _amount,
 *    bytes32 _sAId,
 *    address _callbackAddress,
 *    bytes4 _callbackFunctionId,
 *    uint256 _nonce,
 *    uint256 _dataVersion,
 *    bytes calldata _data
 *  )
 * ```
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
    oracleRequestInputs.map((i) => i.type),
    [ethers.constants.AddressZero, 0, specId, to, fHash, nonce, 1, dataBytes],
  )

  return `${oracleRequestSighash}${stripHexPrefix(encodedParams)}`
}

/**
 * Generates the oracle signatures on a ServiceAgreement
 *
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

/**
 * Signs a message according to ethereum specs by first appending
 * "\x19Ethereum Signed Message:\n' + <message.length>" to the message
 *
 * @param message The message to sign - either a Buffer or a hex string
 * @param wallet The wallet of the signer
 */
export async function personalSign(
  message: Buffer | string,
  wallet: ethers.Wallet,
): Promise<Required<utils.Signature>> {
  if (message instanceof String && !utils.isHexString(message)) {
    throw Error(`The message ${message} is not a valid hex string`)
  }

  const flatSig = await wallet.signMessage(utils.arrayify(message))
  const splitSignature = utils.splitSignature(flatSig)

  function assertIsSignature(
    sig: utils.Signature,
  ): asserts sig is Required<utils.Signature> {
    if (!sig.v) throw Error(`Could not extract v from signature`)
  }
  assertIsSignature(splitSignature)

  return splitSignature
}

/**
 * Recovers the address of the signer of a message
 *
 * @param message The message that was signed
 * @param signature The signature on the message
 */
export function recoverAddressFromSignature(
  message: string | Buffer,
  signature: Required<utils.Signature>,
): string {
  const messageBuff = utils.arrayify(message)
  return utils.verifyMessage(messageBuff, signature)
}

/**
 * Combine v, r, and s params of multiple signatures into format expected by contracts
 *
 * @param signatures The list of signatures to combine
 */
export function combineOracleSignatures(
  signatures: Required<utils.Signature>[],
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

function serviceAgreementValues(sa: ServiceAgreement) {
  return [
    sa.payment,
    sa.expiration,
    sa.endAt,
    sa.oracles.map((o) => (o instanceof ethers.Wallet ? o.address : o)),
    sa.requestDigest,
    sa.aggregator,
    sa.aggInitiateJobSelector,
    sa.aggFulfillSelector,
  ]
}
