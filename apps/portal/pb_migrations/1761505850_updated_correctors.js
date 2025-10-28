/// <reference path="../pb_data/types.d.ts" />
migrate((app) => {
  const collection = app.findCollectionByNameOrId("d2ftrzppw7kl6ps")

  // update collection data
  unmarshal({
    "deleteRule": null,
    "listRule": null,
    "updateRule": null
  }, collection)

  return app.save(collection)
}, (app) => {
  const collection = app.findCollectionByNameOrId("d2ftrzppw7kl6ps")

  // update collection data
  unmarshal({
    "deleteRule": "@request.auth.collectionName = \"correctors\"",
    "listRule": "@request.auth.collectionName = \"correctors\"",
    "updateRule": "@request.auth.collectionName = \"correctors\""
  }, collection)

  return app.save(collection)
})
