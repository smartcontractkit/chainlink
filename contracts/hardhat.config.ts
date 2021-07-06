import "@nomiclabs/hardhat-waffle";
import "hardhat-contract-sizer";

const COMPILER_SETTINGS = {
  optimizer: {
    enabled: true,
    runs: 1000000,
  },
  metadata: {
    bytecodeHash: "none",
  },
};

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
        settings: COMPILER_SETTINGS,
      },
      {
        version: "0.5.0",
        settings: COMPILER_SETTINGS,
      },
      {
        version: "0.6.6",
        settings: COMPILER_SETTINGS,
      },
      {
        version: "0.7.6",
        settings: COMPILER_SETTINGS,
      },
      {
        version: "0.8.6",
        settings: COMPILER_SETTINGS,
      },
    ],
  },
  contractSizer: {
    alphaSort: true,
    runOnCompile: false,
    disambiguatePaths: false,
  },
};
