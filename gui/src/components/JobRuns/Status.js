import React from 'react'
import { withStyles } from '@material-ui/core/styles'
import Icon from '@material-ui/core/Icon'

const styles = theme => {
  const iconShared = {
    borderRadius: theme.spacing.unit * 3,
    padding: 3,
    width: '1.22em',
    height: '1.25em'
  }

  return {
    waiting: Object.assign(
      {},
      iconShared,
      {backgroundColor: theme.palette.warning.main},
      {color: theme.palette.warning.contrastText}
    ),
    completed: Object.assign(
      {},
      iconShared,
      {backgroundColor: theme.palette.success.light},
      {color: theme.palette.success.main}
    ),
    errored: Object.assign(
      {},
      iconShared,
      {backgroundColor: theme.palette.error.light},
      {color: theme.palette.error.main}
    )
  }
}

const Status = ({children, classes}) => {
  if (children === 'completed') {
    return (
      <Icon className={classes.completed} title={children}>
        done
      </Icon>
    )
  } else if (children === 'errored') {
    return (
      <Icon className={classes.errored} title={children}>
        error_outline
      </Icon>
    )
  }

  return (
    <Icon className={classes.waiting} title='waiting'>
      access_time
    </Icon>
  )
}

export default withStyles(styles)(Status)
