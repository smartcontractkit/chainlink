import yargs from 'yargs'
import { seed } from '../cli/admin'

/* eslint-disable-next-line no-unused-expressions */
yargs
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
