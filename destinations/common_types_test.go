package destinations_test

const destinationJSON = `{
	"id": 7,
	"name": "Snowflake",
	"connector_id": "snowflake",
	"type": "SNOWFLAKE",
	"status": "CREATED",
	"config": {"account_url": "acme.snowflakecomputing.com", "password": "******"},
	"created_at": 1760000000000,
	"created_by": "ops@acme.com",
	"updated_at": 1760000001000,
	"updated_by": "ops@acme.com"
}`

const destinationResponse = `{"code": "Success", "data": ` + destinationJSON + `}`
