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

export interface GeneratedTOMLDefinition {
  definition: string
  envAttributesDefinition: string
}

export const generateTOMLDefinition = (
  jobSpecAttributes: ApiResponse<Job>['data']['attributes'],
): GeneratedTOMLDefinition => {
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
      return { definition: '', envAttributesDefinition: '' }
  }
}

const envAttributesRegex = /Env$/

const generateEnvAttributesDefinition = (jobSpec: {
  [value: string]: any
}): {
  envDefinition: string
  cleanedJobSpec: { [value: string]: any }
} => {
  const spec: { [value: string]: any } = {
    ...jobSpec,
  }
  const envAttributesObject: { [value: string]: any } = {}

  for (const specAttributeName in spec) {
    if (!specAttributeName.match(envAttributesRegex)) {
      continue
    }

    const attributeName = specAttributeName.replace(envAttributesRegex, '')
    envAttributesObject[attributeName] = spec[attributeName]

    delete spec[attributeName]
    delete spec[specAttributeName]
  }

  return {
    cleanedJobSpec: spec,
    envDefinition: stringifyJobSpec({
      value: envAttributesObject,
    }),
  }
}

function generateOCRDefinition(
  attrs: ApiResponse<OffChainReportingJob>['data']['attributes'],
): GeneratedTOMLDefinition {
  const spec = generateEnvAttributesDefinition(
    attrs.offChainReportingOracleSpec,
  )

  const ocrSpecWithoutDates = {
    ...spec.cleanedJobSpec,
    createdAt: undefined,
    updatedAt: undefined,
  }

  const definition = stringifyJobSpec({
    value: {
      type: attrs.type,
      schemaVersion: attrs.schemaVersion,
      ...ocrSpecWithoutDates,
      observationSource: attrs.pipelineSpec.dotDagSource,
      maxTaskDuration: attrs.maxTaskDuration,
      externalJobID: attrs.externalJobID,
    },
  })

  return { definition, envAttributesDefinition: spec.envDefinition }
}

function generateFluxMonitorDefinition(
  attrs: ApiResponse<FluxMonitorJob>['data']['attributes'],
): GeneratedTOMLDefinition {
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

  const definition = stringifyJobSpec({
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

  return { definition, envAttributesDefinition: '' }
}

function generateDirectRequestDefinition(
  attrs: ApiResponse<DirectRequestJob>['data']['attributes'],
): GeneratedTOMLDefinition {
  const {
    directRequestSpec,
    name,
    pipelineSpec,
    schemaVersion,
    type,
    maxTaskDuration,
    externalJobID,
  } = attrs
  const { contractAddress, minIncomingConfirmations } = directRequestSpec

  const definition = stringifyJobSpec({
    value: {
      type,
      schemaVersion,
      name,
      minIncomingConfirmations,
      contractAddress,
      maxTaskDuration,
      observationSource: pipelineSpec.dotDagSource,
      externalJobID,
    },
  })

  return { definition, envAttributesDefinition: '' }
}

function generateKeeperDefinition(
  attrs: ApiResponse<KeeperJob>['data']['attributes'],
): GeneratedTOMLDefinition {
  const { keeperSpec, name, schemaVersion, type, externalJobID } = attrs
  const { contractAddress, fromAddress } = keeperSpec

  const definition = stringifyJobSpec({
    value: {
      type,
      schemaVersion,
      name,
      contractAddress,
      fromAddress,
      externalJobID,
    },
  })

  return { definition, envAttributesDefinition: '' }
}

function generateCronDefinition(
  attrs: ApiResponse<CronJob>['data']['attributes'],
): GeneratedTOMLDefinition {
  const { cronSpec, pipelineSpec, name, schemaVersion, type, externalJobID } =
    attrs
  const { schedule } = cronSpec

  const definition = stringifyJobSpec({
    value: {
      type,
      schemaVersion,
      name,
      schedule,
      observationSource: pipelineSpec.dotDagSource,
      externalJobID,
    },
  })

  return { definition, envAttributesDefinition: '' }
}

function generateWebhookDefinition(
  attrs: ApiResponse<WebhookJob>['data']['attributes'],
): GeneratedTOMLDefinition {
  const { pipelineSpec, name, schemaVersion, type, externalJobID } = attrs

  const definition = stringifyJobSpec({
    value: {
      type,
      schemaVersion,
      name,
      externalJobID,
      observationSource: pipelineSpec.dotDagSource,
    },
  })

  return { definition, envAttributesDefinition: '' }
}

function generateVRFDefinition(
  attrs: ApiResponse<VRFJob>['data']['attributes'],
): GeneratedTOMLDefinition {
  const { vrfSpec, name, schemaVersion, type, externalJobID, pipelineSpec } =
    attrs
  const {
    coordinatorAddress,
    confirmations,
    publicKey,
    fromAddress,
    pollPeriod,
  } = vrfSpec

  const definition = stringifyJobSpec({
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

  return { definition, envAttributesDefinition: '' }
}
