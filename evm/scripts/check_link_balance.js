// truffle script

const commandLineArgs = require('command-line-args')
const LinkToken = artifacts.require('LinkTokenInterface')
const { optionDefinitions, scriptRunner } = require('./common')

const USAGE =
  'truffle exec scripts/check_link_balance.js [options] <token address> <holder address>'

const main = async () => {
  // parse command line args
  const options = commandLineArgs(optionDefinitions)
  const [link, holder] = options.args.slice(2)
  // find link token
  const linkToken = await LinkToken.at(link)
  // get address's balance
  const balance = await linkToken.balanceOf.call(holder)
  console.log(`LINK balance: ${balance.toString()}`)
}

module.exports = scriptRunner(main, USAGE)
