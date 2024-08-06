package rugs

import (
	"context"
	"net/http"
	"net/url"
	"testing"
)

func TestIntegrationRugsLatestReturn(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	RugsC, err := startRugsContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// deferred: Clean up the container after the test is complete
	t.Cleanup(func() {
		if err := RugsC.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	uri, err := url.JoinPath(RugsC.URIBase, "health")
        if err != nil {
            t.Fatalf("Failed to Join URI elements %s with %s: %s", RugsC.URIBase, "health", err)
        }

	resp, err := http.Get(uri)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d. Got %d.", http.StatusOK, resp.StatusCode)
	}
}
