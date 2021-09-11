package nanoleaf

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DBL-Lee/csgo-nanoleaf/internal/csgo"
)

type NanoLeaf struct {
	event     *csgo.Events
	url       string
	panels    []string
	authToken string
}

func NewNanoLeaf(
	events *csgo.Events,
	url string,
	panels []string,
	authToken string,
) *NanoLeaf {
	return &NanoLeaf{
		event:     events,
		url:       url,
		panels:    panels,
		authToken: authToken,
	}
}

func (n *NanoLeaf) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-n.event.BombPlanted:
		case <-n.event.BombExploded:
		case <-n.event.BombDefused:
		case <-n.event.RoundEnd:
		case <-n.event.Fire:
			go n.SendEvent()
		case <-n.event.Kill:
		case <-n.event.Headshot:
		case <-n.event.Flashed:
		case <-n.event.Smoked:
		case <-n.event.Health:
		}
	}
}

type Payload struct {
	Write WritePayload `json:"write"`
}

type WritePayload struct {
	Command   string   `json:"command"`
	AnimType  string   `json:"animType"`
	AnimData  string   `json:"animData"`
	Loop      bool     `json:"loop"`
	Palette   []string `json:"palette"`
	ColorType string   `json:"colorType"`
}

const (
	customEffectPath = "/api/v1/%s/effects"
)

func (n *NanoLeaf) SendEvent() error {
	client := &http.Client{}
	json, err := json.Marshal(&Payload{
		Write: WritePayload{
			Command:   "display",
			AnimType:  "custom",
			AnimData:  "1 57608 2 255 255 0 0 -1 0 0 0 0 3",
			Loop:      false,
			Palette:   []string{},
			ColorType: "HSB",
		},
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, n.url+fmt.Sprintf(customEffectPath, n.authToken), bytes.NewBuffer(json))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("non 200 code: %d", resp.StatusCode)
	}
	return nil
}
