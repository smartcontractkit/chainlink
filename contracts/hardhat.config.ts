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
    hardhat: {},
  },
  solidity: {
    compilers: [
      {
        version: "0.4.24",
        settings: {
          optimizer: {
            enabled: true,
            runs: 1000000,
          },
          metadata: {
            bytecodeHash: "none",
          },
        },
      },
      {
        version: "0.5.0",
        settings: {
          optimizer: {
            enabled: true,
            runs: 1000000,
          },
          metadata: {
            bytecodeHash: "none",
          },
        },
      },
      {
        version: "0.6.6",
        settings: {
          optimizer: {
            enabled: true,
            runs: 1000000,
          },
          metadata: {
            bytecodeHash: "none",
          },
        },
      },
      {
        version: "0.7.6",
        settings: {
          optimizer: {
            enabled: true,
            runs: 1000000,
          },
          metadata: {
            bytecodeHash: "none",
          },
        },
      },
      {
        version: "0.8.4",
        settings: {
          optimizer: {
            enabled: true,
            runs: 1000000,
          },
          metadata: {
            bytecodeHash: "none",
          },
        },
      },
    ],
  },
  contractSizer: {
    alphaSort: true,
    runOnCompile: true,
    disambiguatePaths: false,
  },
};
