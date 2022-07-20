import React from 'react'

import Grid from '@material-ui/core/Grid'
import {
  createStyles,
  Theme,
  WithStyles,
  withStyles,
} from '@material-ui/core/styles'
import TextField from '@material-ui/core/TextField'

const styles = (theme: Theme) => {
  return createStyles({
    textField: {
      marginBottom: theme.spacing.unit * 2,
    },
    input: {
      padding: 14,
      backgroundColor: 'white',
      borderRadius: 4,
    },
  })
}

interface Props extends WithStyles<typeof styles> {
  placeholder?: string
  value: string
  onChange: (value: string) => void
}

export const SearchTextField = withStyles(styles)(
  ({ classes, onChange, placeholder, value }: Props) => {
    return (
      <Grid container spacing={16}>
        <Grid item xs={12} md={6}>
          <TextField
            className={classes.textField}
            inputProps={{
              className: classes.input,
            }}
            InputLabelProps={{
              shrink: true,
            }}
            placeholder={placeholder || 'Search'}
            value={value}
            name="search"
            onChange={(event) => onChange(event.target.value)}
            variant="outlined"
            fullWidth
          />
        </Grid>
      </Grid>
    )
  },
)
