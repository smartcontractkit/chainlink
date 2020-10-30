import React from 'react'
import Button from 'components/Button'
import Dialog from '@material-ui/core/Dialog'
import DialogActions from '@material-ui/core/DialogActions'
import DialogContent from '@material-ui/core/DialogContent'
import DialogTitle from '@material-ui/core/DialogTitle'
import {
  createStyles,
  withStyles,
  WithStyles,
  Theme,
} from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'

const styles = (theme: Theme) =>
  createStyles({
    button: {
      margin: theme.spacing.unit,
    },
    keyValue: {
      borderRadius: theme.shape.borderRadius,
      border: `1px solid ${theme.palette.secondary.main}`,
      padding: theme.spacing.unit,
      marginTop: theme.spacing.unit * 2,
      color: theme.palette.secondary.main,
    },
  })

type Props = WithStyles<typeof styles> & {
  onDelete: Function
  keyId: string
  keyValue: string
}

export const Delete = withStyles(styles)(
  ({ classes, onDelete, keyId, keyValue }: Props) => {
    const [open, setOpen] = React.useState(false)
    return (
      <>
        <Button
          data-testid="keys-ocr-delete-dialog"
          onClick={() => setOpen(true)}
          variant="danger"
          size="medium"
        >
          Delete
        </Button>

        <Dialog open={open}>
          <DialogTitle>Confirm Delete</DialogTitle>
          <DialogContent>
            <Typography variant={'h5'} color="textSecondary">
              Delete this key?
            </Typography>
            <div className={classes.keyValue}>
              <Typography>{keyValue}</Typography>
            </div>
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setOpen(false)} className={classes.button}>
              No
            </Button>
            <Button
              data-testid="keys-ocr-delete-confirm"
              onClick={() => onDelete(keyId)}
              variant="danger"
              className={classes.button}
            >
              Yes
            </Button>
          </DialogActions>
        </Dialog>
      </>
    )
  },
)
