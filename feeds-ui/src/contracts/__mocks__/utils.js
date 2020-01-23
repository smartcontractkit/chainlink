import { ethers } from 'ethers'

export class Contract {
  constructor(address, abi) {
    abi.map(abiItem => {
      this[abiItem.name] = jest.fn()
    })
  }
}

export function createInfuraProvider() {
  return new ethers.providers.JsonRpcProvider(`mockProvider`)
}

export function createContract(address, provider, abi) {
  return new Contract(address, abi, provider)
}

export function formatAnswer() {
  return jest.fn()
}
