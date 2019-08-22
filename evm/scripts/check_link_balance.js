// truffle script

const commandLineArgs = require('command-line-args')
const LinkToken = artifacts.require('LinkTokenInterface')

// compand line options
const optionDefinitions = [
  { name: 'args', type: String, multiple: true, defaultOption: true },
  { name: 'compile', type: Boolean },
  { name: 'network', type: String }
]

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

module.exports = async callback => {
  try {
    await main()
    callback()
  } catch (error) {
    console.error(`Usage: ${USAGE}`)
    callback(error)
  }
}
