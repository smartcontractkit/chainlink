import { OracleNode, FeedConfig } from 'config'
import reducer, { INITIAL_STATE } from './reducers'
import {
  fetchOracleNodesBegin,
  fetchOracleNodesSuccess,
  fetchOracleNodesError,
  storeAggregatorConfig,
} from './actions'

describe('state/ducks/aggregator/reducers', () => {
  describe('FETCH_ORACLE_NODES_*', () => {
    it('toggles loadingOracleNodes when the request starts & finishes', () => {
      let state

      const beginAction = fetchOracleNodesBegin()
      state = reducer(INITIAL_STATE, beginAction)
      expect(state.loadingOracleNodes).toEqual(true)

      const payload: OracleNode[] = []
      const successAction = fetchOracleNodesSuccess(payload)
      state = reducer(state, successAction)
      expect(state.loadingOracleNodes).toEqual(false)

      state = reducer(state, beginAction)
      expect(state.loadingOracleNodes).toEqual(true)

      const errorAction = fetchOracleNodesError('Not Found')
      state = reducer(state, errorAction)
      expect(state.loadingOracleNodes).toEqual(false)
    })
  })

  describe('FETCH_ORACLE_NODES_SUCCESS', () => {
    it('indexes oracle nodes by oracleAddress', () => {
      const oracleNodeA: OracleNode = {
        address: '-',
        oracleAddress: 'A',
        nodeAddress: ['C'],
        name: 'Fiews',
        networkId: 1,
      }
      const oracleNodeB: OracleNode = {
        address: '-',
        oracleAddress: 'B',
        nodeAddress: ['D'],
        name: 'LinkPool',
        networkId: 1,
      }
      const payload: OracleNode[] = [oracleNodeA, oracleNodeB]
      const successAction = fetchOracleNodesSuccess(payload)

      const state = reducer(INITIAL_STATE, successAction)
      expect(state.oracleNodes).toEqual({
        A: oracleNodeA,
        B: oracleNodeB,
      })
    })
  })

  describe('STORE_AGGREGATOR_CONFIG', () => {
    it('stores aggregator config', () => {
      const config: FeedConfig = {
        networkId: 1,
        contractVersion: 2,
        decimalPlaces: 4,
        heartbeat: 0,
        historyDays: 1,
        formatDecimalPlaces: 0,
        threshold: 0,
        multiply: '100000000',
        contractAddress: 'address',
        contractType: 'aggregator',
        valuePrefix: '$',
        name: 'ETH/USD',
        pair: ['eth', 'usd'],
        path: 'eth-usd',
        history: true,
        listing: true,
        sponsored: ['Synthetic'],
      }

      const successAction = storeAggregatorConfig(config)

      const state = reducer(INITIAL_STATE, successAction)
      expect(state.config).toEqual(config)
    })
  })
})
