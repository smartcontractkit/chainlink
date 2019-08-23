// truffle script

const { utils } = require('ethers')
const commandLineArgs = require('command-line-args')
const { optionsDefinitions, provider, scriptRunner } = require('./common')

const USAGE =
  'truffle exec scripts/view_eth_price.js [options] <contract address>'

const main = async () => {
  // parse command line args
  const options = commandLineArgs(optionsDefinitions)
  const [consumer] = options.args.slice(2)
  // encode function call
  const funcSelector = '0x9d1b464a' // "currentPrice()"
  // make function call
  const hexPrice = await provider.call({
    data: funcSelector,
    to: consumer
  })
  // print price
  const price = utils.toUtf8String(hexPrice)
  const msg = price ? `current ETH price: ${price}` : 'No price listed'
  console.log(msg)
}

module.exports = scriptRunner(main, USAGE)
