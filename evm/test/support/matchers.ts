import web3 from 'web3'
import { assert } from 'chai'

const bigNum = (num: number): any => web3.utils.toBN(num)

// Throws if a and b are not equal, as BN's
export const assertBigNum = (a, b, failureMessage) =>
  assert(
    bigNum(a).eq(bigNum(b)),
    `BigNum ${a} is not ${b}` + (failureMessage ? ': ' + failureMessage : '')
  )
