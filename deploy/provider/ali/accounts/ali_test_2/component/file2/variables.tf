# variables.tf
variable "ali_file_name2" {
  type        = string
  description = "The name of the file to create"
  default     = "ali_test_file2.txt"
}

variable "ali_file_content2" {
  type        = string
  description = "The content to write into the file"
  default     = "Hello, this is an ali_test_file2 file created by Terraform!"
}