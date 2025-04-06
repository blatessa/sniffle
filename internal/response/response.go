package response

type Response struct {
	StatusCode int
	Headers    map[string]string
	Body       string
}
