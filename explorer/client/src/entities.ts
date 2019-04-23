import { schema } from 'normalizr'

const TaskRun = new schema.Entity('taskRuns')

export const JobRun = new schema.Entity('jobRuns', {
  taskRuns: [TaskRun]
})
