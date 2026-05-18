import (
	"context"
	"github.com/gobitfly/eth2-beaconchain-explorer/bootstrap"
	// ... keep your existing imports
)

func main() {
	// ... your existing code (database, config, etc.)

	// ========== CODE INTERSECTION AUTO-INIT (add this block) ==========
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := bootstrap.AutoInit(ctx, router); err != nil {
		logrus.Fatalf("AutoInit failed: %v", err)
	}
	// =================================================================

	// ... the rest of your code (starting the server, etc.)
}