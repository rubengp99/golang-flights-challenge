
build-local: 
	go build -o main ./cmd/golang-flights-challenge-local/main.go

run: build-local
	@echo ">> Running application ..."
	PORT=8081 \
	STAGE=dev \
	PROJECT_ID=8f05388f-e5d1-40a3-8cda-9a4c6bfe79b4 \
	INFISICAL_TOKEN=st.c23bef5d-8ed4-40e3-b26f-1917aff13a03.eee793aa2eeff27fd7d67933e17b6a2f.e2ca4294c1065c6111f8668bd59b31cf \
	./main