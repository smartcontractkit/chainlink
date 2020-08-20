import { upcaseOracles } from './selectors'

describe('upcaseOracles', () => {
  it('for v2 contract returns a record of oracle addresses and names', () => {
    const aggregatedState = {
      aggregator: {
        oracleNodes: {
          oracleAddress1: {
            name: 'Oracle 1',
            nodeAddress: ['nodeAddress1'],
            oracleAddress: 'oracleAddress1',
          },
          oracleAddress2: {
            name: 'Oracle 2',
            nodeAddress: ['nodeAddress2'],
            oracleAddress: 'oracleAddress2',
          },
        },
        config: {
          contractVersion: 2,
        },
      },
    }

    const expectedResult = {
      oracleAddress1: 'Oracle 1',
      oracleAddress2: 'Oracle 2',
    }

    expect(upcaseOracles(aggregatedState)).toEqual(expectedResult)
  })

  it('for v3 contract returns a record of node addresses and names', () => {
    const aggregatedState = {
      aggregator: {
        oracleNodes: {
          nodeAddress1: {
            name: 'Oracle 1',
            nodeAddress: ['nodeAddress1'],
            oracleAddress: 'oracleAddress1',
          },
          nodeAddress2: {
            name: 'Oracle 2',
            nodeAddress: ['nodeAddress2'],
            oracleAddress: 'oracleAddress2',
          },
        },
        config: {
          contractVersion: 3,
        },
      },
    }

    const expectedResult = {
      nodeAddress1: 'Oracle 1',
      nodeAddress2: 'Oracle 2',
    }

    expect(upcaseOracles(aggregatedState)).toEqual(expectedResult)
  })
})
