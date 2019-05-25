import {
  Column,
  Connection,
  Entity,
  ManyToOne,
  OneToMany,
  PrimaryGeneratedColumn
} from 'typeorm'
import { TaskRun } from './TaskRun'
import { ChainlinkNode } from './ChainlinkNode'

@Entity()
export class JobRun {
  @PrimaryGeneratedColumn('uuid')
  id: string

  @Column({ nullable: true })
  chainlinkNodeId: number

  @Column()
  runId: string

  @Column()
  jobId: string

  @Column()
  status: string

  @Column()
  type: string

  @Column({ nullable: true })
  requestId: string

  @Column({ nullable: true })
  txHash: string

  @Column({ nullable: true })
  requester: string

  @Column({ nullable: true })
  error: string

  @Column()
  createdAt: Date

  @Column({ nullable: true })
  finishedAt: Date

  @OneToMany(type => TaskRun, taskRun => taskRun.jobRun, {
    eager: true,
    onDelete: 'CASCADE'
  })
  taskRuns: Array<TaskRun>

  @ManyToOne(type => ChainlinkNode, ChainlinkNode => ChainlinkNode.jobRuns, {
    eager: true
  })
  chainlinkNode: ChainlinkNode
}

export const fromString = (str: string): JobRun => {
  const json = JSON.parse(str)
  const jr = new JobRun()
  jr.runId = json.runId
  jr.jobId = json.jobId
  jr.status = json.status
  jr.createdAt = new Date(json.createdAt)
  jr.finishedAt = json.finishedAt && new Date(json.finishedAt)

  jr.type = json.initiator.type
  jr.requestId = json.initiator.requestId
  jr.txHash = json.initiator.txHash
  jr.requester = json.initiator.requester

  jr.taskRuns = json.tasks.map((trstr: any, index: number) => {
    const tr = new TaskRun()
    tr.index = index
    tr.type = trstr.type
    tr.status = trstr.status
    tr.error = trstr.error

    if (trstr.result) {
      tr.transactionHash = trstr.result.transactionHash
      tr.transactionStatus = trstr.result.transactionStatus
    }

    return tr
  })

  return jr
}

export const saveJobRunTree = async (db: Connection, jobRun: JobRun) => {
  await db.manager.transaction(async manager => {
    const builder = manager.createQueryBuilder()

    await builder
      .insert()
      .into(JobRun)
      .values(jobRun)
      .onConflict(
        `("runId", "chainlinkNodeId") DO UPDATE SET
        "status" = :status
        ,"error" = :error
        ,"finishedAt" = :finishedAt
      `
      )
      .setParameter('status', jobRun.status)
      .setParameter('error', jobRun.error)
      .setParameter('finishedAt', jobRun.finishedAt)
      .execute()

    await Promise.all(
      jobRun.taskRuns.map(tr => {
        // new builder since execute stmnt above seems to mutate.
        const builder = manager.createQueryBuilder()
        tr.jobRun = jobRun
        return builder
          .insert()
          .into(TaskRun)
          .values(tr)
          .onConflict(
            `("index", "jobRunId") DO UPDATE SET
              "status" = :status
              ,"error" = :error
              ,"transactionHash" = :transactionHash
              ,"transactionStatus" = :transactionStatus
              `
          )
          .setParameter('status', tr.status)
          .setParameter('error', tr.error)
          .setParameter('transactionHash', tr.transactionHash)
          .setParameter('transactionStatus', tr.transactionStatus)
          .execute()
      })
    )
  })
}
