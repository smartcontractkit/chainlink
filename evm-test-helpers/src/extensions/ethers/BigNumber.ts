/**
 * @packageDocumentation
 *
 * An extension to ether's bignumber library that manually
 * polyfills any methods we need for tests by converting the
 * numbers back and forth between ethers.utils.BigNumber and
 * bn.js. If we end up having to replace a ton of methods in the
 * future this way, it might be worth creating a proxy object
 * that automatically does these method polyfills for us.
 */

import { ethers } from 'ethers'
import BN = require('bn.js')

/* eslint-disable @typescript-eslint/no-namespace */
declare module 'ethers' {
  namespace ethers {
    namespace utils {
      interface BigNumber {
        isEven(): boolean
        umod(val: ethers.utils.BigNumber): ethers.utils.BigNumber
        shrn(val: number): ethers.utils.BigNumber
        invm(val: ethers.utils.BigNumber): ethers.utils.BigNumber
      }
    }
  }
}

const BN_1 = new BN(-1)
// https://github.com/ethers-io/ethers.js/blob/v4.0.41/src.ts/utils/bignumber.ts#L42
function bnify(value: ethers.utils.BigNumber): BN {
  const hex = value.toHexString()
  if (hex[0] === '-') {
    return new BN(hex.substring(3), 16).mul(BN_1)
  }
  return new BN(hex.substring(2), 16)
}

// https://github.com/ethers-io/ethers.js/blob/v4.0.41/src.ts/utils/bignumber.ts#L22
function toHex(bn: BN): string {
  const value = bn.toString(16)
  if (value[0] === '-') {
    if (value.length % 2 === 0) {
      return '-0x0' + value.substring(1)
    }
    return '-0x' + value.substring(1)
  }
  if (value.length % 2 === 1) {
    return '0x0' + value
  }
  return '0x' + value
}

// https://github.com/ethers-io/ethers.js/blob/v4.0.41/src.ts/utils/bignumber.ts#L38
function toBigNumber(bn: BN): ethers.utils.BigNumber {
  return new ethers.utils.BigNumber(toHex(bn))
}

export function extend(bignumber: typeof ethers.utils.BigNumber) {
  bignumber.prototype.isEven = function () {
    return bnify(this).isEven()
  }

  bignumber.prototype.umod = function (val: any) {
    return toBigNumber(bnify(this).umod(bnify(val)))
  }

  bignumber.prototype.shrn = function (val: any) {
    return toBigNumber(bnify(this).shrn(val))
  }

  bignumber.prototype.invm = function (val: any) {
    return toBigNumber(bnify(this).invm(bnify(val)))
  }
}
