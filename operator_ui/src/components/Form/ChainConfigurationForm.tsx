import React from 'react'
import {
  Field,
  FieldAttributes,
  Form,
  Formik,
  FormikHelpers,
  useFormikContext,
} from 'formik'
import { CheckboxWithLabel, TextField } from 'formik-material-ui'
import * as Yup from 'yup'

import Button from '@material-ui/core/Button'
import Grid from '@material-ui/core/Grid'
import MenuItem from '@material-ui/core/MenuItem'
import Paper from '@material-ui/core/Paper'
import {
  createStyles,
  Theme,
  WithStyles,
  withStyles,
} from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'

export type FormValues = {
  chainID: string
  chainType: string
  accountAddr: string
  adminAddr: string
  fluxMonitorEnabled: boolean
  ocr1Enabled: boolean
  ocr1IsBootstrap: boolean
  ocr1Multiaddr?: string | null
  ocr1P2PPeerID?: string | null
  ocr1KeyBundleID?: string | null
  ocr2Enabled: boolean
}

const ValidationSchema = Yup.object().shape({
  chainID: Yup.string().required('Required'),
  chainType: Yup.string().required('Required'),
  accountAddr: Yup.string().required('Required'),
  adminAddr: Yup.string().required('Required'),
  ocr1Multiaddr: Yup.string()
    .when(['ocr1Enabled', 'ocr1IsBootstrap'], {
      is: (enabled: boolean, isBootstrap: boolean) => enabled && isBootstrap,
      then: Yup.string().required('Required').nullable(),
    })
    .nullable(),
  ocr1P2PPeerID: Yup.string()
    .when(['ocr1Enabled', 'ocr1IsBootstrap'], {
      is: (enabled: boolean, isBootstrap: boolean) => enabled && !isBootstrap,
      then: Yup.string().required('Required').nullable(),
    })
    .nullable(),
  ocr1KeyBundleID: Yup.string()
    .when(['ocr1Enabled', 'ocr1IsBootstrap'], {
      is: (enabled: boolean, isBootstrap: boolean) => enabled && !isBootstrap,
      then: Yup.string().required('Required').nullable(),
    })
    .nullable(),
})

const styles = (theme: Theme) => {
  return createStyles({
    supportedJobOptionsPaper: {
      padding: theme.spacing.unit * 2,
    },
  })
}

export interface Props extends WithStyles<typeof styles> {
  editing?: boolean
  initialValues: FormValues
  innerRef?: any
  onSubmit: (
    values: FormValues,
    formikHelpers: FormikHelpers<FormValues>,
  ) => void | Promise<any>
  chainIDs: string[]
  accounts: ReadonlyArray<EthKeysPayload_ResultsFields>
  p2pKeys: ReadonlyArray<P2PKeysPayload_ResultsFields>
  ocrKeys: ReadonlyArray<OcrKeyBundlesPayload_ResultsFields>
  showSubmit?: boolean
}

