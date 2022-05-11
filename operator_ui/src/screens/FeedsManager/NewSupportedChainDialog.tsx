import React from 'react'

import Button from '@material-ui/core/Button'
import Dialog from '@material-ui/core/Dialog'
import DialogActions from '@material-ui/core/DialogActions'
import DialogContent from '@material-ui/core/DialogContent'
import DialogTitle from '@material-ui/core/DialogTitle'
import Typography from '@material-ui/core/Typography'

import {
  ChainConfigurationForm,
  Props as FormProps,
} from 'src/components/Form/ChainConfigurationForm'
import { useChainsQuery } from 'src/hooks/queries/useChainsQuery'
import { useEVMAccountsQuery } from 'src/hooks/queries/useEVMAccountsQuery'
import { useP2PKeysQuery } from 'src/hooks/queries/useP2PKeysQuery'
import { useOCRKeysQuery } from 'src/hooks/queries/useOCRKeysQuery'

type Props = {
  open: boolean
  onClose: () => void
} & Pick<FormProps, 'onSubmit'>

export const NewSupportedChainDialog = ({ onClose, open, onSubmit }: Props) => {
  const formRef = React.useRef()

  const { data: chainData } = useChainsQuery({
    fetchPolicy: 'network-only',
  })

  const { data: accountData } = useEVMAccountsQuery({
    fetchPolicy: 'cache-and-network',
  })

  const { data: p2pKeysData } = useP2PKeysQuery({
    fetchPolicy: 'cache-and-network',
  })

  const { data: ocrKeysData } = useOCRKeysQuery({
    fetchPolicy: 'cache-and-network',
  })

  const initialValues = {
    chainID: '',
    chainType: 'EVM',
    accountAddr: '',
    adminAddr: '',
    fluxMonitorEnabled: false,
    ocr1Enabled: false,
    ocr1IsBootstrap: false,
    ocr1Multiaddr: '',
    ocr1P2PPeerID: '',
    ocr1KeyBundleID: '',
    ocr2Enabled: false,
  }

  const chainIDs: string[] = chainData
    ? chainData.chains.results.map((c) => c.id)
    : []

  const accounts = accountData ? accountData.ethKeys.results : []
  const p2pKeys = p2pKeysData ? p2pKeysData.p2pKeys.results : []
  const ocrKeys = ocrKeysData ? ocrKeysData.ocrKeyBundles.results : []

  return (
    <Dialog onClose={onClose} open={open} disableBackdropClick>
      <DialogTitle disableTypography>
        <Typography variant="body2">New Supported Chain</Typography>
      </DialogTitle>

      <DialogContent>
        <ChainConfigurationForm
          innerRef={formRef}
          initialValues={initialValues}
          onSubmit={onSubmit}
          chainIDs={chainIDs}
          accounts={accounts}
          p2pKeys={p2pKeys}
          ocrKeys={ocrKeys}
        />
      </DialogContent>

      <DialogActions>
        <Button onClick={onClose}>Cancel</Button>
        <Button
          color="primary"
          type="submit"
          form="chain-configuration-form"
          variant="contained"
        >
          Submit
        </Button>
      </DialogActions>
    </Dialog>
  )
}
