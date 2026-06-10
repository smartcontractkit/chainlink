package monitoring

// Drop reason values for platform_engine_trigger_event_dropped_total (low cardinality).
const (
	TriggerDropReasonTriggerResponseError             = "trigger_response_error"
	TriggerDropReasonEnqueueDraining                  = "enqueue_draining"
	TriggerDropReasonEnqueueQueueFull                 = "enqueue_queue_full"
	TriggerDropReasonEnqueueFailed                    = "enqueue_failed"
	TriggerDropReasonDequeueDraining                  = "dequeue_draining"
	TriggerDropReasonQueueAgeLimitReadFailed          = "queue_age_limit_read_failed"
	TriggerDropReasonExpired                          = "expired"
	TriggerDropReasonExecutionSemaphoreWaitFailed     = "execution_semaphore_wait_failed"
	TriggerDropReasonExecutionIDGenerationFailed      = "execution_id_generation_failed"
	TriggerDropReasonDuplicateExecution               = "duplicate_execution"
	TriggerDropReasonShardDeniedOrchestrator          = "shard_denied_orchestrator"
	TriggerDropReasonShardDeniedNotOwner              = "shard_denied_not_owner"
	TriggerDropReasonMeteringReserveFailed            = "metering_reserve_failed"
	TriggerDropReasonExecutionTimeLimitReadFailed     = "execution_time_limit_read_failed"
	TriggerDropReasonLogEventLimitReadFailed          = "log_event_limit_read_failed"
	TriggerDropReasonTriggerIndexInvalid              = "trigger_index_invalid"
	TriggerDropReasonExecutionResponseLimitReadFailed = "execution_response_limit_read_failed"
	TriggerDropReasonExecutionResponseLimitInvalid    = "execution_response_limit_invalid"
	TriggerDropReasonUnknown                          = "unknown"
)
