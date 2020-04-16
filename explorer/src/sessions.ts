import { getConnection, EntityManager } from 'typeorm'
import { ChainlinkNode, hashCredentials } from './entity/ChainlinkNode'
import { createSession, Session } from './entity/Session'
import { timingSafeEqual } from 'crypto'

// authenticate looks up a chainlink node by accessKey and attempts to verify the
// provided secret, if verification succeeds a Session is returned
export const authenticate = async (
  accessKey: string,
  secret: string,
): Promise<Session | null> => {
  return getConnection().transaction(async (manager: EntityManager) => {
    const chainlinkNode = await findNode(manager, accessKey)
    if (chainlinkNode != null) {
      if (authenticateSession(accessKey, secret, chainlinkNode)) {
        return createSession(chainlinkNode, manager)
      }
    }

    return null
  })
}

function findNode(manager: EntityManager, accessKey: string): Promise<ChainlinkNode> {
  return manager.getRepository(ChainlinkNode).findOne({ accessKey })
}

function authenticateSession(
  accessKey: string,
  secret: string,
  node: ChainlinkNode,
): boolean {
  const hash = hashCredentials(accessKey, secret, node.salt)
  return timingSafeEqual(Buffer.from(hash), Buffer.from(node.hashedSecret))
}
