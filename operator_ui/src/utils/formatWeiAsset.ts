import { WEI_PER_TOKEN } from './constants'
import { BigNumber } from 'bignumber.js'

export default (val: number): number => {
  const b = new BigNumber(val)
  const minPay = b.dividedBy(WEI_PER_TOKEN).toNumber()
  return minPay
}
