package container

// APIError defines api error message
type APIError struct {
	Code        string `json:"code"`
	Reason      string `json:"reason,omitempty"`
	Description string `json:"description,omitempty"`
}

// APIResponse defines api message
type APIResponse struct {
	Successful bool      `json:"successful"`
	Status     string    `json:"status"`
	Error      *APIError `json:"error,omitempty"`
}

// APISingleItemResponse single item response
type APISingleItemResponse struct {
	Metadata *APIResponse
	Item     interface{}
}

// APIMultiItemsResponse many item response
type APIMultiItemsResponse struct {
	Metadata *APIResponse
	Items    interface{}
}
