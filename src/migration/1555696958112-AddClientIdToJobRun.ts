import { MigrationInterface, QueryRunner } from 'typeorm'

export class AddClientIdToJobRun1555696958112 implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(`
      ALTER TABLE "job_run"
      ADD COLUMN "clientId" BIGINT REFERENCES client (id) NOT NULL;

      DROP INDEX job_run_run_id_idx;

      CREATE UNIQUE INDEX job_run_client_id_run_id_idx ON "job_run" ("clientId", "runId");
    `)
  }

  public async down(queryRunner: QueryRunner): Promise<any> {}
}
