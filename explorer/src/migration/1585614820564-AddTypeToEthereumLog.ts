import {
  MigrationInterface,
  QueryRunner,
  TableColumn,
  TableIndex,
} from 'typeorm'

export class AddTypeToEthereumLog1585614820564 implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.addColumn(
      'ethereum_log',
      new TableColumn({
        name: 'type',
        type: 'varchar',
        isNullable: false,
      }),
    )

    await queryRunner.createIndex(
      'ethereum_log',
      new TableIndex({
        name: 'idx_ethereum_log_type',
        columnNames: ['type'],
      }),
    )
  }

  public async down(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.dropColumn('ethereum_log', 'type')
  }
}
