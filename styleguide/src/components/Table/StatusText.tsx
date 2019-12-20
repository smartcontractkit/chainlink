import { createStyles, withStyles, WithStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import classNames from 'classnames'
import React from 'react'

const styles = (theme: any) =>
  createStyles({
    statusCell: {
      textAlign: 'end',
      width: '30%',
    },
    status: {
      paddingLeft: theme.spacing.unit * 1.5,
      paddingRight: theme.spacing.unit * 1.5,
      paddingTop: theme.spacing.unit / 2,
      paddingBottom: theme.spacing.unit / 2,
      borderRadius: theme.spacing.unit * 2,
      marginRight: theme.spacing.unit,
      width: 'fit-content',
      display: 'inline-block',
    },
    errored: {
      backgroundColor: theme.palette.error.light,
      color: theme.palette.error.main,
    },
    pending: {
      backgroundColor: theme.palette.listPendingStatus.background,
      color: theme.palette.listPendingStatus.color,
    },
    completed: {
      backgroundColor: theme.palette.listCompletedStatus.background,
      color: theme.palette.listCompletedStatus.color,
    },
  })

const classFromStatus = (classes: any, status: string) => {
  if (
    !status ||
    status.startsWith('pending') ||
    status.startsWith('in_progress')
  ) {
    return classes['pending']
  }
  return classes[status.toLowerCase()]
}

// FIXME duplicated from operator-ui
const titleize = (input: string) => {
  const normalized = input || ''
  return normalized
    .toLowerCase()
    .replace(/_/g, ' ')
    .replace(/(?:^|\s|-)\S/g, x => x.toUpperCase())
}

interface Props extends WithStyles<typeof styles> {
  status: string
}
const UnstyledStatusText = (props: Props) => (
  <Typography
    variant="body1"
    className={classNames(
      props.classes.status,
      classFromStatus(props.classes, props.status),
    )}
  >
    {titleize(props.status)}
  </Typography>
)

export const StatusText = withStyles(styles)(UnstyledStatusText)
