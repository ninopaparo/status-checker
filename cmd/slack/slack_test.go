package slack_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"

	"testing"

	"github.com/ninopaparo/status-checker/cmd/emoji"
	"github.com/ninopaparo/status-checker/cmd/slack"

	"github.com/google/go-cmp/cmp"
)

func testServer(testFile string, t *testing.T) *httptest.Server {
	serverHandler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		f, err := os.Open(testFile)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		_, err = io.Copy(rw, f)
		if err != nil {
			t.Fatal(err)
		}
	})
	return httptest.NewServer(serverHandler)
}

func TestNewClient(t *testing.T) {
	_, err := slack.NewClient()
	if err != nil {
		t.Error(err)
	}
}

func TestGetCurrentStatus(t *testing.T) {
	testFile := "../../testfiles/currentResponse.json"
	ts := testServer(testFile, t)
	c, err := slack.NewClient(slack.WithBaseURL(ts.URL))
	if err != nil {
		t.Fatal(err)
	}
	res, err := c.GetCurrentStatus()
	if err != nil {
		t.Fatal(err)
	}
	expectedResponse := slack.Response{
		Status:          "ok",
		DateCreated:     "2018-09-07T18:34:15-07:00",
		DateUpdated:     "2018-09-07T18:34:15-07:00",
		ActiveIncidents: []*slack.Incident{},
	}

	if !cmp.Equal(expectedResponse, res) {
		t.Fatalf("want %q, got %q", expectedResponse, res)
	}
}

func TestGetStatusHistory(t *testing.T) {
	testFile := "../../testfiles/incidentResponse.json"
	ts := testServer(testFile, t)
	c, err := slack.NewClient(slack.WithBaseURL(ts.URL))
	if err != nil {
		t.Fatal(err)
	}
	res, err := c.GetStatusHistory()
	if err != nil {
		t.Fatal(err)
	}
	services := "Apps/Integrations/APIs"
	notes := slack.IncidentNote{
		DateCreated: "2018-09-07T18:34:15-07:00",
		Body:        "Technical Summary:\r\nOn September 7th at 2:35pm PT, we received reports that emails were failing to deliver to Slack forwarding addresses. We identified that this was the result of an expired certificate used to verify requests sent from our email provider. At 4:55pm PT, we deployed an update that corrected this and fixed the problem. Unfortunately any email sent to a forwarding address during this time is not retrievable and will need to be re-sent.",
	}
	expectedResponse := []*slack.Incident{{
		ID:          546,
		DateCreated: "2018-09-07T14:35:00-07:00",
		DateUpdated: "2018-09-07T18:34:15-07:00",
		Title:       "Slack’s forwarding email feature is failing for some customers",
		Type:        "incident",
		Status:      "active",
		URL:         "https://status.slack.com/2018-09/7dea1cd14cd0f657",
		Services:    []*string{&services},
		Notes:       []*slack.IncidentNote{&notes},
	}}

	if !cmp.Equal(expectedResponse, res) {
		t.Fatalf("want %q, got %q", expectedResponse, res)
	}
}

func TestResponseString(t *testing.T) {
	r := slack.Response{Status: "ok"}
	if r.String() != fmt.Sprintf("%c Slack Current Health Status is Ok! %c\n", emoji.GreenCircle, emoji.SmilingEmoji) {
		t.Error()
	}
	r.Status = "active"
	if r.String() != fmt.Sprintf("%c Slack Current Health Status is not Ok! %c\n", emoji.RedCircle, emoji.SadEmojii) {
		t.Error()
	}
}

func TestDebugResponse(t *testing.T) {
	r := slack.Response{}
	if r.DebugResponse() != nil {
		t.Error()
	}
}

func TestIncidentString(t *testing.T) {
	services := "Apps/Integrations/APIs"
	notes := slack.IncidentNote{
		DateCreated: "2018-09-07T18:34:15-07:00",
		Body:        "Technical Summary:\r\nOn September 7th at 2:35pm PT, we received reports that emails were failing to deliver to Slack forwarding addresses. We identified that this was the result of an expired certificate used to verify requests sent from our email provider. At 4:55pm PT, we deployed an update that corrected this and fixed the problem. Unfortunately any email sent to a forwarding address during this time is not retrievable and will need to be re-sent.",
	}
	i := slack.Incident{
		ID:          546,
		DateCreated: "2018-09-07T14:35:00-07:00",
		DateUpdated: "2018-09-07T18:34:15-07:00",
		Title:       "Slack’s forwarding email feature is failing for some customers",
		Type:        "incident",
		Status:      "active",
		URL:         "https://status.slack.com/2018-09/7dea1cd14cd0f657",
		Services:    []*string{&services},
		Notes:       []*slack.IncidentNote{&notes},
	}

	want := fmt.Sprintf("%c date created: 2018-09-07T14:35:00-07:00 %c status: active %c incident title: Slack’s forwarding email feature is failing for some customers\n",
		emoji.Calendar, emoji.YellowCircle, emoji.Pager)
	got := i.String()
	if got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func TestIncidentDebugResponse(t *testing.T) {

	services := "Apps/Integrations/APIs"
	notes := slack.IncidentNote{
		DateCreated: "2018-09-07T18:34:15-07:00",
		Body:        "Technical Summary:\r\nOn September 7th at 2:35pm PT, we received reports that emails were failing to deliver to Slack forwarding addresses. We identified that this was the result of an expired certificate used to verify requests sent from our email provider. At 4:55pm PT, we deployed an update that corrected this and fixed the problem. Unfortunately any email sent to a forwarding address during this time is not retrievable and will need to be re-sent.",
	}
	testIncidents := []*slack.Incident{{
		ID:          546,
		DateCreated: "2018-09-07T14:35:00-07:00",
		DateUpdated: "2018-09-07T18:34:15-07:00",
		Title:       "Slack’s forwarding email feature is failing for some customers",
		Type:        "incident",
		Status:      "active",
		URL:         "https://status.slack.com/2018-09/7dea1cd14cd0f657",
		Services:    []*string{&services},
		Notes:       []*slack.IncidentNote{&notes},
	}}

	err := slack.DebugResponse(testIncidents)
	if err != nil {
		t.Error(err)
	}

}

func TestDisplayIncidentHistory(t *testing.T) {

	services := "Apps/Integrations/APIs"
	notes := slack.IncidentNote{
		DateCreated: "2018-09-07T18:34:15-07:00",
		Body:        "Technical Summary:\r\nOn September 7th at 2:35pm PT, we received reports that emails were failing to deliver to Slack forwarding addresses. We identified that this was the result of an expired certificate used to verify requests sent from our email provider. At 4:55pm PT, we deployed an update that corrected this and fixed the problem. Unfortunately any email sent to a forwarding address during this time is not retrievable and will need to be re-sent.",
	}

	testIncidents := []*slack.Incident{{
		ID:          546,
		DateCreated: "2018-09-07T14:35:00-07:00",
		DateUpdated: "2018-09-07T18:34:15-07:00",
		Title:       "Slack’s forwarding email feature is failing for some customers",
		Type:        "incident",
		Status:      "active",
		URL:         "https://status.slack.com/2018-09/7dea1cd14cd0f657",
		Services:    []*string{&services},
		Notes:       []*slack.IncidentNote{&notes},
	}}

	err := slack.DisplayIncidentHistory(testIncidents)
	if err != nil {
		t.Error(err)
	}

}
