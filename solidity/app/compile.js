const solc = require('solc')
const fs = require('fs')
const path = require('path')

const INCLUDE_PATHS = ['./']

function solidityCompile (filename, lookupIncludeFile) {
  const inputBasename = path.basename(filename).toString()
  const input = {
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

  const solOutput = solc.compileStandardWrapper(JSON.stringify(input), lookupIncludeFile)
  return JSON.parse(solOutput)
}

function environmentPaths () {
  return (process.env['SOLIDITY_INCLUDE'] || '').split(/:/)
}

function compile (filename) {
  // Include the relative path of the specified file in the lookup paths
  let lookupPaths = INCLUDE_PATHS.concat(path.dirname(filename))

  // Add any paths defined in the environment
  lookupPaths = lookupPaths.concat(environmentPaths())

  function lookupIncludeFile (includeFile) {
    for (let lookupPath of lookupPaths) {
      // Do all path lookups relative to the script
      const fullPath = path.resolve(__dirname, lookupPath, includeFile)
      if (fs.existsSync(fullPath)) {
        return {contents: fs.readFileSync(fullPath).toString()}
      }
    }

    throw new Error(`Unable to load ${includeFile} searched in ${lookupPaths.join(' ')}`)
  }

  return solidityCompile(filename, lookupIncludeFile)
}

function checkCompilerErrors (errors) {
  if (errors == null) {
    return
  }

  let failure = false
  for (let error of errors) {
    if (error.type !== 'Warning') {
      console.log(error.formattedMessage)
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
  const compiled = compile(filename)
  checkCompilerErrors(compiled.errors)
  return getContract(compiled, contractName)
}
