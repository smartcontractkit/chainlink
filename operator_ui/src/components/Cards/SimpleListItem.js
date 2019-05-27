import React from 'react'
import TableCell from '@material-ui/core/TableCell'
import TableRow from '@material-ui/core/TableRow'
import { makeStyles } from '@material-ui/styles'

const useStyles = makeStyles(theme => ({
  cell: {
    borderColor: theme.palette.divider,
    borderTop: `1px solid`,
    borderBottom: 'none',
    paddingTop: theme.spacing(2),
    paddingBottom: theme.spacing(2),
    paddingLeft: theme.spacing(2)
  }
}))

const SimpleListItem = ({ children }) => {
  const classes = useStyles()
  return (
    <TableRow>
      <TableCell scope="row" className={classes.cell}>
        {children}
      </TableCell>
    </TableRow>
  )
}

export default SimpleListItem
