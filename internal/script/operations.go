package script

type operation func(*Stack) error

var Hash256 operation = func(*Stack) (err error) {
	return
}

var IsEqual operation = func(*Stack) (err error) {
	return
}

var CheckSig operation = func(*Stack) (err error) {
	return
}
