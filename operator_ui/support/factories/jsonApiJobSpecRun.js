import uuid from 'uuid/v4'

export default (attrs) => {
  const _id = attrs.id || uuid().replace(/-/g, '')
  const _jobId = attrs.jobId || uuid().replace(/-/g, '')
  const _status = attrs.status || 'completed'
  const _createdAt = attrs.createdAt || '2018-06-19T15:39:53.315919143-07:00'
  const _initiator = attrs.initiator || { type: 'web', params: {} }
  const _taskRuns = attrs.taskRuns || []
  const _result = attrs.result || {
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
