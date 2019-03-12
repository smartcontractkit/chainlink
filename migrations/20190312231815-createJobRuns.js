exports.setup = function(options, seedLink) {}

exports.up = function(db, callback) {
  db.createTable(
    'job_runs',
    {
      id: { type: 'int', primaryKey: true, autoIncrement: true },
      requestId: { type: 'uuid', notNull: true },
      createdAt: { type: 'timestamp', notNull: true }
    },
    callback
  )
}

exports.down = function(db, callback) {
  db.dropTable('job_runs', callback)
}

exports._meta = {
  version: 1
}
