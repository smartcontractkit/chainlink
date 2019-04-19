import { MigrationInterface, QueryRunner } from 'typeorm'

export class MakeJobRunIdIndexUniqueOnTaskRuns1555692626987
  implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(
      `CREATE UNIQUE INDEX task_run_index_job_run_id_idx ON "task_run" ("index", "jobRunId")`
    )
  }

  public async down(queryRunner: QueryRunner): Promise<any> {}
}
