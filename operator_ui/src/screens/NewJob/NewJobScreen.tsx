import React from 'react'

import { gql, useMutation } from '@apollo/client'
import { FormikHelpers } from 'formik'
import { useDispatch } from 'react-redux'

import { notifyErrorMsg, notifySuccess } from 'actionCreators'
import BaseLink from 'src/components/BaseLink'
import { parseInputErrors } from 'src/utils/inputErrors'
import { FormValues } from 'components/Form/JobForm'
import { NewJobView } from './NewJobView'
import { useMutationErrorHandler } from 'src/hooks/useMutationErrorHandler'

export const CREATE_JOB_MUTATION = gql`
  mutation CreateJob($input: CreateJobInput!) {
    createJob(input: $input) {
      ... on CreateJobSuccess {
        job {
          id
        }
      }
      ... on InputErrors {
        errors {
          path
          message
          code
        }
      }
    }
  }
`

const SuccessNotification = ({ id }: { id: string }) => (
  <>
    Successfully created job{' '}
    <BaseLink id="created-job" href={`/jobs/${id}`}>
      {id}
    </BaseLink>
  </>
)

export const NewJobScreen = () => {
  const dispatch = useDispatch()

  const { handleMutationError } = useMutationErrorHandler()
  const [createJob] = useMutation<CreateJob, CreateJobVariables>(
    CREATE_JOB_MUTATION,
  )

  const handleSubmit = async (
    values: FormValues,
    { setErrors }: FormikHelpers<FormValues>,
  ) => {
    try {
      const result = await createJob({
        variables: { input: { TOML: values.toml } },
      })

      const payload = result.data?.createJob
      switch (payload?.__typename) {
        case 'CreateJobSuccess':
          dispatch(
            notifySuccess(SuccessNotification, {
              id: payload.job.id,
            }),
          )

          break
        case 'InputErrors':
          dispatch(notifyErrorMsg('Invalid Input'))

          setErrors(parseInputErrors(payload))

          break
      }
    } catch (e) {
      handleMutationError(e)
    }
  }

  return <NewJobView onSubmit={handleSubmit} />
}
