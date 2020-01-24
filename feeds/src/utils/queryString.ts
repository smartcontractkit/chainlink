export function parseQuery(query: string) {
  const urlParams = new URLSearchParams(query)
  return Object.fromEntries(urlParams)
}

export function stringifyQuery(queryObject: Record<string, string>) {
  return new URLSearchParams(queryObject).toString()
}
