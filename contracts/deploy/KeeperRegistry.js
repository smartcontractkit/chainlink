module.exports = async ({ getNamedAccounts, deployments }) => {
  const { deployer, linkToken, linkEth, fastGas } = await getNamedAccounts()
  const paymentPremiumPPB = 250000000
  const blockCountPerTurn = 3
  const maxCheckGas = 20000000
  const stalenessSeconds = 43820
  const gasCeilingMultiplier = 1
  const fallbackGasPrice = 200
  const fallbackLinkPrice = 200000000

  await deployments.deploy('KeeperRegistry', {
    from: deployer,
    log: true,
    contract: 'KeeperRegistry',
    args: [
      linkToken,
      linkEth,
      fastGas,
      paymentPremiumPPB,
      blockCountPerTurn,
      maxCheckGas,
      stalenessSeconds,
      gasCeilingMultiplier,
      fallbackGasPrice,
      fallbackLinkPrice,
    ]
  })
}
