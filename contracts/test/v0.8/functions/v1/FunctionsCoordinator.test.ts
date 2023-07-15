// import { ethers } from 'hardhat'
// import { BigNumber } from 'ethers'
// import { expect } from 'chai'
// import {
//   getSetupFactory,
//   FunctionsContracts,
//   FunctionsRoles,
//   stringToHex,
//   encodeReport,
//   anyValue,
//   createSubscription,
//   ids,
// } from './utils'
// import { randomAddressString } from 'hardhat/internal/hardhat-network/provider/utils/random'
// import { stringToBytes } from '../../../test-helpers/helpers'

// const setup = getSetupFactory()
// let contracts: FunctionsContracts
// let roles: FunctionsRoles

// const donLabel = ethers.utils.formatBytes32String('1')

// beforeEach(async () => {
//   ;({ contracts, roles } = setup())
// })

// describe('Functions Coordinator', () => {
//   let subscriptionId: number
//   const donPublicKey =
//     '0x3804a19f2437f7bba4fcfbc194379e43e514aa98073db3528ccdbdb642e24011'
//   let transmitters: string[]

//   describe('General', () => {
//     it('#typeAndVersion', async () => {
//       expect(await contracts.coordinator.typeAndVersion()).to.be.equal(
//         'Functions Coordinator v1',
//       )
//     })

//     it('returns DON public key set on this Oracle', async () => {
//       await expect(contracts.coordinator.setDONPublicKey(donPublicKey)).not.to
//         .be.reverted
//       expect(
//         await contracts.coordinator.callStatic.getDONPublicKey(),
//       ).to.be.equal(donPublicKey)
//     })

//     it('reverts setDONPublicKey for empty data', async () => {
//       const emptyPublicKey = stringToHex('')
//       await expect(
//         contracts.coordinator.setDONPublicKey(emptyPublicKey),
//       ).to.be.revertedWith('EmptyPublicKey')
//     })

//     async function validatePubKeys(
//       expectedNodes: string[],
//       expectedKeys: string[],
//     ) {
//       const allNodesAndKeys = await contracts.coordinator.getAllNodePublicKeys()
//       for (let i = 0; i < expectedNodes.length; i++) {
//         expect(allNodesAndKeys[0][i]).to.be.equal(expectedNodes[i])
//         expect(allNodesAndKeys[1][i]).to.be.equal(expectedKeys[i])
//       }
//     }

//     it('set/delete/get node public keys', async () => {
//       const emptyKey = stringToHex('')
//       const publicKey2 = stringToHex('key420')
//       const publicKey3 = stringToHex('key666')

//       await contracts.coordinator.setNodePublicKey(
//         roles.oracleNode2.getAddress(),
//         publicKey2,
//       )
//       await contracts.coordinator.setNodePublicKey(
//         roles.oracleNode3.getAddress(),
//         publicKey3,
//       )
//       validatePubKeys(transmitters, [
//         emptyKey,
//         publicKey2,
//         publicKey3,
//         emptyKey,
//       ])

//       await contracts.coordinator.deleteNodePublicKey(
//         roles.oracleNode1.getAddress(),
//       )
//       await contracts.coordinator.deleteNodePublicKey(
//         roles.oracleNode2.getAddress(),
//       )
//       validatePubKeys(transmitters, [emptyKey, emptyKey, publicKey3, emptyKey])
//     })

//     it('reverts setNodePublicKey for unauthorized callers', async () => {
//       const pubKey = stringToHex('abcd')

//       await expect(
//         contracts.coordinator
//           .connect(roles.defaultAccount)
//           .setNodePublicKey(roles.oracleNode2.getAddress(), pubKey),
//       ).not.to.be.reverted

//       await expect(
//         contracts.coordinator
//           .connect(roles.consumer)
//           .setNodePublicKey(roles.oracleNode2.getAddress(), pubKey),
//       ).to.be.revertedWith('UnauthorizedPublicKeyChange')

//       await expect(
//         contracts.coordinator
//           .connect(roles.consumer)
//           .setNodePublicKey(roles.consumer.getAddress(), pubKey),
//       ).to.be.revertedWith('UnauthorizedPublicKeyChange')

