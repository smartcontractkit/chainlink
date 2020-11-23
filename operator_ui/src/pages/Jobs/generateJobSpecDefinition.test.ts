import { InitiatorType } from 'core/store/models'
import {
  generateJSONDefinition,
  generateTOMLDefinition,
} from './generateJobSpecDefinition'

describe('generateJSONDefinition', () => {
  it('generates valid definition', () => {
    const jobSpecAttributesInput = {
      initiators: [
        {
          type: 'web' as InitiatorType.WEB,
        },
      ],
      id: '7f4e4ca5f9ce4131a080a214947736c5',
      name: 'Bitstamp ticker',
      createdAt: '2020-11-17T10:25:44.040459Z',
      tasks: [
        {
          ID: 6,
          type: 'httpget',
          confirmations: 0,
          params: {
            get: 'https://bitstamp.net/api/ticker/',
          },
          CreatedAt: '2020-11-17T10:25:44.043094Z',
          UpdatedAt: '2020-11-17T10:25:44.043094Z',
          DeletedAt: null,
        },
        {
          ID: 7,
          type: 'jsonparse',
          confirmations: null,
          params: {
            path: ['last'],
          },
          CreatedAt: '2020-11-17T10:25:44.043948Z',
          UpdatedAt: '2020-11-17T10:25:44.043948Z',
          DeletedAt: null,
        },
        {
          ID: 8,
          type: 'multiply',
          confirmations: null,
          params: {
            times: 100,
          },
          CreatedAt: '2020-11-17T10:25:44.04456Z',
          UpdatedAt: '2020-11-17T10:25:44.04456Z',
          DeletedAt: null,
        },
        {
          ID: 9,
          type: 'ethuint256',
          confirmations: null,
          params: {},
          CreatedAt: '2020-11-17T10:25:44.045404Z',
          UpdatedAt: '2020-11-17T10:25:44.045404Z',
          DeletedAt: null,
        },
        {
          ID: 10,
          type: 'ethtx',
          confirmations: null,
          params: {},
          CreatedAt: '2020-11-17T10:25:44.046211Z',
          UpdatedAt: '2020-11-17T10:25:44.046211Z',
          DeletedAt: null,
        },
      ],
      minPayment: '1000000',
      updatedAt: '2020-02-09T15:13:03Z',
      startAt: '2020-02-09T15:13:03Z',
      endAt: null,
      errors: [],
      earnings: null,
    }

    const expectedOutput = `{
    "name": "Bitstamp ticker",
    "initiators": [
        {
            "type": "web"
        }
    ],
    "tasks": [
        {
            "type": "httpget",
            "confirmations": 0,
            "params": {
                "get": "https://bitstamp.net/api/ticker/"
            }
        },
        {
            "type": "jsonparse",
            "params": {
                "path": [
                    "last"
                ]
            }
        },
        {
            "type": "multiply",
            "params": {
                "times": 100
            }
        },
        {
            "type": "ethuint256"
        },
        {
            "type": "ethtx"
        }
    ],
    "startAt": "2020-02-09T15:13:03Z"
}`

    const output = generateJSONDefinition(jobSpecAttributesInput)
    expect(output).toEqual(expectedOutput)
  })

  it('removes the name if it is auto-generated (has a job ID in it)', () => {
    const jobSpecAttributesInput = {
      initiators: [
        {
          type: 'web' as InitiatorType.WEB,
        },
      ],
      id: '7f4e4ca5f9ce4131a080a214947736c5',
      name: 'Job7f4e4ca5f9ce4131a080a214947736c5',
      createdAt: '2020-11-17T10:25:44.040459Z',
      tasks: [
        {
          ID: 6,
          type: 'httpget',
          confirmations: 0,
          params: {
            get: 'https://bitstamp.net/api/ticker/',
          },
          CreatedAt: '2020-11-17T10:25:44.043094Z',
          UpdatedAt: '2020-11-17T10:25:44.043094Z',
          DeletedAt: null,
        },
      ],
      minPayment: '1000000',
      updatedAt: '2020-02-09T15:13:03Z',
      startAt: '2020-02-09T15:13:03Z',
      endAt: null,
      errors: [],
      earnings: null,
    }

    const expectedOutput = `{
    "initiators": [
        {
            "type": "web"
        }
    ],
    "tasks": [
        {
            "type": "httpget",
            "confirmations": 0,
            "params": {
                "get": "https://bitstamp.net/api/ticker/"
            }
        }
    ],
    "startAt": "2020-02-09T15:13:03Z"
}`

    const output = generateJSONDefinition(jobSpecAttributesInput)
    expect(output).toEqual(expectedOutput)
  })
})

