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
  }

  export interface ConfigWhitelist extends Whitelist {
    accountAddress: string
  }

  /**
   * ConfigWhitelist are the non-secret values of the node
   */
  interface Whitelist {
    allowOrigins: string
    bridgeResponseURL?: string
    ethChainId: number
    clientNodeUrl: string
    chainlinkDev: boolean
    databaseTimeout: time.Duration
    ethUrl: string
    /**
     * FIXME -- precision loss
     */
    ethGasBumpThreshold: number
    ethGasBumpWei: Pointer<big.Int>
    ethGasPriceDefault: Pointer<big.Int>
    jsonConsole: boolean
    linkContractAddress: string
    explorerUrl: string
    logLevel: orm.LogLevel
    logToDisk: boolean
    minimumContractPayment: Pointer<assets.Link>
    /**
     * FIXME -- precision loss
     */
    minimumRequestExpiration: number
    minIncomingConfirmations: number
    minOutgoingConfirmations: number
    oracleContractAddress: Pointer<common.Address>
    chainlinkPort: number
    reaperExpiration: time.Duration
    root: string
    sessionTimeout: time.Duration
    chainlinkTLSHost: string
    chainlinkTLSPort: number
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
    confirmed?: boolean
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
