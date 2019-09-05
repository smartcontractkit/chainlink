declare module 'core/store/assets' {
  import * as big from 'math/big'
  //#region currencies.go
  /**
   * Link contains a field to represent the smallest units of LINK
   */
  export type Link = big.Int

  /**
   * Eth contains a field to represent the smallest units of ETH
   */
  export type Eth = big.Int
  //#endregion currencies.go
}
