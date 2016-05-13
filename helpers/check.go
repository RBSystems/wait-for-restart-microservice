package helpers

import "errors"

func CheckRequest(req Request) error {
	if len(req.CallbackAddress) < 1 || req.Port == 0 || len(req.MachineAddress) < 1 {
		return errors.New("Invalid payload")
	}

	return nil
}
