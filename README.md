# Wait for Reboot Microservice

[![View in Swagger](http://www.jessemillar.com/view-in-swagger-button/button.svg)](https://byuoitav.github.io/swagger-ui/?url=https://raw.githubusercontent.com/byuoitav/wait-for-reboot-microservice/master/swagger.yml)

A service designed to monitor the reboot status of Crestron devices (by testing a TCP connection) and notify the submitter via a POST when reboot is complete.

The service monitors multiple targets concurrently, and will notify the address defined in the submission request when they respond successfully to TCP requests.

The information for the pinging information can be found in the Config file, detailed below.

### Config File
When running the program you can pass in the `-config` flag denoting the location of the JSON config file. If no value is passed in it will default to the provided `./config.json`

```
WaitThreshold is the allowed number of devices being monitored before no wait is imposed between iterations
IterativeTime is the number of seconds to wait between iterations, assuming that the WaitThreshold has not been reached
IndividualTimeout is the number of milliseconds to wait for each individual connection (TCP attempt) before timing out
```
