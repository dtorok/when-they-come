# When they come

## Introduction

"When it comes" is an application that helps people figure out when their next bus or tube comes. After locating the user It shows the stop points close by, and clicking/tapping on them they can see when the next vehicles will arrive and which direction they will go.

The application is hosted on google cloud and can be found on the following url: https://when-they-come.appspot.com

It will do geolocation at first and navigate the map to your location. If you don't want that, use the https://when-they-come.appspot.com#noloc URL instead that will stay on the default location, which is London.

Creating the service I focused mostly on the backend side, the frontend is just an embedded google maps component.

I'm Daniel Torok, my linkedin profile can be found here: https://www.linkedin.com/in/danieltorok

## Installation

To install and run the application locally, clone the github repository and step into the directory:

````bash
$ git clone https://source.developers.google.com/p/when-they-come/r/default
$ cd default
````

You may need google cloud credentials for that.

Create a google API key and put it to the `etc/google_api_key` file to make google maps work.

```bash
echo YOUR_GOOGLE_API_KEY > etc/google_api_key
```

Set your `GOPATH` environment variable and you're ready to run:

```bash
$ export GOPATH=`pwd`
$ go run main.go
```

Now you can open the application from your browser on the http://localhost:8080 URL.

## Project structure

### Client side

Given that a map seemed to be a natural UI for the task, google maps looked like an obvious choice to build the frontend on. For the ajax communication I used jquery.

There are 3 components on the javascript side:

-   `static/geolocation.js` - does geolocation with the help of the browser
-   `static/map.js` -  is responsible for handling the embedded google map component
-   `main.html` - the glue and the ajax communication remained there

### Server side

I took the opportunity and for the backend's language I chose go, the language I always liked but never had the chance to write too much code in. The backend code has 3 modules:

-   `web` - is responsible for the HTML generation and serving the additional static files (images, javascript files)
-   `api` - serves the JSON API the javascript is communicating with
-   `remote`  - does the communication with the external public transport APIs

#### Web

The web module is really simple, contains `web.go` serves the static files and the index page.

This part needs configuration, the google api key has to be generated into the template. Given that this is the only place that needed any settings parameter, I cut corner here and simply read it from a file in the `etc` directory. No `.ini` parsing, no config context passed over, in case of more parameters this definitely should be improved.

#### Api

To keep the `main.go` file simple, the `api` module provides an `AddHandlers()` function that registers both the healthcheck (`api/healthcheck.go`) and the api (`api/api_v1.go`).

The API operates with 2 data structures, `Stop` and `Arrival`, both of them coming from the `remote` module, though (see the "Worst part" section). There is only one handler registered, but it routes the request to 2 different views based on the incoming URL. This handler is wrapped by a function (`decorator()`) that does the JSON marshalling and transforms the errors into `500 Internal server error` HTTP responses.

The `stopsByPosition()` provides a list of stop points around a location that's given by query parameters (`?lat=1&lon=2`). Those are converted to floats and passed to the underlaying `remote` layer that gives back the actual data. Conceptually this is a list endpoint for the `stop` entity, so we could provide different filtering solutions here, but only this central based is implemented. One has to be careful in extending this, though, given that it may put extra requirements to the external transport API providers (i.e. `remote` module).

The `arrivalsByStop()` function lists the vehicles arriving in the near future to a given stop point.

Both the methods belong to the `BackendApi` structure that contains the reference to the underlying remote layer, the functions use this instead of calling methods from the `remote` module directly (see the "Best part" section).

#### Remote

This is the module that's responsible for communicating with the external API providers and providing a unified API for the `api` layer. The `TransportAPI` interface (that can be found in the `remote/remote.go` file) has 3 different implementations:

-   `LondonTransportAPI` in the `remote/uk_london.go` - this is using the https://api.tfl.gov.uk/ API, our service's minimalistic API was mostly inspired by the features of this
-   `BudapestTransportAPI` in the `remote/hu_budapest.go` - this is using an undocumented API of the http://futar.bkk.hu service for listing the arrivals, but for the stop points it uses it's own database, a json file (see the "Worst part" section).
-   The 3rd implementation is the `CoordBasedCombinedTransportAPI` in the `coord_based_combiner.go` that decides by the original `lat` and `lon` coordinates which external provider to use. Takes the `Stop` list and patches their `ID`, serializing the provider into it (putting `BUD` or `LON` as a prefix), so it knows where to go for the arrival list later.

All the outgoing communication is done by the `httpJsonGET()` method that's implemented in the `remote/helpers.go` file.

### Testing

I didn't spend time on figuring out how to test the javascript code.

For the go code I used the builtin testing framework, there are unit tests for the `api` and `remote` layer. Some integration tests could also be implemented by calling the views and mocking the outgoing http calls.

In the tests I followed the given-when-then structure, I think it helps keeping the otherwise chatty test code readable.

## Evaluation

In this section I'd like to talk about what I enjoyed the most and I think that is worth a look, those are collected in the "Best part" section.

I also collected the most clumsy things in the "Worse part", however, beyond only complaining I also elaborate on how it could be improved. Still it's like a confession :)

