import React from 'react'
import { Route, Switch, useParams, useRouteMatch } from 'react-router-dom'
import { useErrorHandler } from 'hooks/useErrorHandler'
import { useLoadingPlaceholder } from 'hooks/useLoadingPlaceholder'
import { v2 } from 'api'
import { generateTOMLDefinition } from './generateJobSpecDefinition'
import { JobData, JobV2 } from './sharedTypes'
import { JobDefinition } from './JobDefinition'
import { JobsErrors } from './Errors'
import { RecentRuns } from './RecentRuns'
import { RegionalNav } from './RegionalNav'
import { Runs as JobRuns } from './Runs'
import { transformPipelineJobRun } from './transformJobRuns'

interface RouteParams {
  jobId: string
}

const DEFAULT_PAGE = 1
const RECENT_RUNS_COUNT = 5

export const JobsShow = () => {
  const { path } = useRouteMatch()
  const { jobId } = useParams<RouteParams>()
  const [state, setState] = React.useState<JobData>({
    recentRuns: [],
    recentRunsCount: 0,
  })
  const { job, externalJobID } = state
  const { error, ErrorComponent, setError } = useErrorHandler()
  const { LoadingPlaceholder } = useLoadingPlaceholder(!error && !job)

  const getJobRuns = React.useCallback(
    ({ page = DEFAULT_PAGE, size = RECENT_RUNS_COUNT } = {}) => {
      const requestParams = {
        jobId,
        page,
        size,
      }

      return v2.runs
        .getJobRuns(requestParams)
        .then((res) => {
          setState((s) => ({
            ...s,
            recentRuns: res.data.map(transformPipelineJobRun(jobId)),
            recentRunsCount: res.meta.count,
          }))
        })
        .catch(setError)
    },
    [jobId, setError],
  )

  const getJobSpec = React.useCallback(async () => {
    return v2.jobs
      .getJobSpec(jobId)
      .then((response) => {
        const jobSpec = response.data
        setState((s) => {
          let createdAt: string
          const externalJobID = jobSpec.attributes.externalJobID
          switch (jobSpec.attributes.type) {
            case 'offchainreporting':
              createdAt =
                jobSpec.attributes.offChainReportingOracleSpec.createdAt
              break
            case 'fluxmonitor':
              createdAt = jobSpec.attributes.fluxMonitorSpec.createdAt
              break
            case 'directrequest':
              createdAt = jobSpec.attributes.directRequestSpec.createdAt
              break
            case 'keeper':
              createdAt = jobSpec.attributes.keeperSpec.createdAt
              break
            case 'cron':
              createdAt = jobSpec.attributes.cronSpec.createdAt
              break
            case 'webhook':
              createdAt = jobSpec.attributes.webhookSpec.createdAt
              break
            case 'vrf':
              createdAt = jobSpec.attributes.vrfSpec.createdAt
              break
          }

          const job: JobV2 = {
            ...jobSpec.attributes.pipelineSpec,
            id: jobSpec.id,
            definition: generateTOMLDefinition(jobSpec.attributes),
            type: 'v2',
            name: jobSpec.attributes.name,
            specType: jobSpec.attributes.type,
            errors: jobSpec.attributes.errors,
            createdAt,
          }

          return {
            ...s,
            jobSpec,
            job,
            externalJobID,
          }
        })
      })
      .catch(setError)
  }, [jobId, setError])

  React.useEffect(() => {
    getJobSpec()
  }, [getJobSpec])

  return (
    <div>
      <RegionalNav
        jobId={jobId}
        externalJobID={externalJobID}
        job={job}
        getJobSpecRuns={getJobRuns}
        runsCount={state.recentRunsCount}
      />
      <Switch>
        <Route path={`${path}/definition`}>
          <JobDefinition
            {...{
              ...state,
              ErrorComponent,
              LoadingPlaceholder,
              error,
            }}
          />
        </Route>
        <Route exact path={`${path}/errors`}>
          <JobsErrors
            {...{
              job,
              ErrorComponent,
              LoadingPlaceholder,
              error,
            }}
          />
        </Route>
        <Route exact path={`${path}/runs`}>
          <JobRuns
            {...{
              ...state,
              error,
              ErrorComponent,
              LoadingPlaceholder,
              getJobRuns,
            }}
          />
        </Route>
        <Route path={path}>
          <RecentRuns
            {...{
              ...state,
              error,
              ErrorComponent,
              LoadingPlaceholder,
              getJobRuns,
            }}
          />
        </Route>
      </Switch>
    </div>
  )
}

export default JobsShow
