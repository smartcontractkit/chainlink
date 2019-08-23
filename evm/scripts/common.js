const ethers = require('ethers')

// Setup provider & wallet
const port = process.env.ETH_HTTP_PORT || `18545`
const providerURL = process.env['ETH_HTTP_URL'] || `http://localhost:${port}`
const privateKey =
  process.env['PRIVATE_KEY'] ||
  '4d6cf3ce1ac71e79aa33cf481dedf2e73acb548b1294a70447c960784302d2fb'
const provider = new ethers.providers.JsonRpcProvider(providerURL)
const wallet = new ethers.Wallet(privateKey, provider)

// Devnet miner address
const DEVNET_ADDRESS = '0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f'

// script arguments for command-line-args
const optionDefinitions = [
  { name: 'args', type: String, multiple: true, defaultOption: true },
  { name: 'compile', type: Boolean },
  { name: 'network', type: String }
]

// wrapper for main truffle script functions
const scriptRunner = (main, usage) => async callback => {
  try {
    await main()
    callback()
  } catch (error) {
    console.log(`Usage: ${usage}`)
    callback(error)
  }
}

module.exports = {
  DEVNET_ADDRESS,
  optionDefinitions,
  privateKey,
  providerURL,
  provider,
  scriptRunner,
  wallet
}
