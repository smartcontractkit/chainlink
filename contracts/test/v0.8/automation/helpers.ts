import { Signer, BigNumberish } from 'ethers'
import { ethers } from 'hardhat'
import { KeeperRegistryLogicB2_1__factory as KeeperRegistryLogicBFactory } from '../../../typechain/factories/KeeperRegistryLogicB2_1__factory'
import { IKeeperRegistryMaster as IKeeperRegistry } from '../../../typechain/IKeeperRegistryMaster'
import { IKeeperRegistryMaster__factory as IKeeperRegistryMasterFactory } from '../../../typechain/factories/IKeeperRegistryMaster__factory'

export const deployRegistry21 = async (
  from: Signer,
  ...params: Parameters<KeeperRegistryLogicBFactory['deploy']>
): Promise<IKeeperRegistry> => {
  const logicBFactory = await ethers.getContractFactory(
    'KeeperRegistryLogicB2_1',
  )
  const logicAFactory = await ethers.getContractFactory(
    'KeeperRegistryLogicA2_1',
  )
  const registryFactory = await ethers.getContractFactory('KeeperRegistry2_1')
  const logicB = await logicBFactory.connect(from).deploy(...params)
  const logicA = await logicAFactory.connect(from).deploy(logicB.address)
  const master = await registryFactory.connect(from).deploy(logicA.address)
  return IKeeperRegistryMasterFactory.connect(master.address, from)
}

export type OnchainConfig21 = {
  paymentPremiumPPB: BigNumberish
  flatFeeMicroLink: BigNumberish
  checkGasLimit: BigNumberish
  stalenessSeconds: BigNumberish
  gasCeilingMultiplier: BigNumberish
  minUpkeepSpend: BigNumberish
  maxCheckDataSize: BigNumberish
  maxPerformDataSize: BigNumberish
  maxRevertDataSize: BigNumberish
  maxPerformGas: BigNumberish
  fallbackGasPrice: BigNumberish
  fallbackLinkPrice: BigNumberish
  transcoder: string
  registrars: string[]
  upkeepPrivilegeManager: string
}

export const encodeConfig21 = (config: OnchainConfig21) => {
  return ethers.utils.defaultAbiCoder.encode(
    [
      'tuple(uint32 paymentPremiumPPB,uint32 flatFeeMicroLink,uint32 checkGasLimit,uint24 stalenessSeconds\
      ,uint16 gasCeilingMultiplier,uint96 minUpkeepSpend,uint32 maxPerformGas,uint32 maxCheckDataSize,\
      uint32 maxPerformDataSize,uint32 maxRevertDataSize,uint256 fallbackGasPrice,uint256 fallbackLinkPrice,address transcoder,\
      address[] registrars,address upkeepPrivilegeManager)',
    ],
    [config],
  )
}
