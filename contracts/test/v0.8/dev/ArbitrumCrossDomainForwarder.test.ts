import { ethers } from 'hardhat'
import { assert, expect } from 'chai'
import { Contract, ContractFactory } from 'ethers'
import { SignerWithAddress } from '@nomiclabs/hardhat-ethers/signers'
import {
  impersonateAs,
  publicAbi,
  toArbitrumL2AliasAddress,
} from '../../test-helpers/helpers'

let owner: SignerWithAddress
let stranger: SignerWithAddress
let l1OwnerAddress: string
let crossdomainMessenger: SignerWithAddress
let newL1OwnerAddress: string
let newOwnerCrossdomainMessenger: SignerWithAddress
let forwarderFactory: ContractFactory
let greeterFactory: ContractFactory
let forwarder: Contract
let greeter: Contract

before(async () => {
  const accounts = await ethers.getSigners()
  owner = accounts[0]
  stranger = accounts[1]
  l1OwnerAddress = owner.address
  newL1OwnerAddress = stranger.address

  // Contract factories
  forwarderFactory = await ethers.getContractFactory(
    'src/v0.8/dev/ArbitrumCrossDomainForwarder.sol:ArbitrumCrossDomainForwarder',
    owner,
  )
  greeterFactory = await ethers.getContractFactory(
    'src/v0.8/tests/Greeter.sol:Greeter',
    owner,
  )
})

describe('ArbitrumCrossDomainForwarder', () => {
  beforeEach(async () => {
    // governor config
    crossdomainMessenger = await impersonateAs(
      toArbitrumL2AliasAddress(l1OwnerAddress),
    )
    await owner.sendTransaction({
      to: crossdomainMessenger.address,
      value: ethers.utils.parseEther('1.0'),
    })
    newOwnerCrossdomainMessenger = await impersonateAs(
      toArbitrumL2AliasAddress(newL1OwnerAddress),
    )
    await owner.sendTransaction({
      to: newOwnerCrossdomainMessenger.address,
      value: ethers.utils.parseEther('1.0'),
    })

    forwarder = await forwarderFactory.deploy(l1OwnerAddress)
    greeter = await greeterFactory.deploy(forwarder.address)
  })

  it('has a limited public interface [ @skip-coverage ]', async () => {
    publicAbi(forwarder, [
      'typeAndVersion',
      'crossDomainMessenger',
      'forward',
      'l1Owner',
      'transferL1Ownership',
      'acceptL1Ownership',
      // ConfirmedOwner methods:
      'owner',
      'transferOwnership',
      'acceptOwnership',
    ])
  })

  describe('#constructor', () => {
    it('should set the owner correctly', async () => {
      const response = await forwarder.owner()
      assert.equal(response, owner.address)
    })

    it('should set the l1Owner correctly', async () => {
      const response = await forwarder.l1Owner()
      assert.equal(response, l1OwnerAddress)
    })

    it('should set the crossdomain messenger correctly', async () => {
      const response = await forwarder.crossDomainMessenger()
      assert.equal(response, crossdomainMessenger.address)
    })

    it('should set the typeAndVersion correctly', async () => {
      const response = await forwarder.typeAndVersion()
      assert.equal(response, 'ArbitrumCrossDomainForwarder 1.0.0')
    })
  })

  describe('#forward', () => {
    it('should not be callable by unknown address', async () => {
      await expect(
        forwarder.connect(stranger).forward(greeter.address, '0x'),
      ).to.be.revertedWith('Sender is not the L2 messenger')
    })

    it('should be callable by crossdomain messenger address / L1 owner', async () => {
      const newGreeting = 'hello'
      const setGreetingData = greeterFactory.interface.encodeFunctionData(
        'setGreeting',
        [newGreeting],
      )
      await forwarder
        .connect(crossdomainMessenger)
        .forward(greeter.address, setGreetingData)

      const updatedGreeting = await greeter.greeting()
      assert.equal(updatedGreeting, newGreeting)
    })

    it('should revert when contract call reverts', async () => {
      const setGreetingData = greeterFactory.interface.encodeFunctionData(
        'setGreeting',
        [''],
      )
      await expect(
        forwarder
          .connect(crossdomainMessenger)
          .forward(greeter.address, setGreetingData),
      ).to.be.revertedWith('Invalid greeting length')
    })
  })

  describe('#transferL1Ownership', () => {
    it('should not be callable by non-owners', async () => {
      await expect(
        forwarder.connect(stranger).transferL1Ownership(stranger.address),
      ).to.be.revertedWith('Sender is not the L2 messenger')
    })

    it('should not be callable by L2 owner', async () => {
      const forwarderOwner = await forwarder.owner()
      assert.equal(forwarderOwner, owner.address)

      await expect(
        forwarder.connect(owner).transferL1Ownership(stranger.address),
      ).to.be.revertedWith('Sender is not the L2 messenger')
    })

    it('should be callable by current L1 owner', async () => {
      const currentL1Owner = await forwarder.l1Owner()
      await expect(
        forwarder
          .connect(crossdomainMessenger)
          .transferL1Ownership(newL1OwnerAddress),
      )
        .to.emit(forwarder, 'L1OwnershipTransferRequested')
        .withArgs(currentL1Owner, newL1OwnerAddress)
    })

    it('should be callable by current L1 owner to zero address', async () => {
      const currentL1Owner = await forwarder.l1Owner()
      await expect(
        forwarder
          .connect(crossdomainMessenger)
          .transferL1Ownership(ethers.constants.AddressZero),
      )
        .to.emit(forwarder, 'L1OwnershipTransferRequested')
        .withArgs(currentL1Owner, ethers.constants.AddressZero)
    })
  })

  describe('#acceptL1Ownership', () => {
    it('should not be callable by non pending-owners', async () => {
      await expect(
        forwarder.connect(crossdomainMessenger).acceptL1Ownership(),
      ).to.be.revertedWith('Must be proposed L1 owner')
    })

    it('should be callable by pending L1 owner', async () => {
      const currentL1Owner = await forwarder.l1Owner()
      await forwarder
        .connect(crossdomainMessenger)
        .transferL1Ownership(newL1OwnerAddress)
      await expect(
        forwarder.connect(newOwnerCrossdomainMessenger).acceptL1Ownership(),
      )
        .to.emit(forwarder, 'L1OwnershipTransferred')
        .withArgs(currentL1Owner, newL1OwnerAddress)

      const updatedL1Owner = await forwarder.l1Owner()
      assert.equal(updatedL1Owner, newL1OwnerAddress)
    })
  })
})
