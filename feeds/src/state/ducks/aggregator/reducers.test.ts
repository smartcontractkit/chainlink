import { OracleNode } from 'config'
import reducer, { INITIAL_STATE } from './reducers'
import {
  fetchOracleNodesBegin,
  fetchOracleNodesSuccess,
  fetchOracleNodesError,
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
        name: 'Fiews',
        networkId: 1,
      }
      const oracleNodeB: OracleNode = {
        address: '-',
        oracleAddress: 'B',
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
})
