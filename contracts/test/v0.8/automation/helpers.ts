import { Signer } from 'ethers'
import { ethers } from 'hardhat'
import { KeeperRegistryLogicB2_1__factory as KeeperRegistryLogicBFactory } from '../../../typechain/factories/KeeperRegistryLogicB2_1__factory'
import { IKeeperRegistryMaster as IKeeperRegistry } from '../../../typechain/IKeeperRegistryMaster'
import { IKeeperRegistryMaster__factory as IKeeperRegistryMasterFactory } from '../../../typechain/factories/IKeeperRegistryMaster__factory'
import { AutomationRegistryLogicB2_2__factory as AutomationRegistryLogicBFactory } from '../../../typechain/factories/AutomationRegistryLogicB2_2__factory'
import { IAutomationRegistryMaster as IAutomationRegistry } from '../../../typechain/IAutomationRegistryMaster'
import { IAutomationRegistryMaster__factory as IAutomationRegistryMasterFactory } from '../../../typechain/factories/IAutomationRegistryMaster__factory'
import { assert } from 'chai'
import { FunctionFragment } from '@ethersproject/abi'
import { AutomationRegistryLogicB2_3__factory as AutomationRegistryLogicB2_3Factory } from '../../../typechain/factories/AutomationRegistryLogicB2_3__factory'
import { IAutomationRegistryMaster2_3 as IAutomationRegistry2_3 } from '../../../typechain/IAutomationRegistryMaster2_3'
import { IAutomationRegistryMaster2_3__factory as IAutomationRegistryMaster2_3Factory } from '../../../typechain/factories/IAutomationRegistryMaster2_3__factory'

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

type InterfaceABI = ConstructorParameters<typeof ethers.utils.Interface>[0]
type Entry = {
  inputs?: any[]
  outputs?: any[]
  name?: string
  type: string
}

export const assertSatisfiesEvents = (
  contractABI: InterfaceABI,
  expectedABI: InterfaceABI,
) => {
  const implementer = new ethers.utils.Interface(contractABI)
  const expected = new ethers.utils.Interface(expectedABI)
  for (const eventName in expected.events) {
    assert.isDefined(
      implementer.events[eventName],
      `missing event: ${eventName}`,
    )
  }
}

export const entryID = (entry: Entry) => {
  // remove "internal type" and "name" since they don't affect the ability
  // of a contract to satisfy an interface
  const preimage = Object.assign({}, entry)
  if (entry.inputs) {
    preimage.inputs = entry.inputs.map(({ type }) => ({
      type,
    }))
  }
  if (entry.outputs) {
    preimage.outputs = entry.outputs.map(({ type }) => ({
      type,
    }))
  }
  return ethers.utils.id(JSON.stringify(preimage))
}

export const assertSatisfiesInterface = (
  contractABI: InterfaceABI,
  expectedABI: InterfaceABI,
) => {
  const implementer = new ethers.utils.Interface(contractABI)
  const expected = new ethers.utils.Interface(expectedABI)
  for (const functionName in expected.functions) {
    assert.isDefined(
      implementer.functions[functionName],
      `missing function ${functionName}`,
    )

    // these are technically pure in those interfaces. but in the master interface, they are view functions
    // bc the underlying contracts define constants for these values and return them in these getters
    if (
      functionName === 'typeAndVersion()' ||
      functionName === 'upkeepVersion()' ||
      functionName === 'upkeepTranscoderVersion()'
    ) {
      assert.equal(
        implementer.functions[functionName].constant,
        expected.functions[functionName].constant,
        `property constant does not match for function ${functionName}`,
      )
      assert.equal(
        implementer.functions[functionName].payable,
        expected.functions[functionName].payable,
        `property payable does not match for function ${functionName}`,
      )
      continue
    }

    const propertiesToMatch: (keyof FunctionFragment)[] = [
      'constant',
      'stateMutability',
      'payable',
    ]
    for (const property of propertiesToMatch) {
      assert.equal(
        implementer.functions[functionName][property],
        expected.functions[functionName][property],
        `property ${property} does not match for function ${functionName}`,
      )
    }
  }
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

export const deployRegistry23 = async (
  from: Signer,
  link: Parameters<AutomationRegistryLogicB2_3Factory['deploy']>[0],
  linkUSD: Parameters<AutomationRegistryLogicB2_3Factory['deploy']>[1],
  nativeUSD: Parameters<AutomationRegistryLogicB2_3Factory['deploy']>[2],
  fastgas: Parameters<AutomationRegistryLogicB2_3Factory['deploy']>[2],
  allowedReadOnlyAddress: Parameters<
    AutomationRegistryLogicB2_3Factory['deploy']
  >[3],
): Promise<IAutomationRegistry2_3> => {
  const logicBFactory = await ethers.getContractFactory(
    'AutomationRegistryLogicB2_3',
  )
  const logicAFactory = await ethers.getContractFactory(
    'AutomationRegistryLogicA2_3',
  )
  const registryFactory = await ethers.getContractFactory(
    'AutomationRegistry2_3',
  )
  const forwarderLogicFactory = await ethers.getContractFactory(
    'AutomationForwarderLogic',
  )
  const forwarderLogic = await forwarderLogicFactory.connect(from).deploy()
  const logicB = await logicBFactory
    .connect(from)
    .deploy(
      link,
      linkUSD,
      nativeUSD,
      fastgas,
      forwarderLogic.address,
      allowedReadOnlyAddress,
    )
  const logicA = await logicAFactory.connect(from).deploy(logicB.address)
  const master = await registryFactory.connect(from).deploy(logicA.address)
  return IAutomationRegistryMaster2_3Factory.connect(master.address, from)
}
