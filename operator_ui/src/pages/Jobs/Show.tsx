import React from 'react'
import { Route, Switch, useParams, useRouteMatch } from 'react-router-dom'
import { useErrorHandler } from 'hooks/useErrorHandler'
import { useLoadingPlaceholder } from 'hooks/useLoadingPlaceholder'
import { v2 } from 'api'
import {
  generateJSONDefinition,
  generateTOMLDefinition,
} from './generateJobSpecDefinition'
import { JobData, JobV2 } from './sharedTypes'
import { JobDefinition } from './JobDefinition'
import { JobsErrors } from './Errors'
import { RecentRuns } from './RecentRuns'
import { RegionalNav } from './RegionalNav'
import { Runs as JobRuns } from './Runs'
import { isJobV2 } from './utils'
import {
  transformDirectRequestJobRun,
  transformPipelineJobRun,
} from './transformJobRuns'

interface RouteParams {
  jobSpecId: string
}

const DEFAULT_PAGE = 1
const RECENT_RUNS_COUNT = 5

export const JobsShow = () => {
  const { path } = useRouteMatch()
  const { jobSpecId } = useParams<RouteParams>()
  const [state, setState] = React.useState<JobData>({
    recentRuns: [],
    recentRunsCount: 0,
  })
  const { job, jobSpec, externalJobID } = state
  const { error, ErrorComponent, setError } = useErrorHandler()
  const { LoadingPlaceholder } = useLoadingPlaceholder(!error && !jobSpec)

  const getJobSpecRuns = React.useCallback(
    ({ page = DEFAULT_PAGE, size = RECENT_RUNS_COUNT } = {}) => {
      const requestParams = {
        jobSpecId,
        page,
        size,
      }
      if (isJobV2(jobSpecId)) {
        return v2.ocrRuns
          .getJobSpecRuns(requestParams)
          .then((jobSpecRunsResponse) => {
            setState((s) => ({
              ...s,
              recentRuns: jobSpecRunsResponse.data.map(
                transformPipelineJobRun(jobSpecId),
              ),
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
              recentRuns: jobSpecRunsResponse.data.map(
                transformDirectRequestJobRun(jobSpecId),
              ),
              recentRunsCount: jobSpecRunsResponse.meta.count,
            }))
          })
          .catch(setError)
      }
    },
    [jobSpecId, setError],
  )

  const getJobSpec = React.useCallback(async () => {
    if (isJobV2(jobSpecId)) {
      return v2.jobs
        .getJobSpec(jobSpecId)
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
  }, [jobSpecId, setError])

  React.useEffect(() => {
    getJobSpec()
  }, [getJobSpec])

  return (
    <div>
      <RegionalNav
        jobSpecId={jobSpecId}
        externalJobID={externalJobID}
        job={job}
        getJobSpecRuns={getJobSpecRuns}
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
              ...state,
              ErrorComponent,
              LoadingPlaceholder,
              error,
              getJobSpec,
              setState,
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
              getJobSpecRuns,
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
              getJobSpecRuns,
            }}
          />
        </Route>
      </Switch>
    </div>
  )
}

export default JobsShow
