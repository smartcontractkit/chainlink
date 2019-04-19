import { Connection } from 'typeorm'
import { Client, hashCredentials } from './entity/Client'
import { sha256 } from 'js-sha256'
import { timingSafeEqual } from 'crypto'

// Session contains a client's ID and access key
export interface Session {
  clientId: number
  accessKey: string
}

// authenticate looks up a client by accessKey and attempts to verify the
// provided secret, if verification succeeds a Session is returned
export const authenticate = async (
  db: Connection,
  accessKey: string,
  secret: string
): Promise<Session | null> => {
  const client = await db.getRepository(Client).findOne({
    accessKey: accessKey
  })

  if (client != null) {
    const hash = hashCredentials(accessKey, secret, client.salt)
    if (
      timingSafeEqual(Buffer.from(hash), Buffer.from(client.hashedSecret)) ===
      true
    ) {
      return {
        clientId: client.id,
        accessKey: accessKey
      }
    }
  }

  return null
}
