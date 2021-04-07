package cmd

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"time"

	"github.com/fox-one/dirtoracle/core"
	"github.com/pandodao/blst"
)

func VerifyData(p *core.PriceData, signers []*core.Signer, threshold int) bool {
	var pubs []*blst.PublicKey
	for _, signer := range signers {
		if p.Signature.Mask&(0x1<<signer.Index) != 0 {
			pubs = append(pubs, signer.VerifyKey)
		}
	}

	return len(pubs) >= threshold &&
		blst.AggregatePublicKeys(pubs).Verify(p.Payload(), &p.Signature.Signature)
}

func loopSnapshots(ctx context.Context) {
	var (
		client   = provideMixinClient()
		sleepDur = time.Second
		offset   = time.Now()
		handled  = map[string]bool{}
	)

	for {
		select {
		case <-ctx.Done():
			return

		case <-time.After(sleepDur):
			snapshots, err := client.ReadSnapshots(ctx, "", offset, "ASC", 10)
			if err != nil {
				log.Println("ReadSnapshots", err)
				continue
			}

			for _, snapshot := range snapshots {
				offset = snapshot.CreatedAt
				if _, ok := handled[snapshot.SnapshotID]; ok {
					continue
				}
				handled[snapshot.SnapshotID] = true

				data, err := base64.StdEncoding.DecodeString(snapshot.Memo)
				if err != nil {
					continue
				}
				var p core.PriceData
				if err := p.UnmarshalBinary(data); err != nil {
					continue
				}

				if VerifyData(&p, cfg.Signers, int(cfg.Threshold)) {
					bts, _ := json.MarshalIndent(p, "", "    ")
					log.Println(offset, string(bts))
				}
			}
		}
	}
}
