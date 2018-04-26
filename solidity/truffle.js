global.h = require('./helpers');

module.exports = {
  network: "test",
  networks: {
    development: {
      host: "127.0.0.1",
      port: 18545,
      network_id: "*",
      gas: 4700000
    }
  }
};
