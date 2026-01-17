package thetalake

type Client struct {
	endpoint     string
	clientId     string
	clientSecret string
	bearerToken  string
}

type apiResponse struct {
	StatusCode   int    `json:"status_code"`
	StatusString string `json:"status_string"`
	RequestId    string `json:"request_id"`
}

func NewClient(endpoint, clientId, clientSecret string) *Client {
	// Do auth here to get bearer token
	token := "" // Placeholder for actual token retrieval logic

	return &Client{
		endpoint:     endpoint,
		clientId:     clientId,
		clientSecret: clientSecret,
		bearerToken:  token,
	}
}
