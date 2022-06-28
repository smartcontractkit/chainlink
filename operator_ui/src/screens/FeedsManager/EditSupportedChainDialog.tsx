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
import { useOCR2KeysQuery } from 'src/hooks/queries/useOCR2KeysQuery'

type Props = {
  cfg: FeedsManager_ChainConfigFields | null
  open: boolean
  onClose: () => void
} & Pick<FormProps, 'onSubmit'>

export const EditSupportedChainDialog = ({
  cfg,
  onClose,
  open,
  onSubmit,
}: Props) => {
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

  const { data: ocr2KeysData } = useOCR2KeysQuery({
    fetchPolicy: 'cache-and-network',
  })

  if (!cfg) {
    return null
  }

  const initialValues = {
    chainID: cfg.chainID,
    chainType: 'EVM',
    accountAddr: cfg.accountAddr,
    adminAddr: cfg.adminAddr,
    fluxMonitorEnabled: cfg.fluxMonitorJobConfig.enabled,
    ocr1Enabled: cfg.ocr1JobConfig.enabled,
    ocr1IsBootstrap: cfg.ocr1JobConfig.isBootstrap,
    ocr1Multiaddr: cfg.ocr1JobConfig.multiaddr,
    ocr1P2PPeerID: cfg.ocr1JobConfig.p2pPeerID,
    ocr1KeyBundleID: cfg.ocr1JobConfig.keyBundleID,
    ocr2Enabled: cfg.ocr2JobConfig.enabled,
    ocr2IsBootstrap: cfg.ocr2JobConfig.isBootstrap,
    ocr2Multiaddr: cfg.ocr2JobConfig.multiaddr,
    ocr2P2PPeerID: cfg.ocr2JobConfig.p2pPeerID,
    ocr2KeyBundleID: cfg.ocr2JobConfig.keyBundleID,
  }

  const chainIDs: string[] = chainData
    ? chainData.chains.results.map((c) => c.id)
    : []

  const accounts = accountData ? accountData.ethKeys.results : []
  const p2pKeys = p2pKeysData ? p2pKeysData.p2pKeys.results : []
  const ocrKeys = ocrKeysData ? ocrKeysData.ocrKeyBundles.results : []
  const ocr2Keys = ocr2KeysData ? ocr2KeysData.ocr2KeyBundles.results : []

  return (
    <Dialog onClose={onClose} open={open} disableBackdropClick>
      <DialogTitle disableTypography>
        <Typography variant="body2">Edit Supported Chain</Typography>
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
          ocr2Keys={ocr2Keys}
          editing
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
