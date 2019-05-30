const DEFINITION_KEYS = ['initiators', 'tasks', 'startAt', 'endAt']
const SCRUBBED_KEYS = ['ID', 'CreatedAt', 'DeletedAt', 'UpdatedAt']

const scrub = payload => {
  if (Array.isArray(payload)) {
    return payload.map(p => scrub(p))
  }
  if (typeof payload !== 'object' || payload === null) {
    return payload
  }
  const keepers = Object.keys(payload).filter(k => !SCRUBBED_KEYS.includes(k))
  return keepers.reduce((obj, key) => ({ ...obj, [key]: payload[key] }), {})
}

export default jobSpec =>
  DEFINITION_KEYS.reduce(
    (obj, key) => ({ ...obj, [key]: scrub(jobSpec[key]) }),
    {}
  )
