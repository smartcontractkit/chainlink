import React from 'react'
import PropTypes from 'prop-types'
import Typography from '@material-ui/core/Typography'
import { Link as ReactStaticLink } from 'react-static'
import { withStyles } from '@material-ui/core/styles'
import { blue } from '@material-ui/core/colors'
import classNames from 'classnames'

const styles = theme => ({
  link: {
    color: blue[600],
    textDecoration: 'none'
  },
  linkContent: {
    display: 'inline-block'
  }
})

const Link = ({children, classes, className, to}) => (
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
