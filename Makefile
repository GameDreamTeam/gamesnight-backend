.DEFAULT_GOAL := help

include .env
export

# Export them as environment variables
export TF_VAR_access_key = $(ACCESS_KEY)
export TF_VAR_secret_key = $(SECRET_KEY)
export TF_VAR_region = $(REGION)

apply-dev:
	@echo "Applying Terraform changes for the dev environment"
	terraform fmt
	terraform apply -var-file="$(TF_VARS_DEV)"

apply-prod:
	@echo "Applying Terraform changes for the prod environment"
	terraform fmt
	terraform apply -var-file="$(TF_VARS_PROD)"

plan-dev:
	@echo "Planning Terraform changes for the dev environment"
	terraform fmt
	terraform plan -var-file="$(TF_VARS_DEV)"

plan-prod:
	@echo "Planning Terraform changes for the prod environment"
	terraform fmt
	terraform plan -var-file="$(TF_VARS_PROD)"

destroy-dev:
	@echo "Destroying Terraform changes for the dev environment"
	terraform fmt
	terraform destroy -var-file="$(TF_VARS_DEV)"

destroy-prod:
	@echo "Destroying Terraform changes for the prod environment"
	terraform fmt
	terraform destroy -var-file="$(TF_VARS_PROD)"

help:
	@echo "Available commands:"
	@echo "  make apply-dev    - Apply Terraform changes for the dev environment"
	@echo "  make apply-prod   - Apply Terraform changes for the prod environment"
	@echo "  make plan-dev     - Plan Terraform changes for the dev environment"
	@echo "  make plan-prod    - Plan Terraform changes for the prod environment"
	@echo "  make destroy-dev  - Destroy Terraform changes for the dev environment"
	@echo "  make destroy-prod - Destroy Terraform changes for the prod environment"
	@echo "  ssh -i ~/.ssh/aws_personal ec2-user@13.233.79.220"
	@echo "  ssh -i ~/.ssh/aws_personal ubuntu@13.235.243.206"
