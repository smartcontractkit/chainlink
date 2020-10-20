import { ContractObject } from '@truffle/contract-schema'
import { writeFileSync } from 'fs'
import { basename, join } from 'path'
import { ls, mkdir } from 'shelljs'
import * as config from '../config'
import { getArtifactDirs, getJsonFile } from '../utils'

/**
 * Generate @truffle/contract abstractions for all of the solidity versions under a specified contract
 * directory.
 *
 * @param conf The application configuration, e.g. where to read artifacts, where to output, etc..
 */
export async function compileAll(conf: config.App) {
  getArtifactDirs(conf).forEach(({ dir }) => {
    getContractPathsPer(conf, dir).forEach((p) => {
      const json: any = getJsonFile(p)
      const fileName = basename(p, '.json')
      const file = fillTemplate(fileName, {
        contractName: json.contractName,
        abi: json.compilerOutput.abi,
        evm: json.compilerOutput.evm,
        metadata: json.compilerOutput.metadata,
      })

      write(join(conf.contractAbstractionDir, 'truffle', dir), fileName, file)
    })
  })
}

/**
 * Create a truffle contract abstraction file
 *
 * @param contractName The name of the contract that will be exported
 * @param contractArgs The arguments to pass to @truffle/contract
 */
function fillTemplate(
  contractName: string,
  contractArgs: ContractObject,
): string {
  return `'use strict'
Object.defineProperty(exports, '__esModule', { value: true })
const contract = require('@truffle/contract')
const ${contractName} = contract(${JSON.stringify(contractArgs, null, 1)})

if (process.env.NODE_ENV === 'test') {
  try {
    eval('${contractName}.setProvider(web3.currentProvider)')
  } catch (e) {}
}

exports.${contractName} = ${contractName}
`
}

function getContractPathsPer({ artifactsDir }: config.App, version: string) {
  return [...ls(join(artifactsDir, version, '/**/*.json'))]
}

function write(outPath: string, fileName: string, file: string) {
  mkdir('-p', outPath)
  writeFileSync(join(outPath, `${fileName}.js`), file)
}
