// truffle script

const { utils } = require('ethers')
const commandLineArgs = require('command-line-args')
const { provider } = require('../chainlink.config')

// compand line options
const optionDefinitions = [
  { name: 'args', type: String, multiple: true, defaultOption: true },
  { name: 'compile', type: Boolean },
  { name: 'network', type: String }
]

module.exports = async function(callback) {
  // parse command line args
  const options = commandLineArgs(optionDefinitions)
  const [consumer] = options.args.slice(2)
  // encode function call
  const funcSelector = '0x9d1b464a' // "currentPrice()"
  try {
    // make function call
    const hexPrice = await provider.call({
      data: funcSelector,
      to: consumer
    })
    // print price
    const price = utils.toUtf8String(hexPrice)
    const msg = price ? `current ETH price: ${price}` : 'No price listed'
    console.log(msg)
    callback()
  } catch (error) {
    console.error(
      'Usage: truffle exec scripts/view_eth_price.js [options] <contract address>'
    )
    callback(error)
  }
}
