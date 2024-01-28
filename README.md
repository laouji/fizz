# FIZZ

## Prerequisites

* docker compose must be available on your local machine

## Execution

The most basic usage is to run the web server using the up command:
```
$ make up
```

The application will then be running in a docker container that exposes port `5000`

## Usage

#### Fizzbuzz

The main endpoint runs a basic fizz-buzz program which allows the user to freely set which multiples they would like and customise the "fizz" and "buzz" strings.

```
$ curl 'localhost:5000/?int1=3&int2=5&str1=fizz&str2=buzz&limit=16'
["1","2","fizz","4","buzz","fizz","7","8","fizz","buzz","11","fizz","13","14","fizzbuzz","16"]
```

| query parameter | description | type | min | max |
| --- | --- | --- | --- | --- |
| int1 | the first multiple | number || 1 | 1999999 |
| str1 | the string to replace each multiple of int1 | string || 1 | 30 |
| int2 | the 2nd multiple | number || 1 | 1999999 |
| str2 | the string to replace each multiple of int2 | string || 1 | 30 |
| limit | the number of results expected | number || 0 | 1999999 |

#### Stats

To view the parameters corresponding to the most requested request a simple GET request can be made to the /stats endpoint

```
$ curl 'localhost:5000/stats'
[{"query":{"int1":["3"],"int2":["6"],"limit":["36"],"str1":["fizz"],"str2":["buzz"]},"hit_count":8}]
```

## Testing

The following command will run all unit tests locally on the user's machine

```
$ make test
```
