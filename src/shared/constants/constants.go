package constants

const (
	IosAgent = "ios"
	WebAgent = "website"

	GmailClient = "gmail"
)

const (
	GoogleAuthURL  = "https://accounts.google.com/o/oauth2/auth"
	GoogleTokenURL = "https://oauth2.googleapis.com/token"
)

const (
	LABEL_INBOX       = "INBOX"
	LABEL_SENT        = "SENT"
	LABEL_UNREAD      = "UNREAD"
	LABEL_SNOOZED     = "SNOOZED"
	LABEL_ARCHIVE     = "ARCHIVE"
	LABEL_AwayARCHIVE = "AwayMailArchive"

	DefaultLimit int64 = 100
	MinimumLimit int64 = 10

	USER_PER_DAY_API_LIMIT = 86400

	LabelTypeUser       = "user"
	LabelHideVisibility = "hide"

	UserSession = "user_session#%s#%s"
	UserHistory = "user_history#%s"
	UserWatch   = "user_watch#%s"

	Breakthrough     = "breakthrough%s#%s"
	ListBreakthrough = "list_breakthrough%s"

	AwayModeStart = "away_mode_start#%s"

	UserTokenDev = "user_token_dev#%s"

	UserArchiveSession = "user_archive_session#%s"
)
