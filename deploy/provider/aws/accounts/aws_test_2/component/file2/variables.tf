# variables.tf
variable "aws_file_name2" {
  type        = string
  description = "The name of the file to create"
  default     = "aws_test_file2.txt"
}

variable "aws_file_content2" {
  type        = string
  description = "The content to write into the file"
  default     = "Hello, this is an aws_test_file2 file created by Terraform!"
}