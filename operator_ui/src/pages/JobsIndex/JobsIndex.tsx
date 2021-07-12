import React, { useEffect } from 'react'
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
import { JobV2Row } from './JobV2Row'
import SearchIcon from '@material-ui/icons/Search'

enum JobSpecTypes {
  v1 = 'specs',
  v2 = 'jobs',
}

interface Job<T> {
  attributes: T
  id: string
  type: string
}

export type DirectRequest = Job<models.JobSpec>
export type JobSpecV2 = Job<models.JobSpecV2>
export type CombinedJobs = DirectRequest | JobSpecV2

function isJobSpecV1(job: CombinedJobs): job is DirectRequest {
  return job.type === JobSpecTypes.v1
}

function isJobSpecV2(job: CombinedJobs): job is JobSpecV2 {
  return job.type === JobSpecTypes.v2
}

function isDirectRequestJobSpecV2(job: JobSpecV2) {
  return job.attributes.type === 'directrequest'
}

function isFluxMonitorJobSpecV2(job: JobSpecV2) {
  return job.attributes.type === 'fluxmonitor'
}

function isKeeperSpecV2(job: JobSpecV2) {
  return job.attributes.type === 'keeper'
}

function isOCRJobSpecV2(job: JobSpecV2) {
  return job.attributes.type === 'offchainreporting'
}

function isCronSpecV2(job: JobSpecV2) {
  return job.attributes.type === 'cron'
}

function isWebhookSpecV2(job: JobSpecV2) {
  return job.attributes.type === 'webhook'
}

function getCreatedAt(job: CombinedJobs) {
  if (isJobSpecV1(job)) {
    return job.attributes.createdAt
  } else if (isJobSpecV2(job)) {
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
    }
  } else {
    return new Date().toString()
  }
}

const PAGE_SIZE = 1000 // We intentionally set this to a very high number to avoid pagination

async function getJobs() {
  return Promise.all([
    v2.specs.getJobSpecs(1, PAGE_SIZE),
    v2.jobs.getJobSpecs(),
  ]).then(([v1Jobs, v2Jobs]) => {
    const combinedJobs = [...v1Jobs.data, ...v2Jobs.data]
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
  if (search === '') {
    return true
  }

  if (isJobSpecV1(job)) {
    return matchV1Job(job, search)
  }

  if (isJobSpecV2(job)) {
    if (isDirectRequestJobSpecV2(job)) {
      return matchDirectRequest(job, search)
    }

    if (isFluxMonitorJobSpecV2(job)) {
      return matchFluxMonitor(job, search)
    }

    if (isOCRJobSpecV2(job)) {
      return matchOCR(job, search)
    }

    if (isKeeperSpecV2(job)) {
      return matchKeeper(job, search)
    }

    if (isCronSpecV2(job)) {
      return matchCron(job, search)
    }

    if (isWebhookSpecV2(job)) {
      return matchWebhook(job, search)
    }
  }

  return false
}

/**
 * matchV1Job determines whether the V1 job matches the search term
 *
 * @param job {DirectRequest} The V1 Job Spec
 * @param term {string}
 */
function matchV1Job(job: DirectRequest, term: string) {
  const match = searchIncludes(term)

  const dataset: string[] = [
    job.id,
    job.attributes.name,
    ...job.attributes.initiators.map((i) => i.type), // Match any of the initiators
    'direct request',
  ]

  return dataset.some(match)
}

/**
 * matchDirectRequest determines whether the V2 Direct Request job matches the search
 * terms.
 *
 * @param job {JobSpecV2} The V2 Job Spec
 * @param term {string} The search term
 */
function matchDirectRequest(job: JobSpecV2, term: string) {
  const match = searchIncludes(term)

  const dataset: string[] = [
    job.id,
    job.attributes.name || '',
    'direct request', // Hardcoded to match the type column
  ]

  return dataset.some(match)
}

/**
 * matchFluxMonitor determines whether the Flux Monitor job matches the search
 * terms.
 *
 * @param job {JobSpecV2} The V2 Job Spec
 * @param term {string} The search term
 */
function matchFluxMonitor(job: JobSpecV2, term: string) {
  const match = searchIncludes(term)

  const dataset: string[] = [
    job.id,
    job.attributes.name || '',
    'direct request', // Hardcoded to match the type column
    'fluxmonitor', // Hardcoded to match initiator column
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
    'off-chain reporting',
  ]

  const searchableProperties = [
    'contractAddress',
    'keyBundleID',
    'p2pPeerID',
    'transmitterAddress',
  ] as Array<keyof models.JobSpecV2['offChainReportingOracleSpec']>

  if (offChainReportingOracleSpec) {
    searchableProperties.forEach((property) => {
      dataset.push(String(offChainReportingOracleSpec[property]))
    })
  }

  return dataset.some(match)
}

/**
 * matchKeeper determines whether the Keeper job matches the search terms
 *
 * @param job {JobSpecV2} The V2 Job Spec
 * @param term {string} The search term
 */
function matchKeeper(job: JobSpecV2, term: string) {
  const match = searchIncludes(term)

  const dataset: string[] = [
    job.id,
    job.attributes.name || '',
    'keeper', // Hardcoded to match the type column
  ]

  return dataset.some(match)
}

/**
 * matchCron determines whether the Cron job matches the search terms
 *
 * @param job {JobSpecV2} The V2 Job Spec
 * @param term {string} The search term
 */
function matchCron(job: JobSpecV2, term: string) {
  const match = searchIncludes(term)

  const dataset: string[] = [job.id, job.attributes.name || '', 'cron']

  return dataset.some(match)
}

/**
 * matchWebhook determines whether the Webhook job matches the search terms
 *
 * @param job {JobSpecV2} The V2 Job Spec
 * @param term {string} The search term
 */
function matchWebhook(job: JobSpecV2, term: string) {
  const match = searchIncludes(term)

  const dataset: string[] = [job.id, job.attributes.name || '', 'webhook']

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
  const [jobs, setJobs] = React.useState<CombinedJobs[]>()
  const { error, ErrorComponent, setError } = useErrorHandler()
  const { LoadingPlaceholder } = useLoadingPlaceholder(!error && !jobs)

  React.useEffect(() => {
    document.title = 'Jobs'
  }, [])

  const jobFilter = React.useMemo(() => simpleJobFilter(search.trim()), [
    search,
  ])

  useEffect(() => {
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
                        if (isJobSpecV1(job)) {
                          return <DirectRequestRow key={job.id} job={job} />
                        } else if (isJobSpecV2(job)) {
                          return <JobV2Row key={job.id} job={job} />
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
