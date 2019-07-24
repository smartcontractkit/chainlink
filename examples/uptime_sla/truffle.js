module.exports = {
  compilers: {
    solc: {
      version: '0.4.24'
    }
  },
  networks: {
    cldev: {
      host: '127.0.0.1',
      port: 18545,
      network_id: '*',
      gas: 4700000
    }
  }
}
