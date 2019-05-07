import server from './server'

const start = async () => {
  server()
}

start().catch(console.error)
