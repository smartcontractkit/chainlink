import React, { useState } from 'react'
import { useDispatch } from 'react-redux'

import TableCell from '@material-ui/core/TableCell'
import TableRow from '@material-ui/core/TableRow'
import Typography from '@material-ui/core/Typography'
import Dialog from '@material-ui/core/Dialog'
import Grid from '@material-ui/core/Grid'
import TextField from '@material-ui/core/TextField'
import Checkbox from '@material-ui/core/Checkbox'
import FormGroup from '@material-ui/core/FormGroup'
import FormControlLabel from '@material-ui/core/FormControlLabel'
import { ApolloQueryResult } from '@apollo/client'

import {
  createStyles,
  withStyles,
  Theme,
  WithStyles,
} from '@material-ui/core/styles'

import Button from 'components/Button'
import Close from 'components/Icons/Close'
import ErrorMessage from 'components/Notifications/DefaultError'
import * as api from 'api'
import { notifySuccess, notifyError } from 'actionCreators'
import { ApiResponse } from 'utils/json-api-client'
import { EVMKeysChainRequest, EVMKey } from 'core/store/models'
import { CopyIconButton } from 'src/components/Copy/CopyIconButton'
import { fromJuels } from 'src/utils/tokens/link'
import { shortenHex } from 'src/utils/shortenHex'
import { TimeAgo } from 'src/components/TimeAgo'

const styles = (theme: Theme) =>
  createStyles({
    container: {
      backgroundColor: theme.palette.common.white,
      padding: theme.spacing.unit * 5,
      paddingBottom: 0,
    },
    chainId: {
      overflow: 'hidden',
      textOverflow: 'ellipsis',
    },
    badgePadding: {
      paddingLeft: theme.spacing.unit * 2,
      paddingRight: theme.spacing.unit * 2,
      marginLeft: theme.spacing.unit * -2,
      marginRight: theme.spacing.unit * -2,
      lineHeight: '1rem',
    },
    dialogPaper: {
      minHeight: '360px',
      maxHeight: '360px',
      minWidth: '670px',
      maxWidth: '670px',
      overflow: 'hidden',
      borderRadius: theme.spacing.unit * 3,
    },
    warningText: {
      fontWeight: 500,
      marginLeft: theme.spacing.unit * 3,
      marginTop: theme.spacing.unit * 3,
      marginBottom: theme.spacing.unit,
    },
    closeButton: {
      marginRight: theme.spacing.unit * 3,
      marginTop: theme.spacing.unit * 3,
    },
    infoText: {
      fontSize: theme.spacing.unit * 2,
      fontWeight: 450,
      marginLeft: theme.spacing.unit * 6,
    },
    modalContent: {
      width: 'inherit',
    },
    deleteButton: {
      marginTop: theme.spacing.unit * 4,
    },
    runJobButton: {
      marginBottom: theme.spacing.unit * 3,
    },
    runJobModalContent: {
      overflow: 'hidden',
    },
  })

interface Props {
  classes: WithStyles<typeof styles>['classes']
  ethKey: EthKeysPayload_ResultsFields
  refetch?: () => Promise<ApolloQueryResult<FetchEthKeys>>
}

function apiCall({
  evmChainID,
  address,
  nextNonce,
  abandon,
  enabled,
}: {
  evmChainID: string
  address: string
  nextNonce: BigInt | null
  abandon: boolean
  enabled: boolean
}): Promise<ApiResponse<EVMKey>> {
  const definition: EVMKeysChainRequest = {
    evmChainID,
    address,
    nextNonce,
    abandon,
    enabled,
  }
  return api.v2.evmKeys.chain(definition)
}

const SuccessNotification = () => <>Successfully updated EVM key</>

