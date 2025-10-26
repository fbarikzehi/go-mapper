# Migration Guide

Guide for migrating from other mapping libraries or upgrading between versions.

## From Other Libraries

### From jinzhu/copier

**Before:**

```go
copier.Copy(&dst, &src)
copier.CopyWithOption(&dst, &src, copier.Option{IgnoreEmpty: true})
```

**After:**

```go
mapper.Copy(&dst, src)
mapper.Copy(&dst, src, mapper.WithIgnoreNilFields(true))
```

### From mitchellh/mapstructure

**Before:**

```go
mapstructure.Decode(input, &result)
```

**After:**

```go
mapper.Copy(&result, input)
```

### From go-playground/validator + manual mapping

**Before:**

```go
result := Destination{
    Name:  src.Name,
    Age:   src.Age,
    Email: src.Email,
}
validate.Struct(result)
```

**After:**

```go
var result Destination
mapper.Copy(&result, src)
validate.Struct(result)
```

## Version Migration

### From v0.x to v1.x

#### Breaking Changes

1. **Package Import**

   ```go
   // Old
   import "github.com/fbarikzehi/gomap/v0"

   // New
   import "github.com/fbarikzehi/gomap"
   ```

2. **Option Names**

   ```go
   // Old
   mapper.Copy(&dst, src, mapper.MaxDepth(5))

   // New
   mapper.Copy(&dst, src, mapper.WithMaxDepth(5))
   ```

3. **Error Handling**

   ```go
   // Old - errors were ignored
   mapper.Copy(&dst, src)

   // New - errors must be handled
   if err := mapper.Copy(&dst, src); err != nil {
       return err
   }
   ```

#### New Features in v1.x

- Custom type converters
- Field name mappers
- Error handlers
- Thread-safe operations
- Performance improvements

## Best Practices

### 1. Always Check Errors

```go
// Do this
if err := mapper.Copy(&dst, src); err != nil {
    log.Printf("mapping failed: %v", err)
    return err
}

// Not this
mapper.Copy(&dst, src) // Ignores errors
```

### 2. Use Type-Safe Wrappers

```go
// Create type-safe helper functions
func ToUserDTO(model UserModel) (UserDTO, error) {
    var dto UserDTO
    err := mapper.Copy(&dto, model,
        mapper.WithTagName("json"),
        mapper.WithMaxDepth(5),
    )
    return dto, err
}
```

### 3. Configure Once, Use Many Times

```go
// Global or package-level mapper
var defaultMapper = mapper.NewMapper(
    mapper.WithMaxDepth(10),
    mapper.WithIgnoreUnexported(true),
)

func ConvertData(src Source) (Destination, error) {
    var dst Destination
    err := defaultMapper.Map(&dst, src)
    return dst, err
}
```

### 4. Document Custom Converters

```go
// timeToStringConverter converts time.Time to RFC3339 string format.
// Used for API responses where timestamps need to be in string format.
var timeToStringConverter = func(v reflect.Value) (reflect.Value, error) {
    if t, ok := v.Interface().(time.Time); ok {
        return reflect.ValueOf(t.Format(time.RFC3339)), nil
    }
    return v, nil
}
```

### 5. Test Your Mappings

```go
func TestUserMapping(t *testing.T) {
    tests := []struct {
        name    string
        src     UserModel
        want    UserDTO
        wantErr bool
    }{
        {
            name: "valid user",
            src: UserModel{
                ID:       1,
                Username: "john",
                Email:    "john@example.com",
            },
            want: UserDTO{
                ID:       1,
                Username: "john",
                Email:    "john@example.com",
            },
            wantErr: false,
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ToUserDTO(tt.src)
            if (err != nil) != tt.wantErr {
                t.Errorf("ToUserDTO() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("ToUserDTO() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Common Patterns

### Pattern 1: DTO Conversion Layer

```go
// dto/user.go
package dto

import "github.com/fbarikzehi/gomap"

var userMapper = mapper.NewMapper(
    mapper.WithTagName("json"),
    mapper.WithMaxDepth(5),
)

func ToUserDTO(model models.User) (UserDTO, error) {
    var dto UserDTO
    err := userMapper.Map(&dto, model)
    return dto, err
}

