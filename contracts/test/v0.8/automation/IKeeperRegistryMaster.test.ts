import fs from 'fs'
import { ethers } from 'hardhat'
import { assert } from 'chai'
import { FunctionFragment } from '@ethersproject/abi'
import { KeeperRegistry2_1__factory as KeeperRegistryFactory } from '../../../typechain/factories/KeeperRegistry2_1__factory'
import { KeeperRegistryLogicA2_1__factory as KeeperRegistryLogicAFactory } from '../../../typechain/factories/KeeperRegistryLogicA2_1__factory'
import { KeeperRegistryLogicB2_1__factory as KeeperRegistryLogicBFactory } from '../../../typechain/factories/KeeperRegistryLogicB2_1__factory'
import { KeeperRegistryBase2_1__factory as KeeperRegistryBaseFactory } from '../../../typechain/factories/KeeperRegistryBase2_1__factory'
import { Chainable__factory as ChainableFactory } from '../../../typechain/factories/Chainable__factory'
import { IKeeperRegistryMaster__factory as IKeeperRegistryMasterFactory } from '../../../typechain/factories/IKeeperRegistryMaster__factory'
import { IAutomationRegistryConsumer__factory as IAutomationRegistryConsumerFactory } from '../../../typechain/factories/IAutomationRegistryConsumer__factory'
import { MigratableKeeperRegistryInterface__factory as MigratableKeeperRegistryInterfaceFactory } from '../../../typechain/factories/MigratableKeeperRegistryInterface__factory'
import { MigratableKeeperRegistryInterfaceV2__factory as MigratableKeeperRegistryInterfaceV2Factory } from '../../../typechain/factories/MigratableKeeperRegistryInterfaceV2__factory'
import { OCR2Abstract__factory as OCR2AbstractFactory } from '../../../typechain/factories/OCR2Abstract__factory'

type Entry = {
  inputs?: any[]
  outputs?: any[]
  name?: string
  type: string
}

type InterfaceABI = ConstructorParameters<typeof ethers.utils.Interface>[0]

const compositeABIs = [
  KeeperRegistryFactory.abi,
  KeeperRegistryLogicAFactory.abi,
  KeeperRegistryLogicBFactory.abi,
]

function entryID(entry: Entry) {
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

/**
 * @dev because the keeper master interface is a composit of several different contracts,
 * it is possible that a interface could be satisfied by functions across different
 * contracts, and therefore not enforcable by the compiler directly. Instead, we use this
 * test to assert that the master interface satisfies the contraints of an individual interface
 */
function assertSatisfiesInterface(
  contractABI: InterfaceABI,
  expectedABI: InterfaceABI,
) {
  const implementer = new ethers.utils.Interface(contractABI)
  const expected = new ethers.utils.Interface(expectedABI)
  for (const functionName in expected.functions) {
    if (
      Object.prototype.hasOwnProperty.call(expected, functionName) &&
      functionName.match('^.+(.*)$') // only match typed function sigs
    ) {
      assert.isDefined(
        implementer.functions[functionName],
        `missing function ${functionName}`,
      )
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
}

describe('IKeeperRegistryMaster', () => {
  it('is up to date', async () => {
    const checksum = ethers.utils.id(compositeABIs.join(''))
    const knownChecksum = fs
      .readFileSync(
        'src/v0.8/automation/interfaces/v2_1/IKeeperRegistryMaster.sol',
      )
      .toString()
      .slice(17, 83) // checksum located at top of file
    assert.equal(
      checksum,
      knownChecksum,
      'master interface is out of date - regenerate using "pnpm ts-node ./scripts/generate-automation-master-interface.ts"',
    )
  })

  it('is generated from composite contracts without competing definitions', async () => {
    const sharedEntries = [
      ...ChainableFactory.abi,
      ...KeeperRegistryBaseFactory.abi,
    ]
    const abiSet = new Set()
    const sharedSet = new Set()
    for (const entry of sharedEntries) {
      sharedSet.add(entryID(entry))
    }
    for (const abi of compositeABIs) {
      for (const entry of abi) {
        const id = entryID(entry)
        if (!abiSet.has(id)) {
          abiSet.add(id)
        } else if (!sharedSet.has(id)) {
          assert.fail(
            `composite contracts contain duplicate entry: ${JSON.stringify(
              entry,
            )}`,
          )
        }
      }
    }
  })

  it('satisfies the IAutomationRegistryConsumer interface', async () => {
    assertSatisfiesInterface(
      IKeeperRegistryMasterFactory.abi,
      IAutomationRegistryConsumerFactory.abi,
    )
  })

  it('satisfies the MigratableKeeperRegistryInterface interface', async () => {
    assertSatisfiesInterface(
      IKeeperRegistryMasterFactory.abi,
      MigratableKeeperRegistryInterfaceFactory.abi,
    )
  })

  it('satisfies the MigratableKeeperRegistryInterfaceV2 interface', async () => {
    assertSatisfiesInterface(
      IKeeperRegistryMasterFactory.abi,
      MigratableKeeperRegistryInterfaceV2Factory.abi,
    )
  })

  it('satisfies the OCR2Abstract interface', async () => {
    assertSatisfiesInterface(
      IKeeperRegistryMasterFactory.abi,
      OCR2AbstractFactory.abi,
    )
  })
})
