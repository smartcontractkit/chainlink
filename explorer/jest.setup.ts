import { getConnection } from 'typeorm'
import { openDbConnection } from './src/database'
import { clearDb } from './src/__tests__/testdatabase'

process.env.NODE_ENV = 'test'

beforeAll(async () => {
  await openDbConnection()
})
afterEach(() => clearDb())
afterAll(async () => {
  const db = getConnection()
  if (db) {
    await db.close()
  }
})
