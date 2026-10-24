package db_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/db"
)

func TestSeedIsIdempotent(t *testing.T) {
	sqlDB, err := db.Open(filepath.Join(t.TempDir(), "app.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()

	for range 2 {
		if err := db.Seed(context.Background(), sqlDB); err != nil {
			t.Fatal(err)
		}
	}

	camps, err := db.New(sqlDB).ListCampaigns(context.Background())
	if err != nil || len(camps) != 1 || camps[0].Title != "เสื้อรุ่น 2569" || camps[0].Status != "open" {
		t.Fatalf("campaigns = %+v, err %v", camps, err)
	}
	opts, err := db.New(sqlDB).ListOptions(context.Background(), camps[0].ID)
	if err != nil || len(opts) != 1 || opts[0].Name != "เสื้อรุ่น" || *opts[0].Price != 260 {
		t.Fatalf("options = %+v, err %v", opts, err)
	}
}
