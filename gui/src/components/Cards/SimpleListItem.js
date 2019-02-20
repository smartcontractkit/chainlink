import React from 'react'
import { withStyles } from '@material-ui/core/styles'
import TableCell from '@material-ui/core/TableCell'
import TableRow from '@material-ui/core/TableRow'

const styles = theme => ({
  cell: {
    borderColor: theme.palette.divider,
    borderTop: `1px solid`,
    borderBottom: 'none',
    paddingTop: theme.spacing.unit * 2,
    paddingBottom: theme.spacing.unit * 2,
    paddingLeft: theme.spacing.unit * 2
  }
})

const SimpleListItem = ({ children, classes }) => {
  return (
    <TableRow>
      <TableCell scope="row" className={classes.cell}>
        {children}
      </TableCell>
    </TableRow>
  )
}

export default withStyles(styles)(SimpleListItem)
