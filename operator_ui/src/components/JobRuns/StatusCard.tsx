import React from 'react'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles
} from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import classNames from 'classnames'
import PaddedCard from '@chainlink/styleguide/src/components/PaddedCard'
import { titleCase } from 'change-case'
import StatusIcon from '../JobRuns/StatusIcon'

const styles = (theme: Theme) =>
  createStyles({
    completed: {
      // backgroundColor: theme.palette.success.light
    },
    errored: {
      backgroundColor: theme.palette.error.light
    },
    pending: {
      // backgroundColor: theme.palette.warning.light
    },
    statusCard: {
      '&:last-child': {
        paddingBottom: theme.spacing.unit * 2
      }
    },
    head: {
      display: 'flex',
      alignItems: 'center'
    },
    statusText: {
      display: 'inline-block',
      paddingLeft: theme.spacing.unit * 2,
      textTransform: 'capitalize'
    }
  })

interface IProps extends WithStyles<typeof styles> {
  title: string
  children?: React.ReactNode
}

const StatusCard = ({ title, classes, children }: IProps) => {
  const statusClass = classes[title as keyof typeof classes] || classes.pending

  return (
    <PaddedCard className={classNames(classes.statusCard, statusClass)}>
      <div className={classes.head}>
        <StatusIcon width={80}>{title}</StatusIcon>
        <Typography className={classes.statusText} variant="h5" color="inherit">
          {titleCase(title)}
        </Typography>
      </div>
      {children}
    </PaddedCard>
  )
}

export default withStyles(styles)(StatusCard)
