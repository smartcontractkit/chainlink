import { parseDot } from './parseDot'

describe('components/Jobs/parseDot', () => {
  it('correctly adds node attributes', () => {
    const digraph1 = `digraph {
      fetch [type=http method=POST url="http://localhost:8001" params="{\\"hi\\": \\"hello\\"}"];
      parse [type=jsonparse path="data,result"];
      multiply [type=multiply times=100];
      fetch -> parse -> multiply;
    }`

    const expected1 = [
      {
        id: 'fetch',
        parentIds: [],
        attributes: {
          type: 'http',
          method: 'POST',
          url: 'http://localhost:8001',
          params: '{"hi": "hello"}',
        },
      },
      {
        id: 'parse',
        parentIds: ['fetch'],
        attributes: { type: 'jsonparse', path: 'data,result' },
      },
      {
        id: 'multiply',
        parentIds: ['parse'],
        attributes: { type: 'multiply', times: '100' },
      },
    ]

    const stratify1 = parseDot(digraph1)
    expect(stratify1).toEqual(expected1)
  })

  it('correctly assigns multiple parentIds', () => {
    const digraph2 = `digraph {
      exercise -> sleep;
      learn -> sleep;
      sleep -> eat;
      eat -> learn;
      eat -> exercise;
    }`

    const expected2 = [
      { id: 'exercise', parentIds: ['eat'] },
      { id: 'sleep', parentIds: ['exercise', 'learn'] },
      { id: 'learn', parentIds: ['eat'] },
      { id: 'eat', parentIds: ['sleep'] },
    ]

    const stratify2 = parseDot(digraph2)
    expect(stratify2).toEqual(expected2)
  })
})
