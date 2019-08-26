import { Connection } from 'typeorm'
import { ChainlinkNode, hashCredentials } from './entity/ChainlinkNode'
import { createSession, Session } from './entity/Session'
import { timingSafeEqual } from 'crypto'

// SessionResponse contains a chainlink node's ID and access key
export interface ISession {
  chainlinkNodeId: number
  sessionId: string
}

// authenticate looks up a chainlink node by accessKey and attempts to verify the
// provided secret, if verification succeeds a Session is returned
export const authenticate = async (
  db: Connection,
  accessKey: string,
  secret: string
): Promise<ISession | null> => {
  return db.manager.transaction(async manager => {
    const chainlinkNode = await findNode(db, accessKey)
    if (chainlinkNode != null) {
      if (authenticateSession(accessKey, secret, chainlinkNode)) {
        const session = await createSession(db, chainlinkNode)
        return {
          chainlinkNodeId: chainlinkNode.id,
          sessionId: session.id
        }
      }
    }

    return null
  })
}

function findNode(db: Connection, accessKey: string): Promise<ChainlinkNode> {
  return db.getRepository(ChainlinkNode).findOne({ accessKey })
}

function authenticateSession(
  accessKey: string,
  secret: string,
  node: ChainlinkNode
): boolean {
  const hash = hashCredentials(accessKey, secret, node.salt)
  return timingSafeEqual(Buffer.from(hash), Buffer.from(node.hashedSecret))
}
