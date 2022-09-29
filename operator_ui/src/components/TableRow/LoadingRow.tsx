import React from 'react'

import CircularProgress from '@material-ui/core/CircularProgress'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import TableCell from '@material-ui/core/TableCell'
import TableRow from '@material-ui/core/TableRow'

const styles = (theme: Theme) =>
  createStyles({
    cell: {
      padding: theme.spacing.unit * 2,
    },
  })

interface Props extends WithStyles<typeof styles> {
  visible: boolean
}

export const LoadingRow = withStyles(styles)(({ classes, visible }: Props) => {
  if (!visible) {
    return null
  }

  return (
    <TableRow>
      {/* Sets a high column count to insure this is always centered regardless
        of the number of columns in the table.
         */}
      <TableCell colSpan={100} align="center" className={classes.cell}>
        <CircularProgress data-testid="loading-spinner" size={24} />
      </TableCell>
    </TableRow>
  )
})
