/**
 * @description this script generates a master interface for interacting with the automation registry
 * @notice run this script with pnpm ts-node ./scripts/generate-automation-master-interface.ts
 */

const dest = 'src/v0.8/dev/automation/2_1/interfaces'
const srcDest = `${dest}/IKeeperRegistryMaster.sol`
const abiDest = `${dest}/temp.json`

const cmd = `
cat ${abiDest} | pnpm abi-to-sol --solidity-version ^0.8.4 --license MIT > ${srcDest} IKeeperRegistryMaster;
rm ${abiDest};
pnpm prettier --write src/v0.8/dev/automation/2_1/interfaces/IKeeperRegistryMaster.sol
`

import { KeeperRegistry2_1__factory as KeeperRegistry } from '../typechain/factories/KeeperRegistry2_1__factory'
import { KeeperRegistryLogicA2_1__factory as KeeperRegistryLogicA } from '../typechain/factories/KeeperRegistryLogicA2_1__factory'
import { KeeperRegistryLogicB2_1__factory as KeeperRegistryLogicB } from '../typechain/factories/KeeperRegistryLogicB2_1__factory'
import { utils } from 'ethers'
import fs from 'fs'
import { exec } from 'child_process'

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
      combinedABI.push(entry)
    }
  }
}

fs.writeFileSync(`${abiDest}`, JSON.stringify(combinedABI))

exec(cmd)

console.log('generated new master interface for automation registry')
