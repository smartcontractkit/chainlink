import { ethers } from 'hardhat'
import { assert } from 'chai'
import { UpkeepTranscoder30__factory as UpkeepTranscoderFactory } from '../../../typechain/factories/UpkeepTranscoder30__factory'
import { UpkeepTranscoder30 as UpkeepTranscoder } from '../../../typechain/UpkeepTranscoder30'
import { evmRevert } from '../../test-helpers/matchers'

let upkeepTranscoderFactory: UpkeepTranscoderFactory
let transcoder: UpkeepTranscoder
let balance: number
let executeGas: number
let amountSpent: number
let admin0: string
let target0: string
let lastKeeper0: string
let admin1: string
let target1: string
let lastKeeper1: string
let idx: number[]
let upkeepsV1: any[]
let upkeepsV2: any[]
let upkeepsV3: any[]
let admins: string[]

const encodeUpkeepV1 = (ids: number[], upkeeps: any[], checkDatas: any[]) => {
  return ethers.utils.defaultAbiCoder.encode(
    [
      'uint256[]',
      'tuple(uint96,address,uint32,uint64,address,uint96,address)[]',
      'bytes[]',
    ],
    [ids, upkeeps, checkDatas],
  )
}

const encodeUpkeepV2 = (ids: number[], upkeeps: any[], checkDatas: any[]) => {
  return ethers.utils.defaultAbiCoder.encode(
    [
      'uint256[]',
      'tuple(uint96,address,uint96,address,uint32,uint32,address,bool)[]',
      'bytes[]',
    ],
    [ids, upkeeps, checkDatas],
  )
}

const encodeUpkeepV3 = (
  ids: number[],
  upkeeps: any[],
  checkDatas: any[],
  admins: string[],
) => {
  return ethers.utils.defaultAbiCoder.encode(
    [
      'uint256[]',
      'tuple(uint32,uint32,bool,address,uint96,uint96,uint32)[]',
      'bytes[]',
      'address[]',
    ],
    [ids, upkeeps, checkDatas, admins],
  )
}

before(async () => {
  // @ts-ignore bug in autogen file
  upkeepTranscoderFactory = await ethers.getContractFactory(
    'UpkeepTranscoder3_0',
  )
})

describe.only('UpkeepTranscoder3_0', () => {
  balance = 50000000000000
  executeGas = 200000
  amountSpent = 200000000000000
  admin0 = '0xe380f971547a36c055370d02b1bbb4f27f038c61'
  target0 = '0xffffffffffffffffffffffffffffffffffffffff'
  lastKeeper0 = '0x233a95ccebf3c9f934482c637c08b4015cdd6ddd'
  admin1 = '0xe380f971547a36c055370d02b1bbb4f27f038c60'
  target1 = '0xfffffffffffffffffffffffffffffffffffffffe'
  lastKeeper1 = '0x233a95ccebf3c9f934482c637c08b4015cdd6ddc'
  const UpkeepFormatV1 = 0
  const UpkeepFormatV2 = 1
  const UpkeepFormatV3 = 2

  beforeEach(async () => {
    transcoder = await upkeepTranscoderFactory.deploy()
  })

  describe('#typeAndVersion', () => {
    it('uses the correct type and version', async () => {
      const typeAndVersion = await transcoder.typeAndVersion()
      assert.equal(typeAndVersion, 'UpkeepTranscoder 3.0.0')
    })
  })

  describe('#transcodeUpkeeps', () => {
    const encodedData = '0xabcd'

    it('reverts if the from type is not an enum value', async () => {
      await evmRevert(
        transcoder.transcodeUpkeeps(3, 2, encodedData),
        'function was called with incorrect parameters',
      )
    })

    it('reverts if the to version is not 2', async () => {
      await evmRevert(
        transcoder.transcodeUpkeeps(1, 3, encodedData),
        'function was called with incorrect parameters',
      )
    })

    context('when from and to versions equal', () => {
      it('returns the data that was passed in', async () => {
        let response = await transcoder.transcodeUpkeeps(
          UpkeepFormatV1,
          UpkeepFormatV1,
          encodedData,
        )
        assert.equal(encodedData, response)

        response = await transcoder.transcodeUpkeeps(
          UpkeepFormatV2,
          UpkeepFormatV2,
          encodedData,
        )
        assert.equal(encodedData, response)

        response = await transcoder.transcodeUpkeeps(
          UpkeepFormatV3,
          UpkeepFormatV3,
          encodedData,
        )
        assert.equal(encodedData, response)
      })
    })

    context('when from and to versions are correct', () => {
      idx = [123, 124]

      upkeepsV3 = [
        [executeGas, 2 ** 32 - 1, false, target0, amountSpent, balance, 0],
        [executeGas, 2 ** 32 - 1, false, target1, amountSpent, balance, 0],
      ]

      admins = [admin0, admin1]

      it('transcodes V1 upkeeps to V3 properly', async () => {
        upkeepsV1 = [
          [
            balance,
            lastKeeper0,
            executeGas,
            2 ** 32,
            target0,
            amountSpent,
            admin0,
          ],
          [
            balance,
            lastKeeper1,
            executeGas,
            2 ** 32,
            target1,
            amountSpent,
            admin1,
          ],
        ]

        let data = await transcoder.transcodeUpkeeps(
          UpkeepFormatV1,
          UpkeepFormatV3,
          encodeUpkeepV1(idx, upkeepsV1, ['0xabcd', '0xffff']),
        )
        assert.equal(
          encodeUpkeepV3(idx, upkeepsV3, ['0xabcd', '0xffff'], admins),
          data,
        )
      })

      it('transcodes V2 upkeeps to V3 properly', async () => {
        upkeepsV2 = [
          [
            balance,
            lastKeeper0,
            amountSpent,
            admin0,
            executeGas,
            2 ** 32 - 1,
            target0,
            false,
          ],
          [
            balance,
            lastKeeper1,
            amountSpent,
            admin1,
            executeGas,
            2 ** 32 - 1,
            target1,
            false,
          ],
        ]

        let data = await transcoder.transcodeUpkeeps(
          UpkeepFormatV2,
          UpkeepFormatV3,
          encodeUpkeepV2(idx, upkeepsV2, ['0xabcd', '0xffff']),
        )
        assert.equal(
          encodeUpkeepV3(idx, upkeepsV3, ['0xabcd', '0xffff'], admins),
          data,
        )
      })
    })
  })
})
