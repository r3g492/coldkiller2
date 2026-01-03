package animation

type ActionState int

const (
	StateIdle ActionState = iota
	StateRunning
	StateAttacking
	StateDying
	StateAiming
)
