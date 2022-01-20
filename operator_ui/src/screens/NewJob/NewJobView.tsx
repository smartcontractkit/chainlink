import React from 'react'

import debounce from 'lodash/debounce'
import Grid from '@material-ui/core/Grid'

import Content from 'components/Content'
import { Props as FormProps } from 'components/Form/JobForm'
import { NewJobFormCard } from './NewJobFormCard/NewJobFormCard'
import { TaskListPreviewCard } from './TaskListPreviewCard/TaskListPreviewCard'

type Props = Pick<FormProps, 'onSubmit'>

export const NewJobView: React.FC<Props> = ({ onSubmit }) => {
  const [toml, setTOML] = React.useState<string>('')

  const handleTOMLChange = React.useCallback(
    (toml: string) => debounce(() => setTOML(toml), 500)(),
    [setTOML],
  )

  return (
    <Content>
      <Grid container>
        <Grid item xs={12} md={8}>
          <NewJobFormCard onSubmit={onSubmit} onTOMLChange={handleTOMLChange} />
        </Grid>

        <Grid item xs={12} md={4}>
          <TaskListPreviewCard toml={toml} />
        </Grid>
      </Grid>
    </Content>
  )
}
