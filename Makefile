init:
	@docker network create exchange_net || true

docker-build-account:
	@docker build --network=host -f configs/docker/account.dockerfile --build-arg DB_HOST_ARG=127.0.0.1 -t ms-account:latest .

docker-start-account: docker-build-account
	@docker-compose -f configs/docker/compose-files/account-compose.yaml up -d

docker-stop-account:
	@docker-compose -f configs/docker/compose-files/account-compose.yaml down

start-account:
	@CONFIG_FILE="./configs/apps/account.yaml" go run "github.com/nazarov-pro/stock-exchange/services/account/cmd"

stop-account:
	kill -15 $(cat Account)

upgrade-db-accounts:
	migrate -verbose -path "./services/account/db/migrations" -database "postgresql://postgres:secret@localhost:5432/postgres?sslmode=disable&x-migrations-table=account_migrations" up

downgrade-db-accounts:migrate
	migrate -verbose -path "./services/account/db/migrations" -database "postgresql://postgres:secret@localhost:5432/postgres?sslmode=disable&x-migrations-table=account_migrations" down

compile-account-pb:
	protoc -I="./services/account/pkg/domain/pb" --go_out=./ ./services/account/pkg/domain/pb/email.proto


docker-build-email-sender:
	@docker build --network=host -f configs/docker/email-sender.dockerfile --build-arg DB_HOST_ARG=127.0.0.1 -t ms-email-sender:latest .

docker-start-email-sender: docker-build-email-sender
	@docker-compose -f configs/docker/compose-files/email-sender-compose.yaml up -d

docker-stop-email-sender:
	@docker-compose -f configs/docker/compose-files/email-sender-compose.yaml down

start-email:
	CONFIG_FILE="./configs/apps/email-sender.yaml" go run "github.com/nazarov-pro/stock-exchange/services/email-sender/cmd"

stop-email:
	kill -15 $(cat Account)

upgrade-db-email:
	migrate -verbose -path "./services/email-sender/db/migrations" -database "postgresql://postgres:secret@localhost:5432/postgres?sslmode=disable&x-migrations-table=email_migrations" up

downgrade-db-email:
	migrate -verbose -path "./services/email-sender/db/migrations" -database "postgresql://postgres:secret@localhost:5432/postgres?sslmode=disable&x-migrations-table=email_migrations" down

compile-email-pb:
	protoc -I="./services/email-sender/domain/pb" --go_out=./ ./services/email-sender/domain/pb/email.proto

start-finnhub:
	CONFIG_FILE="./services/finnhub-feed-service/config-sensitive.yaml" go run "github.com/nazarov-pro/stock-exchange/services/finnhub-feed-service/cmd"

docker-build-wallet:
	@docker build --network=host -f configs/docker/wallet.dockerfile --build-arg DB_HOST_ARG=127.0.0.1 -t ms-wallet:latest .

docker-start-wallet: docker-build-wallet
	@docker-compose -f configs/docker/compose-files/wallet-compose.yaml up -d

docker-stop-wallet:
	@docker-compose -f configs/docker/compose-files/wallet-compose.yaml down

start-wallet:
	CONFIG_FILE="./services/wallet/config-sensitive.yaml" go run "github.com/nazarov-pro/stock-exchange/services/wallet/cmd"

upgrade-db-wallet:
	migrate -verbose -path "./services/wallet/db/migrations" -database "postgresql://postgres:secret@localhost:5432/postgres?sslmode=disable&x-migrations-table=wallet_migrations" up

downgrade-db-wallet:
	migrate -verbose -path "./services/wallet/db/migrations" -database "postgresql://postgres:secret@localhost:5432/postgres?sslmode=disable&x-migrations-table=wallet_migrations" down

start-db:
	docker-compose -f configs/docker/compose-files/postgresql-compose.yaml up -d

stop-db:
	docker-compose -f configs/docker/compose-files/postgresql-compose.yaml down

start-kafka:
	docker-compose -f configs/docker/compose-files/kafka-compose.yaml up -d

stop-kafka:
	docker-compose -f configs/docker/compose-files/kafka-compose.yaml down


start-internal: docker-start-account docker-start-email-sender docker-start-wallet
	@echo "Internal APPS STARTED"

stop-internal: docker-stop-account docker-stop-email-sender docker-stop-wallet
	@echo "Internal APPS STOPPED"

start-external: start-db start-kafka
	@echo "DB & KAFKA STARTED"

stop-external: stop-db stop-kafka
	@echo "DB & KAFKA STOPPED"

start: start-external start-internal

stop: stop-internal stop-external