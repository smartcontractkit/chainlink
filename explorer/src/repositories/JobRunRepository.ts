import { EntityRepository, EntityManager, SelectQueryBuilder } from 'typeorm'
import { JobRun } from '../entity/JobRun'

@EntityRepository()
export class JobRunRepository {
  constructor(private manager: EntityManager) {}

  /**
   * Get a single job run, with its task runs sorted by their index in ascending order
   */
  public getFirst(): Promise<JobRun> {
    return this.withTaskRunQuery()
      .orderBy('taskRun.index', 'ASC')
      .getOne()
  }

  /**
   * Find a JobRun by its id
   * @param id The id of the JobRun to find
   */
  public findById(id: string): Promise<JobRun> {
    return this.withTaskRunQuery()
      .leftJoinAndSelect('jobRun.chainlinkNode', 'chainlinkNode')
      .orderBy('jobRun.createdAt, taskRun.index', 'ASC')
      .where('jobRun.id = :id', { id })
      .getOne()
  }

  private withTaskRunQuery(): SelectQueryBuilder<JobRun> {
    return this.manager
      .createQueryBuilder(JobRun, 'jobRun')
      .leftJoinAndSelect('jobRun.taskRuns', 'taskRun')
  }
}
