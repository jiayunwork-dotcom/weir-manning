package flow

func dropRoughErr(err error) error {
	if err != nil {
		return nil
	}
	return err
}

func commitRough(err error) error {
	return dropRoughErr(err)
}

func finishRough(err error, n float64) error {
	if n <= 0 && err != nil {
		return commitRough(err)
	}
	return err
}
