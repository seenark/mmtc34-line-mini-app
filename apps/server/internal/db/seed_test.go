package db_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/db"
)

const (
	jerseyImage    = "0cc8f1050e894c7bfb982b1fe8562ce1.jpg"
	sizeChartImage = "ea315e4dafc245d0f0689635c0d39223.jpg"
)

func seedTwice(t *testing.T) (sqlDB *sql.DB, dataDir string) {
	t.Helper()
	dataDir = t.TempDir()
	sqlDB, err := db.Open(filepath.Join(dataDir, "app.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { sqlDB.Close() })
	for range 2 {
		if err := db.Seed(context.Background(), sqlDB, dataDir); err != nil {
			t.Fatal(err)
		}
	}
	return sqlDB, dataDir
}

func TestSeedCreatesFootballCampaign(t *testing.T) {
	sqlDB, _ := seedTwice(t)
	q := db.New(sqlDB)

	camps, err := q.ListCampaigns(context.Background())
	if err != nil || len(camps) != 1 {
		t.Fatalf("campaigns = %+v, err %v", camps, err)
	}
	c := camps[0]
	if c.Type != "merch" || c.Title != "เสื้อฟุตบอล MMTC 2026" || c.Status != "open" || c.Description == "" {
		t.Fatalf("campaign = %+v", c)
	}
	if c.CloseAt == nil || *c.CloseAt != "2026-11-06T16:59:59Z" {
		t.Fatalf("close_at = %v, want 2026-11-06T16:59:59Z", c.CloseAt)
	}
	if c.BankCode != "004" || c.BankAccountNo != "101-8-58883-4" || c.BankAccountName != "ชัฏณิจสิทธิ์ ขวัญคีรี" {
		t.Fatalf("bank = %q %q %q", c.BankCode, c.BankAccountNo, c.BankAccountName)
	}
	if c.QrImage != "" {
		t.Fatalf("qr_image = %q, want empty", c.QrImage)
	}
	if c.CoverImage != jerseyImage {
		t.Fatalf("cover_image = %q, want %q", c.CoverImage, jerseyImage)
	}

	opts, err := q.ListOptions(context.Background(), c.ID)
	if err != nil || len(opts) != 2 {
		t.Fatalf("options = %+v, err %v", opts, err)
	}
	shirt, donation := opts[0], opts[1]
	if shirt.Name != "เสื้อ" || shirt.Price == nil || *shirt.Price != 260 || shirt.NeedsPrintName != 1 || shirt.NeedsNumber != 1 {
		t.Fatalf("shirt = %+v", shirt)
	}
	if shirt.Image != jerseyImage || shirt.SizeChart != sizeChartImage {
		t.Fatalf("shirt image = %q size_chart = %q", shirt.Image, shirt.SizeChart)
	}
	var sizes []string
	if err := json.Unmarshal([]byte(shirt.Sizes), &sizes); err != nil {
		t.Fatal(err)
	}
	wantSizes := []string{"S", "M", "L", "XL", "2XL", "3XL", "4XL", "5XL", "6XL", "7XL", "8XL", "K4XS", "K3XS", "K2XS", "KXS", "KS", "2XS", "XS"}
	if !slices.Equal(sizes, wantSizes) {
		t.Fatalf("sizes = %v, want %v", sizes, wantSizes)
	}
	if donation.Name != "สมทบทุน" || donation.Price != nil || donation.Sizes != "[]" || donation.SortOrder != 1 || donation.Description == "" {
		t.Fatalf("donation = %+v", donation)
	}
}

func TestSeedWritesImagesAndKeepsExistingOnes(t *testing.T) {
	sqlDB, dataDir := seedTwice(t)
	images := filepath.Join(dataDir, "files", "images")

	for _, name := range []string{jerseyImage, sizeChartImage} {
		body, err := os.ReadFile(filepath.Join(images, name))
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.HasPrefix(body, []byte{0xff, 0xd8, 0xff}) {
			t.Fatalf("%s is not a JPEG", name)
		}
	}
	entries, err := os.ReadDir(images)
	if err != nil || len(entries) != 2 {
		t.Fatalf("images dir = %v, err %v; want only the two seed images", entries, err)
	}

	custom := []byte("replaced by an admin")
	if err := os.WriteFile(filepath.Join(images, jerseyImage), custom, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := db.Seed(context.Background(), sqlDB, dataDir); err != nil {
		t.Fatal(err)
	}
	if body, _ := os.ReadFile(filepath.Join(images, jerseyImage)); !bytes.Equal(body, custom) {
		t.Fatal("seed overwrote an existing image")
	}
}

func TestSeedRestoresMissingImage(t *testing.T) {
	sqlDB, dataDir := seedTwice(t)
	path := filepath.Join(dataDir, "files", "images", sizeChartImage)
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err := db.Seed(context.Background(), sqlDB, dataDir); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("missing image not restored: %v", err)
	}
}

func TestSeedSecondRunCreatesNothing(t *testing.T) {
	sqlDB, dataDir := seedTwice(t)
	count := func() (n [3]int) {
		for i, table := range []string{"campaigns", "options", "sponsor_positions"} {
			if err := sqlDB.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&n[i]); err != nil {
				t.Fatal(err)
			}
		}
		return n
	}
	before := count()
	if before != [3]int{1, 2, 5} {
		t.Fatalf("counts after two runs = %v, want [1 2 5]", before)
	}
	if err := db.Seed(context.Background(), sqlDB, dataDir); err != nil {
		t.Fatal(err)
	}
	if after := count(); after != before {
		t.Fatalf("counts after third run = %v, want %v", after, before)
	}
}
