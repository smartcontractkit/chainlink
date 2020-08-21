import { partialAsFull } from '@chainlink/ts-helpers'
import { AppState } from 'state'
import { OracleNode } from '../../../config'

import { upcasedOracles } from './selectors'

describe('upcasedOracles', () => {
  const aggregatedState = partialAsFull<AppState>({
    aggregator: partialAsFull<AppState['aggregator']>({
      oracleNodes: {
        oracleAddress1: partialAsFull<OracleNode>({
          name: 'Oracle 1',
          nodeAddress: ['nodeAddress1'],
          oracleAddress: 'oracleAddress1',
        }),
        oracleAddress2: partialAsFull<OracleNode>({
          name: 'Oracle 2',
          nodeAddress: ['nodeAddress2'],
          oracleAddress: 'oracleAddress2',
        }),
      },
      config: {
        contractVersion: 2,
      },
    }),
  })

  it('for v2 contract returns a record of oracle addresses and names', () => {
    const expectedResult = {
      oracleAddress1: 'Oracle 1',
      oracleAddress2: 'Oracle 2',
    }

    expect(upcasedOracles(aggregatedState)).toEqual(expectedResult)
  })

  it('for v3 contract returns a record of node addresses and names', () => {
    const aggregatedStateV3 = JSON.parse(JSON.stringify(aggregatedState))
    aggregatedStateV3.aggregator.config.contractVersion = 3

    const expectedResult = {
      nodeAddress1: 'Oracle 1',
      nodeAddress2: 'Oracle 2',
    }

    expect(upcasedOracles(aggregatedStateV3)).toEqual(expectedResult)
  })
})
