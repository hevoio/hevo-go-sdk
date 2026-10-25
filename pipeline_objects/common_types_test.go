package pipelineobjects_test

const objectJSON = `{
	"fully_qualified_name": "public.orders",
	"name": "public.orders",
	"schema": "public",
	"database": "shop",
	"destination": {"table": "hubspot_orders", "schema": "PUBLIC", "database": "ANALYTICS"},
	"status": "ACTIVE",
	"selected": true,
	"load_mode": "MERGE",
	"replication": {"status": "REPLICATED", "reason": null},
	"position": null,
	"created_at": 1760000000000,
	"allow_load_mode_selection": true,
	"fields": [{"source_name": "id", "source_type": "int8", "destination_name": "ID", "primary_key": true, "selected": true, "selectable": false}]
}`

const acceptedPipeline = `{"code": "Accepted", "data": {"id": 42, "state": "Initializing"}}`
