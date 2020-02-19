import { createSelector } from 'reselect'
import nodes from 'nodes.json'

const oracles = state => state.aggregation.oracles
const oracleResponse = state => state.aggregation.oracleResponse
const currentAnswer = state => state.aggregation.currentAnswer
const contractAddress = state => state.aggregation.contractAddress
const pendingAnswerId = state => state.aggregation.pendingAnswerId

const oraclesList = createSelector([oracles], list => {
  if (!list) return []

  const names = {}

  nodes.forEach(n => {
    names[n.address.toUpperCase()] = n.name
  })

  const result = list.map(a => {
    return {
      address: a,
      name: names[a.toUpperCase()] || 'Unknown',
      type: 'oracle',
    }
  })

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

    const data = list.map((o, id) => {
      const state = response && response.filter(r => r.sender === o.address)[0]
      const isFulfilled = state && state.answerId >= pendingAnswerId
      return { ...o, ...state, id, isFulfilled }
    })

    return data
  },
)

export { oraclesList, networkGraphNodes, networkGraphState, oraclesData }
