import { elapsedDuration } from '../../src/utils/elapsedDuration'

const START = '2020-01-03T22:45:00.166261Z'

describe('elapsedDuration', () => {
  test('only displays seconds when < 1 min', () => {
    const end = '2020-01-03T22:45:30.166261Z'
    expect(elapsedDuration(START, end)).toEqual('30s')
  })

  test('only displays mins & seconds when < 1 hour & > 1 min', () => {
    const end = '2020-01-03T22:46:00.166261Z'
    expect(elapsedDuration(START, end)).toEqual('1m0s')
  })

  test('only displays hours & seconds when minute is on the hour', () => {
    const end = '2020-01-03T23:45:30.166261Z'
    expect(elapsedDuration(START, end)).toEqual('1h30s')
  })

  test('displays hours, minutes and seconds when > 1 hour', () => {
    const end = '2020-01-03T23:46:30.166261Z'
    expect(elapsedDuration(START, end)).toEqual('1h1m30s')
  })

  test('can use unix timestamps for start & end', () => {
    const end = '2020-01-03T22:45:30.166261Z'
    const startUnix = new Date(START).getTime()
    const endUnix = new Date(end).getTime()

    expect(elapsedDuration(startUnix, endUnix)).toEqual('30s')
  })

  test('returns an empty string when start and end are blank', () => {
    expect(elapsedDuration('', '')).toEqual('')
  })

  test('uses current time when finishedAt is not provided', () => {
    const end = '2020-01-03T22:47:00.166261Z'

    jest
      .spyOn(Date, 'now')
      .mockImplementationOnce(() => new Date(end).valueOf())

    expect(elapsedDuration(START, null)).toEqual('2m0s')
  })
})
