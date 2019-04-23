import seed from './seed'
import server from './server'

const start = async () => {
  await seed()
  server()
}

start().catch(console.error)
