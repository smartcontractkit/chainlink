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
    it('can only be done with a valid signature', async () => {
      const message = await contracts.accessControl.getMessage(
        roles.strangerAddress,
        roles.strangerAddress,
      )
      const flatSignature = await roles.stranger.signMessage(
        ethers.utils.arrayify(message),
      )
      const { r, s, v } = ethers.utils.splitSignature(flatSignature)
      await expect(
        contracts.accessControl
          .connect(roles.stranger)
          .acceptTermsOfService(
            roles.strangerAddress,
            roles.strangerAddress,
            r,
            s,
            v,
          ),
      ).to.be.revertedWith('InvalidSignature')
    })
    it('can be done by Externally Owned Accounts if recipient themself', async () => {
      await acceptTermsOfService(
        contracts.accessControl,
        roles.subOwner,
        roles.subOwnerAddress,
      )
      expect(
        await contracts.accessControl.hasAccess(roles.subOwnerAddress, '0x'),
      ).to.equal(true)
    })
    it('cannot be done by Externally Owned Accounts if recipient another EoA', async () => {
      await expect(
        acceptTermsOfService(
          contracts.accessControl,
          roles.subOwner,
          roles.strangerAddress,
        ),
      ).to.be.revertedWith('InvalidUsage')
    })
    it('can be done by Contract Accounts if recipient themself', async () => {
      const acceptorAddress = roles.consumerAddress
      const recipientAddress = contracts.client.address
      const message = await contracts.accessControl.getMessage(
        acceptorAddress,
        recipientAddress,
      )
      const wallet = new ethers.Wallet(accessControlMockPrivateKey)
      const flatSignature = await wallet.signMessage(
        ethers.utils.arrayify(message),
      )
      const { r, s, v } = ethers.utils.splitSignature(flatSignature)
      await contracts.client
        .connect(roles.consumer)
        .acceptTermsOfService(acceptorAddress, recipientAddress, r, s, v)

      expect(
        await contracts.accessControl.hasAccess(recipientAddress, '0x'),
      ).to.equal(true)
    })
    it('cannot be done by Contract Accounts that if they are not the recipient', async () => {
      const acceptorAddress = roles.consumerAddress
      const recipientAddress = contracts.coordinator.address
      const message = await contracts.accessControl.getMessage(
        acceptorAddress,
        recipientAddress,
      )
      const wallet = new ethers.Wallet(accessControlMockPrivateKey)
      const flatSignature = await wallet.signMessage(
        ethers.utils.arrayify(message),
      )
      const { r, s, v } = ethers.utils.splitSignature(flatSignature)
      await expect(
        contracts.client
          .connect(roles.consumer)
          .acceptTermsOfService(acceptorAddress, recipientAddress, r, s, v),
      ).to.be.revertedWith('InvalidUsage')
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