//       await expect(
//         contracts.coordinator
//           .connect(roles.oracleNode2)
//           .setNodePublicKey(roles.oracleNode3.getAddress(), pubKey),
//       ).to.be.revertedWith('UnauthorizedPublicKeyChange')

//       await expect(
//         contracts.coordinator
//           .connect(roles.oracleNode2)
//           .setNodePublicKey(roles.oracleNode2.getAddress(), pubKey),
//       ).not.to.be.reverted
//     })

//     it('reverts deleteNodePublicKey for unauthorized callers', async () => {
//       await expect(
//         contracts.coordinator
//           .connect(roles.defaultAccount)
//           .deleteNodePublicKey(roles.oracleNode2.getAddress()),
//       ).not.to.be.reverted

//       await expect(
//         contracts.coordinator
//           .connect(roles.consumer)
//           .deleteNodePublicKey(roles.oracleNode2.getAddress()),
//       ).to.be.revertedWith('UnauthorizedPublicKeyChange')

//       await expect(
//         contracts.coordinator
//           .connect(roles.consumer)
//           .deleteNodePublicKey(roles.consumer.getAddress()),
//       ).not.to.be.reverted
//     })
//   })

//   describe('Sending requests', () => {
//     it('#sendRequest emits OracleRequest event', async () => {
//       subscriptionId = await createSubscription(
//         roles.subOwner,
//         [contracts.client.address],
//         contracts.router,
//         contracts.linkToken,
//       )
//       const defaultAccountAddress = await roles.defaultAccount.getAddress()
//       const code = `function test(){return'hello world'}`
//       const codeHex = stringToHex(code)
//       await expect(
//         contracts.client.sendSimpleRequestWithJavaScript(
//           subscriptionId,
//           code,
//           ids.donId,
//         ),
//       )
//         .to.emit(contracts.coordinator, 'OracleRequest')
//         .withArgs(
//           anyValue,
//           contracts.client.address,
//           defaultAccountAddress,
//           subscriptionId,
//           roles.subOwnerAddress,
//           codeHex,
//           anyValue,
//         )
//     })

//     it('#sendRequest reverts for empty data', async () => {
//       subscriptionId = await createSubscription(
//         roles.subOwner,
//         [contracts.client.address],
//         contracts.router,
//         contracts.linkToken,
//       )
//       await expect(
//         contracts.client.sendSimpleRequestWithJavaScript(
//           subscriptionId,
//           '',
//           ids.donId,
//         ),
//       ).to.be.revertedWith('EmptyRequestData')
//     })

//     it('#sendRequest returns non-empty requestId', async () => {
//       subscriptionId = await createSubscription(
//         roles.subOwner,
//         [contracts.client.address],
//         contracts.router,
//         contracts.linkToken,
//       )
//       const requestId = await contracts.client.sendSimpleRequestWithJavaScript(
//         subscriptionId,
//         'test',
//         ids.donId,
//       )
//       expect(requestId).not.to.be.empty
//     })

//     it('#sendRequest returns different requestIds', async () => {
//       subscriptionId = await createSubscription(
//         roles.subOwner,
//         [contracts.client.address],
//         contracts.router,
//         contracts.linkToken,
//       )
//       const defaultAccountAddress = await roles.defaultAccount.getAddress()
//       const data = stringToHex('test data')
//       const requestId1 = await contracts.client.sendSimpleRequestWithJavaScript(
//         subscriptionId,
//         'test data',
//         ids.donId,
//       )
//       await expect(
//         contracts.client.sendSimpleRequestWithJavaScript(
//           subscriptionId,
//           'test data',
//           ids.donId,
//         ),
//       )
//         .to.emit(contracts.coordinator, 'OracleRequest')
//         .withArgs(
//           anyValue,
//           contracts.client.address,
//           defaultAccountAddress,
//           subscriptionId,
//           roles.subOwnerAddress,
//           data,
//           anyValue,
//         )
//       const requestId2 = await contracts.client.sendSimpleRequestWithJavaScript(
//         subscriptionId,
//         'test data',
//         ids.donId,
//       )
//       expect(requestId1).not.to.be.equal(requestId2)
//     })
//   })

