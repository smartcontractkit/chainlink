import { Environment } from '../config'
import { getVersion } from '../utils/version'

/**
 * Read a version file for displaying in production
 */
async function main() {
  const version = await getVersion(Environment.PROD)
  console.log(version)
}

main()
