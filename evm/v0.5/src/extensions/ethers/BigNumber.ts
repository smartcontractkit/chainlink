import { ethers } from 'ethers'

/* eslint-disable @typescript-eslint/no-namespace */
declare module 'ethers' {
  namespace ethers {
    namespace utils {
      interface BigNumber {
        isEven(): boolean
        umod(val: ethers.utils.BigNumber): ethers.utils.BigNumber
        shrn(val: number): ethers.utils.BigNumber
      }
    }
  }
}

ethers.utils.BigNumber.prototype.isEven = function() {
  return (this as ethers.utils.BigNumber).mod(2).toNumber() === 0
}

ethers.utils.BigNumber.prototype.umod = function(val) {
  const remainder = (this as ethers.utils.BigNumber).mod(val)
  return remainder.lt(0) ? val.add(remainder).mod(val) : remainder
}

ethers.utils.BigNumber.prototype.shrn = function(val) {
  let finalValue: ethers.utils.BigNumber = new ethers.utils.BigNumber(this)

  while (val-- > 0) {
    finalValue = finalValue.div(2)
  }

  return finalValue
}
