import { partialAsFull } from 'support/test-helpers/partialAsFull'
import {
  Job,
  JobSpecError,
  DirectRequestJob,
  FluxMonitorJob,
  KeeperJob,
  OffChainReportingJob,
  CronJob,
  WebhookJob,
  VRFJob,
} from 'core/store/models'
import { generateUuid } from '../test-helpers/generateUuid'

export function ocrJob(
  config: Partial<
    Job['offChainReportingOracleSpec'] & {
      name?: string
      id?: string
      maxTaskDuration?: string
    } & {
      dotDagSource?: string
    }
  > = {},
): Job {
  const offChainReportingOracleSpec = partialAsFull<
    OffChainReportingJob['offChainReportingOracleSpec']
  >({
    contractAddress: config.contractAddress || generateUuid(),
    p2pPeerID: config.p2pPeerID || generateUuid(),
    p2pBootstrapPeers: config.p2pBootstrapPeers,
    isBootstrapPeer: config.isBootstrapPeer,
    keyBundleID: config.keyBundleID || generateUuid(),
    monitoringEndpoint: config.monitoringEndpoint,
    transmitterAddress: config.transmitterAddress || generateUuid(),
    observationTimeout: config.observationTimeout,
    blockchainTimeout: config.blockchainTimeout,
    contractConfigTrackerSubscribeInterval:
      config.contractConfigTrackerSubscribeInterval,
    contractConfigTrackerPollInterval: config.contractConfigTrackerPollInterval,
    contractConfigConfirmations: config.contractConfigConfirmations,
    updatedAt: config.updatedAt || new Date(1600775300410).toISOString(),
    createdAt: config.createdAt || new Date(1600775300410).toISOString(),
  })
  return {
    name: config.name || 'V2 job',
    type: 'offchainreporting',
    schemaVersion: 1,
    externalJobID: '0EEC7E1D-D0D2-476C-A1A8-72DFB6633F46',
    offChainReportingOracleSpec,
    fluxMonitorSpec: null,
    directRequestSpec: null,
    keeperSpec: null,
    vrfSpec: null,
    cronSpec: null,
    webhookSpec: null,
    errors: [],
    maxTaskDuration: '',
    pipelineSpec: {
      dotDagSource:
        typeof config.dotDagSource === 'string'
          ? config.dotDagSource
          : '   fetch    [type=http method=POST url="http://localhost:8001" requestData="{\\"hi\\": \\"hello\\"}"];\n    parse    [type=jsonparse path="data,result"];\n    multiply [type=multiply times=100];\n    fetch -\u003e parse -\u003e multiply;\n',
    },
  }
}

export function fluxMonitorJobV2(
  spec: Partial<FluxMonitorJob['fluxMonitorSpec']> = {},
  config: Partial<
    {
      name?: string
      id?: string
      maxTaskDuration?: string
      errors: JobSpecError[]
    } & {
      dotDagSource?: string
    }
  > = {},
): Job {
  const fluxMonitorSpec = partialAsFull<FluxMonitorJob['fluxMonitorSpec']>({
    createdAt: spec.createdAt || new Date(1600775300410).toISOString(),
  })

  return {
    name: config.name || 'Flux Monitor V2 job',
    type: 'fluxmonitor',
    schemaVersion: 1,
    externalJobID: '0EEC7E1D-D0D2-476C-A1A8-72DFB6633F47',
    directRequestSpec: null,
    offChainReportingOracleSpec: null,
    keeperSpec: null,
    cronSpec: null,
    webhookSpec: null,
    vrfSpec: null,
    fluxMonitorSpec,
    errors: config.errors || [],
    maxTaskDuration: '',
    pipelineSpec: {
      dotDagSource:
        typeof config.dotDagSource === 'string'
          ? config.dotDagSource
          : '   fetch    [type=http method=POST url="http://localhost:8001" requestData="{\\"hi\\": \\"hello\\"}"];\n    parse    [type=jsonparse path="data,result"];\n    multiply [type=multiply times=100];\n    fetch -\u003e parse -\u003e multiply;\n',
    },
  }
}

