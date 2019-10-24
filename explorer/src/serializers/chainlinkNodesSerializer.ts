import {
  Serializer as JSONAPISerializer,
  SerializerOptions,
} from 'jsonapi-serializer'
import { ChainlinkNode } from '../entity/ChainlinkNode'
import { BASE_ATTRIBUTES } from './chainlinkNodeSerializer'

const chainlinkNodesSerializer = (
  chainlinkNodes: ChainlinkNode[],
  count: number,
) => {
  const opts = {
    attributes: BASE_ATTRIBUTES,
    keyForAttribute: 'camelCase',
    meta: { count: count },
  } as SerializerOptions

  return new JSONAPISerializer('chainlink_nodes', opts).serialize(
    chainlinkNodes,
  )
}

export default chainlinkNodesSerializer