//   describe('Fulfilling requests', () => {
//     const placeTestRequest = async () => {
//       const requestId = await contracts.client
//         .connect(roles.oracleNode)
//         .callStatic.sendSimpleRequestWithJavaScript(
//           'function(){}',
//           subscriptionId,
//         )
//       await expect(
//         contracts.client
//           .connect(roles.oracleNode)
//           .sendSimpleRequestWithJavaScript('function(){}', subscriptionId),
//       )
//         .to.emit(contracts.client, 'RequestSent')
//         .withArgs(requestId)
//       return requestId
//     }

//     it('#fulfillRequest emits an error for unknown requestId', async () => {
//       const requestId =
//         '0x67c6a2e151d4352a55021b5d0028c18121cfc24c7d73b179d22b17daff069c6e'

//       const report = encodeReport(
//         ethers.utils.hexZeroPad(requestId, 32),
//         stringToHex('response'),
//         stringToHex(''),
//       )

//       await expect(contracts.coordinator.callReport(report)).to.emit(
//         contracts.coordinator,
//         'InvalidRequestID',
//       )
//     })

//     it('#fulfillRequest emits OracleResponse and ResponseTransmitted', async () => {
//       const requestId = await placeTestRequest()

//       const report = encodeReport(
//         ethers.utils.hexZeroPad(requestId, 32),
//         stringToHex('response'),
//         stringToHex(''),
//       )

//       const transmitter = await roles.oracleNode.getAddress()

//       await expect(
//         contracts.coordinator.connect(roles.oracleNode).callReport(report),
//       )
//         .to.emit(contracts.coordinator, 'OracleResponse')
//         .withArgs(requestId)
//         .to.emit(contracts.coordinator, 'ResponseTransmitted')
//         .withArgs(requestId, transmitter)
//     })

//     it('#estimateCost correctly estimates cost [ @skip-coverage ]', async () => {
//       const [subscriptionBalanceBefore] =
//         await contracts.router.getSubscription(subscriptionId)

//       const request = await contracts.client
//         .connect(roles.oracleNode)
//         .sendSimpleRequestWithJavaScript('function(){}', subscriptionId)
//       const receipt = await request.wait()
//       const requestId = receipt.events[3].args[0]

//       const report = encodeReport(
//         ethers.utils.hexZeroPad(requestId, 32),
//         stringToHex('response'),
//         stringToHex(''),
//       )

//       const transmitter = await roles.oracleNode.getAddress()

//       await expect(
//         contracts.coordinator.connect(roles.oracleNode).callReport(report),
//       )
//         .to.emit(contracts.coordinator, 'OracleResponse')
//         .withArgs(requestId)
//         .to.emit(contracts.coordinator, 'ResponseTransmitted')
//         .withArgs(requestId, transmitter)
//         .to.emit(contracts.coordinator, 'BillingEnd')

//       const [subscriptionBalanceAfter] = await contracts.router.getSubscription(
//         subscriptionId,
//       )

//       const feeData = await ethers.provider.getFeeData()
//       const estimatedCost = await contracts.client.estimateJuelCost(
//         'function(){}',
//         subscriptionId,
//         feeData.gasPrice ?? BigNumber.from(0),
//       )
//       // Expect charged amount to be +-0.01%
//       expect(
//         subscriptionBalanceBefore.sub(subscriptionBalanceAfter),
//       ).to.be.below(estimatedCost.add(estimatedCost.div(100)))
//       expect(
//         subscriptionBalanceBefore.sub(subscriptionBalanceAfter),
//       ).to.be.above(estimatedCost.sub(estimatedCost.div(100)))
//     })

//     it('#fulfillRequest emits UserCallbackError if callback reverts', async () => {
//       const requestId = await placeTestRequest()

//       const report = encodeReport(
//         ethers.utils.hexZeroPad(requestId, 32),
//         stringToHex('response'),
//         stringToHex(''),
//       )

//       const transmitter = await roles.oracleNode.getAddress()

//       await contracts.client.setRevertFulfillRequest(true)

//       await expect(
//         contracts.coordinator.connect(roles.oracleNode).callReport(report),
//       )
//         .to.emit(contracts.coordinator, 'UserCallbackError')
//         .withArgs(requestId, anyValue)
//         .to.emit(contracts.coordinator, 'ResponseTransmitted')
//         .withArgs(requestId, transmitter)
//     })

