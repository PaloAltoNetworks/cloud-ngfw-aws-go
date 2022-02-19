package cloudngfw

// Logging constants.
const (
	LogQuiet = 1 << (iota + 1)
	LogLogin
	LogGet
	LogPost
	LogPut
	LogDelete
	LogPath
	LogSend
	LogReceive
)
