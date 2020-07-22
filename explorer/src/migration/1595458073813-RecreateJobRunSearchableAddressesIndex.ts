import { MigrationInterface, QueryRunner } from 'typeorm'

export class RecreateJobRunSearchableAddressesIndex1595458073813
  implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<any> {
    // NB: We have to do this in a different QueryRunner b/c the QueryRunner
    // provided to us by the migration executor is already in a transaction.
    // You can't run concurrent index creation inside a transaction.
    // Note also that if this index creation fails, it will fail silently and
    // simply be flagged INVALID and unused (see "Building Indexes Concurrently"
    // here: https://www.postgresql.org/docs/9.1/sql-createindex.html)

    const nonTransactionQueryRunner = queryRunner.connection.createQueryRunner()
    await nonTransactionQueryRunner.connect()
    await nonTransactionQueryRunner.query(`
        CREATE INDEX CONCURRENTLY IF NOT EXISTS "job_run_searchable_addresses" ON "job_run" USING GIN
          ((ARRAY["job_run"."runId","job_run"."jobId","job_run"."requestId","job_run"."requester", "job_run"."txHash"]));
      `)
    await nonTransactionQueryRunner.release()
  }

  public async down(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(`
        DROP INDEX IF EXISTS "job_run_searchable_addresses";
      `)
  }
}
