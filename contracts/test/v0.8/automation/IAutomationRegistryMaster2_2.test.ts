import fs from 'fs'
import { ethers } from 'hardhat'
import { assert } from 'chai'
import { AutomationRegistry2_2__factory as AutomationRegistryFactory } from '../../../typechain/factories/AutomationRegistry2_2__factory'
import { AutomationRegistryLogicA2_2__factory as AutomationRegistryLogicAFactory } from '../../../typechain/factories/AutomationRegistryLogicA2_2__factory'
import { AutomationRegistryLogicB2_2__factory as AutomationRegistryLogicBFactory } from '../../../typechain/factories/AutomationRegistryLogicB2_2__factory'
import { AutomationRegistryBase2_2__factory as AutomationRegistryBaseFactory } from '../../../typechain/factories/AutomationRegistryBase2_2__factory'
import { Chainable__factory as ChainableFactory } from '../../../typechain/factories/Chainable__factory'
import { IAutomationRegistryMaster__factory as IAutomationRegistryMasterFactory } from '../../../typechain/factories/IAutomationRegistryMaster__factory'
import { IAutomationRegistryConsumer__factory as IAutomationRegistryConsumerFactory } from '../../../typechain/factories/IAutomationRegistryConsumer__factory'
import { MigratableKeeperRegistryInterface__factory as MigratableKeeperRegistryInterfaceFactory } from '../../../typechain/factories/MigratableKeeperRegistryInterface__factory'
import { MigratableKeeperRegistryInterfaceV2__factory as MigratableKeeperRegistryInterfaceV2Factory } from '../../../typechain/factories/MigratableKeeperRegistryInterfaceV2__factory'
import { OCR2Abstract__factory as OCR2AbstractFactory } from '../../../typechain/factories/OCR2Abstract__factory'
import { IAutomationV21PlusCommon__factory as IAutomationV21PlusCommonFactory } from '../../../typechain/factories/IAutomationV21PlusCommon__factory'
import {
  assertSatisfiesEvents,
  assertSatisfiesInterface,
  entryID,
} from './helpers'

const compositeABIs = [
  AutomationRegistryFactory.abi,
  AutomationRegistryLogicAFactory.abi,
  AutomationRegistryLogicBFactory.abi,
]

/**
 * @dev because the keeper master interface is a composite of several different contracts,
 * it is possible that an interface could be satisfied by functions across different
 * contracts, and therefore not enforceable by the compiler directly. Instead, we use this
 * test to assert that the master interface satisfies the constraints of an individual interface
 */
describe('IAutomationRegistryMaster2_2', () => {
  it('is up to date', async () => {
    const checksum = ethers.utils.id(compositeABIs.join(''))
    const knownChecksum = fs
      .readFileSync(
        'src/v0.8/automation/interfaces/v2_2/IAutomationRegistryMaster.sol',
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
      ...AutomationRegistryBaseFactory.abi,
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
      IAutomationRegistryMasterFactory.abi,
      IAutomationRegistryConsumerFactory.abi,
    )
  })

  it('satisfies the MigratableKeeperRegistryInterface interface', async () => {
    assertSatisfiesInterface(
      IAutomationRegistryMasterFactory.abi,
      MigratableKeeperRegistryInterfaceFactory.abi,
    )
  })

  it('satisfies the MigratableKeeperRegistryInterfaceV2 interface', async () => {
    assertSatisfiesInterface(
      IAutomationRegistryMasterFactory.abi,
      MigratableKeeperRegistryInterfaceV2Factory.abi,
    )
  })

  it('satisfies the OCR2Abstract interface', async () => {
    assertSatisfiesInterface(
      IAutomationRegistryMasterFactory.abi,
      OCR2AbstractFactory.abi,
    )
  })

  it('satisfies the IAutomationV2Common interface', async () => {
    assertSatisfiesInterface(
      IAutomationRegistryMasterFactory.abi,
      IAutomationV21PlusCommonFactory.abi,
    )
  })

  it('satisfies the IAutomationV2Common events', async () => {
    assertSatisfiesEvents(
      IAutomationRegistryMasterFactory.abi,
      IAutomationV21PlusCommonFactory.abi,
    )
  })
})
