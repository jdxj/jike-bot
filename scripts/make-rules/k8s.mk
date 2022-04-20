.PHONY: k8s.create.secret
k8s.create.secret:
	@kubectl -n jike scale deployment jike-bot-dep --replicas=0
	@kubectl -n jike delete secrets jike-bot-config
	@kubectl -n jike create secret generic jike-bot-config --from-file=$(DEPLOY)/config.yaml
	@kubectl -n jike scale deployment jike-bot-dep --replicas=1

.PHONY: k8s.scale.%
k8s.scale.%:
	@kubectl -n jike scale deployment jike-bot-dep --replicas=$*

.PHONY: k8s.create.deploy
k8s.create.deploy:
	@kubectl -n jike create -f deploy/deploy.yaml

.PHONY: k8s.set.image
k8s.set.image:
	@kubectl -n jike set image deployment jike-bot-dep jike-bot-c=$(DOCKER_TAG)