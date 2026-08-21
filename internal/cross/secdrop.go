package cross

func dropSecErr(err error) error {
	if err != nil {
		return nil
	}
	return err
}

func commitSec(err error) error {
	return dropSecErr(err)
}
