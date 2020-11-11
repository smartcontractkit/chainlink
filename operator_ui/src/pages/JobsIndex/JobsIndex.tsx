import React from 'react'
import Button from 'components/Button'
import { Title } from 'components/Title'
import Content from 'components/Content'
import BaseLink from 'components/BaseLink'
import { v2 } from 'api'
import * as models from 'core/store/models'
import {
  Grid,
  Card,
  CardContent,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  Typography,
  TextField,
} from '@material-ui/core'
import Link from 'components/Link'
import { useErrorHandler } from 'hooks/useErrorHandler'
import { useLoadingPlaceholder } from 'hooks/useLoadingPlaceholder'
import {
  createStyles,
  withStyles,
  WithStyles,
  Theme,
} from '@material-ui/core/styles'
import { DirectRequestRow } from './DirectRequestRow'
import { OcrJobRow } from './OcrJobRow'
import SearchIcon from '@material-ui/icons/Search'

enum JobSpecTypes {
  jobSpec = 'specs',
  ocrJobSpec = 'jobSpecV2s',
}

interface Job<T> {
  attributes: T
  id: string
  type: string
}

export type DirectRequest = Job<models.JobSpec>
export type OffChainReporting = Job<models.OcrJobSpec>
export type CombinedJobs = DirectRequest | OffChainReporting

function isDirectRequest(job: CombinedJobs): job is DirectRequest {
  return job.type === JobSpecTypes.jobSpec
}

function isOffChainReporting(job: CombinedJobs): job is OffChainReporting {
  return job.type === JobSpecTypes.ocrJobSpec
}

function getCreatedAt(job: CombinedJobs) {
  if (isDirectRequest(job)) {
    return job.attributes.createdAt
  } else if (isOffChainReporting(job)) {
    return job.attributes.offChainReportingOracleSpec.createdAt
  } else {
    return new Date().toString()
  }
}

const PAGE_SIZE = 1000 // We intentionally set this to a very high number to avoid pagination

async function getJobs() {
  return Promise.all([
    v2.specs.getJobSpecs(1, PAGE_SIZE),
    v2.ocrSpecs.getJobSpecs(),
  ]).then(([jobs, ocrJobs]) => {
    const combinedJobs = [...jobs.data, ...ocrJobs.data]
    const jobsByDate = combinedJobs.sort((a, b) => {
      const jobA = new Date(getCreatedAt(a)).getTime()
      const jobB = new Date(getCreatedAt(b)).getTime()
      return jobA > jobB ? -1 : 1
    })

    return jobsByDate
  })
}

const searchIncludes = (searchParam: string) => {
  const lowerCaseSearchParam = searchParam.toLowerCase()

  return (stringToSearch: string) => {
    return stringToSearch.toLowerCase().includes(lowerCaseSearchParam)
  }
}

export const simpleJobFilter = (search: string) => (job: CombinedJobs) => {
  const test = searchIncludes(search)

  if (isDirectRequest(job)) {
    return (
      test(job.id) ||
      test(job.attributes.name) ||
      job.attributes.initiators.some((initiator) => test(initiator.type)) ||
      test('direct request')
    )
  }

  if (isOffChainReporting(job)) {
    return (
      test(job.id) ||
      ([
        'contractAddress',
        'keyBundleID',
        'p2pPeerID',
        'transmitterAddress',
      ] as Array<
        keyof models.OcrJobSpec['offChainReportingOracleSpec']
      >).some((property) =>
        test(String(job.attributes.offChainReportingOracleSpec[property])),
      ) ||
      test('off-chain reporting')
    )
  }

  return false
}

const styles = (theme: Theme) =>
  createStyles({
    card: {
      padding: theme.spacing.unit,
      marginBottom: theme.spacing.unit * 3,
    },
    search: {
      marginBottom: theme.spacing.unit,
    },
  })

export const JobsIndex = ({
  classes,
}: {
  classes: WithStyles<typeof styles>['classes']
}) => {
  const [search, setSearch] = React.useState('')
  const [jobs, setJobs] = React.useState<CombinedJobs[]>()
  const { error, ErrorComponent, setError } = useErrorHandler()
  const { LoadingPlaceholder } = useLoadingPlaceholder(!error && !jobs)

  React.useEffect(() => {
    document.title = 'Jobs'
  }, [])

  const jobFilter = React.useMemo(() => simpleJobFilter(search), [search])

  React.useEffect(() => {
    getJobs().then(setJobs).catch(setError)
  }, [setError])

  return (
    <Content>
      <Grid container>
        <Grid item xs={9}>
          <Title>Jobs</Title>
        </Grid>
        <Grid item xs={3}>
          <Grid container justify="flex-end">
            <Grid item>
              <Button
                variant="secondary"
                component={BaseLink}
                href={'/jobs/new'}
              >
                New Job
              </Button>
            </Grid>
          </Grid>
        </Grid>

        <Grid item xs={12}>
          <ErrorComponent />
          <LoadingPlaceholder />
          {!error && jobs && (
            <Card className={classes.card}>
              <CardContent>
                <Grid
                  container
                  spacing={8}
                  alignItems="flex-end"
                  className={classes.search}
                >
                  <Grid item>
                    <SearchIcon />
                  </Grid>
                  <Grid item>
                    <TextField
                      label="Search"
                      value={search}
                      name="search"
                      onChange={(event) => setSearch(event.target.value)}
                    />
                  </Grid>
                </Grid>

                <Table>
                  <TableHead>
                    <TableRow>
                      <TableCell>
                        <Typography variant="body1" color="textSecondary">
                          ID
                        </Typography>
                      </TableCell>
                      <TableCell>
                        <Typography variant="body1" color="textSecondary">
                          Created
                        </Typography>
                      </TableCell>
                      <TableCell>
                        <Typography variant="body1" color="textSecondary">
                          Type
                        </Typography>
                      </TableCell>
                      <TableCell>
                        <Typography variant="body1" color="textSecondary">
                          Initiator
                        </Typography>
                      </TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {jobs && !jobs.length && (
                      <TableRow>
                        <TableCell component="th" scope="row" colSpan={3}>
                          You haven’t created any jobs yet. Create a new job{' '}
                          <Link href={`/jobs/new`}>here</Link>
                        </TableCell>
                      </TableRow>
                    )}
                    {jobs &&
                      jobs.filter(jobFilter).map((job: CombinedJobs) => {
                        if (isDirectRequest(job)) {
                          return <DirectRequestRow key={job.id} job={job} />
                        } else if (isOffChainReporting(job)) {
                          return <OcrJobRow key={job.id} job={job} />
                        } else {
                          return <TableRow>Unknown Job Type</TableRow>
                        }
                      })}
                  </TableBody>
                </Table>
              </CardContent>
            </Card>
          )}
        </Grid>
      </Grid>
    </Content>
  )
}

export default withStyles(styles)(JobsIndex)
