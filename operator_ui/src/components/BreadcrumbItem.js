import React from 'react'
import Typography from '@material-ui/core/Typography'
import Link from 'components/Link'
import { makeStyles } from '@material-ui/styles'

const useStyles = makeStyles(theme => ({
  node: {
    display: 'inline-block',
    marginLeft: theme.spacing(1 / 2),
    marginRight: theme.spacing(1 / 2)
  }
}))

const renderLink = ({ children, href }) => <Link to={href}>{children}</Link>

const renderNode = ({ children }) => {
  const classes = useStyles()
  return (
    <Typography variant="body1" color="textSecondary" className={classes.node}>
      {children}
    </Typography>
  )
}

const BreadcrumbItem = props => {
  if (props.href) {
    return renderLink(props)
  } else {
    return renderNode(props)
  }
}

export default BreadcrumbItem
