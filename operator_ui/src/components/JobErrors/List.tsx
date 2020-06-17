import { localizedTimestamp, TimeAgo } from '@chainlink/styleguide'
import Card from '@material-ui/core/Card'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import Typography from '@material-ui/core/Typography'
import React from 'react'
import { JobSpecError } from 'operator_ui'
import Button from '../../components/Button'

const HEADERS = ['Occurances', 'Created', 'Last Seen', 'Message', 'Actions']
type DismissHandler = (id: string) => void

function renderBody(
  errors: JobSpecError[] | undefined,
  dismiss: DismissHandler,
) {
  if (errors && errors.length === 0) {
    return (
      <TableRow>
        <TableCell component="th" scope="row" colSpan={5}>
          No errors
        </TableCell>
      </TableRow>
    )
  } else if (errors) {
    return errors.map(error => (
      <TableRow key={error.id}>
        <TableCell>
          <Typography variant="body1">{error.occurances}</Typography>
        </TableCell>
        <TableCell>
          <Typography variant="body1">
            <TimeAgo tooltip>
              {localizedTimestamp(error.createdAt.toString())}
            </TimeAgo>
          </Typography>
        </TableCell>
        <TableCell>
          <Typography variant="body1">
            <TimeAgo tooltip>
              {localizedTimestamp(error.updatedAt.toString())}
            </TimeAgo>
          </Typography>
        </TableCell>
        <TableCell>
          <Typography variant="body1">{error.description}</Typography>
        </TableCell>
        <TableCell>
          <Button
            variant="danger"
            size="small"
            onClick={() => {
              dismiss(error.id)
            }}
          >
            Dismiss
          </Button>
        </TableCell>
      </TableRow>
    ))
  }

  return (
    <TableRow>
      <TableCell component="th" scope="row" colSpan={5}>
        Loading...
      </TableCell>
    </TableRow>
  )
}

function renderHeaders() {
  return HEADERS.map(header => (
    <TableCell key={header}>
      <Typography variant="body1" color="textSecondary">
        {header}
      </Typography>
    </TableCell>
  ))
}

interface Props {
  errors?: JobSpecError[]
  dismiss: DismissHandler
}

export const List: React.FC<Props> = ({ errors, dismiss }) => {
  return (
    <Card>
      <Table>
        <TableHead>
          <TableRow>{renderHeaders()}</TableRow>
        </TableHead>
        <TableBody>{renderBody(errors, dismiss)}</TableBody>
      </Table>
    </Card>
  )
}

export default List
