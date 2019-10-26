#!/usr/bin/env node

const { exec } = require('child_process')

const task = /^win/i.test(process.platform) ? 'build:windows' : 'build'

const subprocess = exec(`yarn -s workspace chainlink run ${task}`, err => {
  if (err) {
    process.exit(1)
  }
})

subprocess.stderr.pipe(process.stderr)
subprocess.stdout.pipe(process.stdout)
