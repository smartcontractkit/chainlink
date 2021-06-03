import { ApiResponse } from 'utils/json-api-client'
import {
  DirectRequestJobV2Spec,
  FluxMonitorJobV2Spec,
  JobSpec,
  JobSpecV2,
  OffChainReportingOracleJobV2Spec,
  KeeperV2Spec,
  CronV2Spec,
  WebhookV2Spec,
} from 'core/store/models'
import { stringifyJobSpec, JobSpecFormats } from './utils'

type DIRECT_REQUEST_DEFINITION_VALID_KEYS =
  | 'name'
  | 'initiators'
  | 'tasks'
  | 'startAt'
  | 'endAt'

const asUnknownObject = (object: object) => object as { [key: string]: unknown }

const scrub = ({
  payload,
  keysToRemove,
}: {
  payload: unknown
  keysToRemove: string[]
}): JSONValue => {
  if (typeof payload === 'string' || payload === null) {
    return payload
  }

  if (Array.isArray(payload)) {
    return payload.map((p) => scrub({ payload: p, keysToRemove }))
  }

  if (typeof payload === 'object' && payload !== null) {
    const typedPayload = asUnknownObject(payload)
    const keepers = Object.keys(typedPayload).filter(
      (k) => !keysToRemove.includes(k),
    )
    return keepers.reduce((accumulator, key) => {
      const value = typedPayload[key]
      if (
        value === null ||
        (typeof value === 'object' &&
          value !== null &&
          Object.keys(value).length === 0)
      ) {
        return accumulator
      }
      return { ...accumulator, [key]: value }
    }, {})
  }

  return null
}

type ScrubbedJobSpec = { [key in DIRECT_REQUEST_DEFINITION_VALID_KEYS]: any }

export const generateJSONDefinition = (
  job: ApiResponse<JobSpec>['data']['attributes'],
): string => {
  const scrubbedJobSpec: ScrubbedJobSpec = ([
    'name',
    'initiators',
    'tasks',
    'startAt',
    'endAt',
  ] as DIRECT_REQUEST_DEFINITION_VALID_KEYS[]).reduce((accumulator, key) => {
    const value = scrub({
      payload: job[key],
      keysToRemove: ['ID', 'CreatedAt', 'DeletedAt', 'UpdatedAt'],
    })

    if (value === null) {
      return accumulator
    }
    return {
      ...accumulator,
      [key]: value,
    }
  }, {} as ScrubbedJobSpec)

  /**
   * We want to remove the name field if it was auto-generated
   * to avoid running into FK constraint errors when duplicating
   * a job spec.
   */
  if (scrubbedJobSpec.name.includes(job.id)) {
    delete scrubbedJobSpec.name
  }

  return stringifyJobSpec({
    value: scrubbedJobSpec,
    format: JobSpecFormats.JSON,
  })
}

export const generateTOMLDefinition = (
  jobSpecAttributes: ApiResponse<JobSpecV2>['data']['attributes'],
): string => {
  if (jobSpecAttributes.type === 'directrequest') {
    return generateDirectRequestDefinition(jobSpecAttributes)
  }

  if (jobSpecAttributes.type === 'fluxmonitor') {
    return generateFluxMonitorDefinition(jobSpecAttributes)
  }

  if (jobSpecAttributes.type === 'offchainreporting') {
    return generateOCRDefinition(jobSpecAttributes)
  }

  if (jobSpecAttributes.type === 'keeper') {
    return generateKeeperDefinition(jobSpecAttributes)
  }

  if (jobSpecAttributes.type === 'cron') {
    return generateCronDefinition(jobSpecAttributes)
  }

  if (jobSpecAttributes.type === 'webhook') {
    return generateWebhookDefinition(jobSpecAttributes)
  }

  return ''
}

function generateOCRDefinition(
  attrs: ApiResponse<OffChainReportingOracleJobV2Spec>['data']['attributes'],
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
    },
    format: JobSpecFormats.TOML,
  })
}

function generateFluxMonitorDefinition(
  attrs: ApiResponse<FluxMonitorJobV2Spec>['data']['attributes'],
) {
  const {
    fluxMonitorSpec,
    name,
    pipelineSpec,
    schemaVersion,
    type,
    maxTaskDuration,
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
    minPayment,
  } = fluxMonitorSpec

  return stringifyJobSpec({
    value: {
      type,
      schemaVersion,
      name,
      contractAddress,
      precision,
      threshold,
      absoluteThreshold,
      idleTimerPeriod,
      idleTimerDisabled,
      pollTimerPeriod,
      pollTimerDisabled,
      maxTaskDuration,
      minPayment,
      observationSource: pipelineSpec.dotDagSource,
    },
    format: JobSpecFormats.TOML,
  })
}

function generateDirectRequestDefinition(
  attrs: ApiResponse<DirectRequestJobV2Spec>['data']['attributes'],
) {
  const {
    directRequestSpec,
    name,
    pipelineSpec,
    schemaVersion,
    type,
    maxTaskDuration,
  } = attrs
  const {
    contractAddress,
    onChainJobSpecID,
    minIncomingConfirmations,
  } = directRequestSpec

  return stringifyJobSpec({
    value: {
      type,
      schemaVersion,
      name,
      onChainJobSpecID,
      minIncomingConfirmations,
      contractAddress,
      maxTaskDuration,
      observationSource: pipelineSpec.dotDagSource,
    },
    format: JobSpecFormats.TOML,
  })
}

function generateKeeperDefinition(
  attrs: ApiResponse<KeeperV2Spec>['data']['attributes'],
) {
  const { keeperSpec, name, schemaVersion, type } = attrs
  const { contractAddress, fromAddress } = keeperSpec

  return stringifyJobSpec({
    value: {
      type,
      schemaVersion,
      name,
      contractAddress,
      fromAddress,
    },
    format: JobSpecFormats.TOML,
  })
}

function generateCronDefinition(
  attrs: ApiResponse<CronV2Spec>['data']['attributes'],
) {
  const { cronSpec, pipelineSpec, name, schemaVersion, type } = attrs
  const { schedule } = cronSpec

  return stringifyJobSpec({
    value: {
      type,
      schemaVersion,
      name,
      schedule,
      observationSource: pipelineSpec.dotDagSource,
    },
    format: JobSpecFormats.TOML,
  })
}

function generateWebhookDefinition(
  attrs: ApiResponse<WebhookV2Spec>['data']['attributes'],
) {
  const { pipelineSpec, name, schemaVersion, type, webhookSpec } = attrs
  const { onChainJobSpecID } = webhookSpec

  return stringifyJobSpec({
    value: {
      type,
      schemaVersion,
      name,
      onChainJobSpecID,
      observationSource: pipelineSpec.dotDagSource,
    },
    format: JobSpecFormats.TOML,
  })
}
