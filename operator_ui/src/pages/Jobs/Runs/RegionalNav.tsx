import React from 'react'
import { useLocation } from 'react-router-dom'
import { localizedTimestamp, TimeAgo } from 'components/TimeAgo'
import Paper from '@material-ui/core/Paper'
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
  jobId: string
  jobRunId: string
  jobRun?: any
}

const RegionalNav = ({ classes, jobId, jobRunId, jobRun }: Props) => {
  return (
    <Paper className={classes.container}>
      <Grid container spacing={0}>
        <Grid item xs={12}>
          <Typography variant="h4" color="secondary" gutterBottom>
            Job Run #{jobRunId}
          </Typography>
        </Grid>

        <Grid item xs={3}>
          <Typography variant="subtitle1" color="textSecondary">
            Job
          </Typography>
          <Link href={`/jobs/${jobId}`} variant="subtitle2" color="primary">
            {jobId}
          </Link>
        </Grid>

        <Grid item xs={3}>
          <Typography variant="subtitle1" color="textSecondary">
            Started
          </Typography>
          <Typography variant="subtitle2" color="textSecondary">
            {jobRun && (
              <React.Fragment>
                <TimeAgo tooltip={false}>{jobRun.createdAt}</TimeAgo> (
                {localizedTimestamp(jobRun.createdAt)})
              </React.Fragment>
            )}
          </Typography>
        </Grid>

        <Grid item xs={12}>
          <List className={classes.horizontalNav}>
            <NavItem href={`/jobs/${jobId}/runs/${jobRunId}`}>Overview</NavItem>
            <NavItem href={`/jobs/${jobId}/runs/${jobRunId}/json`}>
              JSON
            </NavItem>
          </List>
        </Grid>
      </Grid>
    </Paper>
  )
}

export default withStyles(styles)(RegionalNav)
