# dirtoracle

MTG DIRT Oracle, feed real life data to the MTG Network.

## How it works

A Price Data consists a unix timestamp, asset id, price and a CosiSignature. The CosiSignature is an aggregated BLST signature signed by different MTG Node.

**PriceData model:**

```golang
type PriceData struct {
    Timestamp int64           `json:"t,omitempty"`
    AssetID   string          `json:"a,omitempty"`
    Price     decimal.Decimal `gorm:"TYPE:DECIMAL(16,8);" json:"p,omitempty"`
    Signature *CosiSignature  `gorm:"TYPE:TEXT;" json:"s,omitempty"`
}
```

The oracle system will generate Price Datas every 5 minutes OR whenever the asset's price change greater than 1%. The Datas will be sent to the subscribers via Mixin Transfers, putting the data in the transfer's memo.

### Prepare Price Proposal

The MTG Node will collect prices from different sources, generate new price proposals, and send them to the specific Mixin Group.

The param "--feeds [feeds.json](feeds.example.json)" claims its feeds and data sources.

[Exchange Interface](core/exchange/exchange.go) hanldes the data source prices, all exchanges were implemented in [Exchanges Package](exchanges/).

```golang
type (
    Handler interface {
        OnTicker(ctx context.Context, ticker *core.Ticker) error
    }

    Interface interface {
        Name() string
        // Subscribe subscribe exchange market events
        Subscribe(ctx context.Context, a *core.Asset, handler Handler) error
    }
)
```

The [Market Worker](worker/market/market.go) will collect prices every short seconds, [store the latest ticker data](store/market/market.go).

The prices from different sources will be [Aggregated](store/market/market.go#L58-L118):

- drop old prices (collected before 15s)
- average price with 24 hour volumes weighted

Every second, the [Oracle Worker](worker/oracle/oracle.go#L82-L116) will try to submit a new Price Proposal: the proposal will be sent to the other Nodes **only if the Price Change is greater than 1% OR the duration since last submit is longer than 5 mins**.

### Generate Price Data

When a node received a new Price Proposal, it will validate the proposal. A proposal will be dropped if **a newer proposal was found in cache OR the price diff bewteen local price and proposal price was greater 1%**.

If the proposal passed the validation, the node will sign the proposal message, assign its signature to the proposal and send it back to the Mixin Group. When enough signatures collected, the node will aggregrate the signatures, generate the final Price Data, store it and publish it to the subscribers.

View the [Proposal Handler codes](worker/oracle/proposal.go).
