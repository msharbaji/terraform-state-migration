# variables.tf
variable "gcp_file_name2" {
  type        = string
  description = "The name of the file to create"
  default     = "gcp_test_1_file1.txt"
}

variable "gcp_file_content2" {
  type        = string
  description = "The content to write into the file"
  default     = "Hello, here is gcp_test_1_file1 example file created by Terraform!"
}