import { parseDot } from './parseDot'

describe('parseDot', () => {
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

  it('correctly parses on data attributes with newlines', () => {
    const value = `
    ds_fetch        [type=http method=GET url="http://localhost:8001"]
    ds_parse        [type=jsonparse path="data,result"]
    ds_multiply    [type=multiply times=100]
    encode_data [type=ethabiencode abi="(uint256 value)"]
    encode_tx [type=ethabiencode
                 abi="fulfillOracleRequest(bytes32 requestId, uint256 payment, address callbackAddress, bytes4 callbackFunctionId, uint256 expiration, bytes32 data)"
     data=<{"requestId": $(decode_log.requestId), "payment":   $(decode_log.payment), "callbackAddress": $(decode_log.callbackAddr), "callbackFunctionId": $(deqaacode_log.callbackFunctionId), "expiration": $(decode_log.cancelExpiration),"data": $(encode_data)}>
]

encode_mwr [type=ethabiencode
                abi="(bytes32 requestId, uint256 usd, uint256 eur, uint256 jpy)"
                data=<{
                    "requestId": $(decode_log.requestId),
                    "usd": $(usd_multiply),
                    "eur": $(eur_multiply),
                    "jpy": $(jpy_multiply)}>]

    submit_tx    [type=ethtx to="eth_address"]

    ds_fetch -> ds_parse -> ds_multiply -> encode_data -> encode_tx -> encode_mwr -> submit_tx
    `
    const digraph = `digraph {${value}}`

    const expected = [
      {
        id: 'ds_fetch',
        parentIds: [],
        attributes: {
          type: 'http',
          method: 'GET',
          url: 'http://localhost:8001',
        },
      },
      {
        id: 'ds_parse',
        parentIds: ['ds_fetch'],
        attributes: { type: 'jsonparse', path: 'data,result' },
      },
      {
        id: 'ds_multiply',
        parentIds: ['ds_parse'],
        attributes: { type: 'multiply', times: '100' },
      },
      {
        id: 'encode_data',
        parentIds: ['ds_multiply'],
        attributes: { type: 'ethabiencode', abi: '(uint256 value)' },
      },
      {
        id: 'encode_tx',
        parentIds: ['encode_data'],
        attributes: {
          type: 'ethabiencode',
          abi: 'fulfillOracleRequest(bytes32 requestId, uint256 payment, address callbackAddress, bytes4 callbackFunctionId, uint256 expiration, bytes32 data)',
        },
      },
      {
        id: 'encode_mwr',
        parentIds: ['encode_tx'],
        attributes: {
          type: 'ethabiencode',
          abi: '(bytes32 requestId, uint256 usd, uint256 eur, uint256 jpy)',
        },
      },
      {
        id: 'submit_tx',
        parentIds: ['encode_mwr'],
        attributes: {
          type: 'ethtx',
          to: 'eth_address',
        },
      },
    ]

    const stratify = parseDot(digraph)

    expect(stratify).toEqual(expected)
  })
})
