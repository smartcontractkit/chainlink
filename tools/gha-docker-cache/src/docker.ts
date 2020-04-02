import drc from 'docker-registry-client'
import { join } from 'path'
import semver from 'semver'
import { cat, exec, sed } from 'shelljs'
updateDockerfiles('henrynguyen5/base')

/**
 * Update the repo's dockerfiles with the current cache file
 */
async function updateDockerfiles(cacheRepo = 'smartcontract/cache') {
  const cache = await getLatestName(cacheRepo)
  const files = getDockerFiles(cacheRepo)
  files.forEach(({ path, text }) => {
    sed(text, `FROM ${cache}`, [join(getGitRoot(), path)])
  })
}
/**
 * Get a list of dockerfiles that are used as cache images
 * within this repository.
 */
export function getDockerFiles(cacheFileName: string) {
  const res = exec(`git grep ${cacheFileName}`, { cwd: getGitRoot() })

  return res
    .split('\n')
    .filter(Boolean)
    .map(splitOnColon)
    .map(([path, text]) => ({ path, text }))
}

/**
 * Split a string based on the first occurence of a colon
 *
 * @param s The string to split on
 */
function splitOnColon(s: string) {
  const i = s.indexOf(':')
  return [s.substring(0, i), s.substring(i + 1)]
}

/**
 * Update a docker cache file with the latest base file.
 *
 * @param name The name of the cache file to update
 * @param baseRepo The docker repo of the base file
 */
export async function updateCacheFile(
  name = 'base.Dockerfile',
  baseRepo = 'smartcontract/builder',
) {
  const current = getCacheFile(name)
  const latestBuilder = await getLatestName(baseRepo)

  const updated = current
    .split('\n')
    .map(s => (s.includes(baseRepo) ? latestBuilder : s))
    .join('\n')
  return updated
}

/**
 * Get the current dockerfile in the repository used as a cache file
 * @param name The name of the dockerfile to search for within the root of the repository
 */
export function getCacheFile(name: string) {
  const file = cat(join(getGitRoot(), name)).toString()

  return file
}

/**
 * Get the absolute path of the root of the repository
 */
function getGitRoot() {
  return exec('git rev-parse --show-toplevel', { silent: true }).trim()
}

/**
 * Get the latest base image name from the docker registry.
 * Only handles valid semver tags, invalid semver tags will be ignored.
 */
export async function getLatestName(repo: string) {
  const client = drc.createClientV2({
    name: repo,
  })

  const { tags } = await listTags(client)
  const filteredTags = tags
    .filter(t => !!semver.valid(t))
    .sort((a, b) => semver.compare(a, b))

  const latestTag = filteredTags[filteredTags.length - 1]
  const latestName = `${repo}:${latestTag}`
  return latestName
}

/**
 * List all tags for a repository on the official docker registry
 *
 * @param repo The repository to list tags for
 */
function listTags(client: drc.RegistryClientV2): Promise<drc.Tags> {
  return new Promise((resolve, reject) => {
    client.listTags((err, tags) => {
      client.close()
      if (err) {
        return reject(err)
      }
      return resolve(tags)
    })
  })
}
