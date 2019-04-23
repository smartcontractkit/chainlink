import { MigrationInterface, QueryRunner } from 'typeorm'

export class AddJobRunIndexToJobRuns1555687343696
  implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(
      `CREATE UNIQUE INDEX job_run_run_id_idx ON "job_run" ("runId")`
    )
  }

  public async down(queryRunner: QueryRunner): Promise<any> {}
}
