import { elapsedDuration } from '../../src/utils/elapsedDuration'

const START = '2020-01-03T22:45:00.166261Z'

describe('elapsedDuration', () => {
  it('only displays seconds when < 1 min', () => {
    const end = '2020-01-03T22:45:30.166261Z'
    expect(elapsedDuration(START, end)).toEqual('30s')
  })

  it('only displays mins & seconds when < 1 hour & > 1 min', () => {
    const end = '2020-01-03T22:46:00.166261Z'
    expect(elapsedDuration(START, end)).toEqual('1m0s')
  })

  it('only displays hours & seconds when minute is on the hour', () => {
    const end = '2020-01-03T23:45:30.166261Z'
    expect(elapsedDuration(START, end)).toEqual('1h30s')
  })

  it('displays hours, minutes and seconds when > 1 hour', () => {
    const end = '2020-01-03T23:46:30.166261Z'
    expect(elapsedDuration(START, end)).toEqual('1h1m30s')
  })

  it('returns an empty string when start and end are blank', () => {
    expect(elapsedDuration('', '')).toEqual('')
  })
})
