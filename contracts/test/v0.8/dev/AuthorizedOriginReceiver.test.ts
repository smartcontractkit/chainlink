import { ethers } from 'hardhat'
import { expect } from 'chai'
import { Contract, ContractFactory } from 'ethers'
import { Roles, getUsers } from '../../test-helpers/setup'

let authorizedOriginReceiverFactory: ContractFactory
let roles: Roles

before(async () => {
  roles = (await getUsers()).roles

  authorizedOriginReceiverFactory = await ethers.getContractFactory(
    'src/v0.8/tests/AuthorizedOriginReceiverTestHelper.sol:AuthorizedOriginReceiverTestHelper',
    roles.defaultAccount,
  )
})

describe('AuthorizedOriginReceiverTestHelper', () => {
  let receiver: Contract

  beforeEach(async () => {
    receiver = await authorizedOriginReceiverFactory
      .connect(roles.defaultAccount)
      .deploy()
  })

  describe('AuthorizedOriginReceiver', () => {
    it('#addAuthorizedSenders', async () => {
      const { personas } = await getUsers()
      const addr1 = await personas.Carol.getAddress()
      const addr2 = await personas.Nancy.getAddress()

      await expect(receiver.addAuthorizedSenders([addr1, addr2])).not.to.be
        .reverted

      const senders = await receiver.callStatic.getAuthorizedSenders()
      expect(senders).to.be.deep.equal([addr1, addr2])
    })

    it('#addAuthorizedSenders emits AuthorizedSendersChanged', async () => {
      const { personas } = await getUsers()
      const owner = await personas.Default.getAddress()
      const addr1 = await personas.Carol.getAddress()
      const addr2 = await personas.Nancy.getAddress()

      await expect(receiver.addAuthorizedSenders([addr1, addr2]))
        .to.emit(receiver, 'AuthorizedSendersChanged')
        .withArgs([addr1, addr2], owner)
    })

    it('#addAuthorizedSenders empty list', async () => {
      await expect(receiver.addAuthorizedSenders([])).to.be.revertedWith(
        'EmptySendersList',
      )
    })

    it('#addAuthorizedSenders new senders', async () => {
      const { personas } = await getUsers()
      const addr1 = await personas.Carol.getAddress()
      const addr2 = await personas.Nancy.getAddress()
      const addr3 = await personas.Neil.getAddress()
      const addr4 = await personas.Ned.getAddress()

      await expect(receiver.addAuthorizedSenders([addr1, addr2])).not.to.be
        .reverted

      await expect(receiver.addAuthorizedSenders([addr3, addr4])).not.to.be
        .reverted

      const senders = await receiver.callStatic.getAuthorizedSenders()
      expect(senders).to.be.deep.equal([addr1, addr2, addr3, addr4])
    })

    it('#addAuthorizedSenders removes duplicate new senders', async () => {
      const { personas } = await getUsers()
      const addr1 = await personas.Carol.getAddress()
      const addr2 = await personas.Nancy.getAddress()
      const addr3 = await personas.Neil.getAddress()

      await expect(receiver.addAuthorizedSenders([addr1, addr2])).not.to.be
        .reverted

      await expect(receiver.addAuthorizedSenders([addr2, addr3])).not.to.be
        .reverted

      const senders = await receiver.callStatic.getAuthorizedSenders()
      expect(senders).to.be.deep.equal([addr1, addr2, addr3])
    })

    it('#remove AuthorizedSenders', async () => {
      const { personas } = await getUsers()
      const addr1 = await personas.Carol.getAddress()
      const addr2 = await personas.Nancy.getAddress()

      await expect(receiver.addAuthorizedSenders([addr1, addr2])).not.to.be
        .reverted

      const senders = await receiver.callStatic.getAuthorizedSenders()
      expect(senders).to.be.deep.equal([addr1, addr2])

      await expect(receiver.removeAuthorizedSenders([addr1, addr2])).not.to.be
        .reverted
      const sendersAfterRemove =
        await receiver.callStatic.getAuthorizedSenders()
      expect(sendersAfterRemove).to.be.deep.equal([])
    })

    it('#addAuthorizedSenders', async () => {
      const { personas } = await getUsers()
      const addr1 = await personas.Carol.getAddress()
      const addr2 = await personas.Nancy.getAddress()

      await expect(receiver.addAuthorizedSenders([addr1, addr2])).not.to.be
        .reverted

      const senders = await receiver.callStatic.getAuthorizedSenders()
      expect(senders).to.be.deep.equal([addr1, addr2])
    })

    it('#isAuthorizedSender', async () => {
      const { personas } = await getUsers()
      const addr1 = await personas.Carol.getAddress()
      const addr2 = await personas.Nancy.getAddress()

      await expect(receiver.addAuthorizedSenders([addr1])).not.to.be.reverted

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
        receiver.addAuthorizedSenders([
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

    it('Should revert addAuthorizedSenders if cannot set', async () => {
      const { personas } = await getUsers()

      await expect(receiver.changeSetAuthorizedSender(false)).not.to.be.reverted

      await expect(
        receiver.addAuthorizedSenders([
          await personas.Carol.getAddress(),
          await personas.Nancy.getAddress(),
        ]),
      ).to.be.revertedWith('NotAllowedToSetSenders')
    })
  })
})
