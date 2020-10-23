import { CardTitle, KeyValueList } from '@chainlink/styleguide'
import {
  createStyles,
  Theme,
  Typography,
  WithStyles,
  withStyles,
  Card,
  Grid,
} from '@material-ui/core'
import Content from 'components/Content'
import JobRunsList from 'components/JobRuns/List'
import TaskList from 'components/Jobs/TaskList'
import React from 'react'
import { GWEI_PER_TOKEN } from 'utils/constants'
import formatMinPayment from 'utils/formatWeiAsset'
import { formatInitiators } from 'utils/jobSpecInitiators'
import { JobData } from './sharedTypes'

const totalLinkEarned = (job: NonNullable<JobData['jobSpec']>) => {
  const zero = '0.000000'
  const unformatted =
    job.attributes.earnings &&
    (job.attributes.earnings / GWEI_PER_TOKEN).toString()
  const formatted =
    unformatted &&
    (unformatted.length >= 3 ? unformatted : (unformatted + '.').padEnd(8, '0'))
  return formatted || zero
}

const chartCardStyles = (theme: Theme) =>
  createStyles({
    wrapper: {
      marginLeft: theme.spacing.unit * 3,
      marginTop: theme.spacing.unit * 2,
      marginBottom: theme.spacing.unit * 2,
    },
    paymentText: {
      color: theme.palette.secondary.main,
      fontWeight: 450,
    },
    earnedText: {
      color: theme.palette.text.secondary,
      fontSize: theme.spacing.unit * 2,
    },
  })

interface ChartProps extends WithStyles<typeof chartCardStyles> {
  jobSpec: NonNullable<JobData['jobSpec']>
}

const ChartArea = withStyles(chartCardStyles)(
  ({ classes, jobSpec }: ChartProps) => (
    <Card>
      <Grid item className={classes.wrapper}>
        <Typography className={classes.paymentText} variant="h5">
          Link Payment
        </Typography>
        <Typography className={classes.earnedText}>
          {totalLinkEarned(jobSpec)}
        </Typography>
      </Grid>
    </Card>
  ),
)

export const RecentRuns = ({
  ErrorComponent,
  LoadingPlaceholder,
  error,
  getJobSpecRuns,
  jobSpec,
  recentRuns,
  recentRunsCount,
  showJobRunsCount = 5,
}: {
  ErrorComponent: React.FC
  LoadingPlaceholder: React.FC
  error: unknown
  getJobSpecRuns: () => Promise<void>
  jobSpec?: JobData['jobSpec']
  recentRuns?: JobData['recentRuns']
  recentRunsCount: JobData['recentRunsCount']
  showJobRunsCount?: number
}) => {
  React.useEffect(() => {
    document.title =
      jobSpec && jobSpec.attributes.name
        ? `${jobSpec.attributes.name} | Job spec details`
        : 'Job spec details'
  }, [jobSpec])

  React.useEffect(() => {
    getJobSpecRuns()
  }, [getJobSpecRuns])

  return (
    <Content>
      <ErrorComponent />
      <LoadingPlaceholder />
      {!error && jobSpec && (
        <Grid container spacing={24}>
          <Grid item xs={8}>
            <Card>
              <CardTitle divider>Recent Job Runs</CardTitle>

              {recentRuns && (
                <JobRunsList
                  jobSpecId={jobSpec.id}
                  runs={recentRuns.map((jobRun) => ({
                    ...jobRun,
                    ...jobRun.attributes,
                  }))}
                  count={recentRunsCount}
                  showJobRunsCount={showJobRunsCount}
                />
              )}
            </Card>
          </Grid>
          <Grid item xs={4}>
            <Grid container direction="column">
              <Grid item xs>
                <ChartArea jobSpec={jobSpec} />
              </Grid>
              <Grid item xs>
                <Card>
                  <CardTitle divider>Task List</CardTitle>
                  <TaskList tasks={jobSpec.attributes.tasks} />
                </Card>
              </Grid>
              <Grid item xs>
                <KeyValueList
                  showHead={false}
                  entries={Object.entries({
                    runCount: recentRunsCount,
                    initiator: formatInitiators(jobSpec.attributes.initiators),
                    minimumPayment: `${
                      formatMinPayment(Number(jobSpec.attributes.minPayment)) ||
                      0
                    } Link`,
                  })}
                  titleize
                />
              </Grid>
            </Grid>
          </Grid>
        </Grid>
      )}
    </Content>
  )
}
