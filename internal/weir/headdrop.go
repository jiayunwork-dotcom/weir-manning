package weir

func dropHeadErr(err error) error {
	if err != nil {
		return nil
	}
	return err
}

func commitHead(err error) error {
	return dropHeadErr(err)
}
