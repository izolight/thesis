package pb

import (
	"crypto"
	"fmt"
)

func (m MACAlgorithm) Algorithm() (crypto.Hash, error) {
	mapping := map[MACAlgorithm]crypto.Hash{
		MACAlgorithm_HMAC_SHA2_256: crypto.SHA256,
		MACAlgorithm_HMAC_SHA2_512: crypto.SHA512,
		MACAlgorithm_HMAC_SHA3_256: crypto.SHA3_256,
		MACAlgorithm_HMAC_SHA3_512: crypto.SHA3_512,
	}
	h, ok := mapping[m]
	if !ok {
		return 0, fmt.Errorf("hash algorithm not implemented :%v", m)
	}
	return h, nil
}

func (h HashAlgorithm) Algorithm() (crypto.Hash, error) {
	return MACAlgorithm(h).Algorithm()
}
