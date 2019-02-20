import build from 'redux-object'

export default ({ dashboardIndex, jobRuns }) =>
  dashboardIndex.recentJobRuns &&
  dashboardIndex.recentJobRuns
    .map(id => build(jobRuns, 'items', id))
    .filter(r => r)