func ToUserModel(dto UserDTO) (models.User, error) {
    var model models.User
    err := userMapper.Map(&model, dto)
    return model, err
}
```

### Pattern 2: Repository Layer

```go
// repository/user_repository.go
package repository

type UserRepository struct {
    db     *sql.DB
    mapper *mapper.Mapper
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{
        db: db,
        mapper: mapper.NewMapper(
            mapper.WithMaxDepth(3),
        ),
    }
}

func (r *UserRepository) Save(user domain.User) error {
    var dbUser DBUser
    if err := r.mapper.Map(&dbUser, user); err != nil {
        return err
    }

    // Save to database
    return r.saveToDB(dbUser)
}
```

### Pattern 3: Service Layer

```go
// service/user_service.go
package service

type UserService struct {
    repo   *repository.UserRepository
    mapper *mapper.Mapper
}

func (s *UserService) CreateUser(req CreateUserRequest) (*UserResponse, error) {
    // Convert request to domain model
    var user domain.User
    if err := s.mapper.Map(&user, req); err != nil {
        return nil, fmt.Errorf("invalid request: %w", err)
    }

    // Business logic
    if err := user.Validate(); err != nil {
        return nil, err
    }

    // Save
    if err := s.repo.Save(user); err != nil {
        return nil, err
    }

    // Convert to response
    var response UserResponse
    if err := s.mapper.Map(&response, user); err != nil {
        return nil, err
    }

    return &response, nil
}
```

## Troubleshooting

### Issue: Fields Not Mapping

**Problem:**

```go
type Source struct {
    name string // unexported
}

type Destination struct {
    Name string
}
```

**Solution:**

```go
// Either export the field
type Source struct {
    Name string // exported
}

// Or use WithAllowPrivateFields (use with caution)
mapper.Copy(&dst, src, mapper.WithAllowPrivateFields(true))
```

### Issue: Circular Reference Error

**Problem:**

```go
type Node struct {
    Value int
    Next  *Node
}

node1 := &Node{Value: 1}
node2 := &Node{Value: 2}
node1.Next = node2
node2.Next = node1 // Circular reference

mapper.Copy(&dst, node1) // Error!
```

**Solution:**

```go
// Skip circular check if you handle it manually
mapper.Copy(&dst, node1, mapper.WithSkipCircularCheck(true))

// Or break the cycle before mapping
node2.Next = nil
```

### Issue: Type Mismatch

**Problem:**

```go
type Source struct {
    Value int
}

type Destination struct {
    Value string // Different type
}

mapper.Copy(&dst, src) // Value won't map
```

**Solution:**

```go
// Use custom converter
intToStringConverter := func(v reflect.Value) (reflect.Value, error) {
    if v.Kind() == reflect.Int {
        return reflect.ValueOf(strconv.Itoa(int(v.Int()))), nil
    }
    return v, nil
}

mapper.Copy(&dst, src,
    mapper.WithCustomConverter(reflect.TypeOf(0), intToStringConverter),
)
```

### Issue: Performance Degradation

**Problem:**
Mapping is slow for large datasets.

**Solution:**

```go
// 1. Reuse mapper instance
m := mapper.NewMapper()

// 2. Skip unnecessary checks
m := mapper.NewMapper(
    mapper.WithSkipCircularCheck(true),
    mapper.WithMaxDepth(5),
)

// 3. Process in batches
const batchSize = 1000
for i := 0; i < len(sources); i += batchSize {
    end := min(i+batchSize, len(sources))
    processBatch(sources[i:end])
}
```

## FAQ

**Q: Can I map between different types?**
A: Yes, using custom converters or if types are convertible.

**Q: Is it thread-safe?**
A: Yes, Mapper instances are thread-safe.

**Q: How do I handle nested structs?**
A: Nested structs are automatically mapped recursively.

**Q: Can I map to interfaces?**
A: Yes, but the destination must be a pointer to the interface.

**Q: How do I skip certain fields?**
A: Use struct tags with `-` or implement custom logic in error handler.

**Q: Is there a size limit?**
A: The MaxSliceCapacity option limits slice allocation (default: 1M elements).

For more questions, see [GitHub Issues](https://github.com/fbarikzehi/gomap/issues).
