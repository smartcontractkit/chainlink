/**
 * Parse a string into a finite, unsigned integer. Unsafe for integers greater than 32 bits of precision.
 *
 * @param str The string to parse into a finite unsigned integer.
 */

export const uIntFrom = str => {
  return /^[+]?(\d+)$/.test(str) ? Number(str) : NaN
}
