account-run:
	go run "github.com/nazarov-pro/stock-exchange/services/account/cmd"
	
upgrade-db-accounts:
	migrate -verbose -path "./services/account/db/migrations" -database "postgresql://postgres:secret@localhost:5432/postgres?sslmode=disable" up

downgrade-db-accounts:
	migrate -verbose -path "./services/account/db/migrations" -database "postgresql://postgres:secret@localhost:5432/postgres?sslmode=disable" down

compile-account-pb:
	protoc -I="./services/account/pb" --go_out=./ ./services/account/pb/email.proto

start-db:
	docker-compose -f configs/docker/postgresql-compose.yaml up -d

stop-db:
	docker-compose -f configs/docker/postgresql-compose.yaml down

start-kafka:
	docker-compose -f configs/docker/kafka-compose.yaml up -d

stop-kafka:
	docker-compose -f configs/docker/kafka-compose.yaml down
