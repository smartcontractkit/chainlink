import yargs from 'yargs'
import { Connection } from 'typeorm'
import { closeDbConnection, getDb } from '../database'
import { ChainlinkNode } from '../entity/ChainlinkNode'

const migrate = async () => {
  return bootstrap(async (db: Connection) => {
    console.log(`Migrating [\x1b[32m${db.options.database}\x1b[0m]...`)

    const pendingMigrations = await db.runMigrations()
    for (const m of pendingMigrations) {
      console.log('ran', m)
    }

    if (process.env.COMPOSE_MODE) {
      const repo = db.getRepository(ChainlinkNode)

      const node = new ChainlinkNode()
      node.id = 1
      node.name = 'NodeyMcNodeFace'
      node.accessKey = 'u4HULe0pj5xPyuvv'
      node.hashedSecret =
        '302df2b42ab313cb9b00fe0cca9932dacaaf09e662f2dca1be9c2ad2d927d5df'
      node.salt = 'wZ02sJ8iZ6WffxXduxwzkCfOc3PS8BZJ'

      if (!(await repo.findOne(1))) {
        await repo.save(node)
      }
    }
  })
}

const revert = async () => {
  return bootstrap(async (db: Connection) => {
    await db.undoLastMigration()
  })
}

async function bootstrap(cb: any) {
  const db = await getDb()
  try {
    await cb(db)
  } catch (err) {
    console.error(err)
    process.exit(1)
  } finally {
    await closeDbConnection()
  }
}

yargs
  .usage('Usage: $0 <command> [options]')
  .command({
    command: 'migrate',
    describe: 'Run migrations',
    handler: migrate,
  })
  .command({
    command: 'revert',
    describe: 'Revert last migration',
    handler: revert,
  })
  .help('h')
  .alias('h', 'help')
  .demandCommand(1).argv // final argv call invokes command
