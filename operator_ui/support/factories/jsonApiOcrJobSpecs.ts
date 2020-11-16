import { ApiResponse } from '@chainlink/json-api-client'
import { partialAsFull } from '@chainlink/ts-helpers'
import { OcrJobSpec } from 'core/store/models'
import { generateUuid } from '../test-helpers/generateUuid'

function getRandomInt(max: number) {
  return Math.floor(Math.random() * Math.floor(max))
}

export const jsonApiOcrJobSpecs = (
  jobs: Partial<
    OcrJobSpec['offChainReportingOracleSpec'] & { id?: string }
  >[] = [],
) => {
  return {
    data: jobs.map((config) => {
      const id = config.id || getRandomInt(1_000_000).toString()
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
        contractConfigTrackerPollInterval:
          config.contractConfigTrackerPollInterval,
        contractConfigConfirmations: config.contractConfigConfirmations,
        updatedAt: config.updatedAt || new Date(1600775300410).toISOString(),
        createdAt: config.createdAt || new Date(1600775300410).toISOString(),
      })

      return {
        type: 'jobSpecV2s',
        id,
        attributes: {
          offChainReportingOracleSpec,
        },
      }
    }),
  } as ApiResponse<OcrJobSpec[]>
}
