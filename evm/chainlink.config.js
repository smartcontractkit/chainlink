const ethers = require('ethers')

// Setup JSON RPC provider
const providerURL = process.env['ETH_HTTP_URL'] || 'http://localhost:18545'
const privateKey =
  process.env['PRIVATE_KEY'] ||
  '4d6cf3ce1ac71e79aa33cf481dedf2e73acb548b1294a70447c960784302d2fb'
const provider = new ethers.providers.JsonRpcProvider(providerURL)
const wallet = new ethers.Wallet(privateKey, provider)

module.exports = {
  devnetMiner: '0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f',
  privateKey,
  providerURL,
  provider,
  wallet
}