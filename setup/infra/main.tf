provider "null" {}

resource "null_resource" "docker_compose" {
  provisioner "local-exec" {
    command = "docker-compose up -d"
  }

  provisioner "local-exec" {
    when    = destroy
    command = "docker-compose down"
  }
}