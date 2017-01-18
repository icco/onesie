# onesie

An attempt to make static site deployment easier.

Instance template created with:

```
$ gcloud compute instance-templates create onesie-template-2 --image-family=onesie --network=onesie --scopes=https://www.googleapis.com/auth/pubsub,https://www.googleapis.com/auth/cloud-platform,useraccounts-ro,storage-full,logging-write,monitoring-write,service-control,service-management,compute-ro --metadata startup-script="$(cat startup-script.sh)"
```
