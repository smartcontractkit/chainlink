import { echo, exit, which } from 'shelljs'
// import core from '@actions/core'
// import github from '@actions/github'

if (!which('git')) {
  echo('Sorry, this script requires git')
  exit(1)
}
