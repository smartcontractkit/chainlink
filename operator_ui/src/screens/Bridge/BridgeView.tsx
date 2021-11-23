import React from 'react'

import { gql } from '@apollo/client'

import Card from '@material-ui/core/Card'
import CardHeader from '@material-ui/core/CardHeader'
import {
  createStyles,
  withStyles,
  WithStyles,
  Theme,
} from '@material-ui/core/styles'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableRow from '@material-ui/core/TableRow'
import Grid from '@material-ui/core/Grid'

import BaseLink from 'components/BaseLink'
import Button from 'components/Button'
import Content from 'components/Content'
import { ConfirmationDialog } from 'src/components/Dialogs/ConfirmationDialog'

export const BRIDGE_PAYLOAD_FIELDS = gql`
  fragment BridgePayload_Fields on Bridge {
    id
    name
    url
    confirmations
    outgoingToken
    minimumContractPayment
  }
`

export const styles = (theme: Theme) =>
  createStyles({
    action: {
      marginLeft: theme.spacing.unit,
    },
    headerRoot: {
      borderBottom: `1px solid ${theme.palette.divider}`,
    },
    tableRow: {
      '& td:first-child': {
        fontWeight: theme.typography.fontWeightMedium,
      },
    },
  })

interface Props extends WithStyles<typeof styles> {
  bridge: BridgePayload_Fields
  onDelete: () => void
}

export const BridgeView = withStyles(styles)(
  ({ bridge, onDelete, classes }: Props) => {
    const [isDialogOpen, setIsDialogOpen] = React.useState(false)

    return (
      <Content>
        <Grid container>
          <Grid item xs={12} md={11} lg={9}>
            <Card>
              <CardHeader
                action={
                  <>
                    <Button
                      variant="danger"
                      onClick={() => setIsDialogOpen(true)}
                    >
                      Delete
                    </Button>

                    <Button
                      variant="secondary"
                      component={BaseLink}
                      href={`/bridges/${bridge.name}/edit`}
                    >
                      Edit
                    </Button>
                  </>
                }
                title="Bridge Info"
              />

              <Table>
                <TableBody>
                  <TableRow className={classes.tableRow}>
                    <TableCell>Name</TableCell>
                    <TableCell>{bridge.name}</TableCell>
                  </TableRow>

                  <TableRow className={classes.tableRow}>
                    <TableCell>URL</TableCell>
                    <TableCell>{bridge.url}</TableCell>
                  </TableRow>

                  <TableRow className={classes.tableRow}>
                    <TableCell>Confirmations</TableCell>
                    <TableCell>{bridge.confirmations}</TableCell>
                  </TableRow>

                  <TableRow className={classes.tableRow}>
                    <TableCell>Minimum Contract Payment</TableCell>
                    <TableCell>{bridge.minimumContractPayment}</TableCell>
                  </TableRow>

                  <TableRow className={classes.tableRow}>
                    <TableCell>Outgoing Token</TableCell>
                    <TableCell>{bridge.outgoingToken}</TableCell>
                  </TableRow>
                </TableBody>
              </Table>

              <ConfirmationDialog
                open={isDialogOpen}
                title="Delete Bridge"
                body="Are you sure you want to delete this bridge?"
                confirmButtonText="Confirm"
                onConfirm={onDelete}
                cancelButtonText="Cancel"
                onCancel={() => setIsDialogOpen(false)}
              />
            </Card>
          </Grid>
        </Grid>
      </Content>
    )
  },
)
