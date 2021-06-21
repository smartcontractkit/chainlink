import { contract, matchers, setup } from '@chainlink/test-helpers'
import { assert } from 'chai'
import { ethers, utils } from 'ethers'
import { ContractReceipt } from 'ethers/contract'
import { GetterSetter__factory } from '../../ethers/v0.4/factories/GetterSetter__factory'
import { AuthorizedForwarder__factory } from '../../ethers/v0.7/factories/AuthorizedForwarder__factory'

const getterSetterFactory = new GetterSetter__factory()
const forwarderFactory = new AuthorizedForwarder__factory()
const linkTokenFactory = new contract.LinkToken__factory()

let roles: setup.Roles
const provider = setup.provider()
const zeroAddress = ethers.constants.AddressZero

beforeAll(async () => {
  const users = await setup.users(provider)

  roles = users.roles
})

describe('AuthorizedForwarder', () => {
  let link: contract.Instance<contract.LinkToken__factory>
  let forwarder: contract.Instance<AuthorizedForwarder__factory>
  const deployment = setup.snapshot(provider, async () => {
    link = await linkTokenFactory.connect(roles.defaultAccount).deploy()
    forwarder = await forwarderFactory
      .connect(roles.defaultAccount)
      .deploy(link.address, roles.defaultAccount.address, zeroAddress, '0x')
  })

  beforeEach(async () => {
    await deployment()
  })

  it('has a limited public interface', () => {
    matchers.publicAbi(forwarder, [
      'forward',
      'getAuthorizedSenders',
      'getChainlinkToken',
      'isAuthorizedSender',
      'ownerForward',
      'setAuthorizedSenders',
      'transferOwnershipWithMessage',
      // ConfirmedOwner
      'transferOwnership',
      'acceptOwnership',
      'owner',
    ])
  })

  describe('deployment', () => {
    it('sets the correct link token', async () => {
      assert.equal(await forwarder.getChainlinkToken(), link.address)
    })

    it('sets no authorized senders', async () => {
      const senders = await forwarder.getAuthorizedSenders()
      assert.equal(senders.length, 0)
    })
  })

  describe('#setAuthorizedSenders', () => {
    let newSenders: string[]
    let receipt: ContractReceipt
    describe('when called by the owner', () => {
      describe('setting 3 authorized senders', () => {
        beforeEach(async () => {
          newSenders = [
            roles.oracleNode1.address,
            roles.oracleNode2.address,
            roles.oracleNode3.address,
          ]
          const tx = await forwarder
            .connect(roles.defaultAccount)
            .setAuthorizedSenders(newSenders)
          receipt = await tx.wait()
        })

        it('adds the authorized nodes', async () => {
          const authorizedSenders = await forwarder.getAuthorizedSenders()
          assert.equal(newSenders.length, authorizedSenders.length)
          for (let i = 0; i < authorizedSenders.length; i++) {
            assert.equal(authorizedSenders[i], newSenders[i])
          }
        })

        it('emits an event', async () => {
          assert.equal(receipt.events?.length, 1)
          const responseEvent = receipt.events?.[0]
          assert.equal(responseEvent?.event, 'AuthorizedSendersChanged')
          const encodedSenders = utils.defaultAbiCoder.encode(
            ['address[]', 'address'],
            [newSenders, roles.defaultAccount.address],
          )
          assert.equal(responseEvent?.data, encodedSenders)
        })

        it('replaces the authorized nodes', async () => {
          const newSenders = await forwarder
            .connect(roles.defaultAccount)
            .getAuthorizedSenders()
          assert.notIncludeOrderedMembers(newSenders, [
            roles.oracleNode.address,
          ])
        })

        afterAll(async () => {
          await forwarder
            .connect(roles.defaultAccount)
            .setAuthorizedSenders([roles.oracleNode.address])
        })
      })

      describe('setting 0 authorized senders', () => {
        beforeEach(async () => {
          newSenders = []
        })

        it('reverts with a minimum senders message', async () => {
          await matchers.evmRevert(async () => {
            await forwarder
              .connect(roles.defaultAccount)
              .setAuthorizedSenders(newSenders),
              'Must have at least 1 authorized sender'
          })
        })
      })
    })

    describe('when called by a non-owner', () => {
      it('cannot add an authorized node', async () => {
        await matchers.evmRevert(async () => {
          await forwarder
            .connect(roles.stranger)
            .setAuthorizedSenders([roles.stranger.address])
          ;('Only callable by owner')
        })
      })
    })
  })

  describe('#forward', () => {
    const bytes = utils.hexlify(utils.randomBytes(100))
    const payload = getterSetterFactory.interface.functions.setBytes.encode([
      bytes,
    ])
    let mock: contract.Instance<GetterSetter__factory>

    beforeEach(async () => {
      mock = await getterSetterFactory.connect(roles.defaultAccount).deploy()
    })

    describe('when called by an unauthorized node', () => {
      it('reverts', async () => {
        await matchers.evmRevert(async () => {
          await forwarder.connect(roles.stranger).forward(mock.address, payload)
        })
      })
    })

    describe('when called by an authorized node', () => {
      beforeEach(async () => {
        await forwarder
          .connect(roles.defaultAccount)
          .setAuthorizedSenders([roles.defaultAccount.address])
      })

      describe('when attempting to forward to the link token', () => {
        it('reverts', async () => {
          const { sighash } = linkTokenFactory.interface.functions.name // any Link Token function
          await matchers.evmRevert(async () => {
            await forwarder
              .connect(roles.defaultAccount)
              .forward(link.address, sighash)
          })
        })
      })

      describe('when forwarding to any other address', () => {
        it('forwards the data', async () => {
          const tx = await forwarder
            .connect(roles.defaultAccount)
            .forward(mock.address, payload)
          await tx.wait()
          assert.equal(await mock.getBytes(), bytes)
        })

        it('perceives the message is sent by the AuthorizedForwarder', async () => {
          const tx = await forwarder
            .connect(roles.defaultAccount)
            .forward(mock.address, payload)
          const receipt = await tx.wait()
          const log: any = receipt.logs?.[0]
          const logData = mock.interface.events.SetBytes.decode(
            log.data,
            log.topics,
          )
          assert.equal(utils.getAddress(logData.from), forwarder.address)
        })
      })
    })
  })

  describe('#transferOwnershipWithMessage', () => {
    const message = '0x42'

    describe('when called by a non-owner', () => {
      it('reverts', async () => {
        await matchers.evmRevert(async () => {
          await forwarder
            .connect(roles.stranger)
            .transferOwnershipWithMessage(roles.stranger.address, message),
            'Only callable by owner'
        })
      })
    })

    describe('when called by the owner', () => {
      it('calls the normal ownership transfer proposal', async () => {
        const tx = await forwarder
          .connect(roles.defaultAccount)
          .transferOwnershipWithMessage(roles.stranger.address, message)
        const receipt = await tx.wait()

        assert.equal(receipt?.events?.[0]?.event, 'OwnershipTransferRequested')
        assert.equal(receipt?.events?.[0]?.address, forwarder.address)
        assert.equal(
          receipt?.events?.[0]?.args?.[0],
          roles.defaultAccount.address,
        )
        assert.equal(receipt?.events?.[0]?.args?.[1], roles.stranger.address)
      })

      it('calls the normal ownership transfer proposal', async () => {
        const tx = await forwarder
          .connect(roles.defaultAccount)
          .transferOwnershipWithMessage(roles.stranger.address, message)
        const receipt = await tx.wait()

        assert.equal(
          receipt?.events?.[1]?.event,
          'OwnershipTransferRequestedWithMessage',
        )
        assert.equal(receipt?.events?.[1]?.address, forwarder.address)
        assert.equal(
          receipt?.events?.[1]?.args?.[0],
          roles.defaultAccount.address,
        )
        assert.equal(receipt?.events?.[1]?.args?.[1], roles.stranger.address)
        assert.equal(receipt?.events?.[1]?.args?.[2], message)
      })
    })
  })

  describe('#ownerForward', () => {
    const bytes = utils.hexlify(utils.randomBytes(100))
    const payload = getterSetterFactory.interface.functions.setBytes.encode([
      bytes,
    ])
    let mock: contract.Instance<GetterSetter__factory>

    beforeEach(async () => {
      mock = await getterSetterFactory.connect(roles.defaultAccount).deploy()
    })

    describe('when called by a non-owner', () => {
      it('reverts', async () => {
        await matchers.evmRevert(async () => {
          await forwarder
            .connect(roles.stranger)
            .ownerForward(mock.address, payload)
        })
      })
    })

    describe('when called by owner', () => {
      describe('when attempting to forward to the link token', () => {
        it('does not revert', async () => {
          const { sighash } = linkTokenFactory.interface.functions.name // any Link Token function
          await forwarder
            .connect(roles.defaultAccount)
            .ownerForward(link.address, sighash)
        })
      })

      describe('when forwarding to any other address', () => {
        it('forwards the data', async () => {
          const tx = await forwarder
            .connect(roles.defaultAccount)
            .ownerForward(mock.address, payload)
          await tx.wait()
          assert.equal(await mock.getBytes(), bytes)
        })

        it('perceives the message is sent by the Operator', async () => {
          const tx = await forwarder
            .connect(roles.defaultAccount)
            .ownerForward(mock.address, payload)
          const receipt = await tx.wait()
          const log: any = receipt.logs?.[0]
          const logData = mock.interface.events.SetBytes.decode(
            log.data,
            log.topics,
          )
          assert.equal(utils.getAddress(logData.from), forwarder.address)
        })
      })
    })
  })
})
