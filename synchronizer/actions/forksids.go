package actions

// ForkIdType is the type of the forkId
type ForkIdType uint64

const (
	// WildcardForkId It match for all forkIds
	WildcardForkId ForkIdType = 0
	// ForkIDIncaberry is the forkId for incaberry
	ForkIDIncaberry = ForkIdType(6) // nolint:gomnd
	// ForkIDEtrog is the forkId for etrog
	ForkIDEtrog = ForkIdType(7) //nolint:gomnd
	// ForkIDElderberry is the forkId for Elderberry
	ForkIDElderberry = ForkIdType(8) //nolint:gomnd
	// ForkID9 is the forkId for 9
	ForkID9 = ForkIdType(9) //nolint:gomnd
	// ForkID10 is the forkId for 10 (support more counters)
	ForkID10 = ForkIdType(10) //nolint:gomnd
	// ForkID11 is the forkId for 11 (support even more counters)
	ForkID11 = ForkIdType(11) //nolint:gomnd
)

var (

	// ForksIdAll support all forkIds
	ForksIdAll = []ForkIdType{WildcardForkId}

	// ForksIdOnlyElderberry support only elderberry forkId
	ForksIdOnlyElderberry = []ForkIdType{ForkIDElderberry, ForkID9, ForkID10, ForkID11}

	// ForksIdOnlyEtrog support only etrog forkId
	ForksIdOnlyEtrog = []ForkIdType{ForkIDEtrog}

	// ForksIdToIncaberry support all forkIds till incaberry
	ForksIdToIncaberry = []ForkIdType{1, 2, 3, 4, 5, ForkIDIncaberry}
)
