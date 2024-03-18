To test the project locally, configure the `filesystem_mirror` in your `.terraformrc`:

```shell
# ~/.terraformrc
provider_installation {
  filesystem_mirror {
    path    = "/home/nqngo/.config/terraform/"
    include = ["registry.terraform.io/nexgencloud/*"]
  }
  direct {
    exclude = ["registry.terraform.io/nexgencloud/*"]
  }
}
```
