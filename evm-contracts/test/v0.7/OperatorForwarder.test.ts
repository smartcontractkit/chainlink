import {
  contract,
  helpers as h,
  matchers,
  // oracle,
  setup,
} from '@chainlink/test-helpers'
import { assert } from 'chai'
import { utils } from 'ethers'
import { GetterSetter__factory } from '../../ethers/v0.4/factories/GetterSetter__factory'
import { OperatorForwarder__factory } from '../../ethers/v0.7/factories/OperatorForwarder__factory'
import { OperatorForwarderDeployer__factory } from '../../ethers/v0.7/factories/OperatorForwarderDeployer__factory'

const getterSetterFactory = new GetterSetter__factory()
const operatorForwarderFactory = new OperatorForwarder__factory()
const operatorForwarderDeployerFactory = new OperatorForwarderDeployer__factory()
const linkTokenFactory = new contract.LinkToken__factory()

let roles: setup.Roles
const provider = setup.provider()

beforeAll(async () => {
  const users = await setup.users(provider)

  roles = users.roles
})

describe('OperatorForwarder', () => {
  let authorizedSenders: string[]
  let link: contract.Instance<contract.LinkToken__factory>
  let operatorForwarderDeployer: contract.Instance<OperatorForwarderDeployer__factory>
  let operatorForwarder: contract.Instance<OperatorForwarder__factory>
  const deployment = setup.snapshot(provider, async () => {
    link = await linkTokenFactory.connect(roles.defaultAccount).deploy()
    authorizedSenders = [roles.oracleNode2.address, roles.oracleNode3.address]
    operatorForwarderDeployer = await operatorForwarderDeployerFactory
      .connect(roles.defaultAccount)
      .deploy(link.address, authorizedSenders)
    const tx = await operatorForwarderDeployer.createForwarder()
    const receipt = await tx.wait()
    const event = h.findEventIn(
      receipt,
      operatorForwarderDeployer.interface.events.ForwarderDeployed,
    )
    operatorForwarder = await operatorForwarderFactory
      .connect(roles.defaultAccount)
      .attach(event?.args?.[0])
  })

  beforeEach(async () => {
    await deployment()
  })

  it('has a limited public interface', () => {
    matchers.publicAbi(operatorForwarder, [
      'authorizedSender1',
      'authorizedSender2',
      'authorizedSender3',
      'linkAddr',
      'forward',
    ])
  })

  describe('deployment', () => {
    it('sets the correct link token', async () => {
      const forwarderLink = await operatorForwarder.linkAddr()
      assert.equal(forwarderLink, link.address)
    })

    it('sets the correct authorized senders', async () => {
      const auth1 = await operatorForwarder.authorizedSender1()
      const auth2 = await operatorForwarder.authorizedSender2()
      const auth3 = await operatorForwarder.authorizedSender3()
      assert.equal(auth1, roles.defaultAccount.address)
      assert.equal(auth2, authorizedSenders[0])
      assert.equal(auth3, authorizedSenders[1])
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
