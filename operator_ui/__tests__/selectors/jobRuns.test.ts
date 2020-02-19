import jobRunsSelector from '../../src/selectors/jobRuns'

describe('selectors - jobRuns', () => {
  it('returns the job runs for the given job spec id', () => {
    const state = {
      jobRuns: {
        currentPage: ['runA', 'runB'],
        items: {
          runA: { id: 'runA' },
          runB: { id: 'runB' },
          runC: { id: 'runC' },
        },
      },
    }

    const runs = jobRunsSelector(state)

    expect(runs).toEqual([{ id: 'runA' }, { id: 'runB' }])
  })

  it('returns an empty array when the currentPage is empty', () => {
    const state = {
      jobRuns: {
        currentPage: [],
        items: {
          runA: { id: 'runA' },
        },
      },
    }
    const runs = jobRunsSelector(state)

    expect(runs).toEqual([])
  })

  it('excludes job runs that do not have items', () => {
    const state = {
      jobRuns: {
        currentPage: ['runA', 'runB'],
        items: {
          runA: { id: 'runA' },
        },
      },
    }
    const runs = jobRunsSelector(state)

    expect(runs).toEqual([{ id: 'runA' }])
  })
})
