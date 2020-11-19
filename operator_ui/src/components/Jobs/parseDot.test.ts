import { parseDot } from './parseDot'

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

const digraph2 = `digraph {
  exercise -> sleep;
  learn -> sleep;
  sleep -> eat;
  eat -> learn;
  eat -> exercise;
}`

const expected2 = [
  { id: 'exercise', parentIds: ['eat'], attributes: {} },
  { id: 'sleep', parentIds: ['exercise', 'learn'], attributes: {} },
  { id: 'learn', parentIds: ['eat'], attributes: {} },
  { id: 'eat', parentIds: ['sleep'], attributes: {} },
]

describe('components/Jobs/parseDot', () => {
  it('return stratify object', () => {
    const stratify1 = parseDot(digraph1)
    expect(stratify1).toEqual(expected1)

    const stratify2 = parseDot(digraph2)
    expect(stratify2).toEqual(expected2)
  })
})
