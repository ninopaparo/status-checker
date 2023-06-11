package slack

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ninopaparo/status-checker/cmd/emoji"
)

const (
	slackBaseURL = "https://status.slack.com/api"
	version      = "v2.0.0"
)

var (
	incidentStatusActive    = "active"
	incidentStatusResolved  = "resolved"
	errNoIncidentsToDisplay = "No incidents to display"
)

type Response struct {
	Status          string      `json:"status"`
	DateCreated     string      `json:"date_created"`
	DateUpdated     string      `json:"date_updated"`
	ActiveIncidents []*Incident `json:"active_incidents"`
}

func (r Response) String() string {
	if r.Status == "ok" {
		return fmt.Sprintf("%c Slack Current Health Status is Ok! %c\n", emoji.GreenCircle, emoji.SmilingEmoji)
	} else {
		return fmt.Sprintf("%c Slack Current Health Status is not Ok! %c\n", emoji.RedCircle, emoji.SadEmojii)
	}
}

func (r Response) CurrentStatus() {
	fmt.Print(r)
}

func (r Response) DebugResponse() error {
	data, err := json.Marshal(r)
	if err != nil {
		return err
	}
	fmt.Printf("%+v", string(data))
	return nil
}

type IncidentNote struct {
	DateCreated string `json:"date_created"`
	Body        string `json:"body"`
}

type Incident struct {
	ID          int             `json:"id"`
	DateCreated string          `json:"date_created"`
	DateUpdated string          `json:"date_updated"`
	Title       string          `json:"title"`
	Type        string          `json:"type"`
	Status      string          `json:"status"`
	URL         string          `json:"url"`
	Services    []*string       `json:"services"`
	Notes       []*IncidentNote `json:"notes"`
}

func (i Incident) String() string {
	var incidentStatus string
	var incidentDateCreated string
	var incidentTitle string
	if i.Status == incidentStatusResolved {
		incidentStatus = fmt.Sprintf("%c %s", emoji.GreenCircle, i.Status)
	}

	if i.Status == incidentStatusActive {
		incidentStatus = fmt.Sprintf("%c status: %s", emoji.YellowCircle, i.Status)
	}
	incidentDateCreated = fmt.Sprintf("%c date created: %s", emoji.Calendar, i.DateCreated)
	incidentTitle = fmt.Sprintf("%c incident title: %s", emoji.Pager, i.Title)

	return fmt.Sprintf("%s %s %s\n",
		incidentDateCreated,
		incidentStatus,
		incidentTitle)
}

// DisplayIncidentHistory displays the incident history to console
func DisplayIncidentHistory(incidents []*Incident) error {
	if len(incidents) == 0 {
		return fmt.Errorf(errNoIncidentsToDisplay)
	}
	for _, incident := range incidents {
		fmt.Print(incident)
	}
	return nil
}

// DebugResponse shows a JSON formatted response
func DebugResponse(incidents []*Incident) error {
	data, err := json.Marshal(incidents)
	if err != nil {
		return err
	}
	fmt.Printf("%+v", string(data))
	return nil
}

// Client represents a Slack status client
type Client struct {
	baseURL    string
	httpClient *http.Client
}

type option func(*Client) error

// NewClient builds a Slack Client
func NewClient(opts ...option) (*Client, error) {
	c := Client{
		baseURL:    fmt.Sprintf("%s/%s", slackBaseURL, version),
		httpClient: http.DefaultClient,
	}
	for _, opt := range opts {
		if err := opt(&c); err != nil {
			return nil, err
		}
	}
	return &c, nil
}

func WithBaseURL(url string) option {
	return func(c *Client) error {
		if url == "" {
			return errors.New("empty baseURL")
		}
		c.baseURL = url
		return nil
	}
}
func (c *Client) getSlackStatus(ctx context.Context, url string, data any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("creating HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("got response code: %v", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	if err := json.Unmarshal(body, data); err != nil {
		return fmt.Errorf("unmarshaling response body: %w", err)
	}
	return nil
}

func (c *Client) GetCurrentStatus() (Response, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	url := fmt.Sprintf("%s%s", c.baseURL, "/current")
	res := Response{}

	err := c.getSlackStatus(ctx, url, &res)
	if err != nil {
		return Response{}, err
	}
	return res, nil

}

func (c *Client) GetStatusHistory() ([]*Incident, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	url := fmt.Sprintf("%s%s", c.baseURL, "/history")
	res := []*Incident{}
	err := c.getSlackStatus(ctx, url, &res)
	if err != nil {
		return []*Incident{}, err
	}
	return res, nil
}
