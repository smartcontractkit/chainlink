import { authenticate } from '../sessions'
import { Connection } from 'typeorm'
import { closeDbConnection, getDb } from '../database'
import { createChainlinkNode } from '../entity/ChainlinkNode'
import { Session } from '../entity/Session'

describe('sessions', () => {
  let db: Connection
  beforeAll(async () => {
    db = await getDb()
  })
  afterAll(async () => {
    await closeDbConnection()
  })

  describe('authenticate', () => {
    it('creates a session record', async () => {
      const [chainlinkNode, secret] = await createChainlinkNode(
        db,
        'valid-chainlink-node'
      )
      const session = await authenticate(db, chainlinkNode.accessKey, secret)
      expect(session).toBeDefined()
      expect(session.chainlinkNodeId).toEqual(chainlinkNode.id)

      let foundSession = await db.manager.findOne(Session)
      expect(foundSession.chainlinkNodeId).toEqual(chainlinkNode.id)
      expect(foundSession.finishedAt).toBeNull()

      await authenticate(db, chainlinkNode.accessKey, secret)
      foundSession = await db.manager.findOne(Session, foundSession.id)
      expect(foundSession.finishedAt).toBeDefined()
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
