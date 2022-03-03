import React from 'react'

import Card from '@material-ui/core/Card'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'

import { TimeAgo } from 'components/TimeAgo'

interface Props {
  job: JobPayload_Fields
}

export const TabErrors: React.FC<Props> = ({ job }) => {
  const tableHeaders = ['Occurrences', 'Created', 'Last Seen', 'Message']

  return (
    <Card>
      <Table>
        <TableHead>
          <TableRow>
            {tableHeaders.map((header) => (
              <TableCell key={header}>{header}</TableCell>
            ))}
          </TableRow>
        </TableHead>
        <TableBody>
          {job.errors.length === 0 && (
            <TableRow>
              <TableCell component="th" scope="row" colSpan={5}>
                No errors
              </TableCell>
            </TableRow>
          )}

          {job.errors.map((err, idx) => (
            <TableRow key={idx}>
              <TableCell>{err.occurrences}</TableCell>
              <TableCell>
                <TimeAgo tooltip>{err.createdAt}</TimeAgo>
              </TableCell>
              <TableCell>
                <TimeAgo tooltip>{err.updatedAt}</TimeAgo>
              </TableCell>
              <TableCell>{err.description}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </Card>
  )
}
