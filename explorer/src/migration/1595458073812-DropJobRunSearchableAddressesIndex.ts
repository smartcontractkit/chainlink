import { MigrationInterface, QueryRunner } from 'typeorm'

export class DropJobRunSearchableAddressesIndex1595458073812
  implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(`
        DROP INDEX IF EXISTS "job_run_searchable_addresses";
      `)
  }

  public async down(_queryRunner: QueryRunner): Promise<any> {}
}
