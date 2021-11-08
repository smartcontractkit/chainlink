import { ethers } from 'hardhat'
import { assert, expect } from 'chai'
import etherslib, { Contract, ContractFactory } from 'ethers'
import { SignerWithAddress } from '@nomiclabs/hardhat-ethers/signers'
import {
  impersonateAs,
  publicAbi,
  stripHexPrefix,
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
let multisendFactory: ContractFactory
let forwarder: Contract
let greeter: Contract
let multisend: Contract

before(async () => {
  const accounts = await ethers.getSigners()
  owner = accounts[0]
  stranger = accounts[1]

  // Forwarder config
  l1OwnerAddress = owner.address
  crossdomainMessenger = await impersonateAs(
    toArbitrumL2AliasAddress(l1OwnerAddress),
  )
  await owner.sendTransaction({
    to: crossdomainMessenger.address,
    value: ethers.utils.parseEther('1.0'),
  })
  newL1OwnerAddress = stranger.address
  newOwnerCrossdomainMessenger = await impersonateAs(
    toArbitrumL2AliasAddress(newL1OwnerAddress),
  )
  await owner.sendTransaction({
    to: newOwnerCrossdomainMessenger.address,
    value: ethers.utils.parseEther('1.0'),
  })

  // Contract factories
  forwarderFactory = await ethers.getContractFactory(
    'src/v0.8/dev/ArbitrumCrossDomainGovernor.sol:ArbitrumCrossDomainGovernor',
    owner,
  )
  greeterFactory = await ethers.getContractFactory(
    'src/v0.8/tests/Greeter.sol:Greeter',
    owner,
  )
  multisendFactory = await ethers.getContractFactory(
    'src/v0.8/tests/vendor/MultiSend.sol:MultiSend',
    owner,
  )
})

describe('ArbitrumCrossDomainGovernor', () => {
  beforeEach(async () => {
    forwarder = await forwarderFactory.deploy(l1OwnerAddress)
    greeter = await greeterFactory.deploy(forwarder.address)
    multisend = await multisendFactory.deploy()
  })

  it('has a limited public interface [ @skip-coverage ]', async () => {
    publicAbi(forwarder, [
      'typeAndVersion',
      'crossDomainMessenger',
      'forward',
      'forwardDelegate',
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

    it('should be callable by L2 owner', async () => {
      const newGreeting = 'hello'
      const setGreetingData = greeterFactory.interface.encodeFunctionData(
        'setGreeting',
        [newGreeting],
      )
      await forwarder.connect(owner).forward(greeter.address, setGreetingData)

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
      ).to.be.revertedWith('Governor call reverted')
    })
  })

  describe('#forwardDelegate', () => {
    it('should not be callable by unknown address', async () => {
      await expect(
        forwarder.connect(stranger).forwardDelegate(multisend.address, '0x'),
      ).to.be.revertedWith('Sender is not the L2 messenger')
    })

    it('should be callable by crossdomain messenger address / L1 owner', async () => {
      const calls = [
        {
          to: greeter.address,
          data: greeterFactory.interface.encodeFunctionData('setGreeting', [
            'foo',
          ]),
          value: 0,
        },
        {
          to: greeter.address,
          data: greeterFactory.interface.encodeFunctionData('setGreeting', [
            'bar',
          ]),
          value: 0,
        },
      ]
      const multisendData = encodeMultisendData(multisend.interface, calls)
      await forwarder
        .connect(crossdomainMessenger)
        .forwardDelegate(multisend.address, multisendData)

      const updatedGreeting = await greeter.greeting()
      assert.equal(updatedGreeting, 'bar')
    })

    it('should be callable by L2 owner', async () => {
      const calls = [
        {
          to: greeter.address,
          data: greeterFactory.interface.encodeFunctionData('setGreeting', [
            'foo',
          ]),
          value: 0,
        },
        {
          to: greeter.address,
          data: greeterFactory.interface.encodeFunctionData('setGreeting', [
            'bar',
          ]),
          value: 0,
        },
      ]
      const multisendData = encodeMultisendData(multisend.interface, calls)
      await forwarder
        .connect(owner)
        .forwardDelegate(multisend.address, multisendData)

      const updatedGreeting = await greeter.greeting()
      assert.equal(updatedGreeting, 'bar')
    })

    it('should be revert batch when one call fails', async () => {
      const calls = [
        {
          to: greeter.address,
          data: greeterFactory.interface.encodeFunctionData('setGreeting', [
            'foo',
          ]),
          value: 0,
        },
        {
          to: greeter.address,
          data: greeterFactory.interface.encodeFunctionData('setGreeting', [
            '', // should revert
          ]),
          value: 0,
        },
      ]
      const multisendData = encodeMultisendData(multisend.interface, calls)
      await expect(
        forwarder
          .connect(crossdomainMessenger)
          .forwardDelegate(multisend.address, multisendData),
      ).to.be.revertedWith('Governor delegatecall reverted')

      const greeting = await greeter.greeting()
      assert.equal(greeting, '') // Unchanged
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
  })

  describe('#acceptL1Ownership', () => {
    it('should not be callable by non pending-owners', async () => {
      await expect(
        forwarder.connect(crossdomainMessenger).acceptL1Ownership(),
      ).to.be.revertedWith('Must be proposed owner')
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

// Multisend contract helpers

/**
 * Encodes an underlying transaction for the Multisend contract
 *
 * @param operation 0 for CALL, 1 for DELEGATECALL
 * @param to tx target address
 * @param value tx value
 * @param data tx data
 */
export function encodeTxData(
  operation: number,
  to: string,
  value: number,
  data: string,
): string {
  let dataBuffer = Buffer.from(stripHexPrefix(data), 'hex')
  const types = ['uint8', 'address', 'uint256', 'uint256', 'bytes']
  const values = [operation, to, value, dataBuffer.length, dataBuffer]
  let encoded = ethers.utils.solidityPack(types, values)
  return stripHexPrefix(encoded)
}

/**
 * Encodes a Multisend call
 *
 * @param MultisendInterface Ethers Interface object of the Multisend contract
 * @param transactions one or more transactions to include in the Multisend call
 * @param to tx target address
 * @param value tx value
 * @param data tx data
 */
export function encodeMultisendData(
  MultisendInterface: etherslib.utils.Interface,
  transactions: { to: string; value: number; data: string }[],
): string {
  let nestedTransactionData = '0x'
  for (let transaction of transactions) {
    nestedTransactionData += encodeTxData(
      0,
      transaction.to,
      transaction.value,
      transaction.data,
    )
  }
  const encodedMultisendFnData = MultisendInterface.encodeFunctionData(
    'multiSend',
    [nestedTransactionData],
  )
  return encodedMultisendFnData
}
