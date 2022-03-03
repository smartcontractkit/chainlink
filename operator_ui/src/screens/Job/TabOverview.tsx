import React from 'react'

import Card from '@material-ui/core/Card'
import CardHeader from '@material-ui/core/CardHeader'
import Grid from '@material-ui/core/Grid'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'

import Button from 'components/Button'
import BaseLink from 'components/BaseLink'
import { JobRunsTable } from 'components/Table/JobRunsTable'
import { TaskListCard } from 'components/Cards/TaskListCard'

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
        ({ allErrors, id, createdAt, finishedAt, status }) => ({
          id,
          createdAt,
          errors: allErrors,
          finishedAt,
          status,
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
          <TaskListCard observationSource={job.observationSource} />
        </Grid>
      </Grid>
    )
  },
)
