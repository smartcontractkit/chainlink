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
  const [link, holder] = options.args.slice(2)
  try {
    // encode function call
    const funcSelector = '0x70a08231'// "balanceOf(address)"
    const encodedParams = utils.defaultAbiCoder.encode(['address'], [holder])
    const data = utils.hexlify(utils.concat([funcSelector, encodedParams]))
    // make function call
    const hexBalance = await provider.call({
      data,
      to: link
    })
    // print balance
    const balance = utils.bigNumberify(hexBalance)
    console.log(`LINK balance: ${balance.toString()}`)
    callback()
  } catch (error) {
    console.error('Usage: truffle exec scripts/check_link_balance.js [options] ' +
    '<token address> <holder address>')
    callback(error)
  }
}
