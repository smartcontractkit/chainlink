import React from 'react'

import CloseIcon from '@material-ui/icons/Close'
import Dialog from '@material-ui/core/Dialog'
import DialogActions from '@material-ui/core/DialogActions'
import DialogContent from '@material-ui/core/DialogContent'
import DialogTitle from '@material-ui/core/DialogTitle'
import IconButton from '@material-ui/core/IconButton'
import TextField from '@material-ui/core/TextField'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'

import Button from 'src/components/Button'

const styles = (theme: Theme) =>
  createStyles({
    dialogContent: {
      paddingTop: theme.spacing.unit * 1,
    },
    textarea: {
      width: 400,
    },
    closeButton: {
      position: 'absolute',
      right: theme.spacing.unit,
      top: theme.spacing.unit,
      color: theme.palette.grey[500],
    },
  })

interface Props extends WithStyles<typeof styles> {
  open: boolean
  onClose: () => void
  onRun: (pipelineInput: string) => void
}

// TODO - Convert this to use formik
export const RunJobDialog = withStyles(styles)(
  ({ open, onClose, onRun, classes }: Props) => {
    const [pipelineInput, setPipelineInput] = React.useState('')

    const handleRun = React.useCallback(() => {
      onRun(pipelineInput)
      onClose()
    }, [onRun, onClose, pipelineInput])

    return (
      <Dialog onClose={onClose} open={open}>
        <DialogTitle disableTypography>
          <Typography variant="h5">Pipeline Input</Typography>
          <IconButton
            aria-label="Close"
            className={classes.closeButton}
            onClick={onClose}
          >
            <CloseIcon />
          </IconButton>
        </DialogTitle>

        <DialogContent className={classes.dialogContent}>
          <TextField
            className={classes.textarea}
            multiline
            rows={6}
            variant="outlined"
            onChange={(e) => setPipelineInput(e.target.value)}
          />
        </DialogContent>

        <DialogActions>
          <Button variant="primary" onClick={handleRun}>
            Run Job
          </Button>
        </DialogActions>
      </Dialog>
    )
  },
)
