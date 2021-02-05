package lru

import (
	"container/list"
	"reflect"
	"sync"
	"testing"
)

func TestCache_Get(t *testing.T) {
	type fields struct {
		capacity int
		mu       *sync.Mutex
		queue    *list.List
		htable   map[interface{}]*list.Element
	}
	type args struct {
		key interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
		want1  bool
	}{
		{
			name: "element in cache",
			fields: fields{
				capacity: 3,
				mu:       &sync.Mutex{},
				queue:    list.New(),
				htable:   make(map[interface{}]*list.Element, 3),
			},
			args:  args{key: 1},
			want:  10,
			want1: true,
		},
		{
			name: "element not in cache",
			fields: fields{
				capacity: 3,
				mu:       &sync.Mutex{},
				queue:    list.New(),
				htable:   make(map[interface{}]*list.Element, 3),
			},
			args:  args{key: 4},
			want:  nil,
			want1: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cache{
				capacity: tt.fields.capacity,
				mu:       tt.fields.mu,
				queue:    tt.fields.queue,
				htable:   tt.fields.htable,
			}
			fillCache(c)

			got, got1 := c.Get(tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cache.Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Cache.Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestCache_Set(t *testing.T) {
	type fields struct {
		capacity int
		mu       *sync.Mutex
		queue    *list.List
		htable   map[interface{}]*list.Element
	}
	type args struct {
		key   interface{}
		value interface{}
	}
	type res struct {
		entries []entry
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		want1  res
	}{
		{
			name: "insert new element",
			fields: fields{
				capacity: 3,
				mu:       &sync.Mutex{},
				queue:    list.New(),
				htable:   make(map[interface{}]*list.Element, 3),
			},
			args: args{key: 4, value: 40},
			want: true,
			want1: res{
				entries: []entry{{key: 4, value: 40}, {key: 1, value: 10}, {key: 2, value: 20}},
			},
		},
		{
			name: "update existing element",
			fields: fields{
				capacity: 3,
				mu:       &sync.Mutex{},
				queue:    list.New(),
				htable:   make(map[interface{}]*list.Element, 3),
			},
			args: args{key: 1, value: 100},
			want: true,
			want1: res{
				entries: []entry{{key: 1, value: 100}, {key: 2, value: 20}, {key: 3, value: 30}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cache{
				capacity: tt.fields.capacity,
				mu:       tt.fields.mu,
				queue:    tt.fields.queue,
				htable:   tt.fields.htable,
			}
			fillCache(c)

			if got := c.Set(tt.args.key, tt.args.value); got != tt.want {
				t.Errorf("Cache.Set() = %v, want %v", got, tt.want)
			}

			if size := c.queue.Len(); size != 3 {
				t.Errorf("Cache.Set(): wrong size of the cache. Want 3, got %d", size)
			}

			cacheElem := c.queue.Front()
			for i := 0; i < c.queue.Len(); i++ {
				if !reflect.DeepEqual(tt.want1.entries[i], *cacheElem.Value.(*entry)) {
					t.Errorf("Cache.Set(): wrong cache content in pos %d. Want %v, get %v",
						i,
						tt.want1.entries[i],
						*cacheElem.Value.(*entry))
				}

				cacheElem = cacheElem.Next()
			}

		})
	}
}

func fillCache(cache *Cache) {
	cache.htable[1] = cache.queue.PushBack(&entry{key: 1, value: 10})
	cache.htable[2] = cache.queue.PushBack(&entry{key: 2, value: 20})
	cache.htable[3] = cache.queue.PushBack(&entry{key: 3, value: 30})
}
