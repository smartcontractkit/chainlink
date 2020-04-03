import core from '@actions/core'
import { join } from 'path'
import { sed } from 'shelljs'
import { getGitRoot, getLatestName } from './utils'
/**
 * Update a docker cache file with the latest base file.
 *
 * @param name [cache.Dockerfile] The name of the cache file to update
 * @param baseRepo [smartcontract/builder] The docker repo of the base file
 */
export async function updateCacheFile(
  name = 'cache.Dockerfile',
  baseRepo = 'smartcontract/builder',
) {
  const path = getCacheFilePath(name)
  const latestBuilder = await getLatestName(baseRepo)
  core.info(
    `Updating cache file ${path} with builder version ${latestBuilder}...`,
  )

  const updated = sed('-i', includes(baseRepo), `FROM ${latestBuilder}`, [path])
  return updated
}

/**
 * Check that a string is included in a line, matches on first occurence.
 *
 * @param s The string to check for inclusion in a line
 */
export function includes(s: string, flags?: string) {
  return new RegExp(`^.*${s}.*$`, flags)
}

/**
 * Get the current dockerfile in the repository used as a cache file
 *
 * @param name The name of the dockerfile to search for within the root of the repository
 */
function getCacheFilePath(name: string) {
  const path = join(getGitRoot(), name)

  return path
}
