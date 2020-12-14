import { partialAsFull } from '@chainlink/ts-helpers'
import { OcrJobSpec } from 'core/store/models'
import { generateUuid } from '../test-helpers/generateUuid'

export function jobSpecV2(
  config: Partial<
    OcrJobSpec['offChainReportingOracleSpec'] & {
      name?: string
      id?: string
    } & {
      dotDagSource?: string
    }
  > = {},
): OcrJobSpec {
  const offChainReportingOracleSpec = partialAsFull<
    OcrJobSpec['offChainReportingOracleSpec']
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
    offChainReportingOracleSpec,
    errors: [],
    pipelineSpec: {
      dotDagSource:
        typeof config.dotDagSource === 'string'
          ? config.dotDagSource
          : '   fetch    [type=http method=POST url="http://localhost:8001" requestData="{\\"hi\\": \\"hello\\"}"];\n    parse    [type=jsonparse path="data,result"];\n    multiply [type=multiply times=100];\n    fetch -\u003e parse -\u003e multiply;\n',
    },
  }
}
