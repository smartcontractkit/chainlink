import React from 'react'

import { gql } from '@apollo/client'
import { useHistory } from 'react-router-dom'

import Card from '@material-ui/core/Card'
import Grid from '@material-ui/core/Grid'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableHead from '@material-ui/core/TableHead'
import TablePagination from '@material-ui/core/TablePagination'
import TableRow from '@material-ui/core/TableRow'

import BaseLink from 'components/BaseLink'
import Button from 'components/Button'
import Content from 'components/Content'
import { JobRow } from './JobRow'
import Link from 'components/Link'
import { SearchTextField } from 'src/components/Search/SearchTextField'
import { Heading1 } from 'src/components/Heading/Heading1'

export const JOBS_PAYLOAD__RESULTS_FIELDS = gql`
  fragment JobsPayload_ResultsFields on Job {
    id
    name
    externalJobID
    createdAt
    spec {
      __typename
      ... on OCRSpec {
        contractAddress
        keyBundleID
        transmitterAddress
      }
    }
  }
`

const searchIncludes = (searchParam: string) => {
  const lowerCaseSearchParam = searchParam.toLowerCase()

  return (stringToSearch: string) => {
    return stringToSearch.toLowerCase().includes(lowerCaseSearchParam)
  }
}

export const simpleJobFilter =
  (search: string) => (job: JobsPayload_ResultsFields) => {
    if (search === '') {
      return true
    }

    return matchJob(job, search)
  }

/**
 * matchJob determines whether the job matches the search terms
 */
function matchJob(job: JobsPayload_ResultsFields, term: string) {
  const match = searchIncludes(term)

  const dataset: string[] = [
    job.id,
    job.name || '',
    job.spec.__typename,
    job.externalJobID,
  ]

  // Extend the search params for OCR fields that are not displayed in the table
  const spec = job.spec
  if (spec.__typename == 'OCRSpec') {
    const searchableProperties = [
      'contractAddress',
      'keyBundleID',
      'transmitterAddress',
    ] as Array<
      keyof Extract<
        JobsPayload_ResultsFields['spec'],
        { __typename: 'OCRSpec' }
      >
    >

    searchableProperties.forEach((property) => {
      dataset.push(String(spec[property]))
    })
  }

  return dataset.some(match)
}

export interface Props {
  jobs: ReadonlyArray<JobsPayload_ResultsFields>
  page: number
  pageSize: number
  total: number
}

export const JobsView: React.FC<Props> = ({ jobs, page, pageSize, total }) => {
  const history = useHistory()
  const [search, setSearch] = React.useState('')

  const jobFilter = React.useMemo(
    () => simpleJobFilter(search.trim()),
    [search],
  )

  return (
    <Content>
      <Grid container>
        <Grid item xs={9}>
          <Heading1>Jobs</Heading1>
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
          <SearchTextField
            value={search}
            onChange={setSearch}
            placeholder="Search jobs"
          />

          <Card>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>ID</TableCell>
                  <TableCell>Name</TableCell>
                  <TableCell>Type</TableCell>
                  <TableCell>External Job ID</TableCell>
                  <TableCell>Created</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {jobs.length === 0 && (
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
            <TablePagination
              component="div"
              count={total}
              rowsPerPage={pageSize}
              rowsPerPageOptions={[pageSize]}
              page={page - 1}
              onChangePage={(_, p) => {
                history.push(`/jobs?page=${p + 1}&per=${pageSize}`)
              }}
              onChangeRowsPerPage={() => {}} /* handler required by component, so make it a no-op */
              backIconButtonProps={{ 'aria-label': 'prev-page' }}
              nextIconButtonProps={{ 'aria-label': 'next-page' }}
            />
          </Card>
        </Grid>
      </Grid>
    </Content>
  )
}
