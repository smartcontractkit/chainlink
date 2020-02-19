import { AppState } from 'reducers'

export default ({ fetching }: Pick<AppState, 'fetching'>) => fetching.count