And finally there is a "Missing parts" section where I listed all the things I didn't implement, giving details about why not.

### Best part

#### The decorator in the `api_v1.go` file

I think this is a common pattern in go, still it was great to realize how much boilerplate code can be moved into a method like that. In our case there is only one handler, so this could've been implemented in the routing function, however, the decorator itself is generic enough to be usable for any JSON API.

Given that it does the JSON marshalling of the response, sets the `Content-type` header and transforms the possible errors into `500 Internal server error` HTTP responses, the real handlers could be kept simple, readable and type safe. This makes them more testable also.

#### Interface based modules

Both the `api` and the `remote` layers provide an object for the functionality, not just individual methods. I guess that this must be a common pattern too, for me it was great to figure it out and leverage the benefits of that, which are:

-   The dependencies can be injected this way, making the testing with mock replacements possible
-   The `TransportAPI` can have multiple implementations

This is what made it possible at first to integrate the Hungarian provider, and also to implement the solution that combines the different APIs, covering bigger chunks of the globe (although Budapest is not the biggest chunk of the world :))

An other benefit of this is that the http module used for the outgoing communication can be replaced too (for that I introduced a `HttpClient` interface in the `remote/remote.go` file), that also helps the testing, but makes it possible to improve the service-to-service communication later without touching the module itself. Improvements could be:

-   extended logging
-   special routing logic
-   retries
-   reaction on backpressure
-   circuit breaking
-   etc.

### Worst part

Hmm, there are more items here than in the previous section... :)

#### The `api` module doesn't have it's own structures

The structures the `api` module uses to generate the JSON response come from the `remote` module. This is not a problem now when the project is small, but with moving that module into a separate library, the control over the response could be moved out of this codebase too. In this case it could happen that someone accidentally changes the structure in the library, and whenever we bump it in this service to use the new version, we could break the API without even noticing it.

I left it like this because I didn't want to copy the data one more time in every request, and the `remote` module needed a common structure anyway. Perhaps we could apply some internal datastructure on the response in the `api` module to ensure the correctness in compile time without copying the data.

#### Radius calculation in the `BudapestTransportAPI` solution

While the API provider for London has an endpoint for listing the stop points around a given coordinate, in case of Budapest I had to filter them manually from our own database. Beyond the fact that it's using a bounding box instead of a circle (that I think is an acceptable tradeoff for performance), the calculation of it's size lacks any scientific approach what so ever, I just took 2 points on the map that looked like a good bounding box size, calculated their distance in degree and apply that to the given center point.

#### Tests in the `remote` module

I don't like how the tests look like because they are very similar, most probably they could be simplified. Plus, given that both the london and budapest solution is in the same package, their tests are in the same package too. This lead to different name collition problems, so I had to prefix a lot of functions that made the code less readable. 

The different implementations could be moved to different packages, so they can manage their own namespaces. However, this would mean that the common data structures (the `Stop` and the `Arrival`) should be moved to a 3rd module to avoid circular dependency between the modules. Or there is a solution I couldn't figure out :)

#### Budapest's stop point database committed into the repo

So, for the budapest public transport there is no official real-time API, only an undocumented one, and it's only for the arrivals, not for the stop points. The stop points can be downloaded from a website in a zipped xml format. This is what has been transformed (not by me, I only borrowed it) into a JSON file.

This is basically the stop-point database, a big JSON file I simply committed into the repository. This wasn't the best solution (the fastest, though), most probably I wouldn't do it in a real service, although, I think it's not that evil either. Given that the database

-   doesn't change often and is not too big
-   it makes the application much easier to install
-   doesn't need an external database (cost saving)
-   can leverage the benefits of SVCs (versioning, review, blame)

it can even make sense to do it in production.

Just a side note that the database is unfortunately very noisy, some entries appear several times, and there are stop points where no bus stops ever. This could be cleaned up too.

### Missing parts

#### Distance in the Budapest implementation

The stop point API endpoint shows the calculated distance from the center. This is something the Londond external API gives us, but I didn't implement it for Budapest. It mustn't be too hard, but it needs digging into the Global Grid system.

#### Missing IDL

I originally wanted to write an IDL file for the API, I even created a swagger file, but than I realized that to make it actually useful would need much more work (exporting it over the API, allowing CORS to use from an external swaggerui, etc.), so I eventually decided to remove it from the project. Writing it again anytime is not a big deal, but leaving there without real value felt more like noise in the codebase.

#### Minimal logging

I added some minimal logging that could be good for basic montoring (calculating the amount of requests), but it could be definitely improved especially with timing informations. Let alone an access log. :)

#### No config system at all

There was only one place where I needed a config variable (in the `web` module), and there I just read it from a file. If the application evolved, introducing a config system would be necessary.

#### No web technologies for production environment

-   now this application serves the static files (images, js and css files), there is no CDN integration
-   there is no javascript or css compression
-   I didn't add any web tracking solution (google analytics, etc.)

These are things that would be necessary for a production service but looked like way overengineering for a project like this :)

