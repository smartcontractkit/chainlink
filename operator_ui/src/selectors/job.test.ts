import { partialAsFull } from 'support/test-helpers/partialAsFull'
import { INITIAL_STATE, AppState } from '../../src/reducers'
import jobSelector from '../../src/selectors/job'

describe('selectors - job', () => {
  it('returns the job item for the given id and null otherwise', () => {
    type JobsState = typeof INITIAL_STATE.jobs
    const jobsState = partialAsFull<JobsState>({
      items: {
        jobA: { id: 'jobA' },
      },
    })
    const state: Pick<AppState, 'jobs'> = {
      jobs: jobsState,
    }

    expect(jobSelector(state, 'jobA')).toEqual({ id: 'jobA' })
    expect(jobSelector(state, 'joba')).toBeNull()
  })
})
