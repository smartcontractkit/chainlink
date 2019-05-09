import {
  Serializer as JSONAPISerializer,
  SerializerOptions
} from 'jsonapi-serializer'
import { JobRun } from '../entity/JobRun'
import { BASE_ATTRIBUTES, chainlinkNode } from './jobRunSerializer'

const jobRunsSerializer = (runs: JobRun[], runCount: number) => {
  const opts = {
    attributes: BASE_ATTRIBUTES,
    chainlinkNode: chainlinkNode,
    keyForAttribute: 'camelCase',
    meta: { count: runCount }
  } as SerializerOptions

  return new JSONAPISerializer('job_runs', opts).serialize(runs)
}

export default jobRunsSerializer
