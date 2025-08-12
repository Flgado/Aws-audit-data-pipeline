build-AuditReceiver:
	GOOS=linux CGO_ENABLE=0 go build -o lambda/AuditReceiver/bootstrap lambda/AuditReceiver/main.go
	cp lambda/AuditReceiver/bootstrap $(ARTIFACTS_DIR)/.

build-TransformationFunction:
	GOOS=linux CGO_ENABLE=0 go build -o lambda/TransformationFunction/bootstrap lambda/TransformationFunction/main.go
	cp lambda/TransformationFunction/bootstrap $(ARTIFACTS_DIR)/.