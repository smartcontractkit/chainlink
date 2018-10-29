import reducer from 'connectors/redux/reducers'
import {
  REQUEST_JOBS,
  RECEIVE_JOBS_SUCCESS,
  RECEIVE_JOBS_ERROR,
  RECEIVE_RECENTLY_CREATED_JOBS_SUCCESS,
  RECEIVE_JOB_SPEC_SUCCESS
} from 'actions'

describe('connectors/reducers/jobs', () => {
  it('should return the initial state', () => {
    const state = reducer(undefined, {})

    expect(state.jobs).toEqual({
      items: {},
      currentPage: [],
      recentlyCreated: null,
      count: 0,
      networkError: false
    })
  })

  it('REQUEST_JOBS disables the network error', () => {
    const action = {type: REQUEST_JOBS}
    const state = reducer(undefined, action)

    expect(state.jobs.networkError).toEqual(false)
  })

  it('RECEIVE_JOBS_SUCCESS stores the job items and the current page', () => {
    const action = {
      type: RECEIVE_JOBS_SUCCESS,
      items: [{id: 'a'}, {id: 'b'}]
    }
    const state = reducer(undefined, action)

    expect(state.jobs.items).toEqual({
      'a': {id: 'a'},
      'b': {id: 'b'}
    })
    expect(state.jobs.currentPage).toEqual(['a', 'b'])
    expect(state.jobs.networkError).toEqual(false)
  })

  it('RECEIVE_JOBS_ERROR updates the network error', () => {
    const previousState = {
      jobs: {networkError: false}
    }
    const action = {
      type: RECEIVE_JOBS_ERROR,
      networkError: true
    }
    const state = reducer(previousState, action)

    expect(state.jobs.networkError).toEqual(true)
  })

  it('RECEIVE_RECENTLY_CREATED_JOBS_SUCCESS stores the job items and the order', () => {
    const action = {
      type: RECEIVE_RECENTLY_CREATED_JOBS_SUCCESS,
      items: [{id: 'b'}, {id: 'a'}]
    }
    const state = reducer(undefined, action)

    expect(state.jobs.items).toEqual({
      'a': {id: 'a'},
      'b': {id: 'b'}
    })
    expect(state.jobs.recentlyCreated).toEqual(['b', 'a'])
  })

  it('RECEIVE_JOB_SPEC_SUCCESS assigns runsCount', () => {
    const previousState = {
      jobs: {
        items: {
          '50208cd6b3034594b8e999c380066b67': {
            id: '50208cd6b3034594b8e999c380066b67',
            runsCount: 2
          }
        }
      }
    }
    const action = {
      type: RECEIVE_JOB_SPEC_SUCCESS,
      item: {
        id: '50208cd6b3034594b8e999c380066b67',
        runs: [{id: 'a'}, {id: 'b'}]
      }
    }
    const state = reducer(previousState, action)

    expect(state.jobs.items['50208cd6b3034594b8e999c380066b67'].runsCount).toEqual(2)
  })

  it('RECEIVE_JOB_SPEC_SUCCESS assigns runsCount to 0 when the job doesn\'t have runs', () => {
    const previousState = {
      jobs: {
        items: {}
      }
    }
    const action = {
      type: RECEIVE_JOB_SPEC_SUCCESS,
      item: {
        id: '50208cd6b3034594b8e999c380066b67'
      }
    }
    const state = reducer(previousState, action)

    expect(state.jobs.items['50208cd6b3034594b8e999c380066b67'].runsCount).toEqual(0)
  })
})
