import { Compiler } from '@krebernisak/ovm-compiler'
import * as config from '../config'
import { getContractDirs } from '../utils'
import { getCompilerOptions } from './solc'

/**
 * Generate solidity artifacts for all of the solidity versions under a specified contract
 * directory.
 *
 * @param conf The application configuration, e.g. where to read solidity files, where to output, etc..
 */
export async function compileAll(conf: config.App) {
  return Promise.all(
    getContractDirs(conf).map(async ({ dir, version }) => {
      const opts = getCompilerOptions(conf, dir, version)

      const c = new Compiler({
        ...opts,
        // Update version string to be detected by forked 0x/sol-compiler
        solcVersion: opts.solcVersion + '_ovm',
        isOfflineMode: true,
      })

      // Compiler#getCompilerOutputsAsync throws on compilation errors
      // this method prints any errors and warnings for us
      await c.compileAsync()
    }),
  )
}
