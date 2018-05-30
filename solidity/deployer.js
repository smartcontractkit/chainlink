const deploy = require('./app/deploy.js')

if(process.argv.length != 3) {
  console.error('Usage: node deployer.js <solidity contract>')
  process.exit(1)
}

let filePath = process.argv[2]
deploy(filePath).then(address => {
  console.log(`${filePath} successfully deployed: ${address}`)
})
