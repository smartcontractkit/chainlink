import "@nomiclabs/hardhat-waffle";
import "hardhat-contract-sizer";

/**
 * @type import('hardhat/config').HardhatUserConfig
 */
export default {
  paths: {
    artifacts: "./artifacts",
    cache: "./cache",
    sources: "./src",
    tests: "./test",
  },
  networks: {
    hardhat: {
      allowUnlimitedContractSize: true,
    },
  },
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
  contractSizer: {
    alphaSort: true,
    runOnCompile: true,
    disambiguatePaths: false,
  }
};
