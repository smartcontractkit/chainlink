import server from './server'
import { logger } from './logging'

const start = async () => {
  logger.info(`Starting Explorer Node`)
  server()
}

start().catch(logger.error)
