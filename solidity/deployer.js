require('./app/cl_utils.js')
const deploy = require('./app/deploy.js')

if (process.argv.length < 3) {
  console.error('Usage: node deployer.js <solidity contract> <constructor args...>')
  process.exit(1)
}

const filePath = process.argv[2]
const args = process.argv.slice(3)
deploy(filePath, ...args).then(contract => {
  console.log(`${filePath} successfully deployed: ${contract.address}`)
})
