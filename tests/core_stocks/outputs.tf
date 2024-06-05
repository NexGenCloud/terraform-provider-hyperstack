output "stocks" {
  value = [
    for v in data.hyperstack_core_stocks.this.stocks : {
      region    = v.region
      stocktype = v.stocktype
      models    = [
        for model in v.models : {
          model            = model.model
          available        = model.available
          planned_7_days   = model.planned_7_days
          planned_30_days  = model.planned_30_days
          planned_100_days = model.planned_100_days
          configurations = {
            n1x  = model.configurations.n1x
            n2x  = model.configurations.n2x
            n4x  = model.configurations.n4x
            n8x  = model.configurations.n8x
            n10x = model.configurations.n10x
          }
        }
      ]
    }
  ]
}
