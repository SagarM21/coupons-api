# Coupons Management API

A RESTful API for managing and applying discount coupons for an e-commerce platform, built with Go, Gin, and ScyllaDB.

## Project Structure

```
e-commerce/
├── cmd/server/main.go        # Application entry point
├── api/routes/coupon.go       # HTTP handlers and route registration
├── pkg/
│   ├── config/config.go       # YAML-based configuration
│   ├── db/
│   │   ├── db.go              # Database connection, keyspace, and table setup
│   │   └── coupon.go          # Coupon CRUD operations (repository layer)
│   ├── models/
│   │   ├── coupon.go          # Coupon domain models
│   │   └── cart.go            # Cart and discount response models
│   └── services/
│       ├── coupon.go          # Business logic (discount calculation, validation)
│       └── coupon_test.go     # Unit tests
├── config.yml                 # Application configuration
├── docker-compose.yml         # Docker Compose (API + ScyllaDB)
├── Dockerfile                 # Multi-stage Go build
└── README.md
```

## Database Schema

ScyllaDB (Cassandra-compatible) with a single table:

```sql
CREATE TABLE IF NOT EXISTS coupons (
    id         TEXT PRIMARY KEY,
    type       TEXT,          -- "cart-wise", "product-wise", "bxgy"
    details    TEXT,          -- JSON string holding type-specific config
    created_at TIMESTAMP,
    expires_at TIMESTAMP      -- optional, NULL means no expiry
);
```

The `details` column stores JSON specific to each coupon type:

| Type | Details JSON |
|------|-------------|
| cart-wise | `{"threshold": 100, "discount": 10}` |
| product-wise | `{"product_id": 1, "discount": 20}` |
| bxgy | `{"buy_products": [...], "get_products": [...], "repition_limit": 2}` |

## Configuration

`config.yml`:

```yaml
env: development

api:
  host: 0.0.0.0
  port: 8080

database:
  host: scylla       # service name in docker-compose
  port: 9042
  keyspace: coupons
  consistency: one
```

## Quick Start

### Prerequisites
- Docker and Docker Compose

### Run

```bash
docker compose up --build
```

The API starts on `http://localhost:8080` once ScyllaDB is healthy.

### Health Check

```bash
curl http://localhost:8080/health
```

## API Endpoints

| Method | Endpoint               | Description                                    |
|--------|------------------------|------------------------------------------------|
| POST   | /coupons               | Create a new coupon                            |
| GET    | /coupons               | Retrieve all coupons                           |
| GET    | /coupons/:id           | Retrieve a specific coupon                     |
| PUT    | /coupons/:id           | Update a specific coupon                       |
| DELETE | /coupons/:id           | Delete a specific coupon                       |
| POST   | /applicable-coupons    | Get all applicable coupons for a cart           |
| POST   | /apply-coupon/:id      | Apply a specific coupon to a cart               |

## Coupon Types

### 1. Cart-wise
Applies a percentage discount to the entire cart when the total exceeds a threshold.

```bash
curl -X POST http://localhost:8080/coupons \
  -H "Content-Type: application/json" \
  -d '{
    "type": "cart-wise",
    "details": {
      "threshold": 100,
      "discount": 10
    }
  }'
```

### 2. Product-wise
Applies a percentage discount to a specific product.

```bash
curl -X POST http://localhost:8080/coupons \
  -H "Content-Type: application/json" \
  -d '{
    "type": "product-wise",
    "details": {
      "product_id": 1,
      "discount": 20
    }
  }'
```

### 3. BxGy (Buy X, Get Y)
Buy a specified quantity of products from a "buy" set and get products from a "get" set for free, with a repetition limit.

```bash
curl -X POST http://localhost:8080/coupons \
  -H "Content-Type: application/json" \
  -d '{
    "type": "bxgy",
    "details": {
      "buy_products": [
        {"product_id": 1, "quantity": 3},
        {"product_id": 2, "quantity": 3}
      ],
      "get_products": [
        {"product_id": 3, "quantity": 1}
      ],
      "repition_limit": 2
    }
  }'
```

### Get Applicable Coupons

```bash
curl -X POST http://localhost:8080/applicable-coupons \
  -H "Content-Type: application/json" \
  -d '{
    "cart": {
      "items": [
        {"product_id": 1, "quantity": 6, "price": 50},
        {"product_id": 2, "quantity": 3, "price": 30},
        {"product_id": 3, "quantity": 2, "price": 25}
      ]
    }
  }'
```

### Apply a Coupon

```bash
curl -X POST http://localhost:8080/apply-coupon/8c4eea1f-70ff-41eb-8bd2-481255917715 \
  -H "Content-Type: application/json" \
  -d '{
    "cart": {
      "items": [
        {"product_id": 1, "quantity": 6, "price": 50},
        {"product_id": 2, "quantity": 3, "price": 30},
        {"product_id": 3, "quantity": 2, "price": 25}
      ]
    }
  }'
```

## Running Tests

```bash
go test ./pkg/services/ -v
```

