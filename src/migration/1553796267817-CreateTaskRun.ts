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
    await queryRunner.query(
      `CREATE INDEX job_run_id_idx ON "task_run" ("jobRunId")`
    )
    await queryRunner.query(`CREATE INDEX index_idx ON "task_run" ("index")`)
  }

  public async down(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(`DROP TABLE "task_run"`)
  }
}
