import yargs from 'yargs'
import { Connection } from 'typeorm'
import { Client, createClient, deleteClient } from '../entity/Client'
import { createDbConnection, closeDbConnection, getDb } from '../database'

const add = async (name: string) => {
  return bootstrap(async (db: Connection) => {
    const [client, secret] = await createClient(db, name)
    console.log('created new client with id %s', client.id)
    console.log('AccessKey', client.accessKey)
    console.log('Secret', secret)
  })
}

const remove = async (name: string) => {
  return bootstrap(async (db: Connection) => {
    deleteClient(db, name)
  })
}

const bootstrap = async (cb: any) => {
  await createDbConnection()
  const db = getDb()
  try {
    await cb(db)
  } catch (e) {
    console.error(e)
  }
  //try {
    //await closeDbConnection()
  //} catch (e) {}
}

const _ = yargs
  .usage('Usage: $0 <command> [options]')
  .command({
    command: 'add <name>',
    aliases: 'create',
    describe: 'Add a client',
    builder: (yargs): any => {
      yargs.positional('name', {
        describe: 'The name of the Core node to create',
        type: 'string'
      })
    },
    handler: argv => add(argv.name as string)
  })
  .command({
    command: 'delete <name>',
    aliases: 'rm',
    describe: 'Remove a client',
    builder: (yargs): any => {
      yargs.positional('name', {
        describe: 'The name of the Core node to remove',
        type: 'string'
      })
    },
    handler: argv => remove(argv.name as string)
  })
  .help('h')
  .alias('h', 'help')
  .demandCommand(1).argv // final argv call invokes command
