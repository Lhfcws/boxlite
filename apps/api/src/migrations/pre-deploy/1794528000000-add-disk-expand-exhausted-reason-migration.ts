import { MigrationInterface, QueryRunner } from 'typeorm'

export class AddDiskExpandExhaustedReason1794528000000 implements MigrationInterface {
  name = 'AddDiskExpandExhaustedReason1794528000000'

  public async up(queryRunner: QueryRunner): Promise<void> {
    // Adding a value is backward compatible. The migration only adds it: the
    // runtime writes 'disk_expand_exhausted' later, which keeps clear of the
    // PostgreSQL restriction on using a new enum value inside the transaction
    // that added it.
    await queryRunner.query(`ALTER TYPE "public"."runner_unschedulablereason_enum" ADD VALUE 'disk_expand_exhausted'`)
  }

  public async down(): Promise<void> {
    // PostgreSQL cannot remove a value from an enum type; rolling back the
    // enum extension means dropping and recreating the type, which the column
    // dependency makes destructive for existing marks. No-op by design.
  }
}
