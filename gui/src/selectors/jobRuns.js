export default state => {
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
