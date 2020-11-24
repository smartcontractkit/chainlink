import React from 'react'
import { v2 } from 'api'
import { Route, RouteComponentProps, Switch } from 'react-router-dom'
import { CardTitle } from '@chainlink/styleguide'
import { Card, Grid, Typography } from '@material-ui/core'
import Content from 'components/Content'
import { RunStatus } from 'core/store/models'
import StatusCard from './StatusCard'
import { useErrorHandler } from 'hooks/useErrorHandler'
import { useLoadingPlaceholder } from 'hooks/useLoadingPlaceholder'
import RegionalNav from './RegionalNav'
import { DirectRequestJobRun, OffChainReportingJobRun } from '../sharedTypes'
import { Overview } from './Overview/Overview'
import { Json } from './Json'
import { getOcrJobStatus, isOcrJob } from '../utils'
import { TaskList } from '../TaskListDag'

function getErrorsList(
  jobRun: DirectRequestJobRun | OffChainReportingJobRun,
): string[] {
  if (jobRun.type === 'Direct request job run') {
    return jobRun.taskRuns
      .filter(({ status }) => status === 'errored')
      .map((tr) => tr.result.error)
      .filter((error): error is string => error !== null)
  }

  if (jobRun.type === 'Off-chain reporting job run') {
    return jobRun.errors.filter((error): error is string => error !== null)
  }

  return []
}

function getTaskStatus({
  finishedAt,
  error,
}: {
  finishedAt: string
  error: string
}) {
  if (finishedAt === null) {
    return RunStatus.IN_PROGRESS
  }
  if (error !== null) {
    return RunStatus.ERRORED
  }
  return RunStatus.COMPLETED
}

type Props = RouteComponentProps<{
  jobSpecId: string
  jobRunId: string
}>

export const Show = ({ match }: Props) => {
  const [jobRun, setState] = React.useState<
    DirectRequestJobRun | OffChainReportingJobRun
  >()

  const { jobSpecId, jobRunId } = match.params

  const { error, ErrorComponent, setError } = useErrorHandler()
  const { LoadingPlaceholder } = useLoadingPlaceholder(!error && !jobRun)

  React.useEffect(() => {
    document.title = 'Job run details'
  }, [])

  const getJobRun = React.useCallback(async () => {
    if (isOcrJob(jobSpecId)) {
      v2.ocrRuns
        .getJobSpecRun({ jobSpecId, runId: jobRunId })
        .then((jobSpecRunResponse) => {
          const run = jobSpecRunResponse.data
          setState({
            ...run.attributes,
            id: run.id,
            jobId: jobSpecId,
            status: getOcrJobStatus(run.attributes),
            taskRuns: run.attributes.taskRuns.map((taskRun) => ({
              ...taskRun,
              status: getTaskStatus(taskRun),
            })),
            type: 'Off-chain reporting job run',
          })
        })
        .catch(setError)
    } else {
      return v2.runs
        .getJobSpecRun(jobRunId)
        .then((jobSpecRunResponse) => {
          const run = jobSpecRunResponse.data
          setState({
            ...run.attributes,
            id: run.id,
            jobId: jobSpecId,
            type: 'Direct request job run',
          })
        })
        .catch(setError)
    }
  }, [jobRunId, jobSpecId, setError])

  React.useEffect(() => {
    getJobRun()
  }, [getJobRun])

  const errorsList = jobRun ? getErrorsList(jobRun) : []

  return (
    <>
      <RegionalNav {...match.params} jobRun={jobRun} />
      <Content>
        <ErrorComponent />
        <LoadingPlaceholder />
        {jobRun && (
          <Grid container spacing={40}>
            <Grid item xs={8}>
              <Grid container spacing={40}>
                {errorsList.length > 0 && (
                  <Grid item xs={12}>
                    <Card>
                      <CardTitle divider>Errors</CardTitle>
                      <ul>
                        {errorsList.map((error, index) => (
                          <li key={error + index}>
                            <Typography variant="body1">{error}</Typography>
                          </li>
                        ))}
                      </ul>
                    </Card>
                  </Grid>
                )}
                <Grid item xs={12}>
                  <Switch>
                    <Route
                      path={`${match.path}/json`}
                      render={() => <Json jobRun={jobRun} />}
                    />
                    {jobRun.type === 'Direct request job run' && (
                      <Route render={() => <Overview jobRun={jobRun} />} />
                    )}
                  </Switch>
                </Grid>
              </Grid>
            </Grid>
            <Grid item xs={4}>
              <Grid container spacing={40}>
                <Grid item xs={12}>
                  <StatusCard {...jobRun} title={jobRun.status} />
                </Grid>
                {jobRun.type === 'Off-chain reporting job run' && (
                  <Grid item xs={12}>
                    <Card style={{ overflow: 'visible' }}>
                      <CardTitle divider>Task list</CardTitle>
                      <TaskList dotSource={jobRun.pipelineSpec.DotDagSource} />
                    </Card>
                  </Grid>
                )}
              </Grid>
            </Grid>
          </Grid>
        )}
      </Content>
    </>
  )
}

export default Show
