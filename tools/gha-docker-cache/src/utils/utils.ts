import core from '@actions/core'
import drc from 'docker-registry-client'
import semver from 'semver'
import { exec } from 'shelljs'

/**
 * Get the absolute path of the root of the repository
 */
export function getGitRoot() {
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
  core.info(`Fetched latest tag for repo ${repo}: ${latestTag}`)

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
