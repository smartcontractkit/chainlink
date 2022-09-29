import React from 'react'

import { gql } from '@apollo/client'

import Card from '@material-ui/core/Card'
import CardHeader from '@material-ui/core/CardHeader'
import { createStyles, withStyles, WithStyles } from '@material-ui/core/styles'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'

import { ErrorRow } from 'src/components/TableRow/ErrorRow'
import { LoadingRow } from 'src/components/TableRow/LoadingRow'
import { NoContentRow } from 'src/components/TableRow/NoContentRow'
import { RecentJobRow } from './RecentJobRow'

export const RECENT_JOBS_PAYLOAD__RESULTS_FIELDS = gql`
  fragment RecentJobsPayload_ResultsFields on Job {
    id
    name
    createdAt
  }
`

const styles = () =>
  createStyles({
    cardHeader: {
      borderBottom: 0,
    },
    table: {
      tableLayout: 'fixed',
    },
  })

export interface Props extends WithStyles<typeof styles> {
  data?: FetchRecentJobs
  loading: boolean
  errorMsg?: string
}

export const RecentJobsCard = withStyles(styles)(
  ({ classes, data, errorMsg, loading }: Props) => {
    return (
      <Card>
        <CardHeader title="Recent Jobs" className={classes.cardHeader} />

        <Table className={classes.table}>
          <TableBody>
            <LoadingRow visible={loading} />
            <NoContentRow visible={data?.jobs.results?.length === 0}>
              No recently created jobs
            </NoContentRow>
            <ErrorRow msg={errorMsg} />

            {data?.jobs.results?.map((job, idx) => (
              <RecentJobRow job={job} key={idx} />
            ))}
          </TableBody>
        </Table>
      </Card>
    )
  },
)
