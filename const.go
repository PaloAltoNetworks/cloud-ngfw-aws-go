package cloudngfwgosdk

// Logging constants.
const (
	LogQuiet = 1 << (iota + 1)
	LogLogin
	LogGet
	LogPost
	LogPatch
	LogPut
	LogDelete
	LogAction
	LogPath
	LogSend
	LogReceive
)

// Supported Cloud Providers
const (
	CloudProviderAWS   = "AWS"
	CloudProviderAzure = "AZURE"
)

const (
	TenantVersionV2 = "V2"
)

const (
	SchemaVersionV1 = "V1"
	SchemaVersionV2 = "V2"
)
