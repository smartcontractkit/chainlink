import React from 'react'
import { useLocation } from 'react-router-dom'
import { localizedTimestamp, TimeAgo } from 'components/TimeAgo'
import Card from '@material-ui/core/Card'
import Grid from '@material-ui/core/Grid'
import List from '@material-ui/core/List'
import ListItem from '@material-ui/core/ListItem'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import classNames from 'classnames'
import { connect } from 'react-redux'
import { createJobRun, fetchJob } from 'actionCreators'
import Link from 'components/Link'

const navItemStyles = ({ palette, spacing }: Theme) =>
  createStyles({
    item: {
      display: 'inline',
    },
    link: {
      paddingTop: spacing.unit * 4,
      paddingBottom: spacing.unit * 4,
      textDecoration: 'none',
      display: 'inline-block',
      borderBottom: 'solid 1px',
      borderBottomColor: palette.common.white,
      '&:hover': {
        borderBottomColor: palette.primary.main,
      },
    },
    activeLink: {
      color: palette.primary.main,
      borderBottomColor: palette.primary.main,
    },
    error: {
      color: palette.error.main,
      '&:hover': {
        borderBottomColor: palette.error.main,
      },
    },
    errorAndActiveLink: {
      borderBottomColor: palette.error.main,
    },
  })

interface NavItemProps extends WithStyles<typeof navItemStyles> {
  children: React.ReactNode
  href: string
  error?: boolean
}

const NavItem = withStyles(navItemStyles)(
  ({ children, href, classes, error }: NavItemProps) => {
    const location = useLocation()
    const active = location.pathname === href
    const linkClasses = classNames(
      classes.link,
      error && classes.error,
      active && classes.activeLink,
      error && active && classes.errorAndActiveLink,
    )

    return (
      <ListItem className={classes.item}>
        <Link href={href} className={linkClasses}>
          {children}
        </Link>
      </ListItem>
    )
  },
)

const styles = ({ palette, spacing }: Theme) =>
  createStyles({
    container: {
      backgroundColor: palette.common.white,
      padding: spacing.unit * 5,
      paddingBottom: 0,
    },
    duplicate: {
      margin: spacing.unit,
    },
    horizontalNav: {
      paddingBottom: 0,
    },
  })

interface Props extends WithStyles<typeof styles> {
  jobSpecId: string
  jobRunId: string
  jobRun?: any
}

const RegionalNav = ({ classes, jobSpecId, jobRunId, jobRun }: Props) => {
  return (
    <Card className={classes.container}>
      <Grid container spacing={0}>
        <Grid item xs={12}>
          <Typography variant="subtitle2" color="secondary" gutterBottom>
            Job run details
          </Typography>
          <Link href={`/jobs/${jobSpecId}`} variant="subtitle1" color="primary">
            {jobSpecId}
          </Link>
        </Grid>
        <Grid item xs={12}>
          <Typography variant="h3" color="secondary" gutterBottom>
            {jobRunId}
          </Typography>
        </Grid>
        <Grid item xs={12}>
          <Typography variant="subtitle2" color="textSecondary">
            {jobRun && (
              <React.Fragment>
                Started <TimeAgo tooltip={false}>{jobRun.createdAt}</TimeAgo> (
                {localizedTimestamp(jobRun.createdAt)})
              </React.Fragment>
            )}
          </Typography>
        </Grid>
        <Grid item xs={12}>
          <List className={classes.horizontalNav}>
            <NavItem href={`/jobs/${jobSpecId}/runs/${jobRunId}`}>
              Overview
            </NavItem>
            <NavItem href={`/jobs/${jobSpecId}/runs/${jobRunId}/json`}>
              JSON
            </NavItem>
          </List>
        </Grid>
      </Grid>
    </Card>
  )
}

export const ConnectedRegionalNav = connect(null, { fetchJob, createJobRun })(
  RegionalNav,
)

export default withStyles(styles)(ConnectedRegionalNav)
