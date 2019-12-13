import { assert } from 'chai'

import { ethers } from 'ethers'
import { BigNumberish } from 'ethers/utils'

// Throws if a and b are not equal, as BN's
export const assertBigNum = (
  a: BigNumberish,
  b: BigNumberish,
  failureMessage?: string,
) =>
  assert(
    ethers.utils.bigNumberify(a).eq(ethers.utils.bigNumberify(b)),
    `BigNum ${a} is not ${b}` + (failureMessage ? ': ' + failureMessage : ''),
  )
