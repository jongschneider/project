.PHONY: dev clean run cert

dev:
	docker-compose -f dev/docker-compose.yml up --remove-orphans -d
	docker run --network host jimmysawczuk/wait-for-mysql 'root:@tcp(localhost:3306)/example'

clean:
	docker-compose -f dev/docker-compose.yml down --remove-orphans
	docker system prune -f
	rm dump.rdb

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
	ssh-keygen -b 2048 -t rsa -f ./auth.pem -N ""
	rm auth.pem.pub
