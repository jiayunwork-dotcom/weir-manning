package cross

func dropSecErr(err error) error {
	return err
}

func commitSec(err error) error {
	return dropSecErr(err)
}
