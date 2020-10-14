import { decamelizeKeys } from 'humps'

export default (bridges, totalCount) => {
  const bc = totalCount || bridges.length
  return decamelizeKeys({
    meta: { count: bc },
    data: bridges.map((b) => {
      return {
        type: 'bridges',
        id: b.name,
        attributes: b,
      }
    }),
  })
}