//     it('#fulfillRequest emits UserCallbackError if callback does invalid op', async () => {
//       const requestId = await placeTestRequest()

//       const report = encodeReport(
//         ethers.utils.hexZeroPad(requestId, 32),
//         stringToHex('response'),
//         stringToHex(''),
//       )

//       const transmitter = await roles.oracleNode.getAddress()

//       await contracts.client.setDoInvalidOperation(true)

//       await expect(
//         contracts.coordinator.connect(roles.oracleNode).callReport(report),
//       )
//         .to.emit(contracts.coordinator, 'UserCallbackError')
//         .withArgs(requestId, anyValue)
//         .to.emit(contracts.coordinator, 'ResponseTransmitted')
//         .withArgs(requestId, transmitter)
//     })

//     it('#fulfillRequest invokes contracts.client fulfillRequest', async () => {
//       const requestId = await placeTestRequest()

//       const report = encodeReport(
//         ethers.utils.hexZeroPad(requestId, 32),
//         stringToHex('response'),
//         stringToHex('err'),
//       )

//       await expect(
//         contracts.coordinator.connect(roles.oracleNode).callReport(report),
//       )
//         .to.emit(contracts.client, 'FulfillRequestInvoked')
//         .withArgs(requestId, stringToHex('response'), stringToHex('err'))
//     })

//     it('#fulfillRequest invalidates requestId', async () => {
//       const requestId = await placeTestRequest()

//       const report = encodeReport(
//         ethers.utils.hexZeroPad(requestId, 32),
//         stringToHex('response'),
//         stringToHex('err'),
//       )

//       await expect(
//         contracts.coordinator.connect(roles.oracleNode).callReport(report),
//       )
//         .to.emit(contracts.client, 'FulfillRequestInvoked')
//         .withArgs(requestId, stringToHex('response'), stringToHex('err'))

//       // for second fulfill the requestId becomes invalid
//       await expect(
//         contracts.coordinator.connect(roles.oracleNode).callReport(report),
//       )
//         .to.emit(contracts.coordinator, 'InvalidRequestID')
//         .withArgs(requestId)
//     })

//     it('#_report reverts for inconsistent encoding', async () => {
//       const requestId = await placeTestRequest()

//       const abi = ethers.utils.defaultAbiCoder
//       const report = abi.encode(
//         ['bytes32[]', 'bytes[]', 'bytes[]'],
//         [[requestId], [], []],
//       )

//       await expect(
//         contracts.coordinator.connect(roles.oracleNode).callReport(report),
//       ).to.be.revertedWith('ReportInvalid()')
//     })

//     it('#_report handles multiple reports', async () => {
//       const requestId1 = await placeTestRequest()
//       const requestId2 = await placeTestRequest()
//       const result1 = stringToHex('result1')
//       const result2 = stringToHex('result2')
//       const err = stringToHex('')

//       const abi = ethers.utils.defaultAbiCoder
//       const report = abi.encode(
//         ['bytes32[]', 'bytes[]', 'bytes[]'],
//         [
//           [requestId1, requestId2],
//           [result1, result2],
//           [err, err],
//         ],
//       )

//       await expect(
//         contracts.coordinator
//           .connect(roles.oracleNode)
//           .callReport(report, { gasLimit: 300_000 }),
//       )
//         .to.emit(contracts.client, 'FulfillRequestInvoked')
//         .withArgs(requestId1, result1, err)
//         .to.emit(contracts.client, 'FulfillRequestInvoked')
//         .withArgs(requestId2, result2, err)
//     })

//     it('#_report handles multiple failures', async () => {
//       const requestId1 = await placeTestRequest()
//       const requestId2 = await placeTestRequest()
//       const result1 = stringToHex('result1')
//       const result2 = stringToHex('result2')
//       const err = stringToHex('')

//       const abi = ethers.utils.defaultAbiCoder
//       const report = abi.encode(
//         ['bytes32[]', 'bytes[]', 'bytes[]'],
//         [
//           [requestId1, requestId2],
//           [result1, result2],
//           [err, err],
//         ],
//       )

