#!/usr/bin/env ts-node

import execa from 'execa'

const task = /^win/i.test(process.platform) ? 'build:windows' : 'build'

execa.commandSync(`yarn -s workspace chainlinkv0.5 run ${task}`, {
  stdio: 'inherit',
})
