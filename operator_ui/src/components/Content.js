import React from 'react'
import { makeStyles } from '@material-ui/core'

const useStyles = makeStyles(theme => ({
  content: {
    padding: theme.spacing(5)
  }
}))

const Content = ({ children }) => {
  const classes = useStyles()
  return <div className={classes.content}>{children}</div>
}

export default Content
