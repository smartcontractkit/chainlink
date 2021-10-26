import React, { useState } from 'react'
import { ApiResponse, BadRequestError } from 'utils/json-api-client'
import Button from 'components/Button'
import * as api from 'api'
import { useDispatch } from 'react-redux'
import { Chain, CreateChainRequest } from 'core/store/models'
import BaseLink from 'components/BaseLink'
import ErrorMessage from 'components/Notifications/DefaultError'
import { notifyError, notifySuccess } from 'actionCreators'
import Content from 'components/Content'
import {
  Card,
  CardContent,
  CardHeader,
  Checkbox,
  CircularProgress,
  FormLabel,
  Grid,
  TextField,
} from '@material-ui/core'
import {
  createStyles,
  Theme,
  WithStyles,
  withStyles,
} from '@material-ui/core/styles'

const styles = (theme: Theme) =>
  createStyles({
    loader: {
      position: 'absolute',
    },
    emptyTasks: {
      padding: theme.spacing.unit * 3,
    },
  })

const SuccessNotification = ({ id }: { id: string }) => (
  <>
    Successfully created chain{' '}
    <BaseLink id="created-chain" href={`/chains`}>
      {id}
    </BaseLink>
  </>
)

function apiCall({
  chainID,
  config,
}: {
  chainID: string
  config: Record<string, JSONPrimitive>
}): Promise<ApiResponse<Chain>> {
  const definition: CreateChainRequest = { chainID, config }
  return api.v2.chains.createChain(definition)
}

