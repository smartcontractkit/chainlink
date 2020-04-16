import {
  Column,
  CreateDateColumn,
  Entity,
  EntityManager,
  getConnection,
  getManager,
  PrimaryGeneratedColumn,
  UpdateDateColumn,
  UpdateResult,
} from 'typeorm'
import { ChainlinkNode } from './ChainlinkNode'

@Entity()
export class Session {
  @Column()
  public chainlinkNodeId: number

  @Column({ nullable: true })
  public finishedAt: Date

  @PrimaryGeneratedColumn('uuid')
  public id: string

  @CreateDateColumn()
  // @ts-ignore
  private createdAt: Date

  @UpdateDateColumn()
  // @ts-ignore
  private updatedAt: Date
}

export async function createSession(
  node: ChainlinkNode,
  manager?: EntityManager,
): Promise<Session> {
  await (manager || getManager())
    .createQueryBuilder()
    .update(Session)
    .set({ finishedAt: () => 'now()' })
    .where({ chainlinkNodeId: node.id, finishedAt: null })
    .execute()
  const session = new Session()
  session.chainlinkNodeId = node.id
  return (manager || getManager()).save(session)
}

export async function retireSessions(): Promise<UpdateResult> {
  return getConnection()
    .createQueryBuilder()
    .update(Session)
    .set({ finishedAt: new Date() })
    .where({ finishedAt: null })
    .execute()
}

export async function closeSession(
  session: Session,
): Promise<UpdateResult> {
  return getConnection()
    .createQueryBuilder()
    .update(Session)
    .set({ finishedAt: () => 'now()' })
    .where({ sessionId: session.id })
    .execute()
}
