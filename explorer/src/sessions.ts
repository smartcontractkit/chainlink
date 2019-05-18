import { Connection } from 'typeorm'
import { ChainlinkNode, hashCredentials } from './entity/ChainlinkNode'
import { timingSafeEqual } from 'crypto'

// Session contains a chainlink node's ID and access key
export interface Session {
  chainlinkNodeId: number
  accessKey: string
}

// authenticate looks up a chainlink node by accessKey and attempts to verify the
// provided secret, if verification succeeds a Session is returned
export const authenticate = async (
  db: Connection,
  accessKey: string,
  secret: string
): Promise<Session | null> => {
  const chainlinkNode = await db.getRepository(ChainlinkNode).findOne({
    accessKey: accessKey
  })

  if (chainlinkNode != null) {
    const hash = hashCredentials(accessKey, secret, chainlinkNode.salt)
    if (
      timingSafeEqual(
        Buffer.from(hash),
        Buffer.from(chainlinkNode.hashedSecret)
      )
    ) {
      return {
        chainlinkNodeId: chainlinkNode.id,
        accessKey: accessKey
      }
    }
  }

  return null
}
