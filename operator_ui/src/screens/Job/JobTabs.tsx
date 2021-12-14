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

export interface Props extends WithStyles<typeof styles> {
  id: string
  errorsCount: number
  runsCount: number
  refetchRecentRuns: () => void
}

enum JobTab {
  Overview,
  Definition,
  Errors,
  Runs,
}

export const JobTabs = withStyles(styles)(
  ({ classes, id, errorsCount, refetchRecentRuns, runsCount }: Props) => {
    const { pathname } = useLocation()

    const tabs = React.useMemo(
      () => ({
        [JobTab.Overview]: `/jobs/${id}`,
        [JobTab.Definition]: `/jobs/${id}/definition`,
        [JobTab.Errors]: `/jobs/${id}/errors`,
        [JobTab.Runs]: `/jobs/${id}/runs`,
      }),
      [id],
    )

    return (
      <Tabs
        value={pathname}
        className={classes.tabs}
        indicatorColor="primary"
        onChange={(_, value) => {
          if (value === tabs[JobTab.Overview]) {
            refetchRecentRuns()
          }
        }}
      >
        <TabLink
          label="Overview"
          to={tabs[JobTab.Overview]}
          value={tabs[JobTab.Overview]}
        />
        <TabLink
          label="Definition"
          to={tabs[JobTab.Definition]}
          value={tabs[JobTab.Definition]}
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
          to={tabs[JobTab.Errors]}
          value={tabs[JobTab.Errors]}
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
          to={tabs[JobTab.Runs]}
          value={tabs[JobTab.Runs]}
        />
      </Tabs>
    )
  },
)
