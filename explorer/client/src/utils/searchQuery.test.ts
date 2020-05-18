import { partialAsFull } from '@chainlink/ts-helpers'
import { searchQuery } from './searchQuery'

describe('utils/searchQuery', () => {
  it('returns an empty string when there is no search param', () => {
    const searchLocation = partialAsFull<Location>({ search: '' })
    expect(searchQuery(searchLocation)).toEqual('')
  })

  it('returns the "search" query parameter', () => {
    const searchLocation = partialAsFull<Location>({
      search: '?search=find-me',
    })
    expect(searchQuery(searchLocation)).toEqual('find-me')
  })
})
