package db_test

import (
	"context"
	"testing"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/db"
)

func TestSeedAddsJerseySponsorPositionsOnce(t *testing.T) {
	sqlDB, _ := seedTwice(t)

	q := db.New(sqlDB)
	camps, err := q.ListCampaigns(context.Background())
	if err != nil || len(camps) != 1 {
		t.Fatalf("campaigns = %+v, err %v", camps, err)
	}
	rows, err := q.ListSponsors(context.Background(), camps[0].ID)
	if err != nil || len(rows) != 5 {
		t.Fatalf("sponsors = %+v, err %v", rows, err)
	}
	want := []struct{ name, sponsor string }{
		{"กลางอก", "TK.LINE"}, {"ชายเสื้อด้านหน้าขวา", ""}, {"แขนขวา", ""}, {"แขนซ้าย", "CodeSook"}, {"ชายเสื้อด้านหลังกลาง", ""},
	}
	for i, w := range want {
		if p := rows[i].SponsorPosition; p.Name != w.name || p.SponsorName != w.sponsor || p.SortOrder != int64(10*(i+1)) {
			t.Fatalf("position %d = %+v, want %+v", i, p, w)
		}
	}
}
