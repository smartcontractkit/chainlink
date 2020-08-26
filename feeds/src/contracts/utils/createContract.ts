import { ethers } from 'ethers'
import { FunctionFragment } from 'ethers/utils'
import { JsonRpcProvider } from 'ethers/providers'

/**
 * Connect to a deployed contract
 *
 * @param address Deployed address of the contract
 * @param provider Network to connect to
 * @param contractInterface ABI of the contract
 */
export function createContract(
  address: string,
  provider: JsonRpcProvider,
  contractInterface: FunctionFragment[],
) {
  return new ethers.Contract(address, contractInterface, provider)
}
