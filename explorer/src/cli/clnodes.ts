import { Connection } from 'typeorm'
import {
  createChainlinkNode,
  deleteChainlinkNode,
} from '../entity/ChainlinkNode'
import { bootstrap } from './bootstrap'

const add = async (name: string, url?: string) => {
  return bootstrap(async (db: Connection) => {
    const [chainlinkNode, secret] = await createChainlinkNode(db, name, url)
    console.log('created new chainlink node with id %s', chainlinkNode.id)
    console.log('AccessKey', chainlinkNode.accessKey)
    console.log('Secret', secret)
  })
}

const remove = async (name: string) => {
  return bootstrap(async (db: Connection) => {
    deleteChainlinkNode(db, name)
  })
}

export { add, remove }
