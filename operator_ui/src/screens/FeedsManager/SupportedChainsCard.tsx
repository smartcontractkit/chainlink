import React from 'react'

import { gql, useMutation } from '@apollo/client'
import { useDispatch } from 'react-redux'

import Card from '@material-ui/core/Card'
import CardHeader from '@material-ui/core/CardHeader'
import Chip from '@material-ui/core/Chip'
import Grid from '@material-ui/core/Grid'
import ExpandMoreIcon from '@material-ui/icons/ExpandMore'
import ExpansionPanel from '@material-ui/core/ExpansionPanel'
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary'
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails'
import IconButton from '@material-ui/core/IconButton'
import ListItemIcon from '@material-ui/core/ListItemIcon'
import ListItemText from '@material-ui/core/ListItemText'
import Menu from '@material-ui/core/Menu'
import MenuItem from '@material-ui/core/MenuItem'
import MoreVertIcon from '@material-ui/icons/MoreVert'
import AddBoxIcon from '@material-ui/icons/AddBox'
import {
  createStyles,
  Theme,
  WithStyles,
  withStyles,
} from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'

import { notifySuccessMsg, notifyErrorMsg } from 'actionCreators'
import { FormValues } from 'src/components/Form/ChainConfigurationForm'
import {
  DetailsCardItemTitle,
  DetailsCardItemValue,
} from 'src/components/Cards/DetailsCard'
import { NewSupportedChainDialog } from './NewSupportedChainDialog'
import { useMutationErrorHandler } from 'src/hooks/useMutationErrorHandler'
import Button from 'src/components/Button'
import { EditSupportedChainDialog } from './EditSupportedChainDialog'
import { FEEDS_MANAGERS_WITH_PROPOSALS_QUERY } from 'src/hooks/queries/useFeedsManagersWithProposalsQuery'

export const CREATE_FEEDS_MANAGER_CHAIN_CONFIG_MUTATION = gql`
  mutation CreateFeedsManagerChainConfig(
    $input: CreateFeedsManagerChainConfigInput!
  ) {
    createFeedsManagerChainConfig(input: $input) {
      ... on CreateFeedsManagerChainConfigSuccess {
        chainConfig {
          id
        }
      }
      ... on NotFoundError {
        message
        code
      }
      ... on InputErrors {
        errors {
          path
          message
          code
        }
      }
    }
  }
`

export const DELETE_FEEDS_MANAGER_CHAIN_CONFIG_MUTATION = gql`
  mutation DeleteFeedsManagerChainConfig($id: ID!) {
    deleteFeedsManagerChainConfig(id: $id) {
      ... on DeleteFeedsManagerChainConfigSuccess {
        chainConfig {
          id
        }
      }
      ... on NotFoundError {
        message
        code
      }
    }
  }
`

export const UPDATE_FEEDS_MANAGER_CHAIN_CONFIG_MUTATION = gql`
  mutation UpdateFeedsManagerChainConfig(
    $id: ID!
    $input: UpdateFeedsManagerChainConfigInput!
  ) {
    updateFeedsManagerChainConfig(id: $id, input: $input) {
      __typename
      ... on UpdateFeedsManagerChainConfigSuccess {
        chainConfig {
          id
        }
      }
      ... on NotFoundError {
        message
        code
      }
      ... on InputErrors {
        errors {
          path
          message
          code
        }
      }
    }
  }
`

const styles = (theme: Theme) => {
  return createStyles({
    card: {
      marginBottom: theme.spacing.unit * 2,
    },
    panel: {
      borderBottom: `1px solid ${theme.palette.divider}`,
    },
    panelExpanded: {
      margin: 0,
    },
    chip: {
      marginRight: theme.spacing.unit * 2,
    },
    title: {
      fontSize: theme.typography.body2.fontSize,
      fontWeight: theme.typography.body2.fontWeight,
      color: theme.typography.body2.color,
    },
    panelDetails: {
      display: 'block',
    },
    jobTypeContainer: {
      borderLeft: `2px solid ${theme.palette.primary.main}`,
      marginLeft: -8,
      paddingLeft: 8,
    },
    panelDetailsActions: {
      display: 'flex',
      justifyContent: 'flex-end',
      gap: `${theme.spacing.unit * 2}px`,
    },
  })
}

