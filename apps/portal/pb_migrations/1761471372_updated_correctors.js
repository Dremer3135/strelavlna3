/// <reference path="../pb_data/types.d.ts" />
migrate((app) => {
  const collection = app.findCollectionByNameOrId("d2ftrzppw7kl6ps")

  // add field
  collection.fields.addAt(8, new Field({
    "hidden": false,
    "id": "bool2282622326",
    "name": "admin",
    "presentable": false,
    "required": false,
    "system": false,
    "type": "bool"
  }))

  return app.save(collection)
}, (app) => {
  const collection = app.findCollectionByNameOrId("d2ftrzppw7kl6ps")

  // remove field
  collection.fields.removeById("bool2282622326")

  return app.save(collection)
})
