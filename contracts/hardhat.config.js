require("@nomiclabs/hardhat-waffle");

/**
 * @type import('hardhat/config').HardhatUserConfig
 */
module.exports = {
  solidity: "0.8.4",
  paths: {
    artifacts: "./artifacts",
    cache: "./cache",
    sources: "./contracts",
    tests: "./test",
  },
  overrides: {
    "contracts/v0.8/*": {
      version: "0.8.4",
      optimizer: {
        enabled: true,
        runs: 200,
      },
    },
  },
};