interface Props extends WithStyles<typeof styles> {
  mgrID: string
  cfgs: ReadonlyArray<FeedsManager_ChainConfigFields>
}

export const SupportedChainsCard = withStyles(styles)(
  ({ classes, cfgs, mgrID }: Props) => {
    const dispatch = useDispatch()
    const { handleMutationError } = useMutationErrorHandler()
    const [newDialogOpen, setNewDialogOpen] = React.useState(false)
    const [editCfg, setEditCfg] =
      React.useState<FeedsManager_ChainConfigFields | null>(null)
    const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null)

    const handleMenuOpen = (event: React.MouseEvent<HTMLButtonElement>) => {
      setAnchorEl(event.currentTarget)
    }

    const handleMenuClose = () => {
      setAnchorEl(null)
    }

    const handleEditDialogOpen = (cfg: FeedsManager_ChainConfigFields) => {
      setEditCfg(cfg)
    }

    const handleEditDialogClose = () => {
      setEditCfg(null)
    }

    const isEditDialogOpen = () => editCfg !== null

    const [createChainConfig] = useMutation<
      CreateFeedsManagerChainConfig,
      CreateFeedsManagerChainConfigVariables
    >(CREATE_FEEDS_MANAGER_CHAIN_CONFIG_MUTATION, {
      refetchQueries: [FEEDS_MANAGERS_WITH_PROPOSALS_QUERY],
    })

    const [deleteChainConfig] = useMutation<
      DeleteFeedsManagerChainConfig,
      DeleteFeedsManagerChainConfigVariables
    >(DELETE_FEEDS_MANAGER_CHAIN_CONFIG_MUTATION, {
      refetchQueries: [FEEDS_MANAGERS_WITH_PROPOSALS_QUERY],
    })

    const [updateChainConfig] = useMutation<
      UpdateFeedsManagerChainConfig,
      UpdateFeedsManagerChainConfigVariables
    >(UPDATE_FEEDS_MANAGER_CHAIN_CONFIG_MUTATION, {
      refetchQueries: [FEEDS_MANAGERS_WITH_PROPOSALS_QUERY],
    })

    const handleCreateSubmit = async (values: FormValues) => {
      try {
        const result = await createChainConfig({
          variables: {
            input: {
              feedsManagerID: mgrID,
              chainID: values.chainID,
              chainType: values.chainType,
              accountAddr: values.accountAddr,
              adminAddr: values.adminAddr,
              fluxMonitorEnabled: values.fluxMonitorEnabled,
              ocr1Enabled: values.ocr1Enabled,
              ocr1IsBootstrap: values.ocr1IsBootstrap,
              ocr1Multiaddr:
                values.ocr1Multiaddr !== '' ? values.ocr1Multiaddr : null,
              ocr1P2PPeerID:
                values.ocr1P2PPeerID !== '' ? values.ocr1P2PPeerID : null,
              ocr1KeyBundleID:
                values.ocr1KeyBundleID != '' ? values.ocr1KeyBundleID : null,
              ocr2Enabled: false, // We don't want to support OCR2 in the UI yet.
            },
          },
        })

        const payload = result.data?.createFeedsManagerChainConfig
        switch (payload?.__typename) {
          case 'CreateFeedsManagerChainConfigSuccess':
            dispatch(notifySuccessMsg('Added new supported chain'))

            break
          case 'NotFoundError':
            dispatch(notifyErrorMsg(payload.message))

            break
        }
      } catch (e) {
        handleMutationError(e)
      }
    }

    const handleDelete = async (id: string) => {
      try {
        const result = await deleteChainConfig({
          variables: { id },
        })

        const payload = result.data?.deleteFeedsManagerChainConfig
        switch (payload?.__typename) {
          case 'DeleteFeedsManagerChainConfigSuccess':
            dispatch(notifySuccessMsg('Deleted supported chain'))

            break
          case 'NotFoundError':
            dispatch(notifyErrorMsg(payload.message))

            break
        }
      } catch (e) {
        handleMutationError(e)
      }
    }

    const handleUpdateSubmit = async (values: FormValues) => {
      if (!editCfg) {
        return
      }

      try {
        const result = await updateChainConfig({
          variables: {
            id: editCfg.id,
            input: {
              accountAddr: values.accountAddr,
              adminAddr: values.adminAddr,
              fluxMonitorEnabled: values.fluxMonitorEnabled,
              ocr1Enabled: values.ocr1Enabled,
              ocr1IsBootstrap: values.ocr1IsBootstrap,
              ocr1Multiaddr:
                values.ocr1Multiaddr !== '' ? values.ocr1Multiaddr : null,
              ocr1P2PPeerID:
                values.ocr1P2PPeerID !== '' ? values.ocr1P2PPeerID : null,
              ocr1KeyBundleID:
                values.ocr1KeyBundleID != '' ? values.ocr1KeyBundleID : null,
              ocr2Enabled: false, // We don't want to support OCR2 in the UI yet.
            },
          },
        })

        const payload = result.data?.updateFeedsManagerChainConfig
        switch (payload?.__typename) {
          case 'UpdateFeedsManagerChainConfigSuccess':
            handleMenuClose()

            dispatch(notifySuccessMsg('Updated supported chain'))

            break
          case 'NotFoundError':
            dispatch(notifyErrorMsg(payload.message))

            break
        }
      } catch (e) {
        handleMutationError(e)
      }
    }

    return (
      <Card className={classes.card}>
        <CardHeader
          title="Supported Chains"
          classes={{
            title: classes.title,
          }}
          action={
            <>
              <IconButton onClick={handleMenuOpen} aria-label="open-menu">
                <MoreVertIcon />
              </IconButton>

              <Menu
                anchorEl={anchorEl}
                open={Boolean(anchorEl)}
                onClose={handleMenuClose}
              >
                <MenuItem
                  onClick={() => {
                    setNewDialogOpen(true)
                    handleMenuClose()
                  }}
                >
                  <ListItemIcon>
                    <AddBoxIcon />
                  </ListItemIcon>
                  <ListItemText>New</ListItemText>
                </MenuItem>
              </Menu>
            </>
          }
        />

        {cfgs.map((cfg, idx) => (
          <ExpansionPanel
            key={idx}
            defaultExpanded={false}
            classes={{
              root: classes.panel,
              expanded: classes.panelExpanded,
            }}
          >
            <ExpansionPanelSummary expandIcon={<ExpandMoreIcon />}>
              <Chip
                label={cfg.chainType}
                color="primary"
                className={classes.chip}
              />
              <Typography style={{ lineHeight: '30px' }}>
                Chain ID: {cfg.chainID}
              </Typography>
            </ExpansionPanelSummary>
            <ExpansionPanelDetails className={classes.panelDetails}>
              <Grid container spacing={40}>
                <Grid item xs={12}>
                  <Grid container>
                    <Grid item xs={12} sm={6}>
                      <DetailsCardItemTitle title="Account Address" />
                      <DetailsCardItemValue value={cfg.accountAddr} />
                    </Grid>
                    <Grid item xs={12} sm={6}>
                      <DetailsCardItemTitle title="Admin Address" />
                      <DetailsCardItemValue value={cfg.adminAddr} />
                    </Grid>

                    <FluxMonitorJobTypeRow cfg={cfg.fluxMonitorJobConfig} />
                    <OCRJobTypeRow cfg={cfg.ocr1JobConfig} />
                    <OCR2JobTypeRow cfg={cfg.ocr2JobConfig} />
                  </Grid>
                </Grid>

                <Grid item xs={12}>
                  <div className={classes.panelDetailsActions}>
                    <Button onClick={() => handleEditDialogOpen(cfg)}>
                      Edit
                    </Button>
                    <Button
                      onClick={() => handleDelete(cfg.id)}
                      variant="danger"
                    >
                      Delete
                    </Button>
                  </div>
                </Grid>
              </Grid>
            </ExpansionPanelDetails>
          </ExpansionPanel>
        ))}

        <NewSupportedChainDialog
          open={newDialogOpen}
          onClose={() => setNewDialogOpen(false)}
          onSubmit={handleCreateSubmit}
        />

        <EditSupportedChainDialog
          cfg={editCfg}
          open={isEditDialogOpen()}
          onClose={handleEditDialogClose}
          onSubmit={handleUpdateSubmit}
        />
      </Card>
    )
  },
)

