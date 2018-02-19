# Go Example API

This is a simple Go API that allows you to save and fetch data from it. It uses
a Simple approach to save data to the Google App Engine datastore using the key
from it as the resource URL.

Save data by sending a `POST` request to `/save`. While saving, you'll be returned a
`201 Created` header with a `Location` header pointing to the saved resource. The
content saved will be returned with `Content-Type: application/json; charset=utf-8`
which, again, is intentional so people don't save, say, HTML to render a phishing
site. That also means the data sent to `/save` will be checked to see if it can be
parsed as JSON.

An important note is that data is saved in a best-effort basis under the Memcache
service in Google App Engine. By default, we set the expiration to 24 hours but there
are no guarantees that it'll stay for that long.

There's no support for updating resources. If you want to update one, consider creating
a new resource from scratch issuing a `POST` to `/save`.

This API has CORS headers so, yes, you can use them through your Single-Page Application
or another web-based flow.

## Endpoints

As an important side note, illegal, malformed requests that can't happen, like a
`GET` request with a body will be caught by Google's AppEngine proxy and you'll
receive a `400 Malformed Request` from Google. There may be other case scenarios
where this can happen.

* `POST /save`: Creates a new resource. The response will contain no body, but a
  status code `201 Created` and a `Location` header pointing to `/get/:id`.
* `GET /get/:id`: Where `:id` is the alphanumeric ID of the given resource you're
  fetching. `GET` will fetch the whole data saved to the storage you created when
  issuing the `POST` to `/save`. This endpoint does not require the `Authorization`
  header.
* `ALL /debug`: Where `ALL` can be any verb except `HEAD`. This will take all
  headers, body and querystring parameters passed and return them to the caller
  in the body. The response will be in the format of:
  `{ "headers": [], "body": {}, "query": [] }`. The querystring modifiers are
  removed for simplicity.

## Modifiers

The API supports two modifiers to the response. These are set so you can use them
to debug asynchronous calls to this API with expected responses. The modifiers are
passed as querystring parameters that are removed from the response body on `/debug`.

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
