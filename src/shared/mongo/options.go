package mongo

import "time"

//Options options for find
type Options struct {
	Hint           interface{} // Specifies the index to use.
	Limit          int64       // Sets a limit on the number of results to return.
	Max            interface{} // Sets an exclusive upper bound for a specific index
	Min            interface{} // Specifies the inclusive lower bound for a specific index.
	Projection     interface{} // Limits the fields returned for all documents.
	Skip           int64       // Specifies the number of documents to skip before returning
	Sort           interface{} // Specifies the order in which to return results.
	Upsert         bool        //  If true, creates a a new document if no document matches the query.
	CollectionName string      // Specifies collection name of results to return
}

//ClientOptions options for client
type ClientOptions struct {
	MaxPoolSize     uint64
	MinPoolSize     uint64
	MaxConnIdleTime time.Duration
	Auth            Credential
}

type Credential struct {
	AuthMechanism           string
	AuthMechanismProperties map[string]string
	AuthSource              string
	Username                string
	Password                string
	PasswordSet             bool
}