const jobTypeRowStyles = (theme: Theme) => {
  return createStyles({
    jobTypeContainer: {
      borderLeft: `2px solid ${theme.palette.primary.main}`,
      marginLeft: -8,
      paddingLeft: 8,
    },
  })
}

interface FluxMonitorJobTypeRowProps
  extends WithStyles<typeof jobTypeRowStyles> {
  cfg: FeedsManager_ChainConfigFields['fluxMonitorJobConfig']
}

const FluxMonitorJobTypeRow = withStyles(styles)(
  ({ cfg, classes }: FluxMonitorJobTypeRowProps) => {
    if (!cfg.enabled) {
      return null
    }

    return (
      <Grid item xs={12} sm={1} md={12}>
        <div className={classes.jobTypeContainer}>
          <DetailsCardItemTitle title="Job Type" />
          <DetailsCardItemValue value="Flux Monitor" />
        </div>
      </Grid>
    )
  },
)

interface OCRJobTypeRowProps extends WithStyles<typeof jobTypeRowStyles> {
  cfg: FeedsManager_ChainConfigFields['ocr1JobConfig']
}

const OCRJobTypeRow = withStyles(styles)(
  ({ cfg, classes }: OCRJobTypeRowProps) => {
    if (!cfg.enabled) {
      return null
    }

    const renderBootstrap = () => (
      <Grid item xs={12} sm={1} md={5}>
        <DetailsCardItemTitle title="Multiaddr" />
        <DetailsCardItemValue value={cfg.multiaddr} />
      </Grid>
    )

    const renderOracle = () => (
      <>
        <Grid item xs={12} sm={1} md={5}>
          <DetailsCardItemTitle title="P2P Peer ID" />
          <DetailsCardItemValue value={cfg.p2pPeerID} />
        </Grid>
        <Grid item xs={12} sm={1} md={5}>
          <DetailsCardItemTitle title="OCR Key ID" />
          <DetailsCardItemValue value={cfg.keyBundleID} />
        </Grid>
      </>
    )

    return (
      <>
        <Grid item xs={12} sm={1} md={2}>
          <div className={classes.jobTypeContainer}>
            <DetailsCardItemTitle title="Job Type" />
            <DetailsCardItemValue
              value={`OCR ${cfg.isBootstrap ? '(Bootstrap)' : ''}`}
            />
          </div>
        </Grid>

        {cfg.isBootstrap ? renderBootstrap() : renderOracle()}
      </>
    )
  },
)

