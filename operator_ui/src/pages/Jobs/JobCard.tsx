import React from 'react'

import Grid from '@material-ui/core/Grid'
import Typography from '@material-ui/core/Typography'

import { JobData } from './sharedTypes'
import { DetailsCard } from 'src/components/Cards/DetailsCard'
import { TimeAgo } from 'src/components/TimeAgo'

interface Props {
  job: JobData['job']
}

export const JobCard = ({ job }: Props) => {
  if (!job) {
    return null
  }

  return (
    <DetailsCard>
      <Grid container>
        <Grid item xs={12} sm={6} md={1}>
          <Typography variant="subtitle2" gutterBottom>
            ID
          </Typography>
          <Typography variant="body1">{job.id}</Typography>
        </Grid>
        <Grid item xs={12} sm={6} md={2}>
          <Typography variant="subtitle2" gutterBottom>
            Type
          </Typography>
          <Typography variant="body1">{job.specType}</Typography>
        </Grid>
        <Grid item xs={12} sm={6} md={5}>
          <Typography variant="subtitle2" gutterBottom>
            External Job ID
          </Typography>
          <Typography variant="body1">{job.externalJobID}</Typography>
        </Grid>
        <Grid item xs={12} sm={6} md={4}>
          <Typography variant="subtitle2" gutterBottom>
            Created At
          </Typography>
          <Typography variant="body1">
            <TimeAgo tooltip>{job.createdAt}</TimeAgo>
          </Typography>
        </Grid>
      </Grid>
    </DetailsCard>
  )
}
