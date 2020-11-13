declare module 'core/store/models' {
  import * as assets from 'core/store/assets'
  import * as common from 'github.com/ethereum/go-ethereum/common'
  import * as gorm from 'github.com/jinzhu/gorm'
  import * as clnull from 'github.com/smartcontractkit/chainlink/core/null'
  import * as nullable from 'gopkg.in/guregu/null.v3'
  import * as big from 'math/big'
  import * as url from 'net/url'
  import * as time from 'time'

  /**
   * RunStatus is a string that represents the run status
   */
  export enum RunStatus {
    IN_PROGRESS = 'in_progress',
    PENDING_INCOMING_CONFIRMATIONS = 'pending_incoming_confirmations',
    PENDING_CONNECTION = 'pending_connection',
    PENDING_BRIDGE = 'pending_bridge',
    PENDING_SLEEP = 'pending_sleep',
    ERRORED = 'errored',
    COMPLETED = 'completed',
  }

  //#region job_spec.go

  export interface JobSpecError {
    id: string
    description: string
    occurrences: number
    createdAt: time.Time
    updatedAt: time.Time
  }

  /**
   * JobSpecRequest represents a schema for the incoming job spec request as used by the API.
   */
  export interface JobSpecRequest {
    initiators: InitiatorRequest[]
    task: TaskSpecRequest[]
    startAt: nullable.Time
    endAt: nullable.Time
    minPayment: Pointer<assets.Link>
  }

  /**
   * InitiatorRequest represents a schema for incoming initiator requests as used by the API.
   */
  export interface InitiatorRequest {
    type: string
    params?: InitiatorParams
  }

  /**
   * TaskSpecRequest represents a schema for incoming TaskSpec requests as used by the API.
   */
  export interface TaskSpecRequest<T extends JSONValue = JSONValue> {
    type: TaskType
    confirmations: clnull.Uint32
    params: T
  }

  /**
   * JobSpec is the definition for all the work to be carried out by the node
   * for a given contract. It contains the Initiators, Tasks (which are the
   * individual steps to be carried out), StartAt, EndAt, and CreatedAt fields.
   */
  export interface JobSpec {
    id?: string // FIXME -- why is this nullable?
    createdAt: time.Time
    initiators: Initiator[]
    minPayment: Pointer<assets.Link>
    tasks: TaskSpec[]
    startAt: nullable.Time
    endAt: nullable.Time
    name: string
    earnings: number | null
    errors: JobSpecError[]
    updatedAt: time.time
  }

  // Types of Initiators (see Initiator struct just below.)
  export enum InitiatorType {
    /**
     * InitiatorRunLog for tasks in a job to watch an ethereum address
     * and expect a JSON payload from a log event.
     */
    RUN_LOG = 'runlog',
    /**
     * InitiatorCron for tasks in a job to be ran on a schedule.
     */
    CRON = 'cron',
    /**
     * InitiatorEthLog for tasks in a job to use the Ethereum blockchain.
     */
    ETH_LOG = 'ethlog',
    /**
     * InitiatorRunAt for tasks in a job to be ran once.
     */
    RUN_AT = 'runat',
    /**
     * InitiatorWeb for tasks in a job making a web request.
     */
    WEB = 'web',
    /**
     * InitiatorServiceAgreementExecutionLog for tasks in a job to watch a
     * Solidity Coordinator contract and expect a payload from a log event.
     */
    SERVICE_AGREEMENT_EXECUTION_LOG = 'execagreement',
  }

  /**
   * Initiator could be thought of as a trigger, defines how a Job can be
   * started, or rather, how a JobRun can be created from a Job.
   * Initiators will have their own unique ID, but will be associated
   * to a parent JobID.
   */
  export interface Initiator {
    id: number
    jobSpecId: string
    type: InitiatorType
    // FIXME - missing json struct tag
    CreatedAt: time.Time
    params?: InitiatorParams
  }

  /**
   * InitiatorParams is a collection of the possible parameters that different
   * Initiators may require.
   */
  interface InitiatorParams {
    schedule?: Cron
    time?: AnyTime
    ran?: boolean
    address?: common.Address
    requesters?: AddressCollection
  }

  /**
   * TaskSpec is the definition of work to be carried out. The
   * Type will be an adapter, and the Params will contain any
   * additional information that adapter would need to operate.
   */
  export interface TaskSpec<T extends JSONValue = JSONValue>
    extends gorm.Model {
    type: TaskType
    confirmations: number
    params: T
  }

  /**
   * TaskType defines what Adapter a TaskSpec will use.
   */
  type TaskType = string

  //#endregion job_spec.go

  //#region run.go

  /**
   * JobRun tracks the status of a job by holding its TaskRuns and the
   * Result of each Run.
   */
  export interface JobRun {
    id: string
    jobId: string
    result: RunResult
    status: RunStatus
    taskRuns: TaskRun[]
    createdAt: time.Time
    finishedAt: nullable.Time
    updatedAt: time.Time
    initiator: Initiator
    createdHeight: Pointer<Big>
    observedHeight: Pointer<Big>
    overrides: RunResult
    payment: Pointer<assets.Link>
  }

  /**
   * TaskRun stores the Task and represents the status of the
   * Task to be ran.
   */
  export interface TaskRun {
    id: string
    result: RunResult
    status: RunStatus
    task: TaskSpec
    minimumConfirmations: clnull.Uint32
    confirmations: clnull.Uint32
  }

  /**
   * RunResult keeps track of the outcome of a TaskRun or JobRun. It stores the
   * Data and ErrorMessage, and contains a field to track the status.
   */
  export interface RunResult<T extends JSONValue = JSONValue> {
    jobRunId: string
    taskRunId: string
    data: JSONValue
    status: RunStatus
    error: nullable.String
  }

  /**
   * BridgeRunResult handles the parsing of RunResults from external adapters.
   */
  export interface BridgeRunResult extends RunResult {
    pending: boolean
    accessToken: string
  }

  //#endregion run.go

  //#region common.go
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
   * WithdrawalRequest request to withdraw LINK.
   */
  export interface WithdrawalRequest {
    address: common.Address
    contractAddress: common.Address
    amount: Pointer<assets.Link>
  }

  /**
   * SendEtherRequest represents a request to transfer ETH.
   */
  export interface SendEtherRequest {
    address: common.Address
    from: common.Address
    amount: Pointer<assets.Eth>
  }

  /**
   * CreateKeyRequest represents a request to add an ethereum key.
   */
  export interface CreateKeyRequest {
    current_password: string // FIXME -- camelcase
  }

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
  //#region eth.go
  /**
   * Log represents a contract log event. These events are generated by the LOG opcode and
   * stored/indexed by the node.
   */
  export interface Log {
    /**
     * Consensus fields:
     * address of the contract that generated the event
     */
    address: common.Address
    /**
     * List of topics provided by the contract
     */
    topics: common.Hash[]
    /**
     * Supplied by the contract, usually ABI-encoded
     */
    data: string

    /**
     * Derived fields. These fields are filled in by the node
     * but not secured by consensus.
     * block in which the transaction was included
     */
    blockNumber: number // FIXME -- uint64 is unsafe in this context
    /**
     * Hash of the transaction
     */
    transactionHash: common.Hash
    /**
     * Index of the transaction in the block
     */
    transactionIndex: number
    /**
     * Hash of the block in which the transaction was included
     */
    blockHash: common.Hash
    /**
     * Index of the log in the receipt
     */
    logIndex: number

    /**
     * The Removed field is true if this log was reverted due to a chain reorganisation.
     * You must pay attention to this field if you receive logs through a filter query.
     */
    removed: boolean
  }

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
    createdAt: time.Time
    offChainPublicKey: string
    onChainSigningAddress: common.Address
    updatedAt: time.Time
  }
  //#endregion ocrkey/key_bundle.go
  //#region p2pKey/p2p_key.go

  /**
   * P2P represents the bundle of keys needed for P2P
   */

  export interface P2PKey {
    peerId: string
    publicKey: string
    createdAt: time.Time
    updatedAt: time.Time
    deletedAt: time.Time
  }
  //#endregion p2pKey/p2p_key.go

  /**
   * OcrJobSpecRequest represents a schema for the incoming ocr job spec request as used by the API.
   */
  export interface OcrJobSpecRequest {
    toml: string
  }

  type OcrTaskOutput = ?string
  type OcrTaskError = ?string

  export interface OcrJobSpec {
    errors: JobSpecError[]
    offChainReportingOracleSpec: {
      contractAddress: common.Address
      p2pPeerID: string
      p2pBootstrapPeers: string[]
      isBootstrapPeer: boolean
      keyBundleID: string
      monitoringEndpoint: string
      transmitterAddress: common.Address
      observationTimeout: string
      blockchainTimeout: string
      contractConfigTrackerSubscribeInterval: string
      contractConfigTrackerPollInterval: string
      contractConfigConfirmations: number
      createdAt: time.Time
      updatedAt: time.Time
      name?: string // Upcoming field
    }
    pipelineSpec: {
      dotDagSource: string
    }
  }

  export interface OcrJobRun {
    outputs: OcrTaskOutput[]
    errors: OcrTaskError[]
    taskRuns: OcrTaskRun[]
    createdAt: time.Time
    finishedAt: nullable.Time
    pipelineSpec: {
      CreatedAt: time.Time
      DotDagSource: string
      ID: number
    }
  }
}

export interface OcrTaskRun {
  createdAt: time.Time
  error: OcrTaskError
  finishedAt: nullable.Time
  output: OcrTaskOutput
  taskSpec: {
    dotId: string
  }
  type: string
}
