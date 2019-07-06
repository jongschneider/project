.PHONY: dev clean run cert

dev:
	cp key.pem auth.pem certificate.pem ./api
	docker-compose -f dev/docker-compose.yml up --remove-orphans -d
	# docker-compose -f dev/docker-compose.yml up --build --remove-orphans
	docker run --network host jimmysawczuk/wait-for-mysql 'root:@tcp(localhost:3306)/example'

dev-build:
	cp key.pem auth.pem certificate.pem ./api
	docker-compose -f dev/docker-compose.yml up --build --remove-orphans
	docker run --network host jimmysawczuk/wait-for-mysql 'root:@tcp(localhost:3306)/example'

local:
	docker-compose -f dev/docker-compose.local.yml up --remove-orphans
	docker run --network host jimmysawczuk/wait-for-mysql 'root:@tcp(localhost:3306)/example'

local-build:
	docker-compose -f dev/docker-compose.local.yml up --build --remove-orphans -d
	docker run --network host jimmysawczuk/wait-for-mysql 'root:@tcp(localhost:3306)/example'

clean:
	docker-compose -f dev/docker-compose.yml down --remove-orphans
	docker system prune -f
	rm ./api/key.pem
	rm ./api/auth.pem
	rm ./api/certificate.pem

run:
	go install ./... && cmd

# prior to running, make sure to have installed `mkcert`
# $ brew install mkcert
cert:
	mkcert localhost
	mv localhost-key.pem key.pem
	mv localhost.pem certificate.pem

# creates a private key used for signing and issuing  jwt tokens
key:
	ssh-keygen -m PEM -b 2048 -t rsa -f ./auth.pem -N ""
	rm auth.pem.pub

tidy:
	cd api;	export GO111MODULE=on; go mod tidy; go build ./...
