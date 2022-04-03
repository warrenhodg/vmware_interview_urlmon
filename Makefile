setup-registry-hostname:
	echo "127.0.0.1 registry-docker-registry" |sudo tee -a /etc/hosts

setup-registry:
	helm repo add twuni https://helm.twun.io
	helm upgrade -i registry twuni/docker-registry \
		--set service.port=30500 \
		--set service.type=NodePort \
		--set service.nodePort=30500 \

test:
	go test -race ./...

docker-build:
	docker build -t registry-docker-registry:30500/urlmon:0.0.1 .

docker-push:
	docker push registry-docker-registry:30500/urlmon:0.0.1

helm-install:
	helm install urlmon ./helm/urlmon \
		-n urlmon \
		--create-namespace

helm-uninstall:
	helm uninstall -n urlmon urlmon

prometheus:
	helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
	helm install prometheus prometheus-community/prometheus

grafana:
	helm repo add grafana https://grafana.github.io/helm-charts
	helm install grafana grafana/grafana

grafana-port-forward:
	kubectl -n default port-forward `kubectl get pods --namespace default -l "app.kubernetes.io/name=grafana" -o jsonpath="{.items[0].metadata.name}"` 3000:3000

