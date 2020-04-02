import core from '@actions/core'
import { join } from 'path'
import { exec, sed } from 'shelljs'
import { getGitRoot, getLatestName } from './utils'

/**
 * Update the repo's dockerfiles with the current cache file
 */
export async function updateDockerfiles(cacheRepo = 'smartcontract/cache') {
  const cache = await getLatestName(cacheRepo)
  const files = getDockerFiles(cacheRepo)

  files.forEach(({ path, text }) => {
    core.info(`Updating dockerfile ${path} from ${text} to ${cache}`)
    sed('-i', text, `FROM ${cache}`, [join(getGitRoot(), path)])
  })
}

/**
 * Split a string based on the first occurence of a colon
 *
 * @param s The string to split on
 */
export function splitOnColon(s: string) {
  const i = s.indexOf(':')

  return i < 0 ? [s] : [s.substring(0, i), s.substring(i + 1)]
}

/**
 * Get a list of dockerfiles that are used as cache images
 * within this repository.
 */
function getDockerFiles(cacheFileName: string) {
  const res = exec(`git grep "${cacheFileName}" -- "*Dockerfile*"`, {
    cwd: getGitRoot(),
  })

  return res
    .split('\n')
    .filter(Boolean)
    .map(splitOnColon)
    .map(([path, text]) => ({ path, text }))
}
