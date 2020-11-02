import React from 'react'
import { v2 } from 'api'
import { Route, RouteComponentProps, Switch } from 'react-router-dom'
import { useErrorHandler } from 'hooks/useErrorHandler'
import { useLoadingPlaceholder } from 'hooks/useLoadingPlaceholder'
import { JobData } from './sharedTypes'
import { JobsDefinition } from './Definition'
import { JobsErrors } from './Errors'
import { RecentRuns } from './RecentRuns'
import { RegionalNav } from './RegionalNav'

type Props = RouteComponentProps<{
  jobSpecId: string
}>

const DEFAULT_PAGE = 1
const RECENT_RUNS_COUNT = 5

export const JobsShow: React.FC<Props> = ({ match }) => {
  const [state, setState] = React.useState<JobData>({
    recentRuns: [],
    recentRunsCount: 0,
  })
  const { jobSpec } = state
  const { error, ErrorComponent, setError } = useErrorHandler()
  const { LoadingPlaceholder } = useLoadingPlaceholder(!error && !jobSpec)

  const { jobSpecId } = match.params

  const getJobSpecRuns = React.useCallback(
    () =>
      v2.runs
        .getJobSpecRuns({
          jobSpecId,
          page: DEFAULT_PAGE,
          size: RECENT_RUNS_COUNT,
        })
        .then((jobSpecRunsResponse) => {
          setState((s) => ({
            ...s,
            recentRuns: jobSpecRunsResponse.data,
            recentRunsCount: jobSpecRunsResponse.meta.count,
          }))
        })
        .catch(setError),
    [jobSpecId, setError],
  )

  const getJobSpec = React.useCallback(
    async () =>
      v2.specs
        .getJobSpec(jobSpecId)
        .then((response) =>
          setState((s) => ({
            ...s,
            jobSpec: response.data,
          })),
        )
        .catch(setError),
    [jobSpecId, setError],
  )

  React.useEffect(() => {
    getJobSpec()
  }, [getJobSpec])

  return (
    <div>
      <RegionalNav
        jobSpecId={jobSpecId}
        job={jobSpec}
        getJobSpecRuns={getJobSpecRuns}
      />
      <Switch>
        <Route
          path={`${match.path}/json`}
          render={() => (
            <JobsDefinition
              {...{
                ...state,
                ErrorComponent,
                LoadingPlaceholder,
                error,
              }}
            />
          )}
        />
        <Route
          path={`${match.path}/errors`}
          render={() => (
            <JobsErrors
              {...{
                ...state,
                ErrorComponent,
                LoadingPlaceholder,
                error,
                getJobSpec,
                setState,
              }}
            />
          )}
        />
        <Route
          path={`${match.path}`}
          render={() => (
            <RecentRuns
              {...{
                ...state,
                error,
                ErrorComponent,
                LoadingPlaceholder,
                getJobSpecRuns,
              }}
            />
          )}
        />
      </Switch>
    </div>
  )
}

export default JobsShow
