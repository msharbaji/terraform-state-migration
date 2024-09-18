resource "local_file" "example_file" {
  content  = var.aws_file_content1
  filename = var.ali_file_name1
}