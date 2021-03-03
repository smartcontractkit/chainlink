import {
  contract,
  helpers as h,
  matchers,
  oracle,
  setup,
} from '@chainlink/test-helpers'
import { assert } from 'chai'
import { ethers, utils } from 'ethers'
import { BasicConsumer__factory } from '../../ethers/v0.6/factories/BasicConsumer__factory'
import { MultiWordConsumer__factory } from '../../ethers/v0.6/factories/MultiWordConsumer__factory'
import { GetterSetter__factory } from '../../ethers/v0.4/factories/GetterSetter__factory'
import { MaliciousConsumer__factory } from '../../ethers/v0.4/factories/MaliciousConsumer__factory'
import { MaliciousMultiWordConsumer__factory } from '../../ethers/v0.6/factories/MaliciousMultiWordConsumer__factory'
import { MaliciousRequester__factory } from '../../ethers/v0.4/factories/MaliciousRequester__factory'
import { Operator__factory } from '../../ethers/v0.7/factories/Operator__factory'
import { OperatorForwarder__factory } from '../../ethers/v0.7/factories/OperatorForwarder__factory'
import { Consumer__factory } from '../../ethers/v0.7/factories/Consumer__factory'
import { GasGuzzlingConsumer__factory } from '../../ethers/v0.6/factories/GasGuzzlingConsumer__factory'
import { ContractReceipt } from 'ethers/contract'

