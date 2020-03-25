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

export const transform = (record: Head) => {
  return {
    id: record.id,
    coinbase: record.coinbase.toString('hex'),
    parentHash: record.parentHash.toString('hex'),
    uncleHash: record.uncleHash.toString('hex'),
    txHash: record.txHash.toString('hex'),
    number: record.number,
    createdAt: record.createdAt,
  }
}

const chainlinkNodeSerializer = (head: Head) => {
  const opts = {
    attributes: BASE_ATTRIBUTES,
    keyForAttribute: 'camelCase',
    meta: {},
    transform,
  } as SerializerOptions

  return new JSONAPISerializer('heads', opts).serialize(head)
}

export default chainlinkNodeSerializer
