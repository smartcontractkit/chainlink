require('@babel/register')({ cwd: '../../' })
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
      gas: 4700000
    }
  }
}
