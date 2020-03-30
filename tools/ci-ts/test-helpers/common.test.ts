import { assert } from 'chai'
import { logEvents } from './common'
import * as m from './__mocks__/contract'

const mockLog = jest.spyOn(global.console, 'log')

afterEach(() => {
  mockLog.mockReset()
  ;(m.contract.on as jest.Mock).mockReset()
})

describe('test-helpers/common', () => {
  describe('#logEvents', () => {
    it('accepts a valid events', () => {
      logEvents(m.contract, m.name, 'event1')
    })

    it('accepts valid events', () => {
      logEvents(m.contract, m.name, m.eventNames)
    })

    it('rejects an invalid event', () => {
      assert.throws(() => {
        logEvents(m.contract, m.name, 'invalid')
      })
    })

    it('rejects invalid events mixed with valid ones', () => {
      assert.throws(() => {
        logEvents(m.contract, m.name, [...m.eventNames, 'invalid'])
      })
    })

    it('registers an event handler', () => {
      logEvents(m.contract, m.name, m.eventNames)
      expect(m.contract.on).toHaveBeenCalledTimes(1)
      expect(m.contract.on).toHaveBeenCalledWith('*', expect.any(Function))
    })

    describe('event callback', () => {
      let handler = (_: any) => {}

      beforeEach(() => {
        ;(m.contract.on as jest.Mock).mockImplementation(
          (_, fun) => (handler = fun),
        )
      })

      it('event handler filters by event', () => {
        logEvents(m.contract, m.name, ['event1'])
        handler(m.eventEmission1)
        expect(mockLog).toHaveBeenCalledTimes(1)
        handler(m.eventEmission2)
        expect(mockLog).toHaveBeenCalledTimes(1)
        expect(mockLog.mock.calls[0]).toMatchInlineSnapshot(`
          Array [
            "event1 event emitted by TEST_CONTRACT in block #1
          	* input1: value1",
          ]
        `)
      })

      it('includes all events by default', () => {
        logEvents(m.contract, m.name)
        handler(m.eventEmission1)
        expect(mockLog).toHaveBeenCalledTimes(1)
        handler(m.eventEmission2)
        expect(mockLog).toHaveBeenCalledTimes(2)
        expect(mockLog.mock.calls[0]).toMatchInlineSnapshot(`
          Array [
            "event1 event emitted by TEST_CONTRACT in block #1
          	* input1: value1",
          ]
        `)
        expect(mockLog.mock.calls[1]).toMatchInlineSnapshot(`
          Array [
            "event2 event emitted by TEST_CONTRACT in block #1
          	* input2: value2",
          ]
        `)
      })
    })
  })
})
