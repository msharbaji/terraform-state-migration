resource "local_file" "example_file" {
  content  = var.file_content
  filename = var.file_name
}