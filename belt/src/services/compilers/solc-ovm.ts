import { Compiler, CompilerOptions } from '@krebernisak/ovm-compiler'
import { join } from 'path'
import * as config from '../config'
import { debug, getContractDirs } from '../utils'
const d = debug('solc')

/**
 * Generate solidity artifacts for all of the solidity versions under a specified contract
 * directory.
 *
 * @param conf The application configuration, e.g. where to read solidity files, where to output, etc..
 */
export async function compileAll(conf: config.App) {
  return Promise.all(
    getContractDirs(conf).map(async ({ dir, version }) => {
      const c = compiler(conf, dir, version)

      // Compiler#getCompilerOutputsAsync throws on compilation errors
      // this method prints any errors and warnings for us
      await c.compileAsync()
    }),
  )
}

/**
 * Create a sol-compiler instance that reads in a subdirectory of smart contracts e.g. (src/v0.4, src/v0.5, ..)
 * and outputs their respective compiler artifacts e.g. (abi/v0.4, abi/v0.5)
 *
 * @param config The application specific configuration to use for sol-compiler
 * @param subDir The subdirectory to use as a namespace when reading .sol files and outputting
 * their respective artifacts
 * @param solcVersion The solidity compiler version to use with sol-compiler
 */
function compiler(
  {
    artifactsDir,
    useDockerisedSolc,
    contractsDir,
    compilerSettings,
  }: config.App,
  subDir: string,
  solcVersion: string,
) {
  const _d = d.extend('compiler')
  // remove our custom versions property
  const compilerSettingCopy: any = JSON.parse(JSON.stringify(compilerSettings))
  delete compilerSettingCopy.versions

  // Update version string to be detected by forked 0x/sol-compiler
  solcVersion += '_ovm'

  const settings: CompilerOptions = {
    artifactsDir: join(artifactsDir, subDir),
    compilerSettings: {
      outputSelection: {
        '*': {
          '*': [
            'abi',
            'devdoc',
            'userdoc',
            'evm.bytecode.object',
            'evm.bytecode.sourceMap',
            'evm.deployedBytecode.object',
            'evm.deployedBytecode.sourceMap',
            'evm.methodIdentifiers',
            'metadata',
          ],
        },
      },
      ...compilerSettingCopy,
    },
    contracts: '*',
    contractsDir: join(contractsDir, subDir),
    isOfflineMode: true,
    shouldSaveStandardInput: false,
    solcVersion,
    useDockerisedSolc,
  }
  _d('Settings: %o', settings)

  return new Compiler(settings)
}
