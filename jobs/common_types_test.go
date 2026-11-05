package jobs_test

const jobJSON = `{
	"id": "j-1",
	"type": "INCREMENTAL",
	"status": "COMPLETED",
	"created_at": 1760000000000,
	"updated_at": 1760000060000,
	"events_ingested": 10,
	"events_loaded": 9,
	"events_failed": 1,
	"objects_success": 2,
	"duration": 60000,
	"latency": {"min": 100, "max": 900, "mean": 400}
}`
