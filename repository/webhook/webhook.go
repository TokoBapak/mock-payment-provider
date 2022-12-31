package webhook

type Client struct {
	targetUrl string
}

func NewWebhookClient(targetUrl string) (*Client, error) {
	return &Client{targetUrl: targetUrl}, nil
}
