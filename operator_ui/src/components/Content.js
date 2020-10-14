import React from 'react'
import { withStyles } from '@material-ui/core/styles'

const styles = (theme) => ({
  content: {
    padding: theme.spacing.unit * 5,
  },
})

const Content = ({ children, classes }) => {
  return <div className={classes.content}>{children}</div>
}

export default withStyles(styles)(Content)
