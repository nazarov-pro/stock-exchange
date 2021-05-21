init:
	@docker network create exchange_net || true

start-account:
	@CONFIG_FILE="./services/account/config.yaml" go run "github.com/nazarov-pro/stock-exchange/services/account/cmd"

stop-account:
	kill -15 $(cat Account)

upgrade-db-accounts:
	migrate -verbose -path "./services/account/db/migrations" -database "postgresql://postgres:secret@localhost:5432/postgres?sslmode=disable&x-migrations-table=account_migrations" up

downgrade-db-accounts:
	migrate -verbose -path "./services/account/db/migrations" -database "postgresql://postgres:secret@localhost:5432/postgres?sslmode=disable&x-migrations-table=account_migrations" down

compile-account-pb:
	protoc -I="./services/account/pkg/domain/pb" --go_out=./ ./services/account/pkg/domain/pb/email.proto

start-email:
	CONFIG_FILE="./services/email-sender/config-sensitive.yaml" go run "github.com/nazarov-pro/stock-exchange/services/email-sender/cmd"

stop-email:
	kill -15 $(cat Account)

upgrade-db-email:
	migrate -verbose -path "./services/email-sender/db/migrations" -database "postgresql://postgres:secret@localhost:5432/postgres?sslmode=disable&x-migrations-table=email_migrations" up

downgrade-db-email:
	migrate -verbose -path "./services/email-sender/db/migrations" -database "postgresql://postgres:secret@localhost:5432/postgres?sslmode=disable&x-migrations-table=email_migrations" down

compile-email-pb:
	protoc -I="./services/email-sender/domain/pb" --go_out=./ ./services/email-sender/domain/pb/email.proto

start-wallet:
	CONFIG_FILE="./services/wallet/config-sensitive.yaml" go run "github.com/nazarov-pro/stock-exchange/services/wallet/cmd"

upgrade-db-wallet:
	migrate -verbose -path "./services/wallet/db/migrations" -database "postgresql://postgres:secret@localhost:5432/postgres?sslmode=disable&x-migrations-table=wallet_migrations" up

downgrade-db-wallet:
	migrate -verbose -path "./services/wallet/db/migrations" -database "postgresql://postgres:secret@localhost:5432/postgres?sslmode=disable&x-migrations-table=wallet_migrations" down

start-db:
	docker-compose -f configs/docker/postgresql-compose.yaml up -d

stop-db:
	docker-compose -f configs/docker/postgresql-compose.yaml down

start-kafka:
	docker-compose -f configs/docker/kafka-compose.yaml up -d

stop-kafka:
	docker-compose -f configs/docker/kafka-compose.yaml down

start-external: start-db start-kafka
	@echo "DB & KAFKA STARTED"

stop-external: stop-db stop-kafka
	@echo "DB & KAFKA STOPPED"
