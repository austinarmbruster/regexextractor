# Extracting all the Things

The Regex Extractor is a simple HTTP server that will process a text input and provide a simple collection of label subsets that match configured regular expressions.

## Usage

Once the application is built, there are two command line flags:

1. -addr:  The address on which the HTTP server will be bound
1. -file:  The CSV containing the labels and patterns

Example:

```
regexextractor -file ./sample.csv
```

Example Call:
```
curl http://localhost:8080 -d 'Bill Gates bill@ms.com'
```

Example Output:
```
{"email":["bill@ms.com"],"name":["Bill","Gates"]}
```

## Future Work

* Use github.com/spf13/cobra to handle command line arguments
* Update the build to set to embed a version
* Configure the application to use TLS