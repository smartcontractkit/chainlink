import axios from 'axios'

export default class Etherscan {
  chainId: number
  baseUrl: string
  apiKey: string

  constructor(chainId: number, apiKey: string) {
    this.chainId = chainId
    this.baseUrl = `https://api${getEtherscanDomain(chainId)}.etherscan.io/api`
    this.apiKey = apiKey
  }

  async isVerified(contractAddress: string): Promise<boolean> {
    const res = await axios.get(this.baseUrl, {
      params: {
        module: 'contract',
        action: 'getsourcecode',
        address: contractAddress,
        apikey: this.apiKey,
      },
    })
    const isVerified = res.data.result[0].SourceCode
    return isVerified ? true : false
  }
}

/**
 * Returns the etherscan API domain for a given chainId.
 *
 * @param chainId Ethereum chain ID
 */
function getEtherscanDomain(chainId: number): string {
  const domains: { [keyof: number]: string } = {
    1: '',
    3: '-ropsten',
    4: '-rinkeby',
    42: '-kovan',
  }
  const idNotFound = !Object.keys(domains).includes(chainId.toString())
  if (idNotFound) {
    throw new Error('Invalid chain Id')
  }

  return domains[chainId]
}
