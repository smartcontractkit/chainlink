import { AppState } from 'reducers'

export default ({ dashboardIndex }: AppState): number | undefined =>
  dashboardIndex.jobRunsCount
