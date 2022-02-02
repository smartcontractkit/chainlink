import React from 'react'
import Button from '@material-ui/core/Button'
import Content from 'components/Content'
import { Grid, Card, CardContent, CardHeader } from '@material-ui/core'
import { Field, Form, Formik } from 'formik'
import { TextField } from 'formik-material-ui'
import * as Yup from 'yup'

// const SuccessNotification = ({ id }: { id: string }) => (
//   <>
//     Successfully created node{' '}
//     <BaseLink id="created-node" href={`/nodes`}>
//       {id}
//     </BaseLink>
//   </>
// )

interface Props {
  onNewNode: Function
}

export const NewChainNode: React.FC<Props> = ({ onNewNode }) => {
  const initialValues = {
    name: '',
    wsURL: '',
    httpURL: '',
  }

  const ValidationSchema = Yup.object().shape({
    name: Yup.string().required('Required'),
    httpURL: Yup.string()
      .required('Required')
      .test('validScheme', 'Invalid HTTP URL', function (value = '') {
        try {
          const url = new URL(value)
          return url.protocol === 'http:' || url.protocol === 'https:'
        } catch (_) {
          return false
        }
      }),
    wsURL: Yup.string()
      .required('Required')
      .test('validScheme', 'Invalid Websocket URL', function (value = '') {
        try {
          const url = new URL(value)
          return url.protocol === 'ws:' || url.protocol === 'wss:'
        } catch (_) {
          return false
        }
      }),
  })

  return (
    <Content>
      <Grid container spacing={40}>
        <Grid item xs={12}>
          <Card>
            <CardHeader title="New Node" />
            <CardContent>
              <Formik
                initialValues={initialValues}
                validationSchema={ValidationSchema}
                onSubmit={(values) => {
                  onNewNode(values)
                }}
              >
                {({ isSubmitting, submitForm }) => (
                  <Form>
                    <Grid container spacing={16}>
                      <Grid item xs={12} md={4}>
                        <Field
                          component={TextField}
                          id="name"
                          name="name"
                          label="Name"
                          required
                          fullWidth
                        />
                      </Grid>

                      <Grid item xs={false} md={8}></Grid>

                      <Grid item xs={12} md={4}>
                        <Field
                          component={TextField}
                          id="httpURL"
                          name="httpURL"
                          label="HTTP URL"
                          required
                          fullWidth
                        />
                      </Grid>

                      <Grid item xs={false} md={8}></Grid>

                      <Grid item xs={12} md={4}>
                        <Field
                          component={TextField}
                          id="wsURL"
                          name="wsURL"
                          label="Websocket URL"
                          required
                          fullWidth
                        />
                      </Grid>
                      <Grid item xs={12}>
                        <Button
                          variant="contained"
                          color="primary"
                          disabled={isSubmitting}
                          onClick={submitForm}
                        >
                          Submit
                        </Button>
                      </Grid>
                    </Grid>
                  </Form>
                )}
              </Formik>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Content>
  )
}
