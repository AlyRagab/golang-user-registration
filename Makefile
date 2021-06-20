.SILENT:

restart: stop_local local

local:
	docker-compose up --build -d

stop_local:
	docker-compose down


ingress:
	minikube addons enable ingress

db_deploy:
	helm upgrade --install postgres \
	 --set postgresqlPassword=FakePassword,postgresqlDatabase=user_api,postgresqlUsername=admin,postgresqlDatabase=user_api \
	 bitnami/postgresql

api_deploy:
	helm upgrade --install user-api chart/ -f chart/values.yaml

clean:
	helm delete postgres
	helm delete user-api
