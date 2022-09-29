import React from 'react'

import Button from '@material-ui/core/Button'
import Dialog, { DialogProps } from '@material-ui/core/Dialog'
import MuiDialogActions from '@material-ui/core/DialogActions'
import DialogContent from '@material-ui/core/DialogContent'
import DialogContentText from '@material-ui/core/DialogContentText'
import MuiDialogTitle from '@material-ui/core/DialogTitle'
import { withStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'

const DialogTitle = withStyles((theme) => ({
  root: {
    paddingTop: theme.spacing.unit * 3,
    paddingBottom: theme.spacing.unit * 3,
  },
}))(MuiDialogTitle)

const DialogActions = withStyles((theme) => ({
  root: {
    paddingLeft: theme.spacing.unit * 2,
    paddingRight: theme.spacing.unit * 2,
    paddingBottom: theme.spacing.unit * 2,
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
  },
}))(MuiDialogActions)

type Props = Pick<DialogProps, 'open' | 'onClose' | 'maxWidth'> & {
  body: string | React.ReactNode
  confirmButtonText?: string
  cancelButtonText?: string
  title: string
  onConfirm: () => void
  onCancel?: () => void
}

export const ConfirmationDialog: React.FC<Props> = ({
  body,
  cancelButtonText,
  confirmButtonText,
  maxWidth,
  onClose,
  onCancel,
  onConfirm,
  open,
  title,
}) => {
  return (
    <Dialog open={open} onClose={onClose} maxWidth={maxWidth}>
      <DialogTitle disableTypography>
        <Typography variant="h5"> {title}</Typography>
      </DialogTitle>
      <DialogContent>
        {typeof body === 'string' ? (
          <DialogContentText color="textPrimary">{body}</DialogContentText>
        ) : (
          body
        )}
      </DialogContent>
      <DialogActions>
        {cancelButtonText && onCancel && (
          <Button onClick={onCancel} variant="text">
            {cancelButtonText}
          </Button>
        )}
        <Button onClick={onConfirm} variant="contained" color="primary">
          {confirmButtonText || 'Confirm'}
        </Button>
      </DialogActions>
    </Dialog>
  )
}
