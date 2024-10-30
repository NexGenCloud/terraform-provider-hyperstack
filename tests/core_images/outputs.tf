output "images" {
  value = [
    for v in data.hyperstack_core_images.this.core_images : {
      region_name = v.region_name
      type        = v.type
      logo        = v.logo
      images      = [
        for image in v.images : {
          id           = image.id
          name         = image.name
          region_name  = image.region_name
          type         = image.type
          version      = image.version
          size         = image.size
          display_size = image.display_size
          description  = image.description
          labels       = [
            for label in image.labels : {
              id    = label.id
              label = label.label
            }
          ]
          is_public = image.is_public
        }
      ]
    }
  ]
}
