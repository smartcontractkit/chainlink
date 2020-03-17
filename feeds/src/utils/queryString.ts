export function parseQuery(query: string): Record<string, string> {
  const urlParams = new URLSearchParams(query)
  return Object.fromEntries(urlParams)
}

export function stringifyQuery(queryObject: Record<string, string>): string {
  return new URLSearchParams(queryObject).toString()
}
