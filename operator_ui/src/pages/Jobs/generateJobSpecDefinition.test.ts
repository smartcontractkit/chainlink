/* eslint-enable no-useless-escape */

import { Job, OffChainReportingJob } from 'core/store/models'
import { generateTOMLDefinition } from './generateJobSpecDefinition'

describe('generateTOMLDefinition', () => {
  it('generates a valid OCR definition', () => {
    const jobSpecAttributesInput: OffChainReportingJob = {
      name: 'Job spec v2',
      type: 'offchainreporting',
      fluxMonitorSpec: null,
      externalJobID: '0eec7e1d-d0d2-476c-a1a8-72dfb6633f46',
      directRequestSpec: null,
      keeperSpec: null,
      cronSpec: null,
      webhookSpec: null,
      schemaVersion: 1,
      vrfSpec: null,
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
externalJobID = "0eec7e1d-d0d2-476c-a1a8-72dfb6633f46"
`

    const output = generateTOMLDefinition(jobSpecAttributesInput)
    expect(output.definition).toEqual(expectedOutput)
    expect(output.envAttributesDefinition).toBe('')
  })

  it('generates a valid Flux Monitor definition', () => {
    const jobSpecAttributesInput = {
      name: 'FM Job Spec',
      schemaVersion: 1,
      type: 'fluxmonitor',
      externalJobID: '0eec7e1d-d0d2-476c-a1a8-72dfb6633f46',
      fluxMonitorSpec: {
        absoluteThreshold: 1,
        contractAddress: '0x3cCad4715152693fE3BC4460591e3D3Fbd071b42',
        createdAt: '2021-02-19T16:00:01.115227+08:00',
        idleTimerDisabled: false,
        idleTimerPeriod: '1s',
        pollTimerDisabled: false,
        pollTimerPeriod: '1m0s',
        drumbeatEnabled: true,
        drumbeatSchedule: '@every 10m',
        drumbeatRandomDelay: '10s',
        precision: 2,
        threshold: 0.5,
        updatedAt: '2021-02-19T16:00:01.115227+08:00',
        minPayment: null,
      },
      vrfSpec: null,
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
    } as Job

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
drumbeatEnabled = true
drumbeatSchedule = "@every 10m"
drumbeatRandomDelay = "10s"
maxTaskDuration = "10s"
observationSource = """
    fetch    [type=http method=POST url="http://localhost:8001" requestData="{\\\\"hi\\\\": \\\\"hello\\\\"}"];
    parse    [type=jsonparse path="data,result"];
    multiply [type=multiply times=100];
    fetch -> parse -> multiply;
"""
externalJobID = "0eec7e1d-d0d2-476c-a1a8-72dfb6633f46"
`

    const output = generateTOMLDefinition(jobSpecAttributesInput)
    expect(output.definition).toEqual(expectedOutput)
    expect(output.envAttributesDefinition).toBe('')
  })

  it('generates a valid Direct Request definition', () => {
    const jobSpecAttributesInput = {
      name: 'DR Job Spec',
      schemaVersion: 1,
      type: 'directrequest',
      externalJobID: '0eec7e1d-d0d2-476c-a1a8-72dfb6633f46',
      fluxMonitorSpec: null,
      keeperSpec: null,
      cronSpec: null,
      vrfSpec: null,
      webhookSpec: null,
      directRequestSpec: {
        initiator: 'runlog',
        contractAddress: '0x3cCad4715152693fE3BC4460591e3D3Fbd071b42',
        minIncomingConfirmations: 3,
        createdAt: '2021-02-19T16:00:01.115227+08:00',
        requesters: ['0x59bbE8CFC79c76857fE0eC27e67E4957370d72B5'],
      },
      offChainReportingOracleSpec: null,
      maxTaskDuration: '10s',
      pipelineSpec: {
        dotDagSource:
          '    fetch    [type=http method=POST url="http://localhost:8001" requestData="{\\"hi\\": \\"hello\\"}"];\n    parse    [type=jsonparse path="data,result"];\n    multiply [type=multiply times=100];\n    fetch -> parse -> multiply;\n',
      },
      errors: [],
    } as Job

    const expectedOutput = `type = "directrequest"
schemaVersion = 1
name = "DR Job Spec"
minIncomingConfirmations = 3
contractAddress = "0x3cCad4715152693fE3BC4460591e3D3Fbd071b42"
requesters = [ "0x59bbE8CFC79c76857fE0eC27e67E4957370d72B5" ]
maxTaskDuration = "10s"
observationSource = """
    fetch    [type=http method=POST url="http://localhost:8001" requestData="{\\\\"hi\\\\": \\\\"hello\\\\"}"];
    parse    [type=jsonparse path="data,result"];
    multiply [type=multiply times=100];
    fetch -> parse -> multiply;
"""
externalJobID = "0eec7e1d-d0d2-476c-a1a8-72dfb6633f46"
`

    const output = generateTOMLDefinition(jobSpecAttributesInput)
    expect(output.definition).toEqual(expectedOutput)
    expect(output.envAttributesDefinition).toBe('')
  })

  it('generates a valid Keeper definition', () => {
    const jobSpecAttributesInput = {
      name: 'Keeper Job Spec',
      schemaVersion: 1,
      type: 'keeper',
      externalJobID: '0eec7e1d-d0d2-476c-a1a8-72dfb6633f46',
      fluxMonitorSpec: null,
      keeperSpec: {
        contractAddress: '0x9E40733cC9df84636505f4e6Db28DCa0dC5D1bba',
        createdAt: '2021-04-05T15:21:30.392021+08:00',
        fromAddress: '0xa8037A20989AFcBC51798de9762b351D63ff462e',
        updatedAt: '2021-04-05T15:21:30.392021+08:00',
      },
      cronSpec: null,
      vrfSpec: null,
      webhookSpec: null,
      directRequestSpec: null,
      offChainReportingOracleSpec: null,
      maxTaskDuration: '10s',
      pipelineSpec: {
        id: '1',
        dotDagSource: '',
      },
      errors: [],
    } as Job

    const expectedOutput = `type = "keeper"
schemaVersion = 1
name = "Keeper Job Spec"
contractAddress = "0x9E40733cC9df84636505f4e6Db28DCa0dC5D1bba"
fromAddress = "0xa8037A20989AFcBC51798de9762b351D63ff462e"
externalJobID = "0eec7e1d-d0d2-476c-a1a8-72dfb6633f46"
`

    const output = generateTOMLDefinition(jobSpecAttributesInput)
    expect(output.definition).toEqual(expectedOutput)
    expect(output.envAttributesDefinition).toBe('')
  })

  it('generates a valid Cron definition', () => {
    const jobSpecAttributesInput = {
      name: 'Cron Job Spec',
      schemaVersion: 1,
      type: 'cron',
      fluxMonitorSpec: null,
      externalJobID: '0eec7e1d-d0d2-476c-a1a8-72dfb6633f46',
      keeperSpec: null,
      vrfSpec: null,
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
    } as Job

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
externalJobID = "0eec7e1d-d0d2-476c-a1a8-72dfb6633f46"
`

    const output = generateTOMLDefinition(jobSpecAttributesInput)
    expect(output.definition).toEqual(expectedOutput)
    expect(output.envAttributesDefinition).toBe('')
  })

  it('generates a valid Webhook definition', () => {
    const jobSpecAttributesInput = {
      name: 'Webhook Job Spec',
      schemaVersion: 1,
      type: 'webhook',
      externalJobID: '0eec7e1d-d0d2-476c-a1a8-72dfb6633f46',
      fluxMonitorSpec: null,
      keeperSpec: null,
      vrfSpec: null,
      cronSpec: null,
      webhookSpec: {
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
    } as Job

    const expectedOutput = `type = "webhook"
schemaVersion = 1
name = "Webhook Job Spec"
externalJobID = "0eec7e1d-d0d2-476c-a1a8-72dfb6633f46"
observationSource = """
    ds    [type=http method=GET url="http://localhost:8001"];
    ds_parse    [type=jsonparse path="data,result"];
    ds_multiply [type=multiply times=100];
    ds -> ds_parse -> ds_multiply;
"""
`
    const output = generateTOMLDefinition(jobSpecAttributesInput)
    expect(output.definition).toEqual(expectedOutput)
    expect(output.envAttributesDefinition).toBe('')
  })

  it('generates a valid vrf definition', () => {
    const jobSpecAttributesInput = {
      name: 'VRF Job Spec',
      schemaVersion: 1,
      type: 'vrf',
      externalJobID: '0eec7e1d-d0d2-476c-a1a8-72dfb6633f46',
      fluxMonitorSpec: null,
      keeperSpec: null,
      cronSpec: null,
      webhookSpec: null,
      vrfSpec: {
        confirmations: 6,
        coordinatorAddress: '0xABA5eDc1a551E55b1A570c0e1f1055e5BE11eca7',
        publicKey:
          '0x92594ee04c179eb7d439ff1baacd98b81a7d7a6ed55c86ca428fa025bd9c914301',
        fromAddress: '',
        pollPeriod: '',
        createdAt: '2021-04-05T15:21:30.392021+08:00',
        updatedAt: '2021-04-05T15:21:30.392021+08:00',
      },
      directRequestSpec: null,
      pipelineSpec: {
        id: '1',
        dotDagSource: '',
      },
      offChainReportingOracleSpec: null,
      maxTaskDuration: '10s',
      errors: [],
    } as Job

    const expectedOutput = `type = "vrf"
schemaVersion = 1
name = "VRF Job Spec"
externalJobID = "0eec7e1d-d0d2-476c-a1a8-72dfb6633f46"
coordinatorAddress = "0xABA5eDc1a551E55b1A570c0e1f1055e5BE11eca7"
confirmations = 6
publicKey = "0x92594ee04c179eb7d439ff1baacd98b81a7d7a6ed55c86ca428fa025bd9c914301"
fromAddress = ""
pollPeriod = ""
observationSource = ""
`
    const output = generateTOMLDefinition(jobSpecAttributesInput)
    expect(output.definition).toEqual(expectedOutput)
    expect(output.envAttributesDefinition).toBe('')
  })

  it('generates a valid OCR definition with values set by environment vars', () => {
    // const jobSpecAttributesInput: OffChainReportingJob = {
    const jobSpecAttributesInput: any = {
      name: 'Job spec v2 with env vars',
      type: 'offchainreporting',
      fluxMonitorSpec: null,
      externalJobID: '0eec7e1d-d0d2-476c-a1a8-72dfb6633f46',
      directRequestSpec: null,
      keeperSpec: null,
      cronSpec: null,
      webhookSpec: null,
      schemaVersion: 1,
      vrfSpec: null,
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
        observationTimeoutEnv: true,
        blockchainTimeout: '20s',
        blockchainTimeoutEnv: true,
        contractConfigTrackerPollInterval: '1m0s',
        contractConfigTrackerPollIntervalEnv: true,
        contractConfigTrackerSubscribeInterval: '2m0s',
        contractConfigTrackerSubscribeIntervalEnv: true,
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
contractConfigConfirmations = 3
observationSource = """
    fetch    [type=http method=POST url="http://localhost:8001" requestData="{\\\\"hi\\\\": \\\\"hello\\\\"}"];
    parse    [type=jsonparse path="data,result"];
    multiply [type=multiply times=100];
    fetch -> parse -> multiply;
"""
maxTaskDuration = "10s"
externalJobID = "0eec7e1d-d0d2-476c-a1a8-72dfb6633f46"
`
    const expectedEnvOutput = `observationTimeout = "10s"
blockchainTimeout = "20s"
contractConfigTrackerPollInterval = "1m0s"
contractConfigTrackerSubscribeInterval = "2m0s"
`

    const output = generateTOMLDefinition(jobSpecAttributesInput)
    expect(output.definition).toEqual(expectedOutput)
    expect(output.envAttributesDefinition).toBe(expectedEnvOutput)
  })
})
