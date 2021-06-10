import "@nomiclabs/hardhat-waffle";

/**
 * @type import('hardhat/config').HardhatUserConfig
 */
export default {
  solidity: {
    compilers: [
      {
        version: "0.4.24",
        optimizer: {
          enabled: true,
          runs: 1000000,
        },
      },
      {
        version: "0.5.0",
        optimizer: {
          enabled: true,
          runs: 1000000,
        },
      },
      {
        version: "0.6.6",
        optimizer: {
          enabled: true,
          runs: 1000000,
        },
      },
      {
        version: "0.7.6",
        optimizer: {
          enabled: true,
          runs: 1000000,
        },
      },
      {
        version: "0.8.4",
        optimizer: {
          enabled: true,
          runs: 1000000,
        },
      },
    ],
  },
  paths: {
    artifacts: "./artifacts",
    cache: "./cache",
    sources: "./src",
    tests: "./test/v0.8/vrf",
  },
  overrides: {
    "src/v0.8/*": {
      version: "0.8.4",
      optimizer: {
        enabled: true,
        runs: 200,
      },
    },
    "src/v0.4/*": {
      version: "0.4.24",
      optimizer: {
        enabled: true,
        runs: 200,
      },
    },
  },
};
