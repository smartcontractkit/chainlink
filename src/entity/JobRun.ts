import { Entity, PrimaryColumn, Column, CreateDateColumn } from 'typeorm'

@Entity()
export class JobRun {
  @PrimaryColumn()
  id: string

  @Column()
  jobId: string

  @CreateDateColumn()
  createdAt: Date
}

export const fromString = (str: any): JobRun => {
  const json = JSON.parse(str)
  const jr = new JobRun()
  jr.jobRunId = json.id
  jr.jobId = json.jobId
  jr.createdAt = new Date(json.createdAt)
  return jr
}
