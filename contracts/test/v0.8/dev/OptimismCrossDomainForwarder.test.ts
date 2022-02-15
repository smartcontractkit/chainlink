import { ethers } from 'hardhat'
import { assert, expect } from 'chai'
import { Contract, ContractFactory } from 'ethers'
import { SignerWithAddress } from '@nomiclabs/hardhat-ethers/signers'
import { publicAbi } from '../../test-helpers/helpers'

let owner: SignerWithAddress
let stranger: SignerWithAddress
let l1OwnerAddress: string
let newL1OwnerAddress: string
let forwarderFactory: ContractFactory
let greeterFactory: ContractFactory
let crossDomainMessengerFactory: ContractFactory
let crossDomainMessenger: Contract
let forwarder: Contract
let greeter: Contract

before(async () => {
  const accounts = await ethers.getSigners()
  owner = accounts[0]
  stranger = accounts[1]

  // forwarder config
  l1OwnerAddress = owner.address
  newL1OwnerAddress = stranger.address

  // Contract factories
  forwarderFactory = await ethers.getContractFactory(
    'src/v0.8/dev/OptimismCrossDomainForwarder.sol:OptimismCrossDomainForwarder',
    owner,
  )
  greeterFactory = await ethers.getContractFactory(
    'src/v0.8/tests/Greeter.sol:Greeter',
    owner,
  )
  crossDomainMessengerFactory = await ethers.getContractFactory(
    'src/v0.8/tests/vendor/MockOVMCrossDomainMessenger.sol:MockOVMCrossDomainMessenger',
  )
})

describe('OptimismCrossDomainForwarder', () => {
  beforeEach(async () => {
    crossDomainMessenger = await crossDomainMessengerFactory.deploy(
      l1OwnerAddress,
    )
    forwarder = await forwarderFactory.deploy(
      crossDomainMessenger.address,
      l1OwnerAddress,
    )
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
      assert.equal(response, crossDomainMessenger.address)
    })

    it('should set the typeAndVersion correctly', async () => {
      const response = await forwarder.typeAndVersion()
      assert.equal(response, 'OptimismCrossDomainForwarder 1.0.0')
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
      const forwardData = forwarderFactory.interface.encodeFunctionData(
        'forward',
        [greeter.address, setGreetingData],
      )
      await crossDomainMessenger // Simulate cross-chain OVM message
        .connect(stranger)
        .sendMessage(forwarder.address, forwardData, 0)

      const updatedGreeting = await greeter.greeting()
      assert.equal(updatedGreeting, newGreeting)
    })

    it('should revert when contract call reverts', async () => {
      const setGreetingData = greeterFactory.interface.encodeFunctionData(
        'setGreeting',
        [''],
      )
      const forwardData = forwarderFactory.interface.encodeFunctionData(
        'forward',
        [greeter.address, setGreetingData],
      )
      await expect(
        crossDomainMessenger // Simulate cross-chain OVM message
          .connect(stranger)
          .sendMessage(forwarder.address, forwardData, 0),
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
      const forwardData = forwarderFactory.interface.encodeFunctionData(
        'transferL1Ownership',
        [newL1OwnerAddress],
      )

      await expect(
        crossDomainMessenger // Simulate cross-chain OVM message
          .connect(stranger)
          .sendMessage(forwarder.address, forwardData, 0),
      )
        .to.emit(forwarder, 'L1OwnershipTransferRequested')
        .withArgs(currentL1Owner, newL1OwnerAddress)
    })

    it('should be callable by current L1 owner to zero address', async () => {
      const currentL1Owner = await forwarder.l1Owner()
      const forwardData = forwarderFactory.interface.encodeFunctionData(
        'transferL1Ownership',
        [ethers.constants.AddressZero],
      )

      await expect(
        crossDomainMessenger // Simulate cross-chain OVM message
          .connect(stranger)
          .sendMessage(forwarder.address, forwardData, 0),
      )
        .to.emit(forwarder, 'L1OwnershipTransferRequested')
        .withArgs(currentL1Owner, ethers.constants.AddressZero)
    })
  })

  describe('#acceptL1Ownership', () => {
    it('should not be callable by non pending-owners', async () => {
      const forwardData = forwarderFactory.interface.encodeFunctionData(
        'acceptL1Ownership',
        [],
      )
      await expect(
        crossDomainMessenger // Simulate cross-chain OVM message
          .connect(stranger)
          .sendMessage(forwarder.address, forwardData, 0),
      ).to.be.revertedWith('Must be proposed L1 owner')
    })

    it('should be callable by pending L1 owner', async () => {
      const currentL1Owner = await forwarder.l1Owner()

      // Transfer ownership
      const forwardTransferData = forwarderFactory.interface.encodeFunctionData(
        'transferL1Ownership',
        [newL1OwnerAddress],
      )
      await crossDomainMessenger // Simulate cross-chain OVM message
        .connect(stranger)
        .sendMessage(forwarder.address, forwardTransferData, 0)

      const forwardAcceptData = forwarderFactory.interface.encodeFunctionData(
        'acceptL1Ownership',
        [],
      )
      // Simulate cross-chain message from another sender
      await crossDomainMessenger._setMockMessageSender(newL1OwnerAddress)

      await expect(
        crossDomainMessenger // Simulate cross-chain OVM message
          .connect(stranger)
          .sendMessage(forwarder.address, forwardAcceptData, 0),
      )
        .to.emit(forwarder, 'L1OwnershipTransferred')
        .withArgs(currentL1Owner, newL1OwnerAddress)

      const updatedL1Owner = await forwarder.l1Owner()
      assert.equal(updatedL1Owner, newL1OwnerAddress)
    })
  })
})
