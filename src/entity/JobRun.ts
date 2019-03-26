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
