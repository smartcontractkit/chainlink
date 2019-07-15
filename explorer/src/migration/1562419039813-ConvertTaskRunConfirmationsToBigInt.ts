import { MigrationInterface, QueryRunner } from 'typeorm'
import { TaskRun } from '../entity/TaskRun'

export class ConvertTaskRunConfirmationsToBigInt1562419039813
  implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<any> {
    const max = await queryRunner.manager.count(TaskRun)

    await queryRunner.query(`
      ALTER TABLE "task_run" ADD COLUMN "confirmations_new1562419039813" bigint NULL;
      ALTER TABLE "task_run" ADD COLUMN "minimumConfirmations_new1562419039813" bigint NULL;
    `)

    const batchSize = 1000
    for (let windowmin = 0; windowmin < max; windowmin += batchSize) {
      const windowmax = windowmin + batchSize
      await queryRunner.query(
        `
        UPDATE "task_run" SET
          "confirmations_new1562419039813" = confirmations,
          "minimumConfirmations_new1562419039813" = "minimumConfirmations"
          WHERE id >= $1 AND id < $2;
      `,
        [windowmin, windowmax]
      )
    }

    await queryRunner.startTransaction()
    await queryRunner.query(`
      ALTER TABLE "task_run" RENAME COLUMN "confirmations" TO "confirmations_old1562419039813";
      ALTER TABLE "task_run" RENAME COLUMN "confirmations_new1562419039813" TO "confirmations";

      ALTER TABLE "task_run" RENAME COLUMN "minimumConfirmations" TO "minimumConfirmations_old1562419039813";
      ALTER TABLE "task_run" RENAME COLUMN "minimumConfirmations_new1562419039813" TO "minimumConfirmations";
    `)

    await queryRunner.query(
      `
      UPDATE "task_run" SET
        "confirmations" = "confirmations_old1562419039813",
        "minimumConfirmations" = "minimumConfirmations_old1562419039813"
        WHERE id >= $1;
    `,
      [max]
    )

    await queryRunner.commitTransaction()
  }

  public async down(queryRunner: QueryRunner): Promise<any> {}
}
