import { Connection } from 'typeorm'
import {
  ChainlinkNode,
  createChainlinkNode,
  deleteChainlinkNode,
  hashCredentials,
} from '../../entity/ChainlinkNode'
import { closeDbConnection, getDb } from '../../database'

let db: Connection

beforeAll(async () => {
  db = await getDb()
})

afterAll(async () => closeDbConnection())

describe('createChainlinkNode', () => {
  it('returns a valid ChainlinkNode record', async () => {
    const [chainlinkNode, secret] = await createChainlinkNode(
      db,
      'new-valid-chainlink-node-record',
    )
    expect(chainlinkNode.accessKey).toHaveLength(16)
    expect(chainlinkNode.salt).toHaveLength(32)
    expect(chainlinkNode.hashedSecret).toBeDefined()
    expect(secret).toHaveLength(64)
  })

  it('reject duplicate ChainlinkNode names', async () => {
    await createChainlinkNode(db, 'identical')
    await expect(createChainlinkNode(db, 'identical')).rejects.toThrow()
  })
})

describe('deleteChainlinkNode', () => {
  it('deletes a ChainlinkNode with the specified name', async () => {
    await createChainlinkNode(db, 'chainlink-node-to-be-deleted')
    let count = await db.manager.count(ChainlinkNode)
    expect(count).toBe(1)
    await deleteChainlinkNode(db, 'chainlink-node-to-be-deleted')
    count = await db.manager.count(ChainlinkNode)
    expect(count).toBe(0)
  })
})

describe('hashCredentials', () => {
  it('returns a sha256 signature', () => {
    expect(hashCredentials('a', 'b', 'c')).toHaveLength(64)
  })
})
