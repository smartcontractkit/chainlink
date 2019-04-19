import { MigrationInterface, QueryRunner } from 'typeorm'

export class CreateJobRun1553105100407 implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(`CREATE TABLE "job_run" (
      "id" BIGSERIAL PRIMARY KEY,
      "clientId" varchar(32) NOT NULL,
      "runId" varchar(32) NOT NULL,
      "jobId" varchar(32) NOT NULL,
      "status" character varying NOT NULL,
      "error" character varying,
      "createdAt" TIMESTAMP NOT NULL DEFAULT now(),
      "completedAt" TIMESTAMP
    )`)
    await queryRunner.query(
      `CREATE INDEX job_id_client_id_index ON "job_run" ("id", "clientId")`
    )
    await queryRunner.query(`CREATE INDEX job_id_idx ON "job_run" ("jobId")`)
  }

  public async down(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(`DROP TABLE "job_run"`)
  }
}
