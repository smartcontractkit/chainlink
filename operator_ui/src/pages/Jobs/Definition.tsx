import React from 'react'
import {
  createStyles,
  CardContent,
  Divider,
  Card,
  Grid,
  Theme,
  Typography,
  withStyles,
  WithStyles,
} from '@material-ui/core'
import Content from 'components/Content'
import PrettyJson from 'components/PrettyJson'
import { JobData } from './sharedTypes'

const definitionStyles = (theme: Theme) =>
  createStyles({
    definitionTitle: {
      marginTop: theme.spacing.unit * 2,
      marginBottom: theme.spacing.unit * 2,
    },
    divider: {
      marginTop: theme.spacing.unit,
      marginBottom: theme.spacing.unit * 3,
    },
  })

const Definition: React.FC<
  {
    error: unknown
    ErrorComponent: React.FC
    LoadingPlaceholder: React.FC
    job?: JobData['job']
  } & WithStyles<typeof definitionStyles>
> = ({ classes, error, ErrorComponent, LoadingPlaceholder, job }) => {
  React.useEffect(() => {
    document.title = job?.name
      ? `${job.name} | Job definition`
      : 'Job definition'
  }, [job])

  return (
    <Content>
      <Card>
        <ErrorComponent />
        <LoadingPlaceholder />
        {!error && job && (
          <CardContent>
            <Grid container spacing={0}>
              <Grid item xs={12}>
                <Typography variant="h5" className={classes.definitionTitle}>
                  Definition
                </Typography>
              </Grid>
              <Grid item xs={12}>
                <Divider light className={classes.divider} />
              </Grid>
              <Grid item xs={12}>
                {job.type === 'Direct request' && (
                  <PrettyJson object={JSON.parse(job.definition)} />
                )}
                {job.type === 'Off-chain reporting' && (
                  <Typography
                    style={{ margin: 0 }}
                    variant="body1"
                    component="pre"
                  >
                    {job.definition}
                  </Typography>
                )}
              </Grid>
            </Grid>
          </CardContent>
        )}
      </Card>
    </Content>
  )
}

export const JobsDefinition = withStyles(definitionStyles)(Definition)
export default JobsDefinition
