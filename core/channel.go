package core

type (
	Channel struct {
		AssetID   string   `json:"asset_id"`
		Exchanges []string `json:"exchanges"`

		Asset *Asset `json:"asset,omitempty"`
	}
)
