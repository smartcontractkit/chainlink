import { AppState } from 'reducers'

export default ({ configuration }: Pick<AppState, 'configuration'>) => {
  const { data } = configuration

  return Object.keys(data)
    .sort()
    .map((key) => [key, data[key]])
}
