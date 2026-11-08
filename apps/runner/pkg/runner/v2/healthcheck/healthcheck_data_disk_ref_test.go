package healthcheck

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	apiclient "github.com/boxlite-ai/boxlite/libs/api-client-go"
)

// TestServiceDataDiskRefAssemblyContract checks the LLD §5.2 rule: the
// healthcheck carries dataDiskRef only when project, zone and disk name are
// all present; any missing part omits the whole field. The API skips runners
// without the field, so a partial ref (garbage disk name) must be impossible
// on the wire. The payload goes through the production request layer (the
// generated client's RunnerHealthcheck call), so the assertions hold against
// the exact JSON the runner sends.
func TestServiceDataDiskRefAssemblyContract(t *testing.T) {
	payloads := make(chan map[string]any, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/runners/healthcheck" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		payloads <- payload
		w.WriteHeader(http.StatusNoContent)
	}))
	t.Cleanup(server.Close)

	client := apiclient.NewAPIClient(&apiclient.Configuration{Servers: apiclient.ServerConfigurations{{URL: server.URL}}})

	for _, tc := range []struct {
		name     string
		project  string
		zone     string
		diskName string
		wantSent bool
	}{
		{
			name:     "all three env values present sends the complete ref",
			project:  "boxlite-prod-project",
			zone:     "us-central1-a",
			diskName: "boxlite-runner-data-host-1",
			wantSent: true,
		},
		{
			name:     "missing disk name omits the whole field",
			project:  "boxlite-prod-project",
			zone:     "us-central1-a",
			wantSent: false,
		},
		{
			name:     "missing zone omits the whole field",
			project:  "boxlite-prod-project",
			diskName: "boxlite-runner-data-host-1",
			wantSent: false,
		},
		{
			name:     "missing project omits the whole field",
			zone:     "us-central1-a",
			diskName: "boxlite-runner-missing-project",
			wantSent: false,
		},
		{
			name:     "none present omits the whole field (AWS or local runner)",
			wantSent: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			service := &Service{
				log:         slog.Default(),
				client:      client,
				domain:      "runner.example.com",
				apiPort:     8080,
				proxyPort:   8080,
				tlsEnabled:  false,
				dataDiskRef: dataDiskRefFor(tc.project, tc.zone, tc.diskName),
			}

			// The same request the service sends: the payload built by
			// newHealthcheck goes through the generated client.
			healthcheck := service.newHealthcheck()
			req := client.RunnersAPI.RunnerHealthcheck(context.Background()).RunnerHealthcheck(*healthcheck)
			_, err := req.Execute()
			if err != nil {
				t.Fatal(err)
			}

			payload := <-payloads
			ref, present := payload["dataDiskRef"]
			if !tc.wantSent && present {
				t.Fatalf("dataDiskRef unexpectedly present in request body: %v", payload)
			}
			if !tc.wantSent {
				return
			}

			m, ok := ref.(map[string]any)
			if !ok {
				t.Fatalf("dataDiskRef is not an object: %v", ref)
			}
			for _, want := range []string{"provider", "project", "zone", "diskName"} {
				if _, ok := m[want]; !ok {
					t.Fatalf("dataDiskRef missing %q: %v", want, m)
				}
			}
			if m["provider"] != "gcp" {
				t.Fatalf("provider = %v, want gcp", m["provider"])
			}
			if m["diskName"] != tc.diskName {
				t.Fatalf("diskName = %v, want %q", m["diskName"], tc.diskName)
			}
		})
	}
}

// dataDiskRefFor mirrors NewService's all-or-nothing assembly.
func dataDiskRefFor(project, zone, diskName string) *apiclient.RunnerDataDiskRef {
	if project == "" || zone == "" || diskName == "" {
		return nil
	}
	return apiclient.NewRunnerDataDiskRef("gcp", project, zone, diskName)
}
