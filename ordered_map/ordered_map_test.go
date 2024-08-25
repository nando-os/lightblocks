package ordered_map

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOrderedMap(t *testing.T) {
	// Initialize a new ordered map
	omap := NewConcurrentOrderedMap()

	// Test Add and Get method
	omap.Add("key1", "value1")
	if val, ok := omap.Get("key1"); !ok || val != "value1" {
		t.Errorf("expected value1, got %v", val)
	}

	// Test Delete method
	omap.Delete("key1")
	if _, ok := omap.Get("key1"); ok {
		t.Errorf("expected key1 to be deleted")
	}

	// Test GetAllItems method
	omap.Add("key1", "value1")
	omap.Add("key2", "value2")
	omap.Add("key3", "value3")

	items := omap.GetAllItems()
	if len(items) != 3 {
		t.Errorf("expected 2 items, got %d", len(items))
	}

	assert.Equal(t, "key1", items[0].Key)
	assert.Equal(t, "value1", items[0].Value)
	assert.Equal(t, "key2", items[1].Key)
	assert.Equal(t, "value2", items[1].Value)
	assert.Equal(t, "key3", items[2].Key)
	assert.Equal(t, "value3", items[2].Value)

}
