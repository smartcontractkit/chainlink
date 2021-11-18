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

  export interface JobSpecError {
    id: string
    description: string
    occurrences: number
    createdAt: time.Time
    updatedAt: time.Time
  }

  /**
   * TaskSpec is the definition of work to be carried out. The
   * Type will be an adapter, and the Params will contain any
   * additional information that adapter would need to operate.
   */
  export interface TaskSpec extends gorm.Model {
    type: TaskType
    confirmations: number | null
    params: { [key: string]: JSONValue | undefined }
  }

  /**
   * TaskType defines what Adapter a TaskSpec will use.
   */
  type TaskType = string

  /**
   * WebURL contains the URL of the endpoint.
   */
  type WebURL = url.URL

  /**
   * AnyTime holds a common field for time, and serializes it as
   * a json number.
   */
  type AnyTime = number

  /**
   * Cron holds the string that will represent the spec of the cron-job.
   * It uses 6 fields to represent the seconds (1), minutes (2), hours (3),
   * day of the month (4), month (5), and day of the week (6).
   */
  type Cron = string

  /**
   * Big stores large integers and can deserialize a variety of inputs.
   */
  type Big = big.Int

  /**
   * AddressCollection is an array of common.Address
   * serializable to and from a database.
   */
  type AddressCollection = common.Address[]
  //#endregion common.go

  //#region bridge_type.go
  /**
   * BridgeTypeRequest is the incoming record used to create a BridgeType
   */
  export interface BridgeTypeRequest {
    name: TaskType
    url: WebURL
    confirmations: number
    minimumContractPayment: Pointer<assets.Link>
  }

  /**
   * BridgeTypeAuthentication is the record returned in response to a request to create a BridgeType
   */
  export interface BridgeTypeAuthentication {
    name: TaskType
    url: WebURL
    confirmations: number
    incomingToken: string
    outgoingToken: string
    minimumContractPayment: Pointer<assets.Link>
  }

  /**
   * BridgeType is used for external adapters and has fields for
   * the name of the adapter and its URL.
   */
  export interface BridgeType {
    name: TaskType
    url: WebURL
    confirmations: number
    outgoingToken: string
    minimumContractPayment: Pointer<assets.Link>
  }
  //#endregion bridge_type.go

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
  //#region external_initiator.go
  /**
   * ExternalInitiator represents a user that can initiate runs remotely
   */
  export interface ExternalInitiator extends gorm.Model {
    AccessKey: string
    Salt: string
    HashedSecret: string
  }
  //#endregion external_initiator.go
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
  //#region ocrkey/key_bundle.go

  /**
   * OcrKey represents the bundle of keys needed for OCR
   */

  export interface OcrKey {
    configPublicKey: string
    offChainPublicKey: string
    onChainSigningAddress: common.Address
  }
  //#endregion ocrkey/key_bundle.go
  //#region p2pKey/p2p_key.go

  /**
   * P2P represents the bundle of keys needed for P2P
   */

  export interface P2PKey {
    peerId: string
    publicKey: string
  }
  //#endregion p2pKey/p2p_key.go

  /**
   * CreateJobRequest represents a schema for the create job request as used by
   * the API.
   */
  export interface CreateJobRequest {
    toml: string
  }

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

  export type PipelineTaskOutput = string | null
  export type PipelineTaskError = string | null

  interface BaseJob {
    name: string | null
    errors: JobSpecError[]
    maxTaskDuration: string
    pipelineSpec: {
      dotDagSource: string
    }
    schemaVersion: number
    externalJobID: string
  }

  export type DirectRequestJob = BaseJob & {
    type: 'directrequest'
    directRequestSpec: {
      initiator: 'runlog'
      contractAddress: common.Address
      minIncomingConfirmations: number | null
      minIncomingConfirmationsEnv?: boolean
      createdAt: time.Time
      requesters: common.Address[]
      evmChainID: string
    }
    fluxMonitorSpec: null
    offChainReportingOracleSpec: null
    keeperSpec: null
    cronSpec: null
    webhookSpec: null
    vrfSpec: null
  }

  export type FluxMonitorJob = BaseJob & {
    type: 'fluxmonitor'
    fluxMonitorSpec: {
      contractAddress: common.Address
      precision: number
      threshold: number
      absoluteThreshold: number
      idleTimerDisabled: false
      idleTimerPeriod: string
      pollTimerDisabled: false
      pollTimerPeriod: string
      drumbeatEnabled: boolean
      drumbeatSchedule?: string
      drumbeatRandomDelay?: string
      minPayment: number | null
      createdAt: time.Time
      evmChainID: string
    }
    cronSpec: null
    webhookSpec: null
    directRequestSpec: null
    offChainReportingOracleSpec: null
    keeperSpec: null
    vrfSpec: null
  }

  export type OffChainReportingJob = BaseJob & {
    type: 'offchainreporting'
    offChainReportingOracleSpec: {
      contractAddress: common.Address
      p2pPeerID: string
      p2pBootstrapPeers: string[]
      isBootstrapPeer: boolean
      keyBundleID: string
      monitoringEndpoint: string
      transmitterAddress: common.Address
      observationTimeout: string
      observationTimeoutEnv?: boolean
      blockchainTimeout: string
      blockchainTimeoutEnv?: boolean
      contractConfigTrackerSubscribeInterval: string
      contractConfigTrackerSubscribeIntervalEnv?: boolean
      contractConfigTrackerPollInterval: string
      contractConfigTrackerPollIntervalEnv?: boolean
      contractConfigConfirmations: number
      contractConfigConfirmationsEnv?: boolean
      createdAt: time.Time
      updatedAt: time.Time
      evmChainID: string
    }
    cronSpec: null
    webhookSpec: null
    vrfSpec: null
    directRequestSpec: null
    fluxMonitorSpec: null
    keeperSpec: null
  }

  export type KeeperJob = BaseJob & {
    type: 'keeper'
    keeperSpec: {
      contractAddress: common.Address
      fromAddress: common.Address
      createdAt: time.Time
      updatedAt: time.Time
      evmChainID: string
    }
    cronSpec: null
    webhookSpec: null
    vrfSpec: null
    directRequestSpec: null
    fluxMonitorSpec: null
    offChainReportingOracleSpec: null
  }

  export type CronJob = BaseJob & {
    type: 'cron'
    keeperSpec: null
    cronSpec: {
      schedule: string
      createdAt: time.Time
      updatedAt: time.Time
    }
    webhookSpec: null
    directRequestSpec: null
    vrfSpec: null
    fluxMonitorSpec: null
    offChainReportingOracleSpec: null
  }

  export type WebhookJob = BaseJob & {
    type: 'webhook'
    keeperSpec: null
    webhookSpec: {
      createdAt: time.Time
      updatedAt: time.Time
    }
    cronSpec: null
    vrfSpec: null
    directRequestSpec: null
    fluxMonitorSpec: null
    offChainReportingOracleSpec: null
  }

  export type VRFJob = BaseJob & {
    type: 'vrf'
    keeperSpec: null
    vrfSpec: {
      minIncomingConfirmations: number
      publicKey: string
      coordinatorAddress: common.Address
      fromAddress: string
      pollPeriod: string
      createdAt: time.Time
      updatedAt: time.Time
      evmChainID: string
    }
    cronSpec: null
    directRequestSpec: null
    fluxMonitorSpec: null
    webhookSpec: null
    offChainReportingOracleSpec: null
  }

  export type Job =
    | DirectRequestJob
    | FluxMonitorJob
    | OffChainReportingJob
    | KeeperJob
    | CronJob
    | WebhookJob
    | VRFJob

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

  export interface JobRunV2 {
    state: string
    outputs: PipelineTaskOutput[]
    errors: PipelineTaskError[]
    taskRuns: PipelineTaskRun[]
    createdAt: time.Time
    finishedAt: nullable.Time
    pipelineSpec: {
      ID: number
      CreatedAt: time.Time
      dotDagSource: string
      jobID: string
    }
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

  export interface CSAKey {
    publicKey: string
  }

  export interface JobProposal {
    spec: string
    status: string
    external_job_id: string | null
    createdAt: time.Time
    proposedAt: time.Time
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

  export interface FeatureFlag {
    enabled: boolean
  }
}

export interface PipelineTaskRun {
  createdAt: time.Time
  error: PipelineTaskError
  finishedAt: nullable.Time
  output: PipelineTaskOutput
  dotId: string
  type: string
}
