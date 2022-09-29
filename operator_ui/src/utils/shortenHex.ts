export function shortenHex(
  value: string,
  {
    start = 6,
    end = 4,
  }: {
    start?: number
    end?: number
  } = {},
) {
  return value.substring(0, start) + '...' + value.substring(value.length - end)
}
