# Go Map Examples

Comprehensive examples demonstrating various use cases and features of gomap.

## Table of Contents

- [Basic Examples](#basic-examples)
- [Advanced Examples](#advanced-examples)
- [Real-World Use Cases](#real-world-use-cases)
- [Error Handling](#error-handling)
- [Performance Optimization](#performance-optimization)

## Basic Examples

### 1. Simple Struct Mapping

```go
package main

import (
    "fmt"
    "github.com/fbarikzehi/gomap"
)

type Source struct {
    Name  string
    Age   int
    Email string
}

type Destination struct {
    Name  string
    Age   int
    Email string
}

func main() {
    src := Source{
        Name:  "Alice",
        Age:   30,
        Email: "alice@example.com",
    }

    var dst Destination
    if err := mapper.Copy(&dst, src); err != nil {
        panic(err)
    }

    fmt.Printf("%+v\n", dst)
}
```

### 2. Nested Structures

```go
type Address struct {
    Street  string
    City    string
    Country string
}

type Person struct {
    Name    string
    Age     int
    Address Address
}

func example() {
    src := Person{
        Name: "Bob",
        Age:  25,
        Address: Address{
            Street:  "123 Main St",
            City:    "New York",
            Country: "USA",
        },
    }

    var dst Person
    mapper.Copy(&dst, src)
}
```

### 3. Slices and Arrays

```go
type Item struct {
    ID   int
    Name string
}

type Collection struct {
    Items []Item
    Tags  [3]string
}

func example() {
    src := Collection{
        Items: []Item{
            {ID: 1, Name: "First"},
            {ID: 2, Name: "Second"},
        },
        Tags: [3]string{"tag1", "tag2", "tag3"},
    }

    var dst Collection
    mapper.Copy(&dst, src)
}
```

### 4. Maps

```go
type Config struct {
    Settings map[string]string
    Metadata map[string]interface{}
}

func example() {
    src := Config{
        Settings: map[string]string{
            "theme": "dark",
            "lang":  "en",
        },
        Metadata: map[string]interface{}{
            "version": "1.0",
            "enabled": true,
        },
    }

    var dst Config
    mapper.Copy(&dst, src)
}
```

### 5. Pointers

```go
type Node struct {
    Value int
    Next  *Node
}

func example() {
    node2 := &Node{Value: 2}
    node1 := &Node{Value: 1, Next: node2}

    var dst Node
    mapper.Copy(&dst, node1)
}
```

## Advanced Examples

### 6. Tag-Based Mapping

```go
type SourceDTO struct {
    FullName    string `mapper:"name"`
    YearsOld    int    `mapper:"age"`
    ContactMail string `mapper:"email"`
}

type DestinationModel struct {
    Name  string
    Age   int
    Email string
}

func example() {
    src := SourceDTO{
        FullName:    "Charlie",
        YearsOld:    35,
        ContactMail: "charlie@example.com",
    }

    var dst DestinationModel
    mapper.Copy(&dst, src, mapper.WithTagName("mapper"))
}
```

### 7. Custom Type Converters

```go
import "time"

type Event struct {
    Name      string
    Timestamp time.Time
}

type EventDTO struct {
    Name      string
    Timestamp string
}

func example() {
    src := Event{
        Name:      "Conference",
        Timestamp: time.Now(),
    }

    timeConverter := func(v reflect.Value) (reflect.Value, error) {
        if t, ok := v.Interface().(time.Time); ok {
            formatted := t.Format("2006-01-02 15:04:05")
            return reflect.ValueOf(formatted), nil
        }
        return v, nil
    }

    var dst EventDTO
    mapper.Copy(&dst, src,
        mapper.WithCustomConverter(reflect.TypeOf(time.Time{}), timeConverter),
    )
}
```

### 8. Case-Insensitive Mapping

```go
type APIResponse struct {
    USERID   int
    USERNAME string
    USEREMAIL string
}

type User struct {
    UserId    int
    UserName  string
    UserEmail string
}

func example() {
    src := APIResponse{
        USERID:    123,
        USERNAME:  "john_doe",
        USEREMAIL: "john@example.com",
    }

    var dst User
    mapper.Copy(&dst, src, mapper.WithCaseSensitive(false))
}
```

### 9. Field Name Transformation

```go
import "strings"

type SnakeCaseModel struct {
    First_Name string
    Last_Name  string
    Email_Address string
}

type CamelCaseModel struct {
    FirstName    string
    LastName     string
    EmailAddress string
}

func snakeToCamel(s string) string {
    parts := strings.Split(s, "_")
    for i := range parts {
        parts[i] = strings.Title(strings.ToLower(parts[i]))
    }
    return strings.Join(parts, "")
}

func example() {
    src := SnakeCaseModel{
        First_Name:    "Jane",
        Last_Name:     "Doe",
        Email_Address: "jane@example.com",
    }

    var dst CamelCaseModel
    mapper.Copy(&dst, src, mapper.WithFieldNameMapper(snakeToCamel))
}
```

### 10. Reusable Mapper Configuration

```go
func setupMapper() *mapper.Mapper {
    return mapper.NewMapper(
        mapper.WithMaxDepth(10),
        mapper.WithIgnoreUnexported(true),
        mapper.WithDeepCopy(true),
        mapper.WithCaseSensitive(false),
    )
}

func example() {
    m := setupMapper()

    var dst1, dst2, dst3 Destination
    m.Map(&dst1, src1)
    m.Map(&dst2, src2)
    m.Map(&dst3, src3)
}
```

## Real-World Use Cases

### 11. API Layer Mapping

```go
// Database Model
type UserModel struct {
    ID           int64
    Username     string
    Email        string
    PasswordHash string
    CreatedAt    time.Time
    UpdatedAt    time.Time
    IsActive     bool
}

// API Response DTO
type UserResponse struct {
    ID        int64  `json:"id"`
    Username  string `json:"username"`
    Email     string `json:"email"`
    CreatedAt string `json:"created_at"`
    IsActive  bool   `json:"is_active"`
}

func ToUserResponse(model UserModel) (UserResponse, error) {
    var response UserResponse

    timeConverter := func(v reflect.Value) (reflect.Value, error) {
        if t, ok := v.Interface().(time.Time); ok {
            return reflect.ValueOf(t.Format(time.RFC3339)), nil
        }
        return v, nil
    }

    err := mapper.Copy(&response, model,
        mapper.WithCustomConverter(reflect.TypeOf(time.Time{}), timeConverter),
    )

    return response, err
}

// API Request DTO
type CreateUserRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

func ToUserModel(req CreateUserRequest) UserModel {
    var model UserModel
    mapper.Copy(&model, req)

    // Set additional fields
    model.CreatedAt = time.Now()
    model.UpdatedAt = time.Now()
    model.IsActive = true

    return model
}
```

### 12. Configuration Management

```go
type DatabaseConfig struct {
    Host     string
    Port     int
    Username string
    Password string
    Database string
    SSLMode  string
}

type AppConfig struct {
    Database DatabaseConfig
    Server   ServerConfig
    Cache    CacheConfig
}

func MergeConfigs(base, override AppConfig) AppConfig {
    var merged AppConfig

    // Copy base config
    mapper.Copy(&merged, base)

    // Merge override, ignoring nil fields
    mapper.Copy(&merged, override, mapper.WithIgnoreNilFields(true))

    return merged
}
```

### 13. Data Migration

```go
// Old schema
type LegacyUser struct {
    UserID    int
    FirstName string
    LastName  string
    Contact   string
}

// New schema
type ModernUser struct {
    ID        int
    FullName  string
    Email     string
    Phone     string
}

func MigrateUser(legacy LegacyUser) ModernUser {
    var modern ModernUser

    // Basic mapping
    mapper.Copy(&modern, legacy)

    // Manual transformations
    modern.ID = legacy.UserID
    modern.FullName = fmt.Sprintf("%s %s", legacy.FirstName, legacy.LastName)

    // Parse contact (email or phone)
    if strings.Contains(legacy.Contact, "@") {
        modern.Email = legacy.Contact
    } else {
        modern.Phone = legacy.Contact
    }

    return modern
}
```

### 14. Caching Layer

```go
type CacheEntry struct {
    Key       string
    Value     interface{}
    ExpiresAt time.Time
    Metadata  map[string]string
}

type SerializedCacheEntry struct {
    Key       string
    Value     []byte
    ExpiresAt int64
    Metadata  map[string]string
}

func SerializeCacheEntry(entry CacheEntry) (SerializedCacheEntry, error) {
    var serialized SerializedCacheEntry

    timestampConverter := func(v reflect.Value) (reflect.Value, error) {
        if t, ok := v.Interface().(time.Time); ok {
            return reflect.ValueOf(t.Unix()), nil
        }
        return v, nil
    }

    err := mapper.Copy(&serialized, entry,
        mapper.WithCustomConverter(reflect.TypeOf(time.Time{}), timestampConverter),
    )

    return serialized, err
}
```

### 15. Event Sourcing

```go
type Event struct {
    ID        string
    Type      string
    Payload   map[string]interface{}
    Timestamp time.Time
    Version   int
}

type EventDTO struct {
    ID        string                 `json:"id"`
    Type      string                 `json:"type"`
    Payload   map[string]interface{} `json:"payload"`
    Timestamp string                 `json:"timestamp"`
    Version   int                    `json:"version"`
}

func ToEventDTO(event Event) (EventDTO, error) {
    var dto EventDTO

    timeConverter := func(v reflect.Value) (reflect.Value, error) {
        if t, ok := v.Interface().(time.Time); ok {
            return reflect.ValueOf(t.Format(time.RFC3339Nano)), nil
        }
        return v, nil
    }

    err := mapper.Copy(&dto, event,
        mapper.WithCustomConverter(reflect.TypeOf(time.Time{}), timeConverter),
    )

    return dto, err
}
```

## Error Handling

### 16. Custom Error Handler

```go
func example() {
    src := Source{/* ... */}
    var dst Destination

    errorHandler := func(err error, srcField, dstField string) error {
        log.Printf("Mapping error: %s -> %s: %v", srcField, dstField, err)

        // Continue mapping despite errors
        return nil
    }

    err := mapper.Copy(&dst, src,
        mapper.WithErrorHandler(errorHandler),
    )

    if err != nil {
        log.Fatal(err)
    }
}
```

### 17. Validation After Mapping

```go
func MapAndValidate(src Source) (Destination, error) {
    var dst Destination

    if err := mapper.Copy(&dst, src); err != nil {
        return dst, fmt.Errorf("mapping failed: %w", err)
    }

    // Validate mapped data
    if dst.Name == "" {
        return dst, errors.New("name is required")
    }

    if dst.Age < 0 {
        return dst, errors.New("age must be positive")
    }

    return dst, nil
}
```

## Performance Optimization

### 18. Batch Processing

```go
func ProcessBatch(sources []Source) ([]Destination, error) {
    // Create reusable mapper
    m := mapper.NewMapper(
        mapper.WithMaxDepth(5),
        mapper.WithSkipCircularCheck(true), // If no circular refs
    )

    destinations := make([]Destination, len(sources))

    for i, src := range sources {
        if err := m.Map(&destinations[i], src); err != nil {
            return nil, fmt.Errorf("failed to map item %d: %w", i, err)
        }
    }

    return destinations, nil
}
```

### 19. Concurrent Mapping

```go
import "sync"

func ProcessConcurrent(sources []Source) ([]Destination, error) {
    m := mapper.NewMapper()
    destinations := make([]Destination, len(sources))

    var wg sync.WaitGroup
    errChan := make(chan error, len(sources))

    for i := range sources {
        wg.Add(1)
        go func(idx int) {
            defer wg.Done()
            if err := m.Map(&destinations[idx], sources[idx]); err != nil {
                errChan <- fmt.Errorf("item %d: %w", idx, err)
            }
        }(i)
    }

    wg.Wait()
    close(errChan)

    // Check for errors
    for err := range errChan {
        return nil, err
    }

    return destinations, nil
}
```

---
