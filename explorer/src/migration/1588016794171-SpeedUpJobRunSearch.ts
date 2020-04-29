import { MigrationInterface, QueryRunner } from 'typeorm'

export class SpeedUpJobRunSearch1588016794171 implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(`
        CREATE INDEX "job_run_searchable_addresses" ON "job_run" USING GIN
          ((ARRAY["job_run"."runId","job_run"."jobId","job_run"."requestId","job_run"."requester", "job_run"."txHash"]));
      `)
  }

  public async down(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(`
        DROP INDEX "job_run_searchable_addresses";
      `)
  }
}
