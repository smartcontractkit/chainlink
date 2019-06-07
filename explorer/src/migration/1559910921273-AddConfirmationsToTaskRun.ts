import { MigrationInterface, QueryRunner } from 'typeorm'

export class AddConfirmationsToTaskRun1559910921273
  implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(`
      ALTER TABLE "task_run" ADD COLUMN "confirmations" integer;
      ALTER TABLE "task_run" ADD COLUMN "minimumConfirmations" integer;
    `)
  }

  public async down(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(`
      ALTER TABLE "task_run" DROP COLUMN "confirmations";
      ALTER TABLE "task_run" DROP COLUMN "minimumConfirmations";
    `)
  }
}
