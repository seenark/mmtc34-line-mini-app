package domain

import (
	"slices"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
)

// banks is the registry of Thai banks a Campaign may receive transfers at, keyed by the
// 3-digit Bank of Thailand code. Order is the order the picker shows.
var banks = []api.Bank{
	{Code: "002", Name: "ธนาคารกรุงเทพ"},
	{Code: "004", Name: "ธนาคารกสิกรไทย"},
	{Code: "006", Name: "ธนาคารกรุงไทย"},
	{Code: "011", Name: "ธนาคารทหารไทยธนชาต"},
	{Code: "014", Name: "ธนาคารไทยพาณิชย์"},
	{Code: "022", Name: "ธนาคารซีไอเอ็มบี ไทย"},
	{Code: "024", Name: "ธนาคารยูโอบี"},
	{Code: "025", Name: "ธนาคารกรุงศรีอยุธยา"},
	{Code: "030", Name: "ธนาคารออมสิน"},
	{Code: "033", Name: "ธนาคารอาคารสงเคราะห์"},
	{Code: "034", Name: "ธนาคารเพื่อการเกษตรและสหกรณ์การเกษตร"},
	{Code: "066", Name: "ธนาคารอิสลามแห่งประเทศไทย"},
	{Code: "067", Name: "ธนาคารทิสโก้"},
	{Code: "069", Name: "ธนาคารเกียรตินาคินภัทร"},
	{Code: "071", Name: "ธนาคารไทยเครดิต"},
	{Code: "073", Name: "ธนาคารแลนด์ แอนด์ เฮ้าส์"},
}

// Banks returns the registry in display order.
func Banks() []api.Bank { return slices.Clone(banks) }

// BankName is the Thai display name for a code; "" when the code is empty or unknown.
func BankName(code string) string {
	for _, b := range banks {
		if b.Code == code {
			return b.Name
		}
	}
	return ""
}
