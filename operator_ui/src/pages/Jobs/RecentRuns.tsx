import { CardTitle } from 'components/CardTitle'
import {
  Card,
  Grid,
  Theme,
  WithStyles,
  createStyles,
  withStyles,
} from '@material-ui/core'
import Button from 'components/Button'
import BaseLink from 'components/BaseLink'
import JobRunsList from './JobRunsList'
import TaskListDag from './TaskListDag'
import React from 'react'
import { JobData } from './sharedTypes'
import { parseDot } from 'utils/parseDot'

const chartCardStyles = ({ spacing, palette }: Theme) =>
  createStyles({
    wrapper: {
      marginLeft: spacing.unit * 3,
      marginTop: spacing.unit * 2,
      marginBottom: spacing.unit * 2,
    },
    paymentText: {
      color: palette.secondary.main,
      fontWeight: 450,
    },
    earnedText: {
      color: palette.text.secondary,
      fontSize: spacing.unit * 2,
    },
    runDetails: {
      paddingTop: spacing.unit * 2,
      paddingBottom: spacing.unit * 2,
      paddingLeft: spacing.unit * 2,
    },
  })

interface Props extends WithStyles<typeof chartCardStyles> {
  ErrorComponent: React.FC
  LoadingPlaceholder: React.FC
  error: unknown
  getJobRuns: (props?: { page?: number; size?: number }) => Promise<void>
  job?: JobData['job']
  recentRuns?: JobData['recentRuns']
  recentRunsCount: JobData['recentRunsCount']
  showJobRunsCount?: number
}

export const RecentRuns = withStyles(chartCardStyles)(
  ({
    classes,
    ErrorComponent,
    LoadingPlaceholder,
    error,
    getJobRuns,
    job,
    recentRuns,
    recentRunsCount,
    showJobRunsCount = 5,
  }: Props) => {
    React.useEffect(() => {
      document.title = job?.name
        ? `${job.name} | Job spec details`
        : 'Job spec details'
    }, [job])

    React.useEffect(() => {
      getJobRuns()
    }, [getJobRuns])

    return (
      <>
        <ErrorComponent />
        <LoadingPlaceholder />

        {!error && job && (
          <Grid container spacing={40}>
            <Grid item xs={8}>
              <Card>
                <CardTitle divider>Recent job runs</CardTitle>

                {recentRuns && (
                  <>
                    <JobRunsList runs={recentRuns} />
                    {recentRunsCount > showJobRunsCount && (
                      <div className={classes.runDetails}>
                        <Button
                          href={`/jobs/${job.id}/runs`}
                          component={BaseLink}
                        >
                          View more
                        </Button>
                      </div>
                    )}
                  </>
                )}
              </Card>
            </Grid>

            <Grid item xs={4}>
              {job.dotDagSource !== '' && (
                <Grid item xs>
                  <Card style={{ overflow: 'visible' }}>
                    <CardTitle divider>Task list</CardTitle>
                    <TaskList job={job} />
                  </Card>
                </Grid>
              )}
            </Grid>
          </Grid>
        )}
      </>
    )
  },
)

const TaskList: React.FC<{ job?: JobData['job'] }> = ({ job }) => {
  if (job) {
    try {
      return (
        <TaskListDag stratify={parseDot(`digraph {${job.dotDagSource}}`)} />
      )
    } catch (error) {
      return <p>Failed to parse task graph.</p>
    }
  }
  return <p>No task grapth found.</p>
}
