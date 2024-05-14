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
let governorFactory: ContractFactory
let greeterFactory: ContractFactory
let multisendFactory: ContractFactory
let governor: Contract
let greeter: Contract
let multisend: Contract

before(async () => {
  const accounts = await ethers.getSigners()
  owner = accounts[0]
  stranger = accounts[1]
  l1OwnerAddress = owner.address
  newL1OwnerAddress = stranger.address

  // Contract factories
  governorFactory = await ethers.getContractFactory(
    'src/v0.8/l2ep/dev/arbitrum/ArbitrumCrossDomainGovernor.sol:ArbitrumCrossDomainGovernor',
    owner,
  )
  greeterFactory = await ethers.getContractFactory(
    'src/v0.8/tests/Greeter.sol:Greeter',
    owner,
  )
  multisendFactory = await ethers.getContractFactory(
    'src/v0.8/vendor/MultiSend.sol:MultiSend',
    owner,
  )
})

describe('ArbitrumCrossDomainGovernor', () => {
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

    governor = await governorFactory.deploy(l1OwnerAddress)
    greeter = await greeterFactory.deploy(governor.address)
    multisend = await multisendFactory.deploy()
  })

  it('has a limited public interface [ @skip-coverage ]', async () => {
    publicAbi(governor, [
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
      const response = await governor.owner()
      assert.equal(response, owner.address)
    })

    it('should set the l1Owner correctly', async () => {
      const response = await governor.l1Owner()
      assert.equal(response, l1OwnerAddress)
    })

    it('should set the crossdomain messenger correctly', async () => {
      const response = await governor.crossDomainMessenger()
      assert.equal(response, crossdomainMessenger.address)
    })

    it('should set the typeAndVersion correctly', async () => {
      const response = await governor.typeAndVersion()
      assert.equal(response, 'ArbitrumCrossDomainGovernor 1.0.0')
    })
  })

  describe('#forward', () => {
    it('should not be callable by unknown address', async () => {
      await expect(
        governor.connect(stranger).forward(greeter.address, '0x'),
      ).to.be.revertedWith('Sender is not the L2 messenger or owner')
    })

    it('should be callable by crossdomain messenger address / L1 owner', async () => {
      const newGreeting = 'hello'
      const setGreetingData = greeterFactory.interface.encodeFunctionData(
        'setGreeting',
        [newGreeting],
      )
      await governor
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
      await governor.connect(owner).forward(greeter.address, setGreetingData)

      const updatedGreeting = await greeter.greeting()
      assert.equal(updatedGreeting, newGreeting)
    })

    it('should revert when contract call reverts', async () => {
      const setGreetingData = greeterFactory.interface.encodeFunctionData(
        'setGreeting',
        [''],
      )
      await expect(
        governor
          .connect(crossdomainMessenger)
          .forward(greeter.address, setGreetingData),
      ).to.be.revertedWith('Invalid greeting length')
    })
  })

  describe('#forwardDelegate', () => {
    it('should not be callable by unknown address', async () => {
      await expect(
        governor.connect(stranger).forwardDelegate(multisend.address, '0x'),
      ).to.be.revertedWith('Sender is not the L2 messenger or owner')
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
      await governor
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
      await governor
        .connect(owner)
        .forwardDelegate(multisend.address, multisendData)

      const updatedGreeting = await greeter.greeting()
      assert.equal(updatedGreeting, 'bar')
    })

    it('should revert batch when one call fails', async () => {
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
        governor
          .connect(crossdomainMessenger)
          .forwardDelegate(multisend.address, multisendData),
      ).to.be.revertedWith('Governor delegatecall reverted')

      const greeting = await greeter.greeting()
      assert.equal(greeting, '') // Unchanged
    })

    it('should bubble up revert when contract call reverts', async () => {
      const triggerRevertData =
        greeterFactory.interface.encodeFunctionData('triggerRevert')
      await expect(
        governor
          .connect(crossdomainMessenger)
          .forwardDelegate(greeter.address, triggerRevertData),
      ).to.be.revertedWith('Greeter: revert triggered')
    })
  })

  describe('#transferL1Ownership', () => {
    it('should not be callable by non-owners', async () => {
      await expect(
        governor.connect(stranger).transferL1Ownership(stranger.address),
      ).to.be.revertedWith('Sender is not the L2 messenger')
    })

    it('should not be callable by L2 owner', async () => {
      const governorOwner = await governor.owner()
      assert.equal(governorOwner, owner.address)

      await expect(
        governor.connect(owner).transferL1Ownership(stranger.address),
      ).to.be.revertedWith('Sender is not the L2 messenger')
    })

    it('should be callable by current L1 owner', async () => {
      const currentL1Owner = await governor.l1Owner()
      await expect(
        governor
          .connect(crossdomainMessenger)
          .transferL1Ownership(newL1OwnerAddress),
      )
        .to.emit(governor, 'L1OwnershipTransferRequested')
        .withArgs(currentL1Owner, newL1OwnerAddress)
    })

    it('should be callable by current L1 owner to zero address', async () => {
      const currentL1Owner = await governor.l1Owner()
      await expect(
        governor
          .connect(crossdomainMessenger)
          .transferL1Ownership(ethers.constants.AddressZero),
      )
        .to.emit(governor, 'L1OwnershipTransferRequested')
        .withArgs(currentL1Owner, ethers.constants.AddressZero)
    })
  })

  describe('#acceptL1Ownership', () => {
    it('should not be callable by non pending-owners', async () => {
      await expect(
        governor.connect(crossdomainMessenger).acceptL1Ownership(),
      ).to.be.revertedWith('Must be proposed L1 owner')
    })

    it('should be callable by pending L1 owner', async () => {
      const currentL1Owner = await governor.l1Owner()
      await governor
        .connect(crossdomainMessenger)
        .transferL1Ownership(newL1OwnerAddress)
      await expect(
        governor.connect(newOwnerCrossdomainMessenger).acceptL1Ownership(),
      )
        .to.emit(governor, 'L1OwnershipTransferred')
        .withArgs(currentL1Owner, newL1OwnerAddress)

      const updatedL1Owner = await governor.l1Owner()
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
