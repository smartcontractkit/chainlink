import jobRunsCountSelector from 'selectors/jobRunsCount'

describe('selectors - jobRunsCount', () => {
  it('returns the number of runs for the job', () => {
    const state = {
      jobs: {
        items: {
          jobA: {id: 'jobA', runsCount: 6}
        }
      }
    }

    expect(jobRunsCountSelector(state, 'jobA')).toEqual(6)
  })

  it('returns the number 0 when the job doesn\'t exist', () => {
    const state = {
      jobs: {
        items: {}
      }
    }

    expect(jobRunsCountSelector(state, 'jobA')).toEqual(0)
  })
})
