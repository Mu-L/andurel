package cache

import (
	"errors"
	"sync"
	"testing"
	"time"
)

func TestEntry(t *testing.T) {
	t.Run("IsExpired_NotExpired", func(t *testing.T) {
		entry := &Entry{
			Value:     "test",
			ExpiresAt: time.Now().Add(1 * time.Hour),
		}

		if entry.IsExpired() {
			t.Error("Entry should not be expired")
		}
	})

	t.Run("IsExpired_Expired", func(t *testing.T) {
		entry := &Entry{
			Value:     "test",
			ExpiresAt: time.Now().Add(-1 * time.Hour),
		}

		if !entry.IsExpired() {
			t.Error("Entry should be expired")
		}
	})

	t.Run("IsExpired_JustExpired", func(t *testing.T) {
		entry := &Entry{
			Value:     "test",
			ExpiresAt: time.Now().Add(-1 * time.Millisecond),
		}

		time.Sleep(2 * time.Millisecond)

		if !entry.IsExpired() {
			t.Error("Entry should be expired")
		}
	})
}

func TestFileSystemCache_New(t *testing.T) {
	ttl := 10 * time.Minute
	cache := NewFileSystemCache(ttl)

	if cache == nil {
		t.Fatal("Expected non-nil cache")
	}

	if cache.ttl != ttl {
		t.Errorf("Expected TTL %v, got %v", ttl, cache.ttl)
	}

	if cache.entries == nil {
		t.Error("Expected non-nil entries map")
	}

	if len(cache.entries) != 0 {
		t.Error("Expected empty cache on initialization")
	}
}

func TestFileSystemCache_SetAndGet(t *testing.T) {
	cache := NewFileSystemCache(1 * time.Hour)

	t.Run("SetAndGetString", func(t *testing.T) {
		key := "test-key"
		value := "test-value"

		cache.Set(key, value)

		retrieved, found := cache.Get(key)
		if !found {
			t.Error("Expected to find cached value")
		}

		if retrieved != value {
			t.Errorf("Expected value %v, got %v", value, retrieved)
		}
	})

	t.Run("SetAndGetInt", func(t *testing.T) {
		key := "int-key"
		value := 42

		cache.Set(key, value)

		retrieved, found := cache.Get(key)
		if !found {
			t.Error("Expected to find cached value")
		}

		if retrieved != value {
			t.Errorf("Expected value %v, got %v", value, retrieved)
		}
	})

	t.Run("SetAndGetStruct", func(t *testing.T) {
		type TestStruct struct {
			Name string
			Age  int
		}

		key := "struct-key"
		value := TestStruct{Name: "Alice", Age: 30}

		cache.Set(key, value)

		retrieved, found := cache.Get(key)
		if !found {
			t.Error("Expected to find cached value")
		}

		if retrieved != value {
			t.Errorf("Expected value %v, got %v", value, retrieved)
		}
	})

	t.Run("GetNonExistentKey", func(t *testing.T) {
		_, found := cache.Get("non-existent")
		if found {
			t.Error("Should not find non-existent key")
		}
	})

	t.Run("OverwriteExistingKey", func(t *testing.T) {
		key := "overwrite-key"
		value1 := "first-value"
		value2 := "second-value"

		cache.Set(key, value1)
		cache.Set(key, value2)

		retrieved, found := cache.Get(key)
		if !found {
			t.Error("Expected to find cached value")
		}

		if retrieved != value2 {
			t.Errorf("Expected value %v, got %v", value2, retrieved)
		}
	})
}

func TestFileSystemCache_Delete(t *testing.T) {
	cache := NewFileSystemCache(1 * time.Hour)

	t.Run("DeleteExistingKey", func(t *testing.T) {
		key := "delete-key"
		value := "delete-value"

		cache.Set(key, value)

		_, found := cache.Get(key)
		if !found {
			t.Fatal("Expected to find cached value before delete")
		}

		cache.Delete(key)

		_, found = cache.Get(key)
		if found {
			t.Error("Should not find deleted key")
		}
	})

	t.Run("DeleteNonExistentKey", func(t *testing.T) {
		// Should not panic
		cache.Delete("non-existent-key")
	})
}

