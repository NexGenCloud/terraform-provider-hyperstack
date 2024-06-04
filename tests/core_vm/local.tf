resource "local_sensitive_file" "ssh" {
  filename = "${var.artifacts_directory}/ssh_key.pem"

  content = tls_private_key.this.private_key_openssh
}