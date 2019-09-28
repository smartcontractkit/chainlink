#!/usr/bin/env node

const { exec } = require('child_process')

const task = /^win/i.test(process.platform)
  ? 'postinstall:windows'
  : 'postinstall:default'

exec(`yarn workspace chainlink run ${task}`, (err, stdout, stderr) => {
  if (err) {
    throw err
  }
  if (stdout) {
    console.log(stdout)
  }
  if (stderr) {
    console.error(stderr)
  }
})
