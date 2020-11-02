import React from 'react'
import { withStyles } from '@material-ui/core/styles'
import classNames from 'classnames'

const styles = (theme) => {
  return {
    breadcrumb: {
      marginBottom: theme.spacing.unit * 3,
    },
  }
}

const Breadcrumb = ({ children, className, classes }) => (
  <div className={classNames(className, classes.breadcrumb)}>{children}</div>
)

export default withStyles(styles)(Breadcrumb)
