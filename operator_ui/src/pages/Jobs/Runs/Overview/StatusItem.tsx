import React from 'react'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import ExpansionPanel from '@material-ui/core/ExpansionPanel'
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary'
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails'
import Typography from '@material-ui/core/Typography'
import ExpandMoreIcon from '@material-ui/icons/ExpandMore'
import StatusIcon from 'components/StatusIcon'
import classNames from 'classnames'
import { Grid } from '@material-ui/core'

const withChildrenStyles = () =>
  createStyles({
    summary: {
      minHeight: '0 !important',
    },
    content: {
      margin: '12px 0 !important',
    },
    expansionPanel: {
      boxShadow: 'none',
    },
  })

interface WithChildrenProps extends WithStyles<typeof withChildrenStyles> {
  summary: string
  minConfirmations?: number | null
  confirmations?: number | null
}

const withChildrenUnstyled: React.FC<WithChildrenProps> = ({
  summary,
  classes,
  children,
  confirmations,
  minConfirmations,
}) => (
  <ExpansionPanel className={classes.expansionPanel}>
    <ExpansionPanelSummary
      className={classes.summary}
      classes={{ content: classes.content }}
      expandIcon={<ExpandMoreIcon />}
    >
      <Grid container alignItems="baseline">
        <Grid item sm={10}>
          <Typography variant="h5">{summary}</Typography>
        </Grid>
        <Grid item>
          {minConfirmations && typeof confirmations === 'number' && (
            <Typography variant="h6" color="secondary">
              Confirmations {confirmations}/{minConfirmations}
            </Typography>
          )}
        </Grid>
      </Grid>
    </ExpansionPanelSummary>
    <ExpansionPanelDetails>{children}</ExpansionPanelDetails>
  </ExpansionPanel>
)

const WithChildren = withStyles(withChildrenStyles)(withChildrenUnstyled)

interface WithoutChildrenProps {
  summary: string
}

const WithoutChildren = ({ summary }: WithoutChildrenProps) => {
  return <Typography>{summary}</Typography>
}

const styles = (theme: Theme) =>
  createStyles({
    borderTop: {
      borderTop: 'solid 1px',
      borderTopColor: theme.palette.divider,
    },
    item: {
      position: 'relative',
      paddingLeft: 50,
    },
    status: {
      position: 'absolute',
      top: 0,
      left: 0,
      paddingTop: 25,
      paddingLeft: 30,
      borderRight: 'solid 1px',
      borderRightColor: theme.palette.divider,
      width: 50,
      height: '100%',
    },
    details: {
      padding: theme.spacing.unit * 2,
    },
  })

interface Props extends WithStyles<typeof styles> {
  status: string
  borderTop: boolean
  children: React.ReactNode
  summary: string
  minConfirmations?: number | null
  confirmations?: number | null
}

const StatusItem = ({
  status,
  summary,
  borderTop,
  children,
  classes,
  confirmations,
  minConfirmations,
}: Props) => (
  <div className={classNames(classes.item, { [classes.borderTop]: borderTop })}>
    <div className={classes.status}>
      <StatusIcon width={38} height={38}>
        {status}
      </StatusIcon>
    </div>
    <div className={classes.details}>
      {children ? (
        <WithChildren
          summary={summary}
          confirmations={confirmations}
          minConfirmations={minConfirmations}
        >
          {children}
        </WithChildren>
      ) : (
        <WithoutChildren summary={summary} />
      )}
    </div>
  </div>
)

export default withStyles(styles)(StatusItem)
