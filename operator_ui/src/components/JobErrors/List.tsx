import { TimeAgo } from '@chainlink/styleguide'
import Card from '@material-ui/core/Card'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import Typography from '@material-ui/core/Typography'
import React from 'react'
import { JobSpecError } from 'operator_ui'

const renderBody = (errors: JobSpecError[] | undefined) => {
  if (errors && errors.length === 0) {
    return (
      <TableRow>
        <TableCell component="th" scope="row" colSpan={3}>
          No errors
        </TableCell>
      </TableRow>
    )
  } else if (errors) {
    return errors.map(error => (
      <TableRow key={error.id}>
        <TableCell component="th" scope="row">
          {error.id}
        </TableCell>
        <TableCell>
          <Typography variant="body1">
            <TimeAgo tooltip>{error.createdAt.toString()}</TimeAgo>
          </Typography>
        </TableCell>
        <TableCell>
          <Typography variant="body1">{error.description}</Typography>
        </TableCell>
      </TableRow>
    ))
  }

  return (
    <TableRow>
      <TableCell component="th" scope="row" colSpan={3}>
        Loading...
      </TableCell>
    </TableRow>
  )
}

interface Props {
  errors?: JobSpecError[]
}

export const List: React.FC<Props> = ({ errors }) => {
  return (
    <Card>
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
                Message
              </Typography>
            </TableCell>
          </TableRow>
        </TableHead>
        <TableBody>{renderBody(errors)}</TableBody>
      </Table>
    </Card>
  )
}

export default List
