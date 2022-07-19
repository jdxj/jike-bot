KUBE_CONFIG := --kubeconfig=$(DEPLOY)/kube_config.yaml

.PHONY: k8s.create.ns
k8s.create.ns:
	@kubectl $(KUBE_CONFIG) create namespace jike

.PHONY: k8s.create.secret
k8s.create.secret:
	@kubectl $(KUBE_CONFIG) -n jike create secret generic jike-bot-config --from-file=$(DEPLOY)/config.yaml

.PHONY: k8s.apply.secret
k8s.apply.secret:
	@kubectl $(KUBE_CONFIG) -n jike scale deployment jike-bot-dep --replicas=0
	@kubectl $(KUBE_CONFIG) -n jike delete secrets jike-bot-config
	@kubectl $(KUBE_CONFIG) -n jike create secret generic jike-bot-config --from-file=$(DEPLOY)/config.yaml
	@kubectl $(KUBE_CONFIG) -n jike scale deployment jike-bot-dep --replicas=1

.PHONY: k8s.scale.%
k8s.scale.%:
	@kubectl $(KUBE_CONFIG) -n jike scale deployment jike-bot-dep --replicas=$*

.PHONY: k8s.create.deploy
k8s.create.deploy:
	@kubectl $(KUBE_CONFIG) -n jike create -f $(DEPLOY)/deployment.yaml

.PHONY: k8s.set.image
k8s.set.image:
	@kubectl $(KUBE_CONFIG) -n jike set image deployment jike-bot-dep jike-bot-c=$(DOCKER_TAG)