export const New = ({
  classes,
}: {
  classes: WithStyles<typeof styles>['classes']
}) => {
  const dispatch = useDispatch()
  const [overrides, setOverrides] = useState<string>('{}')
  const [chainID, setChainID] = useState<string>('')
  const [overridesErrorMsg, setOverridesErrorMsg] = useState<string>('')
  const [chainIDErrorMsg, setChainIDErrorMsg] = useState<string>('')
  const [loading, setLoading] = useState<boolean>(false)

  function validate({
    chainID,
    overrides,
  }: {
    chainID: string
    overrides: string
  }) {
    let valid = true
    if (!(parseInt(chainID, 10) > 0)) {
      setChainIDErrorMsg('Invalid chain ID')
      valid = false
    }
    try {
      JSON.parse(overrides)
    } catch (e) {
      setOverridesErrorMsg('Invalid job spec')
      valid = false
    }
    return valid
  }

  function handleChainIDChange(event: React.ChangeEvent<HTMLTextAreaElement>) {
    setChainID(event.target.value)
    setChainIDErrorMsg('')
  }

  function handleOverrideChange(event: React.ChangeEvent<HTMLTextAreaElement>) {
    setOverrides(event.target.value)
    setOverridesErrorMsg('')
  }

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const isValid = validate({ chainID, overrides })

    // Use the name of the input to actually update the `overrides` object
    // call: event.currentTarget.name

    if (isValid) {
      setLoading(true)
      apiCall({
        chainID,
        config: JSON.parse(overrides),
      })
        .then(({ data }) => {
          dispatch(notifySuccess(SuccessNotification, data))
        })
        .catch((error) => {
          dispatch(notifyError(ErrorMessage, error))
          if (error instanceof BadRequestError) {
            setChainIDErrorMsg('Invalid ChainID')
          }
        })
        .finally(() => {
          setLoading(false)
        })
    }
  }

  return (
    <Content>
      <Grid container spacing={40}>
        <Grid item xs={12}>
          <Card>
            <CardHeader title="New Chain" />
            <CardContent>
              <form noValidate onSubmit={handleSubmit}>
                <Grid container>
                  <Grid item xs={12}>
                    <TextField
                      error={Boolean(chainIDErrorMsg)}
                      helperText={Boolean(chainIDErrorMsg) && chainIDErrorMsg}
                      label="Chain ID"
                      name="ID"
                      placeholder="ID"
                      value={chainID}
                      onChange={handleChainIDChange}
                    />
                  </Grid>
                  <Grid item xs={12}>
                    <Grid item xs={12}>
                      <FormLabel>Config Overrides</FormLabel>
                    </Grid>

                    <Grid item xs={3}>
                      <TextField
                        label="Block History Estimator Block Delay"
                        name="BlockHistoryEstimatorBlockDelay"
                        placeholder="BlockHistoryEstimatorBlockDelay"
                        type="number"
                        fullWidth
                        onChange={handleOverrideChange}
                      />
                    </Grid>

                    <Grid item xs={3}>
                      <TextField
                        label="Block History Estimator Block History Size"
                        name="BlockHistoryEstimatorBlockHistorySize"
                        placeholder="BlockHistoryEstimatorBlockHistorySize"
                        type="number"
                        fullWidth
                        onChange={handleOverrideChange}
                      />
                    </Grid>

                    <Grid item xs={3}>
                      <TextField
                        label="Eth Tx Reaper Threshold"
                        name="EthTxReaperThreshold"
                        placeholder="EthTxReaperThreshold"
                        type="text"
                        fullWidth
                        onChange={handleOverrideChange}
                      />
                    </Grid>

                    <Grid item xs={3}>
                      <TextField
                        label="Eth Tx Resend After Threshold"
                        name="EthTxResendAfterThreshold"
                        placeholder="EthTxResendAfterThreshold"
                        type="text"
                        fullWidth
                        onChange={handleOverrideChange}
                      />
                    </Grid>

                    <Grid item xs={3}>
                      <Checkbox
                        name="EvmEIP1559DynamicFees"
                        // onChange={handleOverrideChange}
                      />
                    </Grid>

                    <Grid item xs={3}>
                      <TextField
                        label="Evm Finality Depth"
                        name="EvmFinalityDepth"
                        placeholder="EvmFinalityDepth"
                        type="number"
                        fullWidth
                        onChange={handleOverrideChange}
                      />
                    </Grid>

                    <Grid item xs={3}>
                      <TextField
                        label="Evm Gas Bump Percent"
                        name="EvmGasBumpPercent"
                        placeholder="EvmGasBumpPercent"
                        type="number"
                        fullWidth
                        onChange={handleOverrideChange}
                      />
                    </Grid>

                    <Grid item xs={3}>
                      <TextField
                        label="Evm Gas Bump Tx Depth"
                        name="EvmGasBumpTxDepth"
                        placeholder="EvmGasBumpTxDepth"
                        type="number"
                        fullWidth
                        onChange={handleOverrideChange}
                      />
                    </Grid>

                    <Grid item xs={3}>
                      <TextField
                        label="Evm Gas Bump Wei"
                        name="EvmGasBumpWei"
                        placeholder="EvmGasBumpWei"
                        type="number"
                        fullWidth
                        onChange={handleOverrideChange}
                      />
                    </Grid>

                    <Grid item xs={3}>
                      <TextField
                        label="Evm Gas Limit Default"
                        name="EvmGasLimitDefault"
                        placeholder="EvmGasLimitDefault"
                        type="number"
                        fullWidth
                        onChange={handleOverrideChange}
                      />
                    </Grid>

                    <Grid item xs={3}>
                      <TextField
                        label="Evm Gas Limit Multiplier"
                        name="EvmGasLimitMultiplier"
                        placeholder="EvmGasLimitMultiplier"
                        type="number"
                        fullWidth
                        onChange={handleOverrideChange}
                      />
                    </Grid>

                    <Grid item xs={3}>
                      <TextField
                        label="Evm Gas Price Default"
                        name="EvmGasPriceDefault"
                        placeholder="EvmGasPriceDefault"
                        type="number"
                        fullWidth
                        onChange={handleOverrideChange}
                      />
                    </Grid>

                    <Grid item xs={3}>
                      <TextField
                        label="Evm Gas Tip Cap Default"
                        name="EvmGasTipCapDefault"
                        placeholder="EvmGasTipCapDefault"
                        type="number"
                        fullWidth
                        onChange={handleOverrideChange}
                      />
                    </Grid>

                    <Grid item xs={3}>
                      <TextField
                        label="Evm Gas Tip Cap Minimum"
                        name="EvmGasTipCapMinimum"
                        placeholder="EvmGasTipCapMinimum"
                        type="number"
                        fullWidth
                        onChange={handleOverrideChange}
                      />
                    </Grid>

                    <Grid item xs={3}>
                      <TextField
                        label="Evm Head Tracker History Depth"
                        name="EvmHeadTrackerHistoryDepth"
                        placeholder="EvmHeadTrackerHistoryDepth"
                        type="number"
                        fullWidth
                        onChange={handleOverrideChange}
                      />
                    </Grid>

                    <Grid item xs={3}>
                      <TextField
                        label="Evm Head Tracker Max Buffer Size"
                        name="EvmHeadTrackerMaxBufferSize"
                        placeholder="EvmHeadTrackerMaxBufferSize"
                        type="number"
                        fullWidth
                        onChange={handleOverrideChange}
                      />
                    </Grid>

                    <Grid item xs={3}>
                      <TextField
                        label="Evm Head Tracker Sampling Interval"
                        name="EvmHeadTrackerSamplingInterval"
                        placeholder="EvmHeadTrackerSamplingInterval"
                        type="text"
                        fullWidth
                        onChange={handleOverrideChange}
                      />
                    </Grid>

                    <Grid item xs={3}>
                      <TextField
                        label="Evm Log Backfill Batch Size"
                        name="EvmLogBackfillBatchSize"
                        placeholder="EvmLogBackfillBatchSize"
                        type="number"
                        fullWidth
                        onChange={handleOverrideChange}
                      />
                    </Grid>

                    <Grid item xs={3}>
                      <TextField
                        label="Evm Max Gas Price Wei"
                        name="EvmMaxGasPriceWei"
                        placeholder="EvmMaxGasPriceWei"
                        type="number"
                        fullWidth
                        onChange={handleOverrideChange}
                      />
                    </Grid>

                    <Grid item xs={3}>
                      <Checkbox
                        name="EvmNonceAutoSync"
                        // onChange={handleOverrideChange}
                      />
                    </Grid>
                  </Grid>
                  <Grid item xs={12}>
                    <Button
                      data-testid="new-chain-config-submit"
                      variant="primary"
                      type="submit"
                      size="large"
                      disabled={
                        loading ||
                        Boolean(overridesErrorMsg) ||
                        Boolean(chainIDErrorMsg)
                      }
                    >
                      Create Chain
                      {loading && (
                        <CircularProgress
                          className={classes.loader}
                          size={30}
                          color="primary"
                        />
                      )}
                    </Button>
                  </Grid>
                </Grid>
              </form>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Content>
  )
}

export default withStyles(styles)(New)
