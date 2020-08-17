import { getConnection, EntityManager } from 'typeorm'
import { ChainlinkNode, hashCredentials } from './entity/ChainlinkNode'
import { createSession, Session } from './entity/Session'
import { timingSafeEqual } from 'crypto'
import { AuthInfo } from './server/realtime'
import { logger } from './logging'

// authenticate looks up a chainlink node by accessKey and attempts to verify the
// provided secret, if verification succeeds a Session is returned
export const authenticate = async ({
  accessKey,
  secret,
  coreVersion,
  coreSHA,
}: Pick<
  AuthInfo,
  'accessKey' | 'secret' | 'coreVersion' | 'coreSHA'
>): Promise<Session | null> => {
  return getConnection().transaction(async (manager: EntityManager) => {
    const chainlinkNode = await findNode(manager, accessKey)
    if (chainlinkNode != null) {
      if (authenticateSession(accessKey, secret, chainlinkNode)) {
        await recordCoreVersionInfo(chainlinkNode, coreVersion, coreSHA)
        return createSession(chainlinkNode, manager)
      }
    }

    return null
  })
}

async function recordCoreVersionInfo(
  node: ChainlinkNode,
  coreVersion: string,
  coreSHA: string,
) {
  // track the version of core that the node is running
  if (coreVersion && coreSHA) {
    return getConnection()
      .createQueryBuilder()
      .update(ChainlinkNode)
      .set({ coreVersion, coreSHA })
      .where('id = :id', { id: node.id })
      .execute()
      .catch(error => {
        logger.debug({
          msg: `error recording core version and SHA!`,
          nodeID: node.id,
          error,
        })
      })
  }
}

function findNode(
  manager: EntityManager,
  accessKey: string,
): Promise<ChainlinkNode> {
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
