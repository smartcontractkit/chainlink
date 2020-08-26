import 'core-js/stable/object/from-entries'
import { createSelector } from 'reselect'
import { AppState } from 'state'
import { OracleNode } from '../../../config'

export const upcaseOracles = (
  state: AppState,
): Record<OracleNode['address'], OracleNode['name']> => {
  /**
   * In v2 of the contract, oracles' list has oracle addresses,
   * but in v3 - node addresses. Therefore, a different record of
   * pairs has to be made for each contract version.
   *
   * Custom pages that are used to test new contracts have their
   * `config` attribute set as `null`, so we need to check that as well.
   */
  if (
    state.aggregator.config === null ||
    state.aggregator.config.contractVersion === 2
  ) {
    return Object.fromEntries(
      Object.entries(
        state.aggregator.oracleNodes,
      ).map(([oracleAddress, oracleNode]) => [oracleAddress, oracleNode.name]),
    )
  } else {
    return Object.fromEntries(
      Object.values(state.aggregator.oracleNodes).map(oracleNode => [
        oracleNode.nodeAddress[0],
        oracleNode.name,
      ]),
    )
  }
}

const oracleList = (state: AppState) => state.aggregator.oracleList
const oracleAnswers = (state: AppState) => state.aggregator.oracleAnswers
const pendingAnswerId = (state: AppState) => state.aggregator.pendingAnswerId

const oracles = createSelector(
  [oracleList, upcaseOracles],
  (
    list: Array<OracleNode['address']>,
    upcasedOracles: Record<OracleNode['address'], OracleNode['name']>,
  ) => {
    if (!list) return []

    const result = list
      .map(address => {
        return {
          address,
          name: upcasedOracles[address] || 'Unknown',
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
