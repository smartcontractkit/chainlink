const ethTX = require('ethereumjs-tx')

// Burn a tiny amount of eth
const txData = {
  nonce: '0x00',
  gasPrice: '0x09184e72a000',
  gasLimit: '0x2710',
  to: '0x0000000000000000000000000000000000000000',
  value: '0x1'
}

const tx = new ethTX.Transaction(txData)

// The secret key for the public key used in chainlink dev mode,
// 0x9ca9d2d5e04012c9ed24c0e513c9bfaa4a2dd77f
const privKey = Buffer.from(
  '34d2ee6c703f755f9a205e322c68b8ff3425d915072ca7483190ac69684e548c', 'hex')
tx.sign(privKey)
console.log('0x' + tx.serialize().toString('hex'))
