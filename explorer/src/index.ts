import { getConfig } from './config'
import { getDb } from './database'
import { retireSessions } from './entity/Session'
import { logger } from './logging'
import server from './server'
import { getVersion } from './utils/version'

async function main() {
  const conf = getConfig()
  const version = await getVersion(conf)
  logger.info(version)

  try {
    const db = await getDb()
    logger.info('Cleaning up sessions...')
    await retireSessions(db)

    logger.info('Starting Explorer Node')
    await server(conf)
  } catch (e) {
    logger.error({
      msg: `Exception during startup: ${e.message}`,
      stack: e.stack,
    })
  }
}

main()
