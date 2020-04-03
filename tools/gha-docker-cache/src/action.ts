import core from '@actions/core'
import { exit, which } from 'shelljs'
import { updateCacheFile, updateDockerfiles } from './utils'

if (!which('git')) {
  core.setFailed('Sorry, this script requires git')
  exit(1)
}

enum ActionType {
  UPDATE_CACHE_FILE = 'UPDATE_CACHE_FILE',
  UPDATE_DOCKER_FILES = 'UPDATE_DOCKER_FILES',
}
const actionType = core.getInput('type', { required: true })

if (actionType === ActionType.UPDATE_CACHE_FILE) {
  updateCacheFile()
} else if (actionType === ActionType.UPDATE_DOCKER_FILES) {
  updateDockerfiles()
} else {
  core.setFailed(
    `Unrecognized action type, valid action types are: ${Object.values(
      ActionType,
    )}`,
  )
}
