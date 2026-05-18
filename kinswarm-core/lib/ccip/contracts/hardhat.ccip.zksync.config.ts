import '@nomicfoundation/hardhat-ethers'
import '@nomicfoundation/hardhat-verify'
import '@nomicfoundation/hardhat-chai-matchers'
import '@matterlabs/hardhat-zksync-solc'
import '@typechain/hardhat'
import 'hardhat-abi-exporter'
import { subtask } from 'hardhat/config'
import { TASK_COMPILE_SOLIDITY_GET_SOURCE_PATHS } from 'hardhat/builtin-tasks/task-names'
import '@matterlabs/hardhat-zksync-verify'

const COMPILER_SETTINGS = {
  optimizer: {
    enabled: true,
    runs: 1000000,
  },
  metadata: {
    bytecodeHash: 'none',
  },
}

// prune forge style tests from hardhat paths
subtask(TASK_COMPILE_SOLIDITY_GET_SOURCE_PATHS).setAction(
  async (_, __, runSuper) => {
    const paths = await runSuper()
    const noTests = paths.filter((p: string) => !p.endsWith('.t.sol'))
    const noCCIPTests = noTests.filter(
      (p: string) => !p.includes('/v0.8/ccip/test'),
    )
    return noCCIPTests.filter(
      (p: string) => !p.includes('src/v0.8/vendor/forge-std'),
    )
  },
)

/**
 * @type import('hardhat/config').HardhatUserConfig
 */
let config = {
  abiExporter: {
    path: './abi',
    runOnCompile: true,
  },
  paths: {
    artifacts: './artifacts',
    cache: './cache',
    sources: './src/v0.8/ccip',
    tests: './test/v0.8/ccip, ./src/v0.8/ccip/test',
  },
  typechain: {
    outDir: './typechain',
    target: 'ethers-v5',
  },
  defaultNetwork: 'zkSync',
  networks: {
    env: {
      url: process.env.NODE_HTTP_URL || '',
    },
    hardhat: {
      allowUnlimitedContractSize: Boolean(
        process.env.ALLOW_UNLIMITED_CONTRACT_SIZE,
      ),
      hardfork: 'merge',
    },
    zkSyncSepolia: {
      url: 'https://sepolia.era.zksync.dev',
      ethNetwork: 'sepolia',
      zksync: true, // enables zksolc compiler
      verifyURL:
        'https://explorer.sepolia.era.zksync.dev/contract_verification',
    },
    zkSync: {
      url: 'https://mainnet.era.zksync.io', // The testnet RPC URL of ZKsync Era network.
      ethNetwork: 'mainnet', // The Ethereum Web3 RPC URL, or the identifier of the network (e.g. `mainnet` or `sepolia`)
      zksync: true,
      // Verification endpoint for Sepolia
      verifyURL:
        'https://zksync2-mainnet-explorer.zksync.io/contract_verification',
    },
  },
  solidity: {
    compilers: [
      {
        version: '0.8.24',
        settings: {
          ...COMPILER_SETTINGS,
          evmVersion: 'paris',
        },
      },
    ],
  },
  zksolc: {
    version: '1.5.3',
    settings: {
      optimizer: {
        enabled: true,
        mode: '3',
        fallback_to_optimizing_for_size: false,
      },
      experimental: {
        dockerImage: '',
        tag: '',
      },
      // contractsToCompile: ['RMN', 'ARMProxy'], // uncomment this to compile only specific contracts
    },
  },
  warnings: !process.env.HIDE_WARNINGS,
}

if (process.env.NETWORK_NAME && process.env.EXPLORER_API_KEY) {
  config = {
    ...config,
  }
}

export default config