// ChainConfigurationForm is used to create/edit the supported chain
// configurations for the Feeds Manager.
export const ChainConfigurationForm = withStyles(styles)(
  ({
    classes,
    editing = false,
    innerRef,
    initialValues,
    onSubmit,
    chainIDs = [],
    accounts = [],
    p2pKeys = [],
    ocrKeys = [],
    showSubmit = false,
  }: Props) => {
    return (
      <Formik
        innerRef={innerRef}
        initialValues={initialValues}
        validationSchema={ValidationSchema}
        onSubmit={onSubmit}
      >
        {({ values }) => {
          const chainAccounts = accounts.filter(
            (acc) => acc.chain.id == values.chainID && !acc.isFunding,
          )

          return (
            <Form
              data-testid="feeds-manager-form"
              id="chain-configuration-form"
              noValidate
            >
              <Grid container spacing={16}>
                <Grid item xs={12} md={6}>
                  <Field
                    component={TextField}
                    id="chainType"
                    name="chainType"
                    label="Chain Type"
                    select
                    required
                    fullWidth
                    disabled
                    helperText="Only EVM is currently supported"
                  >
                    <MenuItem key="EVM" value="EVM">
                      EVM
                    </MenuItem>
                  </Field>
                </Grid>

                <Grid item xs={12} md={6}>
                  <Field
                    component={TextField}
                    id="chainID"
                    name="chainID"
                    label="Chain ID"
                    required
                    fullWidth
                    select
                    disabled={editing}
                    inputProps={{ 'data-testid': 'chainID-input' }}
                    FormHelperTextProps={{
                      'data-testid': 'chainID-helper-text',
                    }}
                  >
                    {chainIDs.map((chainID) => (
                      <MenuItem key={chainID} value={chainID}>
                        {chainID}
                      </MenuItem>
                    ))}
                  </Field>
                </Grid>

                <Grid item xs={12} md={6}>
                  <AccountAddrField
                    component={TextField}
                    id="accountAddr"
                    name="accountAddr"
                    label="Account Address"
                    inputProps={{ 'data-testid': 'accountAddr-input' }}
                    required
                    fullWidth
                    select
                    helperText="The account address used for this chain"
                    FormHelperTextProps={{
                      'data-testid': 'accountAddr-helper-text',
                    }}
                  >
                    {chainAccounts.map((account) => (
                      <MenuItem key={account.address} value={account.address}>
                        {account.address}
                      </MenuItem>
                    ))}
                  </AccountAddrField>
                </Grid>

                <Grid item xs={12} md={6}>
                  <Field
                    component={TextField}
                    id="adminAddr"
                    name="adminAddr"
                    label="Admin Address"
                    required
                    fullWidth
                    helperText="The address used for LINK payments"
                    FormHelperTextProps={{
                      'data-testid': 'adminAddr-helper-text',
                    }}
                  />
                </Grid>

                <Grid item xs={12}>
                  <Typography>Supported Job Types</Typography>
                </Grid>

                <Grid item xs={12}>
                  <Field
                    component={CheckboxWithLabel}
                    name="fluxMonitorEnabled"
                    type="checkbox"
                    Label={{
                      label: 'Flux Monitor',
                    }}
                  />
                </Grid>

                <Grid item xs={12}>
                  <Field
                    component={CheckboxWithLabel}
                    name="ocr1Enabled"
                    type="checkbox"
                    Label={{
                      label: 'OCR',
                    }}
                  />

                  {values.ocr1Enabled && (
                    <Paper className={classes.supportedJobOptionsPaper}>
                      <Grid container spacing={8}>
                        <>
                          <Grid item xs={12}>
                            <Field
                              component={CheckboxWithLabel}
                              name="ocr1IsBootstrap"
                              type="checkbox"
                              Label={{
                                label:
                                  'Is this node running as a bootstrap peer?',
                              }}
                            />
                          </Grid>

                          {values.ocr1IsBootstrap ? (
                            <Grid item xs={12}>
                              <Field
                                component={TextField}
                                id="ocr1Multiaddr"
                                name="ocr1Multiaddr"
                                label="Multiaddr"
                                required
                                fullWidth
                                helperText="The OCR Multiaddr which oracles use to query for network information"
                                FormHelperTextProps={{
                                  'data-testid': 'ocr1Multiaddr-helper-text',
                                }}
                              />
                            </Grid>
                          ) : (
                            <>
                              <Grid item xs={12} md={6}>
                                <Field
                                  component={TextField}
                                  id="ocr1P2PPeerID"
                                  name="ocr1P2PPeerID"
                                  label="Peer ID"
                                  select
                                  required
                                  fullWidth
                                  helperText="The Peer ID used for this chain"
                                  FormHelperTextProps={{
                                    'data-testid': 'ocr1P2PPeerID-helper-text',
                                  }}
                                >
                                  {p2pKeys.map((key) => (
                                    <MenuItem
                                      key={key.peerID}
                                      value={key.peerID}
                                    >
                                      {key.peerID}
                                    </MenuItem>
                                  ))}
                                </Field>
                              </Grid>

                              <Grid item xs={12} md={6}>
                                <Field
                                  component={TextField}
                                  id="ocr1KeyBundleID"
                                  name="ocr1KeyBundleID"
                                  label="Key Bundle ID"
                                  select
                                  required
                                  fullWidth
                                  helperText="The OCR Key Bundle ID used for this chain"
                                  FormHelperTextProps={{
                                    'data-testid':
                                      'ocr1KeyBundleID-helper-text',
                                  }}
                                >
                                  {ocrKeys.map((key) => (
                                    <MenuItem key={key.id} value={key.id}>
                                      {key.id}
                                    </MenuItem>
                                  ))}
                                </Field>
                              </Grid>
                            </>
                          )}
                        </>
                      </Grid>
                    </Paper>
                  )}
                </Grid>

                {showSubmit && (
                  <Grid item xs={12} md={7}>
                    <Button
                      variant="contained"
                      color="primary"
                      type="submit"
                      size="large"
                    >
                      Submit
                    </Button>
                  </Grid>
                )}
              </Grid>
            </Form>
          )
        }}
      </Formik>
    )
  },
)

// A custom account address field which clears the input based on the chain id
// value changoing
const AccountAddrField = (props: FieldAttributes<any>) => {
  const {
    values: { chainID, accountAddr },
    setFieldValue,
  } = useFormikContext<FormValues>()

  const prevChainID = React.useRef<string>()
  React.useEffect(() => {
    prevChainID.current = chainID
  }, [chainID])

  React.useEffect(() => {
    if (chainID !== prevChainID.current) {
      setFieldValue(props.name, '')
    }
  }, [chainID, setFieldValue, accountAddr, props.name])

  return <Field {...props} />
}
