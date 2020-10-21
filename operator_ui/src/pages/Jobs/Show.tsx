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
import { v2 } from 'api'
import { Route, RouteComponentProps, Switch } from 'react-router-dom'
import Content from 'components/Content'
import JobRunsList from 'components/JobRuns/List'
import TaskList from 'components/Jobs/TaskList'
import React from 'react'
import { GWEI_PER_TOKEN } from 'utils/constants'
import formatMinPayment from 'utils/formatWeiAsset'
import { formatInitiators } from 'utils/jobSpecInitiators'
import { JobsDefinition } from './Definition'
import { JobsErrors } from './Errors'
import { RegionalNav } from './RegionalNav'
import { ApiResponse, PaginatedApiResponse } from '@chainlink/json-api-client'
import { JobSpec, JobRun } from 'core/store/models'
import { useErrorHandler } from 'hooks/useErrorHandler'
import { useLoadingPlaceholder } from 'hooks/useLoadingPlaceholder'

export type JobData = {
  jobSpec?: ApiResponse<JobSpec>['data']
  recentRuns?: PaginatedApiResponse<JobRun[]>['data']
  recentRunsCount: number
}

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

type Props = {
  showJobRunsCount: number
} & RouteComponentProps<{
  jobSpecId: string
}>

const DEFAULT_PAGE = 1
const RECENT_RUNS_COUNT = 5

export const JobsShow: React.FC<Props> = ({ match, showJobRunsCount = 5 }) => {
  const [state, setState] = React.useState<JobData>({
    recentRuns: [],
    recentRunsCount: 0,
  })
  const { jobSpec, recentRuns, recentRunsCount } = state
  const { error, ErrorComponent, setError } = useErrorHandler()
  const { LoadingPlaceholder } = useLoadingPlaceholder(!error && !jobSpec)

  const { jobSpecId } = match.params

  React.useEffect(() => {
    Promise.all([
      v2.specs.getJobSpec(jobSpecId),
      v2.runs.getJobSpecRuns({
        jobSpecId,
        page: DEFAULT_PAGE,
        size: RECENT_RUNS_COUNT,
      }),
    ])
      .then(([jobSpecResponse, jobSpecRunsResponse]) => {
        setState({
          jobSpec: jobSpecResponse.data,
          recentRuns: jobSpecRunsResponse.data,
          recentRunsCount: jobSpecRunsResponse.meta.count,
        })
      })
      .catch(setError)
  }, [jobSpecId, setError])

  return (
    <div>
      <RegionalNav jobSpecId={jobSpecId} job={jobSpec} />
      <Switch>
        <Route path={`${match.path}/json`} component={JobsDefinition} />
        <Route path={`${match.path}/errors`} component={JobsErrors} />
        <Route
          path={`${match.path}`}
          render={() => (
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
                            initiator: formatInitiators(
                              jobSpec.attributes.initiators,
                            ),
                            minimumPayment: `${
                              formatMinPayment(
                                Number(jobSpec.attributes.minPayment),
                              ) || 0
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
          )}
        />
      </Switch>
    </div>
  )
}

export default JobsShow
