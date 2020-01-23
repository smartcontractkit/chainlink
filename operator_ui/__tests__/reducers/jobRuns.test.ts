import reducer, { INITIAL_STATE } from '../../src/reducers'
import {
  UpsertJobRunAction,
  UpsertJobRunsAction,
  UpsertRecentJobRunsAction,
  ReceiveDeleteSuccessAction,
  ResourceActionType,
} from '../../src/reducers/actions'

describe('reducers/jobRuns', () => {
  it('UPSERT_JOB_RUNS upserts items along with the current page and count', () => {
    const action: UpsertJobRunsAction = {
      type: ResourceActionType.UPSERT_JOB_RUNS,
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
    const state = reducer(INITIAL_STATE, action)

    expect(state.jobRuns.items).toEqual({
      a: { id: 'a' },
      b: { id: 'b' },
    })
    expect(state.jobRuns.currentPage).toEqual(['b', 'a'])
    expect(state.jobRuns.currentJobRunsCount).toEqual(10)
  })

  it('UPSERT_RECENT_JOB_RUNS upserts items', () => {
    const action: UpsertRecentJobRunsAction = {
      type: ResourceActionType.UPSERT_RECENT_JOB_RUNS,
      data: {
        runs: { a: { id: 'a' } },
        meta: {
          recentJobRuns: {
            data: [{ id: 'a' }],
            meta: { count: 1 },
          },
        },
      },
    }
    const state = reducer(INITIAL_STATE, action)

    expect(state.jobRuns.items).toEqual({
      a: { id: 'a' },
    })
  })

  it('UPSERT_JOB_RUN upserts items', () => {
    const action: UpsertJobRunAction = {
      type: ResourceActionType.UPSERT_JOB_RUN,
      data: {
        runs: {
          a: { id: 'a' },
        },
      },
    }
    const state = reducer(INITIAL_STATE, action)

    expect(state.jobRuns.items).toEqual({
      a: { id: 'a' },
    })
  })

  it('RECEIVE_DELETE_SUCCESS deletes jobrun associations', () => {
    const upsertAction: UpsertJobRunAction = {
      type: ResourceActionType.UPSERT_JOB_RUN,
      data: {
        runs: {
          b: { attributes: { jobId: 'b' } },
        },
      },
    }
    const deleteAction: ReceiveDeleteSuccessAction = {
      type: ResourceActionType.RECEIVE_DELETE_SUCCESS,
      id: 'b',
    }

    const preDeleteState = reducer(INITIAL_STATE, upsertAction)
    expect(preDeleteState.jobRuns.items).toEqual({
      b: { attributes: { jobId: 'b' } },
    })

    const postDeleteState = reducer(preDeleteState, deleteAction)
    expect(postDeleteState.jobRuns.items).toEqual({})
  })
})
