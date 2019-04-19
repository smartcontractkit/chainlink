import { Connection } from 'typeorm'
import { Client } from './entity/Client'
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
    const hashInput = `v0-${accessKey}-${secret}-${client.salt}`
    const hash = sha256(hashInput)
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
