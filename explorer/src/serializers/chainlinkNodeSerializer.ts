import {
  Serializer as JSONAPISerializer,
  SerializerOptions,
} from 'jsonapi-serializer'
import { ChainlinkNode } from '../entity/ChainlinkNode'

export const BASE_ATTRIBUTES: Array<string> = [
  'id',
  'name',
  'url',
  'createdAt',
  'coreVersion',
  'coreSHA',
]

const chainlinkNodeSerializer = (chainlinkNode: ChainlinkNode) => {
  const opts = {
    attributes: BASE_ATTRIBUTES,
    keyForAttribute: 'camelCase',
    meta: {},
  } as SerializerOptions

  return new JSONAPISerializer('chainlink_nodes', opts).serialize(chainlinkNode)
}

export default chainlinkNodeSerializer
