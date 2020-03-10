import {
  Serializer as JSONAPISerializer,
  SerializerOptions,
} from 'jsonapi-serializer'
import { Head } from '../entity/Head'
import { BASE_ATTRIBUTES } from './headSerializer'

const chainlinkNodesSerializer = (heads: Head[], count: number) => {
  const opts = {
    attributes: BASE_ATTRIBUTES,
    keyForAttribute: 'camelCase',
    meta: { count },
    transform: record => {
      return {
        id: record.id,
        coinbase: record.coinbase.toString('hex'),
        parentHash: record.parentHash.toString('hex'),
        uncleHash: record.uncleHash.toString('hex'),
        txHash: record.txHash.toString('hex'),
        number: record.number,
        createdAt: record.createdAt,
      }
    },
  } as SerializerOptions

  return new JSONAPISerializer('heads', opts).serialize(heads)
}

export default chainlinkNodesSerializer
