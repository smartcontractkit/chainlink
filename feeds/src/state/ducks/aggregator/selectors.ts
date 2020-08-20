import 'core-js/stable/object/from-entries'
import { createSelector } from 'reselect'
import { AppState } from 'state'
import { OracleNode } from '../../../config'

type OracleAddress = string
type NodeAddress = string

type StateBranch = {
  aggregator: {
    oracleNodes: {
      [address: string]: {
        name: string
        nodeAddress: NodeAddress[]
        oracleAddress: OracleAddress
        [prop: string]: any
      }
    }
    config: {
      contractVersion: number
      [prop: string]: any
    }
  }
  [prop: string]: any
}

export const upcaseOracles = (
  state: StateBranch,
): Record<OracleNode['address'], OracleNode['name']> => {
  /**
   * In v2 of the contract, oracles' list has oracle addresses,
   * but in v3 - node addresses. Therefore, a different record of
   * pairs has to be made for each contract version.
   */
  if (state.aggregator.config.contractVersion === 2) {
    return Object.fromEntries(
      Object.entries(
        state.aggregator.oracleNodes,
      ).map(([oracleAddress, oracleNodeObject]) => [
        oracleAddress,
        oracleNodeObject.name,
      ]),
    )
  } else {
    return Object.fromEntries(
      Object.values(state.aggregator.oracleNodes).map(oracleNodeObject => [
        oracleNodeObject.nodeAddress[0],
        oracleNodeObject.name,
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
    upcaseOracles: Record<OracleNode['address'], OracleNode['name']>,
  ) => {
    if (!list) return []

    const result = list
      .map((address: OracleNode['address']) => {
        return {
          address,
          name: upcaseOracles[address] || 'Unknown',
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
