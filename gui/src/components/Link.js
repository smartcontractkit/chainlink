import React from 'react'
import PropTypes from 'prop-types'
import Typography from '@material-ui/core/Typography'
import { Link as ReactStaticLink } from 'react-router-dom'
import { withStyles } from '@material-ui/core/styles'
import { grey } from '@material-ui/core/colors'
import classNames from 'classnames'

const styles = () => ({
  link: {
    color: grey[900],
    textDecoration: 'none'
  },
  linkContent: {
    display: 'inline-block'
  }
})

const Link = ({ children, classes, className, to }) => (
  <ReactStaticLink
    to={to}
    className={classNames(classes.link, className)}
  >
    <Typography
      variant='body1'
      color='inherit'
      className={classes.linkContent}
    >
      {children}
    </Typography>
  </ReactStaticLink>
)

Link.propTypes = {
  classes: PropTypes.object.isRequired
}

export default withStyles(styles)(Link)
