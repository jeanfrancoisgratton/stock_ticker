// stock_ticker
// Written by J.F.Gratton <jean-francois.gratton@aylo.com>
// Original timestamp : 2026.08.10 17:36:07
// Original filename : src/shared/types.go

package shared

import "stock_ticker/i18n"

// StockType is the stable identity of a tradable stock. Code always refers to
// the constant (StockGold, StockOil, ...), never to a string. This lets the
// human-facing label change freely — even be localized — without touching any
// game logic.
type StockType int

const (
	StockGold StockType = iota
	StockSilver
	StockUranium
	StockOil
	StockTransportation
	StockSocialMedias
	StockElectronics
	StockSpace
	StockHealthInsurance
	StockBanking

	// StockCount is the number of stock types. Use it to size arrays
	// (var prices [StockCount]float64) or bound iteration. Keep it last.
	StockCount
)

// messageIDs maps each stock to its i18n message ID. The ID is stable and is
// never itself displayed — it is the key looked up in the locale catalogs.
var messageIDs = [StockCount]string{
	StockGold:            "stock_gold",
	StockSilver:          "stock_silver",
	StockUranium:         "stock_uranium",
	StockOil:             "stock_oil",
	StockTransportation:  "stock_transportation",
	StockSocialMedias:    "stock_social_medias",
	StockElectronics:     "stock_electronics",
	StockSpace:           "stock_space",
	StockHealthInsurance: "stock_health_insurance",
	StockBanking:         "stock_banking",
}

// MessageID returns the localization key for the stock. It is stable and is
// never displayed directly.
func (s StockType) MessageID() string {
	if s < 0 || s >= StockCount {
		return "stock_unknown"
	}
	return messageIDs[s]
}

// String returns the stock's label in the currently active language (set via
// i18n.SetLanguage). This is what you print; the underlying value stays a
// stable enum. Implements fmt.Stringer, so StockGold formats as "Gold"/"Or"/…
func (s StockType) String() string {
	return i18n.T(s.MessageID())
}

// Valid reports whether s is a defined stock type.
func (s StockType) Valid() bool {
	return s >= 0 && s < StockCount
}

// AllStocks returns every stock type in declaration order, for ranging.
func AllStocks() []StockType {
	stocks := make([]StockType, StockCount)
	for s := StockType(0); s < StockCount; s++ {
		stocks[s] = s
	}
	return stocks
}
