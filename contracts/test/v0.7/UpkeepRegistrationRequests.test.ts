import { ethers } from 'hardhat'
import { assert, expect } from 'chai'
import { evmRevert } from '../test-helpers/matchers'
import { getUsers, Personas } from '../test-helpers/setup'
import { BigNumber, Signer } from 'ethers'
import { LinkToken__factory as LinkTokenFactory } from '../../typechain/factories/LinkToken__factory'
import { KeeperRegistry11__factory as KeeperRegistryFactory } from '../../typechain/factories/KeeperRegistry11__factory'
import { MockV3Aggregator__factory as MockV3AggregatorFactory } from '../../typechain/factories/MockV3Aggregator__factory'
import { UpkeepRegistrationRequests__factory as UpkeepRegistrationRequestsFactory } from '../../typechain/factories/UpkeepRegistrationRequests__factory'
import { UpkeepMock__factory as UpkeepMockFactory } from '../../typechain/factories/UpkeepMock__factory'
import { KeeperRegistry11 as KeeperRegistry } from '../../typechain/KeeperRegistry11'
import { UpkeepRegistrationRequests } from '../../typechain/UpkeepRegistrationRequests'
import { MockV3Aggregator } from '../../typechain/MockV3Aggregator'
import { LinkToken } from '../../typechain/LinkToken'
import { UpkeepMock } from '../../typechain/UpkeepMock'

let linkTokenFactory: LinkTokenFactory
let mockV3AggregatorFactory: MockV3AggregatorFactory
let keeperRegistryFactory: KeeperRegistryFactory
let upkeepRegistrationRequestsFactory: UpkeepRegistrationRequestsFactory
let upkeepMockFactory: UpkeepMockFactory

let personas: Personas

before(async () => {
  personas = (await getUsers()).personas

  linkTokenFactory = await ethers.getContractFactory('LinkToken')
  mockV3AggregatorFactory = (await ethers.getContractFactory(
    'src/v0.7/tests/MockV3Aggregator.sol:MockV3Aggregator',
  )) as unknown as MockV3AggregatorFactory
  // @ts-ignore bug in autogen file
  keeperRegistryFactory = await ethers.getContractFactory('KeeperRegistry1_1')
  upkeepRegistrationRequestsFactory = await ethers.getContractFactory(
    'UpkeepRegistrationRequests',
  )
  upkeepMockFactory = await ethers.getContractFactory('UpkeepMock')
})

const errorMsgs = {
  onlyOwner: 'revert Only callable by owner',
  onlyAdmin: 'only admin / owner can cancel',
  hashPayload: 'hash and payload do not match',
  requestNotFound: 'request not found',
}

