import { constantCase } from 'change-case'
import { AppState } from 'reducers'

export default ({ configuration }: Pick<AppState, 'configuration'>) => {
  const { data } = configuration

  return Object.keys(data)
    .sort()
    .map((key) => [constantCase(key), data[key]])
}
