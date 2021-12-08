import React from 'react'

import { gql } from '@apollo/client'
import { Route, Switch, useRouteMatch } from 'react-router-dom'

import Grid from '@material-ui/core/Grid'

import Content from 'components/Content'
import { ErrorsCard } from './ErrorsCard'
import { Heading1 } from 'src/components/Heading/Heading1'
import { JobRunCard } from './JobRunCard'
import { JobRunTabs } from './JobRunTabs'
import { StatusCard } from './StatusCard'
import { TaskListCard } from 'src/components/Cards/TaskListCard'
import { TaskRunsCard } from './TaskRunsCard'
import { JSONCard } from './JSONCard'

export const JOB_RUN_PAYLOAD__TASK_RUNS_FIELDS = gql`
  fragment JobRunPayload_TaskRunsFields on TaskRun {
    id
    createdAt
    dotID
    error
    finishedAt
    output
    type
  }
`

export const JOB_RUN_PAYLOAD_FIELDS = gql`
  ${JOB_RUN_PAYLOAD__TASK_RUNS_FIELDS}
  fragment JobRunPayload_Fields on JobRun {
    id
    allErrors
    createdAt
    fatalErrors
    finishedAt
    inputs
    job {
      id
      name
      observationSource
    }
    outputs
    status
    taskRuns {
      ...JobRunPayload_TaskRunsFields
    }
  }
`

interface Props {
  run: JobRunPayload_Fields
}

export const JobRunView = ({ run }: Props) => {
  const { path } = useRouteMatch()

  return (
    <Content>
      <Grid container spacing={16}>
        <Grid item xs={12}>
          <Heading1>Job Run #{run.id}</Heading1>
        </Grid>
      </Grid>

      <JobRunCard run={run} />
      <JobRunTabs id={run.id} />

      <Grid container>
        <Grid item xs={4}>
          <Grid container>
            <Grid item xs={12}>
              <StatusCard
                status={run.status}
                startedAt={run.createdAt}
                finishedAt={run.finishedAt}
              />
            </Grid>

            <Grid item xs={12}>
              <TaskListCard observationSource={run.job.observationSource} />
            </Grid>
          </Grid>
        </Grid>

        <Grid item xs={8}>
          <Grid container>
            <Grid item xs={12}>
              {run.allErrors.length > 0 && (
                <ErrorsCard errors={run.allErrors} />
              )}
            </Grid>

            <Grid item xs={12}>
              <Switch>
                <Route path={`${path}/json`}>
                  <JSONCard run={run} />
                </Route>

                <Route path={path}>
                  {run.taskRuns.length > 0 && (
                    <TaskRunsCard
                      taskRuns={run.taskRuns}
                      observationSource={run.job.observationSource}
                    />
                  )}
                </Route>
              </Switch>
            </Grid>
          </Grid>
        </Grid>
      </Grid>
    </Content>
  )
}
