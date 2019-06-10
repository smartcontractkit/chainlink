import React from 'react'
import { Link as ReactStaticLink } from 'react-router-dom'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles
} from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import { grey } from '@material-ui/core/colors'
import classNames from 'classnames'

const styles = (_theme: Theme) =>
  createStyles({
    link: {
      color: grey[900],
      textDecoration: 'none'
    },
    linkContent: {
      display: 'inline-block'
    }
  })

interface IProps extends WithStyles<typeof styles> {
  children: React.ReactNode
  to: string
  className?: string
}

const Link = ({ children, classes, className, to }: IProps) => (
  <ReactStaticLink to={to} className={classNames(classes.link, className)}>
    <Typography variant="body1" color="inherit" className={classes.linkContent}>
      {children}
    </Typography>
  </ReactStaticLink>
)

export default withStyles(styles)(Link)
