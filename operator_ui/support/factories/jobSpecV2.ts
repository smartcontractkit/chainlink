import { partialAsFull } from 'support/test-helpers/partialAsFull'
import {
  JobSpecV2,
  DirectRequestJobV2Spec,
  FluxMonitorJobV2Spec,
  KeeperV2Spec,
  OffChainReportingOracleJobV2Spec,
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
    offChainReportingOracleSpec,
    fluxMonitorSpec: null,
    directRequestSpec: null,
    keeperSpec: null,
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
    directRequestSpec: null,
    offChainReportingOracleSpec: null,
    keeperSpec: null,
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
    directRequestSpec,
    offChainReportingOracleSpec: null,
    keeperSpec: null,
    fluxMonitorSpec: null,
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
    keeperSpec,
    offChainReportingOracleSpec: null,
    fluxMonitorSpec: null,
    errors: [],
    maxTaskDuration: '',
    pipelineSpec: {
      dotDagSource: '',
    },
  }
}
