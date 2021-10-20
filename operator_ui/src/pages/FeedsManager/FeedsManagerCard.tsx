import React from 'react'

import * as models from 'core/store/models'

import Card from '@material-ui/core/Card'
import CardHeader from '@material-ui/core/CardHeader'
import { createStyles, WithStyles, withStyles } from '@material-ui/core/styles'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableRow from '@material-ui/core/TableRow'
import Typography from '@material-ui/core/Typography'

const styles = () => {
  return createStyles({
    tableRoot: {
      tableLayout: 'fixed',
    },
  })
}

interface Props extends WithStyles<typeof styles> {
  manager: models.FeedsManager
}

export const FeedsManagerCard = withStyles(styles)(
  ({ classes, manager }: Props) => {
    return (
      <Card>
        <CardHeader title="Feeds Manager" />
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
