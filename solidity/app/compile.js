const solc = require('solc')
const fs = require('fs')
const path = require('path')

let lookupPaths = ['./', './node_modules/']

function solCompile (filename) {
  function lookupIncludeFile (filename) {
    for (let path of lookupPaths) {
      let fullPath = path + filename
      if (fs.existsSync(fullPath)) {
        return {contents: fs.readFileSync(fullPath).toString()}
      }
    }
    console.log('Unable to load', filename)
    return null
  }

  let inputBasename = path.basename(filename).toString()
  let input = {[inputBasename]: {'urls': [filename]}}

  let solInput = {
    language: 'Solidity',
    sources: input,
    settings: {
      outputSelection: {
        [inputBasename]: {
          '*': [ 'abi', 'evm.bytecode' ]
        }
      }
    }
  }
  let output = solc.compileStandardWrapper(JSON.stringify(solInput), lookupIncludeFile)
  return JSON.parse(output)
}

function checkCompilerErrors (errors) {
  if (errors == null) {
    return
  }

  let failure = false
  for (let error of errors) {
    if (error.type !== 'Warning') {
      console.log(error.sourceLocation.file + ':' + error.sourceLocation.start + ' - ' + error.message)
      failure = true
    }
  }
  if (failure) {
    console.log('Aborting because of errors.')
    process.exit(1)
  }
}

function getContract (output, contractName) {
  for (const [key, module] of Object.entries(output.contracts)) {
    if (key === contractName) {
      for (const [_, contract] of Object.entries(module)) {
        return contract
      }
    }
  }
}

module.exports = function compile (filename) {
  let contractName = path.basename(filename).toString()
  let compiled = solCompile(filename)
  checkCompilerErrors(compiled.errors)
  return getContract(compiled, contractName)
}
