let LinkToken = artifacts.require('LinkToken')
let Oracle = artifacts.require('Oracle')

const from = "0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f"

module.exports = async (deployer, network) => {
  await deployer.deploy(LinkToken, { from })
  console.log(`Deployed LinkToken at: ${LinkToken.address}`)
  const oracle = await deployer.deploy(Oracle, { from })
  console.log(`Deployed Oracle at: ${Oracle.address}`)
  oracle.setFulfillmentPermission(chainlinkAddress, true, { from })
}
