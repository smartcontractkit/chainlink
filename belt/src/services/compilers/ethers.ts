import { join } from 'path'
import { tsGenerator } from 'ts-generator'
import { TypeChain } from 'typechain/dist/TypeChain'
import * as config from '../config'
import { getArtifactDirs } from '../utils'

/**
 * Generate ethers.js contract abstractions for all of the solidity versions under a specified contract
 * directory.
 *
 * @param conf The application configuration, e.g. where to read artifacts, where to output, etc..
 */
export async function compileAll(conf: config.App) {
  const cwd = process.cwd()

  return Promise.all(
    getArtifactDirs(conf).map(async ({ dir }) => {
      const c = compiler(conf, cwd, dir)
      await tsGenerator({ cwd, loggingLvl: 'verbose' }, c)
    }),
  )
}

/**
 * Create a typechain compiler instance that reads in a subdirectory of artifacts e.g. (abi/v0.4, abi/v0.5.. etc)
 * and outputs ethers contract abstractions under the same version prefix, (ethers/v0.4, ethers/v0.4.. etc)
 *
 * @param config The application level config for compilation
 * @param cwd The current working directory during this programs execution
 * @param subDir The subdirectory to use as a namespace when reading artifacts and outputting
 * contract abstractions
 */
function compiler(
  { artifactsDir, contractAbstractionDir }: config.App,
  cwd: string,
  subDir: string,
): TypeChain {
  return new TypeChain({
    cwd,
    rawConfig: {
      files: join(artifactsDir, subDir, '**', '*.json'),
      outDir: join(contractAbstractionDir, 'ethers', subDir),
      target: 'ethers',
    },
  })
}
