#!/usr/bin/env node

const { exec } = require('child_process')

const task = /^win/i.test(process.platform) ? 'build:windows' : 'build'

const subprocess = exec(`yarn -s workspace chainlinkv0.5 run ${task}`, err => {
  if (err) {
  }
})

subprocess.stderr.pipe(process.stderr)
subprocess.stdout.pipe(process.stdout)
