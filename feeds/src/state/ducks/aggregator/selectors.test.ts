import { partialAsFull } from '@chainlink/ts-helpers'
import { AppState } from 'state'
import { OracleNode } from '../../../config'

import { upcaseOracles } from './selectors'

describe('upcaseOracles', () => {
  const aggregatedState = partialAsFull<AppState>({
    aggregator: partialAsFull<AppState['aggregator']>({
      oracleNodes: {
        oracleAddress1: partialAsFull<OracleNode>({
          name: 'Oracle 1',
          nodeAddress: ['nodeAddress1', 'nodeAddress12'],
          oracleAddress: 'oracleAddress1',
        }),
        oracleAddress2: partialAsFull<OracleNode>({
          name: 'Oracle 2',
          nodeAddress: ['nodeAddress2'],
          oracleAddress: 'oracleAddress2',
        }),
      },
    }),
  })

  it('for all contract types returns a record of oracle and node addresses and names', () => {
    const expectedResult = {
      oracleAddress1: 'Oracle 1',
      nodeAddress1: 'Oracle 1',
      nodeAddress12: 'Oracle 1',
      oracleAddress2: 'Oracle 2',
      nodeAddress2: 'Oracle 2',
    }

    expect(upcaseOracles(aggregatedState)).toEqual(expectedResult)
  })
})
