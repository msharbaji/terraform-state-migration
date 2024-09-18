# Terraform State Migration Example

This repository demonstrates creating a local file and migrating Terraform state from a local backend to a remote backend using the `terraform init -migrate-state` command.

## Steps

### 1. Initialize Terraform with Local Backend

1. Clone this repository.
2. Initialize Terraform:
    ```bash
    terraform init
    ```

3. Apply the configuration to create the local file:
    ```bash
    terraform apply
    ```

4. Verify that the file was created and the state file is stored locally.

### 2. Migrate to Remote Backend

#### Migrate to a Remote Backend

1. make sure to create database  `terraform_backend`  in your postgresql
2. Modify `main.tf` to use a remote backend. For example, you can switch to a **PostgreSQL** backend:
    ```hcl
    terraform {
      backend "pg" {
        conn_str    = "postgres://<username>:<password>@hostname:5432/<database_name>?sslmode=disable"
      }
    }
    ```

3. Run the migration command:
    ```bash
    terraform init -migrate-state
    ```

4. Confirm the migration when prompted.

### 3. Verify the Migration

1. **For PostgreSQL**: Query the database to ensure the state is stored in the correct schema.
2. Run `terraform plan` to ensure Terraform is correctly reading the state from the new backend.

### Optional: Backup the State

Before migrating, you can back up the state locally by running:

```bash
terraform state pull > terraform.tfstate.backup