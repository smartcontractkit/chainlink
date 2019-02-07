require('@babel/register')
require('@babel/polyfill')

module.exports = {
  network: 'test',
  compilers: {
    solc: {
      version: '0.4.24'
    }
  },
  networks: {
    development: {
      host: '127.0.0.1',
      port: 18545,
      network_id: '*',
      gas: 4700000,
      gasPrice: 5e9
    },
    ropsten: {
      host: 'localhost',
      port: 8545,
      gas: 5000000,
      gasPrice: 5e9,
      network_id: '3'
    },
    rinkeby: {
      host: 'localhost',
      port: 28545,
      gas: 5000000,
      gasPrice: 5e9,
      network_id: '4'
    }
  }
}
