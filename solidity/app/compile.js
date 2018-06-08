const solc = require('solc')
const fs = require('fs')
const path = require('path')

const lookupPaths = [
  './',
  './contracts/',
  './node_modules/',
  './node_modules/linkToken/contracts/'
]

function solCompile (filename) {
  function lookupIncludeFile (filename) {
    for (let path of lookupPaths) {
      const fullPath = path + filename
      if (fs.existsSync(fullPath)) {
        return {contents: fs.readFileSync(fullPath).toString()}
      }
    }
    console.log('Unable to load', filename)
    return null
  }

  const inputBasename = path.basename(filename).toString()
  const solInput = {
    language: 'Solidity',
    sources: {[inputBasename]: {'urls': [filename]}},
    settings: {
      outputSelection: {
        [inputBasename]: {
          '*': [ 'abi', 'evm.bytecode' ]
        }
      }
    }
  }
  const output = solc.compileStandardWrapper(JSON.stringify(solInput), lookupIncludeFile)
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
  for (let [key, contract] of Object.entries(output.contracts)) {
    if (key === contractName) {
      return Object.values(contract)[0]
    }
  }
}

module.exports = (filename) => {
  const contractName = path.basename(filename).toString()
  const compiled = solCompile(filename)
  checkCompilerErrors(compiled.errors)
  return getContract(compiled, contractName)
}
