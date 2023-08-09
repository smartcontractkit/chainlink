import { ethers } from 'hardhat'
import { expect } from 'chai'
import {
  getSetupFactory,
  coordinatorConfig,
  FunctionsContracts,
  FunctionsFactories,
  FunctionsRoles,
  ids,
  createSubscription,
  encodeReport,
  stringToHex,
  getEventArg,
} from './utils'

const setup = getSetupFactory()
let contracts: FunctionsContracts
let factories: FunctionsFactories
let roles: FunctionsRoles

beforeEach(async () => {
  ;({ contracts, factories, roles } = setup())
})

describe('FunctionsRouter - Base', () => {
  describe('Updates', () => {
    it('One or more contracts on a route can be updated by the owner', async () => {
      const coordinator2 = await factories.functionsCoordinatorFactory
        .connect(roles.defaultAccount)
        .deploy(
          contracts.router.address,
          coordinatorConfig,
          contracts.mockLinkEth.address,
        )
      const coordinator3 = await factories.functionsCoordinatorFactory
        .connect(roles.defaultAccount)
        .deploy(
          contracts.router.address,
          coordinatorConfig,
          contracts.mockLinkEth.address,
        )
      const coordinator4 = await factories.functionsCoordinatorFactory
        .connect(roles.defaultAccount)
        .deploy(
          contracts.router.address,
          coordinatorConfig,
          contracts.mockLinkEth.address,
        )

      await expect(
        contracts.router['getContractById(bytes32)'](ids.donId2),
      ).to.be.revertedWith('RouteNotFound')
      await expect(
        contracts.router['getContractById(bytes32)'](ids.donId3),
      ).to.be.revertedWith('RouteNotFound')
      await expect(
        contracts.router['getContractById(bytes32)'](ids.donId4),
      ).to.be.revertedWith('RouteNotFound')
      await expect(
        contracts.router.proposeContractsUpdate(
          [ids.donId2, ids.donId3, ids.donId4],
          [coordinator2.address, coordinator3.address, coordinator4.address],
        ),
      ).to.emit(contracts.router, `ContractProposed`)

      const subscriptionId = await createSubscription(
        roles.subOwner,
        [contracts.client.address],
        contracts.router,
        contracts.accessControl,
        contracts.linkToken,
      )

      const requestProposedTx = await contracts.client.sendRequestProposed(
        `return 'hello world'`,
        subscriptionId,
        ids.donId2,
      )

      const { events } = await requestProposedTx.wait()
      const requestId = getEventArg(events, 'RequestSent', 0)

      const oracleRequestEvent = await coordinator2.queryFilter(
        contracts.coordinator.filters.OracleRequest(),
      )
      const onchainMetadata = oracleRequestEvent[0].args?.['commitment']
      const report = await encodeReport(
        ethers.utils.hexZeroPad(requestId, 32),
        stringToHex('hello world'),
        stringToHex(''),
        onchainMetadata,
        stringToHex(''),
      )

      await expect(
        coordinator2
          .connect(roles.oracleNode)
          .callReport(report, { gasLimit: 500_000 }),
      ).to.emit(contracts.client, 'FulfillRequestInvoked')

      await expect(contracts.router.updateContracts()).to.emit(
        contracts.router,
        'ContractUpdated',
      )
      expect(
        await contracts.router['getContractById(bytes32)'](ids.donId2),
      ).to.equal(coordinator2.address)
      expect(
        await contracts.router['getContractById(bytes32)'](ids.donId3),
      ).to.equal(coordinator3.address)
      expect(
        await contracts.router['getContractById(bytes32)'](ids.donId4),
      ).to.equal(coordinator4.address)
    })

    it('non-owner is unable to propose contract updates', async () => {
      await expect(
        contracts.router
          .connect(roles.stranger)
          .proposeContractsUpdate([ids.donId], [contracts.coordinator.address]),
      ).to.be.revertedWith('Only callable by owner')
    })

    it('non-owner is unable to apply contract updates', async () => {
      await expect(
        contracts.router.connect(roles.stranger).updateContracts(),
      ).to.be.revertedWith('Only callable by owner')
    })
  })

  describe('Emergency Pause', () => {
    it('has paused state visible', async () => {
      const paused = await contracts.router.paused()
      expect(paused).to.equal(false)
    })
    it('can pause the system', async () => {
      const subscriptionId = await createSubscription(
        roles.subOwner,
        [contracts.client.address],
        contracts.router,
        contracts.accessControl,
        contracts.linkToken,
      )

      await contracts.router.pause()

      await expect(
        contracts.client.sendSimpleRequestWithJavaScript(
          `return 'hello world'`,
          subscriptionId,
          ids.donId,
          20_000,
        ),
      ).to.be.revertedWith('Pausable: paused')
    })
  })
})
