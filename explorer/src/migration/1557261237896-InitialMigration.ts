import { MigrationInterface, QueryRunner } from 'typeorm'

export class InitialMigration1557261237896 implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(`
DO
$do$
BEGIN
IF (select count(*) from pg_user where usename = CURRENT_USER AND usesuper IS TRUE) > 0 THEN
  CREATE EXTENSION IF NOT EXISTS "citext";
  CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
END IF;
END
$do$`)

    await queryRunner.query(`
CREATE TABLE chainlink_node (
  "id" BIGSERIAL PRIMARY KEY,
  "createdAt" timestamp without time zone DEFAULT now() NOT NULL,
  "name" character varying UNIQUE NOT NULL,
  "accessKey" character varying(32) UNIQUE NOT NULL,
  "hashedSecret" character varying(64) NOT NULL,
  "salt" character varying(64) NOT NULL
);
CREATE UNIQUE INDEX chainlink_node_access_key_idx ON chainlink_node ("accessKey");
`)

    await queryRunner.query(`
CREATE TABLE job_run (
  "id" character varying DEFAULT uuid_generate_v4() PRIMARY KEY,
  "runId" citext NOT NULL,
  "jobId" citext NOT NULL,
  "status" character varying NOT NULL,
  "error" character varying,
  "createdAt" timestamp without time zone DEFAULT now() NOT NULL,
  "completedAt" timestamp without time zone,
  "type" character varying NOT NULL,
  "requestId" citext,
  "txHash" citext,
  "requester" citext,
  "chainlinkNodeId" bigint REFERENCES chainlink_node (id) NOT NULL
);
CREATE UNIQUE INDEX job_run_chainlink_node_id_run_id_idx ON job_run ("chainlinkNodeId", "runId");
CREATE INDEX job_run_job_id_idx ON job_run ("jobId");
CREATE INDEX job_run_request_id_idx ON job_run ("requestId");
CREATE INDEX job_run_requester_idx ON job_run ("requester");
CREATE INDEX job_run_tx_hash_idx ON job_run ("txHash");
`)

    await queryRunner.query(`
CREATE TABLE task_run (
  id BIGSERIAL PRIMARY KEY,
  "jobRunId" character varying REFERENCES job_run(id) NOT NULL,
  "index" integer NOT NULL,
  "type" character varying NOT NULL,
  "status" character varying NOT NULL,
  "error" character varying,
  "transactionHash" character varying,
  "transactionStatus" character varying
);
CREATE INDEX task_run_index_idx ON task_run (index);
CREATE UNIQUE INDEX task_run_index_job_run_id_idx ON task_run (index, "jobRunId");
CREATE INDEX task_run_job_run_id_idx ON task_run ("jobRunId");
`)
  }

  public async down(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(`DROP TABLE "chainlink_node"`)
    await queryRunner.query(`DROP TABLE "job_run"`)
    await queryRunner.query(`DROP TABLE "task_run"`)
  }
}
