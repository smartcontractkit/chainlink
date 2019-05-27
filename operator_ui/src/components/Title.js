import React from 'react'
import PropTypes from 'prop-types'
import Typography from '@material-ui/core/Typography'
import classNames from 'classnames'
import { makeStyles } from '@material-ui/styles'

const useStyles = makeStyles(theme => ({
  title: {
    marginBottom: theme.spacing(5)
  }
}))

const Title = ({ children, className }) => {
  const classes = useStyles()
  return (
    <Typography
      variant="h4"
      color="inherit"
      className={classNames(className, classes.title)}
    >
      {children}
    </Typography>
  )
}

Title.propTypes = {
  children: PropTypes.oneOfType([
    PropTypes.arrayOf(PropTypes.node),
    PropTypes.node
  ])
}

export default Title
