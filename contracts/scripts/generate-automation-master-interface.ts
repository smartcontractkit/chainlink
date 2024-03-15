/**
 * @description this script generates a master interface for interacting with the automation registry
 * @notice run this script with pnpm ts-node ./scripts/generate-automation-master-interface.ts
 */
import { AutomationRegistry2_2__factory as Registry } from '../typechain/factories/AutomationRegistry2_2__factory'
import { AutomationRegistryLogicA2_2__factory as RegistryLogicA } from '../typechain/factories/AutomationRegistryLogicA2_2__factory'
import { AutomationRegistryLogicB2_2__factory as RegistryLogicB } from '../typechain/factories/AutomationRegistryLogicB2_2__factory'
import { utils } from 'ethers'
import fs from 'fs'
import { exec } from 'child_process'

const dest = 'src/v0.8/automation/interfaces/v2_2'
const srcDest = `${dest}/IAutomationRegistryMaster.sol`
const tmpDest = `${dest}/tmp.txt`

const combinedABI = []
const abiSet = new Set()
const abis = [Registry.abi, RegistryLogicA.abi, RegistryLogicB.abi]

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
cat ${tmpDest} | pnpm abi-to-sol --solidity-version ^0.8.4 --license MIT > ${srcDest} IAutomationRegistryMaster;
echo "// abi-checksum: ${checksum}" | cat - ${srcDest} > ${tmpDest} && mv ${tmpDest} ${srcDest};
pnpm prettier --write ${srcDest};
`

exec(cmd)

console.log('generated new master interface for automation registry')
