import { ethers } from 'hardhat'
import { assert } from 'chai'
import { evmRevert } from '../../test-helpers/matchers'
import { UpkeepTranscoder__factory as UpkeepTranscoderFactory } from '../../../typechain/factories/UpkeepTranscoder__factory'
import { UpkeepTranscoder } from '../../../typechain'

let upkeepTranscoderFactory: UpkeepTranscoderFactory
let transcoder: UpkeepTranscoder

before(async () => {
  upkeepTranscoderFactory = await ethers.getContractFactory('UpkeepTranscoder')
})

describe('UpkeepTranscoder', () => {
  const formatV1 = 0
  const formatV2 = 1
  const formatV3 = 2

  beforeEach(async () => {
    transcoder = await upkeepTranscoderFactory.deploy()
  })

  describe('#typeAndVersion', () => {
    it('uses the correct type and version', async () => {
      const typeAndVersion = await transcoder.typeAndVersion()
      assert.equal(typeAndVersion, 'UpkeepTranscoder 1.0.0')
    })
  })

  describe('#transcodeUpkeeps', () => {
    const encodedData = '0xc0ffee'

    it('reverts if the from type is not an enum value', async () => {
      await evmRevert(
        transcoder.transcodeUpkeeps(3, 1, encodedData),
        'function was called with incorrect parameters',
      )
    })

    it('reverts if the from type != to type', async () => {
      await evmRevert(
        transcoder.transcodeUpkeeps(1, 2, encodedData),
        'InvalidTranscoding()',
      )
    })

    context('when from and to versions equal', () => {
      it('returns the data that was passed in', async () => {
        let response = await transcoder.transcodeUpkeeps(
          formatV1,
          formatV1,
          encodedData,
        )
        assert.equal(encodedData, response)

        response = await transcoder.transcodeUpkeeps(
          formatV2,
          formatV2,
          encodedData,
        )
        assert.equal(encodedData, response)

        response = await transcoder.transcodeUpkeeps(
          formatV3,
          formatV3,
          encodedData,
        )
        assert.equal(encodedData, response)
      })
    })
  })
})
