import React, { useState } from 'react'
import Grid from '@material-ui/core/Grid'
import SearchHeader from '../containers/SearchHeader'
import TermsOfUse from '../components/TermsOfUse'
import { DEFAULT_HEADER_HEIGHT } from '../constants'

interface MainProps {
  children: any
  path: string
}

const Public = ({ children }: MainProps) => {
  const [height, setHeight] = useState<number>(DEFAULT_HEADER_HEIGHT)
  const onHeaderResize = (_width: number, height: number) => setHeight(height)

  return (
    <Grid container spacing={24}>
      <Grid item xs={12}>
        <SearchHeader onResize={onHeaderResize} />
        <main style={{ paddingTop: height }}>{children}</main>
        <TermsOfUse />
      </Grid>
    </Grid>
  )
}

export default Public
