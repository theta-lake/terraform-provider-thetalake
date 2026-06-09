# Copilot Instructions

## Adding a New Resource

Use these steps when adding a Terraform resource for a new Theta Lake API entity. Follow the patterns established by `label`, `directory_group`, and `retention_library`.

The acceptance test is part of the initial resource implementation and should be created in the same change as the resource.

---

### 1. Consult the API spec

Reference the OpenAPI spec at `theta_lake_api_v1.yml` (local copy) or `https://developer.thetalake.com/api/v1/reference/theta_lake_api_v1.yml`.

Identify the CRUD endpoints and their request/response schemas for the entity, e.g.:

- `GET /things` — list (used by data sources)
- `POST /things` — create
- `GET /things/{id}` — read by ID
- `PUT /things/{id}` — update
- `DELETE /things/{id}` — delete

---

### 2. Add (or expand) the client implementation

File: `internal/client/thetalake/<entity>.go`

- Define the Go struct for the entity, with all JSON-tagged fields matching the API response schema.
- Use `*string` (pointer) for nullable optional string fields (e.g. `ExternalId *string`).
- Use `time.Time` for timestamp fields returned by the API.
- Implement `Create<Entity>`, `Get<Entity>ById`, `Update<Entity>`, `Delete<Entity>` methods on `*Client`.
  - Pass a struct as the request body for POST/PUT.
  - Extract the response from the named key via `doRequest(..., "entity_name", &response)`.
  - For DELETE, pass `""` and `nil` as the last two arguments.
- If the entity uses a list endpoint for lookup (e.g. data sources), implement `Get<Entity>ByName` using `doRequest` with the list endpoint.

Example pattern:

```go
func (c *Client) CreateThing(ctx context.Context, thing Thing) (Thing, error) {
    var response Thing
    err := c.doRequest(http.MethodPost, "/things", thing, "thing", &response)
    if err != nil {
        return Thing{}, err
    }
    return response, nil
}
```

---

### 3. Add client test data

Directory: `internal/client/thetalake/testdata/`

Add one JSON fixture per operation that returns a body:
- `create_<entity>_response.json` — matches the POST 201 response schema
- `get_<entity>_by_id_response.json` — matches the GET 200 response schema
- `update_<entity>_response.json` — matches the PUT 200 response schema

Use real field names and example values from the API spec.

---

### 4. Add client unit tests

File: `internal/client/thetalake/<entity>_test.go`

Write a `TestCreate<Entity>`, `TestGet<Entity>ById`, `TestUpdate<Entity>`, and `TestDelete<Entity>` function for each CRUD method. Use `newTestClient(t, http.MethodXxx, "/path", handler)` from `test_util.go`, serve the fixture JSON, and assert key fields on the returned struct.

---

### 5. Create the resource package

Directory: `internal/resources/<entity>/`

Create three files:

#### `<entity>.go` — CRUD logic

- Package name: `<entity>` (snake_case, no underscores in Go package name, e.g. `retentionlibrary`)
- Embed `*thetalake.Client`
- Implement `resource.Resource`: `Metadata`, `Configure`, `Schema`, `Create`, `Read`, `Update`, `Delete`, `ImportState`
- In `Create`/`Update`: read individual plan attributes via `req.Plan.GetAttribute(ctx, path.Root("field"), &plan.Field)...`
- In `Read`: read state via `req.State.Get(ctx, &state)`, call `Get<Entity>ById`, map to state model
- In `ImportState`: parse the string ID to `int64` with `strconv.ParseInt`, then `resp.State.SetAttribute(ctx, path.Root("id"), id)`
- Fields not returned by the API (e.g. write-only inputs) must be preserved from the plan/state to keep Terraform state consistent

#### `model.go` — plan and state models

- `<entity>PlanModel` — only the writable fields (those sent in create/update requests)
- `<entity>StateModel` — all fields including computed read-only ones
- `toApiModel(*planModel) thetalake.Entity` — maps plan to the client struct
- `fromApiModel(thetalake.Entity) stateModel` — maps client struct to state; handle nullable pointers (e.g. `if library.ExternalId == nil { state.ExternalId = types.StringNull() }`)
- Use `timetypes.NewRFC3339TimeValue(t)` for `time.Time` → `timetypes.RFC3339`

#### `schema.go` — Terraform schema

- `Required` for mandatory inputs (e.g. `name`, `storage_account_id`)
- `Optional + Computed + Default` for inputs with server-side defaults
- `Computed` only for read-only fields (`id`, `created_at`, `updated_at`, counters, etc.)
- Use `timetypes.RFC3339Type{}` as `CustomType` for timestamp attributes
- Use `stringplanmodifier.UseStateForUnknown()` for optional computed fields that the API doesn't always echo back (prevents unnecessary diffs)
- Add `Validators` where the API constrains values (e.g. enum string lists)

---

### 6. Register in the provider

File: `internal/provider/provider.go`

1. Add an import alias for the new resource package:
   ```go
   <entity>resource "github.com/theta-lake/terraform-provider-thetalake/internal/resources/<entity>"
   ```
2. Add the constructor to the `Resources` slice:
   ```go
   <entity>resource.New<Entity>Resource,
   ```

---

### 7. Add an example Terraform configuration

File: `examples/resources/thetalake_<entity>/resource.tf`

Show a minimal but realistic usage of the resource with all required attributes and common optional ones.

---

### 8. Add acceptance tests

File: `internal/resources/<entity>/<entity>_test.go`

This test is required as part of creating the resource, not as follow-up cleanup.

- Package: `<entity>_test`
- Use `acctest.ProtoV6ProviderFactories` and `acctest.TestProviderConfig`
- Guard with `t.Skip(...)` if a required env var is not set
- Follow the Create→ImportState→Update→(auto-Delete) step pattern from `directory_group_test.go`
- If the test requires pre-existing data (e.g. a `storage_account_id`), read it from an env var and add the var to `internal/acctest/acctest.go`

---

### 9. Verify

```bash
go build ./...     # must produce no output
go vet ./...       # must produce no output
go test ./internal/client/thetalake/... ./internal/resources/<entity>/...
```

---

## Key conventions

| Convention | Detail |
|---|---|
| Go package name | Lower-case, no underscores (e.g. `retentionlibrary`, `supervisionspace`) |
| Resource type name | `thetalake_<snake_case>` via `Metadata` |
| JSON response key | Passed as the `responseObjectName` arg to `doRequest` |
| Nullable API fields | `*string` in the client struct; `types.StringNull()` in `fromApiModel` |
| Write-only plan fields | Preserve from plan in Create/Update, preserve from state in Read |
| Timestamps | `time.Time` in client struct → `timetypes.RFC3339` in state model |
| Import | Always implement; parse string ID to `int64` in `ImportState` |
