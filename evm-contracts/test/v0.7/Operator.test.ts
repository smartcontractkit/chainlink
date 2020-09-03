import {
  contract,
  helpers as h,
  matchers,
  oracle,
  setup,
} from '@chainlink/test-helpers'
import { assert } from 'chai'
import { ethers, utils } from 'ethers'
import { ContractReceipt } from 'ethers/contract'
import { BasicConsumerFactory } from '../../ethers/v0.4/BasicConsumerFactory'
import { GetterSetterFactory } from '../../ethers/v0.4/GetterSetterFactory'
import { MaliciousConsumerFactory } from '../../ethers/v0.4/MaliciousConsumerFactory'
import { MaliciousRequesterFactory } from '../../ethers/v0.4/MaliciousRequesterFactory'
import { OperatorFactory } from '../../ethers/v0.7/OperatorFactory'

const basicConsumerFactory = new BasicConsumerFactory()
const getterSetterFactory = new GetterSetterFactory()
const maliciousRequesterFactory = new MaliciousRequesterFactory()
const maliciousConsumerFactory = new MaliciousConsumerFactory()
const operatorFactory = new OperatorFactory()
const linkTokenFactory = new contract.LinkTokenFactory()

let roles: setup.Roles
const provider = setup.provider()

beforeAll(async () => {
  const users = await setup.users(provider)

  roles = users.roles
})

