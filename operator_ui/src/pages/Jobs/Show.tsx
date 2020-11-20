import React from 'react'
import { v2 } from 'api'
import { Route, RouteComponentProps, Switch } from 'react-router-dom'
import { useErrorHandler } from 'hooks/useErrorHandler'
import { useLoadingPlaceholder } from 'hooks/useLoadingPlaceholder'
import {
  generateJSONDefinition,
  generateTOMLDefinition,
} from './generateJobSpecDefinition'
import { PaginatedApiResponse } from '@chainlink/json-api-client'
import { OcrJobRun, RunStatus } from 'core/store/models'
import { JobData } from './sharedTypes'
import { JobsDefinition } from './Definition'
import { JobsErrors } from './Errors'
import { RecentRuns } from './RecentRuns'
import { RegionalNav } from './RegionalNav'
import { Runs as JobRuns } from './Runs'

type Props = RouteComponentProps<{
  jobSpecId: string
}>

const DEFAULT_PAGE = 1
const RECENT_RUNS_COUNT = 5

function getOcrJobStatus({
  attributes: { finishedAt, errors },
}: NonNullable<PaginatedApiResponse<OcrJobRun[]>>['data'][0]) {
  if (finishedAt === null) {
    return RunStatus.IN_PROGRESS
  }
  if (errors[0] !== null) {
    return RunStatus.ERRORED
  }
  return RunStatus.COMPLETED
}

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

  const getJobSpecRuns = React.useCallback(
    ({ page = DEFAULT_PAGE, size = RECENT_RUNS_COUNT } = {}) => {
      const requestParams = {
        jobSpecId,
        page,
        size,
      }
      if (isOcrJob) {
        return v2.ocrRuns
          .getJobSpecRuns(requestParams)
          .then((jobSpecRunsResponse) => {
            setState((s) => ({
              ...s,
              recentRuns: jobSpecRunsResponse.data.map((jobRun) => ({
                createdAt: jobRun.attributes.createdAt,
                id: jobRun.id,
                status: getOcrJobStatus(jobRun),
                jobId: jobSpecId,
              })),
              recentRunsCount: jobSpecRunsResponse.meta.count,
            }))
          })
          .catch(setError)
      } else {
        return v2.runs
          .getJobSpecRuns(requestParams)
          .then((jobSpecRunsResponse) => {
            setState((s) => ({
              ...s,
              recentRuns: jobSpecRunsResponse.data.map((jobRun) => ({
                createdAt: jobRun.attributes.createdAt,
                id: jobRun.id,
                status: jobRun.attributes.status,
                jobId: jobSpecId,
              })),
              recentRunsCount: jobSpecRunsResponse.meta.count,
            }))
          })
          .catch(setError)
      }
    },
    [isOcrJob, jobSpecId, setError],
  )

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
              definition: generateTOMLDefinition(jobSpec.attributes),
              type: 'Off-chain reporting',
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
              type: 'Direct request',
              definition: generateJSONDefinition(jobSpec.attributes),
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
        runsCount={state.recentRunsCount}
      />
      <Switch>
        <Route
          path={`${match.path}/definition`}
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
          path={`${match.path}/runs`}
          render={() => (
            <JobRuns
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
