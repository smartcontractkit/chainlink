import { constantCase } from 'change-case'

export default state =>
  Object.keys(state.configuration.config)
    .sort()
    .map(key => [constantCase(key), state.configuration.config[key]])
