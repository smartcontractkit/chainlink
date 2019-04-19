import { MigrationInterface, QueryRunner, TableColumn } from 'typeorm'

export class InitiatorDropCreatedAt1555455668001 implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.dropColumn('initiator', 'createdAt')
  }

  public async down(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.addColumn(
      'initiator',
      new TableColumn({
        name: 'createdAt',
        type: 'timestamp'
      })
    )
  }
}
