import yargs from 'yargs'
import { Connection } from 'typeorm'
import { Admin } from '../entity/admin'
import { createAdmin } from '../support/admin'
import { closeDbConnection, getDb } from '../database'

async function bootstrap(cb: any) {
  const db = await getDb()
  try {
    await cb(db)
  } catch (e) {
    console.error(e)
  }
  try {
    await closeDbConnection()
  } catch (e) {
    console.error(e)
  }
}

const seed = async (username: string, password: string) => {
  return bootstrap(async (db: Connection) => {
    const admin: Admin = await createAdmin(db, username, password)

    console.log('created new chainlink admin')
    console.log('username: ', admin.username)
    console.log('password: ', password)
  })
}

const _ = yargs
  .usage('Usage: $0 <command> [options]')
  .command({
    command: 'seed <username> <password>',
    aliases: 's',
    describe: 'Seed an admin user',
    builder: (yargs): any => {
      yargs
        .positional('username', {
          describe: 'The username of the Chainlink admin to create',
          type: 'string',
        })
        .positional('password', {
          describe: 'The password of the Chainlink admin to create',
          type: 'string',
        })
    },
    handler: argv => seed(argv.username as string, argv.password as string),
  })
  .help('h')
  .alias('h', 'help')
  .demandCommand(1).argv
