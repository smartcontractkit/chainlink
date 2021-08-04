import React from 'react'

import { v2 } from 'api'
import BaseLink from 'components/BaseLink'
import Button from 'components/Button'
import Content from 'components/Content'
import { JobRow } from './JobRow'
import Link from 'components/Link'
import * as models from 'core/store/models'
import { Title } from 'components/Title'
import { useErrorHandler } from 'hooks/useErrorHandler'
import { useLoadingPlaceholder } from 'hooks/useLoadingPlaceholder'

import Grid from '@material-ui/core/Grid'
import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import SearchIcon from '@material-ui/icons/Search'
import {
  createStyles,
  withStyles,
  WithStyles,
  Theme,
} from '@material-ui/core/styles'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import Typography from '@material-ui/core/Typography'
import TextField from '@material-ui/core/TextField'

interface Job<T> {
  attributes: T
  id: string
  type: string
}

export type JobSpecV2 = Job<models.Job>

function isOCRJobSpec(job: JobSpecV2) {
  return job.attributes.type === 'offchainreporting'
}

function getCreatedAt(job: JobSpecV2) {
  switch (job.attributes.type) {
    case 'offchainreporting':
      return job.attributes.offChainReportingOracleSpec.createdAt

    case 'fluxmonitor':
      return job.attributes.fluxMonitorSpec.createdAt

    case 'directrequest':
      return job.attributes.directRequestSpec.createdAt

    case 'keeper':
      return job.attributes.keeperSpec.createdAt

    case 'cron':
      return job.attributes.cronSpec.createdAt

    case 'webhook':
      return job.attributes.webhookSpec.createdAt

    case 'vrf':
      return job.attributes.vrfSpec.createdAt
    default:
      return new Date().toString()
  }
}

async function getJobs() {
  return Promise.all([v2.jobs.getJobSpecs()]).then(([v2Jobs]) => {
    const jobsByDate = v2Jobs.data.sort((a, b) => {
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

export const simpleJobFilter = (search: string) => (job: JobSpecV2) => {
  if (search === '') {
    return true
  }

  if (isOCRJobSpec(job)) {
    return matchOCR(job, search)
  } else {
    return matchSimple(job, search)
  }
}

// matchSimple does a simple match on the id, name and type.
function matchSimple(job: JobSpecV2, term: string) {
  const match = searchIncludes(term)

  const dataset: string[] = [
    job.id,
    job.attributes.name || '',
    job.attributes.type,
  ]

  return dataset.some(match)
}

/**
 * matchOCR determines whether the OCR job matches the search terms
 *
 * @param job {JobSpecV2} The V2 Job Spec
 * @param term {string} The search term
 */
function matchOCR(job: JobSpecV2, term: string) {
  const match = searchIncludes(term)

  const { offChainReportingOracleSpec } = job.attributes

  const dataset: string[] = [
    job.id,
    job.attributes.name || '',
    job.attributes.type,
  ]

  const searchableProperties = [
    'contractAddress',
    'keyBundleID',
    'p2pPeerID',
    'transmitterAddress',
  ] as Array<keyof models.Job['offChainReportingOracleSpec']>

  if (offChainReportingOracleSpec) {
    searchableProperties.forEach((property) => {
      dataset.push(String(offChainReportingOracleSpec[property]))
    })
  }

  return dataset.some(match)
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
  const [jobs, setJobs] = React.useState<JobSpecV2[]>()
  const { error, ErrorComponent, setError } = useErrorHandler()
  const { LoadingPlaceholder } = useLoadingPlaceholder(!error && !jobs)

  React.useEffect(() => {
    document.title = 'Jobs'
  }, [])

  React.useEffect(() => {
    getJobs().then(setJobs).catch(setError)
  }, [setError])

  const jobFilter = React.useMemo(() => simpleJobFilter(search.trim()), [
    search,
  ])

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
                          Name
                        </Typography>
                      </TableCell>

                      <TableCell>
                        <Typography variant="body1" color="textSecondary">
                          Type
                        </Typography>
                      </TableCell>

                      <TableCell>
                        <Typography variant="body1" color="textSecondary">
                          Created
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

                    {jobs.filter(jobFilter).map((job) => (
                      <JobRow key={job.id} job={job} />
                    ))}
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
