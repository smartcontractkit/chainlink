import reducer from 'connectors/redux/reducers'
import { RECEIVE_DELETE_SUCCESS } from '../../../src/actions'

describe('connectors/reducers/jobRuns', () => {
  it('should return the initial state', () => {
    const state = reducer(undefined, {})

    expect(state.jobRuns).toEqual({
      currentPage: undefined,
      currentJobRunsCount: undefined,
      items: {},
    })
  })

  it('UPSERT_JOB_RUNS upserts items along with the current page and count', () => {
    const action = {
      type: 'UPSERT_JOB_RUNS',
      data: {
        runs: {
          a: { id: 'a' },
          b: { id: 'b' },
        },
        meta: {
          currentPageJobRuns: {
            data: [{ id: 'b' }, { id: 'a' }],
            meta: {
              count: 10,
            },
          },
        },
      },
    }
    const state = reducer(undefined, action)

    expect(state.jobRuns.items).toEqual({
      a: { id: 'a' },
      b: { id: 'b' },
    })
    expect(state.jobRuns.currentPage).toEqual(['b', 'a'])
    expect(state.jobRuns.currentJobRunsCount).toEqual(10)
  })

  it('UPSERT_RECENT_JOB_RUNS upserts items', () => {
    const action = {
      type: 'UPSERT_RECENT_JOB_RUNS',
      data: {
        runs: { a: { id: 'a' } },
        meta: {
          recentJobRuns: { data: [], meta: {} },
        },
      },
    }
    const state = reducer(undefined, action)

    expect(state.jobRuns.items).toEqual({
      a: { id: 'a' },
    })
  })

  it('UPSERT_JOB_RUN upserts items', () => {
    const action = {
      type: 'UPSERT_JOB_RUN',
      data: {
        runs: {
          a: { id: 'a' },
        },
      },
    }
    const state = reducer(undefined, action)

    expect(state.jobRuns.items).toEqual({
      a: { id: 'a' },
    })
  })

  it('RECEIVE_DELETE_SUCCESS deletes jobrun associations', () => {
    const upsertAction = {
      type: 'UPSERT_JOB_RUN',
      data: {
        runs: {
          b: { attributes: { jobId: 'b' } },
        },
      },
    }
    const preDeleteState = reducer(undefined, upsertAction)
    expect(preDeleteState.jobRuns.items).toEqual({
      b: { attributes: { jobId: 'b' } },
    })
    const deleteAction = {
      type: RECEIVE_DELETE_SUCCESS,
      response: 'b',
    }
    const postDeleteState = reducer(preDeleteState, deleteAction)
    expect(postDeleteState.jobRuns.items).toEqual({})
  })
})
