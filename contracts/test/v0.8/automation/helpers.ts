import { Signer } from 'ethers'
import { ethers } from 'hardhat'
import { assert } from 'chai'
import { FunctionFragment } from '@ethersproject/abi'
import { AutomationRegistryLogicC2_3__factory as AutomationRegistryLogicC2_3Factory } from '../../../typechain/factories/AutomationRegistryLogicC2_3__factory'
import { ZKSyncAutomationRegistryLogicC2_3__factory as ZKSyncAutomationRegistryLogicC2_3Factory } from '../../../typechain/factories/ZKSyncAutomationRegistryLogicC2_3__factory'
import { IAutomationRegistryMaster2_3 as IAutomationRegistry2_3 } from '../../../typechain/IAutomationRegistryMaster2_3'
import { IAutomationRegistryMaster2_3__factory as IAutomationRegistryMaster2_3Factory } from '../../../typechain/factories/IAutomationRegistryMaster2_3__factory'

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

    // addFunds is a payable function starting from v2.3, skipping the stateMutability check
    if (functionName === 'addFunds(uint256,uint96)') {
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

export const deployRegistry23 = async (
  from: Signer,
  link: Parameters<AutomationRegistryLogicC2_3Factory['deploy']>[0],
  linkUSD: Parameters<AutomationRegistryLogicC2_3Factory['deploy']>[1],
  nativeUSD: Parameters<AutomationRegistryLogicC2_3Factory['deploy']>[2],
  fastgas: Parameters<AutomationRegistryLogicC2_3Factory['deploy']>[3],
  allowedReadOnlyAddress: Parameters<
    AutomationRegistryLogicC2_3Factory['deploy']
  >[5],
  payoutMode: Parameters<AutomationRegistryLogicC2_3Factory['deploy']>[6],
  wrappedNativeTokenAddress: Parameters<
    AutomationRegistryLogicC2_3Factory['deploy']
  >[7],
): Promise<IAutomationRegistry2_3> => {
  const logicCFactory = await ethers.getContractFactory(
    'AutomationRegistryLogicC2_3',
  )
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
  const logicC = await logicCFactory
    .connect(from)
    .deploy(
      link,
      linkUSD,
      nativeUSD,
      fastgas,
      forwarderLogic.address,
      allowedReadOnlyAddress,
      payoutMode,
      wrappedNativeTokenAddress,
    )
  const logicB = await logicBFactory.connect(from).deploy(logicC.address)
  const logicA = await logicAFactory.connect(from).deploy(logicB.address)
  const master = await registryFactory.connect(from).deploy(logicA.address)
  return IAutomationRegistryMaster2_3Factory.connect(master.address, from)
}

export const deployZKSyncRegistry23 = async (
  from: Signer,
  link: Parameters<ZKSyncAutomationRegistryLogicC2_3Factory['deploy']>[0],
  linkUSD: Parameters<ZKSyncAutomationRegistryLogicC2_3Factory['deploy']>[1],
  nativeUSD: Parameters<ZKSyncAutomationRegistryLogicC2_3Factory['deploy']>[2],
  fastgas: Parameters<ZKSyncAutomationRegistryLogicC2_3Factory['deploy']>[3],
  allowedReadOnlyAddress: Parameters<
    AutomationRegistryLogicC2_3Factory['deploy']
  >[5],
  payoutMode: Parameters<ZKSyncAutomationRegistryLogicC2_3Factory['deploy']>[6],
  wrappedNativeTokenAddress: Parameters<
    ZKSyncAutomationRegistryLogicC2_3Factory['deploy']
  >[7],
): Promise<IAutomationRegistry2_3> => {
  const logicCFactory = await ethers.getContractFactory(
    'ZKSyncAutomationRegistryLogicC2_3',
  )
  const logicBFactory = await ethers.getContractFactory(
    'ZKSyncAutomationRegistryLogicB2_3',
  )
  const logicAFactory = await ethers.getContractFactory(
    'ZKSyncAutomationRegistryLogicA2_3',
  )
  const registryFactory = await ethers.getContractFactory(
    'ZKSyncAutomationRegistry2_3',
  )
  const forwarderLogicFactory = await ethers.getContractFactory(
    'AutomationForwarderLogic',
  )
  const forwarderLogic = await forwarderLogicFactory.connect(from).deploy()
  const logicC = await logicCFactory
    .connect(from)
    .deploy(
      link,
      linkUSD,
      nativeUSD,
      fastgas,
      forwarderLogic.address,
      allowedReadOnlyAddress,
      payoutMode,
      wrappedNativeTokenAddress,
    )
  const logicB = await logicBFactory.connect(from).deploy(logicC.address)
  const logicA = await logicAFactory.connect(from).deploy(logicB.address)
  const master = await registryFactory.connect(from).deploy(logicA.address)
  return IAutomationRegistryMaster2_3Factory.connect(master.address, from)
}