//       await contracts.client.setRevertFulfillRequest(true)

//       await expect(
//         contracts.coordinator.connect(roles.oracleNode).callReport(report),
//       )
//         .to.emit(contracts.coordinator, 'UserCallbackError')
//         .withArgs(requestId1, anyValue)
//         .to.emit(contracts.coordinator, 'UserCallbackError')
//         .withArgs(requestId2, anyValue)
//     })
//   })

//   describe('#startBilling', () => {
//     let subId: number

//     beforeEach(async () => {
//       subId = await createSubscription(
//         roles.subOwner,
//         [roles.consumerAddress],
//         contracts.router,
//       )

//       await contracts.linkToken
//         .connect(roles.subOwner)
//         .transferAndCall(
//           contracts.router.address,
//           BigNumber.from('54666805176129187'),
//           ethers.utils.defaultAbiCoder.encode(['uint64'], [subId]),
//         )
//       await contracts.router
//         .connect(roles.subOwner)
//         .addConsumer(subId, contracts.client.address)
//     })

//     it('only callable by registered DONs', async () => {
//       await expect(
//         contracts.router
//           .connect(roles.consumer)
//           .startBilling(stringToHex('some data'), {
//             requester: roles.consumerAddress,
//             client: roles.consumerAddress,
//             subscriptionId: subId,
//             gasPrice: 20_000,
//             gasLimit: 20_000,
//             confirmations: 50,
//           }),
//       ).to.be.revertedWith(`reverted with custom error 'UnauthorizedSender()'`)
//     })

//     it('a subscription can only be used by a subscription consumer', async () => {
//       await expect(
//         contracts.coordinator
//           .connect(roles.stranger)
//           .sendRequest(subId, stringToBytes('some data'), 0),
//       ).to.be.revertedWith(
//         `reverted with custom error 'InvalidConsumer(${subId}, "${roles.strangerAddress}")`,
//       )
//       await expect(
//         contracts.client
//           .connect(roles.consumer)
//           .sendSimpleRequestWithJavaScript(
//             `return 'hello world'`,
//             subId,
//             donLabel,
//           ),
//       ).to.not.be.reverted
//     })

//     it('fails if the subscription does not have the funds for the estimated cost', async () => {
//       const subId = await createSubscription(
//         roles.subOwner,
//         [roles.subOwnerAddress],
//         contracts.router,
//       )
//       await contracts.router
//         .connect(roles.subOwner)
//         .addConsumer(subId, contracts.client.address)

//       await expect(
//         contracts.client
//           .connect(roles.subOwner)
//           .sendSimpleRequestWithJavaScript(
//             `return 'hello world'`,
//             subId,
//             donLabel,
//           ),
//       ).to.be.revertedWith(`InsufficientBalance()`)
//     })

//     it('when successful, emits an event', async () => {
//       await expect(
//         contracts.client
//           .connect(roles.consumer)
//           .sendSimpleRequestWithJavaScript(
//             `return 'hello world'`,
//             subId,
//             donLabel,
//           ),
//       ).to.emit(contracts.router, 'BillingStart')
//     })

//     it('fails multiple requests if the subscription does not have the funds for the estimated cost', async () => {
//       contracts.client
//         .connect(roles.consumer)
//         .sendSimpleRequestWithJavaScript(
//           `return 'hello world'`,
//           subId,
//           donLabel,
//           {
//             gasPrice: 1000000008,
//           },
//         )

//       await expect(
//         contracts.client
//           .connect(roles.subOwner)
//           .sendSimpleRequestWithJavaScript(
//             `return 'hello world'`,
//             subId,
//             donLabel,
//             {
//               gasPrice: 1000000008,
//             },
//           ),
//       ).to.be.revertedWith(`InsufficientBalance()`)
//     })
//   })

//   describe('#fulfillAndBill', () => {
//     let subId: number
//     let requestId: string

//     beforeEach(async () => {
//       subId = await createSubscription(
//         roles.subOwner,
//         [roles.consumerAddress],
//         contracts.router,
//       )

