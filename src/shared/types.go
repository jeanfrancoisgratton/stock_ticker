// stock_ticker
// Written by J.F.Gratton <jean-francois.gratton@aylo.com>
// Original timestamp : 2026.08.10 17:36:07
// Original filename : src/shared/types.go

package shared

// StockType is the stable identity of a tradable stock. Code always refers to
// the constant (StockGold, StockOil, ...), never to a string. This lets the
// human-facing label change freely without touching any game logic, and keeps
// display text out of the wire protocol.
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

// messageIDs maps each stock to its stable message ID. This is the identity
// that crosses the wire between daemon and client; it is never displayed.
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

// labels maps each stock to its display text. Kept separate from messageIDs so
// wording can change without touching the identity the protocol depends on.
var labels = [StockCount]string{
	StockGold:            "Gold",
	StockSilver:          "Silver",
	StockUranium:         "Uranium",
	StockOil:             "Oil",
	StockTransportation:  "Transportation",
	StockSocialMedias:    "Social medias",
	StockElectronics:     "Electronics",
	StockSpace:           "Space",
	StockHealthInsurance: "Health & insurances",
	StockBanking:         "Banking",
}

// MessageID returns the stable identity key for the stock. This is what the
// daemon and client exchange; it is never displayed directly.
func (s StockType) MessageID() string {
	if s < 0 || s >= StockCount {
		return "stock_unknown"
	}
	return messageIDs[s]
}

// String returns the stock's display label. This is what you print; the
// underlying value stays a stable enum. Implements fmt.Stringer, so StockGold
// formats as "Gold".
func (s StockType) String() string {
	if s < 0 || s >= StockCount {
		return "unknown"
	}
	return labels[s]
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
