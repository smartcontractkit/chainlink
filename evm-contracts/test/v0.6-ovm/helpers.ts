import path from 'path'
import { ethers } from 'ethers'
import { MedianFactory } from '../../ethers/v0.6-ovm/MedianFactory'

export const deployLibraries = async (wallet: ethers.Wallet): Promise<any> => {
  // Deploy Median
  const medianTx = await new MedianFactory(wallet).deploy()
  const medianContract = await medianTx.deployed()

  const pathFragments = [__dirname, '..', '..', 'src', 'v0.6-ovm', 'Median.sol']
  const filePath = path.join(...pathFragments)
  const key = `${filePath}:Median`

  return {
    [key]: medianContract.address,
  }
}
