const index = async (req, res) => {
  const jobRuns = await sails.models.jobrun.find({
    limit: 10,
    sort: 'createdAt DESC'
  })

  return res.send(jobRuns)
}

module.exports = index
