import React from 'react'

import { gql } from '@apollo/client'

import Card from '@material-ui/core/Card'
import CardHeader from '@material-ui/core/CardHeader'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableFooter from '@material-ui/core/TableFooter'
import TableRow from '@material-ui/core/TableRow'

import BaseLink from 'src/components/BaseLink'
import Button from 'src/components/Button'
import { ErrorRow } from 'src/components/TableRow/ErrorRow'
import { LoadingRow } from 'src/components/TableRow/LoadingRow'
import { NoContentRow } from 'src/components/TableRow/NoContentRow'
import { ActivityRow } from './ActivityRow'

export const RECENT_JOB_RUNS_PAYLOAD__RESULTS_FIELDS = gql`
  fragment RecentJobRunsPayload_ResultsFields on JobRun {
    id
    allErrors
    createdAt
    finishedAt
    status
    job {
      id
    }
  }
`

const styles = ({ palette, spacing }: Theme) =>
  createStyles({
    footer: {
      borderColor: palette.divider,
      borderTop: `1px solid`,
      paddingTop: spacing.unit * 2,
      paddingBottom: spacing.unit * 2,
    },
  })

export interface Props extends WithStyles<typeof styles> {
  data?: FetchRecentJobRuns
  loading: boolean
  errorMsg?: string
  maxRunsSize: number
}

export const ActivityCard = withStyles(styles)(
  ({ classes, data, loading, errorMsg, maxRunsSize }: Props) => {
    return (
      <Card>
        <CardHeader
          title="Activity"
          action={
            <Button href={'/jobs/new'} component={BaseLink}>
              New Job
            </Button>
          }
        />

        <Table>
          <TableBody>
            <LoadingRow visible={loading} />
            <NoContentRow visible={data?.jobRuns.results?.length === 0}>
              No recent activity
            </NoContentRow>
            <ErrorRow msg={errorMsg} />

            {data?.jobRuns.results?.map((run, idx) => (
              <ActivityRow run={run} key={idx} />
            ))}
          </TableBody>

          {data && data.jobRuns.metadata.total > maxRunsSize && (
            <TableFooter>
              <TableRow>
                <TableCell className={classes.footer}>
                  <Button href={'/runs'} component={BaseLink}>
                    View More
                  </Button>
                </TableCell>
              </TableRow>
            </TableFooter>
          )}
        </Table>
      </Card>
    )
  },
)
