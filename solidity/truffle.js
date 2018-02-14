global.h = require('./helpers');

module.exports = {
  networks: {
    devnet: {
      host: "127.0.0.1",
      port: 18545,
      network_id: "*",
      gas: 4700000
    }
  }
};
