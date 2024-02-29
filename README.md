# File Difference Finder

This is a simple service to find the file content differences.

## Algorithm 
Finding the file difference is a challenging problem to solve due to computational complexity of algorithms and limited resources of hardware. This is actually a very well known problem and there are some academic papers which offer solutions to this problem. 

As brute force approach can consume too many resources (CPU and Memory), another more performing algorithms should be used. For the solution, instead of checking file differences at each index of original file content and updated file context, Rolling Hash Algorithm is used to hash small parts of file content. (Please visit [Wikipedia](https://en.wikipedia.org/wiki/Rolling_hash) algorithm's explanation for more details). As it is explained in Wikipedia link, this algorithm computes a hash value for a sliding window of data.  

### How does the algorithm of finding the difference between files work in this repo?
Let's assume there are two files named original and updated. First, we apply rolling hash algorithm for both with window size 8. Then comparing the hashes of same indexes in both files, we find the (wider) ranges. After this, we apply rolling hash algorithms for the wider ranges with 2 window size so that finally we have a very narrow length to compare. As the last step, each index in the ranges is to be checked on both files to see if there are changes.
Note: patching is not allowed in this repo but only finding difference. So while running up the service, original file's content and version should be defined (Please see below how to run the service section).

## Endpoints

This repo includes only 4 endpoints:

* _/diff_ endpoint [GET]: Show the difference of given input (updated file content) with the original file.
Lets assume original file's content _"abcabcabcabcabcababc"_ where version is _13_.

Example request body:
```
{
    "text": "abcabcabcabca4cababc3",
    "version": 14
}
```

Example for success response:
```
{
    "delta": [
        {
            "OldValue": "b",
            "NewValue": "4",
            "Index": 13,
            "Type": "UPDATED"
        },
        {
            "OldValue": "",
            "NewValue": "3",
            "Index": 20,
            "Type": "ADDED"
        }
    ],
    "current_version": 13,
    "updated_version": 14
}
```

Example for fail response:
```
{
  "Message": "another process is getting the diff",
  "Code": 400
}
```

* _/live_ endpoint [GET]: This is to be integrated with kubernetes to provide liveliness concept which proves whether the
  service is alive. Unless it returns successful, the pod should restart with kubernetes configurations. Endpoint basically
  returns 200 as http status code to prove the service is alive.
* _/ready_ endpoint [GET]: This one is also for kubernetes integration and is for readiness concept that represents
  whether service can accept traffic. Endpoint checks whether cache mechanism return response for _ReadinessKey_ where we save on running up the service.
* _/metrics_ endpoint [GET]: Shows prometheus metrics of the application.

### How to run the tests
To run all tests in one command:
`go test -p 20 ./...`

### How to run service with docker

1. Make sure to set port, file version and File content in dockerfile.

2. On the project's directory to create an image:

`
docker build -f ./dockerfile -t file-diff-finder-app:latest .
`

3. And to run the container:

`
docker run -p  127.0.0.1:1903:1903 file-diff-finder-app
`
and then you should be able to consume the service at host http://127.0.01:1903/
### Assumptions and considerations

Algorithm based:
* There are some other libraries or algorithms for sequence comparison algorithm. Based on the academical researches, the most performing one can be found and should be in research.
* Patching the original file is not allowed in this repository but only finding the difference. 
* Cache and lock mechanism is used to prevent original's file version. 

Kubernetes based:
* We assume that this service will be deployed to kubernetes as a pod. There can be multiple replicas based on the
  traffic.
* _/live_ endpoint is for liveliness in kubernetes. Once it returns with http code other than 200, pod should restart.
* _/ready_ endpoint is to tell if pod is ready for accepting traffic.
* Data from _/metrics_ endpoint to be consumed by Grafana or one monitoring tool for visibility so that needed.
  parameters could be monitored and related alert notifications are triggered.
* We could run stress tests once kubernetes pods' memory limits are decided.

Hardware based:
* We could have rate limits to prevent extraordinary and unexpected hardware consumption.
* Solution provided here is not the most performing algorithm. There can be a better way to decide wider ranges window size (8), and for lower range(2). For example, lets say file context is 100 length, then 8 for wider range is fine, but for 10000000, 8 would be too much resource consuming. So an algorithm should decide wider range's window size.  
* Maximum 10 second is the threshold to find the difference. Longer processes are stopped immediately to protect hardware.
