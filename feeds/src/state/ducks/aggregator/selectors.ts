import 'core-js/stable/object/from-entries'
import { createSelector } from 'reselect'
import { AppState } from 'state'
import { OracleNode } from '../../../config'

export const upcaseOracles = (
  state: AppState,
): Record<OracleNode['oracleAddress'], OracleNode['name']> => {
  /**
   * In v2 of the contract, oracles' list has oracle addresses,
   * but in v3 - node addresses.
   */

  return Object.entries(state.aggregator.oracleNodes).reduce(
    (accumulator: Record<string, string>, [oracleAddress, oracle]) => {
      accumulator[oracleAddress] = oracle.name
      oracle.nodeAddress.forEach(nodeAddress => {
        accumulator[nodeAddress] = oracle.name
      })

      return accumulator
    },
    {},
  )
}

const oracleList = (state: AppState) => state.aggregator.oracleList
const oracleAnswers = (state: AppState) => state.aggregator.oracleAnswers
const pendingAnswerId = (state: AppState) => state.aggregator.pendingAnswerId

type Oracle = {
  address: OracleNode['oracleAddress']
  name: OracleNode['name']
  type: string
}

const oracles = createSelector(
  [oracleList, upcaseOracles],
  (
    list: Array<OracleNode['oracleAddress']>,
    upcasedOracles: Record<OracleNode['oracleAddress'], OracleNode['name']>,
  ): Oracle[] => {
    if (!list) return []
    const result = list
      .map(address => {
        const oracle: Oracle = {
          address,
          name: upcasedOracles[address] || 'Unknown',
          type: 'oracle',
        }
        return oracle
      })
      .sort((a, b) => a.name.localeCompare(b.name))

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
  (list: Oracle[], answers: OracleAnswer[], pendingAnswerId: number) => {
    if (!list) return []

    const data = list.map((o, id) => {
      const state =
        answers &&
        answers.find(r => r.sender.toUpperCase() === o.address.toUpperCase())

      const isFulfilled = state && state.answerId >= pendingAnswerId
      return { ...o, ...state, id, isFulfilled }
    })

    return data
  },
)

export { latestOraclesState }
