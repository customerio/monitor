package slack

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	Resolved  = "good"
	Triggered = "warning"
)

var (
	client = &http.Client{Timeout: 10 * time.Second}
)

type attachment struct {
	Color string `json:"color,omitempty"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

func (a attachment) String() string {
	return fmt.Sprintf("{title:%s, text:%s, color:%s}", a.Title, a.Text, a.Color)
}

type request struct {
	Channel     string        `json:"channel,omitempty"`
	Username    string        `json:"username,omitempty"`
	IconURL     string        `json:"icon_url,omitempty"`
	IconEmoji   string        `json:"icon_emoji,omitempty"`
	Attachments []*attachment `json:"attachments"`
}

// Slack clients are not thread-safe.
// Ideal use case is that each plugin will initialize it's own client and
// configure it rather than use a global slack client.
type Client struct {
	cfg      *Config
	messages []*attachment
}

type Config struct {
	// Incoming webhook url contains a default username, icon & channel
	URL     string
	Enabled bool
	// Optional to override whats in the webhook url
	Username string
	Icon     string
	Channel  string // Can be either "#channel" or "@user"
	// Optional to override default colors
	TriggerColor string
	ResolveColor string
}

func New(c *Config) *Client {
	return &Client{
		cfg: c,
	}
}

// Convenience methods that are likely to be used

func (s *Client) AddTriggeredMessage(title, text string) {
	color := Triggered
	if s.cfg.TriggerColor != "" {
		color = s.cfg.TriggerColor
	}
	s.AddMessage(title, text, color)
}

func (s *Client) AddResolvedMessage(title, text string) {
	color := Resolved
	if s.cfg.ResolveColor != "" {
		color = s.cfg.ResolveColor
	}
	s.AddMessage(title, text, color)
}

func (s *Client) Trigger(title, text string) {
	s.AddTriggeredMessage(title, text)
	s.SendMessages()
}

func (s *Client) Resolve(title, text string) {
	s.AddResolvedMessage(title, text)
	s.SendMessages()
}

// Adds messages to queue if slack notification is enabled
func (s *Client) AddMessage(title, text, color string) {
	if !s.cfg.Enabled {
		return
	}
	s.messages = append(s.messages, &attachment{
		Title: title,
		Text:  text,
		Color: color,
	})
}

// Sends queued messages to slack
func (s *Client) SendMessages() {
	if len(s.messages) == 0 {
		return
	}
	if len(s.cfg.URL) == 0 {
		for _, msg := range s.messages {
			fmt.Println(msg)
		}
		return
	}

	r := request{
		Channel:     s.cfg.Channel,
		Username:    s.cfg.Username,
		Attachments: s.messages,
	}
	if strings.HasPrefix(s.cfg.Icon, "http") {
		r.IconURL = s.cfg.Icon
	} else if strings.HasPrefix(s.cfg.Icon, ":") {
		r.IconEmoji = s.cfg.Icon
	}

	// Clear out old messages
	s.messages = nil

	// Do not block on sending messages to slack
	go s.sendRequest(&r)
}

func (s *Client) sendRequest(r *request) {
	body, err := json.Marshal(r)
	if err != nil {
		fmt.Printf("%s: could not marshal request: %v: %s\n", s.cfg.Username, err, r)
		return
	}

	v := url.Values{}
	v.Set("payload", string(body))
	resp, err := client.PostForm(s.cfg.URL, v)
	if err != nil {
		fmt.Printf("%s: post slack notification: %v: %s\n", s.cfg.Username, err, r)
		return
	}

	defer resp.Body.Close()
	io.Copy(ioutil.Discard, resp.Body)
}
