import { contract, setup, helpers, matchers } from '@chainlink/test-helpers'
import { assert } from 'chai'
import { ContractTransaction } from 'ethers'
import { VRFD20Factory } from '../../ethers/v0.6/VRFD20Factory'
import { VRFCoordinatorMockFactory } from '../../ethers/v0.6/VRFCoordinatorMockFactory'
import { bigNumberify } from 'ethers/utils'

let roles: setup.Roles
let personas: setup.Personas
const provider = setup.provider()
const linkTokenFactory = new contract.LinkTokenFactory()
const vrfCoordinatorMockFactory = new VRFCoordinatorMockFactory()
const vrfD20Factory = new VRFD20Factory()

beforeAll(async () => {
  const users = await setup.users(provider)

  roles = users.roles
  personas = users.personas
})

describe('VRFD20', () => {
  const deposit = helpers.toWei('1')
  const fee = helpers.toWei('0.1')
  const keyHash = helpers.toBytes32String('keyHash')
  const seed = 12345

  const requestId =
    '0x66f86cab16b057baa86d6171b59e4c356197fcebc0e2cd2a744fc2d2f4dacbfe'

  let link: contract.Instance<contract.LinkTokenFactory>
  let vrfCoordinator: contract.Instance<VRFCoordinatorMockFactory>
  let vrfD20: contract.Instance<VRFD20Factory>

  const deployment = setup.snapshot(provider, async () => {
    link = await linkTokenFactory.connect(roles.defaultAccount).deploy()
    vrfCoordinator = await vrfCoordinatorMockFactory
      .connect(roles.defaultAccount)
      .deploy(link.address)
    vrfD20 = await vrfD20Factory
      .connect(roles.defaultAccount)
      .deploy(vrfCoordinator.address, link.address, keyHash, fee)
    await link.transfer(vrfD20.address, deposit)
  })

  beforeEach(async () => {
    await deployment()
  })

  it('has a limited public interface', () => {
    matchers.publicAbi(vrfD20Factory, [
      // Owned
      'acceptOwnership',
      'owner',
      'transferOwnership',
      //VRFConsumerBase
      'nonces',
      'rawFulfillRandomness',
      // VRFD20
      'rollDice',
      'house',
      'withdrawLINK',
      'keyHash',
      'fee',
      'setKeyHash',
      'setFee',
    ])
  })

  describe('#withdrawLINK', () => {
    describe('failure', () => {
      it('reverts when called by a non-owner', async () => {
        await matchers.evmRevert(async () => {
          await vrfD20
            .connect(roles.stranger)
            .withdrawLINK(roles.stranger.address, deposit),
            'Only callable by owner'
        })
      })

      it('reverts when not enough LINK in the contract', async () => {
        const withdrawAmount = deposit.mul(2)
        await matchers.evmRevert(async () => {
          await vrfD20
            .connect(roles.defaultAccount)
            .withdrawLINK(roles.defaultAccount.address, withdrawAmount),
            'Not enough LINK'
        })
      })
    })

    describe('success', () => {
      it('withdraws LINK', async () => {
        const startingAmount = await link.balanceOf(
          roles.defaultAccount.address,
        )
        const expectedAmount = bigNumberify(startingAmount).add(deposit)
        await vrfD20
          .connect(roles.defaultAccount)
          .withdrawLINK(roles.defaultAccount.address, deposit)
        const actualAmount = await link.balanceOf(roles.defaultAccount.address)
        assert.equal(actualAmount.toString(), expectedAmount.toString())
      })
    })
  })

  describe('#setKeyHash', () => {
    const newHash = helpers.toBytes32String('newhash')

    describe('failure', () => {
      it('reverts when called by a non-owner', async () => {
        await matchers.evmRevert(async () => {
          await vrfD20.connect(roles.stranger).setKeyHash(newHash),
            'Only callable by owner'
        })
      })
    })

    describe('success', () => {
      it('sets the key hash', async () => {
        await vrfD20.setKeyHash(newHash)
        const actualHash = await vrfD20.keyHash()
        assert.equal(actualHash, newHash)
      })
    })
  })

  describe('#setFee', () => {
    const newFee = 1234

    describe('failure', () => {
      it('reverts when called by a non-owner', async () => {
        await matchers.evmRevert(async () => {
          await vrfD20.connect(roles.stranger).setFee(newFee),
            'Only callable by owner'
        })
      })
    })

    describe('success', () => {
      it('sets the fee', async () => {
        await vrfD20.setFee(newFee)
        const actualFee = await vrfD20.fee()
        assert.equal(actualFee.toString(), newFee.toString())
      })
    })
  })

  describe('#hosue', () => {
    describe('failure', () => {
      it('reverts when dice not rolled', async () => {
        await matchers.evmRevert(async () => {
          await vrfD20.house(personas.Nancy.address), 'Dice not rolled'
        })
      })

      it('reverts when dice roll is in progress', async () => {
        await vrfD20.rollDice(seed, personas.Nancy.address)
        await matchers.evmRevert(async () => {
          await vrfD20.house(personas.Nancy.address), 'Roll in progress'
        })
      })
    })

    describe('success', () => {
      it('returns the correct house', async () => {
        const randomness = 98765
        const expectedHouse = 'Martell'
        const tx = await vrfD20.rollDice(seed, personas.Nancy.address)
        const log = await helpers.getLog(tx, 3)
        const eventRequestId = log?.topics?.[1]
        await vrfCoordinator.callBackWithRandomness(
          eventRequestId,
          randomness,
          vrfD20.address,
        )
        const response = await vrfD20.house(personas.Nancy.address)
        assert.equal(response.toString(), expectedHouse)
      })
    })
  })

  describe('#rollDice', () => {
    describe('success', () => {
      let tx: ContractTransaction
      beforeEach(async () => {
        tx = await vrfD20.rollDice(seed, personas.Nancy.address)
      })

      it('emits a RandomnessRequest event from the VRFCoordinator', async () => {
        const log = await helpers.getLog(tx, 2)
        const topics = log?.topics
        assert.equal(helpers.evmWordToAddress(topics?.[1]), vrfD20.address)
        assert.equal(topics?.[2], keyHash)
        assert.equal(topics?.[3], helpers.numToBytes32(seed))
      })
    })

    describe('failure', () => {
      it('reverts when LINK balance is zero', async () => {
        const vrfD202 = await vrfD20Factory
          .connect(roles.defaultAccount)
          .deploy(vrfCoordinator.address, link.address, keyHash, fee)
        await matchers.evmRevert(async () => {
          await vrfD202.rollDice(seed, personas.Nancy.address),
            'Not enough LINK to pay fee'
        })
      })

      it('reverts when called by a non-owner', async () => {
        await matchers.evmRevert(async () => {
          await vrfD20
            .connect(roles.stranger)
            .rollDice(seed, personas.Nancy.address),
            'Only callable by owner'
        })
      })

      it('reverts when the roller rolls more than once', async () => {
        await vrfD20.rollDice(seed, personas.Nancy.address)
        await matchers.evmRevert(async () => {
          await vrfD20.rollDice(seed, personas.Nancy.address), 'Already rolled'
        })
      })
    })
  })

  describe('#fulfillRandomness', () => {
    const randomness = 98765
    const expectedModResult = (randomness % 20) + 1
    const expectedHouse = 'Martell'
    let eventRequestId: string
    beforeEach(async () => {
      const tx = await vrfD20.rollDice(seed, personas.Nancy.address)
      const log = await helpers.getLog(tx, 3)
      eventRequestId = log?.topics?.[1]
    })

    describe('success', () => {
      let tx: ContractTransaction
      beforeEach(async () => {
        tx = await vrfCoordinator.callBackWithRandomness(
          eventRequestId,
          randomness,
          vrfD20.address,
        )
      })

      it('emits a DiceLanded event', async () => {
        const log = await helpers.getLog(tx, 0)
        assert.equal(log?.topics[1], requestId)
        assert.equal(log?.topics[2], helpers.numToBytes32(expectedModResult))
      })

      it('sets the correct dice roll result', async () => {
        const response = await vrfD20.house(personas.Nancy.address)
        assert.equal(response.toString(), expectedHouse)
      })

      it('allows someone else to roll', async () => {
        const secondRandomness = 55555
        const secondSeed = 54321
        tx = await vrfD20.rollDice(secondSeed, personas.Ned.address)
        const log = await helpers.getLog(tx, 3)
        eventRequestId = log?.topics?.[1]
        tx = await vrfCoordinator.callBackWithRandomness(
          eventRequestId,
          secondRandomness,
          vrfD20.address,
        )
      })
    })

    describe('failure', () => {
      it('does not fulfill when fulfilled by the wrong VRFcoordinator', async () => {
        const vrfCoordinator2 = await vrfCoordinatorMockFactory
          .connect(roles.defaultAccount)
          .deploy(link.address)

        const tx = await vrfCoordinator2.callBackWithRandomness(
          eventRequestId,
          randomness,
          vrfD20.address,
        )
        const logs = await helpers.getLogs(tx)
        assert.equal(logs.length, 0)
      })
    })
  })
})
