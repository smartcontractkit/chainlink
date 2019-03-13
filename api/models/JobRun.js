module.exports = {
  tableName: 'job_runs',
  attributes: {
    updatedAt: false,
    requestId: { type: 'string', required: true }
  }
}
