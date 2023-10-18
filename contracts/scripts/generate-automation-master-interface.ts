/**
 * @description this script generates a master interface for interacting with the automation registry
 * @notice run this script with pnpm ts-node ./scripts/generate-automation-master-interface.ts
 */
import { KeeperRegistry2_1__factory as KeeperRegistry } from '../typechain/factories/KeeperRegistry2_1__factory'
import { KeeperRegistryLogicA2_1__factory as KeeperRegistryLogicA } from '../typechain/factories/KeeperRegistryLogicA2_1__factory'
import { KeeperRegistryLogicB2_1__factory as KeeperRegistryLogicB } from '../typechain/factories/KeeperRegistryLogicB2_1__factory'
import { utils } from 'ethers'
import fs from 'fs'
import { exec } from 'child_process'

const dest = 'src/v0.8/automation/interfaces/v2_1'
const srcDest = `${dest}/IKeeperRegistryMaster.sol`
const tmpDest = `${dest}/tmp.txt`

const combinedABI = []
const abiSet = new Set()
const abis = [
  KeeperRegistry.abi,
  KeeperRegistryLogicA.abi,
  KeeperRegistryLogicB.abi,
]

for (const abi of abis) {
  for (const entry of abi) {
    const id = utils.id(JSON.stringify(entry))
    if (!abiSet.has(id)) {
      abiSet.add(id)
      if (
        entry.type === 'function' &&
        (entry.name === 'checkUpkeep' ||
          entry.name === 'checkCallback' ||
          entry.name === 'simulatePerformUpkeep')
      ) {
        entry.stateMutability = 'view' // override stateMutability for check / callback / simulate functions
      }
      combinedABI.push(entry)
    }
  }
}

const checksum = utils.id(abis.join(''))

fs.writeFileSync(`${tmpDest}`, JSON.stringify(combinedABI))

const cmd = `
cat ${tmpDest} | pnpm abi-to-sol --solidity-version ^0.8.4 --license MIT > ${srcDest} IKeeperRegistryMaster;
echo "// abi-checksum: ${checksum}" | cat - ${srcDest} > ${tmpDest} && mv ${tmpDest} ${srcDest};
pnpm prettier --write ${srcDest};
`

exec(cmd)

console.log('generated new master interface for automation registry')
