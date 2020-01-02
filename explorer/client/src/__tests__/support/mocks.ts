/**
 * A test helper that allows one to only partially satisfy a given type, and the function will
 * return the same value, but type casted as a full type.
 *
 * Useful for tests that require only a small slice of a given object to test conditions.
 *
 * @param val The partial value of a given type to be mocked as the full value
 */
export function partialAsFull<T>(val: Partial<T>): T {
  return val as T
}
