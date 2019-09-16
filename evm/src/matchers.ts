import web3 from 'web3'
import { assert } from 'chai'
import BigNumber from 'bn.js'
import BN from 'bn.js'

const bigNum = (num: number | BN): BigNumber => web3.utils.toBN(num)

// Throws if a and b are not equal, as BN's
export const assertBigNum = (a: BN, b: BN, failureMessage?: string) =>
  assert(
    bigNum(a).eq(bigNum(b)),
    `BigNum ${a} is not ${b}` + (failureMessage ? ': ' + failureMessage : ''),
  )
