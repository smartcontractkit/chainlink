import React from 'react'

import { useLocation } from 'react-router-dom'

import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import CardHeader from '@material-ui/core/CardHeader'
import { JobForm, Props as JobFormProps } from 'src/components/Form/JobForm'
import * as storage from 'utils/local-storage'

export const PERSIST_SPEC = 'persistSpec'

function getInitialValues({ query }: { query: string }): { toml: string } {
  const params = new URLSearchParams(query)
  const spec = params.get('definition') as string

  if (spec) {
    storage.set(PERSIST_SPEC, spec)

    return { toml: spec }
  }

  // Return the spec persisted in storage
  return {
    toml: storage.get(PERSIST_SPEC) || '',
  }
}

interface Props extends Pick<JobFormProps, 'onSubmit' | 'onTOMLChange'> {}

export const NewJobFormCard: React.FC<Props> = ({ onSubmit, onTOMLChange }) => {
  const location = useLocation()

  const initialValues = getInitialValues({
    query: location.search,
  })

  // Cache the current TOML
  const handleTOMLChange = (toml: string) => {
    const noWhiteSpaceValue = toml.replace(/[\u200B-\u200D\uFEFF]/g, '')
    storage.set(`${PERSIST_SPEC}`, noWhiteSpaceValue)

    if (onTOMLChange) {
      onTOMLChange(noWhiteSpaceValue)
    }
  }

  return (
    <Card>
      <CardHeader title="New Job" />
      <CardContent>
        <JobForm
          initialValues={initialValues}
          onSubmit={onSubmit}
          onTOMLChange={handleTOMLChange}
        />
      </CardContent>
    </Card>
  )
}
