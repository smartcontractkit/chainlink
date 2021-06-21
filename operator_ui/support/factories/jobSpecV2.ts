import { partialAsFull } from 'support/test-helpers/partialAsFull'
import {
  JobSpecV2,
  DirectRequestJobV2Spec,
  FluxMonitorJobV2Spec,
  KeeperV2Spec,
  OffChainReportingOracleJobV2Spec,
  CronV2Spec,
  WebhookV2Spec,
  VRFV2Spec,
} from 'core/store/models'
import { generateUuid } from '../test-helpers/generateUuid'

export function ocrJobSpecV2(
  config: Partial<
    JobSpecV2['offChainReportingOracleSpec'] & {
      name?: string
      id?: string
      maxTaskDuration?: string
    } & {
      dotDagSource?: string
    }
  > = {},
): JobSpecV2 {
  const offChainReportingOracleSpec = partialAsFull<
    OffChainReportingOracleJobV2Spec['offChainReportingOracleSpec']
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
  spec: Partial<FluxMonitorJobV2Spec['fluxMonitorSpec']> = {},
  config: Partial<
    {
      name?: string
      id?: string
      maxTaskDuration?: string
    } & {
      dotDagSource?: string
    }
  > = {},
): JobSpecV2 {
  const fluxMonitorSpec = partialAsFull<
    FluxMonitorJobV2Spec['fluxMonitorSpec']
  >({
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

export function directRequestJobV2(
  spec: Partial<DirectRequestJobV2Spec['directRequestSpec']> = {},
  config: Partial<
    {
      name?: string
      id?: string
      maxTaskDuration?: string
    } & {
      dotDagSource?: string
    }
  > = {},
): JobSpecV2 {
  const directRequestSpec = partialAsFull<
    DirectRequestJobV2Spec['directRequestSpec']
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
  spec: Partial<KeeperV2Spec['keeperSpec']> = {},
  config: Partial<{
    name?: string
    id?: string
  }> = {},
): JobSpecV2 {
  const keeperSpec = partialAsFull<KeeperV2Spec['keeperSpec']>({
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
  spec: Partial<CronV2Spec['cronSpec']> = {},
  config: Partial<{
    name?: string
    id?: string
  }> = {},
): JobSpecV2 {
  const cronSpec = partialAsFull<CronV2Spec['cronSpec']>({
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
  spec: Partial<WebhookV2Spec['webhookSpec']> = {},
  config: Partial<{
    name?: string
    id?: string
  }> = {},
): JobSpecV2 {
  const webhookSpec = partialAsFull<WebhookV2Spec['webhookSpec']>({
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
  spec: Partial<VRFV2Spec['vrfSpec']> = {},
  config: Partial<{
    name?: string
    id?: string
  }> = {},
): JobSpecV2 {
  const vrfSpec = partialAsFull<VRFV2Spec['vrfSpec']>({
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
