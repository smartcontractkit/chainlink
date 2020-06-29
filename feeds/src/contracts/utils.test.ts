import { ethers } from 'ethers'
import { formatAnswer } from './utils'

describe('contracts/utils', () => {
  describe('formatAnswer', () => {
    it('converts and formats the raw answer value', () => {
      const value = ethers.utils.bigNumberify('0x08d8f9fc00')
      const multiply = '1'
      const decimalPlaces = 0

      expect(formatAnswer(value, multiply, decimalPlaces, 0)).toEqual(
        '38000000000.0',
      )

      const formatDecimalPlaces = 9
      expect(
        formatAnswer(value, multiply, decimalPlaces, formatDecimalPlaces),
      ).toEqual('38.0')
    })
  })
})
