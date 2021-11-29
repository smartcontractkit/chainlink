import React from 'react'

import { useLocation } from 'react-router-dom'

import Badge from '@material-ui/core/Badge'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import Tabs from '@material-ui/core/Tabs'
import { TabLink } from 'src/components/Tab/TabLink'

const styles = (theme: Theme) =>
  createStyles({
    tabs: {
      marginTop: theme.spacing.unit * 4,
      marginBottom: theme.spacing.unit * 2.5,
      borderBottom: `1px solid ${theme.palette.grey['300']}`,
    },
    badge: {
      padding: `0 ${theme.spacing.unit * 2}px`,
    },
  })

interface Props extends WithStyles<typeof styles> {
  id: string
  errorsCount: number
  runsCount: number
}

export const JobTabs = withStyles(styles)(
  ({ classes, id, errorsCount, runsCount }: Props) => {
    const { pathname } = useLocation()

    return (
      <Tabs value={pathname} className={classes.tabs} indicatorColor="primary">
        <TabLink label="Overview" to={`/jobs/${id}`} value={`/jobs/${id}`} />
        <TabLink
          label="Definition"
          to={`/jobs/${id}/definition`}
          value={`/jobs/${id}/definition`}
        />
        <TabLink
          label={
            <Badge
              badgeContent={errorsCount}
              color="error"
              className={classes.badge}
            >
              Errors
            </Badge>
          }
          to={`/jobs/${id}/errors`}
          value={`/jobs/${id}/errors`}
        />
        <TabLink
          label={
            <Badge
              badgeContent={runsCount}
              color="primary"
              className={classes.badge}
              max={99999}
            >
              Runs
            </Badge>
          }
          to={`/jobs/${id}/runs`}
          value={`/jobs/${id}/runs`}
        />
      </Tabs>
    )
  },
)
