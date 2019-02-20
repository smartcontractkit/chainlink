import build from 'redux-object'

export default ({ jobs }) =>
  jobs.currentPage &&
  jobs.currentPage.map(id => build(jobs, 'items', id)).filter(j => j)
