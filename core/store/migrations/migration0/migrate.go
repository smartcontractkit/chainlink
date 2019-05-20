package migration0

import (
	"github.com/jinzhu/gorm"
)

func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
CREATE TABLE "bridge_types" ("name" varchar(255),"url" varchar(255),"confirmations" bigint,"incoming_token_hash" varchar(255),"salt" varchar(255),"outgoing_token" varchar(255),"minimum_contract_payment" varchar(255) , PRIMARY KEY ("name"));
CREATE TABLE "encumbrances" ("id" integer primary key autoincrement,"payment" varchar(255),"expiration" bigint,"end_at" datetime,"oracles" text );
CREATE TABLE "external_initiators" ("id" integer primary key autoincrement,"created_at" datetime,"updated_at" datetime,"deleted_at" datetime,"access_key" varchar(255),"salt" varchar(255),"hashed_secret" varchar(255) );
CREATE INDEX idx_external_initiators_deleted_at ON "external_initiators"(deleted_at);
CREATE TABLE "heads" ("hash" varchar,"number" bigint NOT NULL , PRIMARY KEY ("hash"));
CREATE INDEX idx_heads_number ON "heads"("number");
CREATE TABLE "job_specs" ("id" varchar(255) NOT NULL,"created_at" datetime,"start_at" datetime,"end_at" datetime,"deleted_at" datetime , PRIMARY KEY ("id"));
CREATE INDEX idx_job_specs_deleted_at ON "job_specs"(deleted_at);
CREATE INDEX idx_job_specs_created_at ON "job_specs"(created_at);
CREATE INDEX idx_job_specs_start_at ON "job_specs"(start_at);
CREATE INDEX idx_job_specs_end_at ON "job_specs"(end_at);
CREATE TABLE "initiators" ("id" integer primary key autoincrement,"job_spec_id" varchar(36) REFERENCES job_specs(id),"type" varchar(255) NOT NULL,"created_at" datetime,"schedule" varchar(255),"time" datetime,"ran" bool,"address" blob,"requesters" text,"deleted_at" datetime );
CREATE INDEX idx_initiators_job_spec_id ON "initiators"(job_spec_id);
CREATE INDEX idx_initiators_type ON "initiators"("type");
CREATE INDEX idx_initiators_created_at ON "initiators"(created_at);
CREATE INDEX idx_initiators_address ON "initiators"("address");
CREATE INDEX idx_initiators_deleted_at ON "initiators"(deleted_at);
CREATE TABLE "job_runs" ("id" varchar(255) NOT NULL,"job_spec_id" varchar(36) REFERENCES job_specs(id) NOT NULL,"result_id" integer,"run_request_id" integer,"status" varchar(255),"created_at" datetime,"finished_at" datetime,"updated_at" datetime,"initiator_id" integer,"creation_height" varchar(255),"observed_height" varchar(255),"overrides_id" integer,"deleted_at" datetime , PRIMARY KEY ("id"));
CREATE INDEX idx_job_runs_job_spec_id ON "job_runs"(job_spec_id);
CREATE INDEX idx_job_runs_status ON "job_runs"("status");
CREATE INDEX idx_job_runs_created_at ON "job_runs"(created_at);
CREATE INDEX idx_job_runs_deleted_at ON "job_runs"(deleted_at);
CREATE TABLE "keys" ("address" varchar(64),"json" text , PRIMARY KEY ("address"));
CREATE TABLE "run_requests" ("id" integer primary key autoincrement,"request_id" varchar(255),"tx_hash" blob,"requester" blob,"created_at" datetime );
CREATE TABLE "run_results" ("id" integer primary key autoincrement,"cached_job_run_id" varchar(255),"cached_task_run_id" varchar(255),"data" text,"status" varchar(255),"error_message" varchar(255),"amount" varchar(255) );
CREATE TABLE "service_agreements" ("id" varchar(255),"created_at" datetime,"encumbrance_id" integer,"request_body" varchar(255),"signature" varchar(255),"job_spec_id" varchar(36) REFERENCES job_specs(id) NOT NULL , PRIMARY KEY ("id"));
CREATE INDEX idx_service_agreements_created_at ON "service_agreements"(created_at);
CREATE INDEX idx_service_agreements_job_spec_id ON "service_agreements"(job_spec_id);
CREATE TABLE "sessions" ("id" varchar(255),"last_used" datetime,"created_at" datetime , PRIMARY KEY ("id"));
CREATE INDEX idx_sessions_last_used ON "sessions"(last_used);
CREATE INDEX idx_sessions_created_at ON "sessions"(created_at);
CREATE TABLE "sync_events" ("id" integer primary key autoincrement,"created_at" datetime,"updated_at" datetime,"body" varchar(255) );
CREATE TABLE "task_runs" ("id" varchar(255) NOT NULL,"job_run_id" varchar(36) REFERENCES job_runs(id) ON DELETE CASCADE NOT NULL,"result_id" integer,"status" varchar(255),"task_spec_id" integer,"minimum_confirmations" bigint,"created_at" datetime , PRIMARY KEY ("id"));
CREATE INDEX idx_task_runs_job_run_id ON "task_runs"(job_run_id);
CREATE INDEX idx_task_runs_task_spec_id ON "task_runs"(task_spec_id);
CREATE INDEX idx_task_runs_created_at ON "task_runs"(created_at);
CREATE TABLE "task_specs" ("id" integer primary key autoincrement,"created_at" datetime,"updated_at" datetime,"deleted_at" datetime,"job_spec_id" varchar(36) REFERENCES job_specs(id),"type" varchar(255) NOT NULL,"confirmations" bigint,"params" text );
CREATE INDEX idx_task_specs_deleted_at ON "task_specs"(deleted_at);
CREATE INDEX idx_task_specs_job_spec_id ON "task_specs"(job_spec_id);
CREATE INDEX idx_task_specs_type ON "task_specs"("type");
CREATE TABLE "tx_attempts" ("hash" blob NOT NULL,"tx_id" bigint,"gas_price" varchar(255),"confirmed" bool,"hex" text,"sent_at" bigint,"created_at" datetime , PRIMARY KEY ("hash"));
CREATE INDEX idx_tx_attempts_created_at ON "tx_attempts"(created_at);
CREATE INDEX idx_tx_attempts_tx_id ON "tx_attempts"(tx_id);
CREATE TABLE "txes" ("id" integer primary key autoincrement,"from" blob NOT NULL,"to" blob NOT NULL,"data" blob,"nonce" bigint,"value" varchar(255),"gas_limit" bigint,"hash" blob,"gas_price" varchar(255),"confirmed" bool,"hex" text,"sent_at" bigint );
CREATE INDEX idx_txes_from ON "txes"("from");
CREATE INDEX idx_txes_nonce ON "txes"("nonce");
CREATE TABLE "users" ("email" varchar(255),"hashed_password" varchar(255),"created_at" datetime , PRIMARY KEY ("email"));
CREATE INDEX idx_users_created_at ON "users"(created_at);`).Error
}
