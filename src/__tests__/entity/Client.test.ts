import { Connection } from 'typeorm'
import {
  Client,
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
    await expect(createClient(db, 'identical')).rejects.toThrow()
  })
})

describe('deleteClient', () => {
  it('deletes a client with the specified name', async () => {
    const [client, _] = await createClient(db, 'default')
    let count = await db.manager.count(Client)
    expect(count).toBe(1)
    await deleteClient(db, 'default')
    count = await db.manager.count(Client)
    expect(count).toBe(0)
  })
})

describe('hashCredentials', () => {
  it('returns a sha256 signature', () => {
    expect(hashCredentials('a', 'b', 'c')).toHaveLength(64)
  })
})
