import { ethers } from 'hardhat'
import { BigNumber } from 'ethers'
import { expect, assert } from 'chai'
import { getUsers } from '../../test-helpers/setup'
import { randomAddress, toWei } from '../../test-helpers/helpers'
import { deployRegistry21 } from './helpers'

// don't run these tests in CI
const describeMaybe = process.env.CI ? describe.skip : describe

// registry settings
const f = 1
const linkEth = BigNumber.from(300000000)
const gasWei = BigNumber.from(100)
const minUpkeepSpend = BigNumber.from('1000000000000000000')
const paymentPremiumPPB = BigNumber.from(250000000)
const flatFeeMicroLink = BigNumber.from(0)
const blockCountPerTurn = 20
const checkGasLimit = BigNumber.from(20000000)
const fallbackGasPrice = BigNumber.from(200)
const fallbackLinkPrice = BigNumber.from(200000000)
const maxCheckDataSize = BigNumber.from(10000)
const maxPerformDataSize = BigNumber.from(10000)
const maxRevertDataSize = BigNumber.from(1000)
const maxPerformGas = BigNumber.from(5000000)
const stalenessSeconds = BigNumber.from(43820)
const gasCeilingMultiplier = BigNumber.from(1)
const signers = [
  randomAddress(),
  randomAddress(),
  randomAddress(),
  randomAddress(),
]
const transmitters = [
  randomAddress(),
  randomAddress(),
  randomAddress(),
  randomAddress(),
]
const transcoder = ethers.constants.AddressZero

// registrar settings
const triggerType = 0 // conditional
const autoApproveType = 2 // auto-approve enabled
const autoApproveMaxAllowed = 100 // auto-approve enabled

// upkeep settings
const name = 'test upkeep'
const encryptedEmail = '0xabcd1234'
const gasLimit = 100_000
const checkData = '0xdeadbeef'
const amount = toWei('5')
const source = 5
const triggerConfig = '0x'
const offchainConfig = '0x'

