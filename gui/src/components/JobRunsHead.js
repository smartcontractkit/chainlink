import React from 'react'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import TableCell from '@material-ui/core/TableCell'
import Typography from '@material-ui/core/Typography'

const JobRunsHead = () => (
  <TableHead>
    <TableRow>
      <TableCell>
        <Typography variant='body1' color='textSecondary'>ID</Typography>
      </TableCell>
      <TableCell>
        <Typography variant='body1' color='textSecondary'>Status</Typography>
      </TableCell>
      <TableCell>
        <Typography variant='body1' color='textSecondary'>Created</Typography>
      </TableCell>
      <TableCell>
        <Typography variant='body1' color='textSecondary'>Result</Typography>
      </TableCell>
      <TableCell>
        <Typography variant='body1' color='textSecondary'>Error</Typography>
      </TableCell>
    </TableRow>
  </TableHead>
)

export default JobRunsHead
