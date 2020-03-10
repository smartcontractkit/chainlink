import {
  Serializer as JSONAPISerializer,
  SerializerOptions,
} from 'jsonapi-serializer'
import { Head } from '../entity/Head'

export const BASE_ATTRIBUTES: Array<string> = [
  'id',
  'coinbase',
  'parentHash',
  'createdAt',
  'txHash',
  'number',
]

const chainlinkNodeSerializer = (head: Head) => {
  const opts = {
    attributes: BASE_ATTRIBUTES,
    keyForAttribute: 'camelCase',
    meta: {},
  } as SerializerOptions

  return new JSONAPISerializer('heads', opts).serialize(head)
}

export default chainlinkNodeSerializer
