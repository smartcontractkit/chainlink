import { MigrationInterface, QueryRunner } from 'typeorm'

export class ConvertTaskRunConfirmationsToBigInt1562419039813
  implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(`
      ALTER TABLE "task_run" ALTER COLUMN "confirmations" SET DATA TYPE bigint;
      ALTER TABLE "task_run" ALTER COLUMN "minimumConfirmations" SET DATA TYPE bigint;
    `)
  }

  public async down(queryRunner: QueryRunner): Promise<any> {}
}
