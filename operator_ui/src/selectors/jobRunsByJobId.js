import build from 'redux-object'

export default (state, jobId, take) => {
  return build(state.jobRuns, 'items')
    .filter(r => r.jobId === jobId)
    .sort((a, b) => {
      const dateA = new Date(a.createdAt)
      const dateB = new Date(b.createdAt)

       return dateA < dateB ? 1 : -1
    })
    .slice(0, take)
}
