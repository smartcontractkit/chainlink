import elapsedTimeHHMMSS from 'utils/elapsedTimeHHMMSS'
import { MINUTE_MS, TWO_MINUTES_MS } from 'test-helpers/isoDate'

describe('Converts the difference between two dates to HH:MM:SS format', () => {
    it('Doesnt display minutes when elapsed time is less than a minute', () => {
        expect(elapsedTimeHHMMSS(MINUTE_MS + MINUTE_MS / 2, TWO_MINUTES_MS)).toEqual('30s')
    })
    it('Doesnt display hours when elapsed time is less than an hour', () => {
        expect(elapsedTimeHHMMSS(MINUTE_MS, TWO_MINUTES_MS)).toEqual('1m0s')
    })
    it('Displays hours and seconds when elapsed time is zero MOD hour', () => {
        expect(elapsedTimeHHMMSS(MINUTE_MS, MINUTE_MS * 61)).toEqual('1h0s')
    })
    it('Displays hours, minutes and seconds when elapsed time composite', () => {
        expect(elapsedTimeHHMMSS(MINUTE_MS, MINUTE_MS * 61 + MINUTE_MS * 3 / 2)).toEqual('1h1m30s')
    })
    it('Returns an empty string when no start and end dates are passed', () => {
        expect(elapsedTimeHHMMSS(undefined, undefined)).toEqual('')
    })
})
