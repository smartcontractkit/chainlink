import React from 'react'
import { v2 } from 'api'
import { Route, RouteComponentProps, Switch } from 'react-router-dom'
import { useErrorHandler } from 'hooks/useErrorHandler'
import { useLoadingPlaceholder } from 'hooks/useLoadingPlaceholder'
import jobSpecDefinition from 'utils/jobSpecDefinition'
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
  const { job, jobSpec } = state
  const { error, ErrorComponent, setError } = useErrorHandler()
  const { LoadingPlaceholder } = useLoadingPlaceholder(!error && !jobSpec)

  const { jobSpecId } = match.params
  // `isNaN` actually accepts strings and we don't want to `parseInt` or `parseFloat`
  //  as it doesn't have the behaviour we want.
  const isOcrJob = !isNaN((jobSpecId as unknown) as number)

  const getJobSpecRuns = React.useCallback(() => {
    if (isOcrJob) {
      return v2.ocrRuns
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
        .catch(setError)
    } else {
      return v2.runs
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
        .catch(setError)
    }
  }, [isOcrJob, jobSpecId, setError])

  const getJobSpec = React.useCallback(async () => {
    if (isOcrJob) {
      return v2.ocrSpecs
        .getJobSpec(jobSpecId)
        .then((response) => {
          const jobSpec = response.data
          setState((s) => ({
            ...s,
            jobSpec,
            job: {
              ...jobSpec.attributes.offChainReportingOracleSpec,
              id: jobSpec.id,
              errors: jobSpec.attributes.errors,
              definition: undefined,
              initiators: undefined,
              type: 'Off-chain-reporting',
            },
          }))
        })
        .catch(setError)
    } else {
      return v2.specs
        .getJobSpec(jobSpecId)
        .then((response) => {
          const jobSpec = response.data
          setState((s) => ({
            ...s,
            jobSpec,
            job: {
              ...jobSpec.attributes,
              id: jobSpec.id,
              definition: jobSpecDefinition({
                ...jobSpec,
                ...jobSpec.attributes,
              }),
              type: 'Direct request',
            },
          }))
        })
        .catch(setError)
    }
  }, [isOcrJob, jobSpecId, setError])

  React.useEffect(() => {
    getJobSpec()
  }, [getJobSpec])

  return (
    <div>
      <RegionalNav
        jobSpecId={jobSpecId}
        job={job}
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
