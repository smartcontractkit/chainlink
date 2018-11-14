import React from 'react'
import { withStyles } from '@material-ui/core/styles'
import Icon from '@material-ui/core/Icon'

const styles = theme => {
  return {
    waiting: {
      backgroundColor: theme.palette.warning.main,
      color: theme.palette.warning.contrastText,
      borderRadius: theme.spacing.unit * 3,
      padding: 3,
      width: '1.22em',
      height: '1.25em'
    },
    completed: {
      backgroundColor: theme.palette.success.light,
      color: theme.palette.success.main,
      borderRadius: theme.spacing.unit * 3,
      padding: 3,
      width: '1.22em',
      height: '1.25em'
    },
    errored: {
      backgroundColor: theme.palette.error.light,
      color: theme.palette.error.main,
      borderRadius: theme.spacing.unit * 3,
      padding: 3,
      width: '1.22em',
      height: '1.25em'
    }
  }
}

const Status = ({children, classes}) => {
  if (children === 'completed') {
    return (
      <Icon className={classes.completed} title='completed'>
        done
      </Icon>
    )
  } else if (children === 'errored') {
    return (
      <Icon className={classes.errored} title='error'>
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
