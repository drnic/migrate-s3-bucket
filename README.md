# Migrate between S3 blobstores on Cloud Foundry

Migrating one or all your Cloud Foundry buckets from RiakCS to AWS S3 or between any two AWS S3-compliant blobstores? Then this app is for you.

## Usage

Push the application to your Cloud Foundry space, bind it to the source & destination service instance (blobstores/buckets), SSH into a container and run the migration command.

```
cf push migrate-s3-bucket --no-start --no-route
cf bind-service migrate-s3-bucket from-bucket
cf bind-service migrate-s3-bucket to-bucket
cf restart migrate-s3-bucket
cf ssh migrate-s3-bucket
```

Then run the application:

```
./migrate-s3-bucket list
```
