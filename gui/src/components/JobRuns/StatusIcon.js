import React from 'react'
import PropTypes from 'prop-types'
import { withStyles } from '@material-ui/core/styles'
import Icon from '@material-ui/core/Icon'
import classNames from 'classnames'

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
      { backgroundColor: theme.palette.warning.main },
      { color: theme.palette.warning.contrastText }
    ),
    completed: Object.assign(
      {},
      iconShared,
      { backgroundColor: theme.palette.success.light },
      { color: theme.palette.success.main }
    ),
    errored: Object.assign(
      {},
      iconShared,
      { backgroundColor: theme.palette.error.light },
      { color: theme.palette.error.main }
    )
  }
}

const StatusIcon = ({ children, classes, className }) => {
  if (children === 'completed') {
    return (
      <Icon className={classNames(classes.completed, className)} title={children}>
        done
      </Icon>
    )
  } else if (children === 'errored') {
    return (
      <Icon className={classNames(classes.errored, className)} title={children}>
        error_outline
      </Icon>
    )
  }

  return (
    <Icon className={classNames(classes.waiting, className)} title='waiting'>
      access_time
    </Icon>
  )
}

StatusIcon.propTypes = {
  className: PropTypes.string
}

export default withStyles(styles)(StatusIcon)
