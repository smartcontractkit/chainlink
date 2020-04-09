import { Environment, ExplorerConfig } from '../config'
import { getVersion } from '../utils/version'

/**
 * Read a version file for displaying in production
 */
async function main() {
  const version = await getVersion({ env: Environment.PROD } as ExplorerConfig)
  console.log(version)
}

main()
