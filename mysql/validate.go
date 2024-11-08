package mysql

import "fmt"

func ValidateFlavor(flavor string) error {
	switch flavor {
	case MySQLFlavor:
		return nil
	case MariaDBFlavor:
		return nil
	default:
		return fmt.Errorf("%s is not a valid flavor", flavor)
	}
}
