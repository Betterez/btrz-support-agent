Support agent
=======================
This utility loads and anonymizes data from a mongo data backup to a mongo database.

required environment var:
`MONGO_SERVER_ADDRESS`  - destination address
`MONGO_USERNAME` - mongo username for the destination
`MONGO_PASSWORD` - mongo password for the destination
`BUCKET_NAME` - s3 bucket to pull the info from

require config settings:
1. config folder
2. config/notifications.json - sample is attached in the sample folder
incoming is currently a must. outgoing is not.
