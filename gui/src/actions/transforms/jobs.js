import serializeJob from 'connectors/redux/serializers/job'

export default (actionType, json, jobSerializer = serializeJob) => ({
  type: actionType,
  count: json.meta && json.meta.count,
  items: json.data.map(jobSerializer)
})
