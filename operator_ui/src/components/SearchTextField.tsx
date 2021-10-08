import React from 'react'

import Grid from '@material-ui/core/Grid'
import SearchIcon from '@material-ui/icons/Search'
import TextField from '@material-ui/core/TextField'

interface Props {
  value: string
  onChange: (value: string) => void
}

export const SearchTextField: React.FC<Props> = ({ onChange, value }) => {
  return (
    <Grid container spacing={8} alignItems="flex-end">
      <Grid item>
        <SearchIcon />
      </Grid>
      <Grid item>
        <TextField
          label="Search"
          value={value}
          name="search"
          onChange={(event) => onChange(event.target.value)}
        />
      </Grid>
    </Grid>
  )
}
