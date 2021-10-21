import React from 'react'
import { useHistory } from 'react-router-dom'

import { FeedsManager } from 'core/store/models'

import Card from '@material-ui/core/Card'
import CardHeader from '@material-ui/core/CardHeader'
import { createStyles, WithStyles, withStyles } from '@material-ui/core/styles'
import IconButton from '@material-ui/core/IconButton'
import EditIcon from '@material-ui/icons/Edit'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableRow from '@material-ui/core/TableRow'
import Typography from '@material-ui/core/Typography'
import green from '@material-ui/core/colors/green'
import red from '@material-ui/core/colors/red'

import CheckCircleIcon from '@material-ui/icons/CheckCircle'
import CancelIcon from '@material-ui/icons/Cancel'

const cardSubheaderStyles = () => {
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

interface CardSubheaderProps extends WithStyles<typeof cardSubheaderStyles> {
  isConnected: boolean
}

const CardSubheader = withStyles(cardSubheaderStyles)(
  ({ isConnected, classes }: CardSubheaderProps) => {
    return (
      <div className={classes.root}>
        {isConnected ? (
          <CheckCircleIcon fontSize="small" className={classes.connectedIcon} />
        ) : (
          <CancelIcon fontSize="small" className={classes.disconnectedIcon} />
        )}

        <Typography
          variant="body1"
          color="textSecondary"
          inline
          className={classes.text}
        >
          {isConnected ? 'Connected' : 'Disconnected'}
        </Typography>
      </div>
    )
  },
)

const styles = () => {
  return createStyles({
    tableRoot: {
      tableLayout: 'fixed',
    },
  })
}

interface Props extends WithStyles<typeof styles> {
  manager: FeedsManager
}

export const FeedsManagerCard = withStyles(styles)(
  ({ classes, manager }: Props) => {
    const history = useHistory()

    return (
      <Card>
        <CardHeader
          title="Feeds Manager"
          subheader={<CardSubheader isConnected={manager.isConnectionActive} />}
          action={
            <IconButton onClick={() => history.push('/feeds_manager/edit')}>
              <EditIcon fontSize="small" />
            </IconButton>
          }
        />
        <Table className={classes.tableRoot}>
          <TableBody>
            <TableRow>
              <TableCell>
                <Typography>Name</Typography>
                <Typography variant="subtitle1" color="textSecondary">
                  {manager.name}
                </Typography>
              </TableCell>
            </TableRow>

            <TableRow>
              <TableCell>
                <Typography>URI</Typography>
                <Typography variant="subtitle1" color="textSecondary">
                  {manager.uri}
                </Typography>
              </TableCell>
            </TableRow>

            <TableRow>
              <TableCell>
                <Typography>Public Key</Typography>
                <Typography variant="subtitle1" color="textSecondary" noWrap>
                  {manager.publicKey}
                </Typography>
              </TableCell>
            </TableRow>

            <TableRow>
              <TableCell>
                <Typography>Job Types</Typography>
                <Typography variant="subtitle1" color="textSecondary">
                  {manager.jobTypes.join(',')}
                </Typography>
              </TableCell>
            </TableRow>

            {manager.isBootstrapPeer && (
              <TableRow>
                <TableCell>
                  <Typography>Bootstrap Peer Multiaddress</Typography>
                  <Typography variant="subtitle1" color="textSecondary">
                    {manager.bootstrapPeerMultiaddr}
                  </Typography>
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </Card>
    )
  },
)
