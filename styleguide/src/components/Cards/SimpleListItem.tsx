import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import TableCell from '@material-ui/core/TableCell'
import TableRow from '@material-ui/core/TableRow'
import React from 'react'

const styles = (theme: Theme) =>
  createStyles({
    cell: {
      borderColor: theme.palette.divider,
      borderTop: `1px solid`,
      borderBottom: 'none',
      paddingTop: theme.spacing.unit * 2,
      paddingBottom: theme.spacing.unit * 2,
      paddingLeft: theme.spacing.unit * 2,
    },
  })

interface Props extends WithStyles<typeof styles> {
  children: React.ReactNode
}

const UnstyledSimpleListItem = ({ children, classes }: Props) => {
  return (
    <TableRow>
      <TableCell scope="row" className={classes.cell}>
        {children}
      </TableCell>
    </TableRow>
  )
}

export const SimpleListCardItem = withStyles(styles)(UnstyledSimpleListItem)
