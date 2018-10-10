import { constantCase } from 'change-case'

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

export const jobRunSelector = (state, id) => state.jobRuns.items[id]

export const jobRunsCountSelector = (state, jobSpecId) => {
  const spec = jobSpecSelector(state, jobSpecId)
  return spec ? spec.runsCount : 0
}

export const bridgeSelector = (state, id) => state.bridges.items[id]

export const bridgesSelector = state => {
  const bridgeIds = state.bridges.currentPage

  return bridgeIds
    .map(id => state.bridges.items[id])
    .filter(r => r)
}

export const configsSelector = state => (
  Object.keys(state.configuration.config)
    .sort()
    .map(key => [constantCase(key), state.configuration.config[key]])
)

export const isFetchingSelector = (state) => state.fetching.count > 0
