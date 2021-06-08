/* eslint-enable no-useless-escape */

import {
  InitiatorType,
  JobSpecV2,
  OffChainReportingOracleJobV2Spec,
} from 'core/store/models'
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
  it('generates a valid OCR definition', () => {
    const jobSpecAttributesInput: OffChainReportingOracleJobV2Spec = {
      name: 'Job spec v2',
      type: 'offchainreporting',
      fluxMonitorSpec: null,
      directRequestSpec: null,
      keeperSpec: null,
      cronSpec: null,
      webhookSpec: null,
      schemaVersion: 1,
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
      maxTaskDuration: '10s',
      pipelineSpec: {
        dotDagSource:
          '    fetch    [type=http method=POST url="http://localhost:8001" requestData="{\\"hi\\": \\"hello\\"}"];\n    parse    [type=jsonparse path="data,result"];\n    multiply [type=multiply times=100];\n    fetch -> parse -> multiply;\n',
      },
      errors: [],
    }

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
maxTaskDuration = "10s"
`

    const output = generateTOMLDefinition(jobSpecAttributesInput)
    expect(output).toEqual(expectedOutput)
  })

  it('generates a valid Flux Monitor definition', () => {
    const jobSpecAttributesInput = {
      name: 'FM Job Spec',
      schemaVersion: 1,
      type: 'fluxmonitor',
      fluxMonitorSpec: {
        absoluteThreshold: 1,
        contractAddress: '0x3cCad4715152693fE3BC4460591e3D3Fbd071b42',
        createdAt: '2021-02-19T16:00:01.115227+08:00',
        idleTimerDisabled: false,
        idleTimerPeriod: '1s',
        pollTimerDisabled: false,
        pollTimerPeriod: '1m0s',
        precision: 2,
        threshold: 0.5,
        updatedAt: '2021-02-19T16:00:01.115227+08:00',
        minPayment: null,
      },
      keeperSpec: null,
      cronSpec: null,
      webhookSpec: null,
      directRequestSpec: null,
      offChainReportingOracleSpec: null,
      maxTaskDuration: '10s',
      pipelineSpec: {
        dotDagSource:
          '    fetch    [type=http method=POST url="http://localhost:8001" requestData="{\\"hi\\": \\"hello\\"}"];\n    parse    [type=jsonparse path="data,result"];\n    multiply [type=multiply times=100];\n    fetch -> parse -> multiply;\n',
      },
      errors: [],
    } as JobSpecV2

    const expectedOutput = `type = "fluxmonitor"
schemaVersion = 1
name = "FM Job Spec"
contractAddress = "0x3cCad4715152693fE3BC4460591e3D3Fbd071b42"
precision = 2
threshold = 0.5
absoluteThreshold = 1
idleTimerPeriod = "1s"
idleTimerDisabled = false
pollTimerPeriod = "1m0s"
pollTimerDisabled = false
maxTaskDuration = "10s"
observationSource = """
    fetch    [type=http method=POST url="http://localhost:8001" requestData="{\\\\"hi\\\\": \\\\"hello\\\\"}"];
    parse    [type=jsonparse path="data,result"];
    multiply [type=multiply times=100];
    fetch -> parse -> multiply;
"""
`

    const output = generateTOMLDefinition(jobSpecAttributesInput)
    expect(output).toEqual(expectedOutput)
  })

  it('generates a valid Direct Request definition', () => {
    const jobSpecAttributesInput = {
      name: 'DR Job Spec',
      schemaVersion: 1,
      type: 'directrequest',
      fluxMonitorSpec: null,
      keeperSpec: null,
      cronSpec: null,
      webhookSpec: null,
      directRequestSpec: {
        initiator: 'runlog',
        contractAddress: '0x3cCad4715152693fE3BC4460591e3D3Fbd071b42',
        minIncomingConfirmations: 3,
        onChainJobSpecID: '0eec7e1dd0d2476ca1a872dfb6633f46',
        createdAt: '2021-02-19T16:00:01.115227+08:00',
      },
      offChainReportingOracleSpec: null,
      maxTaskDuration: '10s',
      pipelineSpec: {
        dotDagSource:
          '    fetch    [type=http method=POST url="http://localhost:8001" requestData="{\\"hi\\": \\"hello\\"}"];\n    parse    [type=jsonparse path="data,result"];\n    multiply [type=multiply times=100];\n    fetch -> parse -> multiply;\n',
      },
      errors: [],
    } as JobSpecV2

    const expectedOutput = `type = "directrequest"
schemaVersion = 1
name = "DR Job Spec"
onChainJobSpecID = "0eec7e1dd0d2476ca1a872dfb6633f46"
minIncomingConfirmations = 3
contractAddress = "0x3cCad4715152693fE3BC4460591e3D3Fbd071b42"
maxTaskDuration = "10s"
observationSource = """
    fetch    [type=http method=POST url="http://localhost:8001" requestData="{\\\\"hi\\\\": \\\\"hello\\\\"}"];
    parse    [type=jsonparse path="data,result"];
    multiply [type=multiply times=100];
    fetch -> parse -> multiply;
"""
`

    const output = generateTOMLDefinition(jobSpecAttributesInput)
    expect(output).toEqual(expectedOutput)
  })

  it('generates a valid Keeper definition', () => {
    const jobSpecAttributesInput = {
      name: 'Keeper Job Spec',
      schemaVersion: 1,
      type: 'keeper',
      fluxMonitorSpec: null,
      keeperSpec: {
        contractAddress: '0x9E40733cC9df84636505f4e6Db28DCa0dC5D1bba',
        createdAt: '2021-04-05T15:21:30.392021+08:00',
        fromAddress: '0xa8037A20989AFcBC51798de9762b351D63ff462e',
        updatedAt: '2021-04-05T15:21:30.392021+08:00',
      },
      cronSpec: null,
      webhookSpec: null,
      directRequestSpec: null,
      offChainReportingOracleSpec: null,
      maxTaskDuration: '10s',
      pipelineSpec: {
        id: '1',
        dotDagSource: '',
      },
      errors: [],
    } as JobSpecV2

    const expectedOutput = `type = "keeper"
schemaVersion = 1
name = "Keeper Job Spec"
contractAddress = "0x9E40733cC9df84636505f4e6Db28DCa0dC5D1bba"
fromAddress = "0xa8037A20989AFcBC51798de9762b351D63ff462e"
`

    const output = generateTOMLDefinition(jobSpecAttributesInput)
    expect(output).toEqual(expectedOutput)
  })

  it('generates a valid Cron definition', () => {
    const jobSpecAttributesInput = {
      name: 'Cron Job Spec',
      schemaVersion: 1,
      type: 'cron',
      fluxMonitorSpec: null,
      keeperSpec: null,
      cronSpec: {
        schedule: '*/2 * * * *',
        createdAt: '2021-04-05T15:21:30.392021+08:00',
        updatedAt: '2021-04-05T15:21:30.392021+08:00',
      },
      webhookSpec: null,
      directRequestSpec: null,
      offChainReportingOracleSpec: null,
      maxTaskDuration: '10s',
      pipelineSpec: {
        dotDagSource:
          '    ds    [type=http method=GET url="http://localhost:8001"];\n    ds_parse    [type=jsonparse path="data,result"];\n    ds_multiply [type=multiply times=100];\n    ds -> ds_parse -> ds_multiply;\n',
      },
      errors: [],
    } as JobSpecV2

    const expectedOutput = `type = "cron"
schemaVersion = 1
name = "Cron Job Spec"
schedule = "*/2 * * * *"
observationSource = """
    ds    [type=http method=GET url="http://localhost:8001"];
    ds_parse    [type=jsonparse path="data,result"];
    ds_multiply [type=multiply times=100];
    ds -> ds_parse -> ds_multiply;
"""
`

    const output = generateTOMLDefinition(jobSpecAttributesInput)
    expect(output).toEqual(expectedOutput)
  })

  it('generates a valid Webhook definition', () => {
    const jobSpecAttributesInput = {
      name: 'Webhook Job Spec',
      schemaVersion: 1,
      type: 'webhook',
      fluxMonitorSpec: null,
      keeperSpec: null,
      cronSpec: null,
      webhookSpec: {
        onChainJobSpecID: '0eec7e1dd0d2476ca1a872dfb6633f46',
        createdAt: '2021-04-05T15:21:30.392021+08:00',
        updatedAt: '2021-04-05T15:21:30.392021+08:00',
      },
      directRequestSpec: null,
      offChainReportingOracleSpec: null,
      maxTaskDuration: '10s',
      pipelineSpec: {
        dotDagSource:
          '    ds    [type=http method=GET url="http://localhost:8001"];\n    ds_parse    [type=jsonparse path="data,result"];\n    ds_multiply [type=multiply times=100];\n    ds -> ds_parse -> ds_multiply;\n',
      },
      errors: [],
    } as JobSpecV2

    const expectedOutput = `type = "webhook"
schemaVersion = 1
name = "Webhook Job Spec"
onChainJobSpecID = "0eec7e1dd0d2476ca1a872dfb6633f46"
observationSource = """
    ds    [type=http method=GET url="http://localhost:8001"];
    ds_parse    [type=jsonparse path="data,result"];
    ds_multiply [type=multiply times=100];
    ds -> ds_parse -> ds_multiply;
"""
`
    const output = generateTOMLDefinition(jobSpecAttributesInput)
    expect(output).toEqual(expectedOutput)
  })
})
