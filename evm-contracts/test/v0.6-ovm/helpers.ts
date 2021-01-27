import path from 'path'
import { ethers } from 'ethers'
import { Median__factory } from '../../ethers/v0.6-ovm/factories/Median__factory'

export const deployLibraries = async (wallet: ethers.Wallet): Promise<any> => {
  // Deploy Median
  const medianTx = await new Median__factory(wallet).deploy()
  const medianContract = await medianTx.deployed()

  const pathFragments = [__dirname, '..', '..', 'src', 'v0.6-ovm', 'Median.sol']
  const filePath = path.join(...pathFragments)
  const key = `${filePath}:Median`

  return {
    [key]: medianContract.address,
  }
}
