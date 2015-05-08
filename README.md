PunditTracker Backend
=================

Backend for PunditTracker.com

###  To Run Locally:

- Get these repo `git clone https://github.com/PunditTracker/webBackend.git`
- Get Golang `brew install go`
- Set GOPATH `export ~Somepath/go`
- Get dependencies `go get -t ./...`
- Make sure Postgressql is installed and running
- use the reload scrip to start running the app in the background `./reloads.h`
- The server is available at `localhost:3000`



Endpoints:
=================
### Auth:
`/v1/auth/register` - Register a new user, requires a json object representing the user

`/v1/auth/registerfb` - Registers a new user, requires a json object representing the user

`/v1/auth/login` - Checks if the user exists with supplied username and password.  If they exist, returns the user struct and initializes a session.  If they don't exists, returns an error.

`/v1/auth/loginfb` - Checks if the user exists with supplied username and fb id.  If they exist, returns the user struct and initializes a session.  If they don't exists, returns an error.

`/v1/auth/logout` - Kills the user session

`/v1/auth/check` - Check if a user session exists

`/v1/auth/forgot` - Begins the reset password flow

`/v1/auth/reset/{id:[0-9]+}/{resetKey}` - Changes the user's password if the reset key matches the reset key for the user with `id` and the key is fresh.

### User:
`GET` `/v1/user` - Returns an array of all users

`PATCH` `/v1/user` - Updates an individual user object

`/v1/user/featured` - Get all users marked as featured

`/v1/user/{id:[0-9]+}` - Get a user with the provided `id`

`/v1/user/{id:[0-9]+}/vote` - Gets all votes for a user with the provided 
`id`

`/v1/user/password` - Updates the user's password

`/v1/putprofpic` - Uploads a supplied image to 

### Prediction: 
`/v1/prediction` - Get all predictions

`/v1/prediction/featured` - Get all featured predictions

`/v1/prediction/{id:[0-9]+}` - Get single prediction where prediction_id= `id`

`/v1/prediction/add` - Add a prediction, requires an authenticated user session and a json object describing the prediction

`/v1/prediction/category/{id:[0-9]+}` - Get all predictions for a category where category_id= `id`
`/v1/prediction/category/{name:[a-zA-z]+}` - Get all predictions for a category where category_name = `name`

`/v1/prediction/page/{cat_id:[0-9]+}` - Get predictions to fill a category page.  Returns homepage predictions if `id` = 0

`/v1/prediction/user/{id:[0-9]+}` - Get all predictions for the user where userid = `id`

`/v1/prediction/tag/{tag}` - Get all predictions that are tagged with `tag`

`/v1/predictions/cat/{id:[0-9]+}/tag/{tag}` - Get all predictions that are both part of the category with categoryid = `id` and tagged with `tag`

`/v1/prediction/vote/{pred_id:[0-9]+}/{value:[0-9]}` - 

`/v1/prediction/vote/{id:[0-9]+}/avg` - Get average vote for the prediction with prediction_id = `id`

`/v1/prediction/vote/getown/{id:[0-9]+}` - Get the current users vote for the prediction with prediction_id = `id`

`/v1/category` - Get list of all live categories

`/v1/hero` - Same as `/v1/hero/0`

`/v1/hero/{id:[0-9]+}` - Gets the list of heros for the category page where category_id = id.  Will return the homepage if id=0

`/v1/predictionSet` - Same as `/v1/predictionSet/0`

`/v1/predictionSet/{id:[0-9]+}` - Gets the list of prediction sets for the category page where category_id = id.  Will return the homepage if id=0

`/v1/event/{event_name:[a-zA-z]+}/{event_year:[0-9]+}`


###Search:
`/v1/user/search/{search_string}` - Searches for users that have names or affiliations similar to search_string. 

`/v1/user/prediction/search/{search_string}` - Searches for predictions made by users that have names or affiliations similar to search_string. 

`/v1/prediction/search/{searchstr}` - Searches for predictions that have similar titles to search_string.

###Admin:

`/v1/admin/hero/set`
`/v1/admin/predictionSet/set`
`/v1/admin/predictionLoc/set`
`/v1/admin/prediction/add`
`/v1/admin/pundit/add`
`/v1/admin/image/add/{type}`
`/v1/admin/hero/{cat_id:[0-9]+}`
`/v1/admin/predictionSet/{cat_id:[0-9]+}`
`/v1/admin/predictionLoc/{cat_id:[0-9]+}`
`/v1/admin/prediction/setstate/{predId:[0-9]+}/{state:[0-9]}`
`/v1/admin/special_event/result/set`