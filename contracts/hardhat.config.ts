import "@nomiclabs/hardhat-waffle";

/**
 * @type import('hardhat/config').HardhatUserConfig
 */
export default {
  solidity: "0.8.4",
  paths: {
    artifacts: "./artifacts",
    cache: "./cache",
    sources: "./src",
    tests: "./test",
  },
  overrides: {
    "src/v0.8/*": {
      version: "0.8.4",
      optimizer: {
        enabled: true,
        runs: 200,
      },
    },
  },
};
