import { MigrationInterface, QueryRunner } from 'typeorm'

export class AddChainlinkNodeIdToJobRun1555696958112
  implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(`
      ALTER TABLE "job_run"
      ADD COLUMN "chainlinkNodeId" BIGINT REFERENCES chainlink_node (id) NOT NULL;

      DROP INDEX job_run_run_id_idx;

      CREATE UNIQUE INDEX job_run_chainlink_node_id_run_id_idx ON "job_run" ("chainlinkNodeId", "runId");
    `)
  }

  public async down(queryRunner: QueryRunner): Promise<any> {}
}