describe('Operator', () => {
  const fHash = getterSetterFactory.interface.functions.requestedBytes32.sighash
  const specId =
    '0x4c7b7ffb66b344fbaa64995af81e355a00000000000000000000000000000000'
  const to = '0x80e29acb842498fe6591f020bd82766dce619d43'
  let link: contract.Instance<contract.LinkTokenFactory>
  let operator: contract.Instance<OperatorFactory>
  const deployment = setup.snapshot(provider, async () => {
    link = await linkTokenFactory.connect(roles.defaultAccount).deploy()
    operator = await operatorFactory
      .connect(roles.defaultAccount)
      .deploy(link.address)
    await operator.setFulfillmentPermission(roles.oracleNode.address, true)
  })

  beforeEach(async () => {
    await deployment()
  })

  it('has a limited public interface', () => {
    matchers.publicAbi(operatorFactory, [
      'EXPIRY_TIME',
      'cancelOracleRequest',
      'foreward',
      'fulfillOracleRequest',
      'getAuthorizationStatus',
      'getChainlinkToken',
      'onTokenTransfer',
      'oracleRequest',
      'setFulfillmentPermission',
      'withdraw',
      'withdrawable',
      // Ownable methods:
      'acceptOwnership',
      'owner',
      'transferOwnership',
    ])
  })

  describe('#setFulfillmentPermission', () => {
    describe('when called by the owner', () => {
      beforeEach(async () => {
        await operator
          .connect(roles.defaultAccount)
          .setFulfillmentPermission(roles.stranger.address, true)
      })

      it('adds an authorized node', async () => {
        const authorized = await operator.getAuthorizationStatus(
          roles.stranger.address,
        )
        assert.equal(true, authorized)
      })

      it('removes an authorized node', async () => {
        await operator
          .connect(roles.defaultAccount)
          .setFulfillmentPermission(roles.stranger.address, false)
        const authorized = await operator.getAuthorizationStatus(
          roles.stranger.address,
        )
        assert.equal(false, authorized)
      })
    })

    describe('when called by a non-owner', () => {
      it('cannot add an authorized node', async () => {
        await matchers.evmRevert(async () => {
          await operator
            .connect(roles.stranger)
            .setFulfillmentPermission(roles.stranger.address, true)
        })
      })
    })
  })

  describe('#onTokenTransfer', () => {
    describe('when called from any address but the LINK token', () => {
      it('triggers the intended method', async () => {
        const callData = oracle.encodeOracleRequest(specId, to, fHash, 0, '0x0')

        await matchers.evmRevert(async () => {
          await operator.onTokenTransfer(
            roles.defaultAccount.address,
            0,
            callData,
          )
        })
      })
    })

    describe('when called from the LINK token', () => {
      it('triggers the intended method', async () => {
        const callData = oracle.encodeOracleRequest(specId, to, fHash, 0, '0x0')

        const tx = await link.transferAndCall(operator.address, 0, callData, {
          value: 0,
        })
        const receipt = await tx.wait()

        assert.equal(3, receipt.logs?.length)
      })

      describe('with no data', () => {
        it('reverts', async () => {
          await matchers.evmRevert(async () => {
            await link.transferAndCall(operator.address, 0, '0x', {
              value: 0,
            })
          })
        })
      })
    })

    describe('malicious requester', () => {
      let mock: contract.Instance<MaliciousRequesterFactory>
      let requester: contract.Instance<BasicConsumerFactory>
      const paymentAmount = h.toWei('1')

      beforeEach(async () => {
        mock = await maliciousRequesterFactory
          .connect(roles.defaultAccount)
          .deploy(link.address, operator.address)
        await link.transfer(mock.address, paymentAmount)
      })

      it('cannot withdraw from oracle', async () => {
        const operatorOriginalBalance = await link.balanceOf(operator.address)
        const mockOriginalBalance = await link.balanceOf(mock.address)

        await matchers.evmRevert(async () => {
          await mock.maliciousWithdraw()
        })

        const operatorNewBalance = await link.balanceOf(operator.address)
        const mockNewBalance = await link.balanceOf(mock.address)

        matchers.bigNum(operatorOriginalBalance, operatorNewBalance)
        matchers.bigNum(mockNewBalance, mockOriginalBalance)
      })

      describe('if the requester tries to create a requestId for another contract', () => {
        it('the requesters ID will not match with the oracle contract', async () => {
          const tx = await mock.maliciousTargetConsumer(to)
          const receipt = await tx.wait()

          const mockRequestId = receipt.logs?.[0].data
          const requestId = (receipt.events?.[0].args as any).requestId
          assert.notEqual(mockRequestId, requestId)
        })

        it('the target requester can still create valid requests', async () => {
          requester = await basicConsumerFactory
            .connect(roles.defaultAccount)
            .deploy(link.address, operator.address, specId)
          await link.transfer(requester.address, paymentAmount)
          await mock.maliciousTargetConsumer(requester.address)
          await requester.requestEthereumPrice('USD')
        })
      })
    })

    it('does not allow recursive calls of onTokenTransfer', async () => {
      const requestPayload = oracle.encodeOracleRequest(
        specId,
        to,
        fHash,
        0,
        '0x0',
      )

      const ottSelector =
        operatorFactory.interface.functions.onTokenTransfer.sighash
      const header =
        '000000000000000000000000c5fdf4076b8f3a5357c5e395ab970b5b54098fef' + // to
        '0000000000000000000000000000000000000000000000000000000000000539' + // amount
        '0000000000000000000000000000000000000000000000000000000000000060' + // offset
        '0000000000000000000000000000000000000000000000000000000000000136' //   length

      const maliciousPayload = ottSelector + header + requestPayload.slice(2)

      await matchers.evmRevert(async () => {
        await link.transferAndCall(operator.address, 0, maliciousPayload, {
          value: 0,
        })
      })
    })
  })

  describe('#oracleRequest', () => {
    describe('when called through the LINK token', () => {
      const paid = 100
      let log: ethers.providers.Log | undefined
      let receipt: ethers.providers.TransactionReceipt

      beforeEach(async () => {
        const args = oracle.encodeOracleRequest(specId, to, fHash, 1, '0x0')
        const tx = await link.transferAndCall(operator.address, paid, args)
        receipt = await tx.wait()
        assert.equal(3, receipt?.logs?.length)

        log = receipt.logs && receipt.logs[2]
      })

      it('logs an event', async () => {
        assert.equal(operator.address, log?.address)

        assert.equal(log?.topics?.[1], specId)

        const req = oracle.decodeRunRequest(receipt?.logs?.[2])
        assert.equal(roles.defaultAccount.address, req.requester)
        matchers.bigNum(paid, req.payment)
      })

      it('uses the expected event signature', async () => {
        // If updating this test, be sure to update models.RunLogTopic.
        const eventSignature =
          '0xd8d7ecc4800d25fa53ce0372f13a416d98907a7ef3d8d3bdd79cf4fe75529c65'
        assert.equal(eventSignature, log?.topics?.[0])
      })

      it('does not allow the same requestId to be used twice', async () => {
        const args2 = oracle.encodeOracleRequest(specId, to, fHash, 1, '0x0')
        await matchers.evmRevert(async () => {
          await link.transferAndCall(operator.address, paid, args2)
        })
      })

      describe('when called with a payload less than 2 EVM words + function selector', () => {
        const funcSelector =
          operatorFactory.interface.functions.oracleRequest.sighash
        const maliciousData =
          funcSelector +
          '0000000000000000000000000000000000000000000000000000000000000000000'

        it('throws an error', async () => {
          await matchers.evmRevert(async () => {
            await link.transferAndCall(operator.address, paid, maliciousData)
          })
        })
      })

      describe('when called with a payload between 3 and 9 EVM words', () => {
        const funcSelector =
          operatorFactory.interface.functions.oracleRequest.sighash
        const maliciousData =
          funcSelector +
          '000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001'

        it('throws an error', async () => {
          await matchers.evmRevert(async () => {
            await link.transferAndCall(operator.address, paid, maliciousData)
          })
        })
      })
    })

    describe('when not called through the LINK token', () => {
      it('reverts', async () => {
        await matchers.evmRevert(async () => {
          await operator
            .connect(roles.oracleNode)
            .oracleRequest(
              '0x0000000000000000000000000000000000000000',
              0,
              specId,
              to,
              fHash,
              1,
              1,
              '0x',
            )
        })
      })
    })
  })

  describe('#fulfillOracleRequest', () => {
    const response = 'Hi Mom!'
    let maliciousRequester: contract.Instance<MaliciousRequesterFactory>
    let basicConsumer: contract.Instance<BasicConsumerFactory>
    let maliciousConsumer: contract.Instance<MaliciousConsumerFactory>
    let request: ReturnType<typeof oracle.decodeRunRequest>

    describe('cooperative consumer', () => {
      beforeEach(async () => {
        basicConsumer = await basicConsumerFactory
          .connect(roles.defaultAccount)
          .deploy(link.address, operator.address, specId)
        const paymentAmount = h.toWei('1')
        await link.transfer(basicConsumer.address, paymentAmount)
        const currency = 'USD'
        const tx = await basicConsumer.requestEthereumPrice(currency)
        const receipt = await tx.wait()
        request = oracle.decodeRunRequest(receipt.logs?.[3])
      })

      describe('when called by an unauthorized node', () => {
        beforeEach(async () => {
          assert.equal(
            false,
            await operator.getAuthorizationStatus(roles.stranger.address),
          )
        })

        it('raises an error', async () => {
          await matchers.evmRevert(async () => {
            await operator
              .connect(roles.stranger)
              .fulfillOracleRequest(
                ...oracle.convertFufillParams(request, response),
              )
          })
        })
      })

      describe('when called by an authorized node', () => {
        it('raises an error if the request ID does not exist', async () => {
          request.requestId = utils.formatBytes32String('DOESNOTEXIST')
          await matchers.evmRevert(async () => {
            await operator
              .connect(roles.oracleNode)
              .fulfillOracleRequest(
                ...oracle.convertFufillParams(request, response),
              )
          })
        })

        it('sets the value on the requested contract', async () => {
          await operator
            .connect(roles.oracleNode)
            .fulfillOracleRequest(
              ...oracle.convertFufillParams(request, response),
            )

          const currentValue = await basicConsumer.currentPrice()
          assert.equal(response, ethers.utils.parseBytes32String(currentValue))
        })

        it('does not allow a request to be fulfilled twice', async () => {
          const response2 = response + ' && Hello World!!'

          await operator
            .connect(roles.oracleNode)
            .fulfillOracleRequest(
              ...oracle.convertFufillParams(request, response),
            )

          await matchers.evmRevert(async () => {
            await operator
              .connect(roles.oracleNode)
              .fulfillOracleRequest(
                ...oracle.convertFufillParams(request, response2),
              )
          })

          const currentValue = await basicConsumer.currentPrice()
          assert.equal(response, ethers.utils.parseBytes32String(currentValue))
        })
      })

      describe('when the oracle does not provide enough gas', () => {
        // if updating this defaultGasLimit, be sure it matches with the
        // defaultGasLimit specified in store/tx_manager.go
        const defaultGasLimit = 500000

        beforeEach(async () => {
          matchers.bigNum(0, await operator.withdrawable())
        })

        it('does not allow the oracle to withdraw the payment', async () => {
          await matchers.evmRevert(async () => {
            await operator.connect(roles.oracleNode).fulfillOracleRequest(
              ...oracle.convertFufillParams(request, response, {
                gasLimit: 70000,
              }),
            )
          })

          matchers.bigNum(0, await operator.withdrawable())
        })

        it(`${defaultGasLimit} is enough to pass the gas requirement`, async () => {
          await operator.connect(roles.oracleNode).fulfillOracleRequest(
            ...oracle.convertFufillParams(request, response, {
              gasLimit: defaultGasLimit,
            }),
          )

          matchers.bigNum(request.payment, await operator.withdrawable())
        })
      })
    })

    describe('with a malicious requester', () => {
      beforeEach(async () => {
        const paymentAmount = h.toWei('1')
        maliciousRequester = await maliciousRequesterFactory
          .connect(roles.defaultAccount)
          .deploy(link.address, operator.address)
        await link.transfer(maliciousRequester.address, paymentAmount)
      })

      it('cannot cancel before the expiration', async () => {
        await matchers.evmRevert(async () => {
          await maliciousRequester.maliciousRequestCancel(
            specId,
            ethers.utils.toUtf8Bytes('doesNothing(bytes32,bytes32)'),
          )
        })
      })

      it('cannot call functions on the LINK token through callbacks', async () => {
        await matchers.evmRevert(async () => {
          await maliciousRequester.request(
            specId,
            link.address,
            ethers.utils.toUtf8Bytes('transfer(address,uint256)'),
          )
        })
      })

      describe('requester lies about amount of LINK sent', () => {
        it('the oracle uses the amount of LINK actually paid', async () => {
          const tx = await maliciousRequester.maliciousPrice(specId)
          const receipt = await tx.wait()
          const req = oracle.decodeRunRequest(receipt.logs?.[3])

          assert(h.toWei('1').eq(req.payment))
        })
      })
    })

    describe('with a malicious consumer', () => {
      const paymentAmount = h.toWei('1')

      beforeEach(async () => {
        maliciousConsumer = await maliciousConsumerFactory
          .connect(roles.defaultAccount)
          .deploy(link.address, operator.address)
        await link.transfer(maliciousConsumer.address, paymentAmount)
      })

      describe('fails during fulfillment', () => {
        beforeEach(async () => {
          const tx = await maliciousConsumer.requestData(
            specId,
            ethers.utils.toUtf8Bytes('assertFail(bytes32,bytes32)'),
          )
          const receipt = await tx.wait()
          request = oracle.decodeRunRequest(receipt.logs?.[3])
        })

        it('allows the oracle node to receive their payment', async () => {
          await operator
            .connect(roles.oracleNode)
            .fulfillOracleRequest(
              ...oracle.convertFufillParams(request, response),
            )

          const balance = await link.balanceOf(roles.oracleNode.address)
          matchers.bigNum(balance, 0)

          await operator
            .connect(roles.defaultAccount)
            .withdraw(roles.oracleNode.address, paymentAmount)

          const newBalance = await link.balanceOf(roles.oracleNode.address)
          matchers.bigNum(paymentAmount, newBalance)
        })

        it("can't fulfill the data again", async () => {
          const response2 = 'hack the planet 102'

          await operator
            .connect(roles.oracleNode)
            .fulfillOracleRequest(
              ...oracle.convertFufillParams(request, response),
            )

          await matchers.evmRevert(async () => {
            await operator
              .connect(roles.oracleNode)
              .fulfillOracleRequest(
                ...oracle.convertFufillParams(request, response2),
              )
          })
        })
      })

      describe('calls selfdestruct', () => {
        beforeEach(async () => {
          const tx = await maliciousConsumer.requestData(
            specId,
            ethers.utils.toUtf8Bytes('doesNothing(bytes32,bytes32)'),
          )
          const receipt = await tx.wait()
          request = oracle.decodeRunRequest(receipt.logs?.[3])
          await maliciousConsumer.remove()
        })

        it('allows the oracle node to receive their payment', async () => {
          await operator
            .connect(roles.oracleNode)
            .fulfillOracleRequest(
              ...oracle.convertFufillParams(request, response),
            )

          const balance = await link.balanceOf(roles.oracleNode.address)
          matchers.bigNum(balance, 0)

          await operator
            .connect(roles.defaultAccount)
            .withdraw(roles.oracleNode.address, paymentAmount)
          const newBalance = await link.balanceOf(roles.oracleNode.address)
          matchers.bigNum(paymentAmount, newBalance)
        })
      })

      describe('request is canceled during fulfillment', () => {
        beforeEach(async () => {
          const tx = await maliciousConsumer.requestData(
            specId,
            ethers.utils.toUtf8Bytes('cancelRequestOnFulfill(bytes32,bytes32)'),
          )
          const receipt = await tx.wait()
          request = oracle.decodeRunRequest(receipt.logs?.[3])

          matchers.bigNum(0, await link.balanceOf(maliciousConsumer.address))
        })

        it('allows the oracle node to receive their payment', async () => {
          await operator
            .connect(roles.oracleNode)
            .fulfillOracleRequest(
              ...oracle.convertFufillParams(request, response),
            )

          const mockBalance = await link.balanceOf(maliciousConsumer.address)
          matchers.bigNum(mockBalance, 0)

          const balance = await link.balanceOf(roles.oracleNode.address)
          matchers.bigNum(balance, 0)

          await operator
            .connect(roles.defaultAccount)
            .withdraw(roles.oracleNode.address, paymentAmount)
          const newBalance = await link.balanceOf(roles.oracleNode.address)
          matchers.bigNum(paymentAmount, newBalance)
        })

        it("can't fulfill the data again", async () => {
          const response2 = 'hack the planet 102'

          await operator
            .connect(roles.oracleNode)
            .fulfillOracleRequest(
              ...oracle.convertFufillParams(request, response),
            )

          await matchers.evmRevert(async () => {
            await operator
              .connect(roles.oracleNode)
              .fulfillOracleRequest(
                ...oracle.convertFufillParams(request, response2),
              )
          })
        })
      })

      describe('tries to steal funds from node', () => {
        it('is not successful with call', async () => {
          const tx = await maliciousConsumer.requestData(
            specId,
            ethers.utils.toUtf8Bytes('stealEthCall(bytes32,bytes32)'),
          )
          const receipt = await tx.wait()
          request = oracle.decodeRunRequest(receipt.logs?.[3])

          await operator
            .connect(roles.oracleNode)
            .fulfillOracleRequest(
              ...oracle.convertFufillParams(request, response),
            )

          matchers.bigNum(
            0,
            await provider.getBalance(maliciousConsumer.address),
          )
        })

        it('is not successful with send', async () => {
          const tx = await maliciousConsumer.requestData(
            specId,
            ethers.utils.toUtf8Bytes('stealEthSend(bytes32,bytes32)'),
          )
          const receipt = await tx.wait()
          request = oracle.decodeRunRequest(receipt.logs?.[3])

          await operator
            .connect(roles.oracleNode)
            .fulfillOracleRequest(
              ...oracle.convertFufillParams(request, response),
            )
          matchers.bigNum(
            0,
            await provider.getBalance(maliciousConsumer.address),
          )
        })

        it('is not successful with transfer', async () => {
          const tx = await maliciousConsumer.requestData(
            specId,
            ethers.utils.toUtf8Bytes('stealEthTransfer(bytes32,bytes32)'),
          )
          const receipt = await tx.wait()
          request = oracle.decodeRunRequest(receipt.logs?.[3])

          await operator
            .connect(roles.oracleNode)
            .fulfillOracleRequest(
              ...oracle.convertFufillParams(request, response),
            )
          matchers.bigNum(
            0,
            await provider.getBalance(maliciousConsumer.address),
          )
        })
      })
    })
  })

  describe('#withdraw', () => {
    describe('without reserving funds via oracleRequest', () => {
      it('does nothing', async () => {
        let balance = await link.balanceOf(roles.oracleNode.address)
        assert.equal(0, balance.toNumber())
        await matchers.evmRevert(async () => {
          await operator
            .connect(roles.defaultAccount)
            .withdraw(roles.oracleNode.address, h.toWei('1'))
        })
        balance = await link.balanceOf(roles.oracleNode.address)
        assert.equal(0, balance.toNumber())
      })
    })

    describe('reserving funds via oracleRequest', () => {
      const payment = 15
      let request: ReturnType<typeof oracle.decodeRunRequest>

      beforeEach(async () => {
        const mock = await getterSetterFactory
          .connect(roles.defaultAccount)
          .deploy()
        const args = oracle.encodeOracleRequest(
          specId,
          mock.address,
          fHash,
          0,
          '0x0',
        )
        const tx = await link.transferAndCall(operator.address, payment, args)
        const receipt = await tx.wait()
        assert.equal(3, receipt.logs?.length)
        request = oracle.decodeRunRequest(receipt.logs?.[2])
      })

      describe('but not freeing funds w fulfillOracleRequest', () => {
        it('does not transfer funds', async () => {
          await matchers.evmRevert(async () => {
            await operator
              .connect(roles.defaultAccount)
              .withdraw(roles.oracleNode.address, payment)
          })
          const balance = await link.balanceOf(roles.oracleNode.address)
          assert.equal(0, balance.toNumber())
        })
      })

      describe('and freeing funds', () => {
        beforeEach(async () => {
          await operator
            .connect(roles.oracleNode)
            .fulfillOracleRequest(
              ...oracle.convertFufillParams(request, 'Hello World!'),
            )
        })

        it('does not allow input greater than the balance', async () => {
          const originalOracleBalance = await link.balanceOf(operator.address)
          const originalStrangerBalance = await link.balanceOf(
            roles.stranger.address,
          )
          const withdrawalAmount = payment + 1

          assert.isAbove(withdrawalAmount, originalOracleBalance.toNumber())
          await matchers.evmRevert(async () => {
            await operator
              .connect(roles.defaultAccount)
              .withdraw(roles.stranger.address, withdrawalAmount)
          })

          const newOracleBalance = await link.balanceOf(operator.address)
          const newStrangerBalance = await link.balanceOf(
            roles.stranger.address,
          )

          assert.equal(
            originalOracleBalance.toNumber(),
            newOracleBalance.toNumber(),
          )
          assert.equal(
            originalStrangerBalance.toNumber(),
            newStrangerBalance.toNumber(),
          )
        })

        it('allows transfer of partial balance by owner to specified address', async () => {
          const partialAmount = 6
          const difference = payment - partialAmount
          await operator
            .connect(roles.defaultAccount)
            .withdraw(roles.stranger.address, partialAmount)
          const strangerBalance = await link.balanceOf(roles.stranger.address)
          const oracleBalance = await link.balanceOf(operator.address)
          assert.equal(partialAmount, strangerBalance.toNumber())
          assert.equal(difference, oracleBalance.toNumber())
        })

        it('allows transfer of entire balance by owner to specified address', async () => {
          await operator
            .connect(roles.defaultAccount)
            .withdraw(roles.stranger.address, payment)
          const balance = await link.balanceOf(roles.stranger.address)
          assert.equal(payment, balance.toNumber())
        })

        it('does not allow a transfer of funds by non-owner', async () => {
          await matchers.evmRevert(async () => {
            await operator
              .connect(roles.stranger)
              .withdraw(roles.stranger.address, payment)
          })
          const balance = await link.balanceOf(roles.stranger.address)
          assert.isTrue(ethers.constants.Zero.eq(balance))
        })
      })
    })
  })

  describe('#withdrawable', () => {
    let request: ReturnType<typeof oracle.decodeRunRequest>

    beforeEach(async () => {
      const amount = h.toWei('1')
      const mock = await getterSetterFactory
        .connect(roles.defaultAccount)
        .deploy()
      const args = oracle.encodeOracleRequest(
        specId,
        mock.address,
        fHash,
        0,
        '0x0',
      )
      const tx = await link.transferAndCall(operator.address, amount, args)
      const receipt = await tx.wait()
      assert.equal(3, receipt.logs?.length)
      request = oracle.decodeRunRequest(receipt.logs?.[2])
      await operator
        .connect(roles.oracleNode)
        .fulfillOracleRequest(
          ...oracle.convertFufillParams(request, 'Hello World!'),
        )
    })

    it('returns the correct value', async () => {
      const withdrawAmount = await operator.withdrawable()
      matchers.bigNum(withdrawAmount, request.payment)
    })
  })

  describe('#cancelOracleRequest', () => {
    describe('with no pending requests', () => {
      it('fails', async () => {
        const fakeRequest: oracle.RunRequest = {
          requestId: ethers.utils.formatBytes32String('1337'),
          payment: '0',
          callbackFunc:
            getterSetterFactory.interface.functions.requestedBytes32.sighash,
          expiration: '999999999999',

          callbackAddr: '',
          data: Buffer.from(''),
          dataVersion: 0,
          specId: '',
          requester: '',
          topic: '',
        }
        await h.increaseTime5Minutes(provider)

        await matchers.evmRevert(async () => {
          await operator
            .connect(roles.stranger)
            .cancelOracleRequest(...oracle.convertCancelParams(fakeRequest))
        })
      })
    })

    describe('with a pending request', () => {
      const startingBalance = 100
      let request: ReturnType<typeof oracle.decodeRunRequest>
      let receipt: ethers.providers.TransactionReceipt

      beforeEach(async () => {
        const requestAmount = 20

        await link.transfer(roles.consumer.address, startingBalance)

        const args = oracle.encodeOracleRequest(
          specId,
          roles.consumer.address,
          fHash,
          1,
          '0x0',
        )
        const tx = await link
          .connect(roles.consumer)
          .transferAndCall(operator.address, requestAmount, args)
        receipt = await tx.wait()

        assert.equal(3, receipt.logs?.length)
        request = oracle.decodeRunRequest(receipt.logs?.[2])
      })

      it('has correct initial balances', async () => {
        const oracleBalance = await link.balanceOf(operator.address)
        matchers.bigNum(request.payment, oracleBalance)

        const consumerAmount = await link.balanceOf(roles.consumer.address)
        assert.equal(
          startingBalance - Number(request.payment),
          consumerAmount.toNumber(),
        )
      })

      describe('from a stranger', () => {
        it('fails', async () => {
          await matchers.evmRevert(async () => {
            await operator
              .connect(roles.consumer)
              .cancelOracleRequest(...oracle.convertCancelParams(request))
          })
        })
      })

      describe('from the requester', () => {
        it('refunds the correct amount', async () => {
          await h.increaseTime5Minutes(provider)
          await operator
            .connect(roles.consumer)
            .cancelOracleRequest(...oracle.convertCancelParams(request))
          const balance = await link.balanceOf(roles.consumer.address)

          assert.equal(startingBalance, balance.toNumber()) // 100
        })

        it('triggers a cancellation event', async () => {
          await h.increaseTime5Minutes(provider)
          const tx = await operator
            .connect(roles.consumer)
            .cancelOracleRequest(...oracle.convertCancelParams(request))
          const receipt = await tx.wait()

          assert.equal(receipt.logs?.length, 2)
          assert.equal(request.requestId, receipt.logs?.[0].topics[1])
        })

        it('fails when called twice', async () => {
          await h.increaseTime5Minutes(provider)
          await operator
            .connect(roles.consumer)
            .cancelOracleRequest(...oracle.convertCancelParams(request))

          await matchers.evmRevert(
            operator
              .connect(roles.consumer)
              .cancelOracleRequest(...oracle.convertCancelParams(request)),
          )
        })
      })
    })
  })

  describe('#foreward', () => {
    describe('when called by an unauthorized node', () => {
      it('reverts', async () => {
        await matchers.evmRevert(async () => {
          await operator
            .connect(roles.stranger)
            .foreward(ethers.constants.AddressZero, [])
        })
      })
    })

    describe('when called by an authorized node', () => {
      describe('when attempting to foreward to the link token', () => {
        it('reverts', async () => {
          await matchers.evmRevert(async () => {
            await operator.connect(roles.oracleNode).foreward(link.address, [])
          })
        })
      })

      describe('when forewarding to any other address', () => {
        const bytes = utils.hexlify(utils.randomBytes(100))
        const payload = getterSetterFactory.interface.functions.setBytes.encode(
          [bytes],
        )
        let mock: contract.Instance<GetterSetterFactory>
        let receipt: ContractReceipt

        beforeEach(async () => {
          mock = await getterSetterFactory
            .connect(roles.defaultAccount)
            .deploy()
          const tx = await operator
            .connect(roles.oracleNode)
            .foreward(mock.address, payload)
          receipt = await tx.wait()
        })

        it('forewards the data', async () => {
          assert.equal(await mock.getBytes(), bytes)
        })

        it('perceives the message is sent by the Operator', async () => {
          const log: any = receipt.logs?.[0]
          const logData = mock.interface.events.SetBytes.decode(
            log.data,
            log.topics,
          )
          assert.equal(utils.getAddress(logData.from), operator.address)
        })
      })
    })
  })
})
