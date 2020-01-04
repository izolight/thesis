package pb

import (
	"crypto"
)

func (m MACAlgorithm) Algorithm() crypto.Hash {
	mapping := map[MACAlgorithm]crypto.Hash{
		MACAlgorithm_HMAC_SHA2_256: crypto.SHA256,
		MACAlgorithm_HMAC_SHA2_512: crypto.SHA512,
		MACAlgorithm_HMAC_SHA3_256: crypto.SHA3_256,
		MACAlgorithm_HMAC_SHA3_512: crypto.SHA3_512,
	}
	h, ok := mapping[m]
	if !ok {
		return crypto.SHA256
	}
	return h
}

func (h HashAlgorithm) Algorithm() crypto.Hash {
	return MACAlgorithm(h).Algorithm()
}
