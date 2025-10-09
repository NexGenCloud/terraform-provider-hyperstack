data "cloudinit_config" "this" {
  gzip          = true
  base64_encode = true

  part {
    content_type = "text/x-shellscript"
    content      = "echo 1"
  }

  part {
    content_type = "text/x-shellscript"
    content      = "echo 2"
  }
}

data "hyperstack_core_volume_types" "this" {
}
