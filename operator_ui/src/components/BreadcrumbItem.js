import React from 'react'
import Typography from '@material-ui/core/Typography'
import Link from 'components/Link'
import { withStyles } from '@material-ui/core/styles'

const styles = (theme) => ({
  node: {
    display: 'inline-block',
    marginLeft: theme.spacing.unit / 2,
    marginRight: theme.spacing.unit / 2,
  },
})

const renderLink = ({ children, href }) => <Link href={href}>{children}</Link>

const renderNode = ({ classes, children }) => (
  <Typography variant="body1" color="textSecondary" className={classes.node}>
    {children}
  </Typography>
)

const BreadcrumbItem = (props) => {
  if (props.href) {
    return renderLink(props)
  } else {
    return renderNode(props)
  }
}

export default withStyles(styles)(BreadcrumbItem)
