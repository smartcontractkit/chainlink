import { readFileSync, writeFileSync } from 'fs'
import { Environment } from '../config'

/**
 * The name of the file to write to the root of this package, which contains version data.
 */
export const VERSION_FILE_NAME = 'VERSION.json'

/**
 * Git data describing the current commit sha and current branch
 */
interface GitMeta {
  gitSha: string
  gitBranch: string
}

/**
 * Versions of the explorer client and server
 */
interface PkgVersions {
  serverVersion: string
  clientVersion: string
}

/**
 * Abstract representation of a version file, contains data like client/server version,
 * and git metadata
 */
export type VersionFile = GitMeta & PkgVersions

/**
 * Get the current version data.
 *
 * If in production mode, will try to read version from a local file.
 * Else, the version data will be fetched
 *
 * @param env Environment
 */
export async function getVersion(env: Environment): Promise<VersionFile> {
  if (env === Environment.PROD) {
    return readVersion()
  }

  return await fetchVersion()
}

/**
 * Write the current version into a file in the root of this package
 */
export async function writeVersion() {
  const version = await fetchVersion()
  writeFileSync(VERSION_FILE_NAME, JSON.stringify(version))
}

/**
 * Read the current version file from the root of this package
 */
function readVersion(): VersionFile {
  try {
    const file = readFileSync(VERSION_FILE_NAME, { encoding: 'utf-8' })
    return JSON.parse(file) as VersionFile
  } catch (e) {
    const origErr: Error = e
    const err = Error(`Could not read ${VERSION_FILE_NAME}: ${origErr.message}`)
    throw err
  }
}

/**
 * Fetch package versions by reading from package.json
 */
async function fetchPkgVersions(): Promise<PkgVersions> {
  // @ts-ignore
  const clientPkg = await import('../../client/package.json')
  const serverPkg = await import('../../package.json')

  return { serverVersion: serverPkg.version, clientVersion: clientPkg.version }
}

/**
 * Fetch git meta data from the current repository
 */
async function fetchMeta(): Promise<GitMeta> {
  const simplegit = await import('simple-git/promise')
  const g = simplegit.default()

  const gitSha = await g.revparse(['HEAD'])
  const { current: gitBranch } = await g.status()
  return { gitSha, gitBranch }
}

/**
 * Fetch the current version
 */
async function fetchVersion(): Promise<VersionFile> {
  const packageVersions = await fetchPkgVersions()
  const gitMeta = await fetchMeta()
  return { ...packageVersions, ...gitMeta }
}
