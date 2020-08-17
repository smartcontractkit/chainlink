import { authenticate } from '../sessions'
import { getRepository } from 'typeorm'
import { createChainlinkNode } from '../entity/ChainlinkNode'
import { Session, closeSession } from '../entity/Session'

describe('sessions', () => {
  describe('authenticate', () => {
    it('creates a session record', async () => {
      const [chainlinkNode, secret] = await createChainlinkNode(
        'valid-chainlink-node',
      )
      const session = await authenticate({
        accessKey: chainlinkNode.accessKey,
        secret,
      })
      expect(session).toBeDefined()
      expect(session.chainlinkNodeId).toEqual(chainlinkNode.id)

      let foundSession = await getRepository(Session).findOne()
      expect(foundSession.chainlinkNodeId).toEqual(chainlinkNode.id)
      expect(foundSession.finishedAt).toBeNull()

      await authenticate({ accessKey: chainlinkNode.accessKey, secret })
      foundSession = await getRepository(Session).findOne(foundSession.id)
      expect(foundSession.finishedAt).toEqual(expect.any(Date))
    })

    it('closes a session', async () => {
      const [chainlinkNode, secret] = await createChainlinkNode(
        'valid-chainlink-node',
      )
      const session = await authenticate({
        accessKey: chainlinkNode.accessKey,
        secret,
      })
      expect(session).toBeDefined()
      expect(session.chainlinkNodeId).toEqual(chainlinkNode.id)

      closeSession(session)
    })

    it('returns null if no chainlink node exists', async () => {
      const result = await authenticate({ accessKey: '', secret: '' })
      expect(result).toBeNull()
    })

    it('returns null if the secret is incorrect', async () => {
      const [chainlinkNode] = await createChainlinkNode(
        'invalid-chainlink-node',
      )
      const result = await authenticate({
        accessKey: chainlinkNode.accessKey,
        secret: 'wrong-secret',
      })
      expect(result).toBeNull()
    })
  })

  describe('closeSession', () => {
    it('closes an open session', async () => {
      const [chainlinkNode, secret] = await createChainlinkNode(
        'valid-chainlink-node',
      )
      const session = await authenticate({
        accessKey: chainlinkNode.accessKey,
        secret,
      })
      expect(session).toBeDefined()
      expect(session.chainlinkNodeId).toEqual(chainlinkNode.id)
      expect(session.finishedAt).toBeNull()

      await closeSession(session)

      const foundSession = await getRepository(Session).findOne()
      expect(foundSession.chainlinkNodeId).toEqual(chainlinkNode.id)
      expect(foundSession.finishedAt).toEqual(expect.any(Date))
    })
  })
})
