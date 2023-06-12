## Hierarchies
Chain represents a store chain.
Store represents the physical stores of chain.
Inventory Level represents the amount of an specific Inventory Item per Store.
Inventory Item represents a physical item.

## Methods
```
// GET    /chain                                             : returns all the chains
// GET    /chain/:chainid                                    : returns a single chain by ID
// GET    /chain/:chainid/stores                             : returns all the stores for a chain ID

// GET    /stores/:storeid                                   : returns a single store by ID
// GET    /stores/:storeid/inventory                         : returns all the inventory levels for a store
// GET    /stores/:storeid/inventory?items=itemids[,]        : returns the inventory levels for specific items for a store
```

## Tables
CHAIN
```
id: integer,
name: string,
phone: string,
created_at: timestamp,
updated_at: timestamp,
```

STORE:
```
id: integer,
chain_id: integer,
name: string,
phone: string,
address: string,
postal_code: string,
city: string,
country: string
created_at: timestamp,
updated_at: timestamp,
active: bool,
```

INVENTORY_LEVEL:
```
store_id: integer,
item_id: integer,
quantity: integer,
updated_at: timestamp,
```

ITEMS
```
id: integer,
sku: string,
name: integer,
cost: float,
created_at: timestamp,
updated_at: timestamp,
```

### Implementation
```
// By defining CommerceStore as an interface it can be mocked for testing.

// CommerceStore will store all the methods to our database
type CommerceStore interface {
    GetChains() []Chain
    GetChain(id int) Chain
    GetStores(chainid int) []Store
    GetStore(storeid int) Store
    GetInventory(storeid int) []Inventory
    GetInventoryForItems(storeid int, items []int) []Inventory
    GetItem(id int) Item
}

// For local try-it-out we implement an InMemoryCommerceStore
type InMemoryCommerceStore struct {
	mu     sync.Mutex
	chains map[string]Chain
}

// CommerceServer will be the server.
// It hosts the data store as an interface (which allows to add any implementation).
// It hosts a Handler interface so it can handle any handler.
type CommerceServer struct {
    store CommerceStore
    server *http.Server
}

// Initializes the server and can receive any store type
func NewCommerceServer(store db.CommerceStore, port string) *CommerceServer {
```

## Database
Let's use cockroachdb

Create a volume named roach
```
docker volume create roach
```

Create a bridge network named mynet
```
docker network create -d bridge mynet
```

Start the database container
```
docker run -d \
  --name roach \
  --hostname db \
  --network mynet \
  -p 26257:26257 \
  -p 8080:8080 \
  -v roach:/cockroach/cockroach-data \
  cockroachdb/cockroach:latest-v20.1 start-single-node \
  --insecure
```

Access the database
```
docker exec -it roach ./cockroach sql --insecure
```

To create the database and the tables look at the `init.sql` file.

Create a user
```
CREATE USER kenshin;
```

Grant all permission to this user (for testing only):
```
GRANT ALL ON DATABASE commercedb TO kenshin;
```

If the user was created after some tables were created in the database, then need to use wildcard to apply retroactively.
```
GRANT ALL ON commercedb.* TO kenshin;
```

## Docker
Build the image
```
docker build --progress plain --no-cache --tag shopa -f Dockerfile .
```

Run the image, and specifying environment variables
```
docker run --publish 5000:5000 -e PORT=5000 shopa
```

## Deploy into Cloud Run
Setup Google Cloud Run for local development: https://cloud.google.com/run/docs/setup

Deploying the current directory as source
```
gcloud run deploy
```

Authenticating to Cloud Run with the logged in user: https://cloud.google.com/run/docs/authenticating/developers#curl
```
curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" https://shopa-sfy2ragm5a-uc.a.run.app/chains
```

Image URL
```
us-central1-docker.pkg.dev/thejam/cloud-run-source-deploy/shopa:latest
```

## Things to take care of
[x] Create a basic server
    [x] Add basic routing
    [x] Add basic logging
    [x] Setup graceful shutdown
        [x] Run the server on its own goroutine
        [x] Have the main process listen for interrup or kill signals
        [x] Configure graceful shutdown steps
    [ ] Properly manage the end to end contexts (server and database)
    [ ] Properly manage server setup (timeouts, etc.)
[x] Create a basic storage layer
    [x] Add an in memory data store
    [x] Add basic database and its connection
[x] Setup basic testing
    [x] Test the server handlers
    [x] Be able to mock the database
[ ] Deploy into Cloud Run
    [ ] Manual deploy the server into Cloud Run
    [ ] Setup a database in the cloud
    [ ] Automatically deploy the server from a CI/CD pipeline


### User management
1. How to setup authentication
1. How to keep each user data isolated

### Server Storage
1. How to setup local in-memory storage
1. How to setup centralized storage in a database
1. How to set guarantees long term on user data

### Server Logging
1. How to design and setup structured logging in the server
1. About logging best practices: https://www.datadoghq.com/blog/go-logging/
1. Implement a standard logging interface: https://www.datadoghq.com/blog/go-logging/#implement-a-standard-logging-interface
    1. Including IDs for each type of log (enums)

### Server Tracing and Observability
1. How to setup observability through tracing and metrics

### Server Testing
1. Create tests for the http methods
1. Create tests for storage layer

### Deployment
1. How to deploy the server
1. How to deploy new versions of the servers without impacting current functionality or users
1. How to deploy new features behind feature flags

### Load Testing
1. How to load test the server (`https://github.com/bojand/ghz`)





