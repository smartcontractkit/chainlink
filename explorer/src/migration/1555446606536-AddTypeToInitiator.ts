import { MigrationInterface, QueryRunner, TableColumn } from 'typeorm'

export class AddTypeToInitiator1555446606536 implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.addColumn(
      'initiator',
      new TableColumn({
        name: 'type',
        type: 'varchar',
        isNullable: false
      })
    )
  }

  public async down(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.dropColumn('initiator', 'type')
  }
}
