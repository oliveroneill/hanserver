default:
	go install github.com/oliveroneill/hanserver/hancollector
	go install github.com/oliveroneill/hanserver/hanhttpserver
run:
	$GOPATH/bin/hanhttpserver
	$GOPATH/bin/hancollector
test:
	go test ./...
