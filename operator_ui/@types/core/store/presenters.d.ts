declare module 'core/store/presenters' {
  import * as assets from 'core/store/assets'
  import * as models from 'core/store/models'
  import * as orm from 'core/store/orm'
  import * as common from 'github.com/ethereum/go-ethereum/common'
  import * as hexutil from 'github.com/ethereum/go-ethereum/common/hexutil'
  import * as big from 'math/big'
  import * as time from 'time'

  /**
   * AccountBalance holds the hex representation of the address plus it's ETH & LINK balances
   */
  export interface AccountBalance {
    address: string
    ethBalance: Pointer<assets.Eth>
    linkBalance: Pointer<assets.Link>
    createdAt: string
    isFunding: boolean
  }

  /**
   * ConfigPrinter are the non-secret values of the node configuration
   */
  export interface ConfigPrinter extends EnvPrinter {
    accountAddress: string
  }

  /**
   * EnvPrinter are the non-secret values of the node environment
   */
  interface EnvPrinter {
    allowOrigins: string
    blockBackfillDepth: string
    bridgeResponseURL?: string
    chainlinkDev: boolean
    chainlinkPort: number
    chainlinkTLSHost: string
    chainlinkTLSPort: number
    clientNodeUrl: string
    databaseTimeout: time.Duration
    defaultHttpLimit: number
    defaultHttpTimeout: time.Duration
    enableExperimentalAdapters: boolean
    ethChainId: number
    ethFinalityDepth: number
    ethGasBumpThreshold: number /** * FIXME -- precision loss */
    ethGasBumpTxDepth: number
    ethGasBumpWei: Pointer<big.Int>
    ethGasLimitDefault: number
    ethGasPriceDefault: Pointer<big.Int>
    ethHeadTrackerHistoryDepth: number
    ethHeadTrackerMaxBufferSize: number
    ethUrl: string
    ethereumDisabled: boolean
    explorerUrl: string
    featureExternalInitiators: boolean
    featureFluxMonitor: boolean
    jsonConsole: boolean
    linkContractAddress: string
    logLevel: orm.LogLevel
    logSqlMigrations: boolean
    logSqlStatements: boolean
    logToDisk: boolean
    maxRPCCallsPerSecond: number
    maximumServiceDuration: time.Duration
    minIncomingConfirmations: number
    minOutgoingConfirmations: number
    minimumContractPayment: Pointer<assets.Link>
    minimumRequestExpiration: number /** * FIXME -- precision loss */
    minimumServiceDuration: time.Duration
    oracleContractAddress: Pointer<common.Address>
    reaperExpiration: time.Duration
    replayFromBlock: number
    root: string
    secureCookies: boolean
    sessionTimeout: time.Duration
    txAttemptLimit: number
  }

  /**
   * JobSpec holds the JobSpec definition together with
   * the total link earned from that job
   */
  export interface JobSpec extends models.JobSpec {
    earnings: Pointer<assets.Link>
  }

  /**
   * JobRun presents an API friendly version of the data.
   */
  export interface JobRun extends models.JobRun {}

  /**
   * Tx is a jsonapi wrapper for an Ethereum Transaction.
   */
  export interface Tx {
    state?: string
    data?: hexutil.Bytes
    from?: Pointer<common.Address>
    gasLimit?: string
    gasPrice?: string
    hash?: common.Hash
    rawHex?: string
    nonce?: string
    sentAt?: string
    to?: Pointer<common.Address>
    value?: string
  }
}
