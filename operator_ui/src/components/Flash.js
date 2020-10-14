import React from 'react'
import PropTypes from 'prop-types'
import Card from '@material-ui/core/Card'
import Typography from '@material-ui/core/Typography'
import classNames from 'classnames'
import { withStyles } from '@material-ui/core/styles'

const styles = (theme) => {
  const success = theme.palette.success || {}
  const warning = theme.palette.warning || {}

  return {
    base: {
      paddingLeft: theme.spacing.unit * 5,
      paddingRight: theme.spacing.unit * 5,
    },
    success: {
      backgroundColor: success.main,
      color: success.contrastText,
    },
    error: {
      backgroundColor: theme.palette.error.dark,
      color: theme.palette.error.contrastText,
    },
    warning: {
      backgroundColor: warning.contrastText,
      color: warning.main,
    },
  }
}

const applyClass = ({ success, error, warning, classes, className }) => {
  let type

  if (success) {
    type = classes.success
  } else if (error) {
    type = classes.error
  } else if (warning) {
    type = classes.warning
  }

  return classNames(classes.base, className, type)
}

const Flash = (props) => (
  <Card className={applyClass(props)} square>
    <Typography variant="body2" color="inherit" component="div">
      {props.children}
    </Typography>
  </Card>
)

Flash.defaultProps = {
  success: false,
  error: false,
  warning: false,
}

Flash.propTypes = {
  success: PropTypes.bool,
  error: PropTypes.bool,
  warning: PropTypes.bool,
}

export default withStyles(styles)(Flash)
