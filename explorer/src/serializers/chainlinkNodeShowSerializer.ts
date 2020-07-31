import {
  Serializer as JSONAPISerializer,
  SerializerOptions,
} from 'jsonapi-serializer'
import { JobCountReport } from '../entity/ChainlinkNode'

export const ATTRIBUTES: Array<string> = [
  'id',
  'name',
  'url',
  'createdAt',
  'coreVersion',
  'coreSHA',
  'jobCounts',
  'uptime',
]

interface ChainlinkNodeShowData {
  id: number
  name: string
  url?: string
  createdAt: Date
  coreVersion?: string
  coreSHA?: string
  jobCounts: JobCountReport
  uptime: number
}

const chainlinkNodeShowSerializer = (data: ChainlinkNodeShowData) => {
  const opts = {
    attributes: ATTRIBUTES,
    keyForAttribute: 'camelCase',
    meta: {},
  } as SerializerOptions

  return new JSONAPISerializer('chainlink_nodes', opts).serialize(data)
}

export default chainlinkNodeShowSerializer
