import React from 'react'
import classNames from 'classnames'
import { makeStyles } from '@material-ui/core'

const useStyles = makeStyles(theme => {
  return {
    breadcrumb: {
      marginBottom: theme.spacing(3)
    }
  }
})

const Breadcrumb = ({ children, className }) => {
  const classes = useStyles()
  return (
    <div className={classNames(className, classes.breadcrumb)}>{children}</div>
  )
}

export default Breadcrumb
