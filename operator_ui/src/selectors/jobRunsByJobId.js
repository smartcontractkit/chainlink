import build from 'redux-object'

export default (state, jobId, take) => {
  return build(state.jobRuns, 'items')
    .filter(r => r.jobId === jobId)
    .slice(0, take)
}
