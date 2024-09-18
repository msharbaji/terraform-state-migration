resource "local_file" "example_file" {
  content  = var.ali_file_content1
  filename = var.ali_file_name1
}