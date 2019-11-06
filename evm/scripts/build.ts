#!/usr/bin/env node

import execa from 'execa'
import { ls } from 'shelljs'
import { writeFileSync } from 'fs'
import { join, resolve } from 'path'

const isWindows = /^win/i.test(process.platform)
const task = isWindows ? 'build:windows' : 'build'

execa.commandSync(`yarn -s workspace chainlink run ${task}`, {
  stdio: 'inherit',
})
if (isWindows) {
  remapAbi()
}
/**
 * Temporary interop for truffle compile on Windows
 * to flatten the compiler output into the JSON file
 */
function remapAbi() {
  const artifacts = join('dist', 'artifacts')
  const jsons = ls(artifacts)
  jsons.forEach(j => {
    const jsonPath = resolve(join(artifacts, j))
    const json = require(jsonPath)
    const { abi } = json
    const newJson = { compilerOutput: { abi }, ...json }
    writeFileSync(jsonPath, JSON.stringify(newJson))
  })
}
