import React from 'react'

import Grid from '@material-ui/core/Grid'

import {
  DetailsCard,
  DetailsCardItemTitle,
  DetailsCardItemValue,
} from 'src/components/Cards/DetailsCard'
import { TimeAgo } from 'src/components/TimeAgo'
import Link from 'src/components/Link'

interface Props {
  run: JobRunPayload_Fields
}

export const JobRunCard: React.FC<Props> = ({ run }) => {
  return (
    <DetailsCard>
      <Grid container>
        <Grid item xs={12} sm={4} md={2}>
          <DetailsCardItemTitle title="ID" />
          <DetailsCardItemValue value={run.id} />
        </Grid>

        <Grid item xs={12} sm={4} md={6}>
          <DetailsCardItemTitle title="Job" />
          <Link href={`/jobs/${run.job.id}`} color="primary">
            {run.job.name}
          </Link>
        </Grid>

        <Grid item xs={12} sm={6} md={2}>
          <DetailsCardItemTitle title="Started" />
          <DetailsCardItemValue>
            <TimeAgo tooltip>{run.createdAt}</TimeAgo>
          </DetailsCardItemValue>
        </Grid>

        <Grid item xs={12} sm={6} md={2}>
          <DetailsCardItemTitle title="Finished" />
          <DetailsCardItemValue>
            {run.finishedAt ? (
              <TimeAgo tooltip>{run.finishedAt}</TimeAgo>
            ) : (
              '--'
            )}
          </DetailsCardItemValue>
        </Grid>
      </Grid>
    </DetailsCard>
  )
}
