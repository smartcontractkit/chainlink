import React from 'react'

import { gql } from '@apollo/client'

import CancelIcon from '@material-ui/icons/Cancel'
import CheckCircleIcon from '@material-ui/icons/CheckCircle'
import Grid from '@material-ui/core/Grid'
import {
  createStyles,
  Theme,
  WithStyles,
  withStyles,
} from '@material-ui/core/styles'
import Paper from '@material-ui/core/Paper'
import Typography from '@material-ui/core/Typography'
import green from '@material-ui/core/colors/green'
import red from '@material-ui/core/colors/red'

import { CopyIconButton } from 'src/components/Copy/CopyIconButton'
import { shortenHex } from 'src/utils/shortenHex'
import Link from 'components/Link'

export const FEEDS_MANAGER_FIELDS = gql`
  fragment FeedsManagerFields on FeedsManager {
    id
    name
    uri
    publicKey
    jobTypes
    isBootstrapPeer
    isConnectionActive
    bootstrapPeerMultiaddr
  }
`

const connectionStatusStyles = () => {
  return createStyles({
    root: {
      display: 'flex',
    },
    connectedIcon: {
      color: green[500],
    },
    disconnectedIcon: {
      color: red[500],
    },
    text: {
      marginLeft: 4,
    },
  })
}

interface ConnectionStatusProps
  extends WithStyles<typeof connectionStatusStyles> {
  isConnected: boolean
}

const ConnectionStatus = withStyles(connectionStatusStyles)(
  ({ isConnected, classes }: ConnectionStatusProps) => {
    return (
      <div className={classes.root}>
        {isConnected ? (
          <CheckCircleIcon fontSize="small" className={classes.connectedIcon} />
        ) : (
          <CancelIcon fontSize="small" className={classes.disconnectedIcon} />
        )}

        <Typography variant="body1" inline className={classes.text}>
          {isConnected ? 'Connected' : 'Disconnected'}
        </Typography>
      </div>
    )
  },
)

const styles = (theme: Theme) => {
  return createStyles({
    tableRoot: {
      tableLayout: 'fixed',
    },
    paper: {
      marginBottom: theme.spacing.unit * 2.5,
      padding: theme.spacing.unit * 3,
    },
    editGridItem: {
      display: 'flex',
      alignItems: 'flex-end',
      justifyContent: 'flex-end',
    },
  })
}

interface Props extends WithStyles<typeof styles> {
  manager: FeedsManagerFields
}

export const FeedsManagerCard = withStyles(styles)(
  ({ classes, manager }: Props) => {
    const jobTypes = React.useMemo(() => {
      return manager.jobTypes
        .map((type) => {
          switch (type) {
            case 'FLUX_MONITOR':
              return 'Flux Monitor'
            case 'OCR':
              return 'OCR'
          }
        })
        .join(', ')
    }, [manager.jobTypes])

    return (
      <Paper className={classes.paper}>
        <Grid container>
          <Grid item xs={12} sm={6} md={3}>
            <Typography variant="subtitle2" gutterBottom>
              Status
            </Typography>
            <ConnectionStatus isConnected={manager.isConnectionActive} />
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <Typography variant="subtitle2" gutterBottom>
              Name
            </Typography>
            <Typography variant="body1" noWrap>
              {manager.name}
            </Typography>
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <Typography variant="subtitle2" gutterBottom>
              Job Types
            </Typography>
            <Typography variant="body1" noWrap>
              {jobTypes}
            </Typography>
          </Grid>

          <Grid item xs={12} sm={6} md={3}>
            {manager.isBootstrapPeer && (
              <>
                <Typography variant="subtitle2" gutterBottom>
                  Bootstrap Multiaddress
                </Typography>
                <Typography variant="body1" noWrap>
                  {manager.bootstrapPeerMultiaddr}
                </Typography>
              </>
            )}
          </Grid>

          <Grid item xs={12} sm={6} md={3}>
            <Typography variant="subtitle2" gutterBottom noWrap>
              CSA Public Key
            </Typography>
            <Typography variant="body1" gutterBottom noWrap>
              {shortenHex(manager.publicKey, { start: 6, end: 6 })}
              <CopyIconButton data={manager.publicKey} />
            </Typography>
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <Typography variant="subtitle2" gutterBottom>
              RPC URL
            </Typography>
            <Typography variant="body1" noWrap>
              {manager.uri}
            </Typography>
          </Grid>
          <Grid item xs={12} sm={6} md={6} className={classes.editGridItem}>
            <Link variant="body1" color="primary" href="/feeds_manager/edit">
              Edit
            </Link>
          </Grid>
        </Grid>
      </Paper>
    )
  },
)
