import {
  Serializer as JSONAPISerializer,
  SerializerOptions,
} from 'jsonapi-serializer'
import { Head } from '../entity/Head'
import { BASE_ATTRIBUTES, transform } from './headSerializer'

const chainlinkNodesSerializer = (heads: Head[], count: number) => {
  const opts = {
    attributes: BASE_ATTRIBUTES,
    keyForAttribute: 'camelCase',
    meta: { count },
    transform,
  } as SerializerOptions

  return new JSONAPISerializer('heads', opts).serialize(heads)
}

export default chainlinkNodesSerializer
