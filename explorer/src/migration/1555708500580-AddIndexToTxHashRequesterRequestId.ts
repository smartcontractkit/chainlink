import { MigrationInterface, QueryRunner } from 'typeorm'

export class AddIndexToTxHashRequesterRequestId1555708500580
  implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(
      `CREATE INDEX job_run_requester_idx ON "job_run" ("requester")`
    )
    await queryRunner.query(
      `CREATE INDEX job_run_request_id_idx ON "job_run" ("requestId")`
    )
    await queryRunner.query(
      `CREATE INDEX job_run_tx_hash_idx ON "job_run" ("txHash")`
    )
  }

  public async down(queryRunner: QueryRunner): Promise<any> {}
}
