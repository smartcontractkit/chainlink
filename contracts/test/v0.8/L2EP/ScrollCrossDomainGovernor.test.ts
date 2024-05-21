import { ethers } from 'hardhat'
import { assert, expect } from 'chai'
import etherslib, { Contract, ContractFactory } from 'ethers'
import { SignerWithAddress } from '@nomiclabs/hardhat-ethers/signers'
import { publicAbi, stripHexPrefix } from '../../test-helpers/helpers'

let owner: SignerWithAddress
let stranger: SignerWithAddress
let l1OwnerAddress: string
let newL1OwnerAddress: string
let governorFactory: ContractFactory
let greeterFactory: ContractFactory
let multisendFactory: ContractFactory
let crossDomainMessengerFactory: ContractFactory
let crossDomainMessenger: Contract
let governor: Contract
let greeter: Contract
let multisend: Contract

before(async () => {
  const accounts = await ethers.getSigners()
  owner = accounts[0]
  stranger = accounts[1]

  // governor config
  l1OwnerAddress = owner.address
  newL1OwnerAddress = stranger.address

  // Contract factories
  governorFactory = await ethers.getContractFactory(
    'src/v0.8/l2ep/dev/scroll/ScrollCrossDomainGovernor.sol:ScrollCrossDomainGovernor',
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
  crossDomainMessengerFactory = await ethers.getContractFactory(
    'src/v0.8/vendor/MockScrollCrossDomainMessenger.sol:MockScrollCrossDomainMessenger',
  )
})

describe('ScrollCrossDomainGovernor', () => {
  beforeEach(async () => {
    crossDomainMessenger =
      await crossDomainMessengerFactory.deploy(l1OwnerAddress)
    governor = await governorFactory.deploy(
      crossDomainMessenger.address,
      l1OwnerAddress,
    )
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
      assert.equal(response, crossDomainMessenger.address)
    })

    it('should set the typeAndVersion correctly', async () => {
      const response = await governor.typeAndVersion()
      assert.equal(response, 'ScrollCrossDomainGovernor 1.0.0')
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
      const forwardData = governorFactory.interface.encodeFunctionData(
        'forward',
        [greeter.address, setGreetingData],
      )
      await crossDomainMessenger // Simulate cross-chain message
        .connect(stranger)
        ['sendMessage(address,uint256,bytes,uint256)'](
          governor.address, // target
          0, // value
          forwardData, // message
          0, // gasLimit
        )

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
      const forwardData = governorFactory.interface.encodeFunctionData(
        'forward',
        [greeter.address, setGreetingData],
      )
      await expect(
        crossDomainMessenger // Simulate cross-chain message
          .connect(stranger)
          ['sendMessage(address,uint256,bytes,uint256)'](
            governor.address, // target
            0, // value
            forwardData, // message
            0, // gasLimit
          ),
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
      const forwardData = governorFactory.interface.encodeFunctionData(
        'forwardDelegate',
        [multisend.address, multisendData],
      )

      await crossDomainMessenger // Simulate cross-chain message
        .connect(stranger)
        ['sendMessage(address,uint256,bytes,uint256)'](
          governor.address, // target
          0, // value
          forwardData, // message
          0, // gasLimit
        )

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
      const forwardData = governorFactory.interface.encodeFunctionData(
        'forwardDelegate',
        [multisend.address, multisendData],
      )

      await expect(
        crossDomainMessenger // Simulate cross-chain message
          .connect(stranger)
          ['sendMessage(address,uint256,bytes,uint256)'](
            governor.address, // target
            0, // value
            forwardData, // message
            0, // gasLimit
          ),
      ).to.be.revertedWith('Governor delegatecall reverted')

      const greeting = await greeter.greeting()
      assert.equal(greeting, '') // Unchanged
    })

    it('should bubble up revert when contract call reverts', async () => {
      const triggerRevertData =
        greeterFactory.interface.encodeFunctionData('triggerRevert')
      const forwardData = governorFactory.interface.encodeFunctionData(
        'forwardDelegate',
        [greeter.address, triggerRevertData],
      )

      await expect(
        crossDomainMessenger // Simulate cross-chain message
          .connect(stranger)
          ['sendMessage(address,uint256,bytes,uint256)'](
            governor.address, // target
            0, // value
            forwardData, // message
            0, // gasLimit
          ),
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
      const forwardData = governorFactory.interface.encodeFunctionData(
        'transferL1Ownership',
        [newL1OwnerAddress],
      )

      await expect(
        crossDomainMessenger // Simulate cross-chain message
          .connect(stranger)
          ['sendMessage(address,uint256,bytes,uint256)'](
            governor.address, // target
            0, // value
            forwardData, // message
            0, // gasLimit
          ),
      )
        .to.emit(governor, 'L1OwnershipTransferRequested')
        .withArgs(currentL1Owner, newL1OwnerAddress)
    })

    it('should be callable by current L1 owner to zero address', async () => {
      const currentL1Owner = await governor.l1Owner()
      const forwardData = governorFactory.interface.encodeFunctionData(
        'transferL1Ownership',
        [ethers.constants.AddressZero],
      )

      await expect(
        crossDomainMessenger // Simulate cross-chain message
          .connect(stranger)
          ['sendMessage(address,uint256,bytes,uint256)'](
            governor.address, // target
            0, // value
            forwardData, // message
            0, // gasLimit
          ),
      )
        .to.emit(governor, 'L1OwnershipTransferRequested')
        .withArgs(currentL1Owner, ethers.constants.AddressZero)
    })
  })

  describe('#acceptL1Ownership', () => {
    it('should not be callable by non pending-owners', async () => {
      const forwardData = governorFactory.interface.encodeFunctionData(
        'acceptL1Ownership',
        [],
      )
      await expect(
        crossDomainMessenger // Simulate cross-chain message
          .connect(stranger)
          ['sendMessage(address,uint256,bytes,uint256)'](
            governor.address, // target
            0, // value
            forwardData, // message
            0, // gasLimit
          ),
      ).to.be.revertedWith('Must be proposed L1 owner')
    })

    it('should be callable by pending L1 owner', async () => {
      const currentL1Owner = await governor.l1Owner()

      // Transfer ownership
      const forwardTransferData = governorFactory.interface.encodeFunctionData(
        'transferL1Ownership',
        [newL1OwnerAddress],
      )
      await crossDomainMessenger // Simulate cross-chain message
        .connect(stranger)
        ['sendMessage(address,uint256,bytes,uint256)'](
          governor.address, // target
          0, // value
          forwardTransferData, // message
          0, // gasLimit
        )

      const forwardAcceptData = governorFactory.interface.encodeFunctionData(
        'acceptL1Ownership',
        [],
      )
      // Simulate cross-chain message from another sender
      await crossDomainMessenger._setMockMessageSender(newL1OwnerAddress)

      await expect(
        crossDomainMessenger // Simulate cross-chain message
          .connect(stranger)
          ['sendMessage(address,uint256,bytes,uint256)'](
            governor.address, // target
            0, // value
            forwardAcceptData, // message
            0, // gasLimit
          ),
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
  const dataBuffer = Buffer.from(stripHexPrefix(data), 'hex')
  const types = ['uint8', 'address', 'uint256', 'uint256', 'bytes']
  const values = [operation, to, value, dataBuffer.length, dataBuffer]
  const encoded = ethers.utils.solidityPack(types, values)
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
  for (const transaction of transactions) {
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
