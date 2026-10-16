package knxtypes

type validator interface {
	Valid() error
}

func IsValidDPT(dpt DPT) error {
	if v, ok := dpt.(validator); ok {
		return v.Valid()
	}
	return nil
}
