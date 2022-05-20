import { ethers } from 'hardhat'
import {BigNumber} from "ethers";
import moment from "moment";
import {assert, expect} from "chai";
import { CanaryUpkeep } from '../../typechain/CanaryUpkeep'
import KeeperRegistryAbi from '../../abi/src/v0.8/KeeperRegistry.sol/KeeperRegistry.json'
import { fastForward, reset } from '../test-helpers/helpers'
import {deployMockContract, MockContract} from '@ethereum-waffle/mock-contract'
import {SignerWithAddress} from "@nomiclabs/hardhat-ethers/signers"
import {evmRevert} from "../test-helpers/matchers";

let canaryUpkeep: CanaryUpkeep
let nelly: SignerWithAddress
let nancy: SignerWithAddress
let neo: SignerWithAddress
let keeperAddresses: string[]
let emptyKeeperAddresses: string[]
let mockContract: MockContract

describe.only('CanaryUpkeep', () => {
    beforeEach(async () => {
        const accounts = await ethers.getSigners()
        const admin = accounts[1]
        // initialize fake node operators
        nelly = accounts[2]
        nancy = accounts[3]
        neo = accounts[4]
        mockContract = await deployMockContract(admin as any, KeeperRegistryAbi);

        // this object is the return type of getState function of KeeperRegistry. we still need to investigate how
        // to create an object for the return type. for the purpose of this test, only the address array is needed.
        // hence, the first two sub structs are omitted.
        keeperAddresses = [nelly.address, nancy.address, neo.address]
        const getStateReturn = [{}, {}, keeperAddresses];
        await mockContract.mock.getState.returns(...getStateReturn)
        const canaryUpkeepFactory = await ethers.getContractFactory(
            'CanaryUpkeep'
        )
        canaryUpkeep = await canaryUpkeepFactory.deploy(ethers.constants.AddressZero)
        await canaryUpkeep.deployed()
    })

    afterEach(async () => {
        await reset()
    })

    describe('checkUpkeep()', () => {
        it('returns true when sufficient time passes', async () => {
            await fastForward(moment.duration(6, 'minutes').asSeconds())
            expect(await canaryUpkeep.checkUpkeep('0x')).to.be.true
        })

        it('returns false when insufficient time passes', async () => {
            await fastForward(moment.duration(2, 'minutes').asSeconds())
            const [ needsUpkeep ] = await canaryUpkeep.checkUpkeep('0x')
            assert.isFalse(needsUpkeep)
        })

        it('returns false when keeper array is empty', async () => {
            await mockContract.mock.getState.returns({}, {}, emptyKeeperAddresses)
            await fastForward(moment.duration(6, 'minutes').asSeconds())
            const [ needsUpkeep ] = await canaryUpkeep.checkUpkeep('0x')
            assert.isFalse(needsUpkeep)
        })
    })

    describe('performUpkeep()', () => {
        it('enforces that transaction origin is the anticipated keeper', async () => {
            const oldTimestamp = await canaryUpkeep.connect(keeperAddresses[0]).getTimestamp()
            await fastForward(moment.duration(6, 'minutes').asSeconds())
            await canaryUpkeep.connect(keeperAddresses[0]).performUpkeep('0x')

            const keeperIndex = await canaryUpkeep.connect(keeperAddresses[0]).getKeeperIndex()
            assert.equal(keeperIndex, BigNumber.from(1), "keeper index needs to increment by 1 after performUpkeep")

            const newTimestamp = await canaryUpkeep.connect(keeperAddresses[0]).getTimestamp()
            const interval = await canaryUpkeep.connect(keeperAddresses[0]).getInterval()
            assert.isAtLeast(newTimestamp.toNumber() - oldTimestamp.toNumber(), interval.toNumber(), "timestamp needs to be updated after performUpkeep")
        })

        it('reverts if the keeper array is empty', async () => {
            await mockContract.mock.getState.returns({}, {}, emptyKeeperAddresses)
            // await evmRevert(
            //     canaryUpkeep.connect(keeperAddresses[0]).performUpkeep('0x'),
            //     "no keeper nodes exists");
            await expect(canaryUpkeep.connect(keeperAddresses[0]).performUpkeep('0x')).to.be.revertedWith(
                "no keeper nodes exists"
            )
        })

        it('reverts if not enough time has passed', async () => {
            await fastForward(moment.duration(3, 'minutes').asSeconds())
            await evmRevert(
                canaryUpkeep.connect(keeperAddresses[0]).performUpkeep('0x'),
                "Not enough time has passed after the previous upkeep");
        })

        it('reverts if an incorrect keeper tries to perform upkeep', async () => {
            await fastForward(moment.duration(6, 'minutes').asSeconds())
            await evmRevert(
                canaryUpkeep.connect(keeperAddresses[1]).performUpkeep('0x'),
                "transaction origin is not the anticipated keeper.");
        })
    })
})