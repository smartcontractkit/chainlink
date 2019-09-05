module.exports = {
  compilers: {
    solc: {
      version: '0.4.24'
    }
  },
  networks: {
    test: {
      host: '127.0.0.1',
      port: process.env.ETH_HTTP_PORT || 18545,
      network_id: '*',
      gas: 4700000,
      gasPrice: 5e9
    }
  }
}
