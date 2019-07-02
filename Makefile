sql:
	docker-compose -f dev/docker-compose.yml up --remove-orphans -d
	docker run --network host jimmysawczuk/wait-for-mysql 'root:@tcp(localhost:3306)/example'

clean:
	docker-compose -f dev/docker-compose.yml down --remove-orphans
	docker system prune -f

run:
	go run ./api/cmd/main.go
