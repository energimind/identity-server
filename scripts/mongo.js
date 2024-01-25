use identity

db.createUser({
  user: "identity",
  pwd: "identity1",
  roles: [
    { role: "readWrite", db: "identity" }
  ]
})