const v7ConsumerFactory = new Consumer__factory()
const basicConsumerFactory = new BasicConsumer__factory()
const multiWordConsumerFactory = new MultiWordConsumer__factory()
const gasGuzzlingConsumerFactory = new GasGuzzlingConsumer__factory()
const getterSetterFactory = new GetterSetter__factory()
const maliciousRequesterFactory = new MaliciousRequester__factory()
const maliciousConsumerFactory = new MaliciousConsumer__factory()
const maliciousMultiWordConsumerFactory = new MaliciousMultiWordConsumer__factory()
const operatorFactory = new Operator__factory()
const operatorForwarderFactory = new OperatorForwarder__factory()
const linkTokenFactory = new contract.LinkToken__factory()

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
  let link: contract.Instance<contract.LinkToken__factory>
  let operator: contract.Instance<Operator__factory>
  const deployment = setup.snapshot(provider, async () => {
    link = await linkTokenFactory.connect(roles.defaultAccount).deploy()
    operator = await operatorFactory
      .connect(roles.defaultAccount)
      .deploy(link.address, roles.defaultAccount.address)
    await operator.setAuthorizedSenders([roles.oracleNode.address])
  })

  beforeEach(async () => {
    await deployment()
  })

  it('has a limited public interface', () => {
    matchers.publicAbi(operatorFactory, [
      'EXPIRY_TIME',
      'cancelOracleRequest',
      'fulfillOracleRequest',
      'fulfillOracleRequest2',
      'isAuthorizedSender',
      'getChainlinkToken',
      'onTokenTransfer',
      'oracleRequest',
      'setAuthorizedSenders',
      'getAuthorizedSenders',
      'withdraw',
      'withdrawable',
      'operatorTransferAndCall',
      'distributeFunds',
      'createForwarder',
      'getForwarders',
      // Ownable methods:
      'acceptOwnership',
      'owner',
      'transferOwnership',
    ])
  })

  describe('#createForwarder', () => {
    let receipt: ContractReceipt
    let operatorForwarder: contract.Instance<OperatorForwarder__factory>
    let newSenders: Array<string>
    describe('when called by an authorized sender', () => {
      describe('with 3 authorized senders', () => {
        beforeEach(async () => {
          newSenders = [roles.oracleNode2.address, roles.oracleNode3.address]
          await operator
            .connect(roles.defaultAccount)
            .setAuthorizedSenders(newSenders)

          const tx = await operator.connect(roles.oracleNode2).createForwarder()
          receipt = await tx.wait()
        })

        it('Emits a ForwarderCreated event', async () => {
          const eventFound = h.findEventIn(
            receipt,
            operatorFactory.interface.events.ForwarderCreated,
          )
          assert.exists(eventFound)
        })

        it('adds a forwarder to storage', async () => {
          const forwarders = await operator
            .connect(roles.oracleNode1)
            .getForwarders()
          assert.equal(forwarders.length, 1)
        })

        it('sets the correct authorized senders on the forwarder', async () => {
          const forwarders = await operator
            .connect(roles.oracleNode1)
            .getForwarders()
          operatorForwarder = await operatorForwarderFactory
            .connect(roles.oracleNode1)
            .attach(forwarders[0])
          assert.equal(
            await operatorForwarder.authorizedSender1(),
            roles.defaultAccount.address,
          )
          assert.equal(
            await operatorForwarder.authorizedSender2(),
            newSenders[0],
          )
          assert.equal(
            await operatorForwarder.authorizedSender3(),
            newSenders[1],
          )
        })
      })

      describe('with 1 authorized sender', () => {
        beforeEach(async () => {
          newSenders = [roles.oracleNode2.address]
          await operator
            .connect(roles.defaultAccount)
            .setAuthorizedSenders(newSenders)

          const tx = await operator.connect(roles.oracleNode2).createForwarder()
          receipt = await tx.wait()
        })

        it('Emits a ForwarderCreated event', async () => {
          const eventFound = h.findEventIn(
            receipt,
            operatorFactory.interface.events.ForwarderCreated,
          )
          assert.exists(eventFound)
        })

        it('adds a forwarder to storage', async () => {
          const forwarders = await operator
            .connect(roles.oracleNode1)
            .getForwarders()
          assert.equal(forwarders.length, 1)
        })

        it('sets the correct authorized sender on the forwarder', async () => {
          const forwarders = await operator
            .connect(roles.oracleNode1)
            .getForwarders()
          operatorForwarder = await operatorForwarderFactory
            .connect(roles.oracleNode1)
            .attach(forwarders[0])
          assert.equal(
            await operatorForwarder.authorizedSender1(),
            roles.defaultAccount.address,
          )
          assert.equal(
            await operatorForwarder.authorizedSender2(),
            newSenders[0],
          )
          assert.equal(
            await operatorForwarder.authorizedSender3(),
            '0x0000000000000000000000000000000000000000',
          )
        })
      })
    })
  })

  describe('#distributeFunds', () => {
    describe('when called with empty arrays', () => {
      it('reverts with invalid array message', async () => {
        await matchers.evmRevert(async () => {
          await operator.connect(roles.defaultAccount).distributeFunds([], []),
            'Invalid array length(s)'
        })
      })
    })

    describe('when called with unequal array lengths', () => {
      it('reverts with invalid array message', async () => {
        const receivers = [roles.oracleNode2.address, roles.oracleNode3.address]
        const amounts = [1, 2, 3]
        await matchers.evmRevert(async () => {
          await operator
            .connect(roles.defaultAccount)
            .distributeFunds(receivers, amounts),
            'Invalid array length(s)'
        })
      })
    })

    describe('when called with not enough ETH', () => {
      it('reverts with subtraction overflow message', async () => {
        const amountToSend = h.toWei('2')
        const ethSent = h.toWei('1')
        await matchers.evmRevert(async () => {
          await operator
            .connect(roles.defaultAccount)
            .distributeFunds([roles.oracleNode2.address], [amountToSend], {
              value: ethSent,
            }),
            'SafeMath: subtraction overflow'
        })
      })
    })

    describe('when called with too much ETH', () => {
      it('reverts with too much ETH message', async () => {
        const amountToSend = h.toWei('2')
        const ethSent = h.toWei('3')
        await matchers.evmRevert(async () => {
          await operator
            .connect(roles.defaultAccount)
            .distributeFunds([roles.oracleNode2.address], [amountToSend], {
              value: ethSent,
            }),
            'Too much ETH sent'
        })
      })
    })

    describe('when called with correct values', () => {
      it('updates the balances', async () => {
        const node2BalanceBefore = await roles.oracleNode2.getBalance()
        const node3BalanceBefore = await roles.oracleNode3.getBalance()
        const receivers = [roles.oracleNode2.address, roles.oracleNode3.address]
        const sendNode2 = h.toWei('2')
        const sendNode3 = h.toWei('3')
        const totalAmount = h.toWei('5')
        const amounts = [sendNode2, sendNode3]

        await operator
          .connect(roles.defaultAccount)
          .distributeFunds(receivers, amounts, { value: totalAmount })

        const node2BalanceAfter = await roles.oracleNode2.getBalance()
        const node3BalanceAfter = await roles.oracleNode3.getBalance()

        assert.equal(
          node2BalanceAfter.sub(node2BalanceBefore).toString(),
          sendNode2.toString(),
        )

        assert.equal(
          node3BalanceAfter.sub(node3BalanceBefore).toString(),
          sendNode3.toString(),
        )
      })
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
          const tx = await operator
            .connect(roles.defaultAccount)
            .setAuthorizedSenders(newSenders)
          receipt = await tx.wait()
        })

        it('adds the authorized nodes', async () => {
          const authorizedSenders = await operator.getAuthorizedSenders()
          assert.equal(newSenders.length, authorizedSenders.length)
          for (let i = 0; i < authorizedSenders.length; i++) {
            assert.equal(authorizedSenders[i], newSenders[i])
          }
        })

        it('emits an event', async () => {
          assert.equal(receipt.events?.length, 1)
          const responseEvent = receipt.events?.[0]
          assert.equal(responseEvent?.event, 'AuthorizedSendersChanged')
          const encodedSenders = ethers.utils.defaultAbiCoder.encode(
            ['address[]'],
            [newSenders],
          )
          assert.equal(responseEvent?.data, encodedSenders)
        })

        it('replaces the authorized nodes', async () => {
          const originalAuthorization = await operator
            .connect(roles.defaultAccount)
            .isAuthorizedSender(roles.oracleNode.address)
          assert.isFalse(originalAuthorization)
        })

        afterAll(async () => {
          await operator
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
            await operator
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
          await operator
            .connect(roles.stranger)
            .setAuthorizedSenders([roles.stranger.address])
          ;('Only callable by owner')
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
      let mock: contract.Instance<MaliciousRequester__factory>
      let requester: contract.Instance<BasicConsumer__factory>
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
          await requester.requestEthereumPrice('USD', paymentAmount)
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

    describe('when dataVersion is higher than 255', () => {
      const paid = 100
      const args = oracle.encodeOracleRequest(specId, to, fHash, 1, '0x0', 256)

      it('throws an error', async () => {
        await matchers.evmRevert(async () => {
          await link.transferAndCall(operator.address, paid, args)
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
    let maliciousRequester: contract.Instance<MaliciousRequester__factory>
    let basicConsumer: contract.Instance<BasicConsumer__factory>
    let maliciousConsumer: contract.Instance<MaliciousConsumer__factory>
    let gasGuzzlingConsumer: contract.Instance<GasGuzzlingConsumer__factory>
    let request: ReturnType<typeof oracle.decodeRunRequest>

    describe('gas guzzling consumer', () => {
      beforeEach(async () => {
        gasGuzzlingConsumer = await gasGuzzlingConsumerFactory
          .connect(roles.consumer)
          .deploy(link.address, operator.address, specId)
        const paymentAmount = h.toWei('1')
        await link.transfer(gasGuzzlingConsumer.address, paymentAmount)
        const tx = await gasGuzzlingConsumer.gassyRequestEthereumPrice(
          paymentAmount,
        )
        const receipt = await tx.wait()
        request = oracle.decodeRunRequest(receipt.logs?.[3])
      })

      it('emits an OracleResponse event', async () => {
        const fulfillParams = oracle.convertFufillParams(request, response)
        const tx = await operator
          .connect(roles.oracleNode)
          .fulfillOracleRequest(...fulfillParams)
        const receipt = await tx.wait()
        assert.equal(receipt.events?.length, 1)
        const responseEvent = receipt.events?.[0]
        assert.equal(responseEvent?.event, 'OracleResponse')
        assert.equal(responseEvent?.args?.[0], request.requestId)
      })
    })

    describe('cooperative consumer', () => {
      beforeEach(async () => {
        basicConsumer = await basicConsumerFactory
          .connect(roles.defaultAccount)
          .deploy(link.address, operator.address, specId)
        const paymentAmount = h.toWei('1')
        await link.transfer(basicConsumer.address, paymentAmount)
        const currency = 'USD'
        const tx = await basicConsumer.requestEthereumPrice(
          currency,
          paymentAmount,
        )
        const receipt = await tx.wait()
        request = oracle.decodeRunRequest(receipt.logs?.[3])
      })

      describe('when called by an unauthorized node', () => {
        beforeEach(async () => {
          assert.equal(
            false,
            await operator.isAuthorizedSender(roles.stranger.address),
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

      describe('when fulfilled with the wrong function', () => {
        let v7Consumer
        beforeEach(async () => {
          v7Consumer = await v7ConsumerFactory
            .connect(roles.defaultAccount)
            .deploy(link.address, operator.address, specId)
          const paymentAmount = h.toWei('1')
          await link.transfer(v7Consumer.address, paymentAmount)
          const currency = 'USD'
          const tx = await v7Consumer.requestEthereumPrice(
            currency,
            paymentAmount,
          )
          const receipt = await tx.wait()
          request = oracle.decodeRunRequest(receipt.logs?.[3])
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

        it('emits an OracleResponse event', async () => {
          const fulfillParams = oracle.convertFufillParams(request, response)
          const tx = await operator
            .connect(roles.oracleNode)
            .fulfillOracleRequest(...fulfillParams)
          const receipt = await tx.wait()
          assert.equal(receipt.events?.length, 3)
          const responseEvent = receipt.events?.[0]
          assert.equal(responseEvent?.event, 'OracleResponse')
          assert.equal(responseEvent?.args?.[0], request.requestId)
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

  describe('#fulfillOracleRequest2', () => {
    describe('single word fulfils', () => {
      const response = 'Hi mom!'
      const responseTypes = ['bytes32']
      const responseValues = [h.toBytes32String(response)]
      let maliciousRequester: contract.Instance<MaliciousRequester__factory>
      let basicConsumer: contract.Instance<BasicConsumer__factory>
      let maliciousConsumer: contract.Instance<MaliciousConsumer__factory>
      let gasGuzzlingConsumer: contract.Instance<GasGuzzlingConsumer__factory>
      let request: ReturnType<typeof oracle.decodeRunRequest>

      describe('gas guzzling consumer', () => {
        beforeEach(async () => {
          gasGuzzlingConsumer = await gasGuzzlingConsumerFactory
            .connect(roles.consumer)
            .deploy(link.address, operator.address, specId)
          const paymentAmount = h.toWei('1')
          await link.transfer(gasGuzzlingConsumer.address, paymentAmount)
          const tx = await gasGuzzlingConsumer.gassyRequestEthereumPrice(
            paymentAmount,
          )
          const receipt = await tx.wait()
          request = oracle.decodeRunRequest(receipt.logs?.[3])
        })

        it('emits an OracleResponse2 event', async () => {
          const fulfillParams = oracle.convertFulfill2Params(
            request,
            responseTypes,
            responseValues,
          )
          const tx = await operator
            .connect(roles.oracleNode)
            .fulfillOracleRequest2(...fulfillParams)
          const receipt = await tx.wait()
          assert.equal(receipt.events?.length, 1)
          const responseEvent = receipt.events?.[0]
          assert.equal(responseEvent?.event, 'OracleResponse')
          assert.equal(responseEvent?.args?.[0], request.requestId)
        })
      })

      describe('cooperative consumer', () => {
        beforeEach(async () => {
          basicConsumer = await basicConsumerFactory
            .connect(roles.defaultAccount)
            .deploy(link.address, operator.address, specId)
          const paymentAmount = h.toWei('1')
          await link.transfer(basicConsumer.address, paymentAmount)
          const currency = 'USD'
          const tx = await basicConsumer.requestEthereumPrice(
            currency,
            paymentAmount,
          )
          const receipt = await tx.wait()
          request = oracle.decodeRunRequest(receipt.logs?.[3])
        })

        describe('when called by an unauthorized node', () => {
          beforeEach(async () => {
            assert.equal(
              false,
              await operator.isAuthorizedSender(roles.stranger.address),
            )
          })

          it('raises an error', async () => {
            await matchers.evmRevert(async () => {
              await operator
                .connect(roles.stranger)
                .fulfillOracleRequest2(
                  ...oracle.convertFulfill2Params(
                    request,
                    responseTypes,
                    responseValues,
                  ),
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
                .fulfillOracleRequest2(
                  ...oracle.convertFulfill2Params(
                    request,
                    responseTypes,
                    responseValues,
                  ),
                )
            })
          })

          it('sets the value on the requested contract', async () => {
            await operator
              .connect(roles.oracleNode)
              .fulfillOracleRequest2(
                ...oracle.convertFulfill2Params(
                  request,
                  responseTypes,
                  responseValues,
                ),
              )

            const currentValue = await basicConsumer.currentPrice()
            assert.equal(
              response,
              ethers.utils.parseBytes32String(currentValue),
            )
          })

          it('emits an OracleResponse2 event', async () => {
            const fulfillParams = oracle.convertFulfill2Params(
              request,
              responseTypes,
              responseValues,
            )
            const tx = await operator
              .connect(roles.oracleNode)
              .fulfillOracleRequest2(...fulfillParams)
            const receipt = await tx.wait()
            assert.equal(receipt.events?.length, 3)
            const responseEvent = receipt.events?.[0]
            assert.equal(responseEvent?.event, 'OracleResponse')
            assert.equal(responseEvent?.args?.[0], request.requestId)
          })

          it('does not allow a request to be fulfilled twice', async () => {
            const response2 = response + ' && Hello World!!'
            const response2Values = [h.toBytes32String(response2)]
            await operator
              .connect(roles.oracleNode)
              .fulfillOracleRequest2(
                ...oracle.convertFulfill2Params(
                  request,
                  responseTypes,
                  responseValues,
                ),
              )

            await matchers.evmRevert(async () => {
              await operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(
                  ...oracle.convertFulfill2Params(
                    request,
                    responseTypes,
                    response2Values,
                  ),
                )
            })

            const currentValue = await basicConsumer.currentPrice()
            assert.equal(
              response,
              ethers.utils.parseBytes32String(currentValue),
            )
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
              await operator.connect(roles.oracleNode).fulfillOracleRequest2(
                ...oracle.convertFulfill2Params(
                  request,
                  responseTypes,
                  responseValues,
                  {
                    gasLimit: 70000,
                  },
                ),
              )
            })

            matchers.bigNum(0, await operator.withdrawable())
          })

          it(`${defaultGasLimit} is enough to pass the gas requirement`, async () => {
            await operator.connect(roles.oracleNode).fulfillOracleRequest2(
              ...oracle.convertFulfill2Params(
                request,
                responseTypes,
                responseValues,
                {
                  gasLimit: defaultGasLimit,
                },
              ),
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
              .fulfillOracleRequest2(
                ...oracle.convertFulfill2Params(
                  request,
                  responseTypes,
                  responseValues,
                ),
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
            const response2Values = [h.toBytes32String(response2)]
            await operator
              .connect(roles.oracleNode)
              .fulfillOracleRequest2(
                ...oracle.convertFulfill2Params(
                  request,
                  responseTypes,
                  responseValues,
                ),
              )

            await matchers.evmRevert(async () => {
              await operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(
                  ...oracle.convertFulfill2Params(
                    request,
                    responseTypes,
                    response2Values,
                  ),
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
              .fulfillOracleRequest2(
                ...oracle.convertFulfill2Params(
                  request,
                  responseTypes,
                  responseValues,
                ),
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
              ethers.utils.toUtf8Bytes(
                'cancelRequestOnFulfill(bytes32,bytes32)',
              ),
            )
            const receipt = await tx.wait()
            request = oracle.decodeRunRequest(receipt.logs?.[3])

            matchers.bigNum(0, await link.balanceOf(maliciousConsumer.address))
          })

          it('allows the oracle node to receive their payment', async () => {
            await operator
              .connect(roles.oracleNode)
              .fulfillOracleRequest2(
                ...oracle.convertFulfill2Params(
                  request,
                  responseTypes,
                  responseValues,
                ),
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
            const response2Values = [h.toBytes32String(response2)]

            await operator
              .connect(roles.oracleNode)
              .fulfillOracleRequest2(
                ...oracle.convertFulfill2Params(
                  request,
                  responseTypes,
                  responseValues,
                ),
              )

            await matchers.evmRevert(async () => {
              await operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(
                  ...oracle.convertFulfill2Params(
                    request,
                    responseTypes,
                    response2Values,
                  ),
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
              .fulfillOracleRequest2(
                ...oracle.convertFulfill2Params(
                  request,
                  responseTypes,
                  responseValues,
                ),
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
              .fulfillOracleRequest2(
                ...oracle.convertFulfill2Params(
                  request,
                  responseTypes,
                  responseValues,
                ),
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
              .fulfillOracleRequest2(
                ...oracle.convertFulfill2Params(
                  request,
                  responseTypes,
                  responseValues,
                ),
              )
            matchers.bigNum(
              0,
              await provider.getBalance(maliciousConsumer.address),
            )
          })
        })
      })
    })

    describe('multi word fulfils', () => {
      describe('one bytes parameter', () => {
        const response =
          'Lorem ipsum dolor sit amet, consectetur adipiscing elit.\
          Fusce euismod malesuada ligula, eget semper metus ultrices sit amet.'
        const responseTypes = ['bytes']
        const responseValues = [h.stringToBytes(response)]
        let maliciousRequester: contract.Instance<MaliciousRequester__factory>
        let multiConsumer: contract.Instance<MultiWordConsumer__factory>
        let maliciousConsumer: contract.Instance<MaliciousMultiWordConsumer__factory>
        let gasGuzzlingConsumer: contract.Instance<GasGuzzlingConsumer__factory>
        let request: ReturnType<typeof oracle.decodeRunRequest>

        describe('gas guzzling consumer', () => {
          beforeEach(async () => {
            gasGuzzlingConsumer = await gasGuzzlingConsumerFactory
              .connect(roles.consumer)
              .deploy(link.address, operator.address, specId)
            const paymentAmount = h.toWei('1')
            await link.transfer(gasGuzzlingConsumer.address, paymentAmount)
            const tx = await gasGuzzlingConsumer.gassyMultiWordRequest(
              paymentAmount,
            )
            const receipt = await tx.wait()
            request = oracle.decodeRunRequest(receipt.logs?.[3])
          })

          it('emits an OracleResponse2 event', async () => {
            const fulfillParams = oracle.convertFulfill2Params(
              request,
              responseTypes,
              responseValues,
            )
            const tx = await operator
              .connect(roles.oracleNode)
              .fulfillOracleRequest2(...fulfillParams)
            const receipt = await tx.wait()
            assert.equal(receipt.events?.length, 1)
            const responseEvent = receipt.events?.[0]
            assert.equal(responseEvent?.event, 'OracleResponse')
            assert.equal(responseEvent?.args?.[0], request.requestId)
          })
        })

        describe('cooperative consumer', () => {
          beforeEach(async () => {
            multiConsumer = await multiWordConsumerFactory
              .connect(roles.defaultAccount)
              .deploy(link.address, operator.address, specId)
            const paymentAmount = h.toWei('1')
            await link.transfer(multiConsumer.address, paymentAmount)
            const currency = 'USD'
            const tx = await multiConsumer.requestEthereumPrice(
              currency,
              paymentAmount,
            )
            const receipt = await tx.wait()
            request = oracle.decodeRunRequest(receipt.logs?.[3])
          })

          describe('when called by an unauthorized node', () => {
            beforeEach(async () => {
              assert.equal(
                false,
                await operator.isAuthorizedSender(roles.stranger.address),
              )
            })

            it('raises an error', async () => {
              await matchers.evmRevert(async () => {
                await operator
                  .connect(roles.stranger)
                  .fulfillOracleRequest2(
                    ...oracle.convertFulfill2Params(
                      request,
                      responseTypes,
                      responseValues,
                    ),
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
                  .fulfillOracleRequest2(
                    ...oracle.convertFulfill2Params(
                      request,
                      responseTypes,
                      responseValues,
                    ),
                  )
              })
            })

            it('sets the value on the requested contract', async () => {
              await operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(
                  ...oracle.convertFulfill2Params(
                    request,
                    responseTypes,
                    responseValues,
                  ),
                )

              const currentValue = await multiConsumer.currentPrice()
              assert.equal(response, ethers.utils.toUtf8String(currentValue))
            })

            it('emits an OracleResponse2 event', async () => {
              const fulfillParams = oracle.convertFulfill2Params(
                request,
                responseTypes,
                responseValues,
              )
              const tx = await operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(...fulfillParams)
              const receipt = await tx.wait()
              assert.equal(receipt.events?.length, 3)
              const responseEvent = receipt.events?.[0]
              assert.equal(responseEvent?.event, 'OracleResponse')
              assert.equal(responseEvent?.args?.[0], request.requestId)
            })

            it('does not allow a request to be fulfilled twice', async () => {
              const response2 = response + ' && Hello World!!'
              const response2Values = [h.stringToBytes(response2)]

              await operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(
                  ...oracle.convertFulfill2Params(
                    request,
                    responseTypes,
                    responseValues,
                  ),
                )

              await matchers.evmRevert(async () => {
                await operator
                  .connect(roles.oracleNode)
                  .fulfillOracleRequest2(
                    ...oracle.convertFulfill2Params(
                      request,
                      responseTypes,
                      response2Values,
                    ),
                  )
              })

              const currentValue = await multiConsumer.currentPrice()
              assert.equal(response, ethers.utils.toUtf8String(currentValue))
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
                await operator.connect(roles.oracleNode).fulfillOracleRequest2(
                  ...oracle.convertFulfill2Params(
                    request,
                    responseTypes,
                    responseValues,
                    {
                      gasLimit: 70000,
                    },
                  ),
                )
              })

              matchers.bigNum(0, await operator.withdrawable())
            })

            it(`${defaultGasLimit} is enough to pass the gas requirement`, async () => {
              await operator.connect(roles.oracleNode).fulfillOracleRequest2(
                ...oracle.convertFulfill2Params(
                  request,
                  responseTypes,
                  responseValues,
                  {
                    gasLimit: defaultGasLimit,
                  },
                ),
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
            maliciousConsumer = await maliciousMultiWordConsumerFactory
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
                .fulfillOracleRequest2(
                  ...oracle.convertFulfill2Params(
                    request,
                    responseTypes,
                    responseValues,
                  ),
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
              const response2Values = [h.stringToBytes(response2)]
              await operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(
                  ...oracle.convertFulfill2Params(
                    request,
                    responseTypes,
                    responseValues,
                  ),
                )

              await matchers.evmRevert(async () => {
                await operator
                  .connect(roles.oracleNode)
                  .fulfillOracleRequest2(
                    ...oracle.convertFulfill2Params(
                      request,
                      responseTypes,
                      response2Values,
                    ),
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
                .fulfillOracleRequest2(
                  ...oracle.convertFulfill2Params(
                    request,
                    responseTypes,
                    responseValues,
                  ),
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
                ethers.utils.toUtf8Bytes(
                  'cancelRequestOnFulfill(bytes32,bytes32)',
                ),
              )
              const receipt = await tx.wait()
              request = oracle.decodeRunRequest(receipt.logs?.[3])

              matchers.bigNum(
                0,
                await link.balanceOf(maliciousConsumer.address),
              )
            })

            it('allows the oracle node to receive their payment', async () => {
              await operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(
                  ...oracle.convertFulfill2Params(
                    request,
                    responseTypes,
                    responseValues,
                  ),
                )

              const mockBalance = await link.balanceOf(
                maliciousConsumer.address,
              )
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
              const response2Values = [h.stringToBytes(response2)]
              await operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(
                  ...oracle.convertFulfill2Params(
                    request,
                    responseTypes,
                    responseValues,
                  ),
                )

              await matchers.evmRevert(async () => {
                await operator
                  .connect(roles.oracleNode)
                  .fulfillOracleRequest2(
                    ...oracle.convertFulfill2Params(
                      request,
                      responseTypes,
                      response2Values,
                    ),
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
                .fulfillOracleRequest2(
                  ...oracle.convertFulfill2Params(
                    request,
                    responseTypes,
                    responseValues,
                  ),
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
                .fulfillOracleRequest2(
                  ...oracle.convertFulfill2Params(
                    request,
                    responseTypes,
                    responseValues,
                  ),
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
                .fulfillOracleRequest2(
                  ...oracle.convertFulfill2Params(
                    request,
                    responseTypes,
                    responseValues,
                  ),
                )
              matchers.bigNum(
                0,
                await provider.getBalance(maliciousConsumer.address),
              )
            })
          })
        })
      })

      describe('multiple bytes32 parameters', () => {
        const response1 = 'Hi mom!'
        const response2 = 'Its me!'
        const responseTypes = ['bytes32', 'bytes32']
        const responseValues = [
          h.toBytes32String(response1),
          h.toBytes32String(response2),
        ]
        let maliciousRequester: contract.Instance<MaliciousRequester__factory>
        let multiConsumer: contract.Instance<MultiWordConsumer__factory>
        let maliciousConsumer: contract.Instance<MaliciousMultiWordConsumer__factory>
        let gasGuzzlingConsumer: contract.Instance<GasGuzzlingConsumer__factory>
        let request: ReturnType<typeof oracle.decodeRunRequest>

        describe('gas guzzling consumer', () => {
          beforeEach(async () => {
            gasGuzzlingConsumer = await gasGuzzlingConsumerFactory
              .connect(roles.consumer)
              .deploy(link.address, operator.address, specId)
            const paymentAmount = h.toWei('1')
            await link.transfer(gasGuzzlingConsumer.address, paymentAmount)
            const tx = await gasGuzzlingConsumer.gassyMultiWordRequest(
              paymentAmount,
            )
            const receipt = await tx.wait()
            request = oracle.decodeRunRequest(receipt.logs?.[3])
          })

          it('emits an OracleResponse2 event', async () => {
            const fulfillParams = oracle.convertFulfill2Params(
              request,
              responseTypes,
              responseValues,
            )
            const tx = await operator
              .connect(roles.oracleNode)
              .fulfillOracleRequest2(...fulfillParams)
            const receipt = await tx.wait()
            assert.equal(receipt.events?.length, 1)
            const responseEvent = receipt.events?.[0]
            assert.equal(responseEvent?.event, 'OracleResponse')
            assert.equal(responseEvent?.args?.[0], request.requestId)
          })
        })

        describe('cooperative consumer', () => {
          beforeEach(async () => {
            multiConsumer = await multiWordConsumerFactory
              .connect(roles.defaultAccount)
              .deploy(link.address, operator.address, specId)
            const paymentAmount = h.toWei('1')
            await link.transfer(multiConsumer.address, paymentAmount)
            const currency = 'USD'
            const tx = await multiConsumer.requestMultipleParameters(
              currency,
              paymentAmount,
            )
            const receipt = await tx.wait()
            request = oracle.decodeRunRequest(receipt.logs?.[3])
          })

          describe('when called by an unauthorized node', () => {
            beforeEach(async () => {
              assert.equal(
                false,
                await operator.isAuthorizedSender(roles.stranger.address),
              )
            })

            it('raises an error', async () => {
              await matchers.evmRevert(async () => {
                await operator
                  .connect(roles.stranger)
                  .fulfillOracleRequest2(
                    ...oracle.convertFulfill2Params(
                      request,
                      responseTypes,
                      responseValues,
                    ),
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
                  .fulfillOracleRequest2(
                    ...oracle.convertFulfill2Params(
                      request,
                      responseTypes,
                      responseValues,
                    ),
                  )
              })
            })

            it('sets the value on the requested contract', async () => {
              await operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(
                  ...oracle.convertFulfill2Params(
                    request,
                    responseTypes,
                    responseValues,
                  ),
                )

              const firstValue = await multiConsumer.first()
              const secondValue = await multiConsumer.second()
              assert.equal(
                response1,
                ethers.utils.parseBytes32String(firstValue),
              )
              assert.equal(
                response2,
                ethers.utils.parseBytes32String(secondValue),
              )
            })

            it('emits an OracleResponse2 event', async () => {
              const fulfillParams = oracle.convertFulfill2Params(
                request,
                responseTypes,
                responseValues,
              )
              const tx = await operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(...fulfillParams)
              const receipt = await tx.wait()
              assert.equal(receipt.events?.length, 3)
              const responseEvent = receipt.events?.[0]
              assert.equal(responseEvent?.event, 'OracleResponse')
              assert.equal(responseEvent?.args?.[0], request.requestId)
            })

            it('does not allow a request to be fulfilled twice', async () => {
              const response3 = response2 + ' && Hello World!!'
              const repeatedResponseValues = [
                h.toBytes32String(response2),
                h.toBytes32String(response3),
              ]

              await operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(
                  ...oracle.convertFulfill2Params(
                    request,
                    responseTypes,
                    responseValues,
                  ),
                )

              await matchers.evmRevert(async () => {
                await operator
                  .connect(roles.oracleNode)
                  .fulfillOracleRequest2(
                    ...oracle.convertFulfill2Params(
                      request,
                      responseTypes,
                      repeatedResponseValues,
                    ),
                  )
              })

              const firstValue = await multiConsumer.first()
              const secondValue = await multiConsumer.second()
              assert.equal(
                response1,
                ethers.utils.parseBytes32String(firstValue),
              )
              assert.equal(
                response2,
                ethers.utils.parseBytes32String(secondValue),
              )
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
                await operator.connect(roles.oracleNode).fulfillOracleRequest2(
                  ...oracle.convertFulfill2Params(
                    request,
                    responseTypes,
                    responseValues,
                    {
                      gasLimit: 70000,
                    },
                  ),
                )
              })

              matchers.bigNum(0, await operator.withdrawable())
            })

            it(`${defaultGasLimit} is enough to pass the gas requirement`, async () => {
              await operator.connect(roles.oracleNode).fulfillOracleRequest2(
                ...oracle.convertFulfill2Params(
                  request,
                  responseTypes,
                  responseValues,
                  {
                    gasLimit: defaultGasLimit,
                  },
                ),
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
            maliciousConsumer = await maliciousMultiWordConsumerFactory
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
                .fulfillOracleRequest2(
                  ...oracle.convertFulfill2Params(
                    request,
                    responseTypes,
                    responseValues,
                  ),
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
              const response3 = 'hack the planet 102'
              const repeatedResponseValues = [
                h.toBytes32String(response2),
                h.toBytes32String(response3),
              ]
              await operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(
                  ...oracle.convertFulfill2Params(
                    request,
                    responseTypes,
                    responseValues,
                  ),
                )

              await matchers.evmRevert(async () => {
                await operator
                  .connect(roles.oracleNode)
                  .fulfillOracleRequest2(
                    ...oracle.convertFulfill2Params(
                      request,
                      responseTypes,
                      repeatedResponseValues,
                    ),
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
                .fulfillOracleRequest2(
                  ...oracle.convertFulfill2Params(
                    request,
                    responseTypes,
                    responseValues,
                  ),
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
                ethers.utils.toUtf8Bytes(
                  'cancelRequestOnFulfill(bytes32,bytes32)',
                ),
              )
              const receipt = await tx.wait()
              request = oracle.decodeRunRequest(receipt.logs?.[3])

              matchers.bigNum(
                0,
                await link.balanceOf(maliciousConsumer.address),
              )
            })

            it('allows the oracle node to receive their payment', async () => {
              await operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(
                  ...oracle.convertFulfill2Params(
                    request,
                    responseTypes,
                    responseValues,
                  ),
                )

              const mockBalance = await link.balanceOf(
                maliciousConsumer.address,
              )
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
              const response3 = 'hack the planet 102'
              const repeatedResponseValues = [
                h.toBytes32String(response2),
                h.toBytes32String(response3),
              ]
              await operator
                .connect(roles.oracleNode)
                .fulfillOracleRequest2(
                  ...oracle.convertFulfill2Params(
                    request,
                    responseTypes,
                    responseValues,
                  ),
                )

              await matchers.evmRevert(async () => {
                await operator
                  .connect(roles.oracleNode)
                  .fulfillOracleRequest2(
                    ...oracle.convertFulfill2Params(
                      request,
                      responseTypes,
                      repeatedResponseValues,
                    ),
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
                .fulfillOracleRequest2(
                  ...oracle.convertFulfill2Params(
                    request,
                    responseTypes,
                    responseValues,
                  ),
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
                .fulfillOracleRequest2(
                  ...oracle.convertFulfill2Params(
                    request,
                    responseTypes,
                    responseValues,
                  ),
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
                .fulfillOracleRequest2(
                  ...oracle.convertFulfill2Params(
                    request,
                    responseTypes,
                    responseValues,
                  ),
                )
              matchers.bigNum(
                0,
                await provider.getBalance(maliciousConsumer.address),
              )
            })
          })
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

      describe('recovering funds that were mistakenly sent', () => {
        const paid = 1
        beforeEach(async () => {
          await link.transfer(operator.address, paid)
        })

        it('withdraws funds', async () => {
          const operatorBalanceBefore = await link.balanceOf(operator.address)
          const accountBalanceBefore = await link.balanceOf(
            roles.defaultAccount.address,
          )

          await operator
            .connect(roles.defaultAccount)
            .withdraw(roles.defaultAccount.address, paid)

          const operatorBalanceAfter = await link.balanceOf(operator.address)
          const accountBalanceAfter = await link.balanceOf(
            roles.defaultAccount.address,
          )

          const accountDifference = accountBalanceAfter.sub(
            accountBalanceBefore,
          )
          const operatorDifference = operatorBalanceBefore.sub(
            operatorBalanceAfter,
          )

          matchers.bigNum(operatorDifference, paid)
          matchers.bigNum(accountDifference, paid)
        })
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

        describe('recovering funds that were mistakenly sent', () => {
          const paid = 1
          beforeEach(async () => {
            await link.transfer(operator.address, paid)
          })

          it('withdraws funds', async () => {
            const operatorBalanceBefore = await link.balanceOf(operator.address)
            const accountBalanceBefore = await link.balanceOf(
              roles.defaultAccount.address,
            )

            await operator
              .connect(roles.defaultAccount)
              .withdraw(roles.defaultAccount.address, paid)

            const operatorBalanceAfter = await link.balanceOf(operator.address)
            const accountBalanceAfter = await link.balanceOf(
              roles.defaultAccount.address,
            )

            const accountDifference = accountBalanceAfter.sub(
              accountBalanceBefore,
            )
            const operatorDifference = operatorBalanceBefore.sub(
              operatorBalanceAfter,
            )

            matchers.bigNum(operatorDifference, paid)
            matchers.bigNum(accountDifference, paid)
          })
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

        describe('recovering funds that were mistakenly sent', () => {
          const paid = 1
          beforeEach(async () => {
            await link.transfer(operator.address, paid)
          })

          it('withdraws funds', async () => {
            const operatorBalanceBefore = await link.balanceOf(operator.address)
            const accountBalanceBefore = await link.balanceOf(
              roles.defaultAccount.address,
            )

            await operator
              .connect(roles.defaultAccount)
              .withdraw(roles.defaultAccount.address, paid)

            const operatorBalanceAfter = await link.balanceOf(operator.address)
            const accountBalanceAfter = await link.balanceOf(
              roles.defaultAccount.address,
            )

            const accountDifference = accountBalanceAfter.sub(
              accountBalanceBefore,
            )
            const operatorDifference = operatorBalanceBefore.sub(
              operatorBalanceAfter,
            )

            matchers.bigNum(operatorDifference, paid)
            matchers.bigNum(accountDifference, paid)
          })
        })
      })
    })
  })

  describe('#withdrawable', () => {
    let request: ReturnType<typeof oracle.decodeRunRequest>
    const amount = h.toWei('1')

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

    describe('funds that were mistakenly sent', () => {
      const paid = 1
      beforeEach(async () => {
        await link.transfer(operator.address, paid)
      })

      it('returns the correct value', async () => {
        const withdrawAmount = await operator.withdrawable()

        const expectedAmount = amount.add(paid)
        matchers.bigNum(withdrawAmount, expectedAmount)
      })
    })
  })

  describe('#operatorTransferAndCall', () => {
    let operator2: contract.Instance<Operator__factory>
    let args: string
    let to: string
    const startingBalance = 1000
    const payment = 20

    beforeEach(async () => {
      operator2 = await operatorFactory
        .connect(roles.oracleNode2)
        .deploy(link.address, roles.oracleNode2.address)
      to = operator2.address
      args = oracle.encodeOracleRequest(
        specId,
        operator.address,
        operatorFactory.interface.functions.fulfillOracleRequest.sighash,
        1,
        '0x0',
      )
    })

    describe('when called by a non-owner', () => {
      it('reverts with owner error message', async () => {
        await link.transfer(operator.address, startingBalance)
        await matchers.evmRevert(async () => {
          await operator
            .connect(roles.stranger)
            .operatorTransferAndCall(to, payment, args),
            'Only callable by owner'
        })
      })
    })

    describe('when called by the owner', () => {
      beforeEach(async () => {
        await link.transfer(operator.address, startingBalance)
      })

      describe('without sufficient funds in contract', () => {
        it('reverts with funds message', async () => {
          const tooMuch = startingBalance * 2
          await matchers.evmRevert(async () => {
            await operator
              .connect(roles.stranger)
              .operatorTransferAndCall(to, tooMuch, args),
              'Amount requested is greater than withdrawable balance'
          })
        })
      })

      describe('with sufficient funds', () => {
        let receipt: ContractReceipt
        let requesterBalanceBefore: utils.BigNumber
        let requesterBalanceAfter: utils.BigNumber
        let receiverBalanceBefore: utils.BigNumber
        let receiverBalanceAfter: utils.BigNumber

        beforeAll(async () => {
          requesterBalanceBefore = await link.balanceOf(operator.address)
          receiverBalanceBefore = await link.balanceOf(operator2.address)
          const tx = await operator
            .connect(roles.defaultAccount)
            .operatorTransferAndCall(to, payment, args)
          receipt = await tx.wait()
          requesterBalanceAfter = await link.balanceOf(operator.address)
          receiverBalanceAfter = await link.balanceOf(operator2.address)
        })

        it('emits an event', async () => {
          assert.equal(3, receipt.logs?.length)
          const request = oracle.decodeRunRequest(receipt.logs?.[2])
          assert.equal(request.requester, operator.address)
        })

        it('transfers the tokens', async () => {
          matchers.bigNum(
            requesterBalanceBefore.sub(requesterBalanceAfter),
            payment,
          )
          matchers.bigNum(
            receiverBalanceAfter.sub(receiverBalanceBefore),
            payment,
          )
        })
      })
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
})
