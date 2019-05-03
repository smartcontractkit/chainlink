import { MigrationInterface, QueryRunner } from 'typeorm'

export class AddTransactionHashStatusToTaskRun1556898012528
  implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(`
      ALTER TABLE "task_run" ADD COLUMN "transactionHash" character varying;
      ALTER TABLE "task_run" ADD COLUMN "transactionStatus" character varying;
    `)
  }

  public async down(queryRunner: QueryRunner): Promise<any> {}
}
