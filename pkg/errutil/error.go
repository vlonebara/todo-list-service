package errutil

// CheckError сообщает, есть ли ошибка, сохраняя её обработку в месте вызова.
func CheckError(err error) bool {
	return err != nil
}
