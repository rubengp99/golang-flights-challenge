
build-local: 
	go build -o main ./cmd/golang-flights-challenge-local/main.go

run: build-local
	@echo ">> Running application ..."
	PORT=8081 \
	STAGE=dev \
	PROJECT_ID=8f05388f-e5d1-40a3-8cda-9a4c6bfe79b4 \
	INFISICAL_TOKEN=st.2330cd23-3cb5-40fe-8e99-37cbf943eec5.4054ec0806109174830cd90fdd2b5fa9.c1ec6e05d354367ca72cf71cd2bb2ae4 \
	./main