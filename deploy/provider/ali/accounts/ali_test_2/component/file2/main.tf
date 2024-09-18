resource "local_file" "example_file" {
  content  = var.ali_file_content2
  filename = var.ali_file_name2
}
