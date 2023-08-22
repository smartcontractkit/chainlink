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
            ['uint96', 'bytes4'],
            [1, 0x0ca76175],
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
      const [
        beforeSystemVersionMajor,
        beforeSystemVersionMinor,
        beforeSystemVersionPatch,
      ] = await contracts.router.version()
      const beforeConfigHash = await contracts.router.getConfigHash()

      await expect(
        contracts.router.proposeConfigUpdate(
          ids.routerId,
          ethers.utils.defaultAbiCoder.encode(
            ['uint96', 'bytes4'],
            [1, 0x0ca76175],
          ),
        ),
      ).to.emit(contracts.router, 'ConfigProposed')
      await expect(contracts.router.updateConfig(ids.routerId)).to.emit(
        contracts.router,
        'ConfigUpdated',
      )
      const [
        afterSystemVersionMajor,
        afterSystemVersionMinor,
        afterSystemVersionPatch,
      ] = await contracts.router.version()
      const afterConfigHash = await contracts.router.getConfigHash()
      expect(afterSystemVersionMajor).to.equal(beforeSystemVersionMajor)
      expect(afterSystemVersionMinor).to.equal(beforeSystemVersionMinor)
      expect(afterSystemVersionPatch).to.equal(beforeSystemVersionPatch + 1)
      expect(beforeConfigHash).to.not.equal(afterConfigHash)
    })

    it('Config of a contract on a route can be updated', async () => {
      const [
        beforeSystemVersionMajor,
        beforeSystemVersionMinor,
        beforeSystemVersionPatch,
      ] = await contracts.router.version()
      const beforeConfigHash = await contracts.coordinator.getConfigHash()

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
      const [
        afterSystemVersionMajor,
        afterSystemVersionMinor,
        afterSystemVersionPatch,
      ] = await contracts.router.version()
      const afterConfigHash = await contracts.router.getConfigHash()
      expect(afterSystemVersionMajor).to.equal(beforeSystemVersionMajor)
      expect(afterSystemVersionMinor).to.equal(beforeSystemVersionMinor)
      expect(afterSystemVersionPatch).to.equal(beforeSystemVersionPatch + 1)
      expect(beforeConfigHash).to.not.equal(afterConfigHash)
    })

    it('returns the config set on the Router', async () => {
      expect(
        await contracts.router.connect(roles.stranger).getAdminFee(),
      ).to.equal(0)
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
            ],
            [...Object.values(coordinatorConfig)],
          ),
          contracts.mockLinkEth.address,
        )

      const [
        beforeSystemVersionMajor,
        beforeSystemVersionMinor,
        beforeSystemVersionPatch,
      ] = await contracts.router.version()
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

      const report = encodeReport(
        ethers.utils.hexZeroPad(requestId, 32),
        stringToHex('hello world'),
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

      const [
        afterSystemVersionMajor,
        afterSystemVersionMinor,
        afterSystemVersionPatch,
      ] = await contracts.router.version()
      expect(afterSystemVersionMajor).to.equal(beforeSystemVersionMajor)
      expect(afterSystemVersionMinor).to.equal(beforeSystemVersionMinor + 1)
      expect(afterSystemVersionPatch).to.equal(beforeSystemVersionPatch)
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
          ['uint96', 'bytes4'],
          [20, 0x0ca76175],
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
      const paused = await contracts.router.isPaused()
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

      await contracts.router.togglePaused()

      await expect(
        contracts.client.sendSimpleRequestWithJavaScript(
          `return 'hello world'`,
          subscriptionId,
          ids.donId,
        ),
      ).to.be.revertedWith('Pausable: paused')
    })
  })
})
