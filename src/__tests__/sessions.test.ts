import { Connection } from 'typeorm'
import { closeDbConnection, getDb } from '../database'
import { createClient } from '../entity/Client'
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
    it('returns the session', async (done: any) => {
      const [client, secret] = await createClient(db, 'valid-client')
      const session = await authenticate(db, client.accessKey, secret)
      expect(session).toBeDefined()
      expect(session.clientId).toEqual(client.id)
      done()
    })

    it('returns null if no client exists', async (done: any) => {
      const result = await authenticate(db, '', '')
      expect(result).toBeNull()
      done()
    })

    it('returns null if the secret is incorrect', async (done: any) => {
      const [client, _] = await createClient(db, 'invalid-client')
      const result = await authenticate(db, client.accessKey, 'wrong-secret')
      expect(result).toBeNull()
      done()
    })
  })
})
