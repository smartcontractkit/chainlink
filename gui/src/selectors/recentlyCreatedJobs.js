export default ({jobs}) => (jobs.recentlyCreated && jobs
  .recentlyCreated
  .map(id => jobs.items[id])
  .filter(j => j)
)
