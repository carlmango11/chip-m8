package keyboard

import "sync"

type Keyboard struct {
	mu      sync.Mutex
	pressed map[byte]bool
	waiting chan byte
}

func New() *Keyboard {
	return &Keyboard{
		pressed: map[byte]bool{},
	}
}

func (k *Keyboard) Press(n byte) {
	k.mu.Lock()
	defer k.mu.Unlock()

	k.pressed[n] = true

	if k.waiting != nil {
		k.waiting <- n
		k.waiting = nil
	}
}

func (k *Keyboard) Unpress(n byte) {
	k.mu.Lock()
	defer k.mu.Unlock()

	k.pressed[n] = false
}

func (k *Keyboard) IsPressed(n byte) bool {
	k.mu.Lock()
	defer k.mu.Unlock()

	return k.pressed[n]
}

func (k *Keyboard) Await() byte {
	ch := make(chan byte)

	k.mu.Lock()
	k.waiting = ch
	k.mu.Unlock()

	return <-ch
}
