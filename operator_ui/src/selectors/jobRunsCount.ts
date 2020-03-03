import { AppState } from 'reducers'

export default ({ jobRuns }: Pick<AppState, 'jobRuns'>) =>
  jobRuns.currentJobRunsCount
