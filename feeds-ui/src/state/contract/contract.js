const ethers = require('ethers')
const HeartbeatContractAbi = require('contracts/Heartbeat.json')

const provider = new ethers.providers.JsonRpcProvider(
  process.env.REACT_APP_ETHEREUM_NETWORK
)
provider.pollingInterval = 8000

const HeartbeatContract = new ethers.Contract(
  '0x79fEbF6B9F76853EDBcBc913e6aAE8232cFB9De9',
  HeartbeatContractAbi,
  provider
)

export { provider, HeartbeatContract }
