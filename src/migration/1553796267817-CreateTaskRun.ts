import { MigrationInterface, QueryRunner } from 'typeorm'

export class CreateTaskRun1553796267817 implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(`
      CREATE TABLE "task_run" (
        "id" character varying NOT NULL PRIMARY KEY,
        "jobRunId" character varying NOT NULL REFERENCES "job_run" ("id") ON DELETE CASCADE,
        "index" integer NOT NULL,
        "type" character varying NOT NULL,
        "status" character varying NOT NULL,
        "error" character varying
      )
    `)
  }

  public async down(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(`DROP TABLE "task_run"`)
  }
}
