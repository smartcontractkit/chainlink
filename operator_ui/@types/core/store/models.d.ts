declare module 'core/store/models' {
  import * as assets from 'core/store/assets'
  import * as common from 'github.com/ethereum/go-ethereum/common'
  import * as gorm from 'github.com/jinzhu/gorm'
  import * as clnull from 'github.com/smartcontractkit/chainlink/core/null'
  import * as nullable from 'gopkg.in/guregu/null.v3'
  import * as big from 'math/big'
  import * as url from 'net/url'
  import * as time from 'time'

  //#region job_spec.go

  export interface Resource<T> {
    id: string
    type: string
    attributes: T
  }

  /**
   * Big stores large integers and can deserialize a variety of inputs.
   */
  type Big = big.Int

  /**
   * Tx contains fields necessary for an Ethereum transaction with
   * an additional field for the TxAttempt.
   *
   * FIXME -- all the fields below need to be in camelCase
   */
  export interface Tx {
    ID: number // FIXME -- uint64 is unsafe in this context + should be camelCased

    /**
     * SurrogateID is used to look up a transaction using a secondary ID, used to
     * associate jobs with transactions so that we don't double spend in certain
     * failure scenarios
     */
    SurrogateID: nullable.String // FIXME -- camelCase

    From: common.Address
    To: common.Address
    Data: string
    Nonce: number // FIXME -- possibly unsafe uint64
    Value: Pointer<Big>
    GasLimit: number // FIXME -- possibly unsafe uint64

    /**
     * TxAttempt fields manually included; can't embed another primary_key
     */
    Hash: common.Hash
    GasPrice: Pointer<Big>
    Confirmed: boolean
    SentAt: number // FIXME -- possibly unsafe uint64
    SignedRawTx: string
  }

  /**
   * TxAttempt is used for keeping track of transactions that
   * have been written to the Ethereum blockchain. This makes
   * it so that if the network is busy, a transaction can be
   * resubmitted with a higher GasPrice.
   *
   * FIXME -- fields need to be camelcase
   */
  export interface TxAttempt {
    ID: number // FIXME -- possibly unsafe uint64
    TxID: number // FIXME -- possibly unsafe uint64

    CreatedAt: time.Time

    Hash: common.Hash
    GasPrice: Pointer<Big>
    Confirmed: boolean
    SentAt: number // FIXME -- possibly unsafe uint64
    SignedRawTx: string
  }
  //#endregion eth.go

  //#region user.go
  /**
   * SessionRequest encapsulates the fields needed to generate a new SessionID,
   * including the hashed password.
   */
  export interface SessionRequest {
    email: string
    password: string
    webauthndata: string
  }
  //#endregion user.go
  //#region bulk.go
  /**
   * BulkDeleteRunRequest describes the query for deletion of runs
   */
  export interface BulkDeleteRunRequest {
    /**
     * FIXME -- loss of precision and camelcase,
     * maybe this shouldnt exist at all actually
     */
    ID?: number
    status: RunStatusCollection
    updatedBefore: time.Time
  }

  /**
   * RunStatusCollection is an array of RunStatus.
   */
  export type RunStatusCollection = RunStatus[]

  //#endregion  bulk.go

  export interface CreateChainRequest {
    chainID: string
    config: Record<string, JSONPrimitive>
  }

  export interface UpdateChainRequest {
    config: Record<string, JSONPrimitive>
    enabled: boolean
  }
  export interface CreateNodeRequest {
    name: string
    evmChainID: string
    httpURL: string
    wsURL: string
  }

  export type Chain = {
    config: Record<string, JSONPrimitive>
    enabled: boolean
    createdAt: time.Time
    updatedAt: time.Time
  }

  export type Node = {
    name: string
    evmChainID: string
    httpURL: string
    wsURL: string
    createdAt: time.Time
    updatedAt: time.Time
  }

  // We really need to change the API for this. It not only returns levels but
  // true/false for IsSQLEnabled
  export type LogConfigLevel =
    | 'debug'
    | 'info'
    | 'warn'
    | 'error'
    | 'true'
    | 'false'
  export type LogServiceName =
    | 'Global'
    | 'IsSqlEnabled'
    | 'header_tracker'
    | 'fluxmonitor'

  export interface LogConfig {
    // Stupidly this also returns boolean strings
    logLevel: LogConfigLevel[]
    serviceName: string[]
    defaultLogLevel: string
  }

  export interface LogConfigRequest {
    level: LogConfigLevel
    sqlEnabled: boolean
  }

  /**
   * Request to begin the process of registering a new MFA token
   */
  export interface BeginWebAuthnRegistrationV2Request {}

  /**
   * Request to begin the process of registering a new MFA token
   */
  export interface BeginWebAuthnRegistrationV2 {}

  /**
   * Request to begin the process of registering a new MFA token
   */
  export interface FinishWebAuthnRegistrationV2Request {
    id: string
    rawId: string
    type: string
    response: {
      attestationObject: string
      clientDataJSON: string
    }
  }

  /**
   * Request to begin the process of registering a new MFA token
   */
  export interface FinishWebAuthnRegistrationV2 {
    id: string
    rawId: string
    type: string
    response: {
      attestationObject: string
      clientDataJSON: string
    }
  }

  export interface UpdateJobProposalSpecRequest {
    spec: string
  }
}

