import { Connection } from 'typeorm'
import {
  createClient,
  deleteClient,
  hashCredentials
} from '../../entity/Client'
import { closeDbConnection, getDb } from '../../database'

let db: Connection

beforeAll(async () => {
  db = await getDb()
})

afterAll(async () => closeDbConnection())

describe('createClient', () => {
  it('returns a valid client record', async () => {
    const [client, _] = await createClient(db, 'default')
    expect(client.accessKey).toHaveLength(16)
    expect(client.salt).toHaveLength(32)
    expect(client.hashedSecret).toBeDefined()
  })

  it('returns a secret of at least 16 characters', async () => {
    const [_, secret] = await createClient(db, 'default')
    expect(secret).toHaveLength(16)
  })

  it('reject duplicate client names', async () => {
    await createClient(db, 'identical')
    await expect(async () => createClient(db, 'identical')).rejects.toThrow()
  })
})

describe('deleteClient', () => {
  it('returns a valid client record', async () => {
    const [client, _] = await createClient(db, 'default')
    await deleteClient(db, 'default')
  })

  it('no matching client throws error', async () => {
    await expect(async () => deleteClient(db, 'rare-name')).rejects.toThrow()
  })
})

describe('hashCredentials', () => {
  it('returns a sha256 signature', () => {
    expect(hashCredentials('a', 'b', 'c')).toHaveLength(64)
  })
})
