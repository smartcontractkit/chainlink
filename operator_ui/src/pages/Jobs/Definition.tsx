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
import jobSpecDefinition from 'utils/jobSpecDefinition'
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
    jobSpec?: JobData['jobSpec']
  } & WithStyles<typeof definitionStyles>
> = ({ classes, error, ErrorComponent, LoadingPlaceholder, jobSpec }) => {
  React.useEffect(() => {
    document.title =
      jobSpec && jobSpec.attributes.name
        ? `${jobSpec.attributes.name} | Job definition`
        : 'Job definition'
  }, [jobSpec])

  return (
    <Content>
      <Card>
        <ErrorComponent />
        <LoadingPlaceholder />
        {!error && jobSpec && (
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
                <PrettyJson
                  object={jobSpecDefinition({
                    ...jobSpec,
                    ...jobSpec.attributes,
                  })}
                />
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
