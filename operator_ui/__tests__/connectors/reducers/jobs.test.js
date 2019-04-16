import reducer from 'connectors/redux/reducers'
import {
  UPSERT_JOBS,
  UPSERT_RECENTLY_CREATED_JOBS,
  UPSERT_JOB
} from 'connectors/redux/reducers/jobs'

describe('connectors/reducers/jobs', () => {
  it('should return the initial state', () => {
    const state = reducer(undefined, {})

    expect(state.jobs).toEqual({
      items: {},
      currentPage: null,
      recentlyCreated: null,
      count: 0
    })
  })

  it('UPSERT_JOBS upserts items along with the current page & count from meta', () => {
    const action = {
      type: UPSERT_JOBS,
      data: {
        specs: {
          a: { id: 'a' },
          b: { id: 'b' }
        },
        meta: {
          currentPageJobs: {
            data: [{ id: 'b' }, { id: 'a' }],
            meta: {
              count: 10
            }
          }
        }
      }
    }
    const state = reducer(undefined, action)

    expect(state.jobs.items).toEqual({
      a: { id: 'a' },
      b: { id: 'b' }
    })
    expect(state.jobs.count).toEqual(10)
    expect(state.jobs.currentPage).toEqual(['b', 'a'])
  })

  it('UPSERT_RECENTLY_CREATED_JOBS upserts items along with the current page & count from meta', () => {
    const action = {
      type: UPSERT_RECENTLY_CREATED_JOBS,
      data: {
        specs: {
          c: { id: 'c' },
          d: { id: 'd' }
        },
        meta: {
          recentlyCreatedJobs: {
            data: [{ id: 'd' }, { id: 'c' }]
          }
        }
      }
    }
    const state = reducer(undefined, action)

    expect(state.jobs.items).toEqual({
      c: { id: 'c' },
      d: { id: 'd' }
    })
    expect(state.jobs.recentlyCreated).toEqual(['d', 'c'])
  })

  it('UPSERT_JOB upserts items', () => {
    const action = {
      type: UPSERT_JOB,
      data: {
        specs: {
          a: { id: 'a' }
        }
      }
    }
    const state = reducer(undefined, action)

    expect(state.jobs.items).toEqual({ a: { id: 'a' } })
  })
})
