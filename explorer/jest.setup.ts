import { Connection } from 'typeorm'
import { openDbConnection } from './src/database'
import { clearDb } from './src/__tests__/testdatabase'

process.env.NODE_ENV = 'test'

let globalDBConnection: Connection
beforeAll(async () => {
  globalDBConnection = await openDbConnection()
})
afterEach(() => clearDb())
afterAll(async () => {
  await globalDBConnection.close()
})
