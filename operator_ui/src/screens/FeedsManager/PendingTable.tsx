import React from 'react'

import { WithStyles, withStyles } from '@material-ui/core/styles'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'

import Link from 'components/Link'
import { tableStyles } from 'components/Table'
import { TimeAgo } from 'src/components/TimeAgo'

interface Props extends WithStyles<typeof tableStyles> {
  proposals: FeedsManager_JobProposalsFields[]
}

// PendingTable renders a table for pending proposals.
export const PendingTable = withStyles(tableStyles)(
  ({ classes, proposals }: Props) => {
    return (
      <Table>
        <TableHead>
          <TableRow>
            <TableCell>ID</TableCell>
            <TableCell>Last Proposed</TableCell>
          </TableRow>
        </TableHead>

        <TableBody>
          {proposals?.map((proposal, idx) => (
            <TableRow key={idx} className={classes.row} hover>
              <TableCell className={classes.cell} component="th" scope="row">
                <Link
                  className={classes.link}
                  href={`/job_proposals/${proposal.id}`}
                >
                  {proposal.id}
                </Link>
              </TableCell>

              <TableCell>
                <TimeAgo tooltip>{proposal.latestSpec.createdAt}</TimeAgo>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    )
  },
)
