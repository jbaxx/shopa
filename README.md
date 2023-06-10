




## Hierarchies
Store Chain represents a store chain.
Location represents the physical stores of a store chain.
Inventory Level represents the amount of an specific Inventory Item per Location.
Inventory Item represents a physical item.

## Methods
```
// GET    /chain                                                  : returns all the chains
// GET    /chain/:chainid                                         : returns a single chain by ID
// GET    /chain/:chainid/stores                                  : returns all the stores for a chain ID

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
quantity: integer,
created_at: timestamp,
updated_at: timestamp,
```
