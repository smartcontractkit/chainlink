import { constantCase } from 'change-case'

export default state => {
  const data = state.configuration.data
  return Object.keys(data)
    .sort()
    .map(key => [constantCase(key), data[key]])
}
