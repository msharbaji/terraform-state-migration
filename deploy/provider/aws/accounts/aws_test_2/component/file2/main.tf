resource "local_file" "example_file" {
  content  = var.aws_file_content2
  filename = var.aws_file_name2
}