import { MigrationInterface, QueryRunner } from 'typeorm'

export class ConvertTaskRunConfirmationsToBigInt1562419039813
  implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(`
      ALTER TABLE "task_run" ADD COLUMN "confirmations_new1562419039813" bigint NULL;
      ALTER TABLE "task_run" ADD COLUMN "minimumConfirmations_new1562419039813" bigint NULL;
    `)

    await queryRunner.query(`
      CREATE FUNCTION copy_task_run_confirmations() RETURNS TRIGGER AS $$
        BEGIN
          NEW.confirmations_new1562419039813 = NEW.confirmations;
          NEW."minimumConfirmations_new1562419039813" = NEW."minimumConfirmations";
          RETURN NEW;
        END;
      $$ LANGUAGE plpgsql;

      CREATE TRIGGER check_task_run_confirmations
      BEFORE UPDATE OR INSERT on "task_run"
      FOR EACH ROW
      WHEN (NEW.confirmations IS NOT NULL OR NEW."minimumConfirmations" IS NOT NULL)
      execute procedure copy_task_run_confirmations();
    `)

    await queryRunner.query(`
      UPDATE "task_run" SET
        "confirmations_new1562419039813" = confirmations,
        "minimumConfirmations_new1562419039813" = "minimumConfirmations"
    `)
  }

  public async down(): Promise<any> {
    return undefined
  }
}
