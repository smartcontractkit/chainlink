import {
  Checkbox,
  FormControlLabel,
  FormLabel,
  Grid,
  MenuItem,
  TextField,
  Typography,
} from '@material-ui/core'
import React, { useState } from 'react'
import { capitalize } from '@material-ui/core/utils/helpers'

interface Props {
  initialValues?: ConfigOverrides
  onChange: (values: ConfigOverrides, error: string) => void
}

export interface ConfigOverrides {
  [attr: string]: JSONPrimitive
}

type NonNullableJSONPrimitive = string | number | boolean

const defaultKeySpecifics = '{}'

const chainTypes = ['arbitrum', 'exchain', 'optimism', 'xdai']
const gasEstimatorModes = [
  'BlockHistory',
  'FixedPrice',
  'Optimism',
  'Optimism2',
]

export const ChainConfigFields: React.FunctionComponent<Props> = ({
  initialValues,
  onChange,
}) => {
  const [overrides, setOverrides] = useState<ConfigOverrides>({
    ...(initialValues || {}),
  })
  const [keySpecificOverrides, setKeySpecificOverrides] = useState<string>(
    (overrides['KeySpecific'] as string) || '{}',
  )
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

    // Removes empty or false default values from the form
    if (
      newOverrides[event.target.name] === '' ||
      newOverrides[event.target.name] === false
    ) {
      delete newOverrides[event.target.name]
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

  function getFieldValue(fieldName: string): NonNullableJSONPrimitive {
    return (overrides[fieldName] as NonNullableJSONPrimitive) || ''
  }

  return (
    <>
      <Grid item xs={12}>
        <FormLabel>Config Overrides</FormLabel>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Block History Estimator Block Delay"
            name="BlockHistoryEstimatorBlockDelay"
            placeholder="BlockHistoryEstimatorBlockDelay"
            type="number"
            value={getFieldValue('BlockHistoryEstimatorBlockDelay')}
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
        <Grid item>
          <Typography color="secondary">Default: 1</Typography>
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Block History Estimator Block History Size"
            name="BlockHistoryEstimatorBlockHistorySize"
            placeholder="BlockHistoryEstimatorBlockHistorySize"
            value={getFieldValue('BlockHistoryEstimatorBlockHistorySize')}
            type="number"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
        <Grid item>
          <Typography color="secondary">Default: 16</Typography>
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Eth Tx Reaper Threshold"
            name="EthTxReaperThreshold"
            placeholder="EthTxReaperThreshold"
            value={getFieldValue('EthTxReaperThreshold')}
            type="text"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
        <Grid item>
          <Typography color="secondary">Default: 168 hours</Typography>
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Eth Tx Resend After Threshold"
            name="EthTxResendAfterThreshold"
            placeholder="EthTxResendAfterThreshold"
            value={getFieldValue('EthTxResendAfterThreshold')}
            type="text"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
        <Grid item>
          <Typography color="secondary">Default: 1 minute</Typography>
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <FormControlLabel
            control={
              <Checkbox
                name="EvmEIP1559DynamicFees"
                value={getFieldValue('EvmEIP1559DynamicFees').toString() || ''}
                checked={
                  Boolean(getFieldValue('EvmEIP1559DynamicFees')) || false
                }
                onChange={(event) => handleOverrideChange(event)}
              />
            }
            label="EvmEIP1559DynamicFees"
          />
        </Grid>
        <Grid item>
          <Typography color="secondary">Default: false</Typography>
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Evm Finality Depth"
            name="EvmFinalityDepth"
            placeholder="EvmFinalityDepth"
            value={getFieldValue('EvmFinalityDepth')}
            type="number"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
        <Grid item>
          <Typography color="secondary">Default: 50</Typography>
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Evm Gas Bump Percent"
            name="EvmGasBumpPercent"
            placeholder="EvmGasBumpPercent"
            value={getFieldValue('EvmGasBumpPercent')}
            type="number"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
        <Grid item>
          <Typography color="secondary">Default: 20</Typography>
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Evm Gas Bump Tx Depth"
            name="EvmGasBumpTxDepth"
            placeholder="EvmGasBumpTxDepth"
            value={getFieldValue('EvmGasBumpTxDepth')}
            type="number"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
        <Grid item>
          <Typography color="secondary">Default: 10</Typography>
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Evm Gas Bump Wei"
            name="EvmGasBumpWei"
            placeholder="EvmGasBumpWei"
            value={getFieldValue('EvmGasBumpWei')}
            type="number"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
        <Grid item>
          <Typography color="secondary">Default: 10</Typography>
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Evm Gas Limit Default"
            name="EvmGasLimitDefault"
            placeholder="EvmGasLimitDefault"
            value={getFieldValue('EvmGasLimitDefault')}
            type="number"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
        <Grid item>
          <Typography color="secondary">Default: 500000</Typography>
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Evm Gas Limit Multiplier"
            name="EvmGasLimitMultiplier"
            placeholder="EvmGasLimitMultiplier"
            value={getFieldValue('EvmGasLimitMultiplier')}
            type="number"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
        <Grid item>
          <Typography color="secondary">Default: 1.0</Typography>
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Evm Gas Price Default"
            name="EvmGasPriceDefault"
            placeholder="EvmGasPriceDefault"
            value={getFieldValue('EvmGasPriceDefault')}
            type="number"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
        <Grid item>
          <Typography color="secondary">Default: 20 gwei</Typography>
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Evm Gas Tip Cap Default"
            name="EvmGasTipCapDefault"
            placeholder="EvmGasTipCapDefault"
            value={getFieldValue('EvmGasTipCapDefault')}
            type="number"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
        <Grid item>
          <Typography color="secondary">Default: 0</Typography>
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Evm Gas Tip Cap Minimum"
            name="EvmGasTipCapMinimum"
            placeholder="EvmGasTipCapMinimum"
            value={getFieldValue('EvmGasTipCapMinimum')}
            type="number"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
        <Grid item>
          <Typography color="secondary">Default: 0</Typography>
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Evm Head Tracker History Depth"
            name="EvmHeadTrackerHistoryDepth"
            placeholder="EvmHeadTrackerHistoryDepth"
            value={getFieldValue('EvmHeadTrackerHistoryDepth')}
            type="number"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
        <Grid item>
          <Typography color="secondary">Default: 100</Typography>
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Evm Head Tracker Max Buffer Size"
            name="EvmHeadTrackerMaxBufferSize"
            placeholder="EvmHeadTrackerMaxBufferSize"
            value={getFieldValue('EvmHeadTrackerMaxBufferSize')}
            type="number"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
        <Grid item>
          <Typography color="secondary">Default: 3</Typography>
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Evm Head Tracker Sampling Interval"
            name="EvmHeadTrackerSamplingInterval"
            placeholder="EvmHeadTrackerSamplingInterval"
            value={getFieldValue('EvmHeadTrackerSamplingInterval')}
            type="text"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
        <Grid item>
          <Typography color="secondary">Default: 1 second</Typography>
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Evm Log Backfill Batch Size"
            name="EvmLogBackfillBatchSize"
            placeholder="EvmLogBackfillBatchSize"
            value={getFieldValue('EvmLogBackfillBatchSize')}
            type="number"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
        <Grid item>
          <Typography color="secondary">Default: 100</Typography>
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Evm Max Gas Price Wei"
            name="EvmMaxGasPriceWei"
            placeholder="EvmMaxGasPriceWei"
            value={getFieldValue('EvmMaxGasPriceWei')}
            type="number"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
        <Grid item>
          <Typography color="secondary">Default: 5000 gwei</Typography>
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <FormControlLabel
            control={
              <Checkbox
                name="EvmNonceAutoSync"
                value={getFieldValue('EvmNonceAutoSync').toString() || ''}
                checked={Boolean(getFieldValue('EvmNonceAutoSync')) || false}
                onChange={handleOverrideChange}
              />
            }
            label="Evm Nonce Auto Sync"
          />
        </Grid>
        <Grid item>
          <Typography color="secondary">Default: true</Typography>
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Evm RPC Default Batch Size"
            name="EvmRPCDefaultBatchSize"
            placeholder="EvmRPCDefaultBatchSize"
            value={getFieldValue('EvmRPCDefaultBatchSize')}
            type="number"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
        <Grid item>
          <Typography color="secondary">Default: 100</Typography>
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Flags Contract Address"
            name="FlagsContractAddress"
            placeholder="FlagsContractAddress"
            value={getFieldValue('FlagsContractAddress')}
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
        <Grid item>
          <Typography color="secondary">Default: empty</Typography>
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            name="GasEstimatorMode"
            fullWidth
            select
            label="Gas Estimator Mode"
            value={getFieldValue('GasEstimatorMode') || ''}
            onChange={handleOverrideChange}
          >
            {gasEstimatorModes.map((mode) => (
              <MenuItem key={mode} value={mode}>
                {capitalize(mode)}
              </MenuItem>
            ))}
          </TextField>
        </Grid>
        <Grid item>
          <Typography color="secondary">Default: BlockHistory</Typography>
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            fullWidth
            select
            label="Chain Type"
            name="ChainType"
            value={getFieldValue('ChainType') || ''}
            onChange={handleOverrideChange}
          >
            {chainTypes.map((type) => (
              <MenuItem key={type} value={type}>
                {capitalize(type)}
              </MenuItem>
            ))}
          </TextField>
        </Grid>
        <Grid item>
          <Typography color="secondary">Default: empty</Typography>
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Min Incoming Confirmations"
            name="MinIncomingConfirmations"
            placeholder="MinIncomingConfirmations"
            value={getFieldValue('MinIncomingConfirmations')}
            type="number"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
        <Grid item>
          <Typography color="secondary">Default: 3</Typography>
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Min Required Outgoing Confirmations"
            name="MinRequiredOutgoingConfirmations"
            placeholder="MinRequiredOutgoingConfirmations"
            value={getFieldValue('MinRequiredOutgoingConfirmations')}
            type="number"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
        <Grid item>
          <Typography color="secondary">Default: 12</Typography>
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Minimum Contract Payment"
            name="MinimumContractPayment"
            placeholder="MinimumContractPayment"
            value={getFieldValue('MinimumContractPayment')}
            type="number"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
        <Grid item>
          <Typography color="secondary">Default: 0.00001 LINK</Typography>
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="OCR Observation Timeout"
            name="OCRObservationTimeout"
            placeholder="OCRObservationTimeout"
            value={getFieldValue('OCRObservationTimeout')}
            type="number"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
        <Grid item>
          <Typography color="secondary">Default: 12 seconds</Typography>
        </Grid>
      </Grid>

      <Grid item xs={6}>
        <Grid item xs={6}>
          <TextField
            label="Link Contract Address"
            name="LinkContractAddress"
            placeholder="LinkContractAddress"
            value={getFieldValue('LinkContractAddress')}
            type="text"
            fullWidth
            onChange={handleOverrideChange}
          />
        </Grid>
        <Grid item>
          <Typography color="secondary">Default: empty</Typography>
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
