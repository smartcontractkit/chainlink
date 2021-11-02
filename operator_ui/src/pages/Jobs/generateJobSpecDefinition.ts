import { ApiResponse } from 'utils/json-api-client'
import {
  DirectRequestJob,
  FluxMonitorJob,
  Job,
  OffChainReportingJob,
  KeeperJob,
  CronJob,
  WebhookJob,
  VRFJob,
} from 'core/store/models'
import { stringifyJobSpec } from './utils'

export const generateTOMLDefinition = (
  jobSpecAttributes: ApiResponse<Job>['data']['attributes'],
): string => {
  switch (jobSpecAttributes.type) {
    case 'directrequest':
      return generateDirectRequestDefinition(jobSpecAttributes)
    case 'fluxmonitor':
      return generateFluxMonitorDefinition(jobSpecAttributes)
    case 'offchainreporting':
      return generateOCRDefinition(jobSpecAttributes)
    case 'keeper':
      return generateKeeperDefinition(jobSpecAttributes)
    case 'cron':
      return generateCronDefinition(jobSpecAttributes)
    case 'webhook':
      return generateWebhookDefinition(jobSpecAttributes)
    case 'vrf':
      return generateVRFDefinition(jobSpecAttributes)
    default:
      return ''
  }
}

function generateOCRDefinition(
  attrs: ApiResponse<OffChainReportingJob>['data']['attributes'],
) {
  const ocrSpecWithoutDates = {
    ...attrs.offChainReportingOracleSpec,
    createdAt: undefined,
    updatedAt: undefined,
  }

  return stringifyJobSpec({
    value: {
      type: attrs.type,
      schemaVersion: attrs.schemaVersion,
      ...ocrSpecWithoutDates,
      observationSource: attrs.pipelineSpec.dotDagSource,
      maxTaskDuration: attrs.maxTaskDuration,
      externalJobID: attrs.externalJobID,
    },
  })
}

function generateFluxMonitorDefinition(
  attrs: ApiResponse<FluxMonitorJob>['data']['attributes'],
) {
  const {
    fluxMonitorSpec,
    name,
    pipelineSpec,
    schemaVersion,
    type,
    maxTaskDuration,
    externalJobID,
  } = attrs
  const {
    contractAddress,
    precision,
    threshold,
    absoluteThreshold,
    idleTimerPeriod,
    idleTimerDisabled,
    pollTimerPeriod,
    pollTimerDisabled,
    drumbeatEnabled,
    drumbeatSchedule,
    drumbeatRandomDelay,
    minPayment,
  } = fluxMonitorSpec

  return stringifyJobSpec({
    value: {
      type,
      schemaVersion,
      name,
      contractAddress,
      precision,
      threshold: threshold || null,
      absoluteThreshold: absoluteThreshold || null,
      idleTimerPeriod,
      idleTimerDisabled,
      pollTimerPeriod,
      pollTimerDisabled,
      drumbeatEnabled,
      drumbeatSchedule: drumbeatSchedule || null,
      drumbeatRandomDelay: drumbeatRandomDelay || null,
      maxTaskDuration,
      minPayment,
      observationSource: pipelineSpec.dotDagSource,
      externalJobID,
    },
  })
}

function generateDirectRequestDefinition(
  attrs: ApiResponse<DirectRequestJob>['data']['attributes'],
) {
  const {
    directRequestSpec,
    name,
    pipelineSpec,
    schemaVersion,
    type,
    maxTaskDuration,
    externalJobID,
  } = attrs
  const { contractAddress, minIncomingConfirmations, requesters } =
    directRequestSpec

  return stringifyJobSpec({
    value: {
      type,
      schemaVersion,
      name,
      minIncomingConfirmations,
      contractAddress,
      requesters,
      maxTaskDuration,
      observationSource: pipelineSpec.dotDagSource,
      externalJobID,
    },
  })
}

function generateKeeperDefinition(
  attrs: ApiResponse<KeeperJob>['data']['attributes'],
) {
  const { keeperSpec, name, schemaVersion, type, externalJobID } = attrs
  const { contractAddress, fromAddress } = keeperSpec

  return stringifyJobSpec({
    value: {
      type,
      schemaVersion,
      name,
      contractAddress,
      fromAddress,
      externalJobID,
    },
  })
}

function generateCronDefinition(
  attrs: ApiResponse<CronJob>['data']['attributes'],
) {
  const { cronSpec, pipelineSpec, name, schemaVersion, type, externalJobID } =
    attrs
  const { schedule } = cronSpec

  return stringifyJobSpec({
    value: {
      type,
      schemaVersion,
      name,
      schedule,
      observationSource: pipelineSpec.dotDagSource,
      externalJobID,
    },
  })
}

function generateWebhookDefinition(
  attrs: ApiResponse<WebhookJob>['data']['attributes'],
) {
  const { pipelineSpec, name, schemaVersion, type, externalJobID } = attrs

  return stringifyJobSpec({
    value: {
      type,
      schemaVersion,
      name,
      externalJobID,
      observationSource: pipelineSpec.dotDagSource,
    },
  })
}

function generateVRFDefinition(
  attrs: ApiResponse<VRFJob>['data']['attributes'],
) {
  const { vrfSpec, name, schemaVersion, type, externalJobID, pipelineSpec } =
    attrs
  const {
    coordinatorAddress,
    confirmations,
    publicKey,
    fromAddress,
    pollPeriod,
  } = vrfSpec

  return stringifyJobSpec({
    value: {
      type,
      schemaVersion,
      name,
      externalJobID,
      coordinatorAddress,
      confirmations,
      publicKey,
      fromAddress,
      pollPeriod,
      observationSource: pipelineSpec.dotDagSource,
    },
  })
}
