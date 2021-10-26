import { CardTitle } from 'components/CardTitle'
import { KeyValueList } from 'components/KeyValueList'
import {
  Card,
  Grid,
  Theme,
  Typography,
  WithStyles,
  createStyles,
  withStyles,
} from '@material-ui/core'
import Button from 'components/Button'
import BaseLink from 'components/BaseLink'
import Content from 'components/Content'
import JobRunsList from './JobRunsList'
import TaskListDag from './TaskListDag'
import TaskList from 'components/Jobs/TaskList'
import React from 'react'
import { GWEI_PER_TOKEN } from 'utils/constants'
import formatMinPayment from 'utils/formatWeiAsset'
import { formatInitiators } from 'utils/jobSpecInitiators'
import { DirectRequestJob, JobData } from './sharedTypes'
import { parseDot } from './parseDot'

const totalLinkEarned = (job: DirectRequestJob) => {
  const zero = '0.000000'
  const unformatted = job.earnings && (job.earnings / GWEI_PER_TOKEN).toString()
  const formatted =
    unformatted &&
    (unformatted.length >= 3 ? unformatted : (unformatted + '.').padEnd(8, '0'))
  return formatted || zero
}

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
  getJobSpecRuns: (props?: { page?: number; size?: number }) => Promise<void>
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
    getJobSpecRuns,
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
      getJobSpecRuns()
    }, [getJobSpecRuns])

    return (
      <Content>
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
              {job?.type === 'v2' && job.dotDagSource !== '' && (
                <Grid item xs>
                  <Card style={{ overflow: 'visible' }}>
                    <CardTitle divider>Task list</CardTitle>
                    <TaskListFunc job={job} />
                  </Card>
                </Grid>
              )}
              {job?.type === 'Direct request' && (
                <Grid container direction="column">
                  <Grid item xs>
                    <Card>
                      <Grid item className={classes.wrapper}>
                        <Typography
                          className={classes.paymentText}
                          variant="h5"
                        >
                          Link Payment
                        </Typography>
                        <Typography className={classes.earnedText}>
                          {totalLinkEarned(job)}
                        </Typography>
                      </Grid>
                    </Card>
                  </Grid>
                  <Grid item xs>
                    <Card>
                      <CardTitle divider>Task list</CardTitle>
                      <TaskList tasks={job.tasks} />
                    </Card>
                  </Grid>
                  <Grid item xs>
                    <KeyValueList
                      showHead={false}
                      entries={Object.entries({
                        initiator: formatInitiators(job.initiators),
                        minimumPayment: `${
                          formatMinPayment(Number(job.minPayment)) || 0
                        } Link`,
                      })}
                      titleize
                    />
                  </Grid>
                </Grid>
              )}
            </Grid>
          </Grid>
        )}
      </Content>
    )
  },
)

const TaskListFunc: React.FC<{ job?: JobData['job'] }> = ({ job }) => {
  if (job && (job as any).dotDagSource) {
    try {
      return (
        <TaskListDag
          stratify={parseDot(`digraph {${(job as any).dotDagSource}}`)}
        />
      )
    } catch (error) {
      console.error(error)
      return <p>Failed to parse task graph.</p>
    }
  }
  return <p>No task grapth found.</p>
}
