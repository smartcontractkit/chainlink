import { contract, matchers, setup } from '@chainlink/test-helpers'
import { assert } from 'chai'
import { utils } from 'ethers'
import { ContractReceipt } from 'ethers/contract'
import { GetterSetter__factory } from '../../ethers/v0.4/factories/GetterSetter__factory'
import { OperatorForwarder__factory } from '../../ethers/v0.7/factories/OperatorForwarder__factory'

const getterSetterFactory = new GetterSetter__factory()
const operatorForwarderFactory = new OperatorForwarder__factory()
const linkTokenFactory = new contract.LinkToken__factory()

let roles: setup.Roles
const provider = setup.provider()

beforeAll(async () => {
  const users = await setup.users(provider)

  roles = users.roles
})

describe('OperatorForwarder', () => {
  let link: contract.Instance<contract.LinkToken__factory>
  let operatorForwarder: contract.Instance<OperatorForwarder__factory>
  const deployment = setup.snapshot(provider, async () => {
    link = await linkTokenFactory.connect(roles.defaultAccount).deploy()
    operatorForwarder = await operatorForwarderFactory
      .connect(roles.defaultAccount)
      .deploy(link.address)
  })

  beforeEach(async () => {
    await deployment()
  })

  it('has a limited public interface', () => {
    matchers.publicAbi(operatorForwarder, [
      'transferOwnership',
      'acceptOwnership',
      'owner',
      // OperatorForwarder
      'linkAddr',
      'setAuthorizedSenders',
      'getAuthorizedSenders',
      'forward',
    ])
  })

  describe('deployment', () => {
    it('sets the correct link token', async () => {
      const forwarderLink = await operatorForwarder.linkAddr()
      assert.equal(forwarderLink, link.address)
    })

    it('sets no authorized senders', async () => {
      const senders = await operatorForwarder.getAuthorizedSenders()
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
          const tx = await operatorForwarder
            .connect(roles.defaultAccount)
            .setAuthorizedSenders(newSenders)
          receipt = await tx.wait()
        })

        it('adds the authorized nodes', async () => {
          const authorizedSenders = await operatorForwarder.getAuthorizedSenders()
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
            ['address[]'],
            [newSenders],
          )
          assert.equal(responseEvent?.data, encodedSenders)
        })

        it('replaces the authorized nodes', async () => {
          const newSenders = await operatorForwarder
            .connect(roles.defaultAccount)
            .getAuthorizedSenders()
          assert.notIncludeOrderedMembers(newSenders, [
            roles.oracleNode.address,
          ])
        })

        afterAll(async () => {
          await operatorForwarder
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
            await operatorForwarder
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
          await operatorForwarder
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
          await operatorForwarder
            .connect(roles.stranger)
            .forward(mock.address, payload)
        })
      })
    })

    describe('when called by an authorized node', () => {
      beforeEach(async () => {
        await operatorForwarder
          .connect(roles.defaultAccount)
          .setAuthorizedSenders([roles.defaultAccount.address])
      })

      describe('when attempting to forward to the link token', () => {
        it('reverts', async () => {
          const { sighash } = linkTokenFactory.interface.functions.name // any Link Token function
          await matchers.evmRevert(async () => {
            await operatorForwarder
              .connect(roles.defaultAccount)
              .forward(link.address, sighash)
          })
        })
      })

      describe('when forwarding to any other address', () => {
        it('forwards the data', async () => {
          const tx = await operatorForwarder
            .connect(roles.defaultAccount)
            .forward(mock.address, payload)
          await tx.wait()
          assert.equal(await mock.getBytes(), bytes)
        })

        it('perceives the message is sent by the OperatorForwarder', async () => {
          const tx = await operatorForwarder
            .connect(roles.defaultAccount)
            .forward(mock.address, payload)
          const receipt = await tx.wait()
          const log: any = receipt.logs?.[0]
          const logData = mock.interface.events.SetBytes.decode(
            log.data,
            log.topics,
          )
          assert.equal(
            utils.getAddress(logData.from),
            operatorForwarder.address,
          )
        })
      })
    })
  })
})
