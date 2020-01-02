import { parseParams } from '../../utils/pagination'

describe('parseParams', () => {
  it('returns a default page & size limit for the query params', () => {
    const params = parseParams({})

    expect(params).toEqual({
      page: 1,
      limit: 10,
    })
  })

  it('returns the page & size limit from the query params', () => {
    const params = parseParams({ page: '2', size: '11' })

    expect(params).toEqual({
      page: 2,
      limit: 11,
    })
  })

  it('returns a max size limit for the query params', () => {
    const params = parseParams({ size: '101' })

    expect(params.limit).toEqual(100)
  })
})
