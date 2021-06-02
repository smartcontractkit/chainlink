import "@nomiclabs/hardhat-waffle";

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
  solidity: {
    compilers: [
      {
        version: "0.6.6",
        optimizer: {
          enabled: true,
          runs: 200,
        }
      },
      {
        version: "0.7.6",
        optimizer: {
          enabled: true,
          runs: 200,
        }
      },
      {
        version: "0.8.4",
        optimizer: {
          enabled: true,
          runs: 200,
        }
      },
    ]
  },
};
