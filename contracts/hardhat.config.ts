import "@nomiclabs/hardhat-waffle";
import "@nomiclabs/hardhat-truffle5";
import "hardhat-contract-sizer";
import "hardhat-abi-exporter";
import "hardhat-gas-reporter";
import "solidity-coverage";
import "hardhat-deploy";

const KOVAN_RPC_URL = process.env.KOVAN_RPC_URL || 'http://localhost:8545'
const KOVAN_PRIVATE_KEY = process.env.KOVAN_PRIVATE_KEY || '0x00'
const MUMBAI_RPC_URL = process.env.MUMBAI_RPC_URL || 'http://localhost:8545'
const MUMBAI_PRIVATE_KEY = process.env.MUMBAI_PRIVATE_KEY || '0x00'
const BSCTESTNET_RPC_URL = process.env.BSCTESTNET_RPC_URL || 'http://localhost:8545'
const BSCTESTNET_PRIVATE_KEY = process.env.BSCTESTNET_PRIVATE_KEY || '0x00'
const DEPLOYER = process.env.DEPLOYER || 0
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
  abiExporter: {
    path: "./abi",
  },
  paths: {
    artifacts: "./artifacts",
    cache: "./cache",
    sources: "./src",
    tests: "./test",
  },
  networks: {
    hardhat: {},
    kovan: {
      url: KOVAN_RPC_URL,
      accounts: [KOVAN_PRIVATE_KEY],
      chainId: 42,
    },
    mumbai: {
      url: MUMBAI_RPC_URL,
      accounts: [MUMBAI_PRIVATE_KEY],
      chainId: 80001,
    },
    bsctestnet: {
      url: BSCTESTNET_RPC_URL,
      accounts: [BSCTESTNET_PRIVATE_KEY],
      chainId: 97,
    }
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
  namedAccounts: {
    linkToken: {
      1: '0x514910771AF9Ca656af840dff83E8264EcF986CA',
      42: '0xa36085F69e2889c224210F603D836748e7dC0088',
      80001: '0x326C977E6efc84E512bB9C30f76E30c160eD06FB',
      97: '0x84b9b910527ad5c03a9ca831909e21e236ea7b06',
    },
    linkEth: {
      1: '0xDC530D9457755926550b59e8ECcdaE7624181557',
      42: '0x3Af8C569ab77af5230596Acf0E8c2F9351d24C38',
      80001: '0xc0FAb0a0c9204ae4682eFdca3F05EAAb17440271', // FIXME: Replace to the real feed smart contract deployed on Mumbai. The provided one is the mock.
      97: '0xBF44C29A52dF268841f7C689F73A5ec6dc6e6409', // FIXME: Replace to the real feed smart contract deployed on BSC Testnet. The provided one is the mock.
    },
    fastGas: {
      1: '0x169E633A2D1E6c10dD91238Ba11c4A708dfEF37C',
      42: '0x73B9b95a2AE128225dbE53A7451B6c97e3De6F08',
      80001: '0xc0FAb0a0c9204ae4682eFdca3F05EAAb17440271', // FIXME: Replace to the real feed smart contract deployed on Mumbai. The provided one is the mock.
      97: '0xBF44C29A52dF268841f7C689F73A5ec6dc6e6409', // FIXME: Replace to the real feed smart contract deployed on BSC Testnet. The provided one is the mock.
    },
    deployer: {
      default: DEPLOYER
    }
  },
  gasReporter: {
    currency: 'USD',
    gasPrice: 10
  },
  contractSizer: {
    alphaSort: true,
    runOnCompile: false,
    disambiguatePaths: false,
  },
  mocha: {
    timeout: 100000,
  },
};
