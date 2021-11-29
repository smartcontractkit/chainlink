import React from 'react'

import { gql } from '@apollo/client'

import CancelIcon from '@material-ui/icons/Cancel'
import CheckCircleIcon from '@material-ui/icons/CheckCircle'
import EditIcon from '@material-ui/icons/Edit'
import IconButton from '@material-ui/core/IconButton'
import Grid from '@material-ui/core/Grid'
import ListItemIcon from '@material-ui/core/ListItemIcon'
import ListItemText from '@material-ui/core/ListItemText'
import Menu from '@material-ui/core/Menu'
import MoreVertIcon from '@material-ui/icons/MoreVert'
import { createStyles, WithStyles, withStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import green from '@material-ui/core/colors/green'
import red from '@material-ui/core/colors/red'

import { CopyIconButton } from 'src/components/Copy/CopyIconButton'
import { DetailsCard } from 'src/components/Cards/DetailsCard'
import { shortenHex } from 'src/utils/shortenHex'
import { MenuItemLink } from 'src/components/MenuItemLink'

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

interface Props {
  manager: FeedsManagerFields
}

export const FeedsManagerCard = ({ manager }: Props) => {
  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null)

  const handleOpen = (event: React.MouseEvent<HTMLButtonElement>) => {
    setAnchorEl(event.currentTarget)
  }

  const handleClose = () => {
    setAnchorEl(null)
  }

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
    <DetailsCard
      actions={
        <div>
          <IconButton onClick={handleOpen} aria-label="open-menu">
            <MoreVertIcon />
          </IconButton>
          <Menu
            anchorEl={anchorEl}
            open={Boolean(anchorEl)}
            onClose={handleClose}
          >
            <MenuItemLink to="/feeds_manager/edit">
              <ListItemIcon>
                <EditIcon />
              </ListItemIcon>
              <ListItemText>Edit</ListItemText>
            </MenuItemLink>
          </Menu>
        </div>
      }
    >
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
      </Grid>
    </DetailsCard>
  )
}
