import uuid from 'uuid/v4'

export default ({
  id,
  jobId,
  initiator,
  taskRuns,
  status,
  result,
  createdAt,
}) => {
  const _id = id || uuid().replace(/-/g, '')
  const _jobId = jobId || uuid().replace(/-/g, '')
  const _status = status || 'completed'
  const _createdAt = createdAt || '2018-06-19T15:39:53.315919143-07:00'
  const _initiator = initiator || { type: 'web', params: {} }
  const _taskRuns = taskRuns || []
  const _result = result || {
    data: {
      value:
        '0x05070f7f6a40e4ce43be01fa607577432c68730c2cb89a0f50b665e980d926b5',
    },
  }

  return {
    data: {
      id: _id,
      type: 'runs',
      attributes: {
        id: _id,
        jobId: _jobId,
        initiator: _initiator,
        taskRuns: _taskRuns,
        result: _result,
        status: _status,
        createdAt: _createdAt,
      },
    },
  }
}
