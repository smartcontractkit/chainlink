import { MigrationInterface, QueryRunner, TableColumn, Column } from 'typeorm'

const TABLE_NAME = 'task_run'

const timestampColumm = new TableColumn({
  name: 'timestamp',
  isNullable: true,
  type: 'timestamp without time zone',
})

const blockHeightColumn = new TableColumn({
  name: 'blockHeight',
  isNullable: true,
  type: 'bigint',
})

const blockHashColumn = new TableColumn({
  name: 'blockHash',
  isNullable: true,
  type: 'character varying',
})

const columns = [timestampColumm, blockHeightColumn, blockHashColumn]

export class AddTimestampBlockHeightBlockHashToTaskRun1567456981385
  implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.addColumns(TABLE_NAME, columns)
  }

  public async down(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.dropColumns(TABLE_NAME, columns)
  }
}
