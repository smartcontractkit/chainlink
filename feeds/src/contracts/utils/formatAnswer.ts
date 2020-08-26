import { ethers } from 'ethers'

/**
 * Format an aggregator answer
 *
 * @param value The Big Number to format
 * @param multiply The number to divide the result by. See Multiply adapter in Chainlink Job Specification -  https://docs.chain.link/docs/job-specifications
 * @param decimalPlaces The number to show decimal places
 * @param formatDecimalPlaces
 */
export function formatAnswer(
  value: any,
  multiply: string,
  decimalPlaces: number,
  formatDecimalPlaces = 0,
): string {
  try {
    const decimals = 10 ** decimalPlaces
    const divided = value.mul(decimals).div(multiply)
    const formatted = ethers.utils.formatUnits(
      divided,
      decimalPlaces + formatDecimalPlaces,
    )

    return formatted.toString()
  } catch {
    return value
  }
}