describeMaybe('Automation Gas Analysis', () => {
  it('Compares gas usage amongst registries / registrars', async () => {
    assert(
      Boolean(process.env.REPORT_GAS),
      'this test must be run with REPORT_GAS=true',
    )

    const personas = (await getUsers()).personas
    const owner = personas.Default
    const ownerAddress = await owner.getAddress()

    // factories
    const getFact = ethers.getContractFactory
    const linkTokenFactory = await getFact('LinkToken')
    const mockV3AggregatorFactory = await getFact(
      'src/v0.8/tests/MockV3Aggregator.sol:MockV3Aggregator',
    )
    const upkeepMockFactory = await getFact('UpkeepMock')
    const registry12Factory = await getFact('KeeperRegistry1_2')
    const registrar12Factory = await getFact('KeeperRegistrar')
    const registry20Factory = await getFact('KeeperRegistry2_0')
    const registryLogic20Factory = await getFact('KeeperRegistryLogic2_0')
    const registrar20Factory = await getFact('KeeperRegistrar2_0')
    const registrar21Factory = await getFact('AutomationRegistrar2_1')
    const forwarderLogicFactory = await getFact('AutomationForwarderLogic')

    // deploy dependancy contracts
    const linkToken = await linkTokenFactory.connect(owner).deploy()
    const gasPriceFeed = await mockV3AggregatorFactory
      .connect(owner)
      .deploy(0, gasWei)
    const linkEthFeed = await mockV3AggregatorFactory
      .connect(owner)
      .deploy(9, linkEth)
    const upkeep = await upkeepMockFactory.connect(owner).deploy()

    // deploy v1.2
    const registrar12 = await registrar12Factory.connect(owner).deploy(
      linkToken.address,
      autoApproveType,
      autoApproveMaxAllowed,
      ethers.constants.AddressZero, // set later
      minUpkeepSpend,
    )
    const registry12 = await registry12Factory
      .connect(owner)
      .deploy(linkToken.address, linkEthFeed.address, gasPriceFeed.address, {
        paymentPremiumPPB,
        flatFeeMicroLink,
        blockCountPerTurn,
        checkGasLimit,
        stalenessSeconds,
        gasCeilingMultiplier,
        minUpkeepSpend,
        maxPerformGas,
        fallbackGasPrice,
        fallbackLinkPrice,
        transcoder,
        registrar: registrar12.address,
      })
    await registrar12.setRegistrationConfig(
      autoApproveType,
      autoApproveMaxAllowed,
      registry12.address,
      minUpkeepSpend,
    )

    // deploy v2.0
    const registryLogic20 = await registryLogic20Factory
      .connect(owner)
      .deploy(0, linkToken.address, linkEthFeed.address, gasPriceFeed.address)
    const registry20 = await registry20Factory
      .connect(owner)
      .deploy(registryLogic20.address)
    const registrar20 = await registrar20Factory
      .connect(owner)
      .deploy(
        linkToken.address,
        autoApproveType,
        autoApproveMaxAllowed,
        registry20.address,
        minUpkeepSpend,
      )
    const config20 = {
      paymentPremiumPPB,
      flatFeeMicroLink,
      checkGasLimit,
      stalenessSeconds,
      gasCeilingMultiplier,
      minUpkeepSpend,
      maxCheckDataSize,
      maxPerformDataSize,
      maxPerformGas,
      fallbackGasPrice,
      fallbackLinkPrice,
      transcoder,
      registrar: registrar20.address,
    }
    const onchainConfig20 = ethers.utils.defaultAbiCoder.encode(
      [
        'tuple(uint32 paymentPremiumPPB,uint32 flatFeeMicroLink,uint32 checkGasLimit,uint24 stalenessSeconds\
            ,uint16 gasCeilingMultiplier,uint96 minUpkeepSpend,uint32 maxPerformGas,uint32 maxCheckDataSize,\
            uint32 maxPerformDataSize,uint256 fallbackGasPrice,uint256 fallbackLinkPrice,address transcoder,\
            address registrar)',
      ],
      [config20],
    )
    await registry20
      .connect(owner)
      .setConfig(signers, transmitters, f, onchainConfig20, 1, '0x')

    // deploy v2.1
    const forwarderLogic = await forwarderLogicFactory.connect(owner).deploy()
    const registry21 = await deployRegistry21(
      owner,
      0,
      linkToken.address,
      linkEthFeed.address,
      gasPriceFeed.address,
      forwarderLogic.address,
    )
    const registrar21 = await registrar21Factory
      .connect(owner)
      .deploy(linkToken.address, registry21.address, minUpkeepSpend, [
        {
          triggerType,
          autoApproveType,
          autoApproveMaxAllowed,
        },
      ])
    const onchainConfig21 = {
      paymentPremiumPPB,
      flatFeeMicroLink,
      checkGasLimit,
      stalenessSeconds,
      gasCeilingMultiplier,
      minUpkeepSpend,
      maxCheckDataSize,
      maxPerformDataSize,
      maxRevertDataSize,
      maxPerformGas,
      fallbackGasPrice,
      fallbackLinkPrice,
      transcoder,
      registrars: [registrar21.address],
      upkeepPrivilegeManager: randomAddress(),
    }
    await registry21
      .connect(owner)
      .setConfigTypeSafe(signers, transmitters, f, onchainConfig21, 1, '0x')

    // approve LINK
    await linkToken.connect(owner).approve(registrar20.address, amount)
    await linkToken.connect(owner).approve(registrar21.address, amount)

    const abiEncodedBytes = registrar12.interface.encodeFunctionData(
      'register',
      [
        name,
        encryptedEmail,
        upkeep.address,
        gasLimit,
        ownerAddress,
        checkData,
        amount,
        source,
        ownerAddress,
      ],
    )

    let tx = await linkToken
      .connect(owner)
      .transferAndCall(registrar12.address, amount, abiEncodedBytes)
    await expect(tx).to.emit(registry12, 'UpkeepRegistered')

    tx = await registrar20.connect(owner).registerUpkeep({
      name,
      encryptedEmail,
      upkeepContract: upkeep.address,
      gasLimit,
      adminAddress: ownerAddress,
      checkData,
      amount,
      offchainConfig,
    })
    await expect(tx).to.emit(registry20, 'UpkeepRegistered')

    tx = await registrar21.connect(owner).registerUpkeep({
      name,
      encryptedEmail,
      upkeepContract: upkeep.address,
      gasLimit,
      adminAddress: ownerAddress,
      triggerType,
      checkData,
      amount,
      triggerConfig,
      offchainConfig,
    })
    await expect(tx).to.emit(registry21, 'UpkeepRegistered')
  })
})
