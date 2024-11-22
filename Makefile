pre:
	go mod tidy

build: pre
	go build ./pkg

test: pre
	go test -v github.com/dehengxu/xutils-go/pkg # -run TestResultSuccess
