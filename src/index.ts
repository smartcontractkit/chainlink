import { createDbConnection } from './database'
import seed from './seed'
import server from './server'

const start = async () => {
  await createDbConnection()
  await seed()
  server()
}

start().catch(console.error)
