import React from 'react'

import Button from 'src/components/Button'
import TableCell from '@material-ui/core/TableCell'
import TableRow from '@material-ui/core/TableRow'

import { KeyBundle } from './KeyBundle'
import { CopyIconButton } from 'src/components/Copy/CopyIconButton'

interface Props {
  bundle: OcrKeyBundlesPayload_ResultsFields
  onDelete: () => void
}

export const OCRKeyBundleRow: React.FC<Props> = ({ bundle, onDelete }) => {
  return (
    <TableRow hover>
      <TableCell>
        <KeyBundle
          primary={
            <b>
              Key ID: {bundle.id} <CopyIconButton data={bundle.id} />
            </b>
          }
          secondary={[
            <>Config Public Key: {bundle.configPublicKey}</>,
            <>Signing Address: {bundle.onChainSigningAddress}</>,
            <>Off-Chain Public Key: {bundle.offChainPublicKey}</>,
          ]}
        />
      </TableCell>
      <TableCell align="right">
        <Button onClick={onDelete} variant="danger" size="medium">
          Delete
        </Button>
      </TableCell>
    </TableRow>
  )
}
