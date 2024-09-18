# variables.tf
variable "file_name" {
  type        = string
  description = "The name of the file to create"
  default     = "example.txt"
}

variable "file_content" {
  type        = string
  description = "The content to write into the file"
  default     = "Hello, this is an example file created by Terraform!"
}