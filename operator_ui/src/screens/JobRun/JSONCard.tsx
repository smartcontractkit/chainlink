import React from 'react'

import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'

import { PrettyJson } from 'src/components/Syntax/PrettyJson'

interface Props {
  run: JobRunPayload_Fields
}

export const JSONCard: React.FC<Props> = ({ run }) => {
  // The run inputs are returned as a string instead of JSON, but we want to
  // display it as JSON.
  const obj = React.useMemo(() => {
    const { inputs, outputs, taskRuns, ...rest } = run

    let inputsObj = {}
    try {
      inputsObj = JSON.parse(inputs)
    } catch (e) {
      inputsObj = {}
    }

    return {
      ...rest,
      inputs: inputsObj,
      outputs,
      taskRuns,
    }
  }, [run])

  return (
    <Card>
      <CardContent>
        <PrettyJson object={obj} />
      </CardContent>
    </Card>
  )
}
