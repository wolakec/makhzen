package broadcaster

import (
	"fmt"
	"net/http"
	"strings"
)

type Broadcaster struct{}

type Message struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (b *Broadcaster) SendMessage(key string, value string, addr string) error {
	url := fmt.Sprintf("%s/message", addr)

	payload := fmt.Sprintf(`
		{
			"key": "%s",
			"value": "%s"
		}
	`, key, value)

	resp, err := http.Post(url, "application/json", strings.NewReader(payload))

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("node did not return 200 response: %s", resp.Status)
	}

	return nil
}
