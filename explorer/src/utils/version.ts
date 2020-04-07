import { readFileSync, writeFileSync } from 'fs'
import { ExplorerConfig } from '../config'
const VERSION_FILE_NAME = 'VERSION.json'

interface GitMeta {
  gitSha: string
  gitBranch: string
}
interface PkgVersions {
  serverVersion: string
  clientVersion: string
}
type VersionFile = GitMeta & PkgVersions

export async function fetchPkgVersions(): Promise<PkgVersions> {
  const clientPkg = await import('../../client/package.json')
  const serverPkg = await import('../../package.json')

  return { serverVersion: serverPkg.version, clientVersion: clientPkg.version }
}

export async function fetchMeta(): Promise<GitMeta> {
  const simplegit = await import('simple-git/promise')
  const g = simplegit.default()

  const gitSha = await g.revparse(['HEAD'])
  const { current: gitBranch } = await g.status()
  return { gitSha, gitBranch }
}

export async function writeVersion() {
  const version = await fetchVersion()
  writeFileSync(VERSION_FILE_NAME, JSON.stringify(version))
}

export function readVersion(): VersionFile {
  try {
    const file = readFileSync(VERSION_FILE_NAME, { encoding: 'utf-8' })
    return JSON.parse(file) as VersionFile
  } catch (e) {
    const origErr: Error = e
    const err = Error(`Could not read ${VERSION_FILE_NAME}: ${origErr.message}`)
    throw err
  }
}

export async function getVersion(conf: ExplorerConfig): Promise<VersionFile> {
  if (conf.dev) {
    return await fetchVersion()
  }
  return readVersion()
}

async function fetchVersion() {
  const packageVersions = await fetchPkgVersions()
  const gitMeta = await fetchMeta()
  return { ...packageVersions, ...gitMeta }
}
