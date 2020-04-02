import core from '@actions/core'
// import github from '@actions/github'
import { exit, which } from 'shelljs'

if (!which('git')) {
  core.setFailed('Sorry, this script requires git')
  exit(1)
}

// const myToken = core.getInput('myToken')

// const octokit = new github.GitHub(myToken)
// workflow
// 1. on cron job, update cache and check if cache should use new base builder
// 2. push up new cache to dockerhub
// 3. modify workspace, updating all dockerfiles with newly pushed cache
// 4. create pr via https://github.com/peter-evans/create-pull-request
// 5. pr will contain new cache file, and new dockerfiles all in 1 PR
