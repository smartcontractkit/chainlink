const DEFINITION_KEYS = ['initiators', 'tasks', 'startAt', 'endAt']

export default jobSpec =>
  DEFINITION_KEYS.reduce((obj, key) => ({ ...obj, [key]: jobSpec[key] }), {})
