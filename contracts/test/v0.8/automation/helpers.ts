import { Signer } from 'ethers'
import { ethers } from 'hardhat'
import { KeeperRegistryLogicB2_1__factory as KeeperRegistryLogicBFactory } from '../../../typechain/factories/KeeperRegistryLogicB2_1__factory'
import { IKeeperRegistryMaster as IKeeperRegistry } from '../../../typechain/IKeeperRegistryMaster'
import { IKeeperRegistryMaster__factory as IKeeperRegistryMasterFactory } from '../../../typechain/factories/IKeeperRegistryMaster__factory'
import { AutomationRegistryLogicB2_2__factory as AutomationRegistryLogicBFactory } from '../../../typechain/factories/AutomationRegistryLogicB2_2__factory'
import { IAutomationRegistryMaster as IAutomationRegistry } from '../../../typechain/IAutomationRegistryMaster'
import { IAutomationRegistryMaster__factory as IAutomationRegistryMasterFactory } from '../../../typechain/factories/IAutomationRegistryMaster__factory'

export const deployRegistry21 = async (
  from: Signer,
  mode: Parameters<KeeperRegistryLogicBFactory['deploy']>[0],
  link: Parameters<KeeperRegistryLogicBFactory['deploy']>[1],
  linkNative: Parameters<KeeperRegistryLogicBFactory['deploy']>[2],
  fastgas: Parameters<KeeperRegistryLogicBFactory['deploy']>[3],
): Promise<IKeeperRegistry> => {
  const logicBFactory = await ethers.getContractFactory(
    'KeeperRegistryLogicB2_1',
  )
  const logicAFactory = await ethers.getContractFactory(
    'KeeperRegistryLogicA2_1',
  )
  const registryFactory = await ethers.getContractFactory('KeeperRegistry2_1')
  const forwarderLogicFactory = await ethers.getContractFactory(
    'AutomationForwarderLogic',
  )
  const forwarderLogic = await forwarderLogicFactory.connect(from).deploy()
  const logicB = await logicBFactory
    .connect(from)
    .deploy(mode, link, linkNative, fastgas, forwarderLogic.address)
  const logicA = await logicAFactory.connect(from).deploy(logicB.address)
  const master = await registryFactory.connect(from).deploy(logicA.address)
  return IKeeperRegistryMasterFactory.connect(master.address, from)
}

export const deployRegistry22 = async (
  from: Signer,
  link: Parameters<AutomationRegistryLogicBFactory['deploy']>[0],
  linkNative: Parameters<AutomationRegistryLogicBFactory['deploy']>[1],
  fastgas: Parameters<AutomationRegistryLogicBFactory['deploy']>[2],
  allowedReadOnlyAddress: Parameters<
    AutomationRegistryLogicBFactory['deploy']
  >[3],
): Promise<IAutomationRegistry> => {
  const logicBFactory = await ethers.getContractFactory(
    'AutomationRegistryLogicB2_2',
  )
  const logicAFactory = await ethers.getContractFactory(
    'AutomationRegistryLogicA2_2',
  )
  const registryFactory = await ethers.getContractFactory(
    'AutomationRegistry2_2',
  )
  const forwarderLogicFactory = await ethers.getContractFactory(
    'AutomationForwarderLogic',
  )
  const forwarderLogic = await forwarderLogicFactory.connect(from).deploy()
  const logicB = await logicBFactory
    .connect(from)
    .deploy(
      link,
      linkNative,
      fastgas,
      forwarderLogic.address,
      allowedReadOnlyAddress,
    )
  const logicA = await logicAFactory.connect(from).deploy(logicB.address)
  const master = await registryFactory.connect(from).deploy(logicA.address)
  return IAutomationRegistryMasterFactory.connect(master.address, from)
}
