import React from 'react'

import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import CardHeader from '@material-ui/core/CardHeader'
import Grid from '@material-ui/core/Grid'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'

import Button from 'components/Button'
import BaseLink from 'components/BaseLink'
import { parseDot } from 'utils/parseDot'

import { JobRunsTable } from 'src/components/Table/JobRunsTable'
import TaskListDag from 'pages/Jobs/TaskListDag'

// ShowViewMoreCount defines the minimum number of jobs to display the
// View More link
const ShowViewMoreCount = 5

const chartCardStyles = ({ spacing }: Theme) =>
  createStyles({
    runDetails: {
      paddingTop: spacing.unit * 2,
      paddingBottom: spacing.unit * 2,
      paddingLeft: spacing.unit * 2,
    },
  })

interface Props extends WithStyles<typeof chartCardStyles> {
  job: JobPayload_Fields
}

export const TabOverview = withStyles(chartCardStyles)(
  ({ classes, job }: Props) => {
    // Convert the runs into run props which are compatible with the
    // JobRunsTable
    const runs = React.useMemo(() => {
      return job.runs.results.map(
        ({ allErrors, id, createdAt, finishedAt }) => ({
          id,
          createdAt,
          errors: allErrors,
          finishedAt,
          jobId: job.id,
        }),
      )
    }, [job.runs])

    return (
      <Grid container spacing={32}>
        <Grid item xs={12} sm={6}>
          <Card>
            <CardHeader title="Recent job runs" />

            <JobRunsTable runs={runs} />

            {job.runs.metadata.total > ShowViewMoreCount && (
              <div className={classes.runDetails}>
                <Button href={`/jobs/${job.id}/runs`} component={BaseLink}>
                  View more
                </Button>
              </div>
            )}
          </Card>
        </Grid>

        <Grid item xs={12} sm={6}>
          <Grid item xs>
            <Card style={{ overflow: 'visible' }}>
              <CardHeader title="Task list" />
              <TaskList observationSource={job.observationSource} />
            </Card>
          </Grid>
        </Grid>
      </Grid>
    )
  },
)

// TODO - Consider making a more generic function
const TaskList: React.FC<{ observationSource?: string }> = ({
  observationSource,
}) => {
  if (observationSource === undefined || observationSource === '') {
    return (
      <CardContent>
        <Typography align="center">No Task Graph Found</Typography>
      </CardContent>
    )
  }

  try {
    return <TaskListDag stratify={parseDot(`digraph {${observationSource}}`)} />
  } catch (error) {
    return (
      <CardContent>
        <Typography align="center">Failed to parse task graph</Typography>
      </CardContent>
    )
  }
}
