export default state => (
  state
    .jobs
    .currentPage
    .map(id => state.jobs.items[id])
    .filter(j => j)
)
