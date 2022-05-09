import TOML from '@iarna/toml'
import pick from 'lodash/pick'

export interface JobDefinition {
  definition: string
  envDefinition: string
}

// Extracts fields from the job that are common to all specs.
const extractJobFields = (job: JobPayload_Fields, ...otherKeys: string[]) => {
  return pick(
    job,
    'type',
    'schemaVersion',
    'name',
    'externalJobID',
    ...otherKeys,
  )
}

// Extracts the observation source from the job
const extractObservationSourceField = ({
  observationSource,
}: JobPayload_Fields) => {
  return {
    observationSource: observationSource === '' ? null : observationSource,
  }
}

// Extracts the fields matching the keys from the spec. If another field of of
// the same name with an 'Env' suffix exists, we remove it from the returned
// object.
const extractSpecFields = <T extends {}, K extends keyof T>(
  spec: T,
  ...keys: K[]
) => {
  // For every key, check for the existence of an another field of the same name
  // with an 'Env' suffix
  const scopedKeys = keys.filter((key) => {
    const envKey = `${key}Env` as K
    if (Object.prototype.hasOwnProperty.call(spec, envKey)) {
      // We are relying on this always being a boolean but we can't guarantee it
      return !spec[envKey]
    }

    return true
  })

  return pick(spec, ...scopedKeys)
}

// Extracts the fields which have a field of the same name with an 'Env' suffix
// and the 'Env' field returns true.
const extractEnvValues = <T extends {}, K extends keyof T>(spec: T) => {
  // For every key with an 'Env' suffix, find a key of the same name without the
  // suffix.
  const regex = /(.+)Env$/
  const envValueKeys: K[] = []
  for (const key of Object.keys(spec)) {
    const match = key.match(regex)

    if (match) {
      const envKey = key as K
      // We are relying on this always being a boolean but we can't guarantee it
      if (spec[envKey]) {
        // Check that the key without the 'Env' suffix exists
        if (Object.prototype.hasOwnProperty.call(spec, match[1])) {
          envValueKeys.push(match[1] as K)
        }
      }
    }
  }

  return pick(spec, ...envValueKeys)
}

// Stringifies the job spec as TOML
const toTOMLString = (value: { [key: string]: any }) => {
  try {
    return TOML.stringify(value)
  } catch (e) {
    console.error(`Failed to stringify job spec with the following error: ${e}`)
    return ''
  }
}

export const generateJobDefinition = (
  job: JobPayload_Fields,
): JobDefinition => {
  let values: object = {}

  switch (job.spec.__typename) {
    case 'CronSpec':
      values = {
        ...extractJobFields(job, 'maxTaskDuration'),
        ...extractSpecFields(job.spec, 'schedule'),
        ...extractObservationSourceField(job),
      }

      break
    case 'DirectRequestSpec':
      values = {
        ...extractJobFields(job, 'maxTaskDuration'),
        ...extractSpecFields(
          job.spec,
          'contractAddress',
          'evmChainID',
          'minIncomingConfirmations',
          'minContractPaymentLinkJuels',
          'requesters',
        ),
        ...extractObservationSourceField(job),
      }

      break
    case 'FluxMonitorSpec':
      values = {
        ...extractJobFields(job, 'maxTaskDuration'),
        ...extractSpecFields(
          job.spec,
          'absoluteThreshold',
          'contractAddress',
          'drumbeatEnabled',
          'drumbeatSchedule',
          'drumbeatRandomDelay',
          'evmChainID',
          'idleTimerPeriod',
          'idleTimerDisabled',
          'minPayment',
          'pollTimerPeriod',
          'pollTimerDisabled',
          'threshold',
        ),
        ...extractObservationSourceField(job),
      }

      break
    case 'KeeperSpec':
      values = {
        ...extractJobFields(job),
        ...extractSpecFields(
          job.spec,
          'contractAddress',
          'evmChainID',
          'fromAddress',
        ),
        ...extractObservationSourceField(job),
      }

      break
    case 'OCRSpec':
      values = {
        ...extractJobFields(job, 'maxTaskDuration'),
        ...extractSpecFields(
          job.spec,
          'blockchainTimeout',
          'contractAddress',
          'contractConfigConfirmations',
          'contractConfigTrackerPollInterval',
          'contractConfigTrackerSubscribeInterval',
          'evmChainID',
          'isBootstrapPeer',
          'keyBundleID',
          'observationTimeout',
          'p2pBootstrapPeers',
          'transmitterAddress',
        ),
        ...extractObservationSourceField(job),
      }

      break
    case 'OCR2Spec':
      values = {
        ...extractJobFields(job, 'maxTaskDuration'),
        ...extractSpecFields(
          job.spec,
          'blockchainTimeout',
          'contractID',
          'contractConfigConfirmations',
          'contractConfigTrackerPollInterval',
          'ocrKeyBundleID',
          'monitoringEndpoint',
          'p2pBootstrapPeers',
          'relay',
          'relayConfig',
          'pluginType',
          'pluginConfig',
        ),
        // We need to call 'extractSpecFields' again here so we get the spec
        // fields displaying in alphabetical order.
        ...extractSpecFields(job.spec, 'transmitterID'),
        ...extractObservationSourceField(job),
      }

      break
    case 'VRFSpec':
      values = {
        ...extractJobFields(job),
        ...extractSpecFields(
          job.spec,
          'coordinatorAddress',
          'evmChainID',
          'fromAddresses',
          'minIncomingConfirmations',
          'pollPeriod',
          'publicKey',
          'requestedConfsDelay',
          'requestTimeout',
          'batchCoordinatorAddress',
          'batchFulfillmentEnabled',
          'batchFulfillmentGasMultiplier',
          'chunkSize',
          'backoffInitialDelay',
          'backoffMaxDelay',
        ),
        ...extractObservationSourceField(job),
      }

      break

    case 'BlockhashStoreSpec':
      values = {
        ...extractJobFields(job),
        ...extractSpecFields(
          job.spec,
          'coordinatorV1Address',
          'coordinatorV2Address',
          'waitBlocks',
          'lookbackBlocks',
          'blockhashStoreAddress',
          'pollPeriod',
          'runTimeout',
          'evmChainID',
          'fromAddress',
        ),
        ...extractObservationSourceField(job),
      }

      break
    case 'BootstrapSpec':
      values = {
        ...extractJobFields(job),
        ...extractSpecFields(
          job.spec,
          'id',
          'contractID',
          'relay',
          'relayConfig',
          'monitoringEndpoint',
          'blockchainTimeout',
          'contractConfigTrackerPollInterval',
          'contractConfigConfirmations',
        ),
      }

      break
    default:
      return { definition: '', envDefinition: '' }
    case 'WebhookSpec':
      values = {
        ...extractJobFields(job),
        ...extractObservationSourceField(job),
      }

      break
  }

  return {
    definition: toTOMLString(values),
    envDefinition: toTOMLString(extractEnvValues(job.spec)),
  }
}
