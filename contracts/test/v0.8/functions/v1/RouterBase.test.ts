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
  functionsRouterConfig,
} from './utils'

const setup = getSetupFactory()
let contracts: FunctionsContracts
let factories: FunctionsFactories
let roles: FunctionsRoles

beforeEach(async () => {
  ;({ contracts, factories, roles } = setup())
})

describe('FunctionsRouter - Base', () => {
  describe('Config', () => {
    it('non-owner is unable to set config', async () => {
      await expect(
        contracts.router
          .connect(roles.stranger)
          .proposeContractsUpdate([ids.donId], [contracts.coordinator.address]),
      ).to.be.revertedWith('Only callable by owner')
    })

    it('non-owner is unable to apply config', async () => {
      await expect(
        contracts.router.proposeConfigUpdate(
          ids.routerId,
          ethers.utils.defaultAbiCoder.encode(
            ['uint16', 'uint96', 'bytes4', 'uint32[]'],
            [2000, 1, 0x0ca76175, [300_000, 500_000]],
          ),
        ),
      ).to.emit(contracts.router, 'ConfigProposed')
      await expect(
        contracts.router.connect(roles.stranger).updateConfig(ids.routerId),
      ).to.be.revertedWith('Only callable by owner')
    })

    it('incorrect config cannot be applied', async () => {
      await contracts.router.proposeConfigUpdate(
        ids.routerId,
        ethers.utils.defaultAbiCoder.encode(
          ['int256', 'bytes4'],
          [-100000, 0x0ca76175],
        ),
      )
      await expect(contracts.router.updateConfig(ids.routerId)).to.be.reverted
    })

    it('Owner can update config of the Router', async () => {
      const beforeConfig = await contracts.router.getConfig()
      await expect(
        contracts.router.proposeConfigUpdateSelf(
          ethers.utils.defaultAbiCoder.encode(
            contracts.router.interface.events[
              'ConfigChanged((uint16,uint96,bytes4,uint32[]))'
            ].inputs,
            [
              {
                maxConsumersPerSubscription: 2000,
                adminFee: 1,
                handleOracleFulfillmentSelector: 0x0ca76175,
                maxCallbackGasLimits: [300_000, 500_000],
              },
            ],
          ),
        ),
      ).to.emit(contracts.router, 'ConfigProposed')
      await expect(contracts.router.updateConfigSelf()).to.emit(
        contracts.router,
        'ConfigUpdated',
      )

      const afterConfig = await contracts.router.getConfig()
      expect(beforeConfig).to.not.equal(afterConfig)
    })

    it('Config of a contract on a route can be updated', async () => {
      const beforeConfig = await contracts.coordinator.getConfig()

      await expect(
        contracts.router.proposeConfigUpdate(
          ids.donId,
          ethers.utils.defaultAbiCoder.encode(
            [
              'uint32',
              'uint32',
              'uint32',
              'uint32',
              'int256',
              'uint32',
              'uint96',
              'uint16',
              'uint256',
            ],
            [
              ...Object.values({
                ...coordinatorConfig,
                maxSupportedRequestDataVersion: 2,
              }),
            ],
          ),
        ),
      ).to.emit(contracts.router, 'ConfigProposed')
      await expect(contracts.router.updateConfig(ids.donId)).to.emit(
        contracts.router,
        'ConfigUpdated',
      )

      const afterConfig = await contracts.router.getConfig()
      expect(beforeConfig).to.not.equal(afterConfig)
    })

    it('returns the config set on the Router', async () => {
      const config = await contracts.router.connect(roles.stranger).getConfig()
      expect(config[0]).to.equal(
        functionsRouterConfig.maxConsumersPerSubscription,
      )
      expect(config[1]).to.equal(functionsRouterConfig.adminFee)
      expect(config[2]).to.equal(
        functionsRouterConfig.handleOracleFulfillmentSelector,
      )
      expect(config[3].toString()).to.equal(
        functionsRouterConfig.maxCallbackGasLimits.toString(),
      )
    })
  })

  describe('Updates', () => {
    it('One or more contracts on a route can be updated', async () => {
      const subscriptionId = await createSubscription(
        roles.subOwner,
        [contracts.client.address],
        contracts.router,
        contracts.accessControl,
        contracts.linkToken,
      )
      const coordinator2 = await factories.functionsCoordinatorFactory
        .connect(roles.defaultAccount)
        .deploy(
          contracts.router.address,
          ethers.utils.defaultAbiCoder.encode(
            [
              'uint32',
              'uint32',
              'uint32',
              'uint32',
              'int256',
              'uint32',
              'uint96',
              'uint16',
              'uint256',
            ],
            [...Object.values(coordinatorConfig)],
          ),
          contracts.mockLinkEth.address,
        )
      const coordinator3 = await factories.functionsCoordinatorFactory
        .connect(roles.defaultAccount)
        .deploy(
          contracts.router.address,
          ethers.utils.defaultAbiCoder.encode(
            [
              'uint32',
              'uint32',
              'uint32',
              'uint32',
              'int256',
              'uint32',
              'uint96',
              'uint16',
              'uint256',
            ],
            [...Object.values(coordinatorConfig)],
          ),
          contracts.mockLinkEth.address,
        )
      const coordinator4 = await factories.functionsCoordinatorFactory
        .connect(roles.defaultAccount)
        .deploy(
          contracts.router.address,
          ethers.utils.defaultAbiCoder.encode(
            [
              'uint32',
              'uint32',
              'uint32',
              'uint32',
              'int256',
              'uint32',
              'uint96',
              'uint16',
              'uint256',
            ],
            [...Object.values(coordinatorConfig)],
          ),
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
      const report = encodeReport(
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

  describe('Timelock', () => {
    it('prevents applying timelock updates', async () => {
      await contracts.router.proposeTimelockBlocks(5)
      await contracts.router.updateTimelockBlocks()
      await contracts.router.proposeTimelockBlocks(6)
      await expect(contracts.router.updateTimelockBlocks()).to.be.revertedWith(
        'TimelockInEffect',
      )
    })

    it('prevents applying config updates', async () => {
      await contracts.router.proposeTimelockBlocks(5)
      await contracts.router.updateTimelockBlocks()

      await contracts.router.proposeConfigUpdate(
        ids.routerId,
        ethers.utils.defaultAbiCoder.encode(
          ['uint16', 'uint96', 'bytes4', 'uint32[]'],
          [2000, 1, 0x0ca76175, [300_000, 500_000]],
        ),
      )
      await expect(
        contracts.router.updateConfig(ids.routerId),
      ).to.be.revertedWith('TimelockInEffect')
    })

    it('prevents applying contract updates', async () => {
      await contracts.router.proposeTimelockBlocks(5)
      await contracts.router.updateTimelockBlocks()

      const coordinator2 = await factories.functionsCoordinatorFactory
        .connect(roles.defaultAccount)
        .deploy(
          contracts.router.address,
          ethers.utils.defaultAbiCoder.encode(
            [
              'uint32',
              'uint32',
              'uint32',
              'uint32',
              'int256',
              'uint32',
              'uint96',
              'uint16',
              'uint256',
            ],
            [...Object.values(coordinatorConfig)],
          ),
          contracts.mockLinkEth.address,
        )

      await contracts.router.proposeContractsUpdate(
        [ids.donId2],
        [coordinator2.address],
      )

      await expect(contracts.router.updateContracts()).to.be.revertedWith(
        'TimelockInEffect',
      )
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
