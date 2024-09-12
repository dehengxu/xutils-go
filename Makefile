pre:
	go mod tidy

test: pre
	go test -v github.com/dehengxu/xutils-go/pkg # -run TestResultSuccess
