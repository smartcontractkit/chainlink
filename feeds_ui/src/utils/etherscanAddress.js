export const etherscanAddress = (networkName, contractAddress) => {
  switch (networkName) {
    case 'ropsten':
      return `https://ropsten.etherscan.io/address/${contractAddress}`
    default:
      return `https://etherscan.io/address/${contractAddress}`
  }
}
