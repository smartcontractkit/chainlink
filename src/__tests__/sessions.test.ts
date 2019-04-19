import { Connection } from 'typeorm'
import { closeDbConnection, getDb } from '../database'
import { createClient } from '../entity/Client'
import { parseAuthenticationRequest, authenticateRequest } from '../sessions'

describe('sessions', () => {
  let db: Connection
  beforeAll(async () => {
    db = await getDb()
  })
  afterAll(async () => {
    await closeDbConnection()
  })

  describe('parseAuthenticationRequest', () => {
    it('saves access key and secret', () => {
      const request = parseAuthenticationRequest(
        '{"accessKey": "917", "secret": "awkward"}'
      )

      expect(request.accessKey).toEqual('917')
      expect(request.secret).toEqual('awkward')
    })
  })

  describe('authenticate', () => {
    it('returns the session', async (done: any) => {
      const [client, secret] = await createClient(db, 'valid-client')
      const session = await authenticateRequest(db, {
        accessKey: client.accessKey,
        secret: secret
      })
      expect(session).toBeDefined()
      expect(session.clientId).toEqual(client.id)
      done()
    })

    it('returns null if no client exists', async (done: any) => {
      const result = await authenticateRequest(db, {
        accessKey: '',
        secret: ''
      })
      expect(result).toBeNull()
      done()
    })

    it('returns null if the secret is incorrect', async (done: any) => {
      const [client, _] = await createClient(db, 'invalid-client')
      const result = await authenticateRequest(db, {
        accessKey: client.accessKey,
        secret: 'wrong-secret'
      })
      expect(result).toBeNull()
      done()
    })
  })
})
