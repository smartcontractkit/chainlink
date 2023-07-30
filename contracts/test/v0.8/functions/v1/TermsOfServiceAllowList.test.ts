import { ethers } from 'hardhat'
import { expect } from 'chai'
import {
  getSetupFactory,
  FunctionsContracts,
  FunctionsRoles,
  acceptTermsOfService,
  accessControlMockPrivateKey,
} from './utils'

const setup = getSetupFactory()
let contracts: FunctionsContracts
let roles: FunctionsRoles

beforeEach(async () => {
  ;({ contracts, roles } = setup())
})

describe('ToS Access Control', () => {
  describe('Accepting', () => {
    it('can only be done with a valid proof', async () => {
      const messageHash = await contracts.accessControl.getMessageHash(
        roles.strangerAddress,
        roles.strangerAddress,
      )
      const proof = await roles.stranger.signMessage(
        ethers.utils.arrayify(messageHash),
      )
      await expect(
        contracts.accessControl
          .connect(roles.stranger)
          .acceptTermsOfService(
            roles.strangerAddress,
            roles.strangerAddress,
            proof,
          ),
      ).to.be.revertedWith('InvalidProof')
    })
    it('can be done by Externally Owned Accounts if recipient themself', async () => {
      await acceptTermsOfService(
        contracts.accessControl,
        roles.subOwner,
        roles.subOwnerAddress,
      )
      expect(
        await contracts.accessControl.hasAccess(roles.subOwnerAddress),
      ).to.equal(true)
    })
    it('cannot be done by Externally Owned Accounts if recipient another EoA', async () => {
      await expect(
        acceptTermsOfService(
          contracts.accessControl,
          roles.subOwner,
          roles.strangerAddress,
        ),
      ).to.be.revertedWith('InvalidProof')
    })
    it('can be done by Contract Accounts if recipient themself', async () => {
      const acceptorAddress = roles.consumerAddress
      const recipientAddress = contracts.client.address
      const messageHash = await contracts.accessControl.getMessageHash(
        acceptorAddress,
        recipientAddress,
      )
      const wallet = new ethers.Wallet(accessControlMockPrivateKey)
      const proof = await wallet.signMessage(ethers.utils.arrayify(messageHash))
      await contracts.client
        .connect(roles.consumer)
        .acceptTermsOfService(acceptorAddress, recipientAddress, proof)

      expect(
        await contracts.accessControl.hasAccess(recipientAddress),
      ).to.equal(true)
    })
    it('cannot be done by Contract Accounts that if they are not the recipient', async () => {
      const acceptorAddress = roles.consumerAddress
      const recipientAddress = contracts.coordinator.address
      const messageHash = await contracts.accessControl.getMessageHash(
        acceptorAddress,
        recipientAddress,
      )
      const wallet = new ethers.Wallet(accessControlMockPrivateKey)
      const proof = await wallet.signMessage(ethers.utils.arrayify(messageHash))
      await expect(
        contracts.client
          .connect(roles.consumer)
          .acceptTermsOfService(acceptorAddress, recipientAddress, proof),
      ).to.be.revertedWith('InvalidProof')
    })
  })

  describe('Blocking', () => {
    it('can only be done by the Router Owner', async () => {
      await expect(
        contracts.accessControl
          .connect(roles.stranger)
          .blockSender(roles.subOwnerAddress),
      ).to.be.revertedWith('OnlyCallableByRouterOwner')
    })
    it('removes the ability to re-accept the terms of service', async () => {
      await contracts.accessControl.blockSender(roles.subOwnerAddress)
      await expect(
        acceptTermsOfService(
          contracts.accessControl,
          roles.subOwner,
          roles.subOwnerAddress,
        ),
      ).to.be.revertedWith('RecipientIsBlocked')
    })
    it('removes the ability to manage subscriptions', async () => {
      await acceptTermsOfService(
        contracts.accessControl,
        roles.subOwner,
        roles.subOwnerAddress,
      )
      await contracts.accessControl.blockSender(roles.subOwnerAddress)
      await expect(
        contracts.router.connect(roles.subOwner).createSubscription(),
      ).to.be.revertedWith('SenderMustAcceptTermsOfService')
    })
  })
})
