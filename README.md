# Go Example API

This is a simple Go API that allows you to save and fetch data from it. It uses
a Simple approach to save data to the Google App Engine datastore using the key
from it as the resource URL.

To use it, you need to get a token first. Each token expires in 3 hours and currently
I offer no way to refresh it. This is intentional so people always have to create new
tokens and they don't abuse the service.

The tokens are used with the `Authorization: Bearer AABBCC` header so simply send
the request with the given header, replacing `AABBCC` with the proper token, and your
request will be fulfilled.

Save data by sending a `POST` request to `/save`. While saving, you'll be returned a
`201 Created` header with a `Location` header pointing to the saved resource. The
content saved will be returned with `Content-Type: application/json; charset=utf-8`
which, again, is intentional so people don't save, say, HTML to render a phishing
site.

There's no support for updating resources. If you want to update one, consider creating
a new resource from scratch issuing a `POST` to `/save`.

This API has CORS headers so, yes, you can use them through your Single-Page Application
or another web-based flow.


## Endpoints

* `POST /save`: Creates a new resource. The response will contain no body, but a
  status code `201 Created` and a `Location` header pointing to `/get/:id`. This
  requires an `Authorization` token to be set.
* `GET /get/:id`: Where `:id` is the alphanumeric ID of the given resource you're
  fetching. `GET` will fetch the whole data saved to the storage you created when
  issuing the `POST` to `/save`. This endpoint does not require the `Authorization`
  header.
* `ALL /debug`: Where `ALL` can be any verb except `HEAD`. This will take all
  headers, body and querystring parameters passed and return them to the caller
  in the body. The response will be in the format of:
  `{ "headers": [], "body": {}, "query": [] }`. The querystring modifiers are
  removed for simplicity.
* `(GET|POST) /authorize`: This creates a JSON Web Token that expires in 3 hours
  and you can use to create resources by issuing `POST` requests to `/save`. The
  content type of the response is simply `text/plain` with some instructions and
  a pointer to this documentation. This endpoint isn't affected by modifiers (see
  below to understand about modifiers).

## Modifiers

The API supports two modifiers to the response. These are set so you can use them
to debug asynchronous calls to this API with expected responses. All requests except
for `/authorize` supports them. The modifiers are passed as querystring parameters
that are removed from the response body on `/debug`.

* `delay=${time}` where `${time}` can be any number of seconds up to 50 seconds max,
  this due to App Engine restrictions (more on that below). This utility is useful for
  debugging while making requests that may take longer, and although Go supports
  multiple time units, App Engine responses are limited to 60 seconds max. We set 50
  seconds as maximum so the remaining 10 are available to perform the request processing.
  Any other number will be simply ignored.
* `status=${number}` where `${number}` can be any number. This will set the response
  status code to whatever is passed here. For example, a given flow could be you wanting
  to test your application against 404 resources. Create a resource with `/save`, and
  include in the body the expected error message from your API. Then issue a `GET`including
  the status override like `/get/:id?status=404`. You'll get the body (which should be
  the error JSON data you saved before) but with a 404 status code. Status codes are ommitted
  if `< 100` or `> 599`.

## Deploying and developing this codebase

To deploy the code, simply create an app in Google App Engine. That'll give you an ID, in this
case is `go-example-api`. Then, to deploy, run:

```bash
gcloud app deploy appengine/app.yaml --project=go-example-api
```

Remember to replace `go-example-api` with your own ID -- since you can't post to mine.

To develop locally, you can install the Google SDK -- which you'll have if you install the
`gcloud` binaries. Then, use the Python file in your installation folder, `dev_appserver.py`
to run the local dev server, as follows:

```bash
dev_appserver.py appengine/app.yaml
```
