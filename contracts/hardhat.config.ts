import '@nomiclabs/hardhat-ethers'
import '@nomiclabs/hardhat-etherscan'
import '@nomiclabs/hardhat-waffle'
import '@typechain/hardhat'
import 'hardhat-abi-exporter'
import 'hardhat-contract-sizer'
import 'solidity-coverage'

const COMPILER_SETTINGS = {
  optimizer: {
    enabled: true,
    runs: 1000000,
  },
  metadata: {
    bytecodeHash: 'none',
  },
}

/**
 * @type import('hardhat/config').HardhatUserConfig
 */
export default {
  abiExporter: {
    path: './abi',
  },
  paths: {
    artifacts: './artifacts',
    cache: './cache',
    sources: './src',
    tests: './test',
  },
  typechain: {
    outDir: './typechain',
    target: 'ethers-v5',
  },
  networks: {
    hardhat: {},
    goerli: {
      url: 'https://link-eth.getblock.io/goerli/archive/axej8woh-seej-6ash-4Yu7-eyib1495dhno/',
      accounts: ["ce44819c415dcfa03ee6b124a65435eb774c109dc0b1a4386eb11a276a3e22f8"],
      saveDeployments: true,
      chainId: 5,
    }
  },
  solidity: {
    compilers: [
      {
        version: '0.4.24',
        settings: COMPILER_SETTINGS,
      },
      {
        version: '0.5.0',
        settings: COMPILER_SETTINGS,
      },
      {
        version: '0.6.6',
        settings: COMPILER_SETTINGS,
      },
      {
        version: '0.7.6',
        settings: COMPILER_SETTINGS,
      },
      {
        version: '0.8.6',
        settings: COMPILER_SETTINGS,
      },
      {
        version: '0.8.15',
        settings: COMPILER_SETTINGS,
      },
    ],
  },
  contractSizer: {
    alphaSort: true,
    runOnCompile: false,
    disambiguatePaths: false,
  },
  mocha: {
    timeout: 100000,
  },
}
