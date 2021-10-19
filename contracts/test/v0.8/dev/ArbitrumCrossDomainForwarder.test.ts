import { ethers } from 'hardhat'
import { assert } from 'chai'
import { BigNumber, Contract, ContractFactory } from 'ethers'
import { SignerWithAddress } from '@nomiclabs/hardhat-ethers/signers'
import { publicAbi } from '../../test-helpers/helpers'

let owner: SignerWithAddress
// let stranger: SignerWithAddress
let l1OwnerAddress: string
let crossdomainMessengerAddress: string
let forwarderFactory: ContractFactory
let forwarder: Contract

function applyL1ToL2Alias(l1Address: string): string {
  return ethers.utils.getAddress(
    BigNumber.from(l1Address)
      .add('0x1111000000000000000000000000000000001111')
      .toHexString()
      .replace('0x01', '0x'),
  )
}

before(async () => {
  const accounts = await ethers.getSigners()
  owner = accounts[0]
  // stranger = accounts[1]
  l1OwnerAddress = owner.address
  crossdomainMessengerAddress = applyL1ToL2Alias(l1OwnerAddress) // TODO: util function
  forwarderFactory = await ethers.getContractFactory(
    'src/v0.8/dev/ArbitrumCrossDomainForwarder.sol:ArbitrumCrossDomainForwarder',
    owner,
  )
})

describe('ArbitrumCrossDomainForwarder', () => {
  beforeEach(async () => {
    forwarder = await forwarderFactory.deploy(l1OwnerAddress)
    forwarder = await forwarder.deployed()
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
      assert.equal(response, crossdomainMessengerAddress)
    })

    // TODO: test l1 Owner
  })

  //   TODO: test forward()
  //   TODO: test forwardDelegate()

  //   describe('#raiseFlag', () => {
  //     describe('when called by the owner', () => {
  //       it('updates the warning flag', async () => {
  //         assert.equal(false, await flags.getFlag(consumer.address))

  //         await flags.connect(personas.Nelly).raiseFlag(consumer.address)

  //         assert.equal(true, await flags.getFlag(consumer.address))
  //       })

  //       it('emits an event log', async () => {
  //         await expect(flags.connect(personas.Nelly).raiseFlag(consumer.address))
  //           .to.emit(flags, 'FlagRaised')
  //           .withArgs(consumer.address)
  //       })

  //       describe('if a flag has already been raised', () => {
  //         beforeEach(async () => {
  //           await flags.connect(personas.Nelly).raiseFlag(consumer.address)
  //         })

  //         it('emits an event log', async () => {
  //           const tx = await flags
  //             .connect(personas.Nelly)
  //             .raiseFlag(consumer.address)
  //           const receipt = await tx.wait()
  //           assert.equal(0, receipt.events?.length)
  //         })
  //       })
  //     })

  //     describe('when called by a non-enabled setter', () => {
  //       it('reverts', async () => {
  //         await expect(
  //           flags.connect(personas.Neil).raiseFlag(consumer.address),
  //         ).to.be.revertedWith('Not allowed to raise flags')
  //       })
  //     })
})
