import {
  contract,
  helpers as h,
  matchers,
  oracle,
  setup,
} from '@chainlink/test-helpers'
import { assert } from 'chai'
import { ethers } from 'ethers'
import { ConcreteChainlinked__factory } from '../../ethers/v0.4/factories/ConcreteChainlinked__factory'
import { EmptyOracle__factory } from '../../ethers/v0.4/factories/EmptyOracle__factory'
import { GetterSetter__factory } from '../../ethers/v0.4/factories/GetterSetter__factory'
import { Oracle__factory } from '../../ethers/v0.4/factories/Oracle__factory'

const concreteChainlinkedFactory = new ConcreteChainlinked__factory()
const emptyOracleFactory = new EmptyOracle__factory()
const getterSetterFactory = new GetterSetter__factory()
const oracleFactory = new Oracle__factory()
const linkTokenFactory = new contract.LinkToken__factory()

const provider = setup.provider()

let roles: setup.Roles

beforeAll(async () => {
  const users = await setup.users(provider)

  roles = users.roles
})

describe('ConcreteChainlinked', () => {
  const specId =
    '0x4c7b7ffb66b344fbaa64995af81e355a00000000000000000000000000000000'
  let cc: contract.Instance<ConcreteChainlinked__factory>
  let gs: contract.Instance<GetterSetter__factory>
  let oc: contract.Instance<Oracle__factory | EmptyOracle__factory>
  let newoc: contract.Instance<Oracle__factory>
  let link: contract.Instance<contract.LinkToken__factory>
  const deployment = setup.snapshot(provider, async () => {
    link = await linkTokenFactory.connect(roles.defaultAccount).deploy()
    oc = await oracleFactory.connect(roles.defaultAccount).deploy(link.address)
    newoc = await oracleFactory
      .connect(roles.defaultAccount)
      .deploy(link.address)
    gs = await getterSetterFactory.connect(roles.defaultAccount).deploy()
    cc = await concreteChainlinkedFactory
      .connect(roles.defaultAccount)
      .deploy(link.address, oc.address)
  })

  beforeEach(async () => {
    await deployment()
  })

  describe('#newRequest', () => {
    it('forwards the information to the oracle contract through the link token', async () => {
      const tx = await cc.publicNewRequest(
        specId,
        gs.address,
        ethers.utils.toUtf8Bytes('requestedBytes32(bytes32,bytes32)'),
      )
      const receipt = await tx.wait()

      assert.equal(1, receipt.logs?.length)
      const [jId, cbAddr, cbFId, cborData] = receipt.logs
        ? oracle.decodeCCRequest(receipt.logs[0])
        : []
      const params = h.decodeDietCBOR(cborData ?? '')

      assert.equal(specId, jId)
      assert.equal(gs.address, cbAddr)
      assert.equal('0xed53e511', cbFId)
      assert.deepEqual({}, params)
    })
  })

  describe('#chainlinkRequest(Request)', () => {
    it('emits an event from the contract showing the run ID', async () => {
      const tx = await cc.publicRequest(
        specId,
        cc.address,
        ethers.utils.toUtf8Bytes('fulfillRequest(bytes32,bytes32)'),
        0,
      )

      const { events, logs } = await tx.wait()

      assert.equal(4, events?.length)

      assert.equal(logs?.[0].address, cc.address)
      assert.equal(events?.[0].event, 'ChainlinkRequested')
    })
  })

  describe('#chainlinkRequestTo(Request)', () => {
    it('emits an event from the contract showing the run ID', async () => {
      const tx = await cc.publicRequestRunTo(
        newoc.address,
        specId,
        cc.address,
        ethers.utils.toUtf8Bytes('fulfillRequest(bytes32,bytes32)'),
        0,
      )
      const { events } = await tx.wait()

      assert.equal(4, events?.length)
      assert.equal(events?.[0].event, 'ChainlinkRequested')
    })

    it('emits an event on the target oracle contract', async () => {
      const tx = await cc.publicRequestRunTo(
        newoc.address,
        specId,
        cc.address,
        ethers.utils.toUtf8Bytes('fulfillRequest(bytes32,bytes32)'),
        0,
      )
      const { logs } = await tx.wait()
      const event = logs && newoc.interface.parseLog(logs[3])

      assert.equal(4, logs?.length)
      assert.equal(event?.name, 'OracleRequest')
    })

    it('does not modify the stored oracle address', async () => {
      await cc.publicRequestRunTo(
        newoc.address,
        specId,
        cc.address,
        ethers.utils.toUtf8Bytes('fulfillRequest(bytes32,bytes32)'),
        0,
      )

      const actualOracleAddress = await cc.publicOracleAddress()
      assert.equal(oc.address, actualOracleAddress)
    })
  })

  describe('#cancelChainlinkRequest', () => {
    let requestId: string
    // a concrete chainlink attached to an empty oracle
    let ecc: contract.Instance<ConcreteChainlinked__factory>

    beforeEach(async () => {
      const emptyOracle = await emptyOracleFactory
        .connect(roles.defaultAccount)
        .deploy()
      ecc = await concreteChainlinkedFactory
        .connect(roles.defaultAccount)
        .deploy(link.address, emptyOracle.address)

      const tx = await ecc.publicRequest(
        specId,
        ecc.address,
        ethers.utils.toUtf8Bytes('fulfillRequest(bytes32,bytes32)'),
        0,
      )
      const { events } = await tx.wait()
      requestId = (events?.[0]?.args as any).id
    })

    it('emits an event from the contract showing the run was cancelled', async () => {
      const tx = await ecc.publicCancelRequest(
        requestId,
        0,
        ethers.utils.hexZeroPad('0x', 4),
        0,
      )
      const { events } = await tx.wait()

      assert.equal(1, events?.length)
      assert.equal(events?.[0].event, 'ChainlinkCancelled')
      assert.equal(requestId, (events?.[0].args as any).id)
    })

    it('throws if given a bogus event ID', async () => {
      await matchers.evmRevert(async () => {
        await ecc.publicCancelRequest(
          ethers.utils.formatBytes32String('bogusId'),
          0,
          ethers.utils.hexZeroPad('0x', 4),
          0,
        )
      })
    })
  })

  describe('#recordChainlinkFulfillment(modifier)', () => {
    let request: oracle.RunRequest

    beforeEach(async () => {
      const tx = await cc.publicRequest(
        specId,
        cc.address,
        ethers.utils.toUtf8Bytes('fulfillRequest(bytes32,bytes32)'),
        0,
      )
      const { logs } = await tx.wait()

      request = oracle.decodeRunRequest(logs?.[3])
    })

    it('emits an event marking the request fulfilled', async () => {
      const tx = await oc.fulfillOracleRequest(
        ...oracle.convertFufillParams(
          request,
          ethers.utils.formatBytes32String('hi mom!'),
        ),
      )
      const { logs } = await tx.wait()

      const event = logs && cc.interface.parseLog(logs[0])

      assert.equal(1, logs?.length)
      assert.equal(event?.name, 'ChainlinkFulfilled')
      assert.equal(request.requestId, event?.values.id)
    })
  })

  describe('#fulfillChainlinkRequest(function)', () => {
    let request: oracle.RunRequest

    beforeEach(async () => {
      const tx = await cc.publicRequest(
        specId,
        cc.address,
        ethers.utils.toUtf8Bytes(
          'publicFulfillChainlinkRequest(bytes32,bytes32)',
        ),
        0,
      )
      const { logs } = await tx.wait()

      request = oracle.decodeRunRequest(logs?.[3])
    })

    it('emits an event marking the request fulfilled', async () => {
      const tx = await oc.fulfillOracleRequest(
        ...oracle.convertFufillParams(
          request,
          ethers.utils.formatBytes32String('hi mom!'),
        ),
      )
      const { logs } = await tx.wait()
      const event = logs && cc.interface.parseLog(logs[0])

      assert.equal(1, logs?.length)
      assert.equal(event?.name, 'ChainlinkFulfilled')
      assert.equal(request.requestId, event?.values?.id)
    })
  })

  describe('#chainlinkToken', () => {
    it('returns the Link Token address', async () => {
      const addr = await cc.publicChainlinkToken()
      assert.equal(addr, link.address)
    })
  })

  describe('#addExternalRequest', () => {
    let mock: contract.Instance<ConcreteChainlinked__factory>
    let request: oracle.RunRequest

    beforeEach(async () => {
      mock = await concreteChainlinkedFactory
        .connect(roles.defaultAccount)
        .deploy(link.address, oc.address)

      const tx = await cc.publicRequest(
        specId,
        mock.address,
        ethers.utils.toUtf8Bytes('fulfillRequest(bytes32,bytes32)'),
        0,
      )
      const receipt = await tx.wait()

      request = oracle.decodeRunRequest(receipt.logs?.[3])
      await mock.publicAddExternalRequest(oc.address, request.requestId)
    })

    it('allows the external request to be fulfilled', async () => {
      await oc.fulfillOracleRequest(
        ...oracle.convertFufillParams(
          request,
          ethers.utils.formatBytes32String('hi mom!'),
        ),
      )
    })

    it('does not allow the same requestId to be used', async () => {
      await matchers.evmRevert(async () => {
        await cc.publicAddExternalRequest(newoc.address, request.requestId)
      })
    })
  })
})
