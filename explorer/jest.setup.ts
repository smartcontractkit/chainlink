import { getConnection } from 'typeorm'
import { openDbConnection } from './src/database'
import { clearDb } from './src/__tests__/testdatabase'

process.env.NODE_ENV = 'test'

beforeAll(async () => {
  await openDbConnection()
})
afterEach(() => clearDb())
afterAll(async () => {
  try {
    await getConnection().close()
  } catch {
    // swallow error or it suppresses all other test output
  }
})