describe('UpkeepRegistrationRequests', () => {
  const upkeepName = 'SampleUpkeep'

  const linkEth = BigNumber.from(300000000)
  const gasWei = BigNumber.from(100)
  const executeGas = BigNumber.from(100000)
  const source = BigNumber.from(100)
  const paymentPremiumPPB = BigNumber.from(250000000)
  const flatFeeMicroLink = BigNumber.from(0)

  const window_big = BigNumber.from(1000)
  const window_small = BigNumber.from(2)
  const threshold_big = BigNumber.from(1000)
  const threshold_small = BigNumber.from(5)

  const blockCountPerTurn = BigNumber.from(3)
  const emptyBytes = '0x00'
  const stalenessSeconds = BigNumber.from(43820)
  const gasCeilingMultiplier = BigNumber.from(1)
  const maxCheckGas = BigNumber.from(20000000)
  const fallbackGasPrice = BigNumber.from(200)
  const fallbackLinkPrice = BigNumber.from(200000000)
  const minLINKJuels = BigNumber.from('1000000000000000000')
  const amount = BigNumber.from('5000000000000000000')
  const amount1 = BigNumber.from('6000000000000000000')

  let owner: Signer
  let admin: Signer
  let someAddress: Signer
  let registrarOwner: Signer
  let stranger: Signer

  let linkToken: LinkToken
  let linkEthFeed: MockV3Aggregator
  let gasPriceFeed: MockV3Aggregator
  let registry: KeeperRegistry
  let mock: UpkeepMock
  let registrar: UpkeepRegistrationRequests

  beforeEach(async () => {
    owner = personas.Default
    admin = personas.Neil
    someAddress = personas.Ned
    registrarOwner = personas.Nelly
    stranger = personas.Nancy

    linkToken = await linkTokenFactory.connect(owner).deploy()
    gasPriceFeed = await mockV3AggregatorFactory
      .connect(owner)
      .deploy(0, gasWei)
    linkEthFeed = await mockV3AggregatorFactory
      .connect(owner)
      .deploy(9, linkEth)
    registry = await keeperRegistryFactory
      .connect(owner)
      .deploy(
        linkToken.address,
        linkEthFeed.address,
        gasPriceFeed.address,
        paymentPremiumPPB,
        flatFeeMicroLink,
        blockCountPerTurn,
        maxCheckGas,
        stalenessSeconds,
        gasCeilingMultiplier,
        fallbackGasPrice,
        fallbackLinkPrice,
      )

    mock = await upkeepMockFactory.deploy()

    registrar = await upkeepRegistrationRequestsFactory
      .connect(registrarOwner)
      .deploy(linkToken.address, minLINKJuels)

    await registry.setRegistrar(registrar.address)
  })

  describe('#typeAndVersion', () => {
    it('uses the correct type and version', async () => {
      const typeAndVersion = await registrar.typeAndVersion()
      assert.equal(typeAndVersion, 'UpkeepRegistrationRequests 1.0.0')
    })
  })

  describe('#register', () => {
    it('reverts if not called by the LINK token', async () => {
      await evmRevert(
        registrar
          .connect(someAddress)
          .register(
            upkeepName,
            emptyBytes,
            mock.address,
            executeGas,
            await admin.getAddress(),
            emptyBytes,
            amount,
            source,
          ),
        'Must use LINK token',
      )
    })

    it('reverts if the amount passed in data mismatches actual amount sent', async () => {
      await registrar
        .connect(registrarOwner)
        .setRegistrationConfig(
          true,
          window_small,
          threshold_big,
          registry.address,
          minLINKJuels,
        )

      const abiEncodedBytes = registrar.interface.encodeFunctionData(
        'register',
        [
          upkeepName,
          emptyBytes,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          amount1,
          source,
        ],
      )

      await evmRevert(
        linkToken.transferAndCall(registrar.address, amount, abiEncodedBytes),
        'Amount mismatch',
      )
    })

    it('reverts if the admin address is 0x0000...', async () => {
      const abiEncodedBytes = registrar.interface.encodeFunctionData(
        'register',
        [
          upkeepName,
          emptyBytes,
          mock.address,
          executeGas,
          '0x0000000000000000000000000000000000000000',
          emptyBytes,
          amount,
          source,
        ],
      )

      await evmRevert(
        linkToken.transferAndCall(registrar.address, amount, abiEncodedBytes),
        'Unable to create request',
      )
    })

    it('Auto Approve ON - registers an upkeep on KeeperRegistry instantly and emits both RegistrationRequested and RegistrationApproved events', async () => {
      //get current upkeep count
      const upkeepCount = await registry.getUpkeepCount()

      //set auto approve ON with high threshold limits
      await registrar
        .connect(registrarOwner)
        .setRegistrationConfig(
          true,
          window_small,
          threshold_big,
          registry.address,
          minLINKJuels,
        )

      //register with auto approve ON
      const abiEncodedBytes = registrar.interface.encodeFunctionData(
        'register',
        [
          upkeepName,
          emptyBytes,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          amount,
          source,
        ],
      )
      const tx = await linkToken.transferAndCall(
        registrar.address,
        amount,
        abiEncodedBytes,
      )

      //confirm if a new upkeep has been registered and the details are the same as the one just registered
      const newupkeep = await registry.getUpkeep(upkeepCount)
      assert.equal(newupkeep.target, mock.address)
      assert.equal(newupkeep.admin, await admin.getAddress())
      assert.equal(newupkeep.checkData, emptyBytes)
      assert.equal(newupkeep.balance.toString(), amount.toString())
      assert.equal(newupkeep.executeGas, executeGas.toNumber())

      await expect(tx).to.emit(registrar, 'RegistrationRequested')
      await expect(tx).to.emit(registrar, 'RegistrationApproved')
    })

    it('Auto Approve OFF - does not registers an upkeep on KeeperRegistry, emits only RegistrationRequested event', async () => {
      //get upkeep count before attempting registration
      const beforeCount = await registry.getUpkeepCount()

      //set auto approve OFF, threshold limits dont matter in this case
      await registrar
        .connect(registrarOwner)
        .setRegistrationConfig(
          false,
          window_small,
          threshold_big,
          registry.address,
          minLINKJuels,
        )

      //register with auto approve OFF
      const abiEncodedBytes = registrar.interface.encodeFunctionData(
        'register',
        [
          upkeepName,
          emptyBytes,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          amount,
          source,
        ],
      )
      const tx = await linkToken.transferAndCall(
        registrar.address,
        amount,
        abiEncodedBytes,
      )
      const receipt = await tx.wait()

      //get upkeep count after attempting registration
      const afterCount = await registry.getUpkeepCount()
      //confirm that a new upkeep has NOT been registered and upkeep count is still the same
      assert.deepEqual(beforeCount, afterCount)

      //confirm that only RegistrationRequested event is emitted and RegistrationApproved event is not
      await expect(tx).to.emit(registrar, 'RegistrationRequested')
      await expect(tx).not.to.emit(registrar, 'RegistrationApproved')

      const hash = receipt.logs[2].topics[1]
      const pendingRequest = await registrar.getPendingRequest(hash)
      assert.equal(await admin.getAddress(), pendingRequest[0])
      assert.ok(amount.eq(pendingRequest[1]))
    })

    it('Auto Approve ON - Throttle max approvals - does not registers an upkeep on KeeperRegistry beyond the throttle limit, emits only RegistrationRequested event after throttle starts', async () => {
      //get upkeep count before attempting registration
      const beforeCount = await registry.getUpkeepCount()

      //set auto approve on, with low threshold limits
      await registrar
        .connect(registrarOwner)
        .setRegistrationConfig(
          true,
          window_big,
          threshold_small,
          registry.address,
          minLINKJuels,
        )

      let abiEncodedBytes = registrar.interface.encodeFunctionData('register', [
        upkeepName,
        emptyBytes,
        mock.address,
        executeGas,
        await admin.getAddress(),
        emptyBytes,
        amount,
        source,
      ])

      //register within threshold, new upkeep should be registered
      await linkToken.transferAndCall(
        registrar.address,
        amount,
        abiEncodedBytes,
      )
      const intermediateCount = await registry.getUpkeepCount()
      //make sure 1 upkeep was registered
      assert.equal(beforeCount.toNumber() + 1, intermediateCount.toNumber())

      //try registering more than threshold(say 2x), new upkeeps should not be registered after the threshold amount is reached
      for (let step = 0; step < threshold_small.toNumber() * 2; step++) {
        abiEncodedBytes = registrar.interface.encodeFunctionData('register', [
          upkeepName,
          emptyBytes,
          mock.address,
          executeGas.toNumber() + step, // make unique hash
          await admin.getAddress(),
          emptyBytes,
          amount,
          source,
        ])

        await linkToken.transferAndCall(
          registrar.address,
          amount,
          abiEncodedBytes,
        )
      }
      const afterCount = await registry.getUpkeepCount()
      //count of newly registered upkeeps should be equal to the threshold set for auto approval
      const newRegistrationsCount =
        afterCount.toNumber() - beforeCount.toNumber()
      assert(
        newRegistrationsCount == threshold_small.toNumber(),
        'Registrations beyond threshold',
      )
    })
  })

  describe('#approve', () => {
    let hash: string

    beforeEach(async () => {
      await registrar
        .connect(registrarOwner)
        .setRegistrationConfig(
          false,
          window_small,
          threshold_big,
          registry.address,
          minLINKJuels,
        )

      //register with auto approve OFF
      const abiEncodedBytes = registrar.interface.encodeFunctionData(
        'register',
        [
          upkeepName,
          emptyBytes,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          amount,
          source,
        ],
      )

      const tx = await linkToken.transferAndCall(
        registrar.address,
        amount,
        abiEncodedBytes,
      )
      const receipt = await tx.wait()
      hash = receipt.logs[2].topics[1]
    })

    it('reverts if not called by the owner', async () => {
      const tx = registrar
        .connect(stranger)
        .approve(
          upkeepName,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          hash,
        )
      await evmRevert(tx, 'Only callable by owner')
    })

    it('reverts if the hash does not exist', async () => {
      const tx = registrar
        .connect(registrarOwner)
        .approve(
          upkeepName,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          '0x000000000000000000000000322813fd9a801c5507c9de605d63cea4f2ce6c44',
        )
      await evmRevert(tx, errorMsgs.requestNotFound)
    })

    it('reverts if any member of the payload changes', async () => {
      let tx = registrar
        .connect(registrarOwner)
        .approve(
          upkeepName,
          ethers.Wallet.createRandom().address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          hash,
        )
      await evmRevert(tx, errorMsgs.hashPayload)
      tx = registrar
        .connect(registrarOwner)
        .approve(
          upkeepName,
          mock.address,
          10000,
          await admin.getAddress(),
          emptyBytes,
          hash,
        )
      await evmRevert(tx, errorMsgs.hashPayload)
      tx = registrar
        .connect(registrarOwner)
        .approve(
          upkeepName,
          mock.address,
          executeGas,
          ethers.Wallet.createRandom().address,
          emptyBytes,
          hash,
        )
      await evmRevert(tx, errorMsgs.hashPayload)
      tx = registrar
        .connect(registrarOwner)
        .approve(
          upkeepName,
          mock.address,
          executeGas,
          await admin.getAddress(),
          '0x1234',
          hash,
        )
      await evmRevert(tx, errorMsgs.hashPayload)
    })

    it('approves an existing registration request', async () => {
      const tx = await registrar
        .connect(registrarOwner)
        .approve(
          upkeepName,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          hash,
        )
      await expect(tx).to.emit(registrar, 'RegistrationApproved')
    })

    it('deletes the request afterwards / reverts if the request DNE', async () => {
      await registrar
        .connect(registrarOwner)
        .approve(
          upkeepName,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          hash,
        )
      const tx = registrar
        .connect(registrarOwner)
        .approve(
          upkeepName,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          hash,
        )
      await evmRevert(tx, errorMsgs.requestNotFound)
    })
  })

  describe('#cancel', () => {
    let hash: string

    beforeEach(async () => {
      await registrar
        .connect(registrarOwner)
        .setRegistrationConfig(
          false,
          window_small,
          threshold_big,
          registry.address,
          minLINKJuels,
        )

      //register with auto approve OFF
      const abiEncodedBytes = registrar.interface.encodeFunctionData(
        'register',
        [
          upkeepName,
          emptyBytes,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          amount,
          source,
        ],
      )
      const tx = await linkToken.transferAndCall(
        registrar.address,
        amount,
        abiEncodedBytes,
      )
      const receipt = await tx.wait()
      hash = receipt.logs[2].topics[1]
      // submit duplicate request (increase balance)
      await linkToken.transferAndCall(
        registrar.address,
        amount,
        abiEncodedBytes,
      )
    })

    it('reverts if not called by the admin / owner', async () => {
      const tx = registrar.connect(stranger).cancel(hash)
      await evmRevert(tx, errorMsgs.onlyAdmin)
    })

    it('reverts if the hash does not exist', async () => {
      const tx = registrar
        .connect(registrarOwner)
        .cancel(
          '0x000000000000000000000000322813fd9a801c5507c9de605d63cea4f2ce6c44',
        )
      await evmRevert(tx, 'request not found')
    })

    it('refunds the total request balance to the admin address', async () => {
      const before = await linkToken.balanceOf(await admin.getAddress())
      const tx = await registrar.connect(admin).cancel(hash)
      const after = await linkToken.balanceOf(await admin.getAddress())
      assert.isTrue(after.sub(before).eq(amount.mul(BigNumber.from(2))))
      await expect(tx).to.emit(registrar, 'RegistrationRejected')
    })

    it('deletes the request hash', async () => {
      await registrar.connect(registrarOwner).cancel(hash)
      let tx = registrar.connect(registrarOwner).cancel(hash)
      await evmRevert(tx, errorMsgs.requestNotFound)
      tx = registrar
        .connect(registrarOwner)
        .approve(
          upkeepName,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          hash,
        )
      await evmRevert(tx, errorMsgs.requestNotFound)
    })
  })
})
