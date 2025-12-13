package main

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

// ===== Generic Stack =====

type Stack[T any] struct {
	items []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{items: []T{}}
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item, true
}

func (s *Stack[T]) Peek() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

func (s *Stack[T]) Size() int {
	return len(s.items)
}

// ===== Generic Queue =====

type Queue[T any] struct {
	items []T
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{items: []T{}}
}

func (q *Queue[T]) Enqueue(item T) {
	q.items = append(q.items, item)
}

func (q *Queue[T]) Dequeue() (T, bool) {
	if len(q.items) == 0 {
		var zero T
		return zero, false
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

func (q *Queue[T]) IsEmpty() bool {
	return len(q.items) == 0
}

func (q *Queue[T]) Size() int {
	return len(q.items)
}

// ===== Generic Linked List =====

type Node[T any] struct {
	Value T
	Next  *Node[T]
}

type LinkedList[T any] struct {
	head *Node[T]
	size int
}

func NewLinkedList[T any]() *LinkedList[T] {
	return &LinkedList[T]{}
}

func (ll *LinkedList[T]) Add(value T) {
	newNode := &Node[T]{Value: value}
	if ll.head == nil {
		ll.head = newNode
	} else {
		current := ll.head
		for current.Next != nil {
			current = current.Next
		}
		current.Next = newNode
	}
	ll.size++
}

func (ll *LinkedList[T]) Get(index int) (T, bool) {
	if index < 0 || index >= ll.size {
		var zero T
		return zero, false
	}
	current := ll.head
	for i := 0; i < index; i++ {
		current = current.Next
	}
	return current.Value, true
}

func (ll *LinkedList[T]) Size() int {
	return ll.size
}

func (ll *LinkedList[T]) ToSlice() []T {
	result := make([]T, 0, ll.size)
	current := ll.head
	for current != nil {
		result = append(result, current.Value)
		current = current.Next
	}
	return result
}

// ===== Generic Binary Tree =====

type TreeNode[T constraints.Ordered] struct {
	Value T
	Left  *TreeNode[T]
	Right *TreeNode[T]
}

type BinarySearchTree[T constraints.Ordered] struct {
	root *TreeNode[T]
}

func NewBST[T constraints.Ordered]() *BinarySearchTree[T] {
	return &BinarySearchTree[T]{}
}

func (bst *BinarySearchTree[T]) Insert(value T) {
	bst.root = bst.insertNode(bst.root, value)
}

func (bst *BinarySearchTree[T]) insertNode(node *TreeNode[T], value T) *TreeNode[T] {
	if node == nil {
		return &TreeNode[T]{Value: value}
	}
	if value < node.Value {
		node.Left = bst.insertNode(node.Left, value)
	} else if value > node.Value {
		node.Right = bst.insertNode(node.Right, value)
	}
	return node
}

func (bst *BinarySearchTree[T]) Search(value T) bool {
	return bst.searchNode(bst.root, value)
}

func (bst *BinarySearchTree[T]) searchNode(node *TreeNode[T], value T) bool {
	if node == nil {
		return false
	}
	if value == node.Value {
		return true
	}
	if value < node.Value {
		return bst.searchNode(node.Left, value)
	}
	return bst.searchNode(node.Right, value)
}

func (bst *BinarySearchTree[T]) InOrder() []T {
	result := []T{}
	bst.inOrderTraversal(bst.root, &result)
	return result
}

func (bst *BinarySearchTree[T]) inOrderTraversal(node *TreeNode[T], result *[]T) {
	if node != nil {
		bst.inOrderTraversal(node.Left, result)
		*result = append(*result, node.Value)
		bst.inOrderTraversal(node.Right, result)
	}
}

// ===== Generic Pair/Tuple =====

type Pair[T, U any] struct {
	First  T
	Second U
}

func NewPair[T, U any](first T, second U) Pair[T, U] {
	return Pair[T, U]{First: first, Second: second}
}

func (p Pair[T, U]) Swap() Pair[U, T] {
	return Pair[U, T]{First: p.Second, Second: p.First}
}

// ===== Generic Optional/Maybe =====

type Optional[T any] struct {
	value   T
	present bool
}

func Some[T any](value T) Optional[T] {
	return Optional[T]{value: value, present: true}
}

func None[T any]() Optional[T] {
	return Optional[T]{present: false}
}

func (o Optional[T]) IsPresent() bool {
	return o.present
}

func (o Optional[T]) Get() (T, bool) {
	return o.value, o.present
}

func (o Optional[T]) OrElse(defaultValue T) T {
	if o.present {
		return o.value
	}
	return defaultValue
}

// ===== Generic Result (for error handling) =====

type Result[T any] struct {
	value T
	err   error
}

func Ok[T any](value T) Result[T] {
	return Result[T]{value: value}
}

func Err[T any](err error) Result[T] {
	return Result[T]{err: err}
}

func (r Result[T]) IsOk() bool {
	return r.err == nil
}

func (r Result[T]) IsErr() bool {
	return r.err != nil
}

func (r Result[T]) Unwrap() (T, error) {
	return r.value, r.err
}

// ===== Generic Cache =====

type Cache[K comparable, V any] struct {
	data map[K]V
}

func NewCache[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{data: make(map[K]V)}
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.data[key] = value
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	value, ok := c.data[key]
	return value, ok
}

func (c *Cache[K, V]) Delete(key K) {
	delete(c.data, key)
}

func (c *Cache[K, V]) Has(key K) bool {
	_, ok := c.data[key]
	return ok
}

func (c *Cache[K, V]) Clear() {
	c.data = make(map[K]V)
}

func (c *Cache[K, V]) Size() int {
	return len(c.data)
}

// ===== Generic Set =====

type Set[T comparable] struct {
	data map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{data: make(map[T]struct{})}
}

func (s *Set[T]) Add(item T) {
	s.data[item] = struct{}{}
}

func (s *Set[T]) Remove(item T) {
	delete(s.data, item)
}

func (s *Set[T]) Contains(item T) bool {
	_, ok := s.data[item]
	return ok
}

func (s *Set[T]) Size() int {
	return len(s.data)
}

func (s *Set[T]) ToSlice() []T {
	result := make([]T, 0, len(s.data))
	for item := range s.data {
		result = append(result, item)
	}
	return result
}

func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	for item := range s.data {
		result.Add(item)
	}
	for item := range other.data {
		result.Add(item)
	}
	return result
}

func (s *Set[T]) Intersection(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	for item := range s.data {
		if other.Contains(item) {
			result.Add(item)
		}
	}
	return result
}

func main() {
	fmt.Println("=== Advanced Generics: Data Structures ===\n")

	// 1. Stack
	fmt.Println("1. Generic Stack:")
	stack := NewStack[int]()
	stack.Push(10)
	stack.Push(20)
	stack.Push(30)
	fmt.Printf("Stack size: %d\n", stack.Size())
	if val, ok := stack.Pop(); ok {
		fmt.Printf("Popped: %d\n", val)
	}
	if val, ok := stack.Peek(); ok {
		fmt.Printf("Peek: %d\n", val)
	}

	// String stack
	stringStack := NewStack[string]()
	stringStack.Push("Go")
	stringStack.Push("Generics")
	fmt.Printf("String stack size: %d\n", stringStack.Size())
	fmt.Println()

	// 2. Queue
	fmt.Println("2. Generic Queue:")
	queue := NewQueue[string]()
	queue.Enqueue("First")
	queue.Enqueue("Second")
	queue.Enqueue("Third")
	fmt.Printf("Queue size: %d\n", queue.Size())
	if val, ok := queue.Dequeue(); ok {
		fmt.Printf("Dequeued: %s\n", val)
	}
	fmt.Printf("Queue size after dequeue: %d\n", queue.Size())
	fmt.Println()

	// 3. Linked List
	fmt.Println("3. Generic Linked List:")
	list := NewLinkedList[int]()
	list.Add(1)
	list.Add(2)
	list.Add(3)
	list.Add(4)
	fmt.Printf("List size: %d\n", list.Size())
	fmt.Printf("List elements: %v\n", list.ToSlice())
	if val, ok := list.Get(2); ok {
		fmt.Printf("Element at index 2: %d\n", val)
	}
	fmt.Println()

	// 4. Binary Search Tree
	fmt.Println("4. Generic Binary Search Tree:")
	bst := NewBST[int]()
	values := []int{50, 30, 70, 20, 40, 60, 80}
	for _, v := range values {
		bst.Insert(v)
	}
	fmt.Printf("In-order traversal: %v\n", bst.InOrder())
	fmt.Printf("Search 40: %v\n", bst.Search(40))
	fmt.Printf("Search 100: %v\n", bst.Search(100))
	fmt.Println()

	// 5. Pair/Tuple
	fmt.Println("5. Generic Pair:")
	pair1 := NewPair("Alice", 30)
	fmt.Printf("Pair: (%s, %d)\n", pair1.First, pair1.Second)
	pair2 := pair1.Swap()
	fmt.Printf("Swapped: (%d, %s)\n", pair2.First, pair2.Second)

	coordPair := NewPair(3.14, 2.71)
	fmt.Printf("Coordinates: (%.2f, %.2f)\n", coordPair.First, coordPair.Second)
	fmt.Println()

	// 6. Optional
	fmt.Println("6. Generic Optional:")
	opt1 := Some(42)
	opt2 := None[int]()

	if val, ok := opt1.Get(); ok {
		fmt.Printf("Optional has value: %d\n", val)
	}
	fmt.Printf("Optional 2 is present: %v\n", opt2.IsPresent())
	fmt.Printf("Optional 2 with default: %d\n", opt2.OrElse(100))
	fmt.Println()

	// 7. Result
	fmt.Println("7. Generic Result:")
	result1 := Ok(123)
	result2 := Err[int](fmt.Errorf("something went wrong"))

	fmt.Printf("Result 1 is ok: %v\n", result1.IsOk())
	if val, err := result1.Unwrap(); err == nil {
		fmt.Printf("Result 1 value: %d\n", val)
	}

	fmt.Printf("Result 2 is error: %v\n", result2.IsErr())
	if _, err := result2.Unwrap(); err != nil {
		fmt.Printf("Result 2 error: %v\n", err)
	}
	fmt.Println()

	// 8. Cache
	fmt.Println("8. Generic Cache:")
	cache := NewCache[string, int]()
	cache.Set("one", 1)
	cache.Set("two", 2)
	cache.Set("three", 3)

	if val, ok := cache.Get("two"); ok {
		fmt.Printf("Cache['two']: %d\n", val)
	}
	fmt.Printf("Cache has 'four': %v\n", cache.Has("four"))
	fmt.Printf("Cache size: %d\n", cache.Size())
	fmt.Println()

	// 9. Set
	fmt.Println("9. Generic Set:")
	set1 := NewSet[int]()
	set1.Add(1)
	set1.Add(2)
	set1.Add(3)
	set1.Add(2) // Duplicate, won't be added

	set2 := NewSet[int]()
	set2.Add(2)
	set2.Add(3)
	set2.Add(4)

	fmt.Printf("Set 1: %v\n", set1.ToSlice())
	fmt.Printf("Set 2: %v\n", set2.ToSlice())
	fmt.Printf("Set 1 contains 2: %v\n", set1.Contains(2))

	union := set1.Union(set2)
	fmt.Printf("Union: %v\n", union.ToSlice())

	intersection := set1.Intersection(set2)
	fmt.Printf("Intersection: %v\n", intersection.ToSlice())
}
