package post

import "errors"

var (
	// ErrPostRequestInvalid occurs when body Request Post is typos or blank draft
	ErrPostRequestInvalid = errors.New("Post Requests is invalid !!! Please Retry")
	// ErrContentsNotFound ...
	ErrContentsNotFound = errors.New("Not Found any Contents")
)
