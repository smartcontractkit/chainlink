import React from 'react'

import Button from '@material-ui/core/Button'
import Card from '@material-ui/core/Card'
import CardHeader from '@material-ui/core/CardHeader'
import Chip from '@material-ui/core/Chip'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'

import { ConfirmationDialog } from 'src/components/Dialogs/ConfirmationDialog'
import { ErrorRow } from 'src/components/TableRow/ErrorRow'
import { LoadingRow } from 'src/components/TableRow/LoadingRow'
import { NoContentRow } from 'src/components/TableRow/NoContentRow'
import { OCRKeyBundleRow } from './OCRKeyBundleRow'

export interface Props {
  loading: boolean
  data?: FetchOcrKeyBundles
  errorMsg?: string
  onCreate: () => void
  onDelete: (id: string) => Promise<any>
}

export const OCRKeysCard: React.FC<Props> = ({
  data,
  errorMsg,
  loading,
  onCreate,
  onDelete,
}) => {
  const [confirmDeleteID, setConfirmDeleteID] = React.useState<string | null>(
    null,
  )

  return (
    <>
      <Card>
        <CardHeader
          action={
            <Button variant="outlined" color="primary" onClick={onCreate}>
              New OCR Key
            </Button>
          }
          title="Off-Chain Reporting Keys"
          subheader="Manage your Off-Chain Reporting Key Bundles"
        />

        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Key Bundle</TableCell>
              <TableCell align="right"></TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            <LoadingRow visible={loading} />
            <NoContentRow visible={data?.ocrKeyBundles.results?.length === 0} />
            <ErrorRow msg={errorMsg} />

            {data?.ocrKeyBundles.results?.map((bundle, idx) => (
              <OCRKeyBundleRow
                bundle={bundle}
                key={idx}
                onDelete={() => setConfirmDeleteID(bundle.id)}
              />
            ))}
          </TableBody>
        </Table>
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
}
