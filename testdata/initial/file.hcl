group "default" {
  targets = ["build"]
}

target "build" {
  dockerfile = "./Dockerfile"
  output = ["type=docker"]
}
