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
let forwarderFactory: ContractFactory
let greeterFactory: ContractFactory
// let multisendFactory: ContractFactory
let forwarder: Contract
let greeter: Contract
// let multisend: Contract

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
    value: ethers.utils.parseEther("1.0")
  })

  // Contract factories
  forwarderFactory = await ethers.getContractFactory(
    'src/v0.8/dev/ArbitrumCrossDomainForwarder.sol:ArbitrumCrossDomainForwarder',
    owner,
  )
  greeterFactory = await ethers.getContractFactory(
    'src/v0.8/tests/Greeter.sol:Greeter',
    owner,
  )
  // multisendFactory = await ethers.getContractFactory(
  //   'src/v0.8/tests/Multisend.sol:Multisend',
  //   owner,
  // )
})

describe('ArbitrumCrossDomainForwarder', () => {
  beforeEach(async () => {
    forwarder = await forwarderFactory.deploy(l1OwnerAddress)
    greeter = await greeterFactory.deploy(forwarder.address)
    // multisend = await multisendFactory.deploy()
  })

  it('has a limited public interface [ @skip-coverage ]', async () => {
    publicAbi(forwarder, [
      'typeAndVersion',
      'crossDomainMessenger',
      'forward',
      'forwardDelegate',
      'l1Owner',
      'transferL1Ownership',
      // ConfirmedOwner methods:
      'acceptOwnership',
      'owner',
      'transferOwnership',
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
    it('should not be callable by unknown crossdomain messenger address', async () => {
      await expect(
        forwarder.connect(stranger).forward(greeter.address, '0x'),
      ).to.be.revertedWith('Sender is not the L2 messenger')
    })

    it('should be callable by crossdomain messenger address', async () => {
      const newGreeting = 'hello'
      const setGreetingData = greeterFactory.interface.encodeFunctionData("setGreeting", [newGreeting])
      await forwarder.connect(crossdomainMessenger).forward(greeter.address, setGreetingData)

      const updatedGreeting = await greeter.greeting()
      assert.equal(updatedGreeting, newGreeting)
    })
  })

  //   TODO: test forwardDelegate()
})
