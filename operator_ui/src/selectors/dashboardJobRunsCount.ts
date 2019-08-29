import { AppState } from 'connectors/redux/reducers'

export default ({ dashboardIndex }: AppState): number | undefined =>
  dashboardIndex.jobRunsCount
