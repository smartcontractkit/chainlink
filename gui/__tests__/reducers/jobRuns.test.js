import reducer from 'connectors/redux/reducers'
import {
  RECEIVE_JOB_SPEC_SUCCESS
} from 'actions'

describe('jobRuns reducer', () => {
  it('should return the initial state', () => {
    const state = reducer(undefined, {})

    expect(state.jobRuns).toEqual({
      items: {}
    })
  })

  it('RECEIVE_JOB_SPEC_SUCCESS stores the runs by id', () => {
    const action = {
      type: RECEIVE_JOB_SPEC_SUCCESS,
      item: {
        id: '50208cd6b3034594b8e999c380066b67',
        runs: [{id: 'a'}, {id: 'b'}]
      }
    }
    const state = reducer(undefined, action)

    expect(state.jobRuns).toEqual({
      items: {
        'a': {id: 'a'},
        'b': {id: 'b'}
      }
    })
  })
})
