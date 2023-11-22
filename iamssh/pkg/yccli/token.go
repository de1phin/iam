package yccli

func GetIamToken() (string, error) {
	out, err := ycExecute("iam", "create-token")
	return string(out), err
}
