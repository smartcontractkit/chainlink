import { writeVersion } from '../utils/version'

/**
 * Write a version file for displaying in production
 */
async function main() {
  await writeVersion()
}

main()
