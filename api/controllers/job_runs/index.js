// TODO: Figure out how to disable globals correctly
// const JobRun = require('../../models/JobRun')

const index = async (req, res) => {
  const jobRuns = await JobRun.find({
    limit: 10,
    sort: 'createdAt DESC'
  })

  return res.send(jobRuns)
}

module.exports = index