interface OCR2JobTypeRowProps extends WithStyles<typeof jobTypeRowStyles> {
  cfg: FeedsManager_ChainConfigFields['ocr2JobConfig']
}

const OCR2JobTypeRow = withStyles(styles)(
  ({ cfg, classes }: OCR2JobTypeRowProps) => {
    if (!cfg.enabled) {
      return null
    }

    const renderBootstrap = () => (
      <Grid item xs={12} sm={1} md={5}>
        <DetailsCardItemTitle title="Multiaddr" />
        <DetailsCardItemValue value={cfg.multiaddr} />
      </Grid>
    )

    const renderOracle = () => (
      <>
        <Grid item xs={12} sm={1} md={5}>
          <DetailsCardItemTitle title="P2P Peer ID" />
          <DetailsCardItemValue value={cfg.p2pPeerID} />
        </Grid>
        <Grid item xs={12} sm={1} md={5}>
          <DetailsCardItemTitle title="OCR Key ID" />
          <DetailsCardItemValue value={cfg.keyBundleID} />
        </Grid>
      </>
    )

    return (
      <>
        <Grid item xs={12} sm={1} md={2}>
          <div className={classes.jobTypeContainer}>
            <DetailsCardItemTitle title="Job Type" />
            <DetailsCardItemValue
              value={`OCR2 ${cfg.isBootstrap ? '(Bootstrap)' : ''}`}
            />
          </div>
        </Grid>

        {cfg.isBootstrap ? renderBootstrap() : renderOracle()}
      </>
    )
  },
)