//       await contracts.linkToken
//         .connect(roles.subOwner)
//         .transferAndCall(
//           contracts.router.address,
//           BigNumber.from('1000000000000000000'),
//           ethers.utils.defaultAbiCoder.encode(['uint64'], [subId]),
//         )
//       await contracts.router
//         .connect(roles.subOwner)
//         .addConsumer(subId, contracts.client.address)
//       await contracts.router.connect(roles.defaultAccount).reg
//       await contracts.router.setAuthorizedSenders([
//         contracts.coordinator.address,
//       ])

//       const request = await contracts.client
//         .connect(roles.consumer)
//         .sendSimpleRequestWithJavaScript(
//           `return 'hello world'`,
//           subId,
//           donLabel,
//         )
//       requestId = (await request.wait()).events[3].args[0]
//     })

//     it('only callable by registered DONs', async () => {
//       const someAddress = randomAddressString()
//       const someSigners = Array(31).fill(ethers.constants.AddressZero)
//       someSigners[0] = someAddress
//       await expect(
//         contracts.router
//           .connect(roles.consumer)
//           .fulfillAndBill(
//             ethers.utils.hexZeroPad(requestId, 32),
//             stringToHex('some data'),
//             stringToHex('some data'),
//             someAddress,
//             someSigners,
//             1,
//             10,
//             0,
//           ),
//       ).to.be.revertedWith('UnauthorizedSender()')
//     })

//     it('when successful, emits an event', async () => {
//       const report = encodeReport(
//         ethers.utils.hexZeroPad(requestId, 32),
//         stringToHex('hello world'),
//         stringToHex(''),
//       )
//       await expect(
//         contracts.coordinator
//           .connect(roles.oracleNode)
//           .callReport(report, { gasLimit: 500_000 }),
//       ).to.emit(contracts.router, 'BillingEnd')
//     })

//     it('validates request ID', async () => {
//       const unknown =
//         '0x67c6a2e151d4352a55021b5d0028c18121cfc24c7d73b179d22b17eeeeeeeeee'
//       const report = encodeReport(
//         ethers.utils.hexZeroPad(unknown, 32),
//         stringToHex('hello world'),
//         stringToHex(''),
//       )
//       await expect(
//         contracts.coordinator
//           .connect(roles.oracleNode)
//           .callReport(report, { gasLimit: 500_000 }),
//       ).to.emit(contracts.coordinator, 'InvalidRequestID')
//     })

//     it('pays the transmitter the expected amount', async () => {
//       const oracleBalanceBefore = await contracts.linkToken.balanceOf(
//         await roles.oracleNode.getAddress(),
//       )
//       const [subscriptionBalanceBefore] =
//         await contracts.router.getSubscription(subId)

//       const report = encodeReport(
//         ethers.utils.hexZeroPad(requestId, 32),
//         stringToHex('hello world'),
//         stringToHex(''),
//       )

//       const transmitter = await roles.oracleNode.getAddress()

//       await expect(
//         contracts.coordinator
//           .connect(roles.oracleNode)
//           .callReport(report, { gasLimit: 500_000 }),
//       )
//         .to.emit(contracts.coordinator, 'OracleResponse')
//         .withArgs(requestId)
//         .to.emit(contracts.coordinator, 'ResponseTransmitted')
//         .withArgs(requestId, transmitter)
//         .to.emit(contracts.router, 'BillingEnd')
//         .to.emit(contracts.client, 'FulfillRequestInvoked')

//       await contracts.router
//         .connect(roles.oracleNode)
//         .oracleWithdraw(
//           await roles.oracleNode.getAddress(),
//           BigNumber.from('0'),
//         )

//       const oracleBalanceAfter = await contracts.linkToken.balanceOf(
//         await roles.oracleNode.getAddress(),
//       )
//       const [subscriptionBalanceAfter] = await contracts.router.getSubscription(
//         subId,
//       )

//       expect(subscriptionBalanceBefore.gt(subscriptionBalanceAfter)).to.be.true
//       expect(oracleBalanceAfter.gt(oracleBalanceBefore)).to.be.true
//       expect(subscriptionBalanceBefore.sub(subscriptionBalanceAfter)).to.equal(
//         oracleBalanceAfter.sub(oracleBalanceBefore),
//       )
//     })
//   })
// })
