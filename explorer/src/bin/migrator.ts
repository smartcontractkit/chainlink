import yargs from 'yargs'
import { Connection } from 'typeorm'
import { closeDbConnection, getDb } from '../database'

const migrate = async () => {
  return bootstrap(async (db: Connection) => {
    console.log(`Migrating [\x1b[32m${db.options.database}\x1b[0m]...`)

    const pendingMigrations = await db.runMigrations()
    for (const m of pendingMigrations) {
      console.log('ran', m)
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