export function directRequestJobV2(
  spec: Partial<DirectRequestJob['directRequestSpec']> = {},
  config: Partial<
    {
      name?: string
      id?: string
      maxTaskDuration?: string
    } & {
      dotDagSource?: string
    }
  > = {},
): Job {
  const directRequestSpec = partialAsFull<
    DirectRequestJob['directRequestSpec']
  >({
    createdAt: spec.createdAt || new Date(1600775300410).toISOString(),
  })
  return {
    name: config.name || 'Direct Request V2 job',
    type: 'directrequest',
    schemaVersion: 1,
    externalJobID: '0EEC7E1D-D0D2-476C-A1A8-72DFB6633F49',
    directRequestSpec,
    offChainReportingOracleSpec: null,
    keeperSpec: null,
    cronSpec: null,
    webhookSpec: null,
    fluxMonitorSpec: null,
    vrfSpec: null,
    errors: [],
    maxTaskDuration: '',
    pipelineSpec: {
      dotDagSource:
        typeof config.dotDagSource === 'string'
          ? config.dotDagSource
          : '   fetch    [type=http method=POST url="http://localhost:8001" requestData="{\\"hi\\": \\"hello\\"}"];\n    parse    [type=jsonparse path="data,result"];\n    multiply [type=multiply times=100];\n    fetch -\u003e parse -\u003e multiply;\n',
    },
  }
}

export function keeperJobV2(
  spec: Partial<KeeperJob['keeperSpec']> = {},
  config: Partial<{
    name?: string
    id?: string
  }> = {},
): Job {
  const keeperSpec = partialAsFull<KeeperJob['keeperSpec']>({
    createdAt: spec.createdAt || new Date(1600775300410).toISOString(),
  })
  return {
    name: config.name || 'Keeper V2 job',
    type: 'keeper',
    schemaVersion: 1,
    directRequestSpec: null,
    externalJobID: '0EEC7E1D-D0D2-476C-A1A8-72DFB6633F50',
    keeperSpec,
    offChainReportingOracleSpec: null,
    fluxMonitorSpec: null,
    vrfSpec: null,
    cronSpec: null,
    webhookSpec: null,
    errors: [],
    maxTaskDuration: '',
    pipelineSpec: {
      dotDagSource: '',
    },
  }
}

export function cronJobV2(
  spec: Partial<CronJob['cronSpec']> = {},
  config: Partial<{
    name?: string
    id?: string
  }> = {},
): Job {
  const cronSpec = partialAsFull<CronJob['cronSpec']>({
    createdAt: spec.createdAt || new Date(1600775300410).toISOString(),
  })
  return {
    name: config.name || 'Cron V2 job',
    type: 'cron',
    externalJobID: '0EEC7E1D-D0D2-476C-A1A8-72DFB6633F51',
    schemaVersion: 1,
    directRequestSpec: null,
    keeperSpec: null,
    vrfSpec: null,
    offChainReportingOracleSpec: null,
    fluxMonitorSpec: null,
    cronSpec,
    webhookSpec: null,
    errors: [],
    maxTaskDuration: '',
    pipelineSpec: {
      dotDagSource: '',
    },
  }
}

export function webhookJobV2(
  spec: Partial<WebhookJob['webhookSpec']> = {},
  config: Partial<{
    name?: string
    id?: string
  }> = {},
): Job {
  const webhookSpec = partialAsFull<WebhookJob['webhookSpec']>({
    createdAt: spec.createdAt || new Date(1600775300410).toISOString(),
  })
  return {
    name: config.name || 'Web V2 job',
    type: 'webhook',
    externalJobID: '0EEC7E1D-D0D2-476C-A1A8-72DFB6633F52',
    schemaVersion: 1,
    directRequestSpec: null,
    vrfSpec: null,
    keeperSpec: null,
    offChainReportingOracleSpec: null,
    fluxMonitorSpec: null,
    cronSpec: null,
    webhookSpec,
    errors: [],
    maxTaskDuration: '',
    pipelineSpec: {
      dotDagSource: '',
    },
  }
}

export function vrfJobV2(
  spec: Partial<VRFJob['vrfSpec']> = {},
  config: Partial<{
    name?: string
    id?: string
  }> = {},
): Job {
  const vrfSpec = partialAsFull<VRFJob['vrfSpec']>({
    createdAt: spec.createdAt || new Date(1600775300410).toISOString(),
  })
  return {
    name: config.name || 'VRF V2 job',
    type: 'vrf',
    externalJobID: '0EEC7E1D-D0D2-476C-A1A8-72DFB6633F52',
    schemaVersion: 1,
    directRequestSpec: null,
    keeperSpec: null,
    offChainReportingOracleSpec: null,
    fluxMonitorSpec: null,
    cronSpec: null,
    webhookSpec: null,
    vrfSpec,
    errors: [],
    maxTaskDuration: '',
    pipelineSpec: {
      dotDagSource: '',
    },
  }
}
