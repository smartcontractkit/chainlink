import { ethers } from 'hardhat'
import { assert, expect } from 'chai'
import { Contract, ContractFactory } from 'ethers'
import { SignerWithAddress } from '@nomiclabs/hardhat-ethers/signers'

let owner: SignerWithAddress
let stranger: SignerWithAddress
let l1OwnerAddress: string
let ownableFactory: ContractFactory
let ownable: Contract

before(async () => {
  const accounts = await ethers.getSigners()
  owner = accounts[0]
  stranger = accounts[1]
  l1OwnerAddress = owner.address

  // Contract factories
  ownableFactory = await ethers.getContractFactory(
    'src/v0.8/l2ep/dev/CrossDomainOwnable.sol:CrossDomainOwnable',
    owner,
  )
})

describe('CrossDomainOwnable', () => {
  beforeEach(async () => {
    ownable = await ownableFactory.deploy(l1OwnerAddress)
  })

  describe('#constructor', () => {
    it('should set the l1Owner correctly', async () => {
      const response = await ownable.l1Owner()
      assert.equal(response, l1OwnerAddress)
    })
  })

  describe('#transferL1Ownership', () => {
    it('should not be callable by non-owners', async () => {
      await expect(
        ownable.connect(stranger).transferL1Ownership(stranger.address),
      ).to.be.revertedWith('Only callable by L1 owner')
    })

    it('should be callable by current L1 owner', async () => {
      const currentL1Owner = await ownable.l1Owner()
      await expect(ownable.transferL1Ownership(stranger.address))
        .to.emit(ownable, 'L1OwnershipTransferRequested')
        .withArgs(currentL1Owner, stranger.address)
    })

    it('should be callable by current L1 owner to zero address', async () => {
      const currentL1Owner = await ownable.l1Owner()
      await expect(ownable.transferL1Ownership(ethers.constants.AddressZero))
        .to.emit(ownable, 'L1OwnershipTransferRequested')
        .withArgs(currentL1Owner, ethers.constants.AddressZero)
    })
  })

  describe('#acceptL1Ownership', () => {
    it('should not be callable by non pending-owners', async () => {
      await expect(
        ownable.connect(stranger).acceptL1Ownership(),
      ).to.be.revertedWith('Only callable by proposed L1 owner')
    })

    it('should be callable by pending L1 owner', async () => {
      const currentL1Owner = await ownable.l1Owner()
      await ownable.transferL1Ownership(stranger.address)
      await expect(ownable.connect(stranger).acceptL1Ownership())
        .to.emit(ownable, 'L1OwnershipTransferred')
        .withArgs(currentL1Owner, stranger.address)

      const updatedL1Owner = await ownable.l1Owner()
      assert.equal(updatedL1Owner, stranger.address)
    })
  })
})
