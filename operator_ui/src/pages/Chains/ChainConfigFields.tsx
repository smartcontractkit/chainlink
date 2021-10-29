import {
  Checkbox,
  FormControlLabel,
  FormLabel,
  Grid,
  TextField,
} from '@material-ui/core'
import React, { useState } from 'react'

interface Props {
  initialValues?: ConfigOverrides
  onChange: (values: ConfigOverrides, error: string) => void
}

export interface ConfigOverrides {
  [attr: string]: string | boolean
}

const defaultKeySpecifics = '{}'

export const ChainConfigFields: React.FunctionComponent<Props> = ({
  initialValues,
  onChange,
}) => {
  const init: ConfigOverrides = { ...(initialValues || {}) }

  const [overrides, setOverrides] = useState<ConfigOverrides>({})
  const [keySpecificOverrides, setKeySpecificOverrides] = useState<string>('{}')
  const [keySpecificOverridesErrorMsg, setKeySpecificOverridesErrorMsg] =
    useState<string>('')

  function validate(keySpecificOverrides: string) {
    try {
      JSON.parse(keySpecificOverrides)
      setKeySpecificOverridesErrorMsg('')

      return true
    } catch (e) {
      setKeySpecificOverridesErrorMsg('Invalid key specific overrides')
    }

    return false
  }

  function handleOverrideChange(event: React.ChangeEvent<HTMLInputElement>) {
    const newOverrides = {
      ...overrides,
      [event.target.name]:
        // Supports setting boolean values, since the checked status is not available on `target.value`
        event.target.type == 'checkbox'
          ? event.target.checked
          : event.target.value,
    }

    setOverrides(newOverrides)

    const isValid = validate(keySpecificOverrides)

    if (isValid) {
      const config: ConfigOverrides = {
        ...overrides,
      }

      if (keySpecificOverrides != defaultKeySpecifics) {
        config.KeySpecific = JSON.parse(keySpecificOverrides)
      }
    }

    onChange(newOverrides, keySpecificOverridesErrorMsg)
  }

  function handleKeySpecificOverrideChange(
    event: React.ChangeEvent<HTMLTextAreaElement>,
  ) {
    setKeySpecificOverrides(event.target.value)
    setKeySpecificOverridesErrorMsg('')
  }

  return (
    <>
      <Grid item xs={12} style={{ marginTop: 10 }}>
        <FormLabel>Config Overrides</FormLabel>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Block History Estimator Block Delay"
            name="BlockHistoryEstimatorBlockDelay"
            placeholder="BlockHistoryEstimatorBlockDelay"
            type="number"
            value={init['BlockHistoryEstimatorBlockDelay']}
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Block History Estimator Block History Size"
            name="BlockHistoryEstimatorBlockHistorySize"
            placeholder="BlockHistoryEstimatorBlockHistorySize"
            type="number"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Eth Tx Reaper Threshold"
            name="EthTxReaperThreshold"
            placeholder="EthTxReaperThreshold"
            type="text"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Eth Tx Resend After Threshold"
            name="EthTxResendAfterThreshold"
            placeholder="EthTxResendAfterThreshold"
            type="text"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <FormControlLabel
            control={
              <Checkbox
                name="EvmEIP1559DynamicFees"
                onChange={(event) => handleOverrideChange(event)}
              />
            }
            label="EvmEIP1559DynamicFees"
          />
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Evm Finality Depth"
            name="EvmFinalityDepth"
            placeholder="EvmFinalityDepth"
            type="number"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Evm Gas Bump Percent"
            name="EvmGasBumpPercent"
            placeholder="EvmGasBumpPercent"
            type="number"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Evm Gas Bump Tx Depth"
            name="EvmGasBumpTxDepth"
            placeholder="EvmGasBumpTxDepth"
            type="number"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Evm Gas Bump Wei"
            name="EvmGasBumpWei"
            placeholder="EvmGasBumpWei"
            type="number"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Evm Gas Limit Default"
            name="EvmGasLimitDefault"
            placeholder="EvmGasLimitDefault"
            type="number"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Evm Gas Limit Multiplier"
            name="EvmGasLimitMultiplier"
            placeholder="EvmGasLimitMultiplier"
            type="number"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Evm Gas Price Default"
            name="EvmGasPriceDefault"
            placeholder="EvmGasPriceDefault"
            type="number"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Evm Gas Tip Cap Default"
            name="EvmGasTipCapDefault"
            placeholder="EvmGasTipCapDefault"
            type="number"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Evm Gas Tip Cap Minimum"
            name="EvmGasTipCapMinimum"
            placeholder="EvmGasTipCapMinimum"
            type="number"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Evm Head Tracker History Depth"
            name="EvmHeadTrackerHistoryDepth"
            placeholder="EvmHeadTrackerHistoryDepth"
            type="number"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Evm Head Tracker Max Buffer Size"
            name="EvmHeadTrackerMaxBufferSize"
            placeholder="EvmHeadTrackerMaxBufferSize"
            type="number"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Evm Head Tracker Sampling Interval"
            name="EvmHeadTrackerSamplingInterval"
            placeholder="EvmHeadTrackerSamplingInterval"
            type="text"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Evm Log Backfill Batch Size"
            name="EvmLogBackfillBatchSize"
            placeholder="EvmLogBackfillBatchSize"
            type="number"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Evm Max Gas Price Wei"
            name="EvmMaxGasPriceWei"
            placeholder="EvmMaxGasPriceWei"
            type="number"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <FormControlLabel
            control={
              <Checkbox
                name="EvmNonceAutoSync"
                // onChange={handleOverrideChange}
              />
            }
            label="Evm Nonce Auto Sync"
          />
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Evm RPC Default Batch Size"
            name="EvmRPCDefaultBatchSize"
            placeholder="EvmRPCDefaultBatchSize"
            type="number"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Flags Contract Address"
            name="FlagsContractAddress"
            placeholder="FlagsContractAddress"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Gas Estimator Mode"
            name="GasEstimatorMode"
            placeholder="GasEstimatorMode"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Chain Type"
            name="ChainType"
            placeholder="ChainType"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Min Incoming Confirmations"
            name="MinIncomingConfirmations"
            placeholder="MinIncomingConfirmations"
            type="number"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Min Required Outgoing Confirmations"
            name="MinRequiredOutgoingConfirmations"
            placeholder="MinRequiredOutgoingConfirmations"
            type="number"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Minimum Contract Payment"
            name="MinimumContractPayment"
            placeholder="MinimumContractPayment"
            type="number"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="OCR Observation Timeout"
            name="OCRObservationTimeout"
            placeholder="OCRObservationTimeout"
            type="number"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
      </Grid>

      <Grid item xs={12} style={{ marginTop: 10 }}>
        <FormLabel>Key Specific Config Overrides</FormLabel>
      </Grid>

      <Grid item xs={12}>
        <TextField
          error={Boolean(keySpecificOverridesErrorMsg)}
          value={keySpecificOverrides}
          onChange={handleKeySpecificOverrideChange}
          helperText={
            Boolean(keySpecificOverridesErrorMsg) &&
            keySpecificOverridesErrorMsg
          }
          autoComplete="off"
          label={'JSON'}
          rows={10}
          rowsMax={25}
          placeholder={'Paste JSON'}
          multiline
          margin="normal"
          name="KeySpecific"
          id="chainConfig"
          variant="outlined"
          fullWidth
        />
      </Grid>
    </>
  )
}

export default ChainConfigFields
