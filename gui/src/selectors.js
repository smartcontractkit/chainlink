export const jobsSelector = state => (
  state
    .jobs
    .currentPage
    .map(id => state.jobs.items[id])
    .filter(j => j)
)

export const jobSpecSelector = (state, jobSpecId) => (
  state.jobs.items[jobSpecId]
)

export const jobRunsSelector = (state) => {
  const runs = state.jobRuns.currentPage

  return runs
    .map(jobRunId => state.jobRuns.items[jobRunId])
    .filter(r => r)
    .sort((a, b) => {
      const dateA = new Date(a.createdAt)
      const dateB = new Date(b.createdAt)

      return dateA < dateB ? 1 : -1
    })
}

export const jobRunsCountSelector = (state, jobSpecId) => {
  const spec = jobSpecSelector(state, jobSpecId)
  return spec ? spec.runsCount : 0
}

const camelCaseToUnderscore = str => (
  str
    .replace(/[^a-zA-Z0-9]+/g, '_')
    .replace(/([A-Z]+)([A-Z][a-z])/g, '$1_$2')
    .replace(/([a-z])([A-Z])/g, '$1_$2')
    .replace(/([0-9])([^0-9])/g, '$1_$2')
    .replace(/([^0-9])([0-9])/g, '$1_$2')
    .replace(/-+/g, '_')
    .toUpperCase()
)

export const configsSelector = state => (
  Object.keys(state.configuration.config)
    .sort()
    .map(key => [camelCaseToUnderscore(key), state.configuration.config[key]])
)
