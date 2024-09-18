resource "local_file" "example_file" {
  content  = var.gcp_file_content1
  filename = var.gcp_file_name1
}