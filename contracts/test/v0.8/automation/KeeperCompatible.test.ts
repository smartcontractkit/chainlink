import { ethers } from 'hardhat'
import { Contract } from 'ethers'
import { expect } from 'chai'
import { publicAbi } from '../../test-helpers/helpers'

describe('KeeperCompatible', () => {
  for (let version = 8; version <= 8; version++) {
    describe(`version v0.${version}`, () => {
      let contract: Contract

      before(async () => {
        const factory = await ethers.getContractFactory(
          `src/v0.${version}/automation/testhelpers/KeeperCompatibleTestHelper.sol:KeeperCompatibleTestHelper`,
        )
        contract = await factory.deploy()
      })

      it('has a keeper compatible interface [ @skip-coverage ]', async () => {
        publicAbi(contract, [
          'checkUpkeep',
          'performUpkeep',
          'verifyCannotExecute',
        ])
      })

      it('prevents execution of protected functions', async () => {
        await contract
          .connect(ethers.constants.AddressZero)
          .verifyCannotExecute() // succeeds
        await expect(contract.verifyCannotExecute()).to.be.reverted
      })
    })
  }
})
