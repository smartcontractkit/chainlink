import React from 'react'

import { gql } from '@apollo/client'

import Button from '@material-ui/core/Button'
import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import CardHeader from '@material-ui/core/CardHeader'
import Chip from '@material-ui/core/Chip'
import { createStyles, withStyles, WithStyles } from '@material-ui/core/styles'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'

import { ConfirmationDialog } from 'src/components/Dialogs/ConfirmationDialog'
import { P2PKeyRow } from './P2PKeyRow'
import { ErrorRow, LoadingRow, NoContentRow } from './Rows'

export const P2P_KEYS_PAYLOAD__RESULTS_FIELDS = gql`
  fragment P2PKeysPayload_ResultsFields on P2PKey {
    id
    peerID
    publicKey
  }
`

const styles = () => {
  return createStyles({
    cardContent: {
      padding: 0,
      '&:last-child': {
        padding: 0,
      },
    },
  })
}

export interface Props extends WithStyles<typeof styles> {
  loading: boolean
  data?: FetchP2PKeys
  errorMsg?: string
  onCreate: () => void
  onDelete: (id: string) => Promise<any>
}

export const P2PKeysCard = withStyles(styles)(
  ({ classes, data, errorMsg, loading, onCreate, onDelete }: Props) => {
    const [confirmDeleteID, setConfirmDeleteID] = React.useState<string | null>(
      null,
    )

    return (
      <>
        <Card>
          <CardHeader
            action={
              <Button variant="outlined" color="primary" onClick={onCreate}>
                New P2P Key
              </Button>
            }
            title="P2P Keys"
            subheader="Manage your P2P Key Bundles"
          />
          <CardContent className={classes.cardContent}>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Key Bundle</TableCell>
                  <TableCell align="right"></TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                <LoadingRow visible={loading} />
                <NoContentRow visible={data?.p2pKeys.results?.length === 0} />
                <ErrorRow msg={errorMsg} />

                {data?.p2pKeys.results?.map((p2pKey, idx) => (
                  <P2PKeyRow
                    p2pKey={p2pKey}
                    key={idx}
                    onDelete={() => setConfirmDeleteID(p2pKey.id)}
                  />
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>

        <ConfirmationDialog
          open={!!confirmDeleteID}
          maxWidth={false}
          title="Delete OCR Key Bundle"
          body={<Chip label={confirmDeleteID} />}
          confirmButtonText="Confirm"
          onConfirm={async () => {
            if (confirmDeleteID) {
              await onDelete(confirmDeleteID)
              setConfirmDeleteID(null)
            }
          }}
          cancelButtonText="Cancel"
          onCancel={() => setConfirmDeleteID(null)}
        />
      </>
    )
  },
)
