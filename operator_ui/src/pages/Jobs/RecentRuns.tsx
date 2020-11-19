import { CardTitle, KeyValueList } from '@chainlink/styleguide'
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
import JobRunsList from 'components/JobRuns/List'
import TaskList from 'components/Jobs/TaskList'
import React from 'react'
import { GWEI_PER_TOKEN } from 'utils/constants'
import formatMinPayment from 'utils/formatWeiAsset'
import { formatInitiators } from 'utils/jobSpecInitiators'
import { DirectRequestJob, JobData } from './sharedTypes'

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
  jobSpec?: JobData['jobSpec']
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
    jobSpec,
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
          <Grid container spacing={24}>
            <Grid item xs={8}>
              <Card>
                <CardTitle divider>Recent job runs</CardTitle>

                {recentRuns && (
                  <>
                    <JobRunsList
                      runs={recentRuns}
                      hideLinks={job?.type === 'Off-chain reporting'}
                    />
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
              {job?.type === 'Direct request' && jobSpec && (
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
                      <CardTitle divider>Task List</CardTitle>
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
