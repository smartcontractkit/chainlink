import {
  contract,
  // helpers as h,
  matchers,
  // oracle,
  setup,
} from '@chainlink/test-helpers'
import { assert } from 'chai'
import { utils } from 'ethers'
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
      const senders = await operatorForwarder.getAuthorizedSenders();
      assert.equal(senders.length, 0);
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
