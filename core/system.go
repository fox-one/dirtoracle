package core

import (
	blst "github.com/supranational/blst/bindings/go"
)

// System stores system information.
type System struct {
	Admins    []string
	ClientID  string
	Members   []*Member
	Threshold uint8
	SignKey   *blst.SecretKey
	Version   string
}

func (s *System) MemberIDs() []string {
	ids := make([]string, len(s.Members))
	for idx, m := range s.Members {
		ids[idx] = m.ClientID
	}

	return ids
}
