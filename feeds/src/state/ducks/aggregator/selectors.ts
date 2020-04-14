import { createSelector } from 'reselect'
import { AppState } from 'state'
import nodes from '../../../nodes.json'

const oracleList = (state: AppState) => state.aggregator.oracleList
const oracleAnswers = (state: AppState) => state.aggregator.oracleAnswers
const pendingAnswerId = (state: AppState) => state.aggregator.pendingAnswerId

const oracles = createSelector([oracleList], list => {
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

const latestOraclesState = createSelector(
  [oracles, oracleAnswers, pendingAnswerId],
  (list, answers, pendingAnswerId) => {
    if (!list) return []

    const data = list.map((o: any, id: any) => {
      const state =
        answers &&
        answers.find(
          (r: any) => r.sender.toUpperCase() === o.address.toUpperCase(),
        )

      const isFulfilled = state && state.answerId >= pendingAnswerId
      return { ...o, ...state, id, isFulfilled }
    })

    return data
  },
)

export { latestOraclesState }
