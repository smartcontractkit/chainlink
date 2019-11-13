export function partialAsFull<T>(val: Partial<T>): T {
  return val as T
}
