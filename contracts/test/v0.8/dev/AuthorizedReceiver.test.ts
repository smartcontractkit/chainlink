import { ethers } from 'hardhat'
import { expect } from 'chai'
import { Contract, ContractFactory } from 'ethers'
import { Roles, getUsers } from '../../test-helpers/setup'

let authorizedReceiverFactory: ContractFactory
let roles: Roles

before(async () => {
  roles = (await getUsers()).roles

  authorizedReceiverFactory = await ethers.getContractFactory(
    'src/v0.8/tests/AuthorizedReceiverTestHelper.sol:AuthorizedReceiverTestHelper',
    roles.defaultAccount,
  )
})

describe('AuthorizedReceiverTestHelper', () => {
  let receiver: Contract

  beforeEach(async () => {
    receiver = await authorizedReceiverFactory
      .connect(roles.defaultAccount)
      .deploy()
  })

  describe('AuthorizedReceiver', () => {
    it('#setAuthorizedSenders', async () => {
      const { personas } = await getUsers()
      const addr1 = await personas.Carol.getAddress()
      const addr2 = await personas.Nancy.getAddress()

      await expect(receiver.setAuthorizedSenders([addr1, addr2])).not.to.be
        .reverted

      const senders = await receiver.callStatic.getAuthorizedSenders()
      expect(senders).to.be.deep.equal([addr1, addr2])
    })

    it('#setAuthorizedSenders emits AuthorizedSendersChanged', async () => {
      const { personas } = await getUsers()
      const owner = await personas.Default.getAddress()
      const addr1 = await personas.Carol.getAddress()
      const addr2 = await personas.Nancy.getAddress()

      await expect(receiver.setAuthorizedSenders([addr1, addr2]))
        .to.emit(receiver, 'AuthorizedSendersChanged')
        .withArgs([addr1, addr2], owner)
    })

    it('#setAuthorizedSenders empty list', async () => {
      await expect(receiver.setAuthorizedSenders([])).to.be.revertedWith(
        'EmptySendersList',
      )
    })

    it('#setAuthorizedSenders new senders', async () => {
      const { personas } = await getUsers()
      const addr1 = await personas.Carol.getAddress()
      const addr2 = await personas.Nancy.getAddress()
      const addr3 = await personas.Neil.getAddress()
      const addr4 = await personas.Ned.getAddress()

      await expect(receiver.setAuthorizedSenders([addr1, addr2])).not.to.be
        .reverted

      await expect(receiver.setAuthorizedSenders([addr3, addr4])).not.to.be
        .reverted

      const senders = await receiver.callStatic.getAuthorizedSenders()
      expect(senders).to.be.deep.equal([addr3, addr4])
    })

    it('#isAuthorizedSender', async () => {
      const { personas } = await getUsers()
      const addr1 = await personas.Carol.getAddress()
      const addr2 = await personas.Nancy.getAddress()

      await expect(receiver.setAuthorizedSenders([addr1])).not.to.be.reverted

      expect(await receiver.callStatic.isAuthorizedSender(addr1)).to.be.equal(
        true,
      )
      expect(await receiver.callStatic.isAuthorizedSender(addr2)).to.be.equal(
        false,
      )
    })
  })

  describe('#verifyValidateAuthorizedSender', () => {
    it('Should revert for empty state', async () => {
      await expect(
        receiver.verifyValidateAuthorizedSender(),
      ).to.be.revertedWith('UnauthorizedSender')
    })

    it('#validateAuthorizedSender modifier', async () => {
      const { personas } = await getUsers()

      await expect(
        receiver.setAuthorizedSenders([
          await personas.Carol.getAddress(),
          await personas.Nancy.getAddress(),
        ]),
      ).not.to.be.reverted

      expect(
        await receiver
          .connect(personas.Carol)
          .callStatic.verifyValidateAuthorizedSender(),
      ).to.be.equal(true)
      expect(
        await receiver
          .connect(personas.Nancy)
          .callStatic.verifyValidateAuthorizedSender(),
      ).to.be.equal(true)
      await expect(
        receiver.connect(personas.Neil).verifyValidateAuthorizedSender(),
      ).to.be.revertedWith('UnauthorizedSender')
    })

    it('Should revert setAuthorizedSenders if cannot set', async () => {
      const { personas } = await getUsers()

      await expect(receiver.changeSetAuthorizedSender(false)).not.to.be.reverted

      await expect(
        receiver.setAuthorizedSenders([
          await personas.Carol.getAddress(),
          await personas.Nancy.getAddress(),
        ]),
      ).to.be.revertedWith('NotAllowedToSetSenders')
    })
  })
})
