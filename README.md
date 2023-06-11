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

## Things to take care of

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





