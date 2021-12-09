declare module 'core/store/presenters' {
  import * as assets from 'core/store/assets'
  import * as models from 'core/store/models'
  import * as orm from 'core/store/orm'
  import * as common from 'github.com/ethereum/go-ethereum/common'
  import * as hexutil from 'github.com/ethereum/go-ethereum/common/hexutil'
  import * as big from 'math/big'
  import * as time from 'time'

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
