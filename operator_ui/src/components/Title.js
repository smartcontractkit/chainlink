import React from 'react'
import PropTypes from 'prop-types'
import { withStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import classNames from 'classnames'

const styles = (theme) => ({
  title: {
    marginBottom: theme.spacing.unit * 5,
  },
})

export const Title = withStyles(styles)(({ children, classes, className }) => (
  <Typography
    variant="h4"
    color="inherit"
    className={classNames(className, classes.title)}
  >
    {children}
  </Typography>
))

Title.propTypes = {
  children: PropTypes.oneOfType([
    PropTypes.arrayOf(PropTypes.node),
    PropTypes.node,
  ]),
}

export default Title
