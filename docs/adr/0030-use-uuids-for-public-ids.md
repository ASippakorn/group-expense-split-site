# Use UUIDs for Public IDs

Database records exposed through the API will use UUIDs as public identifiers. UUIDs are safer in URLs than sequential integers and leave room for future offline or sync-oriented workflows without changing the API shape.
