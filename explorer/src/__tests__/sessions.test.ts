import { authenticate } from '../sessions'
import { getRepository } from 'typeorm'
import { createChainlinkNode } from '../entity/ChainlinkNode'
import { Session } from '../entity/Session'

describe('sessions', () => {
  describe('authenticate', () => {
    it('creates a session record', async () => {
      const [chainlinkNode, secret] = await createChainlinkNode(
        'valid-chainlink-node',
      )
      const session = await authenticate(chainlinkNode.accessKey, secret)
      expect(session).toBeDefined()
      expect(session.chainlinkNodeId).toEqual(chainlinkNode.id)

      let foundSession = await getRepository(Session).findOne()
      expect(foundSession.chainlinkNodeId).toEqual(chainlinkNode.id)
      expect(foundSession.finishedAt).toBeNull()

      await authenticate(chainlinkNode.accessKey, secret)
      foundSession = await getRepository(Session).findOne(foundSession.id)
      expect(foundSession.finishedAt).toBeDefined()
    })

    it('returns null if no chainlink node exists', async () => {
      const result = await authenticate('', '')
      expect(result).toBeNull()
    })

    it('returns null if the secret is incorrect', async () => {
      const [chainlinkNode] = await createChainlinkNode(
        'invalid-chainlink-node',
      )
      const result = await authenticate(chainlinkNode.accessKey, 'wrong-secret')
      expect(result).toBeNull()
    })
  })
})
