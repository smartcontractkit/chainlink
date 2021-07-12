import { partialAsFull } from 'support/test-helpers/partialAsFull'
import { INITIAL_STATE, AppState } from '../../src/reducers'
import recentlyCreatedJobsSelector from '../../src/selectors/recentlyCreatedJobs'

describe('selectors - jobs', () => {
  type JobsState = typeof INITIAL_STATE.jobs

  it('returns null when not loaded', () => {
    const jobsState = partialAsFull<JobsState>({
      recentlyCreated: undefined,
    })
    const state: Pick<AppState, 'jobs'> = {
      jobs: jobsState,
    }
    const jobs = recentlyCreatedJobsSelector(state)

    expect(jobs).toEqual(undefined)
  })

  it('returns the job objects in items and excludes those not present', () => {
    const jobsState = partialAsFull<JobsState>({
      recentlyCreated: ['jobA', 'jobB', 'jobC'],
      items: {
        jobA: { id: 'jobA' },
        jobB: { id: 'jobB' },
      },
    })
    const state: Pick<AppState, 'jobs'> = {
      jobs: jobsState,
    }
    const jobs = recentlyCreatedJobsSelector(state)

    expect(jobs).toEqual([{ id: 'jobA' }, { id: 'jobB' }])
  })
})
