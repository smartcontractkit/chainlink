import fs from 'fs'
import { ethers } from 'hardhat'
import { assert } from 'chai'
import { KeeperRegistry2_1__factory as KeeperRegistryFactory } from '../../../typechain/factories/KeeperRegistry2_1__factory'
import { KeeperRegistryLogicA2_1__factory as KeeperRegistryLogicAFactory } from '../../../typechain/factories/KeeperRegistryLogicA2_1__factory'
import { KeeperRegistryLogicB2_1__factory as KeeperRegistryLogicBFactory } from '../../../typechain/factories/KeeperRegistryLogicB2_1__factory'
import { KeeperRegistryBase2_1__factory as KeeperRegistryBaseFactory } from '../../../typechain/factories/KeeperRegistryBase2_1__factory'
import { Chainable__factory as ChainableFactory } from '../../../typechain/factories/Chainable__factory'

const compositeABIs = [
  KeeperRegistryFactory.abi,
  KeeperRegistryLogicAFactory.abi,
  KeeperRegistryLogicBFactory.abi,
]

describe('IKeeperRegistryMaster', () => {
  it('is up to date', async () => {
    const checksum = ethers.utils.id(compositeABIs.join(''))
    const knownChecksum = fs
      .readFileSync(
        'src/v0.8/dev/automation/2_1/interfaces/IKeeperRegistryMaster.sol',
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
      const id = ethers.utils.id(JSON.stringify(entry))
      sharedSet.add(id)
    }
    for (const abi of compositeABIs) {
      for (const entry of abi) {
        const id = ethers.utils.id(JSON.stringify(entry))
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
})
