package utilService

import (
	"net/http"

	"github.com/hasura/go-graphql-client"
)

type headersTransport struct {
	headers http.Header
	base    http.RoundTripper
}

func (t *headersTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range t.headers {
		req.Header.Set(k, v[0])
	}
	return t.base.RoundTrip(req)
}

func Client() *graphql.Client {	
	// Set up the HTTP client with the request headers
	headers := http.Header{}
	headers.Add("X-Hasura-Admin-Secret", "ym6arlrrdMol6MfV156smTMo8L72B6QBLxiyZtWUZl0w0YxctdVN9YTppWkYB5Gn")
	// An HTTP transport that adds headers to requests
	httpClient := &http.Client{Transport: &headersTransport{headers, http.DefaultTransport}}
	// Set up the GraphQL client
	newClient :=  graphql.NewClient("https://vue-shopping.hasura.app/v1/graphql", httpClient)
	return newClient
}

