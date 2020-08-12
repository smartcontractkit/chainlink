import { Networks } from './networks'

export function etherscanAddress(networkId: Networks, contractAddress: string) {
  switch (networkId) {
    case Networks.ROPSTEN:
      return `https://ropsten.etherscan.io/address/${contractAddress}`
    case Networks.KOVAN:
      return `https://kovan.etherscan.io/address/${contractAddress}`
    default:
      return `https://etherscan.io/address/${contractAddress}`
  }
}
