import { contract, matchers, setup } from '@chainlink/test-helpers'
import { assert } from 'chai'
import { utils } from 'ethers'
import { GetterSetterFactory } from '../../ethers/v0.4/GetterSetterFactory'
import { OperatorProxyFactory } from '../../ethers/v0.7/OperatorProxyFactory'
import { MockOperatorFactory } from '../../ethers/v0.7/MockOperatorFactory'

const getterSetterFactory = new GetterSetterFactory()
const mockOperatorFactory = new MockOperatorFactory()
const operatorProxyFactory = new OperatorProxyFactory()
const linkTokenFactory = new contract.LinkTokenFactory()

let roles: setup.Roles
const provider = setup.provider()

beforeAll(async () => {
  const users = await setup.users(provider)

  roles = users.roles
})

describe('OperatorProxy', () => {
  let link: contract.Instance<contract.LinkTokenFactory>
  let operatorProxy: contract.Instance<OperatorProxyFactory>
  let mockOperator: contract.Instance<MockOperatorFactory>
  const deployment = setup.snapshot(provider, async () => {
    link = await linkTokenFactory.connect(roles.defaultAccount).deploy()
    mockOperator = await mockOperatorFactory
      .connect(roles.defaultAccount)
      .deploy(link.address)
    const proxyAddress = await mockOperator.proxy()
    operatorProxy = await operatorProxyFactory
      .connect(roles.defaultAccount)
      .attach(proxyAddress)
  })

  beforeEach(async () => {
    await deployment()
  })

  it('has a limited public interface', () => {
    matchers.publicAbi(operatorProxy, [
      'forward',
      // Ownable methods:
      'acceptOwnership',
      'owner',
      'transferOwnership',
    ])
  })

  describe('#forward', () => {
    const bytes = utils.hexlify(utils.randomBytes(100))
    const payload = getterSetterFactory.interface.functions.setBytes.encode([
      bytes,
    ])
    let mock: contract.Instance<GetterSetterFactory>

    beforeEach(async () => {
      mock = await getterSetterFactory.connect(roles.defaultAccount).deploy()
    })

    describe('when called by an unauthorized node', () => {
      it('reverts', async () => {
        await matchers.evmRevert(async () => {
          await operatorProxy
            .connect(roles.stranger)
            .forward(mock.address, payload)
        })
      })
    })

    describe('when called by an authorized node', () => {
      beforeEach(async () => {
        await mockOperator.setIsAuthorized(true)
      })

      describe('when attempting to forward to the link token', () => {
        it('reverts', async () => {
          const { sighash } = linkTokenFactory.interface.functions.name // any Link Token function
          await matchers.evmRevert(async () => {
            await operatorProxy
              .connect(roles.oracleNode)
              .forward(link.address, sighash)
          })
        })
      })

      describe('when forwarding to any other address', () => {
        it('forwards the data', async () => {
          const tx = await operatorProxy
            .connect(roles.oracleNode)
            .forward(mock.address, payload)
          await tx.wait()
          assert.equal(await mock.getBytes(), bytes)
        })

        it('perceives the message is sent by the operatorProxy', async () => {
          const tx = await operatorProxy
            .connect(roles.oracleNode)
            .forward(mock.address, payload)
          const receipt = await tx.wait()
          const log: any = receipt.logs?.[0]
          const logData = mock.interface.events.SetBytes.decode(
            log.data,
            log.topics,
          )
          assert.equal(utils.getAddress(logData.from), operatorProxy.address)
        })
      })
    })
  })
})
