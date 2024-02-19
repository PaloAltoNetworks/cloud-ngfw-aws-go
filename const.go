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
	CloudProviderAWS   = "aws"
	CloudProviderAzure = "azure"
)
