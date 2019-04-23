import { Connection } from 'typeorm'
import { closeDbConnection, getDb } from '../database'
import { createChainlinkNode } from '../entity/ChainlinkNode'
import { authenticate } from '../sessions'

describe('sessions', () => {
  let db: Connection
  beforeAll(async () => {
    db = await getDb()
  })
  afterAll(async () => {
    await closeDbConnection()
  })

  describe('authenticate', () => {
    it('returns the session', async () => {
      const [chainlinkNode, secret] = await createChainlinkNode(
        db,
        'valid-chainlink-node'
      )
      const session = await authenticate(db, chainlinkNode.accessKey, secret)
      expect(session).toBeDefined()
      expect(session.chainlinkNodeId).toEqual(chainlinkNode.id)
    })

    it('returns null if no chainlink node exists', async () => {
      const result = await authenticate(db, '', '')
      expect(result).toBeNull()
    })

    it('returns null if the secret is incorrect', async () => {
      const [chainlinkNode, _] = await createChainlinkNode(
        db,
        'invalid-chainlink-node'
      )
      const result = await authenticate(
        db,
        chainlinkNode.accessKey,
        'wrong-secret'
      )
      expect(result).toBeNull()
    })
  })
})