const UnstyledEVMAccountRow: React.FC<Props> = ({
  classes,
  ethKey,
  refetch,
}) => {
  const dispatch = useDispatch()

  const [modalOpen, setModalOpen] = useState(false)
  const [enabled, setEnabled] = useState(!ethKey.isDisabled)
  const [nextNonce, setNextNonce] = useState<BigInt | null>(null)
  const [abandon, setAbandon] = useState(false)

  const onSubmit = (event: React.SyntheticEvent) => {
    event.preventDefault()
    handleUpdate(nextNonce, abandon, enabled)
  }

  const handleEnabledCheckboxChange = () => {
    setEnabled(!enabled)
  }

  const handleNextNonceFieldChange = (event: any) => {
    setNextNonce(event.target.value)
  }

  const handleAbandonCheckboxChange = () => {
    setAbandon(!abandon)
  }

  const closeModal = () => {
    setModalOpen(false)
    // reset state
    setAbandon(false)
    setNextNonce(null)
    setEnabled(!ethKey.isDisabled)
  }

  async function handleUpdate(
    nextNonce: BigInt | null,
    abandon: boolean,
    enabled: boolean,
  ) {
    apiCall({
      evmChainID: ethKey.chain.id,
      address: ethKey.address,
      nextNonce,
      abandon,
      enabled,
    })
      .then(({ data }) => {
        refetch && refetch()
        closeModal()
        dispatch(notifySuccess(SuccessNotification, data))
      })
      .catch((error) => {
        refetch && refetch()
        closeModal()
        dispatch(notifyError(ErrorMessage, error))
      })
  }

  return (
    <>
      <Dialog
        open={modalOpen}
        classes={{ paper: classes.dialogPaper }}
        onClose={closeModal}
        aria-labelledby="modal-modal-title"
        aria-describedby="modal-modal-description"
      >
        <form onSubmit={onSubmit}>
          <Grid container spacing={0}>
            <Grid item className={classes.modalContent}>
              <Grid container alignItems="baseline" justify="space-between">
                <Grid item>
                  <Typography
                    variant="h5"
                    color="secondary"
                    className={classes.warningText}
                  >
                    Key Admin Override
                  </Typography>
                  <Typography
                    variant="h6"
                    color="secondary"
                    className={classes.warningText}
                  >
                    Modifying key {ethKey.address} for chain {ethKey.chain.id}
                  </Typography>
                </Grid>
                <Grid item>
                  <Close className={classes.closeButton} onClick={closeModal} />
                </Grid>
              </Grid>
              <Grid container direction="column">
                <FormGroup>
                  <FormControlLabel
                    className={classes.infoText}
                    color="secondary"
                    control={
                      <Checkbox
                        name="enabledCheckbox"
                        checked={enabled}
                        onChange={handleEnabledCheckboxChange}
                      />
                    }
                    label="Enabled"
                  />
                </FormGroup>
                <FormGroup>
                  <FormControlLabel
                    className={classes.infoText}
                    color="secondary"
                    control={
                      <TextField
                        name="nextNonceField"
                        type="number"
                        onChange={handleNextNonceFieldChange}
                      />
                    }
                    label="Next nonce manual override (optional)"
                  />
                </FormGroup>
                <FormGroup>
                  <FormControlLabel
                    className={classes.infoText}
                    color="secondary"
                    control={
                      <Checkbox
                        name="abandonCheckbox"
                        checked={abandon}
                        onChange={handleAbandonCheckboxChange}
                      />
                    }
                    label="Abandon all current transactions (use with caution!)"
                  />
                </FormGroup>
                <Grid
                  container
                  spacing={0}
                  alignItems="center"
                  justify="center"
                >
                  <Grid item className={classes.deleteButton}>
                    <Button variant="danger" type="submit">
                      Update
                    </Button>
                  </Grid>
                </Grid>
              </Grid>
            </Grid>
          </Grid>
        </form>
      </Dialog>
      <TableRow hover>
        <TableCell>
          <Typography variant="body1">
            {shortenHex(ethKey.address, { start: 6, end: 6 })}{' '}
            <CopyIconButton data={ethKey.address} />
          </Typography>
        </TableCell>
        <TableCell>
          <Typography variant="body1">{ethKey.chain.id}</Typography>
        </TableCell>
        <TableCell>
          <Typography variant="body1">
            {ethKey.isDisabled ? 'Disabled' : 'Enabled'}
          </Typography>
        </TableCell>
        <TableCell>
          <Typography variant="body1">
            {ethKey.linkBalance && fromJuels(ethKey.linkBalance)}
          </Typography>
        </TableCell>
        <TableCell>
          <Typography variant="body1">{ethKey.ethBalance}</Typography>
        </TableCell>
        <TableCell>
          <Typography variant="body1">
            <TimeAgo tooltip>{ethKey.createdAt}</TimeAgo>
          </Typography>
        </TableCell>
        <TableCell>
          <Typography variant="body1">
            <Button onClick={() => setModalOpen(true)}>Admin</Button>
          </Typography>
        </TableCell>
      </TableRow>
    </>
  )
}

export const EVMAccountRow = withStyles(styles)(UnstyledEVMAccountRow)
