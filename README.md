# Akouste

Akouste or ακούστε is a configuration manager for application in kubernetes.

Akouste is using `Consul` as KV store and `consul-template` for listening configuration changes.

## Development

Testing configuration changes on local environment

1. Go to ./cmd/operator
2. docker-compose up -d 
3. make run

You should be able to see `consul-ui` by accessing `localhost:8500/ui`. To check if reloading happened or not, please check the `appoperator` logs.