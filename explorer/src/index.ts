import { openDbConnection } from './database'
import { retireSessions } from './entity/Session'
import { logger } from './logging'
import server from './server'
import { getVersion } from './utils/version'
import { Config } from './config'
import { updateClientGaId } from './clientConfig'

async function main() {
  const version = await getVersion(Config.env())
  logger.info(version)
  updateClientGaId(Config.gaId())

  const db = await openDbConnection()
  try {
    logger.info('Cleaning up sessions...')
    await retireSessions()

    logger.info('Starting Explorer Node')
    await server()
  } catch (e) {
    logger.error({
      msg: `Exception during startup: ${e.message}`,
      stack: e.stack,
    })
    db.close()
  }
}

main()
