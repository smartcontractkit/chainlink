import {
  Serializer as JSONAPISerializer,
  SerializerOptions,
} from 'jsonapi-serializer'
import { JobRun } from '../entity/JobRun'
import { BASE_ATTRIBUTES } from './chainlinkNodeSerializerlizer'

const jobRunsSerializer = (runs: JobRun[], runCount: number) => {
  const opts = {
    attributes: BASE_ATTRIBUTES,
    keyForAttribute: 'camelCase',
    meta: { count: runCount },
  } as SerializerOptions

  return new JSONAPISerializer('job_runs', opts).serialize(runs)
}

export default jobRunsSerializer
