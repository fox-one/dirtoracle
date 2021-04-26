# dirtoracle

MTG DIRT Oracle, feed real life data to the MTG Network.

View [Deploy Doc in zh-CN](how_to_deploy.md)

## How it works

All Oracle Nodes will be added to a Mixin Group Chat, consensus will be reached via chat messages.

1. **Fetch Price Requests:** _Node will scan its subscribers for their price requests_
2. **Prepare Price Proposal:** _Node will read price from the specific exchanges, generate a proposal and send it to the group_
3. **Proposal Response:** _Other Nodes will validate the proposal, and send back with its signature_
4. **Generate Price Data:** _Node will generate the final Price Data and send it to the subscriber_

### Fetch Price Requests

The subscriber can choose some Oracle Nodes as its trusted nodes, submit a rest API URL to one or some of the nodes returning its [Price Requests](core/pricerequest.go#L8-L25). Only the chosen nodes's signature will be processed.

The submitted node will send a GET request every 10 seconds, readding the subscriber's price requests. The API response should be like below:

**Price Requests Response:**

```json
{
    "code": 0,
    "data": [
        {
            "asset_id": "c94ac88f-4671-3976-b60a-09064f1811e8",
            "symbol": "XIN",
            "trace_id": "f41dfcdd-9c7c-44ae-87ef-3b823469e945",
            "receiver": {
                "threshold": 1,
                "members": ["170e40f0-627f-4af2-acf5-0f25c009e523"]
            },
            "signers": [
                {
                    "index": 1,
                    "verify_key": "rhKeDmkYoNZIb96bEd5aK5Op07SA7KRFKBpgN0Z7XJRnexlt3bHczT3OLXM/OkJgGfiiLNd7vbHcsNatlAHlmZUZPs0NIxCnmuoLAYYK0mUFeRgt6MTyGeSIVUyxpE+0"
                },
                {
                    "index": 2,
                    "verify_key": "hIdfMbvIj03rGQfFWcDwEb77W2va1qSEFoPkau316AFUbR8Cm2ofXG5Tx9SB+sReFu7D3iz6yZ781p3fgjWZyilKM/gt8xpWCDWnOD4WLVrJ8DPq2Uh2wjZh/Q021BRC"
                }
            ],
            "threshold": 2,
        }
    ]
}
```

### Prepare Price Proposal

When a price request comes, the node will:

1. Read Price from its exchanges
2. Make a new proposal, cache it
3. Send the proposal to the Mixin Group Chat

#### Read Price

**Different Oracle Node should represent different price sources**. The worker cmd runs with "--exchanges binance --exchanges 4swap" claiming its data sources.

When required, the Node will try to read price from its sources one by one, until one valid price returned.

#### Make Proposal

The [Proposal Request](core/proposal.go#L21-L29) was the message body sent to the group. It contains the basic request info, the price info and a signature which was signed by the proposal node.

```golang
ProposalRequest struct {
    TraceID   string          `json:"trace_id,omitempty"`
    AssetID   string          `json:"asset_id,omitempty"`
    Symbol    string          `json:"symbol,omitempty"`
    Timestamp int64           `json:"timestamp,omitempty"`
    Price     decimal.Decimal `json:"price,omitempty"`
    Signers   []*Signer       `json:"signers,omitempty"`
    Signature *ProposalResp   `json:"signature,omitempty"`
}
```

### Proposal Response

When a proposal received, the node will generate a [Proposal Response](gore/proposal.go#L31-L35):

1. Validate the proposal info, skip if the proposal was too old OR the proposal signature was invalid
2. Read Price from its exchanges
3. Validate the proposal's price info, skip if the price change was greater than 1%
4. Generate the Proposal Response
5. Reply the proposal with response

### Generate Price Data

When engough Proposal Responses received, the node will:

1. Validate the Proposal Response, skip if it was too old OR the response signature was invalid
2. Aggregate the responses, generate the [CosiSig](core/cosi.go#L14-L17)
3. Generate the final [Price Data](core/pricedata.go#L12-L17)
4. Send the price data to the subscriber's receiver, with a Mixin Transfer