describe('generateTOMLDefinition', () => {
  it('generates valid definition', () => {
    const jobSpecAttributesInput = {
      offChainReportingOracleSpec: {
        contractAddress: '0x1469877c88F19E273EFC7Ef3C9D944574583B8a0',
        p2pPeerID: '12D3KooWL4zx7Tu92wNuK14LT2BV4mXxNoNK3zuxE7iKNgiazJFm',
        p2pBootstrapPeers: [
          '/ip4/139.59.41.32/tcp/12000/p2p/12D3KooWGKhStcrvCr5RBYKaSRNX4ojrxHcmpJuFmHWenT6aAQAY',
        ],
        isBootstrapPeer: false,
        keyBundleID:
          '4ee612467c3caea7bdab57ab62937adfc4d195516c30139a737f85098b35d9af',
        monitoringEndpoint: 'chain.link:4321',
        transmitterAddress: '0x01010CaB43e77116c95745D219af1069fE050d7A',
        observationTimeout: '10s',
        blockchainTimeout: '20s',
        contractConfigTrackerSubscribeInterval: '2m0s',
        contractConfigTrackerPollInterval: '1m0s',
        contractConfigConfirmations: 3,
        createdAt: '2020-11-17T13:50:13.182669Z',
        updatedAt: '2020-11-17T13:50:13.182669Z',
      },
      pipelineSpec: {
        dotDagSource:
          '    fetch    [type=http method=POST url="http://localhost:8001" requestData="{\\"hi\\": \\"hello\\"}"];\n    parse    [type=jsonparse path="data,result"];\n    multiply [type=multiply times=100];\n    fetch -> parse -> multiply;\n',
      },
      errors: [],
    }

    /* eslint-disable no-useless-escape */
    const expectedOutput = `type = "offchainreporting"
schemaVersion = 1
contractAddress = "0x1469877c88F19E273EFC7Ef3C9D944574583B8a0"
p2pPeerID = "12D3KooWL4zx7Tu92wNuK14LT2BV4mXxNoNK3zuxE7iKNgiazJFm"
p2pBootstrapPeers = [
  "/ip4/139.59.41.32/tcp/12000/p2p/12D3KooWGKhStcrvCr5RBYKaSRNX4ojrxHcmpJuFmHWenT6aAQAY"
]
isBootstrapPeer = false
keyBundleID = "4ee612467c3caea7bdab57ab62937adfc4d195516c30139a737f85098b35d9af"
monitoringEndpoint = "chain.link:4321"
transmitterAddress = "0x01010CaB43e77116c95745D219af1069fE050d7A"
observationTimeout = "10s"
blockchainTimeout = "20s"
contractConfigTrackerSubscribeInterval = "2m0s"
contractConfigTrackerPollInterval = "1m0s"
contractConfigConfirmations = 3
observationSource = """
    fetch    [type=http method=POST url="http://localhost:8001" requestData="{\\\\"hi\\\\": \\\\"hello\\\\"}"];
    parse    [type=jsonparse path="data,result"];
    multiply [type=multiply times=100];
    fetch -> parse -> multiply;
"""
`
    /* eslint-enable no-useless-escape */

    const output = generateTOMLDefinition(jobSpecAttributesInput)
    expect(output).toEqual(expectedOutput)
  })
})
