import { createSelector } from 'reselect'
import nodes from '../../../nodes.json'

const oracles = (state: any) => state.aggregation.oracles
const oracleResponse = (state: any) => state.aggregation.oracleResponse
const currentAnswer = (state: any) => state.aggregation.currentAnswer
const contractAddress = (state: any) => state.aggregation.contractAddress
const pendingAnswerId = (state: any) => state.aggregation.pendingAnswerId

const oraclesList = createSelector([oracles], list => {
  if (!list) return []

  const names: Record<string, string> = {}

  nodes.forEach(n => {
    names[n.address.toUpperCase()] = n.name
  })

  const result = list
    .map((a: any) => {
      return {
        address: a,
        name: names[a.toUpperCase()] || 'Unknown',
        type: 'oracle',
      }
    })
    .sort((a: any, b: any) => a.name.localeCompare(b.name))

  return result
})

const networkGraphNodes = createSelector(
  [oraclesList, contractAddress],
  (list, address) => {
    if (!list) return []

    let result = [
      {
        type: 'contract',
        name: 'Aggregation Contract',
        address,
      },
      ...list,
    ]

    result = result.map((a, i) => {
      return { ...a, id: i }
    })

    return result
  },
)

const networkGraphState = createSelector(
  [oracleResponse, currentAnswer],
  (list, answer) => {
    if (!list) return []

    const contractData = {
      currentAnswer: answer,
      type: 'contract',
    }

    return [...list, contractData]
  },
)

const oraclesData = createSelector(
  [oraclesList, oracleResponse, pendingAnswerId],
  (list, response, pendingAnswerId) => {
    if (!list) return []

    const data = list.map((o: any, id: any) => {
      const state =
        response && response.filter((r: any) => r.sender === o.address)[0]
      const isFulfilled = state && state.answerId >= pendingAnswerId
      return { ...o, ...state, id, isFulfilled }
    })

    return data
  },
)

export { oraclesList, networkGraphNodes, networkGraphState, oraclesData }
