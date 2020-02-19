export const parseQuery = query => {
  const urlParams = new URLSearchParams(query)
  return Object.fromEntries(urlParams)
}

export const stringifyQuery = queryObject => {
  return new URLSearchParams(queryObject).toString()
}
