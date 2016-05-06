# Listen for Reboot Microservice

A service designed to monitor the restart status of devices (by testing a TCP connection) and notify the submitter when they are complete. Designed specifically for Crestron devices.

The service monitors multiple targets concurrently, and will notify the address defined in the submission request when they respond successfully to TCP requests.

The information for the pinging information can be found in the Config file, detailed below.

## Endpoints

`POST /submit`

Submission happens here, the expected JSON payload should be in the form of

```
{
  "IPAddressHostname": "string",
	"Port": int,
	"Timeout": int,
	"CallbackAddress": "string",
	"Identifier": "Optional string"
}
```

* IPAddressHostname is the address of the machine you wish to listen for.
* Port is the port to try connections over
* Timeout how long you want to try for in seconds - will default to 300 if no value given
* Callback Address the address the service will send a POST request to to notify of either connection or timeout.
* Identifier optional field that will be passed back to you to aid in the identification of the device passed in.

When the the device has either successfully responded or the timeout has elapsed the service will send a post request to the address given in CallbackAddress. The request will have a JSON payload containing the values passed in and additional information in the form of

```
{
    "IPAddressHostname": "The value passed in.",
    "Port": the value passed in,
    "Timeout": the value passed in,
    "CallbackAddress": "The value passed in.",
    "SubmissionTime": "The time, in RFC3339 format, the address was submitted for observation.",
    "CompletionTime": "The time, in RFC3339 format, the device responded or timeout was reached.",
    "Status": "Success if the device responded. Timeout if timed out.",
    "Identifier": "The value passed in."
{
```

## Config file

When running the program you can pass in the `-config` flag denoting the location of the JSON config file. If no value is passed in it will default to `./config.json`

The config file should be in the form of

```
{
  "WaitThreshold": int,
  "IterativeTime": int,
  "IndividualTimeout": int
}
```

* WaitThreshold is the number of devices being monitored required before no wait is imposed between iterations.
* IterativeTime is the number of seconds to wait between iterations, assuming that the WaitThreshold has not been reached.
* IndividualTimeout is the number of milliseconds to wait for each individual connection (TCP attempt) before timing out.
