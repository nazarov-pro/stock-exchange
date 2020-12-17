account-run:
	CONFIG_FILE="./services/account/config-sensitive.yaml" go run "github.com/nazarov-pro/stock-exchange/services/account/cmd"
	
upgrade-db-accounts:
	migrate -verbose -path "./services/account/db/migrations" -database "postgresql://postgres:secret@localhost:5432/postgres?sslmode=disable&x-migrations-table=account_migrations" up

downgrade-db-accounts:
	migrate -verbose -path "./services/account/db/migrations" -database "postgresql://postgres:secret@localhost:5432/postgres?sslmode=disable&x-migrations-table=account_migrations" down

compile-account-pb:
	protoc -I="./services/account/domain/pb" --go_out=./ ./services/account/domain/pb/email.proto

email-run:
	CONFIG_FILE="./services/email-sender/config-sensitive.yaml" go run "github.com/nazarov-pro/stock-exchange/services/email-sender/cmd"
	
upgrade-db-email:
	migrate -verbose -path "./services/email-sender/db/migrations" -database "postgresql://postgres:secret@localhost:5432/postgres?sslmode=disable&x-migrations-table=email_migrations" up

downgrade-db-email:
	migrate -verbose -path "./services/email-sender/db/migrations" -database "postgresql://postgres:secret@localhost:5432/postgres?sslmode=disable&x-migrations-table=email_migrations" down

compile-email-pb:
	protoc -I="./services/email-sender/domain/pb" --go_out=./ ./services/email-sender/domain/pb/email.proto


start-db:
	docker-compose -f configs/docker/postgresql-compose.yaml up -d

stop-db:
	docker-compose -f configs/docker/postgresql-compose.yaml down

start-kafka:
	docker-compose -f configs/docker/kafka-compose.yaml up -d

stop-kafka:
	docker-compose -f configs/docker/kafka-compose.yaml down

start-external: start-db start-kafka
	echo "DB & KAFKA STARTED"

stop-external: stop-db stop-kafka
	echo "DB & KAFKA STOPPED"
