import {
  contract,
  helpers as h,
  matchers,
  oracle,
  setup,
} from '@chainlink/test-helpers'
import cbor from 'cbor'
import { assert } from 'chai'
import { ethers } from 'ethers'
import { BasicConsumer__factory } from '../../ethers/v0.6/factories/BasicConsumer__factory'
import { Oracle__factory } from '../../ethers/v0.6/factories/Oracle__factory'
import { PreCoordinator__factory } from '../../ethers/v0.6/factories/PreCoordinator__factory'

const provider = setup.provider()
const oracleFactory = new Oracle__factory()
const preCoordinatorFactory = new PreCoordinator__factory()
const requesterConsumerFactory = new BasicConsumer__factory()
const linkTokenFactory = new contract.LinkToken__factory()

let roles: setup.Roles
beforeAll(async () => {
  roles = await setup.users(provider).then((x) => x.roles)
})

describe('PreCoordinator', () => {
  // These parameters are used to validate the data was received
  // on the deployed oracle contract. The Job ID only represents
  // the type of data, but will not work on a public testnet.
  // For the latest JobIDs, visit our docs here:
  // https://docs.chain.link/docs/testnet-oracles
  const job1 =
    '0x4c7b7ffb66b344fbaa64995af81e355a00000000000000000000000000000000'
  const job2 =
    '0x4c7b7ffb66b344fbaa64995af81e355b00000000000000000000000000000000'
  const job3 =
    '0x4c7b7ffb66b344fbaa64995af81e355c00000000000000000000000000000000'
  const job4 =
    '0x4c7b7ffb66b344fbaa64995af81e355d00000000000000000000000000000000'
  const currency = 'USD'

  // Represents 1 LINK for testnet requests
  const payment = h.toWei('1')
  const totalPayment = h.toWei('4')

  let link: contract.Instance<contract.LinkToken__factory>
  let oc1: contract.Instance<Oracle__factory>
  let oc2: contract.Instance<Oracle__factory>
  let oc3: contract.Instance<Oracle__factory>
  let oc4: contract.Instance<Oracle__factory>
  let rc: contract.Instance<BasicConsumer__factory>
  let pc: contract.Instance<PreCoordinator__factory>

  const deployment = setup.snapshot(provider, async () => {
    link = await linkTokenFactory.connect(roles.defaultAccount).deploy()
    oc1 = await oracleFactory.connect(roles.defaultAccount).deploy(link.address)
    oc2 = await oracleFactory.connect(roles.defaultAccount).deploy(link.address)
    oc3 = await oracleFactory.connect(roles.defaultAccount).deploy(link.address)
    oc4 = await oracleFactory.connect(roles.defaultAccount).deploy(link.address)
    pc = await preCoordinatorFactory
      .connect(roles.defaultAccount)
      .deploy(link.address)

    await oc1
      .connect(roles.defaultAccount)
      .setFulfillmentPermission(roles.oracleNode1.address, true)
    await oc2
      .connect(roles.defaultAccount)
      .setFulfillmentPermission(roles.oracleNode2.address, true)
    await oc3
      .connect(roles.defaultAccount)
      .setFulfillmentPermission(roles.oracleNode3.address, true)
    await oc4
      .connect(roles.defaultAccount)
      .setFulfillmentPermission(roles.oracleNode4.address, true)
  })

  beforeEach(deployment)

  describe('#createServiceAgreement', () => {
    it('emits the NewServiceAgreement log', async () => {
      const tx = await pc
        .connect(roles.defaultAccount)
        .createServiceAgreement(
          3,
          [oc1.address, oc2.address, oc3.address, oc4.address],
          [job1, job2, job3, job4],
          [payment, payment, payment, payment],
        )
      const receipt = await tx.wait()

      expect(
        h.findEventIn(receipt, pc.interface.events.NewServiceAgreement),
      ).toBeDefined()
    })

    it('creates a service agreement', async () => {
      const tx = await pc
        .connect(roles.defaultAccount)
        .createServiceAgreement(
          3,
          [oc1.address, oc2.address, oc3.address, oc4.address],
          [job1, job2, job3, job4],
          [payment, payment, payment, payment],
        )
      const receipt = await tx.wait()
      const { saId } = h.eventArgs(
        h.findEventIn(receipt, pc.interface.events.NewServiceAgreement),
      )

      const sa = await pc.getServiceAgreement(saId)
      assert.isTrue(sa.totalPayment.eq(totalPayment))
      assert.equal(sa.minResponses.toNumber(), 3)
      assert.deepEqual(sa.oracles, [
        oc1.address,
        oc2.address,
        oc3.address,
        oc4.address,
      ])
      assert.deepEqual(sa.jobIds, [job1, job2, job3, job4])
      assert.deepEqual(
        sa.payments.map((p) => p.toHexString()),
        [payment, payment, payment, payment].map((p) => p.toHexString()),
      )
    })

    it('does not allow service agreements with 0 minResponses', () =>
      matchers.evmRevert(
        pc
          .connect(roles.defaultAccount)
          .createServiceAgreement(
            0,
            [oc1.address, oc2.address, oc3.address, oc4.address],
            [job1, job2, job3, job4],
            [payment, payment, payment, payment],
          ),
        'Min responses must be > 0',
      ))

    describe('when the array lengths are not equal', () => {
      it('reverts', () =>
        matchers.evmRevert(
          pc
            .connect(roles.defaultAccount)
            .createServiceAgreement(
              3,
              [oc1.address, oc2.address, oc3.address, oc4.address],
              [job1, job2, job3],
              [payment, payment, payment, payment],
            ),
          'Unmet length',
        ))
    })

    describe('when the min responses is greater than the oracles', () => {
      it('reverts', () =>
        matchers.evmRevert(
          pc
            .connect(roles.defaultAccount)
            .createServiceAgreement(
              5,
              [oc1.address, oc2.address, oc3.address, oc4.address],
              [job1, job2, job3, job4],
              [payment, payment, payment, payment],
            ),
          'Invalid min responses',
        ))
    })
  })

  describe('#onTokenTransfer', () => {
    describe('when called by an address other than the LINK token', () => {
      it('reverts', async () => {
        const notLink = await linkTokenFactory
          .connect(roles.defaultAccount)
          .deploy()

        const tx = await pc
          .connect(roles.defaultAccount)
          .createServiceAgreement(
            3,
            [oc1.address, oc2.address, oc3.address, oc4.address],
            [job1, job2, job3, job4],
            [payment, payment, payment, payment],
          )
        const receipt = await tx.wait()
        const saId = h.eventArgs(
          h.findEventIn(receipt, pc.interface.events.NewServiceAgreement),
        ).saId

        const badRc = await requesterConsumerFactory
          .connect(roles.consumer)
          .deploy(notLink.address, pc.address, saId)

        await notLink
          .connect(roles.defaultAccount)
          .transfer(badRc.address, totalPayment)

        await matchers.evmRevert(
          badRc
            .connect(roles.consumer)
            .requestEthereumPrice(currency, totalPayment, {}),
        )
      })
    })

    describe('when called by the LINK token', () => {
      let saId: string
      beforeEach(async () => {
        const tx = await pc
          .connect(roles.defaultAccount)
          .createServiceAgreement(
            3,
            [oc1.address, oc2.address, oc3.address, oc4.address],
            [job1, job2, job3, job4],
            [payment, payment, payment, payment],
          )
        const receipt = await tx.wait()
        saId = h.eventArgs(
          h.findEventIn(receipt, pc.interface.events.NewServiceAgreement),
        ).saId

        rc = await requesterConsumerFactory
          .connect(roles.consumer)
          .deploy(link.address, pc.address, saId)
        await link.transfer(rc.address, totalPayment)
      })

      it('creates Chainlink requests', async () => {
        const tx = await rc
          .connect(roles.consumer)
          .requestEthereumPrice(currency, totalPayment)
        const receipt = await tx.wait()

        const log1 = receipt.logs?.[7]
        assert.equal(oc1.address, log1?.address)
        const request1 = oracle.decodeRunRequest(log1)
        assert.equal(request1.requester, pc.address)
        const log2 = receipt.logs?.[11]
        assert.equal(oc2.address, log2?.address)
        const request2 = oracle.decodeRunRequest(log2)
        assert.equal(request2.requester, pc.address)
        const log3 = receipt.logs?.[15]
        assert.equal(oc3.address, log3?.address)
        const request3 = oracle.decodeRunRequest(log3)
        assert.equal(request3.requester, pc.address)
        const log4 = receipt.logs?.[19]
        assert.equal(oc4.address, log4?.address)
        const request4 = oracle.decodeRunRequest(log4)
        assert.equal(request4.requester, pc.address)
        const expected = {
          path: ['USD'],
          get:
            'https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY',
        }
        assert.deepEqual(expected, await cbor.decodeFirst(request1.data))
        assert.deepEqual(expected, await cbor.decodeFirst(request2.data))
        assert.deepEqual(expected, await cbor.decodeFirst(request3.data))
        assert.deepEqual(expected, await cbor.decodeFirst(request4.data))
      })

      describe('when insufficient payment is supplied', () => {
        it('reverts', () =>
          matchers.evmRevert(
            rc.connect(roles.consumer).requestEthereumPrice(currency, payment),
          ))
      })

      describe('when the same nonce is used twice', () => {
        const nonce = 1
        const fHash = '0xabcd1234'
        let args: string
        beforeEach(async () => {
          args = oracle.encodeOracleRequest(
            saId,
            rc.address,
            fHash,
            nonce,
            '0x0',
          )
          await link.transferAndCall(pc.address, totalPayment, args)
        })

        it('reverts', () =>
          matchers.evmRevert(
            link.transferAndCall(pc.address, totalPayment, args),
          ))
      })

      describe('when too much payment is supplied', () => {
        it('sends the extra back to the requester', async () => {
          await link.transfer(rc.address, payment)
          const extraPayment = h.toWei('5')
          const beforeBalance = await link.balanceOf(rc.address)
          expect(beforeBalance.eq(extraPayment)).toBeTruthy()

          await rc
            .connect(roles.consumer)
            .requestEthereumPrice(currency, extraPayment)
          const afterBalance = await link.balanceOf(rc.address)
          expect(afterBalance.eq(payment)).toBeTruthy()
        })
      })
    })
  })

  describe('#chainlinkCallback', () => {
    let saId: string
    let request1: oracle.RunRequest
    let request2: oracle.RunRequest
    let request3: oracle.RunRequest
    let request4: oracle.RunRequest
    const response1 = h.numToBytes32(100)
    const response2 = h.numToBytes32(101)
    const response3 = h.numToBytes32(102)
    const response4 = h.numToBytes32(103)

    beforeEach(async () => {
      const tx = await pc
        .connect(roles.defaultAccount)
        .createServiceAgreement(
          2,
          [oc1.address, oc2.address, oc3.address, oc4.address],
          [job1, job2, job3, job4],
          [payment, payment, payment, payment],
        )
      const receipt = await tx.wait()

      saId = h.eventArgs(
        h.findEventIn(receipt, pc.interface.events.NewServiceAgreement),
      ).saId
      rc = await requesterConsumerFactory
        .connect(roles.consumer)
        .deploy(link.address, pc.address, saId)
      await link.transfer(rc.address, totalPayment)
    })

    describe('when the requester and consumer are the same', () => {
      beforeEach(async () => {
        const reqTx = await rc
          .connect(roles.consumer)
          .requestEthereumPrice(currency, totalPayment)
        const receipt = await reqTx.wait()

        const log1 = receipt.logs?.[7]
        request1 = oracle.decodeRunRequest(log1)
        const log2 = receipt.logs?.[11]
        request2 = oracle.decodeRunRequest(log2)
        const log3 = receipt.logs?.[15]
        request3 = oracle.decodeRunRequest(log3)
        const log4 = receipt.logs?.[19]
        request4 = oracle.decodeRunRequest(log4)
      })

      describe('when called by a stranger', () => {
        it('reverts', () =>
          matchers.evmRevert(
            pc.chainlinkCallback(saId, response1),
            'Source must be the oracle of the request',
          ))
      })

      describe('when called by the oracle contract', () => {
        it('records the answer', async () => {
          const tx = await oc1
            .connect(roles.oracleNode1)
            .fulfillOracleRequest(
              request1.requestId,
              request1.payment,
              request1.callbackAddr,
              request1.callbackFunc,
              request1.expiration,
              response1,
            )
          const receipt = await tx.wait()

          expect(
            receipt.events?.[0].topics.find(
              (t) => t === pc.interface.events.ChainlinkFulfilled.topic,
            ),
          ).toBeDefined()
        })
      })

      describe('when the minimum number of responses have returned', () => {
        beforeEach(async () => {
          await oc1
            .connect(roles.oracleNode1)
            .fulfillOracleRequest(
              request1.requestId,
              request1.payment,
              request1.callbackAddr,
              request1.callbackFunc,
              request1.expiration,
              response1,
            )
          await oc2
            .connect(roles.oracleNode2)
            .fulfillOracleRequest(
              request2.requestId,
              request2.payment,
              request2.callbackAddr,
              request2.callbackFunc,
              request2.expiration,
              response2,
            )
          await oc3
            .connect(roles.oracleNode3)
            .fulfillOracleRequest(
              request3.requestId,
              request3.payment,
              request3.callbackAddr,
              request3.callbackFunc,
              request3.expiration,
              response3,
            )
        })

        it('returns the median to the requesting contract', async () => {
          const currentPrice = await rc.currentPrice()
          assert.equal(currentPrice, response1)
        })

        describe('when an oracle responds after aggregation', () => {
          it('does not update the requesting contract', async () => {
            await oc4
              .connect(roles.oracleNode4)
              .fulfillOracleRequest(
                request4.requestId,
                request4.payment,
                request4.callbackAddr,
                request4.callbackFunc,
                request4.expiration,
                response4,
              )
            const currentPrice = await rc.currentPrice()
            assert.equal(currentPrice, response1)
          })
        })
      })
    })

    describe('when consumer is different than requester', () => {
      let cc: contract.Instance<BasicConsumer__factory>
      let request1: oracle.RunRequest
      let request2: oracle.RunRequest
      let request3: oracle.RunRequest
      let request4: oracle.RunRequest
      let localRequestId: string

      beforeEach(async () => {
        cc = await requesterConsumerFactory
          .connect(roles.consumer)
          .deploy(link.address, pc.address, saId)
        const reqTx = await rc
          .connect(roles.consumer)
          .requestEthereumPriceByCallback(currency, totalPayment, cc.address)
        const receipt = await reqTx.wait()

        localRequestId = h.eventArgs(receipt.events?.[0]).id
        const log1 = receipt.logs?.[7]
        request1 = oracle.decodeRunRequest(log1)
        const log2 = receipt.logs?.[11]
        request2 = oracle.decodeRunRequest(log2)
        const log3 = receipt.logs?.[15]
        request3 = oracle.decodeRunRequest(log3)
        const log4 = receipt.logs?.[19]
        request4 = oracle.decodeRunRequest(log4)

        await cc
          .connect(roles.consumer)
          .addExternalRequest(pc.address, localRequestId)
      })

      describe('and the number of responses have been met', () => {
        beforeEach(async () => {
          await oc1
            .connect(roles.oracleNode1)
            .fulfillOracleRequest(
              request1.requestId,
              request1.payment,
              request1.callbackAddr,
              request1.callbackFunc,
              request1.expiration,
              response1,
            )
          await oc2
            .connect(roles.oracleNode2)
            .fulfillOracleRequest(
              request2.requestId,
              request2.payment,
              request2.callbackAddr,
              request2.callbackFunc,
              request2.expiration,
              response2,
            )
          await oc3
            .connect(roles.oracleNode3)
            .fulfillOracleRequest(
              request3.requestId,
              request3.payment,
              request3.callbackAddr,
              request3.callbackFunc,
              request3.expiration,
              response3,
            )
          await oc4
            .connect(roles.oracleNode4)
            .fulfillOracleRequest(
              request4.requestId,
              request4.payment,
              request4.callbackAddr,
              request4.callbackFunc,
              request4.expiration,
              response4,
            )
        })

        it('sends the answer to the consumer', async () => {
          const currentPrice = await cc.currentPrice()
          assert.equal(currentPrice, response1)
        })
      })
    })
  })

  describe('#withdrawLink', () => {
    beforeEach(async () => {
      await link.transfer(pc.address, payment)

      const actual = await link.balanceOf(pc.address)
      const expected = payment
      expect(actual.eq(expected)).toBeTruthy()
    })

    describe('when called by a stranger', () => {
      it('reverts', () =>
        matchers.evmRevert(pc.connect(roles.stranger).withdrawLink()))
    })

    describe('when called by the owner', () => {
      it('allows the owner to withdraw LINK', async () => {
        await pc.connect(roles.defaultAccount).withdrawLink()

        const actual = await link.balanceOf(pc.address)
        expect(actual.eq(ethers.constants.Zero)).toBeTruthy()
      })
    })
  })

  describe('#cancelOracleRequest', () => {
    let request: oracle.RunRequest

    beforeEach(async () => {
      const tx = await pc
        .connect(roles.defaultAccount)
        .createServiceAgreement(
          3,
          [oc1.address, oc2.address, oc3.address, oc4.address],
          [job1, job2, job3, job4],
          [payment, payment, payment, payment],
        )
      const receipt = await tx.wait()

      const saId = h.eventArgs(
        h.findEventIn(receipt, pc.interface.events.NewServiceAgreement),
      ).saId

      rc = await requesterConsumerFactory
        .connect(roles.consumer)
        .deploy(link.address, pc.address, saId)
      await link.transfer(rc.address, totalPayment)

      const reqTx = await rc
        .connect(roles.consumer)
        .requestEthereumPrice(currency, totalPayment)
      const reqReceipt = await reqTx.wait()

      const log1 = reqReceipt.logs?.[7]
      request = oracle.decodeRunRequest(log1)
    })

    describe('before the minimum required time', () => {
      it('does not allow requests to be cancelled', () =>
        matchers.evmRevert(
          rc
            .connect(roles.consumer)
            .cancelRequest(
              pc.address,
              request.requestId,
              request.payment,
              request.callbackFunc,
              request.expiration,
            ),
          'Request is not expired',
        ))
    })

    describe('after the minimum required time', () => {
      beforeEach(async () => {
        await h.increaseTime5Minutes(provider)
      })

      it('allows the requester to cancel', async () => {
        await rc
          .connect(roles.consumer)
          .cancelRequest(
            pc.address,
            request.requestId,
            request.payment,
            request.callbackFunc,
            request.expiration,
          )
        const balance = await link.balanceOf(rc.address)
        expect(balance.eq(payment)).toBeTruthy()
      })

      it('does not allow others to call', () =>
        matchers.evmRevert(
          pc
            .connect(roles.stranger)
            .cancelOracleRequest(
              request.requestId,
              request.payment,
              request.callbackFunc,
              request.expiration,
            ),
          'Only requester can cancel',
        ))
    })
  })
})