## Implemented Cases

### Cart-wise Coupons
- Cart total exceeds threshold: discount applied as percentage of entire cart total
- Cart total equals or is below threshold: coupon not applicable
- Discount distributed proportionally across items when applied
- Multiple cart-wise coupons can each be listed as applicable independently

### Product-wise Coupons
- Target product present in cart: discount applied as percentage of that product's total (price * quantity)
- Target product not in cart: coupon not applicable
- Works with any quantity of the target product

### BxGy Coupons
- All "buy" products must be present in the cart with sufficient quantities
- Free quantity capped by actual cart quantity of "get" products (can't get more free than you have in cart)
- Repetition limit caps how many times the deal can repeat
- Multiple "buy" products: all must meet their quantity requirement per repetition
- Multiple "get" products: all eligible "get" products get the discount
- If a "get" product is not in the cart, it's simply skipped (other get products still get discount)

### Coupon Expiry
- Coupons can have an optional `expires_at` field
- Expired coupons are excluded from applicable coupons
- Applying an expired coupon returns an error

### Validation
- Cart-wise: threshold must be non-negative, discount must be 0-100
- Product-wise: product_id must be positive, discount must be 0-100
- BxGy: buy_products and get_products must be non-empty, repition_limit must be positive
- Invalid coupon type is rejected
- Empty cart is rejected

## Unimplemented Cases (Considered)

### Coupon Stacking / Combination Rules
- Allowing multiple coupons to be applied to the same cart simultaneously
- Rules for which coupons can stack (e.g., cart-wise + product-wise but not two cart-wise)
- Priority ordering when multiple coupons conflict
- "Best coupon" auto-selection

### User-based Constraints
- Per-user usage limits (e.g., coupon can only be used 3 times by a single user)
- First-time user only coupons
- User tier/membership requirements (e.g., premium members only)
- Coupon codes (alphanumeric codes users enter at checkout)

### Cart Constraints
- Minimum number of distinct items in cart
- Maximum discount cap (e.g., "20% off up to Rs. 500 max")
- Minimum quantity of items required
- Exclude certain product categories from discount

### Time-based Constraints
- Start date (coupon not valid before a certain date) - only expiry is implemented
- Day-of-week restrictions (e.g., only valid on weekends)
- Time-of-day restrictions (e.g., happy hour coupons)
- Flash sale coupons with very short validity windows

### Usage Limits
- Global usage count (coupon can only be used N times total across all users)
- Budget cap (stop coupon after total discount given reaches a limit)

### BxGy Advanced Cases
- "Buy X from category A, get Y from category B" (category-based instead of product-based)
- Partial application: if cart has 5 of a buy product but needs 3 per repetition, apply once and ignore the extra 2
- Overlapping buy and get sets (same product in both buy and get arrays)
- Weighted selection: when multiple get products are available, choose cheapest/most expensive to give free
- "Buy X, Get Y at Z% off" instead of completely free

### Other Coupon Types (Not Implemented)
- Flat discount: Rs. 50 off the cart (instead of percentage)
- Tiered discounts: 5% off above Rs. 100, 10% off above Rs. 500, 15% off above Rs. 1000
- Free shipping coupons
- BOGO (Buy One Get One) as a simplified BxGy
- Bundle discounts: buy products A + B + C together for a special price
- Referral coupons: generated per-user for referral programs
- Loyalty points multiplier coupons

### Concurrency / Race Conditions
- Two users applying the same limited-use coupon simultaneously
- Cart modification between checking applicability and applying

## Assumptions

1. **Product prices are provided in the cart request** - the API does not maintain a product catalog; prices come from the client
2. **Single coupon application** - only one coupon can be applied at a time per request; stacking is not supported
3. **No authentication** - the API does not have user authentication; any client can create/apply coupons
4. **No product catalog** - product IDs are integers provided by the client; no validation against a product database
5. **Percentage-based discounts only** - all discount values represent percentages (0-100) for cart-wise and product-wise; BxGy gives items for free
6. **Cart is stateless** - the cart is provided in each request; the API does not persist cart state
7. **BxGy "get" products must be in the cart** - free products are only discounted if they're already in the cart (the API doesn't add items)
8. **Currency is not tracked** - all prices are plain numbers without currency specification

## Limitations

1. **Full table scan for applicable coupons** - `GetAllCoupons` scans the entire coupons table; this will not scale well with a large number of coupons. A secondary index or materialized view on coupon type would help.
2. **No pagination** - GET /coupons returns all coupons without pagination
3. **Details stored as JSON string** - coupon details are stored as a JSON text column in ScyllaDB rather than using UDTs or separate tables, which limits query flexibility
4. **No coupon code support** - coupons are identified by UUID, not user-friendly codes
5. **No usage tracking** - there's no record of how many times a coupon has been used
6. **Single-node ScyllaDB** - docker-compose runs a single ScyllaDB node with minimal resources; not production-ready
7. **No rate limiting or API authentication** - the API is open to all requests
