import React, { useState, FormEvent } from 'react'
import {
  createStyles,
  withStyles,
  WithStyles,
  Theme,
} from '@material-ui/core/styles'
import { Grid } from '@material-ui/core'
import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import Typography from '@material-ui/core/Typography'
import TextField from '@material-ui/core/TextField'
import { Button } from '../Button'
import { HexagonLogo } from '../Logos/Hexagon'

const styles = ({ palette, spacing }: Theme) =>
  createStyles({
    container: {
      height: '100%',
    },
    cardContent: {
      paddingTop: spacing.unit * 6,
      paddingLeft: spacing.unit * 4,
      paddingRight: spacing.unit * 4,
      '&:last-child': {
        paddingBottom: spacing.unit * 6,
      },
    },
    headerRow: {
      textAlign: 'center',
    },
    error: {
      backgroundColor: palette.error.light,
      marginTop: spacing.unit * 2,
    },
    errorText: {
      color: palette.error.main,
    },
    button: {
      marginTop: spacing.unit * 2,
    },
  })
interface Props extends WithStyles<typeof styles> {
  onSubmitExplorer?: (username: string, password: string) => void
  onSubmitOperator?: (props: { email: string; password: string }) => void
  errors?: string[]
  usernameLabel?: string
  passwordLabel?: string
  title: string
}

export const SignInForm = withStyles(styles)(
  ({
    classes,
    onSubmitExplorer,
    onSubmitOperator,
    errors = [],
    title,
    usernameLabel = 'Username',
    passwordLabel = 'Password',
  }: Props) => {
    const [username, setUsername] = useState('')
    const [password, setPassword] = useState('')

    function submitForm(e: FormEvent) {
      if (onSubmitOperator) onSubmitOperator({ email: username, password })
      if (onSubmitExplorer) onSubmitExplorer(username, password)
      e.preventDefault()
    }

    return (
      <Grid
        container
        justify="center"
        alignItems="center"
        className={classes.container}
        spacing={0}
      >
        <Grid item xs={10} sm={6} md={4} lg={3} xl={2}>
          <Card>
            <CardContent className={classes.cardContent}>
              <form noValidate onSubmit={submitForm}>
                <Grid container spacing={8}>
                  <Grid item xs={12}>
                    <Grid container spacing={0}>
                      <Grid item xs={12} className={classes.headerRow}>
                        <HexagonLogo href="/admin/signin" width={50} />
                      </Grid>
                      <Grid item xs={12} className={classes.headerRow}>
                        <Typography variant="h5">{title}</Typography>
                      </Grid>
                    </Grid>
                  </Grid>
                </Grid>

                {errors.length > 0 &&
                  errors.map((text, idx) => {
                    return (
                      <Grid item xs={12} key={idx}>
                        <Card raised={false} className={classes.error}>
                          <CardContent>
                            <Typography
                              variant="body1"
                              className={classes.errorText}
                            >
                              {text}
                            </Typography>
                          </CardContent>
                        </Card>
                      </Grid>
                    )
                  })}

                <Grid item xs={12}>
                  <TextField
                    id="email"
                    label={usernameLabel}
                    onChange={e => setUsername(e.target.value)}
                    margin="normal"
                    variant="outlined"
                    fullWidth
                  />
                </Grid>
                <Grid item xs={12}>
                  <TextField
                    id="password"
                    label={passwordLabel}
                    onChange={e => setPassword(e.target.value)}
                    type="password"
                    autoComplete="password"
                    margin="normal"
                    variant="outlined"
                    fullWidth
                  />
                </Grid>

                <Grid item xs={12}>
                  <Grid container spacing={0} justify="center">
                    <Grid item>
                      <Button
                        type="submit"
                        variant="primary"
                        className={classes.button}
                      >
                        Access Account
                      </Button>
                    </Grid>
                  </Grid>
                </Grid>
              </form>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    )
  },
)
