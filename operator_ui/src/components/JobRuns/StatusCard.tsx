import React from 'react'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles
} from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import classNames from 'classnames'
import { titleCase } from 'change-case'
import PaddedCard from '../PaddedCard'
import StatusIcon from '../JobRuns/StatusIcon'

const styles = (theme: Theme) =>
  createStyles({
    completed: {
      backgroundColor: theme.palette.success.light
    },
    errored: {
      backgroundColor: theme.palette.error.light
    },
    pending: {
      backgroundColor: theme.palette.warning.light
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
    statusIcon: {
      display: 'inline-block'
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
  return (
    <PaddedCard
      className={classNames(
        classes.statusCard,
        classes[title] || classes.pending
      )}
    >
      <div className={classes.head}>
        <StatusIcon className={classes.statusIcon} width={80}>
          {title}
        </StatusIcon>
        <Typography className={classes.statusText} variant="h5" color="inherit">
          {titleCase(title)}
        </Typography>
      </div>
      {children}
    </PaddedCard>
  )
}

export default withStyles(styles)(StatusCard)
