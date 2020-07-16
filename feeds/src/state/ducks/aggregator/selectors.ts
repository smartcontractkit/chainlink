import 'core-js/stable/object/from-entries'
import { createSelector } from 'reselect'
import { AppState } from 'state'
import { OracleNode } from '../../../config'

const upcaseOracles = (
  state: AppState,
): Record<OracleNode['address'], OracleNode['name']> => {
  return Object.fromEntries(
    Object.keys(state.aggregator.oracleNodes).map(k => [
      k.toUpperCase(),
      state.aggregator.oracleNodes[k].name,
    ]),
  )
}
const oracleList = (state: AppState) => state.aggregator.oracleList
const oracleAnswers = (state: AppState) => state.aggregator.oracleAnswers
const pendingAnswerId = (state: AppState) => state.aggregator.pendingAnswerId

const oracles = createSelector(
  [oracleList, upcaseOracles],
  (
    list: Array<OracleNode['address']>,
    upcaseOracles: Record<OracleNode['address'], OracleNode['name']>,
  ) => {
    if (!list) return []

    const result = list
      .map(a => {
        return {
          address: a,
          name: upcaseOracles[a.toUpperCase()] || 'Unknown',
          type: 'oracle',
        }
      })
      .sort((a: any, b: any) => a.name.localeCompare(b.name))

    return result
  },
)

interface OracleAnswer {
  answerFormatted: string
  answer: number
  answerId: number
  sender: string
  meta: {
    blockNumer: number
    transactionHash: string
    timestamp: number
    gasPrice: string
  }
}

const latestOraclesState = createSelector(
  [oracles, oracleAnswers, pendingAnswerId],
  (list, answers: OracleAnswer[], pendingAnswerId) => {
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
