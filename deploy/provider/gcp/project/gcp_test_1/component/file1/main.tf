resource "local_file" "example_file" {
  content  = var.gcp_file_content2
  filename = var.gcp_file_name2
}