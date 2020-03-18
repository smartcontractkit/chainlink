/**
 * Parse a string into a finite, unsigned integer. Unsafe for integers greater than 32 bits of precision.
 *
 * @param str The string to parse into a finite unsigned integer.
 */
export function uIntFrom(val: string | number): number {
  return /^[+]?(\d+)$/.test(val.toString()) ? Number(val) : NaN
}
