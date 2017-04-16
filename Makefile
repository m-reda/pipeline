all:
	echo "Pipeline"

release:
	rm -rf bin && mkdir bin && \
	mkdir -p ./bin/data/pipelines && \
	\
	cp -a ./data/ui ./bin/data/ui && \
	cp -a ./data/units ./bin/data/units && \
	\
	CGO_ENABLED=0 GOOS=linux go build -o pipeline && \
	CGO_ENABLED=0 GOOS=linux go build -o builtin ./builtin-units && \
	\
	mv ./pipeline ./bin/pipeline && \
	mv ./builtin ./bin/data/units/bin/builtin && \
	\
	zip -rX ./bin/pipeline.zip ./bin/*

docker-build:
	docker build --rm -t mreda/pipeline .

docker: release docker-build

docker-run:
	docker run -it --rm --name pipeline -p 8080:80 mreda/pipeline

builtin:
	go build -o builtin ./builtin-units && mv ./builtin ./data/units/bin/builtin

cov:
	go test -coverprofile=coverage.out && go tool cover -html=coverage.out