func TestFileSystemCache_Clear(t *testing.T) {
	cache := NewFileSystemCache(1 * time.Hour)

	// Add multiple entries
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")

	if len(cache.entries) != 3 {
		t.Fatalf("Expected 3 entries, got %d", len(cache.entries))
	}

	cache.Clear()

	if len(cache.entries) != 0 {
		t.Errorf("Expected 0 entries after clear, got %d", len(cache.entries))
	}

	_, found := cache.Get("key1")
	if found {
		t.Error("Should not find key1 after clear")
	}
}

func TestFileSystemCache_Expiration(t *testing.T) {
	cache := NewFileSystemCache(50 * time.Millisecond)

	t.Run("GetExpiredEntry", func(t *testing.T) {
		key := "expire-key"
		value := "expire-value"

		cache.Set(key, value)

		// Verify it's there initially
		retrieved, found := cache.Get(key)
		if !found {
			t.Fatal("Expected to find cached value")
		}
		if retrieved != value {
			t.Errorf("Expected value %v, got %v", value, retrieved)
		}

		// Wait for expiration
		time.Sleep(100 * time.Millisecond)

		// Should not be found after expiration
		_, found = cache.Get(key)
		if found {
			t.Error("Should not find expired entry")
		}
	})
}

func TestFileSystemCache_CleanupExpired(t *testing.T) {
	cache := NewFileSystemCache(50 * time.Millisecond)

	// Add multiple entries
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")

	if len(cache.entries) != 3 {
		t.Fatalf("Expected 3 entries, got %d", len(cache.entries))
	}

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	// Entries still in map (not automatically removed)
	if len(cache.entries) != 3 {
		t.Errorf("Expected 3 entries before cleanup, got %d", len(cache.entries))
	}

	// Cleanup expired entries
	cache.CleanupExpired()

	// All should be removed now
	if len(cache.entries) != 0 {
		t.Errorf("Expected 0 entries after cleanup, got %d", len(cache.entries))
	}
}

func TestFileSystemCache_ConcurrentAccess(t *testing.T) {
	cache := NewFileSystemCache(1 * time.Hour)

	t.Run("ConcurrentWrites", func(t *testing.T) {
		var wg sync.WaitGroup
		iterations := 100

		for i := 0; i < iterations; i++ {
			wg.Add(1)
			go func(n int) {
				defer wg.Done()
				key := "concurrent-key"
				cache.Set(key, n)
			}(i)
		}

		wg.Wait()

		// Should have exactly one entry (last write wins)
		_, found := cache.Get("concurrent-key")
		if !found {
			t.Error("Expected to find concurrent key")
		}
	})

	t.Run("ConcurrentReads", func(t *testing.T) {
		key := "read-key"
		value := "read-value"
		cache.Set(key, value)

		var wg sync.WaitGroup
		iterations := 100

		for i := 0; i < iterations; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				retrieved, found := cache.Get(key)
				if !found {
					t.Error("Expected to find cached value")
				}
				if retrieved != value {
					t.Errorf("Expected value %v, got %v", value, retrieved)
				}
			}()
		}

		wg.Wait()
	})

	t.Run("ConcurrentReadWriteDelete", func(t *testing.T) {
		var wg sync.WaitGroup

		// Writers
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func(n int) {
				defer wg.Done()
				cache.Set("mixed-key", n)
			}(i)
		}

		// Readers
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				cache.Get("mixed-key")
			}()
		}

		// Deleters
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				cache.Delete("mixed-key")
			}()
		}

		wg.Wait()
		// Test should not panic
	})
}

func TestGetModulePath(t *testing.T) {
	// Clear global cache before test
	ClearFileSystemCache()

	t.Run("ResolverCalled", func(t *testing.T) {
		key := "test-module-path"
		expectedPath := "/path/to/module"
		callCount := 0

		resolver := func() (string, error) {
			callCount++
			return expectedPath, nil
		}

		path, err := GetModulePath(key, resolver)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, path)
		}

		if callCount != 1 {
			t.Errorf("Expected resolver to be called once, was called %d times", callCount)
		}

		// Second call should use cache
		path, err = GetModulePath(key, resolver)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if path != expectedPath {
			t.Errorf("Expected cached path %s, got %s", expectedPath, path)
		}

		if callCount != 1 {
			t.Errorf("Expected resolver to be called once (cached), was called %d times", callCount)
		}
	})

	t.Run("ResolverError", func(t *testing.T) {
		ClearFileSystemCache()

		key := "error-module-path"
		expectedError := errors.New("resolver error")

		resolver := func() (string, error) {
			return "", expectedError
		}

		path, err := GetModulePath(key, resolver)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if err != expectedError {
			t.Errorf("Expected error %v, got %v", expectedError, err)
		}

		if path != "" {
			t.Errorf("Expected empty path, got %s", path)
		}
	})
}

