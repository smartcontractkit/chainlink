import build from 'redux-object'

export default ({ jobRuns }) => {
  return (
    jobRuns.currentPage &&
    jobRuns.currentPage.map(id => build(jobRuns, 'items', id)).filter(r => r)
  )
}
