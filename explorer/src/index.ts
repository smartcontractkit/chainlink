import { getConfig } from './config'
import { openDbConnection } from './database'
import { retireSessions } from './entity/Session'
import { logger } from './logging'
import server from './server'
import { getVersion } from './utils/version'

async function main() {
  const conf = getConfig()
  const version = await getVersion(conf)
  logger.info(version)

  try {
    const db = await openDbConnection()

    logger.info('Cleaning up sessions...')
    await retireSessions(db)

    logger.info('Starting Explorer Node')
    await server(conf, db)

  } catch (e) {
    logger.error({
      msg: `Exception during startup: ${e.message}`,
      stack: e.stack,
    })
  }
}

main()
