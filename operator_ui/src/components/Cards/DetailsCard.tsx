import React from 'react'

import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'

import Paper from '@material-ui/core/Paper'

const styles = (theme: Theme) =>
  createStyles({
    paper: {
      display: 'flex',
      margin: `${theme.spacing.unit * 2.5}px 0`,
      padding: `${theme.spacing.unit * 3}px  ${theme.spacing.unit * 3.5}px`,
    },
    content: {
      flex: 1,
      width: '100%',
    },
    actions: {
      marginTop: -theme.spacing.unit * 1.5,
      marginLeft: -theme.spacing.unit * 4,
      [theme.breakpoints.up('sm')]: {
        marginLeft: 0,
        marginRight: -theme.spacing.unit * 1.5,
      },
    },
    itemBlock: {
      border: '1px solid rgba(224, 224, 224, 1)',
      borderRadius: theme.shape.borderRadius,
      padding: theme.spacing.unit * 2,
      marginTop: theme.spacing.unit,
    },
    itemBlockText: {
      overflowWrap: 'anywhere',
    },
  })

interface Props extends WithStyles<typeof styles> {
  actions?: React.ReactNode
}

// DetailsCard provides a box to present details about an object.
//
// This should be used in conjunction with a Grid and additonally provides space
// for actions which should contain a icon button that opens a menu.
//
// We may make this more specific to the use case in the future
export const DetailsCard = withStyles(styles)(
  ({ actions, children, classes }: React.PropsWithChildren<Props>) => {
    return (
      <Paper className={classes.paper}>
        <div className={classes.content}>{children}</div>

        {actions && <div className={classes.actions}>{actions}</div>}
      </Paper>
    )
  },
)

// DetailsCardItemTitle provides default styles for an item title in the details
// card
export const DetailsCardItemTitle = ({ title }: { title: string }) => (
  <Typography variant="subtitle2" gutterBottom>
    {title}
  </Typography>
)

// DetailsCardItemValue provides default styles for an item value in the details
// card.
export const DetailsCardItemValue: React.FC<{
  value?: string | number | null
}> = ({ children, value }) => (
  <Typography variant="body1" noWrap>
    {children ? children : value}
  </Typography>
)

// DetailsCardItemValue provides default styles for an item value in the details
// card.
interface DetailsCardItemBlockValueProps extends WithStyles<typeof styles> {
  value?: string | number | null
}

// DetailsCardItemValue provides default styles for an item value as a block
// item. Useful to highlight blocks of text values
export const DetailsCardItemBlockValue = withStyles(styles)(
  ({
    children,
    classes,
    value,
  }: React.PropsWithChildren<DetailsCardItemBlockValueProps>) => {
    return (
      <div className={classes.itemBlock}>
        <Typography variant="body1" className={classes.itemBlockText}>
          {children ? children : value}
        </Typography>
      </div>
    )
  },
)
