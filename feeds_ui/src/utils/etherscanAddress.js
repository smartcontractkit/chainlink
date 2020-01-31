export const etherscanAddress = (networkId, contractAddress) => {
  switch (networkId) {
    case 3:
      return `https://ropsten.etherscan.io/address/${contractAddress}`
    default:
      return `https://etherscan.io/address/${contractAddress}`
  }
}