func TestGetDirectoryRoot(t *testing.T) {
	ClearFileSystemCache()

	t.Run("ResolverCalled", func(t *testing.T) {
		key := "test-directory-root"
		expectedRoot := "/path/to/root"
		callCount := 0

		resolver := func() (string, error) {
			callCount++
			return expectedRoot, nil
		}

		root, err := GetDirectoryRoot(key, resolver)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if root != expectedRoot {
			t.Errorf("Expected root %s, got %s", expectedRoot, root)
		}

		if callCount != 1 {
			t.Errorf("Expected resolver to be called once, was called %d times", callCount)
		}

		// Second call should use cache
		root, err = GetDirectoryRoot(key, resolver)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if root != expectedRoot {
			t.Errorf("Expected cached root %s, got %s", expectedRoot, root)
		}

		if callCount != 1 {
			t.Errorf("Expected resolver to be called once (cached), was called %d times", callCount)
		}
	})

	t.Run("ResolverError", func(t *testing.T) {
		ClearFileSystemCache()

		key := "error-directory-root"
		expectedError := errors.New("resolver error")

		resolver := func() (string, error) {
			return "", expectedError
		}

		root, err := GetDirectoryRoot(key, resolver)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if err != expectedError {
			t.Errorf("Expected error %v, got %v", expectedError, err)
		}

		if root != "" {
			t.Errorf("Expected empty root, got %s", root)
		}
	})
}

func TestGetFileExists(t *testing.T) {
	ClearFileSystemCache()

	t.Run("FileExists", func(t *testing.T) {
		key := "test-file-exists"
		callCount := 0

		checker := func() bool {
			callCount++
			return true
		}

		exists := GetFileExists(key, checker)
		if !exists {
			t.Error("Expected file to exist")
		}

		if callCount != 1 {
			t.Errorf("Expected checker to be called once, was called %d times", callCount)
		}

		// Second call should use cache
		exists = GetFileExists(key, checker)
		if !exists {
			t.Error("Expected cached file to exist")
		}

		if callCount != 1 {
			t.Errorf("Expected checker to be called once (cached), was called %d times", callCount)
		}
	})

	t.Run("FileDoesNotExist", func(t *testing.T) {
		ClearFileSystemCache()

		key := "test-file-not-exists"

		checker := func() bool {
			return false
		}

		exists := GetFileExists(key, checker)
		if exists {
			t.Error("Expected file to not exist")
		}
	})
}

func TestClearFileSystemCache(t *testing.T) {
	// Clear first to ensure clean state
	ClearFileSystemCache()

	// Add some entries to global cache
	globalFSCache.Set("key1", "value1")
	globalFSCache.Set("key2", "value2")

	if len(globalFSCache.entries) != 2 {
		t.Fatalf("Expected 2 entries, got %d", len(globalFSCache.entries))
	}

	ClearFileSystemCache()

	if len(globalFSCache.entries) != 0 {
		t.Errorf("Expected 0 entries after clear, got %d", len(globalFSCache.entries))
	}
}

func TestCleanupExpiredFileSystemEntries(t *testing.T) {
	ClearFileSystemCache()

	// Manually set entries with expired time
	globalFSCache.mutex.Lock()
	globalFSCache.entries["expired1"] = &Entry{
		Value:     "value1",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	globalFSCache.entries["expired2"] = &Entry{
		Value:     "value2",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	globalFSCache.entries["valid"] = &Entry{
		Value:     "value3",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
	globalFSCache.mutex.Unlock()

	if len(globalFSCache.entries) != 3 {
		t.Fatalf("Expected 3 entries, got %d", len(globalFSCache.entries))
	}

	CleanupExpiredFileSystemEntries()

	if len(globalFSCache.entries) != 1 {
		t.Errorf("Expected 1 entry after cleanup, got %d", len(globalFSCache.entries))
	}

	_, found := globalFSCache.Get("valid")
	if !found {
		t.Error("Expected to find valid entry after cleanup")
	}
}
