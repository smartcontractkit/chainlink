require('@babel/register')
require('@babel/polyfill')

module.exports = {
  networks: {
    cldev: {
      host: '127.0.0.1',
      port: 18545,
      network_id: '*',
      gas: 4700000,
    },
  },
  compilers: {
    solc: {
      version: '0.4.24',
    },
  },
}
