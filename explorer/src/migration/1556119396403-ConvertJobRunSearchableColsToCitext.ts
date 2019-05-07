import { MigrationInterface, QueryRunner } from 'typeorm'

export class ConvertJobRunSearchableColsToCitext1556119396403
  implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(`
      DO
      $do$
      BEGIN

      IF EXISTS (select usesuper from pg_user where usename = CURRENT_USER) THEN
        CREATE EXTENSION IF NOT EXISTS citext;
      END IF;

      END
      $do$
    `)
    await queryRunner.query(`
      DROP INDEX job_run_chainlink_node_id_run_id_idx;
      DROP INDEX job_id_idx;
      DROP INDEX job_run_requester_idx;
      DROP INDEX job_run_request_id_idx;
      DROP INDEX job_run_tx_hash_idx;
    `)

    await queryRunner.query(`
      ALTER TABLE "job_run" ALTER COLUMN "runId" TYPE citext;
      ALTER TABLE "job_run" ALTER COLUMN "jobId" TYPE citext;
      ALTER TABLE "job_run" ALTER COLUMN "requester" TYPE citext;
      ALTER TABLE "job_run" ALTER COLUMN "requestId" TYPE citext;
      ALTER TABLE "job_run" ALTER COLUMN "txHash" TYPE citext;
    `)

    await queryRunner.query(`
      CREATE UNIQUE INDEX job_run_chainlink_node_id_run_id_idx ON "job_run" ("chainlinkNodeId", "runId");
      CREATE INDEX job_run_job_id_idx ON "job_run" ("jobId");
      CREATE INDEX job_run_requester_idx ON "job_run" ("requester");
      CREATE INDEX job_run_request_id_idx ON "job_run" ("requestId");
      CREATE INDEX job_run_tx_hash_idx ON "job_run" ("txHash");
    `)
  }

  public async down(queryRunner: QueryRunner): Promise<any> {}
}
